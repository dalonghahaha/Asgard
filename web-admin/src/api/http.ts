 import axios, {
   type AxiosError,
   type AxiosInstance,
   type AxiosRequestConfig,
   type AxiosResponse,
   type InternalAxiosRequestConfig,
 } from 'axios'
 import { ElMessage } from 'element-plus'
 import { useAuthStore } from '@/stores/auth'
 import router from '@/router'
 import type { PageData } from '@/types'

 export type { PageData }

 // T-205 统一 axios 封装：
 //  - baseURL 走 Vite 代理（开发期 /api → http://localhost:12345）
 //  - request 拦截器自动注入 Bearer token
 //  - response 拦截器：code !== 200 走 ElMessage；401 跳 /login

 export interface ApiEnvelope<T = unknown> {
   code: number
   message: string
   data: T
 }

 const http: AxiosInstance = axios.create({
   baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
   timeout: 15000,
 })

 // 简易 loading 计数（避免每个请求都搞一个全屏 spinner）
 let loadingCount = 0
 function showLoading() {
   loadingCount += 1
 }
 function hideLoading() {
   loadingCount = Math.max(0, loadingCount - 1)
 }

 http.interceptors.request.use(
   (config: InternalAxiosRequestConfig) => {
     showLoading()
     const auth = useAuthStore()
     if (auth.token && !config.headers['Authorization']) {
       config.headers['Authorization'] = `Bearer ${auth.token}`
     }
     return config
   },
   (error: AxiosError) => {
     hideLoading()
     return Promise.reject(error)
   },
 )

 http.interceptors.response.use(
  (response: AxiosResponse<ApiEnvelope>) => {
    hideLoading()
    const env = response.data
    if (!env || typeof env.code === 'undefined') {
      // 兼容裸 JSON（理论上后端都包了一层 envelope）
      return response.data as unknown as AxiosResponse
    }
    if (env.code === 200) {
      return env as unknown as AxiosResponse
    }
    if (env.code === 401) {
      const auth = useAuthStore()
      auth.clear()
      router.push({ name: 'login' })
      ElMessage.error(env.message || '登录已失效')
      return Promise.reject(new Error(env.message || 'unauthorized'))
    }
    ElMessage.error(env.message || '请求失败')
    return Promise.reject(new Error(env.message || 'request failed'))
  },
  (error: AxiosError<ApiEnvelope>) => {
    hideLoading()
    if (error.response?.status === 401) {
      const auth = useAuthStore()
      auth.clear()
      router.push({ name: 'login' })
      ElMessage.error('登录已失效')
    } else if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error(error.message || '网络异常')
    }
    return Promise.reject(error)
  },
 )

 export interface RequestOptions {
   silent?: boolean
 }

 // 便捷方法：返回值直接是 data 字段
 export async function get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
   const res = await http.get<ApiEnvelope<T>>(url, config)
   return (res.data as ApiEnvelope<T>).data
 }

 export async function post<T>(
   url: string,
   body?: unknown,
   config?: AxiosRequestConfig,
 ): Promise<T> {
   const res = await http.post<ApiEnvelope<T>>(url, body, config)
   return (res.data as ApiEnvelope<T>).data
 }

 export async function put<T>(
   url: string,
   body?: unknown,
   config?: AxiosRequestConfig,
 ): Promise<T> {
   const res = await http.put<ApiEnvelope<T>>(url, body, config)
   return (res.data as ApiEnvelope<T>).data
 }

 export async function del<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
   const res = await http.delete<ApiEnvelope<T>>(url, config)
   return (res.data as ApiEnvelope<T>).data
 }

 export async function postForm<T>(
   url: string,
   body: Record<string, unknown>,
   config?: AxiosRequestConfig,
 ): Promise<T> {
   const form = new URLSearchParams()
   for (const [k, v] of Object.entries(body)) {
     if (v === undefined || v === null) continue
     form.set(k, String(v))
   }
   const cfg: AxiosRequestConfig = {
     headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
     ...config,
   }
   const res = await http.post<ApiEnvelope<T>>(url, form.toString(), cfg)
   return (res.data as ApiEnvelope<T>).data
 }

 export default http
