<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <div class="w-12 h-12 rounded-xl bg-indigo-500 flex items-center justify-center text-white text-xl font-bold mx-auto mb-4">H</div>
        <h1 class="text-2xl font-bold text-gray-900">Create your account</h1>
        <p class="text-sm text-gray-500 mt-1">Start your free trial of HSMart POS</p>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 mb-1.5">Business Name</label>
            <input id="name" v-model="name" type="text" required class="w-full px-3.5 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm" placeholder="My Store" />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input id="email" v-model="email" type="email" required class="w-full px-3.5 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm" placeholder="you@example.com" />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1.5">Password</label>
            <input id="password" v-model="password" type="password" required minlength="8" class="w-full px-3.5 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm" placeholder="Min 8 chars, letters + numbers" />
          </div>

          <div v-if="errorMessage" class="bg-red-50 border border-red-200 text-red-700 text-sm px-3 py-2 rounded-lg">
            {{ errorMessage }}
          </div>

          <button type="submit" :disabled="loading" class="w-full py-2.5 px-4 bg-slate-800 text-white rounded-lg hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium transition-colors">
            {{ loading ? 'Creating...' : 'Create Account' }}
          </button>
        </form>
      </div>

      <p class="text-sm text-gray-500 text-center mt-4">
        Already have an account?
        <router-link to="/login" class="text-indigo-600 hover:text-indigo-500 font-medium">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { register, getProfile } from '../api/auth'

const router = useRouter()
const auth = useAuthStore()
const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

async function handleSubmit() {
  errorMessage.value = ''
  loading.value = true
  try {
    const data = await register(name.value, email.value, password.value)
    auth.login({ token: data.token, refresh_token: data.refresh_token })
    const profile = await getProfile()
    auth.login({
      token: data.token,
      refresh_token: data.refresh_token,
      user: { id: profile.user_id, email: profile.email, role: profile.role },
      tenant_id: profile.tenant_id,
      role: profile.role,
    })
    router.push('/dashboard')
  } catch (err) {
    errorMessage.value = err.response?.data?.error ?? err.message ?? 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>
