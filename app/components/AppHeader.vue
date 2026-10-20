<script setup lang="ts">
import type { CashStatus } from '~/constants/cash'

defineProps<{
  todayLabel: string
  status: CashStatus
}>()

defineEmits<{
  download: []
  newCashCount: []
}>()
</script>

<template>
  <header class="flex min-h-[68px] items-center justify-between gap-3 border-b border-[#edeff5] bg-white px-[max(4.5%,calc((100%-1120px)/2))] max-[780px]:min-h-[56px] max-[780px]:px-[4%] max-[520px]:gap-[5px] print:hidden">
    <a class="inline-flex w-fit items-center gap-[10px] text-[#222637] no-underline focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" href="#arqueo" aria-label="Arqueo, inicio">
      <span class="grid size-9 place-items-center rounded-[11px] bg-[#7965df] text-white shadow-[0_5px_12px_#7965df38]">
        <svg class="size-[23px]" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <rect x="3.5" y="6" width="17" height="14" rx="3" stroke="currentColor" stroke-width="1.8" />
          <path d="M7.5 6V4.8A1.8 1.8 0 0 1 9.3 3h8.2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
          <circle cx="12" cy="13" r="2.3" stroke="currentColor" stroke-width="1.8" />
        </svg>
      </span>
      <span class="font-sans text-[17px] font-extrabold tracking-[-.75px] text-[#222637]">arqueo</span>
    </a>
    <div class="flex items-center gap-[9px] max-[520px]:gap-[5px]">
      <span class="inline-flex items-center gap-[7px] text-[10px] capitalize text-[#73788a] max-[520px]:hidden">
        <svg class="size-[15px] text-[#9b9fb0]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3.5" y="5" width="17" height="15.5" rx="2.5" stroke="currentColor" stroke-width="1.7" /><path d="M7.5 3.5v3M16.5 3.5v3M3.5 9.5h17" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" /></svg>
        {{ todayLabel }}
      </span>
      <button
        v-if="status === 'closed'"
        type="button"
        class="inline-flex min-h-8 cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#282b3d] px-[11px] text-[10px] font-semibold text-white transition hover:bg-[#393c50] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 max-[520px]:min-h-[30px] max-[520px]:px-2 max-[520px]:text-[9px]"
        @click="$emit('download')"
      >
        <svg class="size-[14px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5v11m0 0 4-4m-4 4-4-4M4.5 16v4h15v-4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
        Descargar PDF
      </button>
      <button v-if="status === 'closed'" type="button" class="inline-flex min-h-8 cursor-pointer items-center justify-center gap-2 rounded-[7px] border border-[#e6e8ef] bg-white px-[11px] text-[10px] font-semibold text-[#565a6c] transition hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 max-[520px]:min-h-[30px] max-[520px]:px-2 max-[520px]:text-[9px]" @click="$emit('newCashCount')">
        Nuevo arqueo
      </button>
    </div>
  </header>
</template>
