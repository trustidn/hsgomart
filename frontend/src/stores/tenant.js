import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getTenantProfile } from '../api/tenant'

const baseURL =
  typeof window !== 'undefined' && window.location.port === '8080'
    ? ''
    : 'http://localhost:8080'

export const useTenantStore = defineStore('tenant', () => {
  const profile = ref(null)
  const loaded = ref(false)

  async function load() {
    if (loaded.value) return profile.value
    try {
      profile.value = await getTenantProfile()
      applyBranding()
    } catch {
      profile.value = null
    }
    loaded.value = true
    return profile.value
  }

  function applyBranding() {
    if (!profile.value) return
    const name = profile.value.name
    if (name) {
      document.title = name
    }
    const logo = profile.value.logo_url
    if (logo) {
      const favicon = document.getElementById('app-favicon')
      if (favicon) {
        favicon.href = `${baseURL}${logo}`
        favicon.type = 'image/png'
      }
    }
  }

  function reset() {
    profile.value = null
    loaded.value = false
  }

  const storeName = () => profile.value?.name || 'HSMart'
  const logoUrl = () => profile.value?.logo_url || ''

  return { profile, loaded, load, reset, storeName, logoUrl, applyBranding }
})
