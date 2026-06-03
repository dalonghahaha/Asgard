 <script setup lang="ts">
 import { computed, h, onMounted } from 'vue'
 import type { Component } from 'vue'
 import { useRoute, useRouter } from 'vue-router'
 import {
   Odometer,
   User,
   Cpu,
   Folder,
   Box,
   AlarmClock,
   Timer,
   Warning,
   Document,
   SwitchButton,
   Setting,
 } from '@element-plus/icons-vue'
 import { useAuthStore } from '@/stores/auth'

 // T-302 主框架：侧边栏 + 顶栏 + 内容区
 const route = useRoute()
 const router = useRouter()
 const auth = useAuthStore()

 const ICONS: Record<string, Component> = {
   Odometer,
   User,
   Cpu,
   Folder,
   Box,
   AlarmClock,
   Timer,
   Warning,
   Document,
 }

 const menus = computed(() =>
   router
     .getRoutes()
     .filter(
       (r) =>
         r.name &&
         r.meta?.title &&
         !r.meta.hidden &&
         typeof r.path === 'string' &&
         r.path.startsWith('/') &&
         !r.path.includes(':'),
     )
     .filter((r) => !r.meta?.admin || auth.isAdmin)
     .map((r) => ({
       name: r.name as string,
       path: r.path,
       title: r.meta?.title as string,
       icon: (r.meta?.icon as string) || 'Odometer',
     })),
 )

 const activeMenu = computed(() => route.path)

 onMounted(() => {
   if (!auth.user && auth.token) {
     auth.fetchInfo().catch(() => auth.clear())
   }
 })

 async function onLogout() {
   await auth.logout()
   router.push({ name: 'login' })
 }

 function onCommand(cmd: string) {
   if (cmd === 'logout') return onLogout()
   if (cmd === 'profile') return router.push({ name: 'profile' })
   if (cmd === 'change-password') return router.push({ name: 'change-password' })
   return undefined
 }
 </script>

 <template>
   <el-container class="main-layout">
     <el-aside width="220px" class="sidebar">
       <div class="logo">Asgard</div>
       <el-menu
         :default-active="activeMenu"
         router
         background-color="#001529"
         text-color="#c0c4cc"
         active-text-color="#ffffff"
       >
         <el-menu-item v-for="m in menus" :key="m.name" :index="m.path">
           <el-icon><component :is="ICONS[m.icon] || Odometer" /></el-icon>
           <span>{{ m.title }}</span>
         </el-menu-item>
       </el-menu>
     </el-aside>

     <el-container>
       <el-header class="topbar">
         <div class="topbar-title">{{ route.meta?.title }}</div>
         <el-dropdown @command="onCommand">
           <span class="user-trigger">
             <el-icon><User /></el-icon>
             {{ auth.user?.nickname || '未登录' }}
             <el-tag v-if="auth.isAdmin" size="small" type="danger" effect="dark">管理员</el-tag>
           </span>
           <template #dropdown>
             <el-dropdown-menu>
               <el-dropdown-item command="profile">
                 <el-icon><Setting /></el-icon> 个人设置
               </el-dropdown-item>
               <el-dropdown-item command="change-password">
                 <el-icon><Document /></el-icon> 修改密码
               </el-dropdown-item>
               <el-dropdown-item divided command="logout">
                 <el-icon><SwitchButton /></el-icon> 退出
               </el-dropdown-item>
             </el-dropdown-menu>
           </template>
         </el-dropdown>
       </el-header>

       <el-main>
         <router-view />
       </el-main>
     </el-container>
   </el-container>
 </template>

 <style lang="scss" scoped>
 .main-layout {
   height: 100vh;
 }
 .sidebar {
   background: #001529;
   color: #fff;
   .logo {
     height: 60px;
     line-height: 60px;
     text-align: center;
     font-size: 20px;
     font-weight: bold;
     color: #fff;
     border-bottom: 1px solid #1f2d3d;
   }
   :deep(.el-menu) {
     border-right: none;
   }
 }
 .topbar {
   display: flex;
   justify-content: space-between;
   align-items: center;
   background: #fff;
   border-bottom: 1px solid #ebeef5;
   padding: 0 16px;
   .topbar-title {
     font-size: 16px;
     font-weight: 600;
   }
   .user-trigger {
     cursor: pointer;
     display: inline-flex;
     align-items: center;
     gap: 6px;
   }
 }
 </style>
