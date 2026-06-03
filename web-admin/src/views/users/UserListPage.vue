 <script setup lang="ts">
 import { onMounted, reactive, ref } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 import { userApi } from '@/api'
 import type { User } from '@/types'
 import { USER_STATUS, statusInfo } from '@/utils/status'
 import { formatTime } from '@/utils/format'
 import { useAuthStore } from '@/stores/auth'

 // T-304 用户管理：列表 + 新建 + 编辑 + 禁用 + 重置密码
 const auth = useAuthStore()
 const list = ref<User[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ nickname: '', email: '', status: -99 })
 const page = ref(1)

 const dialogVisible = ref(false)
 const dialogMode = ref<'create' | 'edit'>('create')
 const form = reactive({
   id: 0,
   nickname: '',
   email: '',
   mobile: '',
   password: '',
   role: 'User' as 'User' | 'Administrator',
   status: 1,
 })

 const resetPwdVisible = ref(false)
 const resetPwdForm = reactive({ id: 0, password: '' })

 async function load() {
   loading.value = true
   try {
     const res = await userApi.list({ ...filter, page: page.value })
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
   Object.assign(form, {
     id: 0,
     nickname: '',
     email: '',
     mobile: '',
     password: '',
     role: 'User',
     status: 1,
   })
   dialogVisible.value = true
 }

 function openEdit(u: User) {
   dialogMode.value = 'edit'
   Object.assign(form, { ...u, password: '' })
   dialogVisible.value = true
 }

 async function submitForm() {
   if (dialogMode.value === 'create') {
     await userApi.create({ ...form })
     ElMessage.success('创建成功')
   } else {
     await userApi.update(form.id, { ...form })
     ElMessage.success('更新成功')
   }
   dialogVisible.value = false
   load()
 }

 async function onForbidden(u: User) {
   await ElMessageBox.confirm(`确定禁用「${u.nickname}」吗？`, '提示', { type: 'warning' })
   await userApi.forbidden(u.id)
   ElMessage.success('已禁用')
   load()
 }

 function openResetPwd(u: User) {
   Object.assign(resetPwdForm, { id: u.id, password: '' })
   resetPwdVisible.value = true
 }

 async function submitResetPwd() {
   if (!resetPwdForm.password) {
     ElMessage.warning('请输入新密码')
     return
   }
   await userApi.resetPassword(resetPwdForm.id, resetPwdForm.password)
   ElMessage.success('密码已重置')
   resetPwdVisible.value = false
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">用户管理</span>
       <el-button v-if="auth.isAdmin" type="primary" @click="openCreate">新建用户</el-button>
     </div>

     <div class="page-toolbar">
       <el-input v-model="filter.nickname" placeholder="昵称" clearable style="width: 200px" />
       <el-input v-model="filter.email" placeholder="邮箱" clearable style="width: 200px" />
       <el-select v-model="filter.status" placeholder="状态" style="width: 140px">
         <el-option label="全部" :value="-99" />
         <el-option v-for="s in USER_STATUS" :key="s.id" :label="s.name" :value="s.id" />
       </el-select>
       <el-button type="primary" @click="onSearch">查询</el-button>
     </div>

     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="nickname" label="昵称" />
       <el-table-column prop="email" label="邮箱" />
       <el-table-column prop="mobile" label="手机" />
       <el-table-column prop="role" label="角色" width="100">
         <template #default="{ row }">
           <el-tag :type="row.role === 'Administrator' ? 'danger' : 'info'">
             {{ row.role === 'Administrator' ? '管理员' : '普通' }}
           </el-tag>
         </template>
       </el-table-column>
       <el-table-column label="状态" width="100">
         <template #default="{ row }">
           <el-tag :type="statusInfo(USER_STATUS, row.status).type as 'success' | 'danger' | 'info'">
             {{ statusInfo(USER_STATUS, row.status).name }}
           </el-tag>
         </template>
       </el-table-column>
       <el-table-column label="创建时间" width="170">
         <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
       </el-table-column>
       <el-table-column v-if="auth.isAdmin" label="操作" width="220">
         <template #default="{ row }">
           <el-button size="small" @click="openEdit(row)">编辑</el-button>
           <el-button size="small" @click="openResetPwd(row)">重置密码</el-button>
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

     <el-dialog v-model="dialogVisible" :title="dialogMode === 'create' ? '新建用户' : '编辑用户'" width="480px">
       <el-form label-width="80px">
         <el-form-item label="昵称">
           <el-input v-model="form.nickname" />
         </el-form-item>
         <el-form-item label="邮箱">
           <el-input v-model="form.email" />
         </el-form-item>
         <el-form-item label="手机号">
           <el-input v-model="form.mobile" />
         </el-form-item>
         <el-form-item v-if="dialogMode === 'create'" label="密码">
           <el-input v-model="form.password" type="password" show-password />
         </el-form-item>
         <el-form-item label="角色">
           <el-radio-group v-model="form.role">
             <el-radio value="User">普通</el-radio>
             <el-radio value="Administrator">管理员</el-radio>
           </el-radio-group>
         </el-form-item>
         <el-form-item label="状态">
           <el-select v-model="form.status" style="width: 100%">
             <el-option label="正常" :value="1" />
             <el-option label="禁用" :value="-1" />
           </el-select>
         </el-form-item>
       </el-form>
       <template #footer>
         <el-button @click="dialogVisible = false">取消</el-button>
         <el-button type="primary" @click="submitForm">提交</el-button>
       </template>
     </el-dialog>

     <el-dialog v-model="resetPwdVisible" title="重置密码" width="400px">
       <el-form label-width="80px">
         <el-form-item label="新密码">
           <el-input v-model="resetPwdForm.password" type="password" show-password />
         </el-form-item>
       </el-form>
       <template #footer>
         <el-button @click="resetPwdVisible = false">取消</el-button>
         <el-button type="primary" @click="submitResetPwd">提交</el-button>
       </template>
     </el-dialog>
   </div>
 </template>
