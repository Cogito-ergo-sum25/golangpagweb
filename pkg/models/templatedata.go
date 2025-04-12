package models

type TemplateData struct {
	StringMap 	map[string]string
	IntMap 		map[string]int
	FloatMap 	map[string]float32
	Data 		map[string]interface{}
	CSRFToken	string
	Flash		string
	Warning		string
	Error		string
	Productos []Producto
	Producto   Producto   // Para editar un solo producto (ej. en /editar/{id})
}