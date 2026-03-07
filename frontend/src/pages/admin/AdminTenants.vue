<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Tenants</h1>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <table v-else class="w-full bg-white rounded-lg shadow text-sm">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-4 py-3 text-left text-gray-600">Name</th>
          <th class="px-4 py-3 text-left text-gray-600">Email</th>
          <th class="px-4 py-3 text-left text-gray-600">Status</th>
          <th class="px-4 py-3 text-left text-gray-600">Plan</th>
          <th class="px-4 py-3 text-right text-gray-600">Users</th>
          <th class="px-4 py-3 text-left text-gray-600">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in tenants" :key="t.id" class="border-t">
          <td class="px-4 py-3 font-medium">{{ t.name }}</td>
          <td class="px-4 py-3">{{ t.email }}</td>
          <td class="px-4 py-3"><span class="px-2 py-0.5 rounded text-xs" :class="t.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">{{ t.status }}</span></td>
          <td class="px-4 py-3">{{ t.plan_name }}</td>
          <td class="px-4 py-3 text-right">{{ t.user_count }}</td>
          <td class="px-4 py-3">
            <button v-if="t.status === 'active'" @click="toggle(t.id, 'suspended')" class="text-red-600 hover:underline text-xs">Suspend</button>
            <button v-else @click="toggle(t.id, 'active')" class="text-green-600 hover:underline text-xs">Activate</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listTenants, updateTenant } from '../../api/admin'

const loading = ref(true)
const tenants = ref([])

async function load() {
  loading.value = true
  try { tenants.value = await listTenants() }
  catch { tenants.value = [] }
  finally { loading.value = false }
}

async function toggle(id, status) {
  await updateTenant(id, status)
  await load()
}

onMounted(load)
</script>
