package handlers

import (
	"log"
	"github.com/Cogito-ergo-sum25/golangpagweb/pkg/models"
)

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

// ExisteID verifica si un ID existe en una tabla específica
func (m *Repository) ExisteID(tabla string, id int) bool {
    var count int
    
    // Mapeo de nombres de tablas y columnas ID
    tablas := map[string]struct {
        nombreTabla  string
        columnaID    string
    }{
        "marca":         {"marcas", "id_marca"},
        "marcas":        {"marcas", "id_marca"}, // Alias para plural
        "tipo":          {"tipos_producto", "id_tipo"},
        "tipos":         {"tipos_producto", "id_tipo"},
        "tipos_producto": {"tipos_producto", "id_tipo"},
        "clasificacion": {"clasificaciones", "id_clasificacion"},
        "clasificaciones": {"clasificaciones", "id_clasificacion"},
        "pais":          {"paises", "id_pais"},
        "paises":        {"paises", "id_pais"},
        "certificacion": {"certificaciones", "id_certificacion"},
        "certificaciones": {"certificaciones", "id_certificacion"},
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