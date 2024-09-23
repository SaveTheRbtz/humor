package server

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secureRand "crypto/rand"
	insecureRand "math/rand/v2"
	"time"

	choicesv1 "github.com/SaveTheRbtz/humor/gen/go/proto"

	"github.com/google/uuid"
)

type Theme struct {
	Text   string  `firestore:"text"`
	Random float64 `firestore:"random"`
}

type Joke struct {
	ThemeID string  `firestore:"theme_id"`
	Text    string  `firestore:"text"`
	Random  float64 `firestore:"random"`
}

type Choice struct {
	ThemeID     string            `firestore:"theme_id"`
	SessionID   string            `firestore:"session_id"`
	LeftJokeID  string            `firestore:"left_joke_id"`
	RightJokeID string            `firestore:"right_joke_id"`
	Winner      *choicesv1.Winner `firestore:"winner,omitempty"`
	Known       *choicesv1.Winner `firestore:"known,omitempty"`
	CreatedAt   time.Time         `firestore:"created_at"`
	RatedAt     *time.Time        `firestore:"rated_at,omitempty"`
}

type Server struct {
	choicesv1.UnimplementedArenaServer

	logger          *zap.Logger
	firestoreClient *firestore.Client
	rand            *insecureRand.Rand
}

func NewServer(firestoreClient *firestore.Client, logger *zap.Logger) (*Server, error) {
	var seedBytes [32]byte
	if _, err := secureRand.Read(seedBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to seed random generator: %w", err)
	}
	return &Server{
		logger:          logger,
		firestoreClient: firestoreClient,
		rand:            insecureRand.New(insecureRand.NewChaCha8(seedBytes)),
	}, nil
}

func (s *Server) GetChoices(ctx context.Context, req *choicesv1.GetChoicesRequest) (*choicesv1.GetChoicesResponse, error) {
	themeQuery := s.firestoreClient.Collection("themes").Query
	themeDoc, err := s.getRandomDocument(ctx, themeQuery, s.rand.Float64())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get random theme: %v", err)
	}
	if themeDoc == nil {
		return nil, status.Error(codes.NotFound, "No themes found")
	}
	var theme Theme
	if err := themeDoc.DataTo(&theme); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decode theme: %v", err)
	}
	s.logger.Debug("GetChoices", zap.String("theme_id", themeDoc.Ref.ID), zap.String("theme_text", theme.Text))

	jokesCollection := s.firestoreClient.Collection("jokes")
	jokesQuery := jokesCollection.Where("theme_id", "==", themeDoc.Ref.ID)

	leftJokeDoc, err := s.getRandomDocument(ctx, jokesQuery, s.rand.Float64())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get first random joke: %v", err)
	}
	if leftJokeDoc == nil {
		return nil, status.Errorf(codes.NotFound, "No jokes found for theme %s (%s)", theme.Text, themeDoc.Ref.ID)
	}
	var leftJoke Joke
	if err := leftJokeDoc.DataTo(&leftJoke); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decode joke: %v", err)
	}

	var rightJokeDoc *firestore.DocumentSnapshot
	for i := 0; i < 3; i++ {
		rightJokeDoc, err = s.getRandomDocument(ctx, jokesQuery, s.rand.Float64())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to get second random joke: %v", err)
		}
		if rightJokeDoc == nil {
			return nil, status.Errorf(codes.NotFound, "No jokes found for theme %s (%s)", theme.Text, themeDoc.Ref.ID)
		}
		if rightJokeDoc.Ref.ID != leftJokeDoc.Ref.ID {
			break
		}
	}
	if rightJokeDoc.Ref.ID == leftJokeDoc.Ref.ID {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Too few jokes found for theme: %s (%s)", theme.Text, themeDoc.Ref.ID))
	}

	var rightJoke Joke
	if err := rightJokeDoc.DataTo(&rightJoke); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decode joke: %v", err)
	}

	id := uuid.New().String()
	choice := Choice{
		SessionID: req.SessionId,

		ThemeID:     themeDoc.Ref.ID,
		LeftJokeID:  leftJokeDoc.Ref.ID,
		RightJokeID: rightJokeDoc.Ref.ID,
		CreatedAt:   time.Now(),
	}

	_, err = s.firestoreClient.Collection("choices").Doc(id).Set(ctx, choice)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save choice: %v", err)
	}

	return &choicesv1.GetChoicesResponse{
		Id:        id,
		Theme:     theme.Text,
		LeftJoke:  leftJoke.Text,
		RightJoke: rightJoke.Text,
	}, nil
}

func (s *Server) getRandomDocument(ctx context.Context, baseQuery firestore.Query, r float64) (*firestore.DocumentSnapshot, error) {
	query := baseQuery.Where("random", ">=", r).OrderBy("random", firestore.Asc).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	if len(docs) > 0 {
		return docs[0], nil
	}

	query = baseQuery.Where("random", "<", r).OrderBy("random", firestore.Asc).Limit(1)
	docs, err = query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get wraped random documents: %w", err)
	}
	if len(docs) > 0 {
		return docs[0], nil
	}

	return nil, nil
}

func (s *Server) RateChoices(ctx context.Context, req *choicesv1.RateChoicesRequest) (*choicesv1.RateChoicesResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is required")
	}
	if req.Winner == choicesv1.Winner_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "Winner is required")
	}
	if req.Known == choicesv1.Winner_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "Known is required")
	}

	_, err := s.firestoreClient.Collection("choices").Doc(req.Id).Update(ctx, []firestore.Update{
		{
			Path:  "winner",
			Value: req.Winner.String(),
		},
		{
			Path:  "known",
			Value: req.Known.String(),
		},
		{
			Path:  "rated_at",
			Value: time.Now(),
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update choice: %v", err)
	}

	return &choicesv1.RateChoicesResponse{}, nil
}
