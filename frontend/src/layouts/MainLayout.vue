<template>
  <div class="min-h-screen flex bg-gray-50 dark:bg-gray-950 transition-colors">
    <!-- Mobile overlay -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 bg-black/40 z-30 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-40 flex flex-col bg-slate-900 dark:bg-gray-900 text-white border-r border-slate-800 dark:border-gray-800 transition-all duration-200 ease-in-out',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
        'lg:translate-x-0 lg:static lg:z-auto',
        collapsed ? 'lg:w-16' : 'lg:w-60',
        'w-60',
      ]"
    >
      <!-- Logo + Collapse toggle -->
      <div class="h-14 flex items-center px-3 border-b border-slate-800 dark:border-gray-800 shrink-0">
        <router-link to="/" class="flex items-center gap-2 flex-1 min-w-0 hover:opacity-80 transition-opacity">
          <img v-if="tenantStore.logoUrl()" :src="logoSrc" class="w-8 h-8 rounded-lg object-cover shrink-0" alt="" />
          <div v-else class="w-8 h-8 rounded-lg bg-indigo-500 flex items-center justify-center text-sm font-bold shrink-0">{{ tenantStore.storeName().charAt(0) }}</div>
          <span v-show="!collapsed" class="font-semibold text-base tracking-tight truncate">{{ tenantStore.storeName() }}</span>
        </router-link>
        <button
          class="hidden lg:flex items-center justify-center w-7 h-7 rounded-md text-slate-400 hover:text-white hover:bg-slate-700 transition-colors shrink-0"
          @click="collapsed = !collapsed"
          :title="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
        >
          <svg :class="['w-4 h-4 transition-transform', collapsed ? 'rotate-180' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
        </button>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 overflow-y-auto py-3 px-2 space-y-0.5">
        <template v-for="group in navGroups" :key="group.label">
          <div v-if="!collapsed" class="px-2 pt-4 pb-1 text-[10px] uppercase tracking-widest text-slate-500 font-semibold">{{ group.label }}</div>
          <router-link
            v-for="item in group.items"
            :key="item.path"
            :to="item.path"
            active-class="!bg-slate-700/80 !text-white"
            :class="[
              'flex items-center gap-3 px-3 py-2 rounded-lg text-slate-300 hover:bg-slate-800 hover:text-white transition-colors text-sm',
              collapsed ? 'justify-center' : '',
            ]"
            :title="collapsed ? item.label : undefined"
            @click="sidebarOpen = false"
          >
            <component :is="item.icon" class="w-[18px] h-[18px] shrink-0" />
            <span v-show="!collapsed" class="truncate">{{ item.label }}</span>
          </router-link>
        </template>
      </nav>
    </aside>

    <!-- Main content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Trial banner -->
      <div v-if="trialDaysLeft !== null && trialDaysLeft <= 7" class="bg-amber-500 text-white text-center text-sm py-2 px-4 font-medium">
        Trial Anda habis dalam {{ trialDaysLeft }} hari.
        <router-link to="/subscription" class="underline ml-1">Upgrade sekarang</router-link>
      </div>

      <!-- Top bar -->
      <header class="h-14 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 flex items-center justify-between px-4 shrink-0 sticky top-0 z-20 transition-colors">
        <div class="flex items-center gap-3">
          <button class="lg:hidden text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 -ml-1 p-1" @click="sidebarOpen = true">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>
          </button>
          <h2 class="text-sm font-medium text-gray-700 dark:text-gray-200 truncate">{{ pageTitle }}</h2>
        </div>
        <div class="flex items-center gap-1.5">
          <!-- Dark mode toggle -->
          <button
            @click="themeStore.toggle()"
            class="p-1.5 rounded-lg text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            :title="themeStore.dark ? 'Light mode' : 'Dark mode'"
          >
            <!-- Sun icon (shown in dark mode) -->
            <svg v-if="themeStore.dark" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="5" stroke-width="2"/><path stroke-linecap="round" stroke-width="2" d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
            <!-- Moon icon (shown in light mode) -->
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/></svg>
          </button>

          <!-- WhatsApp Contact -->
          <div class="relative">
            <button @click="showContactForm = !showContactForm" class="p-1.5 rounded-lg text-gray-400 dark:text-gray-500 hover:text-green-600 dark:hover:text-green-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors" title="Hubungi Admin">
              <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/></svg>
            </button>
            <div v-if="showContactForm" class="fixed inset-0 z-40" @click="showContactForm = false" />
            <Transition
              enter-active-class="transition ease-out duration-150"
              enter-from-class="opacity-0 -translate-y-1"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition ease-in duration-100"
              leave-from-class="opacity-100 translate-y-0"
              leave-to-class="opacity-0 -translate-y-1"
            >
              <div v-if="showContactForm" class="fixed left-4 right-4 top-14 sm:absolute sm:left-auto sm:right-0 sm:top-full sm:w-80 sm:mt-2 mt-2 z-50 rounded-xl bg-white dark:bg-gray-900 shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden" @click.stop>
                  <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 font-medium text-sm text-gray-800 dark:text-gray-200 flex items-center justify-between">
                    <span>Hubungi Admin</span>
                    <button type="button" class="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-500" @click="showContactForm = false" aria-label="Tutup">&times;</button>
                  </div>
                <div v-if="saasStore.whatsappNumber" class="p-4 space-y-3">
                  <textarea v-model="contactMessage" rows="3" placeholder="Tulis pesan atau pertanyaan Anda..."
                    class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-green-500 outline-none resize-none" />
                  <button @click="sendWhatsApp" :disabled="!contactMessage.trim()"
                    class="w-full py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 flex items-center justify-center gap-2">
                    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/></svg>
                    Kirim via WhatsApp
                  </button>
                </div>
                <div v-else class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">
                  <p>Nomor WhatsApp admin belum dikonfigurasi.</p>
                </div>
              </div>
            </Transition>
          </div>

          <!-- Notifications -->
          <div class="relative">
            <button @click="showNotif = !showNotif" class="relative p-1.5 rounded-lg text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" /></svg>
              <span v-if="notifCount > 0" class="absolute top-0 right-0 bg-red-500 text-white text-[10px] min-w-[16px] h-4 rounded-full flex items-center justify-center px-1 font-medium">{{ notifCount > 9 ? '9+' : notifCount }}</span>
            </button>
            <div v-if="showNotif" class="fixed inset-0 z-40" @click="showNotif = false" />
            <Transition
              enter-active-class="transition ease-out duration-150"
              enter-from-class="opacity-0 -translate-y-1"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition ease-in duration-100"
              leave-from-class="opacity-100 translate-y-0"
              leave-to-class="opacity-0 -translate-y-1"
            >
              <div v-if="showNotif" class="fixed left-4 right-4 top-14 sm:absolute sm:left-auto sm:right-0 sm:top-full sm:w-80 sm:mt-2 mt-2 z-50 bg-white dark:bg-gray-900 rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden" @click.stop>
                  <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 font-medium text-sm text-gray-800 dark:text-gray-200 flex items-center justify-between">
                    <span>Notifications</span>
                    <button type="button" class="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-500" @click="showNotif = false" aria-label="Tutup">&times;</button>
                  </div>
                <div class="max-h-64 overflow-auto divide-y divide-gray-100 dark:divide-gray-800">
                  <div v-if="!notifications.length" class="px-4 py-6 text-sm text-gray-400 text-center">No alerts</div>
                  <button
                    v-for="n in notifications"
                    :key="n.id"
                    class="w-full text-left px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors flex items-start gap-3"
                    @click="goToNotif(n)"
                  >
                    <span :class="[n.type === 'low_stock' ? 'bg-amber-100 text-amber-600 dark:bg-amber-900/40 dark:text-amber-400' : 'bg-red-100 text-red-600 dark:bg-red-900/40 dark:text-red-400', 'w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 mt-0.5']">{{ n.type === 'low_stock' ? '!' : '⏰' }}</span>
                    <div>
                      <div class="text-sm font-medium text-gray-800 dark:text-gray-200">{{ n.type === 'low_stock' ? 'Low Stock' : 'Expiring Soon' }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ n.label }}</div>
                    </div>
                  </button>
                </div>
              </div>
            </Transition>
          </div>

          <!-- User menu -->
          <div class="flex items-center gap-2 pl-2 border-l border-gray-200 dark:border-gray-700">
            <div class="w-7 h-7 rounded-full bg-slate-200 dark:bg-gray-700 flex items-center justify-center text-xs font-medium text-slate-600 dark:text-gray-300 uppercase">{{ userInitials }}</div>
            <span class="text-sm text-gray-600 dark:text-gray-400 hidden sm:block max-w-[140px] truncate">{{ userEmail }}</span>
            <button
              type="button"
              class="text-xs text-gray-400 dark:text-gray-500 hover:text-red-600 dark:hover:text-red-400 transition-colors p-1"
              title="Logout"
              @click="handleLogout"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
            </button>
          </div>
        </div>
      </header>

      <main class="flex-1 p-4 lg:p-6 overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTenantStore } from '../stores/tenant'
import { useSaasStore } from '../stores/saas'
import { useThemeStore } from '../stores/theme'
import client, { baseURL } from '../api/client'

const tenantStore = useTenantStore()
const saasStore = useSaasStore()
const themeStore = useThemeStore()

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const sidebarOpen = ref(false)
const collapsed = ref(false)
const showNotif = ref(false)
const showContactForm = ref(false)
const contactMessage = ref('')
const notifications = ref([])
const trialDaysLeft = ref(null)

function sendWhatsApp() {
  const msg = contactMessage.value.trim()
  if (!msg) return
  const raw = saasStore.whatsappNumber.replace(/\D/g, '')
  const phone = raw.startsWith('62') ? raw : '62' + raw.replace(/^0+/, '')
  window.open(`https://wa.me/${phone}?text=${encodeURIComponent(msg)}`, '_blank')
  contactMessage.value = ''
  showContactForm.value = false
}

const IconDashboard = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z' })]) }
const IconBox = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' })]) }
const IconTag = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z' })]) }
const IconArchive = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4' })]) }
const IconTruck = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1M5 17a2 2 0 104 0m-4 0a2 2 0 114 0m6 0a2 2 0 104 0m-4 0a2 2 0 114 0' })]) }
const IconHistory = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })]) }
const IconClipboard = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4' })]) }
const IconChart = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' })]) }
const IconClock = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })]) }
const IconUsers = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' })]) }
const IconCredit = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z' })]) }
const IconCart = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z' })]) }
const IconSettings = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z' })]) }
const IconBook = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253' })]) }

const menu = [
  { group: 'Main', items: [
    { path: '/dashboard', label: 'Dashboard', icon: IconDashboard, roles: ['owner', 'cashier'] },
    { path: '/pos', label: 'POS', icon: IconCart, roles: ['owner', 'cashier'] },
  ]},
  { group: 'Catalog', items: [
    { path: '/products', label: 'Products', icon: IconBox, roles: ['owner'] },
    { path: '/categories', label: 'Categories', icon: IconTag, roles: ['owner'] },
  ]},
  { group: 'Inventory', items: [
    { path: '/inventory', label: 'Inventory', icon: IconArchive, roles: ['owner'] },
    { path: '/purchases', label: 'Purchases', icon: IconTruck, roles: ['owner'] },
    { path: '/inventory-history', label: 'History', icon: IconHistory, roles: ['owner'] },
    { path: '/stock-opname', label: 'Stock Opname', icon: IconClipboard, roles: ['owner'] },
  ]},
  { group: 'Business', items: [
    { path: '/reports', label: 'Reports', icon: IconChart, roles: ['owner'] },
    { path: '/shifts', label: 'Shifts', icon: IconClock, roles: ['owner'] },
    { path: '/users', label: 'Users', icon: IconUsers, roles: ['owner'] },
    { path: '/subscription', label: 'Subscription', icon: IconCredit, roles: ['owner'] },
    { path: '/settings', label: 'Settings', icon: IconSettings, roles: ['owner'] },
  ]},
  { group: 'Help', items: [
    { path: '/documentation', label: 'Dokumentasi', icon: IconBook, roles: ['owner', 'cashier'] },
  ]},
]

const navGroups = computed(() => {
  const role = auth.role || 'cashier'
  return menu
    .map(g => ({ label: g.group, items: g.items.filter(i => i.roles.includes(role)) }))
    .filter(g => g.items.length > 0)
})

const pageTitle = computed(() => {
  const name = route.meta?.title || route.name
  if (typeof name === 'string') return name.charAt(0).toUpperCase() + name.slice(1)
  return 'Dashboard'
})

const userEmail = computed(() => auth.user?.email ?? 'User')
const userInitials = computed(() => {
  const email = auth.user?.email || 'U'
  return email.substring(0, 2).toUpperCase()
})

const notifCount = computed(() => notifications.value.length)
const logoSrc = computed(() => tenantStore.logoUrl() ? `${baseURL}${tenantStore.logoUrl()}` : '')

function handleLogout() {
  auth.logout()
  router.push('/login')
}

function goToNotif(n) {
  showNotif.value = false
  router.push('/inventory')
}

onMounted(async () => {
  await saasStore.load(true)
  tenantStore.load()

  try {
    const [lowRes, expRes] = await Promise.all([
      client.get('/api/inventory/low-stock').catch(() => ({ data: [] })),
      client.get('/api/inventory/expiring', { params: { days: 30 } }).catch(() => ({ data: [] })),
    ])
    const low = Array.isArray(lowRes.data) ? lowRes.data : []
    const exp = Array.isArray(expRes.data) ? expRes.data : []
    notifications.value = [
      ...low.map((p, i) => ({ id: `low-${i}`, type: 'low_stock', label: `${p.product_name} (${p.stock})` })),
      ...exp.map((p, i) => ({ id: `exp-${i}`, type: 'expiring', label: `${p.product_name} exp ${p.expired_at}` })),
    ]
  } catch { /* ignore */ }

  if (auth.role === 'owner') {
    try {
      const { data } = await client.get('/api/subscription')
      if (data.trial_days_left !== undefined && data.trial_days_left !== null) {
        trialDaysLeft.value = data.trial_days_left
      }
    } catch { /* ignore */ }
  }
})
</script>
