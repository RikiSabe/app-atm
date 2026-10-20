<script setup lang="ts">
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

type CategoryId = 'constante' | 'variable'
type MovementType = 'entrada' | 'salida'
type PaymentMethod = 'QR' | 'Efectivo'
type CashStatus = 'opening' | 'open' | 'closing' | 'closed'

interface CashCategory {
  id: CategoryId
  name: string
}

interface CashMovement {
  id: number
  category: CategoryId
  type: MovementType
  paymentMethod: PaymentMethod
  amount: number
  concept: string
  time: string
}

interface PaymentSummary {
  count: number
  entries: number
  exits: number
  net: number
}

const categories: CashCategory[] = [
  {
    id: 'constante',
    name: 'Activo constante'
  },
  {
    id: 'variable',
    name: 'Activo variable'
  }
]

const paymentMethods: PaymentMethod[] = ['QR', 'Efectivo']
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
const movementType = ref<MovementType>('entrada')
const movementPaymentMethod = ref<PaymentMethod>('Efectivo')
const movementAmount = ref<number | null>(null)
const movementConcept = ref('')
const errorMessage = ref('')
const downloadError = ref('')
function localDateValue() {
  const now = new Date()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${now.getFullYear()}-${month}-${day}`
}

const todayLabel = computed(() =>
  new Date(`${date.value}T12:00:00`).toLocaleDateString('es', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
)

const totalExpected = computed(() =>
  categories.reduce((total, category) => total + expectedBalance(category.id), 0)
)
const totalCounted = computed(() =>
  categories.reduce((total, category) => total + closingBalances[category.id], 0)
)
const totalDifference = computed(() => totalCounted.value - totalExpected.value)

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

function paymentSummary(paymentMethod: PaymentMethod): PaymentSummary {
  const methodMovements = movements.value.filter(movement => movement.paymentMethod === paymentMethod)
  const entries = methodMovements
    .filter(movement => movement.type === 'entrada')
    .reduce((total, movement) => total + movement.amount, 0)
  const exits = methodMovements
    .filter(movement => movement.type === 'salida')
    .reduce((total, movement) => total + movement.amount, 0)

  return {
    count: methodMovements.length,
    entries,
    exits,
    net: entries - exits
  }
}

function formatAmount(amount: number) {
  return new Intl.NumberFormat('es', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

function formatSignedAmount(amount: number) {
  return amount < 0
    ? `−$${formatAmount(Math.abs(amount))}`
    : `+$${formatAmount(amount)}`
}

function formatPdfSignedAmount(amount: number) {
  return formatSignedAmount(amount).replace('−', '-')
}

function lastPdfTableY(pdf: jsPDF) {
  return (pdf as jsPDF & { lastAutoTable: { finalY: number } }).lastAutoTable.finalY
}

function formatTime(time: string) {
  return new Date(time).toLocaleTimeString('es', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

function categoryName(categoryId: CategoryId) {
  return categories.find(category => category.id === categoryId)?.name ?? categoryId
}

function beginDay() {
  errorMessage.value = ''
  if (!date.value) {
    errorMessage.value = 'Selecciona la fecha del arqueo para continuar.'
    return
  }
  if (categories.some(category => openingBalances[category.id] < 0)) {
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
  movementConcept.value = ''
}

function removeMovement(id: number) {
  movements.value = movements.value.filter(movement => movement.id !== id)
}

function startClosing() {
  errorMessage.value = ''
  for (const category of categories) {
    closingBalances[category.id] = expectedBalance(category.id)
  }
  status.value = 'closing'
}

function finishClosing() {
  errorMessage.value = ''
  if (categories.some(category => !Number.isFinite(closingBalances[category.id]) || closingBalances[category.id] < 0)) {
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
  openingBalances.constante = 0
  openingBalances.variable = 0
  closingBalances.constante = 0
  closingBalances.variable = 0
  movements.value = []
  errorMessage.value = ''
}

async function downloadReport() {
  downloadError.value = ''

  try {
    const pdf = new jsPDF({ orientation: 'landscape', unit: 'mm', format: 'a4' })

    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(18)
    pdf.text('Arqueo diario', 14, 18)
    pdf.setFont('helvetica', 'normal')
    pdf.setFontSize(10)
    pdf.text(todayLabel.value, 14, 25)

    autoTable(pdf, {
      head: [['Categoría', 'Saldo inicial', 'Entradas', 'Salidas', 'Esperado', 'Contado', 'Diferencia']],
      body: categories.map(category => {
        const categoryMovements = movementsFor(category.id)
        const entries = categoryMovements
          .filter(movement => movement.type === 'entrada')
          .reduce((total, movement) => total + movement.amount, 0)
        const exits = categoryMovements
          .filter(movement => movement.type === 'salida')
          .reduce((total, movement) => total + movement.amount, 0)

        return [
          category.name,
          `$${formatAmount(openingBalances[category.id])}`,
          `$${formatAmount(entries)}`,
          `$${formatAmount(exits)}`,
          `$${formatAmount(expectedBalance(category.id))}`,
          `$${formatAmount(closingBalances[category.id])}`,
          formatPdfSignedAmount(closingBalances[category.id] - expectedBalance(category.id))
        ]
      }),
      foot: [[
        'Total',
        '',
        '',
        '',
        `$${formatAmount(totalExpected.value)}`,
        `$${formatAmount(totalCounted.value)}`,
        formatPdfSignedAmount(totalDifference.value)
      ]],
      startY: 32,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 8, cellPadding: 2.5 },
      headStyles: { fillColor: [118, 97, 220], textColor: 255 },
      footStyles: { fillColor: [245, 244, 252], textColor: [45, 47, 63], fontStyle: 'bold' }
    })

    autoTable(pdf, {
      head: [['Método de pago', 'Movimientos', 'Entradas', 'Salidas', 'Neto']],
      body: paymentMethods.map(paymentMethod => {
        const summary = paymentSummary(paymentMethod)
        return [
          paymentMethod,
          String(summary.count),
          `$${formatAmount(summary.entries)}`,
          `$${formatAmount(summary.exits)}`,
          `$${formatAmount(summary.net)}`
        ]
      }),
      startY: lastPdfTableY(pdf) + 10,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 8, cellPadding: 2.5 },
      headStyles: { fillColor: [118, 97, 220], textColor: 255 }
    })

    autoTable(pdf, {
      head: [['Hora', 'Categoría', 'Tipo', 'Método de pago', 'Concepto', 'Importe']],
      body: movements.value.length
        ? movements.value.map(movement => [
            formatTime(movement.time),
            categoryName(movement.category),
            movement.type,
            movement.paymentMethod,
            movement.concept,
            `$${formatAmount(movement.amount)}`
          ])
        : [['', '', '', '', 'No se registraron movimientos.', '']],
      startY: lastPdfTableY(pdf) + 10,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 8, cellPadding: 2.5, overflow: 'linebreak' },
      headStyles: { fillColor: [118, 97, 220], textColor: 255 }
    })

    const blob = pdf.output('blob')
    const filename = `arqueo-${date.value}.pdf`

    if (typeof File !== 'undefined') {
      const file = new File([blob], filename, { type: 'application/pdf' })

      if (navigator.canShare?.({ files: [file] }) && navigator.share) {
        await navigator.share({ files: [file], title: filename })
        return
      }
    }

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

        <section v-if="status === 'opening'" class="grid grid-cols-[minmax(0,760px)] justify-center max-[780px]:grid-cols-[minmax(0,1fr)]">
          <article class="rounded-xl border border-[#eceef4] bg-white p-[27px_29px_22px] shadow-[0_3px_12px_#282b3c05] max-[520px]:p-[21px_17px_18px]">
            <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-grid size-[19px] place-items-center rounded-md bg-[#f0edff] text-[9px] tracking-normal text-[#725fd2]">01</span> INICIO DE JORNADA</div>
            <h2 class="mb-[5px] mt-[14px] text-[19px] font-bold tracking-[-.5px] text-[#282b3b] max-[520px]:text-[17px]">Empecemos por el principio</h2>
            <p class="mb-[22px] max-w-[440px] text-[11px] leading-[1.7] text-[#898d9c]">Anota con cuánto dinero inicia cada categoría. Puedes agregar movimientos durante el día.</p>

            <label class="mb-[7px] block text-[10px] font-semibold text-[#55596a]" for="cash-date">Fecha del arqueo</label>
            <div class="flex min-h-10 items-center gap-[9px] rounded-[7px] border border-[#e8eaf0] px-[11px] text-[#9296a7]">
              <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3.5" y="5" width="17" height="15.5" rx="2.5" stroke="currentColor" stroke-width="1.7" /><path d="M7.5 3.5v3M16.5 3.5v3M3.5 9.5h17" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" /></svg>
              <input id="cash-date" v-model="date" class="w-full border-0 bg-transparent text-[11px] text-[#333648] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="date" :disabled="status !== 'opening'" />
            </div>

            <div class="my-[22px] flex items-center gap-3 text-[8px] font-bold tracking-[.9px] text-[#a1a5b2] before:h-px before:flex-1 before:bg-[#eff0f5] before:content-[''] after:h-px after:flex-1 after:bg-[#eff0f5] after:content-['']"><span>SALDOS INICIALES</span></div>
            <div class="grid grid-cols-2 gap-[14px] max-[520px]:grid-cols-1 max-[520px]:gap-[10px]">
              <label v-for="category in categories" :key="category.id" class="block">
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
            <button type="button" class="mt-[18px] inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] hover:shadow-[0_5px_14px_#7966de40] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="beginDay">
              Abrir caja
              <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M5 12h14m-6-6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
            </button>
            <p class="mt-3 flex items-center justify-center gap-[6px] text-center text-[9px] text-[#9a9eac]"><svg class="size-[13px] text-[#8c80d8]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5 4.5 6.5v5.1c0 4.3 3 7.2 7.5 8.9 4.5-1.7 7.5-4.6 7.5-8.9V6.5L12 3.5Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" /><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" /></svg> Los datos se mantienen solo mientras esta página esté abierta.</p>
          </article>

        </section>

        <template v-else>
          <section class="grid grid-cols-[1fr_1fr_1.05fr] gap-[15px] max-[780px]:grid-cols-2 max-[780px]:gap-[9px]" aria-label="Resumen de saldos">
            <article v-for="category in categories" :key="category.id" class="min-h-[145px] rounded-xl border border-[#eceef4] bg-white p-[17px] shadow-[0_3px_12px_#282b3c05] max-[520px]:min-h-[157px] max-[520px]:p-[12px_11px_11px]" :class="category.id === 'constante' ? 'border-t-[3px] border-t-[#8b78e3]' : 'border-t-[3px] border-t-[#55bea1]'">
              <div class="text-[10px] font-semibold text-[#777b8b] max-[520px]:text-[9px]">{{ category.name }}</div>
              <div class="mt-[3px] text-[25px] font-bold tracking-[-.9px] text-[#292c3d] max-[520px]:text-[21px]"><span class="mr-0.5 text-base font-medium text-[#8c8f9e] max-[520px]:text-[13px]">$</span>{{ formatAmount(expectedBalance(category.id)) }}</div>
              <div class="mt-2 flex justify-between gap-[5px] text-[8px] text-[#989cab] max-[520px]:flex-wrap max-[520px]:text-[7px]">
                <span>Saldo esperado</span>
                <span>{{ movementsFor(category.id).length }} {{ movementsFor(category.id).length === 1 ? 'movimiento' : 'movimientos' }}</span>
              </div>
              <div v-if="status === 'closed'" class="mt-[7px] flex gap-[10px] text-[8px] text-[#898d9b] max-[520px]:flex-col max-[520px]:gap-[3px]">
                <span class="flex gap-1">Contado <strong class="font-semibold text-[#505466]">${{ formatAmount(closingBalances[category.id]) }}</strong></span>
                <span class="flex gap-1" :class="expectedBalance(category.id) === closingBalances[category.id] ? 'text-[#279978]' : 'text-[#c16b4c]'">
                  Diferencia <strong>{{ formatSignedAmount(closingBalances[category.id] - expectedBalance(category.id)) }}</strong>
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

          <section class="mt-[17px] grid grid-cols-[minmax(0,1.85fr)_minmax(235px,.85fr)] items-start gap-[15px] max-[780px]:grid-cols-1 max-[520px]:mt-[10px] max-[520px]:gap-[10px]">
            <div class="grid min-w-0 gap-[14px]">
              <article v-if="status === 'open'" class="rounded-xl border border-[#eceef4] bg-white p-[17px_18px_15px] shadow-[0_3px_12px_#282b3c05] max-[520px]:px-3">
                <div class="flex items-center justify-between gap-[10px]">
                  <div>
                    <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-grid size-[19px] place-items-center rounded-md bg-[#f0edff] text-[9px] tracking-normal text-[#725fd2]">02</span> MOVIMIENTOS DEL DÍA</div>
                    <h2 class="mt-2 text-[15px] font-bold tracking-[-.3px] text-[#282b3b]">Registra un movimiento</h2>
                  </div>
                </div>
                <form class="mt-[17px] grid grid-cols-[1fr_.8fr_1fr_.85fr_1.3fr_auto] items-end gap-2 max-[1100px]:grid-cols-3 max-[520px]:grid-cols-2" @submit.prevent="addMovement">
                  <label class="min-w-0 max-[1100px]:col-span-1">
                    <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Categoría</span>
                    <select v-model="movementCategory" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
                      <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
                    </select>
                  </label>
                  <label class="min-w-0">
                    <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Método de pago</span>
                    <select v-model="movementPaymentMethod" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
                      <option v-for="paymentMethod in paymentMethods" :key="paymentMethod" :value="paymentMethod">{{ paymentMethod }}</option>
                    </select>
                  </label>
                  <label class="min-w-0">
                    <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Tipo</span>
                    <select v-model="movementType" class="h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]">
                      <option value="entrada">Entrada</option>
                      <option value="salida">Salida</option>
                    </select>
                  </label>
                  <label class="min-w-0">
                    <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Importe</span>
                    <span class="flex min-h-[38px] items-center rounded-[7px] border border-[#e8eaf0] bg-white px-2 focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
                      <span class="mr-[5px] text-[10px] text-[#a3a7b5]">$</span>
                      <input v-model.number="movementAmount" class="w-full min-w-0 border-0 bg-transparent text-[10px] text-[#303345] outline-none placeholder:text-[#b7bac5] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0.01" step="0.01" placeholder="0.00" />
                    </span>
                  </label>
                  <label class="min-w-0 max-[1100px]:col-span-2 max-[520px]:col-span-2">
                    <span class="mb-[7px] block text-[10px] font-semibold text-[#55596a]">Concepto</span>
                    <input v-model="movementConcept" class="block h-[38px] w-full rounded-md border border-[#e8eaf0] bg-white px-[9px] text-[9px] text-[#55596b] outline-none placeholder:text-[9px] placeholder:text-[#b7bac5] focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="text" maxlength="80" placeholder="¿De qué se trata?" />
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
                    <button v-if="status === 'open'" type="button" class="grid size-[25px] flex-none cursor-pointer place-items-center rounded-md border-0 bg-transparent text-[#a1a5b3] transition hover:bg-[#fff0f1] hover:text-[#c34e59] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] max-[520px]:size-[22px]" :aria-label="`Eliminar movimiento ${movement.concept}`" @click="removeMovement(movement.id)">
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
                        <tr v-for="paymentMethod in paymentMethods" :key="paymentMethod" class="text-[#55596a]">
                          <th class="border-b border-[#f0f1f5] px-2 py-[7px] font-semibold">{{ paymentMethod }}</th>
                          <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">{{ paymentSummary(paymentMethod).count }}</td>
                          <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">${{ formatAmount(paymentSummary(paymentMethod).entries) }}</td>
                          <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right">${{ formatAmount(paymentSummary(paymentMethod).exits) }}</td>
                          <td class="border-b border-[#f0f1f5] px-2 py-[7px] text-right font-semibold">${{ formatAmount(paymentSummary(paymentMethod).net) }}</td>
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
                <button type="button" class="inline-flex min-h-[38px] w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#282b3d] px-[15px] text-[9px] font-semibold text-white transition hover:-translate-y-px hover:bg-[#393c50] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="startClosing">
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
                <label v-for="category in categories" :key="category.id" class="mt-[13px] block">
                  <span class="mb-[6px] flex justify-between gap-[5px] text-[9px] font-semibold text-[#535768]">{{ category.name }} <small class="whitespace-nowrap text-[8px] font-normal text-[#9397a6]">Esperado: ${{ formatAmount(expectedBalance(category.id)) }}</small></span>
                  <span class="flex min-h-9 items-center rounded-[7px] border border-[#e8eaf0] bg-white px-[11px] focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
                    <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
                    <input v-model.number="closingBalances[category.id]" class="w-full min-w-0 border-0 bg-transparent text-[10px] text-[#303345] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="0.01" :aria-label="`Importe contado de ${category.name}`" />
                  </span>
                  <small class="mt-1 block text-[8px]" :class="closingBalances[category.id] - expectedBalance(category.id) === 0 ? 'text-[#279978]' : 'text-[#c16b4c]'">
                    Diferencia: {{ formatSignedAmount(closingBalances[category.id] - expectedBalance(category.id)) }}
                  </small>
                </label>
                <p v-if="errorMessage" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ errorMessage }}</p>
                <button type="button" class="mt-[13px] inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="finishClosing">
                  Confirmar cierre
                  <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
                </button>
                <button type="button" class="mx-auto mt-[10px] block cursor-pointer border-0 bg-transparent px-[6px] py-[3px] text-[9px] text-[#777b8b] underline decoration-1 underline-offset-2 hover:text-[#6552cc] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" @click="status = 'open'">Volver a los movimientos</button>
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
                <button type="button" class="inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-[15px] text-[11px] font-semibold text-white shadow-[0_4px_9px_#7966de29] transition hover:-translate-y-px hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2" @click="downloadReport">
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
          </section>
        </template>

      </div>

      <InventoryPanel v-show="activeSection === 'inventory'" />

      <section v-if="activeSection === 'cash' && status === 'closed'" class="hidden p-0 font-sans text-[#252637] print:block" aria-label="Reporte de arqueo">
        <div class="flex items-center justify-between border-b-2 border-[#7661dc] pb-[13px] text-[10px] font-bold tracking-[.7px] text-[#76798a]">
          <div class="text-lg font-extrabold tracking-[-.5px] text-[#292c3d]">arqueo</div>
          <span>REPORTE DE ARQUEO DIARIO</span>
        </div>
        <h1 class="mb-1 mt-[22px] text-2xl font-bold">Arqueo diario</h1>
        <p class="mb-5 text-[11px] capitalize text-[#737789]">{{ todayLabel }}</p>
        <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a] [&_tfoot_th]:bg-[#faf9fd]">
          <thead><tr><th>Categoría</th><th>Saldo inicial</th><th>Entradas</th><th>Salidas</th><th>Esperado</th><th>Contado</th><th>Diferencia</th></tr></thead>
          <tbody>
            <tr v-for="category in categories" :key="category.id">
              <td>{{ category.name }}</td>
              <td>${{ formatAmount(openingBalances[category.id]) }}</td>
              <td>${{ formatAmount(movementsFor(category.id).filter(item => item.type === 'entrada').reduce((sum, item) => sum + item.amount, 0)) }}</td>
              <td>${{ formatAmount(movementsFor(category.id).filter(item => item.type === 'salida').reduce((sum, item) => sum + item.amount, 0)) }}</td>
              <td>${{ formatAmount(expectedBalance(category.id)) }}</td>
              <td>${{ formatAmount(closingBalances[category.id]) }}</td>
              <td>{{ formatSignedAmount(closingBalances[category.id] - expectedBalance(category.id)) }}</td>
            </tr>
          </tbody>
          <tfoot><tr><th>Total</th><th colspan="3"></th><th>${{ formatAmount(totalExpected) }}</th><th>${{ formatAmount(totalCounted) }}</th><th>{{ formatSignedAmount(totalDifference) }}</th></tr></tfoot>
        </table>
        <h2 class="mb-[9px] mt-[23px] text-sm font-bold">Resumen por método de pago</h2>
        <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a]">
          <thead><tr><th>Método</th><th>Movimientos</th><th>Entradas</th><th>Salidas</th><th>Neto</th></tr></thead>
          <tbody>
            <tr v-for="paymentMethod in paymentMethods" :key="paymentMethod">
              <td>{{ paymentMethod }}</td>
              <td>{{ paymentSummary(paymentMethod).count }}</td>
              <td>${{ formatAmount(paymentSummary(paymentMethod).entries) }}</td>
              <td>${{ formatAmount(paymentSummary(paymentMethod).exits) }}</td>
              <td>${{ formatAmount(paymentSummary(paymentMethod).net) }}</td>
            </tr>
          </tbody>
        </table>
        <h2 class="mb-[9px] mt-[23px] text-sm font-bold">Detalle de movimientos</h2>
        <table class="mb-5 w-full border-collapse text-[9px] [&_td]:border [&_td]:border-[#e1e3eb] [&_td]:p-2 [&_td]:text-left [&_th]:border [&_th]:border-[#e1e3eb] [&_th]:p-2 [&_th]:text-left [&_thead_th]:bg-[#f5f4fc] [&_thead_th]:font-bold [&_thead_th]:text-[#55596a] [&_tfoot_th]:bg-[#faf9fd]">
          <thead><tr><th>Hora</th><th>Categoría</th><th>Tipo</th><th>Método de pago</th><th>Concepto</th><th>Importe</th></tr></thead>
          <tbody>
            <tr v-for="movement in movements" :key="movement.id">
              <td>{{ formatTime(movement.time) }}</td><td>{{ categoryName(movement.category) }}</td><td>{{ movement.type }}</td><td>{{ movement.paymentMethod }}</td><td>{{ movement.concept }}</td><td>${{ formatAmount(movement.amount) }}</td>
            </tr>
            <tr v-if="!movements.length"><td colspan="6">No se registraron movimientos.</td></tr>
          </tbody>
        </table>
        <p class="mt-[25px] text-[8px] text-[#858999]">Reporte generado en el dispositivo. Los datos no se almacenan en la aplicación.</p>
      </section>
    </main>
  </div>
</template>
