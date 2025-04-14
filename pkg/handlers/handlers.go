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

// Catalogo is the handler for the catalogo page
func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "catalogo.page.tmpl", &models.TemplateData{})
}


// TODO LO DE INVENTARIO

// PANTALLA INICIO INVENTARIO
func (m *Repository) Inventario(w http.ResponseWriter, r *http.Request) {
    rows, err := m.App.DB.Query(`
        SELECT 
            id_producto, 
            COALESCE(clasificacion, '') as clasificacion,
            COALESCE(marca, '') as marca,
            COALESCE(nombre, '') as nombre,
            COALESCE(sku, '') as sku, 
            COALESCE(cantidad, 0) as cantidad,
            COALESCE(precio_lista, 0) as precio_lista
        FROM productos
        ORDER BY id_producto DESC`)
    
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
            &p.Clasificacion,
            &p.Marca,
            &p.Nombre,
            &p.SKU,
            &p.Cantidad,
            &p.PrecioLista,
        )
        if err != nil {
            http.Error(w, "Error al leer producto: "+err.Error(), http.StatusInternalServerError)
            return
        }
        productos = append(productos, p)
    }

    if err = rows.Err(); err != nil {
        http.Error(w, "Error después de leer productos: "+err.Error(), http.StatusInternalServerError)
        return
    }

    data := &models.TemplateData{
        Productos: productos,
        CSRFToken: nosurf.Token(r),
    }
    
    render.RenderTemplate(w, "inventario.page.tmpl", data)
}

// handler para todo el crear
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
    requiredFields := []string{"marca", "nombre", "sku", "cantidad", "precio_lista", "clasificacion"}
    for _, field := range requiredFields {
        if r.Form.Get(field) == "" {
            http.Error(w, "El campo "+field+" es requerido", http.StatusBadRequest)
            return
        }
    }

    // 4. Convertir tipos numéricos y booleanos
    cantidad, err := strconv.Atoi(r.Form.Get("cantidad"))
    if err != nil {
        http.Error(w, "La cantidad debe ser un número válido", http.StatusBadRequest)
        return
    }

    precioLista, err := strconv.ParseFloat(r.Form.Get("precio_lista"), 64)
    if err != nil {
        http.Error(w, "El precio debe ser un número válido", http.StatusBadRequest)
        return
    }

    precioMinimo, _ := strconv.ParseFloat(r.Form.Get("precio_minimo"), 64) // Opcional
    stockMinimo, _ := strconv.Atoi(r.Form.Get("stock_minimo")) // Opcional
    tiempoEntrega, _ := strconv.Atoi(r.Form.Get("tiempo_entrega")) // Opcional

    // Convertir checkboxes a booleanos
    requiereInstalacion := r.Form.Get("requiere_instalacion") == "on"
    enPromocion := r.Form.Get("en_promocion") == "on"

    // 5. Insertar en la base de datos
    _, err = m.App.DB.Exec(`
        INSERT INTO productos (
            marca, tipo, sku, nombre, descripcion, cantidad,
            imagen_url, ficha_tecnica_url, modelo, codigo_fabricante,
            precio_lista, precio_minimo, clasificacion, serie,
            pais_origen, certificaciones, requiere_instalacion,
            tiempo_entrega, stock_minimo, en_promocion,
            clave_producto_sat, unidad_medida_sat
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        r.Form.Get("marca"),
        r.Form.Get("tipo"),
        r.Form.Get("sku"),
        r.Form.Get("nombre"),
        r.Form.Get("descripcion"),
        cantidad,
        r.Form.Get("imagen_url"),
        r.Form.Get("ficha_tecnica_url"),
        r.Form.Get("modelo"),
        r.Form.Get("codigo_fabricante"),
        precioLista,
        precioMinimo,
        r.Form.Get("clasificacion"),
        r.Form.Get("serie"),
        r.Form.Get("pais_origen"),
        r.Form.Get("certificaciones"),
        requiereInstalacion,
        tiempoEntrega,
        stockMinimo,
        enPromocion,
        r.Form.Get("clave_producto_sat"),
        r.Form.Get("unidad_medida_sat"),
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

// Muestra el formulario de creación
func (m *Repository) MostrarFormularioCrear(w http.ResponseWriter, r *http.Request) {
    data := &models.TemplateData{
        CSRFToken: nosurf.Token(r), // Añade el token CSRF
    }
    render.RenderTemplate(w, "crear.page.tmpl", data)
}

// handler para todo el editar
func (m *Repository) EditarProducto(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    // Validar campos requeridos
    requiredFields := []string{"marca", "nombre", "sku", "cantidad", "precio_lista", "clasificacion"}
    for _, field := range requiredFields {
        if r.Form.Get(field) == "" {
            http.Error(w, "El campo "+field+" es requerido", http.StatusBadRequest)
            return
        }
    }

    // Convertir tipos
    cantidad, err := strconv.Atoi(r.Form.Get("cantidad"))
    if err != nil {
        http.Error(w, "Cantidad debe ser un número válido", http.StatusBadRequest)
        return
    }

    precioLista, err := strconv.ParseFloat(r.Form.Get("precio_lista"), 64)
    if err != nil {
        http.Error(w, "Precio debe ser un número válido", http.StatusBadRequest)
        return
    }

    precioMinimo, _ := strconv.ParseFloat(r.Form.Get("precio_minimo"), 64)
    stockMinimo, _ := strconv.Atoi(r.Form.Get("stock_minimo"))
    tiempoEntrega, _ := strconv.Atoi(r.Form.Get("tiempo_entrega"))
    requiereInstalacion := r.Form.Get("requiere_instalacion") == "on"
    enPromocion := r.Form.Get("en_promocion") == "on"

    // Actualizar en la base de datos
    _, err = m.App.DB.Exec(`
        UPDATE productos SET
            marca = ?,
            tipo = ?,
            sku = ?,
            nombre = ?,
            descripcion = ?,
            cantidad = ?,
            imagen_url = ?,
            ficha_tecnica_url = ?,
            modelo = ?,
            codigo_fabricante = ?,
            precio_lista = ?,
            precio_minimo = ?,
            clasificacion = ?,
            serie = ?,
            pais_origen = ?,
            certificaciones = ?,
            requiere_instalacion = ?,
            tiempo_entrega = ?,
            stock_minimo = ?,
            en_promocion = ?,
            clave_producto_sat = ?,
            unidad_medida_sat = ?,
            updated_at = NOW()
        WHERE id_producto = ?`,
        r.Form.Get("marca"),
        r.Form.Get("tipo"),
        r.Form.Get("sku"),
        r.Form.Get("nombre"),
        r.Form.Get("descripcion"),
        cantidad,
        r.Form.Get("imagen_url"),
        r.Form.Get("ficha_tecnica_url"),
        r.Form.Get("modelo"),
        r.Form.Get("codigo_fabricante"),
        precioLista,
        precioMinimo,
        r.Form.Get("clasificacion"),
        r.Form.Get("serie"),
        r.Form.Get("pais_origen"),
        r.Form.Get("certificaciones"),
        requiereInstalacion,
        tiempoEntrega,
        stockMinimo,
        enPromocion,
        r.Form.Get("clave_producto_sat"),
        r.Form.Get("unidad_medida_sat"),
        id,
    )

    if err != nil {
        log.Println("Error al actualizar producto:", err)
        http.Error(w, "Error al guardar cambios", http.StatusInternalServerError)
        return
    }

    m.App.Session.Put(r.Context(), "flash", "Producto actualizado correctamente")
    http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

// Muestra el formulario de editar
func (m *Repository) MostrarFormularioEditar(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    var producto models.Producto
    err = m.App.DB.QueryRow(`
    SELECT 
        id_producto, marca, tipo, sku, nombre, descripcion, cantidad,
        COALESCE(imagen_url, '') as imagen_url,
        COALESCE(ficha_tecnica_url, '') as ficha_tecnica_url,
        COALESCE(modelo, '') as modelo,
        COALESCE(codigo_fabricante, '') as codigo_fabricante,
        COALESCE(precio_lista, 0) as precio_lista,
        COALESCE(precio_minimo, 0) as precio_minimo,
        COALESCE(clasificacion, '') as clasificacion,
        COALESCE(serie, '') as serie,
        COALESCE(pais_origen, '') as pais_origen,
        COALESCE(certificaciones, '') as certificaciones,
        requiere_instalacion,
        COALESCE(tiempo_entrega, 0) as tiempo_entrega,
        COALESCE(stock_minimo, 0) as stock_minimo,
        en_promocion,
        COALESCE(clave_producto_sat, '') as clave_producto_sat,
        COALESCE(unidad_medida_sat, 'PIEZA') as unidad_medida_sat
    FROM productos 
    WHERE id_producto = ?`, id).Scan(
        &producto.IDProducto,
        &producto.Marca,
        &producto.Tipo,
        &producto.SKU,
        &producto.Nombre,
        &producto.Descripcion,
        &producto.Cantidad,
        &producto.ImagenURL,
        &producto.FichaTecnicaURL,
        &producto.Modelo,
        &producto.CodigoFabricante,
        &producto.PrecioLista,
        &producto.PrecioMinimo,
        &producto.Clasificacion,
        &producto.Serie,
        &producto.PaisOrigen,
        &producto.Certificaciones,
        &producto.RequiereInstalacion,
        &producto.TiempoEntrega,
        &producto.StockMinimo,
        &producto.EnPromocion,
        &producto.ClaveProductoSAT,
        &producto.UnidadMedidaSAT,
    )

    if err != nil {
        log.Println("Error al obtener producto:", err)
        http.Error(w, "Producto no encontrado", http.StatusNotFound)
        return
    }

    data := &models.TemplateData{
        Producto:   producto,
        CSRFToken: nosurf.Token(r),
    }
    
    render.RenderTemplate(w, "editar.page.tmpl", data)
}

// Handler para eliminar producto
func (m *Repository) EliminarProducto(w http.ResponseWriter, r *http.Request) {
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