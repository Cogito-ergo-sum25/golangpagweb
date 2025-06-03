# 🏗️ Sistema de Gestión de Licitaciones - INTEVI

![Go Version](https://img.shields.io/badge/Go-1.24.2-blue?logo=go)
![MySQL Version](https://img.shields.io/badge/MySQL-8.0.41-orange?logo=mysql)
![Bootstrap](https://img.shields.io/badge/Bootstrap-5.1.3-purple?logo=bootstrap)
![FontAwesome](https://img.shields.io/badge/FontAwesome-6.5.1-black?logo=fontawesome)
![Badge en Desarollo](https://img.shields.io/badge/STATUS-EN%20DESAROLLO-green)
<img src="https://badges.ws/badge/Licencia-Personal-red" />


Este proyecto es un sistema integral para administrar **licitaciones públicas**, con un enfoque especial en la organización de productos de inventario, proyectos asociados y procesos licitatorios complejos como aclaraciones, propuestas, fallos y requerimientos técnicos.

---

## 📌 ¿Qué es una licitación?

Una **licitación** es un proceso mediante el cual una entidad pública o privada solicita formalmente la adquisición de bienes o servicios, seleccionando al proveedor mediante una convocatoria competitiva. Este sistema permite controlar y documentar cada etapa de ese proceso.

---

## 🧩 Módulos del sistema

### 🗂 Inventario

CRUD completo de productos con la siguiente estructura:

- Información básica: nombre, modelo, versión, SKU, etc.
- Clasificación y origen: marca, país, tipo, clasificación.
- Documentación: ficha técnica, imagen, código del fabricante.
- Certificaciones: asociadas por medio de tabla intermedia.

### 🏗 Proyectos

Vínculo entre productos e instancias de licitaciones. Cada proyecto agrupa productos específicos con sus cantidades, precios unitarios y especificaciones técnicas para una licitación concreta.

### 📑 Licitaciones

Cada licitación contiene:

- Datos generales: nombre, tipo, estatus, fechas clave.
- **Partidas**: núcleo del proceso, cada una representa una solicitud de bien o servicio.
  - Productos relacionados
  - Requerimientos técnicos (instalación, capacitación, etc.)
  - Aclaraciones por empresa
  - Propuestas de productos por empresa
  - Fallos con evaluación técnica, administrativa y legal

### 🛠 Datos de referencia

Catálogo administrable desde el sistema para mantener los datos auxiliares como:

- Marcas
- Tipos de producto
- Clasificaciones
- Países
- Certificaciones

### 🧾 Entidades y Compañías

Administración de las instituciones que emiten las licitaciones (entidades) y compañías externas que participan.

---

## 🛠 Tecnologías utilizadas

### Backend

- **Lenguaje:** Go 1.24.2
- **Frameworks y librerías:**
  - `github.com/go-chi/chi/v5` – Enrutador ligero
  - `github.com/alexedwards/scs/v2` – Manejo de sesiones
  - `github.com/justinas/nosurf` – Protección CSRF
  - `github.com/go-sql-driver/mysql` – Driver MySQL

### Base de Datos

- **Motor:** MySQL 8.0.41
- Diseño relacional con claves foráneas, timestamps y tablas intermedias.

### Frontend

- **Bootstrap 5.1.3** – Framework CSS
- **FontAwesome 6** – Íconos
- HTML renderizado mediante plantillas Go (`html/template`)

---

## 🧱 Modelo de datos (estructura general)

El sistema se basa en un modelo relacional robusto, diseñado para gestionar el ciclo completo de licitaciones públicas, desde la administración de productos hasta la evaluación de propuestas y adjudicación de partidas. La base de datos está normalizada y estructurada en las siguientes categorías principales:

### 🏷️ Datos base y de referencia
Estas tablas permiten gestionar información esencial para el control de productos y entidades involucradas:

- `marcas`, `tipos_producto`, `clasificaciones`, `paises`, `certificaciones`: características y clasificación de productos.
- `entidades`, `compañias`, `estados_republica`: representan a las instituciones que emiten las licitaciones.
- `empresas_externas`: empresas que participan como proveedores.

### 📦 Inventario de productos
- `productos`: contiene información detallada de cada producto como nombre, modelo, serie, ficha técnica, país de origen, etc.
- `producto_certificaciones`: relación N:M entre productos y certificaciones.

### 📁 Gestión de proyectos
- `proyectos`: representan la agrupación lógica de productos seleccionados para participar en una licitación específica.
- `producto_proyecto`: define qué productos y en qué cantidades están incluidos en un proyecto.

### 📜 Licitaciones y estructura de partidas
- `licitaciones`: entidad principal del sistema. Contiene metadatos como fechas clave, tipo de licitación, entidad emisora, etc.
- `partidas`: componentes individuales de una licitación, cada una con requerimientos específicos de productos.
- `licitacion_partidas`: relación entre una licitación y sus partidas.

### 🔗 Relaciones dinámicas
Estas tablas documentan la interacción de las empresas con cada partida:

- `partida_productos`: productos ofertados por partida con precio y observaciones.
- `requerimientos_partida`: requisitos técnicos o de servicio por partida (mantenimiento, capacitación, instalación, etc.).
- `aclaraciones_partida`: preguntas o solicitudes de aclaración hechas por empresas durante la etapa de junta de aclaraciones.
- `propuestas_partida`: productos propuestos por cada empresa para una partida.
- `fallos_partida`: resultado de evaluación de propuestas, incluye cumplimiento de criterios técnicos, administrativos y legales, así como puntos obtenidos y si resultó ganador.
## 🚧 Estado del proyecto

Actualmente en desarrollo activo. Se han completado:

- CRUD de productos
- Alta de licitaciones y partidas 
- Edición de licitaciones y partidas 
- Alta de proyectos con productos vinculados
- Administración básica de licitaciones y partidas
- Registro de aclaraciones y requerimientos

En desarrollo:

- Registro completo de propuestas y fallos
- Reportes y búsquedas avanzadas
- Autenticación de usuarios
- Calendario con fechas importantes
---

## 📷 Documentación

Video explicativo 20 mayo 2025:
https://www.youtube.com/watch?v=8oavojJYFSY&t=1s

Estructura base de datos:
<img src="static\images\BD.png" alt="">

---

## 📄 Licencia

Proyecto privado de uso interno para fines administrativos. No distribuible sin autorización.

---

## ✍️ Autor

| [<img src="https://scontent.fmex3-2.fna.fbcdn.net/v/t1.6435-9/123342783_1622730034575416_3218249410147359747_n.jpg?_nc_cat=102&ccb=1-7&_nc_sid=6ee11a&_nc_eui2=AeFF9stKC6W5PC1_mo27zeV5c_dbhoV4A1dz91uGhXgDV4Pg734Ea22qZ_fZPZ8SyGai_i6kNCOUv_UtL9IbwUIp&_nc_ohc=RNvPEJHqE_YQ7kNvwG19I6x&_nc_oc=AdlEvOFz1DifdtZ1JR6Of8SXLrajWBnUjr6Omb9YxGFFbIZH_no7WNvwEywn81jAqoArfu-mGSZhMHsBJSf4VCLs&_nc_zt=23&_nc_ht=scontent.fmex3-2.fna&_nc_gid=kkgfZjbji4vmMIQ3JYiheA&oh=00_AfKWrh5zW2oVcn5kxI17O1EyYPqDG6vLSfo_MavqVIB02A&oe=6866F1C2" width=115><br><sub>**Jose Luis Valencia** Ingeniero Biomedico</sub>](https://github.com/Cogito-ergo-sum25) |
| :---: |

 

Contacto: valencia.rivera.jose.luis@gmail.com
