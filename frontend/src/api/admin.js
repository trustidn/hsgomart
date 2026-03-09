import client from './client'

// Tenants
export async function listTenants(status) {
  const params = status ? { status } : {}
  const { data } = await client.get('/admin/tenants', { params })
  return data
}

export async function getTenant(id) {
  const { data } = await client.get(`/admin/tenants/${id}`)
  return data
}

export async function createTenant(payload) {
  const { data } = await client.post('/admin/tenants', payload)
  return data
}

export async function updateTenant(id, payload) {
  const { data } = await client.put(`/admin/tenants/${id}`, payload)
  return data
}

export async function deleteTenant(id) {
  const { data } = await client.delete(`/admin/tenants/${id}`)
  return data
}

export async function resetOwnerPassword(tenantId, newPassword) {
  const { data } = await client.put(`/admin/tenants/${tenantId}/reset-password`, { new_password: newPassword })
  return data
}

// Subscriptions
export async function listSubscriptions() {
  const { data } = await client.get('/admin/subscriptions')
  return data
}

export async function updateSubscription(id, payload) {
  const { data } = await client.put(`/admin/subscriptions/${id}`, payload)
  return data
}

// Plans
export async function listPlans() {
  const { data } = await client.get('/admin/plans')
  return data
}

export async function createPlan(payload) {
  const { data } = await client.post('/admin/plans', payload)
  return data
}

export async function updatePlan(id, payload) {
  const { data } = await client.put(`/admin/plans/${id}`, payload)
  return data
}

export async function deletePlan(id) {
  const { data } = await client.delete(`/admin/plans/${id}`)
  return data
}

// Orders
export async function listOrders(status) {
  const params = status ? { status } : {}
  const { data } = await client.get('/admin/orders', { params })
  return data
}

export async function getOrderDetail(id) {
  const { data } = await client.get(`/admin/orders/${id}`)
  return data
}

export async function approveOrder(id, adminNotes) {
  const { data } = await client.put(`/admin/orders/${id}/approve`, { admin_notes: adminNotes })
  return data
}

export async function rejectOrder(id, adminNotes) {
  const { data } = await client.put(`/admin/orders/${id}/reject`, { admin_notes: adminNotes })
  return data
}

// SaaS Settings
export async function getSettings() {
  const { data } = await client.get('/admin/settings')
  return data
}

export async function updateSettings(payload) {
  const { data } = await client.put('/admin/settings', payload)
  return data
}

export async function uploadSaasLogo(file) {
  const fd = new FormData()
  fd.append('logo', file)
  const { data } = await client.post('/admin/settings/logo', fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return data
}

// Revenue Report
export async function getRevenueReport(from, to) {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  const { data } = await client.get('/admin/reports/revenue', { params })
  return data
}

// Stats
export async function getStats() {
  const { data } = await client.get('/admin/stats')
  return data
}

// Documentation (admin CRUD)
export async function listDocumentation() {
  const { data } = await client.get('/admin/documentation')
  return data
}

export async function createDocumentation(payload) {
  const { data } = await client.post('/admin/documentation', payload)
  return data
}

export async function updateDocumentation(id, payload) {
  const { data } = await client.put(`/admin/documentation/${id}`, payload)
  return data
}

export async function deleteDocumentation(id) {
  const { data } = await client.delete(`/admin/documentation/${id}`)
  return data
}

// Platform Updates (admin CRUD)
export async function listAdminUpdates() {
  const { data } = await client.get('/admin/updates')
  return data
}

export async function createPlatformUpdate(payload) {
  const { data } = await client.post('/admin/updates', payload)
  return data
}

export async function editPlatformUpdate(id, payload) {
  const { data } = await client.put(`/admin/updates/${id}`, payload)
  return data
}

export async function deletePlatformUpdate(id) {
  const { data } = await client.delete(`/admin/updates/${id}`)
  return data
}

// Superadmin User Management
export async function listSuperadmins() {
  const { data } = await client.get('/admin/users')
  return data
}

export async function createSuperadmin(payload) {
  const { data } = await client.post('/admin/users', payload)
  return data
}

export async function updateSuperadmin(id, payload) {
  const { data } = await client.put(`/admin/users/${id}`, payload)
  return data
}

// Public SaaS Info (no auth)
export async function getSaasInfo() {
  const { data } = await client.get('/api/saas-info')
  return data
}

// Public documentation (published only)
export async function getPublicDocumentation() {
  const { data } = await client.get('/api/documentation')
  return data
}

// Public updates
export async function getRecentUpdates() {
  const { data } = await client.get('/api/updates')
  return data
}

export async function getAllUpdates() {
  const { data } = await client.get('/api/updates/all')
  return data
}
