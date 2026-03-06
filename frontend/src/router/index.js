import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('../pages/Dashboard.vue') },
      { path: 'products', name: 'Products', component: () => import('../pages/Products.vue') },
      { path: 'inventory', name: 'Inventory', component: () => import('../pages/Inventory.vue') },
      { path: 'inventory-history', name: 'InventoryHistory', component: () => import('../pages/InventoryHistory.vue') },
      { path: 'categories', name: 'Categories', component: () => import('../pages/Categories.vue') },
      { path: 'pos', name: 'POS', component: () => import('../pages/POS.vue') },
      { path: 'reports', name: 'Reports', component: () => import('../pages/Reports.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    if (to.name === 'Login' && auth.token) return { path: '/dashboard' }
    return true
  }
  if (to.meta.requiresAuth && !auth.token) return { name: 'Login' }
  return true
})

export default router
