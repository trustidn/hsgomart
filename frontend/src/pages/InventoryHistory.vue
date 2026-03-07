<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 mb-4">Inventory History</h1>

    <p v-if="loading" class="text-gray-600">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden divide-y divide-gray-200">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Reference</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="(m, i) in movements" :key="i" class="hover:bg-gray-50">
            <td class="px-4 py-2 text-sm text-gray-800">{{ m.product_name || '—' }}</td>
            <td class="px-4 py-2">
              <span :class="typeClass(m.type)" class="font-medium">{{ m.type || '—' }}</span>
            </td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatQuantity(m.quantity) }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ m.reference || '—' }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ formatDate(m.created_at) }}</td>
          </tr>
          <tr v-if="!movements?.length">
            <td colspan="5" class="px-4 py-4 text-sm text-gray-500 text-center">No movements yet.</td>
          </tr>
        </tbody>
      </table>

      <div v-if="total > 0" class="px-4 py-3 flex items-center justify-between border-t border-gray-200 bg-gray-50">
        <p class="text-sm text-gray-600">
          Menampilkan {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} dari {{ total }}
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
            Sebelumnya
          </button>
          <span class="text-sm text-gray-600">Halaman {{ page }} dari {{ totalPages }}</span>
          <button
            type="button"
            class="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50"
            :disabled="page >= totalPages"
            @click="goToPage(page + 1)"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getMovements } from '../api/inventory'

const movements = ref([])
const total = ref(0)
const loading = ref(true)
const error = ref(null)
const page = ref(1)
const pageSize = ref(20)

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

function formatDate(iso) {
  if (!iso) return '—'
  try {
    const d = new Date(iso)
    return d.toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return iso
  }
}

function formatQuantity(q) {
  const n = Number(q)
  if (n < 0) return n.toString()
  if (n > 0) return `+${n}`
  return '0'
}

function typeClass(type) {
  const t = (type || '').toLowerCase()
  if (t === 'sale') return 'text-red-600'
  if (t === 'purchase') return 'text-green-600'
  if (t === 'adjustment') return 'text-yellow-600'
  return 'text-gray-600'
}

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const data = await getMovements({ limit: pageSize.value, page: page.value })
    movements.value = Array.isArray(data?.movements) ? data.movements : []
    total.value = Number(data?.total) ?? 0
  } catch (err) {
    error.value = 'Failed to load movements.'
  } finally {
    loading.value = false
  }
}

function goToPage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  loadData()
}

onMounted(loadData)
</script>
