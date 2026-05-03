<template>
  <div class="page-container login-page">
    <van-nav-bar title="登录" :left-arrow="true" @click-left="$router.back()" />
    
    <div class="logo-area">
      <van-icon name="shopping-cart" size="80" color="#1989fa" />
      <h2>农产品社区团购</h2>
    </div>

    <div class="form-area">
      <van-cell-group inset>
        <van-field
          v-model="form.username"
          label="用户名"
          placeholder="请输入用户名"
          left-icon="user-o"
        />
        <van-field
          v-model="form.password"
          type="password"
          label="密码"
          placeholder="请输入密码"
          left-icon="lock"
        />
      </van-cell-group>

      <div class="button-area">
        <van-button type="primary" block :loading="loading" @click="handleLogin">
          登录
        </van-button>
        <div class="links">
          <span class="link" @click="$router.push('/register')">立即注册</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!form.username || !form.password) {
    showToast('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const success = await userStore.login(form)
    if (success) {
      const redirect = route.query.redirect || '/home'
      router.replace(redirect)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  background: linear-gradient(to bottom, #e3f2fd, #ffffff);
  min-height: 100vh;
}

.logo-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 0 40px;
}

.logo-area h2 {
  margin-top: 15px;
  font-size: 20px;
  color: #333;
  font-weight: 500;
}

.form-area {
  padding: 0 20px;
}

.button-area {
  margin-top: 30px;
}

.links {
  display: flex;
  justify-content: flex-end;
  margin-top: 15px;
}

.link {
  color: #1989fa;
  font-size: 14px;
  cursor: pointer;
}
</style>
