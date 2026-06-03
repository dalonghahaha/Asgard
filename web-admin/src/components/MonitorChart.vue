 <script setup lang="ts">
 import { computed, provide } from 'vue'
 import VChart, { THEME_KEY } from 'vue-echarts'
 import { use } from 'echarts/core'
 import { CanvasRenderer } from 'echarts/renderers'
 import { LineChart } from 'echarts/charts'
 import {
   GridComponent,
   TooltipComponent,
   LegendComponent,
   TitleComponent,
 } from 'echarts/components'
 import type { MonitorPoint } from '@/types'

 use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent])

 // T-208 ECharts 折线图：CPU + 内存双轴时间序列
 const props = defineProps<{
   points: MonitorPoint[]
   title?: string
   height?: string
 }>()

 const option = computed(() => ({
   title: props.title ? { text: props.title, left: 'center' } : undefined,
   tooltip: { trigger: 'axis' },
   legend: { data: ['CPU (%)', '内存 (MB)'], top: props.title ? 30 : 0 },
   grid: { top: props.title ? 70 : 40, left: 50, right: 50, bottom: 40 },
   xAxis: {
     type: 'category',
     data: props.points.map((p) => p.created_at),
     axisLabel: { rotate: 30, fontSize: 11 },
   },
   yAxis: [
     { type: 'value', name: 'CPU %', position: 'left' },
     { type: 'value', name: 'Memory MB', position: 'right' },
   ],
   series: [
     {
       name: 'CPU (%)',
       type: 'line',
       smooth: true,
       data: props.points.map((p) => Number(p.cpu.toFixed(2))),
       yAxisIndex: 0,
       showSymbol: false,
     },
     {
       name: '内存 (MB)',
       type: 'line',
       smooth: true,
       data: props.points.map((p) => Number(p.memory.toFixed(2))),
       yAxisIndex: 1,
       showSymbol: false,
     },
   ],
 }))

 provide(THEME_KEY, 'light')
 </script>

 <template>
   <div :style="{ width: '100%', height: height || '320px' }">
     <v-chart :option="option" autoresize />
   </div>
 </template>
