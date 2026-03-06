import client from './client'

/**
 * GET /api/reports/inventory - list with product_name, stock (and product_id for adjust)
 */
export async function getInventory() {
  const { data } = await client.get('/api/reports/inventory')
  return data
}

export async function getProductStock(productId) {
  const { data } = await client.get(`/api/products/${productId}/stock`)
  return data
}

export async function adjustStock(productId, payload) {
  const { data } = await client.post(`/api/products/${productId}/adjust-stock`, payload)
  return data
}
