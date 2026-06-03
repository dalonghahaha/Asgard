 <script setup lang="ts">
 // 应用根节点：路由决定渲染哪个 layout/page
 </script>

 <template>
   <el-config-provider>
     <router-view />
   </el-config-provider>
 </template>

 <style lang="scss">
 html,
 body,
 #app {
   height: 100%;
   margin: 0;
   padding: 0;
 }
 </style>
