<template>
  <div class="h-full flex flex-col -m-4 lg:-m-6">
    <!-- Compact top bar -->
    <div class="flex items-center gap-3 px-4 py-2.5 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shrink-0">
      <template v-if="isCashier && currentShift">
        <div class="flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/30 px-2 py-1 rounded-full">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
          Shift open
        </div>
        <button type="button" class="text-xs text-gray-500 dark:text-gray-400 hover:text-amber-600 dark:hover:text-amber-400 transition-colors" @click="shiftModalMode = 'close'">Close shift</button>
      </template>
      <template v-else-if="isCashier && !currentShift">
        <button type="button" class="px-3 py-1.5 text-xs font-medium bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors" @click="shiftModalMode = 'open'">Open Shift</button>
      </template>
      <div class="flex-1" />
      <div class="relative max-w-xs w-full">
        <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input ref="barcodeInputRef" v-model="barcodeInput" type="text" placeholder="Scan / cari barcode..."
          class="w-full pl-8 pr-3 py-1.5 text-sm border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent placeholder-gray-400 dark:placeholder-gray-500"
          @keydown.enter.prevent="onBarcodeEnter" />
      </div>
      <div class="text-[10px] text-gray-400 dark:text-gray-600 hidden xl:flex items-center gap-1.5">
        <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 font-mono">F1</kbd>Cash
        <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 font-mono">F2</kbd>Card
        <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 font-mono">F3</kbd>QRIS
        <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 font-mono">Esc</kbd>Clear
      </div>
    </div>

    <!-- Main content -->
    <div class="flex-1 min-h-0 flex">
      <!-- LEFT: Product grid -->
      <div class="flex-1 flex flex-col min-w-0">
        <div class="px-4 py-2.5 border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shrink-0 flex flex-col sm:flex-row gap-2">
          <div class="relative flex-1 min-w-0">
            <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            <input v-model="searchQuery" type="text" placeholder="Cari produk..."
              class="w-full pl-8 pr-3 py-2 text-sm border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent placeholder-gray-400 dark:placeholder-gray-500" />
          </div>
          <select v-model="posCategoryFilter" class="px-3 py-2 text-sm border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent min-w-[140px]">
            <option value="">Semua Kategori</option>
            <option v-for="c in posCategories" :key="posCatId(c)" :value="posCatId(c)">{{ posCatName(c) }}</option>
          </select>
        </div>
        <div class="flex-1 overflow-auto p-3 pb-44 lg:pb-3 bg-gray-50 dark:bg-gray-950">
          <p v-if="productsLoading" class="text-gray-400 dark:text-gray-600 text-sm py-8 text-center">Loading...</p>
          <p v-else-if="productsError" class="text-red-500 text-sm py-4 text-center">{{ productsError }}</p>
          <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-2">
            <button
              v-for="p in filteredProducts" :key="productId(p)" type="button"
              :class="['group rounded-xl p-3 text-left transition-all active:scale-[0.98]', clickedProductId === productId(p) ? 'animate-pos-product-click ring-2 ring-indigo-400 dark:ring-indigo-500 bg-indigo-50 dark:bg-indigo-900/40 border border-indigo-200 dark:border-indigo-800' : 'bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 hover:border-indigo-300 dark:hover:border-indigo-700 hover:shadow-sm']"
              @click="onProductClick(p)"
            >
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200 truncate group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">{{ productName(p) }}</div>
              <div class="text-xs text-gray-400 dark:text-gray-600 truncate mt-0.5">{{ productSku(p) || '—' }}</div>
              <div class="text-sm font-semibold text-gray-900 dark:text-white mt-1.5">{{ formatPrice(productPrice(p)) }}</div>
            </button>
          </div>
          <p v-if="!productsLoading && !productsError && !filteredProducts.length" class="text-gray-400 dark:text-gray-600 text-sm py-8 text-center">Tidak ada produk ditemukan</p>
        </div>
      </div>

      <!-- RIGHT: Cart panel -->
      <div class="w-80 xl:w-96 shrink-0 flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-800 hidden lg:flex">
        <!-- Cart header -->
        <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-800 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/></svg>
            <span class="text-sm font-semibold text-gray-800 dark:text-gray-200">Cart</span>
            <span v-if="cartItems.length" class="text-[10px] bg-indigo-100 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400 px-1.5 py-0.5 rounded-full font-medium">{{ cartItems.length }}</span>
          </div>
          <button v-if="cartItems.length" @click="confirmClearCart" class="text-[10px] text-gray-400 hover:text-red-500 dark:hover:text-red-400 transition-colors uppercase tracking-wider font-medium">Clear</button>
        </div>

        <!-- Cart items -->
        <div class="flex-1 overflow-auto min-h-0">
          <div v-if="!cartItems.length" class="flex flex-col items-center justify-center h-full text-gray-300 dark:text-gray-700">
            <svg class="w-12 h-12 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/></svg>
            <span class="text-xs">Keranjang kosong</span>
          </div>
          <div v-else class="divide-y divide-gray-100 dark:divide-gray-800">
            <div v-for="item in cartItems" :key="item.product_id" class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <div class="text-sm font-medium text-gray-800 dark:text-gray-200 truncate">{{ item.name }}</div>
                  <div class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ formatPrice(item.price) }} x {{ item.quantity }}</div>
                </div>
                <div class="text-sm font-semibold text-gray-900 dark:text-white shrink-0">{{ formatPrice(item.price * item.quantity) }}</div>
              </div>
              <div class="flex items-center gap-1 mt-2">
                <button type="button" class="w-7 h-7 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center text-xs transition-colors" :disabled="item.quantity <= 1" @click="changeQuantity(item.product_id, -1)">−</button>
                <span class="w-8 text-center text-sm font-medium text-gray-800 dark:text-gray-200">{{ item.quantity }}</span>
                <button type="button" class="w-7 h-7 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center text-xs transition-colors" @click="changeQuantity(item.product_id, 1)">+</button>
                <button type="button" class="px-1.5 py-1 rounded border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 text-[10px] font-medium transition-colors" @click="changeQuantity(item.product_id, 5)">+5</button>
                <button type="button" class="px-1.5 py-1 rounded border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 text-[10px] font-medium transition-colors" @click="changeQuantity(item.product_id, 10)">+10</button>
                <div class="flex-1" />
                <button type="button" class="text-gray-300 dark:text-gray-600 hover:text-red-500 dark:hover:text-red-400 transition-colors" @click="removeFromCart(item.product_id)">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Cart footer -->
        <div class="border-t border-gray-200 dark:border-gray-800 p-4 shrink-0 space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-500 dark:text-gray-400">Total</span>
            <span class="text-xl font-bold text-gray-900 dark:text-white">{{ formatPrice(totalAmount) }}</span>
          </div>
          <p v-if="isCashier && !currentShift" class="text-xs text-amber-600 dark:text-amber-400">Buka shift untuk mulai transaksi</p>
          <div class="grid grid-cols-5 gap-1.5">
            <button v-for="pm in paymentMethods" :key="pm.value" type="button"
              class="flex flex-col items-center justify-center py-2.5 rounded-xl text-[10px] font-semibold transition-all disabled:opacity-30 disabled:cursor-not-allowed active:scale-95"
              :class="pm.class"
              :disabled="!cartItems.length || !canCheckout"
              @click="openCheckoutModal(pm.value)"
              :title="pm.label">
              <svg class="w-5 h-5 mb-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" v-html="pm.icon" />
              <span>{{ pm.label }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile bottom bar: integrated cart + payment (bukan overlay) -->
    <Teleport to="body">
      <div v-if="cartItems.length" class="lg:hidden fixed inset-x-0 bottom-0 z-30 safe-bottom">
        <!-- Single panel: cart list (expandable) + payment buttons selalu terlihat -->
        <div class="bg-indigo-50 dark:bg-gray-800 border-t border-indigo-200 dark:border-gray-700 rounded-t-2xl shadow-2xl flex flex-col max-h-[85vh]">
          <!-- Cart list (expandable, bagian dari panel, bukan overlay) -->
          <Transition
            enter-active-class="transition-all duration-200 ease-out"
            enter-from-class="max-h-0 opacity-0"
            enter-to-class="max-h-[50vh] opacity-100"
            leave-active-class="transition-all duration-150 ease-in"
            leave-from-class="max-h-[50vh] opacity-100"
            leave-to-class="max-h-0 opacity-0"
          >
            <div v-if="mobileCartExpanded" class="overflow-hidden shrink-0">
              <button type="button" class="flex justify-center pt-2 pb-1 w-full touch-none" @click="mobileCartExpanded = false">
                <span class="w-10 h-1 rounded-full bg-indigo-300 dark:bg-gray-500" />
              </button>
              <div class="px-4 py-2 flex items-center justify-between border-b border-indigo-200/60 dark:border-gray-700">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/></svg>
                  <span class="text-sm font-semibold text-gray-800 dark:text-gray-200">Keranjang</span>
                  <span class="text-[10px] bg-indigo-200 dark:bg-indigo-900/60 text-indigo-700 dark:text-indigo-300 px-1.5 py-0.5 rounded-full font-medium">{{ cartItems.length }}</span>
                </div>
                <button type="button" class="text-[10px] text-gray-400 dark:text-gray-500 hover:text-red-500 dark:hover:text-red-400 uppercase font-semibold tracking-wide" @click="confirmClearCart">Clear</button>
              </div>
              <div class="overflow-auto max-h-[45vh] overscroll-contain border-b border-indigo-200/40 dark:border-gray-700">
                <div class="divide-y divide-indigo-100 dark:divide-gray-700">
                  <div v-for="item in cartItems" :key="item.product_id" class="px-4 py-3 bg-white/60 dark:bg-gray-800/80">
                    <div class="flex items-start justify-between gap-2">
                      <div class="min-w-0">
                        <div class="text-sm font-medium text-gray-800 dark:text-gray-200 truncate">{{ item.name }}</div>
                        <div class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ formatPrice(item.price) }} × {{ item.quantity }}</div>
                      </div>
                      <div class="text-sm font-semibold text-gray-900 dark:text-white shrink-0">{{ formatPrice(item.price * item.quantity) }}</div>
                    </div>
                    <div class="flex items-center gap-1 mt-2">
                      <button type="button" class="w-8 h-8 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center text-sm transition-colors bg-white dark:bg-gray-900" :disabled="item.quantity <= 1" @click="changeQuantity(item.product_id, -1)">−</button>
                      <span class="w-8 text-center text-sm font-medium text-gray-800 dark:text-gray-200">{{ item.quantity }}</span>
                      <button type="button" class="w-8 h-8 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center text-sm transition-colors bg-white dark:bg-gray-900" @click="changeQuantity(item.product_id, 1)">+</button>
                      <button type="button" class="px-2 py-1 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 text-[10px] font-medium transition-colors bg-white dark:bg-gray-900" @click="changeQuantity(item.product_id, 5)">+5</button>
                      <button type="button" class="px-2 py-1 rounded-lg border border-gray-200 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 text-[10px] font-medium transition-colors bg-white dark:bg-gray-900" @click="changeQuantity(item.product_id, 10)">+10</button>
                      <div class="flex-1" />
                      <button type="button" class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-red-500 dark:hover:text-red-400 transition-colors" @click="removeFromCart(item.product_id)">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </Transition>

          <!-- Summary + payment buttons (selalu terlihat, bisa klik meskipun cart terbuka) -->
          <div class="shrink-0 bg-indigo-50 dark:bg-gray-800">
          <!-- Tap to expand cart list -->
          <button type="button" class="w-full px-4 py-2.5 flex items-center justify-between touch-none" @click="mobileCartExpanded = !mobileCartExpanded">
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/></svg>
              <span class="text-sm font-semibold text-gray-800 dark:text-gray-200">{{ cartItems.length }} item</span>
              <span class="text-base font-bold text-gray-900 dark:text-white">{{ formatPrice(totalAmount) }}</span>
            </div>
            <svg :class="['w-4 h-4 text-gray-400 dark:text-gray-500 transition-transform', mobileCartExpanded ? 'rotate-180' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"/></svg>
          </button>
          <!-- Payment buttons (always visible) -->
          <div class="px-4 pb-4 pt-1">
            <p v-if="isCashier && !currentShift" class="text-xs text-amber-600 dark:text-amber-400 mb-2">Buka shift untuk mulai transaksi</p>
            <div class="grid grid-cols-5 gap-1.5">
              <button v-for="pm in paymentMethods" :key="pm.value" type="button"
                class="flex flex-col items-center justify-center py-2.5 rounded-xl text-[10px] font-semibold transition-all disabled:opacity-30 active:scale-95"
                :class="pm.class"
                :disabled="!canCheckout"
                @click="openCheckoutModal(pm.value)"
                :title="pm.label">
                <svg class="w-5 h-5 mb-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" v-html="pm.icon" />
                <span>{{ pm.label }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      </div>
    </Teleport>

    <!-- Checkout modal -->
    <Teleport to="body">
      <div v-if="showCheckoutModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" @click.self="showCheckoutModal = false">
        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100 dark:border-gray-800">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Checkout</h2>
          </div>
          <form @submit.prevent="submitCheckout" class="p-6 space-y-4">
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5 uppercase tracking-wider">Nama (opsional)</label>
              <input v-model="checkoutForm.customer_name" type="text" placeholder="Nama pelanggan"
                class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-gray-50 dark:bg-gray-800 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5 uppercase tracking-wider">No HP (opsional)</label>
              <input v-model="checkoutForm.customer_phone" type="tel" placeholder="08xxxxxxxxxx - untuk WhatsApp struk"
                class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-gray-50 dark:bg-gray-800 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5 uppercase tracking-wider">Metode Pembayaran</label>
              <select v-model="checkoutForm.payment_method"
                class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-gray-50 dark:bg-gray-800 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500">
                <option value="cash">Cash</option>
                <option value="card">Card</option>
                <option value="qris">QRIS</option>
                <option value="ewallet">E-Wallet</option>
                <option value="transfer">Transfer</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5 uppercase tracking-wider">Jumlah Bayar</label>
              <input v-model.number="checkoutForm.paid_amount" type="number" min="0" step="0.01" required
                class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-gray-50 dark:bg-gray-800 text-sm text-gray-800 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500" />
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Total: {{ formatPrice(totalAmount) }}</p>
            </div>
            <p v-if="checkoutError" class="text-xs text-red-500">{{ checkoutError }}</p>
            <div class="flex gap-2 pt-2">
              <button type="button" class="flex-1 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl text-sm font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors" @click="showCheckoutModal = false">Batal</button>
              <button type="submit" :disabled="checkoutSubmitting" class="flex-1 py-2.5 bg-indigo-600 text-white rounded-xl text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition-colors">
                {{ checkoutSubmitting ? 'Memproses...' : 'Bayar' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Shift modal -->
    <ShiftModal v-if="shiftModalMode" :mode="shiftModalMode" @done="onShiftModalDone" @close="onShiftModalClose" />

    <!-- Toast: stok habis / stok rendah (auto-hide) -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 translate-y-2"
      >
        <div
          v-if="toast.show"
          class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 px-4 py-3 rounded-lg shadow-lg text-sm font-medium max-w-[90vw] sm:max-w-sm"
          :class="toast.type === 'error' ? 'bg-red-600 text-white dark:bg-red-700' : 'bg-amber-500 text-white dark:bg-amber-600'"
        >
          {{ toast.message }}
        </div>
      </Transition>
    </Teleport>

    <!-- Receipt modal -->
    <Teleport to="body">
      <div v-if="showReceiptModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" @click.self="closeReceipt">
        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl max-w-sm w-full max-h-[90vh] overflow-auto">
          <div class="p-5">
            <Receipt v-if="receiptData" :store-name="receiptData.storeName" :date="receiptData.date" :transaction-id="receiptData.transactionId" :cashier="receiptData.cashier" :items="receiptData.items" :total="receiptData.total" :paid-amount="receiptData.paidAmount" :change="receiptData.change" />
          </div>
          <div class="p-4 border-t border-gray-100 dark:border-gray-800 flex flex-wrap gap-2 justify-end">
            <button type="button" class="px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-800 text-sm text-gray-700 dark:text-gray-300 transition-colors" @click="printReceipt">Print</button>
            <button type="button" class="px-3 py-2 border border-emerald-300 dark:border-emerald-700 text-emerald-600 dark:text-emerald-400 rounded-xl hover:bg-emerald-50 dark:hover:bg-emerald-900/20 text-sm transition-colors" @click="downloadReceiptPDF">PDF</button>
            <button type="button" class="px-3 py-2 bg-emerald-600 text-white rounded-xl hover:bg-emerald-700 text-sm transition-colors" @click="shareWhatsApp">WhatsApp</button>
            <button type="button" class="px-3 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 text-sm transition-colors" @click="closeReceipt">Selesai</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useTenantStore } from '../stores/tenant'
import { getProducts, getProductByBarcode, getCategories } from '../api/products'
import { getProductStock, getLowStock } from '../api/inventory'
import { checkout as checkoutApi } from '../api/pos'
import { getCurrentShift } from '../api/shifts'
import { generateReceiptPDF, buildReceiptText } from '../utils/receipt-pdf'
import Receipt from '../components/Receipt.vue'
import ShiftModal from '../components/ShiftModal.vue'

const tenantStore = useTenantStore()
const auth = useAuthStore()
const isCashier = computed(() => auth.role === 'cashier')
const currentShift = ref(null)
const shiftModalMode = ref(null)
const lowStockList = ref([])
const toast = ref({ show: false, message: '', type: 'error' })
let toastTimer = null

function showToast(message, type = 'error') {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { show: true, message, type }
  toastTimer = setTimeout(() => {
    toast.value = { ...toast.value, show: false }
    toastTimer = null
  }, 3000)
}
const canCheckout = computed(() => !isCashier.value || !!currentShift.value)

const paymentMethods = [
  { value: 'cash',     label: 'Cash',     class: 'bg-emerald-600 hover:bg-emerald-700 text-white',  icon: '<rect x="2" y="6" width="20" height="12" rx="2"/><line x1="6" y1="12" x2="6.01" y2="12" stroke-width="2"/><circle cx="12" cy="12" r="2"/><line x1="18" y1="12" x2="18.01" y2="12" stroke-width="2"/>' },
  { value: 'card',     label: 'Card',     class: 'bg-blue-600 hover:bg-blue-700 text-white',        icon: '<rect x="1" y="4" width="22" height="16" rx="2"/><line x1="1" y1="10" x2="23" y2="10"/>' },
  { value: 'qris',     label: 'QRIS',     class: 'bg-violet-600 hover:bg-violet-700 text-white',    icon: '<rect x="2" y="2" width="6" height="6" rx="1"/><rect x="16" y="2" width="6" height="6" rx="1"/><rect x="2" y="16" width="6" height="6" rx="1"/><rect x="12" y="12" width="4" height="4" rx="0.5"/><rect x="18" y="18" width="4" height="4" rx="0.5"/><rect x="12" y="18" width="4" height="4" rx="0.5"/><rect x="18" y="12" width="4" height="4" rx="0.5"/>' },
  { value: 'ewallet',  label: 'E-Wallet', class: 'bg-amber-500 hover:bg-amber-600 text-white',      icon: '<path d="M21 12V7a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h7"/><path d="M13 17l3 3 5-5"/>' },
  { value: 'transfer', label: 'Transfer', class: 'bg-slate-600 hover:bg-slate-700 text-white',      icon: '<path d="M3 7h14l-4-4"/><path d="M21 17H7l4 4"/><line x1="3" y1="12" x2="21" y2="12" stroke-dasharray="2 2"/>' },
]

const SCAN_DEBOUNCE_MS = 200
const MIN_BARCODE_LENGTH = 8

const products = ref([])
const posCategories = ref([])
const posCategoryFilter = ref('')
const clickedProductId = ref(null)
const barcodeInput = ref('')
const barcodeInputRef = ref(null)
const filteredProducts = ref([])
const productsLoading = ref(true)
const productsError = ref(null)
const searchQuery = ref('')
const cartItems = ref([])
const showCheckoutModal = ref(false)
const checkoutForm = ref({ payment_method: 'cash', paid_amount: 0, customer_name: '', customer_phone: '' })
const checkoutError = ref('')
const checkoutSubmitting = ref(false)
const showReceiptModal = ref(false)
const receiptData = ref(null)
const mobileCartExpanded = ref(false)
let scanDebounceTimer = null

function productId(p) { return p?.id ?? p?.ID ?? '' }
function productName(p) { return p?.name ?? p?.Name ?? '' }
function productSku(p) { return p?.sku ?? p?.SKU ?? '' }
function productPrice(p) {
  const n = p?.sell_price ?? p?.SellPrice ?? p?.price ?? p?.Price
  return typeof n === 'number' ? n : Number(n) || 0
}

function formatPrice(value) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value ?? 0)
}

function posCatId(c) { return c?.id ?? c?.ID ?? '' }
function posCatName(c) { return c?.name ?? c?.Name ?? '' }

function filterProducts() {
  let list = [...(products.value || [])]
  const q = (searchQuery.value || '').trim().toLowerCase()
  if (q) {
    list = list.filter((p) => {
      const name = (productName(p) || '').toLowerCase()
      const sku = (productSku(p) || '').toLowerCase()
      return name.includes(q) || sku.includes(q)
    })
  }
  const catId = posCategoryFilter.value
  if (catId) {
    list = list.filter((p) => (p?.category_id ?? p?.CategoryID ?? '') === catId)
  }
  filteredProducts.value = list
}

async function onProductClick(p) {
  const ok = await addToCartWithStockCheck(p)
  if (ok) {
    clickedProductId.value = productId(p)
    setTimeout(() => { clickedProductId.value = null }, 400)
  }
}

watch([searchQuery, posCategoryFilter], filterProducts)

const totalAmount = computed(() => cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0))

function webAudioBeep(freq = 800, duration = 0.1) {
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain); gain.connect(ctx.destination)
    osc.frequency.value = freq; osc.type = 'sine'
    gain.gain.setValueAtTime(0.15, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + duration)
    osc.start(ctx.currentTime); osc.stop(ctx.currentTime + duration)
  } catch (_) {}
}

function playBeep() { new Audio('/beep.mp3').play().catch(() => webAudioBeep(800)) }
function playErrorBeep() { new Audio('/beep-error.mp3').play().catch(() => webAudioBeep(400, 0.2)) }

function addToCart(p) {
  const id = productId(p)
  const existing = cartItems.value.find((i) => i.product_id === id)
  if (existing) { existing.quantity += 1 }
  else { cartItems.value.push({ product_id: id, name: productName(p), price: productPrice(p), quantity: 1 }) }
  playBeep()
}

async function addToCartWithStockCheck(p) {
  const id = productId(p)
  try {
    const res = await getProductStock(id)
    const stock = res?.stock ?? 0
    if (stock <= 0) { showToast('Stok habis', 'error'); playErrorBeep(); return false }
    const existing = cartItems.value.find((i) => i.product_id === id)
    if ((existing ? existing.quantity + 1 : 1) > stock) { showToast('Stok tidak cukup', 'error'); playErrorBeep(); return false }
    addToCart(p)
    if (lowStockList.value.some((l) => l.product_id === id)) showToast(`⚠ Sisa stok: ${stock}`, 'warning')
    return true
  } catch (err) { showToast(err.response?.data?.error ?? 'Gagal cek stok', 'error'); playErrorBeep(); return false }
}

async function lookupAndAddBarcode(barcode) {
  const code = (barcode || '').trim()
  if (!code || code.length < MIN_BARCODE_LENGTH) return
  try {
    const product = await getProductByBarcode(code)
    if (await addToCartWithStockCheck(product)) { barcodeInput.value = ''; nextTick(() => barcodeInputRef.value?.focus()) }
  } catch (err) {
    showToast(err.response?.status === 404 ? 'Produk tidak ditemukan' : (err.response?.data?.error ?? 'Gagal mencari produk'), 'error')
    playErrorBeep()
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
  scanDebounceTimer = setTimeout(() => { scanDebounceTimer = null; lookupAndAddBarcode(barcode) }, SCAN_DEBOUNCE_MS)
})

function changeQuantity(productId, delta) {
  const item = cartItems.value.find((i) => i.product_id === productId)
  if (!item) return
  item.quantity = Math.max(0, item.quantity + delta)
  if (item.quantity <= 0) cartItems.value = cartItems.value.filter((i) => i.product_id !== productId)
  else playBeep()
}

function removeFromCart(productId) { cartItems.value = cartItems.value.filter((i) => i.product_id !== productId) }

function openCheckoutModal(method) {
  if (!canCheckout.value) return
  checkoutForm.value = { payment_method: method, paid_amount: totalAmount.value, customer_name: '', customer_phone: '' }
  checkoutError.value = ''
  showCheckoutModal.value = true
}

async function refreshCurrentShift() {
  if (!isCashier.value) return
  try { currentShift.value = (await getCurrentShift())?.shift ?? null } catch { currentShift.value = null }
}

function onShiftModalDone() { refreshCurrentShift(); shiftModalMode.value = null }
function onShiftModalClose() { shiftModalMode.value = null; refreshCurrentShift() }

async function submitCheckout() {
  checkoutError.value = ''; checkoutSubmitting.value = true
  try {
    const result = await checkoutApi({
      items: cartItems.value.map((i) => ({ product_id: i.product_id, quantity: i.quantity })),
      payment_method: checkoutForm.value.payment_method,
      paid_amount: checkoutForm.value.paid_amount,
      customer_name: checkoutForm.value.customer_name?.trim() || '',
      customer_phone: checkoutForm.value.customer_phone?.trim() || '',
    })
    showCheckoutModal.value = false
    receiptData.value = {
      storeName: tenantStore.storeName(), date: new Date(),
      transactionId: result?.transaction_id ?? '', cashier: result?.cashier ?? '',
      items: cartItems.value.map((i) => ({ ...i })),
      total: result?.total ?? totalAmount.value,
      paidAmount: checkoutForm.value.paid_amount, change: result?.change ?? 0,
      customerName: checkoutForm.value.customer_name?.trim() || '',
      customerPhone: checkoutForm.value.customer_phone?.trim() || '',
    }
    showReceiptModal.value = true
  } catch (err) { checkoutError.value = err.response?.data?.error ?? 'Checkout gagal.' }
  finally { checkoutSubmitting.value = false }
}

function printReceipt() { window.print() }
function downloadReceiptPDF() {
  if (!receiptData.value) return
  generateReceiptPDF(receiptData.value).save(`receipt-${receiptData.value.transactionId || 'pos'}.pdf`)
}
function shareWhatsApp() {
  if (!receiptData.value) return
  const text = encodeURIComponent(buildReceiptText(receiptData.value))
  const phone = (receiptData.value.customerPhone || '').trim().replace(/\D/g, '')
  const waPhone = phone ? (phone.startsWith('62') ? phone : '62' + phone.replace(/^0+/, '')) : ''
  const url = waPhone
    ? `https://wa.me/${waPhone}?text=${text}`
    : `https://wa.me/?text=${text}`
  window.open(url, '_blank')
}
function closeReceipt() { showReceiptModal.value = false; receiptData.value = null; cartItems.value = []; nextTick(() => barcodeInputRef.value?.focus()) }

function confirmClearCart() {
  if (!cartItems.value.length) return
  if (window.confirm('Kosongkan keranjang?')) { cartItems.value = [] }
}

const shortcutMap = { F1: 'cash', F2: 'card', F3: 'qris', F4: 'ewallet', F5: 'transfer' }
function handleKeyShortcuts(e) {
  if (showCheckoutModal.value || showReceiptModal.value) return
  const method = shortcutMap[e.key]
  if (method) { e.preventDefault(); if (cartItems.value.length && canCheckout.value) openCheckoutModal(method) }
  else if (e.key === 'Escape') confirmClearCart()
}

watch(() => auth.role, async (role) => {
  if (role !== 'cashier') return
  try {
    const shift = (await getCurrentShift())?.shift ?? null
    currentShift.value = shift
    if (shift) shiftModalMode.value = null
  } catch { currentShift.value = null }
}, { immediate: true })

onMounted(async () => {
  window.addEventListener('keydown', handleKeyShortcuts)
  productsLoading.value = true; productsError.value = null
  try {
    const [productsData, lowData, catsData] = await Promise.all([getProducts(), getLowStock().catch(() => []), getCategories().catch(() => [])])
    products.value = Array.isArray(productsData) ? productsData : []
    posCategories.value = Array.isArray(catsData) ? catsData : []
    filterProducts()
    lowStockList.value = Array.isArray(lowData) ? lowData : []
    if (isCashier.value) {
      const shift = (await getCurrentShift().catch(() => null))?.shift ?? null
      currentShift.value = shift
      if (shift) shiftModalMode.value = null
    }
  } catch { productsError.value = 'Gagal memuat produk' }
  finally { productsLoading.value = false }
  await nextTick(); barcodeInputRef.value?.focus()
})

onUnmounted(() => { window.removeEventListener('keydown', handleKeyShortcuts) })
</script>

<style scoped>
@keyframes pos-product-click {
  0% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0.6); }
  50% { box-shadow: 0 0 0 8px rgba(99, 102, 241, 0); }
  100% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0); }
}
.animate-pos-product-click {
  animation: pos-product-click 0.4s ease-out;
}
</style>
