<template>
  <div class="space-y-4">
    <div class="flex flex-col sm:flex-row sm:flex-wrap sm:items-center sm:justify-between gap-3 sm:gap-4">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-200">Reports</h1>
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-sm text-gray-600 dark:text-gray-400 shrink-0">Period:</span>
        <div class="flex flex-wrap items-center gap-2">
          <button
            v-for="opt in dateFilterOptions"
            :key="opt.value"
            type="button"
            :class="[dateFilter === opt.value ? 'bg-slate-600 text-white' : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700']"
            class="px-3 py-1.5 rounded text-sm font-medium"
            @click="setDateFilter(opt.value)"
          >
            {{ opt.label }}
          </button>
          <template v-if="dateFilter === 'custom'">
            <input v-model="customFrom" type="date" class="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 w-full sm:w-auto" />
            <span class="text-gray-400 dark:text-gray-400">to</span>
            <input v-model="customTo" type="date" class="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 w-full sm:w-auto" />
          </template>
        </div>
      </div>
    </div>

    <div class="border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
      <nav class="flex gap-2 sm:gap-4 min-w-max sm:min-w-0 py-1">
        <button
          v-for="t in tabs"
          :key="t.id"
          type="button"
          :class="[activeTab === t.id ? 'border-slate-600 text-slate-600 dark:text-slate-300' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200']"
          class="py-2 px-1 border-b-2 font-medium text-sm"
          @click="activeTab = t.id"
        >
          {{ t.label }}
        </button>
      </nav>
    </div>

    <div v-if="loading" class="text-gray-500 dark:text-gray-400 py-8">Loading...</div>
    <div v-else-if="error" class="text-red-600 py-4">{{ error }}</div>

    <!-- Sales -->
    <template v-else-if="activeTab === 'sales'">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
        <div class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="text-sm text-gray-500 dark:text-gray-400">Total Revenue</div>
          <div class="text-xl font-semibold text-gray-800 dark:text-gray-200">{{ formatPrice(salesSummary?.total_sales ?? 0) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="text-sm text-gray-500 dark:text-gray-400">Total Transactions</div>
          <div class="text-xl font-semibold text-gray-800 dark:text-gray-200">{{ salesSummary?.total_transactions ?? 0 }}</div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex flex-wrap items-center justify-between gap-2">
          <span class="font-medium text-gray-800 dark:text-gray-200">Payment Methods</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Method', 'Transactions', 'Revenue'], paymentsTableRows, 'payment-methods')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Payment Methods', ['Method', 'Transactions', 'Revenue'], paymentsTableRows, 'payment-methods')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="(row, i) in paymentsReport" :key="i" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.method }}</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>{{ row.transactions }} trans</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</span>
            </div>
          </div>
          <p v-if="!paymentsReport?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Method</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Transactions</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in paymentsReport" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.method }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.transactions }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</td>
              </tr>
              <tr v-if="!paymentsReport?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-wrap gap-2">
          <span class="font-medium text-gray-800 dark:text-gray-200">Ringkasan per Hari</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(salesDailyHeaders, salesDailyRows, 'ringkasan-per-hari')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Ringkasan per Hari', salesDailyHeaders, salesDailyRows, 'ringkasan-per-hari')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="row in salesDaily" :key="row.date" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.date }}</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>{{ row.total_transactions }} trans</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.total_sales) }}</span>
            </div>
          </div>
          <p v-if="!salesDaily?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Date</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Transactions</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="row in salesDaily" :key="row.date" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.date }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.total_transactions }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.total_sales) }}</td>
              </tr>
              <tr v-if="!salesDaily?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden mb-6">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-wrap gap-2">
          <span class="font-medium text-gray-800 dark:text-gray-200">Penjualan per Jam</span>
          <div class="flex items-center gap-2">
            <label class="text-sm text-gray-600 dark:text-gray-400">Tanggal:</label>
            <input v-model="hourlyDate" type="date" class="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100" @change="loadSalesHourly" />
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="row in salesHourly" :key="row.hour" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.hour }}:00</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>{{ row.transactions }} trans</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</span>
            </div>
          </div>
          <p v-if="hourlyDate && !salesHourly?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this date.</p>
          <p v-if="!hourlyDate" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">Pilih tanggal untuk melihat penjualan per jam.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Hour</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Transactions</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="row in salesHourly" :key="row.hour" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.hour }}:00</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.transactions }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</td>
              </tr>
              <tr v-if="hourlyDate && !salesHourly?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this date.</td>
              </tr>
              <tr v-if="!hourlyDate">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">Pilih tanggal untuk melihat penjualan per jam.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row sm:flex-wrap gap-2 justify-between">
          <span class="font-medium text-gray-800 dark:text-gray-200">Detail Transaksi</span>
          <div class="flex flex-wrap items-center gap-2">
            <span class="text-sm text-gray-500 dark:text-gray-400">Showing {{ salesTransactionsShowing.from }}–{{ salesTransactionsShowing.to }} of {{ salesTransactionsShowing.total }}</span>
            <div class="flex gap-1">
              <button type="button" class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 disabled:opacity-50" :disabled="salesTransactionsPage <= 1" @click="loadSalesTransactionsPage(salesTransactionsPage - 1)">Prev</button>
              <button type="button" class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 disabled:opacity-50" :disabled="salesTransactionsPage >= Math.ceil((salesTransactionsTotal || 0) / TRANSACTION_PAGE_SIZE)" @click="loadSalesTransactionsPage(salesTransactionsPage + 1)">Next</button>
            </div>
            <div class="flex gap-2">
              <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportDetailTransaksiExcel">Export Excel</button>
              <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportDetailTransaksiPdf">Export PDF</button>
            </div>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="row in salesTransactions" :key="row.id" class="p-4">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ row.created_at }}</p>
            <p class="font-medium text-gray-800 dark:text-gray-200 mt-0.5 font-mono text-xs break-all max-w-full" :title="row.id">#{{ row.id }}</p>
            <p class="text-sm text-gray-600 dark:text-gray-400">{{ row.customer_name || '—' }} {{ row.customer_phone ? '· ' + row.customer_phone : '' }}</p>
            <div class="flex items-center justify-between mt-2">
              <span class="text-sm font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.total_amount) }}</span>
              <button type="button" class="text-sm text-slate-600 dark:text-slate-400 hover:underline" @click="openReceiptModal(row.id)">View</button>
            </div>
          </div>
          <p v-if="!salesTransactions?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No transactions for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Date & Time</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Transaction ID</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Cashier</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Customer</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">No HP</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Amount</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Receipt</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="row in salesTransactions" :key="row.id" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.created_at }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 font-mono text-xs break-all max-w-[140px]" :title="row.id">{{ row.id }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ row.cashier || '—' }}</td>
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.customer_name || '—' }}</td>
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.customer_phone || '—' }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.total_amount) }}</td>
                <td class="px-4 py-2 text-right">
                  <button type="button" class="text-sm text-slate-600 dark:text-slate-400 hover:underline dark:hover:text-slate-300" @click="openReceiptModal(row.id)">View</button>
                </td>
              </tr>
              <tr v-if="!salesTransactions?.length">
                <td colspan="7" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No transactions for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Profit & Loss -->
    <template v-else-if="activeTab === 'profit'">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
        <div class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="text-sm text-gray-500 dark:text-gray-400">Revenue</div>
          <div class="text-xl font-semibold text-gray-800 dark:text-gray-200">{{ formatPrice(profitSummary?.revenue ?? 0) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="text-sm text-gray-500 dark:text-gray-400">Cost</div>
          <div class="text-xl font-semibold text-gray-800 dark:text-gray-200">{{ formatPrice(profitSummary?.cost ?? 0) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800 p-4">
          <div class="text-sm text-gray-500 dark:text-gray-400">Profit</div>
          <div class="text-xl font-semibold" :class="(profitSummary?.profit ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'">{{ formatPrice(profitSummary?.profit ?? 0) }}</div>
        </div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <span class="font-medium text-gray-800 dark:text-gray-200">Tabel Produk</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(profitTableHeaders, profitTableRows, 'profit-produk')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Profit - Tabel Produk', profitTableHeaders, profitTableRows, 'profit-produk')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="(row, i) in profitRows" :key="i" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.product_name }}</p>
            <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>Qty: {{ row.quantity_sold }}</span>
              <span>Rev: {{ formatPrice(row.revenue) }}</span>
              <span>Cost: {{ formatPrice(row.cost) }}</span>
            </div>
            <p class="mt-1 text-sm font-medium" :class="row.profit >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">Profit: {{ formatPrice(row.profit) }}</p>
          </div>
          <p v-if="!profitRows?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Quantity</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Cost</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Profit</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in profitRows" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.quantity_sold }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.cost) }}</td>
                <td class="px-4 py-2 text-sm text-right" :class="row.profit >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">{{ formatPrice(row.profit) }}</td>
              </tr>
              <tr v-if="!profitRows?.length">
                <td colspan="5" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Top Products -->
    <template v-else-if="activeTab === 'products'">
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <span class="font-medium text-gray-800 dark:text-gray-200">Top Products</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Product', 'Quantity Sold', 'Revenue'], topProductsRows, 'top-products')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Top Products', ['Product', 'Quantity Sold', 'Revenue'], topProductsRows, 'top-products')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="(row, i) in topProducts" :key="i" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.product_name }}</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>{{ row.quantity_sold }} sold</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</span>
            </div>
          </div>
          <p v-if="!topProducts?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Quantity Sold</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in topProducts" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.quantity_sold }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</td>
              </tr>
              <tr v-if="!topProducts?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Product Margin -->
    <template v-else-if="activeTab === 'margin'">
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="r in marginRows" :key="r.product_id" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ r.product_name }}</p>
            <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>Rev: {{ formatPrice(r.revenue) }}</span>
              <span>COGS: {{ formatPrice(r.cogs) }}</span>
            </div>
            <p class="mt-1 text-sm" :class="r.margin >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">Margin: {{ formatPrice(r.margin) }} ({{ r.margin_pct }}%)</p>
          </div>
          <p v-if="!marginRows.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">COGS</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Margin</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Margin %</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="r in marginRows" :key="r.product_id" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
              <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ r.product_name }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(r.revenue) }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(r.cogs) }}</td>
              <td class="px-4 py-2 text-sm text-right" :class="r.margin >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">{{ formatPrice(r.margin) }}</td>
              <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ r.margin_pct }}%</td>
            </tr>
            <tr v-if="!marginRows.length"><td colspan="5" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data</td></tr>
          </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Inventory -->
    <template v-else-if="activeTab === 'inventory'">
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <span class="font-medium text-gray-800 dark:text-gray-200">Inventory</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(inventoryTableHeaders, inventoryTableRows, 'inventory')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Inventory', inventoryTableHeaders, inventoryTableRows, 'inventory')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="row in inventoryRows" :key="row.product_id" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.product_name }}</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>Stock: {{ row.stock }}</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.inventory_value ?? 0) }}</span>
            </div>
          </div>
          <p v-if="!inventoryRows?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No inventory data.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Product</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Stock</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Inventory Value</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="row in inventoryRows" :key="row.product_id" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.product_name }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.stock }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.inventory_value ?? 0) }}</td>
              </tr>
              <tr v-if="!inventoryRows?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No inventory data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Cashiers -->
    <template v-else-if="activeTab === 'cashiers'">
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <span class="font-medium text-gray-800 dark:text-gray-200">Cashiers</span>
          <div class="flex gap-2">
            <button type="button" class="px-2 py-1.5 bg-green-600 text-white rounded text-xs hover:bg-green-700" @click="exportCardExcel(['Cashier', 'Transactions', 'Revenue'], cashiersTableRows, 'cashiers')">Export Excel</button>
            <button type="button" class="px-2 py-1.5 bg-red-600 text-white rounded text-xs hover:bg-red-700" @click="exportCardPdf('Cashiers', ['Cashier', 'Transactions', 'Revenue'], cashiersTableRows, 'cashiers')">Export PDF</button>
          </div>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="(row, i) in cashiersRows" :key="i" class="p-4">
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.cashier }}</p>
            <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>{{ row.transactions }} trans</span>
              <span class="font-medium text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</span>
            </div>
          </div>
          <p v-if="!cashiersRows?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Cashier</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Transactions</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Revenue</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in cashiersRows" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.cashier }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ row.transactions }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.revenue) }}</td>
              </tr>
              <tr v-if="!cashiersRows?.length">
                <td colspan="3" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No data for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Shifts (cash reconciliation) -->
    <template v-else-if="activeTab === 'shifts'">
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow border border-gray-200 dark:border-gray-800 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700">
          <span class="font-medium text-gray-800 dark:text-gray-200">Shift reconciliation</span>
        </div>
        <div class="sm:hidden divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="(row, i) in shiftsRows" :key="i" class="p-4">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ row.date }}</p>
            <p class="font-medium text-gray-800 dark:text-gray-200">{{ row.cashier }}</p>
            <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-sm text-gray-600 dark:text-gray-400 mt-1">
              <span>Open: {{ formatPrice(row.opening) }}</span>
              <span>Sales: {{ formatPrice(row.sales) }}</span>
              <span>Actual: {{ formatPrice(row.actual) }}</span>
            </div>
            <p class="mt-1 text-sm font-medium" :class="row.difference !== 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-600 dark:text-gray-400'">Diff: {{ formatPrice(row.difference) }}</p>
          </div>
          <p v-if="!shiftsRows?.length" class="p-4 text-sm text-gray-500 dark:text-gray-400 text-center">No closed shifts for this period.</p>
        </div>
        <div class="hidden sm:block overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Date</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Cashier</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Opening</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Sales</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Expected</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Actual</th>
                <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Difference</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="(row, i) in shiftsRows" :key="i" class="hover:bg-gray-50 dark:hover:bg-gray-800 dark:bg-gray-800">
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.date }}</td>
                <td class="px-4 py-2 text-sm text-gray-800 dark:text-gray-200">{{ row.cashier }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.opening) }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.sales) }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.expected) }}</td>
                <td class="px-4 py-2 text-sm text-right text-gray-800 dark:text-gray-200">{{ formatPrice(row.actual) }}</td>
                <td class="px-4 py-2 text-sm text-right" :class="row.difference !== 0 ? 'text-amber-600 dark:text-amber-400 font-medium' : 'text-gray-800 dark:text-gray-200'">{{ formatPrice(row.difference) }}</td>
              </tr>
              <tr v-if="!shiftsRows?.length">
                <td colspan="7" class="px-4 py-4 text-sm text-gray-500 dark:text-gray-400 text-center">No closed shifts for this period.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Receipt Modal -->
    <div
      v-if="showReceipt"
      class="fixed inset-0 z-20 flex items-center justify-center bg-black/50 p-4"
      @click.self="showReceipt = false"
    >
      <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl max-w-sm w-full max-h-[90vh] overflow-auto">
        <div class="p-4" v-if="receiptLoading">
          <p class="text-gray-500 dark:text-gray-400 text-center py-8">Loading receipt...</p>
        </div>
        <div class="p-4" v-else-if="receiptError">
          <p class="text-red-600 text-center py-4">{{ receiptError }}</p>
        </div>
        <div class="p-4" v-else-if="receiptObj">
          <Receipt
            :store-name="tenantStore.storeName()"
            :date="receiptObj.transaction.created_at"
            :transaction-id="receiptObj.transaction.id"
            :cashier="receiptObj.transaction.cashier"
            :items="receiptObj.items.map(i => ({ name: i.product_name, quantity: i.quantity, price: i.price, product_id: i.product_name }))"
            :total="receiptObj.transaction.total_amount"
            :paid-amount="receiptObj.payments.reduce((s, p) => s + p.amount, 0)"
            :change="Math.max(0, receiptObj.payments.reduce((s, p) => s + p.amount, 0) - receiptObj.transaction.total_amount)"
          />
        </div>
        <div class="p-4 border-t border-gray-200 dark:border-gray-700 flex flex-wrap gap-2 justify-end">
          <button type="button" class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300 text-sm" @click="printReceiptFromReport">Print</button>
          <button type="button" class="px-3 py-2 border border-green-500 dark:border-green-600 text-green-700 dark:text-green-400 rounded-md hover:bg-green-50 dark:hover:bg-green-900/20 text-sm" @click="downloadReceiptPdfFromReport">PDF</button>
          <button type="button" class="px-3 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm" @click="shareWhatsAppFromReport">WhatsApp</button>
          <button type="button" class="px-3 py-2 bg-slate-600 text-white rounded-md hover:bg-slate-700 text-sm" @click="showReceipt = false">Close</button>
        </div>
      </div>
    </div>
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
  getShiftsReport,
  getProductMargin,
} from '../api/reports'
import { formatPrice, formatDate, formatDateTime } from '../utils'
import { getReceipt } from '../api/receipt'
import { generateReceiptPDF, buildReceiptText } from '../utils/receipt-pdf'
import { useTenantStore } from '../stores/tenant'
import Receipt from '../components/Receipt.vue'
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'

const tenantStore = useTenantStore()

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
  { id: 'margin', label: 'Product Margin' },
  { id: 'inventory', label: 'Inventory' },
  { id: 'cashiers', label: 'Cashiers' },
  { id: 'shifts', label: 'Shifts' },
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
const shiftsRows = ref([])
const marginRows = ref([])

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

const salesTransactionsHeaders = ['Date & Time', 'Transaction ID', 'Cashier', 'Customer', 'No HP', 'Amount']
const salesTransactionsRows = computed(() => (salesTransactions.value || []).map((r) => [r.created_at, r.id, r.cashier || '—', r.customer_name || '—', r.customer_phone || '—', formatPrice(r.total_amount)]))
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
  return rows.map((r) => [r.created_at, r.id, r.cashier || '—', r.customer_name || '—', r.customer_phone || '—', formatPrice(r.total_amount)])
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

async function loadShifts() {
  const data = await getShiftsReport(dateRange.value)
  shiftsRows.value = Array.isArray(data) ? data : []
}

async function loadMargin() {
  const data = await getProductMargin(dateRange.value)
  marginRows.value = Array.isArray(data) ? data : []
}

async function loadTab() {
  loading.value = true
  error.value = null
  try {
    if (activeTab.value === 'sales') await loadSales()
    else if (activeTab.value === 'profit') await loadProfit()
    else if (activeTab.value === 'products') await loadProducts()
    else if (activeTab.value === 'inventory') await loadInventory()
    else if (activeTab.value === 'margin') await loadMargin()
    else if (activeTab.value === 'cashiers') await loadCashiers()
    else if (activeTab.value === 'shifts') await loadShifts()
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

const showReceipt = ref(false)
const receiptObj = ref(null)
const receiptLoading = ref(false)
const receiptError = ref('')

async function openReceiptModal(txnId) {
  showReceipt.value = true
  receiptLoading.value = true
  receiptError.value = ''
  receiptObj.value = null
  try {
    receiptObj.value = await getReceipt(txnId)
  } catch (e) {
    receiptError.value = e.response?.data?.error ?? 'Failed to load receipt.'
  } finally {
    receiptLoading.value = false
  }
}

function receiptDataFromObj() {
  if (!receiptObj.value) return null
  const r = receiptObj.value
  const paid = r.payments.reduce((s, p) => s + p.amount, 0)
  return {
    storeName: tenantStore.storeName(),
    date: r.transaction.created_at,
    transactionId: r.transaction.id,
    cashier: r.transaction.cashier,
    items: r.items.map(i => ({ name: i.product_name, quantity: i.quantity, price: i.price })),
    total: r.transaction.total_amount,
    paidAmount: paid,
    change: Math.max(0, paid - r.transaction.total_amount),
  }
}

function printReceiptFromReport() { window.print() }

function downloadReceiptPdfFromReport() {
  const data = receiptDataFromObj()
  if (!data) return
  const doc = generateReceiptPDF(data)
  doc.save(`receipt-${data.transactionId || 'txn'}.pdf`)
}

function shareWhatsAppFromReport() {
  const data = receiptDataFromObj()
  if (!data) return
  const text = buildReceiptText(data)
  window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, '_blank')
}

async function exportCardExcel(headers, rows, slug) {
  await tenantStore.load()
  const p = tenantStore.profile
  const storeName = p?.name || tenantStore.storeName()
  const address = p?.address || '—'
  const description = p?.description || '—'
  const printDate = new Date().toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
  const fmt = (d) => (d ? d.split('-').reverse().join('/') : '—')
  const periodFrom = fmt(dateRange.value.from)
  const periodTo = fmt(dateRange.value.to)
  const periodLabel = `Periode Laporan: ${periodFrom} - ${periodTo}`

  const data = Array.isArray(rows) ? rows : (rows?.value ?? [])
  const wb = XLSX.utils.book_new()
  const headerRows = [
    ['Toko / Store', storeName],
    ['Alamat', address],
    ['Deskripsi', description],
    ['Tanggal Cetak', printDate],
    [periodLabel, ''],
    [],
    headers,
    ...data,
  ]
  const ws = XLSX.utils.aoa_to_sheet(headerRows)
  ws['!cols'] = headers.map((_, i) => ({ wch: Math.min(Math.max(12, 8 + (data.reduce((m, r) => Math.max(m, String(r[i] ?? '').length), 0))), 50) }))
  XLSX.utils.book_append_sheet(wb, ws, slug)
  XLSX.writeFile(wb, `report-${slug}-${dateRangeLabel.value.replace(/\s/g, '-')}.xlsx`)
}

async function exportCardPdf(title, headers, rows, slug, opts = {}) {
  await tenantStore.load()
  const p = tenantStore.profile
  const storeName = p?.name || tenantStore.storeName()
  const address = p?.address || '—'
  const description = p?.description || '—'
  const printDate = new Date().toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
  const fmt = (d) => (d ? d.split('-').reverse().join('/') : '—')
  const periodFrom = fmt(dateRange.value.from)
  const periodTo = fmt(dateRange.value.to)
  const periodLabel = `Periode Laporan: ${periodFrom} - ${periodTo}`

  const data = Array.isArray(rows) ? rows : (rows?.value ?? [])
  const doc = new jsPDF({ orientation: 'landscape' })
  const pageW = doc.internal.pageSize.getWidth()
  const margin = 14
  const contentW = pageW - margin * 2
  let y = 12

  doc.setFontSize(10)
  doc.text(`Toko: ${storeName}`, margin, y)
  y += 6
  doc.text(`Alamat: ${address}`, margin, y)
  y += 6
  doc.text(`Deskripsi: ${description}`, margin, y)
  y += 6
  doc.text(`Tanggal Cetak: ${printDate}`, margin, y)
  y += 6
  doc.text(periodLabel, margin, y)
  y += 8

  doc.setFontSize(14)
  doc.text(`${title} (${dateRangeLabel.value})`, margin, y)
  y += 8

  doc.setFontSize(10)
  if (headers.length && data.length) {
    const colWidths = opts.colWidths || headers.map(() => contentW / headers.length)
    const rightAlignCols = opts.rightAlignCols || []
    const cellPadding = 2
    const lineH = 5
    const totalColW = colWidths.reduce((a, b) => a + b, 0)
    const scale = contentW / totalColW
    const scaledWidths = colWidths.map((w) => w * scale)

    const wrapText = (text, w) => doc.splitTextToSize(String(text), Math.max(5, w - cellPadding * 2))

    doc.setFillColor(240, 240, 240)
    doc.rect(margin, y, contentW, 8, 'F')
    headers.forEach((h, i) => {
      const xStart = margin + scaledWidths.slice(0, i).reduce((a, b) => a + b, 0)
      const cw = scaledWidths[i]
      doc.text(String(h), rightAlignCols.includes(i) ? xStart + cw - cellPadding : xStart + cellPadding, y + 5.5, { align: rightAlignCols.includes(i) ? 'right' : 'left' })
    })
    y += 8
    data.slice(0, 25).forEach((row) => {
      const cellLines = row.map((cell, i) => wrapText(String(cell), scaledWidths[i]))
      const maxLines = Math.max(1, ...cellLines.map((l) => l.length))
      const rowHeight = maxLines * lineH + 2
      headers.forEach((_, i) => {
        const xStart = margin + scaledWidths.slice(0, i).reduce((a, b) => a + b, 0)
        const cw = scaledWidths[i]
        const maxCellW = Math.max(5, cw - cellPadding * 2)
        const lines = cellLines[i] || ['']
        const isRight = rightAlignCols.includes(i)
        lines.forEach((line, li) => {
          doc.text(line, isRight ? xStart + cw - cellPadding : xStart + cellPadding, y + 4 + li * lineH, { align: isRight ? 'right' : 'left', maxWidth: maxCellW })
        })
      })
      y += rowHeight
    })
    if (data.length > 25) doc.text(`... and ${data.length - 25} more rows`, margin, y + 5)
  } else {
    doc.text('No data to export.', margin, y + 5)
  }
  doc.save(`report-${slug}-${dateRangeLabel.value.replace(/\s/g, '-')}.pdf`)
}

async function exportDetailTransaksiExcel() {
  const rows = await fetchAllTransactionsForExport()
  await exportCardExcel(salesTransactionsHeaders, rows, 'detail-transaksi')
}

async function exportDetailTransaksiPdf() {
  const rows = await fetchAllTransactionsForExport()
  const colWidths = [38, 42, 35, 40, 30, 35]
  await exportCardPdf('Detail Transaksi', salesTransactionsHeaders, rows, 'detail-transaksi', {
    colWidths,
    rightAlignCols: [5],
  })
}
</script>

