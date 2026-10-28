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
  expectedBalances: Record<CategoryId, number>
  movementCounts: Record<CategoryId, number>
  closingBalances: Record<CategoryId, number>
  totalExpected: number
  totalCounted: number
  totalDifference: number
}>()
</script>

<template>
  <section class="grid grid-cols-[1fr_1fr_1.05fr] gap-[15px] max-[780px]:grid-cols-2 max-[780px]:gap-[9px]" aria-label="Resumen de saldos">
    <article v-for="category in CATEGORIES" :key="category.id" class="min-h-[145px] rounded-xl border border-[#eceef4] bg-white p-[17px] shadow-[0_3px_12px_#282b3c05] max-[520px]:min-h-[157px] max-[520px]:p-[12px_11px_11px]" :class="category.id === 'constante' ? 'border-t-[3px] border-t-[#8b78e3]' : 'border-t-[3px] border-t-[#55bea1]'">
      <div class="text-[10px] font-semibold text-[#777b8b] max-[520px]:text-[9px]">{{ category.name }}</div>
      <div class="mt-[3px] text-[25px] font-bold tracking-[-.9px] text-[#292c3d] max-[520px]:text-[21px]"><span class="mr-0.5 text-base font-medium text-[#8c8f9e] max-[520px]:text-[13px]">$</span>{{ formatAmount(expectedBalances[category.id]) }}</div>
      <div class="mt-2 flex justify-between gap-[5px] text-[8px] text-[#989cab] max-[520px]:flex-wrap max-[520px]:text-[7px]">
        <span>Saldo esperado</span>
        <span>{{ movementCounts[category.id] }} {{ movementCounts[category.id] === 1 ? 'movimiento' : 'movimientos' }}</span>
      </div>
      <div v-if="status === 'closed'" class="mt-[7px] flex gap-[10px] text-[8px] text-[#898d9b] max-[520px]:flex-col max-[520px]:gap-[3px]">
        <span class="flex gap-1">Contado <strong class="font-semibold text-[#505466]">${{ formatAmount(closingBalances[category.id]) }}</strong></span>
        <span class="flex gap-1" :class="expectedBalances[category.id] === closingBalances[category.id] ? 'text-[#279978]' : 'text-[#c16b4c]'">
          Diferencia <strong>{{ formatSignedAmount(closingBalances[category.id] - expectedBalances[category.id]) }}</strong>
        </span>
      </div>
    </article>
    <article class="min-h-[167px] rounded-[11px] bg-[#2a2d40] p-[18px_18px_15px] text-white max-[780px]:col-span-2 max-[780px]:min-h-[135px] max-[520px]:min-h-[124px] max-[520px]:p-[14px]">
      <div class="text-[8px] font-bold tracking-[.8px] text-[#c0bed3]">{{ status === 'closed' ? 'TOTAL CONTADO' : 'SALDO TOTAL ESPERADO' }}</div>
      <div class="mt-[13px] text-[27px] font-bold tracking-[-.7px] max-[520px]:text-2xl">${{ formatAmount(status === 'closed' ? totalCounted : totalExpected) }}</div>
      <div class="mt-[10px] flex items-center gap-[6px] text-[9px] text-[#c1c2d0]">
        <span class="size-[6px] rounded-full bg-[#71d5ae]"></span>
        {{ status === 'closed' ? 'Diferencia total' : 'En ambas categorías' }}
        <strong v-if="status === 'closed'" class="font-semibold text-white">{{ formatSignedAmount(totalDifference) }}</strong>
      </div>
    </article>
  </section>
</template>
