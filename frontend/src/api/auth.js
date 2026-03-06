import client from './client'

/**
 * POST /auth/login
 * @param {string} email
 * @param {string} password
 * @returns {Promise<{ token: string, user?: { id, email, role }, tenant_id?: string }>}
 */
export async function login(email, password) {
  const { data } = await client.post('/auth/login', { email, password })
  return data
}

/**
 * GET /auth/profile (requires Authorization header)
 * @returns {Promise<{ user_id, tenant_id, email, role }>}
 */
export async function getProfile() {
  const { data } = await client.get('/auth/profile')
  return data
}
