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
    path: '/register',
    name: 'Register',
    component: () => import('../pages/Register.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('../pages/Dashboard.vue'), meta: { title: 'Dashboard', roles: ['owner', 'cashier'] } },
      { path: 'products', name: 'Products', component: () => import('../pages/Products.vue'), meta: { title: 'Products', roles: ['owner'] } },
      { path: 'inventory', name: 'Inventory', component: () => import('../pages/Inventory.vue'), meta: { title: 'Inventory', roles: ['owner'] } },
      { path: 'inventory-history', name: 'InventoryHistory', component: () => import('../pages/InventoryHistory.vue'), meta: { title: 'Inventory History', roles: ['owner'] } },
      { path: 'categories', name: 'Categories', component: () => import('../pages/Categories.vue'), meta: { title: 'Categories', roles: ['owner'] } },
      { path: 'purchases', name: 'Purchases', component: () => import('../pages/Purchases.vue'), meta: { title: 'Purchases', roles: ['owner'] } },
      { path: 'purchases/:id', name: 'PurchaseDetail', component: () => import('../pages/Purchases.vue'), meta: { title: 'Purchase Detail', roles: ['owner'] } },
      { path: 'pos', name: 'POS', component: () => import('../pages/POS.vue'), meta: { title: 'Point of Sale', roles: ['owner', 'cashier'] } },
      { path: 'reports', name: 'Reports', component: () => import('../pages/Reports.vue'), meta: { title: 'Reports', roles: ['owner'] } },
      { path: 'shifts', name: 'Shifts', component: () => import('../pages/Shifts.vue'), meta: { title: 'Shifts', roles: ['owner'] } },
      { path: 'users', name: 'Users', component: () => import('../pages/Users.vue'), meta: { title: 'Users', roles: ['owner'] } },
      { path: 'stock-opname', name: 'StockOpname', component: () => import('../pages/StockOpname.vue'), meta: { title: 'Stock Opname', roles: ['owner'] } },
      { path: 'subscription', name: 'Subscription', component: () => import('../pages/Subscription.vue'), meta: { title: 'Subscription', roles: ['owner'] } },
      { path: 'settings', name: 'Settings', component: () => import('../pages/Settings.vue'), meta: { title: 'Settings', roles: ['owner'] } },
    ],
  },
  {
    path: '/admin',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'AdminDashboard', component: () => import('../pages/admin/AdminDashboard.vue'), meta: { roles: ['superadmin'] } },
      { path: 'tenants', name: 'AdminTenants', component: () => import('../pages/admin/AdminTenants.vue'), meta: { roles: ['superadmin'] } },
      { path: 'subscriptions', name: 'AdminSubscriptions', component: () => import('../pages/admin/AdminSubscriptions.vue'), meta: { roles: ['superadmin'] } },
      { path: 'plans', name: 'AdminPlans', component: () => import('../pages/admin/AdminPlans.vue'), meta: { roles: ['superadmin'] } },
      { path: 'orders', name: 'AdminOrders', component: () => import('../pages/admin/AdminOrders.vue'), meta: { roles: ['superadmin'] } },
      { path: 'revenue', name: 'AdminRevenue', component: () => import('../pages/admin/AdminRevenue.vue'), meta: { roles: ['superadmin'] } },
      { path: 'settings', name: 'AdminSettings', component: () => import('../pages/admin/AdminSettings.vue'), meta: { roles: ['superadmin'] } },
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
    if ((to.name === 'Login' || to.name === 'Register') && auth.token) {
      return { path: auth.role === 'superadmin' ? '/admin/dashboard' : '/dashboard' }
    }
    return true
  }
  if (to.meta.requiresAuth && !auth.token) return { name: 'Login' }
  const role = auth.role || 'cashier'
  if (to.meta.roles && to.meta.roles.length && !to.meta.roles.includes(role)) {
    if (role === 'superadmin') return { name: 'AdminDashboard' }
    return { name: 'Dashboard' }
  }
  return true
})

export default router
