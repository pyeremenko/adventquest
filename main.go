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

func mainHandler(w http.ResponseWriter, r *http.Request) {
	day, err := strconv.Atoi(mux.Vars(r)["day"])
	if err != nil {
		response.BadRequest(w, response.Err("the day is invalid", "invalid_day"))
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

	connectionString := getConnectionString()
	log.Printf("Connecting to postgres: %v\n", connectionString)
	pg, err = sql.Open("postgres", connectionString)
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
