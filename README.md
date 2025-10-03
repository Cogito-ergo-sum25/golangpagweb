# 🏗️ Sistema de Gestión de Licitaciones - INTEVI

![Go Version](https://img.shields.io/badge/Go-1.24.2-blue?logo=go)
![MySQL Version](https://img.shields.io/badge/MySQL-8.0.41-orange?logo=mysql)
![Bootstrap](https://img.shields.io/badge/Bootstrap-5.1.3-purple?logo=bootstrap)
![FontAwesome](https://img.shields.io/badge/FontAwesome-6.5.1-black?logo=fontawesome)
![Badge en Desarrollo](https://img.shields.io/badge/STATUS-PRIMERA%20VERSI%C3%93N-green)
<img src="https://badges.ws/badge/Licencia-Personal-red" />


Este proyecto es un sistema integral para administrar **licitaciones públicas**, con un enfoque especial en la organización de productos de inventario, proyectos asociados y procesos licitatorios complejos como aclaraciones, propuestas, fallos y requerimientos técnicos.

---

## 📌 ¿Qué es una licitación?

Una **licitación** es un proceso mediante el cual una entidad pública o privada solicita formalmente la adquisición de bienes o servicios, seleccionando al proveedor mediante una convocatoria competitiva. En México, estos procesos a menudo se consultan a través de plataformas como ComprasMX o se rigen por normativas específicas para asegurar transparencia y eficiencia. Este sistema permite controlar y documentar cada etapa de ese proceso.

---

## 🧩 Módulos del sistema

### 🗂 Inventario

CRUD completo de productos con la siguiente estructura:

- Información básica: nombre, modelo, versión, SKU, etc.
- Clasificación y origen: marca, país, tipo, clasificación.
- Documentación: ficha técnica, imagen, código del fabricante.
- Certificaciones: asociadas por medio de tabla intermedia.

### 🏗 Proyectos

Este módulo agrupa las licitaciones ganadas o proyectos privados en los que se está trabajando actualmente. Cada proyecto vincula productos específicos con sus cantidades, precios unitarios y especificaciones técnicas para una licitación o proyecto concreto.

### 📑 Licitaciones

Cada licitación contiene:

- Datos generales: nombre, tipo, estatus, fechas clave.
- **Partidas**: núcleo del proceso, cada una representa una solicitud de equipo médico.
  - Productos relacionados de INTEVI
  - Requerimientos técnicos (instalación, capacitación, etc.)
  - Propuestas de productos por empresa externa
  - Fallos con evaluación técnica, administrativa y legal

### 🔐 Autenticación

Permite la creación y gestión de usuarios con nombre, correo y contraseña. Las contraseñas se encriptan utilizando bcrypt. El acceso al sistema requiere una cuenta. Se contempla la implementación de niveles de usuario con restricciones futuras.

### 🗓 Calendario

Integrado con FullCalendar JS, este módulo muestra las fechas clave de las licitaciones. Ofrece vistas anual, mensual y semanal, alertando sobre la proximidad de eventos y el tiempo transcurrido desde ellos, con la posibiliad de editar las licitaciones y acceder a las mismas, tiene la intención de ser la vista principal de administración.

### 📦 Catálogo de Productos

Un módulo dedicado a la visualización de los productos de INTEVI, permitiendo un acceso rápido a sus datos, incluyendo fichas técnicas y otra documentación relevante. Cuenta con filtros para una búsqueda eficiente de productos.

### ⚙️ Menu de Opciones

Este es el centro de administración de catálogos y entidades del sistema. Desde aquí se gestionan todos los datos maestros necesarios para el funcionamiento de la aplicación:
- Datos de Referencia: Permite la gestión de catálogos auxiliares como:
  - Marcas
  - Tipos de producto
  - Clasificaciones
  - Países
  - Certificaciones
  - Compañías
- Empresas y Productos Externos: Administración de las empresas de la competencia y el catálogo de productos que estas ofertan.
- Entidades: Administración de las instituciones (hospitales, clínicas, etc.) a las que se les da servicio.
- Usuarios: Gestión de los perfiles y credenciales de los usuarios que pueden acceder al sistema.

Nota: Las opciones para administrar "Productos" (el inventario interno) y "Proyectos" desde este menú están planificadas y visibles, pero su funcionalidad aún no ha sido implementada.

## 🛠 Tecnologías utilizadas

### Backend

- **Lenguaje:** Go 1.24.2
- **Frameworks y librerías:**
  - `github.com/go-chi/chi/v5` – Enrutador ligero
  - `github.com/alexedwards/scs/v2` – Manejo de sesiones
  - `github.com/justinas/nosurf` – Protección CSRF
  - `github.com/go-sql-driver/mysql` – Driver MySQL
  - `golang.org/x/crypto/bcrypt` – Para encriptación de contraseñas

### Base de Datos

- **Motor:** MySQL 8.0.41
- Diseño relacional con claves foráneas, timestamps y tablas intermedias.

### Frontend

- **Bootstrap 5.1.3** – Framework CSS
- **FontAwesome 6** – Íconos
- **FullCalendar@6.1.7** – Librería JavaScript para calendario
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

| [<img src="https://scontent.fmex3-2.fna.fbcdn.net/v/t1.6435-9/123342783_1622730034575416_3218249410147359747_n.jpg?_nc_cat=102&ccb=1-7&_nc_sid=6ee11a&_nc_eui2=AeFF9stKC6W5PC1_mo27zeV5c_dbhoV4A1dz91uGhXgDV4Pg734Ea22qZ_fZPZ8SyGai_i6kNCOUv_UtL9IbwUIp&_nc_ohc=EN8ZyLP8YIoQ7kNvwGCTT7C&_nc_oc=AdmQhLcLbUEv5GXib4q4xMyiYtJaJgLhhj_yrMUReo9KSmsk9gaCGNeJsVDxMFjKb67WxMcApST9-IrZqbfRNHui&_nc_zt=23&_nc_ht=scontent.fmex3-2.fna&_nc_gid=7o4d5kgeQ-Vt4kuFkHcHzg&oh=00_AffbZ5iwyuXzVxJcjURCy4h8DzMFvxVUqRsFqt9p9giSRw&oe=6906AF82" width=115><br><sub>**Jose Luis Valencia** Ingeniero Biomedico</sub>](https://github.com/Cogito-ergo-sum25) |
| :---: |

Contacto: valencia.rivera.jose.luis@gmail.com