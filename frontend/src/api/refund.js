import client from './client'

export async function createRefund(transactionId, reason) {
  const { data } = await client.post('/api/pos/refund', { transaction_id: transactionId, reason })
  return data
}

export async function listRefunds(limit = 50, offset = 0) {
  const { data } = await client.get('/api/refunds', { params: { limit, offset } })
  return data
}
