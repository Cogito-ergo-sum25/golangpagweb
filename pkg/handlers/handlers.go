package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	render.RenderTemplate(w, "home/home.page.tmpl", &models.TemplateData{})
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
    
    render.RenderTemplate(w, "inventario/inventario.page.tmpl", data)
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
	
	render.RenderTemplate(w, "inventario/crear-producto.page.tmpl", data)
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
    
    render.RenderTemplate(w, "inventario/editar-producto.page.tmpl", data)
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

// TODO LO DE PROYECTOS
func (m *Repository) ProyectosVista(w http.ResponseWriter, r *http.Request) {
	/*proyectos, err := m.ObtenerProyectosConRelaciones()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener proyectos: "+err.Error())
		log.Println("Error al obtener proyectos:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}*/

	data := &models.TemplateData{
		//Proyectos: proyectos,
		CSRFToken: nosurf.Token(r),
	}

	render.RenderTemplate(w, "proyectos/proyectos-vista.page.tmpl", data)
}

func (m *Repository) MostrarNuevoProyecto(w http.ResponseWriter, r *http.Request) {
    // Obtener productos
    productos, err := m.ObtenerTodosProductos()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Obtener licitaciones
    licitaciones, err := m.ObtenerLicitacionesParaSelect()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo licitaciones: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Productos:    productos,
        Licitaciones: licitaciones,
        CSRFToken:    nosurf.Token(r),
    }
    render.RenderTemplate(w, "proyectos/nuevo-proyecto.page.tmpl", data)
}

func (m *Repository) NuevoProyecto(w http.ResponseWriter, r *http.Request) {
    // 1. Verificar método HTTP
    if r.Method != http.MethodPost {
        m.App.Session.Put(r.Context(), "error", "Método no permitido")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 2. Parsear formulario
    err := r.ParseForm()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 3. Validar campos requeridos
    requiredFields := []string{"nombre", "descripcion", "id_licitacion", "fecha_inicio"}
    for _, field := range requiredFields {
        if r.Form.Get(field) == "" {
            m.App.Session.Put(r.Context(), "error", "El campo "+field+" es requerido")
            http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
            return
        }
    }

    // 4. Convertir IDs a enteros
    idLicitacion, err := strconv.Atoi(r.Form.Get("id_licitacion"))
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID de licitación inválido")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 5. Validar que la licitación exista
    if !m.ExisteID("licitaciones", idLicitacion) {
        m.App.Session.Put(r.Context(), "error", "La licitación seleccionada no existe")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 6. Procesar fechas
    fechaInicio, err := time.Parse("2006-01-02", r.Form.Get("fecha_inicio"))
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Formato de fecha de inicio inválido")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    var fechaFin time.Time
    if fechaFinStr := r.Form.Get("fecha_fin"); fechaFinStr != "" {
        fechaFin, err = time.Parse("2006-01-02", fechaFinStr)
        if err != nil {
            m.App.Session.Put(r.Context(), "error", "Formato de fecha fin inválido")
            http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
            return
        }
    }

    // 7. Insertar el proyecto
    stmt := `
        INSERT INTO proyectos (
            nombre, descripcion, id_licitacion, 
            fecha_inicio, fecha_fin, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?)`

    result, err := m.App.DB.Exec(stmt,
        r.Form.Get("nombre"),
        r.Form.Get("descripcion"),
        idLicitacion,
        fechaInicio,
        fechaFin,
        time.Now(),
        time.Now(),
    )

    if err != nil {
        log.Println("Error al insertar proyecto:", err)
        m.App.Session.Put(r.Context(), "error", "Error al crear el proyecto")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 8. Obtener ID del proyecto insertado
    idProyecto, err := result.LastInsertId()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener ID del proyecto")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 9. Procesar productos seleccionados
    if err := m.procesarProductosProyecto(r, idProyecto); err != nil {
        log.Println("Error al procesar productos:", err)
        m.App.Session.Put(r.Context(), "error", "Error al guardar los productos del proyecto")
        http.Redirect(w, r, "/proyectos/nuevo", http.StatusSeeOther)
        return
    }

    // 10. Redireccionar con mensaje de éxito
    m.App.Session.Put(r.Context(), "flash", "Proyecto creado exitosamente!")
    http.Redirect(w, r, fmt.Sprintf("/proyectos/%d", idProyecto), http.StatusSeeOther)
}

// Función auxiliar para procesar productos del proyecto
func (m *Repository) procesarProductosProyecto(r *http.Request, idProyecto int64) error {
    // Obtener todos los valores del formulario para productos
    productos := make([]struct {
        IDProducto      string
        Cantidad       string
        PrecioUnitario string
        Especificaciones string
    }, 0)

    // Recorrer todos los productos del formulario
    for key, values := range r.Form {
        if strings.HasPrefix(key, "productos[") && strings.Contains(key, "][id_producto]") {
            // Extraer el índice del producto (ej: "productos[0][id_producto]" -> "0")
            idx := strings.Split(key, "[")[1]
            idx = strings.Split(idx, "]")[0]

            producto := struct {
                IDProducto      string
                Cantidad       string
                PrecioUnitario string
                Especificaciones string
            }{
                IDProducto:      values[0],
                Cantidad:       r.Form.Get(fmt.Sprintf("productos[%s][cantidad]", idx)),
                PrecioUnitario: r.Form.Get(fmt.Sprintf("productos[%s][precio_unitario]", idx)),
                Especificaciones: r.Form.Get(fmt.Sprintf("productos[%s][especificaciones]", idx)),
            }

            // Solo agregar si se seleccionó un producto
            if producto.IDProducto != "" {
                productos = append(productos, producto)
            }
        }
    }

    // Insertar cada producto asociado al proyecto
    for _, p := range productos {
        idProducto, err := strconv.Atoi(p.IDProducto)
        if err != nil {
            continue // Saltar si el ID no es válido
        }

        cantidad, err := strconv.Atoi(p.Cantidad)
        if err != nil || cantidad <= 0 {
            continue // Saltar si la cantidad no es válida
        }

        precioUnitario, err := strconv.ParseFloat(p.PrecioUnitario, 64)
        if err != nil || precioUnitario <= 0 {
            continue // Saltar si el precio no es válido
        }

        // Verificar que el producto exista
        if !m.ExisteID("productos", idProducto) {
            continue
        }

        // Insertar en la tabla de relación
        _, err = m.App.DB.Exec(`
            INSERT INTO productos_proyecto (
                id_proyecto, id_producto, cantidad, 
                precio_unitario, especificaciones, created_at
            ) VALUES (?, ?, ?, ?, ?, ?)`,
            idProyecto,
            idProducto,
            cantidad,
            precioUnitario,
            p.Especificaciones,
            time.Now(),
        )
        if err != nil {
            return err
        }
    }

    return nil
}

// TODO LO CATALOGO

func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
    // Obtener parámetros de los filtros
    marca := r.URL.Query().Get("marca")
    clasificacion := r.URL.Query().Get("clasificacion")
    busqueda := r.URL.Query().Get("busqueda")

    // Obtener todos los productos
    productos, err := m.ObtenerTodosProductos()
    if err != nil {
        http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Aplicar filtros en memoria
    var filteredProductos []models.Producto
    for _, p := range productos {
        if marca != "" && p.Marca != marca {
            continue
        }
        if clasificacion != "" && p.Clasificacion != clasificacion {
            continue
        }
        if busqueda != "" && !strings.Contains(strings.ToLower(p.Nombre), strings.ToLower(busqueda)) {
            continue
        }
        filteredProductos = append(filteredProductos, p)
    }

    // Obtener opciones para los selectores de filtro
    marcas, err := m.obtenerDatosUnicos("SELECT DISTINCT nombre FROM marcas WHERE nombre != ''")
    if err != nil {
        http.Error(w, "Error al obtener marcas: "+err.Error(), http.StatusInternalServerError)
        return
    }
    clasificaciones, err := m.obtenerDatosUnicos("SELECT DISTINCT nombre FROM clasificaciones WHERE nombre != ''")
    if err != nil {
        http.Error(w, "Error al obtener clasificaciones: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Preparar datos para la plantilla
    data := &models.TemplateData{
        Productos: filteredProductos,
        Data: map[string]interface{}{
            "Marcas":         marcas,
            "Clasificaciones": clasificaciones,
            "Filtros": map[string]interface{}{
                "Marca":        marca,
                "Clasificacion": clasificacion,
                "Busqueda":     busqueda,
            },
        },
        CSRFToken: nosurf.Token(r),
    }

    render.RenderTemplate(w, "catalogo/catalogo.page.tmpl", data)
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

func (m *Repository) ProductoDetalles(w http.ResponseWriter, r *http.Request) {
    // Obtener el ID del producto desde la URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Error al convertir ID: %v", err)
        m.App.Session.Put(r.Context(), "error", "ID inválido")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }

    // Obtener el producto con JOIN para Marca y Clasificación
    var producto models.Producto
    err = m.App.DB.QueryRow(`
        SELECT
            p.id_producto, p.id_marca, p.id_tipo, p.id_clasificacion, p.id_pais_origen,
            p.sku, p.nombre, p.nombre_corto, p.modelo, p.version, p.serie,
            p.codigo_fabricante, p.descripcion, p.imagen_url, p.ficha_tecnica_url,
            m.nombre AS marca, c.nombre AS clasificacion
        FROM productos p
        LEFT JOIN marcas m ON p.id_marca = m.id_marca
        LEFT JOIN clasificaciones c ON p.id_clasificacion = c.id_clasificacion
        WHERE p.id_producto = ?`, id).Scan(
        &producto.IDProducto, &producto.IDMarca, &producto.IDTipo,
        &producto.IDClasificacion, &producto.IDPaisOrigen,
        &producto.SKU, &producto.Nombre, &producto.NombreCorto,
        &producto.Modelo, &producto.Version, &producto.Serie,
        &producto.CodigoFabricante, &producto.Descripcion,
        &producto.ImagenURL, &producto.FichaTecnicaURL,
        &producto.Marca, &producto.Clasificacion,
    )
    if err != nil {
        log.Printf("Error al obtener producto %d: %v", id, err)
        m.App.Session.Put(r.Context(), "error", "Producto no encontrado")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }

    // Obtener datos para los selects (si es necesario mostrarlos en la vista)
    marcas, err := m.ObtenerMarcas()
    if err != nil {
        log.Printf("Error al obtener marcas: %v", err)
        m.App.Session.Put(r.Context(), "error", "Error al cargar datos")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }
    tipos, err := m.ObtenerTiposProducto()
    if err != nil {
        log.Printf("Error al obtener tipos: %v", err)
        m.App.Session.Put(r.Context(), "error", "Error al cargar datos")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }
    clasificaciones, err := m.ObtenerClasificaciones()
    if err != nil {
        log.Printf("Error al obtener clasificaciones: %v", err)
        m.App.Session.Put(r.Context(), "error", "Error al cargar datos")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }
    paises, err := m.ObtenerPaises()
    if err != nil {
        log.Printf("Error al obtener países: %v", err)
        m.App.Session.Put(r.Context(), "error", "Error al cargar datos")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }
    certificaciones, err := m.ObtenerCertificaciones()
    if err != nil {
        log.Printf("Error al obtener certificaciones: %v", err)
        m.App.Session.Put(r.Context(), "error", "Error al cargar datos")
        http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
        return
    }

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
            if err := rows.Scan(&c.IDCertificacion, &c.Nombre); err != nil {
                log.Printf("Error al escanear certificación: %v", err)
                continue
            }
            certs = append(certs, c)
        }
    } else {
        log.Printf("Error al obtener certificaciones del producto %d: %v", id, err)
    }

    // Preparar datos para la plantilla
    data := &models.TemplateData{
        Producto:                producto,
        Marcas:                 marcas,
        TiposProducto:          tipos,
        Clasificaciones:        clasificaciones,
        Paises:                 paises,
        Certificaciones:        certificaciones,
        CertificacionesProducto: certs,
        CSRFToken:              nosurf.Token(r),
    }

    // Log para depurar datos enviados
    log.Printf("Producto: ID=%d, Nombre=%s, Marca=%s, Clasificacion=%s, Modelo=%s, Version=%s, Serie=%s, CodigoFabricante=%s",
        producto.IDProducto, producto.Nombre, producto.Marca, producto.Clasificacion,
        producto.Modelo, producto.Version, producto.Serie, producto.CodigoFabricante)

    // Renderizar la plantilla
    render.RenderTemplate(w, "catalogo/producto-detalle.page.tmpl", data)
}




// TODO LO DE LICITACIONES
func (m *Repository) Licitaciones(w http.ResponseWriter, r *http.Request) {
    licitaciones, err := m.ObtenerTodasLicitaciones()
    if err != nil {
        http.Error(w, "No se pudieron obtener las licitaciones", http.StatusInternalServerError)
        return
    }

    data := &models.TemplateData{
        Licitaciones: licitaciones,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/licitaciones.page.tmpl", data)
}

func (m *Repository) MostrarNuevaLicitacion(w http.ResponseWriter, r *http.Request) {
    // Obtener todas las entidades para el select
    entidades, err := m.ObtenerTodasEntidades()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo entidades: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Entidades:  entidades,
        CSRFToken:  nosurf.Token(r),
        Data: map[string]interface{}{
            "tipos": []string{"Directa", "Apoyo", "Estudio de mercado", "Adjudicación directa", "Producto no adecuado", "No solicitan productos INTEVI"},
            "caracter": []string{"Internacional - Cobertura Tratados", "Nacional", "Internacional Abierto"},
            "criterio": []string{"Mixta", "Binario", "Puntos y porcentajes"},
            "estatus": []string{"Vigente", "En aclaraciones", "Presentada", "Finalizada"},
            
        },
    }
    render.RenderTemplate(w, "licitaciones/nueva-licitacion.page.tmpl", data)
}

func (m *Repository) CrearNuevaLicitacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/licitaciones/nueva", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "No se pudo procesar el formulario", http.StatusBadRequest)
		return
	}

    idEntidad, err := strconv.Atoi(r.FormValue("id_entidad"))
    if err != nil {
        http.Error(w, "ID de entidad inválido", http.StatusBadRequest)
        return
    }

	// Parsear campos del formulario
	licitacion := models.Licitacion{
		IDEntidad:            idEntidad,
		NumContratacion:      r.FormValue("num_contratacion"),
		Caracter:             r.FormValue("caracter"),
		Nombre:               r.FormValue("nombre"),
		Estatus:              r.FormValue("estatus"),
		Tipo:                 r.FormValue("tipo"),
		FechaJunta:           parseDate(r.FormValue("fecha_junta")),
		FechaPropuestas:      parseDate(r.FormValue("fecha_propuestas")),
		FechaFallo:           parseDate(r.FormValue("fecha_fallo")),
		FechaEntrega:         parseDate(r.FormValue("fecha_entrega")),
		TiempoEntrega:        r.FormValue("tiempo_entrega"),
		Revisada:             r.FormValue("revisada") == "on",
		Intevi:               r.FormValue("intevi") == "on",
		ObservacionesGenerales: r.FormValue("observaciones_generales"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
        CriterioEvaluacion: r.FormValue("criterio_evaluacion"),
	}
	// Insertar en la base de datos
	err = m.InsertarLicitacion(licitacion)
	if err != nil {
        fmt.Println("ERROR:", err) 
		http.Error(w, "Error al insertar la licitación", http.StatusInternalServerError)
		return
	}

	// Redirigir o mostrar mensaje de éxito
	http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
}

func (m *Repository) MostrarFormularioEditarLicitacion(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    licitacion, err := m.ObtenerLicitacionPorID(id)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "No se pudo obtener la licitación: "+err.Error())
        http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
        return
    }

    entidades, err := m.ObtenerTodasEntidades()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo entidades: "+err.Error())
        http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Licitacion: licitacion,
        Entidades:  entidades,
        CSRFToken:  nosurf.Token(r),
        Data: map[string]interface{}{
            "tipos": []string{"Directa", "Apoyo", "Estudio de mercado", "Adjudicación directa", "Producto no adecuado", "No solicitan productos INTEVI"},
            "caracter": []string{"Internacional - Cobertura Tratados", "Nacional", "Internacional Abierto"},
            "criterio": []string{"Mixta", "Binario", "Puntos y porcentajes"},
            "estatus": []string{"Vigente", "En aclaraciones","Presentada", "Finalizada"},
        },
    }

    render.RenderTemplate(w, "licitaciones/editar-licitacion.page.tmpl", data)
}

func (m *Repository) EditarLicitacion(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
        return
    }

    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Formulario inválido", http.StatusBadRequest)
        return
    }

    idEntidad, err := strconv.Atoi(r.FormValue("id_entidad"))
    if err != nil {
        http.Error(w, "ID de entidad inválido", http.StatusBadRequest)
        return
    }

    licitacion := models.Licitacion{
        IDLicitacion:         id,
        IDEntidad:            idEntidad,
        NumContratacion:      r.FormValue("num_contratacion"),
        Caracter:             r.FormValue("caracter"),
        Nombre:               r.FormValue("nombre"),
        Estatus:              r.FormValue("estatus"),
        Tipo:                 r.FormValue("tipo"),
        FechaJunta:           parseDate(r.FormValue("fecha_junta")),
        FechaPropuestas:      parseDate(r.FormValue("fecha_propuestas")),
        FechaFallo:           parseDate(r.FormValue("fecha_fallo")),
        FechaEntrega:         parseDate(r.FormValue("fecha_entrega")),
        TiempoEntrega:        r.FormValue("tiempo_entrega"),
        Revisada:             r.FormValue("revisada") == "on",
        Intevi:               r.FormValue("intevi") == "on",
        ObservacionesGenerales: r.FormValue("observaciones_generales"),
        UpdatedAt:            time.Now(),
        CriterioEvaluacion: r.FormValue("criterio_evaluacion"),

    }
    err = m.ActualizarLicitacion(licitacion)
    if err != nil {
        fmt.Println("ERROR:", err) 
        http.Error(w, "Error al actualizar la licitación", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
}

// CALENDARIO
func (m *Repository) Calendario(w http.ResponseWriter, r *http.Request) {
    licitaciones, err := m.ObtenerTodasLicitaciones()
    if err != nil {
        http.Error(w, "No se pudieron obtener las licitaciones", http.StatusInternalServerError)
        return
    }

    // Prepara eventos JS
    type Evento struct {
        Title         string `json:"title"`
        Start         string `json:"start"`
        Color         string            `json:"color,omitempty"` 
        ExtendedProps map[string]string `json:"extendedProps"`
    }

    var eventos []Evento
    for _, l := range licitaciones {
        eventos = append(eventos,
            Evento{
                Title: "📋 Junta: " + l.Nombre,
                Start: l.FechaJunta.Format("2006-01-02"),
                Color: "#0d6efd",
                ExtendedProps: map[string]string{
                    "num":     l.NumContratacion,
                    "tipo":    l.Tipo,
                    "estatus": l.Estatus,
                    "id":      fmt.Sprintf("%d", l.IDLicitacion),
                },
            },
            Evento{
                Title: "📑 Propuestas: " + l.Nombre,
                Start: l.FechaPropuestas.Format("2006-01-02"),
                Color: "#6610f2",
                ExtendedProps: map[string]string{
                    "num":     l.NumContratacion,
                    "tipo":    l.Tipo,
                    "estatus": l.Estatus,
                    "id":      fmt.Sprintf("%d", l.IDLicitacion),
                },
            },
            Evento{
                Title: "⛔ Fallo: " + l.Nombre,
                Start: l.FechaFallo.Format("2006-01-02"),
                Color: "#d66666",
                ExtendedProps: map[string]string{
                    "num":     l.NumContratacion,
                    "tipo":    l.Tipo,
                    "estatus": l.Estatus,
                    "id":      fmt.Sprintf("%d", l.IDLicitacion),
                },
            },
            Evento{
                Title: "📦 Entrega: " + l.Nombre,
                Start: l.FechaEntrega.Format("2006-01-02"),
                Color: "#28a745",
                ExtendedProps: map[string]string{
                    "num":     l.NumContratacion,
                    "tipo":    l.Tipo,
                    "estatus": l.Estatus,
                    "id":      fmt.Sprintf("%d", l.IDLicitacion),
                },
            },
        )
    }

    data := &models.TemplateData{
    Licitaciones: licitaciones,
    CSRFToken:    nosurf.Token(r),
    Data: map[string]interface{}{
        "Eventos": eventos,
        },
    }


    render.RenderTemplate(w, "calendario/calendario-vista.page.tmpl", data)
}



// TODO LO DE PARTIDAS
func (m *Repository) MostrarPartidasPorID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID de licitación inválido", http.StatusBadRequest)
		return
	}

	partidas, err := m.ObtenerPartidasPorLicitacionID(id)
	if err != nil {
		log.Printf("Error al obtener partidas para licitación %d: %v", id, err)
		http.Error(w, "Error al obtener las partidas", http.StatusInternalServerError)
		return
	}

    licitacion, err := m.ObtenerLicitacionPorID(id)
	if err != nil {
		http.Error(w, "Error al obtener licitación", http.StatusInternalServerError)
		return
	}

	data := &models.TemplateData{
		CSRFToken: nosurf.Token(r),
		Partidas:  partidas,
        Licitacion: licitacion,
	}

	render.RenderTemplate(w, "licitaciones/mostrar-partidas.page.tmpl", data)
}

func (m *Repository) MostrarAclaracionesLicitacion(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID de licitación inválido", http.StatusBadRequest)
		return
	}

    licitacion, err := m.ObtenerLicitacionPorID(id)
	if err != nil {
		http.Error(w, "Error al obtener licitación", http.StatusInternalServerError)
		return
	}

    partidas, err := m.ObtenerPartidasPorLicitacionID(id)
	if err != nil {
		log.Printf("Error al obtener partidas para licitación %d: %v", id, err)
		http.Error(w, "Error al obtener las partidas", http.StatusInternalServerError)
		return
	}

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    aclaraciones, err := m.ObtenerAclaracionesPorLicitacionID(id)
    if err != nil {
        log.Printf("Error al obtener aclaraciones de la licitación %d: %v", id, err)
        http.Error(w, "Error al obtener aclaraciones", http.StatusInternalServerError)
        return
    }


	data := &models.TemplateData{
		CSRFToken: nosurf.Token(r),
        Partidas:  partidas,
        Empresas:   empresas,
        Licitacion: licitacion,
        AclaracionesLicitacion: aclaraciones,
	}

	render.RenderTemplate(w, "licitaciones/aclaraciones-licitacion.page.tmpl", data)
}

func (m *Repository) MostrarNuevaAclaracionGeneral(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    licitacion, err := m.ObtenerLicitacionPorID(idPartida)
	if err != nil {
		http.Error(w, "Error al obtener licitación", http.StatusInternalServerError)
		return
	}

    // Obtener partida
    partidas, err := m.ObtenerPartidasPorLicitacionID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Empresas:   empresas,
        Partidas:    partidas,
        Licitacion: licitacion,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/nueva-aclaracion-licitacion.page.tmpl", data)
}

func (m *Repository) AgregarEmpresaExternaContextoAclaraciones(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    nombre := r.FormValue("nombre")
    idLicitacion := r.FormValue("id_licitacion")

    if nombre == "" || idLicitacion == "" {
        m.App.Session.Put(r.Context(), "error", "Todos los campos son obligatorios")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := m.AgregarEmpresaNueva(nombre)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al agregar la empresa")
        fmt.Println("Error al agregar empresa:", err)
        http.Redirect(w, r, "/nueva-aclaracion-general/"+idLicitacion, http.StatusSeeOther)
        return
    }

    m.App.Session.Put(r.Context(), "success", "Empresa agregada correctamente")
    http.Redirect(w, r, "/nueva-aclaracion-general/"+idLicitacion, http.StatusSeeOther)
}

func (m *Repository) CrearNuevaAclaracionGeneral(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "No se pudo procesar el formulario", http.StatusBadRequest)
		return
	}

	idLicitacion, err := strconv.Atoi(r.FormValue("id_licitacion"))
	if err != nil {
		http.Error(w, "ID de licitación inválido", http.StatusBadRequest)
		return
	}

	// id_partida puede ser opcional
	var partida *models.Partida
	if idPartidaStr := r.FormValue("id_partida"); idPartidaStr != "" {
		idPartida, err := strconv.Atoi(idPartidaStr)
		if err != nil {
			http.Error(w, "ID de partida inválido", http.StatusBadRequest)
			return
		}
		partida = &models.Partida{IDPartida: idPartida}
	}

	idEmpresa, err := strconv.Atoi(r.FormValue("id_empresa"))
	if err != nil {
		http.Error(w, "Empresa inválida", http.StatusBadRequest)
		return
	}

	// Campos opcionales
	fichaID, _ := strconv.Atoi(r.FormValue("ficha_tecnica_id"))
	puntosID, _ := strconv.Atoi(r.FormValue("id_puntos_tecnicos_modif"))

	// Aquí podrías agregar lógica extra si deseas distinguir preguntas técnicas o no
	aclaracion := models.AclaracionesLicitacion{
		IDLicitacion:           idLicitacion,
		Partida:                partida,
		IDEmpresa:              idEmpresa,
		Pregunta:               r.FormValue("pregunta"),
		Observaciones:          r.FormValue("observaciones"),
		FichaTecnicaID:         fichaID,
		IDPuntosTecnicosModif:  puntosID,
	}

	err = m.InsertarAclaracionGeneral(aclaracion)
	if err != nil {
		http.Error(w, "Error al guardar la aclaración: "+err.Error(), http.StatusInternalServerError)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Aclaración registrada correctamente")
	http.Redirect(w, r, fmt.Sprintf("/aclaraciones-licitacion/%d", idLicitacion), http.StatusSeeOther)
}

func (m *Repository) MostrarNuevaPartida(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id") // <- Extrae ID desde URL
    idLicitacion, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    licitacion, err := m.ObtenerLicitacionPorID(idLicitacion)
    if err != nil {
        http.Error(w, "No se pudo obtener la licitación", http.StatusInternalServerError)
        return
    }

    productos, err := m.ObtenerTodosProductos()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Productos:   productos,
        Licitacion:  licitacion,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/nueva-partida.page.tmpl", data)
}

func (m *Repository) MostrarEditarPartida(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener partida
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    // Obtener id_licitacion desde tabla intermedia
    idLicitacion, err := m.ObtenerIDLicitacionPorPartida(idPartida)
    if err != nil {
        http.Error(w, "No se pudo encontrar la licitación asociada", http.StatusInternalServerError)
        return
    }

    // Obtener la licitación
    licitacion, err := m.ObtenerLicitacionPorID(idLicitacion)
    if err != nil {
        http.Error(w, "No se pudo obtener la licitación", http.StatusInternalServerError)
        return
    }

    // Render
    data := &models.TemplateData{
        Partida:    partida,
        Licitacion: licitacion,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/editar-partida.page.tmpl", data)
}

func (m *Repository) EditarPartida(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
        return
    }

    idStr := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Formulario inválido", http.StatusBadRequest)
        return
    }

    // Parsear los campos del formulario
    numPartidaConv, _ := strconv.Atoi(r.FormValue("numero_partida_convocatoria"))
    cantidad, _ := strconv.Atoi(r.FormValue("cantidad"))
    cantidadMinima, _ := strconv.Atoi(r.FormValue("cantidad_minima"))
    cantidadMaxima, _ := strconv.Atoi(r.FormValue("cantidad_maxima"))
    garantia, _ := strconv.Atoi(r.FormValue("garantia"))
    fechaEntrega := parseDate(r.FormValue("fecha_de_entrega"))

    partida := models.Partida{
        IDPartida:              idPartida,
        NumPartidaConvocatoria: numPartidaConv,
        NombreDescripcion:      r.FormValue("nombre_descripcion"),
        Cantidad:               cantidad,
        CantidadMinima:         cantidadMinima,
        CantidadMaxima:         cantidadMaxima,
        NoFichaTecnica:         r.FormValue("no_ficha_tecnica"),
        TipoDeBien:             r.FormValue("tipo_de_bien"),
        ClaveCompendio:         r.FormValue("clave_compendio"),
        ClaveCucop:             r.FormValue("clave_cucop"),
        UnidadMedida:           r.FormValue("unidad_medida"),
        DiasDeEntrega:          r.FormValue("días_de_entrega"),
        FechaDeEntrega:         fechaEntrega,
        Garantia:               garantia,
        UpdatedAt:              time.Now(),
    }

    err = m.ActualizarPartida(partida)
    if err != nil {
        fmt.Println("Error al actualizar partida:", err)
        http.Error(w, "Error al actualizar la partida", http.StatusInternalServerError)
        return
    }

    idLicitacion, err := m.ObtenerIDLicitacionPorIDPartida(idPartida)
    if err != nil {
        fmt.Println("Error al obtener ID licitación de la partida:", err)
        http.Error(w, "Error interno", http.StatusInternalServerError)
        return
    }


    http.Redirect(w, r, fmt.Sprintf("/mostrar-partidas/%d", idLicitacion), http.StatusSeeOther)
}

func (m *Repository) CrearNuevaPartida(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    // Validaciones y conversiones seguras
    idLicitacion, err := strconv.Atoi(r.FormValue("id_licitacion"))
    if err != nil {
        http.Error(w, "ID de licitación inválido", http.StatusBadRequest)
        return
    }

    numeroPartida, err := strconv.Atoi(r.FormValue("numero_partida_convocatoria"))
    if err != nil {
        http.Error(w, "Número de partida inválido", http.StatusBadRequest)
        return
    }

    cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
    if err != nil {
        http.Error(w, "Cantidad inválida", http.StatusBadRequest)
        return
    }

    cantidadMinima, err := strconv.Atoi(r.FormValue("cantidad_minima"))
    if err != nil {
        http.Error(w, "Cantidad mínima inválida", http.StatusBadRequest)
        return
    }

    cantidadMaxima, err := strconv.Atoi(r.FormValue("cantidad_maxima"))
    if err != nil {
        http.Error(w, "Cantidad máxima inválida", http.StatusBadRequest)
        return
    }

    garantia, err := strconv.Atoi(r.FormValue("garantia"))
    if err != nil {
        http.Error(w, "Garantía inválida", http.StatusBadRequest)
        return
    }

    fechaEntrega, err := time.Parse("2006-01-02", r.FormValue("fecha_de_entrega"))
    if err != nil {
        http.Error(w, "Fecha de entrega inválida", http.StatusBadRequest)
        return
    }

    // Construcción del modelo
    partida := models.Partida{
        NumPartidaConvocatoria: numeroPartida,
        NombreDescripcion:      r.FormValue("nombre_descripcion"),
        Cantidad:               cantidad,
        CantidadMinima:         cantidadMinima,
        CantidadMaxima:         cantidadMaxima,
        NoFichaTecnica:         r.FormValue("no_ficha_tecnica"),
        TipoDeBien:             r.FormValue("tipo_de_bien"),
        ClaveCompendio:         r.FormValue("clave_compendio"),
        ClaveCucop:             r.FormValue("clave_cucop"),
        UnidadMedida:           r.FormValue("unidad_medida"),
        DiasDeEntrega:          r.FormValue("dias_de_entrega"),
        FechaDeEntrega:         fechaEntrega,
        Garantia:               garantia,
    }

    // Insertar partida
    idPartida, err := m.InsertarPartida(partida)
    if err != nil {
        fmt.Println("ERROR INSERTANDO PARTIDA:", err) // <--- AÑADE ESTO
        http.Error(w, "Error al insertar partida", http.StatusInternalServerError)
        return
    }

    // Insertar relación con licitación
    err = m.InsertarLicitacionPartida(idLicitacion, idPartida)
    if err != nil {
        http.Error(w, "Error al vincular partida con licitación", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/mostrar-partidas/%d", idLicitacion), http.StatusSeeOther)
}

func (m *Repository) ObtenerRequerimientos(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	req, err := m.ObtenerOCrearRequerimientos(id)
	if err != nil {
		log.Println("Error al obtener o crear requerimientos:", err)
		http.Error(w, "Error al obtener requerimientos", http.StatusInternalServerError)
		return
	}

	data := &models.TemplateData{
		Requerimientos: req,
		CSRFToken:      nosurf.Token(r),
	}

	render.RenderTemplate(w, "licitaciones/mostrar-partidas.page.tmpl", data)
}

func (m *Repository) ObtenerRequerimientosJSON(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	req, err := m.ObtenerOCrearRequerimientos(id)
	if err != nil {
		log.Println("Error al obtener o crear requerimientos:", err)
		http.Error(w, "Error al obtener requerimientos", http.StatusInternalServerError)
		return
	}

	// Devolver como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
	"requiere_mantenimiento": req.RequiereMantenimiento,
	"requiere_instalacion": req.RequiereInstalacion,
	"requiere_puesta_marcha": req.RequierePuestaEnMarcha,
	"requiere_capacitacion": req.RequiereCapacitacion,
	"requiere_visita_previa": req.RequiereVisitaPrevia,
	"fecha_visita": req.FechaVisita.Format("2006-01-02"),
	"comentarios_visita": req.ComentariosVisita,
	"requiere_muestra": req.RequiereMuestra,
	"fecha_muestra": req.FechaMuestra.Format("2006-01-02"),
	"comentarios_muestra": req.ComentariosMuestra,
    "fecha_entrega": req.FechaEntrega.Format("2006-01-02"),
    "comentarios_entrega": req.ComentariosEntrega,
    })

}

func (m *Repository) GuardarRequerimientos(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        log.Println("Error al parsear el formulario:", err)
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    idPartidaStr := r.FormValue("id_partida")
    idPartida, err := strconv.Atoi(idPartidaStr)
    if err != nil {
        log.Println("ID de partida inválido:", err)
        http.Error(w, "ID de partida inválido", http.StatusBadRequest)
        return
    }

    // Parsear checkboxes (están presentes solo si están marcados)
    mantenimiento := r.FormValue("requiere_mantenimiento") == "on"
    instalacion := r.FormValue("requiere_instalacion") == "on"
    puestaMarcha := r.FormValue("requiere_puesta_marcha") == "on"
    capacitacion := r.FormValue("requiere_capacitacion") == "on"
    visitaPrevia := r.FormValue("requiere_visita_previa") == "on"
    fechaVisita := r.FormValue("fecha_visita")
    comentariosVisita := r.FormValue("comentarios_visita")
    requiereMuestra := r.FormValue("requiere_muestra") == "on"
    fechaMuestra := r.FormValue("fecha_muestra")
    comentariosMuestra := r.FormValue("comentarios_muestra")
    fechaEntrega := r.FormValue("fecha_entrega")
    comentariosEntrega := r.FormValue("comentarios_entrega")


    // Asegurar que el registro existe (lo crea si no)
    _, err = m.ObtenerOCrearRequerimientos(idPartida)
    if err != nil {
        log.Println("Error al obtener o crear requerimientos:", err)
        http.Error(w, "Error interno al preparar requerimientos", http.StatusInternalServerError)
        return
    }

    // Actualizar requerimientos
    update := `
	UPDATE requerimientos_partida
	SET 
		requiere_mantenimiento = ?,
		requiere_instalacion = ?,
		requiere_puesta_marcha = ?,
		requiere_capacitacion = ?,
		requiere_visita_previa = ?,
		fecha_visita = ?,
		comentarios_visita = ?,
		requiere_muestra = ?,
		fecha_muestra = ?,
		comentarios_muestra = ?,
        fecha_entrega = ?,
        comentarios_entrega = ?,
		updated_at = NOW()
	WHERE id_partida = ?;
    `

    _, err = m.App.DB.Exec(update,
        mantenimiento,
        instalacion,
        puestaMarcha,
        capacitacion,
        visitaPrevia,
        fechaVisita,
        comentariosVisita,
        requiereMuestra,
        fechaMuestra,
        comentariosMuestra,
        fechaEntrega,
        comentariosEntrega,
        idPartida,
    )

    if err != nil {
        log.Println("Error al actualizar requerimientos:", err)
        http.Error(w, "Error al guardar los requerimientos", http.StatusInternalServerError)
        return
    }

    // Puedes redirigir o devolver 200 según cómo llames este handler (fetch o formulario clásico)
    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (m *Repository) MostrarAclaraciones(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener aclaracion
    aclaraciones, err := m.ObtenerAclaracionesPorPartidaID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    // Obtener partida
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    // Render
    data := &models.TemplateData{
        Aclaraciones: aclaraciones,
        Partida:      partida,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/aclaraciones.page.tmpl", data)
}

func (m *Repository) MostrarNuevaAclaracion(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener partida
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Empresas:   empresas,
        Partida:    partida,
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/nueva-aclaracion.page.tmpl", data)
}

func (m *Repository) CrearNuevaAclaracion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "No se pudo procesar el formulario", http.StatusBadRequest)
		return
	}

	idPartida, err := strconv.Atoi(r.FormValue("id_partida"))
	if err != nil {
		http.Error(w, "ID de partida inválido", http.StatusBadRequest)
		return
	}

	idEmpresa, err := strconv.Atoi(r.FormValue("id_empresa"))
	if err != nil {
		http.Error(w, "ID de empresa inválido", http.StatusBadRequest)
		return
	}

	fichaTecnicaID, _ := strconv.Atoi(r.FormValue("ficha_tecnica_id")) // puede venir vacío
	idPuntosTecnicosModif, _ := strconv.Atoi(r.FormValue("id_puntos_tecnicos_modif")) // también opcional

	aclaracion := models.AclaracionesPartida{
		Pregunta:        r.FormValue("pregunta"),
		Observaciones:   r.FormValue("observaciones"),
		FichaTecnica:    fichaTecnicaID,
		IDPuntosTecnico: idPuntosTecnicosModif,
		Partida:         &models.Partida{IDPartida: idPartida},
		Empresa:         &models.Empresas{IDEmpresa: idEmpresa},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = m.InsertarAclaracion(aclaracion)
	if err != nil {
		fmt.Println("ERROR al insertar aclaración:", err)
		http.Error(w, "Error al insertar la aclaración", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/aclaraciones/%d", idPartida), http.StatusSeeOther)
}

func (m *Repository) AgregarEmpresaExternaContexto(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    nombre := r.FormValue("nombre")
    idPartida := r.FormValue("id_partida")

    if nombre == "" || idPartida == "" {
        m.App.Session.Put(r.Context(), "error", "Todos los campos son obligatorios")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := m.AgregarEmpresaNueva(nombre)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al agregar la empresa")
        fmt.Println("Error al agregar empresa:", err)
        http.Redirect(w, r, "/nueva-aclaracion/"+idPartida, http.StatusSeeOther)
        return
    }

    m.App.Session.Put(r.Context(), "success", "Empresa agregada correctamente")
    http.Redirect(w, r, "/nueva-aclaracion/"+idPartida, http.StatusSeeOther)
}

func (m *Repository) MostrarProductosPartida(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }
    // Obtener partida
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }

    productosPartida, err := m.ObtenerProductosDePartida(idPartida)
    if err != nil {
        http.Error(w, "Error al obtener productos de la partida", http.StatusInternalServerError)
        return
    }


    data := &models.TemplateData{

        Partida:      partida,
        ProductosPartida: productosPartida,
        CSRFToken: nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/productos-partida.page.tmpl", data)
}

func (m *Repository) MostrarNuevoProductoPartida(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }
    // Obtener partida
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
        return
    }
    // Obtener todos los productos
    productos, err := m.ObtenerTodosProductos()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener productos: "+err.Error())
        log.Println("Error al obtener productos:", err)
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Productos:  productos,
        Partida:      partida,
        CSRFToken: nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/nuevo-producto-partida.page.tmpl", data)
}

func (m *Repository) CrearNuevoProductoPartida(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        m.App.Session.Put(r.Context(), "error", "Método no permitido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := r.ParseForm()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    idPartidaStr := r.Form.Get("id_partida")
    idPartida, err := strconv.Atoi(idPartidaStr)
    if err != nil || idPartida <= 0 {
        m.App.Session.Put(r.Context(), "error", "ID de partida inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Validar que la partida exista
    if !m.ExisteID("partidas", idPartida) {
        m.App.Session.Put(r.Context(), "error", "La partida no existe")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Procesar productos
    err = m.procesarProductosPartida(r, idPartida)
    if err != nil {
        log.Println("Error al guardar productos de la partida:", err)
        m.App.Session.Put(r.Context(), "error", "Error al guardar productos de la partida")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    m.App.Session.Put(r.Context(), "flash", "Productos guardados exitosamente")
    http.Redirect(w, r, fmt.Sprintf("/productos-partida/%d", idPartida), http.StatusSeeOther)
}

func (m *Repository) procesarProductosPartida(r *http.Request, idPartida int) error {
    productos := make([]struct {
        IDProducto     string
        PrecioOfertado string
        Observaciones  string
    }, 0)

    for key, values := range r.Form {
        if strings.HasPrefix(key, "productos[") && strings.Contains(key, "][id_producto]") {
            idx := strings.Split(key, "[")[1]
            idx = strings.Split(idx, "]")[0]

            producto := struct {
                IDProducto     string
                PrecioOfertado string
                Observaciones  string
            }{
                IDProducto:     values[0],
                PrecioOfertado: r.Form.Get(fmt.Sprintf("productos[%s][precio_ofertado]", idx)),
                Observaciones:  r.Form.Get(fmt.Sprintf("productos[%s][observaciones]", idx)),
            }

            if producto.IDProducto != "" {
                productos = append(productos, producto)
            }
        }
    }

    for _, p := range productos {
        idProducto, err := strconv.Atoi(p.IDProducto)
        if err != nil {
            continue
        }

        precio, err := strconv.ParseFloat(p.PrecioOfertado, 64)
        if err != nil || precio <= 0 {
            continue
        }

        if !m.ExisteID("productos", idProducto) {
            continue
        }

        _, err = m.App.DB.Exec(`
            INSERT INTO partida_productos (
                id_partida, id_producto, precio_ofertado,
                observaciones, created_at, updated_at
            ) VALUES (?, ?, ?, ?, ?, ?)`,
            idPartida,
            idProducto,
            precio,
            p.Observaciones,
            time.Now(),
            time.Now(),
        )
        if err != nil {
            return err
        }
    }

    return nil
}

func (m *Repository) EditarProductoPartida(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error al leer formulario", http.StatusBadRequest)
        return
    }

    idStr := r.Form.Get("id_partida_producto")
    precioStr := r.Form.Get("precio_ofertado")
    observaciones := r.Form.Get("observaciones")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    precio, err := strconv.ParseFloat(precioStr, 64)
    if err != nil {
        http.Error(w, "Precio inválido", http.StatusBadRequest)
        return
    }

    p := models.PartidaProductos{
        IDPartidaProducto: id,
        PrecioOfertado:    precio,
        Observaciones:     observaciones,
    }

    err = m.ActualizarProductoPartida(p)
    if err != nil {
        http.Error(w, "Error actualizando producto", http.StatusInternalServerError)
        return
    }
    idPartida, err := m.ObtenerIDPartidaPorIDPartidaProducto(p.IDPartidaProducto)
    if err != nil {
        http.Error(w, "No se pudo obtener la partida para redireccionar", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/productos-partida/"+strconv.Itoa(idPartida), http.StatusSeeOther)

}

func (m *Repository) EliminarProductoPartida(w http.ResponseWriter, r *http.Request) {
    // Obtener el ID de la relación partida_producto
    idParam := chi.URLParam(r, "id")
    idPartidaProducto, err := strconv.Atoi(idParam)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Primero obtenemos el id_partida para redireccionar después
    var idPartida int
    err = m.App.DB.QueryRow("SELECT id_partida FROM partida_productos WHERE id_partida_producto = ?", idPartidaProducto).Scan(&idPartida)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "No se pudo encontrar la relación partida-producto")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Eliminar el producto de la partida
    err = m.EliminarProductoDePartida(idPartidaProducto)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al eliminar el producto de la partida")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    m.App.Session.Put(r.Context(), "flash", "Producto eliminado de la partida exitosamente")
    http.Redirect(w, r, fmt.Sprintf("/productos-partida/%d", idPartida), http.StatusSeeOther)
}

func (m *Repository) MostrarPropuestas(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idPartida, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Obtener propuestas
	propuestas, err := m.ObtenerPropuestasPorPartidaID(idPartida)
	if err != nil {
		http.Error(w, "No se pudo obtener las propuestas", http.StatusInternalServerError)
		return
	}

	// Cargar el fallo en cada propuesta
	for i := range propuestas {
		fallo, err := m.ObtenerOCrearFalloPorPropuestaID(propuestas[i].IDPropuesta)
		if err != nil {
			// Puedes decidir asignar nil o un fallo vacío
			propuestas[i].Fallo = nil
			// o loguear el error
			log.Println("Error cargando fallo para propuesta", propuestas[i].IDPropuesta, ":", err)
		} else {
			propuestas[i].Fallo = fallo
		}
	}

	// Obtener partida
	partida, err := m.ObtenerPartidaPorID(idPartida)
	if err != nil {
		http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
		return
	}

	// Render
	data := &models.TemplateData{
		PropuestasPartida: propuestas,
		Partida:           partida,
		CSRFToken:         nosurf.Token(r),
	}

	render.RenderTemplate(w, "licitaciones/propuestas.page.tmpl", data)
}

func (m *Repository) MostrarNuevaPropuesta(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener la partida para mostrar información
    partida, err := m.ObtenerPartidaPorID(idPartida)
    if err != nil {
        http.Error(w, "No se pudo cargar la partida", http.StatusInternalServerError)
        return
    }

    marcas, err := m.ObtenerMarcas()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener marcas")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

    paises, err := m.ObtenerPaises()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener países")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }


    productos, _ := m.ObtenerTodosProductosExternos()

    data := &models.TemplateData{
        Partida:   partida,
        ProductosExternos: productos,
        CSRFToken: nosurf.Token(r),
        Marcas: marcas,
        Paises: paises,
        Empresas:   empresas,

    }

    render.RenderTemplate(w, "licitaciones/nueva-propuesta.page.tmpl", data)
}

func (m *Repository) CrearNuevaPropuesta(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPartida, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    err = r.ParseForm()
    if err != nil {
        http.Error(w, "No se pudo procesar el formulario", http.StatusBadRequest)
        return
    }


    propuesta := models.PropuestasPartida{
        IDPartida:         idPartida,
        IDEmpresa:         atoi(r.FormValue("id_empresa")),
        IDProductoExterno: atoi(r.FormValue("id_producto_externo")),
        PrecioOfertado:    atof(r.FormValue("precio_ofertado")),
        PrecioMin:         atof(r.FormValue("precio_min")),
        PrecioMax:         atof(r.FormValue("precio_max")),
        Observaciones:     r.FormValue("observaciones"),
    }

    err = m.InsertarPropuestaPartida(propuesta)
    if err != nil {
        http.Error(w, "Error al guardar la propuesta", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/propuestas/%d", idPartida), http.StatusSeeOther)
}

func (m *Repository) NuevoProductoExternoContexto(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parseando formulario", http.StatusInternalServerError)
        return
    }

    nombre := r.FormValue("nombre")
    modelo := r.FormValue("modelo")
    idMarca, _ := strconv.Atoi(r.FormValue("id_marca"))
    idPais, _ := strconv.Atoi(r.FormValue("id_pais_origen"))
    idEmpresa, _ := strconv.Atoi(r.FormValue("id_empresa_externa"))
    observaciones := r.FormValue("observaciones")
    idPartida := r.FormValue("id_partida") // Para redirigir al regresar

    _, err = m.App.DB.Exec(`
        INSERT INTO productos_externos 
        (nombre, modelo, id_marca, id_pais_origen, id_empresa_externa, observaciones, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
    `, nombre, modelo, idMarca, idPais, idEmpresa, observaciones)

    if err != nil {
        http.Error(w, "Error insertando producto externo", http.StatusInternalServerError)
        return
    }

    // Redirige al formulario de nueva propuesta con el contexto de la partida
    http.Redirect(w, r, "/nueva-propuesta/"+idPartida, http.StatusSeeOther)
}

func (m *Repository) MostrarEditarPropuesta(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    idPropuesta, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Obtener la propuesta con todos los datos relacionados
    propuesta, err := m.ObtenerPropuestaPorID(idPropuesta)
    if err != nil {
        http.Error(w, "No se pudo cargar la propuesta", http.StatusInternalServerError)
        return
    }

    // Obtener la partida asociada a la propuesta
    partida, err := m.ObtenerPartidaPorID(propuesta.IDPartida) // Asegúrate de tener esta función
    if err != nil {
        http.Error(w, "No se pudo cargar la partida", http.StatusInternalServerError)
        return
    }

    // Obtener catálogo para selects
    marcas, err := m.ObtenerMarcas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener marcas")
        http.Redirect(w, r, "/inventario", http.StatusSeeOther)
        return
    }

    paises, err := m.ObtenerPaises()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener países")
        http.Redirect(w, r, "/inventario", http.StatusSeeOther)
        return
    }

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo empresas")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    productos, _ := m.ObtenerTodosProductosExternos()

    data := &models.TemplateData{
        Propuesta: propuesta,  
        ProductosExternos: productos,
        Partida:   partida,
        CSRFToken: nosurf.Token(r),
        Marcas:    marcas,
        Paises:    paises,
        Empresas:  empresas,
    }

    render.RenderTemplate(w, "licitaciones/editar-propuesta.page.tmpl", data)
}

func (m *Repository) EditarPropuesta(w http.ResponseWriter, r *http.Request) {
    // Obtener el ID de la propuesta de la URL
    idParam := chi.URLParam(r, "id")
    idPropuesta, err := strconv.Atoi(idParam)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID de propuesta inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Parsear el formulario
    err = r.ParseForm()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Obtener los valores del formulario
    idProductoExterno, err := strconv.Atoi(r.Form.Get("id_producto_externo"))
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID de producto inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    idEmpresa, err := strconv.Atoi(r.Form.Get("id_empresa"))
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "ID de empresa inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    precioOfertado, err := strconv.ParseFloat(r.Form.Get("precio_ofertado"), 64)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Precio ofertado inválido")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Los campos precio_min y precio_max son opcionales
    var precioMin, precioMax float64
    if r.Form.Get("precio_min") != "" {
        precioMin, err = strconv.ParseFloat(r.Form.Get("precio_min"), 64)
        if err != nil {
            m.App.Session.Put(r.Context(), "error", "Precio mínimo inválido")
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
    }

    if r.Form.Get("precio_max") != "" {
        precioMax, err = strconv.ParseFloat(r.Form.Get("precio_max"), 64)
        if err != nil {
            m.App.Session.Put(r.Context(), "error", "Precio máximo inválido")
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
    }

    // Crear la estructura de propuesta actualizada
    propuestaActualizada := models.PropuestasPartida{
        IDPropuesta:       idPropuesta,
        IDProductoExterno: idProductoExterno,
        IDEmpresa:         idEmpresa,
        PrecioOfertado:    precioOfertado,
        PrecioMin:         precioMin,
        PrecioMax:         precioMax,
        Observaciones:     r.Form.Get("observaciones"),
    }

    // Actualizar en la base de datos
    err = m.ActualizarPropuesta(propuestaActualizada)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al actualizar la propuesta: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Obtener el ID de partida para redireccionar
    propuesta, err := m.ObtenerPropuestaPorID(idPropuesta)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener datos de la propuesta")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Redireccionar a la lista de propuestas de esta partida
    m.App.Session.Put(r.Context(), "flash", "Propuesta actualizada exitosamente")
    http.Redirect(w, r, fmt.Sprintf("/propuestas/%d", propuesta.IDPartida), http.StatusSeeOther)
}

func (m *Repository) ObtenerFalloPropuesta(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	fallo, err := m.ObtenerOCrearFalloPorPropuestaID(id)
	if err != nil {
		log.Println("Error al obtener o crear fallo:", err)
		http.Error(w, "Error al obtener fallo", http.StatusInternalServerError)
		return
	}

	data := &models.TemplateData{
		Fallo:     fallo,
		CSRFToken: nosurf.Token(r),
	}

	render.RenderTemplate(w, "licitaciones/fallo-propuesta.page.tmpl", data)
}

func (m *Repository) ObtenerFalloPropuestaJSON(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idPropuesta, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	fallo, err := m.ObtenerOCrearFallo(idPropuesta)
	if err != nil {
		log.Println("Error al obtener o crear fallo:", err)
		http.Error(w, "Error al obtener el fallo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cumple_legal": fallo.CumpleLegal,
		"cumple_administrativo": fallo.CumpleAdministrativo,
		"cumple_tecnico": fallo.CumpleTecnico,
		"puntos_obtenidos": fallo.PuntosObtenidos,
		"ganador": fallo.Ganador,
		"observaciones": fallo.Observaciones,
	})
}

func (m *Repository) GuardarFalloPropuesta(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error al parsear el formulario:", err)
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	idPropuestaStr := r.FormValue("id_propuesta")
	idPropuesta, err := strconv.Atoi(idPropuestaStr)
	if err != nil {
		log.Println("ID de propuesta inválido:", err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Checkboxes y otros valores
	cumpleLegal := r.FormValue("cumple_legal") == "on"
	cumpleAdmin := r.FormValue("cumple_administrativo") == "on"
	cumpleTecnico := r.FormValue("cumple_tecnico") == "on"
	puntosStr := r.FormValue("puntos_obtenidos")
    if puntosStr == "" {
        puntosStr = "0"
    }
    puntos, err := strconv.Atoi(puntosStr)
    if err != nil {
        log.Println("Puntos inválidos:", err)
        http.Error(w, "Puntos inválidos", http.StatusBadRequest)
        return
    }


	ganador := r.FormValue("ganador") == "on"
	observaciones := r.FormValue("observaciones")

	// Asegurar que el fallo exista (lo crea si no)
	_, err = m.ObtenerOCrearFallo(idPropuesta)
	if err != nil {
		log.Println("Error al obtener o crear fallo:", err)
		http.Error(w, "Error al preparar fallo", http.StatusInternalServerError)
		return
	}

	update := `
	UPDATE fallos_propuesta
	SET 
		cumple_legal = ?,
		cumple_administrativo = ?,
		cumple_tecnico = ?,
		puntos_obtenidos = ?,
		ganador = ?,
		observaciones = ?,
		updated_at = NOW()
	WHERE id_propuesta = ?;
	`

	_, err = m.App.DB.Exec(update,
		cumpleLegal,
		cumpleAdmin,
		cumpleTecnico,
		puntos,
		ganador,
		observaciones,
		idPropuesta,
	)

	if err != nil {
		log.Println("Error al guardar fallo:", err)
		http.Error(w, "Error al guardar fallo", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}





// TODO LO DE OPCIONES
func (m *Repository) Opciones(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "opciones/opciones.page.tmpl", &models.TemplateData{})
}

// DATOS REFERENCIA
func (m *Repository) DatosReferencia(w http.ResponseWriter, r *http.Request) {
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

    compañias, err := m.ObtenerCompañias()
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
        Compañias: compañias,
		CSRFToken:       nosurf.Token(r),
	}

    render.RenderTemplate(w, "opciones/datos-referencia.page.tmpl", data)
}

func (m *Repository) AgregarDato(w http.ResponseWriter, r *http.Request) {
    // Recibe la tabla seleccionada
    tabla := r.FormValue("tabla")
    var err error

    // Dependiendo de la tabla, realizamos diferentes operaciones
    switch tabla {
    case "marcas":
        nombre := r.FormValue("nombre")
        if nombre == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre de la marca es obligatorio")
            fmt.Println("Error: El nombre de la marca es obligatorio")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarMarca(nombre)        
    case "tipos_producto":
        nombre := r.FormValue("nombre")
        if nombre == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre del tipo de producto es obligatorio")
            fmt.Println("Error: El nombre del tipo de producto es obligatorio")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarTipoProducto(nombre)

    case "clasificaciones":
        nombre := r.FormValue("nombre")
        if nombre == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre de la clasificación es obligatorio")
            fmt.Println("Error: El nombre de la clasificación es obligatorio")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarClasificacion(nombre)

    case "paises":
        nombre := r.FormValue("nombre")
        codigo := r.FormValue("codigo")
        if nombre == "" || codigo == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre y el código del país son obligatorios")
            fmt.Println("Error: El nombre y el código del país son obligatorios")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarPais(nombre, codigo)

    case "certificaciones":
        nombre := r.FormValue("nombre")
        organismoEmisor := r.FormValue("organismo_emisor")
        if nombre == "" || organismoEmisor == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre y el organismo emisor son obligatorios")
            fmt.Println("Error: El nombre y el organismo emisor son obligatorios")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarCertificacion(nombre, organismoEmisor)

    case "compañias":
        nombre := r.FormValue("nombre")
        tipo := r.FormValue("tipo")
        if nombre == "" || tipo == "" {
            m.App.Session.Put(r.Context(), "error", "El nombre y el tipo son obligatorios")
            fmt.Println("Error: El nombre y el tipo son obligatorios")
            http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
            return
        }
        err = m.AgregarCompañia(nombre, tipo)

    default:
        m.App.Session.Put(r.Context(), "error", "Tabla no válida")
        fmt.Println("Error: Tabla no válida")
        http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
        return
    }

    // Verificamos si ocurrió algún error durante el proceso
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al agregar el dato")
        fmt.Println("Error al agregar el dato:", err)
        http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
        return
    }

    // Si todo salió bien
    m.App.Session.Put(r.Context(), "success", "Dato agregado correctamente")
    http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
}

func (m *Repository) EliminarDatoReferencia(w http.ResponseWriter, r *http.Request) {
    // Obtener el nombre de la tabla desde la URL
    tabla := r.FormValue("tabla")
    id := chi.URLParam(r, "id")

    var err error
    switch tabla {
    case "marcas":
        err = m.EliminarMarca(id)
    case "tipos_producto":
        err = m.EliminarTipoProducto(id)
    case "clasificaciones":
        err = m.EliminarClasificacion(id)
    case "paises":
        err = m.EliminarPais(id)
    case "certificaciones":
        err = m.EliminarCertificacion(id)
    case "compañias":
        err = m.EliminarCompañia(id)
    default:
        fmt.Println("Error: tabla no válida")
        m.App.Session.Put(r.Context(), "error", "Tabla no válida")
        http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
        return
    }

    if err != nil {
        fmt.Println("Error al eliminar el dato:", err)
        m.App.Session.Put(r.Context(), "error", "Error al eliminar el dato")
        http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
        return
    }
    m.App.Session.Put(r.Context(), "success", "Dato eliminado correctamente")
    http.Redirect(w, r, "/datos-referencia", http.StatusSeeOther)
}

// ENTIDADES
func (m *Repository) Entidades(w http.ResponseWriter, r *http.Request) {
    // Obtener todas las entidades con información de estado
    entidades, err := m.ObtenerTodasEntidades()
    if err != nil {
        log.Println("Error al obtener entidades:", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    // Obtener lista de estados para filtros (opcional)
    estados, err := m.ObtenerEstados()
    if err != nil {
        log.Println("Error al obtener estados:", err)
        // No es crítico, podemos continuar
    }

    data := &models.TemplateData{
        Data: map[string]interface{}{
            "Entidades": entidades,
            "Estados":   estados,
            
        },
        CSRFToken:       nosurf.Token(r),
    }

    render.RenderTemplate(w, "opciones/entidades-opciones.page.tmpl", data)
}

func (m *Repository) MostrarNuevaEntidad(w http.ResponseWriter, r *http.Request) {
    estados, err := m.ObtenerEstados()
    if err != nil {
        log.Println("Error al obtener estados:", err)
        // Manejar error según necesites
    }

    compañias, err := m.ObtenerCompañias()
    if err != nil {
        log.Println("Error al obtener compañías:", err)
    }

    data := &models.TemplateData{
        Data: map[string]interface{}{
            "Estados":   estados,
            "Compañias": compañias,
        },
        CSRFToken:       nosurf.Token(r),
    }
	render.RenderTemplate(w, "opciones/entidades-nueva.page.tmpl", data)
}

func (m *Repository) CrearEntidad(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Printf("ERROR ParseForm: %v\n", err)
        log.Println("Error al parsear formulario:", err)
        m.App.Session.Put(r.Context(), "error", "Datos inválidos")
        http.Redirect(w, r, "/crear-entidad", http.StatusSeeOther)
        return
    }
    // Validar campos requeridos
    required := map[string]string{
        "nombre":      "Nombre",
        "tipo":        "Tipo",
        "id_compañia": "Compañía",
        "estado":      "Estado",
    }
    
    for field, name := range required {
        val := r.Form.Get(field)
        if val == "" {
            fmt.Printf("FALTA Campo requerido: %s\n", name)
            m.App.Session.Put(r.Context(), "error", name+" es requerido")
            http.Redirect(w, r, "/crear-entidad", http.StatusSeeOther)
            return
        }
    }

    // Convertir id_compañia
    idCompañiaStr := r.Form.Get("id_compañia")    
    idCompañia, err := strconv.Atoi(idCompañiaStr)
    if err != nil {
        fmt.Printf("ERROR Conversión id_compañia: %v\n", err)
        m.App.Session.Put(r.Context(), "error", "ID de compañía inválido")
        http.Redirect(w, r, "/crear-entidad", http.StatusSeeOther)
        return
    }
    // Validar estado
    estado := r.Form.Get("estado")
    if len(estado) != 2 {
        fmt.Println("ERROR: Longitud de estado incorrecta")
        m.App.Session.Put(r.Context(), "error", "El estado debe tener 2 caracteres (ej: '09')")
        http.Redirect(w, r, "/crear-entidad", http.StatusSeeOther)
        return
    }
    // Construir estructura
    entidad := models.Entidad{
        Nombre:       r.Form.Get("nombre"),
        Compañia:     models.Compañias{IDCompañia: idCompañia},
        Estado:       models.EstadosRepublica{ClaveEstado: estado},
        Municipio:    r.Form.Get("municipio"),
        CodigoPostal: r.Form.Get("codigo_postal"),
        Direccion:    r.Form.Get("direccion"),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    // Insertar en BD
    if err := m.InsertarEntidad(entidad); err != nil {
        fmt.Printf("ERROR InsertarEntidad: %v\n", err)
        log.Println("Error al guardar entidad:", err)
        m.App.Session.Put(r.Context(), "error", "Error al guardar. ¿La compañía existe? Detalle: "+err.Error())
        http.Redirect(w, r, "/crear-entidad", http.StatusSeeOther)
        return
    }
    m.App.Session.Put(r.Context(), "flash", "¡Entidad creada exitosamente!")
    http.Redirect(w, r, "/datos-entidades", http.StatusSeeOther)
}

// EMPRESAS EXTERNAS Y PRODUCTOS EXTERNOS
func (m *Repository) EmpresasExternas(w http.ResponseWriter, r *http.Request) {
    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        Empresas:   empresas,
        CSRFToken:       nosurf.Token(r),
    }

    render.RenderTemplate(w, "opciones/empresas-externas.page.tmpl", data)
}

func (m *Repository) AgregarEmpresaExterna(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
        http.Redirect(w, r, "/datos-empresas-externas", http.StatusSeeOther)
        return
    }

    nombre := r.FormValue("nombre")
    if nombre == "" {
        m.App.Session.Put(r.Context(), "error", "El nombre de la empresa es obligatorio")
        http.Redirect(w, r, "/datos-empresas-externas", http.StatusSeeOther)
        return
    }

    err := m.AgregarEmpresaNueva(nombre)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al agregar la empresa")
        fmt.Println("Error al agregar empresa:", err)
        http.Redirect(w, r, "/datos-empresas-externas", http.StatusSeeOther)
        return
    }

    m.App.Session.Put(r.Context(), "success", "Empresa agregada correctamente")
    http.Redirect(w, r, "/datos-empresas-externas", http.StatusSeeOther)
}

func (m *Repository) ProductosExternos(w http.ResponseWriter, r *http.Request) {
    marcas, err := m.ObtenerMarcas()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener marcas")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

    paises, err := m.ObtenerPaises()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener países")
		http.Redirect(w, r, "/inventario", http.StatusSeeOther)
		return
	}

    empresas, err := m.ObtenerTodasEmpresas()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo empresas: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }


    productos, err := m.ObtenerTodosProductosExternos()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error obteniendo productos externos: "+err.Error())
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := &models.TemplateData{
        ProductosExternos: productos,
        Marcas: marcas,
        CSRFToken: nosurf.Token(r),
        Paises: paises,
        Empresas:   empresas,

    }

    render.RenderTemplate(w, "opciones/productos-externos.page.tmpl", data)
}

func (m *Repository) NuevoProductoExternoContextoMenu(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parseando formulario", http.StatusInternalServerError)
        return
    }

    nombre := r.FormValue("nombre")
    modelo := r.FormValue("modelo")
    idMarca, _ := strconv.Atoi(r.FormValue("id_marca"))
    idPais, _ := strconv.Atoi(r.FormValue("id_pais_origen"))
    idEmpresa, _ := strconv.Atoi(r.FormValue("id_empresa_externa"))
    observaciones := r.FormValue("observaciones")

    _, err = m.App.DB.Exec(`
        INSERT INTO productos_externos 
        (nombre, modelo, id_marca, id_pais_origen, id_empresa_externa, observaciones, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
    `, nombre, modelo, idMarca, idPais, idEmpresa, observaciones)

    if err != nil {
        http.Error(w, "Error insertando producto externo", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/productos-externos", http.StatusSeeOther)
}