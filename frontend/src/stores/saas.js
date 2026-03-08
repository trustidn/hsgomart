import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getSaasInfo } from '../api/admin'

const baseURL =
  typeof window !== 'undefined' && window.location.port === '8080'
    ? ''
    : 'http://localhost:8080'

export const useSaasStore = defineStore('saas', () => {
  const info = ref(null)
  const loaded = ref(false)

  async function load(force = false) {
    if (loaded.value && !force) return info.value
    try {
      info.value = await getSaasInfo()
    } catch {
      info.value = null
    }
    loaded.value = true
    return info.value
  }

  function applyBranding() {
    const name = platformName.value
    const tag = tagline.value
    document.title = tag ? `${name} — ${tag}` : name

    const logo = logoUrl.value
    if (logo) {
      const favicon = document.getElementById('app-favicon')
      if (favicon) {
        favicon.href = logo.startsWith('http') ? logo : `${baseURL}${logo}`
        favicon.type = 'image/png'
      }
    }
  }

  const platformName = computed(() => info.value?.saas_name || 'HSMart POS')
  const tagline = computed(() => info.value?.tagline || '')
  const logoUrl = computed(() => info.value?.logo_url || '')
  const logoSrc = computed(() => {
    if (!logoUrl.value) return ''
    return logoUrl.value.startsWith('http') ? logoUrl.value : `${baseURL}${logoUrl.value}`
  })
  const whatsappNumber = computed(() => info.value?.whatsapp_number || '')

  return { info, loaded, load, applyBranding, platformName, tagline, logoUrl, logoSrc, whatsappNumber }
})
