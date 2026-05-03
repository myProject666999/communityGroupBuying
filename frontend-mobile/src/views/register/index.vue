<template>
  <div class="page-container register-page">
    <van-nav-bar title="注册" :left-arrow="true" @click-left="$router.back()" />

    <div class="logo-area">
      <van-icon name="user-plus" size="60" color="#1989fa" />
      <h2>新用户注册</h2>
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
          v-model="form.nickname"
          label="昵称"
          placeholder="请输入昵称"
          left-icon="contact"
        />
        <van-field
          v-model="form.phone"
          label="手机号"
          placeholder="请输入手机号"
          left-icon="phone-o"
        />
        <van-field
          v-model="form.password"
          type="password"
          label="密码"
          placeholder="请输入密码"
          left-icon="lock"
        />
        <van-field
          v-model="form.confirm_password"
          type="password"
          label="确认密码"
          placeholder="请再次输入密码"
          left-icon="lock"
        />
      </van-cell-group>

      <div class="button-area">
        <van-button type="primary" block :loading="loading" @click="handleRegister">
          注册
        </van-button>
        <div class="links">
          <span class="link" @click="$router.push('/login')">已有账号？立即登录</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const form = reactive({
  username: '',
  nickname: '',
  phone: '',
  password: '',
  confirm_password: ''
})

const handleRegister = async () => {
  if (!form.username || !form.password) {
    showToast('请填写必要信息')
    return
  }
  
  if (form.password !== form.confirm_password) {
    showToast('两次密码输入不一致')
    return
  }
  
  loading.value = true
  try {
    const success = await userStore.register(form)
    if (success) {
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  background: linear-gradient(to bottom, #e3f2fd, #ffffff);
  min-height: 100vh;
}

.logo-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 0 30px;
}

.logo-area h2 {
  margin-top: 15px;
  font-size: 18px;
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
  justify-content: center;
  margin-top: 15px;
}

.link {
  color: #1989fa;
  font-size: 14px;
  cursor: pointer;
}
</style>
