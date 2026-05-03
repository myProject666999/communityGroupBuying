import axios from 'axios'
import { showToast } from 'vant'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response) => {
    const res = response.data
    return res
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          showToast('登录已过期，请重新登录')
          router.push('/login')
          break
        case 403:
          showToast('无权限访问')
          break
        case 404:
          showToast('请求资源不存在')
          break
        case 500:
          showToast('服务器错误')
          break
        default:
          showToast(error.response.data?.msg || '请求失败')
      }
    } else {
      showToast('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default request
