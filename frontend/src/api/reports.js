import client from './client'

function dateParams(from, to) {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  return params
}

export async function getSalesSummary(params) {
  const { data } = await client.get('/api/reports/sales', { params: dateParams(params?.from, params?.to) })
  return data
}

export async function getSalesDaily(params) {
  const { data } = await client.get('/api/reports/sales/daily', { params: dateParams(params?.from, params?.to) })
  return data
}

export async function getSalesTransactions(params) {
  const { data } = await client.get('/api/reports/sales/transactions', { params: dateParams(params?.from, params?.to) })
  return data
}

export async function getProfitReport(params) {
  const { data } = await client.get('/api/reports/profit', { params: dateParams(params?.from, params?.to) })
  return data
}

export async function getTopProducts(params) {
  const { data } = await client.get('/api/reports/products', { params: dateParams(params?.from, params?.to) })
  return data
}

export async function getInventoryReport() {
  const { data } = await client.get('/api/reports/inventory')
  return data
}

export async function getCashiersReport(params) {
  const { data } = await client.get('/api/reports/cashiers', { params: dateParams(params?.from, params?.to) })
  return data
}
