<template>
  <div class="page-container profile-page">
    <van-nav-bar title="我的" fixed />

    <div class="user-info">
      <div class="avatar-area" v-if="userStore.isLoggedIn">
        <van-avatar
          :src="userStore.userInfo.avatar || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=default%20user%20avatar&image_size=square'"
          size="60"
        />
        <div class="user-text">
          <div class="nickname">{{ userStore.userInfo.nickname || userStore.userInfo.username }}</div>
          <div class="integral">积分: {{ userInfo.integral || 0 }}</div>
        </div>
      </div>
      <div class="avatar-area" v-else @click="$router.push('/login')">
        <van-avatar icon="user-o" size="60" />
        <div class="user-text">
          <div class="nickname">点击登录</div>
          <div class="hint">登录后享受更多服务</div>
        </div>
      </div>
    </div>

    <van-cell-group>
      <van-cell title="我的订单" is-link @click="$router.push('/order')">
        <template #icon>
          <van-icon name="orders-o" color="#1989fa" />
        </template>
      </van-cell>
    </van-cell-group>

    <div class="order-tabs">
      <div class="tab-item" @click="goOrder('pending')">
        <van-icon name="todo-list" size="28" />
        <span>待付款</span>
      </div>
      <div class="tab-item" @click="goOrder('paid')">
        <van-icon name="gift-o" size="28" />
        <span>待发货</span>
      </div>
      <div class="tab-item" @click="goOrder('shipped')">
        <van-icon name="logistics" size="28" />
        <span>待收货</span>
      </div>
      <div class="tab-item" @click="goOrder('completed')">
        <van-icon name="completed" size="28" />
        <span>已完成</span>
      </div>
    </div>

    <van-cell-group>
      <van-cell title="收货地址" is-link @click="$router.push('/address')">
        <template #icon>
          <van-icon name="location-o" color="#1989fa" />
        </template>
      </van-cell>
      <van-cell title="我的收藏" is-link @click="$router.push('/favorite')">
        <template #icon>
          <van-icon name="star-o" color="#ff976a" />
        </template>
      </van-cell>
      <van-cell title="账户充值" is-link @click="$router.push('/recharge')">
        <template #icon>
          <van-icon name="wallet-o" color="#07c160" />
        </template>
      </van-cell>
    </van-cell-group>

    <van-cell-group>
      <van-cell title="社区分享" is-link @click="$router.push('/community')">
        <template #icon>
          <van-icon name="friends-o" color="#1989fa" />
        </template>
      </van-cell>
      <van-cell title="农产品知识" is-link @click="$router.push('/knowledge')">
        <template #icon>
          <van-icon name="bookmark-o" color="#1989fa" />
        </template>
      </van-cell>
      <van-cell title="资讯中心" is-link @click="$router.push('/news')">
        <template #icon>
          <van-icon name="newspaper-o" color="#1989fa" />
        </template>
      </van-cell>
      <van-cell title="论坛社区" is-link @click="$router.push('/forum')">
        <template #icon>
          <van-icon name="chat-o" color="#1989fa" />
        </template>
      </van-cell>
    </van-cell-group>

    <van-cell-group v-if="userStore.isLoggedIn">
      <van-cell title="个人设置" is-link @click="$router.push('/settings')">
        <template #icon>
          <van-icon name="setting-o" color="#999" />
        </template>
      </van-cell>
      <van-cell title="退出登录" is-link @click="handleLogout">
        <template #icon>
          <van-icon name="logout" color="#f56c6c" />
        </template>
      </van-cell>
    </van-cell-group>

    <van-tabbar v-model="activeTab">
      <van-tabbar-item icon="home-o" @click="$router.push('/home')">首页</van-tabbar-item>
      <van-tabbar-item icon="apps-o" @click="$router.push('/category')">分类</van-tabbar-item>
      <van-tabbar-item icon="shopping-cart-o" :badge="cartCount" @click="$router.push('/cart')">购物车</van-tabbar-item>
      <van-tabbar-item icon="user-o" @click="$router.push('/profile')">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { useUserStore } from '@/store/user'
import { getUserInfo } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref(3)
const userInfo = ref({ integral: 0 })

const cartCount = computed(() => {
  const count = localStorage.getItem('cartCount')
  return count ? parseInt(count) : 0
})

const goOrder = (status) => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  router.push(`/order?status=${status}`)
}

const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '提示',
      message: '确定要退出登录吗？'
    })
    userStore.logout()
  } catch (error) {
    if (error !== 'cancel') {
      showToast('操作失败')
    }
  }
}

const fetchUserInfo = async () => {
  if (userStore.isLoggedIn) {
    try {
      const res = await getUserInfo()
      if (res.code === 200) {
        userInfo.value = res.data
        userStore.setUserInfo(res.data)
      }
    } catch (error) {
      console.error('获取用户信息失败', error)
    }
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.profile-page {
  padding-top: 46px;
  padding-bottom: 50px;
  background: #f5f5f5;
}

.user-info {
  background: linear-gradient(135deg, #1989fa 0%, #41b883 100%);
  padding: 20px;
}

.avatar-area {
  display: flex;
  align-items: center;
}

.user-text {
  margin-left: 15px;
  color: #fff;
}

.nickname {
  font-size: 18px;
  font-weight: bold;
}

.integral, .hint {
  font-size: 12px;
  margin-top: 5px;
  opacity: 0.9;
}

.order-tabs {
  display: flex;
  background: #fff;
  padding: 15px 0;
  margin-bottom: 10px;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #333;
  font-size: 12px;
}

.tab-item span {
  margin-top: 5px;
}
</style>
