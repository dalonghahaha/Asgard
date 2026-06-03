 <script setup lang="ts">
 import { onMounted, ref } from 'vue'
 import { useRoute, useRouter } from 'vue-router'
 import { appApi } from '@/api'
 import type { App } from '@/types'
 import { APP_STATUS, statusInfo } from '@/utils/status'
 import { formatTime } from '@/utils/format'

 // T-307 应用详情：展示 + 跳日志 + 跳监控 + 跳归档
 const route = useRoute()
 const router = useRouter()
 const app = ref<App | null>(null)

 async function load() {
   const id = Number(route.params.id)
   app.value = await appApi.show(id)
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page" v-if="app">
     <el-page-header @back="router.back()">
       <template #content>
         <span>应用详情 - {{ app.name }}</span>
       </template>
     </el-page-header>

     <el-descriptions :column="2" border style="margin-top: 16px">
       <el-descriptions-item label="ID">{{ app.id }}</el-descriptions-item>
       <el-descriptions-item label="状态">
         <el-tag :type="statusInfo(APP_STATUS, app.status).type as 'success' | 'danger' | 'info' | 'warning'">
           {{ statusInfo(APP_STATUS, app.status).name }}
         </el-tag>
       </el-descriptions-item>
       <el-descriptions-item label="分组">{{ app.group_name }}</el-descriptions-item>
       <el-descriptions-item label="实例">{{ app.agent_name }}</el-descriptions-item>
       <el-descriptions-item label="执行目录" :span="2">{{ app.dir }}</el-descriptions-item>
       <el-descriptions-item label="程序">{{ app.program }}</el-descriptions-item>
       <el-descriptions-item label="参数">{{ app.args }}</el-descriptions-item>
       <el-descriptions-item label="stdout" :span="2">{{ app.std_out }}</el-descriptions-item>
       <el-descriptions-item label="stderr" :span="2">{{ app.std_err }}</el-descriptions-item>
       <el-descriptions-item label="自动重启">{{ app.auto_restart ? '是' : '否' }}</el-descriptions-item>
       <el-descriptions-item label="监控">{{ app.is_monitor ? '是' : '否' }}</el-descriptions-item>
       <el-descriptions-item label="创建时间" :span="2">{{ formatTime(app.created_at) }}</el-descriptions-item>
     </el-descriptions>

     <div class="page-toolbar" style="margin-top: 16px">
       <el-button @click="router.push({ name: 'monitor', params: { kind: 'app', id: String(app.id) } })">
         监控
       </el-button>
       <el-button @click="router.push({ name: 'archive', params: { kind: 'app', id: String(app.id) } })">
         归档
       </el-button>
       <el-button @click="router.push({ name: 'log', params: { kind: 'app', id: String(app.id), type: 'out' } })">
         stdout
       </el-button>
       <el-button @click="router.push({ name: 'log', params: { kind: 'app', id: String(app.id), type: 'err' } })">
         stderr
       </el-button>
     </div>
   </div>
 </template>
