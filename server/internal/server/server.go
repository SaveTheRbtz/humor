package server

import (
	"context"
	"fmt"
	"sort"

	"go.uber.org/zap"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	insecureRand "math/rand/v2"
	"time"

	choicesv1 "github.com/SaveTheRbtz/humor/gen/go/proto"

	"github.com/seanhagen/bradleyterry"
	"google.golang.org/api/iterator"

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

func NewServer(firestoreClient *firestore.Client, logger *zap.Logger) (*Server, error) {
	randomThemeGetter, err := NewRandomDocumentGetter[Theme](
		firestoreClient,
		firestoreClient.Collection("themes").Query,
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
	// Step 1: Fetch all choices and collect joke IDs
	allChoicesDocs := s.firestoreClient.Collection("choices").Documents(ctx)
	choices := []Choice{}
	jokeIDs := make(map[string]struct{})

	for {
		doc, err := allChoicesDocs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get choices: %w", err)
		}

		var choice Choice
		if err := doc.DataTo(&choice); err != nil {
			return nil, fmt.Errorf("failed to convert choice to struct: %w", err)
		}

		if choice.Winner == nil {
			continue
		}
		if *choice.Winner != choicesv1.Winner_LEFT && *choice.Winner != choicesv1.Winner_RIGHT {
			continue
		}

		jokeIDs[choice.LeftJokeID] = struct{}{}
		jokeIDs[choice.RightJokeID] = struct{}{}

		choices = append(choices, choice)
	}

	// Step 2: Fetch all jokes in a single batch
	jokeDocRefs := make([]*firestore.DocumentRef, 0, len(jokeIDs))
	for jokeID := range jokeIDs {
		docRef := s.firestoreClient.Collection("jokes").Doc(jokeID)
		jokeDocRefs = append(jokeDocRefs, docRef)
	}

	jokeDocs, err := s.firestoreClient.GetAll(ctx, jokeDocRefs)
	if err != nil {
		return nil, fmt.Errorf("failed to get jokes: %w", err)
	}

	jokeMap := make(map[string]Joke)
	for _, docSnap := range jokeDocs {
		joke := Joke{}
		if err := docSnap.DataTo(&joke); err != nil {
			return nil, fmt.Errorf("failed to convert joke to struct: %w", err)
		}
		jokeMap[docSnap.Ref.ID] = joke
	}

	allPairs := []bradleyterry.Pair{}
	for _, choice := range choices {
		leftJoke, ok := jokeMap[choice.LeftJokeID]
		if !ok {
			continue
		}
		rightJoke, ok := jokeMap[choice.RightJokeID]
		if !ok {
			continue
		}

		var winnerModel, loserModel string
		switch *choice.Winner {
		case choicesv1.Winner_LEFT:
			winnerModel = leftJoke.Model
			loserModel = rightJoke.Model
		case choicesv1.Winner_RIGHT:
			winnerModel = rightJoke.Model
			loserModel = leftJoke.Model
		default:
			continue
		}

		allPairs = append(allPairs, bradleyterry.Pair{
			Winner: winnerModel,
			Loser:  loserModel,
		})
	}
	s.logger.Debug("GetLeaderboard", zap.Int("allPairs", len(allPairs)))

	// map[string]float64 - ModelID -> score
	bt := bradleyterry.Model(allPairs)
	leaderboard := []*choicesv1.LeaderboardEntry{}
	for model, score := range bt {
		leaderboard = append(leaderboard, &choicesv1.LeaderboardEntry{
			Model:            model,
			BradleyterrScore: score,
		})
	}
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].BradleyterrScore > leaderboard[j].BradleyterrScore
	})

	return &choicesv1.GetLeaderboardResponse{
		Entries: leaderboard,
	}, nil
}
