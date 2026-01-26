<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <div class="layout-sidebar" :class="{ collapsed: isCollapsed }">
      <div class="logo">
        <span v-if="!isCollapsed">用户中心</span>
        <span v-else class="logo-collapsed">UC</span>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="dark-menu"
        :collapse="isCollapsed"
        router
        background-color="transparent"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <template v-for="menu in userStore.menuRoutes" :key="menu.menuId || menu.path">
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
            <template #title>
              <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
              <span>{{ menu.menuName }}</span>
            </template>
            <el-menu-item 
              v-for="child in menu.children" 
              :key="child.menuId || child.path" 
              :index="child.path.startsWith('/') ? child.path : menu.path + '/' + child.path"
            >
              <el-icon v-if="child.icon"><component :is="child.icon" /></el-icon>
              <template #title>{{ child.menuName }}</template>
            </el-menu-item>
          </el-sub-menu>

          <el-menu-item v-else :index="menu.path">
            <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
            <template #title>{{ menu.menuName }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </div>
    
    <!-- 主区域 -->
    <div class="layout-main">
      <!-- 头部 -->
      <div class="layout-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="toggleCollapse">
            <Expand v-if="isCollapsed" />
            <Fold v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.meta.title !== '首页'">
              {{ $route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" icon="UserFilled" />
              <span class="user-name">{{ userStore.userInfo?.nickName || userStore.userInfo?.userName }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <!-- 标签页 -->
      <TagsView />

      <!-- 内容区 -->
      <div class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <keep-alive :include="cachedViews">
              <component :is="Component" :key="$route.fullPath" />
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useTagsViewStore } from '@/stores/tagsView'
import TagsView from '@/components/TagsView/index.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const tagsViewStore = useTagsViewStore()

const isCollapsed = ref(false)
const cachedViews = computed(() => tagsViewStore.cachedViews)

const activeMenu = computed(() => route.path)

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

const handleCommand = (command) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
    })
  }
}
</script>

<style scoped>
.layout-container {
  display: flex;
  height: 100vh;
  width: 100vw;
}

.layout-sidebar {
  width: 220px;
  height: 100%;
  background-color: #304156;
  color: #fff;
  transition: width 0.3s;
  overflow-y: auto;
  flex-shrink: 0;
}

.layout-sidebar.collapsed {
  width: 64px;
}

.logo {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  background-color: #2b3648;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #f0f2f5;
}

.layout-header {
  height: 50px;
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.collapse-btn:hover {
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #606266;
}

.user-name {
  font-size: 14px;
}

.layout-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
