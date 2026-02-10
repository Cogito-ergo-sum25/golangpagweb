package main

import (
	"net/http"

	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	// Middlewares globales que se aplican a TODAS las peticiones
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// --- RUTAS PÚBLICAS ---
	// Estas rutas no requieren que el usuario haya iniciado sesión.
	mux.Get("/login", handlers.Repo.ShowLoginPage)
	mux.Post("/login", handlers.Repo.PostLoginPage)

	// Los archivos estáticos (CSS, JS) también deben ser públicos
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// --- RUTAS PROTEGIDAS ---
	// Usamos un grupo para aplicar el middleware Auth a todas las rutas de la aplicación.
	mux.Group(func(r chi.Router) {
		// ¡Este es nuestro guardián! Se aplica a todo lo que esté dentro de este bloque.
		r.Use(Auth)

		// A partir de aquí, todas las rutas requieren una sesión activa.
		r.Get("/", handlers.Repo.Home)

		// INVENTARIO
		r.Get("/inventario", handlers.Repo.Inventario)
		r.Get("/crear-producto", handlers.Repo.MostrarFormularioCrear)
		r.Post("/crear-producto", handlers.Repo.CrearProducto)
		r.Get("/editar-producto/{id}", handlers.Repo.MostrarFormularioEditar)
		r.Post("/editar-producto/{id}", handlers.Repo.EditarProducto)
		r.Post("/eliminar/{id}", handlers.Repo.EliminarProducto)

		// PROYECTOS
		r.Get("/proyectos-vista", handlers.Repo.ProyectosVista)
		r.Get("/nuevo-proyecto", handlers.Repo.MostrarNuevoProyecto)
		r.Post("/nuevo-proyecto", handlers.Repo.NuevoProyecto)
		r.Get("/proyectos/editar/{id}", handlers.Repo.MostrarFormularioEditarProyecto)
		r.Post("/proyectos/editar/{id}", handlers.Repo.ProcesarFormularioEditarProyecto)

		// LICITACIONES
		r.Get("/licitaciones", handlers.Repo.Licitaciones)
		r.Get("/nueva-licitacion", handlers.Repo.MostrarNuevaLicitacion)
		r.Post("/nueva-licitacion", handlers.Repo.CrearNuevaLicitacion)
		r.Get("/editar-licitacion/{id}", handlers.Repo.MostrarFormularioEditarLicitacion)
		r.Post("/editar-licitacion/{id}", handlers.Repo.EditarLicitacion)

		// CALENDARIO
		r.Get("/calendario", handlers.Repo.Calendario)

		// LICITACIONES - ARCHIVOS
		mux.Get("/archivos-licitacion/{id}", handlers.Repo.GetArchivosLicitacion)
		mux.Post("/guardar-enlace-licitacion", handlers.Repo.PostGuardarEnlace)
		mux.Post("/eliminar-archivo-licitacion", handlers.Repo.PostEliminarEnlace)

		// PARTIDAS
		r.Get("/mostrar-partidas/{id}", handlers.Repo.MostrarPartidasPorID)
		r.Get("/nueva-partida/{id}", handlers.Repo.MostrarNuevaPartida)
		r.Post("/nueva-partida", handlers.Repo.CrearNuevaPartida)
		r.Get("/editar-partida/{id}", handlers.Repo.MostrarEditarPartida)
		r.Post("/editar-partida/{id}", handlers.Repo.EditarPartida)
		r.Post("/eliminar-partida/{id}", handlers.Repo.PostEliminarPartida)

		

		// ACLARACIONES-LICITACION
		r.Get("/aclaraciones-licitacion/{id}", handlers.Repo.MostrarAclaracionesLicitacion)
		r.Get("/nueva-aclaracion-general/{id}", handlers.Repo.MostrarNuevaAclaracionGeneral)
		r.Post("/datos-empresas-externas-nueva-contexto-aclaraciones", handlers.Repo.AgregarEmpresaExternaContextoAclaraciones)
		r.Post("/nueva-aclaracion-licitacion", handlers.Repo.CrearNuevaAclaracionGeneral)
		r.Get("/editar-aclaracion/{id}", handlers.Repo.MostrarFormularioEditarAclaracion)
		r.Post("/editar-aclaracion/{id}", handlers.Repo.ProcesarFormularioEditarAclaracion)
		
		


		// PRODUCTOS PARTIDA
		r.Get("/productos-partida/{id}", handlers.Repo.MostrarProductosPartida)
		r.Get("/nuevo-producto-partida/{id}", handlers.Repo.MostrarNuevoProductoPartida)
		r.Post("/nuevo-producto-partida", handlers.Repo.CrearNuevoProductoPartida)
		r.Post("/editar-producto-partida", handlers.Repo.EditarProductoPartida)
		r.Post("/eliminar-producto-partida/{id}", handlers.Repo.EliminarProductoPartida)

		// ACLARACIONES
		r.Get("/aclaraciones/{id}", handlers.Repo.MostrarAclaraciones)
		r.Get("/nueva-aclaracion/{id}", handlers.Repo.MostrarNuevaAclaracion)
		r.Post("/datos-empresas-externas-nueva-contexto", handlers.Repo.AgregarEmpresaExternaContexto)
		r.Post("/nueva-aclaracion", handlers.Repo.CrearNuevaAclaracion)

		// REQUERIMIENTOS
		r.Get("/requerimientos/{id}", handlers.Repo.ObtenerRequerimientos)
		r.Get("/requerimientos-json/{id}", handlers.Repo.ObtenerRequerimientosJSON)
		r.Post("/guardar-requerimientos", handlers.Repo.GuardarRequerimientos)

		// PROPUESTAS
		r.Get("/propuestas/{id}", handlers.Repo.MostrarPropuestas)
		r.Get("/nueva-propuesta/{id}", handlers.Repo.MostrarNuevaPropuesta)
		r.Get("/api/productos-externos/buscar", handlers.Repo.BuscarProductosExternosJSON)
		r.Post("/api/marcas", handlers.Repo.CrearMarcaJSON)
		r.Post("/api/empresas-externas", handlers.Repo.CrearEmpresaExternaJSON)
		r.Post("/nueva-propuesta/{id}", handlers.Repo.CrearNuevaPropuesta)
		r.Post("/nuevo-producto-externo-contexto", handlers.Repo.NuevoProductoExternoContexto)
		r.Get("/editar-propuesta/{id}", handlers.Repo.MostrarEditarPropuesta)
		r.Post("/editar-propuesta/{id}", handlers.Repo.EditarPropuesta)


		// FALLOS
		r.Get("/fallo/{id}", handlers.Repo.ObtenerFalloPropuesta)
		r.Get("/fallo-json/{id}", handlers.Repo.ObtenerFalloPropuestaJSON)
		r.Post("/guardar-fallo", handlers.Repo.GuardarFalloPropuesta)

		// CATALOGO
		r.Get("/catalogo", handlers.Repo.Catalogo)
		r.Get("/producto/{id}", handlers.Repo.ProductoDetalles) // Para la vista detallada

		// OPCIONES
		r.Get("/opciones", handlers.Repo.Opciones)
			// OPCIONES DE DATOS REFERENCIA
			r.Get("/datos-referencia", handlers.Repo.DatosReferencia)
			r.Post("/datos-referencia", handlers.Repo.AgregarDato)
			r.Post("/eliminar-referencia/{id}", handlers.Repo.EliminarDatoReferencia)
			// OPCIONES DE ENTIDADES
			r.Get("/datos-entidades", handlers.Repo.Entidades)
			r.Get("/crear-entidad", handlers.Repo.MostrarNuevaEntidad)
			r.Post("/crear-entidad", handlers.Repo.CrearEntidad)
			// OPCIONES DE EMPRESAS EXTERNAS Y PRODUCTOS COMPETENCIA
			r.Get("/datos-empresas-externas", handlers.Repo.EmpresasExternas)
			r.Post("/datos-empresas-externas-nueva", handlers.Repo.AgregarEmpresaExterna)
			r.Get("/productos-externos", handlers.Repo.ProductosExternos)
			r.Post("/nuevo-producto-externo-contexto-menu", handlers.Repo.NuevoProductoExternoContextoMenu)
			// ADMINISTRACIÓN DE USUARIOS
			r.Get("/usuarios", handlers.Repo.MostrarUsuarios)
			r.Post("/eliminar-usuario/{id}", handlers.Repo.EliminarUsuario)
			r.Post("/crear-usuario", handlers.Repo.CrearUsuario)
			r.Post("/editar-usuario/{id}", handlers.Repo.EditarUsuario)      


		// Ruta para cerrar sesión
		r.Get("/logout", handlers.Repo.Logout)
	})

	return mux
}