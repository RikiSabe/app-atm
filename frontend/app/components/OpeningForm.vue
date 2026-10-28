<script setup lang="ts">
import { CATEGORIES, type CategoryId } from '~/constants/cash'

const props = defineProps<{
  date: string
  openingBalances: Record<CategoryId, number>
  errorMessage: string
}>()

const emit = defineEmits<{
  'update:date': [value: string]
  begin: []
}>()

const dateModel = computed({
  get: () => props.date,
  set: (value: string) => emit('update:date', value)
})
</script>

<template>
  <section class="grid grid-cols-[minmax(0,760px)] justify-center max-[780px]:grid-cols-[minmax(0,1fr)]">
    <article class="rounded-xl border border-[#eceef4] bg-white p-[27px_29px_22px] shadow-[0_3px_12px_#282b3c05] max-[520px]:p-[21px_17px_18px]">
      <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-grid size-[19px] place-items-center rounded-md bg-[#f0edff] text-[9px] tracking-normal text-[#725fd2]">01</span> INICIO DE JORNADA</div>
      <h2 class="mb-[5px] mt-[14px] text-[19px] font-bold tracking-[-.5px] text-[#282b3b] max-[520px]:text-[17px]">Empecemos por el principio</h2>
      <p class="mb-[22px] max-w-[440px] text-[11px] leading-[1.7] text-[#898d9c]">Anota con cuánto dinero inicia cada categoría. Puedes agregar movimientos durante el día.</p>

      <label class="mb-[7px] block text-[10px] font-semibold text-[#55596a]" for="cash-date">Fecha del arqueo</label>
      <div class="flex min-h-10 items-center gap-[9px] rounded-[7px] border border-[#e8eaf0] px-[11px] text-[#9296a7]">
        <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3.5" y="5" width="17" height="15.5" rx="2.5" stroke="currentColor" stroke-width="1.7" /><path d="M7.5 3.5v3M16.5 3.5v3M3.5 9.5h17" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" /></svg>
        <input id="cash-date" v-model="dateModel" class="w-full border-0 bg-transparent text-[11px] text-[#333648] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="date" />
      </div>

      <div class="my-[22px] flex items-center gap-3 text-[8px] font-bold tracking-[.9px] text-[#a1a5b2] before:h-px before:flex-1 before:bg-[#eff0f5] before:content-[''] after:h-px after:flex-1 after:bg-[#eff0f5] after:content-['']"><span>SALDOS INICIALES</span></div>
      <div class="grid grid-cols-2 gap-[14px] max-[520px]:grid-cols-1 max-[520px]:gap-[10px]">
        <label v-for="category in CATEGORIES" :key="category.id" class="block">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">{{ category.name }}</span>
          <span class="flex min-h-10 items-center rounded-[7px] border border-[#e8eaf0] bg-white px-[11px]">
            <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
            <input
              v-model.number="openingBalances[category.id]"
              class="w-full min-w-0 border-0 bg-transparent text-[11px] text-[#303345] outline-none placeholder:text-[#b7bac5] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]"
              type="number"
              min="0"
              step="0.01"
              placeholder="0.00"
              :aria-label="`Saldo inicial de ${category.name}`"
            />
          </span>
        </label>
      </div>

      <p v-if="errorMessage" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ errorMessage }}</p>
      <button type="button" class="mt-[18px] inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] hover:shadow-[0_5px_14px_#7966de40] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="$emit('begin')">
        Abrir caja
        <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M5 12h14m-6-6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
      </button>
      <p class="mt-3 flex items-center justify-center gap-[6px] text-center text-[9px] text-[#9a9eac]"><svg class="size-[13px] text-[#8c80d8]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5 4.5 6.5v5.1c0 4.3 3 7.2 7.5 8.9 4.5-1.7 7.5-4.6 7.5-8.9V6.5L12 3.5Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" /><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" /></svg> Los datos se mantienen solo mientras esta página esté abierta.</p>
    </article>
  </section>
</template>
