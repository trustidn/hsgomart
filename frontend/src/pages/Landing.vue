<template>
  <div class="min-h-screen bg-white dark:bg-gray-950 text-gray-900 dark:text-gray-100 transition-colors">
    <!-- Header -->
    <header class="sticky top-0 z-50 bg-white/80 dark:bg-gray-950/80 backdrop-blur-lg border-b border-gray-100 dark:border-gray-800">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <img v-if="saas.logoSrc" :src="saas.logoSrc" class="w-9 h-9 rounded-xl object-cover" alt="" />
          <div v-else class="w-9 h-9 rounded-xl bg-indigo-500 flex items-center justify-center text-white text-lg font-bold">{{ saas.platformName.charAt(0) }}</div>
          <span class="text-lg font-bold tracking-tight">{{ saas.platformName }}</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="themeStore.toggle()"
            class="p-2 rounded-lg text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            :title="themeStore.dark ? 'Light mode' : 'Dark mode'"
          >
            <svg v-if="themeStore.dark" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="5" stroke-width="2"/><path stroke-linecap="round" stroke-width="2" d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/></svg>
          </button>
          <template v-if="auth.isAuthenticated">
            <router-link :to="dashboardPath" class="px-5 py-2 text-sm font-medium rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition-colors">
              Dashboard
            </router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">Masuk</router-link>
            <router-link to="/register" class="px-5 py-2 text-sm font-medium rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition-colors">Daftar Gratis</router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- Hero -->
    <section class="relative overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-br from-indigo-50 via-white to-emerald-50 dark:from-indigo-950/30 dark:via-gray-950 dark:to-emerald-950/20" />
      <div class="relative max-w-6xl mx-auto px-4 sm:px-6 pt-20 pb-24 sm:pt-28 sm:pb-32">
        <div class="max-w-3xl mx-auto text-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-100 dark:bg-indigo-900/40 text-indigo-700 dark:text-indigo-300 text-xs font-medium mb-6">
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>
            Solusi POS Terjangkau untuk UMKM
          </div>
          <h1 class="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight leading-tight">
            Kelola Usaha Anda <br class="hidden sm:block" />
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-emerald-500 dark:from-indigo-400 dark:to-emerald-400">Tanpa Ribet</span>
          </h1>
          <p class="mt-6 text-lg sm:text-xl text-gray-600 dark:text-gray-400 leading-relaxed max-w-2xl mx-auto">
            Sistem kasir &amp; manajemen toko lengkap yang bisa diakses dari <strong class="text-gray-800 dark:text-gray-200">HP, tablet, atau laptop</strong>. 
            Tanpa perlu beli perangkat khusus. Harga mulai dari <strong class="text-gray-800 dark:text-gray-200">Rp 0</strong> untuk masa uji coba.
          </p>
          <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-3">
            <router-link v-if="auth.isAuthenticated" :to="dashboardPath" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-500/25 transition-all hover:shadow-indigo-500/40">
              Buka Dashboard
            </router-link>
            <template v-else>
              <router-link to="/register" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-500/25 transition-all hover:shadow-indigo-500/40">
                Coba Gratis 7 Hari
              </router-link>
              <router-link to="/login" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
                Sudah Punya Akun? Masuk
              </router-link>
            </template>
          </div>
        </div>
      </div>
    </section>

    <!-- Social proof strip -->
    <section class="border-y border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/50">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 py-8 flex flex-wrap items-center justify-center gap-x-10 gap-y-4 text-center">
        <div v-for="stat in stats" :key="stat.label">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ stat.value }}</div>
          <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ stat.label }}</div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="py-20 sm:py-28">
      <div class="max-w-6xl mx-auto px-4 sm:px-6">
        <div class="text-center max-w-2xl mx-auto mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold tracking-tight">Semua yang Anda Butuhkan, <span class="text-indigo-600 dark:text-indigo-400">Satu Platform</span></h2>
          <p class="mt-4 text-gray-600 dark:text-gray-400 text-lg">Fitur lengkap setara sistem POS mahal, dengan harga yang ramah di kantong UMKM.</p>
        </div>
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="f in features" :key="f.title"
            class="group relative bg-white dark:bg-gray-900 border border-gray-100 dark:border-gray-800 rounded-2xl p-6 hover:shadow-lg hover:border-indigo-200 dark:hover:border-indigo-800 transition-all">
            <div class="w-11 h-11 rounded-xl flex items-center justify-center text-lg mb-4"
              :class="f.bgClass">
              <span v-html="f.icon" />
            </div>
            <h3 class="font-semibold text-lg mb-2 text-gray-900 dark:text-white">{{ f.title }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Why us -->
    <section class="py-20 sm:py-28 bg-gray-50 dark:bg-gray-900/50 border-y border-gray-100 dark:border-gray-800">
      <div class="max-w-6xl mx-auto px-4 sm:px-6">
        <div class="text-center max-w-2xl mx-auto mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold tracking-tight">Kenapa Memilih <span class="text-indigo-600 dark:text-indigo-400">{{ saas.platformName }}</span>?</h2>
        </div>
        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-8">
          <div v-for="r in reasons" :key="r.title" class="text-center">
            <div class="w-14 h-14 rounded-2xl mx-auto flex items-center justify-center text-2xl mb-4" :class="r.bgClass">{{ r.emoji }}</div>
            <h3 class="font-semibold text-base mb-1 text-gray-900 dark:text-white">{{ r.title }}</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ r.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="py-20 sm:py-28">
      <div class="max-w-3xl mx-auto px-4 sm:px-6 text-center">
        <h2 class="text-3xl sm:text-4xl font-bold tracking-tight mb-4">Siap Digitalkan Usaha Anda?</h2>
        <p class="text-lg text-gray-600 dark:text-gray-400 mb-10 max-w-xl mx-auto">Daftarkan toko Anda sekarang dan nikmati semua fitur secara gratis selama masa uji coba. Tanpa kartu kredit, tanpa komitmen.</p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
          <router-link v-if="auth.isAuthenticated" :to="dashboardPath" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-500/25 transition-all">
            Buka Dashboard
          </router-link>
          <template v-else>
            <router-link to="/register" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-500/25 transition-all">
              Mulai Gratis Sekarang
            </router-link>
            <router-link to="/login" class="w-full sm:w-auto px-8 py-3.5 text-base font-semibold rounded-xl border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
              Masuk
            </router-link>
          </template>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-gray-100 dark:border-gray-800 py-8">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-gray-500 dark:text-gray-400">
        <div class="flex items-center gap-2">
          <img v-if="saas.logoSrc" :src="saas.logoSrc" class="w-5 h-5 rounded object-cover" alt="" />
          <span>&copy; {{ new Date().getFullYear() }} {{ saas.platformName }}. All rights reserved.</span>
        </div>
        <div class="flex items-center gap-4">
          <router-link to="/login" class="hover:text-gray-900 dark:hover:text-white transition-colors">Masuk</router-link>
          <router-link to="/register" class="hover:text-gray-900 dark:hover:text-white transition-colors">Daftar</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSaasStore } from '../stores/saas'
import { useThemeStore } from '../stores/theme'

const auth = useAuthStore()
const saas = useSaasStore()
const themeStore = useThemeStore()

const dashboardPath = computed(() => auth.role === 'superadmin' ? '/admin/dashboard' : '/dashboard')

onMounted(async () => {
  await saas.load()
  saas.applyBranding()
})

const stats = [
  { value: '100%', label: 'Berbasis Web' },
  { value: '24/7', label: 'Akses Kapan Saja' },
  { value: '0', label: 'Perangkat Khusus' },
  { value: '< 5 menit', label: 'Setup Toko' },
]

const features = [
  {
    icon: '&#128179;',
    title: 'Kasir / POS',
    desc: 'Proses transaksi cepat dengan barcode scanner, metode pembayaran fleksibel, dan struk digital via WhatsApp.',
    bgClass: 'bg-indigo-100 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400',
  },
  {
    icon: '&#128230;',
    title: 'Manajemen Produk',
    desc: 'Kelola ribuan produk dengan kategori, barcode ganda, dan harga beli/jual. Import data dengan mudah.',
    bgClass: 'bg-emerald-100 dark:bg-emerald-900/40 text-emerald-600 dark:text-emerald-400',
  },
  {
    icon: '&#128200;',
    title: 'Laporan Lengkap',
    desc: 'Pantau penjualan harian, profit margin, performa kasir, dan tren produk dalam satu dashboard.',
    bgClass: 'bg-amber-100 dark:bg-amber-900/40 text-amber-600 dark:text-amber-400',
  },
  {
    icon: '&#128451;',
    title: 'Inventory & Stok',
    desc: 'Lacak stok real-time, peringatan stok rendah, stock opname, dan riwayat pergerakan barang.',
    bgClass: 'bg-sky-100 dark:bg-sky-900/40 text-sky-600 dark:text-sky-400',
  },
  {
    icon: '&#128101;',
    title: 'Multi-User & Shift',
    desc: 'Tambahkan kasir dengan hak akses berbeda. Kelola shift kasir dengan rekap transaksi per shift.',
    bgClass: 'bg-rose-100 dark:bg-rose-900/40 text-rose-600 dark:text-rose-400',
  },
  {
    icon: '&#127961;',
    title: 'White Label',
    desc: 'Tampilkan nama toko, logo, dan branding Anda sendiri. Struk dan invoice otomatis menggunakan identitas toko Anda.',
    bgClass: 'bg-violet-100 dark:bg-violet-900/40 text-violet-600 dark:text-violet-400',
  },
]

const reasons = [
  { emoji: '💰', title: 'Harga UMKM-Friendly', desc: 'Mulai dari gratis. Paket berbayar jauh lebih murah dari sistem POS konvensional.', bgClass: 'bg-emerald-100 dark:bg-emerald-900/40' },
  { emoji: '📱', title: 'Akses dari Mana Saja', desc: 'Cukup buka browser di HP, tablet, atau laptop. Tidak perlu install aplikasi.', bgClass: 'bg-sky-100 dark:bg-sky-900/40' },
  { emoji: '⚡', title: 'Setup Instan', desc: 'Daftar, dan langsung gunakan. Data contoh sudah tersedia untuk Anda pelajari.', bgClass: 'bg-amber-100 dark:bg-amber-900/40' },
  { emoji: '🔒', title: 'Aman & Terpercaya', desc: 'Data terenkripsi, backup otomatis, dan hak akses per pengguna untuk keamanan toko Anda.', bgClass: 'bg-rose-100 dark:bg-rose-900/40' },
]
</script>
