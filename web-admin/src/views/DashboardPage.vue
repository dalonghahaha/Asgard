 <script setup lang="ts">
 import { ref } from 'vue'
 import { agentApi, appApi, jobApi, timingApi, monitorApi } from '@/api'
 import type { Agent, MonitorPoint } from '@/types'
 import MonitorChart from '@/components/MonitorChart.vue'

 // T-303 仪表盘：四类计数 + agent CPU/内存图表
 const agents = ref<Agent[]>([])
 const counts = ref({ apps: 0, jobs: 0, timings: 0 })
 const selectedAgent = ref<number>(0)
 const points = ref<MonitorPoint[]>([])

 async function load() {
   const [a, ac, jc, tc] = await Promise.all([
     agentApi.list({ page: 1 }).then((r) => r.list),
     appApi.list({ page: 1 }).then((r) => r.total),
     jobApi.list({ page: 1 }).then((r) => r.total),
     timingApi.list({ page: 1 }).then((r) => r.total),
   ])
   agents.value = (a as Agent[]).filter((x: Agent) => x.status === 1)
   counts.value = { apps: ac, jobs: jc, timings: tc }
   if (agents.value.length > 0 && agents.value[0]) {
     await selectAgent(agents.value[0].id)
   }
 }

async function selectAgent(id: number) {
  selectedAgent.value = id
  points.value = (await monitorApi.agent(id, 50)) as MonitorPoint[]
}

 load()
 </script>

 <template>
   <div class="dashboard">
     <el-row :gutter="16">
       <el-col :span="6">
         <el-card>
           <div class="metric">
             <div class="label">实例</div>
             <div class="value">{{ agents.length }}</div>
           </div>
         </el-card>
       </el-col>
       <el-col :span="6">
         <el-card>
           <div class="metric">
             <div class="label">应用</div>
             <div class="value">{{ counts.apps }}</div>
           </div>
         </el-card>
       </el-col>
       <el-col :span="6">
         <el-card>
           <div class="metric">
             <div class="label">计划任务</div>
             <div class="value">{{ counts.jobs }}</div>
           </div>
         </el-card>
       </el-col>
       <el-col :span="6">
         <el-card>
           <div class="metric">
             <div class="label">定时任务</div>
             <div class="value">{{ counts.timings }}</div>
           </div>
         </el-card>
       </el-col>
     </el-row>

     <el-card class="chart-card">
       <template #header>
         <div class="chart-header">
           <span>实例 CPU/内存</span>
           <el-select
             v-model="selectedAgent"
             placeholder="选择实例"
             style="width: 240px"
             @change="selectAgent"
           >
             <el-option
               v-for="a in agents"
               :key="a.id"
               :label="a.alias || `${a.ip}:${a.port}`"
               :value="a.id"
             />
           </el-select>
         </div>
       </template>
       <MonitorChart :points="points" height="360px" />
     </el-card>
   </div>
 </template>

 <style lang="scss" scoped>
 .dashboard {
   display: flex;
   flex-direction: column;
   gap: 16px;
 }
 .metric {
   text-align: center;
   .label {
     font-size: 14px;
     color: #909399;
     margin-bottom: 8px;
   }
   .value {
     font-size: 32px;
     font-weight: 600;
     color: #1f6feb;
   }
 }
 .chart-card {
   .chart-header {
     display: flex;
     justify-content: space-between;
     align-items: center;
   }
 }
 </style>
