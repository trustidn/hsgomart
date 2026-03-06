import client from './client'

export async function listPurchases() {
  const { data } = await client.get('/api/purchases')
  return data?.purchases ?? []
}

export async function getPurchase(id) {
  const { data } = await client.get(`/api/purchases/${id}`)
  return data
}

export async function createPurchase(body) {
  const { data } = await client.post('/api/purchases', body)
  return data
}
