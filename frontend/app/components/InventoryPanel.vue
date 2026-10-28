<script setup lang="ts">
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

interface InventoryItem {
  id: string
  name: string
  quantity: number
  net_price: number
  sale_price: number
}

const inventoryItems = ref<InventoryItem[]>([])
const inventoryName = ref('')
const inventoryQuantity = ref<number | null>(null)
const inventoryNetPrice = ref<number | null>(null)
const inventorySalePrice = ref<number | null>(null)
const inventoryError = ref('')
const inventoryLoadError = ref('')
const inventoryReportError = ref('')
const editingInventoryId = ref<string | null>(null)
const isLoadingInventory = ref(true)
const isSavingInventory = ref(false)
const deletingInventoryId = ref<string | null>(null)
const apiBase = useRuntimeConfig().public.apiBase.replace(/\/+$/, '')

function refreshAfterCheckout() {
  void loadInventory()
}

onMounted(() => {
  window.addEventListener('inventory:updated', refreshAfterCheckout)
  void loadInventory()
})

onUnmounted(() => {
  window.removeEventListener('inventory:updated', refreshAfterCheckout)
})

const inventoryQuantityTotal = computed(() =>
  inventoryItems.value.reduce((total, item) => total + item.quantity, 0)
)
const inventoryNetValueTotal = computed(() =>
  inventoryItems.value.reduce((total, item) => total + item.quantity * item.net_price, 0)
)
const inventorySaleValueTotal = computed(() =>
  inventoryItems.value.reduce((total, item) => total + item.quantity * item.sale_price, 0)
)

function formatAmount(amount: number) {
  return new Intl.NumberFormat('es', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

function resetInventoryForm() {
  inventoryName.value = ''
  inventoryQuantity.value = null
  inventoryNetPrice.value = null
  inventorySalePrice.value = null
  editingInventoryId.value = null
  inventoryError.value = ''
}

function apiErrorMessage(error: unknown, fallback: string) {
  if (error && typeof error === 'object' && 'data' in error) {
    const data = error.data
    if (data && typeof data === 'object' && 'error' in data && typeof data.error === 'string') {
      return data.error
    }
  }
  return error instanceof Error ? error.message : fallback
}

function lastPdfTableY(pdf: jsPDF) {
  return (pdf as jsPDF & { lastAutoTable: { finalY: number } }).lastAutoTable.finalY
}

async function downloadInventoryReport() {
  inventoryReportError.value = ''

  try {
    const pdf = new jsPDF({ orientation: 'portrait', unit: 'mm', format: 'a4' })
    const generatedAt = new Date()
    const reportDate = generatedAt.toLocaleDateString('es', {
      day: '2-digit',
      month: 'long',
      year: 'numeric'
    })
    const netValue = inventoryNetValueTotal.value
    const saleValue = inventorySaleValueTotal.value

    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(18)
    pdf.text('Reporte de inventario', 14, 18)
    pdf.setFont('helvetica', 'normal')
    pdf.setFontSize(10)
    pdf.text(`Generado el ${reportDate}`, 14, 25)

    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(13)
    pdf.text('Resumen', 14, 36)
    autoTable(pdf, {
      body: [
        ['Artículos', String(inventoryItems.value.length)],
        ['Unidades en existencia', String(inventoryQuantityTotal.value)],
        ['Valor total al costo neto', `$${formatAmount(netValue)}`],
        ['Valor total a precio de venta', `$${formatAmount(saleValue)}`],
        ['Margen bruto potencial', `$${formatAmount(saleValue - netValue)}`]
      ],
      startY: 40,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 9, cellPadding: 3 },
      columnStyles: { 1: { halign: 'right' } },
      headStyles: { fillColor: [118, 97, 220], textColor: 255 },
      margin: { left: 14, right: 14 }
    })

    const tableStartY = lastPdfTableY(pdf) + 10
    if (tableStartY > 245) {
      pdf.addPage()
    }
    const inventoryTableStartY = tableStartY > 245 ? 20 : tableStartY
    pdf.setFont('helvetica', 'bold')
    pdf.setFontSize(13)
    pdf.text('Detalle de artículos', 14, inventoryTableStartY)
    autoTable(pdf, {
      head: [['Artículo', 'Unidades', 'Neto unit.', 'Venta unit.', 'Total neto', 'Total venta']],
      body: inventoryItems.value.length
        ? inventoryItems.value.map(item => [
            item.name,
            String(item.quantity),
            `$${formatAmount(item.net_price)}`,
            `$${formatAmount(item.sale_price)}`,
            `$${formatAmount(item.quantity * item.net_price)}`,
            `$${formatAmount(item.quantity * item.sale_price)}`
          ])
        : [['No hay artículos registrados.', '', '', '', '', '']],
      foot: [[
        'Totales',
        String(inventoryQuantityTotal.value),
        '',
        '',
        `$${formatAmount(netValue)}`,
        `$${formatAmount(saleValue)}`
      ]],
      startY: inventoryTableStartY + 4,
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 8, cellPadding: 2.5, overflow: 'linebreak' },
      headStyles: { fillColor: [118, 97, 220], textColor: 255 },
      footStyles: { fillColor: [245, 244, 252], textColor: [45, 47, 63], fontStyle: 'bold' },
      columnStyles: {
        0: { cellWidth: 'auto' },
        1: { cellWidth: 18, halign: 'right' },
        2: { cellWidth: 25, halign: 'right' },
        3: { cellWidth: 25, halign: 'right' },
        4: { cellWidth: 27, halign: 'right' },
        5: { cellWidth: 27, halign: 'right' }
      },
      margin: { left: 14, right: 14 }
    })

    pdf.setFont('helvetica', 'normal')
    pdf.setFontSize(8)
    pdf.setTextColor(133, 137, 153)
    pdf.text('Reporte generado desde el inventario de la aplicación.', 14, 285)

    const blob = pdf.output('blob')
    const filename = `inventario-${generatedAt.toISOString().slice(0, 10)}.pdf`
    const objectUrl = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = objectUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.setTimeout(() => URL.revokeObjectURL(objectUrl), 1000)
  } catch (error) {
    inventoryReportError.value = error instanceof Error
      ? `No se pudo generar el reporte: ${error.message}`
      : 'No se pudo generar el reporte por un error desconocido.'
  }
}

async function loadInventory() {
  isLoadingInventory.value = true
  inventoryLoadError.value = ''
  try {
    inventoryItems.value = await $fetch<InventoryItem[]>(`${apiBase}/inventory`)
  } catch (error) {
    inventoryLoadError.value = apiErrorMessage(error, 'No se pudo cargar el inventario.')
  } finally {
    isLoadingInventory.value = false
  }
}

async function saveInventoryItem() {
  inventoryError.value = ''
  const name = inventoryName.value.trim()
  const quantity = inventoryQuantity.value
  const netPrice = inventoryNetPrice.value
  const salePrice = inventorySalePrice.value

  if (!name) {
    inventoryError.value = 'Escribe el nombre del artículo.'
    return
  }
  if (typeof quantity !== 'number' || !Number.isInteger(quantity) || quantity < 0) {
    inventoryError.value = 'La cantidad debe ser un número entero igual o mayor que cero.'
    return
  }
  if (typeof netPrice !== 'number' || !Number.isFinite(netPrice) || netPrice < 0) {
    inventoryError.value = 'El precio neto debe ser un importe válido igual o mayor que cero.'
    return
  }
  if (typeof salePrice !== 'number' || !Number.isFinite(salePrice) || salePrice < 0) {
    inventoryError.value = 'El precio de venta debe ser un importe válido igual o mayor que cero.'
    return
  }

  isSavingInventory.value = true
  try {
    const payload = { name, quantity, net_price: netPrice, sale_price: salePrice }
    if (editingInventoryId.value !== null) {
      const updatedItem = await $fetch<InventoryItem>(
        `${apiBase}/inventory/${encodeURIComponent(editingInventoryId.value)}`,
        { method: 'PUT', body: payload }
      )
      const itemIndex = inventoryItems.value.findIndex(item => item.id === updatedItem.id)
      if (itemIndex !== -1) {
        inventoryItems.value[itemIndex] = updatedItem
      }
    } else {
      const createdItem = await $fetch<InventoryItem>(`${apiBase}/inventory`, {
        method: 'POST',
        body: payload
      })
      inventoryItems.value.push(createdItem)
    }
    resetInventoryForm()
  } catch (error) {
    inventoryError.value = apiErrorMessage(error, 'No se pudo guardar el artículo.')
  } finally {
    isSavingInventory.value = false
  }
}

function editInventoryItem(item: InventoryItem) {
  editingInventoryId.value = item.id
  inventoryName.value = item.name
  inventoryQuantity.value = item.quantity
  inventoryNetPrice.value = item.net_price
  inventorySalePrice.value = item.sale_price
  inventoryError.value = ''
}

async function removeInventoryItem(item: InventoryItem) {
  if (!window.confirm(`¿Eliminar "${item.name}" del inventario?`)) {
    return
  }
  inventoryError.value = ''
  deletingInventoryId.value = item.id
  try {
    await $fetch(`${apiBase}/inventory/${encodeURIComponent(item.id)}`, { method: 'DELETE' })
    inventoryItems.value = inventoryItems.value.filter(existingItem => existingItem.id !== item.id)
    if (editingInventoryId.value === item.id) {
      resetInventoryForm()
    }
  } catch (error) {
    inventoryError.value = apiErrorMessage(error, 'No se pudo eliminar el artículo.')
  } finally {
    deletingInventoryId.value = null
  }
}
</script>

<template>
  <section class="mx-auto w-[min(1120px,91%)] py-[38px] max-[780px]:w-[92%] max-[780px]:pt-7 min-[1500px]:pt-12 print:hidden" aria-labelledby="inventory-title">
    <div class="mb-[25px]">
      <div class="flex items-center gap-2 text-[9px] font-bold tracking-[1.05px] text-[#9296a7]"><span class="inline-block h-0.5 w-[17px] rounded-sm bg-[#8875e5]"></span> CONTROL DE PRODUCTOS</div>
      <h1 id="inventory-title" class="my-[9px] mb-1 text-[clamp(29px,3vw,37px)] font-bold leading-[1.15] tracking-[-1.7px] text-[#232637] max-[520px]:text-[29px]">Inventario<span class="text-[#7967de]">.</span></h1>
      <p class="m-0 text-xs text-[#8a8e9e]">Administra tus artículos, existencias y precios.</p>
    </div>

    <section class="mb-4 grid grid-cols-4 gap-[12px] max-[780px]:grid-cols-2 max-[520px]:grid-cols-1 max-[520px]:gap-2" aria-label="Resumen del inventario">
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">ARTÍCULOS</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">{{ inventoryItems.length }}</div>
      </article>
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">UNIDADES EN EXISTENCIA</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">{{ inventoryQuantityTotal }}</div>
      </article>
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">VALOR AL COSTO (NETO)</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">${{ formatAmount(inventoryNetValueTotal) }}</div>
      </article>
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">VALOR A PRECIO DE VENTA</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">${{ formatAmount(inventorySaleValueTotal) }}</div>
      </article>
    </section>

    <section class="grid grid-cols-[minmax(0,0.85fr)_minmax(0,1.65fr)] items-start gap-[15px] max-[780px]:grid-cols-1" aria-label="Administrar artículos">
      <form class="rounded-xl border border-[#eceef4] bg-white p-[19px] shadow-[0_3px_12px_#282b3c05]" @submit.prevent="saveInventoryItem">
        <h2 class="mb-1 text-[15px] font-bold text-[#282b3b]">{{ editingInventoryId === null ? 'Agregar artículo' : 'Editar artículo' }}</h2>
        <p class="mb-4 text-[10px] text-[#898d9c]">Registra nombre, existencias, precio neto y precio de venta.</p>
        <label class="mb-3 block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Nombre</span>
          <input v-model="inventoryName" class="h-10 w-full rounded-[7px] border border-[#e8eaf0] px-[11px] text-[11px] text-[#303345] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="text" maxlength="120" placeholder="Ej. Bebida" />
        </label>
        <label class="mb-3 block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Cantidad</span>
          <input v-model.number="inventoryQuantity" class="h-10 w-full rounded-[7px] border border-[#e8eaf0] px-[11px] text-[11px] text-[#303345] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="1" placeholder="0" />
        </label>
        <label class="mb-3 block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Precio neto del producto</span>
          <span class="flex h-10 items-center rounded-[7px] border border-[#e8eaf0] px-[11px] focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
            <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
            <input v-model.number="inventoryNetPrice" class="w-full border-0 bg-transparent text-[11px] text-[#303345] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="0.01" placeholder="0.00" />
          </span>
        </label>
        <label class="block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Precio de venta</span>
          <span class="flex h-10 items-center rounded-[7px] border border-[#e8eaf0] px-[11px] focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
            <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
            <input v-model.number="inventorySalePrice" class="w-full border-0 bg-transparent text-[11px] text-[#303345] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="0.01" placeholder="0.00" />
          </span>
        </label>
        <p v-if="inventoryError" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ inventoryError }}</p>
        <div class="mt-4 flex gap-2">
          <button type="submit" :disabled="isLoadingInventory || Boolean(inventoryLoadError) || isSavingInventory" class="inline-flex min-h-10 flex-1 cursor-pointer items-center justify-center rounded-[7px] bg-[#7966de] px-3 text-[10px] font-semibold text-white transition hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2 disabled:cursor-wait disabled:opacity-60">
            {{ isSavingInventory ? 'Guardando…' : editingInventoryId === null ? 'Agregar artículo' : 'Guardar cambios' }}
          </button>
          <button v-if="editingInventoryId !== null" type="button" class="min-h-10 cursor-pointer rounded-[7px] border border-[#e6e8ef] bg-white px-3 text-[10px] font-semibold text-[#565a6c] hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" @click="resetInventoryForm">Cancelar</button>
        </div>
      </form>

      <div class="overflow-hidden rounded-xl border border-[#eceef4] bg-white shadow-[0_3px_12px_#282b3c05]">
        <div class="flex items-center justify-between border-b border-[#eff0f5] px-4 py-3">
          <h2 class="text-[13px] font-bold text-[#282b3b]">Artículos</h2>
          <div class="flex items-center gap-2">
            <span class="text-[9px] text-[#9296a7]">{{ inventoryItems.length }} {{ inventoryItems.length === 1 ? 'artículo' : 'artículos' }}</span>
            <button type="button" aria-label="Sacar reporte de todo el inventario" title="Sacar reporte de todo el inventario" :disabled="isLoadingInventory || Boolean(inventoryLoadError)" class="inline-flex min-h-8 cursor-pointer items-center justify-center gap-1.5 rounded-[7px] bg-[#282b3d] px-2.5 text-[9px] font-semibold text-white transition hover:bg-[#393c50] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] disabled:cursor-wait disabled:opacity-60" @click="downloadInventoryReport">
              <svg class="size-[13px]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5v11m0 0 4-4m-4 4-4-4M4.5 16v4h15v-4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
              Sacar reporte
            </button>
          </div>
        </div>
        <p v-if="inventoryReportError" class="border-b border-[#eff0f5] px-4 py-2 text-[10px] text-[#c34e59]" role="alert">{{ inventoryReportError }}</p>
        <div v-if="inventoryLoadError" class="px-4 py-4 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">
          <p>{{ inventoryLoadError }}</p>
          <button type="button" class="mt-2 cursor-pointer rounded-md border border-[#f1dfe1] bg-white px-2 py-1 font-semibold text-[#b14d59]" @click="loadInventory">Reintentar</button>
        </div>
        <div v-else-if="isLoadingInventory" class="flex min-h-[180px] items-center justify-center px-5 text-[10px] text-[#898d9c]" role="status">Cargando inventario…</div>
        <div v-else-if="inventoryItems.length" class="overflow-x-auto">
          <table class="w-full min-w-[700px] border-collapse text-left text-[10px]">
            <thead class="bg-[#fafafd] text-[#777b8b]">
              <tr>
                <th class="px-3 py-2 font-semibold">Nombre</th>
                <th class="px-3 py-2 text-right font-semibold">Cantidad</th>
                <th class="px-3 py-2 text-right font-semibold">Precio neto</th>
                <th class="px-3 py-2 text-right font-semibold">Precio venta</th>
                <th class="px-3 py-2 text-right font-semibold">Valor neto</th>
                <th class="px-3 py-2 text-right font-semibold">Valor venta</th>
                <th class="px-3 py-2 text-right font-semibold">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in inventoryItems" :key="item.id" class="border-t border-[#f0f1f5] text-[#55596a]">
                <th class="px-3 py-3 font-semibold">{{ item.name }}</th>
                <td class="px-3 py-3 text-right">{{ item.quantity }}</td>
                <td class="px-3 py-3 text-right">${{ formatAmount(item.net_price) }}</td>
                <td class="px-3 py-3 text-right">${{ formatAmount(item.sale_price) }}</td>
                <td class="px-3 py-3 text-right font-semibold">${{ formatAmount(item.quantity * item.net_price) }}</td>
                <td class="px-3 py-3 text-right font-semibold">${{ formatAmount(item.quantity * item.sale_price) }}</td>
                <td class="px-3 py-2">
                  <div class="flex justify-end gap-1">
                    <button type="button" class="min-h-8 cursor-pointer rounded-md border border-[#e6e8ef] bg-white px-2 text-[9px] font-semibold text-[#565a6c] hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" :aria-label="`Editar ${item.name}`" @click="editInventoryItem(item)">Editar</button>
                    <button type="button" :disabled="deletingInventoryId === item.id" class="min-h-8 cursor-pointer rounded-md border border-[#f1dfe1] bg-white px-2 text-[9px] font-semibold text-[#b14d59] hover:bg-[#fff5f5] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] disabled:cursor-wait disabled:opacity-60" :aria-label="`Eliminar ${item.name}`" @click="removeInventoryItem(item)">{{ deletingInventoryId === item.id ? 'Eliminando…' : 'Eliminar' }}</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="flex min-h-[180px] flex-col items-center justify-center px-5 text-center">
          <strong class="text-[11px] text-[#55596a]">Tu inventario está vacío</strong>
          <p class="mb-0 mt-1 text-[9px] text-[#999dac]">Agrega un artículo desde el formulario para comenzar.</p>
        </div>
      </div>
    </section>
    <p class="mt-4 text-center text-[9px] text-[#9296a4]">Los artículos se guardan en la base de datos del inventario.</p>
  </section>
</template>
