<template>
  <div class="page-container order-page">
    <van-nav-bar title="我的订单" :left-arrow="true" @click-left="$router.back()" fixed />

    <van-tabs v-model="activeTab" sticky>
      <van-tab title="全部" />
      <van-tab title="待付款" />
      <van-tab title="待发货" />
      <van-tab title="待收货" />
      <van-tab title="已完成" />
    </van-tabs>

    <div class="order-content">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div v-for="order in orders" :key="order.id" class="order-card">
          <van-cell-group>
            <van-cell :title="`订单号: ${order.order_no}`" :value="getOrderStatusText(order.status)">
              <template #label>{{ order.created_at }}</template>
            </van-cell>
            
            <van-cell v-for="item in order.items" :key="item.id">
              <template #default>
                <div class="order-item">
                  <img
                    :src="item.product_image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetable&image_size=square'"
                    class="product-img"
                  />
                  <div class="product-info">
                    <div class="product-name">{{ item.product_name }}</div>
                    <div class="product-price">
                      <span>¥{{ item.price }}</span>
                      <span>x{{ item.quantity }}</span>
                    </div>
                  </div>
                </div>
              </template>
            </van-cell>

            <van-cell>
              <template #default>
                <div class="order-footer">
                  <div class="total">
                    共{{ order.items.length }}件商品
                    <span class="total-price">实付: ¥{{ order.total_amount }}</span>
                  </div>
                  <div class="actions">
                    <van-button
                      v-if="order.status === 0"
                      size="small"
                      type="danger"
                      @click="payOrder(order)"
                    >去支付</van-button>
                    <van-button
                      v-if="order.status === 0"
                      size="small"
                      @click="cancelOrder(order)"
                    >取消订单</van-button>
                    <van-button
                      v-if="order.status === 2"
                      size="small"
                      type="primary"
                      @click="receiveOrder(order)"
                    >确认收货</van-button>
                    <van-button
                      v-if="order.status === 3"
                      size="small"
                      @click="goComment(order)"
                    >去评价</van-button>
                  </div>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </div>

        <van-empty v-if="orders.length === 0 && !loading" description="暂无订单" />
      </van-list>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getOrderList, cancelOrder as apiCancelOrder, payOrder as apiPayOrder, receiveOrder as apiReceiveOrder } from '@/api/order'

const route = useRoute()
const router = useRouter()

const activeTab = ref(0)
const orders = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = ref(10)

const statusMap = {
  0: { text: '待付款', value: 'pending' },
  1: { text: '待发货', value: 'paid' },
  2: { text: '待收货', value: 'shipped' },
  3: { text: '已完成', value: 'completed' },
  4: { text: '已取消', value: 'cancelled' }
}

const getOrderStatusText = (status) => {
  return statusMap[status]?.text || '未知'
}

const getStatusParam = () => {
  const statusParams = [null, 'pending', 'paid', 'shipped', 'completed']
  return statusParams[activeTab.value]
}

const onLoad = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    
    const statusParam = getStatusParam()
    if (statusParam) {
      params.status = statusParam
    }
    
    const res = await getOrderList(params)
    if (res.code === 200) {
      const newOrders = res.data.list || []
      orders.value = [...orders.value, ...newOrders]
      
      if (newOrders.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (error) {
    console.error('加载订单失败', error)
  } finally {
    loading.value = false
  }
}

const refreshOrders = () => {
  page.value = 1
  orders.value = []
  finished.value = false
  onLoad()
}

watch(activeTab, () => {
  refreshOrders()
})

const payOrder = async (order) => {
  try {
    await showConfirmDialog({
      title: '支付确认',
      message: `确定支付 ¥${order.total_amount} 吗？`
    })
    const res = await apiPayOrder(order.id)
    if (res.code === 200) {
      showToast('支付成功')
      refreshOrders()
    } else {
      showToast(res.msg || '支付失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('支付失败')
    }
  }
}

const cancelOrder = async (order) => {
  try {
    await showConfirmDialog({
      title: '取消订单',
      message: '确定要取消该订单吗？'
    })
    const res = await apiCancelOrder(order.id)
    if (res.code === 200) {
      showToast('已取消')
      refreshOrders()
    } else {
      showToast(res.msg || '取消失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('取消失败')
    }
  }
}

const receiveOrder = async (order) => {
  try {
    await showConfirmDialog({
      title: '确认收货',
      message: '确认已收到商品吗？'
    })
    const res = await apiReceiveOrder(order.id)
    if (res.code === 200) {
      showToast('确认成功')
      refreshOrders()
    } else {
      showToast(res.msg || '操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('操作失败')
    }
  }
}

const goComment = (order) => {
  router.push(`/order/${order.id}/comment`)
}

onMounted(() => {
  if (route.query.status) {
    const statusMap = { pending: 1, paid: 2, shipped: 3, completed: 4 }
    activeTab.value = statusMap[route.query.status] || 0
  }
  refreshOrders()
})
</script>

<style scoped>
.order-page {
  padding-top: 46px;
  padding-bottom: 20px;
}

.order-card {
  margin-bottom: 10px;
}

.order-item {
  display: flex;
  width: 100%;
}

.product-img {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
}

.product-info {
  flex: 1;
  margin-left: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
  display: flex;
  justify-content: space-between;
}

.product-price span:first-child {
  color: #f56c6c;
}

.order-footer {
  width: 100%;
}

.total {
  font-size: 12px;
  color: #999;
  margin-bottom: 10px;
}

.total-price {
  color: #f56c6c;
  font-size: 14px;
  font-weight: bold;
  margin-left: 10px;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
