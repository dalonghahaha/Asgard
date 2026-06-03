 <script setup lang="ts">
 import { reactive, ref } from 'vue'
 import { useRoute, useRouter } from 'vue-router'
 import { ElMessage } from 'element-plus'
 import { useAuthStore } from '@/stores/auth'

 // T-301 登录页：支持用户名/邮箱/手机号 + 密码
 const form = reactive({ username: '', password: '' })
 const loading = ref(false)
 const auth = useAuthStore()
 const router = useRouter()
 const route = useRoute()

 async function onSubmit() {
   if (!form.username || !form.password) {
     ElMessage.warning('请输入用户名和密码')
     return
   }
   loading.value = true
   try {
     await auth.login(form.username, form.password)
     ElMessage.success('登录成功')
     const redirect = (route.query.redirect as string) || '/'
     router.push(redirect)
   } finally {
     loading.value = false
   }
 }
 </script>

 <template>
   <div class="login-page">
     <el-card class="login-card">
       <h2>Asgard 管理控制台</h2>
       <el-form :model="form" label-position="top" @submit.prevent="onSubmit">
         <el-form-item label="用户名 / 邮箱 / 手机号">
           <el-input v-model="form.username" placeholder="请输入" autofocus />
         </el-form-item>
         <el-form-item label="密码">
           <el-input v-model="form.password" type="password" show-password placeholder="请输入" />
         </el-form-item>
         <el-button type="primary" :loading="loading" class="submit" @click="onSubmit">
           登录
         </el-button>
       </el-form>
     </el-card>
   </div>
 </template>

 <style lang="scss" scoped>
 .login-page {
   height: 100vh;
   display: flex;
   align-items: center;
   justify-content: center;
   background: linear-gradient(135deg, #1f6feb 0%, #001529 100%);
 }
 .login-card {
   width: 360px;
   h2 {
     text-align: center;
     margin-bottom: 24px;
   }
 }
 .submit {
   width: 100%;
 }
 </style>
