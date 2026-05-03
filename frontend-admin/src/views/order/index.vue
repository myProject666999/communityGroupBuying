<template>
  <div class="order-container">
    <el-card>
      <template #header>
        <span>订单管理</span>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="待付款" value="pending" />
            <el-option label="待发货" value="paid" />
            <el-option label="待收货" value="shipped" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column prop="total_amount" label="订单金额" width="120">
          <template #default="{ row }">
            ¥{{ row.total_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" min-width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 1"
              type="success"
              link
              @click="handleShip(row)"
            >发货</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="detailVisible"
      title="订单详情"
      width="700px"
    >
      <div class="order-detail" v-if="currentOrder">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.order_no }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getStatusType(currentOrder.status)">{{ getStatusText(currentOrder.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户">{{ currentOrder.username }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ currentOrder.created_at }}</el-descriptions-item>
          <el-descriptions-item label="订单金额">¥{{ currentOrder.total_amount }}</el-descriptions-item>
          <el-descriptions-item label="支付时间" v-if="currentOrder.paid_at">{{ currentOrder.paid_at }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>收货地址</el-divider>
        <div class="address-info">
          <p><strong>收货人：</strong>{{ currentOrder.address?.name }}</p>
          <p><strong>联系电话：</strong>{{ currentOrder.address?.phone }}</p>
          <p><strong>收货地址：</strong>{{ currentOrder.address?.province }}{{ currentOrder.address?.city }}{{ currentOrder.address?.district }}{{ currentOrder.address?.detail }}</p>
        </div>

        <el-divider>商品信息</el-divider>
        <el-table :data="currentOrder.items || []" style="width: 100%">
          <el-table-column label="商品图片" width="100">
            <template #default="{ row }">
              <el-image
                :src="row.product_image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=default%20product&image_size=square'"
                style="width: 60px; height: 60px"
                fit="cover"
              />
            </template>
          </el-table-column>
          <el-table-column prop="product_name" label="商品名称" min-width="200" />
          <el-table-column prop="price" label="单价" width="100">
            <template #default="{ row }">¥{{ row.price }}</template>
          </el-table-column>
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column label="小计" width="100">
            <template #default="{ row }">¥{{ (row.price * row.quantity).toFixed(2) }}</template>
          </el-table-column>
        </el-table>

        <el-divider>备注</el-divider>
        <p>{{ currentOrder.remark || '无' }}</p>
      </div>
    </el-dialog>

    <el-dialog
      v-model="shipVisible"
      title="订单发货"
      width="500px"
    >
      <el-form
        ref="shipFormRef"
        :model="shipForm"
        :rules="shipRules"
        label-width="100px"
      >
        <el-form-item label="物流公司" prop="shipping_company">
          <el-select v-model="shipForm.shipping_company" placeholder="请选择物流公司" style="width: 100%">
            <el-option label="顺丰速运" value="顺丰速运" />
            <el-option label="京东物流" value="京东物流" />
            <el-option label="中通快递" value="中通快递" />
            <el-option label="圆通速递" value="圆通速递" />
            <el-option label="申通快递" value="申通快递" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号" prop="shipping_no">
          <el-input v-model="shipForm.shipping_no" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" @click="submitShip" :loading="shipLoading">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const shipLoading = ref(false)
const detailVisible = ref(false)
const shipVisible = ref(false)
const shipFormRef = ref(null)
const currentOrder = ref(null)

const searchForm = reactive({
  order_no: '',
  status: null
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const tableData = ref([
  {
    id: 1,
    order_no: 'ORD202401150001',
    username: '小明',
    total_amount: 128.5,
    status: 1,
    created_at: '2024-01-15 10:30:00',
    items: [
      { product_name: '新鲜有机青菜 500g', price: 9.9, quantity: 2, product_image: '' },
      { product_name: '红富士苹果 5斤装', price: 29.9, quantity: 3, product_image: '' }
    ],
    address: { name: '小明', phone: '138****8001', province: '北京市', city: '北京市', district: '朝阳区', detail: '某某街道123号' }
  },
  {
    id: 2,
    order_no: 'ORD202401150002',
    username: '小红',
    total_amount: 256.0,
    status: 2,
    created_at: '2024-01-14 15:20:00',
    items: [{ product_name: '农家土鸡蛋 30枚', price: 39.9, quantity: 5, product_image: '' }],
    address: { name: '小红', phone: '138****8002', province: '上海市', city: '上海市', district: '浦东新区', detail: '某某路456号' }
  }
])

const shipForm = reactive({
  shipping_company: '',
  shipping_no: ''
})

const shipRules = {
  shipping_company: [{ required: true, message: '请选择物流公司', trigger: 'change' }],
  shipping_no: [{ required: true, message: '请输入物流单号', trigger: 'blur' }]
}

const statusMap = {
  0: { text: '待付款', type: 'warning' },
  1: { text: '待发货', type: 'info' },
  2: { text: '待收货', type: 'primary' },
  3: { text: '已完成', type: 'success' },
  4: { text: '已取消', type: 'danger' }
}

const getStatusText = (status) => statusMap[status]?.text || '未知'
const getStatusType = (status) => statusMap[status]?.type || ''

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/order/list',
      method: 'get',
      params: {
        ...searchForm,
        page: pagination.page,
        page_size: pagination.page_size
      }
    })
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取订单列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.order_no = ''
  searchForm.status = null
  pagination.page = 1
  fetchData()
}

const handleView = (row) => {
  currentOrder.value = row
  detailVisible.value = true
}

const handleShip = (row) => {
  currentOrder.value = row
  shipForm.shipping_company = ''
  shipForm.shipping_no = ''
  shipVisible.value = true
}

const submitShip = async () => {
  const valid = await shipFormRef.value.validate().catch(() => false)
  if (!valid) return

  shipLoading.value = true
  try {
    const res = await request({
      url: `/admin/order/ship/${currentOrder.value.id}`,
      method: 'put',
      data: shipForm
    })
    if (res.code === 200) {
      ElMessage.success('发货成功')
      shipVisible.value = false
      fetchData()
    } else {
      ElMessage.error(res.msg || '发货失败')
    }
  } catch (error) {
    ElMessage.error('发货失败')
  } finally {
    shipLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.order-detail {
  padding: 10px 0;
}

.address-info {
  line-height: 2;
}
</style>
