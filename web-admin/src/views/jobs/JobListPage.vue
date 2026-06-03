 <script setup lang="ts">
 import { computed, onMounted, reactive, ref } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 import { useRouter } from 'vue-router'
 import { agentApi, groupApi, jobApi } from '@/api'
 import type { Agent, Group, Job } from '@/types'
 import { JOB_STATUS, statusInfo } from '@/utils/status'
 import { formatTime } from '@/utils/format'

 // T-308 计划任务
 const router = useRouter()
 const list = ref<Job[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ name: '', group_id: 0, agent_id: 0, status: -99 })
 const page = ref(1)
 const groups = ref<Group[]>([])
 const agents = ref<Agent[]>([])
 const selected = ref<number[]>([])

 const dialogVisible = ref(false)
 const dialogMode = ref<'create' | 'edit'>('create')
 const form = reactive({
   id: 0,
   group_id: 0,
   agent_id: 0,
   name: '',
   dir: '',
   program: '',
   args: '',
   std_out: '',
   std_err: '',
   spec: '',
   timeout: 0,
   is_monitor: 0,
 })

 async function load() {
   loading.value = true
   try {
     const res = await jobApi.list({ ...filter, page: page.value })
     list.value = res.list
     total.value = res.total
   } finally {
     loading.value = false
   }
 }

 async function loadOptions() {
   const [g, a] = await Promise.all([groupApi.list({ page: 1 }), agentApi.list({ page: 1 })])
   groups.value = g.list
   agents.value = a.list
 }

 function onSearch() {
   page.value = 1
   load()
 }

 function openCreate() {
   dialogMode.value = 'create'
   Object.assign(form, {
     id: 0,
     group_id: groups.value[0]?.id || 0,
     agent_id: agents.value[0]?.id || 0,
     name: '',
     dir: '',
     program: '',
     args: '',
     std_out: '',
     std_err: '',
     spec: '0 * * * *',
     timeout: 0,
     is_monitor: 0,
   })
   dialogVisible.value = true
 }

 function openEdit(j: Job) {
   dialogMode.value = 'edit'
   Object.assign(form, j)
   dialogVisible.value = true
 }

 async function submitForm() {
   if (dialogMode.value === 'create') await jobApi.create({ ...form })
   else await jobApi.update(form.id, { ...form })
   ElMessage.success('已保存')
   dialogVisible.value = false
   load()
 }

 async function onAction(action: 'start' | 'restart' | 'pause' | 'delete', ids: number[]) {
   await ElMessageBox.confirm(`确定${actionLabel(action)} ${ids.length} 个计划任务？`, '提示', {
     type: 'warning',
   })
   await jobApi.batch(action, ids)
   ElMessage.success('已提交')
   selected.value = []
   load()
 }

 async function onSingle(action: 'start' | 'restart' | 'pause' | 'delete', j: Job) {
   await ElMessageBox.confirm(`确定${actionLabel(action)}「${j.name}」？`, '提示', { type: 'warning' })
   if (action === 'start') await jobApi.start(j.id)
   else if (action === 'restart') await jobApi.restart(j.id)
   else if (action === 'pause') await jobApi.pause(j.id)
   else await jobApi.remove(j.id)
   ElMessage.success('已提交')
   load()
 }

 function actionLabel(a: string) {
   return { start: '启动', restart: '重启', pause: '暂停', delete: '删除' }[a] || a
 }

 const hasSelection = computed(() => selected.value.length > 0)

 onMounted(async () => {
   await loadOptions()
   await load()
 })
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">计划任务</span>
       <el-button type="primary" @click="openCreate">新建计划任务</el-button>
     </div>

     <div class="page-toolbar">
       <el-input v-model="filter.name" placeholder="名称" clearable style="width: 200px" />
       <el-select v-model="filter.status" placeholder="状态" style="width: 140px">
         <el-option label="全部" :value="-99" />
         <el-option v-for="s in JOB_STATUS" :key="s.id" :label="s.name" :value="s.id" />
       </el-select>
       <el-button type="primary" @click="onSearch">查询</el-button>
       <el-divider direction="vertical" />
       <el-button :disabled="!hasSelection" @click="onAction('start', selected)">批量启动</el-button>
       <el-button :disabled="!hasSelection" @click="onAction('restart', selected)">批量重启</el-button>
       <el-button :disabled="!hasSelection" @click="onAction('pause', selected)">批量暂停</el-button>
       <el-button :disabled="!hasSelection" type="danger" @click="onAction('delete', selected)">批量删除</el-button>
     </div>

     <el-table
       v-loading="loading"
       :data="list"
       border
       @selection-change="(rows: Job[]) => (selected = rows.map((r) => r.id))"
     >
       <el-table-column type="selection" width="48" />
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="name" label="名称" min-width="160" />
       <el-table-column prop="group_name" label="分组" />
       <el-table-column prop="agent_name" label="实例" />
       <el-table-column prop="spec" label="cron" width="120" />
       <el-table-column label="状态" width="100">
         <template #default="{ row }">
           <el-tag :type="statusInfo(JOB_STATUS, row.status).type as 'success' | 'danger' | 'info' | 'warning'">
             {{ statusInfo(JOB_STATUS, row.status).name }}
           </el-tag>
         </template>
       </el-table-column>
       <el-table-column label="操作" width="280" fixed="right">
         <template #default="{ row }">
           <el-button size="small" @click="onSingle('start', row)">启动</el-button>
           <el-button size="small" @click="onSingle('restart', row)">重启</el-button>
           <el-button size="small" @click="onSingle('pause', row)">暂停</el-button>
           <el-button size="small" @click="openEdit(row)">编辑</el-button>
           <el-button size="small" @click="jobApi.copy(row.id).then(load)">复制</el-button>
           <el-button size="small" type="danger" @click="onSingle('delete', row)">删除</el-button>
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

     <el-dialog v-model="dialogVisible" :title="dialogMode === 'create' ? '新建计划任务' : '编辑计划任务'" width="600px">
       <el-form label-width="100px">
         <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
         <el-form-item label="分组">
           <el-select v-model="form.group_id" style="width: 100%">
             <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
           </el-select>
         </el-form-item>
         <el-form-item label="实例">
           <el-select v-model="form.agent_id" style="width: 100%">
             <el-option
               v-for="a in agents"
               :key="a.id"
               :label="a.alias || `${a.ip}:${a.port}`"
               :value="a.id"
             />
           </el-select>
         </el-form-item>
         <el-form-item label="执行目录"><el-input v-model="form.dir" /></el-form-item>
         <el-form-item label="程序"><el-input v-model="form.program" /></el-form-item>
         <el-form-item label="参数"><el-input v-model="form.args" /></el-form-item>
         <el-form-item label="stdout"><el-input v-model="form.std_out" /></el-form-item>
         <el-form-item label="stderr"><el-input v-model="form.std_err" /></el-form-item>
         <el-form-item label="cron 表达式">
           <el-input v-model="form.spec" placeholder="0 * * * *" />
         </el-form-item>
         <el-form-item label="超时(秒)">
           <el-input-number v-model="form.timeout" :min="0" />
         </el-form-item>
       </el-form>
       <template #footer>
         <el-button @click="dialogVisible = false">取消</el-button>
         <el-button type="primary" @click="submitForm">提交</el-button>
       </template>
     </el-dialog>
   </div>
 </template>
