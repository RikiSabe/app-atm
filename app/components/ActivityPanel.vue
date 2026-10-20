<script setup lang="ts">
import { computed } from 'vue'
import {
  CATEGORIES,
  PAYMENT_METHODS,
  categoryName,
  formatAmount,
  formatTime,
  paymentSummary,
  type CashMovement,
  type CashStatus,
  type CategoryId,
  type MovementType,
  type PaymentMethod
} from '~/constants/cash'

const props = defineProps<{
  status: CashStatus
  movements: CashMovement[]
  movementCategory: CategoryId
  movementType: MovementType
  movementPaymentMethod: PaymentMethod
  movementAmount: number | null
  movementConcept: string
  errorMessage: string
}>()

const emit = defineEmits<{
  'update:movementCategory': [value: CategoryId]
  'update:movementType': [value: MovementType]
  'update:movementPaymentMethod': [value: PaymentMethod]
  'update:movementAmount': [value: number | null]
  'update:movementConcept': [value: string]
  addMovement: []
  removeMovement: [id: number]
}>()

const movementCategoryModel = computed({
  get: () => props.movementCategory,
  set: (value: CategoryId) => emit('update:movementCategory', value)
})
const movementTypeModel = computed({
  get: () => props.movementType,
  set: (value: MovementType) => emit('update:movementType', value)
})
const movementPaymentMethodModel = computed({
  get: () => props.movementPaymentMethod,
  set: (value: PaymentMethod) => emit('update:movementPaymentMethod', value)
})
const movementAmountModel = computed({
  get: () => props.movementAmount,
  set: (value: number | null) => emit('update:movementAmount', value)
})
const movementConceptModel = computed({
  get: () => props.movementConcept,
  set: (value: string) => emit('update:movementConcept', value)
})
</script>

<template>
  <div class="grid min-w-0 gap-[14px]">
    <article v-if="status === 'open'" class="rounded-xl border border-[#eceef4] bg-white p-[17px_18px_15px] shadow-[0_3px_12px_#282b3c05] max-[520px]:px-3">
      <div class="flex items-center justify-between gap-[10px]">
        <div>
          <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-grid size-[19px] place-items-center rounded-md bg-[#f0edff] text-[9px] tracking-normal text-[#725fd2]">02</span> MOVIMIENTOS DEL DÍA</div>
          <h2 class="mt-2 text-[15px] font-bold tracking-[-.3px] text-[#282b3b]">Registra un movimiento</h2>
        </div>
      </div>
      <form class="mt-[17px] grid grid-cols-[1fr_.8fr_1fr_.85fr_1.3fr_auto] items-end gap-2 max-[1100px]:grid-cols-3 max-[520px]:grid-cols-2" @submit.prevent="$emit('addMovement')">
        <label class="min-w-0 max-[1100px]:col-span-1">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Categoría</span>
          <select v-model="movementCategoryModel" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
            <option v-for="category in CATEGORIES" :key="category.id" :value="category.id">{{ category.name }}</option>
          </select>
        </label>
        <label class="min-w-0">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Método de pago</span>
          <select v-model="movementPaymentMethodModel" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
            <option v-for="paymentMethod in PAYMENT_METHODS" :key="paymentMethod" :value="paymentMethod">{{ paymentMethod }}</option>
          </select>
        </label>
        <label class="min-w-0">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Tipo</span>
          <select v-model="movementTypeModel" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
            <option value="entrada">Entrada</option>
            <option value="salida">Salida</option>
          </select>
        </label>
        <label class="min-w-0">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Importe</span>
          <span class="flex min-h-[38px] items-center rounded-[7px] border border-[#e8eaf0] bg-white px-2 focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
            <span class="mr-[5px] text-[10px] text-[#a3a7b5]">$</span>
            <input v-model.number="movementAmountModel" class="w-full min-w-0 border-0 bg-transparent text-[10px] text-[#303345] outline-none placeholder:text-[#b7bac5] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0.01" step="0.01" placeholder="0.00" />
          </span>
        </label>
        <label class="min-w-0 max-[1100px]:col-span-2 max-[520px]:col-span-2">
          <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Concepto</span>
          <input v-model="movementConceptModel" class="block h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none placeholder:text-[9px] placeholder:text-[#b7bac5] focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="text" maxlength="80" placeholder="¿De qué se trata?" />
        </label>
        <button class="inline-flex min-h-[38px] cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[11px] text-[9px] font-semibold text-white transition hover:-translate-y-px hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 max-[1100px]:col-start-3 max-[520px]:col-span-2 max-[520px]:col-start-auto" type="submit" aria-label="Agregar movimiento">
          <svg class="size-[14px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 5v14m-7-7h14" stroke="currentColor" stroke-width="2" stroke-linecap="round" /></svg>
          <span>Agregar</span>
        </button>
      </form>
      <p v-if="errorMessage" class="mb-0 mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ errorMessage }}</p>
    </article>

    <article class="overflow-hidden rounded-xl border border-[#eceef4] bg-white px-[18px] pt-[17px] shadow-[0_3px_12px_#282b3c05] max-[520px]:px-3">
      <div class="flex items-center justify-between gap-[10px]">
        <div>
          <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]">REGISTRO</div>
          <h2 class="mt-2 text-[15px] font-bold tracking-[-.3px] text-[#282b3b]">{{ status === 'closed' ? 'Resumen de movimientos' : 'Actividad de hoy' }}</h2>
        </div>
        <span class="rounded-[14px] border border-[#eeeff4] bg-white px-2 py-[5px] text-[8px] text-[#828697]">{{ movements.length }} {{ movements.length === 1 ? 'registro' : 'registros' }}</span>
      </div>

      <div v-if="movements.length" class="mt-[14px] -mx-[18px] max-[520px]:-mx-3">
        <div v-for="movement in movements" :key="movement.id" class="flex min-h-[53px] items-center gap-[10px] border-t border-[#f0f1f5] px-[18px] py-2 max-[520px]:gap-[6px] max-[520px]:px-[10px]">
          <span class="grid size-[29px] flex-none place-items-center rounded-lg" :class="movement.type === 'entrada' ? 'bg-[#eaf7f2] text-[#299c7e]' : 'bg-[#fff3ec] text-[#d1845d]'">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path :d="movement.type === 'entrada' ? 'M12 19V5m-7 7 7-7 7 7' : 'M12 5v14m-7-7 7 7 7-7'" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
          </span>
          <span class="grid min-w-0 flex-1 gap-[3px]">
            <strong class="overflow-hidden text-ellipsis whitespace-nowrap text-[10px] font-semibold text-[#3f4253]">{{ movement.concept }}</strong>
            <small class="text-[8px] text-[#999dac]">{{ categoryName(movement.category) }} <span class="px-0.5">·</span> {{ movement.paymentMethod }} <span class="px-0.5">·</span> {{ formatTime(movement.time) }}</small>
          </span>
          <span class="min-w-10 text-center text-[8px] capitalize max-[520px]:hidden" :class="movement.type === 'entrada' ? 'text-[#379c81]' : 'text-[#bf805b]'">{{ movement.type }}</span>
          <strong class="min-w-[84px] whitespace-nowrap text-right text-[10px] text-[#515466] max-[520px]:min-w-[68px] max-[520px]:text-[9px]" :class="movement.type === 'entrada' ? 'text-[#358f77]' : 'text-[#bf7956]'">{{ movement.type === 'entrada' ? '+' : '−' }}${{ formatAmount(movement.amount) }}</strong>
          <button v-if="status === 'open'" type="button" class="grid size-[25px] flex-none cursor-pointer place-items-center rounded-md border-0 bg-transparent text-[#a1a5b3] transition hover:bg-[#fff0f1] hover:text-[#c34e59] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] max-[520px]:size-[22px]" :aria-label="`Eliminar movimiento ${movement.concept}`" @click="$emit('removeMovement', movement.id)">
            <svg class="size-[15px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 7h16m-10 4v6m4-6v6M6 7l1 13h10l1-13M9 7V4h6v3" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" /></svg>
          </button>
        </div>
      </div>
      <div class="border-t border-[#f0f1f5] px-0 py-[13px]">
        <h3 class="mb-[9px] text-[11px] font-bold text-[#3f4253]">Resumen por método de pago</h3>
        <div class="overflow-x-auto">
          <table class="w-full min-w-[430px] border-collapse text-left text-[9px]">
            <thead class="text-[#85899a]">
              <tr>
                <th class="border-b border-[#eceef4] px-2 py-[7px] font-semibold">Método</th>
                <th class="border-b border-[#eceef4] px-2 py-[7px] text-right font-semibold">Movimientos</th>
                <th class="border-b border-[#eceef4] px-2 py-[7px] text-right font-semibold">Entradas</th>
                <th class="border-b border-[#eceef4] px-2 py-[7px] text-right font-semibold">Salidas</th>
                <th class="border-b border-[#eceef4] px-2 py-[7px] text-right font-semibold">Neto</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="paymentMethod in PAYMENT_METHODS" :key="paymentMethod" class="text-[#55596a]">
                <th class="border-b border-[#f0f1f5] px-2 py-[7px] font-semibold">{{ paymentMethod }}</th>
                <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">{{ paymentSummary(movements, paymentMethod).count }}</td>
                <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">${{ formatAmount(paymentSummary(movements, paymentMethod).entries) }}</td>
                <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">${{ formatAmount(paymentSummary(movements, paymentMethod).exits) }}</td>
                <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right font-semibold">${{ formatAmount(paymentSummary(movements, paymentMethod).net) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div v-if="!movements.length" class="flex min-h-[142px] flex-col items-center justify-center border-t border-[#f0f1f5] text-center">
        <span class="mb-[9px] grid size-[34px] place-items-center rounded-[10px] bg-[#f3f1ff] text-[#8d7be0]">
          <svg class="size-[18px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M5 5.5h14v15H5zM8 9h8m-8 4h8m-8 4h4M8 3.5v4m8-4v4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
        </span>
        <strong class="text-[10px] text-[#55596a]">Aún no hay movimientos</strong>
        <p class="mb-0 mt-1 text-[9px] text-[#999dac]">Cuando registres entradas o salidas, aparecerán aquí.</p>
      </div>
    </article>
  </div>
</template>
