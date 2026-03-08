<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-white">Products</h1>
      <div class="flex gap-2">
        <button
          type="button"
          class="px-3 py-2 text-sm bg-slate-600 text-white rounded-md hover:bg-slate-700"
          @click="showCategoryModal = true"
        >
          Add Category
        </button>
        <button
          type="button"
          class="px-3 py-2 text-sm bg-slate-700 text-white rounded-md hover:bg-slate-600"
          @click="openProductModal()"
        >
          Add Product
        </button>
      </div>
    </div>

    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading...</p>
    <p v-else-if="error" class="text-red-600 dark:text-red-400">{{ error }}</p>

    <div v-else>
      <!-- Mobile: Card layout -->
      <div class="sm:hidden space-y-3">
        <div v-for="p in products" :key="p.id" class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="flex items-start justify-between gap-2 mb-2">
            <h3 class="font-medium text-gray-900 dark:text-white truncate flex-1">{{ p.name }}</h3>
            <span class="text-xs px-2 py-0.5 rounded-full shrink-0"
              :class="p.status === 'active' ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300' : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'">
              {{ p.status }}
            </span>
          </div>
          <div class="flex flex-wrap gap-x-4 gap-y-1 text-sm text-gray-600 dark:text-gray-400 mb-3">
            <span v-if="p.sku">SKU: {{ p.sku }}</span>
            <span>{{ p.unit || 'pcs' }}</span>
            <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatCurrency(p.sell_price) }}</span>
          </div>
          <div class="flex gap-2 pt-2 border-t border-gray-100 dark:border-gray-800">
            <button type="button" class="flex-1 py-1.5 text-sm text-slate-600 dark:text-slate-400 border border-slate-300 dark:border-gray-600 rounded-md" @click="openEditModal(p)">Edit</button>
            <button type="button" class="flex-1 py-1.5 text-sm text-red-600 dark:text-red-400 border border-red-200 dark:border-red-800 rounded-md" @click="confirmDeleteProduct(p)">Delete</button>
          </div>
        </div>
        <p v-if="!products?.length" class="py-8 text-sm text-gray-500 dark:text-gray-400 text-center">No products yet. Add one above.</p>
      </div>

      <!-- Desktop: Table layout -->
      <div class="hidden sm:block bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product Name</th>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">SKU</th>
              <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Sell Price</th>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Unit</th>
              <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Status</th>
              <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="p in products" :key="p.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ p.name }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ p.sku }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ formatCurrency(p.sell_price) }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ p.unit || 'pcs' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ p.status }}</td>
              <td class="px-4 py-2 text-right">
                <button type="button" class="text-sm text-slate-600 dark:text-slate-400 hover:underline mr-2" @click="openEditModal(p)">Edit</button>
                <button type="button" class="text-sm text-red-600 dark:text-red-400 hover:underline" @click="confirmDeleteProduct(p)">Delete</button>
              </td>
            </tr>
            <tr v-if="!products?.length">
              <td colspan="6" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No products yet. Add one above.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Category modal -->
    <div
      v-if="showCategoryModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showCategoryModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">Add Category</h2>
        <form @submit.prevent="handleCreateCategory">
          <div class="mb-4">
            <label for="cat-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category name</label>
            <input
              id="cat-name"
              v-model="categoryForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              placeholder="e.g. Beverages"
            />
          </div>
          <p v-if="categoryError" class="text-sm text-red-600 dark:text-red-400 mb-2">{{ categoryError }}</p>
          <div class="flex gap-2 justify-end">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showCategoryModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="categorySaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ categorySaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add Product modal -->
    <div
      v-if="showProductModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showProductModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-md border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">Add Product</h2>
        <form @submit.prevent="handleCreateProduct">
          <div class="space-y-4">
            <div>
              <label for="prod-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Name</label>
              <input
                id="prod-name"
                v-model="productForm.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="Product name"
              />
            </div>
            <div>
              <label for="prod-sku" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">SKU</label>
              <input
                id="prod-sku"
                v-model="productForm.sku"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. SKU-001"
              />
            </div>
            <div>
              <label for="prod-category" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label>
              <select
                id="prod-category"
                v-model="productForm.category_id"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="">— None —</option>
                <option v-for="c in categories" :key="categoryId(c)" :value="categoryId(c)">{{ categoryName(c) }}</option>
              </select>
            </div>
            <div>
              <label for="prod-cost" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Cost price</label>
              <input
                id="prod-cost"
                v-model.number="productForm.cost_price"
                type="number"
                min="0"
                step="0.01"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
            <div>
              <label for="prod-sell" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sell price</label>
              <input
                id="prod-sell"
                v-model.number="productForm.sell_price"
                type="number"
                min="0"
                step="0.01"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
            <div>
              <label for="prod-unit" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Unit</label>
              <select id="prod-unit" v-model="productForm.unit" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500">
                <option value="pcs">pcs</option>
                <option value="kg">kg</option>
                <option value="liter">liter</option>
                <option value="box">box</option>
                <option value="pack">pack</option>
                <option value="bottle">bottle</option>
                <option value="sachet">sachet</option>
                <option value="meter">meter</option>
                <option value="dozen">dozen</option>
              </select>
            </div>
            <div>
              <label for="prod-barcode" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Barcode (optional)</label>
              <input
                id="prod-barcode"
                v-model="productForm.barcode"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. 8991002101010"
              />
            </div>
            <div>
              <label for="prod-threshold" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Low stock threshold</label>
              <input
                id="prod-threshold"
                v-model.number="productForm.low_stock_threshold"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="10"
              />
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Alert when stock is at or below this value.</p>
            </div>
          </div>
          <p v-if="productError" class="text-sm text-red-600 mt-2">{{ productError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showProductModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="productSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ productSaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Product modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showEditModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-md border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">Edit Product</h2>
        <form v-if="editForm" @submit.prevent="handleUpdateProduct">
          <div class="space-y-4">
            <div>
              <label for="edit-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Name</label>
              <input
                id="edit-name"
                v-model="editForm.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="Product name"
              />
            </div>
            <div>
              <label for="edit-sku" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">SKU</label>
              <input
                id="edit-sku"
                v-model="editForm.sku"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. SKU-001"
              />
            </div>
            <div>
              <label for="edit-category" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label>
              <select
                id="edit-category"
                v-model="editForm.category_id"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="">— None —</option>
                <option v-for="c in categories" :key="categoryId(c)" :value="categoryId(c)">{{ categoryName(c) }}</option>
              </select>
            </div>
            <div>
              <label for="edit-cost" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Cost price</label>
              <input
                id="edit-cost"
                v-model.number="editForm.cost_price"
                type="number"
                min="0"
                step="0.01"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
            <div>
              <label for="edit-sell" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sell price</label>
              <input
                id="edit-sell"
                v-model.number="editForm.sell_price"
                type="number"
                min="0"
                step="0.01"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
            <div>
              <label for="edit-unit" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Unit</label>
              <select id="edit-unit" v-model="editForm.unit" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500">
                <option value="pcs">pcs</option>
                <option value="kg">kg</option>
                <option value="liter">liter</option>
                <option value="box">box</option>
                <option value="pack">pack</option>
                <option value="bottle">bottle</option>
                <option value="sachet">sachet</option>
                <option value="meter">meter</option>
                <option value="dozen">dozen</option>
              </select>
            </div>
            <div>
              <label for="edit-barcode" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Barcode</label>
              <input
                id="edit-barcode"
                v-model="editForm.barcode"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. 8991002101010"
              />
            </div>
            <div>
              <label for="edit-threshold" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Low stock threshold</label>
              <input
                id="edit-threshold"
                v-model.number="editForm.low_stock_threshold"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="10"
              />
            </div>
          </div>
          <p v-if="editError" class="text-sm text-red-600 mt-2">{{ editError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showEditModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="editSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ editSaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete confirm -->
    <div
      v-if="showDeleteConfirm"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showDeleteConfirm = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Delete Product</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-4">Delete this product?</p>
        <p v-if="deleteProductName" class="text-sm font-medium text-gray-800 dark:text-white mb-4">"{{ deleteProductName }}"</p>
        <p v-if="deleteError" class="text-sm text-red-600 mb-2">{{ deleteError }}</p>
        <div class="flex gap-2 justify-end">
          <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showDeleteConfirm = false">
            Cancel
          </button>
          <button type="button" class="px-3 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50" :disabled="deleteSaving" @click="handleDeleteProduct">
            {{ deleteSaving ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProducts, getProduct, createProduct, updateProduct, deleteProduct, getCategories, createCategory, addBarcode } from '../api/products'
import { formatCurrency } from '../utils'

const products = ref([])
const categories = ref([])
const loading = ref(true)
const error = ref(null)

const showCategoryModal = ref(false)
const showProductModal = ref(false)
const showEditModal = ref(false)
const showDeleteConfirm = ref(false)
const categoryForm = ref({ name: '' })
const productForm = ref({
  name: '',
  sku: '',
  category_id: '',
  cost_price: 0,
  sell_price: 0,
  unit: 'pcs',
  barcode: '',
  low_stock_threshold: 10,
})
const editForm = ref(null)
const editProductId = ref('')
const editInitialBarcode = ref('')
const categorySaving = ref(false)
const productSaving = ref(false)
const editSaving = ref(false)
const deleteSaving = ref(false)
const categoryError = ref('')
const productError = ref('')
const editError = ref('')
const deleteError = ref('')
const deleteProductId = ref('')
const deleteProductName = ref('')

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const [prods, cats] = await Promise.all([getProducts(), getCategories()])
    products.value = Array.isArray(prods) ? prods : []
    categories.value = Array.isArray(cats) ? cats : []
  } catch (err) {
    error.value = 'Failed to load products.'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

// Backend may return ID/Name (Go) or id/name (json tag)
function categoryId(c) {
  return c?.id ?? c?.ID ?? ''
}
function categoryName(c) {
  return c?.name ?? c?.Name ?? ''
}

function openProductModal() {
  productForm.value = { name: '', sku: '', category_id: '', cost_price: 0, sell_price: 0, unit: 'pcs', barcode: '', low_stock_threshold: 10 }
  productError.value = ''
  showProductModal.value = true
}

async function openEditModal(p) {
  editError.value = ''
  editProductId.value = p.id ?? p.ID ?? ''
  try {
    const product = await getProduct(editProductId.value)
    const barcode = product.barcode ?? product.Barcode ?? ''
    editForm.value = {
      name: product.name ?? product.Name ?? '',
      sku: product.sku ?? product.SKU ?? '',
      category_id: product.category_id ?? product.CategoryID ?? '',
      cost_price: product.cost_price ?? product.CostPrice ?? 0,
      sell_price: product.sell_price ?? product.SellPrice ?? 0,
      status: product.status ?? product.Status ?? 'active',
      unit: product.unit ?? product.Unit ?? 'pcs',
      barcode,
      low_stock_threshold: product.low_stock_threshold ?? 10,
    }
    editInitialBarcode.value = barcode
    showEditModal.value = true
  } catch (err) {
    editError.value = err.response?.data?.error ?? 'Failed to load product.'
  }
}

function confirmDeleteProduct(p) {
  deleteProductId.value = p.id ?? p.ID ?? ''
  deleteProductName.value = p.name ?? p.Name ?? ''
  deleteError.value = ''
  showDeleteConfirm.value = true
}

async function handleCreateCategory() {
  categoryError.value = ''
  categorySaving.value = true
  try {
    await createCategory({ name: categoryForm.value.name })
    showCategoryModal.value = false
    categoryForm.value.name = ''
    await loadData()
  } catch (err) {
    categoryError.value = err.response?.data?.error ?? 'Failed to create category.'
  } finally {
    categorySaving.value = false
  }
}

async function handleCreateProduct() {
  productError.value = ''
  productSaving.value = true
  try {
    const payload = {
      name: productForm.value.name,
      sku: productForm.value.sku || undefined,
      cost_price: Number(productForm.value.cost_price) || 0,
      sell_price: Number(productForm.value.sell_price) || 0,
      unit: productForm.value.unit || 'pcs',
      low_stock_threshold: Math.max(0, Number(productForm.value.low_stock_threshold) || 10),
    }
    if (productForm.value.category_id) payload.category_id = productForm.value.category_id
    const created = await createProduct(payload)
    const productId = created?.id ?? created?.ID
    if (productId && (productForm.value.barcode || '').trim()) {
      try {
        await addBarcode(productId, productForm.value.barcode.trim())
      } catch (e) {
        productError.value = e.response?.data?.error ?? 'Product created but barcode failed.'
        productSaving.value = false
        return
      }
    }
    showProductModal.value = false
    await loadData()
  } catch (err) {
    productError.value = err.response?.data?.error ?? 'Failed to create product.'
  } finally {
    productSaving.value = false
  }
}

async function handleUpdateProduct() {
  if (!editForm.value || !editProductId.value) return
  editError.value = ''
  editSaving.value = true
  try {
    const payload = {
      name: editForm.value.name,
      sku: editForm.value.sku ?? '',
      cost_price: Number(editForm.value.cost_price) || 0,
      sell_price: Number(editForm.value.sell_price) || 0,
      status: editForm.value.status || 'active',
      unit: editForm.value.unit || 'pcs',
      low_stock_threshold: Math.max(0, Number(editForm.value.low_stock_threshold) ?? 10),
    }
    if (editForm.value.category_id) payload.category_id = editForm.value.category_id
    else payload.category_id = null
    await updateProduct(editProductId.value, payload)
    const newBarcode = (editForm.value.barcode || '').trim()
    if (newBarcode && newBarcode !== editInitialBarcode.value) {
      try {
        await addBarcode(editProductId.value, newBarcode)
      } catch (e) {
        editError.value = e.response?.data?.error ?? 'Product updated but barcode failed.'
        editSaving.value = false
        return
      }
    }
    showEditModal.value = false
    await loadData()
  } catch (err) {
    editError.value = err.response?.data?.error ?? 'Failed to update product.'
  } finally {
    editSaving.value = false
  }
}

async function handleDeleteProduct() {
  deleteError.value = ''
  deleteSaving.value = true
  try {
    await deleteProduct(deleteProductId.value)
    showDeleteConfirm.value = false
    await loadData()
  } catch (err) {
    deleteError.value = err.response?.data?.error ?? 'Failed to delete product.'
  } finally {
    deleteSaving.value = false
  }
}
</script>
