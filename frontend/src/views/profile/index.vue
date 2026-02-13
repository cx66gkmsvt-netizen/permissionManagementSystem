<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <div class="page-card user-card">
          <div class="user-avatar" @click="triggerUpload">
            <el-avatar :size="100" :src="avatarUrl" icon="UserFilled" />
            <div class="avatar-overlay">
              <el-icon :size="24"><Camera /></el-icon>
              <span>更换头像</span>
            </div>
            <input ref="fileInputRef" type="file" accept="image/jpg,image/jpeg,image/png,image/gif" style="display: none" @change="handleAvatarChange" />
          </div>
          <h3>{{ userStore.userInfo?.nickName || userStore.userInfo?.userName }}</h3>
          <p class="user-role">
            <el-tag v-for="role in userStore.roles" :key="role" type="primary" style="margin: 2px">
              {{ role }}
            </el-tag>
          </p>
          <el-descriptions :column="1" border style="margin-top: 20px">
            <el-descriptions-item label="用户名">{{ userStore.userInfo?.userName }}</el-descriptions-item>
            <el-descriptions-item label="所属部门">{{ userStore.userInfo?.dept?.deptName || '未分配' }}</el-descriptions-item>
            <el-descriptions-item label="登录时间">{{ formatDate(userStore.userInfo?.loginDate) }}</el-descriptions-item>
            <el-descriptions-item label="登录IP">{{ userStore.userInfo?.loginIp || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
      
      <el-col :span="16">
        <div class="page-card">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本信息" name="info">
              <el-form ref="infoFormRef" :model="infoForm" :rules="infoRules" label-width="100px" style="max-width: 500px">
                <el-form-item label="昵称" prop="nickName">
                  <el-input v-model="infoForm.nickName" placeholder="请输入昵称" />
                </el-form-item>
                <el-form-item label="手机号" prop="phone">
                  <el-input v-model="infoForm.phone" placeholder="请输入手机号" />
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="infoForm.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="infoLoading" @click="handleUpdateInfo">保存修改</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            
            <el-tab-pane label="修改密码" name="password">
              <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="100px" style="max-width: 500px">
                <el-form-item label="旧密码" prop="oldPassword">
                  <el-input v-model="pwdForm.oldPassword" type="password" placeholder="请输入旧密码" show-password />
                </el-form-item>
                <el-form-item label="新密码" prop="newPassword">
                  <el-input v-model="pwdForm.newPassword" type="password" placeholder="请输入新密码" show-password />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirmPassword">
                  <el-input v-model="pwdForm.confirmPassword" type="password" placeholder="请确认新密码" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="pwdLoading" @click="handleUpdatePassword">修改密码</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Camera } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword, uploadAvatar } from '@/api/profile'

const userStore = useUserStore()

const activeTab = ref('info')
const infoLoading = ref(false)
const pwdLoading = ref(false)
const fileInputRef = ref()

const avatarUrl = computed(() => {
  return userStore.userInfo?.avatar || ''
})

const triggerUpload = () => {
  fileInputRef.value?.click()
}

const handleAvatarChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning('头像文件不能超过2MB')
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  try {
    await uploadAvatar(formData)
    ElMessage.success('头像更新成功')
    await userStore.fetchUserInfo()
  } catch {
    ElMessage.error('上传失败')
  } finally {
    fileInputRef.value.value = ''
  }
}

const infoFormRef = ref()
const pwdFormRef = ref()

const infoForm = reactive({
  nickName: '',
  phone: '',
  email: ''
})

const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const infoRules = {
  nickName: [{ required: true, message: '请输入昵称', trigger: 'blur' }]
}

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== pwdForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const pwdRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

onMounted(() => {
  if (userStore.userInfo) {
    infoForm.nickName = userStore.userInfo.nickName || ''
    infoForm.phone = userStore.userInfo.phone || ''
    infoForm.email = userStore.userInfo.email || ''
  }
})

const handleUpdateInfo = async () => {
  await infoFormRef.value.validate()
  infoLoading.value = true
  try {
    await updateProfile(infoForm)
    ElMessage.success('修改成功')
    // 刷新用户信息
    await userStore.fetchUserInfo()
  } finally {
    infoLoading.value = false
  }
}

const handleUpdatePassword = async () => {
  await pwdFormRef.value.validate()
  pwdLoading.value = true
  try {
    await updatePassword({
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword
    })
    ElMessage.success('密码修改成功')
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
  } finally {
    pwdLoading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped>
.profile-container {
  padding: 0;
}

.user-card {
  text-align: center;
}

.user-avatar {
  margin-bottom: 16px;
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.3s;
}

.user-avatar:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay span {
  margin-top: 4px;
}

.user-card h3 {
  margin: 0 0 8px;
  font-size: 20px;
}

.user-role {
  margin: 0 0 20px;
}
</style>
