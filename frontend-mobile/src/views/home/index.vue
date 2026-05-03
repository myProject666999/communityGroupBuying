<template>
  <div class="page-container home-page">
    <van-nav-bar title="农产品社区团购" fixed />
    
    <van-search
      v-model="keyword"
      placeholder="搜索农产品"
      show-action
      @search="onSearch"
      @cancel="onSearch"
    />

    <van-swipe class="banner" indicator-color="white" autoplay="3000">
      <van-swipe-item v-for="item in banners" :key="item.id">
        <img :src="item.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetables%20and%20fruits%20banner&image_size=landscape_16_9'" alt="banner" class="banner-img" />
      </van-swipe-item>
    </van-swipe>

    <van-grid :column-num="5" class="category-grid">
      <van-grid-item v-for="item in categories" :key="item.id" @click="goCategory(item)">
        <van-icon name="shop-o" size="28" />
        <text>{{ item.name }}</text>
      </van-grid-item>
    </van-grid>

    <div class="section-title">
      <span>推荐商品</span>
    </div>

    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <div class="product-list">
        <van-card
          v-for="item in products"
          :key="item.id"
          :thumb="item.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20green%20vegetable&image_size=square'"
          :title="item.name"
          :price="item.price"
          :origin-price="item.market_price"
          @click="goProduct(item)"
        >
          <template #tags>
            <van-tag v-if="item.is_new" type="warning">新品</van-tag>
            <van-tag v-if="item.is_recommend" type="primary">推荐</van-tag>
          </template>
          <template #bottom>
            <van-button size="mini" type="primary" @click.stop="addToCart(item)">加入购物车</van-button>
          </template>
        </van-card>
      </div>
    </van-list>

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
import { showToast } from 'vant'
import { getBannerList, getCategoryList, getProductList, getRecommendProducts } from '@/api/product'
import { addCart } from '@/api/cart'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const keyword = ref('')
const banners = ref([])
const categories = ref([])
const products = ref([])
const loading = ref(false)
const finished = ref(false)
const activeTab = ref(0)
const page = ref(1)
const pageSize = ref(10)

const cartCount = computed(() => {
  const count = localStorage.getItem('cartCount')
  return count ? parseInt(count) : 0
})

const goCategory = (item) => {
  router.push(`/category?id=${item.id}`)
}

const goProduct = (item) => {
  router.push(`/product/${item.id}`)
}

const addToCart = async (item) => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  try {
    const res = await addCart({ product_id: item.id, quantity: 1 })
    if (res.code === 200) {
      showToast('已加入购物车')
      const count = parseInt(localStorage.getItem('cartCount') || '0') + 1
      localStorage.setItem('cartCount', count.toString())
    } else {
      showToast(res.msg || '添加失败')
    }
  } catch (error) {
    showToast('添加失败')
  }
}

const onSearch = () => {
  page.value = 1
  products.value = []
  finished.value = false
  onLoad()
}

const onLoad = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }
    
    const res = await getProductList(params)
    if (res.code === 200) {
      const newProducts = res.data.list || []
      products.value = [...products.value, ...newProducts]
      
      if (newProducts.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (error) {
    console.error('加载商品失败', error)
  } finally {
    loading.value = false
  }
}

const fetchBanners = async () => {
  try {
    const res = await getBannerList(1)
    if (res.code === 200) {
      banners.value = res.data.length > 0 ? res.data : [
        { id: 1, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetables%20and%20fruits%20promotion%20banner&image_size=landscape_16_9' }
      ]
    }
  } catch (error) {
    console.error('加载轮播图失败', error)
    banners.value = [
      { id: 1, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetables%20and%20fruits%20promotion%20banner&image_size=landscape_16_9' }
    ]
  }
}

const fetchCategories = async () => {
  try {
    const res = await getCategoryList()
    if (res.code === 200) {
      categories.value = res.data.length > 0 ? res.data : [
        { id: 1, name: '蔬菜' },
        { id: 2, name: '水果' },
        { id: 3, name: '肉类' },
        { id: 4, name: '海鲜' },
        { id: 5, name: '粮油' }
      ]
    }
  } catch (error) {
    console.error('加载分类失败', error)
    categories.value = [
      { id: 1, name: '蔬菜' },
      { id: 2, name: '水果' },
      { id: 3, name: '肉类' },
      { id: 4, name: '海鲜' },
      { id: 5, name: '粮油' }
    ]
  }
}

const fetchRecommendProducts = async () => {
  try {
    const res = await getRecommendProducts(20)
    if (res.code === 200 && res.data.length > 0) {
      products.value = res.data
      return
    }
  } catch (error) {
    console.error('加载推荐商品失败', error)
  }
  
  products.value = [
    { id: 1, name: '新鲜有机青菜 500g', price: 9.9, market_price: 15.9, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20green%20vegetable%20lettuce&image_size=square', is_recommend: 1 },
    { id: 2, name: '红富士苹果 5斤装', price: 29.9, market_price: 39.9, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=red%20fuji%20apples%20fresh&image_size=square', is_new: 1 },
    { id: 3, name: '农家土鸡蛋 30枚', price: 39.9, market_price: 49.9, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20farm%20eggs%20basket&image_size=square', is_recommend: 1 },
    { id: 4, name: '精选五花肉 500g', price: 35.9, market_price: 45.9, image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20pork%20belly%20meat&image_size=square' }
  ]
}

onMounted(() => {
  fetchBanners()
  fetchCategories()
  fetchRecommendProducts()
})
</script>

<style scoped>
.home-page {
  padding-top: 92px;
  padding-bottom: 50px;
}

.banner {
  height: 180px;
}

.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.category-grid {
  padding: 10px 0;
  background: #fff;
  margin-bottom: 10px;
}

.category-grid text {
  display: block;
  font-size: 12px;
  color: #333;
  margin-top: 5px;
}

.section-title {
  background: #fff;
  padding: 12px 16px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
}

.product-list {
  padding: 10px;
}
</style>
