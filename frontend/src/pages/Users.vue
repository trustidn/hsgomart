<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-white">Users</h1>
      <button
        type="button"
        class="px-3 py-2 text-sm bg-slate-600 text-white rounded-md hover:bg-slate-700"
        @click="openCreateModal"
      >
        Add user
      </button>
    </div>

    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading...</p>
    <p v-else-if="error" class="text-red-600 dark:text-red-400">{{ error }}</p>

    <div v-else>
      <div class="sm:hidden space-y-3">
        <div v-for="u in users" :key="u.id" class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <p class="font-medium text-gray-900 dark:text-white">{{ u.name || '—' }}</p>
          <p class="text-sm text-gray-600 dark:text-gray-400 truncate">{{ u.email || '—' }}</p>
          <div class="flex items-center gap-2 mt-2">
            <span class="text-xs px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300">{{ u.role || '—' }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ u.status || '—' }}</span>
          </div>
          <div class="flex gap-2 mt-3">
            <button type="button" class="text-sm text-slate-600 dark:text-slate-400 hover:underline" @click="openEditModal(u)">Edit</button>
            <button type="button" class="text-sm text-red-600 dark:text-red-400 hover:underline disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isCurrentUser(u.id)" @click="confirmDelete(u)">Delete</button>
          </div>
        </div>
        <p v-if="!users?.length" class="py-8 text-sm text-gray-500 dark:text-gray-400 text-center">No users yet.</p>
      </div>
      <div class="hidden sm:block bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Name</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Email</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Role</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Status</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="u in users" :key="u.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ u.name || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ u.email || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ u.role || '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ u.status || '—' }}</td>
              <td class="px-4 py-2 text-right">
                <button type="button" class="text-sm text-slate-600 dark:text-slate-400 hover:underline mr-2" @click="openEditModal(u)">Edit</button>
                <button type="button" class="text-sm text-red-600 dark:text-red-400 hover:underline disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isCurrentUser(u.id)" @click="confirmDelete(u)">Delete</button>
              </td>
            </tr>
            <tr v-if="!users?.length">
              <td colspan="5" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No users yet.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add / Edit user modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showModal = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 w-full max-w-sm border border-gray-200 dark:border-gray-800">
        <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">{{ editingId ? 'Edit user' : 'Add user' }}</h2>
        <form @submit.prevent="submitForm">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Name</label>
              <input v-model="form.name" type="text" required class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="e.g. Kasir 1" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
              <input v-model="form.email" type="email" required class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="kasir@store.com" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Password</label>
              <input v-model="form.password" type="password" :required="!editingId" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="Leave blank to keep current" />
              <p v-if="editingId" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Leave blank to keep current password.</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Role</label>
              <select v-model="form.role" required class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">
                <option value="owner">owner</option>
                <option value="cashier">cashier</option>
              </select>
            </div>
            <div v-if="editingId">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label>
              <select v-model="form.status" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">
                <option value="active">active</option>
                <option value="inactive">inactive</option>
              </select>
            </div>
          </div>
          <p v-if="formError" class="text-sm text-red-600 mt-2">{{ formError }}</p>
          <div class="flex gap-2 justify-end mt-4">
            <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="showModal = false">Cancel</button>
            <button type="submit" :disabled="saving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Save' }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { getUsers, createUser, updateUser, deleteUser } from '../api/users'

const auth = useAuthStore()
const users = ref([])
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const editingId = ref(null)
const saving = ref(false)
const formError = ref('')

const form = ref({
  name: '',
  email: '',
  password: '',
  role: 'cashier',
  status: 'active',
})

const isCurrentUser = (id) => (auth.user?.id && auth.user.id === id) || (auth.user?.id === id)

async function loadUsers() {
  loading.value = true
  error.value = null
  try {
    users.value = await getUsers()
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to load users.'
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingId.value = null
  form.value = { name: '', email: '', password: '', role: 'cashier', status: 'active' }
  formError.value = ''
  showModal.value = true
}

function openEditModal(u) {
  editingId.value = u.id
  form.value = { name: u.name ?? '', email: u.email ?? '', password: '', role: u.role ?? 'cashier', status: u.status ?? 'active' }
  formError.value = ''
  showModal.value = true
}

async function submitForm() {
  formError.value = ''
  saving.value = true
  try {
    if (editingId.value) {
      const payload = { name: form.value.name, email: form.value.email, role: form.value.role, status: form.value.status }
      if (form.value.password) payload.password = form.value.password
      await updateUser(editingId.value, payload)
    } else {
      await createUser({
        name: form.value.name,
        email: form.value.email,
        password: form.value.password,
        role: form.value.role,
      })
    }
    showModal.value = false
    await loadUsers()
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Failed to save.'
  } finally {
    saving.value = false
  }
}

function confirmDelete(u) {
  if (isCurrentUser(u.id)) return
  if (!window.confirm(`Delete user "${u.email || u.name}"?`)) return
  deleteUser(u.id).then(() => loadUsers()).catch((e) => { error.value = e.response?.data?.error ?? 'Failed to delete.' })
}

onMounted(loadUsers)
</script>
