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

	expectedScores := []float64{1.0, 2.0, 0.5}
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

	scores, vOut, iterations, err := NewmanRank(winMatrix, tieMatrix, 1.0, 1e-6, 1000)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expectedScores := []float64{1.0, 2.0, 0.5}
	expectedV := 0.5

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
