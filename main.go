package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"adventquest/response"
)

type Task struct {
	Day  int
	Link string
}

var pg *sql.DB

var links = []string{
	"",
	"",
	"https://user-images.githubusercontent.com/696437/49771691-7fcb4000-fd14-11e8-9a3c-358129b18a7d.png",
	"https://user-images.githubusercontent.com/696437/100855679-a6ac4f00-34b4-11eb-9fd2-260c38047795.JPG",
	"https://user-images.githubusercontent.com/696437/100976900-63acb300-356a-11eb-8920-2d7f15b7882b.JPG",
	"https://user-images.githubusercontent.com/696437/100976969-7fb05480-356a-11eb-93d4-21a9da6b9318.JPG",
	"https://user-images.githubusercontent.com/696437/100977011-8d65da00-356a-11eb-8c22-45c4130c4b17.JPG",
	"https://user-images.githubusercontent.com/696437/100855679-a6ac4f00-34b4-11eb-9fd2-260c38047795.JPG",
	"https://user-images.githubusercontent.com/696437/100976900-63acb300-356a-11eb-8920-2d7f15b7882b.JPG",
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	day, err := strconv.Atoi(mux.Vars(r)["day"])
	if err != nil {
		response.BadRequest(w, response.Err("the day is invalid", "invalid_day"))
		return
	}

	if day <= 0 || day >= len(links) {
		response.NotFound(w, response.Err("the day not found", "day_not_found"))
		return
	}

	task, err := getTask(day)
	if err != nil {
		response.InternalError(w, response.Err("can't fetch the day", "day_fetching_error"))
		return
	}
	if task == nil {
		response.NotFound(w, response.Err("the day not found", "day_not_found"))
		return
	}

	http.Redirect(w, r, task.Link, 301)
}

func main() {
	err := connect()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/go/{day}", mainHandler)
	http.ListenAndServe(getPort(), r)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func getConnectionString() string {
	var connStringVariableName = os.Getenv("CONN_STRING")
	if connStringVariableName == "" {
		return "user=developer dbname=adventquest password=developer host=localhost port=5432 sslmode=disable"
	}

	return os.Getenv(connStringVariableName)
}

func connect() error {
	var err error

	pg, err = sql.Open("postgres", getConnectionString())
	return err
}

func getTask(day int) (*Task, error) {
	rows, err := pg.Query("SELECT day, link FROM tasks WHERE day = $1", day)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	task := Task{}
	if rows.Next() {
		if err := rows.Scan(&task.Day, &task.Link); err != nil {
			return nil, err
		}
		return &task, nil
	}
	return nil, nil
}
