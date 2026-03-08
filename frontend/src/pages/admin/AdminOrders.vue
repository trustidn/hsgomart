<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Subscription Orders</h1>
      <div class="flex gap-2">
        <button v-for="s in statuses" :key="s.value" @click="filterStatus = s.value"
          class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors"
          :class="filterStatus === s.value ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700'">
          {{ s.label }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400 dark:text-gray-500">Loading...</div>

    <div v-else class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 overflow-hidden">
      <div v-if="!filteredOrders.length" class="p-8 text-center text-gray-400 dark:text-gray-500 text-sm">No orders found.</div>
      <table v-else class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-xs uppercase">
          <tr>
            <th class="px-4 py-3 text-left">Invoice</th>
            <th class="px-4 py-3 text-left">Tenant</th>
            <th class="px-4 py-3 text-left">Plan</th>
            <th class="px-4 py-3 text-right">Amount</th>
            <th class="px-4 py-3 text-center">Status</th>
            <th class="px-4 py-3 text-left">Date</th>
            <th class="px-4 py-3 text-center">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
          <tr v-for="o in filteredOrders" :key="o.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
            <td class="px-4 py-3 font-mono text-xs text-gray-800 dark:text-gray-200">{{ o.invoice_number }}</td>
            <td class="px-4 py-3 text-gray-800 dark:text-gray-200">{{ o.tenant_name }}</td>
            <td class="px-4 py-3 text-gray-800 dark:text-gray-200">{{ o.plan_name }}</td>
            <td class="px-4 py-3 text-right text-gray-800 dark:text-gray-200">Rp {{ Number(o.amount).toLocaleString('id-ID') }}</td>
            <td class="px-4 py-3 text-center">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="statusClass(o.status)">
                {{ statusLabel(o.status) }}
              </span>
            </td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400">{{ formatDateTime(o.created_at) }}</td>
            <td class="px-4 py-3 text-center">
              <button @click="openDetail(o)" class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300 text-xs font-medium">View</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Detail Modal -->
    <Teleport to="body">
      <div v-if="detail" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="detail = null">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-w-lg p-6 space-y-4 max-h-[90vh] overflow-y-auto">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Order Detail</h3>

          <div class="grid grid-cols-2 gap-3 text-sm">
            <div><span class="text-gray-500 dark:text-gray-400">Invoice:</span> <strong class="font-mono">{{ detail.invoice_number }}</strong></div>
            <div><span class="text-gray-500 dark:text-gray-400">Status:</span> <span class="ml-1 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="statusClass(detail.status)">{{ statusLabel(detail.status) }}</span></div>
            <div><span class="text-gray-500 dark:text-gray-400">Tenant:</span> <span class="text-gray-800 dark:text-gray-200">{{ detail.tenant_name }}</span></div>
            <div><span class="text-gray-500 dark:text-gray-400">Plan:</span> <span class="text-gray-800 dark:text-gray-200">{{ detail.plan_name }}</span></div>
            <div><span class="text-gray-500 dark:text-gray-400">Amount:</span> Rp {{ Number(detail.amount).toLocaleString('id-ID') }}</div>
            <div><span class="text-gray-500 dark:text-gray-400">Created:</span> {{ formatDateTime(detail.created_at) }}</div>
            <div v-if="detail.paid_at"><span class="text-gray-500 dark:text-gray-400">Paid:</span> {{ formatDateTime(detail.paid_at) }}</div>
            <div v-if="detail.reviewed_at"><span class="text-gray-500 dark:text-gray-400">Reviewed:</span> {{ formatDateTime(detail.reviewed_at) }}</div>
          </div>

          <div v-if="detail.notes" class="bg-gray-50 dark:bg-gray-800 rounded-lg p-3 text-xs text-gray-700 dark:text-gray-300">
            <p class="font-medium">Notes:</p>
            <p>{{ detail.notes }}</p>
          </div>

          <div v-if="detail.admin_notes" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-3 text-xs text-yellow-800 dark:text-yellow-300">
            <p class="font-medium">Admin Notes:</p>
            <p>{{ detail.admin_notes }}</p>
          </div>

          <!-- Payment Proof -->
          <div v-if="detail.payment_proof_url">
            <p class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Payment Proof</p>
            <img :src="proofSrc(detail.payment_proof_url)" alt="Payment proof" class="w-full max-h-64 object-contain rounded-lg border border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-800" />
          </div>

          <!-- Actions for pending_review -->
          <div v-if="detail.status === 'pending_review'" class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Admin Notes (optional for approve, required for reject)</label>
              <textarea v-model="adminNotes" rows="2" class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none resize-none bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" placeholder="Notes..." />
            </div>
            <div class="flex gap-3">
              <button @click="handleReject" :disabled="processing || !adminNotes.trim()" class="flex-1 py-2 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700 transition-colors disabled:opacity-50">
                {{ processing ? 'Processing...' : 'Reject' }}
              </button>
              <button @click="handleApprove" :disabled="processing" class="flex-1 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 transition-colors disabled:opacity-50">
                {{ processing ? 'Processing...' : 'Approve' }}
              </button>
            </div>
            <p v-if="actionError" class="text-xs text-red-600">{{ actionError }}</p>
          </div>

          <button @click="detail = null" class="w-full py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">Close</button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import client, { baseURL } from '../../api/client'
import { formatDateTime } from '../../utils'

const loading = ref(true)
const orders = ref([])
const detail = ref(null)
const adminNotes = ref('')
const processing = ref(false)
const actionError = ref('')
const filterStatus = ref('')

const statuses = [
  { value: '', label: 'All' },
  { value: 'pending_review', label: 'Pending Review' },
  { value: 'pending_payment', label: 'Pending Payment' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
]

const filteredOrders = computed(() => {
  if (!filterStatus.value) return orders.value
  return orders.value.filter(o => o.status === filterStatus.value)
})

function proofSrc(url) {
  return url ? `${baseURL}${url}` : ''
}

async function load() {
  try {
    const { data } = await client.get('/admin/orders')
    orders.value = Array.isArray(data) ? data : []
  } catch { /* ignore */ } finally {
    loading.value = false
  }
}

function openDetail(o) {
  detail.value = o
  adminNotes.value = ''
  actionError.value = ''
}

async function handleApprove() {
  processing.value = true
  actionError.value = ''
  try {
    await client.put(`/admin/orders/${detail.value.id}/approve`, { admin_notes: adminNotes.value })
    detail.value = null
    loading.value = true
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.error ?? 'Failed to approve'
  } finally {
    processing.value = false
  }
}

async function handleReject() {
  processing.value = true
  actionError.value = ''
  try {
    await client.put(`/admin/orders/${detail.value.id}/reject`, { admin_notes: adminNotes.value })
    detail.value = null
    loading.value = true
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.error ?? 'Failed to reject'
  } finally {
    processing.value = false
  }
}

function statusClass(status) {
  const map = {
    pending_payment: 'bg-yellow-100 text-yellow-700',
    pending_review: 'bg-blue-100 text-blue-700',
    approved: 'bg-green-100 text-green-700',
    rejected: 'bg-red-100 text-red-700',
  }
  return map[status] || 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300'
}

function statusLabel(status) {
  const map = {
    pending_payment: 'Pending Payment',
    pending_review: 'Pending Review',
    approved: 'Approved',
    rejected: 'Rejected',
  }
  return map[status] || status
}

onMounted(load)
</script>
