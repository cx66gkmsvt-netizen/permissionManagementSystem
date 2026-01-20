<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
          <div class="stat-icon">
            <el-icon :size="40"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.userCount }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)">
          <div class="stat-icon">
            <el-icon :size="40"><UserFilled /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.roleCount }}</div>
            <div class="stat-label">角色总数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)">
          <div class="stat-icon">
            <el-icon :size="40"><OfficeBuilding /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.deptCount }}</div>
            <div class="stat-label">部门总数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)">
          <div class="stat-icon">
            <el-icon :size="40"><Menu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.menuCount }}</div>
            <div class="stat-label">菜单总数</div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <div class="page-card" style="margin-top: 20px">
      <h3 style="margin-bottom: 20px">欢迎使用用户中心管理系统</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="当前用户">
          {{ userStore.userInfo?.nickName || userStore.userInfo?.userName }}
        </el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag v-for="role in userStore.roles" :key="role" type="primary" style="margin-right: 5px">
            {{ role }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属部门">
          {{ userStore.userInfo?.dept?.deptName || '未分配' }}
        </el-descriptions-item>
        <el-descriptions-item label="登录时间">
          {{ formatDate(userStore.userInfo?.loginDate) }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
    
    <div class="page-card">
      <h3 style="margin-bottom: 20px">系统功能</h3>
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="feature-card" @click="$router.push('/system/user')">
            <el-icon :size="32" color="#667eea"><User /></el-icon>
            <div class="feature-title">用户管理</div>
            <div class="feature-desc">管理系统用户</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="feature-card" @click="$router.push('/system/role')">
            <el-icon :size="32" color="#f5576c"><UserFilled /></el-icon>
            <div class="feature-title">角色管理</div>
            <div class="feature-desc">配置角色权限</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="feature-card" @click="$router.push('/system/menu')">
            <el-icon :size="32" color="#4facfe"><Menu /></el-icon>
            <div class="feature-title">菜单管理</div>
            <div class="feature-desc">管理系统菜单</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="feature-card" @click="$router.push('/system/dept')">
            <el-icon :size="32" color="#43e97b"><OfficeBuilding /></el-icon>
            <div class="feature-title">部门管理</div>
            <div class="feature-desc">管理组织架构</div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getStats } from '@/api/stats'

const userStore = useUserStore()

const stats = ref({
  userCount: 0,
  roleCount: 0,
  deptCount: 0,
  menuCount: 0
})

onMounted(async () => {
  try {
    const res = await getStats()
    stats.value = res.data
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
})

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  padding: 20px;
  border-radius: 12px;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 32px;
  font-weight: 600;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
}

.feature-card {
  padding: 30px;
  text-align: center;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.feature-card:hover {
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transform: translateY(-5px);
}

.feature-title {
  font-size: 16px;
  font-weight: 500;
  margin: 12px 0 8px;
}

.feature-desc {
  color: #909399;
  font-size: 13px;
}
</style>
