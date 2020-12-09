package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	l "github.com/sirupsen/logrus"

	"adventquest/model"
	"adventquest/response"
)

type ActionFilter func(http.Handler, Route) http.Handler

func (app *Application) DefaultMiddleware(next http.Handler, route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		request := fmt.Sprintf("%s %s %s", req.Method, route.Name, req.RequestURI)
		log := l.WithFields(l.Fields{"request": request})

		log.Info("started to handle an endpoint")

		if origin := req.Header.Get("Origin"); origin != "" {
			allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Request-ID"
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		}

		if req.Method == "OPTIONS" {
			return
		}

		ctxWithFields := context.WithValue(req.Context(), "log", log)
		rWithLog := req.WithContext(ctxWithFields)

		next.ServeHTTP(w, rWithLog)
	})
}

func (app *Application) AuthMiddleware(next http.Handler, route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := app.GetAuthToken(req)
		if err != nil {
			l.WithError(err).Debug("invalid auth token")
			response.Unauthorized(w, response.Err("failed to authorize", "invalid_token"))
			return
		}
		ctxWithToken := context.WithValue(req.Context(), "token", token)
		rWithToken := req.WithContext(ctxWithToken)

		next.ServeHTTP(w, rWithToken)
	})
}

func (app *Application) AuthAdminMiddleware(next http.Handler, route Route) http.Handler {
	return app.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := req.Context().Value("token").(string)
		if !ok {
			response.Unauthorized(w, response.Err("failed to authorize", "unauthorized"))
			return
		}

		if token != "sup3rT0k3n" {
			response.Unauthorized(w, response.Err("not enough privileges", "not_admin"))
			return
		}

		ctxWithAdmin := context.WithValue(req.Context(), "admin", true)
		rWithAdmin := req.WithContext(ctxWithAdmin)

		next.ServeHTTP(w, rWithAdmin)
	}), route)
}

func (app *Application) NullableInputMiddleware(factory model.InputModelFactory) ActionFilter {
	return func(next http.Handler, route Route) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log := req.Context().Value("log").(*l.Entry)
			input := factory.Create()

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.WithError(err).Error("failed to get input")
				response.BadRequest(w, response.Err(err.Error(), "failed_to_get_input"))
				return
			}
			if len(body) > 0 {
				err = json.Unmarshal(body, input)
				if err != nil {
					log.WithError(err).Error("failed to parse input")
					response.BadRequest(w, response.Err(err.Error(), "failed_to_parse_input"))
					return
				}
			}

			ctxWithJwt := context.WithValue(req.Context(), "input", input)
			rWithJwt := req.WithContext(ctxWithJwt)

			next.ServeHTTP(w, rWithJwt)
		})
	}
}

func (app *Application) InputMiddleware(factory model.InputModelFactory) ActionFilter {
	return func(next http.Handler, route Route) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log := req.Context().Value("log").(*l.Entry)

			input := factory.Create()
			err := json.NewDecoder(req.Body).Decode(input)
			if err != nil {
				log.WithError(err).Error("failed to parse input")
				response.BadRequest(w, response.Err(err.Error(), "failed_to_parse_input"))
				return
			}

			ctxWithJwt := context.WithValue(req.Context(), "input", input)
			rWithJwt := req.WithContext(ctxWithJwt)

			next.ServeHTTP(w, rWithJwt)
		})
	}
}
