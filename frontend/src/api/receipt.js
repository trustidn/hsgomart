import client from './client'

export async function getReceipt(transactionId) {
  const { data } = await client.get(`/api/pos/receipt/${transactionId}`)
  return data
}
