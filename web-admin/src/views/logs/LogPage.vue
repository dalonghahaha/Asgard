 <script setup lang="ts">
 import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
 import { useRoute } from 'vue-router'
 import TerminalLog from '@/components/TerminalLog.vue'
 import { logApi } from '@/api'

 // T-312 日志查看：xterm 实时 + 历史拉取
 const route = useRoute()
 const kind = String(route.params.kind)
 const id = Number(route.params.id)
 const type = String(route.params.type) // out | err

 const sseUrl = computed(() => `/api/v1/sse/${type === 'err' ? 'err_log' : 'out_log'}/${kind}?${kind}_id=${id}&interval=3`)
 const history = ref<string[]>([])
 const loadingHistory = ref(false)
 let es: EventSource | null = null

 async function loadHistory() {
   loadingHistory.value = true
   try {
     if (kind === 'app') {
       const r = type === 'err' ? await logApi.appErr(id, 200) : await logApi.appOut(id, 200)
       history.value = r.content
     } else if (kind === 'job') {
       const r = type === 'err' ? await logApi.jobErr(id, 200) : await logApi.jobOut(id, 200)
       history.value = r.content
     } else if (kind === 'timing') {
       const r = type === 'err' ? await logApi.timingErr(id, 200) : await logApi.timingOut(id, 200)
       history.value = r.content
     }
   } finally {
     loadingHistory.value = false
   }
 }

 function startStream() {
   if (typeof EventSource === 'undefined') return
   es = new EventSource(sseUrl.value, { withCredentials: true })
 }

 function stopStream() {
   es?.close()
   es = null
 }

 onMounted(async () => {
   await loadHistory()
   startStream()
 })

 onBeforeUnmount(stopStream)
 </script>

 <template>
   <div class="asgard-page">
     <div class="page-header">
       <span class="page-title">日志 - {{ kind }} #{{ id }} ({{ type }})</span>
       <div>
         <el-button @click="loadHistory">拉取历史</el-button>
         <el-button @click="startStream" type="primary">开始实时</el-button>
         <el-button @click="stopStream">停止实时</el-button>
       </div>
     </div>

     <el-card v-loading="loadingHistory">
       <template #header>历史日志（最近 200 行）</template>
       <pre class="log-block">{{ history.join('\n') }}</pre>
     </el-card>

     <el-card style="margin-top: 16px">
       <template #header>实时日志（SSE）</template>
       <TerminalLog :url="sseUrl" height="360px" />
     </el-card>
   </div>
 </template>

 <style lang="scss" scoped>
 .log-block {
   background: #1e1e1e;
   color: #d4d4d4;
   padding: 12px;
   margin: 0;
   font-family: Menlo, Monaco, Consolas, monospace;
   font-size: 12px;
   white-space: pre-wrap;
   max-height: 360px;
   overflow: auto;
 }
 </style>
