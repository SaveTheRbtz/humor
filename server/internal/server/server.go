package server

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"

	insecureRandExp "golang.org/x/exp/rand"

	choicesv1 "github.com/SaveTheRbtz/humor/gen/go/proto"

	"google.golang.org/api/iterator"

	"github.com/google/uuid"

	"gonum.org/v1/gonum/stat/sampleuv"

	// embed
	_ "embed"
)

// TODO(rbtz): automate rank generation and open-source policy geneeration.

//go:embed "top-jokes.txt"
var topJokesData string

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
	Active bool    `firestore:"active"`
}

type Joke struct {
	Theme   string  `firestore:"theme"`
	ThemeID string  `firestore:"theme_id"`
	Text    string  `firestore:"text"`
	Random  float64 `firestore:"random"`
	Model   string  `firestore:"model"`
	Policy  string  `firestore:"policy"`
	Active  bool    `firestore:"active"`
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

var modelWeights struct {
	ModelWeights []float64 `firestore:"model_weights"`
	Shape        []int     `firestore:"shape"`
	Models       []string  `firestore:"models"`
	CreatedAt    time.Time `firestore:"created_at"`
}

type Server struct {
	choicesv1.UnimplementedArenaServer

	logger          *zap.Logger
	firestoreClient *firestore.Client
	rand            *insecureRandExp.Rand
	source          insecureRandExp.Source
	themeGetter     *randomDocumentGetterImpl[Theme]

	modelWeightsCache          [][]float64
	modelWeightsCacheTimestamp time.Time
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
	source := insecureRandExp.NewSource(uint64(time.Now().UnixNano()))

	return &Server{
		logger:          logger,
		firestoreClient: firestoreClient,
		themeGetter:     randomThemeGetter,
		rand:            insecureRandExp.New(source),
		source:          source,
	}, nil
}

func (s *Server) getAllModelWeights(ctx context.Context) ([][]float64, error) {
	if s.modelWeightsCache != nil && time.Since(s.modelWeightsCacheTimestamp) < time.Minute {
		return s.modelWeightsCache, nil
	}

	modelWeightsCollection := s.firestoreClient.Collection("model_weights")
	query := modelWeightsCollection.OrderBy("created_at", firestore.Desc).Limit(1)

	docIter := query.Documents(ctx)
	docSnap, err := docIter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get model weights: %w", err)
	}

	if err := docSnap.DataTo(&modelWeights); err != nil {
		return nil, fmt.Errorf("failed to parse model weights document: %w", err)
	}

	// reshape model weights
	modelWeightsMatrix := make([][]float64, modelWeights.Shape[0])
	for i := range modelWeightsMatrix {
		modelWeightsMatrix[i] = modelWeights.ModelWeights[i*modelWeights.Shape[1] : (i+1)*modelWeights.Shape[1]]
	}

	s.modelWeightsCache = modelWeightsMatrix
	s.modelWeightsCacheTimestamp = time.Now()

	return modelWeightsMatrix, nil
}

func (s *Server) getJokesForModel(ctx context.Context, model string, jokes []Joke, jokeDocs []*firestore.DocumentSnapshot) ([]Joke, []*firestore.DocumentSnapshot, error) {
	modelWeightsMatrix, err := s.getAllModelWeights(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all model weights: %w", err)
	}

	leftModelWeights := make([]float64, 0)
	foundModelID := -1
	for modedID, m := range modelWeights.Models {
		if m == model {
			leftModelWeights = modelWeightsMatrix[modedID]
			foundModelID = modedID
			break
		}
	}
	if foundModelID == -1 {
		return nil, nil, fmt.Errorf("model not found in model weights: %s, %v", model, modelWeights.Models)
	}

	rightModelIdx, ok := sampleuv.NewWeighted(leftModelWeights, s.source).Take()
	if !ok {
		return nil, nil, fmt.Errorf("failed to sample model: %w", err)
	}
	rightModel := modelWeights.Models[rightModelIdx]

	var filteredJokes []Joke
	var filteredJokeDocs []*firestore.DocumentSnapshot
	for i, joke := range jokes {
		if joke.Model != rightModel {
			continue
		}
		if !joke.Active {
			continue
		}
		filteredJokes = append(filteredJokes, joke)
		filteredJokeDocs = append(filteredJokeDocs, jokeDocs[i])
	}
	if len(filteredJokes) < 1 {
		return nil, nil, fmt.Errorf("no jokes found for model: %s", rightModel)
	}

	s.logger.Debug("Filtered jokes count", zap.Int("count", len(filteredJokes)))

	return filteredJokes, filteredJokeDocs, nil
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

	// Get all jokes for a theme
	allJokeDocs, err := s.firestoreClient.Collection("jokes").Query.Where("theme", "==", theme.Text).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get jokes: %w", err)
	}

	// Collect all jokes
	jokes := make([]Joke, 0)
	jokeDocs := make([]*firestore.DocumentSnapshot, 0)
	for _, doc := range allJokeDocs {
		if active_raw, ok := doc.Data()["active"]; ok {
			if active := active_raw.(bool); active {
				var joke Joke
				if err := doc.DataTo(&joke); err != nil {
					return nil, fmt.Errorf("failed to parse joke document: %w", err)
				}
				jokes = append(jokes, joke)
				jokeDocs = append(jokeDocs, doc)
			}
		}
	}
	if len(jokes) < 2 {
		return nil, status.Errorf(codes.NotFound, "Not enough jokes found: %d, theme: %s", len(jokes), theme.Text)
	}
	leftJokeID := s.rand.Int63n(int64(len(jokes)))
	leftJokeDoc := jokeDocs[leftJokeID]
	leftJoke := jokes[leftJokeID]

	filteredJokes, filteredJokeDocs, err := s.getJokesForModel(ctx, leftJoke.Model, jokes, jokeDocs)
	if err != nil {
		filteredJokes = jokes
		filteredJokeDocs = jokeDocs
		s.logger.Warn("GetChoices failed to get jokes for model", zap.String("model", leftJoke.Model), zap.String("theme", theme.Text), zap.Error(err))
	}

	rightJokeID := s.rand.Int63n(int64(len(filteredJokes)))
	rightJokeDoc := filteredJokeDocs[rightJokeID]
	rightJoke := filteredJokes[rightJokeID]

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

func (s *Server) GetTopJokes(
	ctx context.Context,
	req *choicesv1.GetTopJokesRequest,
) (*choicesv1.GetTopJokesResponse, error) {
	entries := make([]*choicesv1.TopJokesEntry, 0, 10)
	// read lines from topJokesData
	for i, line := range strings.Split(topJokesData, "\n") {
		if line == "" {
			continue
		}
		entries = append(entries, &choicesv1.TopJokesEntry{
			Rank: uint64(i + 1),
			Text: line,
		})
	}

	return &choicesv1.GetTopJokesResponse{
		Entries: entries,
	}, nil
}

func (s *Server) Close() error {
	return s.firestoreClient.Close()
}
