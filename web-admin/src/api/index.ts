 import { get, post, put, del, postForm } from './http'
 import type {
   Agent,
   App,
   Archive,
   ExceptionRow,
   Group,
   Job,
   LoginResponse,
   MonitorPoint,
   OperationRow,
   PageData,
   Timing,
   User,
 } from '@/types'

 // —— 鉴权 ——
 export const authApi = {
   login: (username: string, password: string) =>
     postForm<LoginResponse>('/auth/login', { username, password }),
   info: () => get<User>('/auth/info'),
   logout: () => post<void>('/auth/logout'),
   changePassword: (password: string) => postForm<void>('/auth/change_password', { password }),
 }

 // —— 用户 ——
 export const userApi = {
   list: (params: Record<string, unknown>) => get<PageData<User>>('/users', { params }),
   show: (id: number) => get<User>(`/users/${id}`),
   create: (body: Partial<User> & { password: string }) => post<User>('/users', body),
   update: (id: number, body: Partial<User>) => put<User>(`/users/${id}`, body),
   forbidden: (id: number) => post<void>(`/users/${id}/forbidden`),
   resetPassword: (id: number, password: string) =>
     postForm<void>(`/users/${id}/reset_password`, { password }),
 }

 // —— 实例 ——
 export const agentApi = {
   list: (params: Record<string, unknown>) => get<PageData<Agent>>('/agents', { params }),
   show: (id: number) => get<Agent>(`/agents/${id}`),
   update: (id: number, body: { alias: string }) => put<Agent>(`/agents/${id}`, body),
   forbidden: (id: number) => post<void>(`/agents/${id}/forbidden`),
 }

 // —— 分组 ——
 export const groupApi = {
   list: (params: Record<string, unknown>) => get<PageData<Group>>('/groups', { params }),
   create: (body: { name: string; status?: number }) => post<Group>('/groups', body),
   update: (id: number, body: { name?: string; status?: number }) =>
     put<Group>(`/groups/${id}`, body),
   remove: (id: number) => del<void>(`/groups/${id}`),
 }

 // —— 应用 ——
 export const appApi = {
   list: (params: Record<string, unknown>) => get<PageData<App>>('/apps', { params }),
   show: (id: number) => get<App>(`/apps/${id}`),
   create: (body: Partial<App>) => post<{ id: number }>('/apps', body),
   update: (id: number, body: Partial<App>) => put<void>(`/apps/${id}`, body),
   copy: (id: number) => post<{ id: number }>(`/apps/${id}/copy`),
   start: (id: number) => post<void>(`/apps/${id}/start`),
   restart: (id: number) => post<void>(`/apps/${id}/restart`),
   pause: (id: number) => post<void>(`/apps/${id}/pause`),
   remove: (id: number) => del<void>(`/apps/${id}`),
   batch: (action: 'start' | 'restart' | 'pause' | 'delete', ids: number[]) =>
     post<void>(`/apps/batch_${action}`, { ids }),
 }

 // —— 计划任务 ——
 export const jobApi = {
   list: (params: Record<string, unknown>) => get<PageData<Job>>('/jobs', { params }),
   show: (id: number) => get<Job>(`/jobs/${id}`),
   create: (body: Partial<Job>) => post<{ id: number }>('/jobs', body),
   update: (id: number, body: Partial<Job>) => put<void>(`/jobs/${id}`, body),
   copy: (id: number) => post<{ id: number }>(`/jobs/${id}/copy`),
   start: (id: number) => post<void>(`/jobs/${id}/start`),
   restart: (id: number) => post<void>(`/jobs/${id}/restart`),
   pause: (id: number) => post<void>(`/jobs/${id}/pause`),
   remove: (id: number) => del<void>(`/jobs/${id}`),
   batch: (action: 'start' | 'restart' | 'pause' | 'delete', ids: number[]) =>
     post<void>(`/jobs/batch_${action}`, { ids }),
 }

 // —— 定时任务 ——
 export const timingApi = {
   list: (params: Record<string, unknown>) => get<PageData<Timing>>('/timings', { params }),
   show: (id: number) => get<Timing>(`/timings/${id}`),
   create: (body: Partial<Timing>) => post<{ id: number }>('/timings', body),
   update: (id: number, body: Partial<Timing>) => put<void>(`/timings/${id}`, body),
   copy: (id: number) => post<{ id: number }>(`/timings/${id}/copy`),
   start: (id: number) => post<void>(`/timings/${id}/start`),
   restart: (id: number) => post<void>(`/timings/${id}/restart`),
   pause: (id: number) => post<void>(`/timings/${id}/pause`),
   remove: (id: number) => del<void>(`/timings/${id}`),
   batch: (action: 'start' | 'restart' | 'pause' | 'delete', ids: number[]) =>
     post<void>(`/timings/batch_${action}`, { ids }),
 }

 // —— 监控 ——
 export const monitorApi = {
   agent: (agent_id: number, size = 100) =>
     get<MonitorPoint[]>('/monitor/agent', { params: { agent_id, size } }),
   app: (app_id: number, size = 100) =>
     get<MonitorPoint[]>('/monitor/app', { params: { app_id, size } }),
   job: (job_id: number, size = 100) =>
     get<MonitorPoint[]>('/monitor/job', { params: { job_id, size } }),
   timing: (timing_id: number, size = 100) =>
     get<MonitorPoint[]>('/monitor/timing', { params: { timing_id, size } }),
 }

 // —— 归档 ——
 export const archiveApi = {
   app: (app_id: number, page = 1) =>
     get<PageData<Archive>>('/archives/app', { params: { app_id, page } }),
   job: (job_id: number, page = 1) =>
     get<PageData<Archive>>('/archives/job', { params: { job_id, page } }),
   timing: (timing_id: number, page = 1) =>
     get<PageData<Archive>>('/archives/timing', { params: { timing_id, page } }),
 }

 // —— 日志 ——
 export const logApi = {
   appOut: (app_id: number, lines = 50) =>
     get<{ app_id: number; path: string; content: string[] }>('/out_logs/app/data', {
       params: { app_id, lines },
     }),
   appErr: (app_id: number, lines = 50) =>
     get<{ app_id: number; path: string; content: string[] }>('/err_logs/app/data', {
       params: { app_id, lines },
     }),
   jobOut: (job_id: number, lines = 50) =>
     get<{ job_id: number; path: string; content: string[] }>('/out_logs/job/data', {
       params: { job_id, lines },
     }),
   jobErr: (job_id: number, lines = 50) =>
     get<{ job_id: number; path: string; content: string[] }>('/err_logs/job/data', {
       params: { job_id, lines },
     }),
   timingOut: (timing_id: number, lines = 50) =>
     get<{ timing_id: number; path: string; content: string[] }>('/out_logs/timing/data', {
      params: { timing_id, lines },
    }),
   timingErr: (timing_id: number, lines = 50) =>
     get<{ timing_id: number; path: string; content: string[] }>('/err_logs/timing/data', {
      params: { timing_id, lines },
    }),
 }

 // —— 异常 + 操作日志 ——
 export const exceptionApi = {
   list: (params: Record<string, unknown>) =>
     get<PageData<ExceptionRow>>('/exceptions', { params }),
 }

 export const operationApi = {
   list: (params: Record<string, unknown>) =>
     get<PageData<OperationRow>>('/operations', { params }),
 }
