<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409EFF">
              <el-icon :size="24"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_users || 0 }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67C23A">
              <el-icon :size="24"><Goods /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_products || 0 }}</div>
              <div class="stat-label">商品总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #E6A23C">
              <el-icon :size="24"><List /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_orders || 0 }}</div>
              <div class="stat-label">订单总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #F56C6C">
              <el-icon :size="24"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ stats.total_amount || 0 }}</div>
              <div class="stat-label">销售总额</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>销售趋势</span>
          </template>
          <div ref="chartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>待处理订单</span>
          </template>
          <el-table :data="pendingOrders" style="width: 100%" max-height="300">
            <el-table-column prop="order_no" label="订单号" width="160" />
            <el-table-column prop="total_amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ row.total_amount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>商品分类统计</span>
          </template>
          <div ref="pieChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新用户</span>
          </template>
          <el-table :data="latestUsers" style="width: 100%">
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="nickname" label="昵称" width="120" />
            <el-table-column prop="phone" label="手机号" width="130" />
            <el-table-column prop="created_at" label="注册时间" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import request from '@/utils/request'

const chartRef = ref(null)
const pieChartRef = ref(null)

const stats = ref({
  total_users: 128,
  total_products: 56,
  total_orders: 324,
  total_amount: 25680.5
})

const pendingOrders = ref([
  { id: 1, order_no: 'ORD202401150001', total_amount: 128.5, status: 0 },
  { id: 2, order_no: 'ORD202401150002', total_amount: 256.0, status: 1 },
  { id: 3, order_no: 'ORD202401150003', total_amount: 89.9, status: 2 }
])

const latestUsers = ref([
  { id: 1, username: 'user001', nickname: '小明', phone: '138****8001', created_at: '2024-01-15 10:30' },
  { id: 2, username: 'user002', nickname: '小红', phone: '138****8002', created_at: '2024-01-14 15:20' },
  { id: 3, username: 'user003', nickname: '小刚', phone: '138****8003', created_at: '2024-01-13 09:15' }
])

const getStatusText = (status) => {
  const statusMap = { 0: '待付款', 1: '待发货', 2: '待收货', 3: '已完成', 4: '已取消' }
  return statusMap[status] || '未知'
}

const getStatusType = (status) => {
  const typeMap = { 0: 'warning', 1: 'info', 2: 'primary', 3: 'success', 4: 'danger' }
  return typeMap[status] || ''
}

const initLineChart = () => {
  if (!chartRef.value) return
  const chart = echarts.init(chartRef.value)
  const option = {
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月']
    },
    yAxis: {
      type: 'value',
      name: '销售额(元)'
    },
    series: [{
      data: [8200, 9320, 9010, 9340, 12900, 13300, 15200],
      type: 'line',
      smooth: true,
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(64, 158, 255, 0.8)' },
          { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
        ])
      },
      lineStyle: { color: '#409EFF', width: 2 }
    }]
  }
  chart.setOption(option)
}

const initPieChart = () => {
  if (!pieChartRef.value) return
  const chart = echarts.init(pieChartRef.value)
  const option = {
    tooltip: { trigger: 'item' },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 14, fontWeight: 'bold' }
      },
      labelLine: { show: false },
      data: [
        { value: 1048, name: '蔬菜' },
        { value: 735, name: '水果' },
        { value: 580, name: '肉类' },
        { value: 484, name: '海鲜' },
        { value: 300, name: '粮油' }
      ]
    }]
  }
  chart.setOption(option)
}

const fetchStats = async () => {
  try {
    const res = await request({
      url: '/admin/dashboard/stats',
      method: 'get'
    })
    if (res.code === 200) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

onMounted(() => {
  fetchStats()
  nextTick(() => {
    initLineChart()
    initPieChart()
  })
})
</script>

<style scoped>
.stat-card {
  cursor: pointer;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  margin-left: 16px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 4px;
}

.chart-container {
  height: 300px;
}
</style>
