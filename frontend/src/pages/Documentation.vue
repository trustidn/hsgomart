<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Dokumentasi</h1>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Memuat dokumentasi...</p>
    <p v-else-if="!docs.length" class="text-gray-500 dark:text-gray-400 py-8 text-center">Belum ada dokumentasi.</p>

    <div v-else class="space-y-3">
      <div v-for="doc in docs" :key="doc.id"
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl overflow-hidden">
        <button
          @click="toggle(doc.id)"
          class="w-full flex items-center justify-between px-5 py-4 text-left hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
        >
          <span class="font-semibold text-gray-900 dark:text-white">{{ doc.title }}</span>
          <svg :class="['w-5 h-5 text-gray-400 dark:text-gray-500 transition-transform', open === doc.id ? 'rotate-180' : '']"
            fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
          </svg>
        </button>
        <div v-show="open === doc.id" class="px-5 pb-5 border-t border-gray-100 dark:border-gray-800">
          <div class="prose prose-sm dark:prose-invert max-w-none pt-4 text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{{ doc.content }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPublicDocumentation } from '../api/admin'

const docs = ref([])
const loading = ref(true)
const open = ref(null)

function toggle(id) {
  open.value = open.value === id ? null : id
}

onMounted(async () => {
  try {
    const data = await getPublicDocumentation()
    docs.value = Array.isArray(data) ? data : []
    if (docs.value.length) open.value = docs.value[0].id
  } catch { /* ignore */ }
  loading.value = false
})
</script>
