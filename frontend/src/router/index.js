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
      { path: 'dashboard', name: 'Dashboard', component: () => import('../pages/Dashboard.vue'), meta: { roles: ['owner', 'cashier'] } },
      { path: 'products', name: 'Products', component: () => import('../pages/Products.vue'), meta: { roles: ['owner'] } },
      { path: 'inventory', name: 'Inventory', component: () => import('../pages/Inventory.vue'), meta: { roles: ['owner'] } },
      { path: 'inventory-history', name: 'InventoryHistory', component: () => import('../pages/InventoryHistory.vue'), meta: { roles: ['owner'] } },
      { path: 'categories', name: 'Categories', component: () => import('../pages/Categories.vue'), meta: { roles: ['owner'] } },
      { path: 'purchases', name: 'Purchases', component: () => import('../pages/Purchases.vue'), meta: { roles: ['owner'] } },
      { path: 'purchases/:id', name: 'PurchaseDetail', component: () => import('../pages/Purchases.vue'), meta: { roles: ['owner'] } },
      { path: 'pos', name: 'POS', component: () => import('../pages/POS.vue'), meta: { roles: ['owner', 'cashier'] } },
      { path: 'reports', name: 'Reports', component: () => import('../pages/Reports.vue'), meta: { roles: ['owner'] } },
      { path: 'users', name: 'Users', component: () => import('../pages/Users.vue'), meta: { roles: ['owner'] } },
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
  const role = auth.role || 'cashier'
  if (to.meta.roles && to.meta.roles.length && !to.meta.roles.includes(role)) {
    return { name: 'Dashboard' }
  }
  return true
})

export default router
