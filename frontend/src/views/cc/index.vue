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
        
        <!-- 组织架构筛选 -->
        <el-form-item label="军团">
          <el-select
            v-model="queryParams.legionId"
            placeholder="请选择军团"
            clearable
            filterable
            style="width: 150px"
            @change="handleQueryLegionChange"
          >
            <el-option
              v-for="item in legionOptions"
              :key="item.id"
              :label="item.legionName"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="团队">
          <el-select
            v-model="queryParams.teamId"
            placeholder="请选择团队"
            clearable
            filterable
            style="width: 150px"
            @change="handleQueryTeamChange"
          >
            <el-option
              v-for="item in queryTeamOptions"
              :key="item.id"
              :label="item.teamName"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="战队">
          <el-select
            v-model="queryParams.squadId"
            placeholder="请选择战队"
            clearable
            filterable
            style="width: 150px"
          >
            <el-option
              v-for="item in querySquadOptions"
              :key="item.id"
              :label="item.squadName"
              :value="item.id"
            />
          </el-select>
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
      <el-table-column prop="id" label="CCID" width="80" align="center" />
      <el-table-column prop="name" label="姓名" width="100" />
      <el-table-column prop="nickName" label="昵称" width="100" />
      <el-table-column prop="mobile" label="手机号" width="120" />
      <el-table-column prop="wechat" label="微信号" width="120" />
      <el-table-column prop="cno1" label="座席号" width="100" />
      <el-table-column prop="cloudAccount1" label="云客账号" width="120" />
      <el-table-column prop="legionName" label="所属军团" width="120" show-overflow-tooltip />
      <el-table-column prop="teamName" label="所属团队" width="120" show-overflow-tooltip />
      <el-table-column prop="squadName" label="所属战队" width="120" show-overflow-tooltip />
      <el-table-column prop="balanceYuan" label="个人资金" width="120" align="right">
        <template #default="{ row }">
          ¥ {{ row.balanceYuan?.toFixed(2) }}
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
            <el-form-item label="座席号" prop="cno1">
              <el-input v-model="form.cno1" placeholder="请输入座席号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="云客账号" prop="cloudAccount1">
              <el-input v-model="form.cloudAccount1" placeholder="请输入云客账号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属军团" prop="legionId">
              <el-select
                v-model="form.legionId"
                placeholder="请选择军团"
                filterable
                style="width: 100%"
                @change="handleFormLegionChange"
              >
                <el-option
                  v-for="item in legionOptions"
                  :key="item.id"
                  :label="item.legionName"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属团队" prop="teamId">
              <el-select
                v-model="form.teamId"
                placeholder="请选择团队"
                filterable
                style="width: 100%"
                @change="handleFormTeamChange"
              >
                <el-option
                  v-for="item in formTeamOptions"
                  :key="item.id"
                  :label="item.teamName"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属战队" prop="squadId">
              <el-select
                v-model="form.squadId"
                placeholder="请选择战队"
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="item in formSquadOptions"
                  :key="item.id"
                  :label="item.squadName"
                  :value="item.id"
                />
              </el-select>
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
defineOptions({
  name: 'CC'
})
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listCC, getCC, createCC, updateCC, deleteCC } from '@/api/cc'
import { listAllLegion } from '@/api/legion'
import { listAllCCTeam } from '@/api/ccTeam'
import { listAllCCSquad } from '@/api/ccSquad'

const loading = ref(false)
const ccList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)

// 组织架构数据
const legionOptions = ref([])
const allTeams = ref([])
const allSquads = ref([])

const queryTeamOptions = ref([])
const querySquadOptions = ref([])

const formTeamOptions = ref([])
const formSquadOptions = ref([])

const formRef = ref()

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  mobile: '',
  legionId: null,
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
  cno1: '',
  cloudAccount1: '',
  legionId: null,
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
  legionId: [{ required: true, message: '请选择所属军团', trigger: 'change' }],
  teamId: [{ required: true, message: '请选择所属团队', trigger: 'change' }]
}

onMounted(() => {
  initOrgData()
  getList()
})

const initOrgData = async () => {
  const [legionRes, teamRes, squadRes] = await Promise.all([
    listAllLegion(),
    listAllCCTeam(),
    listAllCCSquad()
  ])
  legionOptions.value = legionRes.data || []
  allTeams.value = teamRes.data || []
  allSquads.value = squadRes.data || []
}

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

// 搜索栏级联
const handleQueryLegionChange = (val) => {
  queryParams.teamId = null
  queryParams.squadId = null
  queryTeamOptions.value = []
  querySquadOptions.value = []
  if (val) {
    queryTeamOptions.value = allTeams.value.filter(item => item.legionId === val)
  }
}

const handleQueryTeamChange = (val) => {
  queryParams.squadId = null
  querySquadOptions.value = []
  if (val) {
    querySquadOptions.value = allSquads.value.filter(item => item.teamId === val)
  }
}

// 表单级联
const handleFormLegionChange = (val) => {
  form.teamId = null
  form.squadId = null
  formTeamOptions.value = []
  formSquadOptions.value = []
  if (val) {
    formTeamOptions.value = allTeams.value.filter(item => item.legionId === val)
  }
}

const handleFormTeamChange = (val) => {
  form.squadId = null
  formSquadOptions.value = []
  if (val) {
    formSquadOptions.value = allSquads.value.filter(item => item.teamId === val)
  }
}

const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.mobile = ''
  queryParams.legionId = null
  queryParams.teamId = null
  queryParams.squadId = null
  queryParams.status = ''
  queryTeamOptions.value = []
  querySquadOptions.value = []
  handleQuery()
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.nickName = ''
  form.mobile = ''
  form.wechat = ''
  form.cno1 = ''
  form.cloudAccount1 = ''
  form.legionId = null
  form.teamId = null
  form.squadId = null
  form.status = '0'
  formTeamOptions.value = []
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
    cno1: data.cno1,
    cloudAccount1: data.cloudAccount1,
    legionId: data.legionId,
    teamId: data.teamId,
    squadId: data.squadId,
    status: data.status
  })
  
  // 初始化级联选项
  if (data.legionId) {
    formTeamOptions.value = allTeams.value.filter(item => item.legionId === data.legionId)
  }
  if (data.teamId) {
    formSquadOptions.value = allSquads.value.filter(item => item.teamId === data.teamId)
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
