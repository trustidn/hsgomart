<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Subscriptions</h1>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <table v-else class="w-full bg-white rounded-lg shadow text-sm">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-4 py-3 text-left text-gray-600">Tenant</th>
          <th class="px-4 py-3 text-left text-gray-600">Plan</th>
          <th class="px-4 py-3 text-center text-gray-600">Status</th>
          <th class="px-4 py-3 text-left text-gray-600">End Date</th>
          <th class="px-4 py-3 text-center text-gray-600">Days Left</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in subs" :key="s.id" class="border-t hover:bg-gray-50">
          <td class="px-4 py-3 font-medium">{{ s.tenant_name }}</td>
          <td class="px-4 py-3">{{ s.plan_name }}</td>
          <td class="px-4 py-3 text-center">
            <span class="px-2 py-0.5 rounded text-xs font-medium"
              :class="s.status === 'active' ? 'bg-green-100 text-green-700' : s.status === 'trial' ? 'bg-amber-100 text-amber-700' : 'bg-red-100 text-red-700'">
              {{ s.status }}
            </span>
          </td>
          <td class="px-4 py-3">{{ s.end_date ?? '-' }}</td>
          <td class="px-4 py-3 text-center">
            <span v-if="s.days_remaining !== null && s.days_remaining !== undefined" class="text-xs font-medium"
              :class="s.days_remaining <= 0 ? 'text-red-600' : s.days_remaining <= 7 ? 'text-amber-600' : 'text-green-600'">
              {{ s.days_remaining }} days
            </span>
            <span v-else class="text-xs text-gray-400">-</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listSubscriptions } from '../../api/admin'

const loading = ref(true)
const subs = ref([])

onMounted(async () => {
  try { subs.value = await listSubscriptions() }
  catch { subs.value = [] }
  finally { loading.value = false }
})
</script>
