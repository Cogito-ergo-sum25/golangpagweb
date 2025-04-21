package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
)

// Función helper para verificar certificaciones seleccionadas
func isCertSelected(certID int, productCerts []models.Certificacion) bool {
	for _, pc := range productCerts {
		if pc.IDCertificacion == certID {
			return true
		}
	}
	return false
}

// Mapa de funciones para los templates
var functions = template.FuncMap{
	"isCertSelected": isCertSelected, // Registramos la función
}

var app *config.AppConfig

// NewTemplates configura el paquete render
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData añade datos comunes a todos los templates
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renderiza un template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Println("Error creando cache de templates:", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println("No se pudo encontrar el template:", tmpl)
		http.Error(w, "Template no encontrado", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err = t.Execute(buf, td)
	if err != nil {
		log.Println("Error ejecutando template:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error escribiendo template al navegador:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// CreateTemplateCache crea un cache de templates
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		
		// Parseamos el template con las funciones registradas
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}