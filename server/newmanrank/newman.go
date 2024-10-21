package newmanrank

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// NewmanRank computes rankings from pairwise comparisons using the algorithm
// described by M. E. J. Newman in "Efficient computation of rankings from pairwise comparisons".
// The function handles both cases with and without ties.
//
// Parameters:
//   - winMatrix: N x N matrix where winMatrix[i][j] is the number of times i beats j.
//   - tieMatrix: N x N matrix where tieMatrix[i][j] is the number of ties between i and j.
//     If there are no ties, tieMatrix can be nil.
//   - v: Initial value for the parameter ν (nu) in the algorithm (use 0 if no ties).
//   - tolerance: Convergence tolerance.
//   - limit: Maximum number of iterations.
//
// Returns:
// - scores: Slice of ranking scores for each player.
// - vOut: Estimated value of ν (nu) parameter.
// - iterations: Number of iterations performed.
// - err: Error if computation fails.
func NewmanRank(
	winMatrix, tieMatrix *mat.Dense,
	v, tolerance float64,
	limit int,
) (scores []float64, vOut float64, iterations int, err error) {
	n, _ := winMatrix.Dims()
	if tieMatrix != nil {
		nt, _ := tieMatrix.Dims()
		if nt != n {
			return nil, 0, 0, errors.New("winMatrix and tieMatrix must have the same dimensions")
		}
	}

	if n == 0 {
		// Return default scores if matrices are empty
		scores = make([]float64, n)
		for i := range scores {
			scores[i] = 1.0
		}
		vOut = v
		return
	}

	// Initialize winTieHalf = winMatrix + tieMatrix / 2
	winTieHalf := mat.NewDense(n, n, nil)
	if tieMatrix != nil {
		winTieHalf.Apply(func(i, j int, _ float64) float64 {
			return winMatrix.At(i, j) + tieMatrix.At(i, j)/2.0
		}, winMatrix)
	} else {
		winTieHalf.CloneFrom(winMatrix)
	}

	// Initialize scores to ones
	scoresVec := mat.NewVecDense(n, nil)
	for i := 0; i < n; i++ {
		scoresVec.SetVec(i, 1.0)
	}

	converged := false
	iterations = 0
	scoresNewVec := mat.NewVecDense(n, nil)
	vNew := v

	for !converged && iterations < limit {
		iterations++
		if math.IsNaN(vNew) {
			v = tolerance
		} else {
			v = vNew
		}

		// Prepare variables
		scoresData := scoresVec.RawVector().Data

		// Compute sqrtScoresOuter = sqrt(scores_i * scores_j)
		sqrtScoresOuter := mat.NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				sqrtScoresOuter.Set(i, j, math.Sqrt(scoresData[i]*scoresData[j]))
			}
		}

		// Compute sumScores = scores_i + scores_j
		sumScores := mat.NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				sumScores.Set(i, j, scoresData[i]+scoresData[j])
			}
		}

		// Compute sqrtDivScoresOuterT = sqrt(scores_i / scores_j)
		sqrtDivScoresOuterT := mat.NewDense(n, n, nil)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if scoresData[j] == 0 {
					sqrtDivScoresOuterT.Set(i, j, 0)
				} else {
					sqrtDivScoresOuterT.Set(i, j, math.Sqrt(scoresData[i]/scoresData[j]))
				}
			}
		}

		// Compute commonDenominator = sumScores + 2 * v * sqrtScoresOuter
		commonDenominator := mat.NewDense(n, n, nil)
		// temp = 2 * v * sqrtScoresOuter
		temp := mat.NewDense(n, n, nil)
		temp.Scale(2*v, sqrtScoresOuter)
		// commonDenominator = sumScores + temp
		commonDenominator.Add(sumScores, temp)

		// Update scores
		for i := 0; i < n; i++ {
			numSum := 0.0
			denSum := 0.0
			for j := 0; j < n; j++ {
				denominator := commonDenominator.At(i, j)
				if denominator == 0 {
					continue
				}
				// Numerator for scores update
				num := winTieHalf.At(i, j) * (scoresData[j] + v*sqrtScoresOuter.At(i, j)) / denominator
				numSum += num

				// Denominator for scores update
				den := winTieHalf.At(j, i) * (1 + v*sqrtDivScoresOuterT.At(i, j)) / commonDenominator.At(j, i)
				denSum += den
			}
			if denSum == 0 {
				scoresNewVec.SetVec(i, tolerance)
			} else {
				val := numSum / denSum
				if math.IsNaN(val) || math.IsInf(val, 0) {
					val = tolerance
				}
				scoresNewVec.SetVec(i, val)
			}
		}

		// Update v if ties are present
		if tieMatrix != nil {
			vNumerator := 0.0
			vDenominator := 0.0
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					denominator := commonDenominator.At(i, j)
					if denominator == 0 {
						continue
					}
					vNumerator += tieMatrix.At(i, j) * sumScores.At(i, j) / denominator
					vDenominator += winMatrix.At(i, j) * sqrtScoresOuter.At(i, j) / denominator
				}
			}
			vNumerator /= 2.0
			vDenominator *= 2.0
			if vDenominator == 0 {
				vNew = tolerance
			} else {
				vNew = vNumerator / vDenominator
				if math.IsNaN(vNew) || math.IsInf(vNew, 0) {
					vNew = tolerance
				}
			}
		}

		// Check for convergence
		diffVec := mat.NewVecDense(n, nil)
		diffVec.SubVec(scoresNewVec, scoresVec)
		diffNorm := mat.Norm(diffVec, 2)
		converged = diffNorm < tolerance

		// Update scoresVec
		scoresVec.CloneFromVec(scoresNewVec)
	}

	// Extract scores
	scores = make([]float64, n)
	copy(scores, scoresVec.RawVector().Data)
	vOut = v
	if iterations == limit && !converged {
		err = errors.New("did not converge within the specified limit")
	}
	return
}

// WinnerType represents the outcome of a comparison.
type WinnerType int

const (
	LeftWinner WinnerType = iota
	RightWinner
	TieWinner
)

func (w WinnerType) String() string {
	switch w {
	case LeftWinner:
		return "left"
	case RightWinner:
		return "right"
	case TieWinner:
		return "tie"
	default:
		return "unknown"
	}
}

// Comparison represents a single pairwise comparison.
type Comparison struct {
	Left   string
	Right  string
	Winner WinnerType
}

// BuildMatrices constructs the winMatrix and tieMatrix from a slice of Comparisons.
// It returns the matrices along with mappings between player names and indices.
func BuildMatrices(comparisons []Comparison) (winMatrix, tieMatrix *mat.Dense, nameToIndex map[string]int, indexToName []string, err error) {
	// Collect all unique player names
	playerSet := make(map[string]struct{})
	for _, comp := range comparisons {
		playerSet[comp.Left] = struct{}{}
		playerSet[comp.Right] = struct{}{}
	}

	// Create mappings between player names and indices
	indexToName = make([]string, 0, len(playerSet))
	for name := range playerSet {
		indexToName = append(indexToName, name)
	}
	sort.Strings(indexToName) // Sort for consistency

	nameToIndex = make(map[string]int, len(indexToName))
	for idx, name := range indexToName {
		nameToIndex[name] = idx
	}

	n := len(indexToName)
	winData := make([]float64, n*n)
	tieData := make([]float64, n*n)

	for _, comp := range comparisons {
		i, okI := nameToIndex[comp.Left]
		j, okJ := nameToIndex[comp.Right]
		if !okI || !okJ {
			return nil, nil, nil, nil, fmt.Errorf("unknown player name in comparison: %v vs %v", comp.Left, comp.Right)
		}

		if i == j {
			return nil, nil, nil, nil, fmt.Errorf("comparison between the same player: %v", comp.Left)
		}

		switch comp.Winner {
		case LeftWinner:
			winData[i*n+j] += 1
		case RightWinner:
			winData[j*n+i] += 1
		case TieWinner:
			tieData[i*n+j] += 1
			tieData[j*n+i] += 1
		default:
			return nil, nil, nil, nil, fmt.Errorf("invalid winner value: %v", comp.Winner)
		}
	}

	winMatrix = mat.NewDense(n, n, winData)

	// Only create tieMatrix if there are ties
	hasTies := false
	for _, v := range tieData {
		if v != 0 {
			hasTies = true
			break
		}
	}
	if hasTies {
		tieMatrix = mat.NewDense(n, n, tieData)
	} else {
		tieMatrix = nil
	}

	return winMatrix, tieMatrix, nameToIndex, indexToName, nil
}
