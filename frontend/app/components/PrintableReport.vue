<script setup lang="ts">
import {
  CATEGORIES,
  PAYMENT_METHODS,
  formatAmount,
  formatDateTime,
  formatSignedAmount,
  movementDetails,
  paymentSummary,
  type CashMovement,
  type CategoryId
} from '~/constants/cash'

const props = defineProps<{
  todayLabel: string
  movements: CashMovement[]
  openingBalances: Record<CategoryId, number>
  closingBalances: Record<CategoryId, number>
  expectedBalances: Record<CategoryId, number>
  totalExpected: number
  totalCounted: number
  totalDifference: number
}>()

const detailedMovements = computed(() =>
  movementDetails(props.movements, props.openingBalances)
)

function categoryMovements(categoryId: CategoryId) {
  return props.movements.filter(movement => movement.category === categoryId)
}

function movementTotal(categoryId: CategoryId, type: CashMovement['type']) {
  return categoryMovements(categoryId)
    .filter(movement => movement.type === type)
    .reduce((total, movement) => total + movement.amount, 0)
}

function itemDetails(categoryId: CategoryId) {
  const details = categoryMovements(categoryId).map(movement => movement.concept)
  return details.length ? details.join(', ') : 'Sin movimientos'
}
</script>

<template>
  <section class="hidden p-0 font-sans text-[#252637] print:block" aria-label="Reporte de arqueo">
    <div class="flex items-center justify-between border-b-2 border-[#7661dc] pb-[13px] text-[10px] font-bold tracking-[.7px] text-[#76798a]">
      <div class="text-lg font-extrabold tracking-[-.5px] text-[#292c3d]">arqueo</div>
      <span>REPORTE DE ARQUEO DIARIO</span>
    </div>
    <h1 class="mb-1 mt-[22px] text-2xl font-bold">Arqueo diario</h1>
    <p class="mb-5 text-[11px] capitalize text-[#737789]">{{ todayLabel }}</p>
    <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a] [&_tfoot_th]:bg-[#faf9fd]">
      <thead><tr><th>Categoría</th><th>Saldo inicial</th><th>Entradas</th><th>Salidas</th><th>Esperado</th><th>Contado</th><th>Diferencia</th></tr></thead>
      <tbody>
        <tr v-for="category in CATEGORIES" :key="category.id">
          <td>{{ itemDetails(category.id) }}</td>
          <td>${{ formatAmount(openingBalances[category.id]) }}</td>
          <td>${{ formatAmount(movementTotal(category.id, 'entrada')) }}</td>
          <td>${{ formatAmount(movementTotal(category.id, 'salida')) }}</td>
          <td>${{ formatAmount(expectedBalances[category.id]) }}</td>
          <td>${{ formatAmount(closingBalances[category.id]) }}</td>
          <td>{{ formatSignedAmount(closingBalances[category.id] - expectedBalances[category.id]) }}</td>
        </tr>
      </tbody>
      <tfoot><tr><th>Total</th><th colspan="3"></th><th>${{ formatAmount(totalExpected) }}</th><th>${{ formatAmount(totalCounted) }}</th><th>{{ formatSignedAmount(totalDifference) }}</th></tr></tfoot>
    </table>
    <h2 class="mb-[9px] mt-[23px] text-sm font-bold">Resumen por método de pago</h2>
    <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a]">
      <thead><tr><th>Método</th><th>Movimientos</th><th>Entradas</th><th>Salidas</th><th>Neto</th></tr></thead>
      <tbody>
        <tr v-for="paymentMethod in PAYMENT_METHODS" :key="paymentMethod">
          <td>{{ paymentMethod }}</td>
          <td>{{ paymentSummary(movements, paymentMethod).count }}</td>
          <td>${{ formatAmount(paymentSummary(movements, paymentMethod).entries) }}</td>
          <td>${{ formatAmount(paymentSummary(movements, paymentMethod).exits) }}</td>
          <td>${{ formatAmount(paymentSummary(movements, paymentMethod).net) }}</td>
        </tr>
      </tbody>
    </table>
    <h2 class="mb-[9px] mt-[23px] text-sm font-bold">Detalle individual de movimientos</h2>
    <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a] [&_tfoot_th]:bg-[#faf9fd]">
      <thead><tr><th>N.º</th><th>Fecha y hora</th><th>Detalle del ítem</th><th>Tipo</th><th>Método de pago</th><th>Importe</th><th>Saldo después</th></tr></thead>
      <tbody>
        <tr v-for="movement in detailedMovements" :key="movement.id">
          <td>{{ movement.number }}</td><td>{{ formatDateTime(movement.time) }}</td><td>{{ movement.concept }}</td><td>{{ movement.type }}</td><td>{{ movement.paymentMethod }}</td><td>{{ formatSignedAmount(movement.type === 'entrada' ? movement.amount : -movement.amount) }}</td><td>${{ formatAmount(movement.balanceAfter) }}</td>
        </tr>
        <tr v-if="!movements.length"><td colspan="8">No se registraron movimientos.</td></tr>
      </tbody>
    </table>
    <p class="mt-[25px] text-[8px] text-[#858999]">Reporte generado en el dispositivo. Los datos no se almacenan en la aplicación.</p>
  </section>
</template>
