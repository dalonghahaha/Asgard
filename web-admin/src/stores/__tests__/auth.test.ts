 import { describe, it, expect, beforeEach } from 'vitest'
 import { setActivePinia, createPinia } from 'pinia'
 import { useAuthStore } from '../auth'

 // T-211 Pinia auth store 单元测试：token/user 持久化与清理
 describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('starts logged out', () => {
    const auth = useAuthStore()
    expect(auth.isLogin).toBe(false)
    expect(auth.user).toBeNull()
    expect(auth.isAdmin).toBe(false)
  })

  it('clear removes token and user', () => {
    const auth = useAuthStore()
    localStorage.setItem('asgard_token', 'abc')
    auth.clear()
    expect(auth.token).toBe('')
    expect(auth.user).toBeNull()
    expect(localStorage.getItem('asgard_token')).toBeNull()
  })

  it('isAdmin reflects role', () => {
    const auth = useAuthStore()
    auth.user = { id: 1, nickname: 'a', role: 'Administrator' } as never
    expect(auth.isAdmin).toBe(true)
    auth.user = { id: 1, nickname: 'a', role: 'User' } as never
    expect(auth.isAdmin).toBe(false)
  })
 })
