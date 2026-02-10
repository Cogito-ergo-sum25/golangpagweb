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
	"golang.org/x/crypto/bcrypt"
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

// AUTENTICACIÓN
// En tu archivo de handlers (ej. pkg/handlers/handlers.go)

// ShowLoginPage renderiza la página de inicio de sesión.
func (m *Repository) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
    data := &models.TemplateData{
        CSRFToken: nosurf.Token(r),
    }
    render.RenderTemplate(w, "auth/login.page.tmpl", data)
}

// PostLoginPage procesa el envío del formulario de inicio de sesión.
func (m *Repository) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		// ... manejar error ...
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	
    // AQUÍ ESTÁ LA MAGIA: Llamamos a nuestra nueva función helper.
	id, hashedPassword, err := m.Authenticate(email)
	if err != nil {
        // Si hay un error (ej. "usuario no encontrado"), redirigimos.
		m.App.Session.Put(r.Context(), "error", "Credenciales inválidas.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

    // El resto de la lógica es la misma: comparar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Credenciales inválidas.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

    // Si todo es correcto, guardamos en la sesión y redirigimos.
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "¡Bienvenido!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout destruye la sesión del usuario.
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	
	m.App.Session.Put(r.Context(), "flash", "Has cerrado sesión correctamente.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home/home.page.tmpl", &models.TemplateData{})
}

// TODO LO DE INVENTARIO

// PANTALLA INICIO INVENTARIO
func (m *Repository) Inventario(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener parámetros de la URL
    marca := r.URL.Query().Get("marca")
    clasificacion := r.URL.Query().Get("clasificacion")
    busqueda := r.URL.Query().Get("busqueda")

    // 2. Obtener todos los productos (reutilizando tu lógica de consulta)
    // Aquí puedes llamar a una función que haga el SELECT con los JOINs
    productos, err := m.ObtenerProductosParaInventario() 
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al obtener productos")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // 3. Aplicar filtros en memoria (Tu lógica de Catalogo)
    var filteredProductos []models.Producto
    for _, p := range productos {
        if marca != "" && p.Marca != marca {
            continue
        }
        if clasificacion != "" && p.Clasificacion != clasificacion {
            continue
        }
        // Buscamos en Nombre o SKU
        if busqueda != "" {
            term := strings.ToLower(busqueda)
            if !strings.Contains(strings.ToLower(p.Nombre), term) && 
               !strings.Contains(strings.ToLower(p.SKU), term) {
                continue
            }
        }
        filteredProductos = append(filteredProductos, p)
    }

    // 4. Obtener listas para los selectores (reutilizando tu función obtenerDatosUnicos)
    marcas, _ := m.obtenerDatosUnicos("SELECT DISTINCT nombre FROM marcas WHERE nombre != ''")
    clasificaciones, _ := m.obtenerDatosUnicos("SELECT DISTINCT nombre FROM clasificaciones WHERE nombre != ''")

    // 5. Preparar datos para la plantilla
    data := &models.TemplateData{
        Productos: filteredProductos,
        Data: map[string]interface{}{
            "Marcas":         marcas,
            "Clasificaciones": clasificaciones,
            "Filtros": map[string]interface{}{
                "Marca":         marca,
                "Clasificacion": clasificacion,
                "Busqueda":      busqueda,
            },
        },
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
	// 1. Llamamos al helper para obtener los proyectos con el estatus de su licitación.
	todosLosProyectos, err := m.ObtenerProyectosParaVista()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al obtener proyectos: "+err.Error())
		log.Println("Error al obtener proyectos:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 2. Creamos dos listas vacías para la clasificación.
	var proyectosVigentes []models.Proyecto
	var proyectosArchivados []models.Proyecto

	// 3. Iteramos y clasificamos basándonos en el texto del estatus.
	for _, p := range todosLosProyectos {
		// Convertimos el estatus a minúsculas para una comparación segura.
		// Si el estatus es "vigente", va a la primera lista.
		if strings.ToLower(p.Estatus) == "vigente" {
			proyectosVigentes = append(proyectosVigentes, p)
		} else {
			// Cualquier otro estatus ("Cancelado", "Adjudicado", "Desierto", etc.) se considera archivado.
			proyectosArchivados = append(proyectosArchivados, p)
		}
	}

	// 4. Pasamos ambas listas a la plantilla.
	data := make(map[string]interface{})
	data["ProyectosVigentes"] = proyectosVigentes
	data["ProyectosArchivados"] = proyectosArchivados

    render.RenderTemplate(w, "proyectos/proyectos-vista.page.tmpl", &models.TemplateData{
        Data:      data,
        CSRFToken: nosurf.Token(r),
    })
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
	// ... tus validaciones de formulario (que están perfectas) ...
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario")
		http.Redirect(w, r, "/nuevo-proyecto", http.StatusSeeOther)
		return
	}

	requiredFields := []string{"nombre", "descripcion", "id_licitacion", "fecha_inicio"}
	for _, field := range requiredFields {
		if r.Form.Get(field) == "" {
			m.App.Session.Put(r.Context(), "error", "El campo "+field+" es requerido")
			http.Redirect(w, r, "/nuevo-proyecto", http.StatusSeeOther)
			return
		}
	}

	idLicitacion, _ := strconv.Atoi(r.Form.Get("id_licitacion"))
	fechaInicio, _ := time.Parse("2006-01-02", r.Form.Get("fecha_inicio"))

	proyecto := models.Proyecto{
		Nombre:       r.Form.Get("nombre"),
		Descripcion:  r.Form.Get("descripcion"),
		IDLicitacion: idLicitacion,
		FechaInicio:  fechaInicio,
	}

	if fechaFinStr := r.Form.Get("fecha_fin"); fechaFinStr != "" {
		proyecto.FechaFin, err = time.Parse("2006-01-02", fechaFinStr)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "Formato de fecha fin inválido")
			http.Redirect(w, r, "/nuevo-proyecto", http.StatusSeeOther)
			return
		}
	}

	// Llamamos al helper para que realice la inserción
	_, err = m.CrearProyectoEnDB(proyecto)
	if err != nil {
		log.Println("Error al llamar a CrearProyectoEnDB:", err)
		m.App.Session.Put(r.Context(), "error", "Error interno al crear el proyecto")
		http.Redirect(w, r, "/nuevo-proyecto", http.StatusSeeOther)
		return
	}

	// ¡CORREGIDO! Redirigimos a la vista general de proyectos.
	m.App.Session.Put(r.Context(), "flash", "¡Proyecto creado exitosamente!")
	http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
}

func (m *Repository) MostrarFormularioEditarProyecto(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el ID del proyecto desde la URL (ej: /proyectos/editar/123)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "ID de proyecto inválido.")
		http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
		return
	}

	// 1. Obtenemos los datos del proyecto específico que se va a editar.
	proyecto, err := m.ObtenerProyectoPorID(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "No se encontró el proyecto.")
		http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
		return
	}

	// 2. Obtenemos la lista de todas las licitaciones para llenar el menú desplegable.
	licitaciones, err := m.ObtenerLicitacionesParaSelect()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error obteniendo licitaciones.")
		http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
		return
	}

	// 3. Preparamos los datos y los pasamos a la plantilla HTML.
	data := make(map[string]interface{})
	data["Proyecto"] = proyecto
	data["Licitaciones"] = licitaciones

    render.RenderTemplate(w, "proyectos/editar-proyecto.page.tmpl", &models.TemplateData{
        Data:      data,
        CSRFToken: nosurf.Token(r),
    })
}

// ProcesarFormularioEditarProyecto (POST) procesa la actualización de un proyecto.
func (m *Repository) ProcesarFormularioEditarProyecto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "ID de proyecto inválido.")
		http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al procesar el formulario.")
		http.Redirect(w, r, fmt.Sprintf("/proyectos/editar/%d", id), http.StatusSeeOther)
		return
	}

	// Creamos un struct 'Proyecto' con los datos actualizados del formulario.
	idLicitacion, _ := strconv.Atoi(r.Form.Get("id_licitacion"))
	fechaInicio, _ := time.Parse("2006-01-02", r.Form.Get("fecha_inicio"))

	proyecto := models.Proyecto{
		IDProyecto:   id,
		Nombre:       r.Form.Get("nombre"),
		Descripcion:  r.Form.Get("descripcion"),
		IDLicitacion: idLicitacion,
		FechaInicio:  fechaInicio,
	}

	// Manejamos la fecha de fin opcional.
	if fechaFinStr := r.Form.Get("fecha_fin"); fechaFinStr != "" {
		proyecto.FechaFin, err = time.Parse("2006-01-02", fechaFinStr)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "Formato de fecha fin inválido.")
			http.Redirect(w, r, fmt.Sprintf("/proyectos/editar/%d", id), http.StatusSeeOther)
			return
		}
	}

	// Llamamos al helper para que guarde los cambios en la base de datos.
	if err := m.ActualizarProyectoEnDB(proyecto); err != nil {
		log.Println("Error al actualizar proyecto:", err)
		m.App.Session.Put(r.Context(), "error", "Error interno al actualizar el proyecto.")
		http.Redirect(w, r, fmt.Sprintf("/proyectos/editar/%d", id), http.StatusSeeOther)
		return
	}

	// Si todo sale bien, redirigimos a la lista principal con un mensaje de éxito.
	m.App.Session.Put(r.Context(), "flash", "¡Proyecto actualizado exitosamente!")
	http.Redirect(w, r, "/proyectos-vista", http.StatusSeeOther)
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

func (m *Repository) GetArchivosLicitacion(w http.ResponseWriter, r *http.Request) {
    // 1. Identificar de qué licitación hablamos
    idStr := chi.URLParam(r, "id")
    idLicitacion, _ := strconv.Atoi(idStr)

    // 2. Llamar al Helper 1: Datos de la Licitación
    lic, err := m.ObtenerLicitacionPorID(idLicitacion)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "No pudimos encontrar esa licitación")
        http.Redirect(w, r, "/licitaciones", http.StatusSeeOther)
        return
    }

    // 3. Llamar al Helper 2: Lista de archivos
    archivos, err := m.ObtenerArchivosLicitacion(idLicitacion)
    if err != nil {
        // Si no hay archivos o falla, mandamos una lista vacía para que el template no sufra
        archivos = []models.ArchivoLicitacion{}
    }

    // 4. Cargar la vista con todo listo
    data := &models.TemplateData{
        Licitacion: lic,
        Archivos:   archivos, // El campo que acabamos de agregar al struct
        CSRFToken:  nosurf.Token(r),
    }

    render.RenderTemplate(w, "licitaciones/archivos.page.tmpl", data)
}

func (m *Repository) PostGuardarEnlace(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()

    idLicitacion, _ := strconv.Atoi(r.FormValue("id_licitacion"))
    nombre := r.FormValue("nombre_archivo")
    link := r.FormValue("link_servidor")
    tipo := r.FormValue("tipo_archivo") // Capturamos el tipo del form
    comentarios := r.FormValue("comentarios")

    query := `INSERT INTO archivos_licitacion (id_licitacion, nombre_archivo, link_servidor, tipo_archivo, comentarios)
              VALUES (?, ?, ?, ?, ?)`

    _, err := m.App.DB.Exec(query, idLicitacion, nombre, link, tipo, comentarios)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Error al guardar: " + err.Error())
    }

    http.Redirect(w, r, fmt.Sprintf("/archivos-licitacion/%d", idLicitacion), http.StatusSeeOther)
}

func (m *Repository) PostEliminarEnlace(w http.ResponseWriter, r *http.Request) {
    _ = r.ParseForm()
    idArchivo, _ := strconv.Atoi(r.FormValue("id_archivo"))
    idLicitacion, _ := strconv.Atoi(r.FormValue("id_licitacion"))

    query := `DELETE FROM archivos_licitacion WHERE id_archivo = ?`
    _, err := m.App.DB.Exec(query, idArchivo)

    if err != nil {
        m.App.Session.Put(r.Context(), "error", "No se pudo eliminar")
    }

    http.Redirect(w, r, fmt.Sprintf("/archivos-licitacion/%d", idLicitacion), http.StatusSeeOther)
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
    
    aclaraciones, err := m.ObtenerAclaracionesPorLicitacionID(id)
    if err != nil {
        log.Printf("Error al obtener aclaraciones de la licitación %d: %v", id, err)
        http.Error(w, "Error al obtener aclaraciones", http.StatusInternalServerError)
        return
    }

    // Agrupar aclaraciones por ID de partida.
    // Usamos tu tipo de modelo exacto: models.AclaracionesLicitacion
    aclaracionesPorPartida := make(map[int][]models.AclaracionesLicitacion) // <-- CAMBIO AQUÍ
    for _, a := range aclaraciones {
        aclaracionesPorPartida[a.IDPartida] = append(aclaracionesPorPartida[a.IDPartida], a)
    }

    // Crear la lista de empresas únicas a partir de las aclaraciones obtenidas.
    // Usamos tu tipo de modelo exacto: models.Empresas
    empresasMap := make(map[int]models.Empresas) // <-- CAMBIO AQUÍ
    for _, a := range aclaraciones {
        // Asumo que tu struct Empresas tiene un campo IDEmpresa
        if a.Empresa != nil {
             empresasMap[a.Empresa.IDEmpresa] = *a.Empresa
        }
    }
    var empresasUnicas []models.Empresas // <-- CAMBIO AQUÍ
    for _, empresa := range empresasMap {
        empresasUnicas = append(empresasUnicas, empresa)
    }

    // Preparamos los datos para el template
    data := make(map[string]interface{})
    data["Licitacion"] = licitacion
    data["Partidas"] = partidas
    data["AclaracionesPorPartida"] = aclaracionesPorPartida
    data["EmpresasUnicas"] = empresasUnicas

    render.RenderTemplate(w, "licitaciones/aclaraciones-licitacion.page.tmpl", &models.TemplateData{
        Data:      data,
        CSRFToken: nosurf.Token(r),
    })
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
	fichaID := r.FormValue("ficha_tecnica_id")
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

func (m *Repository) MostrarFormularioEditarAclaracion(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener el ID de la aclaración desde la URL
    idParam := chi.URLParam(r, "id")
    idAclaracion, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID de aclaración inválido", http.StatusBadRequest)
        return
    }

    // 2. Obtener la aclaración específica por su ID (necesitarás crear esta función)
    aclaracion, err := m.ObtenerAclaracionPorID(idAclaracion)
    if err != nil {
        log.Println("Error al obtener la aclaración:", err)
        http.Error(w, "No se pudo encontrar la aclaración", http.StatusNotFound)
        return
    }

    // 3. Obtener datos adicionales para los dropdowns (partidas, empresas)
    licitacion, _ := m.ObtenerLicitacionPorID(aclaracion.IDLicitacion)
    partidas, _ := m.ObtenerPartidasPorLicitacionID(aclaracion.IDLicitacion)
    empresas, _ := m.ObtenerTodasEmpresas()

    // 4. Empaquetar todo y renderizar el nuevo template de edición
    data := make(map[string]interface{})
    data["Aclaracion"] = aclaracion
    data["Licitacion"] = licitacion
    data["Partidas"] = partidas
    data["Empresas"] = empresas

    render.RenderTemplate(w, "licitaciones/editar-aclaracion.page.tmpl", &models.TemplateData{
        Data: data,
        CSRFToken: nosurf.Token(r),
    })
}

func (m *Repository) ProcesarFormularioEditarAclaracion(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener el ID de la aclaración desde la URL
    idParam := chi.URLParam(r, "id")
    idAclaracion, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "ID de aclaración inválido", http.StatusBadRequest)
        return
    }

    // 2. Procesar el formulario
    if err := r.ParseForm(); err != nil {
        http.Error(w, "No se pudo procesar el formulario", http.StatusBadRequest)
        return
    }

    // 3. Crear un objeto AclaracionesLicitacion con los datos actualizados
    // Nota: El id_licitacion lo tomamos de un campo oculto para la redirección
    idLicitacion, _ := strconv.Atoi(r.FormValue("id_licitacion"))
    idPartida, _ := strconv.Atoi(r.FormValue("id_partida"))

    aclaracionActualizada := models.AclaracionesLicitacion{
        IDAclaracionLicitacion: idAclaracion, // El ID de la aclaración a actualizar
        IDLicitacion:           idLicitacion,
        IDPartida:              idPartida,
        IDEmpresa:              toInt(r.FormValue("id_empresa")),
        Pregunta:               r.FormValue("pregunta"),
        Observaciones:          r.FormValue("observaciones"),
        FichaTecnicaID:         r.FormValue("ficha_tecnica_id"),
        IDPuntosTecnicosModif:  toInt(r.FormValue("id_puntos_tecnicos_modif")),
    }

    // 4. Llamar a la función de la base de datos para actualizar (necesitarás crearla)
    err = m.ActualizarAclaracion(aclaracionActualizada)
    if err != nil {
        log.Println("Error al actualizar la aclaración:", err)
        http.Error(w, "Error al guardar los cambios", http.StatusInternalServerError)
        return
    }

    // 5. Redirigir de vuelta a la página de listado
    m.App.Session.Put(r.Context(), "flash", "Aclaración actualizada correctamente")
    http.Redirect(w, r, fmt.Sprintf("/aclaraciones-licitacion/%d", idLicitacion), http.StatusSeeOther)
}

// Pequeña función helper para convertir strings a int, por si acaso
func toInt(s string) int {
    i, _ := strconv.Atoi(s)
    return i
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
		http.Error(w, "ID de partida inválido", http.StatusBadRequest)
		return
	}

	// Usamos la función de la capa de DB que solo lee y maneja NULLs
	req, err := m.ObtenerRequerimientosPorPartidaID(id)
	if err != nil {
		log.Println("Error al obtener requerimientos desde la BD:", err)
		http.Error(w, "Error interno al obtener requerimientos", http.StatusInternalServerError)
		return
	}

	// Función auxiliar para formatear fechas o devolver un string vacío si la fecha es "cero"
	formatDate := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02")
	}

	// Construimos el mapa que se convertirá en JSON
	responseMap := map[string]interface{}{
		"requiere_mantenimiento": req.RequiereMantenimiento,
		"requiere_instalacion":   req.RequiereInstalacion,
		"requiere_puesta_marcha": req.RequierePuestaEnMarcha,
		"requiere_capacitacion":  req.RequiereCapacitacion,
		"requiere_visita_previa": req.RequiereVisitaPrevia,
		"fecha_visita":           formatDate(req.FechaVisita),
		"comentarios_visita":     req.ComentariosVisita,
		"requiere_muestra":       req.RequiereMuestra,
		"fecha_muestra":          formatDate(req.FechaMuestra),
		"comentarios_muestra":    req.ComentariosMuestra,
		"fecha_entrega":          formatDate(req.FechaEntrega),
		"comentarios_entrega":    req.ComentariosEntrega,
	}

	// Enviamos la respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseMap); err != nil {
		log.Println("Error al codificar respuesta JSON:", err)
	}
}

func (m *Repository) GuardarRequerimientos(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error al parsear el formulario:", err)
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	idPartida, err := strconv.Atoi(r.FormValue("id_partida"))
	if err != nil {
		log.Println("ID de partida inválido:", err)
		http.Error(w, "ID de partida inválido", http.StatusBadRequest)
		return
	}

	// Parsear checkboxes (si no están marcados, no vienen en el form, por eso `== "on"`)
	mantenimiento := r.FormValue("requiere_mantenimiento") == "on"
	instalacion := r.FormValue("requiere_instalacion") == "on"
	puestaMarcha := r.FormValue("requiere_puesta_marcha") == "on"
	capacitacion := r.FormValue("requiere_capacitacion") == "on"
	visitaPrevia := r.FormValue("requiere_visita_previa") == "on"
	requiereMuestra := r.FormValue("requiere_muestra") == "on"

	// Parsear campos de texto
	comentariosVisita := r.FormValue("comentarios_visita")
	comentariosMuestra := r.FormValue("comentarios_muestra")
	comentariosEntrega := r.FormValue("comentarios_entrega")

	// Manejo especial para fechas: si el string está vacío, usamos 'nil' para que se guarde como NULL
	var fechaVisita, fechaMuestra, fechaEntrega interface{}
	if val := r.FormValue("fecha_visita"); val != "" {
		fechaVisita = val
	} else {
		fechaVisita = nil
	}
	if val := r.FormValue("fecha_muestra"); val != "" {
		fechaMuestra = val
	} else {
		fechaMuestra = nil
	}
	if val := r.FormValue("fecha_entrega"); val != "" {
		fechaEntrega = val
	} else {
		fechaEntrega = nil
	}

	// Query de UPSERT: Inserta una nueva fila o actualiza la existente si la 'id_partida' ya existe
	// (Esto funciona gracias a la restricción UNIQUE en la columna 'id_partida')
	query := `
		INSERT INTO requerimientos_partida (
			id_partida, requiere_mantenimiento, requiere_instalacion, requiere_puesta_marcha, 
			requiere_capacitacion, requiere_visita_previa, fecha_visita, comentarios_visita, 
			requiere_muestra, fecha_muestra, comentarios_muestra, fecha_entrega, comentarios_entrega,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			requiere_mantenimiento = VALUES(requiere_mantenimiento),
			requiere_instalacion = VALUES(requiere_instalacion),
			requiere_puesta_marcha = VALUES(requiere_puesta_marcha),
			requiere_capacitacion = VALUES(requiere_capacitacion),
			requiere_visita_previa = VALUES(requiere_visita_previa),
			fecha_visita = VALUES(fecha_visita),
			comentarios_visita = VALUES(comentarios_visita),
			requiere_muestra = VALUES(requiere_muestra),
			fecha_muestra = VALUES(fecha_muestra),
			comentarios_muestra = VALUES(comentarios_muestra),
			fecha_entrega = VALUES(fecha_entrega),
			comentarios_entrega = VALUES(comentarios_entrega),
			updated_at = NOW();
	`

	_, err = m.App.DB.Exec(query,
		idPartida,
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
	)

	if err != nil {
		log.Println("Error al guardar (UPSERT) requerimientos:", err)
		http.Error(w, "Error al guardar los requerimientos", http.StatusInternalServerError)
		return
	}

	// Agregamos un mensaje flash para notificar al usuario y lo redirigimos a la página anterior
	m.App.Session.Put(r.Context(), "flash", "Requerimientos guardados con éxito.")
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
	// 1. Obtener el ID de la partida desde la URL
	idParam := chi.URLParam(r, "id")
	idPartida, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID de partida inválido", http.StatusBadRequest)
		return
	}

	// 2. Obtener todas las propuestas para esa partida
	propuestas, err := m.ObtenerPropuestasPorPartidaID(idPartida)
	if err != nil {
		log.Println("Error al obtener propuestas:", err)
		http.Error(w, "No se pudo obtener las propuestas", http.StatusInternalServerError)
		return
	}

	// 3. Cargar el fallo para cada propuesta, SÓLO SI EXISTE
	for i := range propuestas {
		// Usamos la nueva función que solo lee y no crea
		fallo, err := m.ObtenerFalloPorPropuestaID(propuestas[i].IDPropuesta)
		if err != nil {
			// Este error solo debe ocurrir si hay un problema con la base de datos
			log.Printf("Error al cargar fallo para propuesta %d: %v", propuestas[i].IDPropuesta, err)
			propuestas[i].Fallo = nil // Aseguramos que el fallo sea nil en caso de error
		} else {
			// Asignamos el resultado. Será 'nil' si no se encontró, o el puntero al fallo si se encontró.
			propuestas[i].Fallo = fallo
		}
	}

	// 4. Obtener la información de la partida para el título y el botón "Volver"
	partida, err := m.ObtenerPartidaPorID(idPartida)
	if err != nil {
		log.Println("Error al obtener la partida:", err)
		http.Error(w, "No se pudo obtener la partida", http.StatusInternalServerError)
		return
	}

	// 5. Preparar los datos y renderizar la plantilla
	templateData := &models.TemplateData{
        PropuestasPartida: propuestas, // Directamente aquí
        Partida:           partida,      // Directamente aquí
        CSRFToken:         nosurf.Token(r),
    }


    render.RenderTemplate(w, "licitaciones/propuestas.page.tmpl", templateData)
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

func (m *Repository) BuscarProductosExternosJSON(w http.ResponseWriter, r *http.Request) {
    // Obtenemos el término de búsqueda de la URL (?q=...)
    query := r.URL.Query().Get("q")

    // Si la búsqueda es muy corta, no devolvemos nada
    if len(query) < 2 {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("[]")) // Devolvemos un array JSON vacío
        return
    }

    // Llamamos a la nueva función de la base de datos
    productos, err := m.BuscarProductosExternosPorNombre(query)
    if err != nil {
        log.Println("Error buscando productos:", err)
        http.Error(w, "Error interno", http.StatusInternalServerError)
        return
    }

    // Devolvemos los resultados como JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(productos); err != nil {
        log.Println("Error codificando JSON de productos:", err)
    }
}

func (m *Repository) CrearMarcaJSON(w http.ResponseWriter, r *http.Request) {
	// 1. Parsear el formulario enviado por fetch
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// 2. Obtener y validar el nombre
	nombre := r.FormValue("nombre")
	if nombre == "" {
		http.Error(w, "El nombre no puede estar vacío", http.StatusBadRequest)
		return
	}

	// 3. Llamar a la función de DB que inserta y devuelve el objeto completo
	nuevaMarca, err := m.InsertarMarcaYDevolver(nombre)
	if err != nil {
		log.Println("Error al insertar marca:", err)
		http.Error(w, "Error al guardar la marca", http.StatusInternalServerError)
		return
	}

	// 4. Devolver el nuevo objeto como JSON con un status 201 (Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(nuevaMarca); err != nil {
		log.Println("Error al codificar JSON de nueva marca:", err)
	}
}

func (m *Repository) CrearEmpresaExternaJSON(w http.ResponseWriter, r *http.Request) {
	// 1. Parsear el formulario enviado por fetch
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// 2. Obtener y validar el nombre
	nombre := r.FormValue("nombre")
	if nombre == "" {
		http.Error(w, "El nombre no puede estar vacío", http.StatusBadRequest)
		return
	}

	// 3. Llamar a la función de DB que inserta y devuelve el objeto completo
	nuevaEmpresa, err := m.InsertarEmpresaExternaYDevolver(nombre)
	if err != nil {
		log.Println("Error al insertar empresa externa:", err)
		http.Error(w, "Error al guardar la empresa", http.StatusInternalServerError)
		return
	}

	// 4. Devolver el nuevo objeto como JSON con un status 201 (Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(nuevaEmpresa); err != nil {
		log.Println("Error al codificar JSON de nueva empresa:", err)
	}
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

    // Usamos la nueva función que SOLO lee
    fallo, err := m.ObtenerFalloPorPropuestaID(id)
    if err != nil {
        log.Println("Error al obtener fallo:", err)
        http.Error(w, "Error al obtener fallo", http.StatusInternalServerError)
        return
    }

    // Si no existe, creamos uno vacío en memoria para pasarlo a la plantilla
    if fallo == nil {
        fallo = &models.FallosPropuesta{
            IDPropuesta: id, // Es útil pasar el ID para el formulario
        }
    }

    data := make(map[string]interface{})
    data["Fallo"] = fallo
    
    render.RenderTemplate(w, "licitaciones/fallo-propuesta.page.tmpl", &models.TemplateData{
        Data:      data,
        CSRFToken: nosurf.Token(r),
    })
}

func (m *Repository) ObtenerFalloPropuestaJSON(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    idPropuesta, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    // Usamos la nueva función que SOLO lee y puede devolver nil
    fallo, err := m.ObtenerFalloPorPropuestaID(idPropuesta)
    if err != nil {
        log.Println("Error al obtener el fallo:", err)
        http.Error(w, "Error al obtener el fallo", http.StatusInternalServerError)
        return
    }

    // SI LA FUNCIÓN DEVUELVE nil, creamos un objeto vacío para enviar al frontend
    if fallo == nil {
        fallo = &models.FallosPropuesta{} 
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "cumple_legal":          fallo.CumpleLegal,
        "cumple_administrativo": fallo.CumpleAdministrativo,
        "cumple_tecnico":        fallo.CumpleTecnico,
        "puntos_obtenidos":      fallo.PuntosObtenidos,
        "ganador":               fallo.Ganador,
        "observaciones":         fallo.Observaciones,
    })
}
func (m *Repository) GuardarFalloPropuesta(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    idPropuesta, _ := strconv.Atoi(r.FormValue("id_propuesta"))
    if idPropuesta == 0 {
        http.Error(w, "ID de propuesta inválido", http.StatusBadRequest)
        return
    }
    
    // Parseo de valores (tu código actual es correcto)
    cumpleLegal := r.FormValue("cumple_legal") == "on"
    cumpleAdmin := r.FormValue("cumple_administrativo") == "on"
    cumpleTecnico := r.FormValue("cumple_tecnico") == "on"
    puntos, _ := strconv.Atoi(r.FormValue("puntos_obtenidos"))
    ganador := r.FormValue("ganador") == "on"
    observaciones := r.FormValue("observaciones")

    // Query de UPSERT: Crea el fallo si no existe, o lo actualiza si ya existe.
    query := `
        INSERT INTO fallos_propuesta (
            id_propuesta, cumple_legal, cumple_administrativo, cumple_tecnico,
            puntos_obtenidos, ganador, observaciones, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
        ON DUPLICATE KEY UPDATE
            cumple_legal = VALUES(cumple_legal),
            cumple_administrativo = VALUES(cumple_administrativo),
            cumple_tecnico = VALUES(cumple_tecnico),
            puntos_obtenidos = VALUES(puntos_obtenidos),
            ganador = VALUES(ganador),
            observaciones = VALUES(observaciones),
            updated_at = NOW();
    `

    _, err = m.App.DB.Exec(query,
        idPropuesta,
        cumpleLegal,
        cumpleAdmin,
        cumpleTecnico,
        puntos,
        ganador,
        observaciones,
    )

    if err != nil {
        log.Println("Error al guardar fallo (UPSERT):", err)
        http.Error(w, "Error al guardar fallo", http.StatusInternalServerError)
        return
    }
    
    // Redirigimos a la página anterior (la lista de propuestas)
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

// USUARIOS
func (m *Repository) MostrarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := m.ObtenerTodosLosUsuarios()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "No se pudieron cargar los usuarios.")
		http.Redirect(w, r, "opciones/opciones.page.tmpl", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["Usuarios"] = usuarios

    render.RenderTemplate(w, "opciones/usuarios.page.tmpl", &models.TemplateData{
        Data:      data,
        CSRFToken: nosurf.Token(r),
    })
}

func (m *Repository) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Formulario inválido.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	nivel, _ := strconv.Atoi(r.Form.Get("nivel_acceso"))
	usuario := models.Usuario{
		Nombre:      r.Form.Get("nombre"),
		Email:       r.Form.Get("email"),
		NivelAcceso: nivel,
	}
	password := r.Form.Get("password")

	if password == "" {
		m.App.Session.Put(r.Context(), "error", "La contraseña no puede estar vacía.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	// Llama al HELPER de la base de datos (el que ya tenías en helpers.go)
	err = m.CrearUsuarioUnico(usuario, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "No se pudo crear el usuario. El email podría ya estar en uso.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Usuario creado con éxito.")
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func (m *Repository) EditarUsuario(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Formulario inválido.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	nivel, _ := strconv.Atoi(r.Form.Get("nivel_acceso"))
	usuario := models.Usuario{
		ID:          id,
		Nombre:      r.Form.Get("nombre"),
		Email:       r.Form.Get("email"),
		NivelAcceso: nivel,
	}
	password := r.Form.Get("password")

	// Llama al HELPER de la base de datos (el que pusiste en el archivo de handlers por error)
	err = m.ActualizarUsuario(usuario, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "No se pudo actualizar el usuario.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Usuario actualizado con éxito.")
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func (m *Repository) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "ID de usuario inválido.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	sesionUserID := m.App.Session.GetInt(r.Context(), "user_id")
	if id == sesionUserID {
		m.App.Session.Put(r.Context(), "error", "No puedes eliminar tu propia cuenta.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	// Llama al HELPER de la base de datos (tu función se llama EliminarUsuarioUnico)
	err = m.EliminarUsuarioUnico(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error al eliminar el usuario.")
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Usuario eliminado correctamente.")
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}