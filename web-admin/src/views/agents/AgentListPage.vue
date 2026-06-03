 <script setup lang="ts">
 import { onMounted, reactive, ref } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 import { agentApi } from '@/api'
 import type { Agent } from '@/types'
 import { AGENT_STATUS, statusInfo } from '@/utils/status'
 import { formatTime } from '@/utils/format'
 import { useAuthStore } from '@/stores/auth'

 // T-305 实例管理：列表 + 改别名 + 禁用
 const auth = useAuthStore()
 const list = ref<Agent[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ alias: '', status: -99 })
 const page = ref(1)

 const editVisible = ref(false)
 const editForm = reactive({ id: 0, alias: '' })

 async function load() {
   loading.value = true
   try {
     const res = await agentApi.list({ ...filter, page: page.value })
     list.value = res.list
     total.value = res.total
   } finally {
     loading.value = false
   }
 }

 function onSearch() {
   page.value = 1
   load()
 }

 function openEdit(a: Agent) {
   Object.assign(editForm, { id: a.id, alias: a.alias })
   editVisible.value = true
 }

 async function submitEdit() {
   await agentApi.update(editForm.id, { alias: editForm.alias })
   ElMessage.success('已更新')
   editVisible.value = false
   load()
 }

 async function onForbidden(a: Agent) {
   await ElMessageBox.confirm(
     `禁用「${a.alias || a.ip}」会级联停用该实例上所有 app/job/timing，确定吗？`,
     '提示',
     { type: 'warning' },
   )
   await agentApi.forbidden(a.id)
   ElMessage.success('已禁用')
   load()
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">实例管理</span>
     </div>

     <div class="page-toolbar">
       <el-input v-model="filter.alias" placeholder="别名" clearable style="width: 200px" />
       <el-select v-model="filter.status" placeholder="状态" style="width: 140px">
         <el-option label="全部" :value="-99" />
         <el-option v-for="s in AGENT_STATUS" :key="s.id" :label="s.name" :value="s.id" />
       </el-select>
       <el-button type="primary" @click="onSearch">查询</el-button>
     </div>

     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="alias" label="别名" />
       <el-table-column label="地址">
         <template #default="{ row }">{{ row.ip }}:{{ row.port }}</template>
       </el-table-column>
       <el-table-column label="状态" width="100">
         <template #default="{ row }">
           <el-tag :type="statusInfo(AGENT_STATUS, row.status).type as 'success' | 'danger' | 'info'">
             {{ statusInfo(AGENT_STATUS, row.status).name }}
           </el-tag>
         </template>
       </el-table-column>
       <el-table-column label="注册时间" width="170">
         <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
       </el-table-column>
       <el-table-column v-if="auth.isAdmin" label="操作" width="200">
         <template #default="{ row }">
           <el-button size="small" @click="openEdit(row)">改别名</el-button>
           <el-button
             size="small"
             type="danger"
             :disabled="row.status === -1"
             @click="onForbidden(row)"
           >
             禁用
           </el-button>
         </template>
       </el-table-column>
     </el-table>

     <div class="pagination-bar">
       <el-pagination
         v-model:current-page="page"
         :total="total"
         :page-size="20"
         layout="prev, pager, next, total"
         @current-change="load"
       />
     </div>

     <el-dialog v-model="editVisible" title="编辑别名" width="400px">
       <el-form label-width="60px">
         <el-form-item label="别名">
           <el-input v-model="editForm.alias" />
         </el-form-item>
       </el-form>
       <template #footer>
         <el-button @click="editVisible = false">取消</el-button>
         <el-button type="primary" @click="submitEdit">提交</el-button>
       </template>
     </el-dialog>
   </div>
 </template>
