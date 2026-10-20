<script setup lang="ts">
interface InventoryItem {
  id: number
  name: string
  quantity: number
  price: number
}

const inventoryItems = ref<InventoryItem[]>([])
const inventoryName = ref('')
const inventoryQuantity = ref<number | null>(null)
const inventoryPrice = ref<number | null>(null)
const inventoryError = ref('')
const editingInventoryId = ref<number | null>(null)
let nextInventoryId = 0

const inventoryQuantityTotal = computed(() =>
  inventoryItems.value.reduce((total, item) => total + item.quantity, 0)
)
const inventoryValueTotal = computed(() =>
  inventoryItems.value.reduce((total, item) => total + item.quantity * item.price, 0)
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
  inventoryPrice.value = null
  editingInventoryId.value = null
  inventoryError.value = ''
}

function saveInventoryItem() {
  inventoryError.value = ''
  const name = inventoryName.value.trim()
  const quantity = inventoryQuantity.value
  const price = inventoryPrice.value

  if (!name) {
    inventoryError.value = 'Escribe el nombre del artículo.'
    return
  }
  if (typeof quantity !== 'number' || !Number.isInteger(quantity) || quantity < 0) {
    inventoryError.value = 'La cantidad debe ser un número entero igual o mayor que cero.'
    return
  }
  if (typeof price !== 'number' || !Number.isFinite(price) || price < 0) {
    inventoryError.value = 'El precio debe ser un importe válido igual o mayor que cero.'
    return
  }

  if (editingInventoryId.value !== null) {
    const itemIndex = inventoryItems.value.findIndex(item => item.id === editingInventoryId.value)
    if (itemIndex === -1) {
      inventoryError.value = 'No se encontró el artículo que intentas editar.'
      return
    }
    inventoryItems.value[itemIndex] = {
      id: editingInventoryId.value,
      name,
      quantity,
      price
    }
  } else {
    inventoryItems.value.push({
      id: ++nextInventoryId,
      name,
      quantity,
      price
    })
  }

  resetInventoryForm()
}

function editInventoryItem(item: InventoryItem) {
  editingInventoryId.value = item.id
  inventoryName.value = item.name
  inventoryQuantity.value = item.quantity
  inventoryPrice.value = item.price
  inventoryError.value = ''
}

function removeInventoryItem(item: InventoryItem) {
  if (!window.confirm(`¿Eliminar "${item.name}" del inventario?`)) {
    return
  }
  inventoryItems.value = inventoryItems.value.filter(existingItem => existingItem.id !== item.id)
  if (editingInventoryId.value === item.id) {
    resetInventoryForm()
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

    <section class="mb-4 grid grid-cols-3 gap-[12px] max-[520px]:grid-cols-1 max-[520px]:gap-2" aria-label="Resumen del inventario">
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">ARTÍCULOS</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">{{ inventoryItems.length }}</div>
      </article>
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">UNIDADES EN EXISTENCIA</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">{{ inventoryQuantityTotal }}</div>
      </article>
      <article class="rounded-xl border border-[#eceef4] bg-white p-[15px] shadow-[0_3px_12px_#282b3c05]">
        <div class="text-[9px] font-bold tracking-[.7px] text-[#9296a7]">VALOR DEL INVENTARIO</div>
        <div class="mt-2 text-[22px] font-bold text-[#292c3d]">${{ formatAmount(inventoryValueTotal) }}</div>
      </article>
    </section>

    <section class="grid grid-cols-[minmax(0,1fr)_minmax(0,1.5fr)] items-start gap-[15px] max-[780px]:grid-cols-1" aria-label="Administrar artículos">
      <form class="rounded-xl border border-[#eceef4] bg-white p-[19px] shadow-[0_3px_12px_#282b3c05]" @submit.prevent="saveInventoryItem">
        <h2 class="mb-1 text-[15px] font-bold text-[#282b3b]">{{ editingInventoryId === null ? 'Agregar artículo' : 'Editar artículo' }}</h2>
        <p class="mb-4 text-[10px] text-[#898d9c]">Registra nombre, cantidad disponible y precio unitario.</p>
        <label class="mb-3 block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Nombre</span>
          <input v-model="inventoryName" class="h-10 w-full rounded-[7px] border border-[#e8eaf0] px-[11px] text-[11px] text-[#303345] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="text" maxlength="120" placeholder="Ej. Bebida" />
        </label>
        <label class="mb-3 block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Cantidad</span>
          <input v-model.number="inventoryQuantity" class="h-10 w-full rounded-[7px] border border-[#e8eaf0] px-[11px] text-[11px] text-[#303345] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="1" placeholder="0" />
        </label>
        <label class="block">
          <span class="mb-[6px] block text-[10px] font-semibold text-[#55596a]">Precio unitario</span>
          <span class="flex h-10 items-center rounded-[7px] border border-[#e8eaf0] px-[11px] focus-within:border-[#b3a8f1] focus-within:shadow-[0_0_0_2px_#eeeaff]">
            <span class="mr-2 text-[11px] text-[#a3a7b5]">$</span>
            <input v-model.number="inventoryPrice" class="w-full border-0 bg-transparent text-[11px] text-[#303345] outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" type="number" min="0" step="0.01" placeholder="0.00" />
          </span>
        </label>
        <p v-if="inventoryError" class="mt-3 text-[10px] leading-[1.5] text-[#c34e59]" role="alert">{{ inventoryError }}</p>
        <div class="mt-4 flex gap-2">
          <button type="submit" class="inline-flex min-h-10 flex-1 cursor-pointer items-center justify-center rounded-[7px] bg-[#7966de] px-3 text-[10px] font-semibold text-white transition hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff] focus-visible:outline-offset-2">
            {{ editingInventoryId === null ? 'Agregar artículo' : 'Guardar cambios' }}
          </button>
          <button v-if="editingInventoryId !== null" type="button" class="min-h-10 cursor-pointer rounded-[7px] border border-[#e6e8ef] bg-white px-3 text-[10px] font-semibold text-[#565a6c] hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" @click="resetInventoryForm">Cancelar</button>
        </div>
      </form>

      <div class="overflow-hidden rounded-xl border border-[#eceef4] bg-white shadow-[0_3px_12px_#282b3c05]">
        <div class="flex items-center justify-between border-b border-[#eff0f5] px-4 py-3">
          <h2 class="text-[13px] font-bold text-[#282b3b]">Artículos</h2>
          <span class="text-[9px] text-[#9296a7]">{{ inventoryItems.length }} {{ inventoryItems.length === 1 ? 'artículo' : 'artículos' }}</span>
        </div>
        <div v-if="inventoryItems.length" class="overflow-x-auto">
          <table class="w-full min-w-[490px] border-collapse text-left text-[10px]">
            <thead class="bg-[#fafafd] text-[#777b8b]">
              <tr>
                <th class="px-3 py-2 font-semibold">Nombre</th>
                <th class="px-3 py-2 text-right font-semibold">Cantidad</th>
                <th class="px-3 py-2 text-right font-semibold">Precio</th>
                <th class="px-3 py-2 text-right font-semibold">Valor</th>
                <th class="px-3 py-2 text-right font-semibold">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in inventoryItems" :key="item.id" class="border-t border-[#f0f1f5] text-[#55596a]">
                <th class="px-3 py-3 font-semibold">{{ item.name }}</th>
                <td class="px-3 py-3 text-right">{{ item.quantity }}</td>
                <td class="px-3 py-3 text-right">${{ formatAmount(item.price) }}</td>
                <td class="px-3 py-3 text-right font-semibold">${{ formatAmount(item.quantity * item.price) }}</td>
                <td class="px-3 py-2">
                  <div class="flex justify-end gap-1">
                    <button type="button" class="min-h-8 cursor-pointer rounded-md border border-[#e6e8ef] bg-white px-2 text-[9px] font-semibold text-[#565a6c] hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" :aria-label="`Editar ${item.name}`" @click="editInventoryItem(item)">Editar</button>
                    <button type="button" class="min-h-8 cursor-pointer rounded-md border border-[#f1dfe1] bg-white px-2 text-[9px] font-semibold text-[#b14d59] hover:bg-[#fff5f5] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]" :aria-label="`Eliminar ${item.name}`" @click="removeInventoryItem(item)">Eliminar</button>
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
    <p class="mt-4 text-center text-[9px] text-[#9296a4]">El inventario se mantiene solo mientras esta página esté abierta.</p>
  </section>
</template>
