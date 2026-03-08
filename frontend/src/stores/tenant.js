import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getTenantProfile } from '../api/tenant'
import { useSaasStore } from './saas'

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
    const saas = useSaasStore()

    const name = profile.value?.name || saas.platformName
    if (name) {
      document.title = name
    }

    const logo = profile.value?.logo_url || saas.logoUrl
    if (logo) {
      const favicon = document.getElementById('app-favicon')
      if (favicon) {
        const src = logo.startsWith('http') ? logo : `${baseURL}${logo}`
        favicon.href = src
        favicon.type = 'image/png'
      }
    }
  }

  function reset() {
    profile.value = null
    loaded.value = false
  }

  const storeName = () => {
    if (profile.value?.name) return profile.value.name
    const saas = useSaasStore()
    return saas.platformName || 'HSMart'
  }

  const logoUrl = () => {
    if (profile.value?.logo_url) return profile.value.logo_url
    const saas = useSaasStore()
    return saas.logoUrl || ''
  }

  return { profile, loaded, load, reset, storeName, logoUrl, applyBranding }
})
