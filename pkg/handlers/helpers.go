package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

//AUTENTICACIÓN
func (m *Repository) Authenticate(email string) (int, string, error) {
    // Establecemos un timeout para la consulta para evitar que se quede colgada.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string // Usaremos el nombre de columna 'pass' de tu tabla

    // QueryRowContext es ideal porque esperamos solo una fila (o ninguna).
	row := m.App.DB.QueryRowContext(ctx, "SELECT id, pass FROM usuarios WHERE email = ?", email)
	
    err := row.Scan(&id, &hashedPassword)
	if err != nil {
        // Si el error es sql.ErrNoRows, significa que el usuario no fue encontrado.
        // Esto no es un error del sistema, sino un fallo de autenticación esperado.
		if err == sql.ErrNoRows {
			return 0, "", errors.New("usuario no encontrado")
		}
        // Para cualquier otro error, sí es un problema del sistema.
		log.Println("Error al escanear usuario:", err)
		return 0, "", err
	}

	return id, hashedPassword, nil
}



// FUNCIONES EXTRA DEL RENDER

// ACTUALIZADORES
func (m *Repository) ActualizarLicitacion(l models.Licitacion) error {
	query := `
		UPDATE licitaciones SET 
			id_entidad=?, num_contratacion=?, caracter=?, nombre=?, estatus=?, tipo=?, 
			fecha_junta=?, fecha_propuestas=?, fecha_fallo=?, fecha_entrega=?, 
			tiempo_entrega=?, revisada=?, intevi=?, observaciones_generales=?, 
			updated_at=?, criterio_evaluacion=?
		WHERE id_licitacion=?`

	_, err := m.App.DB.Exec(query,
		l.IDEntidad, l.NumContratacion, l.Caracter, l.Nombre, l.Estatus, l.Tipo,
		l.FechaJunta, l.FechaPropuestas, l.FechaFallo, l.FechaEntrega,
		l.TiempoEntrega, l.Revisada, l.Intevi, l.ObservacionesGenerales,
		l.UpdatedAt, l.CriterioEvaluacion, l.IDLicitacion)

	return err
}

func (m *Repository) ActualizarPartida(p models.Partida) error {
	query := `
        UPDATE partidas SET 
            numero_partida_convocatoria = ?, 
            nombre_descripcion = ?, 
            cantidad = ?, 
            cantidad_minima = ?, 
            cantidad_maxima = ?, 
            no_ficha_tecnica = ?, 
            tipo_de_bien = ?, 
            clave_compendio = ?, 
            clave_cucop = ?, 
            unidad_medida = ?, 
            dias_de_entrega = ?, 
            fecha_de_entrega = ?, 
            garantia = ?,
            updated_at = ?
        WHERE id_partida = ?
    `

	_, err := m.App.DB.Exec(query,
		p.NumPartidaConvocatoria,
		p.NombreDescripcion,
		p.Cantidad,
		p.CantidadMinima,
		p.CantidadMaxima,
		p.NoFichaTecnica,
		p.TipoDeBien,
		p.ClaveCompendio,
		p.ClaveCucop,
		p.UnidadMedida,
		p.DiasDeEntrega,
		p.FechaDeEntrega,
		p.Garantia,
		p.UpdatedAt,
		p.IDPartida,
	)

	return err
}

func (m *Repository) ActualizarProductoPartida(p models.PartidaProductos) error {
    query := `
        UPDATE partida_productos
        SET precio_ofertado = ?, observaciones = ?, updated_at = NOW()
        WHERE id_partida_producto = ?;
    `

    _, err := m.App.DB.Exec(query, p.PrecioOfertado, p.Observaciones, p.IDPartidaProducto)
    return err
}

func (m *Repository) ActualizarPropuesta(propuesta models.PropuestasPartida) error {
    query := `
        UPDATE propuestas_partida 
        SET 
            id_producto_externo = ?,
            id_empresa = ?,
            precio_ofertado = ?,
            precio_min = ?,
            precio_max = ?,
            observaciones = ?
        WHERE id_propuesta = ?`

    _, err := m.App.DB.Exec(
        query,
        propuesta.IDProductoExterno,
        propuesta.IDEmpresa,
        propuesta.PrecioOfertado,
        propuesta.PrecioMin,
        propuesta.PrecioMax,
        propuesta.Observaciones,
        propuesta.IDPropuesta,
    )

    return err
}

func (m *Repository) ActualizarUsuario(u models.Usuario, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if password != "" {
		// Si se proveyó una nueva contraseña, la hasheamos y la actualizamos
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return err
		}
		query := `UPDATE usuarios SET nombre = ?, email = ?, nivel_acceso = ?, pass = ? WHERE id = ?`
		_, err = m.App.DB.ExecContext(ctx, query, u.Nombre, u.Email, u.NivelAcceso, string(hashedPassword), u.ID)
		return err
	} else {
		// Si no hay contraseña nueva, actualizamos todo lo demás
		query := `UPDATE usuarios SET nombre = ?, email = ?, nivel_acceso = ? WHERE id = ?`
		_, err := m.App.DB.ExecContext(ctx, query, u.Nombre, u.Email, u.NivelAcceso, u.ID)
		return err
	}
}

func (m *Repository) ActualizarProyectoEnDB(proyecto models.Proyecto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var fechaFinParaDB interface{}
	// Si la fecha de fin no es la fecha "cero" de Go, la usamos.
	// Si es la fecha cero, la dejamos como nil para que se inserte NULL en la BD.
	if !proyecto.FechaFin.IsZero() {
		fechaFinParaDB = proyecto.FechaFin
	} else {
		fechaFinParaDB = nil
	}

	stmt := `
		UPDATE proyectos SET
			nombre = ?,
			descripcion = ?,
			id_licitacion = ?,
			fecha_inicio = ?,
			fecha_fin = ?,
			updated_at = ?
		WHERE id_proyecto = ?`

	_, err := m.App.DB.ExecContext(ctx, stmt,
		proyecto.Nombre,
		proyecto.Descripcion,
		proyecto.IDLicitacion,
		proyecto.FechaInicio,
		fechaFinParaDB,
		time.Now(),
		proyecto.IDProyecto,
	)

	return err
}

func (m *Repository) ActualizarAclaracion(a models.AclaracionesLicitacion) error {
    query := `
        UPDATE aclaraciones_licitacion SET
            id_partida = ?, 
            id_empresa = ?, 
            pregunta = ?, 
            observaciones = ?,
            ficha_tecnica_id = ?, 
            id_puntos_tecnicos_modif = ?, 
            updated_at = NOW()
        WHERE id_aclaracion_licitacion = ?
    `

    // Preparamos las variables para los campos que pueden ser NULL,
    // usando la misma lógica que en la función de Insertar.
    var (
        idPartida interface{} = nil
        ftID      interface{} = nil
        ptID      interface{} = nil
    )

    if a.IDPartida != 0 {
        idPartida = a.IDPartida
    }
    if a.FichaTecnicaID != "" {
        ftID = a.FichaTecnicaID
    }
    if a.IDPuntosTecnicosModif != 0 {
        ptID = a.IDPuntosTecnicosModif
    }

    _, err := m.App.DB.Exec(query,
        idPartida,
        a.IDEmpresa,
        a.Pregunta,
        a.Observaciones,
        ftID,
        ptID,
        a.IDAclaracionLicitacion, // El ID para el WHERE
    )

    return err
}

func (m *Repository) ActualizarCatalogo(c models.ProductoCatalogo) error {
	query := `
		UPDATE producto_catalogos 
		SET nombre_version = ?, archivo_url = ?, descripcion = ?, updated_at = NOW()
		WHERE id_catalogo = ?`

	_, err := m.App.DB.Exec(query, c.NombreVersion, c.ArchivoURL, c.Descripcion, c.IDCatalogo)
	return err
}






// GETTERS

func (m *Repository) ObtenerInventarioPorID(idProducto int) (models.ProductoInventario, error) {
    var inv models.ProductoInventario
    
    // Usamos COALESCE para evitar que los NULL rompan el Scan
    query := `
        SELECT 
            COALESCE(unidad_base, 'PIEZA'), 
			COALESCE(unidad_medida_almacen, 'PIEZAS'), 
			COALESCE(metodo_costeo, 'COSTO PROMEDIO'),
            COALESCE(largo, 0.0), 
            COALESCE(ancho, 0.0), 
            COALESCE(alto, 0.0), 
            COALESCE(peso, 0.0), 
            COALESCE(volumen, 0.0),
            COALESCE(requiere_pesaje, 0), 
            COALESCE(considerar_compra_programada, 1), 
            COALESCE(produccion_fabricacion, 0),
            COALESCE(ventas_sin_existencia, 0), 
            COALESCE(maneja_serie, 0), 
            COALESCE(maneja_lote, 0), 
            COALESCE(maneja_fecha_caducidad, 0), 
            COALESCE(lote_automatico, 0)
        FROM producto_inventario 
        WHERE id_producto = ?`
    
    err := m.App.DB.QueryRow(query, idProducto).Scan(
        &inv.UnidadBase, 
        &inv.UnidadMedidaAlmacen, 
        &inv.MetodoCosteo, 
        &inv.Largo, 
        &inv.Ancho, 
        &inv.Alto, 
        &inv.Peso, 
        &inv.Volumen,
        &inv.RequierePesaje, 
        &inv.ConsiderarCompraProgramada, 
        &inv.ProduccionFabricacion,
        &inv.VentasSinExistencia, 
        &inv.ManejaSerie, 
        &inv.ManejaLote, 
        &inv.ManejaFechaCaducidad, 
        &inv.LoteAutomatico,
    )
    
    if err != nil {
        return inv, err
    }
    
    inv.IDProducto = idProducto
    return inv, nil
}

// ObtenerIEPSPorID busca los datos fiscales por ID de producto
func (m *Repository) ObtenerIEPSPorID(idProducto int) (models.IEPS, error) {
	var ieps models.IEPS
	query := `
		SELECT 
			COALESCE(tipo_producto, ''), 
			COALESCE(clave_producto, ''), 
			COALESCE(empaque, ''), 
			COALESCE(unidad_medida, ''), 
			COALESCE(presentacion, 0.00)
		FROM producto_IEPS 
		WHERE id_producto = ?`

	err := m.App.DB.QueryRow(query, idProducto).Scan(
		&ieps.TipoProducto, 
		&ieps.ClaveProducto, 
		&ieps.Empaque, 
		&ieps.UnidadMedida, 
		&ieps.Presentacion,
	)

	ieps.IDProducto = idProducto
	return ieps, err
}

func (m *Repository) ObtenerComercioExteriorPorID(id int) (models.ComercioExterior, error) {
    var ce models.ComercioExterior
    query := `
        SELECT 
            id_producto, 
            COALESCE(modelo, ''), 
            COALESCE(sub_modelo, ''), 
            COALESCE(fraccion_arancelaria, ''), 
            COALESCE(unidad_medida_aduana, ''), 
            COALESCE(factor_conversion_umt, 0.0)
        FROM producto_comercio_exterior 
        WHERE id_producto = ?`

    err := m.App.DB.QueryRow(query, id).Scan(
        &ce.IDProducto,
        &ce.Modelo,
        &ce.SubModelo,
        &ce.FraccionArancelaria,
        &ce.UnidadMedidaAduana,
        &ce.FactorConversionUMT,
    )

    return ce, err
}

func (m *Repository) ObtenerCatalogosPorProductoID(idProducto int) ([]models.ProductoCatalogo, error) {
    var catalogos []models.ProductoCatalogo
    
    // Cambiamos created_at por updated_at para reflejar la última actividad
    query := `
        SELECT 
            c.id_catalogo, 
            c.id_producto, 
            COALESCE(c.id_partida_producto, 0), 
            c.nombre_version, 
            c.archivo_url, 
            c.descripcion, 
            c.updated_at,
            COALESCE(l.num_contratacion, 'Catálogo General') as num_contratacion,
            COALESCE(p.numero_partida_convocatoria, 0) as num_partida
        FROM producto_catalogos c
        LEFT JOIN licitaciones l ON c.id_licitacion = l.id_licitacion
        LEFT JOIN partida_productos pp ON c.id_partida_producto = pp.id_partida_producto
        LEFT JOIN partidas p ON pp.id_partida = p.id_partida
        WHERE c.id_producto = ?
        ORDER BY c.updated_at DESC`

    rows, err := m.App.DB.Query(query, idProducto)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var c models.ProductoCatalogo
        var numContratacion string
        var numPartida int

        // Escaneamos c.UpdatedAt en lugar de CreatedAt
        err := rows.Scan(
            &c.IDCatalogo, 
            &c.IDProducto, 
            &c.IDPartidaProducto,
            &c.NombreVersion, 
            &c.ArchivoURL, 
            &c.Descripcion, 
            &c.UpdatedAt,
            &numContratacion, 
            &numPartida,
        )
        if err != nil {
            return nil, err
        }
        
        // Mapeamos los datos de contexto para las tablas de Intevi
        c.ContextoLicitacion = numContratacion
        c.ContextoPartida = numPartida
        
        catalogos = append(catalogos, c)
    }

    return catalogos, nil
}

func (m *Repository) ObtenerTodasLasPartidasDelProducto(idProducto int) ([]models.PartidaProductos, error) {
    var relaciones []models.PartidaProductos

    query := `
        SELECT 
            pp.id_partida_producto, 
            p.numero_partida_convocatoria, 
            p.nombre_descripcion,
            l.num_contratacion
        FROM partida_productos pp
        INNER JOIN partidas p ON pp.id_partida = p.id_partida
        INNER JOIN licitacion_partidas lp ON p.id_partida = lp.id_partida
        INNER JOIN licitaciones l ON lp.id_licitacion = l.id_licitacion
        WHERE pp.id_producto = ?
        ORDER BY l.created_at DESC, p.numero_partida_convocatoria ASC`

    rows, err := m.App.DB.Query(query, idProducto)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var pp models.PartidaProductos
        var p models.Partida
        var l models.Licitacion

        err := rows.Scan(
            &pp.IDPartidaProducto, 
            &p.NumPartidaConvocatoria, 
            &p.NombreDescripcion,
            &l.NumContratacion,
        )
        if err != nil {
            return nil, err
        }
        pp.Partida = &p
        pp.NumeroContratacion = l.NumContratacion
        relaciones = append(relaciones, pp)
    }
    return relaciones, nil
}

func (m *Repository) ObtenerProductoPorID(id int) (models.Producto, error) {
	var p models.Producto
	
	query := `
		SELECT 
			p.id_producto, p.sku, p.nombre, p.nombre_corto, 
			p.modelo, p.version, p.serie, p.codigo_fabricante, 
			p.descripcion, p.imagen_url, p.ficha_tecnica_url,
			m.nombre as marca,
			p.id_marca, p.id_tipo, p.id_clasificacion, p.id_pais_origen
		FROM productos p
		LEFT JOIN marcas m ON p.id_marca = m.id_marca
		WHERE p.id_producto = ?`

	err := m.App.DB.QueryRow(query, id).Scan(
		&p.IDProducto, &p.SKU, &p.Nombre, &p.NombreCorto,
		&p.Modelo, &p.Version, &p.Serie, &p.CodigoFabricante,
		&p.Descripcion, &p.ImagenURL, &p.FichaTecnicaURL,
		&p.Marca, // Este campo debe estar en tu struct Producto como string
		&p.IDMarca, &p.IDTipo, &p.IDClasificacion, &p.IDPaisOrigen,
	)

	if err != nil {
		return p, err
	}

	return p, nil
}

// ObtenerMarcas devuelve todas las marcas para los selects
func (m *Repository) ObtenerMarcas() ([]models.Marca, error) {
	var marcas []models.Marca
	rows, err := m.App.DB.Query("SELECT id_marca, nombre FROM marcas ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener marcas:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var marca models.Marca
		err := rows.Scan(&marca.IDMarca, &marca.Nombre)
		if err != nil {
			log.Println("Error al escanear marca:", err)
			continue
		}
		marcas = append(marcas, marca)
	}

	return marcas, nil
}

// ObtenerTiposProducto devuelve todos los tipos de producto
func (m *Repository) ObtenerTiposProducto() ([]models.TipoProducto, error) {
	var tipos []models.TipoProducto
	rows, err := m.App.DB.Query("SELECT id_tipo, nombre FROM tipos_producto ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener tipos de producto:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tipo models.TipoProducto
		err := rows.Scan(&tipo.IDTipo, &tipo.Nombre)
		if err != nil {
			log.Println("Error al escanear tipo de producto:", err)
			continue
		}
		tipos = append(tipos, tipo)
	}

	return tipos, nil
}

// ObtenerClasificaciones devuelve todas las clasificaciones
func (m *Repository) ObtenerClasificaciones() ([]models.Clasificacion, error) {
	var clasificaciones []models.Clasificacion
	rows, err := m.App.DB.Query("SELECT id_clasificacion, nombre FROM clasificaciones ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener clasificaciones:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var clasif models.Clasificacion
		err := rows.Scan(&clasif.IDClasificacion, &clasif.Nombre)
		if err != nil {
			log.Println("Error al escanear clasificación:", err)
			continue
		}
		clasificaciones = append(clasificaciones, clasif)
	}

	return clasificaciones, nil
}

// ObtenerPaises devuelve todos los países
func (m *Repository) ObtenerPaises() ([]models.Pais, error) {
	var paises []models.Pais
	rows, err := m.App.DB.Query("SELECT id_pais, nombre, codigo FROM paises ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener países:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pais models.Pais
		err := rows.Scan(&pais.IDPais, &pais.Nombre, &pais.Codigo)
		if err != nil {
			log.Println("Error al escanear país:", err)
			continue
		}
		paises = append(paises, pais)
	}

	return paises, nil
}

// ObtenerCertificaciones devuelve todas las certificaciones
func (m *Repository) ObtenerCertificaciones() ([]models.Certificacion, error) {
	var certificaciones []models.Certificacion
	rows, err := m.App.DB.Query(`
		SELECT id_certificacion, nombre, organismo_emisor 
		FROM certificaciones 
		ORDER BY nombre`)
	if err != nil {
		log.Println("Error al obtener certificaciones:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cert models.Certificacion
		err := rows.Scan(&cert.IDCertificacion, &cert.Nombre, &cert.OrganismoEmisor)
		if err != nil {
			log.Println("Error al escanear certificación:", err)
			continue
		}
		certificaciones = append(certificaciones, cert)
	}

	return certificaciones, nil
}

func (m *Repository) ObtenerEstados() ([]models.EstadosRepublica, error) {
	var estados []models.EstadosRepublica

	query := `
        SELECT clave_estado, nombre 
        FROM estados_republica 
        ORDER BY nombre
    `

	rows, err := m.App.DB.Query(query)
	if err != nil {
		log.Println("Error al obtener estados:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var estado models.EstadosRepublica
		err := rows.Scan(&estado.ClaveEstado, &estado.NombreEstado)
		if err != nil {
			log.Println("Error al escanear estado:", err)
			continue
		}
		estados = append(estados, estado)
	}

	return estados, nil
}

func (m *Repository) ObtenerCompañias() ([]models.Compañias, error) {
	var compañias []models.Compañias

	query := `
        SELECT id_compañia, nombre, tipo
        FROM compañias
        ORDER BY nombre
    `

	rows, err := m.App.DB.Query(query)
	if err != nil {
		log.Println("Error al obtener compañias:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var compañia models.Compañias
		err := rows.Scan(&compañia.IDCompañia, &compañia.Nombre, &compañia.Tipo)
		if err != nil {
			log.Println("Error al escanear compañia:", err)
			continue
		}
		compañias = append(compañias, compañia)
	}

	return compañias, nil
}

// ObtenerTodosProductos devuelve todos los productos
func (m *Repository) ObtenerTodosProductos() ([]models.Producto, error) {
	query := `SELECT p.id_producto, p.sku, m.nombre as marca, c.nombre as clasificacion,
              p.nombre_corto, p.modelo, p.nombre, p.version, p.serie,
              p.codigo_fabricante, p.descripcion, p.imagen_url, p.ficha_tecnica_url
              FROM productos p
              LEFT JOIN marcas m ON p.id_marca = m.id_marca
              LEFT JOIN clasificaciones c ON p.id_clasificacion = c.id_clasificacion
              ORDER BY p.id_producto DESC`

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(
			&p.IDProducto, &p.SKU, &p.Marca, &p.Clasificacion,
			&p.NombreCorto, &p.Modelo, &p.Nombre, &p.Version,
			&p.Serie, &p.CodigoFabricante, &p.Descripcion, &p.ImagenURL, &p.FichaTecnicaURL,
		)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, nil
}

func (m *Repository) ObtenerProductosParaInventario() ([]models.Producto, error) {
	// Preparamos la consulta con todos los JOINs necesarios
	query := `
    SELECT 
        p.id_producto,
        p.sku, 
        COALESCE(m.nombre, '') as marca,
        COALESCE(c.nombre, '') as clasificacion,
        COALESCE(t.nombre, '') as tipo,
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
    LEFT JOIN tipos_producto t ON p.id_tipo = t.id_tipo
    ORDER BY p.id_producto DESC`

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(
			&p.IDProducto,
			&p.SKU,
			&p.Marca,
			&p.Clasificacion,
			&p.Tipo,
			&p.NombreCorto,
			&p.Modelo,
			&p.Nombre,
			&p.Version,
			&p.Serie,
			&p.CodigoFabricante,
			&p.Descripcion,
		)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productos, nil
}

func (m *Repository) ObtenerLicitacionesParaSelect() ([]models.Licitacion, error) {
	query := `SELECT id_licitacion, nombre, num_contratacion FROM licitaciones ORDER BY id_licitacion DESC`

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var licitaciones []models.Licitacion
	for rows.Next() {
		var l models.Licitacion
		err := rows.Scan(&l.IDLicitacion, &l.Nombre, &l.NumContratacion)
		if err != nil {
			return nil, err
		}
		licitaciones = append(licitaciones, l)
	}
	return licitaciones, nil
}

func (m *Repository) ObtenerProyectosParaVista() ([]models.Proyecto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var proyectos []models.Proyecto

	// ¡CONSULTA ACTUALIZADA! Ahora también seleccionamos l.estatus
	query := `
		SELECT 
			p.id_proyecto, p.nombre, p.fecha_inicio, p.fecha_fin,
			l.num_contratacion, l.estatus
		FROM proyectos p
		JOIN licitaciones l ON p.id_licitacion = l.id_licitacion
		ORDER BY p.created_at DESC
	`

	rows, err := m.App.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Proyecto
		var fechaFin sql.NullTime

		// ¡SCAN ACTUALIZADO! Añadimos el campo para el estatus
		err := rows.Scan(
			&p.IDProyecto,
			&p.Nombre,
			&p.FechaInicio,
			&fechaFin,
			&p.NumContratacion,
			&p.Estatus, // Escaneamos el estatus de la licitación
		)
		if err != nil {
			return nil, err
		}

		if fechaFin.Valid {
			p.FechaFin = fechaFin.Time
		}

		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

func (m *Repository) ObtenerProyectoPorID(id int) (*models.Proyecto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p models.Proyecto
	var fechaFin sql.NullTime // Usamos sql.NullTime para manejar fechas que pueden ser NULL

	query := `
		SELECT 
			p.id_proyecto, p.nombre, p.descripcion, p.id_licitacion, 
			p.fecha_inicio, p.fecha_fin
		FROM proyectos p
		WHERE p.id_proyecto = ?
	`
	row := m.App.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&p.IDProyecto,
		&p.Nombre,
		&p.Descripcion,
		&p.IDLicitacion,
		&p.FechaInicio,
		&fechaFin, // Escaneamos en nuestra variable especial
	)
	if err != nil {
		return nil, err
	}

	// Si la fecha de fin es válida (no era NULL en la BD), la asignamos al struct.
	if fechaFin.Valid {
		p.FechaFin = fechaFin.Time
	}

	return &p, nil
}

func (m *Repository) ObtenerTodasEntidades() ([]models.Entidad, error) {
	query := `
        SELECT 
            e.id_entidad, 
            e.nombre, 
            c.id_compañia,
            COALESCE(c.nombre, '') AS nombre_compañia,
            COALESCE(c.tipo, '') AS tipo_compañia,
            e.estado AS estado_clave,
            COALESCE(er.nombre, '') AS nombre_estado,
            COALESCE(e.municipio, '') AS municipio,
            COALESCE(e.codigo_postal, '') AS codigo_postal,
            COALESCE(e.direccion, '') AS direccion,
            e.created_at,
            e.updated_at
        FROM entidades e
        LEFT JOIN estados_republica er ON e.estado = er.clave_estado
        LEFT JOIN compañias c ON e.id_compañia = c.id_compañia
        ORDER BY e.nombre
    `

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entidades []models.Entidad
	for rows.Next() {
		var e models.Entidad

		err := rows.Scan(
			&e.IDEntidad,
			&e.Nombre,
			&e.Compañia.IDCompañia,
			&e.Compañia.Nombre,
			&e.Compañia.Tipo,
			&e.Estado.ClaveEstado,
			&e.Estado.NombreEstado,
			&e.Municipio,
			&e.CodigoPostal,
			&e.Direccion,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entidades = append(entidades, e)
	}
	return entidades, nil
}

func (m *Repository) ObtenerTodasLicitaciones() ([]models.Licitacion, error) {
	query := `
        SELECT 
            l.id_licitacion, l.id_entidad,
            l.num_contratacion, l.caracter, l.nombre, l.estatus, l.tipo,
            l.fecha_junta, l.fecha_propuestas, l.fecha_fallo, l.fecha_entrega,
            l.tiempo_entrega, l.revisada, l.intevi,
            l.observaciones_generales,
            e.nombre AS entidad_nombre,
            COALESCE(e.municipio, '-') AS municipio,
            er.nombre AS estado_nombre,
            c.tipo AS compania_tipo,
            l.created_at, l.updated_at,criterio_evaluacion
        FROM licitaciones l
        LEFT JOIN entidades e ON l.id_entidad = e.id_entidad
        LEFT JOIN estados_republica er ON e.estado = er.clave_estado
        LEFT JOIN compañias c ON e.id_compañia = c.id_compañia
        ORDER BY l.id_licitacion DESC;
    `

	rows, err := m.App.DB.Query(query)
	if err != nil {
		log.Println("Error ejecutando la consulta de licitaciones:", err)
		return nil, err
	}
	defer rows.Close()

	var licitaciones []models.Licitacion
	for rows.Next() {
		var l models.Licitacion
		err := rows.Scan(
			&l.IDLicitacion,
			&l.IDEntidad,
			&l.NumContratacion,
			&l.Caracter,
			&l.Nombre,
			&l.Estatus,
			&l.Tipo,
			&l.FechaJunta,
			&l.FechaPropuestas,
			&l.FechaFallo,
			&l.FechaEntrega,
			&l.TiempoEntrega,
			&l.Revisada,
			&l.Intevi,
			&l.ObservacionesGenerales,
			&l.EntidadNombre,
			&l.EntidadMunicipio,
			&l.EstadoNombre,
			&l.CompaniaTipo,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.CriterioEvaluacion,
		)
		if err != nil {
			log.Println("Error escaneando fila de licitación:", err)
			return nil, err
		}
		licitaciones = append(licitaciones, l)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error al iterar rows:", err)
		return nil, err
	}

	return licitaciones, nil
}

func (m *Repository) ObtenerLicitacionPorID(id int) (models.Licitacion, error) {
	var l models.Licitacion

	query := `
        SELECT 
            l.id_licitacion, l.id_entidad,
            l.num_contratacion, l.caracter, l.nombre, l.estatus, l.tipo,
            l.fecha_junta, l.fecha_propuestas, l.fecha_fallo, l.fecha_entrega,
            l.tiempo_entrega, l.revisada, l.intevi,
            l.observaciones_generales,
            e.nombre AS entidad_nombre,
            COALESCE(e.municipio, '-') AS municipio,
            er.nombre AS estado_nombre,
            c.tipo AS compania_tipo,
            l.created_at, l.updated_at, l.criterio_evaluacion
        FROM licitaciones l
        LEFT JOIN entidades e ON l.id_entidad = e.id_entidad
        LEFT JOIN estados_republica er ON e.estado = er.clave_estado
        LEFT JOIN compañias c ON e.id_compañia = c.id_compañia
        WHERE l.id_licitacion = ?
        LIMIT 1
    `

	row := m.App.DB.QueryRow(query, id)

	err := row.Scan(
		&l.IDLicitacion,
		&l.IDEntidad,
		&l.NumContratacion,
		&l.Caracter,
		&l.Nombre,
		&l.Estatus,
		&l.Tipo,
		&l.FechaJunta,
		&l.FechaPropuestas,
		&l.FechaFallo,
		&l.FechaEntrega,
		&l.TiempoEntrega,
		&l.Revisada,
		&l.Intevi,
		&l.ObservacionesGenerales,
		&l.EntidadNombre,
		&l.EntidadMunicipio,
		&l.EstadoNombre,
		&l.CompaniaTipo,
		&l.CreatedAt,
		&l.UpdatedAt,
		&l.CriterioEvaluacion,
	)

	return l, err
}

func (m *Repository) ObtenerPartidasPorLicitacionID(idLicitacion int) ([]models.Partida, error) {
	query := `
        SELECT 
            p.id_partida,
            p.numero_partida_convocatoria,
            p.nombre_descripcion,
            p.cantidad,
            p.cantidad_minima,
            p.cantidad_maxima,
            p.no_ficha_tecnica,
            p.tipo_de_bien,
            p.clave_compendio,
            p.clave_cucop,
            p.unidad_medida,
            p.dias_de_entrega,
            p.fecha_de_entrega,
            p.garantia,
            p.created_at,
            p.updated_at
        FROM partidas p
        INNER JOIN licitacion_partidas lp ON lp.id_partida = p.id_partida
        WHERE lp.id_licitacion = ?
    `

	rows, err := m.App.DB.Query(query, idLicitacion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partidas []models.Partida

	for rows.Next() {
		var p models.Partida
		err := rows.Scan(
			&p.IDPartida,
			&p.NumPartidaConvocatoria,
			&p.NombreDescripcion,
			&p.Cantidad,
			&p.CantidadMinima,
			&p.CantidadMaxima,
			&p.NoFichaTecnica,
			&p.TipoDeBien,
			&p.ClaveCompendio,
			&p.ClaveCucop,
			&p.UnidadMedida,
			&p.DiasDeEntrega,
			&p.FechaDeEntrega,
			&p.Garantia,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		partidas = append(partidas, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return partidas, nil
}

func (m *Repository) ObtenerPartidaPorID(idPartida int) (models.Partida, error) {
    var p models.Partida
    
    // Variables para campos que pueden ser NULL
    var (
        noFichaTecnica  sql.NullString
        claveCompendio  sql.NullString
        claveCucop      sql.NullString
        fechaDeEntrega  sql.NullTime
    )

    // ESTA ES LA CONSULTA CORREGIDA CON JOIN
    // 1. Seleccionamos campos de la tabla 'p' (partidas) y el id_licitacion de la tabla 'lp' (licitacion_partidas)
    // 2. Unimos 'partidas' (p) con 'licitacion_partidas' (lp) donde los id_partida coincidan
    // 3. Filtramos por el id_partida que nos interesa
    query := `
        SELECT 
            p.id_partida, 
            lp.id_licitacion, -- ¡Obtenido de la tabla intermedia!
            p.numero_partida_convocatoria, 
            p.nombre_descripcion,
            p.cantidad, p.cantidad_minima, p.cantidad_maxima, p.no_ficha_tecnica,
            p.tipo_de_bien, p.clave_compendio, p.clave_cucop, p.unidad_medida,
            p.dias_de_entrega, p.fecha_de_entrega, p.garantia, p.created_at, p.updated_at
        FROM partidas p
        JOIN licitacion_partidas lp ON p.id_partida = lp.id_partida
        WHERE p.id_partida = ?
        LIMIT 1; -- Agregamos LIMIT 1 por si una partida estuviera en múltiples licitaciones
    `

    err := m.App.DB.QueryRow(query, idPartida).Scan(
        &p.IDPartida,
        &p.IDLicitacion, // Ahora esto funcionará
        &p.NumPartidaConvocatoria,
        &p.NombreDescripcion,
        &p.Cantidad,
        &p.CantidadMinima,
        &p.CantidadMaxima,
        &noFichaTecnica,
        &p.TipoDeBien,
        &claveCompendio,
        &claveCucop,
        &p.UnidadMedida,
        &p.DiasDeEntrega,
        &fechaDeEntrega,
        &p.Garantia,
        &p.CreatedAt,
        &p.UpdatedAt,
    )
    if err != nil {
        return models.Partida{}, err
    }

    // Transferir valores solo si no son NULL
    if noFichaTecnica.Valid { p.NoFichaTecnica = noFichaTecnica.String }
    if claveCompendio.Valid { p.ClaveCompendio = claveCompendio.String }
    if claveCucop.Valid { p.ClaveCucop = claveCucop.String }
    if fechaDeEntrega.Valid { p.FechaDeEntrega = fechaDeEntrega.Time }

    return p, nil
}

func (m *Repository) ObtenerIDLicitacionPorPartida(idPartida int) (int, error) {
	var idLicitacion int
	query := `SELECT id_licitacion FROM licitacion_partidas WHERE id_partida = ? LIMIT 1`

	err := m.App.DB.QueryRow(query, idPartida).Scan(&idLicitacion)
	return idLicitacion, err
}

func (m *Repository) ObtenerIDLicitacionPorIDPartida(idPartida int) (int, error) {
	var idLicitacion int
	query := `
        SELECT id_licitacion 
        FROM licitacion_partidas 
        WHERE id_partida = ? 
        LIMIT 1
    `
	err := m.App.DB.QueryRow(query, idPartida).Scan(&idLicitacion)
	return idLicitacion, err
}

func (m *Repository) ObtenerOCrearRequerimientos(idPartida int) (models.RequerimientosPartida, error) {
	var r models.RequerimientosPartida

	query := `
		SELECT 
			id_requerimientos, requiere_mantenimiento, requiere_instalacion, requiere_puesta_marcha, 
			requiere_capacitacion, requiere_visita_previa, fecha_visita, comentarios_visita, 
			requiere_muestra, fecha_muestra, comentarios_muestra, fecha_entrega, comentarios_entrega
		FROM requerimientos_partida
		WHERE id_partida = ?
		LIMIT 1`

	// 1. Declaramos variables que SÍ aceptan NULLs
	var (
		fechaVisita        sql.NullTime
		comentariosVisita  sql.NullString
		fechaMuestra       sql.NullTime
		comentariosMuestra sql.NullString
		fechaEntrega       sql.NullTime
		comentariosEntrega sql.NullString
	)

	row := m.App.DB.QueryRow(query, idPartida)
	// 2. Escaneamos en las variables "nullable"
	err := row.Scan(
		&r.IDRequerimientos,
		&r.RequiereMantenimiento,
		&r.RequiereInstalacion,
		&r.RequierePuestaEnMarcha,
		&r.RequiereCapacitacion,
		&r.RequiereVisitaPrevia,
		&fechaVisita, // Escaneamos en la variable que acepta NULL
		&comentariosVisita,
		&r.RequiereMuestra,
		&fechaMuestra,
		&comentariosMuestra,
		&fechaEntrega,
		&comentariosEntrega,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// La lógica para crear un registro si no existe es correcta.
			// (Aunque es mejor usar UPSERT en el guardado y que esta función solo lea).
			insert := `INSERT INTO requerimientos_partida (id_partida) VALUES (?)`
			res, err := m.App.DB.Exec(insert, idPartida)
			if err != nil {
				return r, err
			}
			lastID, _ := res.LastInsertId()
			r.IDRequerimientos = int(lastID)
			// Devolvemos el struct 'r' vacío, lo cual es correcto.
			return r, nil
		}
		// Si es cualquier otro error (como un error de conexión), lo retornamos.
		return r, err
	}
	
	// 3. Si el escaneo fue exitoso, transferimos los valores válidos al struct final.
	if fechaVisita.Valid {
		r.FechaVisita = fechaVisita.Time
	}
	if comentariosVisita.Valid {
		r.ComentariosVisita = comentariosVisita.String
	}
	if fechaMuestra.Valid {
		r.FechaMuestra = fechaMuestra.Time
	}
	if comentariosMuestra.Valid {
		r.ComentariosMuestra = comentariosMuestra.String
	}
	if fechaEntrega.Valid {
		r.FechaEntrega = fechaEntrega.Time
	}
	if comentariosEntrega.Valid {
		r.ComentariosEntrega = comentariosEntrega.String
	}

	return r, nil
}

func (m *Repository) ObtenerRequerimientosPorPartidaID(idPartida int) (models.RequerimientosPartida, error) {
	var req models.RequerimientosPartida

	// Variables especiales para poder escanear valores que pueden ser NULL en la BD
	var (
		comentariosVisita  sql.NullString
		comentariosMuestra sql.NullString
		comentariosEntrega sql.NullString
		fechaVisita        sql.NullTime
		fechaMuestra       sql.NullTime
		fechaEntrega       sql.NullTime
	)

	query := `
		SELECT 
			id_requerimientos, requiere_mantenimiento, requiere_instalacion, 
			requiere_puesta_marcha, requiere_capacitacion, requiere_visita_previa,
			fecha_visita, comentarios_visita, requiere_muestra, fecha_muestra,
			comentarios_muestra, fecha_entrega, comentarios_entrega
		FROM requerimientos_partida
		WHERE id_partida = ?
	`
	
	row := m.App.DB.QueryRow(query, idPartida)
	err := row.Scan(
		&req.IDRequerimientos,
		&req.RequiereMantenimiento,
		&req.RequiereInstalacion,
		&req.RequierePuestaEnMarcha,
		&req.RequiereCapacitacion,
		&req.RequiereVisitaPrevia,
		&fechaVisita, // Escanea en la variable que acepta NULL
		&comentariosVisita,
		&req.RequiereMuestra,
		&fechaMuestra,
		&comentariosMuestra,
		&fechaEntrega,
		&comentariosEntrega,
	)

	if err != nil {
		// Si el error es 'sql.ErrNoRows', no hay requerimientos. 
		// Devolvemos el struct vacío (req) y un error nulo. ¡Esto es correcto!
		if err == sql.ErrNoRows {
			return req, nil
		}
		// Si es cualquier otro error, sí es un problema real.
		log.Println("Error al escanear requerimientos:", err)
		return req, err
	}

	// Si el escaneo fue exitoso, transferimos los valores al struct final
	// solo si la base de datos no contenía NULL.
	if fechaVisita.Valid {
		req.FechaVisita = fechaVisita.Time
	}
	if comentariosVisita.Valid {
		req.ComentariosVisita = comentariosVisita.String
	}
	if fechaMuestra.Valid {
		req.FechaMuestra = fechaMuestra.Time
	}
	if comentariosMuestra.Valid {
		req.ComentariosMuestra = comentariosMuestra.String
	}
	if fechaEntrega.Valid {
		req.FechaEntrega = fechaEntrega.Time
	}
	if comentariosEntrega.Valid {
		req.ComentariosEntrega = comentariosEntrega.String
	}
	
	return req, nil
}

func (m *Repository) ObtenerAclaracionesPorPartidaID(idPartida int) ([]models.AclaracionesPartida, error) {
	query := `
		SELECT 
			a.id_aclaracion,
			a.pregunta,
			a.observaciones,
			a.ficha_tecnica_id,
			a.id_puntos_tecnicos_modif,
			a.created_at,
			a.updated_at,

			e.id_empresa,
			e.nombre,
			e.created_at AS empresa_created_at,
			e.updated_at AS empresa_updated_at

		FROM aclaraciones_partida a
		JOIN empresas_externas e ON a.id_empresa = e.id_empresa
		WHERE a.id_partida = ?
	`

	rows, err := m.App.DB.Query(query, idPartida)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aclaraciones []models.AclaracionesPartida

	for rows.Next() {
		var a models.AclaracionesPartida
		var empresa models.Empresas

		err := rows.Scan(
			&a.IDAclaracion,
			&a.Pregunta,
			&a.Observaciones,
			&a.FichaTecnica,
			&a.IDPuntosTecnico,
			&a.CreatedAt,
			&a.UpdatedAt,
			&empresa.IDEmpresa,
			&empresa.Nombre,
			&empresa.CreatedAt,
			&empresa.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		a.Empresa = &empresa
		aclaraciones = append(aclaraciones, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aclaraciones, nil
}

func (m *Repository) ObtenerAclaracionesPorLicitacionID(idLicitacion int) ([]models.AclaracionesLicitacion, error) {
    query := `
        SELECT 
            a.id_aclaracion_licitacion, a.id_licitacion, a.id_partida, a.id_empresa,
            a.pregunta, a.observaciones, a.ficha_tecnica_id, a.id_puntos_tecnicos_modif,
            a.pregunta_tecnica, a.created_at, a.updated_at,
            e.id_empresa, e.nombre, e.created_at AS empresa_created_at, e.updated_at AS empresa_updated_at,
            p.id_partida, p.numero_partida_convocatoria, p.nombre_descripcion
        FROM aclaraciones_licitacion a
        JOIN empresas_externas e ON a.id_empresa = e.id_empresa
        LEFT JOIN partidas p ON a.id_partida = p.id_partida
        WHERE a.id_licitacion = ?
    `

    rows, err := m.App.DB.Query(query, idLicitacion)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var aclaraciones []models.AclaracionesLicitacion

    for rows.Next() {
        var a models.AclaracionesLicitacion
        var empresa models.Empresas

        var (
            idPartida, puntosTecnicosID sql.NullInt64
            // --- CAMBIO 1: fichaTecnicaID ahora es sql.NullString ---
            fichaTecnicaID    sql.NullString 
            numConvocatoria   sql.NullInt64
            nombreDescripcion sql.NullString
        )

        err := rows.Scan(
            &a.IDAclaracionLicitacion, &a.IDLicitacion, &idPartida, &a.IDEmpresa,
            &a.Pregunta, &a.Observaciones, &fichaTecnicaID, &puntosTecnicosID,
            &a.PreguntaTecnica, &a.CreatedAt, &a.UpdatedAt,
            &empresa.IDEmpresa, &empresa.Nombre, &empresa.CreatedAt, &empresa.UpdatedAt,
            &idPartida, // nuevamente porque hicimos LEFT JOIN
            &numConvocatoria, &nombreDescripcion,
        )
        if err != nil {
            return nil, err
        }

        if idPartida.Valid {
            a.IDPartida = int(idPartida.Int64)
            a.Partida = &models.Partida{
                IDPartida:              int(idPartida.Int64),
                NumPartidaConvocatoria: int(numConvocatoria.Int64),
                NombreDescripcion:      nombreDescripcion.String,
            }
        }

        // --- CAMBIO 2: Se asigna el valor como string ---
        if fichaTecnicaID.Valid {
            a.FichaTecnicaID = fichaTecnicaID.String
        }
        
        if puntosTecnicosID.Valid {
            a.IDPuntosTecnicosModif = int(puntosTecnicosID.Int64)
        }

        a.Empresa = &empresa
        aclaraciones = append(aclaraciones, a)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return aclaraciones, nil
}

func (m *Repository) ObtenerAclaracionPorID(idAclaracion int) (models.AclaracionesLicitacion, error) {
    var a models.AclaracionesLicitacion
    var empresa models.Empresas

    query := `
        SELECT 
            a.id_aclaracion_licitacion, a.id_licitacion, a.id_partida, a.id_empresa,
            a.pregunta, a.observaciones, a.ficha_tecnica_id, a.id_puntos_tecnicos_modif,
            a.pregunta_tecnica, a.created_at, a.updated_at,
            e.id_empresa, e.nombre,
            p.id_partida, p.numero_partida_convocatoria, p.nombre_descripcion
        FROM aclaraciones_licitacion a
        LEFT JOIN empresas_externas e ON a.id_empresa = e.id_empresa
        LEFT JOIN partidas p ON a.id_partida = p.id_partida
        WHERE a.id_aclaracion_licitacion = ?
    `

    // Variables para manejar campos que pueden ser NULL
    var (
        idPartida         sql.NullInt64
        puntosTecnicosID  sql.NullInt64
        fichaTecnicaID    sql.NullString
        numConvocatoria   sql.NullInt64
        nombreDescripcion sql.NullString
    )

    // Usamos QueryRow porque esperamos una sola fila.
    row := m.App.DB.QueryRow(query, idAclaracion)
    err := row.Scan(
        &a.IDAclaracionLicitacion, &a.IDLicitacion, &idPartida, &a.IDEmpresa,
        &a.Pregunta, &a.Observaciones, &fichaTecnicaID, &puntosTecnicosID,
        &a.PreguntaTecnica, &a.CreatedAt, &a.UpdatedAt,
        &empresa.IDEmpresa, &empresa.Nombre,
        &idPartida, // Se escanea de nuevo porque es del JOIN de partidas
        &numConvocatoria, &nombreDescripcion,
    )
    if err != nil {
        return a, err
    }

    // Se asignan los valores de los campos opcionales
    if idPartida.Valid {
        a.IDPartida = int(idPartida.Int64)
        a.Partida = &models.Partida{
            IDPartida:              int(idPartida.Int64),
            NumPartidaConvocatoria: int(numConvocatoria.Int64),
            NombreDescripcion:      nombreDescripcion.String,
        }
    }

    if fichaTecnicaID.Valid {
        a.FichaTecnicaID = fichaTecnicaID.String
    }

    if puntosTecnicosID.Valid {
        a.IDPuntosTecnicosModif = int(puntosTecnicosID.Int64)
    }

    a.Empresa = &empresa
    
    return a, nil
}

func (m *Repository) ObtenerTodasEmpresas() ([]models.Empresas, error) {
	query := `
		SELECT id_empresa, nombre, created_at, updated_at
		FROM empresas_externas
	`

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empresas []models.Empresas

	for rows.Next() {
		var e models.Empresas
		err := rows.Scan(
			&e.IDEmpresa,
			&e.Nombre,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		empresas = append(empresas, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return empresas, nil
}

func (m *Repository) ObtenerProductosDePartida(idPartida int) ([]models.PartidaProductos, error) {
    query := `
        SELECT 
            pp.id_partida_producto, pp.id_producto, p.nombre, p.modelo, p.sku,
            pp.precio_ofertado, pp.observaciones, pp.created_at,
            COALESCE(pc.id_catalogo, 0),
            COALESCE(pc.nombre_version, ''),
            COALESCE(pc.archivo_url, ''),
            COALESCE(pc.descripcion, '')
        FROM partida_productos pp
        INNER JOIN productos p ON pp.id_producto = p.id_producto
        LEFT JOIN producto_catalogos pc ON pp.id_partida_producto = pc.id_partida_producto
        WHERE pp.id_partida = ?
        ORDER BY pp.created_at DESC`

    rows, err := m.App.DB.Query(query, idPartida)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var productos []models.PartidaProductos
    for rows.Next() {
        var p models.PartidaProductos
        var prod models.Producto
        
        err := rows.Scan(
            &p.IDPartidaProducto, &p.IDProducto, &prod.Nombre, &prod.Modelo, &prod.SKU,
            &p.PrecioOfertado, &p.Observaciones, &p.CreatedAt,
            &p.IDCatalogo, &p.NombreVersion, &p.ArchivoURL, &p.Descripcion,
        )
        if err != nil {
            return nil, err
        }

        p.Producto = &prod
        p.TieneCatalogo = p.IDCatalogo > 0
        productos = append(productos, p)
    }
    return productos, nil
}

func (m *Repository) ObtenerIDPartidaPorIDPartidaProducto(idPartidaProducto int) (int, error) {
    var idPartida int
    query := `SELECT id_partida FROM partida_productos WHERE id_partida_producto = ?`
    err := m.App.DB.QueryRow(query, idPartidaProducto).Scan(&idPartida)
    if err != nil {
        return 0, err
    }
    return idPartida, nil
}

func (m *Repository) ObtenerPropuestasPorPartidaID(idPartida int) ([]models.PropuestasPartida, error) {
	query := `
		SELECT 
			pp.id_propuesta,
			pp.id_partida,
			pp.id_empresa,
			pp.id_producto_externo,
			pp.precio_ofertado,
			pp.precio_min,
			pp.precio_max,
			pp.observaciones,
			pp.created_at,
			pp.updated_at,

			-- Empresa que hace la propuesta
			e.id_empresa,
			e.nombre,

			-- Producto externo
			pe.id_producto,
			pe.nombre,
			pe.modelo,
			pe.observaciones,
			pe.id_marca,
			pe.id_pais_origen,
			pe.id_empresa_externa,

			-- Marca
			m.id_marca,
			m.nombre,

			-- País
			p.id_pais,
			p.nombre,

			-- Empresa externa (fabricante/distribuidor del producto)
			ee.id_empresa,
			ee.nombre

		FROM propuestas_partida pp
		LEFT JOIN empresas_externas e ON pp.id_empresa = e.id_empresa
		LEFT JOIN productos_externos pe ON pp.id_producto_externo = pe.id_producto
		LEFT JOIN marcas m ON pe.id_marca = m.id_marca
		LEFT JOIN paises p ON pe.id_pais_origen = p.id_pais
		LEFT JOIN empresas_externas ee ON pe.id_empresa_externa = ee.id_empresa
		WHERE pp.id_partida = ?
	`

	rows, err := m.App.DB.Query(query, idPartida)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var propuestas []models.PropuestasPartida

	for rows.Next() {
		var (
			p          models.PropuestasPartida
			empresa    models.Empresas
			producto   models.ProductosExternos
			marca      models.Marca
			pais       models.Pais
			empresaExt models.Empresas
		)

		err := rows.Scan(
			&p.IDPropuesta,
			&p.IDPartida,
			&p.IDEmpresa,
			&p.IDProductoExterno,
			&p.PrecioOfertado,
			&p.PrecioMin,
			&p.PrecioMax,
			&p.Observaciones,
			&p.CreatedAt,
			&p.UpdatedAt,

			// Empresa que hace la propuesta
			&empresa.IDEmpresa,
			&empresa.Nombre,

			// Producto externo
			&producto.IDProducto,
			&producto.Nombre,
			&producto.Modelo,
			&producto.Observaciones,
			&producto.IDMarca,
			&producto.IDPaisOrigen,
			&producto.IDEmpresaExterna,

			// Marca
			&marca.IDMarca,
			&marca.Nombre,

			// País
			&pais.IDPais,
			&pais.Nombre,

			// Empresa externa (fabricante/distribuidor)
			&empresaExt.IDEmpresa,
			&empresaExt.Nombre,
		)
		if err != nil {
			return nil, err
		}

		producto.Marca = &marca
		producto.PaisOrigen = &pais
		producto.EmpresaExterna = &empresaExt

		p.Empresa = &empresa
		p.ProductoExterno = &producto

		propuestas = append(propuestas, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return propuestas, nil
}

func (m *Repository) ObtenerTodosProductosExternos() ([]models.ProductosExternos, error) {
    query := `
        SELECT 
            pe.id_producto, pe.nombre, pe.modelo, pe.observaciones,
            pe.id_marca, m.nombre AS marca_nombre,
            pe.id_pais_origen, p.nombre AS pais_nombre,
            pe.id_empresa_externa, e.nombre AS empresa_nombre
        FROM productos_externos pe
        LEFT JOIN marcas m ON pe.id_marca = m.id_marca
        LEFT JOIN paises p ON pe.id_pais_origen = p.id_pais
        LEFT JOIN empresas_externas e ON pe.id_empresa_externa = e.id_empresa
        ORDER BY pe.nombre
    `
    rows, err := m.App.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var productos []models.ProductosExternos
    for rows.Next() {
        var p models.ProductosExternos
        // Suponiendo que tu struct ProductosExternos tiene los campos correctos y punteros para marca, pais y empresa
        var marcaNombre, paisNombre, empresaNombre string
        err := rows.Scan(
            &p.IDProducto,
            &p.Nombre,
            &p.Modelo,
            &p.Observaciones,
            &p.IDMarca,
            &marcaNombre,
            &p.IDPaisOrigen,
            &paisNombre,
            &p.IDEmpresaExterna,
            &empresaNombre,
        )
        if err != nil {
            return nil, err
        }

        // Asignar los nombres a las subestructuras, si existen
        p.Marca = &models.Marca{Nombre: marcaNombre}
        p.PaisOrigen = &models.Pais{Nombre: paisNombre}
        p.EmpresaExterna = &models.Empresas{Nombre: empresaNombre}

        productos = append(productos, p)
    }
    return productos, nil
}

func (m *Repository) ObtenerPropuestaPorID(idPropuesta int) (models.PropuestasPartida, error) {
	query := `
	SELECT 
		pp.id_propuesta,
		pp.id_partida,
		pp.id_empresa,
		pp.id_producto_externo,
		pp.precio_ofertado,
		pp.precio_min,
		pp.precio_max,
		pp.observaciones,

		pe.id_producto,
		pe.nombre,
		pe.modelo,
		pe.observaciones,
		pe.id_marca,
		mar.nombre AS marca_nombre,
		pe.id_pais_origen,
		pa.nombre AS pais_nombre,
		pe.id_empresa_externa,
		ee.nombre AS empresa_nombre

	FROM propuestas_partida pp
	JOIN productos_externos pe ON pp.id_producto_externo = pe.id_producto
	LEFT JOIN marcas mar ON pe.id_marca = mar.id_marca
	LEFT JOIN paises pa ON pe.id_pais_origen = pa.id_pais
	LEFT JOIN empresas_externas ee ON pe.id_empresa_externa = ee.id_empresa
	WHERE pp.id_propuesta = ?
	LIMIT 1
	`

	var propuesta models.PropuestasPartida
	var producto models.ProductosExternos
	var marcaNombre, paisNombre, empresaNombre string

	err := m.App.DB.QueryRow(query, idPropuesta).Scan(
		&propuesta.IDPropuesta,
		&propuesta.IDPartida,
		&propuesta.IDEmpresa,
		&propuesta.IDProductoExterno,
		&propuesta.PrecioOfertado,
		&propuesta.PrecioMin,
		&propuesta.PrecioMax,
		&propuesta.Observaciones,

		&producto.IDProducto,
		&producto.Nombre,
		&producto.Modelo,
		&producto.Observaciones,
		&producto.IDMarca,
		&marcaNombre,
		&producto.IDPaisOrigen,
		&paisNombre,
		&producto.IDEmpresaExterna,
		&empresaNombre,
	)

	if err != nil {
		return propuesta, err
	}

	// Asignar subestructuras al producto
	producto.Marca = &models.Marca{IDMarca: producto.IDMarca, Nombre: marcaNombre}
	producto.PaisOrigen = &models.Pais{IDPais: producto.IDPaisOrigen, Nombre: paisNombre}
	producto.EmpresaExterna = &models.Empresas{IDEmpresa: producto.IDEmpresaExterna, Nombre: empresaNombre}

	// Asignar el producto a la propuesta
	propuesta.ProductoExterno = &producto

	return propuesta, nil
}

func (m *Repository) ObtenerFalloPorPropuestaID(idPropuesta int) (*models.FallosPropuesta, error) {
    var f models.FallosPropuesta

    query := `SELECT id_fallo, id_propuesta, cumple_legal, cumple_administrativo, cumple_tecnico, 
                     puntos_obtenidos, ganador, observaciones 
              FROM fallos_propuesta 
              WHERE id_propuesta = ?`

    row := m.App.DB.QueryRow(query, idPropuesta)
    err := row.Scan(
        &f.IDFallo,
        &f.IDPropuesta,
        &f.CumpleLegal,
        &f.CumpleAdministrativo,
        &f.CumpleTecnico,
        &f.PuntosObtenidos,
        &f.Ganador,
        &f.Observaciones,
    )

    if err != nil {
        // Si el error es 'sql.ErrNoRows', significa que no hay fallo.
        // Esto NO es un error del sistema. Devolvemos (nil, nil).
        if err == sql.ErrNoRows {
            return nil, nil
        }
        // Si es cualquier otro error, sí es un problema.
        return nil, err
    }

    // Si encontramos un fallo, devolvemos un puntero a él.
    return &f, nil
}

func (m *Repository) ObtenerOCrearFallo(idPropuesta int) (*models.FallosPropuesta, error) {
	var fallo models.FallosPropuesta

	query := `
	SELECT 
		id_fallo, id_propuesta, cumple_legal, cumple_administrativo, cumple_tecnico,
		puntos_obtenidos, ganador, observaciones, created_at, updated_at
	FROM fallos_propuesta
	WHERE id_propuesta = ?
	LIMIT 1;
	`

	row := m.App.DB.QueryRow(query, idPropuesta)

	err := row.Scan(
		&fallo.IDFallo,
		&fallo.IDPropuesta,
		&fallo.CumpleLegal,
		&fallo.CumpleAdministrativo,
		&fallo.CumpleTecnico,
		&fallo.PuntosObtenidos,
		&fallo.Ganador,
		&fallo.Observaciones,
		&fallo.CreatedAt,
		&fallo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Si no existe, lo creamos
		insert := `
		INSERT INTO fallos_propuesta (id_propuesta, cumple_legal, cumple_administrativo, cumple_tecnico, puntos_obtenidos, ganador, observaciones, created_at, updated_at)
		VALUES (?, false, false, false, 0, false, '', NOW(), NOW())
		`

		res, err := m.App.DB.Exec(insert, idPropuesta)
		if err != nil {
			return nil, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		fallo.IDFallo = int(lastID)
		fallo.IDPropuesta = idPropuesta
		fallo.CumpleLegal = false
		fallo.CumpleAdministrativo = false
		fallo.CumpleTecnico = false
		fallo.PuntosObtenidos = 0
		fallo.Ganador = false
		fallo.Observaciones = ""
		fallo.CreatedAt = time.Now()
		fallo.UpdatedAt = time.Now()

		return &fallo, nil
	} else if err != nil {
		return nil, err
	}

	return &fallo, nil
}

func (m *Repository) ObtenerTodosLosUsuarios() ([]models.Usuario, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var usuarios []models.Usuario

	query := `SELECT id, nombre, email, nivel_acceso FROM usuarios ORDER BY nombre`
	rows, err := m.App.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.Usuario
		err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &u.NivelAcceso)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	return usuarios, nil
}

func (m *Repository) ObtenerUsuarioPorID(id int) (*models.Usuario, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.Usuario
	query := `SELECT id, nombre, email, nivel_acceso FROM usuarios WHERE id = ?`
	row := m.App.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&u.ID, &u.Nombre, &u.Email, &u.NivelAcceso)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *Repository) BuscarProductosExternosPorNombre(nombre string) ([]models.ProductosExternos, error) {
    var productos []models.ProductosExternos
    
    searchTerm := "%" + nombre + "%" 

    query := `
        SELECT 
            pe.id_producto, pe.nombre, pe.modelo, pe.id_empresa_externa,
            m.nombre as marca_nombre, 
            pa.nombre as pais_nombre, 
            e.nombre as empresa_nombre
        FROM productos_externos pe
        LEFT JOIN marcas m ON pe.id_marca = m.id_marca
        LEFT JOIN paises pa ON pe.id_pais_origen = pa.id_pais
        LEFT JOIN empresas_externas e ON pe.id_empresa_externa = e.id_empresa
        WHERE pe.nombre LIKE ? OR pe.modelo LIKE ?
        LIMIT 10`

    rows, err := m.App.DB.Query(query, searchTerm, searchTerm)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var p models.ProductosExternos
        
        // 1. Variables para recibir los nombres de las tablas unidas.
        //    Usamos sql.NullString para manejar el caso de que un LEFT JOIN no encuentre coincidencia.
        var marcaNombre, paisNombre, empresaNombre sql.NullString

        // 2. Escaneamos en los campos directos de 'p' y en las variables locales.
        err := rows.Scan(
            &p.IDProducto,
            &p.Nombre,
            &p.Modelo,
            &p.IDEmpresaExterna,
            &marcaNombre,
            &paisNombre,
            &empresaNombre,
        )
        if err != nil {
            return nil, err
        }

        // 3. Construimos los structs anidados manualmente.
        //    Primero inicializamos los punteros...
        p.Marca = &models.Marca{}
        p.PaisOrigen = &models.Pais{}
        p.EmpresaExterna = &models.Empresas{}

        //    ...luego asignamos los valores si no son NULL.
        if marcaNombre.Valid {
            p.Marca.Nombre = marcaNombre.String
        }
        if paisNombre.Valid {
            p.PaisOrigen.Nombre = paisNombre.String
        }
        if empresaNombre.Valid {
            p.EmpresaExterna.Nombre = empresaNombre.String
        }
        
        productos = append(productos, p)
    }

    return productos, nil
}

func (m *Repository) ObtenerArchivosLicitacion(idLicitacion int) ([]models.ArchivoLicitacion, error) {
	// Usamos un contexto con tiempo límite para seguridad (opcional pero recomendado)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var archivos []models.ArchivoLicitacion

	// Consulta SQL ordenada por fecha de subida (lo más nuevo arriba)
	query := `
		SELECT 
			id_archivo, 
			id_licitacion, 
			nombre_archivo, 
			link_servidor, 
			tipo_archivo, 
			comentarios, 
			created_at, 
			updated_at
		FROM archivos_licitacion
		WHERE id_licitacion = ?
		ORDER BY created_at DESC`

	rows, err := m.App.DB.QueryContext(ctx, query, idLicitacion)
	if err != nil {
		return archivos, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.ArchivoLicitacion
		err := rows.Scan(
			&a.IDArchivo,
			&a.IDLicitacion,
			&a.NombreArchivo,
			&a.LinkServidor,
			&a.TipoArchivo,
			&a.Comentarios,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return archivos, err
		}
		archivos = append(archivos, a)
	}

	if err = rows.Err(); err != nil {
		return archivos, err
	}

	return archivos, nil
}

func (m *Repository) ObtenerCatalogosPorLicitacionID(idLicitacion int) ([]models.ProductoCatalogo, error) {
	var catalogos []models.ProductoCatalogo

	// Unimos catálogos con productos y partidas para dar contexto completo
	query := `
		SELECT 
			c.id_catalogo, c.id_producto, c.id_licitacion, COALESCE(c.id_partida_producto, 0), 
			c.nombre_version, c.archivo_url, c.descripcion, c.updated_at,
			p.nombre as nombre_producto,
			COALESCE(pa.numero_partida_convocatoria, 0) as num_partida
		FROM producto_catalogos c
		INNER JOIN productos p ON c.id_producto = p.id_producto
		LEFT JOIN partida_productos pp ON c.id_partida_producto = pp.id_partida_producto
		LEFT JOIN partidas pa ON pp.id_partida = pa.id_partida
		WHERE c.id_licitacion = ?
		ORDER BY pa.numero_partida_convocatoria ASC, c.updated_at DESC`

	rows, err := m.App.DB.Query(query, idLicitacion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.ProductoCatalogo
		// Usamos campos temporales para el JOIN que guardaremos en el contexto del struct
		var nombreProd string
		var numPartida int

		err := rows.Scan(
			&c.IDCatalogo, &c.IDProducto, &c.IDLicitacion, &c.IDPartidaProducto,
			&c.NombreVersion, &c.ArchivoURL, &c.Descripcion, &c.UpdatedAt,
			&nombreProd, &numPartida,
		)
		if err != nil {
			return nil, err
		}
		
		c.NombreProducto = nombreProd // Asegúrate de agregar este campo al struct o usarlo en el Map
		c.ContextoPartida = numPartida
		catalogos = append(catalogos, c)
	}

	return catalogos, nil
}

func (m *Repository) ObtenerPartidasConProductoPorLicitacion(idLicitacion int) ([]models.PartidaProductos, error) {
    var lista []models.PartidaProductos

    query := `
        SELECT 
			pp.id_partida_producto, 
			pp.id_producto, 
			p.numero_partida_convocatoria, 
			p.nombre_descripcion,
			prod.nombre AS nombre_producto
		FROM partida_productos pp
		INNER JOIN partidas p ON pp.id_partida = p.id_partida
		INNER JOIN productos prod ON pp.id_producto = prod.id_producto
		-- Necesitamos unir con la tabla intermedia para filtrar por licitación
		INNER JOIN licitacion_partidas lp ON p.id_partida = lp.id_partida
		WHERE lp.id_licitacion = ?
		ORDER BY p.numero_partida_convocatoria ASC;`

    rows, err := m.App.DB.Query(query, idLicitacion)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var pp models.PartidaProductos
        var p models.Partida
        var pr models.Producto
        err := rows.Scan(&pp.IDPartidaProducto, &pp.IDProducto, &p.NumPartidaConvocatoria, &p.NombreDescripcion, &pr.Nombre)
        if err != nil {
            return nil, err
        }
        pp.Partida = &p
        pp.Producto = &pr // Para mostrar qué producto es en el select
        lista = append(lista, pp)
    }
    return lista, nil
}

// obtenerIDProductoDesdePP busca el ID del producto vinculado a una relación Partida-Producto
func (m *Repository) obtenerIDProductoDesdePP(idPartidaProducto int) int {
	var idProducto int
	query := "SELECT id_producto FROM partida_productos WHERE id_partida_producto = ?"
	
	err := m.App.DB.QueryRow(query, idPartidaProducto).Scan(&idProducto)
	if err != nil {
		log.Println("Error al obtener ID de producto desde PartidaProducto:", err)
		return 0
	}
	
	return idProducto
}











// SETTERS

// AgregarMarca inserta una nueva marca en la base de datos
func (m *Repository) AgregarMarca(nombre string) error {
	_, err := m.App.DB.Exec("INSERT INTO marcas (nombre) VALUES (?)", nombre)
	return err
}

func (m *Repository) InsertarMarcaYDevolver(nombre string) (models.Marca, error) {
    query := "INSERT INTO marcas (nombre, created_at, updated_at) VALUES (?, NOW(), NOW())"
    res, err := m.App.DB.Exec(query, nombre)
    if err != nil {
        return models.Marca{}, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return models.Marca{}, err
    }
    return models.Marca{IDMarca: int(id), Nombre: nombre}, nil
}

// AgregarTipoProducto inserta un nuevo tipo de producto en la base de datos
func (m *Repository) AgregarTipoProducto(nombre string) error {
	_, err := m.App.DB.Exec("INSERT INTO tipos_producto (nombre) VALUES (?)", nombre)
	return err
}

// AgregarClasificacion inserta una nueva clasificación en la base de datos
func (m *Repository) AgregarClasificacion(nombre string) error {
	_, err := m.App.DB.Exec("INSERT INTO clasificaciones (nombre) VALUES (?)", nombre)
	return err
}

// AgregarPais inserta un nuevo país en la base de datos
func (m *Repository) AgregarPais(nombre, codigo string) error {
	_, err := m.App.DB.Exec("INSERT INTO paises (nombre, codigo) VALUES (?, ?)", nombre, codigo)
	return err
}

// AgregarCertificacion inserta una nueva certificación en la base de datos
func (m *Repository) AgregarCertificacion(nombre, organismoEmisor string) error {
	_, err := m.App.DB.Exec("INSERT INTO certificaciones (nombre, organismo_emisor) VALUES (?, ?)", nombre, organismoEmisor)
	return err
}

func (m *Repository) AgregarCompañia(nombre, tipo string) error {
	_, err := m.App.DB.Exec("INSERT INTO compañias (nombre, tipo) VALUES (?, ?)", nombre, tipo)
	return err
}

func (m *Repository) InsertarEntidad(entidad models.Entidad) error {
	query := `
        INSERT INTO entidades (
            nombre, id_compañia, estado, 
            municipio, codigo_postal, direccion,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := m.App.DB.Exec(query,
		entidad.Nombre,
		entidad.Compañia.IDCompañia,
		entidad.Estado.ClaveEstado,
		entidad.Municipio,
		entidad.CodigoPostal,
		entidad.Direccion,
		entidad.CreatedAt,
		entidad.UpdatedAt,
	)

	return err
}

func (m *Repository) InsertarLicitacion(licitacion models.Licitacion) error {
	query :=
		`INSERT INTO licitaciones (
        id_entidad, num_contratacion, caracter, nombre,
        estatus, tipo, fecha_junta, fecha_propuestas,
        fecha_fallo, fecha_entrega, tiempo_entrega,
        revisada, intevi, observaciones_generales,
        created_at, updated_at, criterio_evaluacion
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.App.DB.Exec(query,
		licitacion.IDEntidad,
		licitacion.NumContratacion,
		licitacion.Caracter,
		licitacion.Nombre,
		licitacion.Estatus,
		licitacion.Tipo,
		licitacion.FechaJunta,
		licitacion.FechaPropuestas,
		licitacion.FechaFallo,
		licitacion.FechaEntrega,
		licitacion.TiempoEntrega,
		licitacion.Revisada,
		licitacion.Intevi,
		licitacion.ObservacionesGenerales,
		licitacion.CreatedAt,
		licitacion.UpdatedAt,
		licitacion.CriterioEvaluacion,
	)
	return err
}

func (m *Repository) InsertarPartida(p models.Partida) (int, error) {
	query := `
        INSERT INTO partidas (
            numero_partida_convocatoria, nombre_descripcion, cantidad, cantidad_minima,
            cantidad_maxima, no_ficha_tecnica, tipo_de_bien, clave_compendio,
            clave_cucop, unidad_medida, dias_de_entrega, fecha_de_entrega,
            garantia, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
    `

	result, err := m.App.DB.Exec(
		query,
		p.NumPartidaConvocatoria,
		p.NombreDescripcion,
		p.Cantidad,
		p.CantidadMinima,
		p.CantidadMaxima,
		p.NoFichaTecnica,
		p.TipoDeBien,
		p.ClaveCompendio,
		p.ClaveCucop,
		p.UnidadMedida,
		p.DiasDeEntrega,
		p.FechaDeEntrega,
		p.Garantia,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (m *Repository) InsertarLicitacionPartida(idLicitacion, idPartida int) error {
	query := `
        INSERT INTO licitacion_partidas (
            id_licitacion, id_partida, created_at, updated_at
        ) VALUES (?, ?, NOW(), NOW());
    `
	_, err := m.App.DB.Exec(query, idLicitacion, idPartida)
	return err
}

func (m *Repository) InsertarAclaracion(a models.AclaracionesPartida) error {
	query := `
		INSERT INTO aclaraciones_partida (
			pregunta, observaciones, ficha_tecnica_id, id_puntos_tecnicos_modif,
			id_partida, id_empresa, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.App.DB.Exec(query,
		a.Pregunta,
		a.Observaciones,
		a.FichaTecnica,
		a.IDPuntosTecnico,
		a.Partida.IDPartida,
		a.Empresa.IDEmpresa,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (m *Repository) AgregarEmpresaNueva(nombre string) error {
	query := `INSERT INTO empresas_externas (nombre, created_at, updated_at) VALUES (?, NOW(), NOW())`
	_, err := m.App.DB.Exec(query, nombre)
	return err
}

func (m *Repository) InsertarEmpresaExternaYDevolver(nombre string) (models.Empresas, error) {
    query := `INSERT INTO empresas_externas (nombre, created_at, updated_at) VALUES (?, NOW(), NOW())`
    res, err := m.App.DB.Exec(query, nombre)
    if err != nil {
        return models.Empresas{}, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return models.Empresas{}, err
    }
    return models.Empresas{IDEmpresa: int(id), Nombre: nombre}, nil
}

func (m *Repository) InsertarPropuestaPartida(p models.PropuestasPartida) error {
    query := `
        INSERT INTO propuestas_partida 
        (id_partida, id_empresa, id_producto_externo, precio_ofertado, precio_min, precio_max, observaciones, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
    `
    _, err := m.App.DB.Exec(query,
        p.IDPartida,
        p.IDEmpresa,
        p.IDProductoExterno,
        p.PrecioOfertado,
        p.PrecioMin,
        p.PrecioMax,
        p.Observaciones,
    )
    if err != nil {
        // log para debug (puedes usar log.Println o fmt.Println)
        fmt.Println("Error InsertarPropuestaPartida:", err)
    }
    return err
}

func (m *Repository) InsertarAclaracionGeneral(a models.AclaracionesLicitacion) error {
    query := `
        INSERT INTO aclaraciones_licitacion (
            id_licitacion, id_partida, id_empresa, 
            pregunta, observaciones, ficha_tecnica_id, 
            id_puntos_tecnicos_modif, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

    // Campos que pueden ser NULL
    var (
        idPartida interface{} = nil
        ftID      interface{} = nil
        ptID      interface{} = nil
    )

    if a.Partida != nil && a.Partida.IDPartida != 0 {
        idPartida = a.Partida.IDPartida
    }

    // --- CAMBIO CLAVE: Comprobar si el string no está vacío ---
    if a.FichaTecnicaID != "" {
        ftID = a.FichaTecnicaID
    }
    
    if a.IDPuntosTecnicosModif != 0 {
        ptID = a.IDPuntosTecnicosModif
    }

    _, err := m.App.DB.Exec(query,
        a.IDLicitacion,
        idPartida,
        a.IDEmpresa,
        a.Pregunta,
        a.Observaciones,
        ftID,
        ptID,
    )

    return err
}

func (m *Repository) InsertarFalloPropuesta(f models.FallosPropuesta) error {
	query := `
		INSERT INTO fallos_propuesta (
			id_propuesta,
			cumple_legal,
			cumple_administrativo,
			cumple_tecnico,
			puntos_obtenidos,
			ganador,
			observaciones,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := m.App.DB.Exec(
		query,
		f.IDPropuesta,
		f.CumpleLegal,
		f.CumpleAdministrativo,
		f.CumpleTecnico,
		f.PuntosObtenidos,
		f.Ganador,
		f.Observaciones,
	)

	return err
}

func (m *Repository) CrearUsuarioUnico(u models.Usuario, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO usuarios (nombre, email, pass, nivel_acceso) VALUES (?, ?, ?, ?)`
	_, err = m.App.DB.ExecContext(ctx, query, u.Nombre, u.Email, string(hashedPassword), u.NivelAcceso)
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) CrearProyectoEnDB(proyecto models.Proyecto) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparamos la fecha de fin para que pueda ser NULL en la base de datos.
	var fechaFinParaDB interface{}
	// El método IsZero() comprueba si la fecha es '0001-01-01', el valor por defecto.
	if !proyecto.FechaFin.IsZero() {
		fechaFinParaDB = proyecto.FechaFin // Si hay fecha, la usamos.
	} else {
		fechaFinParaDB = nil // Si no hay fecha, usamos nil.
	}

	stmt := `
		INSERT INTO proyectos (
			nombre, descripcion, id_licitacion, 
			fecha_inicio, fecha_fin, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := m.App.DB.ExecContext(ctx, stmt,
		proyecto.Nombre,
		proyecto.Descripcion,
		proyecto.IDLicitacion,
		proyecto.FechaInicio,
		fechaFinParaDB, // Usamos nuestra variable especial aquí
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Repository) InsertarCatalogo(c models.ProductoCatalogo) error {
    query := `
        INSERT INTO producto_catalogos (
            id_producto, id_licitacion, id_partida_producto, 
            nombre_version, archivo_url, descripcion
        ) VALUES (?, ?, ?, ?, ?, ?)`

    // Manejo de Nulos: Si es 0, enviamos nil para que MySQL guarde NULL
    var lic interface{} = nil
    if c.IDLicitacion > 0 {
        lic = c.IDLicitacion
    }

    var part interface{} = nil
    if c.IDPartidaProducto > 0 {
        part = c.IDPartidaProducto
    }

    _, err := m.App.DB.Exec(query, 
        c.IDProducto, 
        lic, 
        part, 
        c.NombreVersion, 
        c.ArchivoURL, 
        c.Descripcion,
    )

    return err
}






//BORRADORES

func (m *Repository) EliminarMarca(id string) error {
	// Ejecutar la consulta SQL para eliminar la marca con el ID especificado
	_, err := m.App.DB.Exec("DELETE FROM marcas WHERE id_marca = ?", id)
	return err
}

func (m *Repository) EliminarTipoProducto(id string) error {
	_, err := m.App.DB.Exec("DELETE FROM tipos_producto WHERE id_tipo = ?", id)
	return err
}

func (m *Repository) EliminarClasificacion(id string) error {
	_, err := m.App.DB.Exec("DELETE FROM clasificaciones WHERE id_clasificacion = ?", id)
	return err
}

func (m *Repository) EliminarPais(id string) error {
	_, err := m.App.DB.Exec("DELETE FROM paises WHERE id_pais = ?", id)
	return err
}

func (m *Repository) EliminarCertificacion(id string) error {
	_, err := m.App.DB.Exec("DELETE FROM certificaciones WHERE id_certificacion = ?", id)
	return err
}

func (m *Repository) EliminarCompañia(id string) error {
	_, err := m.App.DB.Exec("DELETE FROM compañias WHERE id_compañia = ?", id)
	return err
}

func (m *Repository) EliminarProductoDePartida(idPartidaProducto int) error {
    _, err := m.App.DB.Exec("DELETE FROM partida_productos WHERE id_partida_producto = ?", idPartidaProducto)
    return err
}

func (m *Repository) EliminarUsuarioUnico(id int) error { // <- ¡Asegúrate que devuelva 'error'!
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM usuarios WHERE id = ?`
	// La operación ExecContext devuelve un error que debemos capturar.
	_, err := m.App.DB.ExecContext(ctx, query, id)
	
	// Devolvemos ese error (que puede ser 'nil' si todo salió bien).
	return err 
}

func (m *Repository) EliminarPartida(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "DELETE FROM partidas WHERE id_partida = ?"

    _, err := m.App.DB.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    return nil
}

func (m *Repository) EliminarCatalogoUnico(idCatalogo int) error {
	query := `DELETE FROM producto_catalogos WHERE id_catalogo = ?`
	_, err := m.App.DB.Exec(query, idCatalogo)
	return err
}

//UPSERTS

func (m *Repository) UpsertInventario(inv models.ProductoInventario) error {
	query := `
		INSERT INTO producto_inventario (
			id_producto, 
			unidad_base, 
			unidad_medida_almacen, 
			metodo_costeo, 
			largo, 
			ancho, 
			alto, 
			peso, 
			volumen,
			requiere_pesaje, 
			considerar_compra_programada, 
			produccion_fabricacion,
			ventas_sin_existencia, 
			maneja_serie, 
			maneja_lote, 
			maneja_fecha_caducidad, 
			lote_automatico
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			unidad_base = VALUES(unidad_base),
			unidad_medida_almacen = VALUES(unidad_medida_almacen),
			metodo_costeo = VALUES(metodo_costeo),
			largo = VALUES(largo),
			ancho = VALUES(ancho),
			alto = VALUES(alto),
			peso = VALUES(peso),
			volumen = VALUES(volumen),
			requiere_pesaje = VALUES(requiere_pesaje),
			considerar_compra_programada = VALUES(considerar_compra_programada),
			produccion_fabricacion = VALUES(produccion_fabricacion),
			ventas_sin_existencia = VALUES(ventas_sin_existencia),
			maneja_serie = VALUES(maneja_serie),
			maneja_lote = VALUES(maneja_lote),
			maneja_fecha_caducidad = VALUES(maneja_fecha_caducidad),
			lote_automatico = VALUES(lote_automatico)`

	// Ejecutamos pasando los 17 parámetros en el orden exacto del INSERT
	_, err := m.App.DB.Exec(query,
		inv.IDProducto,                 // 1
		inv.UnidadBase,                 // 2
		inv.UnidadMedidaAlmacen,        // 3
		inv.MetodoCosteo,               // 4
		inv.Largo,                      // 5
		inv.Ancho,                      // 6
		inv.Alto,                       // 7
		inv.Peso,                       // 8
		inv.Volumen,                    // 9
		inv.RequierePesaje,             // 10
		inv.ConsiderarCompraProgramada, // 11
		inv.ProduccionFabricacion,      // 12
		inv.VentasSinExistencia,        // 13
		inv.ManejaSerie,                // 14
		inv.ManejaLote,                 // 15
		inv.ManejaFechaCaducidad,       // 16
		inv.LoteAutomatico,             // 17
	)

	return err
}

// UpsertIEPS inserta o actualiza la configuración fiscal
func (m *Repository) UpsertIEPS(ieps models.IEPS) error {
	query := `
		INSERT INTO producto_IEPS (
			id_producto, tipo_producto, clave_producto, empaque, unidad_medida, presentacion
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			tipo_producto = VALUES(tipo_producto),
			clave_producto = VALUES(clave_producto),
			empaque = VALUES(empaque),
			unidad_medida = VALUES(unidad_medida),
			presentacion = VALUES(presentacion)`

	_, err := m.App.DB.Exec(query,
		ieps.IDProducto,
		ieps.TipoProducto,
		ieps.ClaveProducto,
		ieps.Empaque,
		ieps.UnidadMedida,
		ieps.Presentacion,
	)

	return err
}

func (m *Repository) UpsertComercioExterior(ce models.ComercioExterior) error {
    query := `
        INSERT INTO producto_comercio_exterior (
            id_producto, modelo, sub_modelo, fraccion_arancelaria, 
            unidad_medida_aduana, factor_conversion_umt
        ) VALUES (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            modelo = VALUES(modelo),
            sub_modelo = VALUES(sub_modelo),
            fraccion_arancelaria = VALUES(fraccion_arancelaria),
            unidad_medida_aduana = VALUES(unidad_medida_aduana),
            factor_conversion_umt = VALUES(factor_conversion_umt)`

    _, err :=  m.App.DB.Exec(query,
        ce.IDProducto,
        ce.Modelo,
        ce.SubModelo,
        ce.FraccionArancelaria,
        ce.UnidadMedidaAduana,
        ce.FactorConversionUMT,
    )
    return err
}

//FUNCIONES AUXILIARES

// ExisteID verifica si un ID existe en una tabla específica
func (m *Repository) ExisteID(tabla string, id int) bool {
	var count int

	// Mapeo de nombres de tablas y columnas ID
	tablas := map[string]struct {
		nombreTabla string
		columnaID   string
	}{
		"marca":           {"marcas", "id_marca"},
		"marcas":          {"marcas", "id_marca"},
		"tipo":            {"tipos_producto", "id_tipo"},
		"tipos":           {"tipos_producto", "id_tipo"},
		"tipos_producto":  {"tipos_producto", "id_tipo"},
		"clasificacion":   {"clasificaciones", "id_clasificacion"},
		"clasificaciones": {"clasificaciones", "id_clasificacion"},
		"pais":            {"paises", "id_pais"},
		"paises":          {"paises", "id_pais"},
		"certificacion":   {"certificaciones", "id_certificacion"},
		"certificaciones": {"certificaciones", "id_certificacion"},
		"licitacion":      {"licitaciones", "id_licitacion"},
		"licitaciones":    {"licitaciones", "id_licitacion"},
		"producto":        {"productos", "id_producto"},
		"productos":       {"productos", "id_producto"},
		"proyecto":        {"proyectos", "id_proyecto"},
		"proyectos":       {"proyectos", "id_proyecto"},
		"partida":         {"partidas", "id_partida"},
		"partidas":        {"partidas", "id_partida"},
	}

	config, ok := tablas[tabla]
	if !ok {
		log.Printf("Tabla no configurada: %s", tabla)
		return false
	}

	query := "SELECT COUNT(*) FROM " + config.nombreTabla + " WHERE " + config.columnaID + " = ?"
	err := m.App.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		log.Printf("Error al verificar ID en tabla %s: %v\n", config.nombreTabla, err)
		return false
	}
	return count > 0
}

// Función helper para verificar certificaciones seleccionadas
func isCertSelected(certID int, productCerts []models.Certificacion) bool {
	for _, pc := range productCerts {
		if pc.IDCertificacion == certID {
			return true
		}
	}
	return false
}

func parseDate(value string) time.Time {
	t, _ := time.Parse("2006-01-02", value)
	return t
}
func atoi(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

func atof(s string) float64 {
    f, _ := strconv.ParseFloat(s, 64)
    return f
}


