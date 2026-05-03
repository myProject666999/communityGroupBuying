import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/category',
    name: 'Category',
    component: () => import('@/views/category/index.vue'),
    meta: { title: '分类' }
  },
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('@/views/cart/index.vue'),
    meta: { title: '购物车', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/index.vue'),
    meta: { title: '我的' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/product/:id',
    name: 'ProductDetail',
    component: () => import('@/views/product-detail/index.vue'),
    meta: { title: '商品详情' }
  },
  {
    path: '/order',
    name: 'Order',
    component: () => import('@/views/order/index.vue'),
    meta: { title: '我的订单', requiresAuth: true }
  },
  {
    path: '/order/create',
    name: 'OrderCreate',
    component: () => import('@/views/order-create/index.vue'),
    meta: { title: '确认订单', requiresAuth: true }
  },
  {
    path: '/address',
    name: 'Address',
    component: () => import('@/views/address/index.vue'),
    meta: { title: '收货地址', requiresAuth: true }
  },
  {
    path: '/favorite',
    name: 'Favorite',
    component: () => import('@/views/favorite/index.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: () => import('@/views/knowledge/index.vue'),
    meta: { title: '农产品知识' }
  },
  {
    path: '/community',
    name: 'Community',
    component: () => import('@/views/community/index.vue'),
    meta: { title: '社区分享' }
  },
  {
    path: '/recharge',
    name: 'Recharge',
    component: () => import('@/views/recharge/index.vue'),
    meta: { title: '充值', requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/settings/index.vue'),
    meta: { title: '设置', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = to.meta.title || '农产品社区团购'
  
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
})

export default router
