<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Documentation</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors">+ Add Article</button>
    </div>

    <p v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Loading...</p>
    <p v-else-if="!docs.length" class="text-gray-500 dark:text-gray-400 py-8 text-center">No documentation articles yet.</p>

    <div v-else class="space-y-3">
      <div v-for="doc in docs" :key="doc.id"
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl overflow-hidden">
        <div class="flex items-center gap-2 px-5 py-4">
          <button
            @click="toggle(doc.id)"
            class="flex-1 flex items-center justify-between text-left hover:opacity-80 transition-opacity min-w-0"
          >
            <div class="flex items-center gap-2 min-w-0">
              <span class="font-semibold text-gray-900 dark:text-white truncate">{{ doc.title }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full font-medium shrink-0"
                :class="doc.is_published ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300' : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'">
                {{ doc.is_published ? 'Published' : 'Draft' }}
              </span>
              <span v-if="doc.visibility === 'admin'" class="text-xs px-2 py-0.5 rounded-full font-medium bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300 shrink-0">Admin Only</span>
            </div>
            <svg :class="['w-5 h-5 text-gray-400 dark:text-gray-500 transition-transform shrink-0 ml-2', open === doc.id ? 'rotate-180' : '']"
              fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>
          <div class="flex items-center gap-1 shrink-0 ml-2">
            <button @click="openEdit(doc)" class="p-1.5 rounded-lg text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors" title="Edit">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
            </button>
            <button @click="handleDelete(doc.id)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors" title="Delete">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
            </button>
          </div>
        </div>
        <div v-show="open === doc.id" class="px-5 pb-5 border-t border-gray-100 dark:border-gray-800">
          <div class="prose prose-sm dark:prose-invert max-w-none pt-4 text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{{ doc.content }}</div>
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
              <textarea v-model="form.content" rows="16" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none font-mono leading-relaxed" />
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sort Order</label>
                <input v-model.number="form.sort_order" type="number" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Visibility</label>
                <select v-model="form.visibility" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 outline-none">
                  <option value="all">All Users</option>
                  <option value="admin">Admin Only</option>
                </select>
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
const open = ref(null)
const showModal = ref(false)
const editId = ref(null)
const saving = ref(false)
const formError = ref('')
const form = ref({ title: '', content: '', sort_order: 0, is_published: true, visibility: 'all' })

function toggle(id) {
  open.value = open.value === id ? null : id
}

async function load() {
  try {
    const data = await listDocumentation()
    docs.value = Array.isArray(data) ? data : []
    if (docs.value.length && open.value === null) open.value = docs.value[0].id
  } catch { /* ignore */ }
  loading.value = false
}

function openCreate() {
  editId.value = null
  form.value = { title: '', content: '', sort_order: docs.value.length, is_published: true, visibility: 'all' }
  formError.value = ''
  showModal.value = true
}

function openEdit(doc) {
  editId.value = doc.id
  form.value = { title: doc.title, content: doc.content, sort_order: doc.sort_order, is_published: doc.is_published, visibility: doc.visibility || 'all' }
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
