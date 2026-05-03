import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userLogin, userRegister, getUserInfo } from '@/api/user'
import { showToast } from 'vant'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUserInfo(info) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  async function login(form) {
    try {
      const res = await userLogin(form)
      if (res.code === 200) {
        setToken(res.data.token)
        setUserInfo({
          user_id: res.data.user_id,
          username: res.data.username,
          nickname: res.data.nickname,
          avatar: res.data.avatar
        })
        showToast('登录成功')
        return true
      } else {
        showToast(res.msg || '登录失败')
        return false
      }
    } catch (error) {
      showToast('登录失败')
      return false
    }
  }

  async function register(form) {
    try {
      const res = await userRegister(form)
      if (res.code === 200) {
        showToast('注册成功，请登录')
        return true
      } else {
        showToast(res.msg || '注册失败')
        return false
      }
    } catch (error) {
      showToast('注册失败')
      return false
    }
  }

  async function fetchUserInfo() {
    try {
      const res = await getUserInfo()
      if (res.code === 200) {
        setUserInfo(res.data)
        return res.data
      }
    } catch (error) {
      console.error('获取用户信息失败', error)
    }
    return null
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    showToast('已退出登录')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    setToken,
    setUserInfo,
    login,
    register,
    fetchUserInfo,
    logout
  }
})
