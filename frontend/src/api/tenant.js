import client from './client'

export async function getTenantProfile() {
  const { data } = await client.get('/api/tenant/profile')
  return data
}

export async function updateTenantProfile(payload) {
  const { data } = await client.put('/api/tenant/profile', payload)
  return data
}

export async function uploadLogo(file) {
  const formData = new FormData()
  formData.append('logo', file)
  const { data } = await client.post('/api/tenant/logo', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return data
}
