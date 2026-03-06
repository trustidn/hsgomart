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
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ m.quantity }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ m.reference || '—' }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ formatDate(m.created_at) }}</td>
          </tr>
          <tr v-if="!movements?.length">
            <td colspan="5" class="px-4 py-4 text-sm text-gray-500 text-center">No movements yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getMovements } from '../api/inventory'

const movements = ref([])
const loading = ref(true)
const error = ref(null)

function formatDate(iso) {
  if (!iso) return '—'
  try {
    const d = new Date(iso)
    return d.toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return iso
  }
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
    const data = await getMovements()
    movements.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = 'Failed to load movements.'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
