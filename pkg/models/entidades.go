package models

import "time"

type Entidad struct {
    IDEntidad      int       `json:"id_entidad"`
    Nombre         string    `json:"nombre"`
    Tipo           string    `json:"tipo"`
	Compañia 	   Compañias `json:"compañias"`

	Estado         EstadosRepublica `json:"estados"`
	Municipio	   string 	 `json:"municipio"`
	CodigoPostal   string 	 `json:"codigo_postal"`
	Direccion      string 	 `json:"direccion"`

    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type EstadosRepublica struct {
	ClaveEstado   string	`json:"estado_clave"`
	NombreEstado  string	`json:"nombre_estado"`
	CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Compañias struct {
	IDCompañia int 		`json:"id_compañia"`
	Nombre     string	`json:"nombre_compañia"`
	Tipo       string	`json:"tipo"`
	CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}