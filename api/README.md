# API de inventario

Backend REST en Go siguiendo la organización de modelos, DTO, controladores y rutas del proyecto de referencia. Los datos se almacenan en la base PostgreSQL de Supabase.

## Requisitos

- Go 1.25.2 o compatible.
- Un proyecto Supabase con acceso a su base de datos PostgreSQL.

## Configuración y ejecución

Desde esta carpeta, copia `.env.example` a `.env`. En Supabase, abre **Connect**, selecciona **Session pooler** y pega la URI PostgreSQL en `SUPABASE_DB_URL`. La conexión exige TLS; la API lo activa si la URI no incluye `sslmode`. Mantén esta URL privada y no la subas al repositorio. Al arrancar, GORM crea o actualiza la tabla `inventory_items` en el proyecto. Si la tabla existente todavía tiene la columna `price`, la migración copia ese valor a los dos nuevos precios y luego elimina la columna antigua.

```powershell
Copy-Item .env.example .env
go mod download
go run .
```

La API escucha en `http://localhost:5665` de forma predeterminada. `CORS_ALLOWED_ORIGINS` acepta una lista separada por comas; si no se configura, permite `http://localhost:3000` y el frontend de producción `https://app-atm.vercel.app`. En producción, configura en el servicio del backend los dominios exactos del frontend. Si defines esta variable, reemplaza la lista predeterminada, así que debe incluir `https://app-atm.vercel.app` (y cualquier otro dominio usado por el frontend), separado por comas. No se habilitan credenciales CORS porque la API no usa cookies.

Configura `NUXT_PUBLIC_API_BASE` en el frontend con la URL base terminada en `/api`. El valor predeterminado local es `http://localhost:5665/api`.

## Rutas

Todas las rutas están bajo `/api`:

| Método | Ruta | Acción |
| --- | --- | --- |
| `GET` | `/inventory` | Listar artículos |
| `GET` | `/inventory/for-sale` | Listar solo nombre, existencias y precio de venta para el selector de productos |
| `POST` | `/inventory` | Crear un artículo |
| `GET` | `/inventory/{id}` | Obtener un artículo |
| `PUT` | `/inventory/{id}` | Reemplazar nombre, cantidad y precio |
| `DELETE` | `/inventory/{id}` | Eliminar un artículo |
| `POST` | `/inventory/consume` | Descontar productos al confirmar un cierre de caja |

El cuerpo de creación y actualización es `{ "name": "Bebida", "quantity": 12, "net_price": 1.5, "sale_price": 2.5 }`. `net_price` es el costo neto unitario y `sale_price` el precio de venta unitario. Los nombres deben tener entre 1 y 120 caracteres; cantidad y ambos precios deben ser mayores o iguales a cero, y los precios aceptan hasta dos decimales. La respuesta incluye el ID UUID del artículo.

El descuento de inventario recibe `{ "checkout_id": "<uuid>", "items": [{ "id": "<uuid>", "quantity": 2 }] }`. Se aplica dentro de una transacción: si falta stock, no se descuenta ningún artículo y responde `409 Conflict`. El identificador de cierre hace idempotentes los reintentos, evitando descontar dos veces.
