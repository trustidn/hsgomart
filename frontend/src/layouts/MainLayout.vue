<template>
  <div class="min-h-screen flex bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-56 bg-slate-800 text-white flex flex-col shrink-0">
      <div class="p-4 font-semibold text-lg border-b border-slate-700">
        HSMart POS
      </div>
      <nav class="flex-1 p-2 space-y-1">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          active-class="bg-slate-600 text-white"
          class="block px-3 py-2 rounded-md text-slate-200 hover:bg-slate-700"
        >
          {{ item.label }}
        </router-link>
      </nav>
    </aside>

    <!-- Main content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Topbar -->
      <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-4 shrink-0">
        <span class="text-gray-600">HSMart SaaS</span>
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-500">{{ userEmail }}</span>
          <button
            type="button"
            class="text-sm text-red-600 hover:underline"
            @click="handleLogout"
          >
            Logout
          </button>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 p-4 overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const navItems = [
  { path: '/dashboard', label: 'Dashboard' },
  { path: '/products', label: 'Products' },
  { path: '/inventory', label: 'Inventory' },
  { path: '/inventory-history', label: 'Inventory History' },
  { path: '/categories', label: 'Categories' },
  { path: '/purchases', label: 'Purchases' },
  { path: '/pos', label: 'POS' },
  { path: '/reports', label: 'Reports' },
]

const userEmail = computed(() => auth.user?.email ?? 'User')

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
