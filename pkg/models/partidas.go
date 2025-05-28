package models

import "time"

type Partida struct {
	IDPartida              int       `json:"id_partida"`
    IDLicitacion           int       `json:"id_licitacion"` // <-- Agregado aquí
    NumPartidaConvocatoria int       `json:"numero_partida_convocatoria"`
    NombreDescripcion      string    `json:"nombre_descripcion"`
    Cantidad               int       `json:"cantidad"`
	CantidadMinima         int       `json:"cantidad_minima"`
	CantidadMaxima         int       `json:"cantidad_maxima"`
	NoFichaTecnica         string    `json:"no_ficha_tecnica"`
	TipoDeBien             string    `json:"tipo_de_bien"`
    ClaveCompendio         string    `json:"clave_compendio"`
    ClaveCucop             string    `json:"clave_cucop"`
    UnidadMedida           string    `json:"unidad_medida"`
    DiasDeEntrega          string    `json:"días_de_entrega"`
    FechaDeEntrega         time.Time `json:"fecha_de_entrega"`
    Garantia               int       `json:"garantia"`

	// Relaciones
    Licitaciones          []LicitacionPartida `json:"licitaciones"`

	CreatedAt             time.Time  `json:"created_at"`
    UpdatedAt             time.Time  `json:"updated_at"`
}

// LicitacionPartida representa la relación entre licitaciones y partidas
type LicitacionPartida struct {
    IDLicitacionPartida     int       `json:"id_licitacion_partida"`
    IDLicitacion            int       `json:"id_licitacion"`
    IDPartida               int       `json:"id_partida"`
    
    // Datos para mostrar (pueden ser omitidos en JSON si están vacíos)
    Partida                *Partida   `json:"partida"`
    Licitacion             *Licitacion `json:"licitacion"`
    
    CreatedAt              time.Time  `json:"created_at"`
    UpdatedAt              time.Time  `json:"updated_at"`
}

type RequerimientosPartida struct {
    IDRequerimientos        int       `json:"id_requerimientos"`
    IDLicitacion            int       `json:"id_licitacion"`
    RequiereMantenimiento   bool      `json:"requiere_mantenimiento"`
    RequiereInstalacion     bool      `json:"requiere_instalacion"`
    RequierePuestaEnMarcha  bool      `json:"requiere_puesta_marcha"`
    RequiereCapacitacion    bool      `json:"requiere_capacitacion"`
    RequiereVisitaPrevia    bool      `json:"requiere_visita_previa"`
    
    Partida                 *Partida  `json:"partida"`

    CreatedAt               time.Time  `json:"created_at"`
    UpdatedAt               time.Time  `json:"updated_at"`
}