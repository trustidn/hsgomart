<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Stock Opname</h1>
      <button @click="startNew" :disabled="starting" class="px-4 py-2 bg-slate-700 text-white rounded hover:bg-slate-600 text-sm disabled:opacity-50">
        {{ starting ? 'Creating...' : '+ New Opname' }}
      </button>
    </div>

    <div v-if="loading" class="text-gray-500">Loading...</div>
    <div v-else>
      <table class="w-full bg-white rounded-lg shadow text-sm">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-gray-600">Date</th>
            <th class="px-4 py-3 text-left text-gray-600">Status</th>
            <th class="px-4 py-3 text-left text-gray-600">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="op in opnames" :key="op.id" class="border-t">
            <td class="px-4 py-3">{{ new Date(op.created_at).toLocaleDateString() }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 rounded text-xs font-medium"
                :class="op.status === 'completed' ? 'bg-green-100 text-green-700' : op.status === 'submitted' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'">
                {{ op.status }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button @click="viewOpname(op.id)" class="text-slate-600 hover:underline text-xs">View</button>
            </td>
          </tr>
          <tr v-if="!opnames.length"><td colspan="3" class="px-4 py-6 text-center text-gray-400">No stock opnames yet</td></tr>
        </tbody>
      </table>
    </div>

    <div v-if="detail" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[80vh] overflow-auto p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold">Opname Detail ({{ detail.opname.status }})</h2>
          <button @click="detail = null" class="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
        </div>
        <table class="w-full text-sm mb-4">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-3 py-2 text-left">Product</th>
              <th class="px-3 py-2 text-right">System</th>
              <th class="px-3 py-2 text-right">Actual</th>
              <th class="px-3 py-2 text-right">Diff</th>
              <th class="px-3 py-2 text-left">Notes</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in detail.items" :key="item.id" class="border-t">
              <td class="px-3 py-2">{{ item.product_id }}</td>
              <td class="px-3 py-2 text-right">{{ item.system_stock }}</td>
              <td class="px-3 py-2 text-right">{{ item.actual_stock }}</td>
              <td class="px-3 py-2 text-right" :class="item.difference < 0 ? 'text-red-600' : item.difference > 0 ? 'text-green-600' : ''">{{ item.difference }}</td>
              <td class="px-3 py-2">{{ item.notes }}</td>
            </tr>
          </tbody>
        </table>
        <button v-if="detail.opname.status === 'submitted'" @click="approve(detail.opname.id)" :disabled="approving" class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-500 text-sm disabled:opacity-50">
          {{ approving ? 'Approving...' : 'Approve & Adjust Stock' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listOpnames, startOpname, getOpname, approveOpname } from '../api/opname'

const loading = ref(true)
const opnames = ref([])
const starting = ref(false)
const detail = ref(null)
const approving = ref(false)

async function load() {
  loading.value = true
  try {
    const data = await listOpnames()
    opnames.value = Array.isArray(data) ? data : []
  } catch { opnames.value = [] }
  finally { loading.value = false }
}

async function startNew() {
  starting.value = true
  try {
    await startOpname()
    await load()
  } finally { starting.value = false }
}

async function viewOpname(id) {
  try {
    detail.value = await getOpname(id)
  } catch { /* ignore */ }
}

async function approve(id) {
  approving.value = true
  try {
    await approveOpname(id)
    detail.value = null
    await load()
  } finally { approving.value = false }
}

onMounted(load)
</script>
