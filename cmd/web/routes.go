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


	// INVENTARIO
	mux.Get("/inventario",handlers.Repo.Inventario)
	mux.Get("/crear-producto", handlers.Repo.MostrarFormularioCrear)  // Muestra el formulario
    mux.Post("/crear-producto", handlers.Repo.CrearProducto)        // Procesa el formulario
	mux.Get("/editar-producto/{id}", handlers.Repo.MostrarFormularioEditar)   // Muestra formulario de edición (GET)
	mux.Post("/editar-producto/{id}", handlers.Repo.EditarProducto)            // Procesa la edición (POST)
	mux.Post("/eliminar/{id}", handlers.Repo.EliminarProducto)   // Elimina un producto (POST)

	// PROYECTOS
	mux.Get("/proyectos-vista",handlers.Repo.ProyectosVista)
	mux.Get("/nuevo-proyecto",handlers.Repo.MostrarNuevoProyecto)
	mux.Post("/nuevo-proyecto",handlers.Repo.NuevoProyecto)


	// LICITACIONES
	mux.Get("/licitaciones",handlers.Repo.Licitaciones)
	mux.Get("/nueva-licitacion",handlers.Repo.MostrarNuevaLicitacion)
	

	// CATALOGO 
	//mux.Get("/catalogo", handlers.Repo.Catalogo) 
	//mux.Get("/producto/{id}", handlers.Repo.VerProducto)// Para la vista detallada


	// OPCIONES
	mux.Get("/opciones",handlers.Repo.Opciones)
	// OPCIONES DE DATOS REFERENCIA
		mux.Get("/datos-referencia",handlers.Repo.DatosReferencia)
		mux.Post("/datos-referencia", handlers.Repo.AgregarDato)
		mux.Post("/eliminar-referencia/{id}", handlers.Repo.EliminarDatoReferencia)   // Elimina un producto (POST)
	// OPCIONES DE ENTIDADES
		mux.Get("/datos-entidades",handlers.Repo.Entidades)
		mux.Get("/crear-entidad",handlers.Repo.MostrarNuevaEntidad)
		mux.Post("/crear-entidad", handlers.Repo.CrearEntidad)
	    	
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
}