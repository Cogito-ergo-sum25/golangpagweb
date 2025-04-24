package models

import "time"

type Licitacion struct {
    ID             int       `json:"id_licitacion"`
    IDEntidad      int       `json:"id_entidad"`
    Nombre         string    `json:"nombre"`
    NumContratacion string   `json:"num_contratacion"`
    Estatus        string    `json:"estatus"`
    FechaPropuestas time.Time `json:"fecha_propuestas"`
    
    // Datos entidad para mostrar
    EntidadNombre  string    `json:"entidad_nombre"`
    
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type Entidad struct {
    ID             int       `json:"id_entidad"`
    Nombre         string    `json:"nombre"`
    Tipo           string    `json:"tipo"`
    CreatedAt      time.Time `json:"created_at"`
}