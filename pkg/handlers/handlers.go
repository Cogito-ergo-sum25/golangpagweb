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

// TODO LO DE INVENTARIO

// PANTALLA INICIO INVENTARIO
func (m *Repository) Inventario(w http.ResponseWriter, r *http.Request) {
    rows, err := m.App.DB.Query(`
    SELECT 
        p.id_producto,
        p.sku, 
        m.nombre as marca,
        c.nombre as clasificacion,
        p.nombre_corto,
        p.modelo,
        p.nombre,
        p.version,
        p.serie,
        p.codigo_fabricante,
        p.descripcion
    FROM productos p
    LEFT JOIN marcas m ON p.id_marca = m.id_marca
    LEFT JOIN clasificaciones c ON p.id_clasificacion = c.id_clasificacion
    ORDER BY p.id_producto DESC`)

    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al consultar productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    defer rows.Close()

    var productos []models.Producto
    for rows.Next() {
        var p models.Producto
        err := rows.Scan(
            &p.IDProducto,
            &p.SKU,
            &p.Marca,       // Nuevo campo
            &p.Clasificacion, // Nuevo campo
            &p.NombreCorto,
            &p.Modelo,
            &p.Nombre,
            &p.Version,
            &p.Serie,
            &p.CodigoFabricante,
            &p.Descripcion,
        )
        if err != nil {
            m.App.Session.Put(r.Context(), "error", "Error al leer producto: "+err.Error())
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
        productos = append(productos, p)
    }

    if err = rows.Err(); err != nil {
        m.App.Session.Put(r.Context(), "error", "Error después de leer productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Productos: productos,
        CSRFToken: nosurf.Token(r),
    }
    
    render.RenderTemplate(w, "inventario.page.tmpl", data)
}

// Muestra el formulario de creación
func (m *Repository) MostrarFormularioCrear(w http.ResponseWriter, r *http.Request) {
	// Obtener datos para los selects
	marcas, err := m.ObtenerMarcas()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener marcas")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	tipos, err := m.ObtenerTiposProducto()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener tipos de producto")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	clasificaciones, err := m.ObtenerClasificaciones()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener clasificaciones")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	paises, err := m.ObtenerPaises()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener países")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	certificaciones, err := m.ObtenerCertificaciones()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener certificaciones")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	data := &models.TemplateData{
		Marcas:          marcas,
		TiposProducto:   tipos,
		Clasificaciones: clasificaciones,
		Paises:          paises,
		Certificaciones: certificaciones,
		CSRFToken:       nosurf.Token(r),
	}
	
	render.RenderTemplate(w, "crear-producto.page.tmpl", data)
}

// handler para crear producto
func (m *Repository) CrearProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		m.App.Session.Put(r.Context(), "error", "Método no permitido")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
		http.Redirect(w, r, "/crear-producto", http.StatusSeeOther)
		return
	}

	// Validar campos requeridos
	requiredFields := []string{"id_marca", "id_tipo", "id_clasificacion", "id_pais_origen", "nombre", "sku"}
	for _, field := range requiredFields {
		if r.Form.Get(field) == "" {
			m.App.Session.Put(r.Context(), "error", "El campo "+field+" es requerido")
			http.Redirect(w, r, "/crear-producto", http.StatusSeeOther)
			return
		}
	}

	// Convertir IDs a enteros
	idMarca, _ := strconv.Atoi(r.Form.Get("id_marca"))
	idTipo, _ := strconv.Atoi(r.Form.Get("id_tipo"))
	idClasificacion, _ := strconv.Atoi(r.Form.Get("id_clasificacion"))
	idPaisOrigen, _ := strconv.Atoi(r.Form.Get("id_pais_origen"))

	// Validar que los IDs existan
	if !m.ExisteID("marcas", idMarca) || !m.ExisteID("tipos_producto", idTipo) || 
	   !m.ExisteID("clasificaciones", idClasificacion) || !m.ExisteID("paises", idPaisOrigen) {
		m.App.Session.Put(r.Context(), "error", "Uno o más IDs no son válidos")
		http.Redirect(w, r, "/crear-producto", http.StatusSeeOther) 
		return
	}

	// Insertar el producto
	stmt := `
		INSERT INTO productos (
			id_marca, id_tipo, id_clasificacion, id_pais_origen,
			sku, nombre, nombre_corto, modelo, version, serie,
			codigo_fabricante, descripcion, imagen_url, ficha_tecnica_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.App.DB.Exec(stmt,
		idMarca,
		idTipo,
		idClasificacion,
		idPaisOrigen,
		r.Form.Get("sku"),
		r.Form.Get("nombre"),
		r.Form.Get("nombre_corto"),
		r.Form.Get("modelo"),
		r.Form.Get("version"),
		r.Form.Get("serie"),
		r.Form.Get("codigo_fabricante"),
		r.Form.Get("descripcion"),
		r.Form.Get("imagen_url"),
		r.Form.Get("ficha_tecnica_url"),
	)

	if err != nil {
		log.Println("Error al insertar producto:", err)
		m.App.Session.Put(r.Context(), "error", "Error al crear el producto")
		http.Redirect(w, r, "/crear-producto", http.StatusSeeOther)
		return
	}

	// Obtener ID del producto insertado
	idProducto, err := result.LastInsertId()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener ID del producto")
		http.Redirect(w, r, "/crear-producto", http.StatusSeeOther)
		return
	}

	// Procesar certificaciones seleccionadas
	certificaciones := r.Form["certificaciones"]
	for _, idCertStr := range certificaciones {
		idCert, err := strconv.Atoi(idCertStr)
		if err != nil {
			continue
		}

		if !m.ExisteID("certificacion", idCert) {
			continue
		}

		_, err = m.App.DB.Exec(
			"INSERT INTO producto_certificaciones (id_producto, id_certificacion) VALUES (?, ?)",
			idProducto,
			idCert,
		)
		if err != nil {
			log.Println("Error al insertar certificación:", err)
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Producto creado exitosamente")
	http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

// Muestra el formulario de edición
func (m *Repository) MostrarFormularioEditar(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID inválido")
        http.Redirect(w, r, "/inventario", http.StatusSeeOther)
        return
    }

    // Obtener el producto
    var producto models.Producto
    err = m.App.DB.QueryRow(`
        SELECT 
            p.id_producto, p.id_marca, p.id_tipo, p.id_clasificacion, p.id_pais_origen,
            p.sku, p.nombre, p.nombre_corto, p.modelo, p.version, p.serie,
            p.codigo_fabricante, p.descripcion, p.imagen_url, p.ficha_tecnica_url
        FROM productos p
        WHERE p.id_producto = ?`, id).Scan(
        &producto.IDProducto, &producto.IDMarca, &producto.IDTipo,
        &producto.IDClasificacion, &producto.IDPaisOrigen,
        &producto.SKU, &producto.Nombre, &producto.NombreCorto,
        &producto.Modelo, &producto.Version, &producto.Serie,
        &producto.CodigoFabricante, &producto.Descripcion,
        &producto.ImagenURL, &producto.FichaTecnicaURL,
    )

    if err != nil {
        log.Println("Error al obtener producto:", err)
        m.App.Session.Put(r.Context(), "error", "Producto no encontrado")
        http.Redirect(w, r, "/inventario", http.StatusSeeOther)
        return
    }

    // Obtener datos para los selects
    marcas, _ := m.ObtenerMarcas()
    tipos, _ := m.ObtenerTiposProducto()
    clasificaciones, _ := m.ObtenerClasificaciones()
    paises, _ := m.ObtenerPaises()
    certificaciones, _ := m.ObtenerCertificaciones()

    // Obtener certificaciones del producto
    var certs []models.Certificacion
    rows, err := m.App.DB.Query(`
        SELECT c.id_certificacion, c.nombre 
        FROM producto_certificaciones pc
        JOIN certificaciones c ON pc.id_certificacion = c.id_certificacion
        WHERE pc.id_producto = ?`, id)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var c models.Certificacion
            if err := rows.Scan(&c.IDCertificacion, &c.Nombre); err == nil {
                certs = append(certs, c)
            }
        }
    }

    data := &models.TemplateData{
        Producto:         producto,
        Marcas:          marcas,
        TiposProducto:   tipos,
        Clasificaciones: clasificaciones,
        Paises:          paises,
        Certificaciones: certificaciones,
        CertificacionesProducto: certs,
        CSRFToken:       nosurf.Token(r),
    }
    
    render.RenderTemplate(w, "editar-producto.page.tmpl", data)
}

// Handler para actualizar producto
func (m *Repository) EditarProducto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "ID inválido")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
		http.Redirect(w, r, "/editar-producto/"+idStr, http.StatusSeeOther)
		return
	}

	// Validar campos requeridos
	requiredFields := []string{"id_marca", "id_tipo", "id_clasificacion", "id_pais_origen", "nombre", "sku"}
	for _, field := range requiredFields {
		if r.Form.Get(field) == "" {
			m.App.Session.Put(r.Context(), "error", "El campo "+field+" es requerido")
			http.Redirect(w, r, "/editar-producto/"+idStr, http.StatusSeeOther)
			return
		}
	}

	// Convertir IDs
	idMarca, _ := strconv.Atoi(r.Form.Get("id_marca"))
	idTipo, _ := strconv.Atoi(r.Form.Get("id_tipo"))
	idClasificacion, _ := strconv.Atoi(r.Form.Get("id_clasificacion"))
	idPaisOrigen, _ := strconv.Atoi(r.Form.Get("id_pais_origen"))

	// Validar que los IDs existan
	if !m.ExisteID("marca", idMarca) || !m.ExisteID("tipo", idTipo) || 
	   !m.ExisteID("clasificacion", idClasificacion) || !m.ExisteID("pais", idPaisOrigen) {
		m.App.Session.Put(r.Context(), "error", "Uno o más IDs no son válidos")
		http.Redirect(w, r, "/editar-producto/"+idStr, http.StatusSeeOther)
		return
	}

	// Actualizar producto
	_, err = m.App.DB.Exec(`
		UPDATE productos SET
			id_marca = ?,
			id_tipo = ?,
			id_clasificacion = ?,
			id_pais_origen = ?,
			sku = ?,
			nombre = ?,
			nombre_corto = ?,
			modelo = ?,
			version = ?,
			serie = ?,
			codigo_fabricante = ?,
			descripcion = ?,
			imagen_url = ?,
			ficha_tecnica_url = ?,
			updated_at = NOW()
		WHERE id_producto = ?`,
		idMarca,
		idTipo,
		idClasificacion,
		idPaisOrigen,
		r.Form.Get("sku"),
		r.Form.Get("nombre"),
		r.Form.Get("nombre_corto"),
		r.Form.Get("modelo"),
		r.Form.Get("version"),
		r.Form.Get("serie"),
		r.Form.Get("codigo_fabricante"),
		r.Form.Get("descripcion"),
		r.Form.Get("imagen_url"),
		r.Form.Get("ficha_tecnica_url"),
		id,
	)

	if err != nil {
		log.Println("Error al actualizar producto:", err)
		m.App.Session.Put(r.Context(), "error", "Error al guardar cambios")
		http.Redirect(w, r, "/editar-producto/"+idStr, http.StatusSeeOther)
		return
	}

	// Manejar certificaciones
	// Primero, borrar todas las certificaciones existentes para este producto
	_, err = m.App.DB.Exec("DELETE FROM producto_certificaciones WHERE id_producto = ?", id)
	if err != nil {
		log.Println("Error al borrar certificaciones:", err)
	}

	// Luego, insertar las nuevas
	certificaciones := r.Form["certificaciones"]
	for _, idCertStr := range certificaciones {
		idCert, err := strconv.Atoi(idCertStr)
		if err != nil {
			continue
		}
		
		if !m.ExisteID("certificacion", idCert) {
			continue
		}
		
		_, err = m.App.DB.Exec(
			"INSERT INTO producto_certificaciones (id_producto, id_certificacion) VALUES (?, ?)", 
			id, idCert)
		if err != nil {
			log.Println("Error al insertar certificación:", err)
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Producto actualizado correctamente")
	http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

// Handler para eliminar producto
func (m *Repository) EliminarProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		m.App.Session.Put(r.Context(), "error", "Método no permitido")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "ID inválido")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	// Primero eliminar las relaciones en producto_certificaciones
	_, err = m.App.DB.Exec("DELETE FROM producto_certificaciones WHERE id_producto = ?", id)
	if err != nil {
		log.Println("Error al borrar certificaciones:", err)
	}

	// Luego eliminar el producto
	_, err = m.App.DB.Exec("DELETE FROM productos WHERE id_producto = ?", id)
	if err != nil {
		log.Println("Error al borrar producto:", err)
		m.App.Session.Put(r.Context(), "error", "Error al eliminar el producto")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Producto eliminado")
	http.Redirect(w, r, "/inventario", http.StatusSeeOther)
}

// TODO LO CATALOGO
/*

func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
    // Obtener parámetros de los filtros
    marca := r.URL.Query().Get("marca")
    clasificacion := r.URL.Query().Get("clasificacion")
    busqueda := r.URL.Query().Get("busqueda")
    enPromocion := r.URL.Query().Get("en_promocion") == "true"

    // Consulta SQL con filtros
    query := "SELECT id_producto, marca, nombre, imagen_url, precio_lista, en_promocion FROM productos WHERE 1=1"
    var args []interface{}

    if marca != "" {
        query += " AND marca = ?"
        args = append(args, marca)
    }

    if clasificacion != "" {
        query += " AND clasificacion = ?"
        args = append(args, clasificacion)
    }

    if busqueda != "" {
        query += " AND nombre LIKE ?"
        args = append(args, "%"+busqueda+"%")
    }

    if enPromocion {
        query += " AND en_promocion = TRUE"
    }

    // Ejecutar consulta
    rows, err := m.App.DB.Query(query, args...)
    if err != nil {
        http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var productos []models.Producto
    for rows.Next() {
        var p models.Producto
        var promocion bool
        err := rows.Scan(&p.IDProducto, &p.Marca, &p.Nombre, &p.ImagenURL, &p.PrecioLista, &promocion)
        if err != nil {
            http.Error(w, "Error al leer producto: "+err.Error(), http.StatusInternalServerError)
            return
        }
        p.EnPromocion = promocion
        productos = append(productos, p)
    }

    // Obtener opciones para los selectores de filtro
    marcas, _ := m.obtenerDatosUnicos("SELECT DISTINCT marca FROM productos WHERE marca != ''")
    clasificaciones, _ := m.obtenerDatosUnicos("SELECT DISTINCT clasificacion FROM productos WHERE clasificacion != ''")

    // Preparar datos para la plantilla usando tu estructura TemplateData
    data := &models.TemplateData{
        Productos: productos, // Esto ya está en tu struct
        Data: map[string]interface{}{
            "Marcas":          marcas,
            "Clasificaciones": clasificaciones,
            "Filtros": map[string]interface{}{
                "Marca":         marca,
                "Clasificacion": clasificacion,
                "Busqueda":      busqueda,
                "EnPromocion":   enPromocion,
            },
        },
        CSRFToken: nosurf.Token(r),
    }

    render.RenderTemplate(w, "catalogo.page.tmpl", data)
}

func (m *Repository) obtenerDatosUnicos(query string) ([]string, error) {
    rows, err := m.App.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var resultados []string
    for rows.Next() {
        var valor string
        if err := rows.Scan(&valor); err != nil {
            return nil, err
        }
        resultados = append(resultados, valor)
    }
    return resultados, nil
}

func (m *Repository) VerProducto(w http.ResponseWriter, r *http.Request) {
    // Obtener ID del producto de la URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID de producto inválido", http.StatusBadRequest)
        return
    }

    // Consulta SQL para obtener todos los detalles del producto
    var producto models.Producto
    err = m.App.DB.QueryRow(`
        SELECT 
            id_producto, marca, tipo, sku, nombre, descripcion, cantidad,
            COALESCE(imagen_url, '') as imagen_url,
            COALESCE(ficha_tecnica_url, '') as ficha_tecnica_url,
            COALESCE(modelo, '') as modelo,
            COALESCE(codigo_fabricante, '') as codigo_fabricante,
            precio_lista,
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
            COALESCE(unidad_medida_sat, '') as unidad_medida_sat
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
        if err == sql.ErrNoRows {
            http.Error(w, "Producto no encontrado", http.StatusNotFound)
        } else {
            log.Println("Error al obtener producto:", err)
            http.Error(w, "Error interno al obtener el producto", http.StatusInternalServerError)
        }
        return
    }

    // Preparar datos para la plantilla
    data := &models.TemplateData{
        Producto:   producto,
        CSRFToken: nosurf.Token(r),
    }
        render.RenderTemplate(w, "ver_producto.page.tmpl", data)
}

*/