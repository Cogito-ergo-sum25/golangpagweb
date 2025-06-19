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
	mux.Post("/nueva-licitacion",handlers.Repo.CrearNuevaLicitacion)
	mux.Get("/editar-licitacion/{id}", handlers.Repo.MostrarFormularioEditarLicitacion)
	mux.Post("/editar-licitacion/{id}", handlers.Repo.EditarLicitacion)


	// PARTIDAS
	mux.Get("/mostrar-partidas/{id}",handlers.Repo.MostrarPartidasPorID)
		mux.Get("/nueva-partida/{id}", handlers.Repo.MostrarNuevaPartida)
		mux.Post("/nueva-partida",handlers.Repo.CrearNuevaPartida)
		mux.Get("/editar-partida/{id}", handlers.Repo.MostrarEditarPartida)
		mux.Post("/editar-partida/{id}", handlers.Repo.EditarPartida)

	// ACLARACIONES-LICITACION
	mux.Get("/aclaraciones-licitacion/{id}", handlers.Repo.MostrarAclaracionesLicitacion)
	mux.Get("/nueva-aclaracion-general/{id}", handlers.Repo.MostrarNuevaAclaracionGeneral)
	mux.Post("/datos-empresas-externas-nueva-contexto-aclaraciones", handlers.Repo.AgregarEmpresaExternaContextoAclaraciones)
	mux.Post("/nueva-aclaracion-licitacion", handlers.Repo.CrearNuevaAclaracionGeneral)

	
	// PRODUCTOS PARTIDA
	mux.Get("/productos-partida/{id}", handlers.Repo.MostrarProductosPartida)
	mux.Get("/nuevo-producto-partida/{id}", handlers.Repo.MostrarNuevoProductoPartida)
	mux.Post("/nuevo-producto-partida", handlers.Repo.CrearNuevoProductoPartida)
	mux.Post("/editar-producto-partida", handlers.Repo.EditarProductoPartida)
	mux.Post("/eliminar-producto-partida/{id}", handlers.Repo.EliminarProductoPartida)




	// ACLARACIONES
	mux.Get("/aclaraciones/{id}", handlers.Repo.MostrarAclaraciones)
	mux.Get("/nueva-aclaracion/{id}", handlers.Repo.MostrarNuevaAclaracion)
	mux.Post("/datos-empresas-externas-nueva-contexto", handlers.Repo.AgregarEmpresaExternaContexto)
	mux.Post("/nueva-aclaracion", handlers.Repo.CrearNuevaAclaracion)

	// REQUERIMIENTOS
	mux.Get("/requerimientos/{id}", handlers.Repo.ObtenerRequerimientos)
	mux.Get("/requerimientos-json/{id}", handlers.Repo.ObtenerRequerimientosJSON)
	mux.Post("/guardar-requerimientos", handlers.Repo.GuardarRequerimientos)

	// PROPUESTAS
	mux.Get("/propuestas/{id}", handlers.Repo.MostrarPropuestas)
	mux.Get("/nueva-propuesta/{id}", handlers.Repo.MostrarNuevaPropuesta)
	mux.Post("/nueva-propuesta/{id}", handlers.Repo.CrearNuevaPropuesta)
	mux.Post("/nuevo-producto-externo-contexto", handlers.Repo.NuevoProductoExternoContexto)
	mux.Get("/editar-propuesta/{id}", handlers.Repo.MostrarEditarPropuesta)
	mux.Post("/editar-propuesta/{id}", handlers.Repo.EditarPropuesta)

	// FALLOS
	mux.Get("/fallo/{id}", handlers.Repo.ObtenerFalloPropuesta)
	mux.Get("/fallo-json/{id}", handlers.Repo.ObtenerFalloPropuestaJSON)
	mux.Post("/guardar-fallo", handlers.Repo.GuardarFalloPropuesta)







	
	

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
	// OPCIONES DE EMPRESAS EXTERNAS Y PRODUCTOS COMPETENCIA
		mux.Get("/datos-empresas-externas",handlers.Repo.EmpresasExternas)
		mux.Post("/datos-empresas-externas-nueva", handlers.Repo.AgregarEmpresaExterna)
		mux.Post("/crear-entidad", handlers.Repo.CrearEntidad)
		mux.Get("/productos-externos",handlers.Repo.ProductosExternos)
		mux.Post("/nuevo-producto-externo-contexto-menu", handlers.Repo.NuevoProductoExternoContextoMenu)
	

	    	
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
}