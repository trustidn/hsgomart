import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import './style.css'
import App from './App.vue'

const app = createApp(App)
app.use(createPinia())
app.use(router)

app.config.errorHandler = (err, instance, info) => {
  console.error('Vue error:', err, info)
}

const el = document.getElementById('app')
if (!el) {
  document.body.innerHTML = '<p style="padding:1rem;font-family:sans-serif;">#app not found. Check index.html.</p>'
} else {
  try {
    app.mount('#app')
  } catch (err) {
    console.error('Mount failed:', err)
    el.innerHTML = '<p style="padding:1rem;font-family:sans-serif;color:#c00;">App failed to load. Open DevTools (F12) → Console for errors.</p>'
  }
}
