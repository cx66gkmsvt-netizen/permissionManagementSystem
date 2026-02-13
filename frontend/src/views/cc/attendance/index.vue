<template>
  <div class="page-card">
    <!-- 搜索表单 -->
    <el-form :inline="true" :model="queryParams" class="table-actions">
      <div>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 260px"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item label="CC姓名">
          <el-input v-model="queryParams.name" placeholder="请输入" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item label="军团">
          <el-select v-model="queryParams.legionId" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="item in legionOptions" :key="item.id" :label="item.legionName" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
          <el-button type="success" icon="Download" @click="handleExport" :loading="exportLoading">下载</el-button>
        </el-form-item>
      </div>
    </el-form>

    <!-- 在班统计 -->
    <div class="stats-bar" v-if="stats">
      <el-tag type="success">在班: {{ stats['1'] || 0 }}</el-tag>
      <el-tag type="info">休班: {{ stats['2'] || 0 }}</el-tag>
      <el-tag type="warning">请假: {{ stats['3'] || 0 }}</el-tag>
    </div>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="list" border style="width: 100%">
      <el-table-column prop="ccId" label="CCID" width="100" />
      <el-table-column prop="ccName" label="姓名" width="100" fixed="left" />
      <el-table-column prop="nickName" label="昵称" width="100" />
      <el-table-column prop="squadName" label="战队" width="100" />
      <el-table-column prop="teamName" label="团队" width="100" />
      <el-table-column prop="legionName" label="军团" width="100" />
      <el-table-column 
        v-for="date in displayDates" 
        :key="date" 
        :label="formatDateLabel(date)" 
        width="80"
        align="center"
      >
        <template #default="{ row }">
          <el-select
            v-model="row[date]"
            size="small"
            placeholder="--"
            style="width: 70px"
            @change="handleStatusChange(row, date, $event)"
          >
            <el-option label="在" value="1" />
            <el-option label="休" value="2" />
            <el-option label="假" value="3" />
          </el-select>
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
  </div>
</template>

<script setup>
defineOptions({ name: 'Attendance' })

import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAttendanceList, updateAttendance, getAttendanceStats, exportAttendance } from '@/api/attendance'
import { listAllLegion } from '@/api/legion'

const loading = ref(false)
const exportLoading = ref(false)
const list = ref([])
const total = ref(0)
const legionOptions = ref([])
const stats = ref(null)

// 日期范围
const today = new Date()
const dateRange = ref([
  formatDateStr(new Date(today.getTime() - 6 * 24 * 60 * 60 * 1000)),
  formatDateStr(today)
])

const queryParams = reactive({
  pageNum: 1,
  pageSize: 50,
  name: '',
  legionId: null
})

// 显示的日期列表（仅在搜索时更新）
const displayDates = ref([])

function buildDates(range) {
  if (!range || range.length !== 2) return []
  const dates = []
  const start = new Date(range[0])
  const end = new Date(range[1])
  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    dates.push(formatDateStr(new Date(d)))
  }
  return dates
}

onMounted(() => {
  loadLegions()
  displayDates.value = buildDates(dateRange.value)
  getList()
  loadStats()
})

const loadLegions = async () => {
  const res = await listAllLegion()
  legionOptions.value = res.data || []
}

const getList = async () => {
  if (displayDates.value.length === 0) return
  loading.value = true
  try {
    const params = {
      ...queryParams,
      dates: displayDates.value
    }
    const res = await getAttendanceList(params)
    list.value = res.data || []
    total.value = list.value.length
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  const date = formatDateStr(today)
  const res = await getAttendanceStats(date)
  stats.value = res.data || {}
}

const handleQuery = () => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    ElMessage.warning('请选择日期范围')
    return
  }
  const start = new Date(dateRange.value[0])
  const end = new Date(dateRange.value[1])
  const diffDays = Math.round((end - start) / (24 * 60 * 60 * 1000)) + 1
  if (diffDays > 31) {
    ElMessage.warning('日期范围最多不超过31天')
    return
  }
  displayDates.value = buildDates(dateRange.value)
  queryParams.pageNum = 1
  getList()
  loadStats()
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.legionId = null
  handleQuery()
}



const handleStatusChange = async (row, date, status) => {
  try {
    await updateAttendance(row.ccId, date, { status })
    ElMessage.success('更新成功')
    loadStats()
  } catch (error) {
    ElMessage.error('更新失败')
    getList() // 刷新恢复
  }
}

const disabledDate = (date) => {
  return date > new Date()
}

const formatDateLabel = (dateStr) => {
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  return `${month}/${day}\n${weekDays[date.getDay()]}`
}

function formatDateStr(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const handleExport = async () => {
  if (displayDates.value.length === 0) {
    ElMessage.warning('请选择日期范围')
    return
  }
  exportLoading.value = true
  try {
    const params = {
      ...queryParams,
      dates: displayDates.value
    }
    const res = await exportAttendance(params)
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `在班记录_${dateRange.value[0]}_${dateRange.value[1]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    exportLoading.value = false
  }
}
</script>

<style scoped>
.stats-bar {
  margin-bottom: 16px;
  display: flex;
  gap: 12px;
}
</style>
