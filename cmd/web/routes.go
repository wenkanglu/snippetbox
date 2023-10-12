package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wenkanglu/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recoverPanic)
	// TODO: adding metrics fails the tests - fix!
	// r.Use(app.addRequestMetrics)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))
	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowed(w)
	}))

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fileServer)

	r.Get("/ping", ping)

	r.Group(func(r chi.Router) {
		r.Use(noSurf)
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(app.authenticate)

		r.Get("/", app.home)

		r.Get("/snippet/view/{id}", app.snippetView)

		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)

		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/snippet/create", app.snippetCreate)
			r.Post("/snippet/create", app.snippetCreatePost)

			r.Post("/user/logout", app.userLogoutPost)
		})
	})

	return r
}

func (app *application) metricsRoutes(handler http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Handle("/metrics", handler)

	return r
}
