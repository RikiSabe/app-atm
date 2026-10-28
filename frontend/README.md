# Frontend

Aplicación web construida con Nuxt.

## Desarrollo

```bash
npm ci
npm run dev
```

## Compilación

```bash
npm run build
```

Para mantener el despliegue automático con Vercel, configura el **Root Directory** del proyecto como `frontend`.

El inventario usa la API Go. En desarrollo se conecta por defecto a `http://localhost:5665/api`; en producción, a `https://app-atm-srcb.vercel.app/api`. Puedes sobrescribir la URL en Vercel con la variable `NUXT_PUBLIC_API_BASE`, incluyendo `/api` al final.

La API Go se conecta a PostgreSQL de Supabase en el propio servicio del backend mediante `SUPABASE_DB_URL`. Usa la cadena PostgreSQL de **Connect → Session pooler** (`postgres://...`), no el dominio HTTPS `*.supabase.co`; mantenla como variable privada del backend y no la configures en el frontend.
