<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 mb-4">Dashboard</h1>

    <p v-if="loading" class="text-gray-600">Loading dashboard...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <template v-else>
      <!-- Low Stock Alert -->
      <div v-if="lowStock.length > 0" class="mb-6">
        <router-link to="/inventory" class="block bg-amber-50 border border-amber-200 rounded-lg p-4 hover:bg-amber-100 transition">
          <h2 class="text-lg font-semibold text-amber-800 mb-2 flex items-center gap-2">
            <span>⚠</span> Low Stock
          </h2>
          <ul class="text-sm text-amber-900 space-y-1">
            <li v-for="p in lowStock" :key="p.product_id">{{ p.product_name }} ({{ p.stock }})</li>
          </ul>
          <p class="text-xs text-amber-700 mt-2">Click to open Inventory →</p>
        </router-link>
      </div>

      <!-- Metric cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
        <div class="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p class="text-sm font-medium text-gray-500 uppercase tracking-wide">Total Sales</p>
          <p class="text-2xl font-semibold text-gray-800 mt-1">
            {{ formatCurrency(salesSummary?.total_sales) }}
          </p>
        </div>
        <div class="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p class="text-sm font-medium text-gray-500 uppercase tracking-wide">Total Transactions</p>
          <p class="text-2xl font-semibold text-gray-800 mt-1">
            {{ salesSummary?.total_transactions ?? 0 }}
          </p>
        </div>
      </div>

      <!-- Top Products table -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-3">Top Products</h2>
        <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
                <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity Sold</th>
                <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(row, i) in topProducts" :key="i" class="hover:bg-gray-50">
                <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ row.quantity_sold }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatCurrency(row.revenue) }}</td>
              </tr>
              <tr v-if="!topProducts?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Inventory table -->
      <div>
        <h2 class="text-lg font-semibold text-gray-800 mb-3">Inventory Summary</h2>
        <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
                <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Stock</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(row, i) in inventorySummary" :key="i" class="hover:bg-gray-50">
                <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ row.stock }}</td>
              </tr>
              <tr v-if="!inventorySummary?.length">
                <td colspan="2" class="px-4 py-4 text-sm text-gray-500 text-center">No data</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSalesSummary, getTopProducts, getInventorySummary } from '../api/reports'
import { getLowStock } from '../api/inventory'
import { formatCurrency } from '../utils'

const loading = ref(true)
const error = ref(null)
const salesSummary = ref(null)
const topProducts = ref([])
const inventorySummary = ref([])
const lowStock = ref([])

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const [sales, products, inventory, low] = await Promise.all([
      getSalesSummary(),
      getTopProducts(),
      getInventorySummary(),
      getLowStock().catch(() => []),
    ])
    salesSummary.value = sales
    topProducts.value = Array.isArray(products) ? products : []
    inventorySummary.value = Array.isArray(inventory) ? inventory : []
    lowStock.value = Array.isArray(low) ? low : []
  } catch (err) {
    error.value = 'Failed to load dashboard data.'
  } finally {
    loading.value = false
  }
})
</script>
