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

El inventario usa la API Go de `../api`. En desarrollo se conecta por defecto a `http://localhost:5665/api`; configura `NUXT_PUBLIC_API_BASE` para usar otra URL.
