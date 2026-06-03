 <script setup lang="ts">
 defineProps<{ title: string; message: string; backTo?: string }>()
 </script>

 <template>
   <el-result icon="warning" :title="title" :sub-title="message">
     <template #extra>
       <el-button type="primary" @click="$router.push(backTo || '/')">返回首页</el-button>
     </template>
   </el-result>
 </template>

 <style scoped>
 .el-result {
   padding: 80px 0;
 }
 </style>
