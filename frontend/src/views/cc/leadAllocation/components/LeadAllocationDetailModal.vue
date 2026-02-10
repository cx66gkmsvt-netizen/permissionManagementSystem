<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="650px"
    @close="handleClose"
    align-center
  >
    <div v-loading="loading" class="detail-container">
      <!-- Card 1: Basic Info -->
      <div class="info-card">
        <div class="info-row">
          <div class="info-item">
            <el-icon><Clock /></el-icon>
            <span class="label">当天日期</span>
            <span class="value">{{ detail.date }}</span>
          </div>
          <div class="info-item">
            <el-icon><Aim /></el-icon>
            <span class="label">例子分配规则</span>
            <span class="value">{{ detail.allocationRule || '-' }}</span>
          </div>
        </div>
        <div class="info-row">
          <div class="info-item">
            <el-icon><Timer /></el-icon>
            <span class="label">分配时间</span>
            <span class="value">{{ detail.allocationTime || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- Card 2: Allocation Status -->
      <div class="info-card">
        <div class="info-row">
          <div class="info-item">
            <el-icon><CircleCheck /></el-icon>
            <span class="label">是否分配例子</span>
            <el-tag :type="detail.isAllocated === '1' ? 'success' : 'info'" round>
              {{ detail.isAllocated === '1' ? '是' : '否' }}
            </el-tag>
          </div>
        </div>
        <div class="info-row">
          <div class="info-item">
            <el-icon><QuestionFilled /></el-icon>
            <span class="label">分配原因</span>
            <span class="value">{{ detail.allocationReason || '' }}</span>
          </div>
        </div>
      </div>

      <!-- Card 3: Attendance & Stats -->
      <div class="info-card bg-light">
        <div class="status-row">
          <el-icon color="#67c23a"><CircleCheckFilled /></el-icon>
          <span class="status-text">{{ getAttendanceText(detail.attendanceStatus) }}</span>
        </div>
        <div class="status-row" v-if="detail.lastWorkDay">
          <el-icon color="#67c23a"><CircleCheckFilled /></el-icon>
          <span class="status-text">上个工作日({{ detail.lastWorkDay }})通时通次达标</span>
        </div>

        <div class="stats-box" v-if="detail.lastWorkDay">
          <div class="stats-row">
            <div class="stats-left">
              <el-icon color="#67c23a"><CircleCheckFilled /></el-icon>
              <span>有效通话次数{{ detail.callCountTarget }}次</span>
            </div>
            <div class="stats-right success">
              已通话{{ detail.lastWorkDayCallCnt }}次
            </div>
          </div>
          <div class="stats-row">
             <div class="stats-left">
              <el-icon color="#67c23a"><CircleCheckFilled /></el-icon>
              <span>有效通话时长{{ formatDuration(detail.callDurationTarget) }}</span>
            </div>
            <div class="stats-right success">
              已通话{{ formatDuration(detail.lastWorkDayCallDur) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Card 4: Allocation Result (Standard vs Super) -->
      <div class="info-card">
        <div class="allocation-header">
          <el-icon><Calendar /></el-icon>
          <span>预计/透支/实发例子</span>
          <span class="allocation-nums">{{ detail.expectedAllocation }}/{{ detail.overdraft }}/{{ detail.actualAllocation }}</span>
        </div>

        <!-- Super CC Table -->
        <el-table 
          v-if="detail.levelBreakdown && detail.levelBreakdown.length > 0" 
          :data="detail.levelBreakdown" 
          border 
          style="margin-top: 10px;"
          :header-cell-style="{ textAlign: 'center', background: '#f5f7fa' }"
          :cell-style="{ textAlign: 'center' }"
        >
          <el-table-column prop="name" label="等级" />
          <el-table-column prop="predicted" label="预计" />
          <el-table-column prop="overdraft" label="透支" />
          <el-table-column prop="actual" label="实发" />
        </el-table>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getLeadAllocationSingleDetail } from '@/api/leadAllocation'
import { Clock, Aim, Timer, CircleCheck, QuestionFilled, CircleCheckFilled, Calendar } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  ccId: Number,
  date: String
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(false)
const loading = ref(false)
const detail = ref({})
const title = ref('')

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.ccId) {
    fetchDetail()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const fetchDetail = async () => {
  loading.value = true
  try {
    const res = await getLeadAllocationSingleDetail(props.ccId, props.date)
    detail.value = res.data || {}
    title.value = `${detail.value.ccName} (${detail.value.ccId}) 例子分配详情`
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  visible.value = false
}

const getAttendanceText = (status) => {
  const map = { '1': '在班', '2': '休班', '3': '请假' }
  return map[status] || '-'
}

const formatDuration = (seconds) => {
  if (!seconds) return '00:00:00'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}
</script>

<style scoped>
.detail-container {
  padding: 10px;
}
.info-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05);
}
.bg-light {
  background: #f9fafe;
  border: none;
}
.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  gap: 20px;
}
.info-row:last-child {
  margin-bottom: 0;
}
.info-item {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #606266;
}
.info-item .el-icon {
  margin-right: 6px;
  font-size: 16px;
  color: #909399;
}
.info-item .label {
  margin-right: 10px;
  color: #909399;
}
.info-item .value {
  color: #303133;
  font-weight: 500;
}
.status-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}
.status-row .el-icon {
  margin-right: 8px;
  font-size: 18px;
}
.stats-box {
  background: #fff;
  border-radius: 6px;
  padding: 12px;
  margin-top: 10px;
}
.stats-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
}
.stats-row:last-child {
  margin-bottom: 0;
}
.stats-left {
  display: flex;
  align-items: center;
  color: #606266;
}
.stats-left .el-icon {
  margin-right: 6px;
}
.stats-right.success {
  color: #67c23a;
}
.allocation-header {
  display: flex;
  align-items: center;
  font-size: 15px;
  color: #606266;
  margin-bottom: 5px;
}
.allocation-nums {
  margin-left: auto;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
</style>
