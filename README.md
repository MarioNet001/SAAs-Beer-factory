# Sistema de Gestión Cervecera (Beer Management System)

Plataforma integral para la gestión de cervecerías, optimizando la producción, inventario y creación de etiquetas mediante IA.

## Visión del Producto
Convertirnos en la solución líder para cervecerías en América Latina, ofreciendo una plataforma en la nube que optimice la producción, garantice trazabilidad y potencie la rentabilidad.

## Stack Tecnológico
- **Lenguaje:** Go 1.24+
- **Base de Datos:** PostgreSQL
- **Infraestructura:** Docker & Docker Compose
- **Arquitectura:** Hexagonal (Ports & Adapters / Screaming Architecture)

## Estructura del Proyecto
- `src/domain/`: Lógica de negocio (Modelos, Interfaces, Servicios).
- `src/infrastructure/`: Implementación de adaptadores (Repositorios Postgres, etc.).
- `cmd/api/`: Punto de entrada de la aplicación.
- `migrations/`: Scripts de evolución de la base de datos SQL.

## Getting Started

### Prerrequisitos
- Docker & Docker Compose instalados.

### Levantando el entorno local
1. Clonar el repositorio.
2. Levantar los servicios con Docker Compose:
   ```bash
   docker-compose up --build
   ```
3. El API estará disponible en `http://localhost:8080`.

### Variables de Entorno
- `DB_URL`: URL de conexión a Postgres (ej: `postgres://user:password@db:5432/beer_db?sslmode=disable`)
- `PORT`: Puerto de escucha del servidor (default: 8080)

## Git Workflow
Este proyecto utiliza un flujo de trabajo profesional basado en Pull Requests:
1. Crear una rama: `git checkout -b feat/nombre-tarea`
2. Hacer cambios y commits (usando Conventional Commits).
3. Subir rama y crear Pull Request:
   ```bash
   git push origin feat/nombre-tarea
   gh pr create --title "feat: nombre-tarea" --body "Descripción..."
   ```

## Módulos Implementados
- **Inventario (FR-04):** Gestión de stock, auditoría y validación transaccional.
- **Recetas (FR-01):** Catálogo de recetas, versionado y validación de ingredientes.

---
*Licencia MIT*
