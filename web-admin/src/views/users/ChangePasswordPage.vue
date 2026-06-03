 <script setup lang="ts">
 import { ref } from 'vue'
 import { useRouter } from 'vue-router'
 import { ElMessage } from 'element-plus'
 import { useAuthStore } from '@/stores/auth'

 const auth = useAuthStore()
 const router = useRouter()
 const password = ref('')
 const loading = ref(false)

 async function submit() {
   if (!password.value) {
     ElMessage.warning('请输入新密码')
     return
   }
   loading.value = true
   try {
     await auth.changePassword(password.value)
     ElMessage.success('密码已修改，请重新登录')
     router.push({ name: 'login' })
   } finally {
     loading.value = false
   }
 }
 </script>

 <template>
   <div class="asgard-page" style="max-width: 480px">
     <el-card>
       <template #header>修改密码</template>
       <el-form label-width="80px">
         <el-form-item label="新密码">
           <el-input v-model="password" type="password" show-password />
         </el-form-item>
         <el-button type="primary" :loading="loading" @click="submit">提交</el-button>
       </el-form>
     </el-card>
   </div>
 </template>
