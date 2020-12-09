package routing

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	l "github.com/sirupsen/logrus"

	"adventquest/model"
	"adventquest/response"
)

func (app *Application) GoToTaskHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("log").(*l.Entry)

	day, err := strconv.Atoi(mux.Vars(r)["day"])
	if err != nil {
		log.Error("failed to parse the day")
		response.BadRequest(w, response.Err("the day is invalid", "invalid_day"))
		return
	}

	task, err := app.getTask(day)
	if err != nil {
		log.Error("failed to fetch the task")
		response.InternalError(w, response.Err("can't fetch the day's task", "day_fetching_error"))
		return
	}
	if task == nil {
		response.NotFound(w, response.Err("the day's task not found", "day_not_found"))
		return
	}

	http.Redirect(w, r, task.Link, 301)
}

func (app *Application) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

}

// TODO: pass PG here
func (app *Application) getTask(day int) (*model.Task, error) {
	rows, err := app.Pg.Query("SELECT day, link FROM tasks WHERE day = $1", day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	task := model.Task{}
	if rows.Next() {
		if err := rows.Scan(&task.Day, &task.Link); err != nil {
			return nil, err
		}
		return &task, nil
	}
	return nil, nil
}
