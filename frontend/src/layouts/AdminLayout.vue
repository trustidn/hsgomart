<template>
  <div class="min-h-screen flex bg-gray-100">
    <aside class="w-56 bg-gray-900 text-white flex flex-col shrink-0">
      <div class="p-4 font-semibold text-lg border-b border-gray-700">HSMart Admin</div>
      <nav class="flex-1 p-2 space-y-1">
        <router-link to="/admin/dashboard" active-class="bg-gray-700" class="block px-3 py-2 rounded-md text-gray-200 hover:bg-gray-700">Dashboard</router-link>
        <router-link to="/admin/tenants" active-class="bg-gray-700" class="block px-3 py-2 rounded-md text-gray-200 hover:bg-gray-700">Tenants</router-link>
        <router-link to="/admin/subscriptions" active-class="bg-gray-700" class="block px-3 py-2 rounded-md text-gray-200 hover:bg-gray-700">Subscriptions</router-link>
      </nav>
    </aside>
    <div class="flex-1 flex flex-col min-w-0">
      <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-4">
        <span class="text-gray-600">Super Admin</span>
        <button @click="handleLogout" class="text-sm text-red-600 hover:underline">Logout</button>
      </header>
      <main class="flex-1 p-4 overflow-auto"><RouterView /></main>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
const router = useRouter()
const auth = useAuthStore()
function handleLogout() { auth.logout(); router.push('/login') }
</script>
