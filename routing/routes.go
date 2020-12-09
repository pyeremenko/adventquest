package routing

import (
	"adventquest/model"
)

func (app *Application) routes() []Route {
	return Routes{
		Route{
			"CreateTaskHandler", []string{"POST"}, "/task",
			app.CreateTaskHandler,
			[]ActionFilter{
				app.AuthAdminMiddleware,
				app.InputMiddleware(new(model.CreateTaskInputFactory)),
			},
		},
		Route{
			"GoToTask", []string{"GET"}, "/go/{day}",
			app.GoToTaskHandler,
			[]ActionFilter{},
		},
	}
}
