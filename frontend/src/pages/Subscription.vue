<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <h1 class="text-xl font-semibold text-gray-900">Subscription</h1>

    <div v-if="loading" class="text-center py-12 text-gray-400">Loading...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-600 text-sm">{{ error }}</div>

    <template v-else>
      <!-- Current Plan Card -->
      <div v-if="subscription" class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between flex-wrap gap-4">
          <div>
            <p class="text-sm text-gray-500">Current Plan</p>
            <h2 class="text-lg font-semibold text-gray-900">{{ subscription.plan_name }}</h2>
            <p class="text-sm mt-1">
              Status:
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="{
                  'bg-amber-100 text-amber-700': subscription.status === 'trial',
                  'bg-green-100 text-green-700': subscription.status === 'active',
                  'bg-red-100 text-red-700': subscription.status === 'expired',
                }">{{ subscription.status }}</span>
            </p>
            <p v-if="trialDaysLeft !== null" class="text-sm text-amber-600 mt-1">Trial ends in {{ trialDaysLeft }} days</p>
          </div>
          <div class="text-right text-sm text-gray-500">
            <p>Max Users: {{ plan?.max_users }}</p>
            <p>Max Products: {{ plan?.max_products }}</p>
          </div>
        </div>
      </div>

      <!-- Plans -->
      <div>
        <h2 class="text-base font-semibold text-gray-800 mb-3">Available Plans</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div v-for="p in plans" :key="p.id" class="bg-white rounded-xl border-2 p-5 transition-all"
            :class="plan?.id === p.id ? 'border-indigo-500 shadow-sm' : 'border-gray-200 hover:border-gray-300'">
            <h3 class="text-base font-bold text-gray-900">{{ p.name }}</h3>
            <p class="text-2xl font-bold text-gray-800 mt-2">
              {{ p.price === 0 ? 'Free' : 'Rp ' + p.price.toLocaleString('id-ID') }}
              <span class="text-sm font-normal text-gray-400">/mo</span>
            </p>
            <ul class="text-sm text-gray-600 mt-3 space-y-1">
              <li>{{ p.max_users }} Users</li>
              <li>{{ p.max_products }} Products</li>
            </ul>
            <div v-if="plan?.id === p.id" class="mt-4 text-center text-sm text-indigo-600 font-medium">Current Plan</div>
            <button v-else-if="p.price > 0" @click="startOrder(p)" class="mt-4 w-full py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors text-sm font-medium">
              Upgrade
            </button>
          </div>
        </div>
      </div>

      <!-- Order History -->
      <div>
        <h2 class="text-base font-semibold text-gray-800 mb-3">Order History</h2>
        <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
          <div v-if="!orders.length" class="p-6 text-center text-gray-400 text-sm">No orders yet.</div>
          <table v-else class="w-full text-sm">
            <thead class="bg-gray-50 text-gray-500 text-xs uppercase">
              <tr>
                <th class="px-4 py-3 text-left">Invoice</th>
                <th class="px-4 py-3 text-left">Plan</th>
                <th class="px-4 py-3 text-right">Amount</th>
                <th class="px-4 py-3 text-center">Status</th>
                <th class="px-4 py-3 text-center">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="o in orders" :key="o.id">
                <td class="px-4 py-3 font-mono text-xs">{{ o.invoice_number }}</td>
                <td class="px-4 py-3">{{ planName(o.plan_id) }}</td>
                <td class="px-4 py-3 text-right">Rp {{ Number(o.amount).toLocaleString('id-ID') }}</td>
                <td class="px-4 py-3 text-center">
                  <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="statusClass(o.status)">{{ statusLabel(o.status) }}</span>
                </td>
                <td class="px-4 py-3 text-center">
                  <button v-if="o.status === 'pending_payment' || o.status === 'rejected'" @click="openUpload(o)"
                    class="text-indigo-600 hover:text-indigo-800 text-xs font-medium">Upload Proof</button>
                  <span v-else-if="o.status === 'pending_review'" class="text-xs text-gray-400">Waiting review</span>
                  <span v-else-if="o.status === 'approved'" class="text-xs text-green-600">Approved</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Order Modal -->
      <Teleport to="body">
        <div v-if="showOrderModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showOrderModal = false">
          <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6 space-y-4">
            <h3 class="text-lg font-semibold text-gray-900">Order Subscription</h3>
            <div class="bg-gray-50 rounded-lg p-4 text-sm space-y-1">
              <p><span class="text-gray-500">Plan:</span> <strong>{{ selectedPlan?.name }}</strong></p>
              <p><span class="text-gray-500">Price:</span> <strong>Rp {{ Number(selectedPlan?.price).toLocaleString('id-ID') }}</strong> / month</p>
            </div>
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 text-xs text-blue-800">
              <p class="font-medium mb-1">Transfer to:</p>
              <p>Bank BCA — 1234567890</p>
              <p>a/n HSMart Payment</p>
            </div>
            <p class="text-xs text-gray-500">After creating the order, please transfer and upload the payment proof.</p>
            <div class="flex gap-3">
              <button @click="showOrderModal = false" class="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
              <button @click="confirmOrder" :disabled="creatingOrder" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
                {{ creatingOrder ? 'Creating...' : 'Create Order' }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Upload Proof Modal -->
      <Teleport to="body">
        <div v-if="showUploadModal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="showUploadModal = false">
          <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6 space-y-4">
            <h3 class="text-lg font-semibold text-gray-900">Upload Payment Proof</h3>
            <p class="text-sm text-gray-500">Invoice: <strong>{{ uploadOrder?.invoice_number }}</strong></p>
            <div v-if="uploadOrder?.status === 'rejected' && uploadOrder?.admin_notes" class="bg-red-50 border border-red-200 rounded-lg p-3 text-xs text-red-700">
              <p class="font-medium">Rejected:</p>
              <p>{{ uploadOrder.admin_notes }}</p>
            </div>
            <div>
              <label class="cursor-pointer inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                Choose File
                <input type="file" class="hidden" accept=".png,.jpg,.jpeg,.webp,.pdf" @change="handleProofFile" />
              </label>
              <span v-if="proofFile" class="text-xs text-gray-500 ml-2">{{ proofFile.name }}</span>
            </div>
            <div class="flex gap-3">
              <button @click="showUploadModal = false" class="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
              <button @click="submitProof" :disabled="!proofFile || uploading" class="flex-1 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50">
                {{ uploading ? 'Uploading...' : 'Upload' }}
              </button>
            </div>
            <p v-if="uploadError" class="text-xs text-red-600">{{ uploadError }}</p>
          </div>
        </div>
      </Teleport>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSubscription, listPlans, createOrder, uploadPaymentProof, getOrders } from '../api/subscription'

const loading = ref(true)
const error = ref('')
const subscription = ref(null)
const plan = ref(null)
const plans = ref([])
const orders = ref([])
const trialDaysLeft = ref(null)

const showOrderModal = ref(false)
const selectedPlan = ref(null)
const creatingOrder = ref(false)

const showUploadModal = ref(false)
const uploadOrder = ref(null)
const proofFile = ref(null)
const uploading = ref(false)
const uploadError = ref('')

async function load() {
  try {
    const [subData, plansData, ordersData] = await Promise.all([
      getSubscription(),
      listPlans(),
      getOrders().catch(() => []),
    ])
    subscription.value = subData.subscription
    plan.value = subData.plan
    trialDaysLeft.value = subData.trial_days_left ?? null
    plans.value = Array.isArray(plansData) ? plansData : []
    orders.value = Array.isArray(ordersData) ? ordersData : []
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to load subscription'
  } finally {
    loading.value = false
  }
}

function planName(planId) {
  return plans.value.find(p => p.id === planId)?.name || `Plan #${planId}`
}

function startOrder(p) {
  selectedPlan.value = p
  showOrderModal.value = true
}

async function confirmOrder() {
  creatingOrder.value = true
  try {
    await createOrder(selectedPlan.value.id)
    showOrderModal.value = false
    await load()
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to create order'
  } finally {
    creatingOrder.value = false
  }
}

function openUpload(order) {
  uploadOrder.value = order
  proofFile.value = null
  uploadError.value = ''
  showUploadModal.value = true
}

function handleProofFile(e) {
  proofFile.value = e.target.files?.[0] || null
}

async function submitProof() {
  if (!proofFile.value || !uploadOrder.value) return
  uploading.value = true
  uploadError.value = ''
  try {
    await uploadPaymentProof(uploadOrder.value.id, proofFile.value)
    showUploadModal.value = false
    await load()
  } catch (e) {
    uploadError.value = e.response?.data?.error ?? 'Upload failed'
  } finally {
    uploading.value = false
  }
}

function statusClass(status) {
  const map = {
    pending_payment: 'bg-yellow-100 text-yellow-700',
    pending_review: 'bg-blue-100 text-blue-700',
    approved: 'bg-green-100 text-green-700',
    rejected: 'bg-red-100 text-red-700',
  }
  return map[status] || 'bg-gray-100 text-gray-700'
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
