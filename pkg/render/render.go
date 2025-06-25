package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

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

func formatPrice(price float64) string {
    s := fmt.Sprintf("%.2f", price)
    parts := strings.Split(s, ".")
    intPart := parts[0]
    decPart := parts[1]

    // Insertar comas
    n := len(intPart)
    if n <= 3 {
        return intPart + "." + decPart
    }

    var result strings.Builder
    pre := n % 3
    if pre > 0 {
        result.WriteString(intPart[:pre])
        if n > pre {
            result.WriteString(",")
        }
    }

    for i := pre; i < n; i += 3 {
        result.WriteString(intPart[i : i+3])
        if i+3 < n {
            result.WriteString(",")
        }
    }

    return result.String() + "." + decPart
}

// Función que convierte cualquier valor a JSON
func toJson(v interface{}) template.JS {
    a, _ := json.Marshal(v)
    return template.JS(a)
}


// Mapa de funciones para los templates
var functions = template.FuncMap{
	"isCertSelected": isCertSelected,
	"formatPrice": formatPrice,
	"toJson": toJson,
	"safeJS": func(s string) template.JS {
		return template.JS(s)
	},
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

	// Normalizar la clave para buscarla en el mapa
	tmpl = filepath.ToSlash(tmpl) // ← esto estandariza "home\home.page.tmpl" a "home/home.page.tmpl"

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

	// Buscar todos los archivos .page.tmpl en subdirectorios
	pages, err := filepath.Glob("./templates/**/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Crear template con funciones
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Agregar layouts si existen
		layouts, err := filepath.Glob("./templates/base/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/base/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		// Clave relativa: sin "templates/", usando `/`
		relPath, err := filepath.Rel("templates", page)
		if err != nil {
			return myCache, err
		}
		relPath = filepath.ToSlash(relPath) // estandariza a "home/home.page.tmpl"

		myCache[relPath] = ts
	}

	return myCache, nil
}
