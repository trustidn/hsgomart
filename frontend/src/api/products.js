import client from './client'

export async function getProducts() {
  const { data } = await client.get('/api/products')
  return data
}

export async function createProduct(payload) {
  const { data } = await client.post('/api/products', payload)
  return data
}

export async function getCategories() {
  const { data } = await client.get('/api/categories')
  return data
}

export async function createCategory(payload) {
  const { data } = await client.post('/api/categories', payload)
  return data
}

export async function addBarcode(productId, barcode) {
  const { data } = await client.post(`/api/products/${productId}/barcodes`, { barcode })
  return data
}
