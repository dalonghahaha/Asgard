 // 状态常量/映射：与后端 constants/constant.go 保持一致

 export const USER_STATUS = [
   { id: 0, name: '未审核', type: 'info' },
   { id: 1, name: '正常', type: 'success' },
   { id: -1, name: '禁用', type: 'danger' },
 ]

 export const AGENT_STATUS = [
   { id: 1, name: '在线', type: 'success' },
   { id: 0, name: '离线', type: 'info' },
   { id: -1, name: '禁用', type: 'danger' },
 ]

 export const GROUP_STATUS = [
   { id: 1, name: '启用', type: 'success' },
   { id: 0, name: '未启用', type: 'info' },
   { id: -1, name: '已删除', type: 'danger' },
 ]

 export const APP_STATUS = [
   { id: 1, name: '运行中', type: 'success' },
   { id: 2, name: '暂停', type: 'warning' },
   { id: 0, name: '停止', type: 'info' },
   { id: -2, name: '未知', type: 'warning' },
   { id: -3, name: '失败', type: 'danger' },
   { id: -1, name: '已删除', type: 'danger' },
 ]

 export const JOB_STATUS = APP_STATUS

 export const TIMING_STATUS = [
   { id: 1, name: '运行中', type: 'success' },
   { id: 2, name: '暂停', type: 'warning' },
   { id: 0, name: '停止', type: 'info' },
   { id: 3, name: '已完成', type: 'success' },
   { id: -2, name: '未知', type: 'warning' },
   { id: -1, name: '已删除', type: 'danger' },
 ]

 export function statusInfo(list: { id: number; name: string; type: string }[], id: number) {
   return list.find((s) => s.id === id) || { id, name: String(id), type: 'info' }
 }
