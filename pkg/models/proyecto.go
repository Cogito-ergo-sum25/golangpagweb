package models

import "time"

type Proyecto struct {
    IDProyecto    int       `json:"id_proyecto"`
    Nombre        string    `json:"nombre"`
    Descripcion   string    `json:"descripcion"`
    FechaInicio   time.Time `json:"fecha_inicio"`
	FechaFin   time.Time 	`json:"fecha_fin"`
    
    // Relaciones para mostrar
    IDLicitacion  int       `json:"id_licitacion"`
    LicitacionNombre string `json:"licitacion_nombre"`
    NumContratacion  string `json:"num_contratacion"`
    FechaJunta       time.Time
    FechaPropuestas  time.Time
    FechaFallo       time.Time
    FechaEntrega     time.Time
    Lugar            string    `json:"lugar"`
    EstadoLicitacion string
    EntidadNombre    string `json:"entidad_nombre"`
    
    // Productos asociados
    Productos       []ProductoProyecto `json:"productos"`
    
    CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProductoProyecto representa la relación muchos-a-muchos entre productos y proyectos
type ProductoProyecto struct {
    IDProducto      int       `json:"id_producto"`
    IDProyecto      int       `json:"id_proyecto"`
    Cantidad        int       `json:"cantidad"`
    PrecioUnitario  float64   `json:"precio_unitario"`
    Especificaciones string   `json:"especificaciones"`
    
    // Campos para mostrar
    ProductoNombre string    `json:"producto_nombre"`
    SKU            string    `json:"sku"`
    ImagenURL      string    `json:"imagen_url"`
	Modelo         string    `json:"modelo"`

}