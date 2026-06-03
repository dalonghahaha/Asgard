 <script setup lang="ts">
 import { onMounted, reactive, ref } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 import { groupApi } from '@/api'
 import type { Group } from '@/types'
 import { GROUP_STATUS, statusInfo } from '@/utils/status'
 import { formatTime } from '@/utils/format'

 // T-306 分组管理
 const list = ref<Group[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ name: '', status: -99 })
 const page = ref(1)

 const dialogVisible = ref(false)
 const dialogMode = ref<'create' | 'edit'>('create')
 const form = reactive({ id: 0, name: '', status: 1 })

 async function load() {
   loading.value = true
   try {
     const res = await groupApi.list({ ...filter, page: page.value })
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

 function openCreate() {
   dialogMode.value = 'create'
   Object.assign(form, { id: 0, name: '', status: 1 })
   dialogVisible.value = true
 }

 function openEdit(g: Group) {
   dialogMode.value = 'edit'
   Object.assign(form, { id: g.id, name: g.name, status: g.status })
   dialogVisible.value = true
 }

 async function submitForm() {
   if (dialogMode.value === 'create') await groupApi.create({ ...form })
   else await groupApi.update(form.id, { ...form })
   ElMessage.success('已保存')
   dialogVisible.value = false
   load()
 }

 async function onDelete(g: Group) {
   await ElMessageBox.confirm(`确定删除「${g.name}」吗？`, '提示', { type: 'warning' })
   await groupApi.remove(g.id)
   ElMessage.success('已删除')
   load()
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">分组管理</span>
       <el-button type="primary" @click="openCreate">新建分组</el-button>
     </div>

     <div class="page-toolbar">
       <el-input v-model="filter.name" placeholder="名称" clearable style="width: 200px" />
       <el-select v-model="filter.status" placeholder="状态" style="width: 140px">
         <el-option label="全部" :value="-99" />
         <el-option v-for="s in GROUP_STATUS" :key="s.id" :label="s.name" :value="s.id" />
       </el-select>
       <el-button type="primary" @click="onSearch">查询</el-button>
     </div>

     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="name" label="名称" />
       <el-table-column label="状态" width="100">
         <template #default="{ row }">
           <el-tag :type="statusInfo(GROUP_STATUS, row.status).type as 'success' | 'danger' | 'info'">
             {{ statusInfo(GROUP_STATUS, row.status).name }}
           </el-tag>
         </template>
       </el-table-column>
       <el-table-column label="创建时间" width="170">
         <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
       </el-table-column>
       <el-table-column label="操作" width="180">
         <template #default="{ row }">
           <el-button size="small" @click="openEdit(row)">编辑</el-button>
           <el-button
             size="small"
             type="danger"
             :disabled="row.status === -1"
             @click="onDelete(row)"
           >
             删除
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

     <el-dialog v-model="dialogVisible" :title="dialogMode === 'create' ? '新建分组' : '编辑分组'" width="400px">
       <el-form label-width="60px">
         <el-form-item label="名称">
           <el-input v-model="form.name" />
         </el-form-item>
         <el-form-item label="状态">
           <el-radio-group v-model="form.status">
             <el-radio :value="1">启用</el-radio>
             <el-radio :value="0">未启用</el-radio>
           </el-radio-group>
         </el-form-item>
       </el-form>
       <template #footer>
         <el-button @click="dialogVisible = false">取消</el-button>
         <el-button type="primary" @click="submitForm">提交</el-button>
       </template>
     </el-dialog>
   </div>
 </template>
