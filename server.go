package main

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func initServer(host string, useLog bool) {
	tokenAuth = jwtauth.New("HS256", sKey, nil)

	r := chi.NewRouter()

	if useLog {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsHandler)

	// Frontend
	r.Get("/", indexHandler)
	r.Get("/favicon.ico", serveFileHandler)
	FileServer(r, "/static", http.Dir(filepath.Join("public", "static")))
	r.NotFound(indexHandler)

	// Auth
	r.Group(func(r chi.Router) {
		r.Post("/login", login)
	})

	// REST API
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Use(render.SetContentType(render.ContentTypeJSON))

		r.Route("/api/v1/proxies", func(r chi.Router) {
			r.Get("/all", listProxies)
			r.Get("/work", listWorkProxies)
			r.Get("/anon", listAnonProxies)
			r.Get("/counts", getCounts)
		})

		r.Route("/api/v1/proxy", func(r chi.Router) {
			r.Get("/{id}", getProxy)
		})
	})

	err := http.ListenAndServe(host, r)
	errmsg("ListenAndServe", err)
}
