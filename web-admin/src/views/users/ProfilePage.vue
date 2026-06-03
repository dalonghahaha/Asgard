 <script setup lang="ts">
 import { reactive } from 'vue'
 import { ElMessage } from 'element-plus'
 import { useAuthStore } from '@/stores/auth'

 // 个人设置
 const auth = useAuthStore()
 const form = reactive({
   nickname: auth.user?.nickname || '',
   email: auth.user?.email || '',
   mobile: auth.user?.mobile || '',
 })

 function save() {
   // 后端目前没有专门的 profile 更新接口，复用 userApi.update
   // 这里只更新本地 store，避免在没接口时假报成功
   if (auth.user) {
     auth.user.nickname = form.nickname
     auth.user.email = form.email
     auth.user.mobile = form.mobile
   }
   ElMessage.success('已保存（本地预览）')
 }
 </script>

 <template>
   <div class="asgard-page" style="max-width: 480px">
     <el-card>
       <template #header>个人设置</template>
       <el-form label-width="80px">
         <el-form-item label="昵称">
           <el-input v-model="form.nickname" />
         </el-form-item>
         <el-form-item label="邮箱">
           <el-input v-model="form.email" />
         </el-form-item>
         <el-form-item label="手机号">
           <el-input v-model="form.mobile" />
         </el-form-item>
         <el-button type="primary" @click="save">保存</el-button>
       </el-form>
     </el-card>
   </div>
 </template>
