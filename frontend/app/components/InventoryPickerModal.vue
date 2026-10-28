<script setup lang="ts">
interface SellableInventoryItem {
  id: string
  name: string
  quantity: number
  sale_price: number
}

const props = defineProps<{
  reservedQuantities: Record<string, number>
}>()

const emit = defineEmits<{
  close: []
  select: [item: SellableInventoryItem, quantity: number]
}>()

const items = ref<SellableInventoryItem[]>([])
const search = ref('')
const selectedItemId = ref('')
const selectedQuantity = ref(1)
const isLoading = ref(true)
const loadError = ref('')
const apiBase = useRuntimeConfig().public.apiBase.replace(/\/+$/, '')

const availableItems = computed(() => items.value.map(item => ({
  ...item,
  availableQuantity: Math.max(0, item.quantity - (props.reservedQuantities[item.id] ?? 0))
})))

const filteredItems = computed(() => {
  const term = search.value.trim().toLocaleLowerCase()
  return availableItems.value.filter(item =>
    item.availableQuantity > 0 && item.name.toLocaleLowerCase().includes(term)
  )
})

const selectedItem = computed(() =>
  availableItems.value.find(item => item.id === selectedItemId.value)
)

onMounted(loadItems)

async function loadItems() {
  isLoading.value = true
  loadError.value = ''
  try {
    items.value = await $fetch<SellableInventoryItem[]>(`${apiBase}/inventory/for-sale`)
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : 'No se pudo cargar el inventario.'
  } finally {
    isLoading.value = false
  }
}

function chooseItem(item: SellableInventoryItem) {
  selectedItemId.value = item.id
  selectedQuantity.value = 1
}

function addSelectedItem() {
  const item = selectedItem.value
  if (!item || !Number.isInteger(selectedQuantity.value) ||
      selectedQuantity.value < 1 || selectedQuantity.value > item.availableQuantity) {
    return
  }
  emit('select', item, selectedQuantity.value)
}
</script>

<template>
  <div class="fixed inset-0 z-[70] grid place-items-center bg-[#17182a80] p-4" @click.self="emit('close')">
    <section
      class="flex max-h-[min(680px,90vh)] w-full max-w-[540px] flex-col overflow-hidden rounded-xl border border-[#eceef4] bg-white shadow-[0_20px_70px_#15162945]"
      role="dialog"
      aria-modal="true"
      aria-labelledby="inventory-picker-title"
    >
      <header class="flex items-start justify-between gap-3 border-b border-[#eff0f5] px-5 py-4">
        <div>
          <div class="text-[9px] font-bold tracking-[1.05px] text-[#9296a7]">IMPORTAR DEL INVENTARIO</div>
          <h2 id="inventory-picker-title" class="mb-1 mt-2 text-[16px] font-bold text-[#282b3b]">Buscar producto</h2>
          <p class="m-0 text-[10px] text-[#898d9c]">El precio de venta se agregará como entrada al movimiento.</p>
        </div>
        <button type="button" class="grid size-8 flex-none cursor-pointer place-items-center rounded-md border border-[#e8eaf0] bg-white text-[#777b8b] hover:bg-[#f8f8fb]" aria-label="Cerrar" @click="emit('close')">×</button>
      </header>

      <div class="border-b border-[#eff0f5] px-5 py-3">
        <label class="block">
          <span class="sr-only">Buscar producto por nombre</span>
          <input v-model="search" autofocus class="h-10 w-full rounded-[7px] border border-[#e8eaf0] px-3 text-[11px] text-[#303345] outline-none focus:border-[#b3a8f1] focus:shadow-[0_0_0_2px_#eeeaff]" type="search" placeholder="Escribe el nombre del producto" />
        </label>
      </div>

      <div v-if="loadError" class="px-5 py-5 text-[10px] text-[#c34e59]" role="alert">
        <p>{{ loadError }}</p>
        <button type="button" class="cursor-pointer rounded-md border border-[#f1dfe1] bg-white px-2 py-1 font-semibold text-[#b14d59]" @click="loadItems">Reintentar</button>
      </div>
      <div v-else-if="isLoading" class="grid min-h-32 place-items-center text-[10px] text-[#898d9c]" role="status">Cargando productos…</div>
      <div v-else-if="filteredItems.length" class="min-h-0 flex-1 overflow-y-auto p-3">
        <button
          v-for="item in filteredItems"
          :key="item.id"
          type="button"
          class="mb-2 flex w-full cursor-pointer items-center justify-between gap-3 rounded-lg border px-3 py-3 text-left transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]"
          :class="selectedItemId === item.id ? 'border-[#a99be9] bg-[#f7f5ff]' : 'border-[#eceef4] bg-white hover:bg-[#fafafd]'"
          @click="chooseItem(item)"
        >
          <span class="min-w-0">
            <strong class="block truncate text-[11px] text-[#3f4253]">{{ item.name }}</strong>
            <small class="mt-1 block text-[9px] text-[#9296a7]">{{ item.availableQuantity }} disponibles</small>
          </span>
          <strong class="flex-none text-[11px] text-[#358f77]">${{ new Intl.NumberFormat('es', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(item.sale_price) }}</strong>
        </button>
      </div>
      <div v-else class="grid min-h-32 place-items-center px-5 text-center text-[10px] text-[#898d9c]">
        {{ search ? 'No hay productos disponibles que coincidan con la búsqueda.' : 'No hay productos con existencias disponibles.' }}
      </div>

      <footer v-if="selectedItem" class="border-t border-[#eff0f5] bg-[#fafafd] p-4">
        <div class="mb-3 flex items-center justify-between gap-3">
          <div class="min-w-0">
            <strong class="block truncate text-[10px] text-[#3f4253]">{{ selectedItem.name }}</strong>
            <small class="text-[9px] text-[#898d9c]">${{ new Intl.NumberFormat('es', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(selectedItem.sale_price) }} por unidad</small>
          </div>
          <label class="flex items-center gap-2 text-[9px] font-semibold text-[#55596a]">
            Cantidad
            <input v-model.number="selectedQuantity" class="h-9 w-[76px] rounded-md border border-[#e8eaf0] bg-white px-2 text-[10px] outline-none focus:border-[#b3a8f1]" type="number" min="1" :max="selectedItem.availableQuantity" step="1" />
          </label>
        </div>
        <button type="button" class="inline-flex min-h-10 w-full cursor-pointer items-center justify-center gap-2 rounded-[7px] bg-[#7966de] px-3 text-[10px] font-semibold text-white transition hover:bg-[#6955d0] disabled:cursor-not-allowed disabled:opacity-50" :disabled="!Number.isInteger(selectedQuantity) || selectedQuantity < 1 || selectedQuantity > selectedItem.availableQuantity" @click="addSelectedItem">
          Agregar al movimiento · ${{ new Intl.NumberFormat('es', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(selectedItem.sale_price * selectedQuantity) }}
        </button>
      </footer>
    </section>
  </div>
</template>
