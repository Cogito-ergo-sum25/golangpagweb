# Details

Date : 2026-02-18 14:56:44

Directory c:\\Users\\josec\\OneDrive\\Documentos\\AUTISMO\\GOLANGUEANDO

Total : 65 files,  11731 codes, 644 comments, 1793 blanks, all 14168 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [.env](/.env) | Dotenv | 4 | 0 | 0 | 4 |
| [Dockerfile](/Dockerfile) | Docker | 13 | 12 | 10 | 35 |
| [GOLANGUEANDO.code-workspace](/GOLANGUEANDO.code-workspace) | JSON with Comments | 8 | 0 | 0 | 8 |
| [README.md](/README.md) | Markdown | 107 | 0 | 45 | 152 |
| [cmd/web/main.go](/cmd/web/main.go) | Go | 69 | 3 | 21 | 93 |
| [cmd/web/middleware.go](/cmd/web/middleware.go) | Go | 28 | 10 | 8 | 46 |
| [cmd/web/routes.go](/cmd/web/routes.go) | Go | 111 | 29 | 34 | 174 |
| [docker-compose.yml](/docker-compose.yml) | YAML | 26 | 4 | 3 | 33 |
| [go.mod](/go.mod) | Go Module File | 13 | 0 | 5 | 18 |
| [go.sum](/go.sum) | Go Checksum File | 14 | 0 | 1 | 15 |
| [pkg/config/config.go](/pkg/config/config.go) | Go | 15 | 1 | 3 | 19 |
| [pkg/database/database.go](/pkg/database/database.go) | Go | 30 | 3 | 6 | 39 |
| [pkg/handlers/handlers.go](/pkg/handlers/handlers.go) | Go | 3,060 | 262 | 548 | 3,870 |
| [pkg/handlers/helpers.go](/pkg/handlers/helpers.go) | Go | 2,374 | 101 | 389 | 2,864 |
| [pkg/models/entidades.go](/pkg/models/entidades.go) | Go | 27 | 0 | 6 | 33 |
| [pkg/models/licitaciones.go](/pkg/models/licitaciones.go) | Go | 55 | 3 | 11 | 69 |
| [pkg/models/partidas.go](/pkg/models/partidas.go) | Go | 127 | 11 | 29 | 167 |
| [pkg/models/producto.go](/pkg/models/producto.go) | Go | 100 | 11 | 20 | 131 |
| [pkg/models/proyecto.go](/pkg/models/proyecto.go) | Go | 34 | 4 | 8 | 46 |
| [pkg/models/templatedata.go](/pkg/models/templatedata.go) | Go | 53 | 23 | 25 | 101 |
| [pkg/models/usuarios.go](/pkg/models/usuarios.go) | Go | 7 | 0 | 1 | 8 |
| [pkg/render/render.go](/pkg/render/render.go) | Go | 130 | 13 | 31 | 174 |
| [static/css/styles.css](/static/css/styles.css) | PostCSS | 397 | 35 | 97 | 529 |
| [templates/auth/login.page.tmpl](/templates/auth/login.page.tmpl) | HTML | 57 | 2 | 7 | 66 |
| [templates/base/base.layout.tmpl](/templates/base/base.layout.tmpl) | HTML | 110 | 9 | 15 | 134 |
| [templates/calendario/calendario-vista.page.tmpl](/templates/calendario/calendario-vista.page.tmpl) | HTML | 148 | 3 | 17 | 168 |
| [templates/catalogo/catalogo.page.tmpl](/templates/catalogo/catalogo.page.tmpl) | HTML | 68 | 6 | 4 | 78 |
| [templates/catalogo/producto-detalle.page.tmpl](/templates/catalogo/producto-detalle.page.tmpl) | HTML | 45 | 4 | 0 | 49 |
| [templates/home/home.page.tmpl](/templates/home/home.page.tmpl) | HTML | 100 | 5 | 8 | 113 |
| [templates/inventario/crear-producto.page.tmpl](/templates/inventario/crear-producto.page.tmpl) | HTML | 180 | 22 | 24 | 226 |
| [templates/inventario/editar-producto.page.tmpl](/templates/inventario/editar-producto.page.tmpl) | HTML | 212 | 0 | 20 | 232 |
| [templates/inventario/gestion-catalogos.page.tmpl](/templates/inventario/gestion-catalogos.page.tmpl) | HTML | 218 | 0 | 17 | 235 |
| [templates/inventario/gestion-comercio-exterior.page.tmpl](/templates/inventario/gestion-comercio-exterior.page.tmpl) | HTML | 64 | 0 | 9 | 73 |
| [templates/inventario/gestion-ieps.page.tmpl](/templates/inventario/gestion-ieps.page.tmpl) | HTML | 95 | 0 | 10 | 105 |
| [templates/inventario/gestion-inventario.page.tmpl](/templates/inventario/gestion-inventario.page.tmpl) | HTML | 135 | 0 | 9 | 144 |
| [templates/inventario/inventario.page.tmpl](/templates/inventario/inventario.page.tmpl) | HTML | 92 | 0 | 8 | 100 |
| [templates/inventario/ver\_producto.page.tmpl](/templates/inventario/ver_producto.page.tmpl) | HTML | 93 | 6 | 8 | 107 |
| [templates/licitaciones/aclaraciones-licitacion.page.tmpl](/templates/licitaciones/aclaraciones-licitacion.page.tmpl) | HTML | 157 | 0 | 14 | 171 |
| [templates/licitaciones/aclaraciones.page.tmpl](/templates/licitaciones/aclaraciones.page.tmpl) | HTML | 27 | 0 | 0 | 27 |
| [templates/licitaciones/archivos.page.tmpl](/templates/licitaciones/archivos.page.tmpl) | HTML | 79 | 0 | 3 | 82 |
| [templates/licitaciones/editar-aclaracion.page.tmpl](/templates/licitaciones/editar-aclaracion.page.tmpl) | HTML | 63 | 0 | 8 | 71 |
| [templates/licitaciones/editar-licitacion.page.tmpl](/templates/licitaciones/editar-licitacion.page.tmpl) | HTML | 187 | 0 | 24 | 211 |
| [templates/licitaciones/editar-partida.page.tmpl](/templates/licitaciones/editar-partida.page.tmpl) | HTML | 94 | 0 | 8 | 102 |
| [templates/licitaciones/editar-propuesta.page.tmpl](/templates/licitaciones/editar-propuesta.page.tmpl) | HTML | 84 | 0 | 11 | 95 |
| [templates/licitaciones/licitacion-catalogos.page.tmpl](/templates/licitaciones/licitacion-catalogos.page.tmpl) | HTML | 142 | 0 | 9 | 151 |
| [templates/licitaciones/licitaciones.page.tmpl](/templates/licitaciones/licitaciones.page.tmpl) | HTML | 152 | 0 | 12 | 164 |
| [templates/licitaciones/mostrar-partidas.page.tmpl](/templates/licitaciones/mostrar-partidas.page.tmpl) | HTML | 184 | 2 | 21 | 207 |
| [templates/licitaciones/nueva-aclaracion-licitacion.page.tmpl](/templates/licitaciones/nueva-aclaracion-licitacion.page.tmpl) | HTML | 89 | 9 | 7 | 105 |
| [templates/licitaciones/nueva-aclaracion.page.tmpl](/templates/licitaciones/nueva-aclaracion.page.tmpl) | HTML | 77 | 8 | 13 | 98 |
| [templates/licitaciones/nueva-licitacion.page.tmpl](/templates/licitaciones/nueva-licitacion.page.tmpl) | HTML | 197 | 6 | 22 | 225 |
| [templates/licitaciones/nueva-partida.page.tmpl](/templates/licitaciones/nueva-partida.page.tmpl) | HTML | 88 | 8 | 20 | 116 |
| [templates/licitaciones/nueva-propuesta.page.tmpl](/templates/licitaciones/nueva-propuesta.page.tmpl) | HTML | 274 | 0 | 27 | 301 |
| [templates/licitaciones/nuevo-producto-partida.page.tmpl](/templates/licitaciones/nuevo-producto-partida.page.tmpl) | HTML | 148 | 0 | 13 | 161 |
| [templates/licitaciones/productos-partida.page.tmpl](/templates/licitaciones/productos-partida.page.tmpl) | HTML | 121 | 1 | 15 | 137 |
| [templates/licitaciones/propuestas.page.tmpl](/templates/licitaciones/propuestas.page.tmpl) | HTML | 183 | 1 | 24 | 208 |
| [templates/opciones/datos-referencia.page.tmpl](/templates/opciones/datos-referencia.page.tmpl) | HTML | 304 | 6 | 30 | 340 |
| [templates/opciones/empresas-externas.page.tmpl](/templates/opciones/empresas-externas.page.tmpl) | HTML | 61 | 1 | 2 | 64 |
| [templates/opciones/entidades-nueva.page.tmpl](/templates/opciones/entidades-nueva.page.tmpl) | HTML | 66 | 8 | 14 | 88 |
| [templates/opciones/entidades-opciones.page.tmpl](/templates/opciones/entidades-opciones.page.tmpl) | HTML | 35 | 0 | 3 | 38 |
| [templates/opciones/opciones.page.tmpl](/templates/opciones/opciones.page.tmpl) | HTML | 42 | 0 | 2 | 44 |
| [templates/opciones/productos-externos.page.tmpl](/templates/opciones/productos-externos.page.tmpl) | HTML | 103 | 1 | 6 | 110 |
| [templates/opciones/usuarios.page.tmpl](/templates/opciones/usuarios.page.tmpl) | HTML | 134 | 4 | 10 | 148 |
| [templates/proyectos/editar-proyecto.page.tmpl](/templates/proyectos/editar-proyecto.page.tmpl) | HTML | 52 | 2 | 6 | 60 |
| [templates/proyectos/nuevo-proyecto.page.tmpl](/templates/proyectos/nuevo-proyecto.page.tmpl) | HTML | 52 | 1 | 7 | 60 |
| [templates/proyectos/proyectos-vista.page.tmpl](/templates/proyectos/proyectos-vista.page.tmpl) | HTML | 117 | 4 | 5 | 126 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)