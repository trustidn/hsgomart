import { defineStore } from 'pinia'

const TOKEN_KEY = 'hsmart_token'
const USER_KEY = 'hsmart_user'
const TENANT_ID_KEY = 'hsmart_tenant_id'

function loadFromStorage(key, parse = false) {
  try {
    const v = localStorage.getItem(key)
    if (v == null) return null
    return parse ? JSON.parse(v) : v
  } catch {
    return null
  }
}

function saveToStorage(key, value) {
  if (value == null) localStorage.removeItem(key)
  else localStorage.setItem(key, typeof value === 'string' ? value : JSON.stringify(value))
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: loadFromStorage(TOKEN_KEY),
    user: loadFromStorage(USER_KEY, true),
    tenant_id: loadFromStorage(TENANT_ID_KEY),
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    login(payload) {
      this.token = payload.token ?? this.token
      this.user = payload.user ?? null
      this.tenant_id = payload.tenant_id ?? null
      saveToStorage(TOKEN_KEY, this.token)
      saveToStorage(USER_KEY, this.user)
      saveToStorage(TENANT_ID_KEY, this.tenant_id)
    },
    logout() {
      this.token = null
      this.user = null
      this.tenant_id = null
      saveToStorage(TOKEN_KEY, null)
      saveToStorage(USER_KEY, null)
      saveToStorage(TENANT_ID_KEY, null)
    },
  },
})
