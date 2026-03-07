import client from './client'

export async function getSubscription() {
  const { data } = await client.get('/api/subscription')
  return data
}

export async function listPlans() {
  const { data } = await client.get('/api/subscription/plans')
  return data
}

export async function changePlan(planId) {
  const { data } = await client.put('/api/subscription', { plan_id: planId })
  return data
}

export async function createOrder(planId) {
  const { data } = await client.post('/api/subscription/order', { plan_id: planId })
  return data
}

export async function uploadPaymentProof(orderId, file) {
  const formData = new FormData()
  formData.append('payment_proof', file)
  const { data } = await client.post(`/api/subscription/order/${orderId}/payment`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return data
}

export async function getOrders() {
  const { data } = await client.get('/api/subscription/orders')
  return data
}

export async function getOrder(orderId) {
  const { data } = await client.get(`/api/subscription/order/${orderId}`)
  return data
}
