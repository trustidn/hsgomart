<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Documentation</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors">+ Add Article</button>
    </div>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Loading...</p>

    <div v-else class="space-y-3">
      <div v-if="!docs.length" class="text-gray-500 dark:text-gray-400 py-8 text-center">No documentation articles yet.</div>
      <div v-for="doc in docs" :key="doc.id"
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-5">
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xs px-2 py-0.5 rounded-full font-medium"
                :class="doc.is_published ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300' : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'">
                {{ doc.is_published ? 'Published' : 'Draft' }}
              </span>
              <span class="text-xs text-gray-400 dark:text-gray-500">Order: {{ doc.sort_order }}</span>
            </div>
            <h3 class="font-semibold text-gray-900 dark:text-white">{{ doc.title }}</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 line-clamp-2">{{ doc.content?.substring(0, 200) }}{{ doc.content?.length > 200 ? '...' : '' }}</p>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <button @click="openEdit(doc)" class="text-indigo-600 dark:text-indigo-400 hover:underline text-sm">Edit</button>
            <button @click="handleDelete(doc.id)" class="text-red-600 dark:text-red-400 hover:underline text-sm">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">{{ editId ? 'Edit Article' : 'New Article' }}</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Title</label>
              <input v-model="form.title" type="text" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Content</label>
              <textarea v-model="form.content" rows="12" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none font-mono" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sort Order</label>
                <input v-model.number="form.sort_order" type="number" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none" />
              </div>
              <div class="flex items-end pb-1">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="form.is_published" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
                  <span class="text-sm text-gray-700 dark:text-gray-300">Published</span>
                </label>
              </div>
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
import { listDocumentation, createDocumentation, updateDocumentation, deleteDocumentation } from '../../api/admin'

const docs = ref([])
const loading = ref(true)
const showModal = ref(false)
const editId = ref(null)
const saving = ref(false)
const formError = ref('')
const form = ref({ title: '', content: '', sort_order: 0, is_published: true })

async function load() {
  try {
    const data = await listDocumentation()
    docs.value = Array.isArray(data) ? data : []
  } catch { /* ignore */ }
  loading.value = false
}

function openCreate() {
  editId.value = null
  form.value = { title: '', content: '', sort_order: docs.value.length, is_published: true }
  formError.value = ''
  showModal.value = true
}

function openEdit(doc) {
  editId.value = doc.id
  form.value = { title: doc.title, content: doc.content, sort_order: doc.sort_order, is_published: doc.is_published }
  formError.value = ''
  showModal.value = true
}

async function handleSave() {
  if (!form.value.title.trim()) { formError.value = 'Title is required'; return }
  saving.value = true
  formError.value = ''
  try {
    if (editId.value) {
      await updateDocumentation(editId.value, form.value)
    } else {
      await createDocumentation(form.value)
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
  if (!confirm('Delete this article?')) return
  try {
    await deleteDocumentation(id)
    await load()
  } catch { /* ignore */ }
}

onMounted(load)
</script>
