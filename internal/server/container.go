package server

import "victorina/internal/model"

type container struct {
	Questions *questions
	Game      *game
}

type questions struct {
	List []questionSchema
	Form questionSchema
}

type game struct {
	Scores  []int
	Current questionSchema
}

func containerWithQuestions(qs []model.Question) container {
	return container{
		Questions: &questions{
			List: questionsToSchema(qs),
		},
	}
}

func containerWithGame(scores []int) container {
	return container{
		Game: &game{
			Scores: scores,
		},
	}
}
