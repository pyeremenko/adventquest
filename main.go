package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"adventquest/routing"
)

func main() {
	pg, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	app := &routing.Application{Pg: pg, SuperToken: getAdminToken()}

	http.ListenAndServe(getPort(), app.Run())
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info("no PORT environment variable detected, defaulting to " + port)
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

func getAdminToken() string {
	var token = os.Getenv("SUPERTOKEN")
	if token == "" {
		token = "supertoken"
		log.Info("no SUPERTOKEN environment variable detected, defaulting to " + token)
	}
	return token
}

func connect() (*sql.DB, error) {
	connectionString := getConnectionString()
	log.Infof("Connecting to postgres: %v\n", connectionString)
	return sql.Open("postgres", connectionString)
}
