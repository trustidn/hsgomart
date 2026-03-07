import client from './client'

export async function getUsers() {
  const { data } = await client.get('/api/users')
  return data
}

export async function createUser(payload) {
  const { data } = await client.post('/api/users', payload)
  return data
}

export async function updateUser(id, payload) {
  const { data } = await client.put(`/api/users/${id}`, payload)
  return data
}

export async function deleteUser(id) {
  await client.delete(`/api/users/${id}`)
}
