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
	
	// Datos de productos
	Productos []Producto
	Producto  Producto
	
	// Listas para formularios
	Marcas          []Marca
	TiposProducto   []TipoProducto
	Clasificaciones []Clasificacion
	Paises          []Pais
	Certificaciones []Certificacion
	
	// Para edición de productos
	CertificacionesProducto []Certificacion // Certificaciones asignadas al producto
	TodasCertificaciones   []Certificacion // Todas las certificaciones disponibles
}