<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-semibold text-gray-800">Products</h1>
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

    <p v-if="loading" class="text-gray-600">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product Name</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">SKU</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Sell Price</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="p in products" :key="p.id" class="hover:bg-gray-50">
            <td class="px-4 py-2 text-sm text-gray-800">{{ p.name }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ p.sku }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatCurrency(p.sell_price) }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ p.status }}</td>
          </tr>
          <tr v-if="!products?.length">
            <td colspan="4" class="px-4 py-4 text-sm text-gray-500 text-center">No products yet. Add one above.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Category modal -->
    <div
      v-if="showCategoryModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showCategoryModal = false"
    >
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Add Category</h2>
        <form @submit.prevent="handleCreateCategory">
          <div class="mb-4">
            <label for="cat-name" class="block text-sm font-medium text-gray-700 mb-1">Category name</label>
            <input
              id="cat-name"
              v-model="categoryForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
              placeholder="e.g. Beverages"
            />
          </div>
          <p v-if="categoryError" class="text-sm text-red-600 mb-2">{{ categoryError }}</p>
          <div class="flex gap-2 justify-end">
            <button type="button" class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md" @click="showCategoryModal = false">
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
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Add Product</h2>
        <form @submit.prevent="handleCreateProduct">
          <div class="space-y-4">
            <div>
              <label for="prod-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <input
                id="prod-name"
                v-model="productForm.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="Product name"
              />
            </div>
            <div>
              <label for="prod-sku" class="block text-sm font-medium text-gray-700 mb-1">SKU</label>
              <input
                id="prod-sku"
                v-model="productForm.sku"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="e.g. SKU-001"
              />
            </div>
            <div>
              <label for="prod-category" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <select
                id="prod-category"
                v-model="productForm.category_id"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
              >
                <option value="">— None —</option>
                <option v-for="c in categories" :key="categoryId(c)" :value="categoryId(c)">{{ categoryName(c) }}</option>
              </select>
            </div>
            <div>
              <label for="prod-cost" class="block text-sm font-medium text-gray-700 mb-1">Cost price</label>
              <input
                id="prod-cost"
                v-model.number="productForm.cost_price"
                type="number"
                min="0"
                step="0.01"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
            <div>
              <label for="prod-sell" class="block text-sm font-medium text-gray-700 mb-1">Sell price</label>
              <input
                id="prod-sell"
                v-model.number="productForm.sell_price"
                type="number"
                min="0"
                step="0.01"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
                placeholder="0"
              />
            </div>
          </div>
          <p v-if="productError" class="text-sm text-red-600 mt-2">{{ productError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button type="button" class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md" @click="showProductModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="productSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ productSaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProducts, createProduct, getCategories, createCategory } from '../api/products'
import { formatCurrency } from '../utils'

const products = ref([])
const categories = ref([])
const loading = ref(true)
const error = ref(null)

const showCategoryModal = ref(false)
const showProductModal = ref(false)
const categoryForm = ref({ name: '' })
const productForm = ref({
  name: '',
  sku: '',
  category_id: '',
  cost_price: 0,
  sell_price: 0,
})
const categorySaving = ref(false)
const productSaving = ref(false)
const categoryError = ref('')
const productError = ref('')

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
  productForm.value = { name: '', sku: '', category_id: '', cost_price: 0, sell_price: 0 }
  productError.value = ''
  showProductModal.value = true
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
    }
    if (productForm.value.category_id) payload.category_id = productForm.value.category_id
    await createProduct(payload)
    showProductModal.value = false
    await loadData()
  } catch (err) {
    productError.value = err.response?.data?.error ?? 'Failed to create product.'
  } finally {
    productSaving.value = false
  }
}
</script>
