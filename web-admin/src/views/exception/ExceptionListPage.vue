 <script setup lang="ts">
 import { onMounted, reactive, ref } from 'vue'
 import { exceptionApi } from '@/api'
 import type { ExceptionRow } from '@/types'
 import { formatTime } from '@/utils/format'

 // T-313 异常记录
 const list = ref<ExceptionRow[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ page: 1, type: 0 })

 async function load() {
   loading.value = true
   try {
     const r = await exceptionApi.list({
       page: filter.page,
       type: filter.type || undefined,
     })
     list.value = r.list
     total.value = r.total
   } finally {
     loading.value = false
   }
 }

 onMounted(load)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">异常记录</span>
       <el-select v-model="filter.type" placeholder="类型" style="width: 200px" @change="load">
         <el-option label="全部" :value="0" />
         <el-option label="实例" :value="1" />
         <el-option label="应用" :value="2" />
         <el-option label="计划任务" :value="3" />
         <el-option label="定时任务" :value="4" />
         <el-option label="分组" :value="5" />
         <el-option label="用户" :value="6" />
       </el-select>
     </div>

     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="type" label="类型" width="120" />
       <el-table-column prop="name" label="对象" />
       <el-table-column prop="desc" label="描述" />
       <el-table-column label="时间" width="170">
         <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
       </el-table-column>
     </el-table>

     <div class="pagination-bar">
       <el-pagination
         v-model:current-page="filter.page"
         :total="total"
         :page-size="20"
         layout="prev, pager, next, total"
         @current-change="load"
       />
     </div>
   </div>
 </template>
