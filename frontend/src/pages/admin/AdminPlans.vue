<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-200">Subscription Plans</h1>
      <button @click="openCreate" class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors">Add Plan</button>
    </div>

    <div v-if="loading" class="text-gray-400 dark:text-gray-500 py-8 text-center">Loading...</div>
    <table v-else class="w-full bg-white dark:bg-gray-900 rounded-lg shadow text-sm">
      <thead class="bg-gray-50 dark:bg-gray-800">
        <tr>
          <th class="px-4 py-3 text-left text-gray-600 dark:text-gray-400">Name</th>
          <th class="px-4 py-3 text-right text-gray-600 dark:text-gray-400">Price</th>
          <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Duration</th>
          <th class="px-4 py-3 text-right text-gray-600 dark:text-gray-400">Max Users</th>
          <th class="px-4 py-3 text-right text-gray-600 dark:text-gray-400">Max Products</th>
          <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Tenants</th>
          <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Status</th>
          <th class="px-4 py-3 text-center text-gray-600 dark:text-gray-400">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="p in plans" :key="p.id" class="border-t dark:border-gray-700">
          <td class="px-4 py-3">
            <div class="font-medium text-gray-800 dark:text-gray-200">{{ p.name }}</div>
            <div v-if="p.description" class="text-xs text-gray-400 mt-0.5">{{ p.description }}</div>
          </td>
          <td class="px-4 py-3 text-right text-gray-800 dark:text-gray-200">{{ p.price === 0 ? 'Free' : 'Rp ' + Number(p.price).toLocaleString('id-ID') }}</td>
          <td class="px-4 py-3 text-center text-gray-800 dark:text-gray-200">{{ formatDuration(p.duration_days) }}</td>
          <td class="px-4 py-3 text-right text-gray-800 dark:text-gray-200">{{ p.max_users }}</td>
          <td class="px-4 py-3 text-right text-gray-800 dark:text-gray-200">{{ p.max_products }}</td>
          <td class="px-4 py-3 text-center text-gray-800 dark:text-gray-200">{{ p.tenant_count }}</td>
          <td class="px-4 py-3 text-center">
            <span class="px-2 py-0.5 rounded text-xs" :class="p.is_active ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400'">
              {{ p.is_active ? 'Active' : 'Inactive' }}
            </span>
          </td>
          <td class="px-4 py-3 text-center space-x-2">
            <button @click="openEdit(p)" class="text-indigo-600 dark:text-indigo-400 hover:underline text-xs">Edit</button>
            <button v-if="p.is_active" @click="toggleActive(p.id, false)" class="text-red-600 dark:text-red-400 hover:underline text-xs">Deactivate</button>
            <button v-else @click="toggleActive(p.id, true)" class="text-green-600 dark:text-green-400 hover:underline text-xs">Activate</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-w-md p-6 space-y-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ editing ? 'Edit Plan' : 'New Plan' }}</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Name</label>
              <input v-model="form.name" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Price (IDR)</label>
                <input v-model.number="form.price" type="number" min="0" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Duration (days)</label>
                <input v-model.number="form.duration_days" type="number" min="1" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Max Users</label>
                <input v-model.number="form.max_users" type="number" min="1" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Max Products</label>
                <input v-model.number="form.max_products" type="number" min="1" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">Description</label>
              <textarea v-model="form.description" rows="2" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"></textarea>
            </div>
          </div>
          <p v-if="formError" class="text-xs text-red-600">{{ formError }}</p>
          <div class="flex gap-3">
            <button @click="showModal = false" class="flex-1 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">Cancel</button>
            <button @click="submitForm" :disabled="submitting" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
              {{ submitting ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPlans, createPlan, updatePlan, deletePlan } from '../../api/admin'

const loading = ref(true)
const plans = ref([])
const showModal = ref(false)
const editing = ref(null)
const submitting = ref(false)
const formError = ref('')

const form = ref({ name: '', price: 0, duration_days: 30, max_users: 5, max_products: 100, description: '' })

async function load() {
  loading.value = true
  try { plans.value = await listPlans() }
  catch { plans.value = [] }
  finally { loading.value = false }
}

function openCreate() {
  editing.value = null
  form.value = { name: '', price: 0, duration_days: 30, max_users: 5, max_products: 100, description: '' }
  formError.value = ''
  showModal.value = true
}

function openEdit(p) {
  editing.value = p.id
  form.value = { name: p.name, price: p.price, duration_days: p.duration_days || 30, max_users: p.max_users, max_products: p.max_products, description: p.description || '' }
  formError.value = ''
  showModal.value = true
}

async function submitForm() {
  if (!form.value.name) { formError.value = 'Name is required'; return }
  submitting.value = true
  formError.value = ''
  try {
    if (editing.value) {
      await updatePlan(editing.value, form.value)
    } else {
      await createPlan(form.value)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Failed to save'
  } finally { submitting.value = false }
}

function formatDuration(days) {
  if (!days) return '30 days'
  if (days === 365) return '1 year'
  if (days === 180) return '6 months'
  if (days === 90) return '3 months'
  if (days === 60) return '2 months'
  if (days === 30) return '1 month'
  return `${days} days`
}

async function toggleActive(id, active) {
  try {
    if (active) {
      await updatePlan(id, { is_active: true })
    } else {
      await deletePlan(id)
    }
    await load()
  } catch {}
}

onMounted(load)
</script>
