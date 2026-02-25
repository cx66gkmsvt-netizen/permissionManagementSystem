<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="680px"
    @close="handleClose"
    align-center
    class="custom-dialog"
  >
    <div v-loading="loading" class="lanhu-container">
      
      <!-- Card 1: 基本信息 -->
      <div class="lanhu-card basic-info-card">
        <div class="info-row">
          <div class="info-item">
            <img class="icon-img" src="./assets/img/SketchPngf0d5d2c697915be87c13c1f78fab1425805653a0b0c9545526c78096201082ff.png" />
            <span class="label">当天日期</span>
            <span class="value">{{ detail.date }}</span>
          </div>
          <div class="dot-divider"></div>
          <div class="info-item">
            <span class="label">例子分配规则</span>
            <span class="value">{{ getRuleText(detail.allocationRule) }}</span>
          </div>
        </div>
        <div class="info-row" style="margin-top: 16px;">
          <div class="info-item">
            <img class="icon-img" src="./assets/img/SketchPng416ce1ac6bbdbc70fc231f19b030681b7313244885a02b967e15a009a04259d0.png" />
            <span class="label">分配时间</span>
            <span class="value">{{ detail.allocationTime || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- Card 2: 分配原因/未分配原因 & 达标情况 -->
      <div class="lanhu-card status-card">
        <div class="status-header">
          <div class="indicator-line"></div>
          <div class="status-title-col">
            <span class="title-text">是否分配例子</span>
            <span class="title-text sub-title">分配原因 / 未分配原因</span>
          </div>
          <div class="is-allocated-badge">
            <template v-if="detail.isAllocated === '1'">
              <img class="badge-icon" src="./assets/img/SketchPngefff77f9e10a217cfa460924397166a0ee3a1a1fde29dfac81c64f8134056e17.png" />
              <span class="badge-text success">是</span>
            </template>
            <template v-else>
              <el-icon color="#f56c6c"><CircleCloseFilled /></el-icon>
              <span class="badge-text danger">否</span>
            </template>
          </div>
        </div>
        
        <!-- 分配原因/未分配原因详情列表 -->
        <div class="reason-list">
          
          <!-- 1. 在班/不在班 -->
          <div class="reason-item">
            <img class="reason-icon" :src="getAttendanceIcon(detail.attendanceStatus)" />
            <span class="reason-text" :class="{'text-danger': detail.attendanceStatus !== '1'}">{{ getAttendanceText(detail.attendanceStatus) }}</span>
          </div>

          <!-- 2. 其他基于规则 Reasons -->
          
          <!-- 规则A: 节假日全部发放 -->
          <div v-if="detail.allocationRule === 'A'" class="reason-item">
            <img class="reason-icon" src="./assets/img/SketchPng4f26407c9019a19219c9900898ae4f7f0a0e6fdf52346b59703c112994620808.png" />
            <span class="reason-text">节假日全部发放</span>
          </div>
          
          <!-- 规则B/C达标判断 -->
          <div v-else-if="detail.lastWorkDay" class="reason-item with-details">
            <div class="reason-item-header">
              <img class="reason-icon" :src="detail.lastWorkDayReach ? './assets/img/SketchPng4f26407c9019a19219c9900898ae4f7f0a0e6fdf52346b59703c112994620808.png' : defaultErrorIcon" />
              <span class="reason-text" :class="{'text-danger': !detail.lastWorkDayReach}">上个工作日({{ detail.lastWorkDay }})通时通次{{ detail.lastWorkDayReach ? '达标' : '未达标' }}</span>
            </div>
            
            <div class="target-stats-box">
               <div class="stats-row">
                 <div class="stats-left">
                   <img class="reason-icon" :src="detail.lastWorkDayCallCnt >= detail.callCountTarget ? './assets/img/SketchPng4f26407c9019a19219c9900898ae4f7f0a0e6fdf52346b59703c112994620808.png' : defaultErrorIcon" />
                   <span>有效通话次数{{ detail.callCountTarget }}次</span>
                 </div>
                 <div class="stats-right" :class="detail.lastWorkDayCallCnt >= detail.callCountTarget ? 'success' : 'danger'">
                   已通话{{ detail.lastWorkDayCallCnt }}次
                 </div>
               </div>
               <div class="stats-row">
                 <div class="stats-left">
                   <img class="reason-icon" :src="detail.lastWorkDayCallDur >= detail.callDurationTarget ? './assets/img/SketchPng4f26407c9019a19219c9900898ae4f7f0a0e6fdf52346b59703c112994620808.png' : defaultErrorIcon" />
                   <span>有效通话时长{{ formatDuration(detail.callDurationTarget) }}</span>
                 </div>
                 <div class="stats-right" :class="detail.lastWorkDayCallDur >= detail.callDurationTarget ? 'success' : 'danger'">
                   已通话{{ formatDuration(detail.lastWorkDayCallDur) }}
                 </div>
               </div>
            </div>
          </div>

          <!-- 具体原因显示 -->
          <div v-if="detail.allocationReason" class="reason-item">
            <el-icon color="#909399" class="reason-icon"><InfoFilled /></el-icon>
            <span class="reason-text">{{ detail.allocationReason }}</span>
          </div>

        </div>
      </div>

      <!-- Card 3: 预计/透支/实发例子 -->
      <div class="lanhu-card final-card">
        <div class="final-header">
          <div class="final-left">
            <img class="icon-img" src="./assets/img/SketchPng8bc4c234f0e89e23ce65283b8cb9f59ad9f7a10499c337e1cc262adfd56c82a7.png" />
            <span class="label">预计/透支/实发例子</span>
          </div>
          <div class="final-right">
            {{ detail.expectedAllocation }}/{{ detail.overdraft }}/{{ detail.actualAllocation }}
          </div>
        </div>

        <!-- Super CC 等级明细表格 -->
        <div v-if="detail.levelBreakdown && detail.levelBreakdown.length > 0" class="super-cc-table">
          <el-table 
            :data="detail.levelBreakdown" 
            border 
            size="small"
            :header-cell-style="{ textAlign: 'center', background: '#f5f7fa', color: '#333' }"
            :cell-style="{ textAlign: 'center' }"
          >
            <el-table-column prop="name" label="等级" />
            <el-table-column prop="predicted" label="预计分配" />
            <el-table-column prop="overdraft" label="透支" />
            <el-table-column prop="actual" label="实际分配" />
          </el-table>
        </div>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getLeadAllocationSingleDetail } from '@/api/leadAllocation'
import { CircleCloseFilled, InfoFilled } from '@element-plus/icons-vue'
import successIcon from './assets/img/SketchPng4f26407c9019a19219c9900898ae4f7f0a0e6fdf52346b59703c112994620808.png'

const defaultErrorIcon = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="%23f56c6c" d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm0 393.664L407.328 353.312 353.312 407.328 457.664 512 353.312 616.672l54.016 54.016L512 566.336l104.672 104.672 54.016-54.016L566.336 512l104.32-104.32-54.016-54.016L512 457.664z"/></svg>'

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
    title.value = `${detail.value.ccName}（${detail.value.ccId}）例子分配详情`
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  visible.value = false
}

const getAttendanceText = (status) => {
  const map = { '1': '在班', '2': '休班', '3': '请假' }
  return map[status] || '未知状态'
}

const getAttendanceIcon = (status) => {
  return status === '1' ? successIcon : defaultErrorIcon
}

const getRuleText = (rule) => {
  const map = { 'A': '节假日', 'B': '工作日（无例子补偿）', 'C': '工作日' }
  return map[rule] || rule || '-'
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
/* 弹窗及容器 */
:deep(.custom-dialog .el-dialog__header) {
  border-bottom: 1px solid #ebeef5;
  margin-right: 0;
  padding-bottom: 15px;
}
:deep(.custom-dialog .el-dialog__title) {
  font-weight: 600;
  color: #222;
}

.lanhu-container {
  padding: 5px;
  background-color: #f7f8fa; /* 背景色 */
  border-radius: 4px;
}

.lanhu-card {
  background-color: #fff;
  border-radius: 12px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.basic-info-card {
  padding: 20px;
  background: url(./assets/img/SketchPng51a0ce7f2f6aeb91c05a875bd87cac5f83c84d00d7dcc2f589556d4fa9293ec4.png) no-repeat right bottom;
  background-size: cover;
}

.info-row {
  display: flex;
  align-items: center;
}

.info-item {
  display: flex;
  align-items: center;
}

.icon-img {
  width: 16px;
  height: 16px;
  margin-right: 8px;
}

.label {
  color: #888;
  font-size: 15px;
  margin-right: 12px;
}

.value {
  color: #222;
  font-size: 15px;
  font-weight: 500;
}

.dot-divider {
  width: 8px;
  height: 8px;
  background-color: #d7d3ff;
  border-radius: 50%;
  margin: 0 24px;
}

/* 分配状态/原因 Card */
.status-card {
  padding: 20px;
  background: url(./assets/img/SketchPnge1f38266a947abe3b34265d96ddee112768ded71df1746dff3a66f7edfe4eea4.png) no-repeat right top;
  background-size: cover;
}

.status-header {
  display: flex;
  position: relative;
  padding-bottom: 15px;
  border-bottom: 1px dashed #ebeef5;
  margin-bottom: 15px;
  align-items: flex-start;
}

.indicator-line {
  width: 4px;
  height: 48px;
  background-color: #d7d3ff;
  border-radius: 2px;
  margin-right: 12px;
}

.status-title-col {
  display: flex;
  flex-direction: column;
}

.title-text {
  color: #888;
  font-size: 15px;
  line-height: 24px;
}

.sub-title {
  font-size: 14px;
  margin-top: 5px;
}

.is-allocated-badge {
  position: absolute;
  right: 10px;
  top: 0;
  display: flex;
  align-items: center;
  background: url(./assets/img/SketchPng4f243a825f045f39521ec069d21756e43c9bab2fef416d3cd96f3fdd47dab8fa.png) no-repeat center;
  background-size: contain;
  padding: 6px 12px;
}

.badge-icon {
  width: 14px;
  height: 14px;
  margin-right: 4px;
}

.badge-text {
  font-size: 14px;
  font-weight: 600;
}

.badge-text.success { color: #1db954; }
.badge-text.danger { color: #f56c6c; margin-left: 4px; }

/* 详情列表 */
.reason-list {
  padding-left: 16px;
}

.reason-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.reason-item:last-child {
  margin-bottom: 0;
}

.reason-icon {
  width: 16px;
  height: 16px;
  margin-right: 12px;
}

.reason-text {
  color: #222;
  font-size: 15px;
}

.text-danger {
  color: #f56c6c !important;
}

/* 达标子项详情 */
.with-details {
  flex-direction: column;
  align-items: flex-start;
}

.reason-item-header {
  display: flex;
  align-items: center;
}

.target-stats-box {
  margin-top: 12px;
  margin-left: 28px;
  width: calc(100% - 28px);
  background: #f9fafe;
  border-radius: 8px;
  padding: 15px;
}

.stats-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.stats-row:last-child {
  margin-bottom: 0;
}

.stats-left {
  display: flex;
  align-items: center;
  color: #222;
  font-size: 14px;
}

.stats-right {
  font-size: 14px;
  font-weight: 500;
}

.stats-right.success { color: #1db954; }
.stats-right.danger { color: #f56c6c; }

/* Card 3 结果 */
.final-card {
  padding: 0;
  overflow: hidden;
}

.final-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: url(./assets/img/SketchPngab4e29fe8aa970c5a605141740764e5bd6dcaf7eb6dd782b5a65139e8213eec8.png) no-repeat left bottom;
  background-size: cover;
}

.final-left {
  display: flex;
  align-items: center;
}

.final-right {
  font-size: 18px;
  font-weight: 600;
  color: #222;
}

.super-cc-table {
  padding: 0 20px 20px 20px;
}
</style>
