package models

type Producto struct {
    IDProducto  int    `json:"id_producto"`
    Marca       string `json:"marca"`
    Tipo        string `json:"tipo"`
    SKU         string `json:"sku"`
    Nombre      string `json:"nombre"`
    Descripcion string `json:"descripcion"`
    Cantidad    int    `json:"cantidad"`
}