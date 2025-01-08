package model

import (
	"math/rand"
	"slices"

	"github.com/google/uuid"
)

// Game stores process state.
type Game struct {
	scores map[int][]uuid.UUID
}

// Reset resets game process state to initial.
func (g *Game) Reset(questions []AssignedQuestion) {
	g.scores = make(map[int][]uuid.UUID)
	for _, question := range questions {
		g.scores[question.Score] = append(g.scores[question.Score], question.ID)
	}
}

// NextQuestion returns next question with given score.
func (g *Game) NextQuestion(score int) (uuid.UUID, int, error) {
	switch x := len(g.scores[score]); x {
	case 0:
		return uuid.Nil, 0, ErrNotFound
	case 1:
		out := g.scores[score][0]
		delete(g.scores, score)
		return out, 0, nil
	default:
		i := rand.Intn(x)
		out := g.scores[score][i]
		g.scores[score] = append(g.scores[score][:i], g.scores[score][i+1:]...)
		return out, x - 1, nil
	}
}

// Scores returns sorted deduplicated list of questions' scores.
func (g *Game) Scores() []int {
	out := make([]int, 0, len(g.scores))
	for score := range g.scores {
		out = append(out, score)
	}
	slices.Sort(out)
	return out
}
