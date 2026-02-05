<template>
  <div class="page-card">
    <!-- 搜索表单 -->
    <el-form :inline="true" :model="queryParams" class="table-actions">
      <div>
        <el-form-item label="日期">
          <el-date-picker
            v-model="currentDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 140px"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item label="CC姓名">
          <el-input v-model="queryParams.ccName" placeholder="请输入" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item label="军团">
          <el-select v-model="queryParams.legionId" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="item in legionOptions" :key="item.id" :label="item.legionName" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="在班状态">
          <el-select v-model="queryParams.attendanceStatus" placeholder="全部" clearable style="width: 100px">
            <el-option label="在班" value="1" />
            <el-option label="休班" value="2" />
            <el-option label="请假" value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </div>
      <div>
        <el-button type="success" @click="handleBatchAllocate">批量分配</el-button>
        <el-button type="warning" @click="handleBatchNotAllocate">批量不分配</el-button>
      </div>
    </el-form>

    <!-- 统计 -->
    <div class="stats-bar" v-if="stats">
      <el-tag>预计分配: {{ stats.totalExpected || 0 }}</el-tag>
      <el-tag type="success">实际分配: {{ stats.totalActual || 0 }}</el-tag>
      <el-tag type="danger">透支: {{ stats.totalOverdraft || 0 }}</el-tag>
    </div>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="list" border @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" />
      <el-table-column prop="ccName" label="姓名" width="100" />
      <el-table-column prop="nickName" label="昵称" width="100" />
      <el-table-column prop="squadName" label="战队" width="100" />
      <el-table-column prop="teamName" label="团队" width="100" />
      <el-table-column prop="legionName" label="军团" width="100" />
      <el-table-column prop="attendanceStatus" label="在班状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.attendanceStatus)" size="small">
            {{ getStatusText(row.attendanceStatus) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="expectedAllocation" label="预计分配" width="90" align="right" />
      <el-table-column prop="actualAllocation" label="实际分配" width="90" align="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleShowDetail(row)">{{ row.actualAllocation || 0 }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="overdraft" label="透支" width="70" align="right">
        <template #default="{ row }">
          <span :class="row.overdraft > 0 ? 'text-danger' : ''">{{ row.overdraft || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="isAllocated" label="是否分配" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.isAllocated === '1' ? 'success' : 'info'" size="small">
            {{ row.isAllocated === '1' ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="allocationRule" label="规则" width="80" />
      <el-table-column prop="allocationReason" label="原因/备注" min-width="150" />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="queryParams.pageNum"
      v-model:page-size="queryParams.pageSize"
      :page-sizes="[20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next"
      style="margin-top: 20px; justify-content: flex-end"
      @size-change="getList"
      @current-change="getList"
    />

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editDialogVisible" title="编辑分配" width="500px">
      <el-form ref="editFormRef" :model="editForm" label-width="100px">
        <el-form-item label="预计分配">
          <el-input-number v-model="editForm.expectedAllocation" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="实际分配">
          <el-input-number v-model="editForm.actualAllocation" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="透支">
          <el-input-number v-model="editForm.overdraft" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="是否分配">
          <el-radio-group v-model="editForm.isAllocated">
            <el-radio value="1">是</el-radio>
            <el-radio value="0">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="分配规则">
          <el-select v-model="editForm.allocationRule" style="width: 100%">
            <el-option label="A-节假日" value="A" />
            <el-option label="B-工作日无补偿" value="B" />
            <el-option label="C-工作日" value="C" />
          </el-select>
        </el-form-item>
        <el-form-item label="原因/备注">
          <el-input v-model="editForm.allocationReason" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleEditSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 批量操作弹窗 -->
    <el-dialog v-model="batchDialogVisible" :title="batchDialogTitle" width="400px">
      <el-form label-width="80px">
        <el-form-item label="原因">
          <el-input v-model="batchReason" type="textarea" placeholder="请输入原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleBatchSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 分配详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" :title="detailTitle" width="700px">
      <el-table v-loading="detailLoading" :data="detailList" border max-height="400px">
        <el-table-column prop="allocationDate" label="日期" width="120" />
        <el-table-column prop="expectedAllocation" label="预计分配" width="100" align="right" />
        <el-table-column prop="actualAllocation" label="实际分配" width="100" align="right" />
        <el-table-column prop="overdraft" label="透支" width="80" align="right">
          <template #default="{ row }">
            <span :class="row.overdraft > 0 ? 'text-danger' : ''">{{ row.overdraft || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="isAllocated" label="是否分配" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.isAllocated === '1' ? 'success' : 'info'" size="small">
              {{ row.isAllocated === '1' ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="allocationRule" label="规则" width="100" />
        <el-table-column prop="allocationReason" label="原因/备注" min-width="120" />
      </el-table>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({ name: 'LeadAllocation' })

import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getLeadAllocationList, updateLeadAllocation, batchUpdateIsAllocated, getLeadAllocationStats, getLeadAllocationDetail } from '@/api/leadAllocation'
import { listAllLegion } from '@/api/legion'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const legionOptions = ref([])
const stats = ref(null)
const selectedRows = ref([])

const today = new Date()
const currentDate = ref(formatDateStr(today))

const queryParams = reactive({
  pageNum: 1,
  pageSize: 50,
  ccName: '',
  legionId: null,
  attendanceStatus: ''
})

// 编辑
const editDialogVisible = ref(false)
const editFormRef = ref()
const currentRow = ref(null)
const editForm = reactive({
  expectedAllocation: 0,
  actualAllocation: 0,
  overdraft: 0,
  isAllocated: '0',
  allocationRule: '',
  allocationReason: ''
})
const submitLoading = ref(false)

const batchDialogVisible = ref(false)
const batchDialogTitle = ref('')
const batchReason = ref('')
const batchIsAllocated = ref('1')

// 详情弹窗
const detailDialogVisible = ref(false)
const detailLoading = ref(false)
const detailTitle = ref('')
const detailList = ref([])

onMounted(() => {
  loadLegions()
  getList()
  loadStats()
})

const loadLegions = async () => {
  const res = await listAllLegion()
  legionOptions.value = res.data || []
}

const getList = async () => {
  loading.value = true
  try {
    const params = { ...queryParams, date: currentDate.value }
    const res = await getLeadAllocationList(params)
    list.value = res.data?.rows || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  const res = await getLeadAllocationStats(currentDate.value)
  stats.value = res.data || {}
}

const handleQuery = () => { queryParams.pageNum = 1; getList() }
const resetQuery = () => { Object.assign(queryParams, { ccName: '', legionId: null, attendanceStatus: '' }); handleQuery() }
const handleDateChange = () => { getList(); loadStats() }
const handleSelectionChange = (rows) => { selectedRows.value = rows }

const handleEdit = (row) => {
  currentRow.value = row
  Object.assign(editForm, {
    expectedAllocation: row.expectedAllocation || 0,
    actualAllocation: row.actualAllocation || 0,
    overdraft: row.overdraft || 0,
    isAllocated: row.isAllocated || '0',
    allocationRule: row.allocationRule || '',
    allocationReason: row.allocationReason || ''
  })
  editDialogVisible.value = true
}

const handleEditSubmit = async () => {
  submitLoading.value = true
  try {
    await updateLeadAllocation(currentRow.value.ccId, currentDate.value, editForm)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    getList()
    loadStats()
  } finally {
    submitLoading.value = false
  }
}

const handleBatchAllocate = () => {
  if (selectedRows.value.length === 0) { ElMessage.warning('请先选择CC'); return }
  batchDialogTitle.value = '批量分配'
  batchIsAllocated.value = '1'
  batchReason.value = ''
  batchDialogVisible.value = true
}

const handleBatchNotAllocate = () => {
  if (selectedRows.value.length === 0) { ElMessage.warning('请先选择CC'); return }
  batchDialogTitle.value = '批量不分配'
  batchIsAllocated.value = '0'
  batchReason.value = ''
  batchDialogVisible.value = true
}

const handleBatchSubmit = async () => {
  submitLoading.value = true
  try {
    await batchUpdateIsAllocated({
      ccIds: selectedRows.value.map(r => r.ccId),
      date: currentDate.value,
      isAllocated: batchIsAllocated.value,
      reason: batchReason.value
    })
    ElMessage.success('操作成功')
    batchDialogVisible.value = false
    getList()
    loadStats()
  } finally {
    submitLoading.value = false
  }
}

const handleShowDetail = async (row) => {
  detailTitle.value = `${row.ccName} - 分配详情`
  detailDialogVisible.value = true
  detailLoading.value = true
  detailList.value = []
  try {
    // 获取最近30天的分配记录
    const endDate = currentDate.value
    const startDate = formatDateStr(new Date(new Date(endDate).getTime() - 29 * 24 * 60 * 60 * 1000))
    const res = await getLeadAllocationDetail(row.ccId, { startDate, endDate })
    detailList.value = res.data || []
  } catch (error) {
    ElMessage.error('获取详情失败')
  } finally {
    detailLoading.value = false
  }
}

const getStatusType = (status) => ({ '1': 'success', '2': 'info', '3': 'warning' })[status] || 'info'
const getStatusText = (status) => ({ '1': '在班', '2': '休班', '3': '请假' })[status] || '-'

function formatDateStr(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}
</script>

<style scoped>
.stats-bar { margin-bottom: 16px; display: flex; gap: 12px; }
.text-danger { color: #f56c6c; font-weight: 600; }
</style>
