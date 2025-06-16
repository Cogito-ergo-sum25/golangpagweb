package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
)

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

// GETTERS

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
              p.codigo_fabricante, p.descripcion
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
			&p.Serie, &p.CodigoFabricante, &p.Descripcion,
		)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
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

func (m *Repository) ObtenerProyectosConRelaciones() ([]models.Proyecto, error) {
	query := `
        SELECT 
            p.id_proyecto, p.nombre, p.descripcion, p.fecha_inicio, 
            COALESCE(p.fecha_fin, CAST('1970-01-01' AS DATE)) AS fecha_fin,
            p.created_at, p.updated_at,
            l.id_licitacion, l.nombre as licitacion_nombre, l.num_contratacion,
            COALESCE(l.lugar, 'indefinido') AS lugar,
            COALESCE(l.fecha_junta, CAST('1970-01-01' AS DATE)) AS fecha_junta,
            COALESCE(l.fecha_propuestas, CAST('1970-01-01' AS DATE)) AS fecha_propuestas,
            COALESCE(l.fecha_fallo, CAST('1970-01-01' AS DATE)) AS fecha_fallo,
            COALESCE(l.fecha_entrega, CAST('1970-01-01' AS DATE)) AS fecha_entrega,
            COALESCE(l.estado, 'indefinido') AS estado,
            e.nombre as entidad_nombre,
            pp.id_producto, pp.cantidad, pp.precio_unitario, pp.especificaciones,
            pr.nombre as producto_nombre, pr.sku, pr.imagen_url, pr.modelo
        FROM 
            proyectos p
        LEFT JOIN 
            licitaciones l ON p.id_licitacion = l.id_licitacion
        LEFT JOIN 
            entidades e ON l.id_entidad = e.id_entidad
        LEFT JOIN 
            producto_proyecto pp ON p.id_proyecto = pp.id_proyecto
        LEFT JOIN 
            productos pr ON pp.id_producto = pr.id_producto
        ORDER BY 
            p.id_proyecto, pp.id_producto
    `

	rows, err := m.App.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proyectos []models.Proyecto
	var currentProyectoID int
	var currentProyecto *models.Proyecto

	for rows.Next() {
		var p models.Proyecto
		var pp models.ProductoProyecto

		err := rows.Scan(
			&p.IDProyecto, &p.Nombre, &p.Descripcion, &p.FechaInicio, &p.FechaFin,
			&p.CreatedAt, &p.UpdatedAt,
			&p.IDLicitacion, &p.LicitacionNombre, &p.NumContratacion,
			&p.Lugar,
			&p.FechaJunta, &p.FechaPropuestas, &p.FechaFallo, &p.FechaEntrega,
			&p.EstadoLicitacion,
			&p.EntidadNombre,
			&pp.IDProducto, &pp.Cantidad, &pp.PrecioUnitario, &pp.Especificaciones,
			&pp.ProductoNombre, &pp.SKU, &pp.ImagenURL, &pp.Modelo,
		)

		if err != nil {
			return nil, err
		}

		if p.IDProyecto != currentProyectoID {
			if currentProyecto != nil {
				proyectos = append(proyectos, *currentProyecto)
			}
			currentProyectoID = p.IDProyecto
			currentProyecto = &p
			currentProyecto.Productos = []models.ProductoProyecto{}
		}

		if pp.IDProducto != 0 {
			currentProyecto.Productos = append(currentProyecto.Productos, pp)
		}
	}

	if currentProyecto != nil {
		proyectos = append(proyectos, *currentProyecto)
	}

	return proyectos, nil
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
	query := `
        SELECT 
            id_partida,
            
            numero_partida_convocatoria,
            nombre_descripcion,
            cantidad,
            cantidad_minima,
            cantidad_maxima,
            no_ficha_tecnica,
            tipo_de_bien,
            clave_compendio,
            clave_cucop,
            unidad_medida,
            dias_de_entrega,
            fecha_de_entrega,
            garantia,
            created_at,
            updated_at
        FROM partidas
        WHERE id_partida = ?
    `

	var p models.Partida

	err := m.App.DB.QueryRow(query, idPartida).Scan(
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
		return models.Partida{}, err
	}

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
		id_requerimientos,
		requiere_mantenimiento,
		requiere_instalacion,
		requiere_puesta_marcha,
		requiere_capacitacion,
		requiere_visita_previa
	FROM requerimientos_partida
	WHERE id_partida = ?
	LIMIT 1`

	row := m.App.DB.QueryRow(query, idPartida)
	err := row.Scan(
		&r.IDRequerimientos,
		&r.RequiereMantenimiento,
		&r.RequiereInstalacion,
		&r.RequierePuestaEnMarcha,
		&r.RequiereCapacitacion,
		&r.RequiereVisitaPrevia,
	)

	if err == sql.ErrNoRows {
		// No existe, lo creamos por defecto
		insert := `
		INSERT INTO requerimientos_partida (
			id_partida,
			requiere_mantenimiento,
			requiere_instalacion,
			requiere_puesta_marcha,
			requiere_capacitacion,
			requiere_visita_previa
		) VALUES (?, false, false, false, false, false)`

		res, err := m.App.DB.Exec(insert, idPartida)
		if err != nil {
			return r, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return r, err
		}

		// Devolver el registro recién creado
		r.IDRequerimientos = int(lastID)
		r.RequiereMantenimiento = false
		r.RequiereInstalacion = false
		r.RequierePuestaEnMarcha = false
		r.RequiereCapacitacion = false
		r.RequiereVisitaPrevia = false

		return r, nil
	}

	if err != nil {
		return r, err
	}

	return r, nil
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
            pp.id_partida_producto,
            pp.id_producto,
            p.nombre,
            p.modelo,
            p.sku,
            pp.precio_ofertado,
            pp.observaciones,
            pp.created_at
        FROM partida_productos pp
        INNER JOIN productos p ON pp.id_producto = p.id_producto
        WHERE pp.id_partida = ?
        ORDER BY pp.created_at DESC;
    `

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
            &p.IDPartidaProducto,
            &p.IDProducto,
            &prod.Nombre,
            &prod.Modelo,
            &prod.SKU,
            &p.PrecioOfertado,
            &p.Observaciones,
            &p.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        p.Producto = &prod
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










// SETTERS

// AgregarMarca inserta una nueva marca en la base de datos
func (m *Repository) AgregarMarca(nombre string) error {
	_, err := m.App.DB.Exec("INSERT INTO marcas (nombre) VALUES (?)", nombre)
	return err
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