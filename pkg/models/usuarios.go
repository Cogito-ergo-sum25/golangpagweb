package models

type Usuario struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Email       string `json:"email"`
	NivelAcceso int    `json:"nivel_acceso"`
}