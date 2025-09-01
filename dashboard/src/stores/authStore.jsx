// TODO: Implement auth state with Jotai during development

import { atom } from 'jotai'

// Auth atoms
export const userAtom = atom(null)
export const isAuthenticatedAtom = atom((get) => {
  const user = get(userAtom)
  return user !== null
})

export const tenantAtom = atom(null)
export const permissionsAtom = atom([])

// Auth actions
export const loginAtom = atom(
  null,
  async (get, set, { email, password }) => {
    // TODO: Implement login logic
    try {
      // const response = await authApi.login(email, password)
      // set(userAtom, response.user)
      // set(tenantAtom, response.tenant)
      // set(permissionsAtom, response.permissions)
      // localStorage.setItem('token', response.token)
    } catch (error) {
      throw error
    }
  }
)

export const logoutAtom = atom(
  null,
  (get, set) => {
    // TODO: Implement logout logic
    set(userAtom, null)
    set(tenantAtom, null)
    set(permissionsAtom, [])
    // localStorage.removeItem('token')
  }
)

export const initAuthAtom = atom(
  null,
  async (get, set) => {
    // TODO: Implement auth initialization
    // Check for stored token and validate
    // const token = localStorage.getItem('token')
    // if (token) {
    //   try {
    //     const user = await authApi.validateToken(token)
    //     set(userAtom, user)
    //   } catch (error) {
    //     localStorage.removeItem('token')
    //   }
    // }
  }
)