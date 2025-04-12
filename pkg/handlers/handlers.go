package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/config"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/render"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
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
	rows, err := m.App.DB.Query(`SELECT id_producto, marca, tipo, sku, nombre, descripcion, cantidad FROM productos`)
	if err != nil {
        http.Error(w, "Error al consultar productos: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
	var productos []models.Producto
    for rows.Next() {
        var p models.Producto
        err := rows.Scan(
            &p.IDProducto,
            &p.Marca,
            &p.Tipo,
            &p.SKU,
            &p.Nombre,
            &p.Descripcion,
            &p.Cantidad,
        )
        if err != nil {
            http.Error(w, "Error al leer producto: "+err.Error(), http.StatusInternalServerError)
            return
        }
        productos = append(productos, p)
    }

    // Verificar errores después de iterar
    if err = rows.Err(); err != nil {
        http.Error(w, "Error después de leer productos: "+err.Error(), http.StatusInternalServerError)
        return
    }

    data := &models.TemplateData{
        Productos: productos,
        CSRFToken: nosurf.Token(r), // ¡Añade esto!
    }
	render.RenderTemplate(w, "inventario.page.tmpl", data)
}

// Catalogo is the handler for the catalogo page
func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "catalogo.page.tmpl", &models.TemplateData{})
}

// Muestra el formulario de creación
func (m *Repository) MostrarFormularioCrear(w http.ResponseWriter, r *http.Request) {
    data := &models.TemplateData{
        CSRFToken: nosurf.Token(r), // Añade el token CSRF
    }
    render.RenderTemplate(w, "crear.page.tmpl", data)
}

func (m *Repository) CrearProducto(w http.ResponseWriter, r *http.Request) {
    // 1. Verificar método HTTP
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    // 2. Parsear el formulario
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario: "+err.Error(), http.StatusBadRequest)
        return
    }

    // 3. Validar campos requeridos
    requiredFields := []string{"marca", "nombre", "cantidad"}
    for _, field := range requiredFields {
        if r.Form.Get(field) == "" {
            http.Error(w, "El campo "+field+" es requerido", http.StatusBadRequest)
            return
        }
    }

    // 4. Convertir cantidad
    cantidad, err := strconv.Atoi(r.Form.Get("cantidad"))
    if err != nil {
        http.Error(w, "La cantidad debe ser un número válido", http.StatusBadRequest)
        return
    }

    // 5. Insertar en la base de datos
    _, err = m.App.DB.Exec(`INSERT INTO productos (marca, tipo, sku, nombre, descripcion, cantidad) VALUES (?, ?, ?, ?, ?, ?)`,
        r.Form.Get("marca"),
        r.Form.Get("tipo"),
        r.Form.Get("sku"),
        r.Form.Get("nombre"),
        r.Form.Get("descripcion"),
        cantidad,
    )

    if err != nil {
        log.Println("Error al insertar producto:", err)
        http.Error(w, "Error interno al guardar el producto", http.StatusInternalServerError)
        return
    }

    // 6. Redirigir con mensaje de éxito
    m.App.Session.Put(r.Context(), "flash", "Producto creado exitosamente")
    http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

func (m *Repository) MostrarFormularioEditar(w http.ResponseWriter, r *http.Request) {
    // Extraer el ID del producto de la URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener el producto de la base de datos
    var producto models.Producto
    err = m.App.DB.QueryRow(`
        SELECT id_producto, marca, tipo, sku, nombre, descripcion, cantidad 
        FROM productos 
        WHERE id_producto = ?`, id).Scan(
            &producto.IDProducto,
            &producto.Marca,
            &producto.Tipo,
            &producto.SKU,
            &producto.Nombre,
            &producto.Descripcion,
            &producto.Cantidad,
    )
    if err != nil {
        http.Error(w, "Producto no encontrado", http.StatusNotFound)
        return
    }

    // Renderizar el formulario de edición con los datos del producto
    data := &models.TemplateData{
        Producto: producto,  // Pasas un solo producto
        CSRFToken: nosurf.Token(r), // Asegúrate de incluir el token CSRF
    }
    
    render.RenderTemplate(w, "editar.page.tmpl", data)
}

func (m *Repository) EditarProducto(w http.ResponseWriter, r *http.Request) {
    // Extraer el ID del producto
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Parsear el formulario
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    // Validar campos requeridos
    if r.Form.Get("marca") == "" || r.Form.Get("nombre") == "" || r.Form.Get("cantidad") == "" {
        http.Error(w, "Campos requeridos faltantes", http.StatusBadRequest)
        return
    }

    cantidad, err := strconv.Atoi(r.Form.Get("cantidad"))
    if err != nil {
        http.Error(w, "Cantidad debe ser un número", http.StatusBadRequest)
        return
    }

    // Actualizar en la base de datos
    _, err = m.App.DB.Exec(`
        UPDATE productos 
        SET marca = ?, tipo = ?, sku = ?, nombre = ?, descripcion = ?, cantidad = ?, updated_at = NOW()
        WHERE id_producto = ?`,
        r.Form.Get("marca"),
        r.Form.Get("tipo"),
        r.Form.Get("sku"),
        r.Form.Get("nombre"),
        r.Form.Get("descripcion"),
        cantidad,
        id,
    )
    if err != nil {
        log.Println("Error al actualizar producto:", err)
        http.Error(w, "Error al guardar cambios", http.StatusInternalServerError)
        return
    }

    // Redirigir al listado con mensaje de éxito
    m.App.Session.Put(r.Context(), "flash", "Producto actualizado correctamente")
    http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

func (m *Repository) EliminarProducto(w http.ResponseWriter, r *http.Request) {
    log.Println("Token CSRF recibido:", r.FormValue("csrf_token")) // Verifica si llega
    log.Println("Token CSRF esperado:", nosurf.Token(r))          // Compara con el válido
    // Verifica el método HTTP
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    // Valida el token CSRF (si usas middleware, esto ya se hace automáticamente)
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Ejecuta el DELETE en la DB
    _, err = m.App.DB.Exec("DELETE FROM productos WHERE id_producto = ?", id)
    if err != nil {
        log.Println("Error al borrar:", err)
        http.Error(w, "Error interno", http.StatusInternalServerError)
        return
    }

    // Redirige con mensaje de éxito
    m.App.Session.Put(r.Context(), "flash", "Producto eliminado")
    http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}