<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Admin Dashboard</h1>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white rounded-lg shadow p-5">
        <p class="text-sm text-gray-500">Total Tenants</p>
        <p class="text-3xl font-bold text-gray-800">{{ stats.total_tenants }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5">
        <p class="text-sm text-gray-500">Total Transactions</p>
        <p class="text-3xl font-bold text-gray-800">{{ stats.total_transactions }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-5">
        <p class="text-sm text-gray-500">Total Revenue</p>
        <p class="text-3xl font-bold text-gray-800">Rp {{ Number(stats.total_revenue || 0).toLocaleString() }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStats } from '../../api/admin'

const loading = ref(true)
const stats = ref({})

onMounted(async () => {
  try { stats.value = await getStats() }
  catch { stats.value = {} }
  finally { loading.value = false }
})
</script>
