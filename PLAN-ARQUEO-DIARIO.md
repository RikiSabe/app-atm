# Plan: arqueo diario de caja

## Objetivo

Crear una aplicación web pequeña, completamente en el frontend, para abrir, controlar y cerrar el arqueo diario de dos categorías: **Activo constante** y **Activo variable**.

## Alcance funcional

1. **Apertura**
   - Elegir la fecha del arqueo.
   - Ingresar el saldo inicial de cada categoría.
   - Validar que la fecha exista y los importes iniciales no sean negativos.
2. **Movimientos**
   - Registrar entradas y salidas por categoría con importe, concepto y método de pago (QR o Efectivo).
   - Mostrar fecha/hora, categoría, tipo, método de pago, concepto e importe.
   - Permitir eliminar movimientos mientras la caja siga abierta.
   - Calcular por categoría el saldo esperado: saldo inicial + entradas − salidas.
3. **Cierre**
   - Solicitar el efectivo contado para cada categoría.
   - Comparar lo contado con el saldo esperado y mostrar diferencias por categoría y total.
   - Permitir volver al registro de movimientos antes de confirmar el cierre.
4. **Reporte**
   - Mostrar un resumen por método de pago con cantidad de movimientos, entradas, salidas y neto.
   - Generar un PDF descargable en formato vertical, con ingresos y egresos en secciones separadas, concepto, método de pago, fecha y hora por movimiento, y sus totales.
   - Crear el PDF como Blob en el navegador y descargarlo directamente al dispositivo.
5. **Privacidad**
   - Mantener el estado solo en memoria durante la sesión de la página.
   - No usar base de datos, backend, `localStorage` ni `sessionStorage`.
   - Advertir que recargar o cerrar la página elimina los datos.

## Implementación

- Usar Nuxt y Vue ya presentes en el proyecto, junto con jsPDF y jsPDF-AutoTable para generar el reporte.
- Construir la experiencia en `app/app.vue`, con diseño adaptable a escritorio y móvil.
- Mantener categorías, apertura, movimientos y cierre en estado reactivo local.
- Generar y descargar el PDF en el navegador desde un Blob, sin depender del diálogo de impresión.

## Criterios de aceptación

- Se puede abrir un arqueo para las dos categorías y añadir movimientos.
- Las entradas aumentan y las salidas disminuyen correctamente el saldo esperado.
- El cierre solicita importes contados y calcula diferencias positivas, negativas y cero.
- Se puede descargar un PDF con el detalle del arqueo cerrado.
- No se persisten datos al recargar la aplicación ni se transmiten a un backend.
