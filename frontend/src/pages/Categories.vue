<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-white">Categories</h1>
      <button
        type="button"
        class="px-3 py-2 text-sm bg-slate-600 text-white rounded-md hover:bg-slate-700"
        @click="openAddModal()"
      >
        Add Category
      </button>
    </div>

    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>

    <div v-else class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Category Name</th>
            <th scope="col" class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Action</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="c in categories" :key="catId(c)" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
            <td class="px-4 py-2 text-sm text-gray-800 dark:text-white">{{ catName(c) }}</td>
            <td class="px-4 py-2 text-right space-x-2">
              <button type="button" class="text-sm text-slate-600 hover:underline" @click="openEditModal(c)">
                Edit
              </button>
              <button type="button" class="text-sm text-red-600 hover:underline" @click="confirmDelete(c)">
                Delete
              </button>
            </td>
          </tr>
          <tr v-if="!categories?.length">
            <td colspan="2" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No categories yet. Add one above.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Category modal -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showAddModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">Add Category</h2>
        <form @submit.prevent="handleCreateCategory">
          <div class="mb-4">
            <label for="add-cat-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category name</label>
            <input
              id="add-cat-name"
              v-model="addForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              placeholder="e.g. Beverages"
            />
          </div>
          <p v-if="addError" class="text-sm text-red-600 mb-2">{{ addError }}</p>
          <div class="flex gap-2 justify-end">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showAddModal = false">Cancel</button>
            <button type="submit" :disabled="addSaving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ addSaving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Category modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showEditModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">Edit Category</h2>
        <form @submit.prevent="handleUpdateCategory">
          <div class="mb-4">
            <label for="edit-cat-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category name</label>
            <input
              id="edit-cat-name"
              v-model="editForm.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-slate-500"
              placeholder="e.g. Beverages"
            />
          </div>
          <p v-if="editError" class="text-sm text-red-600 mb-2">{{ editError }}</p>
          <div class="flex gap-2 justify-end">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showEditModal = false">Cancel</button>
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
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Delete Category</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-4">Delete this category?</p>
        <p v-if="deleteCategoryName" class="text-sm font-medium text-gray-800 dark:text-white mb-4">"{{ deleteCategoryName }}"</p>
        <p v-if="deleteError" class="text-sm text-red-600 mb-2">{{ deleteError }}</p>
        <div class="flex gap-2 justify-end">
          <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showDeleteConfirm = false">
            Cancel
          </button>
          <button type="button" class="px-3 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50" :disabled="deleteSaving" @click="handleDeleteCategory">
            {{ deleteSaving ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, createCategory, updateCategory, deleteCategory } from '../api/products'

const categories = ref([])
const loading = ref(true)
const error = ref(null)

const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteConfirm = ref(false)
const addForm = ref({ name: '' })
const editForm = ref({ id: '', name: '' })
const addError = ref('')
const editError = ref('')
const deleteError = ref('')
const addSaving = ref(false)
const editSaving = ref(false)
const deleteSaving = ref(false)
const deleteCategoryId = ref('')
const deleteCategoryName = ref('')

function catId(c) {
  return c?.id ?? c?.ID ?? ''
}
function catName(c) {
  return c?.name ?? c?.Name ?? ''
}

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const data = await getCategories()
    categories.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = 'Failed to load categories.'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

function openAddModal() {
  addForm.value = { name: '' }
  addError.value = ''
  showAddModal.value = true
}

function openEditModal(cat) {
  editForm.value = { id: catId(cat), name: catName(cat) }
  editError.value = ''
  showEditModal.value = true
}

async function handleCreateCategory() {
  addError.value = ''
  addSaving.value = true
  try {
    await createCategory({ name: addForm.value.name })
    showAddModal.value = false
    await loadData()
  } catch (err) {
    addError.value = err.response?.data?.error ?? 'Failed to create category.'
  } finally {
    addSaving.value = false
  }
}

async function handleUpdateCategory() {
  editError.value = ''
  editSaving.value = true
  try {
    await updateCategory(editForm.value.id, { name: editForm.value.name })
    showEditModal.value = false
    await loadData()
  } catch (err) {
    editError.value = err.response?.data?.error ?? 'Failed to update category.'
  } finally {
    editSaving.value = false
  }
}

function confirmDelete(cat) {
  deleteCategoryId.value = catId(cat)
  deleteCategoryName.value = catName(cat)
  deleteError.value = ''
  showDeleteConfirm.value = true
}

async function handleDeleteCategory() {
  deleteError.value = ''
  deleteSaving.value = true
  try {
    await deleteCategory(deleteCategoryId.value)
    showDeleteConfirm.value = false
    await loadData()
  } catch (err) {
    deleteError.value = err.response?.data?.error ?? 'Failed to delete category.'
  } finally {
    deleteSaving.value = false
  }
}
</script>
