<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Platform Updates</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors">+ New Update</button>
    </div>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Loading...</p>

    <div v-else class="space-y-3">
      <div v-if="!items.length" class="text-gray-500 dark:text-gray-400 py-8 text-center">No updates yet.</div>
      <div v-for="item in items" :key="item.id"
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-5">
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">{{ item.created_at }}</p>
            <h3 class="font-semibold text-gray-900 dark:text-white">{{ item.title }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1 whitespace-pre-wrap">{{ item.content }}</p>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <button @click="openEdit(item)" class="text-indigo-600 dark:text-indigo-400 hover:underline text-sm">Edit</button>
            <button @click="handleDelete(item.id)" class="text-red-600 dark:text-red-400 hover:underline text-sm">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-2xl w-full max-w-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">{{ editId ? 'Edit Update' : 'New Update' }}</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Title</label>
              <input v-model="form.title" type="text" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none" placeholder="e.g. Fitur baru: Stock Opname" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Content</label>
              <textarea v-model="form.content" rows="6" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none" placeholder="Describe the update..." />
            </div>
            <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>
            <div class="flex gap-3 pt-2">
              <button @click="showModal = false" class="flex-1 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800">Cancel</button>
              <button @click="handleSave" :disabled="saving" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50">
                {{ saving ? 'Saving...' : 'Save' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listAdminUpdates, createPlatformUpdate, editPlatformUpdate, deletePlatformUpdate } from '../../api/admin'

const items = ref([])
const loading = ref(true)
const showModal = ref(false)
const editId = ref(null)
const saving = ref(false)
const formError = ref('')
const form = ref({ title: '', content: '' })

async function load() {
  try {
    const data = await listAdminUpdates()
    items.value = Array.isArray(data) ? data : []
  } catch { /* ignore */ }
  loading.value = false
}

function openCreate() {
  editId.value = null
  form.value = { title: '', content: '' }
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  editId.value = item.id
  form.value = { title: item.title, content: item.content }
  formError.value = ''
  showModal.value = true
}

async function handleSave() {
  if (!form.value.title.trim() || !form.value.content.trim()) { formError.value = 'Title and content are required'; return }
  saving.value = true
  formError.value = ''
  try {
    if (editId.value) {
      await editPlatformUpdate(editId.value, form.value)
    } else {
      await createPlatformUpdate(form.value)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Failed to save'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (!confirm('Delete this update?')) return
  try {
    await deletePlatformUpdate(id)
    await load()
  } catch { /* ignore */ }
}

onMounted(load)
</script>
