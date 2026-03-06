import client from './client'

/**
 * GET /api/reports/sales?from=YYYY-MM-DD&to=YYYY-MM-DD
 * @param {string} [from] - optional start date
 * @param {string} [to] - optional end date
 */
export async function getSalesSummary(from, to) {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  const { data } = await client.get('/api/reports/sales', { params })
  return data
}

/**
 * GET /api/reports/products?from=YYYY-MM-DD&to=YYYY-MM-DD
 */
export async function getTopProducts(from, to) {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  const { data } = await client.get('/api/reports/products', { params })
  return data
}

/**
 * GET /api/reports/inventory
 */
export async function getInventorySummary() {
  const { data } = await client.get('/api/reports/inventory')
  return data
}
