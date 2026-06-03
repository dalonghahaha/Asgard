 <script setup lang="ts">
 import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
 import { useRoute } from 'vue-router'
 import { monitorApi } from '@/api'
 import MonitorChart from '@/components/MonitorChart.vue'

 // T-310 监控：折线图 + 实时 SSE
 const route = useRoute()
 const kind = String(route.params.kind)
 const id = Number(route.params.id)

 const points = ref<{ cpu: number; memory: number; created_at: string }[]>([])
 let es: EventSource | null = null

 const title = computed(() => `监控 - ${kind} #${id}`)

 async function load() {
   if (kind === 'agent') points.value = await monitorApi.agent(id, 100)
   else if (kind === 'app') points.value = await monitorApi.app(id, 100)
   else if (kind === 'job') points.value = await monitorApi.job(id, 100)
   else if (kind === 'timing') points.value = await monitorApi.timing(id, 100)
 }

 function connect() {
   if (typeof EventSource === 'undefined') return
   const url = `/api/v1/sse/monitor/${kind}?${kind}_id=${id}&interval=5`
   es = new EventSource(url, { withCredentials: true })
   es.addEventListener('point', (ev) => {
     try {
       const data = JSON.parse((ev as MessageEvent).data)
       points.value.push(data)
       if (points.value.length > 200) points.value.shift()
     } catch {
       // ignore parse errors
     }
   })
 }

 onMounted(async () => {
   await load()
   connect()
 })

 onBeforeUnmount(() => {
   es?.close()
 })
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">{{ title }}</span>
       <el-button @click="load">刷新</el-button>
     </div>
     <el-card>
       <MonitorChart :points="points" height="500px" />
     </el-card>
   </div>
 </template>
