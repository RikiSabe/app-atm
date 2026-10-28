# Monorepo

Este repositorio contiene el frontend Nuxt y la API de inventario en Go:

- `frontend/`: aplicacion web Nuxt, desplegada por Vercel.
- `api/`: API REST en Go para administrar el inventario con PostgreSQL.

## Desarrollo del frontend

```bash
cd frontend
npm ci
npm run dev
```

Para compilar localmente:

```bash
cd frontend
npm run build
```

En Vercel, configura **Root Directory** como `frontend`. El despliegue automatico desde GitHub puede mantenerse habilitado sin agregar credenciales de despliegue al workflow.

## API de inventario

Consulta [`api/README.md`](./api/README.md) para conectar la API con Supabase, iniciarla y definir `NUXT_PUBLIC_API_BASE` en el frontend. El arqueo de caja continúa funcionando en memoria; el backend persiste únicamente los artículos del inventario.

## GitHub Actions

Los workflows de GitHub Actions compilan el frontend y ejecutan pruebas/vet del backend cuando cambian sus respectivos directorios. No despliegan la aplicacion; Vercel sigue encargandose del despliegue del frontend.
