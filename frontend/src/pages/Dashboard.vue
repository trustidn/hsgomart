<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
      <span class="text-sm text-gray-400 dark:text-gray-500">{{ new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }) }}</span>
    </div>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Loading dashboard...</p>
    <p v-else-if="error" class="text-red-600 dark:text-red-400 py-4">{{ error }}</p>

    <template v-else>
      <!-- Subscription Info -->
      <div v-if="sub" class="mb-6">
        <router-link to="/subscription" class="block bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-6 hover:border-indigo-300 dark:hover:border-indigo-600 transition-colors">
          <div class="flex items-center justify-between flex-wrap gap-4">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Current Plan</p>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ sub.plan_name }}</h2>
              <p class="text-sm mt-1">
                Status:
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="{
                    'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300': sub.status === 'trial',
                    'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300': sub.status === 'active',
                    'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300': sub.status === 'expired',
                  }">{{ sub.status }}</span>
              </p>
              <p v-if="subDays !== null" class="text-sm mt-1" :class="subDays <= 7 ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'">
                {{ subDays }} days remaining
              </p>
            </div>
            <div class="text-right text-sm text-gray-500 dark:text-gray-400">
              <p>Max Users: {{ subPlan?.max_users }}</p>
              <p>Max Products: {{ subPlan?.max_products }}</p>
            </div>
          </div>
        </router-link>
      </div>

      <!-- Latest Updates (right below subscription) -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-semibold text-gray-800 dark:text-white">Update Terbaru</h2>
          <router-link v-if="recentUpdates.length" to="/updates" class="text-sm text-indigo-600 dark:text-indigo-400 hover:underline">Lihat Semua</router-link>
        </div>
        <div v-if="recentUpdates.length" class="space-y-3">
          <div v-for="u in recentUpdates" :key="u.id" class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-lg bg-indigo-50 dark:bg-indigo-900/40 flex items-center justify-center shrink-0 mt-0.5">
                <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-2">
                  <h3 class="font-medium text-sm text-gray-900 dark:text-white">{{ u.title }}</h3>
                  <span class="text-xs text-gray-400 dark:text-gray-500 shrink-0">{{ u.created_at }}</span>
                </div>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ u.content }}</p>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-6 text-center">
          <p class="text-sm text-gray-400 dark:text-gray-500">Belum ada update terbaru dari admin.</p>
        </div>
      </div>

      <!-- Low Stock Alert -->
      <div v-if="lowStock.length > 0" class="mb-6">
        <router-link to="/inventory" class="block bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-4 hover:bg-amber-100 dark:hover:bg-amber-900/30 transition">
          <h2 class="text-lg font-semibold text-amber-800 dark:text-amber-200 mb-2 flex items-center gap-2">
            <span>⚠</span> Low Stock
          </h2>
          <ul class="text-sm text-amber-900 dark:text-amber-100 space-y-1">
            <li v-for="p in lowStock" :key="p.product_id">{{ p.product_name }} ({{ p.stock }})</li>
          </ul>
          <p class="text-xs text-amber-700 dark:text-amber-300 mt-2">Click to open Inventory</p>
        </router-link>
      </div>

      <!-- Expiring Alert -->
      <div v-if="expiring.length > 0" class="mb-6">
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <h2 class="text-lg font-semibold text-red-800 dark:text-red-200 mb-2">Expiring Soon</h2>
          <ul class="text-sm text-red-900 dark:text-red-100 space-y-1">
            <li v-for="p in expiring" :key="p.batch_id">{{ p.product_name }} - exp {{ p.expired_at }} ({{ p.remaining }} left)</li>
          </ul>
        </div>
      </div>

      <!-- Metric cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-900 rounded-xl p-5 border border-gray-200 dark:border-gray-800">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Total Sales</p>
            <span class="w-9 h-9 rounded-lg bg-indigo-50 dark:bg-indigo-900/40 flex items-center justify-center">
              <svg class="w-5 h-5 text-indigo-500 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            </span>
          </div>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-2">{{ formatCurrency(salesSummary?.total_sales) }}</p>
          <p v-if="compare" class="text-sm mt-1" :class="compare.change_pct >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">
            {{ compare.change_pct >= 0 ? '+' : '' }}{{ compare.change_pct.toFixed(1) }}% vs last month
          </p>
        </div>
        <div class="bg-white dark:bg-gray-900 rounded-xl p-5 border border-gray-200 dark:border-gray-800">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Total Transactions</p>
            <span class="w-9 h-9 rounded-lg bg-emerald-50 dark:bg-emerald-900/40 flex items-center justify-center">
              <svg class="w-5 h-5 text-emerald-500 dark:text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
            </span>
          </div>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-2">{{ salesSummary?.total_transactions ?? 0 }}</p>
        </div>
      </div>

      <!-- Charts -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-900 rounded-xl p-5 border border-gray-200 dark:border-gray-800">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">Sales (7 days)</h3>
          <canvas ref="salesChartRef" height="200"></canvas>
        </div>
        <div class="bg-white dark:bg-gray-900 rounded-xl p-5 border border-gray-200 dark:border-gray-800">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">Top Products</h3>
          <canvas ref="productsChartRef" height="200"></canvas>
        </div>
      </div>

      <!-- Top Products table -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-3">Top Products</h2>
        <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product Name</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Qty Sold</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in topProducts" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ row.quantity_sold }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ formatCurrency(row.revenue) }}</td>
              </tr>
              <tr v-if="!topProducts?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Inventory table -->
      <div>
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-3">Inventory Summary</h2>
        <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product Name</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Stock</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in inventorySummary" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ row.stock }}</td>
              </tr>
              <tr v-if="!inventorySummary?.length">
                <td colspan="2" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { Chart, registerables } from 'chart.js'
import { getSalesSummary, getTopProducts, getInventorySummary, getSalesDaily, getSalesCompare } from '../api/reports'
import { getLowStock } from '../api/inventory'
import { getSubscription } from '../api/subscription'
import { getRecentUpdates } from '../api/admin'
import { formatCurrency } from '../utils'
import client from '../api/client'

Chart.register(...registerables)

const loading = ref(true)
const error = ref(null)
const salesSummary = ref(null)
const topProducts = ref([])
const inventorySummary = ref([])
const lowStock = ref([])
const expiring = ref([])
const compare = ref(null)
const sub = ref(null)
const subPlan = ref(null)
const subDays = ref(null)
const recentUpdates = ref([])

const salesChartRef = ref(null)
const productsChartRef = ref(null)

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const [sales, products, inventory, low, expData, cmp, subData, updatesData] = await Promise.all([
      getSalesSummary(),
      getTopProducts(),
      getInventorySummary(),
      getLowStock().catch(() => []),
      client.get('/api/inventory/expiring', { params: { days: 30 } }).then(r => r.data).catch(() => []),
      getSalesCompare().catch(() => null),
      getSubscription().catch(() => null),
      getRecentUpdates().catch(() => []),
    ])
    salesSummary.value = sales
    topProducts.value = Array.isArray(products) ? products : []
    inventorySummary.value = Array.isArray(inventory) ? inventory : []
    lowStock.value = Array.isArray(low) ? low : []
    expiring.value = Array.isArray(expData) ? expData : []
    compare.value = cmp
    recentUpdates.value = Array.isArray(updatesData) ? updatesData : []
    if (subData) {
      sub.value = subData.subscription
      subPlan.value = subData.plan
      subDays.value = subData.days_remaining ?? subData.trial_days_left ?? null
    }

    await nextTick()
    renderCharts(products)
  } catch (err) {
    error.value = 'Failed to load dashboard data.'
  } finally {
    loading.value = false
  }
})

async function renderCharts(products) {
  try {
    const dailyData = await getSalesDaily().catch(() => [])
    const daily = Array.isArray(dailyData) ? dailyData : []
    if (salesChartRef.value && daily.length) {
      new Chart(salesChartRef.value, {
        type: 'line',
        data: {
          labels: daily.map(d => d.date),
          datasets: [{
            label: 'Sales',
            data: daily.map(d => d.total_sales),
            borderColor: '#475569',
            backgroundColor: 'rgba(71,85,105,0.1)',
            fill: true,
            tension: 0.3,
          }],
        },
        options: { responsive: true, plugins: { legend: { display: false } }, scales: { y: { beginAtZero: true } } },
      })
    }
  } catch { /* chart optional */ }

  try {
    const prods = Array.isArray(products) ? products.slice(0, 5) : []
    if (productsChartRef.value && prods.length) {
      new Chart(productsChartRef.value, {
        type: 'bar',
        data: {
          labels: prods.map(p => p.product_name),
          datasets: [{
            label: 'Revenue',
            data: prods.map(p => p.revenue),
            backgroundColor: ['#475569', '#64748b', '#94a3b8', '#cbd5e1', '#e2e8f0'],
          }],
        },
        options: { responsive: true, plugins: { legend: { display: false } }, scales: { y: { beginAtZero: true } } },
      })
    }
  } catch { /* chart optional */ }
}
</script>
