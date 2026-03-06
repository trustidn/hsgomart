<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <h1 class="text-2xl font-semibold text-gray-800">Reports</h1>
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-sm text-gray-600">Period:</span>
        <button
          v-for="opt in dateFilterOptions"
          :key="opt.value"
          type="button"
          :class="[dateFilter === opt.value ? 'bg-slate-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']"
          class="px-3 py-1.5 rounded text-sm font-medium"
          @click="setDateFilter(opt.value)"
        >
          {{ opt.label }}
        </button>
        <template v-if="dateFilter === 'custom'">
          <input v-model="customFrom" type="date" class="px-2 py-1.5 border rounded text-sm" />
          <span class="text-gray-400">to</span>
          <input v-model="customTo" type="date" class="px-2 py-1.5 border rounded text-sm" />
        </template>
      </div>
    </div>

    <div class="border-b border-gray-200">
      <nav class="flex gap-4">
        <button
          v-for="t in tabs"
          :key="t.id"
          type="button"
          :class="[activeTab === t.id ? 'border-slate-600 text-slate-600' : 'border-transparent text-gray-500 hover:text-gray-700']"
          class="py-2 px-1 border-b-2 font-medium text-sm"
          @click="activeTab = t.id"
        >
          {{ t.label }}
        </button>
      </nav>
    </div>

    <div v-if="loading" class="text-gray-500 py-8">Loading...</div>
    <div v-else-if="error" class="text-red-600 py-4">{{ error }}</div>

    <!-- Sales -->
    <template v-else-if="activeTab === 'sales'">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <div class="text-sm text-gray-500">Total Revenue</div>
          <div class="text-xl font-semibold text-gray-800">{{ formatPrice(salesSummary?.total_sales ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <div class="text-sm text-gray-500">Total Transactions</div>
          <div class="text-xl font-semibold text-gray-800">{{ salesSummary?.total_transactions ?? 0 }}</div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between">
          <span class="font-medium text-gray-800">Payment Methods</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Method', 'Transactions', 'Revenue'], paymentsTableRows, 'payment-methods')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Payment Methods', ['Method', 'Transactions', 'Revenue'], paymentsTableRows, 'payment-methods')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Method</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Transactions</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="(row, i) in paymentsReport" :key="i" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.method }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.transactions }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.revenue) }}</td>
            </tr>
            <tr v-if="!paymentsReport?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between flex-wrap gap-2">
          <span class="font-medium text-gray-800">Ringkasan per Hari</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(salesDailyHeaders, salesDailyRows, 'ringkasan-per-hari')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Ringkasan per Hari', salesDailyHeaders, salesDailyRows, 'ringkasan-per-hari')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Transactions</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="row in salesDaily" :key="row.date" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.date }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.total_transactions }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.total_sales) }}</td>
            </tr>
            <tr v-if="!salesDaily?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between flex-wrap gap-2">
          <span class="font-medium text-gray-800">Penjualan per Jam</span>
          <div class="flex items-center gap-2">
            <label class="text-sm text-gray-600">Tanggal:</label>
            <input v-model="hourlyDate" type="date" class="px-2 py-1.5 border rounded text-sm" @change="loadSalesHourly" />
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Hour</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Transactions</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="row in salesHourly" :key="row.hour" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.hour }}:00</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.transactions }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.revenue) }}</td>
            </tr>
            <tr v-if="hourlyDate && !salesHourly?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this date.</td>
            </tr>
            <tr v-if="!hourlyDate">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">Pilih tanggal untuk melihat penjualan per jam.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between flex-wrap gap-2">
          <span class="font-medium text-gray-800">Detail Transaksi</span>
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-500">Showing {{ salesTransactionsShowing.from }}–{{ salesTransactionsShowing.to }} of {{ salesTransactionsShowing.total }} transactions</span>
            <div class="flex gap-1">
              <button type="button" class="px-2 py-1 border rounded text-sm disabled:opacity-50" :disabled="salesTransactionsPage <= 1" @click="loadSalesTransactionsPage(salesTransactionsPage - 1)">Prev</button>
              <button type="button" class="px-2 py-1 border rounded text-sm disabled:opacity-50" :disabled="salesTransactionsPage >= Math.ceil((salesTransactionsTotal || 0) / TRANSACTION_PAGE_SIZE)" @click="loadSalesTransactionsPage(salesTransactionsPage + 1)">Next</button>
            </div>
            <div class="flex gap-2">
              <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportDetailTransaksiExcel">Export Excel</button>
              <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportDetailTransaksiPdf">Export PDF</button>
            </div>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Date & Time</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Transaction ID</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Cashier</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Amount</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="row in salesTransactions" :key="row.id" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.created_at }}</td>
              <td class="px-4 py-2 text-sm text-gray-600">{{ row.id }}</td>
              <td class="px-4 py-2 text-sm text-gray-600">{{ row.cashier || '—' }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.total_amount) }}</td>
            </tr>
            <tr v-if="!salesTransactions?.length">
              <td colspan="4" class="px-4 py-4 text-sm text-gray-500 text-center">No transactions for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Profit & Loss -->
    <template v-else-if="activeTab === 'profit'">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <div class="text-sm text-gray-500">Revenue</div>
          <div class="text-xl font-semibold text-gray-800">{{ formatPrice(profitSummary?.revenue ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <div class="text-sm text-gray-500">Cost</div>
          <div class="text-xl font-semibold text-gray-800">{{ formatPrice(profitSummary?.cost ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <div class="text-sm text-gray-500">Profit</div>
          <div class="text-xl font-semibold" :class="(profitSummary?.profit ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'">{{ formatPrice(profitSummary?.profit ?? 0) }}</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between">
          <span class="font-medium text-gray-800">Tabel Produk</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(profitTableHeaders, profitTableRows, 'profit-produk')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Profit - Tabel Produk', profitTableHeaders, profitTableRows, 'profit-produk')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Cost</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Profit</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="(row, i) in profitRows" :key="i" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.quantity_sold }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.revenue) }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.cost) }}</td>
              <td class="px-4 py-2 text-sm text-right" :class="row.profit >= 0 ? 'text-green-600' : 'text-red-600'">{{ formatPrice(row.profit) }}</td>
            </tr>
            <tr v-if="!profitRows?.length">
              <td colspan="5" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Top Products -->
    <template v-else-if="activeTab === 'products'">
      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between">
          <span class="font-medium text-gray-800">Top Products</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Product', 'Quantity Sold', 'Revenue'], topProductsRows, 'top-products')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Top Products', ['Product', 'Quantity Sold', 'Revenue'], topProductsRows, 'top-products')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity Sold</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="(row, i) in topProducts" :key="i" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.quantity_sold }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.revenue) }}</td>
            </tr>
            <tr v-if="!topProducts?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Inventory -->
    <template v-else-if="activeTab === 'inventory'">
      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between">
          <span class="font-medium text-gray-800">Inventory</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(inventoryTableHeaders, inventoryTableRows, 'inventory')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Inventory', inventoryTableHeaders, inventoryTableRows, 'inventory')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Product</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Stock</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Inventory Value</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="row in inventoryRows" :key="row.product_id" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.product_name }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.stock }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.inventory_value ?? 0) }}</td>
            </tr>
            <tr v-if="!inventoryRows?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No inventory data.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Cashiers -->
    <template v-else-if="activeTab === 'cashiers'">
      <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
        <div class="p-3 border-b border-gray-200 flex items-center justify-between">
          <span class="font-medium text-gray-800">Cashiers</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Cashier', 'Transactions', 'Revenue'], cashiersTableRows, 'cashiers')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Cashiers', ['Cashier', 'Transactions', 'Revenue'], cashiersTableRows, 'cashiers')">Export PDF</button>
          </div>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Cashier</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Transactions</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="(row, i) in cashiersRows" :key="i" class="hover:bg-gray-50">
              <td class="px-4 py-2 text-sm text-gray-800">{{ row.cashier }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ row.transactions }}</td>
              <td class="px-4 py-2 text-sm text-right">{{ formatPrice(row.revenue) }}</td>
            </tr>
            <tr v-if="!cashiersRows?.length">
              <td colspan="3" class="px-4 py-4 text-sm text-gray-500 text-center">No data for this period.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import {
  getSalesSummary,
  getSalesDaily,
  getSalesTransactions,
  getSalesHourly,
  getPaymentsReport,
  getProfitReport,
  getTopProducts,
  getInventoryReport,
  getCashiersReport,
} from '../api/reports'
import { formatPrice } from '../utils'
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'

const dateFilterOptions = [
  { value: 'today', label: 'Today' },
  { value: 'week', label: 'This Week' },
  { value: 'month', label: 'This Month' },
  { value: 'custom', label: 'Custom Range' },
]

const tabs = [
  { id: 'sales', label: 'Sales' },
  { id: 'profit', label: 'Profit & Loss' },
  { id: 'products', label: 'Top Products' },
  { id: 'inventory', label: 'Inventory' },
  { id: 'cashiers', label: 'Cashiers' },
]

const dateFilter = ref('month')
const customFrom = ref('')
const customTo = ref('')
const activeTab = ref('sales')
const loading = ref(false)
const error = ref(null)

const TRANSACTION_PAGE_SIZE = 20
const salesSummary = ref(null)
const salesDaily = ref([])
const salesTransactions = ref([])
const salesTransactionsTotal = ref(0)
const salesTransactionsPage = ref(1)
const salesHourly = ref([])
const hourlyDate = ref('')
const paymentsReport = ref([])
const profitSummary = ref(null)
const profitRows = ref([])
const topProducts = ref([])
const inventoryRows = ref([])
const cashiersRows = ref([])

const dateRange = computed(() => {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  if (dateFilter.value === 'today') {
    return { from: `${y}-${m}-${d}`, to: `${y}-${m}-${d}` }
  }
  if (dateFilter.value === 'week') {
    const day = now.getDay()
    const diff = now.getDate() - day + (day === 0 ? -6 : 1)
    const monday = new Date(now)
    monday.setDate(diff)
    const from = monday.toISOString().slice(0, 10)
    return { from, to: `${y}-${m}-${d}` }
  }
  if (dateFilter.value === 'month') {
    return { from: `${y}-${m}-01`, to: `${y}-${m}-${d}` }
  }
  return { from: customFrom.value || undefined, to: customTo.value || undefined }
})

const dateRangeLabel = computed(() => {
  if (dateFilter.value === 'today') return 'Today'
  if (dateFilter.value === 'week') return 'This Week'
  if (dateFilter.value === 'month') return 'This Month'
  return `${dateRange.value.from || '?'} to ${dateRange.value.to || '?'}`
})

const activeTabLabel = computed(() => tabs.find((t) => t.id === activeTab.value)?.label ?? activeTab.value)

const salesDailyHeaders = ['Date', 'Transactions', 'Revenue']
const salesDailyRows = computed(() => (salesDaily.value || []).map((r) => [r.date, r.total_transactions, formatPrice(r.total_sales)]))

const salesTransactionsHeaders = ['Date & Time', 'Transaction ID', 'Cashier', 'Amount']
const salesTransactionsRows = computed(() => (salesTransactions.value || []).map((r) => [r.created_at, r.id, r.cashier || '—', formatPrice(r.total_amount)]))
const salesTransactionsShowing = computed(() => {
  const total = salesTransactionsTotal.value
  const page = salesTransactionsPage.value
  const limit = TRANSACTION_PAGE_SIZE
  if (total === 0) return { from: 0, to: 0, total: 0 }
  const from = (page - 1) * limit + 1
  const to = Math.min(page * limit, total)
  return { from, to, total }
})

const profitTableHeaders = ['Product', 'Quantity', 'Revenue', 'Cost', 'Profit']
const profitTableRows = computed(() => (profitRows.value || []).map((r) => [r.product_name, r.quantity_sold, formatPrice(r.revenue), formatPrice(r.cost), formatPrice(r.profit)]))

const topProductsRows = computed(() => (topProducts.value || []).map((r) => [r.product_name, r.quantity_sold, formatPrice(r.revenue)]))

const inventoryTableHeaders = ['Product', 'Stock', 'Inventory Value']
const inventoryTableRows = computed(() => (inventoryRows.value || []).map((r) => [r.product_name, r.stock, formatPrice(r.inventory_value ?? 0)]))

const cashiersTableRows = computed(() => (cashiersRows.value || []).map((r) => [r.cashier, r.transactions, formatPrice(r.revenue)]))

const paymentsTableRows = computed(() => (paymentsReport.value || []).map((r) => [r.method, r.transactions, formatPrice(r.revenue)]))

function setDateFilter(value) {
  dateFilter.value = value
  const now = new Date()
  if (value === 'custom') {
    customFrom.value = now.toISOString().slice(0, 10)
    customTo.value = now.toISOString().slice(0, 10)
  }
}

async function loadSales() {
  const range = dateRange.value
  if (!hourlyDate.value && range.to) hourlyDate.value = range.to
  const [summary, daily, payments, transactionsResp] = await Promise.all([
    getSalesSummary(range),
    getSalesDaily(range),
    getPaymentsReport(range),
    getSalesTransactions({ ...range, page: salesTransactionsPage.value, limit: TRANSACTION_PAGE_SIZE }),
  ])
  salesSummary.value = summary
  salesDaily.value = Array.isArray(daily) ? daily : []
  paymentsReport.value = Array.isArray(payments) ? payments : []
  salesTransactionsTotal.value = transactionsResp?.total ?? 0
  salesTransactions.value = Array.isArray(transactionsResp?.rows) ? transactionsResp.rows : []
  if (hourlyDate.value) await loadSalesHourly()
}

async function loadSalesHourly() {
  if (!hourlyDate.value) return
  try {
    const data = await getSalesHourly(hourlyDate.value)
    salesHourly.value = Array.isArray(data) ? data : []
  } catch {
    salesHourly.value = []
  }
}

async function loadSalesTransactionsPage(page) {
  if (page < 1) return
  const totalPages = Math.ceil(salesTransactionsTotal.value / TRANSACTION_PAGE_SIZE)
  if (totalPages > 0 && page > totalPages) return
  salesTransactionsPage.value = page
  const resp = await getSalesTransactions({
    ...dateRange.value,
    page: salesTransactionsPage.value,
    limit: TRANSACTION_PAGE_SIZE,
  })
  salesTransactions.value = Array.isArray(resp?.rows) ? resp.rows : []
}

async function fetchAllTransactionsForExport() {
  const resp = await getSalesTransactions({
    ...dateRange.value,
    page: 1,
    limit: 0,
  })
  const rows = Array.isArray(resp?.rows) ? resp.rows : []
  return rows.map((r) => [r.created_at, r.id, r.cashier || '—', formatPrice(r.total_amount)])
}

async function loadProfit() {
  const res = await getProfitReport(dateRange.value)
  profitSummary.value = res?.summary ?? null
  profitRows.value = Array.isArray(res?.rows) ? res.rows : []
}

async function loadProducts() {
  const data = await getTopProducts(dateRange.value)
  topProducts.value = Array.isArray(data) ? data : []
}

async function loadInventory() {
  const data = await getInventoryReport()
  inventoryRows.value = Array.isArray(data) ? data : []
}

async function loadCashiers() {
  const data = await getCashiersReport(dateRange.value)
  cashiersRows.value = Array.isArray(data) ? data : []
}

async function loadTab() {
  loading.value = true
  error.value = null
  try {
    if (activeTab.value === 'sales') await loadSales()
    else if (activeTab.value === 'profit') await loadProfit()
    else if (activeTab.value === 'products') await loadProducts()
    else if (activeTab.value === 'inventory') await loadInventory()
    else if (activeTab.value === 'cashiers') await loadCashiers()
  } catch (err) {
    error.value = err.response?.data?.error ?? 'Failed to load report.'
  } finally {
    loading.value = false
  }
}

watch(activeTab, loadTab)
watch(dateRange, () => {
  salesTransactionsPage.value = 1
  loadTab()
}, { deep: true })

onMounted(() => {
  if (dateFilter.value === 'custom' && !customFrom.value) setDateFilter('custom')
  loadTab()
})

function exportCardExcel(headers, rows, slug) {
  const data = Array.isArray(rows) ? rows : (rows?.value ?? [])
  const wb = XLSX.utils.book_new()
  const ws = XLSX.utils.aoa_to_sheet([headers, ...data])
  XLSX.utils.book_append_sheet(wb, ws, slug)
  XLSX.writeFile(wb, `report-${slug}-${dateRangeLabel.value.replace(/\s/g, '-')}.xlsx`)
}

function exportCardPdf(title, headers, rows, slug) {
  const data = Array.isArray(rows) ? rows : (rows?.value ?? [])
  const doc = new jsPDF({ orientation: 'landscape' })
  doc.setFontSize(14)
  doc.text(`${title} (${dateRangeLabel.value})`, 14, 15)
  doc.setFontSize(10)
  if (headers.length && data.length) {
    const colWidth = 270 / headers.length
    let y = 25
    doc.setFillColor(240, 240, 240)
    doc.rect(14, y, 270, 8, 'F')
    headers.forEach((h, i) => doc.text(String(h), 14 + i * colWidth + 2, y + 5.5))
    y += 8
    data.slice(0, 25).forEach((row) => {
      row.forEach((cell, i) => doc.text(String(cell), 14 + i * colWidth + 2, y + 5))
      y += 7
    })
    if (data.length > 25) doc.text(`... and ${data.length - 25} more rows`, 14, y + 5)
  } else {
    doc.text('No data to export.', 14, 25)
  }
  doc.save(`report-${slug}-${dateRangeLabel.value.replace(/\s/g, '-')}.pdf`)
}

async function exportDetailTransaksiExcel() {
  const rows = await fetchAllTransactionsForExport()
  exportCardExcel(salesTransactionsHeaders, rows, 'detail-transaksi')
}

async function exportDetailTransaksiPdf() {
  const rows = await fetchAllTransactionsForExport()
  exportCardPdf('Detail Transaksi', salesTransactionsHeaders, rows, 'detail-transaksi')
}
</script>

