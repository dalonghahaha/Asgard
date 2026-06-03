 import { defineStore } from 'pinia'
 import { ref, computed } from 'vue'
 import { authApi } from '@/api'
 import type { User } from '@/types'

 const TOKEN_KEY = 'asgard_token'
 const USER_KEY = 'asgard_user'

 // T-206 全局状态：token + 当前用户 + 角色判定
 export const useAuthStore = defineStore('auth', () => {
   const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
   const user = ref<User | null>(loadUser())

   const isLogin = computed(() => !!token.value)
   const isAdmin = computed(() => user.value?.role === 'Administrator')

   async function login(username: string, password: string) {
     const res = await authApi.login(username, password)
     token.value = res.token
     user.value = res.user
     localStorage.setItem(TOKEN_KEY, res.token)
     localStorage.setItem(USER_KEY, JSON.stringify(res.user))
     return res
   }

   async function fetchInfo() {
     const info = await authApi.info()
     user.value = info
     localStorage.setItem(USER_KEY, JSON.stringify(info))
     return info
   }

   async function logout() {
     try {
       await authApi.logout()
     } catch {
       // 即便后端失败也要清掉本地状态
     }
     clear()
   }

   async function changePassword(password: string) {
     await authApi.changePassword(password)
     clear()
   }

   function clear() {
     token.value = ''
     user.value = null
     localStorage.removeItem(TOKEN_KEY)
     localStorage.removeItem(USER_KEY)
   }

   return {
     token,
     user,
     isLogin,
     isAdmin,
     login,
     fetchInfo,
     logout,
     changePassword,
     clear,
   }
 })

 function loadUser(): User | null {
   const raw = localStorage.getItem(USER_KEY)
   if (!raw) return null
   try {
     return JSON.parse(raw) as User
   } catch {
     return null
   }
 }
