package models

type TemplateData struct {
	// Mapeos genéricos
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // Para datos dinámicos

	// Seguridad
	CSRFToken string

	// Mensajes al usuario
	Flash   string
	Warning string
	Error   string

	// --- Datos de productos ---
	Productos []Producto
	Producto  Producto

	// --- Datos de proyectos ---
	Proyectos []Proyecto
	Proyecto  Proyecto

	// --- Datos de licitaciones ---
	Licitaciones []Licitacion
	Licitacion   Licitacion

	// --- Datos de entidades ---
	Entidades []Entidad
	Entidad   Entidad

	// --- Datos de estados ---
	Estados []EstadosRepublica
	Estado  EstadosRepublica

	// --- Datos de compañias ---
	Compañias []Compañias
	Compañia  Compañias

	// --- Datos de partidas ---
	Partidas []Partida
	Partida  Partida

	// --- Datos de empresas externas ---
	Empresas []Empresas
	Empresa  Empresas

	// --- Datos de archivos de licitación ---
	Archivos []ArchivoLicitacion
	Archivo  ArchivoLicitacion

	// --- Datos de catálogos de productos ---
    Catalogos []ProductoCatalogo // Para listar todas las versiones
    Catalogo  ProductoCatalogo   // Para edición o visualización individual

	// --- Datos de requerimientos ---
	Requerimientos RequerimientosPartida

	// --- Datos de aclaraciones ---
	Aclaraciones []AclaracionesPartida

	// --- Datos de aclaraciones por licitacion ---
	AclaracionesLicitacion []AclaracionesLicitacion

	// --- Productos asociados a una partida ---
	ProductosPartida []PartidaProductos

	// Productos externos asociados a una propuesta
	ProductosExternos []ProductosExternos 

	// Propuestas asociadas a una partida
	Propuesta PropuestasPartida // Propuesta individual
	PropuestasPartida []PropuestasPartida 

	// Fallos de propuestas
	Fallo *FallosPropuesta
	Fallos []FallosPropuesta

	PartidasProducto []PartidaProductos

	// Listas para formularios
	Marcas          []Marca
	TiposProducto   []TipoProducto
	Clasificaciones []Clasificacion
	Paises          []Pais
	Certificaciones []Certificacion

	// Para edición de productos
	CertificacionesProducto []Certificacion // Certificaciones asignadas al producto
	TodasCertificaciones    []Certificacion // Todas las certificaciones disponibles

	// Filtros comunes
	Filtro struct {
		Estatus        string
		EntidadID      int
		TipoLicitacion string
	}
}
