<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 mb-4">Inventory</h1>

    <p v-if="loading" class="text-gray-600">Loading inventory...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden divide-y divide-gray-200">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Stock</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Inventory Value</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Action</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="(row, i) in inventory" :key="row.product_id || i" class="hover:bg-gray-50">
            <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ row.stock }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatCurrency(inventoryValue(row)) }}</td>
            <td class="px-4 py-2">
              <span v-if="(row.stock ?? 0) < 10" class="text-red-600 font-semibold">LOW STOCK</span>
              <span v-else class="text-green-600">OK</span>
            </td>
            <td class="px-4 py-2 text-right space-x-2">
              <button
                type="button"
                class="text-sm px-2 py-1 border border-gray-300 rounded hover:bg-gray-100"
                @click="quickAdjust(row, 10)"
              >
                +10
              </button>
              <button
                type="button"
                class="text-sm px-2 py-1 border border-gray-300 rounded hover:bg-gray-100"
                @click="quickAdjust(row, -10)"
              >
                -10
              </button>
              <button
                type="button"
                class="text-sm text-slate-600 hover:text-slate-800 font-medium"
                @click="openAdjustModal(row)"
              >
                Adjust
              </button>
            </td>
          </tr>
          <tr v-if="!inventory?.length">
            <td colspan="5" class="px-4 py-4 text-sm text-gray-500 text-center">No inventory data.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Adjust Stock modal -->
    <div
      v-if="showAdjustModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showAdjustModal = false"
    >
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-2">Adjust Stock</h2>
        <p v-if="adjustProductName" class="text-sm text-gray-600 mb-4">{{ adjustProductName }}</p>
        <form @submit.prevent="handleAdjustStock">
          <div class="space-y-4">
            <div>
              <label for="adj-quantity" class="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
              <input
                id="adj-quantity"
                v-model.number="adjustForm.quantity"
                type="number"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. 10 or -5"
              />
              <p class="text-xs text-gray-500 mt-1">Use positive to add, negative to subtract.</p>
            </div>
            <div>
              <label for="adj-type" class="block text-sm font-medium text-gray-700 mb-1">Type</label>
              <select
                id="adj-type"
                v-model="adjustForm.type"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="adjustment">Adjustment</option>
                <option value="purchase">Purchase</option>
                <option value="return">Return</option>
              </select>
            </div>
            <div>
              <label for="adj-reference" class="block text-sm font-medium text-gray-700 mb-1">Reference</label>
              <input
                id="adj-reference"
                v-model="adjustForm.reference"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. manual adjustment"
              />
            </div>
          </div>
          <p v-if="adjustError" class="text-sm text-red-600 mt-2">{{ adjustError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button type="button" class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md" @click="showAdjustModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="adjustSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ adjustSaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getInventory, adjustStock } from '../api/inventory'

const inventory = ref([])
const loading = ref(true)
const error = ref(null)

const showAdjustModal = ref(false)
const adjustProductId = ref('')
const adjustProductName = ref('')
const adjustForm = ref({ quantity: 0, type: 'adjustment', reference: '' })
const adjustSaving = ref(false)
const adjustError = ref('')

function formatCurrency(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value ?? 0)
}

function inventoryValue(row) {
  const stock = row?.stock ?? 0
  const cost = row?.cost_price ?? 0
  return stock * cost
}

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const data = await getInventory()
    inventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = 'Failed to load inventory.'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

function openAdjustModal(row) {
  adjustProductId.value = row.product_id ?? row.product_Id ?? ''
  adjustProductName.value = row.product_name ?? row.product_Name ?? ''
  adjustForm.value = { quantity: 0, type: 'adjustment', reference: '' }
  adjustError.value = ''
  showAdjustModal.value = true
}

async function quickAdjust(row, delta) {
  const productId = row.product_id ?? row.product_Id ?? ''
  if (!productId) return
  const quantity = delta > 0 ? delta : -delta
  try {
    await adjustStock(productId, {
      quantity: delta > 0 ? quantity : -quantity,
      type: 'adjustment',
      reference: 'quick adjust',
    })
    await loadData()
  } catch (err) {
    error.value = err.response?.data?.error ?? 'Failed to adjust stock.'
  }
}

async function handleAdjustStock() {
  if (!adjustProductId.value) return
  adjustError.value = ''
  adjustSaving.value = true
  try {
    await adjustStock(adjustProductId.value, {
      quantity: Number(adjustForm.value.quantity),
      type: adjustForm.value.type,
      reference: adjustForm.value.reference || undefined,
    })
    showAdjustModal.value = false
    await loadData()
  } catch (err) {
    adjustError.value = err.response?.data?.error ?? 'Failed to adjust stock.'
  } finally {
    adjustSaving.value = false
  }
}
</script>
