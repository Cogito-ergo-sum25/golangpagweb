package models

import "time"

type Partida struct {
	IDPartida              int       `json:"id_partida"`
	IDLicitacion           int       `json:"id_licitacion"`
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
	Licitaciones []LicitacionPartida `json:"licitaciones"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LicitacionPartida representa la relación entre licitaciones y partidas
type LicitacionPartida struct {
	IDLicitacionPartida int `json:"id_licitacion_partida"`
	IDLicitacion        int `json:"id_licitacion"`
	IDPartida           int `json:"id_partida"`

	// Datos para mostrar (pueden ser omitidos en JSON si están vacíos)
	Partida    *Partida    `json:"partida"`
	Licitacion *Licitacion `json:"licitacion"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PartidaProductos representa la relación mucho a muchos entre partidas y productos
// con el precio ofertado y observaciones
type PartidaProductos struct {
	IDPartidaProducto int     `json:"id_partida_producto"`
	IDProducto        int     `json:"id_producto"`
	IDPartida         int     `json:"id_partida"`
	PrecioOfertado    float64 `json:"precio_ofertado"`
	Observaciones     string  `json:"observaciones"`

	Partida  *Partida  `json:"partida"`
	Producto *Producto `json:"producto"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Empresas struct {
	IDEmpresa int       `json:"id_empresa"`
	Nombre    string    `json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequerimientosPartida struct {
	IDRequerimientos       int  `json:"id_requerimientos"`
	IDLicitacion           int  `json:"id_licitacion"`
	RequiereMantenimiento  bool `json:"requiere_mantenimiento"`
	RequiereInstalacion    bool `json:"requiere_instalacion"`
	RequierePuestaEnMarcha bool `json:"requiere_puesta_marcha"`
	RequiereCapacitacion   bool `json:"requiere_capacitacion"`
	RequiereVisitaPrevia   bool `json:"requiere_visita_previa"`
	FechaVisita    time.Time `json:"fecha_visita"`
	ComentariosVisita	string `json:"comentarios_visita"`
	RequiereMuestra bool `json:"requiere_muestra_producto"`
	FechaMuestra time.Time `json:"fecha_muestra"`
	ComentariosMuestra string `json:"comentarios_muestra"`
	FechaEntrega time.Time `json:"fecha_entrega"`
	ComentariosEntrega string `json:"comentarios_entrega"`

	Partida *Partida `json:"partida"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AclaracionesPartida struct {
	IDAclaracion    int    `json:"id_aclaracion"`
	Pregunta        string `json:"pregunta"`
	Observaciones   string `json:"observaciones"`
	FichaTecnica    int    `json:"ficha_tecnica_id"`
	IDPuntosTecnico int    `json:"id_puntos_tecnicos_modif"`

	Partida *Partida  `json:"id_partida"`
	Empresa *Empresas `json:"id_empresa"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductosExternos representa los productos externos asociados a una propuesta
// que pueden ser parte de una partida en una licitación.
type ProductosExternos struct {
	IDProducto       int     `json:"id_producto"`
	IDMarca          int     `json:"id_marca"`
	IDPaisOrigen     int     `json:"id_pais_origen"`
	IDEmpresaExterna int     `json:"id_empresa_externa"`
	Nombre           string  `json:"nombre"`
	Modelo           string  `json:"modelo"`
	PrecioOfertado   float64 `json:"precio_ofertado"`
	Observaciones    string  `json:"observaciones"`

	Marca          *Marca    `json:"marca"`
	PaisOrigen     *Pais     `json:"pais_origen"`
	EmpresaExterna *Empresas `json:"empresa_externa"`

	Partida  *Partida  `json:"partida"`
	Producto *Producto `json:"producto"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PropuestasPartida representa las propuestas realizadas por empresas para una partida específica
// en una licitación. Incluye el precio ofertado, precios mínimo y máximo, y observaciones.
type PropuestasPartida struct {
	IDPropuesta       int     `json:"id_propuesta"`
	IDPartida         int     `json:"id_partida"`
	IDEmpresa         int     `json:"id_empresa"`
	IDProductoExterno int     `json:"id_producto_externo"`
	PrecioOfertado    float64 `json:"precio_ofertado"`
	PrecioMin         float64 `json:"precio_min"`
	PrecioMax         float64 `json:"precio_max"`
	Observaciones     string  `json:"observaciones"`

	ProductoExterno *ProductosExternos `json:"producto_externo"`
	Partida         *Partida           `json:"partida"`
	Empresa         *Empresas          `json:"empresa"`
	Fallo *FallosPropuesta `json:"fallo"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FallosPropuesta representa el resultado de la evaluación de una propuesta
type FallosPropuesta struct {
	IDFallo              int       `json:"id_fallo"`
	IDPropuesta          int       `json:"id_propuesta"`
	CumpleLegal          bool      `json:"cumple_legal"`
	CumpleAdministrativo bool      `json:"cumple_administrativo"`
	CumpleTecnico        bool      `json:"cumple_tecnico"`
	PuntosObtenidos      int       `json:"puntos_obtenidos"`
	Ganador              bool      `json:"ganador"`
	Observaciones        string    `json:"observaciones"`

	Propuesta            *PropuestasPartida `json:"propuesta"`

	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
};
