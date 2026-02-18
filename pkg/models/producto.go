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

    // Campos de CROL
    Inventario *ProductoInventario `json:"inventario,omitempty"`
    IEPS       *IEPS             `json:"ieps,omitempty"`
    ComercioExterior *ComercioExterior `json:"comercio_exterior,omitempty"`
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

type ProductoInventario struct {
    IDProducto                  int     `json:"id_producto"`
    
    // Unidades y Costeo
    UnidadBase                  string  `json:"unidad_base"`
    UnidadMedidaAlmacen         string  `json:"unidad_medida_almacen"`
    MetodoCosteo                string  `json:"metodo_costeo"`
    
    // Dimensiones Físicas
    Largo                       float64 `json:"largo"`
    Ancho                       float64 `json:"ancho"`
    Alto                        float64 `json:"alto"`
    Peso                        float64 `json:"peso"`
    Volumen                     float64 `json:"volumen"`
    
    // Control Logístico
    RequierePesaje              bool    `json:"requiere_pesaje"`
    ConsiderarCompraProgramada  bool    `json:"considerar_compra_programada"`
    ProduccionFabricacion       bool    `json:"produccion_fabricacion"`
    VentasSinExistencia         bool    `json:"ventas_sin_existencia"`
    ManejaSerie                 bool    `json:"maneja_serie"`
    ManejaLote                  bool    `json:"maneja_lote"`
    ManejaFechaCaducidad        bool    `json:"maneja_fecha_caducidad"`
    LoteAutomatico              bool    `json:"lote_automatico"`
}

type IEPS struct {
    IDProducto     int
    TipoProducto   string
    ClaveProducto  string
    Empaque        string
    UnidadMedida   string
    Presentacion   float64
}

type ComercioExterior struct {
    IDProducto          int     `json:"id_producto"`
    Modelo              string  `json:"modelo"`
    SubModelo           string  `json:"sub_modelo"`
    FraccionArancelaria string  `json:"fraccion_arancelaria"`
    UnidadMedidaAduana  string  `json:"unidad_medida_aduana"`
    FactorConversionUMT float64 `json:"factor_conversion_umt"`
}

type ProductoCatalogo struct {
    IDCatalogo        int       `json:"id_catalogo"`
    IDProducto        int       `json:"id_producto"`
    IDLicitacion      int       `json:"id_licitacion"`
    IDPartidaProducto int       `json:"id_partida_producto"` 
    NombreProducto    string    `json:"nombre_producto"`
    NombreVersion     string    `json:"nombre_version"`
    ArchivoURL        string    `json:"archivo_url"`
    Descripcion       string    `json:"descripcion"`
    UpdatedAt         time.Time `json:"updated_at"`
    CreatedAt         time.Time `json:"created_at"`
    ContextoLicitacion string    // Para mostrar el NumContratacion
    ContextoPartida    int       // Para mostrar el número de partida
}