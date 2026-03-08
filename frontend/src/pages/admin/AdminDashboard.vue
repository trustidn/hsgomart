<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-6">Admin Dashboard</h1>
    <div v-if="loading" class="text-gray-500 dark:text-gray-400">Loading...</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <p class="text-sm text-gray-500 dark:text-gray-400">Total Tenants</p>
        <p class="text-3xl font-bold text-gray-800 dark:text-gray-200">{{ stats.total_tenants }}</p>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <p class="text-sm text-gray-500 dark:text-gray-400">Active Tenants</p>
        <p class="text-3xl font-bold text-green-600">{{ stats.active_tenants }}</p>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <p class="text-sm text-gray-500 dark:text-gray-400">Subscription Revenue</p>
        <p class="text-3xl font-bold text-gray-800 dark:text-gray-200">Rp {{ Number(stats.total_revenue || 0).toLocaleString('id-ID') }}</p>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <p class="text-sm text-gray-500 dark:text-gray-400">Pending Orders</p>
        <p class="text-3xl font-bold" :class="stats.pending_orders > 0 ? 'text-amber-600' : 'text-gray-800 dark:text-gray-200'">{{ stats.pending_orders }}</p>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <p class="text-sm text-gray-500 dark:text-gray-400">Expiring in 7 Days</p>
        <p class="text-3xl font-bold" :class="stats.expiring_in_7d > 0 ? 'text-red-600' : 'text-gray-800 dark:text-gray-200'">{{ stats.expiring_in_7d }}</p>
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
