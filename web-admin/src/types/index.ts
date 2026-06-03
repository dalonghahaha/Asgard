 // 通用类型定义
 export type ID = number

 // 分页响应（与后端 {list, total, page, page_size} 契约一致）
 export interface PageData<T> {
   list: T[]
   total: number
   page: number
   page_size: number
 }

 export type UserRole = 'Administrator' | 'User'

 export interface User {
   id: ID
   nickname: string
   email: string
   mobile: string
   avatar: string
   role: UserRole
   status: number
   created_at: string
 }

 export interface Agent {
   id: ID
   alias: string
   ip: string
   port: string
   master: string
   status: number
   created_at: string
 }

 export interface Group {
   id: ID
   name: string
   status: number
   creator: ID
   created_at: string
 }

 export interface CmdBase {
   group_id: ID
   agent_id: ID
   name: string
   dir: string
   program: string
   args: string
   std_out: string
   std_err: string
   is_monitor: number
 }

 export interface App extends CmdBase {
   id: ID
   auto_restart: number
   status: number
   group_name: string
   agent_name: string
   created_at: string
 }

 export interface Job extends CmdBase {
   id: ID
   spec: string
   timeout: number
   status: number
   group_name: string
   agent_name: string
   created_at: string
 }

 export interface Timing extends CmdBase {
   id: ID
   time: string
   timeout: number
   status: number
   group_name: string
   agent_name: string
   created_at: string
 }

 export interface Archive {
   id: ID
   uuid: string
   pid: number
   begin_time: string
   end_time: string
   status: number
   signal: string
 }

 export interface MonitorPoint {
   cpu: number
   memory: number
   created_at: string
 }

 export interface ExceptionRow {
   id: ID
   type: string
   name: string
   desc: string
   created_at: string
 }

 export interface OperationRow {
   id: ID
   desc: string
   created_at: string
 }

 export interface LoginResponse {
   token: string
   expires_at: number
   user: User
 }
