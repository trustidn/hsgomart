<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <span class="text-sm text-gray-400">{{ new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }) }}</span>
    </div>

    <p v-if="loading" class="text-gray-500 py-8 text-center">Loading dashboard...</p>
    <p v-else-if="error" class="text-red-600 py-4">{{ error }}</p>

    <template v-else>
      <!-- Subscription Info -->
      <div v-if="sub" class="mb-6">
        <router-link to="/subscription" class="block bg-white rounded-xl border border-gray-200 p-6 hover:border-indigo-300 transition-colors">
          <div class="flex items-center justify-between flex-wrap gap-4">
            <div>
              <p class="text-sm text-gray-500">Current Plan</p>
              <h2 class="text-lg font-semibold text-gray-900">{{ sub.plan_name }}</h2>
              <p class="text-sm mt-1">
                Status:
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="{
                    'bg-amber-100 text-amber-700': sub.status === 'trial',
                    'bg-green-100 text-green-700': sub.status === 'active',
                    'bg-red-100 text-red-700': sub.status === 'expired',
                  }">{{ sub.status }}</span>
              </p>
              <p v-if="subDays !== null" class="text-sm mt-1" :class="subDays <= 7 ? 'text-red-600' : 'text-gray-500'">
                {{ subDays }} days remaining
              </p>
            </div>
            <div class="text-right text-sm text-gray-500">
              <p>Max Users: {{ subPlan?.max_users }}</p>
              <p>Max Products: {{ subPlan?.max_products }}</p>
            </div>
          </div>
        </router-link>
      </div>

      <!-- Low Stock Alert -->
      <div v-if="lowStock.length > 0" class="mb-6">
        <router-link to="/inventory" class="block bg-amber-50 border border-amber-200 rounded-lg p-4 hover:bg-amber-100 transition">
          <h2 class="text-lg font-semibold text-amber-800 mb-2 flex items-center gap-2">
            <span>⚠</span> Low Stock
          </h2>
          <ul class="text-sm text-amber-900 space-y-1">
            <li v-for="p in lowStock" :key="p.product_id">{{ p.product_name }} ({{ p.stock }})</li>
          </ul>
          <p class="text-xs text-amber-700 mt-2">Click to open Inventory</p>
        </router-link>
      </div>

      <!-- Expiring Alert -->
      <div v-if="expiring.length > 0" class="mb-6">
        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
          <h2 class="text-lg font-semibold text-red-800 mb-2">Expiring Soon</h2>
          <ul class="text-sm text-red-900 space-y-1">
            <li v-for="p in expiring" :key="p.batch_id">{{ p.product_name }} - exp {{ p.expired_at }} ({{ p.remaining }} left)</li>
          </ul>
        </div>
      </div>

      <!-- Metric cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
        <div class="bg-white rounded-xl p-5 border border-gray-200">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Total Sales</p>
            <span class="w-9 h-9 rounded-lg bg-indigo-50 flex items-center justify-center">
              <svg class="w-5 h-5 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            </span>
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ formatCurrency(salesSummary?.total_sales) }}</p>
          <p v-if="compare" class="text-sm mt-1" :class="compare.change_pct >= 0 ? 'text-green-600' : 'text-red-600'">
            {{ compare.change_pct >= 0 ? '+' : '' }}{{ compare.change_pct.toFixed(1) }}% vs last month
          </p>
        </div>
        <div class="bg-white rounded-xl p-5 border border-gray-200">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Total Transactions</p>
            <span class="w-9 h-9 rounded-lg bg-emerald-50 flex items-center justify-center">
              <svg class="w-5 h-5 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
            </span>
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ salesSummary?.total_transactions ?? 0 }}</p>
        </div>
      </div>

      <!-- Charts -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
        <div class="bg-white rounded-xl p-5 border border-gray-200">
          <h3 class="text-sm font-semibold text-gray-700 mb-4">Sales (7 days)</h3>
          <canvas ref="salesChartRef" height="200"></canvas>
        </div>
        <div class="bg-white rounded-xl p-5 border border-gray-200">
          <h3 class="text-sm font-semibold text-gray-700 mb-4">Top Products</h3>
          <canvas ref="productsChartRef" height="200"></canvas>
        </div>
      </div>

      <!-- Top Products table -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-3">Top Products</h2>
        <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Qty Sold</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
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
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Stock</th>
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
import { ref, onMounted, nextTick } from 'vue'
import { Chart, registerables } from 'chart.js'
import { getSalesSummary, getTopProducts, getInventorySummary, getSalesDaily, getSalesCompare } from '../api/reports'
import { getLowStock } from '../api/inventory'
import { getSubscription } from '../api/subscription'
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

const salesChartRef = ref(null)
const productsChartRef = ref(null)

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const [sales, products, inventory, low, expData, cmp, subData] = await Promise.all([
      getSalesSummary(),
      getTopProducts(),
      getInventorySummary(),
      getLowStock().catch(() => []),
      client.get('/api/inventory/expiring', { params: { days: 30 } }).then(r => r.data).catch(() => []),
      getSalesCompare().catch(() => null),
      getSubscription().catch(() => null),
    ])
    salesSummary.value = sales
    topProducts.value = Array.isArray(products) ? products : []
    inventorySummary.value = Array.isArray(inventory) ? inventory : []
    lowStock.value = Array.isArray(low) ? low : []
    expiring.value = Array.isArray(expData) ? expData : []
    compare.value = cmp
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
