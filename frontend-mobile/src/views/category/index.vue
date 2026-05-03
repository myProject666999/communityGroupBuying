<template>
  <div class="page-container category-page">
    <van-nav-bar title="分类" fixed />

    <div class="content">
      <van-sidebar v-model="activeCategory" class="sidebar">
        <van-sidebar-item
          v-for="item in categories"
          :key="item.id"
          :title="item.name"
          @click="selectCategory(item)"
        />
      </van-sidebar>

      <div class="product-content">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
        >
          <div class="product-grid">
            <div
              class="product-item"
              v-for="item in products"
              :key="item.id"
              @click="goProduct(item)"
            >
              <img
                :src="item.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetable&image_size=square'"
                class="product-img"
              />
              <div class="product-info">
                <div class="product-name">{{ item.name }}</div>
                <div class="product-price">
                  <span class="price">¥{{ item.price }}</span>
                  <span class="market-price" v-if="item.market_price">¥{{ item.market_price }}</span>
                </div>
              </div>
            </div>
          </div>
        </van-list>
      </div>
    </div>

    <van-tabbar v-model="activeTab">
      <van-tabbar-item icon="home-o" @click="$router.push('/home')">首页</van-tabbar-item>
      <van-tabbar-item icon="apps-o" @click="$router.push('/category')">分类</van-tabbar-item>
      <van-tabbar-item icon="shopping-cart-o" :badge="cartCount" @click="$router.push('/cart')">购物车</van-tabbar-item>
      <van-tabbar-item icon="user-o" @click="$router.push('/profile')">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getCategoryList, getProductList } from '@/api/product'

const router = useRouter()
const route = useRoute()

const activeCategory = ref(0)
const activeTab = ref(1)
const categories = ref([])
const products = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = ref(20)
const currentCategoryId = ref(null)

const cartCount = computed(() => {
  const count = localStorage.getItem('cartCount')
  return count ? parseInt(count) : 0
})

const goProduct = (item) => {
  router.push(`/product/${item.id}`)
}

const selectCategory = (item) => {
  currentCategoryId.value = item.id
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
    
    if (currentCategoryId.value) {
      params.category_id = currentCategoryId.value
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

const fetchCategories = async () => {
  try {
    const res = await getCategoryList()
    if (res.code === 200 && res.data.length > 0) {
      categories.value = res.data
      if (route.query.id) {
        const index = res.data.findIndex(c => c.id === parseInt(route.query.id))
        if (index >= 0) {
          activeCategory.value = index
          currentCategoryId.value = res.data[index].id
        }
      } else {
        currentCategoryId.value = res.data[0].id
      }
      onLoad()
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

watch(activeCategory, (val) => {
  if (categories.value[val]) {
    selectCategory(categories.value[val])
  }
})

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-page {
  padding-top: 46px;
  padding-bottom: 50px;
}

.content {
  display: flex;
  height: calc(100vh - 96px);
}

.sidebar {
  width: 25%;
  height: 100%;
  overflow-y: auto;
}

.product-content {
  width: 75%;
  padding: 10px;
  overflow-y: auto;
}

.product-grid {
  display: flex;
  flex-wrap: wrap;
}

.product-item {
  width: 50%;
  padding: 5px;
  box-sizing: border-box;
}

.product-img {
  width: 100%;
  height: 100px;
  object-fit: cover;
  border-radius: 4px;
}

.product-info {
  padding: 5px;
}

.product-name {
  font-size: 12px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  margin-top: 5px;
}

.price {
  font-size: 14px;
  color: #f56c6c;
  font-weight: bold;
}

.market-price {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
  margin-left: 5px;
}
</style>
