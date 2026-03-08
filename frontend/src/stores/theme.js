import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const dark = ref(localStorage.getItem('theme') === 'dark')

  function apply() {
    document.documentElement.classList.toggle('dark', dark.value)
  }

  function toggle() {
    dark.value = !dark.value
  }

  watch(dark, (v) => {
    localStorage.setItem('theme', v ? 'dark' : 'light')
    apply()
  })

  apply()

  return { dark, toggle }
})
