import client from './client'

export async function getProducts() {
  const { data } = await client.get('/api/products')
  return data
}

export async function getProductByBarcode(barcode) {
  const { data } = await client.get(`/api/products/barcode/${encodeURIComponent(barcode)}`)
  return data
}

export async function getProduct(id) {
  const { data } = await client.get(`/api/products/${id}`)
  return data
}

export async function createProduct(payload) {
  const { data } = await client.post('/api/products', payload)
  return data
}

export async function updateProduct(id, payload) {
  const { data } = await client.put(`/api/products/${id}`, payload)
  return data
}

export async function deleteProduct(id) {
  await client.delete(`/api/products/${id}`)
}

export async function getCategories() {
  const { data } = await client.get('/api/categories')
  return data
}

export async function createCategory(payload) {
  const { data } = await client.post('/api/categories', payload)
  return data
}

export async function updateCategory(id, payload) {
  const { data } = await client.put(`/api/categories/${id}`, payload)
  return data
}

export async function deleteCategory(id) {
  await client.delete(`/api/categories/${id}`)
}

export async function addBarcode(productId, barcode) {
  const { data } = await client.post(`/api/products/${productId}/barcodes`, { barcode })
  return data
}
