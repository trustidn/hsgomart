<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Subscription</h1>

    <div v-if="loading" class="text-gray-500">Loading...</div>
    <div v-else-if="error" class="text-red-600">{{ error }}</div>
    <div v-else>
      <div v-if="subscription" class="bg-white rounded-lg shadow p-6 mb-6">
        <h2 class="text-lg font-semibold mb-2">Current Plan: {{ subscription.plan_name }}</h2>
        <p class="text-sm text-gray-600">Status: <span class="font-medium" :class="subscription.status === 'trial' ? 'text-amber-600' : 'text-green-600'">{{ subscription.status }}</span></p>
        <p v-if="trialDaysLeft !== null" class="text-sm text-amber-600 mt-1">Trial ends in {{ trialDaysLeft }} days</p>
        <p class="text-sm text-gray-500 mt-1">Max Users: {{ plan?.max_users }} | Max Products: {{ plan?.max_products }}</p>
      </div>

      <h2 class="text-lg font-semibold mb-3">Available Plans</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div v-for="p in plans" :key="p.id" class="bg-white rounded-lg shadow p-5 border-2" :class="plan?.id === p.id ? 'border-slate-600' : 'border-transparent'">
          <h3 class="text-lg font-bold">{{ p.name }}</h3>
          <p class="text-2xl font-bold text-slate-700 mt-2">{{ p.price === 0 ? 'Free' : 'Rp ' + p.price.toLocaleString() }}<span class="text-sm font-normal text-gray-500">/mo</span></p>
          <ul class="text-sm text-gray-600 mt-3 space-y-1">
            <li>{{ p.max_users }} Users</li>
            <li>{{ p.max_products }} Products</li>
          </ul>
          <button v-if="plan?.id !== p.id" @click="selectPlan(p.id)" :disabled="changing" class="mt-4 w-full py-2 bg-slate-700 text-white rounded hover:bg-slate-600 disabled:opacity-50 text-sm">
            {{ changing ? 'Changing...' : 'Select Plan' }}
          </button>
          <p v-else class="mt-4 text-center text-sm text-green-600 font-medium">Current Plan</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSubscription, listPlans, changePlan } from '../api/subscription'

const loading = ref(true)
const error = ref('')
const subscription = ref(null)
const plan = ref(null)
const plans = ref([])
const trialDaysLeft = ref(null)
const changing = ref(false)

async function load() {
  try {
    const [subData, plansData] = await Promise.all([getSubscription(), listPlans()])
    subscription.value = subData.subscription
    plan.value = subData.plan
    trialDaysLeft.value = subData.trial_days_left ?? null
    plans.value = Array.isArray(plansData) ? plansData : []
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to load subscription'
  } finally {
    loading.value = false
  }
}

async function selectPlan(planId) {
  changing.value = true
  try {
    await changePlan(planId)
    await load()
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Failed to change plan'
  } finally {
    changing.value = false
  }
}

onMounted(load)
</script>
