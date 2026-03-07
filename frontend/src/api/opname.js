import client from './client'

export async function startOpname() {
  const { data } = await client.post('/api/inventory/opname')
  return data
}

export async function submitOpnameItems(opnameId, items) {
  const { data } = await client.put(`/api/inventory/opname/${opnameId}`, items)
  return data
}

export async function approveOpname(opnameId) {
  const { data } = await client.post(`/api/inventory/opname/${opnameId}/approve`)
  return data
}

export async function getOpname(opnameId) {
  const { data } = await client.get(`/api/inventory/opname/${opnameId}`)
  return data
}

export async function listOpnames(limit = 50, offset = 0) {
  const { data } = await client.get('/api/inventory/opnames', { params: { limit, offset } })
  return data
}
