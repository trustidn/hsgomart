<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Business Profile</h1>

    <div v-if="loading" class="text-center py-12 text-gray-400">Loading...</div>

    <div v-else class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 divide-y divide-gray-100 dark:divide-gray-800">
      <!-- Logo Section -->
      <div class="p-6">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Logo</label>
        <div class="flex items-center gap-5">
          <div class="w-20 h-20 rounded-xl bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-800 flex items-center justify-center overflow-hidden shrink-0">
            <img v-if="currentLogo" :src="logoSrc" alt="Logo" class="w-full h-full object-cover" />
            <span v-else class="text-gray-400 text-2xl font-bold">{{ form.name?.charAt(0) || 'H' }}</span>
          </div>
          <div>
            <label class="cursor-pointer inline-flex items-center gap-2 px-4 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
              Upload Logo
              <input type="file" class="hidden" accept=".png,.jpg,.jpeg,.webp" @change="handleLogoUpload" />
            </label>
            <p class="text-xs text-gray-400 mt-1">PNG, JPG, WEBP. Max 2MB.</p>
            <p v-if="logoUploading" class="text-xs text-indigo-600 mt-1">Uploading...</p>
            <p v-if="logoError" class="text-xs text-red-600 mt-1">{{ logoError }}</p>
          </div>
        </div>
      </div>

      <!-- Form -->
      <form @submit.prevent="saveProfile" class="p-6 space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Business Name</label>
          <input v-model="form.name" type="text" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" placeholder="Your business name" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Phone</label>
          <input v-model="form.phone" type="text" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" placeholder="08xxxxxxxxxx" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Address</label>
          <textarea v-model="form.address" rows="2" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none resize-none" placeholder="Business address" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description</label>
          <textarea v-model="form.description" rows="3" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none resize-none" placeholder="Brief description of your business" />
        </div>

        <div class="flex items-center gap-3">
          <button type="submit" :disabled="saving" class="px-5 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50">
            {{ saving ? 'Saving...' : 'Save Changes' }}
          </button>
          <span v-if="saveSuccess" class="text-sm text-green-600">Saved successfully!</span>
          <span v-if="saveError" class="text-sm text-red-600">{{ saveError }}</span>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTenantProfile, updateTenantProfile, uploadLogo } from '../api/tenant'

const loading = ref(true)
const saving = ref(false)
const saveSuccess = ref(false)
const saveError = ref('')
const logoUploading = ref(false)
const logoError = ref('')
const currentLogo = ref('')

const form = ref({
  name: '',
  phone: '',
  address: '',
  description: '',
})

const baseURL = typeof window !== 'undefined' && window.location.port === '8080' ? '' : 'http://localhost:8080'
const logoSrc = computed(() => currentLogo.value ? `${baseURL}${currentLogo.value}` : '')

onMounted(async () => {
  try {
    const data = await getTenantProfile()
    form.value.name = data.name || ''
    form.value.phone = data.phone || ''
    form.value.address = data.address || ''
    form.value.description = data.description || ''
    currentLogo.value = data.logo_url || ''
  } catch {
    saveError.value = 'Failed to load profile'
  } finally {
    loading.value = false
  }
})

async function handleLogoUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  logoError.value = ''
  logoUploading.value = true
  try {
    const data = await uploadLogo(file)
    currentLogo.value = data.logo_url
  } catch (err) {
    logoError.value = err.response?.data?.error || 'Upload failed'
  } finally {
    logoUploading.value = false
  }
}

async function saveProfile() {
  saving.value = true
  saveSuccess.value = false
  saveError.value = ''
  try {
    const data = await updateTenantProfile(form.value)
    form.value.name = data.name || ''
    form.value.phone = data.phone || ''
    form.value.address = data.address || ''
    form.value.description = data.description || ''
    saveSuccess.value = true
    setTimeout(() => (saveSuccess.value = false), 3000)
  } catch (err) {
    saveError.value = err.response?.data?.error || 'Failed to save'
  } finally {
    saving.value = false
  }
}
</script>
