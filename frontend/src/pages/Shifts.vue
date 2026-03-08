<template>
  <div>
    <h1 class="text-2xl font-semibold text-gray-800 mb-4">Shifts</h1>
    <p v-if="loading" class="text-gray-600">Loading...</p>
    <p v-else-if="error" class="text-red-600">{{ error }}</p>
    <div v-else class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Opened At</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">User ID</th>
            <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Opening Cash</th>
            <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Closing Cash</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Closed At</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="s in shifts" :key="s.id" class="hover:bg-gray-50">
            <td class="px-4 py-2 text-sm text-gray-800">{{ formatDateTime(s.opened_at) }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ s.user_id }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ formatPrice(s.opening_cash) }}</td>
            <td class="px-4 py-2 text-sm text-gray-600 text-right">{{ s.closing_cash != null ? formatPrice(s.closing_cash) : '—' }}</td>
            <td class="px-4 py-2 text-sm text-gray-600">{{ s.closed_at ? formatDateTime(s.closed_at) : '—' }}</td>
            <td class="px-4 py-2">
              <span :class="s.status === 'open' ? 'text-green-600' : 'text-gray-600'" class="font-medium">{{ s.status }}</span>
            </td>
          </tr>
          <tr v-if="!shifts?.length">
            <td colspan="6" class="px-4 py-6 text-sm text-gray-500 text-center">
              <p class="mb-2">Belum ada data shift yang tampil.</p>
              <p class="text-xs max-w-md mx-auto">Pastikan Anda login sebagai <strong>owner</strong> dari tenant yang sama dengan kasir yang membuka shift. Jika kasir sudah buka/tutup shift tapi daftar kosong, coba logout lalu login lagi dengan akun owner yang dipakai saat membuat user kasir.</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="mt-4 text-sm text-gray-500">For cash reconciliation report (Expected vs Actual), see Reports → Shifts tab.</p>
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
