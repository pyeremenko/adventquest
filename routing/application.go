package routing

import (
	"database/sql"
	"errors"
	"net/http"
	"net/textproto"

	"github.com/gorilla/mux"
)

var (
	AuthHeader = "AUTHORIZATION"
)

type Application struct {
	Pg         *sql.DB
	SuperToken string
}

func (app *Application) Run() *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(false)

	for _, route := range app.routes() {
		muxRouter.
			Methods(append(route.Method, "OPTIONS")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(app.middleware(route, app.ChainFilters(app.DefaultMiddleware, route.ActionFilters...)...))
	}

	// TODO:
	// muxRouter.HandleFunc("/", app.health)

	return muxRouter
}

func (app *Application) ChainFilters(filter ActionFilter, filters ...ActionFilter) []ActionFilter {
	result := make([]ActionFilter, len(filters)+1)
	result[0] = filter
	copy(result[1:], filters[:])
	return result
}

func (app *Application) GetAuthToken(request *http.Request) (string, error) {
	prioritizedKeys := []string{
		AuthHeader, textproto.CanonicalMIMEHeaderKey(AuthHeader),
	}

	for _, headerKey := range prioritizedKeys {
		tokens, ok := request.Header[headerKey]
		if ok && len(tokens) >= 1 {
			return tokens[0], nil
		}
	}

	return "", errors.New("unable to find an authorization header")
}

func (app *Application) middleware(route Route, middlewares ...ActionFilter) http.Handler {
	var h http.Handler = route.HandlerFunc
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		h = mw(h, route)
	}
	return h
}
