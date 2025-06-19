# Details

Date : 2025-06-18 20:01:10

Directory c:\\Users\\josec\\OneDrive\\Documentos\\AUTISMO\\GOLANGUEANDO

Total : 49 files,  8056 codes, 587 comments, 1302 blanks, all 9945 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [GOLANGUEANDO.code-workspace](/GOLANGUEANDO.code-workspace) | JSON with Comments | 8 | 0 | 0 | 8 |
| [README.md](/README.md) | Markdown | 107 | 0 | 56 | 163 |
| [cmd/web/main.go](/cmd/web/main.go) | Go | 59 | 3 | 18 | 80 |
| [cmd/web/middleware.go](/cmd/web/middleware.go) | Go | 18 | 2 | 5 | 25 |
| [cmd/web/routes.go](/cmd/web/routes.go) | Go | 74 | 17 | 38 | 129 |
| [go.mod](/go.mod) | Go Module File | 9 | 0 | 4 | 13 |
| [go.sum](/go.sum) | Go Checksum File | 10 | 0 | 1 | 11 |
| [pkg/config/config.go](/pkg/config/config.go) | Go | 15 | 1 | 3 | 19 |
| [pkg/database/database.go](/pkg/database/database.go) | Go | 30 | 3 | 6 | 39 |
| [pkg/handlers/handlers.go](/pkg/handlers/handlers.go) | Go | 2,186 | 306 | 382 | 2,874 |
| [pkg/handlers/helpers.go](/pkg/handlers/helpers.go) | Go | 1,533 | 36 | 246 | 1,815 |
| [pkg/models/entidades.go](/pkg/models/entidades.go) | Go | 27 | 0 | 6 | 33 |
| [pkg/models/licitaciones.go](/pkg/models/licitaciones.go) | Go | 45 | 3 | 11 | 59 |
| [pkg/models/partidas.go](/pkg/models/partidas.go) | Go | 122 | 10 | 29 | 161 |
| [pkg/models/producto.go](/pkg/models/producto.go) | Go | 48 | 7 | 13 | 68 |
| [pkg/models/proyecto.go](/pkg/models/proyecto.go) | Go | 33 | 4 | 8 | 45 |
| [pkg/models/templatedata.go](/pkg/models/templatedata.go) | Go | 48 | 21 | 22 | 91 |
| [pkg/render/render.go](/pkg/render/render.go) | Go | 127 | 13 | 31 | 171 |
| [static/css/styles.css](/static/css/styles.css) | PostCSS | 312 | 23 | 79 | 414 |
| [templates/base/base.layout.tmpl](/templates/base/base.layout.tmpl) | HTML | 104 | 8 | 11 | 123 |
| [templates/catalogo/catalogo.page.tmpl](/templates/catalogo/catalogo.page.tmpl) | HTML | 82 | 7 | 5 | 94 |
| [templates/home/home.page.tmpl](/templates/home/home.page.tmpl) | HTML | 59 | 2 | 5 | 66 |
| [templates/inventario/crear-producto.page.tmpl](/templates/inventario/crear-producto.page.tmpl) | HTML | 177 | 22 | 23 | 222 |
| [templates/inventario/editar-producto.page.tmpl](/templates/inventario/editar-producto.page.tmpl) | HTML | 187 | 22 | 23 | 232 |
| [templates/inventario/inventario.page.tmpl](/templates/inventario/inventario.page.tmpl) | HTML | 39 | 0 | 0 | 39 |
| [templates/inventario/ver\_producto.page.tmpl](/templates/inventario/ver_producto.page.tmpl) | HTML | 93 | 6 | 8 | 107 |
| [templates/licitaciones/aclaraciones-licitacion.page.tmpl](/templates/licitaciones/aclaraciones-licitacion.page.tmpl) | HTML | 45 | 0 | 4 | 49 |
| [templates/licitaciones/aclaraciones.page.tmpl](/templates/licitaciones/aclaraciones.page.tmpl) | HTML | 27 | 0 | 0 | 27 |
| [templates/licitaciones/editar-licitacion.page.tmpl](/templates/licitaciones/editar-licitacion.page.tmpl) | HTML | 163 | 6 | 16 | 185 |
| [templates/licitaciones/editar-partida.page.tmpl](/templates/licitaciones/editar-partida.page.tmpl) | HTML | 88 | 8 | 20 | 116 |
| [templates/licitaciones/editar-propuesta.page.tmpl](/templates/licitaciones/editar-propuesta.page.tmpl) | HTML | 84 | 0 | 11 | 95 |
| [templates/licitaciones/licitaciones.page.tmpl](/templates/licitaciones/licitaciones.page.tmpl) | HTML | 152 | 0 | 12 | 164 |
| [templates/licitaciones/mostrar-partidas.page.tmpl](/templates/licitaciones/mostrar-partidas.page.tmpl) | HTML | 169 | 2 | 20 | 191 |
| [templates/licitaciones/nueva-aclaracion-licitacion.page.tmpl](/templates/licitaciones/nueva-aclaracion-licitacion.page.tmpl) | HTML | 89 | 9 | 7 | 105 |
| [templates/licitaciones/nueva-aclaracion.page.tmpl](/templates/licitaciones/nueva-aclaracion.page.tmpl) | HTML | 77 | 8 | 13 | 98 |
| [templates/licitaciones/nueva-licitacion.page.tmpl](/templates/licitaciones/nueva-licitacion.page.tmpl) | HTML | 169 | 6 | 17 | 192 |
| [templates/licitaciones/nueva-partida.page.tmpl](/templates/licitaciones/nueva-partida.page.tmpl) | HTML | 88 | 8 | 20 | 116 |
| [templates/licitaciones/nueva-propuesta.page.tmpl](/templates/licitaciones/nueva-propuesta.page.tmpl) | HTML | 140 | 1 | 15 | 156 |
| [templates/licitaciones/nuevo-producto-partida.page.tmpl](/templates/licitaciones/nuevo-producto-partida.page.tmpl) | HTML | 105 | 2 | 2 | 109 |
| [templates/licitaciones/productos-partida.page.tmpl](/templates/licitaciones/productos-partida.page.tmpl) | HTML | 96 | 1 | 10 | 107 |
| [templates/licitaciones/propuestas.page.tmpl](/templates/licitaciones/propuestas.page.tmpl) | HTML | 176 | 1 | 24 | 201 |
| [templates/opciones/datos-referencia.page.tmpl](/templates/opciones/datos-referencia.page.tmpl) | HTML | 304 | 6 | 30 | 340 |
| [templates/opciones/empresas-externas.page.tmpl](/templates/opciones/empresas-externas.page.tmpl) | HTML | 61 | 1 | 2 | 64 |
| [templates/opciones/entidades-nueva.page.tmpl](/templates/opciones/entidades-nueva.page.tmpl) | HTML | 66 | 8 | 14 | 88 |
| [templates/opciones/entidades-opciones.page.tmpl](/templates/opciones/entidades-opciones.page.tmpl) | HTML | 35 | 0 | 3 | 38 |
| [templates/opciones/opciones.page.tmpl](/templates/opciones/opciones.page.tmpl) | HTML | 42 | 0 | 2 | 44 |
| [templates/opciones/productos-externos.page.tmpl](/templates/opciones/productos-externos.page.tmpl) | HTML | 103 | 1 | 6 | 110 |
| [templates/proyectos/nuevo-proyecto.page.tmpl](/templates/proyectos/nuevo-proyecto.page.tmpl) | HTML | 158 | 3 | 19 | 180 |
| [templates/proyectos/proyectos-vista.page.tmpl](/templates/proyectos/proyectos-vista.page.tmpl) | HTML | 67 | 0 | 2 | 69 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)