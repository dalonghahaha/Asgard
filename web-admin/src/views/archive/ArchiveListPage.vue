 <script setup lang="ts">
 import { onMounted, ref } from 'vue'
 import { useRoute } from 'vue-router'
 import { archiveApi } from '@/api'
 import type { Archive } from '@/types'
 import { formatTime } from '@/utils/format'

 // T-311 归档列表
 const route = useRoute()
 const kind = String(route.params.kind)
 const id = Number(route.params.id)
 const list = ref<Archive[]>([])
 const total = ref(0)
 const page = ref(1)
 const loading = ref(false)

 async function load() {
   loading.value = true
   try {
     if (kind === 'app') {
       const r = await archiveApi.app(id, page.value)
       list.value = r.list
       total.value = r.total
     } else if (kind === 'job') {
       const r = await archiveApi.job(id, page.value)
       list.value = r.list
       total.value = r.total
     } else if (kind === 'timing') {
       const r = await archiveApi.timing(id, page.value)
       list.value = r.list
       total.value = r.total
     }
   } finally {
     loading.value = false
   }
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">归档列表 - {{ kind }} #{{ id }}</span>
     </div>
     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="uuid" label="UUID" min-width="200" />
       <el-table-column prop="pid" label="PID" width="80" />
       <el-table-column label="开始" width="170">
         <template #default="{ row }">{{ formatTime(row.begin_time) }}</template>
       </el-table-column>
       <el-table-column label="结束" width="170">
         <template #default="{ row }">{{ formatTime(row.end_time) }}</template>
       </el-table-column>
       <el-table-column prop="status" label="退出码" width="80" />
       <el-table-column prop="signal" label="信号" />
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
   </div>
 </template>
