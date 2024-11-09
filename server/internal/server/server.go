package server

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	insecureRand "math/rand/v2"
	"time"

	choicesv1 "github.com/SaveTheRbtz/humor/gen/go/proto"

	"google.golang.org/api/iterator"

	"github.com/google/uuid"
)

type leaderboardEntry struct {
	Model         string  `firestore:"model"`
	Votes         int64   `firestore:"votes"`
	NewmanScore   float64 `firestore:"newman_score"`
	NewmanCiLower float64 `firestore:"newman_ci_lower"`
	NewmanCiUpper float64 `firestore:"newman_ci_upper"`
	EloScore      float64 `firestore:"elo_score"`
	EloCiLower    float64 `firestore:"elo_ci_lower"`
	EloCiUpper    float64 `firestore:"elo_ci_upper"`
}

type themesCache struct {
	topics    []string
	timestamp time.Time
}

type Theme struct {
	Text   string  `firestore:"text"`
	Random float64 `firestore:"random"`
}

type Joke struct {
	ThemeID string  `firestore:"theme_id"`
	Text    string  `firestore:"text"`
	Random  float64 `firestore:"random"`
	Model   string  `firestore:"model"`
	Policy  string  `firestore:"policy"`
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
	themeGetter     *randomDocumentGetterImpl[Theme]
}

func NewServer(
	firestoreClient *firestore.Client,
	logger *zap.Logger,
) (*Server, error) {
	randomThemeGetter, err := NewRandomDocumentGetter[Theme](
		firestoreClient,
		firestoreClient.Collection("themes").Query,
		time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create random theme getter: %w", err)
	}

	return &Server{
		logger:          logger,
		firestoreClient: firestoreClient,
		themeGetter:     randomThemeGetter,
	}, nil
}

func (s *Server) GetChoices(
	ctx context.Context,
	req *choicesv1.GetChoicesRequest,
) (*choicesv1.GetChoicesResponse, error) {
	themes, themeDocs, err := s.themeGetter.GetRandomDocuments(ctx, 1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get random theme: %v", err)
	}
	theme, themeDoc := themes[0], themeDocs[0]
	s.logger.Debug("GetChoices", zap.String("theme_id", themeDoc.Ref.ID), zap.String("theme_text", theme.Text))

	jokeGetter, err := NewRandomDocumentGetter[Joke](
		s.firestoreClient,
		s.firestoreClient.Collection("jokes").Query.Where("theme_id", "==", themeDoc.Ref.ID),
		time.Duration(0),
	)
	jokes, jokeDocs, err := jokeGetter.GetRandomDocuments(ctx, 2)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get random jokes: %v", err)
	}
	leftJoke, leftJokeDoc := jokes[0], jokeDocs[0]
	rightJoke, rightJokeDoc := jokes[1], jokeDocs[1]
	s.logger.Debug("GetChoices",
		zap.String("left_joke_id", leftJokeDoc.Ref.ID),
		zap.String("left_joke_text", leftJoke.Text),
		zap.String("right_joke_id", rightJokeDoc.Ref.ID),
		zap.String("right_joke_text", rightJoke.Text),
	)

	id := uuid.New().String()
	noWinner := choicesv1.Winner_UNSPECIFIED
	choice := Choice{
		SessionID: req.SessionId,

		ThemeID:     themeDoc.Ref.ID,
		LeftJokeID:  leftJokeDoc.Ref.ID,
		RightJokeID: rightJokeDoc.Ref.ID,
		CreatedAt:   time.Now(),

		Winner: &noWinner,
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
			Value: req.Winner.Number(),
		},
		{
			Path:  "known",
			Value: req.Known.Number(),
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

func (s *Server) GetLeaderboard(
	ctx context.Context,
	req *choicesv1.GetLeaderboardRequest,
) (*choicesv1.GetLeaderboardResponse, error) {
	leaderboardCollection := s.firestoreClient.Collection("leaderboard")
	query := leaderboardCollection.OrderBy("created_at", firestore.Desc).Limit(1)
	docIter := query.Documents(ctx)
	docSnap, err := docIter.Next()
	if err == iterator.Done {
		return nil, status.Error(codes.NotFound, "No leaderboard found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	var leaderboardDoc struct {
		Leaderboard []leaderboardEntry `firestore:"leaderboard"`
		CreatedAt   time.Time          `firestore:"created_at"`
	}
	if err := docSnap.DataTo(&leaderboardDoc); err != nil {
		return nil, fmt.Errorf("failed to parse leaderboard document: %w", err)
	}

	// Map the entries to choicesv1.LeaderboardEntry
	entries := make([]*choicesv1.LeaderboardEntry, 0, len(leaderboardDoc.Leaderboard))
	for _, entryData := range leaderboardDoc.Leaderboard {
		entry := &choicesv1.LeaderboardEntry{
			Model:         entryData.Model,
			Votes:         uint64(entryData.Votes),
			NewmanScore:   entryData.NewmanScore,
			NewmanCILower: entryData.NewmanCiLower,
			NewmanCIUpper: entryData.NewmanCiUpper,
			EloScore:      entryData.EloScore,
			EloCILower:    entryData.EloCiLower,
			EloCIUpper:    entryData.EloCiUpper,
		}

		entries = append(entries, entry)
	}

	return &choicesv1.GetLeaderboardResponse{
		Entries: entries,
	}, nil
}
