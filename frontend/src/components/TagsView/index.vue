<template>
  <div class="tags-view-container">
    <el-scrollbar ref="scrollPaneRef" class="tags-view-wrapper" @scroll="handleScroll">
      <router-link
        v-for="tag in visitedViews"
        :key="tag.path"
        :data-path="tag.path"
        :class="isActive(tag) ? 'active' : ''"
        :to="{ path: tag.path, query: tag.query, fullPath: tag.fullPath }"
        class="tags-view-item"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ tag.title }}
        <el-icon v-if="!isAffix(tag)" class="close-icon" @click.prevent.stop="closeSelectedTag(tag)">
          <Close />
        </el-icon>
      </router-link>
    </el-scrollbar>
    <ul v-show="visible" :style="{ left: left + 'px', top: top + 'px' }" class="contextmenu">
      <li @click="refreshSelectedTag(selectedTag)">刷新</li>
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)">关闭</li>
      <li @click="closeOthersTags">关闭其他</li>
      <li @click="closeAllTags(selectedTag)">关闭所有</li>
    </ul>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTagsViewStore } from '@/stores/tagsView'
import { Close } from '@element-plus/icons-vue'

const visible = ref(false)
const top = ref(0)
const left = ref(0)
const selectedTag = ref({})
const affixTags = ref([])
const scrollPaneRef = ref(null)

const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()

const visitedViews = computed(() => tagsViewStore.visitedViews)
const routes = computed(() => router.getRoutes())

watch(route, () => {
  addTags()
  moveToCurrentTag()
})

watch(visible, (value) => {
  if (value) {
    document.body.addEventListener('click', closeMenu)
  } else {
    document.body.removeEventListener('click', closeMenu)
  }
})

onMounted(() => {
  initTags()
  addTags()
})

function isActive(tag) {
  return tag.path === route.path
}

function isAffix(tag) {
  return tag.meta && tag.meta.affix
}

function filterAffixTags(routes, basePath = '/') {
  let tags = []
  routes.forEach(route => {
    if (route.meta && route.meta.affix) {
      const tagPath = basePath + route.path // handle nested paths if needed
      tags.push({
        fullPath: tagPath,
        path: tagPath,
        name: route.name,
        meta: { ...route.meta }
      })
    }
    if (route.children) {
      const tempTags = filterAffixTags(route.children, basePath + route.path)
      if (tempTags.length >= 1) {
        tags = [...tags, ...tempTags]
      }
    }
  })
  return tags
}

function initTags() {
  // simplified affix tags init
  const affixTags = [] 
  // You can traverse router routes to find affix tags
  for (const tag of affixTags) {
    if (tag.name) {
      tagsViewStore.addVisitedView(tag)
    }
  }
}

function addTags() {
  const { name } = route
  if (name && name !== 'Login') {
    tagsViewStore.addView(route)
  }
  return false
}

function moveToCurrentTag() {
  nextTick(() => {
    // Scroll handling logic can be added here
  })
}

function refreshSelectedTag(view) {
  tagsViewStore.delCachedView(view).then(() => {
    const { fullPath } = view
    nextTick(() => {
      router.replace({
        path: '/redirect' + fullPath
      }).catch((err) => {
          // Just reload normally if redirect route doesn't exist
           router.go(0)
      })
    })
  })
}

function closeSelectedTag(view) {
  tagsViewStore.delView(view).then(({ visitedViews }) => {
    if (isActive(view)) {
      toLastView(visitedViews, view)
    }
  })
}

function closeOthersTags() {
  router.push(selectedTag.value)
  tagsViewStore.delOthersViews(selectedTag.value).then(() => {
    moveToCurrentTag()
  })
}

function closeAllTags(view) {
  tagsViewStore.delAllViews().then(({ visitedViews }) => {
    if (affixTags.value.some(tag => tag.path === view.path)) {
      return
    }
    toLastView(visitedViews, view)
  })
}

function toLastView(visitedViews, view) {
  const latestView = visitedViews.slice(-1)[0]
  if (latestView) {
    router.push(latestView.fullPath)
  } else {
    // now the default is to redirect to the home page if there is no tags-view,
    // you can adjust it according to your needs.
    if (view.name === 'Dashboard') {
      // to reload home page
      router.replace({ path: '/redirect' + view.fullPath })
    } else {
      router.push('/')
    }
  }
}

function openMenu(tag, e) {
  const menuMinWidth = 105
  const offsetLeft = e.target.getBoundingClientRect().left // container margin left
  const offsetWidth = e.target.closest('.tags-view-container').offsetWidth // container width
  const maxLeft = offsetWidth - menuMinWidth // left boundary
  const leftPos = e.clientX - offsetLeft + 15 // 15: margin right

  if (leftPos > maxLeft) {
    left.value = maxLeft
  } else {
    left.value = leftPos
  }

  // Adjust top based on container position, or simple use clientY
  // For simplicity, just using clientX/Y but relative to container would be better
  // Let's use fixed positioning relative to viewport or absolute to container
  // For context menu, absolute position relative to body or container is common.
  // Here e.clientX is viewport.
  left.value = e.clientX
  top.value = e.clientY
  
  visible.value = true
  selectedTag.value = tag
}

function closeMenu() {
  visible.value = false
}

function handleScroll() {
  closeMenu()
}
</script>

<style scoped>
.tags-view-container {
  height: 34px;
  width: 100%;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);
}

.tags-view-wrapper {
  .tags-view-item {
    display: inline-block;
    position: relative;
    cursor: pointer;
    height: 26px;
    line-height: 26px;
    border: 1px solid #d8dce5;
    color: #495060;
    background: #fff;
    padding: 0 8px;
    font-size: 12px;
    margin-left: 5px;
    margin-top: 4px;
    &:first-of-type {
      margin-left: 15px;
    }
    &:last-of-type {
      margin-right: 15px;
    }
    &.active {
      background-color: #409eff;
      color: #fff;
      border-color: #409eff;
      &::before {
        content: '';
        background: #fff;
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        position: relative;
        margin-right: 2px;
      }
    }
  }
}

.contextmenu {
  margin: 0;
  background: #fff;
  z-index: 3000;
  position: fixed; /* Fixed to viewport */
  list-style-type: none;
  padding: 5px 0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  color: #333;
  box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, .3);
  li {
    margin: 0;
    padding: 7px 16px;
    cursor: pointer;
    &:hover {
      background: #eee;
    }
  }
}

.close-icon {
    width: 12px;
    height: 12px;
    vertical-align: -2px;
    border-radius: 50%;
    text-align: center;
    transition: all .3s cubic-bezier(.645, .045, .355, 1);
    transform-origin: 100% 50%;
    margin-left: 3px;
    &:hover {
      background-color: #b4bccc;
      color: #fff;
    }
}
</style>
