<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Revenue Report</h1>

    <!-- Filters -->
    <div class="bg-white rounded-xl border border-gray-200 p-4 mb-6 flex flex-wrap items-end gap-4">
      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1">From</label>
        <input v-model="fromDate" type="date" class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
      </div>
      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1">To</label>
        <input v-model="toDate" type="date" class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" />
      </div>
      <button @click="load" class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors">Filter</button>
      <button @click="exportCSV" v-if="report.orders?.length" class="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors">Export CSV</button>
    </div>

    <div v-if="loading" class="text-gray-400 py-8 text-center">Loading...</div>
    <template v-else>
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white rounded-xl border border-gray-200 p-5">
          <p class="text-sm text-gray-500">Total Revenue</p>
          <p class="text-2xl font-bold text-gray-800">Rp {{ Number(report.total_revenue || 0).toLocaleString('id-ID') }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-5">
          <p class="text-sm text-gray-500">Total Orders</p>
          <p class="text-2xl font-bold text-gray-800">{{ report.total_orders || 0 }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-5">
          <p class="text-sm text-gray-500">Average per Order</p>
          <p class="text-2xl font-bold text-gray-800">Rp {{ Number(report.avg_per_order || 0).toLocaleString('id-ID') }}</p>
        </div>
      </div>

      <!-- Plan Breakdown -->
      <div v-if="report.plan_breakdown?.length" class="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <h2 class="text-base font-semibold text-gray-800 mb-3">Revenue by Plan</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div v-for="pb in report.plan_breakdown" :key="pb.plan_name" class="bg-gray-50 rounded-lg p-3">
            <p class="text-sm font-medium text-gray-700">{{ pb.plan_name }}</p>
            <p class="text-lg font-bold text-gray-800">Rp {{ Number(pb.total || 0).toLocaleString('id-ID') }}</p>
            <p class="text-xs text-gray-400">{{ pb.order_count }} order(s)</p>
          </div>
        </div>
      </div>

      <!-- Detail Table -->
      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div v-if="!report.orders?.length" class="p-8 text-center text-gray-400 text-sm">No approved orders in this period.</div>
        <table v-else class="w-full text-sm">
          <thead class="bg-gray-50 text-gray-500 text-xs uppercase">
            <tr>
              <th class="px-4 py-3 text-left">Invoice</th>
              <th class="px-4 py-3 text-left">Tenant</th>
              <th class="px-4 py-3 text-left">Plan</th>
              <th class="px-4 py-3 text-right">Amount</th>
              <th class="px-4 py-3 text-left">Approved At</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="o in report.orders" :key="o.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 font-mono text-xs">{{ o.invoice }}</td>
              <td class="px-4 py-3">{{ o.tenant_name }}</td>
              <td class="px-4 py-3">{{ o.plan_name }}</td>
              <td class="px-4 py-3 text-right font-medium">Rp {{ Number(o.amount).toLocaleString('id-ID') }}</td>
              <td class="px-4 py-3 text-gray-500">{{ formatDateTime(o.approved_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getRevenueReport } from '../../api/admin'
import { formatDateTime } from '../../utils'

const loading = ref(true)
const fromDate = ref('')
const toDate = ref('')
const report = ref({})

async function load() {
  loading.value = true
  try {
    report.value = await getRevenueReport(fromDate.value, toDate.value)
  } catch { report.value = {} }
  finally { loading.value = false }
}

function exportCSV() {
  if (!report.value.orders?.length) return
  const headers = ['Invoice', 'Tenant', 'Plan', 'Amount', 'Approved At']
  const rows = report.value.orders.map(o => [o.invoice, o.tenant_name, o.plan_name, o.amount, o.approved_at])
  const csv = [headers, ...rows].map(r => r.map(c => `"${c}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `revenue_${fromDate.value || 'all'}_${toDate.value || 'all'}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(load)
</script>
