package models

type TemplateData struct {
	// Mapeos genéricos
	StringMap  map[string]string
	IntMap     map[string]int
	FloatMap   map[string]float32
	Data       map[string]interface{} // Para datos dinámicos
	
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
	Estado   EstadosRepublica

	// --- Datos de compañias ---
    Compañias []Compañias
	Compañia   Compañias

	// --- Datos de partidas ---
    Partidas []Partida
    Partida   Partida

	// --- Datos de requerimientos ---
	Requerimientos RequerimientosPartida

	
	// Listas para formularios
	Marcas          []Marca
	TiposProducto   []TipoProducto
	Clasificaciones []Clasificacion
	Paises          []Pais
	Certificaciones []Certificacion
	
	// Para edición de productos
	CertificacionesProducto []Certificacion // Certificaciones asignadas al producto
	TodasCertificaciones   []Certificacion // Todas las certificaciones disponibles

	// Filtros comunes
    Filtro struct {
        Estatus string
        EntidadID int
        TipoLicitacion string
    }
}