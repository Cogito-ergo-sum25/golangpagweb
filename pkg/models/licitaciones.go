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
    CriterioEvaluacion     string `json:"criterio_evaluacion"`
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

type AclaracionesLicitacion struct {
    IDAclaracionLicitacion int       `json:"id_aclaracion_licitacion"`
    IDLicitacion           int       `json:"id_licitacion"`
    IDPartida              int       `json:"id_partida,omitempty"` // Puede ser NULL si es una aclaración general
    IDEmpresa              int       `json:"id_empresa"`
    Pregunta               string    `json:"pregunta"`
    Observaciones          string    `json:"observaciones"`
    FichaTecnicaID        string           `json:"ficha_tecnica_id,omitempty"` // Puede ser NULL si no aplica
    IDPuntosTecnicosModif int       `json:"id_puntos_tecnicos_modif,omitempty"` // Puede ser NULL si no aplica
    PreguntaTecnica       bool      `json:"pregunta_tecnica"` // TRUE si es una pregunta técnica, FALSE si es general
    CreatedAt             time.Time `json:"created_at"`
    UpdatedAt             time.Time `json:"updated_at"`

    // Relaciones con otras entidades
    Licitacion *Licitacion `json:"licitacion,omitempty"` // Relación con la licitación
    Partida    *Partida    `json:"partida,omitempty"`    // Relación
    Empresa    *Empresas   `json:"empresa,omitempty"`    // Relación con la empresa que hace la pregunta
}
