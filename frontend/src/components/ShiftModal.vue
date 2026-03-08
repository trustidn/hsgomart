<template>
  <div class="fixed inset-0 z-20 flex items-center justify-center bg-black/50 p-4" @click.self="$emit('close')">
    <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl w-full max-w-sm p-6 border border-gray-200 dark:border-gray-800">
      <h2 class="text-lg font-semibold text-gray-800 dark:text-white mb-4">{{ mode === 'open' ? 'Open Shift' : 'Close Shift' }}</h2>

      <form v-if="mode === 'open'" @submit.prevent="submitOpen">
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Opening cash (IDR)</label>
          <input
            v-model.number="openForm.opening_cash"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
            placeholder="e.g. 200000"
          />
        </div>
        <p v-if="openError" class="text-sm text-red-600 mb-2">{{ openError }}</p>
        <div class="flex gap-2 justify-end">
          <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="$emit('close')">Cancel</button>
          <button type="submit" :disabled="saving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">Open Shift</button>
        </div>
      </form>

      <template v-else-if="mode === 'close'">
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Closing cash (IDR)</label>
          <input
            v-model.number="closeForm.closing_cash"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
            placeholder="e.g. 1250000"
          />
        </div>
        <p v-if="closeError" class="text-sm text-red-600 mb-2">{{ closeError }}</p>
        <div v-if="closeResult" class="mb-4 p-3 bg-gray-50 dark:bg-gray-800 rounded text-sm">
          <div class="flex justify-between"><span>Expected</span><span>{{ formatPrice(closeResult.expected_cash) }}</span></div>
          <div class="flex justify-between"><span>Actual</span><span>{{ formatPrice(closeResult.actual_cash) }}</span></div>
          <div class="flex justify-between font-medium" :class="closeResult.difference !== 0 ? 'text-amber-600' : ''">
            <span>Difference</span><span>{{ formatPrice(closeResult.difference) }}</span>
          </div>
        </div>
        <div class="flex gap-2 justify-end">
          <button type="button" class="px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md" @click="$emit('close')">{{ closeResult ? 'Done' : 'Cancel' }}</button>
          <button v-if="!closeResult" type="button" :disabled="saving" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50" @click="submitClose">Close Shift</button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { openShift, closeShift } from '../api/shifts'

const props = defineProps({
  mode: { type: String, required: true }, // 'open' | 'close'
})
const emit = defineEmits(['close', 'done'])

const saving = ref(false)
const openError = ref('')
const closeError = ref('')
const closeResult = ref(null)
const openForm = ref({ opening_cash: 200000 })
const closeForm = ref({ closing_cash: 0 })

watch(() => props.mode, (m) => {
  openError.value = ''
  closeError.value = ''
  closeResult.value = null
})

function formatPrice(v) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v ?? 0)
}

async function submitOpen() {
  saving.value = true
  openError.value = ''
  try {
    await openShift({ opening_cash: openForm.value.opening_cash })
    emit('done')
    emit('close')
  } catch (e) {
    const msg = e.response?.data?.error ?? e.message ?? 'Failed to open shift.'
    openError.value = typeof msg === 'string' ? msg : 'Failed to open shift.'
  } finally {
    saving.value = false
  }
}

async function submitClose() {
  saving.value = true
  closeError.value = ''
  closeResult.value = null
  try {
    const res = await closeShift({ closing_cash: closeForm.value.closing_cash })
    closeResult.value = res
    emit('done')
  } catch (e) {
    closeError.value = e.response?.data?.error ?? 'Failed to close shift.'
  } finally {
    saving.value = false
  }
}
</script>
