import client from './client'

/**
 * GET /api/reports/inventory - list with product_id, product_name, stock, cost_price
 */
export async function getInventory() {
  const { data } = await client.get('/api/reports/inventory')
  return data
}

/**
 * GET /api/inventory/movements - stock movement history (optional product_id, limit, page)
 * Returns { movements: [], total: number }
 */
export async function getMovements(params = {}) {
  const { data } = await client.get('/api/inventory/movements', { params })
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

/**
 * GET /api/inventory/low-stock - products with stock <= threshold
 */
export async function getLowStock() {
  const { data } = await client.get('/api/inventory/low-stock')
  return data
}
