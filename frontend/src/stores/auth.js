import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    user: null,
    tenant_id: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    login(payload) {
      // Payload from API: { token, user?, tenant_id? }. Call after successful auth API response.
      this.token = payload.token
      this.user = payload.user ?? null
      this.tenant_id = payload.tenant_id ?? null
    },
    logout() {
      this.token = null
      this.user = null
      this.tenant_id = null
    },
  },
})
