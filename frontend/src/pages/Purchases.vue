<template>
  <div>
    <!-- List view -->
    <template v-if="!detailId">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-2xl font-semibold text-gray-800">Purchases</h1>
        <button
          type="button"
          class="px-3 py-2 text-sm bg-slate-600 text-white rounded-md hover:bg-slate-700"
          @click="openCreateModal"
        >
          Create Purchase
        </button>
      </div>

      <p v-if="loading" class="text-gray-600">Loading...</p>
      <p v-else-if="error" class="text-red-600">{{ error }}</p>
      <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Nama Produk</th>

              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Supplier</th>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Invoice</th>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Keterangan</th>
              <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
              <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="p in purchases" :key="p.id" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ formatDate(p.created_at) }}</td>
                            <td class="px-4 py-2 text-sm text-gray-600 max-w-[220px] truncate" :title="(p.product_names || []).join(', ')">{{ (p.product_names && p.product_names.length) ? p.product_names.join(', ') : '—' }}</td>

              <td class="px-4 py-2 text-sm text-gray-800">{{ p.supplier_name || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600">{{ p.invoice_number || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 max-w-[200px] truncate" :title="p.notes || ''">{{ p.notes || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatPrice(p.total_amount) }}</td>
              <td class="px-4 py-2 text-right">
                <router-link :to="`/purchases/${p.id}`" class="text-sm text-slate-600 hover:underline">View</router-link>
              </td>
            </tr>
            <tr v-if="!purchases?.length">
              <td colspan="7" class="px-4 py-4 text-sm text-gray-500 text-center">No purchases yet. Create one above.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Create Purchase modal -->
      <div
        v-if="showCreateModal"
        class="fixed inset-0 z-10 flex items-center justify-center bg-black/50 p-4"
        @click.self="showCreateModal = false"
      >
        <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col">
          <div class="p-4 border-b border-gray-200">
            <h2 class="text-lg font-semibold text-gray-800">Create Purchase</h2>
          </div>
          <form class="flex-1 overflow-auto p-4" @submit.prevent="submitPurchase">
            <div class="grid grid-cols-2 gap-4 mb-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Supplier</label>
                <input v-model="form.supplier_name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="e.g. PT Indofood" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Invoice number</label>
                <input v-model="form.invoice_number" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="e.g. INV-001" />
              </div>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Keterangan</label>
              <textarea v-model="form.notes" rows="2" class="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="Informasi tambahan (opsional)" />
            </div>

            <div class="mb-2 flex items-center justify-between">
              <label class="block text-sm font-medium text-gray-700">Items</label>
              <button type="button" class="text-sm text-slate-600 hover:underline" @click="addItem">+ Add row</button>
            </div>
            <table class="min-w-full border border-gray-200 rounded overflow-hidden mb-4">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
                  <th class="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase w-24">Qty</th>
                  <th class="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase w-32">Cost price</th>
                  <th class="px-3 py-2 text-right text-xs font-medium text-gray-500 uppercase w-32">Subtotal</th>
                  <th class="w-10"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="(row, i) in form.items" :key="i">
                  <td class="px-3 py-2">
                    <select v-model="row.product_id" required class="w-full px-2 py-1.5 border border-gray-300 rounded text-sm" @change="onProductChange(i)">
                      <option value="">Select product</option>
                      <option v-for="prod in products" :key="prod.id" :value="prod.id">{{ prod.name }}</option>
                    </select>
                  </td>
                  <td class="px-3 py-2 text-right">
                    <input v-model.number="row.quantity" type="number" min="1" class="w-full px-2 py-1.5 border border-gray-300 rounded text-sm text-right" @input="updateSubtotal(i)" />
                  </td>
                  <td class="px-3 py-2 text-right">
                    <input v-model.number="row.cost_price" type="number" min="0" step="0.01" class="w-full px-2 py-1.5 border border-gray-300 rounded text-sm text-right" @input="updateSubtotal(i)" />
                  </td>
                  <td class="px-3 py-2 text-sm text-right text-gray-600">{{ formatPrice(row.quantity * row.cost_price) }}</td>
                  <td class="px-2 py-2">
                    <button type="button" class="text-red-600 hover:underline text-sm" @click="removeItem(i)">Remove</button>
                  </td>
                </tr>
              </tbody>
            </table>
            <p v-if="form.items.length === 0" class="text-sm text-amber-600 mb-2">Add at least one item.</p>
            <div class="flex items-center justify-between border-t border-gray-200 pt-4">
              <p class="font-medium text-gray-800">Total purchase value: {{ formatPrice(totalPurchaseValue) }}</p>
              <div class="flex gap-2">
                <button type="button" class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md" @click="showCreateModal = false">Cancel</button>
                <button type="submit" :disabled="createSaving || form.items.length === 0" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">Save</button>
              </div>
            </div>
          </form>
          <p v-if="createError" class="px-4 pb-4 text-sm text-red-600">{{ createError }}</p>
        </div>
      </div>
    </template>

    <!-- Detail view -->
    <template v-else>
      <div class="mb-4">
        <router-link to="/purchases" class="text-sm text-slate-600 hover:underline">&larr; Back to list</router-link>
      </div>
      <p v-if="detailLoading" class="text-gray-600">Loading...</p>
      <p v-else-if="detailError" class="text-red-600">{{ detailError }}</p>
      <template v-else-if="detail">
        <h1 class="text-2xl font-semibold text-gray-800 mb-4">Purchase detail</h1>
        <div class="bg-white rounded-lg shadow border border-gray-200 p-4 mb-6">
          <dl class="grid grid-cols-2 gap-2 text-sm">
            <dt class="text-gray-500">Supplier</dt>
            <dd class="text-gray-800">{{ detail.purchase?.supplier_name || '—' }}</dd>
            <dt class="text-gray-500">Invoice</dt>
            <dd class="text-gray-800">{{ detail.purchase?.invoice_number || '—' }}</dd>
            <dt class="text-gray-500">Date</dt>
            <dd class="text-gray-800">{{ formatDate(detail.purchase?.created_at) }}</dd>
            <dt class="text-gray-500">Total</dt>
            <dd class="text-gray-800 font-medium">{{ formatPrice(detail.purchase?.total_amount) }}</dd>
            <template v-if="detail.purchase?.notes">
              <dt class="text-gray-500">Keterangan</dt>
              <dd class="text-gray-800">{{ detail.purchase.notes }}</dd>
            </template>
          </dl>
        </div>
        <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Cost price</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(item, i) in detail.items" :key="i">
                <td class="px-4 py-2 text-sm text-gray-800">{{ item.product_name || '—' }}</td>
                <td class="px-4 py-2 text-sm text-right">{{ item.quantity }}</td>
                <td class="px-4 py-2 text-sm text-right">{{ formatPrice(item.cost_price) }}</td>
                <td class="px-4 py-2 text-sm text-right">{{ formatPrice(item.subtotal) }}</td>
              </tr>
              <tr v-if="!detail.items?.length">
                <td colspan="4" class="px-4 py-4 text-sm text-gray-500 text-center">No items.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listPurchases, createPurchase, getPurchase } from '../api/purchases'
import { getProducts } from '../api/products'
import { formatPrice, formatDate } from '../utils'

const route = useRoute()
const router = useRouter()

const purchases = ref([])
const loading = ref(false)
const error = ref(null)
const detailId = computed(() => route.params.id || null)
const detail = ref(null)
const detailLoading = ref(false)
const detailError = ref(null)
const products = ref([])
const showCreateModal = ref(false)
const createSaving = ref(false)
const createError = ref(null)

const form = ref({
  supplier_name: '',
  invoice_number: '',
  notes: '',
  items: [{ product_id: '', quantity: 1, cost_price: 0 }],
})

// formatDate imported from utils

const totalPurchaseValue = computed(() => {
  return form.value.items.reduce((sum, row) => sum + (row.quantity || 0) * (row.cost_price || 0), 0)
})

async function loadPurchases() {
  loading.value = true
  error.value = null
  try {
    purchases.value = await listPurchases()
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to load purchases.'
  } finally {
    loading.value = false
  }
}

async function loadDetail() {
  if (!detailId.value) return
  detailLoading.value = true
  detailError.value = null
  try {
    detail.value = await getPurchase(detailId.value)
  } catch (e) {
    detailError.value = e.response?.data?.error ?? 'Failed to load purchase.'
  } finally {
    detailLoading.value = false
  }
}

async function loadProducts() {
  try {
    const data = await getProducts()
    products.value = Array.isArray(data) ? data : []
  } catch {
    products.value = []
  }
}

watch(detailId, (id) => {
  if (id) loadDetail()
  else detail.value = null
}, { immediate: true })

function addItem() {
  form.value.items.push({ product_id: '', quantity: 1, cost_price: 0 })
}

function removeItem(i) {
  form.value.items.splice(i, 1)
}

function onProductChange(i) {
  const prod = products.value.find(p => p.id === form.value.items[i].product_id)
  if (prod != null && (form.value.items[i].cost_price == null || form.value.items[i].cost_price === 0)) {
    form.value.items[i].cost_price = prod.cost_price ?? 0
  }
  updateSubtotal(i)
}

function updateSubtotal() {}

function openCreateModal() {
  form.value = { supplier_name: '', invoice_number: '', notes: '', items: [{ product_id: '', quantity: 1, cost_price: 0 }] }
  createError.value = null
  showCreateModal.value = true
  loadProducts()
}

async function submitPurchase() {
  const items = form.value.items.filter(r => r.product_id && r.quantity > 0 && (r.cost_price ?? 0) >= 0)
  if (items.length === 0) {
    createError.value = 'Add at least one item with product, quantity and cost price.'
    return
  }
  createSaving.value = true
  createError.value = null
  try {
    const payload = {
      supplier_name: form.value.supplier_name,
      invoice_number: form.value.invoice_number,
      notes: form.value.notes || undefined,
      items: items.map(r => ({ product_id: r.product_id, quantity: r.quantity, cost_price: r.cost_price ?? 0 })),
    }
    await createPurchase(payload)
    showCreateModal.value = false
    await loadPurchases()
    router.push('/purchases')
  } catch (e) {
    createError.value = e.response?.data?.error ?? 'Failed to create purchase.'
  } finally {
    createSaving.value = false
  }
}

loadPurchases()
</script>
