package main

import (
	"log"
	"net/http"
	"os"
)

const TheOnlyLink = "https://user-images.githubusercontent.com/696437/49771691-7fcb4000-fd14-11e8-9a3c-358129b18a7d.png"

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, TheOnlyLink, 301)
}

func main() {
	http.HandleFunc("/go/2", handler)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
