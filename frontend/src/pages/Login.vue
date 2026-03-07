<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-6 rounded-lg shadow-md w-full max-w-sm">
      <h1 class="text-xl font-semibold text-gray-800 mb-4">Login</h1>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            autocomplete="email"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500 focus:border-transparent"
            placeholder="user@example.com"
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="current-password"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500 focus:border-transparent"
            placeholder="••••••••"
          />
        </div>

        <p v-if="errorMessage" class="text-sm text-red-600">
          {{ errorMessage }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2 px-4 bg-slate-700 text-white rounded-md hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Signing in...' : 'Sign in' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { login as apiLogin, getProfile } from '../api/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

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
    auth.login({ token })
    const profile = await getProfile()
    auth.login({
      token,
      user: { id: profile.user_id, email: profile.email, role: profile.role },
      tenant_id: profile.tenant_id,
      role: profile.role,
    })
    router.push('/dashboard')
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
