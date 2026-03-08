<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-200 mb-4">Shifts</h1>
    <p v-if="loading" class="text-gray-600 dark:text-gray-400">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>
    <div v-else>
      <div class="sm:hidden space-y-3">
        <div v-for="s in shifts" :key="s.id" class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="flex items-start justify-between gap-2 mb-2">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDateTime(s.opened_at) }}</span>
            <span :class="s.status === 'open' ? 'text-green-600' : 'text-gray-600 dark:text-gray-400'" class="font-medium text-sm">{{ s.status }}</span>
          </div>
          <p class="text-sm text-gray-800 dark:text-gray-200 font-medium">User: {{ s.user_id }}</p>
          <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-sm text-gray-600 dark:text-gray-400 mt-1">
            <span>Buka: {{ formatPrice(s.opening_cash) }}</span>
            <span>Tutup: {{ s.closing_cash != null ? formatPrice(s.closing_cash) : '—' }}</span>
          </div>
          <p v-if="s.closed_at" class="text-xs text-gray-500 dark:text-gray-400 mt-1">Closed: {{ formatDateTime(s.closed_at) }}</p>
        </div>
        <div v-if="!shifts?.length" class="p-4 rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900">
          <p class="text-sm text-gray-500 dark:text-gray-400 text-center mb-2">Belum ada data shift yang tampil.</p>
          <p class="text-xs max-w-md mx-auto text-gray-400 dark:text-gray-500 text-center">Pastikan Anda login sebagai <strong>owner</strong> dari tenant yang sama dengan kasir yang membuka shift.</p>
        </div>
      </div>
      <div class="hidden sm:block bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Opened At</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">User ID</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Opening Cash</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Closing Cash</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Closed At</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="s in shifts" :key="s.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ formatDateTime(s.opened_at) }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ s.user_id }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ formatPrice(s.opening_cash) }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 text-right">{{ s.closing_cash != null ? formatPrice(s.closing_cash) : '—' }}</td>
              <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ s.closed_at ? formatDateTime(s.closed_at) : '—' }}</td>
              <td class="px-4 py-2">
                <span :class="s.status === 'open' ? 'text-green-600' : 'text-gray-600 dark:text-gray-400'" class="font-medium">{{ s.status }}</span>
              </td>
            </tr>
            <tr v-if="!shifts?.length">
              <td colspan="6" class="px-4 py-6 text-sm text-gray-500 dark:text-gray-400 text-center">
                <p class="mb-2">Belum ada data shift yang tampil.</p>
                <p class="text-xs max-w-md mx-auto">Pastikan Anda login sebagai <strong>owner</strong> dari tenant yang sama dengan kasir yang membuka shift. Jika kasir sudah buka/tutup shift tapi daftar kosong, coba logout lalu login lagi dengan akun owner yang dipakai saat membuat user kasir.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <p class="mt-4 text-sm text-gray-500 dark:text-gray-400">For cash reconciliation report (Expected vs Actual), see Reports → Shifts tab.</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { listShifts } from '../api/shifts'
import { formatDateTime, formatPrice } from '../utils'

const auth = useAuthStore()
const shifts = ref([])
const loading = ref(true)
const error = ref(null)

async function load() {
  if (!auth.token) {
    error.value = 'Silakan login sebagai owner untuk melihat daftar shift.'
    loading.value = false
    return
  }
  loading.value = true
  error.value = null
  try {
    const data = await listShifts({ limit: 50 })
    const arr = Array.isArray(data?.shifts) ? data.shifts : []
    shifts.value = arr
    if (import.meta.env.DEV && (arr.length === 0 && data != null)) {
      console.debug('[Shifts] API response:', data, '| parsed shifts count:', arr.length)
    }
  } catch (e) {
    const msg = e.response?.data?.error ?? e.message ?? 'Gagal memuat data shift.'
    error.value = msg === 'missing authorization header' ? 'Sesi habis atau belum login. Silakan login lagi.' : msg
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
