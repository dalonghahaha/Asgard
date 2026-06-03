 import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
 import { useAuthStore } from '@/stores/auth'

 const MainLayout = () => import('@/layouts/MainLayout.vue')

 // T-204 路由表：与后端 web/router.go 11 个 group 一一对应
 //  - meta.title: 顶栏/面包屑标题
 //  - meta.requiresAuth: true 时走鉴权守卫
 //  - meta.admin: true 时仅管理员可见
 const routes: RouteRecordRaw[] = [
   {
     path: '/login',
     name: 'login',
     component: () => import('@/views/auth/LoginPage.vue'),
     meta: { title: '登录', requiresAuth: false },
   },
   {
     path: '/',
     component: MainLayout,
     redirect: { name: 'dashboard' },
     meta: { requiresAuth: true },
     children: [
       {
         path: 'dashboard',
         name: 'dashboard',
         component: () => import('@/views/DashboardPage.vue'),
         meta: { title: '首页', icon: 'Odometer' },
       },
       {
         path: 'users',
         name: 'users',
         component: () => import('@/views/users/UserListPage.vue'),
         meta: { title: '用户管理', icon: 'User', admin: true },
       },
       {
         path: 'agents',
         name: 'agents',
         component: () => import('@/views/agents/AgentListPage.vue'),
         meta: { title: '实例管理', icon: 'Cpu' },
       },
       {
         path: 'groups',
         name: 'groups',
         component: () => import('@/views/groups/GroupListPage.vue'),
         meta: { title: '分组管理', icon: 'Folder' },
       },
       {
         path: 'apps',
         name: 'apps',
         component: () => import('@/views/apps/AppListPage.vue'),
         meta: { title: '应用管理', icon: 'Box' },
       },
       {
         path: 'apps/:id',
         name: 'app-show',
         component: () => import('@/views/apps/AppShowPage.vue'),
         meta: { title: '应用详情', hidden: true },
       },
       {
         path: 'jobs',
         name: 'jobs',
         component: () => import('@/views/jobs/JobListPage.vue'),
         meta: { title: '计划任务', icon: 'AlarmClock' },
       },
       {
         path: 'timings',
         name: 'timings',
         component: () => import('@/views/timings/TimingListPage.vue'),
         meta: { title: '定时任务', icon: 'Timer' },
       },
       {
         path: 'monitor/:kind/:id',
         name: 'monitor',
         component: () => import('@/views/monitor/MonitorPage.vue'),
         meta: { title: '监控', hidden: true },
       },
       {
         path: 'archive/:kind/:id',
         name: 'archive',
         component: () => import('@/views/archive/ArchiveListPage.vue'),
         meta: { title: '归档', hidden: true },
       },
       {
         path: 'logs/:kind/:id/:type(out|err)',
         name: 'log',
         component: () => import('@/views/logs/LogPage.vue'),
         meta: { title: '日志', hidden: true },
       },
       {
         path: 'exceptions',
         name: 'exceptions',
         component: () => import('@/views/exception/ExceptionListPage.vue'),
         meta: { title: '异常记录', icon: 'Warning' },
       },
       {
         path: 'operations',
         name: 'operations',
         component: () => import('@/views/operation/OperationListPage.vue'),
         meta: { title: '操作日志', icon: 'Document' },
       },
       {
         path: 'profile',
         name: 'profile',
         component: () => import('@/views/users/ProfilePage.vue'),
         meta: { title: '个人设置', hidden: true },
       },
       {
         path: 'change-password',
         name: 'change-password',
         component: () => import('@/views/users/ChangePasswordPage.vue'),
         meta: { title: '修改密码', hidden: true },
       },
     ],
   },
   {
     path: '/no-login',
     name: 'no-login',
     component: () => import('@/views/error/NoLoginPage.vue'),
     meta: { title: '未登录' },
   },
   {
     path: '/auth-fail',
     name: 'auth-fail',
     component: () => import('@/views/error/AuthFailPage.vue'),
     meta: { title: '登录失效' },
   },
   {
     path: '/admin-only',
     name: 'admin-only',
     component: () => import('@/views/error/AdminOnlyPage.vue'),
     meta: { title: '权限不足' },
   },
   {
     path: '/forbidden',
     name: 'forbidden',
     component: () => import('@/views/error/ForbiddenPage.vue'),
     meta: { title: '已被禁用' },
   },
   {
     path: '/:pathMatch(.*)*',
     name: 'not-found',
     component: () => import('@/views/error/NotFoundPage.vue'),
     meta: { title: '页面不存在' },
   },
 ]

 const router = createRouter({
   history: createWebHistory(),
   routes,
 })

 router.beforeEach((to) => {
   const auth = useAuthStore()
   if (to.meta.requiresAuth !== false && !auth.isLogin) {
     return { name: 'login', query: { redirect: to.fullPath } }
   }
   if (to.meta.admin && !auth.isAdmin) {
     return { name: 'admin-only' }
   }
   return undefined
 })

 router.afterEach((to) => {
   const base = 'Asgard 管理控制台'
   document.title = to.meta.title ? `${to.meta.title as string} - ${base}` : base
 })

 export default router
