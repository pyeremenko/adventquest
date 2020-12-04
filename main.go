package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"adventquest/response"
)

var links = []string{
	"",
	"",
	"https://user-images.githubusercontent.com/696437/49771691-7fcb4000-fd14-11e8-9a3c-358129b18a7d.png",
	"https://user-images.githubusercontent.com/696437/100855679-a6ac4f00-34b4-11eb-9fd2-260c38047795.JPG",
	"https://user-images.githubusercontent.com/696437/100976900-63acb300-356a-11eb-8920-2d7f15b7882b.JPG",
	"https://user-images.githubusercontent.com/696437/100976969-7fb05480-356a-11eb-93d4-21a9da6b9318.JPG",
	"https://user-images.githubusercontent.com/696437/100977011-8d65da00-356a-11eb-8c22-45c4130c4b17.JPG",
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

	http.Redirect(w, r, links[day], 301)
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
