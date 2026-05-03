<template>
  <div class="page-container recharge-page">
    <van-nav-bar title="账户充值" :left-arrow="true" @click-left="$router.back()" fixed />

    <div class="balance-area">
      <div class="balance-label">账户余额</div>
      <div class="balance-value">¥{{ userInfo.balance || 0 }}</div>
    </div>

    <div class="amount-area">
      <div class="title">选择充值金额</div>
      <div class="amount-grid">
        <div
          class="amount-item"
          :class="{ active: selectedAmount === item.value }"
          v-for="item in amountOptions"
          :key="item.value"
          @click="selectedAmount = item.value"
        >
          <div class="amount">¥{{ item.value }}</div>
          <div class="gift" v-if="item.gift">送¥{{ item.gift }}</div>
        </div>
      </div>
    </div>

    <div class="custom-amount">
      <van-field
        v-model="customAmount"
        type="number"
        label="自定义金额"
        placeholder="请输入充值金额"
        @focus="selectedAmount = null"
      />
    </div>

    <div class="pay-method">
      <div class="title">支付方式</div>
      <van-cell-group>
        <van-cell
          title="微信支付"
          is-link
          @click="payMethod = 'wechat'"
        >
          <template #icon>
            <van-icon name="wechat" color="#07c160" size="24" />
          </template>
          <template #right-icon>
            <van-icon :name="payMethod === 'wechat' ? 'checked' : 'circle'" :color="payMethod === 'wechat' ? '#1989fa' : '#c8c9cc'" />
          </template>
        </van-cell>
        <van-cell
          title="支付宝"
          is-link
          @click="payMethod = 'alipay'"
        >
          <template #icon>
            <van-icon name="alipay" color="#1677ff" size="24" />
          </template>
          <template #right-icon>
            <van-icon :name="payMethod === 'alipay' ? 'checked' : 'circle'" :color="payMethod === 'alipay' ? '#1989fa' : '#c8c9cc'" />
          </template>
        </van-cell>
      </van-cell-group>
    </div>

    <div class="submit-area">
      <van-button
        type="primary"
        block
        :loading="loading"
        @click="handleRecharge"
      >
        充值 ¥{{ totalAmount }}
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import { getUserInfo } from '@/api/user'

const userInfo = ref({})
const selectedAmount = ref(50)
const customAmount = ref('')
const payMethod = ref('wechat')
const loading = ref(false)

const amountOptions = [
  { value: 50, gift: 0 },
  { value: 100, gift: 5 },
  { value: 200, gift: 15 },
  { value: 500, gift: 50 },
  { value: 1000, gift: 120 }
]

const totalAmount = computed(() => {
  if (selectedAmount.value) {
    return selectedAmount.value
  }
  return parseFloat(customAmount.value) || 0
})

const fetchUserInfo = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 200) {
      userInfo.value = res.data
    }
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
}

const handleRecharge = async () => {
  if (totalAmount.value <= 0) {
    showToast('请选择充值金额')
    return
  }
  
  loading.value = true
  try {
    showToast('充值功能需要对接支付平台')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.recharge-page {
  padding-top: 46px;
  padding-bottom: 80px;
  background: #f5f5f5;
}

.balance-area {
  background: linear-gradient(135deg, #1989fa 0%, #41b883 100%);
  padding: 30px 20px;
  color: #fff;
  text-align: center;
}

.balance-label {
  font-size: 14px;
  opacity: 0.9;
}

.balance-value {
  font-size: 36px;
  font-weight: bold;
  margin-top: 10px;
}

.amount-area,
.pay-method {
  background: #fff;
  padding: 15px;
  margin-top: 10px;
}

.title {
  font-size: 14px;
  color: #333;
  font-weight: bold;
  margin-bottom: 15px;
}

.amount-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.amount-item {
  width: calc(33.33% - 7px);
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px 10px;
  text-align: center;
  cursor: pointer;
}

.amount-item.active {
  border-color: #1989fa;
  background: rgba(25, 137, 250, 0.05);
}

.amount {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.gift {
  font-size: 12px;
  color: #f56c6c;
  margin-top: 5px;
}

.custom-amount {
  background: #fff;
  margin-top: 10px;
}

.submit-area {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 15px;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}
</style>
