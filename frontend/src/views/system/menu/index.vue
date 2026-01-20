<template>
  <div class="page-card">
    <!-- 工具栏 -->
    <div class="table-actions">
      <div>
        <el-button v-permission="'system:menu:add'" type="primary" icon="Plus" @click="handleAdd()">新增</el-button>
        <el-button type="info" icon="Sort" @click="toggleExpandAll">展开/折叠</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      v-if="refreshTable"
      v-loading="loading"
      :data="menuList"
      row-key="menuId"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children' }"
      border
    >
      <el-table-column prop="menuName" label="菜单名称" width="200" />
      <el-table-column prop="icon" label="图标" width="80" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.icon">
            <component :is="row.icon" />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="menuType" label="类型" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.menuType === 'M'" type="primary">目录</el-tag>
          <el-tag v-else-if="row.menuType === 'C'" type="success">菜单</el-tag>
          <el-tag v-else type="warning">按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column prop="perms" label="权限标识" width="200" />
      <el-table-column prop="path" label="路由地址" width="150" />
      <el-table-column prop="component" label="组件路径" width="200" />
      <el-table-column prop="visible" label="可见" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.visible === '0' ? 'success' : 'info'">
            {{ row.visible === '0' ? '显示' : '隐藏' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'system:menu:edit'" type="primary" link icon="Edit" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button v-permission="'system:menu:add'" type="success" link icon="Plus" @click="handleAdd(row)">
            新增
          </el-button>
          <el-button v-permission="'system:menu:remove'" type="danger" link icon="Delete" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="上级菜单">
              <el-tree-select
                v-model="form.parentId"
                :data="menuOptions"
                :props="{ label: 'menuName', value: 'menuId', children: 'children' }"
                placeholder="请选择上级菜单"
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="菜单类型" prop="menuType">
              <el-radio-group v-model="form.menuType">
                <el-radio value="M">目录</el-radio>
                <el-radio value="C">菜单</el-radio>
                <el-radio value="F">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType !== 'F'">
            <el-form-item label="菜单图标">
              <el-input v-model="form.icon" placeholder="请输入图标名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="form.menuName" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="form.menuType !== 'F'">
          <el-col :span="12">
            <el-form-item label="路由地址" prop="path">
              <el-input v-model="form.path" placeholder="请输入路由地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType === 'C'">
            <el-form-item label="组件路径" prop="component">
              <el-input v-model="form.component" placeholder="请输入组件路径" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12" v-if="form.menuType !== 'M'">
            <el-form-item label="权限标识" prop="perms">
              <el-input v-model="form.perms" placeholder="如: system:user:list" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType !== 'F'">
            <el-form-item label="显示状态" prop="visible">
              <el-radio-group v-model="form.visible">
                <el-radio value="0">显示</el-radio>
                <el-radio value="1">隐藏</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listMenu, getMenu, createMenu, updateMenu, deleteMenu } from '@/api/menu'

const loading = ref(false)
const menuList = ref([])
const menuOptions = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const isExpandAll = ref(true)
const refreshTable = ref(true)

const formRef = ref()

const form = reactive({
  menuId: null,
  parentId: 0,
  menuName: '',
  menuType: 'M',
  path: '',
  component: '',
  perms: '',
  icon: '',
  sort: 0,
  visible: '0',
  status: '0'
})

const rules = {
  menuName: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  menuType: [{ required: true, message: '请选择菜单类型', trigger: 'change' }]
}

onMounted(() => {
  getList()
})

const getList = async () => {
  loading.value = true
  try {
    const res = await listMenu()
    menuList.value = res.data || []
    menuOptions.value = [{ menuId: 0, menuName: '根目录', children: res.data }]
  } finally {
    loading.value = false
  }
}

const toggleExpandAll = () => {
  refreshTable.value = false
  isExpandAll.value = !isExpandAll.value
  setTimeout(() => {
    refreshTable.value = true
  }, 0)
}

const resetForm = () => {
  form.menuId = null
  form.parentId = 0
  form.menuName = ''
  form.menuType = 'M'
  form.path = ''
  form.component = ''
  form.perms = ''
  form.icon = ''
  form.sort = 0
  form.visible = '0'
  form.status = '0'
}

const handleAdd = (row) => {
  resetForm()
  if (row) {
    form.parentId = row.menuId
  }
  dialogTitle.value = '新增菜单'
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  resetForm()
  const res = await getMenu(row.menuId)
  const menu = res.data
  Object.assign(form, {
    menuId: menu.menuId,
    parentId: menu.parentId,
    menuName: menu.menuName,
    menuType: menu.menuType,
    path: menu.path,
    component: menu.component,
    perms: menu.perms,
    icon: menu.icon,
    sort: menu.sort,
    visible: menu.visible,
    status: menu.status
  })
  dialogTitle.value = '编辑菜单'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (form.menuId) {
      await updateMenu(form.menuId, form)
      ElMessage.success('更新成功')
    } else {
      await createMenu(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定删除该菜单吗？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteMenu(row.menuId)
      ElMessage.success('删除成功')
      getList()
    })
}
</script>
