<script setup lang="ts">
import type { CashStatus } from '~/constants/cash'

defineProps<{ status: CashStatus }>()

const statusLabels: Record<CashStatus, string> = {
  opening: 'Por iniciar',
  open: 'Caja abierta',
  closing: 'En cierre',
  closed: 'Caja cerrada'
}
</script>

<template>
  <section class="mb-[25px] flex items-end justify-between gap-[18px] max-[520px]:items-start">
    <div>
      <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-block h-0.5 w-[17px] rounded-sm bg-[#8875e5]"></span> CONTROL DE CAJA</div>
      <h1 class="my-[9px] mb-1 text-[clamp(29px,3vw,37px)] font-bold leading-[1.15] tracking-[-1.7px] text-[#232637] max-[520px]:text-[29px]">Arqueo diario<span class="text-[#7967de]">.</span></h1>
      <p class="m-0 text-xs text-[#8a8e9e] max-[520px]:max-w-[235px] max-[520px]:text-[10px] max-[520px]:leading-[1.5]">Cada movimiento en su lugar. Tu control del efectivo, siempre ordenado.</p>
    </div>
    <div class="inline-flex items-center gap-[7px] whitespace-nowrap rounded-full border border-[#e6e8ef] bg-white px-[10px] py-[7px] text-[10px] font-semibold text-[#777c8d] max-[520px]:mt-[19px] max-[520px]:px-2 max-[520px]:py-[6px] max-[520px]:text-[8px]">
      <span class="inline-block size-[6px] rounded-full" :class="status === 'opening' ? 'bg-[#e5aa53]' : status === 'closing' ? 'bg-[#8b77e4]' : 'bg-[#29ac84]'"></span>
      {{ statusLabels[status] }}
    </div>
  </section>
</template>
