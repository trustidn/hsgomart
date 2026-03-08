<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 dark:text-white mb-4">Inventory</h1>

    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading inventory...</p>
    <p v-else-if="error" class="text-red-600 dark:text-red-400">{{ error }}</p>

    <div v-else class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden divide-y divide-gray-200 dark:divide-gray-700">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product Name</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Stock</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Inventory Value</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Status</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Action</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="(row, i) in inventory" :key="row.product_id || i" class="hover:bg-gray-50 dark:hover:bg-gray-800">
            <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ row.stock }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ formatCurrency(inventoryValue(row)) }}</td>
            <td class="px-4 py-2">
              <span v-if="(row.stock ?? 0) < 10" class="text-red-600 dark:text-red-400 font-semibold">LOW STOCK</span>
              <span v-else class="text-green-600 dark:text-green-400">OK</span>
            </td>
            <td class="px-4 py-2 text-right">
              <button
                type="button"
                class="text-sm text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 font-medium"
                @click="openAdjustModal(row)"
              >
                Adjust
              </button>
            </td>
          </tr>
          <tr v-if="!inventory?.length">
            <td colspan="5" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No inventory data.</td>
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
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Koreksi stok (pengurangan)</h2>
        <p v-if="adjustProductName" class="text-sm text-gray-600 dark:text-gray-400 mb-4">{{ adjustProductName }}</p>
        <p class="text-xs text-amber-700 dark:text-amber-300 bg-amber-50 dark:bg-amber-900/30 border border-amber-200 dark:border-amber-800 rounded px-2 py-1.5 mb-4">Penambahan stok hanya melalui menu Purchase. Adjust hanya untuk pengurangan.</p>
        <form @submit.prevent="handleAdjustStock">
          <div class="space-y-4">
            <div>
              <label for="adj-quantity" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Jumlah pengurangan</label>
              <input
                id="adj-quantity"
                v-model.number="adjustForm.quantity"
                type="number"
                required
                min="1"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="contoh: 8"
              />
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Isi angka jumlah yang akan dikurangi (misal 8).</p>
            </div>
            <div>
              <label for="adj-type" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Type</label>
              <select
                id="adj-type"
                v-model="adjustForm.type"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="adjustment">Adjustment</option>
                <option value="purchase">Purchase</option>
                <option value="return">Return</option>
              </select>
            </div>
            <div>
              <label for="adj-reason" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Alasan (audit)</label>
              <select
                id="adj-reason"
                v-model="adjustForm.reason"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="">— Pilih alasan —</option>
                <option value="expired product">Kadaluarsa</option>
                <option value="damaged item">Barang rusak</option>
                <option value="manual correction">Koreksi manual</option>
                <option value="stock audit">Stock opname</option>
              </select>
            </div>
            <div>
              <label for="adj-reference" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Reference</label>
              <input
                id="adj-reference"
                v-model="adjustForm.reference"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. manual adjustment"
              />
            </div>
          </div>
          <p v-if="adjustError" class="text-sm text-red-600 mt-2">{{ adjustError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button v-if="!showAdjustConfirm" type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showAdjustModal = false">
              Batal
            </button>
            <template v-if="showAdjustConfirm">
              <p class="text-sm text-gray-700 dark:text-gray-300 mr-auto">Yakin kurangi stok <strong>{{ adjustProductName }}</strong> sebanyak <strong>{{ adjustForm.quantity }}</strong>?</p>
              <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showAdjustConfirm = false">
                Batal
              </button>
              <button type="button" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700" @click="confirmAdjustStock">
                Ya, kurangi
              </button>
            </template>
            <button v-else type="submit" :disabled="adjustSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ adjustSaving ? 'Menyimpan...' : 'Lanjut' }}
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
const showAdjustConfirm = ref(false)
const adjustProductId = ref('')
const adjustProductName = ref('')
const adjustForm = ref({ quantity: 0, type: 'adjustment', reason: '', reference: '' })
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
  adjustForm.value = { quantity: 0, type: 'adjustment', reason: '', reference: '' }
  adjustError.value = ''
  showAdjustConfirm.value = false
  showAdjustModal.value = true
}

async function doAdjustStock() {
  const qty = Math.abs(Number(adjustForm.value.quantity)) || 0
  if (qty < 1) {
    adjustError.value = 'Jumlah pengurangan minimal 1.'
    return
  }
  adjustSaving.value = true
  adjustError.value = ''
  try {
    await adjustStock(adjustProductId.value, {
      quantity: -qty,
      type: adjustForm.value.type,
      reason: adjustForm.value.reason || undefined,
      reference: adjustForm.value.reference || undefined,
    })
    showAdjustModal.value = false
    showAdjustConfirm.value = false
    await loadData()
  } catch (err) {
    adjustError.value = err.response?.data?.error ?? 'Gagal mengoreksi stok.'
  } finally {
    adjustSaving.value = false
  }
}

function handleAdjustStock() {
  if (!adjustProductId.value) return
  const qty = Number(adjustForm.value.quantity)
  if (!qty || qty < 1) {
    adjustError.value = 'Isi jumlah pengurangan (angka positif, misal 8).'
    return
  }
  adjustError.value = ''
  showAdjustConfirm.value = true
}

function confirmAdjustStock() {
  doAdjustStock()
}
</script>
