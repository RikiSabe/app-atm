<script setup lang="ts">
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'
import {
  CATEGORIES,
  formatAmount,
  formatDateTime,
  formatSignedAmount,
  localDateValue,
  type CashMovement,
  type CashStatus,
  type CategoryId,
  type PaymentMethod
} from '~/constants/cash'

const activeSection = ref<'cash' | 'inventory'>('cash')
const status = ref<CashStatus>('opening')
const date = ref(localDateValue())
const openingBalances = reactive<Record<CategoryId, number>>({
  constante: 0,
  variable: 0
})
const closingBalances = reactive<Record<CategoryId, number>>({
  constante: 0,
  variable: 0
})
const movements = ref<CashMovement[]>([])
const movementCategory = ref<CategoryId>('variable')
const movementType = ref<'entrada' | 'salida'>('entrada')
const movementPaymentMethod = ref<PaymentMethod>('Efectivo')
const movementAmount = ref<number | null>(null)
const movementConcept = ref('')
const keepMovementContext = ref(false)
const errorMessage = ref('')
const downloadError = ref('')
const toastMessage = ref('')
let toastTimeout: ReturnType<typeof setTimeout> | undefined

const todayLabel = computed(() =>
  new Date(`${date.value}T12:00:00`).toLocaleDateString('es', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
)

function expectedBalance(categoryId: CategoryId) {
  return movements.value
    .filter(movement => movement.category === categoryId)
    .reduce((balance, movement) => {
      return balance + (movement.type === 'entrada' ? movement.amount : -movement.amount)
    }, openingBalances[categoryId])
}

function movementsFor(categoryId: CategoryId) {
  return movements.value.filter(movement => movement.category === categoryId)
}

const expectedBalances = computed<Record<CategoryId, number>>(() => ({
  constante: expectedBalance('constante'),
  variable: expectedBalance('variable')
}))
const movementCounts = computed<Record<CategoryId, number>>(() => ({
  constante: movementsFor('constante').length,
  variable: movementsFor('variable').length
}))
const totalExpected = computed(() =>
  CATEGORIES.reduce((total, category) => total + expectedBalance(category.id), 0)
)
const totalCounted = computed(() =>
  CATEGORIES.reduce((total, category) => total + closingBalances[category.id], 0)
)
const totalDifference = computed(() => totalCounted.value - totalExpected.value)

function showToast(message: string) {
  toastMessage.value = message
  if (toastTimeout) {
    clearTimeout(toastTimeout)
  }
  toastTimeout = setTimeout(() => {
    toastMessage.value = ''
    toastTimeout = undefined
  }, 2800)
}

onUnmounted(() => {
  if (toastTimeout) {
    clearTimeout(toastTimeout)
  }
})

function beginDay() {
  errorMessage.value = ''
  if (!date.value) {
    errorMessage.value = 'Selecciona la fecha del arqueo para continuar.'
    return
  }
  if (CATEGORIES.some(category => openingBalances[category.id] < 0)) {
    errorMessage.value = 'Los saldos iniciales no pueden ser negativos.'
    return
  }
  status.value = 'open'
}

function addMovement() {
  errorMessage.value = ''
  const amount = Number(movementAmount.value)
  if (!movementConcept.value.trim()) {
    errorMessage.value = 'Escribe un concepto para identificar el movimiento.'
    return
  }
  if (!Number.isFinite(amount) || amount <= 0) {
    errorMessage.value = 'El importe debe ser mayor que cero.'
    return
  }

  movements.value.unshift({
    id: Date.now(),
    category: movementCategory.value,
    type: movementType.value,
    paymentMethod: movementPaymentMethod.value,
    amount,
    concept: movementConcept.value.trim(),
    time: new Date().toISOString()
  })
  movementAmount.value = null
  if (!keepMovementContext.value) {
    movementConcept.value = ''
  }
  showToast('Movimiento agregado')
}

function removeMovement(id: number) {
  const movementExists = movements.value.some(movement => movement.id === id)
  movements.value = movements.value.filter(movement => movement.id !== id)
  if (movementExists) {
    showToast('Movimiento eliminado')
  }
}

function startClosing() {
  errorMessage.value = ''
  for (const category of CATEGORIES) {
    closingBalances[category.id] = expectedBalance(category.id)
  }
  status.value = 'closing'
}

function finishClosing() {
  errorMessage.value = ''
  if (CATEGORIES.some(category => !Number.isFinite(closingBalances[category.id]) || closingBalances[category.id] < 0)) {
    errorMessage.value = 'Los importes contados deben ser válidos y no pueden ser negativos.'
    return
  }
  status.value = 'closed'
}

function newCashCount() {
  if (!window.confirm('¿Iniciar un nuevo arqueo? El arqueo actual se eliminará de esta sesión.')) {
    return
  }
  status.value = 'opening'
  date.value = localDateValue()
  for (const category of CATEGORIES) {
    openingBalances[category.id] = 0
    closingBalances[category.id] = 0
  }
  movements.value = []
  errorMessage.value = ''
}

function lastPdfTableY(pdf: jsPDF) {
  return (pdf as jsPDF & { lastAutoTable: { finalY: number } }).lastAutoTable.finalY
}

function formatPdfSignedAmount(amount: number) {
  return formatSignedAmount(amount).replace('−', '-')
}

async function downloadReport() {
  downloadError.value = ''

  try {
    const pdf = new jsPDF({ orientation: 'portrait', unit: 'mm', format: 'a4' })
    const entries = movements.value.filter(movement => movement.type === 'entrada')
    const exits = movements.value.filter(movement => movement.type === 'salida')
    const totalEntries = entries.reduce((total, movement) => total + movement.amount, 0)
    const totalExits = exits.reduce((total, movement) => total + movement.amount, 0)

    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(18)
    pdf.text('Arqueo diario', 14, 18)
    pdf.setFont('helvetica', 'normal')
    pdf.setFontSize(10)
    pdf.text(todayLabel.value, 14, 25)

    function addMovementSection(title: string, sectionMovements: CashMovement[], total: number, startY: number) {
      pdf.setFont('helvetica', 'bold')
      pdf.setFontSize(13)
      pdf.text(title, 14, startY)

      autoTable(pdf, {
        head: [['Concepto - método de pago', 'Fecha y hora', 'Importe']],
        body: sectionMovements.length
          ? sectionMovements.map(movement => [
              `${movement.concept} - (${movement.paymentMethod})`,
              formatDateTime(movement.time),
              `$${formatAmount(movement.amount)}`
            ])
          : [[`No se registraron ${title.toLowerCase()}.`, '', '']],
        foot: [[`Total ${title.toLowerCase()}`, '', `$${formatAmount(total)}`]],
        startY: startY + 4,
        theme: 'grid',
        styles: { font: 'helvetica', fontSize: 9, cellPadding: 3, overflow: 'linebreak' },
        columnStyles: {
          0: { cellWidth: 'auto' },
          1: { cellWidth: 48 },
          2: { cellWidth: 32, halign: 'right' }
        },
        headStyles: { fillColor: [118, 97, 220], textColor: 255 },
        footStyles: { fillColor: [245, 244, 252], textColor: [45, 47, 63], fontStyle: 'bold' },
        margin: { left: 14, right: 14 }
      })
    }

    addMovementSection('Ingresos', entries, totalEntries, 36)

    let nextSectionY = lastPdfTableY(pdf) + 12
    if (nextSectionY > 245) {
      pdf.addPage()
      nextSectionY = 20
    }
    addMovementSection('Egresos', exits, totalExits, nextSectionY)

    let totalsY = lastPdfTableY(pdf) + 12
    if (totalsY > 245) {
      pdf.addPage()
      totalsY = 20
    }
    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(13)
    pdf.text('Totales', 14, totalsY)
    autoTable(pdf, {
      body: [
        ['Total ingresos', `$${formatAmount(totalEntries)}`],
        ['Total egresos', `$${formatAmount(totalExits)}`],
        ['Neto del día', formatPdfSignedAmount(totalEntries - totalExits)]
      ],
      startY: totalsY + 4,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 9, cellPadding: 3 },
      columnStyles: { 1: { halign: 'right' } },
      margin: { left: 14, right: 14 }
    })

    const blob = pdf.output('blob')
    const filename = `arqueo-${date.value}.pdf`
    const objectUrl = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = objectUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.setTimeout(() => URL.revokeObjectURL(objectUrl), 1000)
  } catch (error) {
    if (error instanceof Error && error.name === 'AbortError') {
      return
    }

    downloadError.value = error instanceof Error
      ? `No se pudo descargar el PDF: ${error.message}`
      : 'No se pudo descargar el PDF por un error desconocido.'
  }
}
</script>

<template>
  <div class="min-h-screen min-w-[320px] bg-[#f7f8fc] font-sans text-[#222637] antialiased [font-synthesis:none] [text-rendering:optimizeLegibility] print:bg-white">
    <main id="arqueo" class="min-h-screen print:m-0">
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
        <nav class="flex items-center gap-1 rounded-[9px] bg-[#f5f4fa] p-1" aria-label="Secciones principales">
          <button
            type="button"
            class="min-h-8 cursor-pointer rounded-[7px] border-0 px-3 text-[10px] font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] max-[520px]:px-2 max-[520px]:text-[9px]"
            :class="activeSection === 'cash' ? 'bg-white text-[#5c4bc0] shadow-sm' : 'bg-transparent text-[#777b8b] hover:text-[#393c4e]'"
            :aria-current="activeSection === 'cash' ? 'page' : undefined"
            @click="activeSection = 'cash'"
          >
            Arqueo
          </button>
          <button
            type="button"
            class="min-h-8 cursor-pointer rounded-[7px] border-0 px-3 text-[10px] font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] max-[520px]:px-2 max-[520px]:text-[9px]"
            :class="activeSection === 'inventory' ? 'bg-white text-[#5c4bc0] shadow-sm' : 'bg-transparent text-[#777b8b] hover:text-[#393c4e]'"
            :aria-current="activeSection === 'inventory' ? 'page' : undefined"
            @click="activeSection = 'inventory'"
          >
            Inventario
          </button>
        </nav>
        <div class="flex items-center gap-[9px] max-[520px]:gap-[5px]">
          <span class="inline-flex items-center gap-[7px] text-[10px] capitalize text-[#73788a] max-[520px]:hidden">
            <svg class="size-[15px] text-[#9b9fb0]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3.5" y="5" width="17" height="15.5" rx="2.5" stroke="currentColor" stroke-width="1.7" /><path d="M7.5 3.5v3M16.5 3.5v3M3.5 9.5h17" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" /></svg>
            {{ todayLabel }}
          </span>
          <button
            v-if="status === 'closed'"
            type="button"
            class="inline-flex min-h-8 cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#282b3d] px-[11px] text-[10px] font-semibold text-white transition hover:bg-[#393c50] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 max-[520px]:min-h-[30px] max-[520px]:px-2 max-[520px]:text-[9px]"
            @click="downloadReport"
          >
            <svg class="size-[14px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5v11m0 0 4-4m-4 4-4-4M4.5 16v4h15v-4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
            Descargar PDF
          </button>
          <button v-if="status === 'closed'" type="button" class="inline-flex min-h-8 cursor-pointer items-center justify-center gap-2 rounded-[7px] border border-[#e6e8ef] bg-white px-[11px] text-[10px] font-semibold text-[#565a6c] transition hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 max-[520px]:min-h-[30px] max-[520px]:px-2 max-[520px]:text-[9px]" @click="newCashCount">
            Nuevo arqueo
          </button>
        </div>
      </header>

      <div v-show="activeSection === 'cash'" class="mx-auto w-[min(1120px,91%)] py-[38px] max-[780px]:w-[92%] max-[780px]:pt-7 min-[1500px]:pt-12 print:hidden">
        <section class="mb-[25px] flex items-end justify-between gap-[18px] max-[520px]:items-start">
          <div>
            <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-block h-0.5 w-[17px] rounded-sm bg-[#8875e5]"></span> CONTROL DE CAJA</div>
            <h1 class="my-[9px] mb-1 text-[clamp(29px,3vw,37px)] font-bold leading-[1.15] tracking-[-1.7px] text-[#232637] max-[520px]:text-[29px]">Arqueo diario<span class="text-[#7967de]">.</span></h1>
            <p class="m-0 text-xs text-[#8a8e9e] max-[520px]:max-w-[235px] max-[520px]:text-[10px] max-[520px]:leading-[1.5]">Cada movimiento en su lugar. Tu control del efectivo, siempre ordenado.</p>
          </div>
          <div class="inline-flex items-center gap-[7px] whitespace-nowrap rounded-full border border-[#e6e8ef] bg-white px-[10px] py-[7px] text-[10px] font-semibold text-[#777c8d] max-[520px]:mt-[19px] max-[520px]:px-2 max-[520px]:py-[6px] max-[520px]:text-[8px]">
            <span class="inline-block size-[6px] rounded-full" :class="status === 'opening' ? 'bg-[#e5aa53]' : status === 'closing' ? 'bg-[#8b77e4]' : 'bg-[#29ac84]'"></span>
            {{ status === 'opening' ? 'Por iniciar' : status === 'open' ? 'Caja abierta' : status === 'closing' ? 'En cierre' : 'Caja cerrada' }}
          </div>
        </section>

        <OpeningForm
          v-if="status === 'opening'"
          v-model:date="date"
          :opening-balances="openingBalances"
          :error-message="errorMessage"
          @begin="beginDay"
        />

        <template v-else>
          <BalanceSummary
            :status="status"
            :expected-balances="expectedBalances"
            :movement-counts="movementCounts"
            :closing-balances="closingBalances"
            :total-expected="totalExpected"
            :total-counted="totalCounted"
            :total-difference="totalDifference"
          />

          <section class="mt-[17px] grid grid-cols-[minmax(0,1.85fr)_minmax(235px,.85fr)] items-start gap-[15px] max-[780px]:grid-cols-1 max-[520px]:mt-[10px] max-[520px]:gap-[10px]">
            <ActivityPanel
              v-model:movement-category="movementCategory"
              v-model:movement-type="movementType"
              v-model:movement-payment-method="movementPaymentMethod"
              v-model:movement-amount="movementAmount"
              v-model:movement-concept="movementConcept"
              v-model:keep-movement-context="keepMovementContext"
              :status="status"
              :movements="movements"
              :error-message="errorMessage"
              @add-movement="addMovement"
              @remove-movement="removeMovement"
            />
            <ClosingPanel
              :status="status"
              :closing-balances="closingBalances"
              :expected-balances="expectedBalances"
              :total-expected="totalExpected"
              :total-difference="totalDifference"
              :error-message="errorMessage"
              :download-error="downloadError"
              @start-closing="startClosing"
              @finish-closing="finishClosing"
              @back-to-movements="status = 'open'"
              @download="downloadReport"
            />
          </section>
        </template>
      </div>

      <InventoryPanel v-show="activeSection === 'inventory'" />

      <PrintableReport
        v-if="activeSection === 'cash' && status === 'closed'"
        :today-label="todayLabel"
        :movements="movements"
        :opening-balances="openingBalances"
        :closing-balances="closingBalances"
        :expected-balances="expectedBalances"
        :total-expected="totalExpected"
        :total-counted="totalCounted"
        :total-difference="totalDifference"
      />
    </main>
    <div
      v-if="toastMessage"
      class="fixed bottom-5 right-5 z-50 flex items-center gap-2 rounded-lg border bg-white px-4 py-3 text-sm font-medium shadow-[0_8px_28px_#202a241f] max-[520px]:bottom-3 max-[520px]:right-3"
      :class="toastMessage === 'Movimiento eliminado' ? 'border-[#f0c5c8] text-[#c34e59]' : 'border-[#b8e4d1] text-[#21835f]'"
      role="status"
      aria-live="polite"
    >
      <svg class="size-4 flex-none" viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="m5 12 4 4L19 6" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
      {{ toastMessage }}
    </div>
  </div>
</template>
