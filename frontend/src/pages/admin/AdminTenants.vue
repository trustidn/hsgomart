<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Tenants</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors">Add Tenant</button>
    </div>

    <!-- Filters -->
    <div class="flex gap-2 mb-4">
      <button v-for="s in statusFilters" :key="s.value" @click="filterStatus = s.value; load()"
        class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors"
        :class="filterStatus === s.value ? 'bg-indigo-600 text-white' : 'bg-white border border-gray-300 text-gray-600 hover:bg-gray-50'">
        {{ s.label }}
      </button>
    </div>

    <div v-if="loading" class="text-gray-400 py-8 text-center">Loading...</div>
    <div v-else class="bg-white rounded-lg shadow overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-gray-600">Name</th>
            <th class="px-4 py-3 text-left text-gray-600">Email</th>
            <th class="px-4 py-3 text-center text-gray-600">Status</th>
            <th class="px-4 py-3 text-left text-gray-600">Plan</th>
            <th class="px-4 py-3 text-center text-gray-600">Subscription</th>
            <th class="px-4 py-3 text-center text-gray-600">Days Left</th>
            <th class="px-4 py-3 text-right text-gray-600">Users</th>
            <th class="px-4 py-3 text-center text-gray-600">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tenants" :key="t.id" class="border-t hover:bg-gray-50">
            <td class="px-4 py-3 font-medium">{{ t.name }}</td>
            <td class="px-4 py-3 text-gray-500">{{ t.email }}</td>
            <td class="px-4 py-3 text-center">
              <span class="px-2 py-0.5 rounded text-xs font-medium" :class="tenantStatusClass(t.status)">{{ t.status }}</span>
            </td>
            <td class="px-4 py-3">{{ t.plan_name }}</td>
            <td class="px-4 py-3 text-center">
              <span v-if="t.sub_status" class="px-2 py-0.5 rounded text-xs font-medium" :class="subStatusClass(t.sub_status)">{{ t.sub_status }}</span>
              <span v-else class="text-xs text-gray-400">—</span>
            </td>
            <td class="px-4 py-3 text-center">
              <span v-if="t.days_remaining !== null && t.days_remaining !== undefined" class="text-xs font-medium" :class="daysClass(t.days_remaining)">
                {{ t.days_remaining }} days
              </span>
              <span v-else class="text-xs text-gray-400">—</span>
            </td>
            <td class="px-4 py-3 text-right">{{ t.user_count }}</td>
            <td class="px-4 py-3 text-center space-x-1">
              <button @click="openEdit(t)" class="text-indigo-600 hover:underline text-xs">Edit</button>
              <button @click="confirmDelete(t)" class="text-red-600 hover:underline text-xs">Delete</button>
            </td>
          </tr>
          <tr v-if="!tenants.length">
            <td colspan="8" class="px-4 py-8 text-center text-gray-400">No tenants found.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-lg p-6 space-y-4">
          <h3 class="text-lg font-semibold text-gray-900">{{ isEdit ? 'Edit Tenant' : 'New Tenant' }}</h3>
          <div class="space-y-3">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-600 mb-1">Name</label>
                <input v-model="form.name" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-600 mb-1">Email</label>
                <input v-model="form.email" type="email" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
              </div>
            </div>
            <div v-if="!isEdit">
              <label class="block text-sm font-medium text-gray-600 mb-1">Password</label>
              <input v-model="form.password" type="password" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" placeholder="Min 8 chars, uppercase, number" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-600 mb-1">Phone</label>
                <input v-model="form.phone" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-600 mb-1">Status</label>
                <select v-model="form.status" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
                  <option value="active">Active</option>
                  <option value="suspended">Suspended</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 mb-1">Subscription Plan</label>
              <select v-model="form.plan_id" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
                <option :value="null">— Trial (default) —</option>
                <option v-for="p in availablePlans" :key="p.id" :value="p.id">{{ p.name }} — Rp {{ Number(p.price).toLocaleString('id-ID') }}</option>
              </select>
            </div>
          </div>
          <p v-if="formError" class="text-xs text-red-600">{{ formError }}</p>
          <div class="flex gap-3">
            <button @click="showModal = false" class="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
            <button @click="submitForm" :disabled="submitting" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
              {{ submitting ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirm -->
    <Teleport to="body">
      <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showDeleteConfirm = false">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-sm p-6 space-y-4">
          <h3 class="text-lg font-semibold text-gray-900">Delete Tenant</h3>
          <p class="text-sm text-gray-600">Are you sure you want to delete <strong>{{ deleteTarget?.name }}</strong>? This will deactivate the tenant and all associated subscriptions.</p>
          <div class="flex gap-3">
            <button @click="showDeleteConfirm = false" class="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
            <button @click="doDelete" :disabled="deleting" class="flex-1 py-2 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700 transition-colors disabled:opacity-50">
              {{ deleting ? 'Deleting...' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listTenants, createTenant, updateTenant, deleteTenant, listPlans } from '../../api/admin'

const loading = ref(true)
const tenants = ref([])
const availablePlans = ref([])
const filterStatus = ref('')
const statusFilters = [
  { label: 'All', value: '' },
  { label: 'Active', value: 'active' },
  { label: 'Suspended', value: 'suspended' },
]

const showModal = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const submitting = ref(false)
const formError = ref('')
const form = ref({ name: '', email: '', password: '', phone: '', status: 'active', plan_id: null })

const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

async function load() {
  loading.value = true
  try {
    const [t, p] = await Promise.all([listTenants(filterStatus.value), listPlans()])
    tenants.value = t || []
    availablePlans.value = (p || []).filter(x => x.is_active)
  } catch { tenants.value = [] }
  finally { loading.value = false }
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { name: '', email: '', password: '', phone: '', status: 'active', plan_id: null }
  formError.value = ''
  showModal.value = true
}

function openEdit(t) {
  isEdit.value = true
  editId.value = t.id
  form.value = { name: t.name, email: t.email, phone: t.phone || '', status: t.status, plan_id: t.plan_id || null }
  formError.value = ''
  showModal.value = true
}

async function submitForm() {
  if (!form.value.name || !form.value.email) { formError.value = 'Name and email required'; return }
  if (!isEdit.value && !form.value.password) { formError.value = 'Password required'; return }
  submitting.value = true
  formError.value = ''
  try {
    if (isEdit.value) {
      const payload = { name: form.value.name, email: form.value.email, phone: form.value.phone, status: form.value.status }
      if (form.value.plan_id) payload.plan_id = form.value.plan_id
      await updateTenant(editId.value, payload)
    } else {
      const payload = { name: form.value.name, email: form.value.email, password: form.value.password }
      if (form.value.plan_id) payload.plan_id = form.value.plan_id
      await createTenant(payload)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Failed to save'
  } finally { submitting.value = false }
}

function confirmDelete(t) {
  deleteTarget.value = t
  showDeleteConfirm.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await deleteTenant(deleteTarget.value.id)
    showDeleteConfirm.value = false
    await load()
  } catch {}
  finally { deleting.value = false }
}

function tenantStatusClass(s) {
  const map = { active: 'bg-green-100 text-green-700', suspended: 'bg-red-100 text-red-700' }
  return map[s] || 'bg-gray-100 text-gray-600'
}

function subStatusClass(s) {
  const map = { active: 'bg-green-100 text-green-700', trial: 'bg-amber-100 text-amber-700', expired: 'bg-red-100 text-red-700' }
  return map[s] || 'bg-gray-100 text-gray-600'
}

function daysClass(d) {
  if (d <= 0) return 'text-red-600'
  if (d <= 7) return 'text-amber-600'
  return 'text-green-600'
}

onMounted(load)
</script>
