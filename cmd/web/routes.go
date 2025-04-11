package main

import (
	"net/http"

	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	
)

func routes(app *config.AppConfig) http.Handler {
	mux:= chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/",handlers.Repo.Home)
	mux.Get("/catalogo",handlers.Repo.Catalogo)
	mux.Get("/inventario",handlers.Repo.Inventario)
	mux.Get("/crear",handlers.Repo.Crear)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
}