<template>
  <div class="page-card">
    <!-- 搜索表单 -->
    <el-form :inline="true" :model="queryParams" class="table-actions">
      <div>
        <el-form-item label="角色名称">
          <el-input v-model="queryParams.roleName" placeholder="请输入角色名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="权限字符">
          <el-input v-model="queryParams.roleKey" placeholder="请输入权限字符" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </div>
      <div>
        <el-button v-permission="'system:role:add'" type="primary" icon="Plus" @click="handleAdd">新增</el-button>
      </div>
    </el-form>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="roleList" border>
      <el-table-column type="index" label="序号" width="60" align="center" />
      <el-table-column prop="roleName" label="角色名称" width="150" />
      <el-table-column prop="roleKey" label="权限字符" width="150" />
      <el-table-column prop="dataScope" label="数据范围" width="150">
        <template #default="{ row }">
          <el-tag :type="getDataScopeType(row.dataScope)">
            {{ getDataScopeLabel(row.dataScope) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === '0' ? 'success' : 'danger'">
            {{ row.status === '0' ? '正常' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.createTime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'system:role:edit'" type="primary" link icon="Edit" @click="handleEdit(row)" :disabled="row.roleKey === 'admin'">
            编辑
          </el-button>
          <el-button v-permission="'system:role:remove'" type="danger" link icon="Delete" @click="handleDelete(row)" :disabled="row.roleKey === 'admin'">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="queryParams.pageNum"
      v-model:page-size="queryParams.pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 20px; justify-content: flex-end"
      @size-change="getList"
      @current-change="getList"
    />

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色名称" prop="roleName">
              <el-input v-model="form.roleName" placeholder="请输入角色名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限字符" prop="roleKey">
              <el-input v-model="form.roleKey" placeholder="请输入权限字符" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio value="0">正常</el-radio>
                <el-radio value="1">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="数据范围" prop="dataScope">
              <el-select v-model="form.dataScope" placeholder="请选择" style="width: 100%">
                <el-option label="全部数据" value="1" />
                <el-option label="自定义数据" value="2" />
                <el-option label="本部门及以下" value="3" />
                <el-option label="仅本部门" value="4" />
                <el-option label="仅本人" value="5" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="菜单权限" prop="menuIds">
          <div class="tree-border">
            <el-tree
              ref="menuTreeRef"
              :data="menuOptions"
              :props="{ label: 'menuName', children: 'children' }"
              node-key="menuId"
              show-checkbox
              default-expand-all
            />
          </div>
        </el-form-item>
        <el-form-item v-if="form.dataScope === '2'" label="数据权限" prop="deptIds">
          <div class="tree-border">
            <el-tree
              ref="deptTreeRef"
              :data="deptOptions"
              :props="{ label: 'deptName', children: 'children' }"
              node-key="deptId"
              show-checkbox
              default-expand-all
            />
          </div>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listRole, getRole, createRole, updateRole, deleteRole } from '@/api/role'
import { listMenu } from '@/api/menu'
import { listDept } from '@/api/dept'

const loading = ref(false)
const roleList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const menuOptions = ref([])
const deptOptions = ref([])

const formRef = ref()
const menuTreeRef = ref()
const deptTreeRef = ref()

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  roleName: '',
  roleKey: ''
})

const form = reactive({
  roleId: null,
  roleName: '',
  roleKey: '',
  dataScope: '1',
  sort: 0,
  status: '0',
  remark: '',
  menuIds: [],
  deptIds: []
})

const rules = {
  roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  roleKey: [{ required: true, message: '请输入权限字符', trigger: 'blur' }]
}

const dataScopeOptions = {
  '1': { label: '全部数据', type: 'primary' },
  '2': { label: '自定义数据', type: 'warning' },
  '3': { label: '本部门及以下', type: 'success' },
  '4': { label: '仅本部门', type: 'info' },
  '5': { label: '仅本人', type: 'danger' }
}

const getDataScopeLabel = (scope) => dataScopeOptions[scope]?.label || '未知'
const getDataScopeType = (scope) => dataScopeOptions[scope]?.type || ''

onMounted(() => {
  getList()
  getMenus()
  getDepts()
})

const getList = async () => {
  loading.value = true
  try {
    const res = await listRole(queryParams)
    roleList.value = res.data.rows || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const getMenus = async () => {
  const res = await listMenu()
  menuOptions.value = res.data || []
}

const getDepts = async () => {
  const res = await listDept()
  deptOptions.value = res.data || []
}

const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

const resetQuery = () => {
  queryParams.roleName = ''
  queryParams.roleKey = ''
  handleQuery()
}

const resetForm = () => {
  form.roleId = null
  form.roleName = ''
  form.roleKey = ''
  form.dataScope = '1'
  form.sort = 0
  form.status = '0'
  form.remark = ''
  form.menuIds = []
  form.deptIds = []
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增角色'
  dialogVisible.value = true
  nextTick(() => {
    menuTreeRef.value?.setCheckedKeys([])
    deptTreeRef.value?.setCheckedKeys([])
  })
}

const handleEdit = async (row) => {
  resetForm()
  const res = await getRole(row.roleId)
  const role = res.data.role
  const menuIds = res.data.menuIds || []
  const deptIds = res.data.deptIds || []
  
  Object.assign(form, {
    roleId: role.roleId,
    roleName: role.roleName,
    roleKey: role.roleKey,
    dataScope: role.dataScope,
    sort: role.sort,
    status: role.status,
    remark: role.remark
  })
  
  dialogTitle.value = '编辑角色'
  dialogVisible.value = true
  
  nextTick(() => {
    menuTreeRef.value?.setCheckedKeys(menuIds)
    deptTreeRef.value?.setCheckedKeys(deptIds)
  })
}

const handleSubmit = async () => {
  await formRef.value.validate()
  
  // 获取选中的菜单和部门
  form.menuIds = menuTreeRef.value?.getCheckedKeys() || []
  form.deptIds = deptTreeRef.value?.getCheckedKeys() || []
  
  submitLoading.value = true
  try {
    if (form.roleId) {
      await updateRole(form.roleId, form)
      ElMessage.success('更新成功')
    } else {
      await createRole(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定删除该角色吗？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteRole(row.roleId)
      ElMessage.success('删除成功')
      getList()
    })
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>
