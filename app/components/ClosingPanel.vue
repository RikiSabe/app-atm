<script setup lang="ts">
import {
  CATEGORIES,
  formatAmount,
  formatSignedAmount,
  type CashStatus,
  type CategoryId
} from '~/constants/cash'

defineProps<{
  status: CashStatus
  closingBalances: Record<CategoryId, number>
  expectedBalances: Record<CategoryId, number>
  totalExpected: number
  totalDifference: number
  errorMessage: string
  downloadError: string
}>()

defineEmits<{
  startClosing: []
  finishClosing: []
  backToMovements: []
  download: []
}>()
</script>

<template>
  <aside class="grid min-w-0 gap-[14px]">
    <article v-if="status === 'open'" class="rounded-xl border border-[#eceef4] bg-white p-[17px] shadow-[0_3px_12px_#282b3c05]">
      <div class="mb-[13px] grid size-[30px] place-items-center rounded-[9px] bg-[#f0edff] text-[#7462d4]">
        <svg class="size-[17px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 12h16m-6-6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </div>
      <div class="text-[8px] font-bold tracking-[1.05px] text-[#9296a7]">FIN DE JORNADA</div>
      <h3 class="mb-[5px] mt-2 text-[15px] font-bold tracking-[-.2px] text-[#313446]">¿Listo para cerrar?</h3>
      <p class="m-0 text-[9px] leading-[1.7] text-[#8c90a0]">Controla el efectivo de cada categoría y compáralo con el saldo esperado.</p>
      <div class="my-[15px] flex items-center justify-between gap-2 rounded-[7px] border border-[#eff0f5] bg-[#fafafd] px-[10px] py-[11px] text-[9px] text-[#878b9b]">
        <span>Saldo esperado</span>
        <strong class="text-[11px] text-[#393c4e]">${{ formatAmount(totalExpected) }}</strong>
      </div>
      <button type="button" class="inline-flex min-h-[38px] w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#282b3d] px-[15px] text-[9px] font-semibold text-white transition hover:-translate-y-px hover:bg-[#393c50] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="$emit('startClosing')">
        Hacer cierre de caja
        <svg class="size-[14px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M5 12h14m-6-6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </button>
      <div class="mt-3 flex items-start gap-[6px] text-[8px] leading-[1.5] text-[#999dac]"><span class="text-[#8775df]">✳</span> Podrás revisar la diferencia antes de exportar.</div>
    </article>

    <article v-else-if="status === 'closing'" class="rounded-xl border border-[#eceef4] bg-white p-[17px] shadow-[0_3px_12px_#282b3c05]">
      <div class="mb-[13px] grid size-[30px] place-items-center rounded-[9px] bg-[#f0edff] text-[#7462d4]">
        <svg class="size-[17px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 12h16m-6-6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </div>
      <div class="text-[8px] font-bold tracking-[1.05px] text-[#9296a7]">CONFIRMA EL CONTEO</div>
      <h3 class="mb-[5px] mt-2 text-[15px] font-bold tracking-[-.2px] text-[#313446]">¿Cuánto hay en caja?</h3>
      <p class="m-0 text-[9px] leading-[1.7] text-[#8c90a0]">Ingresa el efectivo que contaste en cada categoría para ver las diferencias.</p>
      <label v-for="category in CATEGORIES" :key="category.id" class="mt-[13px] block">
        <span class="mb-[6px] flex justify-between gap-[5px] text-[9px] font-semibold text-[#535768]">{{ category.name }} <small class="whitespace-nowrap text-[8px] font-normal text-[#9397a6]">Esperado: ${{ formatAmount(expectedBalances[category.id]) }}</small></span>
        <span class="flex min-h-9 items-center rounded-[7px] border border-[#e8eaf0] bg-white px-[11px] focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
          <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
          <input v-model.number="closingBalances[category.id]" class="w-full min-w-0 border-0 bg-transparent text-[10px] text-[#303345] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="0.01" :aria-label="`Importe contado de ${category.name}`" />
        </span>
        <small class="mt-1 block text-[8px]" :class="closingBalances[category.id] - expectedBalances[category.id] === 0 ? 'text-[#279978]' : 'text-[#c16b4c]'">
          Diferencia: {{ formatSignedAmount(closingBalances[category.id] - expectedBalances[category.id]) }}
        </small>
      </label>
      <p v-if="errorMessage" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ errorMessage }}</p>
      <button type="button" class="mt-[13px] inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="$emit('finishClosing')">
        Confirmar cierre
        <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </button>
      <button type="button" class="mx-auto mt-[10px] block cursor-pointer border-0 bg-transparent px-[6px] py-[3px] text-[9px] text-[#777b8b] underline decoration-1 underline-offset-2 hover:text-[#6552cc] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" @click="$emit('backToMovements')">Volver a los movimientos</button>
    </article>

    <article v-else class="rounded-xl border border-[#e4f0eb] bg-gradient-to-br from-white via-white to-[#f4fbf8] p-[17px] shadow-[0_3px_12px_#282b3c05]">
      <div class="mb-[13px] grid size-[30px] place-items-center rounded-full bg-[#e5f6ef] text-[#258c70]">
        <svg class="size-[17px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </div>
      <div class="text-[8px] font-bold tracking-[1.05px] text-[#9296a7]">CIERRE COMPLETADO</div>
      <h3 class="mb-[5px] mt-2 text-[15px] font-bold tracking-[-.2px] text-[#313446]">El día está cerrado.</h3>
      <p class="m-0 text-[9px] leading-[1.7] text-[#8c90a0]">Descarga el resumen del arqueo en PDF para conservarlo.</p>
      <div class="my-[15px] flex items-center justify-between gap-2 rounded-[7px] border border-[#eff0f5] bg-[#fafafd] px-[10px] py-[11px] text-[9px] text-[#878b9b]">
        <span>Diferencia total</span>
        <strong :class="totalDifference === 0 ? 'text-[#279978]' : 'text-[#c16b4c]'">{{ formatSignedAmount(totalDifference) }}</strong>
      </div>
      <button type="button" class="inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="$emit('download')">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5v11m0 0 4-4m-4 4-4-4M4.5 16v4h15v-4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
        Descargar PDF
      </button>
      <p v-if="downloadError" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ downloadError }}</p>
      <div v-else class="mt-3 flex items-start gap-[6px] text-[8px] leading-[1.5] text-[#999dac]"><span class="text-[#8775df]">✳</span> En iPhone, elige “Guardar en Archivos” en el menú que aparece.</div>
    </article>

    <div class="text-[9px] text-[#9296a4]">
      Los datos se borran al recargar o cerrar esta página.
    </div>
  </aside>
</template>
