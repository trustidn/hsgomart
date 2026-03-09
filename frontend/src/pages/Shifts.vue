<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-200 mb-4">Shifts</h1>
    <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">Laporan rekonsiliasi kas (Expected vs Actual) per shift yang sudah ditutup.</p>

    <!-- Date filter -->
    <div class="flex flex-wrap items-center gap-3 mb-4">
      <span class="text-sm text-gray-600 dark:text-gray-400">Periode:</span>
      <input v-model="dateFrom" type="date" class="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
      <span class="text-gray-400">s/d</span>
      <input v-model="dateTo" type="date" class="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
      <button type="button" class="px-3 py-2 text-sm bg-slate-600 text-white rounded-md hover:bg-slate-700" @click="load">Muat</button>
    </div>

    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>
    <div v-else>
      <div class="sm:hidden space-y-3">
        <div v-for="(row, i) in shiftsRows" :key="i" class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ row.date }}</p>
          <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.cashier }}</p>
          <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-sm text-gray-600 dark:text-gray-400 mt-1">
            <span>Open: {{ formatPrice(row.opening) }}</span>
            <span>Sales: {{ formatPrice(row.sales) }}</span>
            <span>Actual: {{ formatPrice(row.actual) }}</span>
          </div>
          <p class="mt-1 text-sm font-medium" :class="diffClass(row.difference)">{{ formatPrice(row.difference) }}</p>
        </div>
        <div v-if="!shiftsRows?.length" class="p-4 rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900">
          <p class="text-sm text-gray-500 dark:text-gray-400 text-center">Tidak ada shift yang ditutup pada periode ini.</p>
        </div>
      </div>
      <div class="hidden sm:block bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Date</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Cashier</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Opening</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Sales</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Expected</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Actual</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Difference</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="(row, i) in shiftsRows" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.date }}</td>
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.cashier }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-600 dark:text-gray-400">{{ formatPrice(row.opening) }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-600 dark:text-gray-400">{{ formatPrice(row.sales) }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-600 dark:text-gray-400">{{ formatPrice(row.expected) }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-600 dark:text-gray-400">{{ formatPrice(row.actual) }}</td>
              <td class="px-4 py-2 text-sm text-right font-medium rounded" :class="diffClass(row.difference)">{{ formatPrice(row.difference) }}</td>
            </tr>
            <tr v-if="!shiftsRows?.length">
              <td colspan="7" class="px-4 py-6 text-sm text-gray-500 dark:text-gray-400 text-center">Tidak ada shift yang ditutup pada periode ini.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { getShiftsReport } from '../api/reports'
import { formatPrice } from '../utils'

const auth = useAuthStore()
const shiftsRows = ref([])
const loading = ref(true)
const error = ref(null)
const dateFrom = ref('')
const dateTo = ref('')

function defaultDateRange() {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  const first = `${y}-${m}-01`
  return { from: first, to: `${y}-${m}-${d}` }
}

function diffClass(diff) {
  if (diff > 0) return 'text-green-700 dark:text-green-400 bg-green-100 dark:bg-green-900/30'
  if (diff < 0) return 'text-red-700 dark:text-red-400 bg-red-100 dark:bg-red-900/30'
  return 'text-gray-800 dark:text-gray-200'
}

async function load() {
  if (!auth.token) {
    error.value = 'Silakan login sebagai owner untuk melihat daftar shift.'
    loading.value = false
    return
  }
  const range = defaultDateRange()
  if (!dateFrom.value) dateFrom.value = range.from
  if (!dateTo.value) dateTo.value = range.to
  loading.value = true
  error.value = null
  try {
    const data = await getShiftsReport({ from: dateFrom.value, to: dateTo.value })
    shiftsRows.value = Array.isArray(data) ? data : []
  } catch (e) {
    const msg = e.response?.data?.error ?? e.message ?? 'Gagal memuat data shift.'
    error.value = msg === 'missing authorization header' ? 'Sesi habis atau belum login. Silakan login lagi.' : msg
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
