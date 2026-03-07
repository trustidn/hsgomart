import { defineStore } from 'pinia'

const TOKEN_KEY = 'hsmart_token'
const REFRESH_TOKEN_KEY = 'hsmart_refresh_token'
const USER_KEY = 'hsmart_user'
const TENANT_ID_KEY = 'hsmart_tenant_id'
const ROLE_KEY = 'hsmart_role'

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
    refreshToken: loadFromStorage(REFRESH_TOKEN_KEY),
    user: loadFromStorage(USER_KEY, true),
    tenant_id: loadFromStorage(TENANT_ID_KEY),
    role: loadFromStorage(ROLE_KEY),
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    login(payload) {
      this.token = payload.token ?? this.token
      this.refreshToken = payload.refresh_token ?? this.refreshToken
      this.user = payload.user ?? null
      this.tenant_id = payload.tenant_id ?? null
      this.role = payload.role ?? payload.user?.role ?? null
      saveToStorage(TOKEN_KEY, this.token)
      saveToStorage(REFRESH_TOKEN_KEY, this.refreshToken)
      saveToStorage(USER_KEY, this.user)
      saveToStorage(TENANT_ID_KEY, this.tenant_id)
      saveToStorage(ROLE_KEY, this.role)
    },
    setTokens(token, refreshToken) {
      this.token = token
      this.refreshToken = refreshToken
      saveToStorage(TOKEN_KEY, token)
      saveToStorage(REFRESH_TOKEN_KEY, refreshToken)
    },
    logout() {
      this.token = null
      this.refreshToken = null
      this.user = null
      this.tenant_id = null
      this.role = null
      saveToStorage(TOKEN_KEY, null)
      saveToStorage(REFRESH_TOKEN_KEY, null)
      saveToStorage(USER_KEY, null)
      saveToStorage(TENANT_ID_KEY, null)
      saveToStorage(ROLE_KEY, null)
    },
  },
})
