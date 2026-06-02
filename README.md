```
██╗   ██╗████████╗ █████╗ ███████╗██╗  ██╗███████╗██████╗
██║   ██║╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗
██║   ██║   ██║   ███████║███████╗█████╔╝ █████╗  ██████╔╝
██║   ██║   ██║   ██╔══██║╚════██║██╔═██╗ ██╔══╝  ██╔══██╗
╚██████╔╝   ██║   ██║  ██║███████║██║  ██╗███████╗██║  ██║
 ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
```

**utasker** · offline-first, project-aware CLI task manager · by [U-SITE](https://u-site.app)

Gestor de tareas desde la terminal que asocia automáticamente las tareas al directorio del proyecto donde trabajas. Sin daemons, sin base de datos, sin nube — solo JSON.

---

## Características

- **Offline-first** — datos en `~/.tasker/tasks.json`, sin internet
- **Proyectos automáticos** — tareas asociadas al directorio actual
- **Salida a color** — iconos, prioridades, etiquetas, vencimientos
- **Prioridades 1–5** — código de colores
- **Etiquetas** — organización con tags
- **Fechas límite** — vencidas en rojo
- **Vista panorámica** — `utasker board` agrupa por estado
- **Exportación** — a texto plano
- **Sin dependencias** — binario único de Go

---

## Instalación

### Desde el código fuente

```sh
git clone https://github.com/U-SITE-SAS-BIC/utasker.git
cd utasker
make build
make install          # a ~/go/bin/ (recomendado)
# o
sudo make sudo-install   # a /usr/local/bin/
```

### Vía Go install

```sh
go install github.com/U-SITE-SAS-BIC/utasker@latest
```

### Binarios precompilados

Descarga desde [releases](https://github.com/U-SITE-SAS-BIC/utasker/releases).

---

## Inicio rápido

```sh
# Inicia seguimiento en tu proyecto
cd mi-proyecto
utasker init web-app

# Agrega tareas
utasker add "Implementar login" -p 4 -d "Usar JWT" -t "backend,auth" --due 2026-07-15
utasker add "Diseñar landing" -p 3 -t "frontend"
utasker add "Arreglar bug" -p 5 -t "urgente"

# Lista del proyecto actual
utasker list

# Marca completada
utasker done 1

# Panorama completo
utasker board

# Exportar
utasker export -f tareas.txt

# Información del proyecto
utasker about
utasker version
```

---

## Comandos

| Comando | Descripción |
|---------|-------------|
| `utasker init [nombre]` | Inicia seguimiento en el directorio actual |
| `utasker add <título>` | Agrega una tarea |
| `utasker list` | Lista tareas del proyecto actual |
| `utasker board` | Panorama completo agrupado por estado |
| `utasker show <id>` | Detalle de tarea |
| `utasker done <id>` | Marca como completada |
| `utasker undo <id>` | Reabre tarea |
| `utasker edit <id>` | Edita tarea |
| `utasker delete <id>` | Elimina tarea |
| `utasker project [nombre]` | Muestra o cambia proyecto |
| `utasker export` | Exporta a texto plano |
| `utasker about` | Información del proyecto |
| `utasker version` | Versión detallada |

### Banderas

| Bandera | Se usa con | Descripción |
|---------|------------|-------------|
| `-a, --all` | `list`, `board`, `export` | Todos los proyectos |
| `-s, --status` | `list`, `export` | Filtra por estado |
| `-p, --priority` | `add`, `edit` | Prioridad 1–5 |
| `-d, --desc` | `add`, `edit` | Descripción |
| `-t, --tags` | `add`, `edit` | Etiquetas separadas por coma |
| `--due` | `add`, `edit` | Fecha límite |
| `--project` | cualquiera | Sobrescribe proyecto activo |
| `-f, --file` | `export` | Archivo de salida |

---

## Data

```
~/.tasker/tasks.json           # todas las tareas
mi-proyecto/.task-project       # {"project": "nombre"}
```

---

## Licencia

Apache 2.0 — ver [LICENSE](LICENSE).

---

## Hecho por

[U-SITE](https://u-site.app) · segundo open source ✨
