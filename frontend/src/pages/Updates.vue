<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Update Terbaru</h1>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Memuat...</p>
    <p v-else-if="!items.length" class="text-gray-500 dark:text-gray-400 py-8 text-center">Belum ada update.</p>

    <div v-else class="space-y-4">
      <div v-for="item in items" :key="item.id"
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-5">
        <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">{{ item.created_at }}</p>
        <h3 class="font-semibold text-gray-900 dark:text-white">{{ item.title }}</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-2 whitespace-pre-wrap">{{ item.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAllUpdates } from '../api/admin'

const items = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await getAllUpdates()
    items.value = Array.isArray(data) ? data : []
  } catch { /* ignore */ }
  loading.value = false
})
</script>
