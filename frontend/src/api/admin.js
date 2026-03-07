import client from './client'

export async function listTenants() {
  const { data } = await client.get('/admin/tenants')
  return data
}

export async function getTenant(id) {
  const { data } = await client.get(`/admin/tenants/${id}`)
  return data
}

export async function updateTenant(id, status) {
  const { data } = await client.put(`/admin/tenants/${id}`, { status })
  return data
}

export async function listSubscriptions() {
  const { data } = await client.get('/admin/subscriptions')
  return data
}

export async function updateSubscription(id, payload) {
  const { data } = await client.put(`/admin/subscriptions/${id}`, payload)
  return data
}

export async function getStats() {
  const { data } = await client.get('/admin/stats')
  return data
}
