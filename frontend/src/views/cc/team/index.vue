<template>
  <div class="page-card">
    <!-- 搜索表单 -->
    <el-form :inline="true" :model="queryParams" class="table-actions">
      <div>
        <el-form-item label="团队名称">
          <el-input v-model="queryParams.teamName" placeholder="请输入团队名称" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="业务类型">
          <el-select v-model="queryParams.businessType" placeholder="请选择" clearable style="width: 120px">
            <el-option label="业务A" value="业务A" />
            <el-option label="业务B" value="业务B" />
            <el-option label="业务C" value="业务C" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属军团">
          <el-select v-model="queryParams.legionId" placeholder="请选择" clearable style="width: 150px">
            <el-option v-for="item in legionOptions" :key="item.id" :label="item.legionName" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </div>
      <div>
        <el-button type="primary" icon="Plus" @click="handleAdd">添加团队</el-button>
      </div>
    </el-form>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="list" border>
      <el-table-column prop="id" label="团队ID" width="100" align="center" />
      <el-table-column prop="performanceRank" label="排名" width="80" align="center" />
      <el-table-column label="当月业绩" width="150" align="right">
        <template #default="{ row }">
          ¥ {{ formatAmount(row.monthlyPerformance) }}
        </template>
      </el-table-column>
      <el-table-column label="团队余额" width="150" align="right">
        <template #default="{ row }">
          ¥ {{ formatAmount(row.balance) }}
        </template>
      </el-table-column>

      <el-table-column prop="teamName" label="团队名称" min-width="120" />
      <el-table-column prop="businessType" label="业务类型" width="100" />
      <el-table-column prop="leaderName" label="团长" width="100">
        <template #default="{ row }">{{ row.leaderName || '-' }}</template>
      </el-table-column>
      <el-table-column prop="legionName" label="军团名称" width="120">
        <template #default="{ row }">{{ row.legionName || '-' }}</template>
      </el-table-column>
      <el-table-column prop="legionLeaderName" label="军团长" width="100">
        <template #default="{ row }">{{ row.legionLeaderName || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="primary" link @click="handleFund(row)">资金</el-button>
          <el-button type="info" link @click="handleLogs(row)">修改记录</el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="团队名称" prop="teamName">
          <el-input v-model="form.teamName" placeholder="不超过10个字符" maxlength="10" />
        </el-form-item>
        <el-form-item label="业务类型" prop="businessType">
          <el-select v-model="form.businessType" placeholder="必选" style="width: 100%">
            <el-option label="业务A" value="业务A" />
            <el-option label="业务B" value="业务B" />
            <el-option label="业务C" value="业务C" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.id" label="团长">
          <el-select v-model="form.leaderId" placeholder="请选择团长" clearable style="width: 100%">
            <el-option v-for="item in leaderOptions" :key="item.userId" :label="item.nickName" :value="item.userId" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属军团">
          <el-select v-model="form.legionId" placeholder="请选择所属军团" clearable style="width: 100%">
            <el-option v-for="item in legionOptions" :key="item.id" :label="item.legionName" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="showStructAdjust" label="架构调整时间" prop="structAdjustDate">
          <el-date-picker v-model="form.structAdjustDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">提交</el-button>
      </template>
    </el-dialog>

    <!-- 交易确认弹窗 -->
    <el-dialog v-model="confirmDialogVisible" title="交易确认" width="400px">
      <p>{{ confirmMessage }}</p>
      <el-form ref="confirmFormRef" :model="confirmForm" :rules="confirmRules" label-width="100px">
        <el-form-item label="交易金额" prop="amount">
          <el-input-number v-model="confirmForm.amount" :min="1" :precision="0" style="width: 100%" placeholder="必填" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="confirmDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleConfirmSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 修改记录弹窗 -->
    <el-dialog v-model="logsDialogVisible" title="修改记录" width="700px">
      <div v-loading="logsLoading" class="logs-container">
        <div v-for="log in logs" :key="log.id" class="log-item">
          <div class="log-header">
            <span class="operator">{{ log.operatorName }}</span>
            <span class="time">{{ formatDate(log.createTime) }}</span>
          </div>
          <div class="log-content">{{ log.content }}</div>
        </div>
        <el-empty v-if="logs.length === 0" description="暂无记录" />
      </div>
    </el-dialog>

    <!-- 资金管理弹窗 -->
    <el-dialog v-model="fundDialogVisible" :title="`团队资金（${currentRow?.teamName}）`" width="800px">
      <div class="fund-header">
        <div class="balance-info">
          <span class="label">余额：</span>
          <span class="value">¥ {{ formatAmount(fundInfo.balance) }}</span>
        </div>
        <div class="fund-actions">
          <el-button type="primary" @click="handleEditBalance">编辑余额</el-button>
          <el-button type="success" @click="handleRecharge">充值</el-button>
          <el-button type="warning" @click="handleTransfer">转账</el-button>
        </div>
      </div>
      
      <el-tabs v-model="billType" @tab-change="loadBills">
        <el-tab-pane label="非流水" name="non_flow" />
        <el-tab-pane label="流水" name="flow" />
        <el-tab-pane label="全部" name="all" />
      </el-tabs>
      
      <el-table v-loading="billsLoading" :data="bills" border max-height="400">
        <el-table-column prop="createTime" label="时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createTime) }}</template>
        </el-table-column>
        <el-table-column prop="logType" label="类别" width="150">
          <template #default="{ row }">{{ formatLogType(row.logType) }}</template>
        </el-table-column>
        <el-table-column prop="content" label="说明" />
      </el-table>
    </el-dialog>

    <!-- 编辑余额弹窗 -->
    <el-dialog v-model="editBalanceDialogVisible" title="编辑余额" width="400px">
      <el-form ref="editBalanceFormRef" :model="editBalanceForm" :rules="editBalanceRules" label-width="100px">
        <el-form-item label="修改类别" prop="editType">
          <el-select v-model="editBalanceForm.editType" style="width: 100%">
            <el-option label="增减" value="adjust" />
            <el-option label="设置为" value="set" />
          </el-select>
        </el-form-item>
        <el-form-item label="修改余额" prop="amount">
          <el-input-number v-model="editBalanceForm.amount" :precision="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="修改原因" prop="reason">
          <el-input v-model="editBalanceForm.reason" type="textarea" placeholder="必填" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editBalanceDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleEditBalanceSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 充值/转账弹窗 (与Legion类似，省略重复代码) -->
    <el-dialog v-model="rechargeDialogVisible" title="充值" width="400px">
      <el-form ref="rechargeFormRef" :model="rechargeForm" :rules="rechargeRules" label-width="100px">
        <el-form-item label="团长余额">
          <span>¥ {{ formatAmount(leaderBalance) }}</span>
        </el-form-item>
        <el-form-item label="充值金额" prop="amount">
          <el-input-number v-model="rechargeForm.amount" :min="1" :precision="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleRechargeSubmit">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="transferDialogVisible" title="转账" width="400px">
      <el-form ref="transferFormRef" :model="transferForm" :rules="transferRules" label-width="100px">
        <el-form-item label="转账金额" prop="amount">
          <el-input-number v-model="transferForm.amount" :min="1" :precision="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="收账人" prop="recipientId">
          <el-select v-model="transferForm.recipientId" placeholder="请选择收账人" filterable style="width: 100%">
            <el-option v-for="item in ccOptions" :key="item.id" :label="item.nickName" :value="item.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="transferDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleTransferSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({ name: 'CCTeam' })

import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listCCTeam, getCCTeam, createCCTeam, updateCCTeam, getCCTeamLogs, getCCTeamFund, editCCTeamFund, rechargeCCTeam, transferCCTeam, getCCTeamBills } from '@/api/ccTeam'
import { listAllLegion } from '@/api/legion'
import { listCC } from '@/api/cc'
import { listUser } from '@/api/user'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const legionOptions = ref([])
const ccOptions = ref([])
const leaderOptions = ref([])

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  teamName: '',
  businessType: '',
  legionId: null
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const formRef = ref()

const form = reactive({
  id: null,
  teamName: '',
  businessType: '',
  leaderId: null,
  legionId: null,
  structAdjustDate: '',
  transactionAmount: null
})

const originalLegionId = ref(null)
const showStructAdjust = computed(() => form.legionId !== originalLegionId.value && originalLegionId.value !== null)

const rules = {
  teamName: [{ required: true, message: '请输入团队名称', trigger: 'blur' }],
  businessType: [{ required: true, message: '请选择业务类型', trigger: 'change' }]
}

// 交易确认
const confirmDialogVisible = ref(false)
const confirmFormRef = ref()
const confirmForm = reactive({ amount: null })
const confirmRules = { amount: [{ required: true, message: '请输入交易金额', trigger: 'blur' }] }
const confirmMessage = ref('')
const confirmType = ref('') // 'add_team' | 'promote_leader'

// 日志和资金相关状态
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logs = ref([])
const fundDialogVisible = ref(false)
const currentRow = ref(null)
const fundInfo = reactive({ balance: 0 })
const billType = ref('non_flow')
const billsLoading = ref(false)
const bills = ref([])
const leaderBalance = ref(0)

const editBalanceDialogVisible = ref(false)
const editBalanceFormRef = ref()
const editBalanceForm = reactive({ editType: 'adjust', amount: 0, reason: '' })
const editBalanceRules = {
  editType: [{ required: true }],
  amount: [{ required: true }],
  reason: [{ required: true }]
}

const rechargeDialogVisible = ref(false)
const rechargeFormRef = ref()
const rechargeForm = reactive({ amount: null })
const rechargeRules = { amount: [{ required: true }] }

const transferDialogVisible = ref(false)
const transferFormRef = ref()
const transferForm = reactive({ amount: null, recipientId: null })
const transferRules = {
  amount: [{ required: true }],
  recipientId: [{ required: true }]
}

onMounted(() => {
  loadOptions()
  getList()
})

const loadOptions = async () => {
  const [legionRes, ccRes] = await Promise.all([
    listAllLegion(),
    listCC({ pageSize: 9999, isBlocked: '0' })
  ])
  legionOptions.value = legionRes.data || []
  ccOptions.value = ccRes.data.rows || []
}

const getList = async () => {
  loading.value = true
  try {
    const res = await listCCTeam(queryParams)
    list.value = res.data.rows || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const handleQuery = () => { queryParams.pageNum = 1; getList() }
const resetQuery = () => { Object.assign(queryParams, { teamName: '', businessType: '', legionId: null }); handleQuery() }

const resetForm = () => {
  Object.assign(form, { id: null, teamName: '', businessType: '', leaderId: null, legionId: null, structAdjustDate: '', transactionAmount: null })
  originalLegionId.value = null
}

const handleAdd = () => { resetForm(); dialogTitle.value = '添加团队'; dialogVisible.value = true }

const handleEdit = async (row) => {
  resetForm()
  const res = await getCCTeam(row.id)
  const data = res.data
  Object.assign(form, { id: data.id, teamName: data.teamName, businessType: data.businessType, leaderId: data.leaderId, legionId: data.legionId })
  originalLegionId.value = data.legionId
  // 加载可选团长：必须是CC团长角色 (从用户表加载)
  const userRes = await listUser({ roleKey: 'cc_team_leader', pageSize: 1000 })
  leaderOptions.value = userRes.data.rows || []
  dialogTitle.value = `编辑团队（${data.id}）`
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  
  // 新增团队且有军团，需要交易确认
  if (!form.id && form.legionId) {
    confirmMessage.value = '确认添加团队吗？'
    confirmType.value = 'add_team'
    confirmForm.amount = null
    confirmDialogVisible.value = true
    return
  }
  
  // 修改团长，需要交易确认
  const original = list.value.find(x => x.id === form.id)
  if (form.id && form.leaderId && form.leaderId !== original?.leaderId) {
    confirmMessage.value = '确认修改团长吗？'
    confirmType.value = 'promote_leader'
    confirmForm.amount = null
    confirmDialogVisible.value = true
    return
  }
  
  submitAction()
}

const handleConfirmSubmit = async () => {
  await confirmFormRef.value.validate()
  form.transactionAmount = confirmForm.amount * 100
  submitAction()
}

const submitAction = async () => {
  submitLoading.value = true
  try {
    if (form.id) {
      await updateCCTeam(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createCCTeam(form)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    confirmDialogVisible.value = false
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleLogs = async (row) => {
  currentRow.value = row
  logsLoading.value = true
  logsDialogVisible.value = true
  try {
    const res = await getCCTeamLogs(row.id)
    logs.value = res.data || []
  } finally {
    logsLoading.value = false
  }
}

const handleFund = async (row) => {
  currentRow.value = row
  fundDialogVisible.value = true
  billType.value = 'non_flow'
  const res = await getCCTeamFund(row.id)
  Object.assign(fundInfo, res.data)
  if (row.leaderId) {
    const ccRes = await listCC({ ccId: row.leaderId })
    leaderBalance.value = ccRes.data.rows?.[0]?.balance || 0
  } else {
    leaderBalance.value = 0
  }
  loadBills()
}

const loadBills = async () => {
  billsLoading.value = true
  try {
    const res = await getCCTeamBills(currentRow.value.id, billType.value)
    bills.value = res.data || []
  } finally {
    billsLoading.value = false
  }
}

const handleEditBalance = () => { Object.assign(editBalanceForm, { editType: 'adjust', amount: 0, reason: '' }); editBalanceDialogVisible.value = true }

const handleEditBalanceSubmit = async () => {
  await editBalanceFormRef.value.validate()
  if (editBalanceForm.editType === 'set' && editBalanceForm.amount < 0) { ElMessage.error('余额不可以操作至负数'); return }
  submitLoading.value = true
  try {
    await editCCTeamFund(currentRow.value.id, { editType: editBalanceForm.editType, amount: editBalanceForm.amount * 100, reason: editBalanceForm.reason })
    ElMessage.success('修改成功')
    editBalanceDialogVisible.value = false
    handleFund(currentRow.value)
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleRecharge = () => {
  if (!currentRow.value.leaderId) { ElMessage.warning('团队没有团长，无法充值'); return }
  rechargeForm.amount = null
  rechargeDialogVisible.value = true
}

const handleRechargeSubmit = async () => {
  await rechargeFormRef.value.validate()
  submitLoading.value = true
  try {
    await rechargeCCTeam(currentRow.value.id, { amount: rechargeForm.amount * 100 })
    ElMessage.success('充值成功')
    rechargeDialogVisible.value = false
    handleFund(currentRow.value)
    getList()
  } finally {
    submitLoading.value = false
  }
}

const handleTransfer = () => { Object.assign(transferForm, { amount: null, recipientId: null }); transferDialogVisible.value = true }

const handleTransferSubmit = async () => {
  await transferFormRef.value.validate()
  submitLoading.value = true
  try {
    await transferCCTeam(currentRow.value.id, { amount: transferForm.amount * 100, recipientId: transferForm.recipientId })
    ElMessage.success('转账成功')
    transferDialogVisible.value = false
    handleFund(currentRow.value)
    getList()
  } finally {
    submitLoading.value = false
  }
}

const formatAmount = (amount) => amount ? (amount / 100).toFixed(2) : '0.00'
const formatDate = (date) => date ? new Date(date).toLocaleString('zh-CN') : '-'
const formatLogType = (type) => ({
  'team_balance_edit': '余额修改',
  'team_recharge': '充值',
  'team_transfer': '转账',
  'add_squad': '添加战队',
  'promote_team_leader': '晋升团长'
})[type] || type
</script>

<style scoped>
.logs-container { max-height: 400px; overflow-y: auto; }
.log-item { padding: 12px; border-bottom: 1px solid #eee; }
.log-item:last-child { border-bottom: none; }
.log-header { display: flex; justify-content: space-between; margin-bottom: 8px; font-size: 12px; color: #999; }
.log-content { color: #333; }
.fund-header { display: flex; justify-content: space-between; align-items: center; padding: 16px; background: #f5f7fa; border-radius: 4px; margin-bottom: 16px; }
.balance-info .label { color: #666; }
.balance-info .value { font-size: 24px; font-weight: bold; color: #409eff; }
.fund-actions { display: flex; gap: 8px; }
</style>
