import client from './client'

export async function checkout(payload) {
  const { data } = await client.post('/api/pos/checkout', payload)
  return data
}
