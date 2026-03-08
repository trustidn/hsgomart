import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const baseURL = import.meta.env.VITE_API_URL ||
  (typeof window !== 'undefined' && window.location.port === '8080'
    ? ''
    : 'http://localhost:8080')

const client = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

let isRefreshing = false
let failedQueue = []

function processQueue(error, token = null) {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error)
    else resolve(token)
  })
  failedQueue = []
}

client.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 402) {
      if (error.response?.data?.read_only) {
        return Promise.reject(error)
      }
      if (!window.location.pathname.startsWith('/subscription')) {
        window.location.href = '/subscription'
      }
      return Promise.reject(error)
    }
    const originalRequest = error.config
    if (error.response?.status === 401 && !originalRequest._retry) {
      const auth = useAuthStore()
      if (!auth.refreshToken || originalRequest.url === '/auth/refresh') {
        auth.logout()
        window.location.href = '/login'
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return client(originalRequest)
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const { data } = await axios.post(`${baseURL}/auth/refresh`, {
          refresh_token: auth.refreshToken,
        })
        auth.setTokens(data.token, data.refresh_token)
        processQueue(null, data.token)
        originalRequest.headers.Authorization = `Bearer ${data.token}`
        return client(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        auth.logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
    return Promise.reject(error)
  },
)

export default client
