package newmanrank

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNewmanRank_NoTies(t *testing.T) {
	winData := []float64{
		0, 3, 2,
		1, 0, 4,
		0, 0, 0,
	}
	winMatrix := mat.NewDense(3, 3, winData)

	scores, vOut, iterations, err := NewmanRank(winMatrix, nil, 0, 1e-6, 1000)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expectedScores := []float64{5, 1.66, 0}
	expectedV := 0.0

	for i, score := range scores {
		if math.Abs(score-expectedScores[i]) > 1e-1 {
			t.Errorf("Score[%d]: expected %v, got %v", i, expectedScores[i], score)
		}
	}

	if math.Abs(vOut-expectedV) > 1e-6 {
		t.Errorf("Expected vOut %v, got %v", expectedV, vOut)
	}

	t.Logf("Converged in %d iterations", iterations)
}

func TestNewmanRank_WithTies(t *testing.T) {
	winData := []float64{
		0, 3, 2,
		1, 0, 4,
		0, 0, 0,
	}
	tieData := []float64{
		0, 1, 0,
		1, 0, 0,
		2, 0, 0,
	}
	winMatrix := mat.NewDense(3, 3, winData)
	tieMatrix := mat.NewDense(3, 3, tieData)

	scores, vOut, iterations, err := NewmanRank(winMatrix, tieMatrix, 1.0, 1e-2, 1000)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expectedScores := []float64{1.73, 1.55, 0.29}
	expectedV := 0.25

	for i, score := range scores {
		if math.Abs(score-expectedScores[i]) > 1e-1 {
			t.Errorf("Score[%d]: expected %v, got %v", i, expectedScores[i], score)
		}
	}

	if math.Abs(vOut-expectedV) > 1e-1 {
		t.Errorf("Expected vOut %v, got %v", expectedV, vOut)
	}

	t.Logf("Converged in %d iterations", iterations)
}

func TestNewmanRank_LimitExceeded(t *testing.T) {
	winData := []float64{
		0, 1,
		0, 0,
	}
	winMatrix := mat.NewDense(2, 2, winData)

	_, _, _, err := NewmanRank(winMatrix, nil, 0, 1e-12, 1)
	if err == nil {
		t.Fatalf("Expected an error due to non-convergence")
	}
}

// https://github.com/dustalov/evalica
func TestNewmanRank_FoodPreferences_WithBuildMatrices(t *testing.T) {
	// Convert CSV data into comparisons
	comparisonData := []Comparison{
		{"Pizza", "Sushi", LeftWinner},
		{"Burger", "Pasta", RightWinner},
		{"Tacos", "Pizza", LeftWinner},
		{"Sushi", "Tacos", RightWinner},
		{"Burger", "Pizza", LeftWinner},
		{"Pasta", "Sushi", RightWinner},
		{"Tacos", "Burger", LeftWinner},
		{"Pizza", "Pasta", TieWinner},
		{"Sushi", "Burger", RightWinner},
		{"Pasta", "Tacos", LeftWinner},
		{"Burger", "Sushi", RightWinner},
		{"Tacos", "Pizza", LeftWinner},
		{"Sushi", "Pasta", TieWinner},
		{"Burger", "Tacos", LeftWinner},
		{"Pizza", "Burger", RightWinner},
		{"Pasta", "Tacos", RightWinner},
		{"Burger", "Sushi", RightWinner},
		{"Tacos", "Pasta", LeftWinner},
		{"Pizza", "Tacos", RightWinner},
		{"Sushi", "Burger", LeftWinner},
		{"Pasta", "Pizza", RightWinner},
		{"Burger", "Sushi", RightWinner},
		{"Tacos", "Pasta", LeftWinner},
		{"Pizza", "Sushi", LeftWinner},
		{"Sushi", "Burger", RightWinner},
		{"Pasta", "Tacos", RightWinner},
		{"Tacos", "Pizza", LeftWinner},
		{"Burger", "Pasta", RightWinner},
		{"Sushi", "Tacos", LeftWinner},
		{"Pizza", "Burger", RightWinner},
	}

	// Expected scores
	expectedScores := map[string]float64{
		"Tacos":  2.6652107229187254,
		"Sushi":  1.0906271591531849,
		"Burger": 0.8296600883232509,
		"Pasta":  0.7101537820870093,
		"Pizza":  0.536812511441843,
	}

	// Build matrices
	winMatrix, tieMatrix, _, indexToName, err := BuildMatrices(comparisonData)
	if err != nil {
		t.Fatalf("Error building matrices: %v", err)
	}

	// Run NewmanRank
	scores, vOut, iterations, err := NewmanRank(winMatrix, tieMatrix, 1.0, 1e-3, 10000)
	if err != nil {
		t.Fatalf("Error computing rankings: %v", err)
	}

	// Map indices back to item names
	for idx, score := range scores {
		item := indexToName[idx]
		expectedScore := expectedScores[item]
		if math.Abs(score-expectedScore)/expectedScore >= 0.1 {
			t.Errorf("Score for %s: expected %v, got %v", item, expectedScore, score)
		} else {
			t.Logf("Score for %s: expected %v, got %v", item, expectedScore, score)
		}
	}

	t.Logf("Estimated ν: %v", vOut)
	t.Logf("Converged in %d iterations", iterations)
}
