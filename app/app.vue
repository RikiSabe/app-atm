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
      <AppHeader
        :today-label="todayLabel"
        :status="status"
        @download="downloadReport"
        @new-cash-count="newCashCount"
      />

      <div class="mx-auto w-[min(1120px,91%)] py-[38px] max-[780px]:w-[92%] max-[780px]:pt-7 min-[1500px]:pt-12 print:hidden">
        <PageHeading :status="status" />

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

      <PrintableReport
        v-if="status === 'closed'"
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
