<template>
  <div class="page-card">
    <!-- 工具栏 -->
    <div class="table-actions">
      <div>
        <el-button v-permission="'system:dept:add'" type="primary" icon="Plus" @click="handleAdd()">新增</el-button>
        <el-button type="info" icon="Sort" @click="toggleExpandAll">展开/折叠</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      v-if="refreshTable"
      v-loading="loading"
      :data="deptList"
      row-key="deptId"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children' }"
      border
    >
      <el-table-column prop="deptName" label="部门名称" width="250" />
      <el-table-column prop="sort" label="排序" width="100" align="center" />
      <el-table-column prop="leader" label="负责人" width="150" />
      <el-table-column prop="phone" label="联系电话" width="150" />
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column prop="status" label="状态" width="100" align="center">
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
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'system:dept:edit'" type="primary" link icon="Edit" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button v-permission="'system:dept:add'" type="success" link icon="Plus" @click="handleAdd(row)">
            新增
          </el-button>
          <el-button v-permission="'system:dept:remove'" type="danger" link icon="Delete" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="上级部门" prop="parentId">
              <el-tree-select
                v-model="form.parentId"
                :data="deptOptions"
                :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
                placeholder="请选择上级部门"
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门名称" prop="deptName">
              <el-input v-model="form.deptName" placeholder="请输入部门名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="负责人" prop="leader">
              <el-input v-model="form.leader" placeholder="请输入负责人" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
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
import { listDept, getDept, createDept, updateDept, deleteDept } from '@/api/dept'

const loading = ref(false)
const deptList = ref([])
const deptOptions = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const isExpandAll = ref(true)
const refreshTable = ref(true)

const formRef = ref()

const form = reactive({
  deptId: null,
  parentId: 0,
  deptName: '',
  sort: 0,
  leader: '',
  phone: '',
  email: '',
  status: '0'
})

const rules = {
  deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
}

onMounted(() => {
  getList()
})

const getList = async () => {
  loading.value = true
  try {
    const res = await listDept()
    deptList.value = res.data || []
    deptOptions.value = [{ deptId: 0, deptName: '根部门', children: res.data }]
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
  form.deptId = null
  form.parentId = 0
  form.deptName = ''
  form.sort = 0
  form.leader = ''
  form.phone = ''
  form.email = ''
  form.status = '0'
}

const handleAdd = (row) => {
  resetForm()
  if (row) {
    form.parentId = row.deptId
  }
  dialogTitle.value = '新增部门'
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  resetForm()
  const res = await getDept(row.deptId)
  const dept = res.data
  Object.assign(form, {
    deptId: dept.deptId,
    parentId: dept.parentId,
    deptName: dept.deptName,
    sort: dept.sort,
    leader: dept.leader,
    phone: dept.phone,
    email: dept.email,
    status: dept.status
  })
  dialogTitle.value = '编辑部门'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (form.deptId) {
      await updateDept(form.deptId, form)
      ElMessage.success('更新成功')
    } else {
      await createDept(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定删除该部门吗？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteDept(row.deptId)
      ElMessage.success('删除成功')
      getList()
    })
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>
