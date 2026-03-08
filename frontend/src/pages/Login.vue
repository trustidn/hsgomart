<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-gray-950 px-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <img v-if="saas.logoSrc" :src="saas.logoSrc" class="w-14 h-14 rounded-xl object-cover mx-auto mb-4" alt="" />
        <div v-else class="w-14 h-14 rounded-xl bg-indigo-500 flex items-center justify-center text-white text-2xl font-bold mx-auto mb-4">{{ saas.platformName.charAt(0) }}</div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ saas.platformName }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ saas.tagline || 'Sign in to your account' }}</p>
      </div>

      <div class="bg-white dark:bg-gray-900 p-6 rounded-xl shadow-sm border border-gray-200 dark:border-gray-800">
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              autocomplete="email"
              class="w-full px-3.5 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-3.5 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm"
              placeholder="Enter your password"
            />
          </div>

          <div v-if="errorMessage" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-sm px-3 py-2 rounded-lg">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 bg-slate-800 text-white rounded-lg hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium transition-colors"
          >
            {{ loading ? 'Signing in...' : 'Sign in' }}
          </button>
        </form>
      </div>

      <p class="text-sm text-gray-500 dark:text-gray-400 text-center mt-4">
        Don't have an account?
        <router-link to="/register" class="text-indigo-600 hover:text-indigo-500 font-medium">Create one</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSaasStore } from '../stores/saas'
import { login as apiLogin, getProfile } from '../api/auth'

const router = useRouter()
const auth = useAuthStore()
const saas = useSaasStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

onMounted(async () => {
  await saas.load()
  saas.applyBranding()
})

async function handleSubmit() {
  errorMessage.value = ''
  loading.value = true
  try {
    const data = await apiLogin(email.value, password.value)
    const token = data.token
    if (!token) {
      errorMessage.value = 'Invalid response from server.'
      return
    }
    auth.login({ token, refresh_token: data.refresh_token })
    const profile = await getProfile()
    auth.login({
      token,
      refresh_token: data.refresh_token,
      user: { id: profile.user_id, email: profile.email, role: profile.role },
      tenant_id: profile.tenant_id,
      role: profile.role,
    })
    router.push(profile.role === 'superadmin' ? '/admin/dashboard' : '/dashboard')
  } catch (err) {
    const msg = err.response?.data?.error ?? err.message
    errorMessage.value = (err.response?.status === 401 || (msg && msg.toLowerCase().includes('invalid')))
      ? 'Invalid email or password.'
      : (msg || 'Invalid email or password.')
  } finally {
    loading.value = false
  }
}
</script>
