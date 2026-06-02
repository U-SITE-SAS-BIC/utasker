# task

**Gestor de tareas CLI, offline-first y con detección automática de proyecto.**

Task es un gestor de tareas desde la terminal que asocia automáticamente las tareas al directorio del proyecto donde trabajas. Sin daemons, sin base de datos, sin nube — solo archivos JSON en `~/.tasker/` y un archivo `.task-project` en tu proyecto.

---

## Características

- **Offline-first** — todos los datos viven en `~/.tasker/tasks.json`, sin internet
- **Proyectos automáticos** — las tareas se asocian al directorio donde estás
- **Salida a color** — iconos de estado, niveles de prioridad, etiquetas, vencimientos
- **Prioridades** — escala 1–5 con código de colores
- **Etiquetas** — organiza con tags separados por coma
- **Fechas límite** — tareas vencidas se resaltan en rojo
- **Vista panorámica** — `task board` agrupa por estado con totales
- **Exportación** — a texto plano para compartir o respaldar
- **Sin dependencias** — binario único de Go, cero runtime

---

## Instalación

### Desde el código fuente

```sh
git clone https://github.com/U-SITE-SAS-BIC/utasker.git
cd tasker
go build -o task .
sudo cp task /usr/local/bin/
```

### Vía Go install

```sh
go install github.com/U-SITE-SAS-BIC/utasker@latest
```

### Binarios precompilados

Descarga desde [releases](https://github.com/U-SITE-SAS-BIC/utasker/releases) para Linux, macOS y Windows.

---

## Inicio rápido

```sh
# Inicia el seguimiento en tu proyecto
cd mi-proyecto
task init web-app

# Agrega tareas
task add "Implementar login" -p 4 -d "Usar JWT" -t "backend,auth" --due 2026-07-15
task add "Diseñar landing page" -p 3 -t "frontend"
task add "Arreglar bug de navegación" -p 5 -t "urgente"

# Lista las tareas del proyecto actual
task list

# Marca como completada
task done 1

# Vista panorámica
task board

# Exportar
task export -f tareas.txt
```

---

## Uso

### Comandos

| Comando | Descripción |
|---------|-------------|
| `task init [nombre]` | Inicia seguimiento en el directorio actual |
| `task add <título>` | Agrega una nueva tarea |
| `task list` | Lista tareas (filtradas al proyecto actual) |
| `task board` | Panorama completo agrupado por estado |
| `task show <id>` | Muestra detalle de una tarea |
| `task done <id>` | Marca tarea como completada |
| `task undo <id>` | Reabre una tarea completada |
| `task edit <id>` | Edita campos de una tarea |
| `task delete <id>` | Elimina una tarea permanentemente |
| `task project [nombre]` | Muestra o cambia el proyecto actual |
| `task export` | Exporta tareas a texto plano |

### Banderas

| Bandera | Se usa con | Descripción |
|---------|------------|-------------|
| `-a, --all` | `list`, `board`, `export` | Muestra tareas de todos los proyectos |
| `-s, --status` | `list`, `export` | Filtra por estado (`pending`, `done`, `cancelled`) |
| `-p, --priority` | `add`, `edit` | Nivel de prioridad 1–5 |
| `-d, --desc` | `add`, `edit` | Descripción |
| `-t, --tags` | `add`, `edit` | Etiquetas separadas por coma |
| `--due` | `add`, `edit` | Fecha límite (YYYY-MM-DD) |
| `--project` | cualquiera | Sobrescribe el proyecto activo |
| `-f, --file` | `export` | Archivo de salida |

### Asociación a proyecto

Crea un archivo `.task-project` en el directorio de tu proyecto:

```sh
cd mi-proyecto
task init
# o
task init nombre-del-proyecto
```

Una vez iniciado, `task list` y los demás comandos filtran automáticamente a ese proyecto. El nombre del proyecto se guarda en `.task-project` (JSON).

Para ver tareas de todos los proyectos:

```sh
task list -a
task board -a
```

---

## Almacenamiento

Todas las tareas se guardan en un solo archivo JSON:

```
~/.tasker/tasks.json
```

La asociación al proyecto se guarda por directorio:

```
mi-proyecto/.task-project    # {"project": "nombre-del-proyecto"}
```

---

## Desarrollo

### Requisitos

- Go 1.21+

### Compilar

```sh
go build -o task .
```

### Probar

```sh
go test ./...
```

---

## Licencia

Apache 2.0 — ver [LICENSE](LICENSE).

---

## Contribuir

1. Haz fork del repositorio
2. Crea una rama (`git checkout -b feature/mi-feature`)
3. Haz commit (`git commit -am 'Agrega mi feature'`)
4. Empuja a la rama (`git push origin feature/mi-feature`)
5. Abre un Pull Request
