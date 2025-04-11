package handlers

import (
	"net/http"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// Inventario is the handler for the inventario page
func (m *Repository) Inventario(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "inventario.page.tmpl", &models.TemplateData{})
}

// Catalogo is the handler for the catalogo page
func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "catalogo.page.tmpl", &models.TemplateData{})
}

// Crear is the handler for the crear page
func (m *Repository) Crear(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "crear.page.tmpl", &models.TemplateData{})
}
