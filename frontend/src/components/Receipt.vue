<template>
  <div id="receipt-print" class="receipt bg-white text-black p-4 max-w-sm mx-auto font-mono text-sm">
    <div class="text-center border-b border-dashed border-gray-400 pb-2 mb-2">
      <div class="font-bold text-base">{{ storeName }}</div>
      <div class="text-xs text-gray-600">{{ dateFormatted }}</div>
      <div v-if="transactionId" class="text-xs text-gray-500 mt-1">Transaction: #{{ transactionId }}</div>
      <div v-if="cashier" class="text-sm mt-1">Cashier: {{ cashier }}</div>
    </div>
    <table class="w-full text-left border-b border-dashed border-gray-400 mb-2">
      <thead>
        <tr class="border-b border-gray-300">
          <th class="py-1 pr-2">Product</th>
          <th class="py-1 text-right w-12">Qty</th>
          <th class="py-1 text-right">Price</th>
          <th class="py-1 text-right">Subtotal</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.product_id" class="border-b border-gray-200">
          <td class="py-1 pr-2">{{ item.name }}</td>
          <td class="py-1 text-right">{{ item.quantity }}</td>
          <td class="py-1 text-right">{{ formatPrice(item.price) }}</td>
          <td class="py-1 text-right">{{ formatPrice(item.price * item.quantity) }}</td>
        </tr>
      </tbody>
    </table>
    <div class="space-y-1 border-t border-dashed border-gray-400 pt-2">
      <div class="flex justify-between">
        <span>Total</span>
        <span class="font-semibold">{{ formatPrice(total) }}</span>
      </div>
      <div class="flex justify-between">
        <span>Payment</span>
        <span>{{ formatPrice(paidAmount) }}</span>
      </div>
      <div class="flex justify-between">
        <span>Change</span>
        <span>{{ formatPrice(change) }}</span>
      </div>
    </div>
    <div class="text-center text-xs text-gray-500 mt-4 pt-2 border-t border-dashed border-gray-400">
      Thank you
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { parseDate, formatPrice } from '../utils'

const props = defineProps({
  storeName: { type: String, default: 'HSMart' },
  date: { type: [String, Date], default: () => new Date() },
  transactionId: { type: String, default: '' },
  cashier: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  total: { type: Number, default: 0 },
  paidAmount: { type: Number, default: 0 },
  change: { type: Number, default: 0 },
})

const dateFormatted = computed(() => {
  const d = parseDate(props.date)
  if (!d || isNaN(d)) return '—'
  return d.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
})
</script>

<style scoped>
.receipt {
  font-size: 12px;
}
</style>

<style>
@media print {
  body * {
    visibility: hidden;
  }
  #receipt-print,
  #receipt-print * {
    visibility: visible;
  }
  #receipt-print {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    max-width: none;
    box-shadow: none;
  }
}
</style>
