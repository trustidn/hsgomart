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
