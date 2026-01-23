<template>
  <div class="page-card">
    <!-- 搜索表单 -->
    <el-form :inline="true" :model="queryParams" class="table-actions">
      <div>
        <el-form-item label="姓名">
          <el-input v-model="queryParams.name" placeholder="请输入姓名" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="queryParams.mobile" placeholder="请输入手机号" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="军团">
          <el-tree-select
            v-model="queryParams.teamId"
            :data="teamOptions"
            :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
            placeholder="请选择军团"
            check-strictly
            clearable
            style="width: 150px"
            @change="handleQueryTeamChange"
          />
        </el-form-item>
        <el-form-item label="战队">
          <el-tree-select
            v-model="queryParams.squadId"
            :data="querySquadOptions"
            :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
            placeholder="请选择战队"
            check-strictly
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 120px">
            <el-option label="正常" value="0" />
            <el-option label="停用" value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </div>
      <div>
        <el-button type="primary" icon="Plus" @click="handleAdd">新增CC</el-button>
      </div>
    </el-form>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="ccList" border>
      <el-table-column type="index" label="序号" width="60" align="center" />
      <el-table-column prop="name" label="姓名" width="100" />
      <el-table-column prop="nickName" label="昵称" width="100" />
      <el-table-column prop="mobile" label="手机号" width="120" />
      <el-table-column prop="wechat" label="微信号" width="120" />
      <el-table-column prop="cno" label="座席号" width="100" />
      <el-table-column prop="cloudAccount" label="云客账号" width="120" />
      <el-table-column label="所属军团" width="120">
        <template #default="{ row }">
          {{ getDeptName(row.teamId) }}
        </template>
      </el-table-column>
      <el-table-column label="所属战队" width="120">
        <template #default="{ row }">
          {{ getDeptName(row.squadId) }}
        </template>
      </el-table-column>
      <el-table-column prop="balance" label="个人资金" width="120">
        <template #default="{ row }">
          ¥ {{ row.balance }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === '0' ? 'success' : 'danger'">
            {{ row.status === '0' ? '正常' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.createTime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link icon="Edit" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button type="danger" link icon="Delete" @click="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="form.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称" prop="nickName">
              <el-input v-model="form.nickName" placeholder="请输入昵称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="手机号" prop="mobile">
              <el-input v-model="form.mobile" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="微信号" prop="wechat">
              <el-input v-model="form.wechat" placeholder="请输入微信号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="座席号" prop="cno">
              <el-input v-model="form.cno" placeholder="请输入座席号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="云客账号" prop="cloudAccount">
              <el-input v-model="form.cloudAccount" placeholder="请输入云客账号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属军团" prop="teamId">
              <el-tree-select
                v-model="form.teamId"
                :data="teamOptions"
                :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
                placeholder="请选择军团"
                check-strictly
                style="width: 100%"
                @change="handleFormTeamChange"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属战队" prop="squadId">
              <el-tree-select
                v-model="form.squadId"
                :data="formSquadOptions"
                :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
                placeholder="请选择战队"
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="0">正常</el-radio>
            <el-radio value="1">停用</el-radio>
          </el-radio-group>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listCC, getCC, createCC, updateCC, deleteCC } from '@/api/cc'
import { listDept } from '@/api/dept'

const loading = ref(false)
const ccList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)

// 部门数据
const deptList = ref([])
const teamOptions = ref([])
const querySquadOptions = ref([])
const formSquadOptions = ref([])

const formRef = ref()

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  mobile: '',
  teamId: null,
  squadId: null,
  status: ''
})

const form = reactive({
  id: null,
  name: '',
  nickName: '',
  mobile: '',
  wechat: '',
  cno: '',
  cloudAccount: '',
  teamId: null,
  squadId: null,
  status: '0'
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  teamId: [{ required: true, message: '请选择所属军团', trigger: 'change' }]
}

onMounted(() => {
  getDepts()
  getList()
})

const getList = async () => {
  loading.value = true
  try {
    const res = await listCC(queryParams)
    ccList.value = res.data.rows || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const getDepts = async () => {
  const res = await listDept()
  deptList.value = res.data || []
  teamOptions.value = deptList.value // 假设顶层是军团
}

// 查找部门名称
const getDeptName = (deptId) => {
  if (!deptId) return '-'
  const findDept = (depts) => {
    for (const d of depts) {
      if (d.deptId === deptId) return d.deptName
      if (d.children) {
        const name = findDept(d.children)
        if (name) return name
      }
    }
    return null
  }
  return findDept(deptList.value) || deptId
}

// 搜索栏军团变化
const handleQueryTeamChange = (val) => {
  queryParams.squadId = null
  querySquadOptions.value = []
  if (val) {
    const team = findDeptInTree(teamOptions.value, val)
    if (team && team.children) {
      querySquadOptions.value = team.children
    }
  }
}

// 表单军团变化
const handleFormTeamChange = (val) => {
  form.squadId = null
  formSquadOptions.value = []
  if (val) {
    const team = findDeptInTree(teamOptions.value, val)
    if (team && team.children) {
      formSquadOptions.value = team.children
    }
  }
}

const findDeptInTree = (tree, id) => {
  for (const node of tree) {
    if (node.deptId === id) return node
    if (node.children) {
      const res = findDeptInTree(node.children, id)
      if (res) return res
    }
  }
  return null
}

const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.mobile = ''
  queryParams.teamId = null
  queryParams.squadId = null
  queryParams.status = ''
  querySquadOptions.value = []
  handleQuery()
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.nickName = ''
  form.mobile = ''
  form.wechat = ''
  form.cno = ''
  form.cloudAccount = ''
  form.teamId = null
  form.squadId = null
  form.status = '0'
  formSquadOptions.value = []
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增CC'
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  resetForm()
  const res = await getCC(row.id)
  const data = res.data
  Object.assign(form, {
    id: data.id,
    name: data.name,
    nickName: data.nickName,
    mobile: data.mobile,
    wechat: data.wechat,
    cno: data.cno,
    cloudAccount: data.cloudAccount,
    teamId: data.teamId,
    squadId: data.squadId,
    status: data.status
  })
  // 初始化战队选项
  if (data.teamId) {
    const team = findDeptInTree(teamOptions.value, data.teamId)
    if (team && team.children) {
      formSquadOptions.value = team.children
    }
  }
  dialogTitle.value = '编辑CC'
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (form.id) {
      await updateCC(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createCC(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定删除该CC吗？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteCC(row.id)
      ElMessage.success('删除成功')
      getList()
    })
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>
