package models

type Producto struct {
    IDProducto          int     `json:"id_producto"`
    Marca              string  `json:"marca"`
    Tipo               string  `json:"tipo"`
    SKU                string  `json:"sku"`
    Nombre             string  `json:"nombre"`
    Descripcion        string  `json:"descripcion"`
    Cantidad           int     `json:"cantidad"`
    ImagenURL          string  `json:"imagen_url"`
    FichaTecnicaURL    string  `json:"ficha_tecnica_url"`
    Modelo             string  `json:"modelo"`
    CodigoFabricante   string  `json:"codigo_fabricante"`
    PrecioLista        float64 `json:"precio_lista"`
    PrecioMinimo       float64 `json:"precio_minimo"`
    Clasificacion      string  `json:"clasificacion"`
    Serie              string  `json:"serie"`
    PaisOrigen         string  `json:"pais_origen"`
    Certificaciones    string  `json:"certificaciones"`
    RequiereInstalacion bool   `json:"requiere_instalacion"`
    TiempoEntrega      int     `json:"tiempo_entrega"`
    StockMinimo        int     `json:"stock_minimo"`
    EnPromocion        bool    `json:"en_promocion"`
    ClaveProductoSAT   string  `json:"clave_producto_sat"`
    UnidadMedidaSAT    string  `json:"unidad_medida_sat"`
}

