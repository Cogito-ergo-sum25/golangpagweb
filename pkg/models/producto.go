package models

import "time"

type Producto struct {
    // IDs y relaciones
    IDProducto        int       `json:"id_producto"`
    IDMarca          int       `json:"id_marca"`
    IDTipo           int       `json:"id_tipo"`
    IDClasificacion  int       `json:"id_clasificacion"`
    IDPaisOrigen     int       `json:"id_pais_origen"`
    
    // Información básica
    SKU              string    `json:"sku"`
    Nombre           string    `json:"nombre"`
    NombreCorto      string    `json:"nombre_corto"`
    Modelo           string    `json:"modelo"`
    Version          string    `json:"version"`
    Serie            string    `json:"serie"`
    CodigoFabricante string    `json:"codigo_fabricante"`
    Descripcion      string    `json:"descripcion"`
    
    // URLs e imágenes
    ImagenURL        string    `json:"imagen_url"`
    FichaTecnicaURL  string    `json:"ficha_tecnica_url"`
    
    // Campos de control
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    
    // Campos para mostrar (no se guardan en BD)
    Marca            string    `json:"marca"`             // Nombre de la marca (JOIN)
    Tipo             string    `json:"tipo"`              // Nombre del tipo (JOIN)
    Clasificacion    string    `json:"clasificacion"`     // Nombre clasificación (JOIN)
    PaisOrigen       string    `json:"pais_origen"`       // Nombre país (JOIN)
    
    // Certificaciones (para mostrar)
    Certificaciones  []Certificacion `json:"certificaciones"`
}

// Modelos auxiliares para las relaciones
type Marca struct {
    IDMarca   int    `json:"id_marca"`
    Nombre    string `json:"nombre"`
}

type TipoProducto struct {
    IDTipo    int    `json:"id_tipo"`
    Nombre    string `json:"nombre"`
}

type Clasificacion struct {
    IDClasificacion int    `json:"id_clasificacion"`
    Nombre         string `json:"nombre"`
}

type Pais struct {
    IDPais int    `json:"id_pais"`
    Nombre string `json:"nombre"`
    Codigo string `json:"codigo"`
}

type Certificacion struct {
    IDCertificacion int    `json:"id_certificacion"`
    Nombre          string `json:"nombre"`
    OrganismoEmisor string `json:"organismo_emisor"`
}
