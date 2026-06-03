 <script setup lang="ts">
 import { onMounted, reactive, ref } from 'vue'
 import { operationApi } from '@/api'
 import type { OperationRow } from '@/types'
 import { formatTime } from '@/utils/format'

 // T-314 操作日志
 const list = ref<OperationRow[]>([])
 const total = ref(0)
 const loading = ref(false)
 const filter = reactive({ page: 1, user_id: 0, type: 0 })

 async function load() {
   loading.value = true
   try {
     const r = await operationApi.list({
       page: filter.page,
       user_id: filter.user_id || undefined,
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
       <span class="page-title">操作日志</span>
       <el-space>
         <el-input
           v-model.number="filter.user_id"
           placeholder="用户ID"
           clearable
           style="width: 120px"
         />
         <el-select v-model="filter.type" placeholder="对象类型" style="width: 160px" @change="load">
           <el-option label="全部" :value="0" />
           <el-option label="实例" :value="1" />
           <el-option label="应用" :value="2" />
           <el-option label="计划任务" :value="3" />
           <el-option label="定时任务" :value="4" />
           <el-option label="分组" :value="5" />
           <el-option label="用户" :value="6" />
         </el-select>
         <el-button type="primary" @click="load">查询</el-button>
       </el-space>
     </div>

     <el-table v-loading="loading" :data="list" border>
       <el-table-column prop="id" label="ID" width="80" />
       <el-table-column prop="desc" label="操作" />
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
