package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"cloud.google.com/go/firestore"
	choicesv1 "github.com/SaveTheRbtz/humor/gen/go/proto"
	serverImpl "github.com/SaveTheRbtz/humor/server/internal/server"
	"go.uber.org/zap"
	healthgrpc "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableStacktrace = true
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal("Failed to create logger", zap.Error(err))
	}
	defer logger.Sync()

	firestoreClient, err := firestore.NewClient(ctx, "humor-arena")
	if err != nil {
		logger.Fatal("Failed to create Firestore client", zap.Error(err))
	}
	defer firestoreClient.Close()

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			logger.Fatal("Failed to listen", zap.Error(err))
		}
		grpcServer := grpc.NewServer()

		healthServer := healthgrpc.NewServer()
		healthServer.SetServingStatus("grpc.health.v1.Health", grpc_health_v1.HealthCheckResponse_SERVING)
		grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

		choicesServer, err := serverImpl.NewServer(firestoreClient, logger)
		if err != nil {
			logger.Fatal("Failed to create server", zap.Error(err))
		}
		choicesv1.RegisterArenaServer(grpcServer, choicesServer)

		reflection.Register(grpcServer)
		logger.Info("Serving gRPC on localhost:9090")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = choicesv1.RegisterArenaHandlerFromEndpoint(ctx, grpcMux, "localhost:9090", opts)
	if err != nil {
		logger.Fatal("Failed to register gRPC gateway", zap.Error(err))
	}
	logger.Info("Serving HTTP on localhost:8080")

	mainMux := http.NewServeMux()
	mainMux.Handle("/v1/", allowCORS(grpcMux))
	mainMux.Handle("/", http.FileServer(http.Dir("./static")))

	if err := http.ListenAndServe(":8080", mainMux); err != nil {
		logger.Fatal("Failed to serve HTTP", zap.Error(err))
	}
}
