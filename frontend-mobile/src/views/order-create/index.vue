<template>
  <div class="page-container order-create-page">
    <van-nav-bar title="确认订单" :left-arrow="true" @click-left="$router.back()" fixed />

    <div class="content" v-if="items.length > 0">
      <van-cell-group>
        <van-cell title="收货地址" is-link @click="selectAddress">
          <template #default>
            <div v-if="selectedAddress" class="address-info">
              <div class="address-header">
                <span class="name">{{ selectedAddress.name }}</span>
                <span class="phone">{{ selectedAddress.phone }}</span>
              </div>
              <div class="address-detail">
                {{ selectedAddress.province }}{{ selectedAddress.city }}{{ selectedAddress.district }}{{ selectedAddress.detail }}
              </div>
            </div>
            <span v-else class="no-address">请选择收货地址</span>
          </template>
        </van-cell>
      </van-cell-group>

      <van-cell-group>
        <van-cell title="商品清单" />
        <van-cell v-for="item in items" :key="item.id || item.product_id">
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
      </van-cell-group>

      <van-cell-group>
        <van-field
          v-model="form.remark"
          type="textarea"
          label="备注"
          placeholder="选填，请输入订单备注"
          maxlength="100"
          show-word-limit
        />
      </van-cell-group>
    </div>

    <van-submit-bar
      :price="totalAmount * 100"
      button-text="提交订单"
      @submit="submitOrder"
    >
      <span>共{{ items.length }}件商品</span>
    </van-submit-bar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getAddressList } from '@/api/address'
import { createOrder } from '@/api/order'

const route = useRoute()
const router = useRouter()

const items = ref([])
const selectedAddress = ref(null)
const form = ref({
  remark: ''
})

const totalAmount = computed(() => {
  return items.value.reduce((total, item) => total + item.price * item.quantity, 0)
})

const selectAddress = () => {
  router.push({
    path: '/address',
    query: { from: 'order' }
  })
}

const submitOrder = async () => {
  if (!selectedAddress.value) {
    showToast('请选择收货地址')
    return
  }
  
  try {
    const orderItems = items.value.map(item => ({
      product_id: item.product_id,
      quantity: item.quantity,
      price: item.price
    }))
    
    const res = await createOrder({
      address_id: selectedAddress.value.id,
      items: orderItems,
      remark: form.value.remark
    })
    
    if (res.code === 200) {
      showToast('订单创建成功')
      router.replace({
        path: '/order',
        query: { status: 'pending' }
      })
    } else {
      showToast(res.msg || '下单失败')
    }
  } catch (error) {
    showToast('下单失败')
  }
}

onMounted(() => {
  if (route.query.items) {
    try {
      items.value = JSON.parse(route.query.items)
    } catch (error) {
      console.error('解析商品数据失败', error)
    }
  }
  
  getAddressList().then(res => {
    if (res.code === 200 && res.data.length > 0) {
      const defaultAddress = res.data.find(a => a.is_default) || res.data[0]
      selectedAddress.value = defaultAddress
    }
  }).catch(() => {})
})
</script>

<style scoped>
.order-create-page {
  padding-top: 46px;
  padding-bottom: 50px;
  background: #f5f5f5;
}

.content {
  padding: 10px 0;
}

.address-info {
  text-align: right;
}

.address-header {
  margin-bottom: 5px;
}

.name {
  font-weight: bold;
  color: #333;
}

.phone {
  margin-left: 10px;
  color: #666;
}

.address-detail {
  font-size: 12px;
  color: #999;
}

.no-address {
  color: #999;
  font-size: 14px;
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
</style>
