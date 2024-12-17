package server

import "victorina/internal/model"

type container struct {
	Questions *questions
}

type questions struct {
	List []questionSchema
	Form questionSchema
}

func containerWithQuestions(qs []model.Question) container {
	return container{
		Questions: &questions{
			List: questionsToSchema(qs),
		},
	}
}
