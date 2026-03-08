<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <h1 class="text-2xl font-bold text-gray-800">Platform Settings</h1>

    <div v-if="loading" class="text-gray-400 py-8 text-center">Loading...</div>
    <template v-else>
      <!-- Logo -->
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h2 class="text-base font-semibold text-gray-800 mb-4">Platform Logo</h2>
        <div class="flex items-center gap-6">
          <div class="w-20 h-20 rounded-xl bg-gray-100 border border-gray-200 flex items-center justify-center overflow-hidden shrink-0">
            <img v-if="form.logo_url" :src="logoSrc" class="w-full h-full object-cover" alt="Logo" />
            <span v-else class="text-2xl font-bold text-gray-400">{{ (form.saas_name || 'S').charAt(0) }}</span>
          </div>
          <div>
            <label class="cursor-pointer inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
              Upload Logo
              <input type="file" class="hidden" accept=".png,.jpg,.jpeg,.webp" @change="handleLogo" />
            </label>
            <p class="text-xs text-gray-400 mt-1">PNG, JPG, WebP. Max 2MB.</p>
          </div>
        </div>
      </div>

      <!-- General -->
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h2 class="text-base font-semibold text-gray-800 mb-4">General</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Platform Name</label>
            <input v-model="form.saas_name" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Tagline</label>
            <input v-model="form.tagline" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
        </div>
      </div>

      <!-- Bank Info -->
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h2 class="text-base font-semibold text-gray-800 mb-4">Payment Information</h2>
        <p class="text-xs text-gray-400 mb-4">This will be displayed to tenants on the subscription payment page.</p>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Bank Name</label>
            <input v-model="form.bank_name" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" placeholder="BCA" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Account Number</label>
            <input v-model="form.bank_account" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" placeholder="1234567890" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Account Holder</label>
            <input v-model="form.bank_holder" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" placeholder="PT Example" />
          </div>
        </div>
      </div>

      <!-- Contact -->
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h2 class="text-base font-semibold text-gray-800 mb-4">Contact Information</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Email</label>
            <input v-model="form.contact_email" type="email" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Phone</label>
            <input v-model="form.contact_phone" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
          </div>
        </div>
      </div>

      <!-- Save -->
      <div class="flex justify-end">
        <button @click="save" :disabled="saving" class="px-6 py-2.5 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
          {{ saving ? 'Saving...' : 'Save Settings' }}
        </button>
      </div>

      <p v-if="msg" class="text-sm text-green-600">{{ msg }}</p>
      <p v-if="errMsg" class="text-sm text-red-600">{{ errMsg }}</p>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getSettings, updateSettings, uploadSaasLogo } from '../../api/admin'

const loading = ref(true)
const saving = ref(false)
const msg = ref('')
const errMsg = ref('')

const form = ref({
  saas_name: '',
  tagline: '',
  logo_url: '',
  bank_name: '',
  bank_account: '',
  bank_holder: '',
  contact_email: '',
  contact_phone: '',
})

const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const logoSrc = computed(() => {
  if (!form.value.logo_url) return ''
  return form.value.logo_url.startsWith('http') ? form.value.logo_url : apiBase + form.value.logo_url
})

async function load() {
  try {
    const data = await getSettings()
    Object.keys(form.value).forEach(k => { if (data[k] !== undefined) form.value[k] = data[k] || '' })
  } catch { errMsg.value = 'Failed to load settings' }
  finally { loading.value = false }
}

async function handleLogo(e) {
  const file = e.target.files?.[0]
  if (!file) return
  try {
    const res = await uploadSaasLogo(file)
    form.value.logo_url = res.logo_url
    msg.value = 'Logo uploaded'
    setTimeout(() => msg.value = '', 3000)
  } catch { errMsg.value = 'Failed to upload logo' }
}

async function save() {
  saving.value = true
  msg.value = ''
  errMsg.value = ''
  try {
    await updateSettings(form.value)
    msg.value = 'Settings saved successfully'
    setTimeout(() => msg.value = '', 3000)
  } catch { errMsg.value = 'Failed to save settings' }
  finally { saving.value = false }
}

onMounted(load)
</script>
