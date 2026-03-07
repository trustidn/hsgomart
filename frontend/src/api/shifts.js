import client from './client'

export async function openShift(payload) {
  const { data } = await client.post('/api/shifts/open', payload)
  return data
}

export async function closeShift(payload) {
  const { data } = await client.post('/api/shifts/close', payload)
  return data
}

export async function getCurrentShift() {
  const { data } = await client.get('/api/shifts/current')
  return data
}

export async function listShifts(params = {}) {
  const { data } = await client.get('/api/shifts', { params })
  return data
}
