<script setup lang="ts">
defineProps<{
  src: string
  title: string
  filename: string
}>()

defineEmits<{
  close: []
}>()
</script>

<template>
  <div class="fixed inset-0 z-[100] flex flex-col bg-[#171923]/80 p-3 backdrop-blur-sm sm:p-6" role="dialog" aria-modal="true" :aria-label="`Vista previa: ${title}`">
    <section class="mx-auto flex min-h-0 w-full max-w-6xl flex-1 flex-col overflow-hidden rounded-xl bg-white shadow-2xl">
      <header class="flex min-h-14 items-center justify-between gap-3 border-b border-[#eceef4] px-4 sm:px-5">
        <h2 class="truncate text-sm font-semibold text-[#282b3b]">{{ title }}</h2>
        <div class="flex shrink-0 items-center gap-2">
          <a
            :href="src"
            :download="filename"
            class="inline-flex min-h-9 items-center rounded-md bg-[#7966de] px-3 text-xs font-semibold text-white transition hover:bg-[#6955d0] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]"
          >
            Descargar
          </a>
          <button
            type="button"
            class="inline-flex min-h-9 items-center rounded-md border border-[#e6e8ef] bg-white px-3 text-xs font-semibold text-[#565a6c] transition hover:bg-[#f8f8fb] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[#b4a8ff]"
            aria-label="Cerrar vista previa"
            @click="$emit('close')"
          >
            Cerrar
          </button>
        </div>
      </header>
      <iframe
        class="min-h-0 w-full flex-1 bg-[#e6e7eb]"
        :src="src"
        :title="title"
      />
    </section>
  </div>
</template>
