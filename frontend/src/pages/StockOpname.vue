<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-white">Stock Opname</h1>
      <button @click="startNew" :disabled="starting"
        class="px-4 py-2 bg-slate-700 text-white rounded-lg hover:bg-slate-600 text-sm font-medium disabled:opacity-50 transition-colors">
        {{ starting ? 'Creating...' : '+ Opname Baru' }}
      </button>
    </div>

    <div v-if="loading" class="text-gray-500 dark:text-gray-400 py-8 text-center">Loading...</div>

    <!-- Opname List -->
    <div v-else class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-xs uppercase">
          <tr>
            <th class="px-4 py-3 text-left">Tanggal</th>
            <th class="px-4 py-3 text-center">Status</th>
            <th class="px-4 py-3 text-center">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-800">
          <tr v-for="op in opnames" :key="op.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
            <td class="px-4 py-3 text-gray-800 dark:text-gray-200">{{ formatDate(op.created_at) }}</td>
            <td class="px-4 py-3 text-center">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="statusClass(op.status)">
                {{ statusLabel(op.status) }}
              </span>
            </td>
            <td class="px-4 py-3 text-center space-x-3">
              <button @click="openDetail(op)" class="text-indigo-600 hover:text-indigo-800 text-xs font-medium">
                {{ op.status === 'draft' ? 'Input Stok' : op.status === 'submitted' ? 'Review & Approve' : 'Lihat Detail' }}
              </button>
              <button v-if="op.status === 'draft'" @click="confirmDelete(op)" class="text-red-500 hover:text-red-700 text-xs font-medium">
                Hapus
              </button>
            </td>
          </tr>
          <tr v-if="!opnames.length">
            <td colspan="3" class="px-4 py-10 text-center text-gray-400 dark:text-gray-500">
              Belum ada data stock opname. Klik "+ Opname Baru" untuk memulai.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Detail / Edit Modal -->
    <Teleport to="body">
      <div v-if="modal" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4" @click.self="closeModal">
        <div class="bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col border border-gray-200 dark:border-gray-800">
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Stock Opname</h2>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                {{ formatDate(modal.opname.created_at) }} &mdash;
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="statusClass(modal.opname.status)">
                  {{ statusLabel(modal.opname.status) }}
                </span>
              </p>
            </div>
            <button @click="closeModal" class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 text-2xl leading-none">&times;</button>
          </div>

          <!-- Search (draft mode) -->
          <div v-if="modal.opname.status === 'draft'" class="px-6 pt-4">
            <input v-model="searchQuery" type="text" placeholder="Cari produk..."
              class="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" />
          </div>

          <!-- Table -->
          <div class="flex-1 overflow-auto px-6 py-4">
            <div v-if="loadingDetail" class="py-8 text-center text-gray-400 dark:text-gray-500">Memuat produk...</div>
            <table v-else class="w-full text-sm">
              <thead class="bg-gray-50 dark:bg-gray-800 sticky top-0">
                <tr>
                  <th class="px-3 py-2.5 text-left text-gray-600 dark:text-gray-400 text-xs uppercase">Produk</th>
                  <th class="px-3 py-2.5 text-right text-gray-600 dark:text-gray-400 text-xs uppercase w-24">Stok Sistem</th>
                  <th class="px-3 py-2.5 text-right text-gray-600 dark:text-gray-400 text-xs uppercase w-28">
                    {{ modal.opname.status === 'draft' ? 'Stok Aktual' : 'Aktual' }}
                  </th>
                  <th class="px-3 py-2.5 text-right text-gray-600 dark:text-gray-400 text-xs uppercase w-20">Selisih</th>
                  <th class="px-3 py-2.5 text-left text-gray-600 dark:text-gray-400 text-xs uppercase w-40">Catatan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-gray-800">
                <tr v-for="item in filteredItems" :key="item.product_id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
                  <td class="px-3 py-2 text-gray-800 dark:text-gray-200 font-medium">{{ item.product_name }}</td>
                  <td class="px-3 py-2 text-right text-gray-600 dark:text-gray-400">{{ item.system_stock }}</td>
                  <td class="px-3 py-2 text-right">
                    <input v-if="modal.opname.status === 'draft'"
                      v-model.number="item.actual_stock" type="number" min="0"
                      class="w-20 border border-gray-300 dark:border-gray-600 rounded px-2 py-1 text-sm text-right bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
                      @input="updateDiff(item)" />
                    <span v-else>{{ item.actual_stock }}</span>
                  </td>
                  <td class="px-3 py-2 text-right font-medium"
                    :class="item.difference < 0 ? 'text-red-600' : item.difference > 0 ? 'text-blue-600' : 'text-gray-400'">
                    {{ item.difference > 0 ? '+' : '' }}{{ item.difference }}
                  </td>
                  <td class="px-3 py-2">
                    <input v-if="modal.opname.status === 'draft'"
                      v-model="item.notes" type="text" placeholder="—"
                      class="w-full border border-gray-300 dark:border-gray-600 rounded px-2 py-1 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none" />
                    <span v-else class="text-gray-500 text-xs">{{ item.notes || '—' }}</span>
                  </td>
                </tr>
<tr v-if="!filteredItems.length">
                <td colspan="5" class="px-3 py-6 text-center text-gray-400 dark:text-gray-500">Tidak ada produk ditemukan</td>
              </tr>
              </tbody>
            </table>
          </div>

          <!-- Summary & Actions -->
          <div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 rounded-b-xl">
            <div class="flex items-center justify-between">
              <div class="text-xs text-gray-500 dark:text-gray-400 space-x-4">
                <span>Total produk: <strong>{{ modalItems.length }}</strong></span>
                <span>Selisih kurang: <strong class="text-red-600">{{ shortageCount }}</strong></span>
                <span>Selisih lebih: <strong class="text-blue-600">{{ surplusCount }}</strong></span>
                <span>Cocok: <strong class="text-green-600">{{ matchCount }}</strong></span>
              </div>
              <div class="flex gap-2">
                <button @click="closeModal"
                  class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
                  Tutup
                </button>
                <button v-if="modal.opname.status === 'draft'" @click="deleteFromModal"
                  class="px-4 py-2 border border-red-300 dark:border-red-700 text-red-600 dark:text-red-400 rounded-lg text-sm font-medium hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                  Hapus Draft
                </button>
                <button v-if="modal.opname.status === 'draft'" @click="submitItems" :disabled="submitting"
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition-colors">
                  {{ submitting ? 'Menyimpan...' : 'Submit Opname' }}
                </button>
                <button v-if="modal.opname.status === 'submitted'" @click="approveOpnameAction" :disabled="approving"
                  class="px-4 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 disabled:opacity-50 transition-colors">
                  {{ approving ? 'Memproses...' : 'Approve & Sesuaikan Stok' }}
                </button>
              </div>
            </div>
            <p v-if="actionError" class="text-xs text-red-600 mt-2">{{ actionError }}</p>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { listOpnames, startOpname, getOpname, submitOpnameItems, approveOpname, deleteOpname } from '../api/opname'
import { getInventory } from '../api/inventory'
import { formatDate } from '../utils'

const loading = ref(true)
const opnames = ref([])
const starting = ref(false)
const modal = ref(null)
const modalItems = ref([])
const loadingDetail = ref(false)
const submitting = ref(false)
const approving = ref(false)
const actionError = ref('')
const searchQuery = ref('')

const filteredItems = computed(() => {
  if (!searchQuery.value) return modalItems.value
  const q = searchQuery.value.toLowerCase()
  return modalItems.value.filter(it => it.product_name.toLowerCase().includes(q))
})

const shortageCount = computed(() => modalItems.value.filter(it => it.difference < 0).length)
const surplusCount = computed(() => modalItems.value.filter(it => it.difference > 0).length)
const matchCount = computed(() => modalItems.value.filter(it => it.difference === 0).length)

function statusClass(s) {
  return {
    draft: 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300',
    submitted: 'bg-blue-100 text-blue-700',
    completed: 'bg-green-100 text-green-700',
  }[s] || 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300'
}

function statusLabel(s) {
  return { draft: 'Draft', submitted: 'Submitted', completed: 'Completed' }[s] || s
}

function updateDiff(item) {
  item.difference = (item.actual_stock ?? 0) - item.system_stock
}

async function load() {
  loading.value = true
  try {
    const data = await listOpnames()
    opnames.value = Array.isArray(data) ? data : []
  } catch { opnames.value = [] }
  finally { loading.value = false }
}

async function startNew() {
  starting.value = true
  try {
    const op = await startOpname()
    await load()
    await openDetail(op)
  } finally { starting.value = false }
}

async function openDetail(op) {
  loadingDetail.value = true
  actionError.value = ''
  searchQuery.value = ''
  modal.value = { opname: op }
  modalItems.value = []

  try {
    if (op.status === 'draft') {
      const inventory = await getInventory()
      modalItems.value = (Array.isArray(inventory) ? inventory : []).map(inv => ({
        product_id: inv.product_id,
        product_name: inv.product_name,
        system_stock: inv.stock,
        actual_stock: inv.stock,
        difference: 0,
        notes: '',
      }))
    } else {
      const detail = await getOpname(op.id)
      const items = detail.items || []
      modalItems.value = items.map(it => ({
        product_id: it.product_id,
        product_name: it.product_name || it.product_id,
        system_stock: it.system_stock,
        actual_stock: it.actual_stock,
        difference: it.difference,
        notes: it.notes || '',
      }))
      if (detail.opname) {
        modal.value = { opname: detail.opname }
      }
    }
  } catch (e) {
    actionError.value = 'Gagal memuat data: ' + (e.response?.data?.error || e.message)
  } finally {
    loadingDetail.value = false
  }
}

async function submitItems() {
  submitting.value = true
  actionError.value = ''
  try {
    const changedItems = modalItems.value.filter(it => it.difference !== 0)
    if (!changedItems.length) {
      actionError.value = 'Tidak ada perubahan stok. Ubah kolom "Stok Aktual" pada produk yang selisih.'
      submitting.value = false
      return
    }
    const payload = changedItems.map(it => ({
      product_id: it.product_id,
      actual_stock: it.actual_stock,
      notes: it.notes,
    }))
    await submitOpnameItems(modal.value.opname.id, payload)
    modal.value.opname.status = 'submitted'
    await load()
    await openDetail(modal.value.opname)
  } catch (e) {
    actionError.value = e.response?.data?.error || 'Gagal submit opname'
  } finally {
    submitting.value = false
  }
}

async function approveOpnameAction() {
  approving.value = true
  actionError.value = ''
  try {
    await approveOpname(modal.value.opname.id)
    modal.value = null
    modalItems.value = []
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.error || 'Gagal approve opname'
  } finally {
    approving.value = false
  }
}

async function deleteFromModal() {
  if (!modal.value || !confirm('Hapus opname draft ini?')) return
  try {
    await deleteOpname(modal.value.opname.id)
    closeModal()
    await load()
  } catch (e) {
    actionError.value = e.response?.data?.error || 'Gagal menghapus opname'
  }
}

async function confirmDelete(op) {
  if (!confirm('Hapus opname draft ini?')) return
  try {
    await deleteOpname(op.id)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal menghapus opname')
  }
}

function closeModal() {
  modal.value = null
  modalItems.value = []
  searchQuery.value = ''
  actionError.value = ''
}

onMounted(load)
</script>
