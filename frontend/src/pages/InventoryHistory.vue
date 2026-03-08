<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 mb-4">Inventory History</h1>

    <!-- Filters -->
    <div class="bg-white rounded-lg shadow border border-gray-200 p-4 mb-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 items-end">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Product</label>
          <select v-model="filters.product_id" class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm">
            <option value="">Semua</option>
            <option v-for="p in productOptions" :key="p.product_id" :value="p.product_id">{{ p.product_name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
          <select v-model="filters.type" class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm">
            <option value="">Semua</option>
            <option value="purchase">Purchase</option>
            <option value="sale">Sale</option>
            <option value="adjustment">Adjustment</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Dari tanggal</label>
          <input v-model="filters.from_date" type="date" class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Sampai tanggal</label>
          <input v-model="filters.to_date" type="date" class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm" />
        </div>
        <div class="flex gap-2">
          <button type="button" class="px-3 py-2 bg-slate-600 text-white rounded-md text-sm hover:bg-slate-700" @click="applyFilters">
            Terapkan
          </button>
          <button type="button" class="px-3 py-2 border border-gray-300 rounded-md text-sm hover:bg-gray-100" @click="resetFilters">
            Reset
          </button>
        </div>
      </div>
    </div>

    <p v-if="loading" class="text-gray-600">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden divide-y divide-gray-200">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase cursor-help" title="Stock after this movement">
              Stock After
            </th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Reference</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Reason</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="(m, i) in movements" :key="i" class="hover:bg-gray-50">
            <td class="px-4 py-2 text-sm text-gray-600">{{ formatDate(m.created_at) }}</td>
            <td class="px-4 py-2 text-sm text-gray-800">{{ m.product_name || '—' }}</td>
            <td class="px-4 py-2">
              <span :class="typeBadgeClass(m.type)" class="inline-flex px-2 py-0.5 rounded text-xs font-medium">{{ m.type || '—' }}</span>
            </td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right font-medium">{{ formatQuantity(m.quantity) }}</td>
            <td class="px-4 py-2 text-right" :title="'Stock after this movement'">
              <span :class="(m.stock_after ?? 0) < 0 ? 'text-red-600 font-semibold' : 'text-gray-600'">{{ m.stock_after ?? '—' }}</span>
            </td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ m.reference || '—' }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ m.reason || '—' }}</td>
          </tr>
          <tr v-if="!movements?.length">
            <td colspan="7" class="px-4 py-4 text-sm text-gray-500 text-center">No movements yet.</td>
          </tr>
        </tbody>
      </table>

      <div v-if="total > 0" class="px-4 py-3 flex flex-wrap items-center justify-between gap-2 border-t border-gray-200 bg-gray-50">
        <p class="text-sm text-gray-600">
          Showing {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} of {{ total }}
        </p>
        <div class="flex items-center gap-2">
          <select v-model.number="pageSize" class="text-sm border border-gray-300 rounded px-2 py-1" @change="goToPage(1)">
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
          </select>
          <button
            type="button"
            class="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50"
            :disabled="page <= 1"
            @click="goToPage(page - 1)"
          >
            Previous
          </button>
          <span class="text-sm text-gray-600">Page {{ page }} of {{ totalPages }}</span>
          <button
            type="button"
            class="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50"
            :disabled="page >= totalPages"
            @click="goToPage(page + 1)"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getMovements, getInventory } from '../api/inventory'
import { formatDate } from '../utils'

const movements = ref([])
const total = ref(0)
const loading = ref(true)
const error = ref(null)
const page = ref(1)
const pageSize = ref(20)
const productOptions = ref([])

const filters = ref({
  product_id: '',
  type: '',
  from_date: '',
  to_date: '',
})

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

// formatDate imported from utils

function formatQuantity(q) {
  const n = Number(q)
  if (n < 0) return n.toString()
  if (n > 0) return `+${n}`
  return '0'
}

// Color labels: purchase → green, sale → red, adjustment → yellow
function typeBadgeClass(type) {
  const t = (type || '').toLowerCase()
  if (t === 'purchase') return 'bg-green-100 text-green-800'
  if (t === 'sale') return 'bg-red-100 text-red-800'
  if (t === 'adjustment') return 'bg-yellow-100 text-yellow-800'
  return 'bg-gray-100 text-gray-800'
}

function buildParams() {
  const params = { limit: pageSize.value, page: page.value }
  if (filters.value.product_id) params.product_id = filters.value.product_id
  if (filters.value.type) params.type = filters.value.type
  if (filters.value.from_date) params.from_date = filters.value.from_date
  if (filters.value.to_date) params.to_date = filters.value.to_date
  return params
}

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const data = await getMovements(buildParams())
    movements.value = Array.isArray(data?.movements) ? data.movements : []
    total.value = Number(data?.total) ?? 0
  } catch (err) {
    error.value = 'Failed to load movements.'
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  page.value = 1
  loadData()
}

function resetFilters() {
  filters.value = { product_id: '', type: '', from_date: '', to_date: '' }
  page.value = 1
  loadData()
}

function goToPage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  loadData()
}

async function loadProductOptions() {
  try {
    const data = await getInventory()
    productOptions.value = Array.isArray(data) ? data : []
  } catch {
    productOptions.value = []
  }
}

onMounted(() => {
  loadProductOptions()
  loadData()
})
</script>
