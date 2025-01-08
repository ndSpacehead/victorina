package server

import "victorina/internal/model"

type container struct {
	Questions *questions
	Scenarios *scenarios
	Game      *game
}

type questions struct {
	List []questionSchema
	Form questionSchema
}

type scenarios struct {
	List      []scenarioSchema
	Form      scenarioSchema
	Questions scenariosQuestions
}

type scenariosQuestions struct {
	Form         scenarioQuestionSchema
	AssignedList []assignedQuestionSchema
	FreeList     []questionSchema
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

func containerWithScenario(scs []model.Scenario) container {
	return container{
		Scenarios: &scenarios{
			List: scenariosToSchema(scs),
		},
	}
}

func containerWithGame(name string, scores []int) container {
	return container{
		Game: &game{
			Scores: scores,
		},
	}
}
