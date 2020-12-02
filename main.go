package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const TheOnlyLink = "https://user-images.githubusercontent.com/696437/49771691-7fcb4000-fd14-11e8-9a3c-358129b18a7d.png"

var links = []string{
	"",
	"",
	"https://user-images.githubusercontent.com/696437/49771691-7fcb4000-fd14-11e8-9a3c-358129b18a7d.png",
	"https://user-images.githubusercontent.com/696437/100855679-a6ac4f00-34b4-11eb-9fd2-260c38047795.JPG",
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	day, _ := strconv.Atoi(mux.Vars(r)["day"])
	//TODO: handle err
	if day > 0 && day < len(links) {
		http.Redirect(w, r, links[day], 301)
	}
}

func main() {
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
