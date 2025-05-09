package models

import "time"

type Licitacion struct {
    IDLicitacion    int       `json:"id_licitacion"`
    IDEntidad       int       `json:"id_entidad"`
    NumContratacion string   `json:"num_contratacion"`
    Caracter         string    `json:"caracter"`
    Nombre           string    `json:"nombre"`
    Estatus          string    `json:"estatus"`
    Tipo             string    `json:"tipo"`
    Lugar            string    `json:"lugar"`
    FechaJunta       time.Time `json:"fecha_junta"` 
    FechaPropuestas  time.Time `json:"fecha_propuestas"`
    FechaFallo       time.Time `json:"fecha_fallo"`
    FechaEntrega     time.Time `json:"fecha_entrega"`
    TiempoEntrega    string    `json:"tiempo_entrega"`
    Revisada         bool      `json:"revisada"`
    Intevi           bool      `json:"intevi"`
    Estado           string    `json:"estado"`
    ObservacionesGenerales string `json:"observaciones_generales"`

    // Datos entidad para mostrar
    EntidadNombre    string    `json:"entidad_nombre"`
    EntidadTipo      string    `json:"entidad_tipo"`
    EntidadMunicipio string
    EstadoNombre     string
    CompaniaTipo     string


    // Relación con proyecto (1:1)
    Proyecto         *Proyecto `json:"proyecto,omitempty"`

    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
