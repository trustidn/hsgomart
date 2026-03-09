<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-200">Superadmin Users</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors">Tambah Superadmin</button>
    </div>

    <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">Kelola akun superadmin: ubah email, password, atau tambah superadmin baru.</p>

    <div v-if="loading" class="text-gray-400 dark:text-gray-500 py-8 text-center">Loading...</div>
    <div v-else class="bg-white dark:bg-gray-900 rounded-lg shadow overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th class="px-4 py-3 text-left text-gray-600 dark:text-gray-400">Nama</th>
            <th class="px-4 py-3 text-left text-gray-600 dark:text-gray-400">Email</th>
            <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Status</th>
            <th class="px-4 py-3 text-left text-gray-600 dark:text-gray-400">Dibuat</th>
            <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="border-t dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
            <td class="px-4 py-3 font-medium text-gray-800 dark:text-gray-200">{{ u.name }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400">{{ u.email }}</td>
            <td class="px-4 py-3 text-center">
              <span class="px-2 py-0.5 rounded text-xs font-medium" :class="statusClass(u.status)">{{ u.status }}</span>
            </td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400">{{ u.created_at }}</td>
            <td class="px-4 py-3 text-center">
              <button @click="openEdit(u)" class="text-indigo-600 dark:text-indigo-400 hover:underline text-xs">Edit</button>
            </td>
          </tr>
          <tr v-if="!users.length">
            <td colspan="5" class="px-4 py-8 text-center text-gray-400 dark:text-gray-500">Belum ada superadmin.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-w-md p-6 space-y-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Tambah Superadmin</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Nama</label>
              <input v-model="createForm.name" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="Nama lengkap" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Email</label>
              <input v-model="createForm.email" type="email" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="email@example.com" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Password</label>
              <input v-model="createForm.password" type="password" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="Min. 8 karakter, huruf besar + angka" />
            </div>
          </div>
          <p v-if="createError" class="text-xs text-red-600">{{ createError }}</p>
          <div class="flex gap-3">
            <button @click="showCreateModal = false" class="flex-1 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">Batal</button>
            <button @click="submitCreate" :disabled="submitting" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
              {{ submitting ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Edit Modal -->
    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showEditModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-w-md p-6 space-y-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Edit Superadmin</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Nama</label>
              <input v-model="editForm.name" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Email</label>
              <input v-model="editForm.email" type="email" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Password Baru (kosongkan jika tidak diubah)</label>
              <input v-model="editForm.password" type="password" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="Min. 8 karakter, huruf besar + angka" />
            </div>
          </div>
          <p v-if="editError" class="text-xs text-red-600">{{ editError }}</p>
          <div class="flex gap-3">
            <button @click="showEditModal = false" class="flex-1 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">Batal</button>
            <button @click="submitEdit" :disabled="submitting" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
              {{ submitting ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listSuperadmins, createSuperadmin, updateSuperadmin } from '../../api/admin'

const loading = ref(true)
const users = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editId = ref(null)
const submitting = ref(false)
const createError = ref('')
const editError = ref('')
const createForm = ref({ name: '', email: '', password: '' })
const editForm = ref({ name: '', email: '', password: '' })

function statusClass(s) {
  const map = { active: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' }
  return map[s] || 'bg-gray-100 text-gray-600 dark:text-gray-400'
}

async function load() {
  loading.value = true
  try {
    const data = await listSuperadmins()
    users.value = Array.isArray(data) ? data : []
  } catch (e) {
    users.value = []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  createForm.value = { name: '', email: '', password: '' }
  createError.value = ''
  showCreateModal.value = true
}

async function submitCreate() {
  if (!createForm.value.name || !createForm.value.email || !createForm.value.password) {
    createError.value = 'Nama, email, dan password wajib diisi'
    return
  }
  submitting.value = true
  createError.value = ''
  try {
    await createSuperadmin(createForm.value)
    showCreateModal.value = false
    await load()
  } catch (e) {
    createError.value = e.response?.data?.error ?? 'Gagal menambah superadmin'
  } finally {
    submitting.value = false
  }
}

function openEdit(u) {
  editId.value = u.id
  editForm.value = { name: u.name, email: u.email, password: '' }
  editError.value = ''
  showEditModal.value = true
}

async function submitEdit() {
  if (!editForm.value.name || !editForm.value.email) {
    editError.value = 'Nama dan email wajib diisi'
    return
  }
  submitting.value = true
  editError.value = ''
  try {
    const payload = { name: editForm.value.name, email: editForm.value.email }
    if (editForm.value.password) payload.password = editForm.value.password
    await updateSuperadmin(editId.value, payload)
    showEditModal.value = false
    await load()
  } catch (e) {
    editError.value = e.response?.data?.error ?? 'Gagal memperbarui superadmin'
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>
