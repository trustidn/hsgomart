<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center gap-4 mb-4 shrink-0">
      <h1 class="text-2xl font-semibold text-gray-800">POS</h1>
      <div class="flex-1 max-w-xs">
        <input
          ref="barcodeInputRef"
          v-model="barcodeInput"
          type="text"
          placeholder="Scan barcode..."
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
          @keydown.enter.prevent="onBarcodeEnter"
        />
        <p class="text-xs text-gray-500 mt-0.5">Scan or type barcode (min 8 chars, auto-detected after 200ms)</p>
        <p v-if="barcodeError" class="text-sm text-red-600 mt-1">{{ barcodeError }}</p>
      </div>
    </div>
    <p class="text-xs text-gray-500 mb-3">
      <span class="font-medium text-gray-600">Shortcuts:</span>
      <kbd class="px-1.5 py-0.5 rounded bg-gray-100 border border-gray-300 font-mono text-xs">F1</kbd> Cash
      <span class="mx-1.5 text-gray-400">·</span>
      <kbd class="px-1.5 py-0.5 rounded bg-gray-100 border border-gray-300 font-mono text-xs">F2</kbd> Card
      <span class="mx-1.5 text-gray-400">·</span>
      <kbd class="px-1.5 py-0.5 rounded bg-gray-100 border border-gray-300 font-mono text-xs">Esc</kbd> Clear cart
    </p>

    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- LEFT: Product search list -->
      <div class="lg:col-span-2 flex flex-col min-h-0 bg-white rounded-lg shadow border border-gray-200">
        <div class="p-3 border-b border-gray-200 shrink-0">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search product..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </div>
        <div class="flex-1 overflow-auto p-3">
          <p v-if="productsLoading" class="text-gray-500">Loading products...</p>
          <p v-else-if="productsError" class="text-red-600">{{ productsError }}</p>
          <ul v-else class="grid grid-cols-2 sm:grid-cols-3 gap-2">
            <li
              v-for="p in filteredProducts"
              :key="productId(p)"
              class="border border-gray-200 rounded-lg p-3 hover:bg-slate-50 cursor-pointer transition-colors"
              @click="addToCartWithStockCheck(p)"
            >
              <div class="font-medium text-gray-800 truncate">{{ productName(p) }}</div>
              <div class="text-sm text-gray-500 truncate">{{ productSku(p) || '—' }}</div>
              <div class="text-sm font-semibold text-slate-600 mt-1">{{ formatPrice(productPrice(p)) }}</div>
            </li>
          </ul>
          <p v-if="!productsLoading && !productsError && !filteredProducts.length" class="text-gray-500">No products match your search.</p>
        </div>
      </div>

      <!-- RIGHT: Cart panel -->
      <div class="flex flex-col w-full lg:w-80 shrink-0 bg-white rounded-lg shadow border border-gray-200">
        <div class="p-3 border-b border-gray-200 font-semibold text-gray-800">Cart</div>
        <div class="flex-1 overflow-auto p-3 min-h-0">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 text-left text-gray-500 uppercase text-xs">
                <th class="pb-2 pr-2">Product</th>
                <th class="pb-2 pr-2 text-right">Qty</th>
                <th class="pb-2 text-right">Subtotal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in cartItems" :key="item.product_id" class="border-b border-gray-100 align-top">
                <td class="py-2 pr-2">
                  <div class="font-medium text-gray-800">{{ item.name }}</div>
                  <div class="flex flex-wrap items-center gap-1 mt-1">
                    <button
                      type="button"
                      class="w-6 h-6 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 flex items-center justify-center text-xs font-medium"
                      :disabled="item.quantity <= 1"
                      @click="changeQuantity(item.product_id, -1)"
                    >
                      −
                    </button>
                    <span class="min-w-[1.5rem] text-center">{{ item.quantity }}</span>
                    <button type="button" class="px-1.5 py-0.5 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 text-xs" @click="changeQuantity(item.product_id, 1)">+1</button>
                    <button type="button" class="px-1.5 py-0.5 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 text-xs" @click="changeQuantity(item.product_id, 5)">+5</button>
                    <button type="button" class="px-1.5 py-0.5 rounded border border-gray-300 text-gray-600 hover:bg-gray-100 text-xs" @click="changeQuantity(item.product_id, 10)">+10</button>
                    <button
                      type="button"
                      class="text-red-600 hover:underline text-xs ml-1"
                      @click="removeFromCart(item.product_id)"
                    >
                      Remove
                    </button>
                  </div>
                </td>
                <td class="py-2 pr-2 text-right">{{ item.quantity }}</td>
                <td class="py-2 text-right font-medium">{{ formatPrice(item.price * item.quantity) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-if="!cartItems.length" class="text-gray-500 py-4">Cart is empty.</p>
        </div>
        <div class="p-3 border-t border-gray-200 shrink-0">
          <div class="flex justify-between items-center text-lg font-semibold text-gray-800 mb-3">
            <span>Total</span>
            <span>{{ formatPrice(totalAmount) }}</span>
          </div>
          <div class="flex gap-2">
            <button
              type="button"
              class="flex-1 py-2.5 bg-slate-600 text-white rounded-md hover:bg-slate-700 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="!cartItems.length"
              @click="openCheckoutModal('cash')"
            >
              CASH
            </button>
            <button
              type="button"
              class="flex-1 py-2.5 bg-slate-700 text-white rounded-md hover:bg-slate-800 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="!cartItems.length"
              @click="openCheckoutModal('card')"
            >
              CARD
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Checkout modal -->
    <div
      v-if="showCheckoutModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50"
      @click.self="showCheckoutModal = false"
    >
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Checkout</h2>
        <form @submit.prevent="submitCheckout">
          <div class="mb-4">
            <label for="payment-method" class="block text-sm font-medium text-gray-700 mb-1">Payment method</label>
            <select
              id="payment-method"
              v-model="checkoutForm.payment_method"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
            >
              <option value="cash">Cash</option>
              <option value="card">Card</option>
            </select>
          </div>
          <div class="mb-4">
            <label for="paid-amount" class="block text-sm font-medium text-gray-700 mb-1">Paid amount</label>
            <input
              id="paid-amount"
              v-model.number="checkoutForm.paid_amount"
              type="number"
              min="0"
              step="0.01"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-slate-500"
              placeholder="0"
            />
            <p class="text-xs text-gray-500 mt-1">Total: {{ formatPrice(totalAmount) }}</p>
          </div>
          <p v-if="checkoutError" class="text-sm text-red-600 mb-2">{{ checkoutError }}</p>
          <div class="flex gap-2 justify-end">
            <button type="button" class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md" @click="showCheckoutModal = false">
              Cancel
            </button>
            <button type="submit" :disabled="checkoutSubmitting" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 disabled:opacity-50">
              {{ checkoutSubmitting ? 'Processing...' : 'Complete' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Receipt modal (after checkout) -->
    <div
      v-if="showReceiptModal"
      class="fixed inset-0 z-10 flex items-center justify-center bg-black/50 p-4"
      @click.self="closeReceipt"
    >
      <div class="bg-white rounded-lg shadow-xl max-w-sm w-full max-h-[90vh] overflow-auto">
        <div class="p-4">
          <Receipt
            v-if="receiptData"
            :store-name="receiptData.storeName"
            :date="receiptData.date"
            :transaction-id="receiptData.transactionId"
            :items="receiptData.items"
            :total="receiptData.total"
            :paid-amount="receiptData.paidAmount"
            :change="receiptData.change"
          />
        </div>
        <div class="p-4 border-t border-gray-200 flex gap-2 justify-end">
          <button type="button" class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50" @click="printReceipt">
            Print Receipt
          </button>
          <button type="button" class="px-4 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700" @click="closeReceipt">
            Done
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { getProducts, getProductByBarcode } from '../api/products'
import { getProductStock } from '../api/inventory'
import { checkout as checkoutApi } from '../api/pos'
import Receipt from '../components/Receipt.vue'

const SCAN_DEBOUNCE_MS = 200
const MIN_BARCODE_LENGTH = 8

const products = ref([])
const barcodeInput = ref('')
const barcodeInputRef = ref(null)
const barcodeError = ref('')
const filteredProducts = ref([])
const productsLoading = ref(true)
const productsError = ref(null)
const searchQuery = ref('')

const cartItems = ref([])

const showCheckoutModal = ref(false)
const checkoutForm = ref({ payment_method: 'cash', paid_amount: 0 })
const checkoutError = ref('')
const checkoutSubmitting = ref(false)
const showReceiptModal = ref(false)
const receiptData = ref(null)
let scanDebounceTimer = null

function productId(p) {
  return p?.id ?? p?.ID ?? ''
}
function productName(p) {
  return p?.name ?? p?.Name ?? ''
}
function productSku(p) {
  return p?.sku ?? p?.SKU ?? ''
}
function productPrice(p) {
  const n = p?.sell_price ?? p?.SellPrice ?? p?.price ?? p?.Price
  return typeof n === 'number' ? n : Number(n) || 0
}

function formatPrice(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value ?? 0)
}

function filterProducts() {
  const q = (searchQuery.value || '').trim().toLowerCase()
  if (!q) {
    filteredProducts.value = [...products.value]
    return
  }
  filteredProducts.value = products.value.filter((p) => {
    const name = (productName(p) || '').toLowerCase()
    const sku = (productSku(p) || '').toLowerCase()
    return name.includes(q) || sku.includes(q)
  })
}

watch(searchQuery, filterProducts)

const totalAmount = computed(() =>
  cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
)

function webAudioBeep(freq = 800, duration = 0.1) {
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.frequency.value = freq
    osc.type = 'sine'
    gain.gain.setValueAtTime(0.15, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + duration)
    osc.start(ctx.currentTime)
    osc.stop(ctx.currentTime + duration)
  } catch (_) {}
}

function playBeep() {
  const audio = new Audio('/beep.mp3')
  audio.play().catch(() => webAudioBeep(800))
}

function playErrorBeep() {
  const audio = new Audio('/beep-error.mp3')
  audio.play().catch(() => webAudioBeep(400, 0.2))
}

function addToCart(p) {
  const id = productId(p)
  const existing = cartItems.value.find((i) => i.product_id === id)
  if (existing) {
    existing.quantity += 1
  } else {
    cartItems.value.push({
      product_id: id,
      name: productName(p),
      price: productPrice(p),
      quantity: 1,
    })
  }
  playBeep()
}

/** Returns true if product was added, false if blocked by stock or error. */
async function addToCartWithStockCheck(p) {
  const id = productId(p)
  barcodeError.value = ''
  try {
    const res = await getProductStock(id)
    const stock = res?.stock ?? 0
    if (stock <= 0) {
      barcodeError.value = 'Insufficient stock'
      playErrorBeep()
      return false
    }
    const existing = cartItems.value.find((i) => i.product_id === id)
    const nextQty = existing ? existing.quantity + 1 : 1
    if (nextQty > stock) {
      barcodeError.value = 'Insufficient stock'
      playErrorBeep()
      return false
    }
    addToCart(p)
    return true
  } catch (err) {
    barcodeError.value = err.response?.data?.error ?? 'Failed to check stock.'
    playErrorBeep()
    return false
  }
}

async function lookupAndAddBarcode(barcode) {
  const code = (barcode || '').trim()
  if (!code || code.length < MIN_BARCODE_LENGTH) return
  barcodeError.value = ''
  try {
    const product = await getProductByBarcode(code)
    const added = await addToCartWithStockCheck(product)
    if (added) {
      barcodeInput.value = ''
      nextTick(() => barcodeInputRef.value?.focus())
    }
  } catch (err) {
    if (err.response?.status === 404) {
      barcodeError.value = 'Product not found'
      playErrorBeep()
    } else {
      barcodeError.value = err.response?.data?.error ?? 'Failed to find product.'
      playErrorBeep()
    }
  }
}


async function onBarcodeEnter() {
  const barcode = (barcodeInput.value || '').trim()
  if (barcode.length < MIN_BARCODE_LENGTH) return
  await lookupAndAddBarcode(barcode)
}

watch(barcodeInput, (val) => {
  if (scanDebounceTimer) clearTimeout(scanDebounceTimer)
  const barcode = (val || '').trim()
  if (barcode.length < MIN_BARCODE_LENGTH) return
  scanDebounceTimer = setTimeout(() => {
    scanDebounceTimer = null
    lookupAndAddBarcode(barcode)
  }, SCAN_DEBOUNCE_MS)
})

function changeQuantity(productId, delta) {
  const item = cartItems.value.find((i) => i.product_id === productId)
  if (!item) return
  item.quantity = Math.max(0, item.quantity + delta)
  if (item.quantity <= 0) {
    cartItems.value = cartItems.value.filter((i) => i.product_id !== productId)
  }
}

function removeFromCart(productId) {
  cartItems.value = cartItems.value.filter((i) => i.product_id !== productId)
}

function openCheckoutModal(method) {
  checkoutForm.value = { payment_method: method, paid_amount: totalAmount.value }
  checkoutError.value = ''
  showCheckoutModal.value = true
}

async function submitCheckout() {
  checkoutError.value = ''
  checkoutSubmitting.value = true
  try {
    const result = await checkoutApi({
      items: cartItems.value.map((i) => ({ product_id: i.product_id, quantity: i.quantity })),
      payment_method: checkoutForm.value.payment_method,
      paid_amount: checkoutForm.value.paid_amount,
    })
    showCheckoutModal.value = false
    receiptData.value = {
      storeName: 'HSMart',
      date: new Date(),
      transactionId: result?.transaction_id ?? '',
      items: cartItems.value.map((i) => ({ ...i })),
      total: result?.total ?? totalAmount.value,
      paidAmount: checkoutForm.value.paid_amount,
      change: result?.change ?? 0,
    }
    showReceiptModal.value = true
  } catch (err) {
    checkoutError.value = err.response?.data?.error ?? 'Checkout failed. Please try again.'
  } finally {
    checkoutSubmitting.value = false
  }
}

function printReceipt() {
  window.print()
}

function closeReceipt() {
  showReceiptModal.value = false
  receiptData.value = null
  cartItems.value = []
  nextTick(() => barcodeInputRef.value?.focus())
}

function confirmClearCart() {
  if (!cartItems.value.length) return
  if (window.confirm('Clear cart?')) {
    cartItems.value = []
    barcodeError.value = ''
  }
}

function handleKeyShortcuts(e) {
  if (showCheckoutModal.value || showReceiptModal.value) return
  if (e.key === 'F1') {
    e.preventDefault()
    if (cartItems.value.length) openCheckoutModal('cash')
  } else if (e.key === 'F2') {
    e.preventDefault()
    if (cartItems.value.length) openCheckoutModal('card')
  } else if (e.key === 'Escape') {
    confirmClearCart()
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleKeyShortcuts)
  productsLoading.value = true
  productsError.value = null
  try {
    const data = await getProducts()
    products.value = Array.isArray(data) ? data : []
    filterProducts()
  } catch (err) {
    productsError.value = 'Failed to load products.'
  } finally {
    productsLoading.value = false
  }
  await nextTick()
  barcodeInputRef.value?.focus()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyShortcuts)
})
</script>
