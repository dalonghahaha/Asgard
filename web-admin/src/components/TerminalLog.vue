 <script setup lang="ts">
 import { onMounted, onBeforeUnmount, ref, watch } from 'vue'
 import { Terminal } from 'xterm'
 import { FitAddon } from 'xterm-addon-fit'
 import 'xterm/css/xterm.css'

 // T-209 xterm.js 终端组件：用 EventSource 订阅 SSE 日志流，
 // 把每条 log 事件直接 write 进 xterm。
 const props = defineProps<{
   url: string
   height?: string
 }>()

 const elRef = ref<HTMLDivElement | null>(null)
 let term: Terminal | null = null
 let fit: FitAddon | null = null
 let es: EventSource | null = null

 onMounted(() => {
   if (!elRef.value) return
   term = new Terminal({
     convertEol: true,
     fontFamily: 'Menlo, Monaco, Consolas, monospace',
     fontSize: 12,
     theme: { background: '#1e1e1e' },
     cursorBlink: false,
     disableStdin: true,
   })
   fit = new FitAddon()
   term.loadAddon(fit)
   term.open(elRef.value)
   fit.fit()
   window.addEventListener('resize', resize)
   connect()
 })

 onBeforeUnmount(() => {
   window.removeEventListener('resize', resize)
   es?.close()
   term?.dispose()
   term = null
 })

 watch(
   () => props.url,
   () => {
     es?.close()
     term?.clear()
     connect()
   },
 )

 function resize() {
   fit?.fit()
 }

 function connect() {
   if (!props.url || typeof EventSource === 'undefined') return
   es = new EventSource(props.url, { withCredentials: true })
   es.addEventListener('log', (ev) => {
     term?.writeln((ev as MessageEvent).data)
   })
   es.addEventListener('error', (ev) => {
     term?.writeln('\x1b[31m[stream error]\x1b[0m')
     // EventSource 会自动重连，这里不主动重连
     void ev
   })
 }
 </script>

 <template>
   <div ref="elRef" class="terminal" :style="{ height: height || '480px' }" />
 </template>

 <style scoped>
 .terminal {
   width: 100%;
   background: #1e1e1e;
   padding: 8px;
   border-radius: 4px;
 }
 </style>
