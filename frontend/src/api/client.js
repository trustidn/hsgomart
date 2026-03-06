import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// When app is served from backend (e.g. :8080), use same origin. In dev (e.g. :5173), use backend URL.
const baseURL =
  typeof window !== 'undefined' && window.location.port === '8080'
    ? ''
    : 'http://localhost:8080'

const client = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

client.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

export default client
