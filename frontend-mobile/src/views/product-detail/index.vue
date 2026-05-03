<template>
  <div class="page-container product-page">
    <van-nav-bar
      title="商品详情"
      :left-arrow="true"
      @click-left="$router.back()"
      fixed
    />

    <div class="content">
      <van-swipe class="banner" indicator-color="white" autoplay="3000">
        <van-swipe-item>
          <img
            :src="product.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetable&image_size=square_hd'"
            class="banner-img"
          />
        </van-swipe-item>
      </van-swipe>

      <van-cell-group>
        <van-cell>
          <template #default>
            <div class="price-area">
              <span class="price">¥{{ product.price }}</span>
              <span class="market-price" v-if="product.market_price">¥{{ product.market_price }}</span>
              <van-tag v-if="product.is_new" type="warning" class="tag">新品</van-tag>
              <van-tag v-if="product.is_recommend" type="primary" class="tag">推荐</van-tag>
            </div>
            <div class="product-name">{{ product.name }}</div>
            <div class="sales-info">
              <span>销量: {{ product.sales || 0 }}</span>
              <span>库存: {{ product.stock || 0 }}</span>
            </div>
          </template>
        </van-cell>
      </van-cell-group>

      <van-cell-group>
        <van-cell title="商品介绍" />
        <div class="description" v-html="product.description"></div>
      </van-cell-group>

      <div class="section-title">商品评价</div>
      <van-list
        v-model:loading="commentLoading"
        :finished="commentFinished"
        finished-text="没有更多了"
        @load="loadComments"
      >
        <van-cell-group v-if="comments.length > 0">
          <van-cell v-for="item in comments" :key="item.id">
            <template #default>
              <div class="comment-item">
                <div class="comment-header">
                  <van-avatar :size="32">
                    <van-icon name="user-o" />
                  </van-avatar>
                  <div class="user-info">
                    <div class="username">{{ item.nickname || '用户' }}</div>
                    <van-rate v-model="item.rating" readonly size="14" />
                  </div>
                  <div class="date">{{ item.created_at }}</div>
                </div>
                <div class="comment-content">{{ item.content }}</div>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
        <van-empty v-else-if="!commentLoading" description="暂无评价" />
      </van-list>
    </div>

    <van-goods-action>
      <van-goods-action-icon
        icon="star-o"
        text="收藏"
        @click="toggleFavorite"
      />
      <van-goods-action-icon
        icon="shopping-cart-o"
        text="购物车"
        :badge="cartCount"
        @click="$router.push('/cart')"
      />
      <van-goods-action-button type="warning" text="加入购物车" @click="addToCart" />
      <van-goods-action-button type="danger" text="立即购买" @click="buyNow" />
    </van-goods-action>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getProductDetail } from '@/api/product'
import { addCart } from '@/api/cart'
import { addFavorite, deleteFavorite, getFavoriteList } from '@/api/favorite'
import { getCommentList } from '@/api/comment'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const product = ref({})
const isFavorited = ref(false)
const favoriteId = ref(null)

const comments = ref([])
const commentLoading = ref(false)
const commentFinished = ref(false)
const commentPage = ref(1)
const pageSize = ref(10)

const cartCount = computed(() => {
  const count = localStorage.getItem('cartCount')
  return count ? parseInt(count) : 0
})

const fetchProduct = async () => {
  try {
    const res = await getProductDetail(route.params.id)
    if (res.code === 200) {
      product.value = res.data
    }
  } catch (error) {
    console.error('获取商品详情失败', error)
    product.value = {
      id: route.params.id,
      name: '新鲜有机蔬菜 500g',
      price: 9.9,
      market_price: 15.9,
      sales: 100,
      stock: 50,
      description: '<p>这是一款新鲜的有机蔬菜，来自农场直供，绿色健康，营养丰富。</p><p>适合炒菜、凉拌等多种烹饪方式。</p>'
    }
  }
}

const checkFavorite = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const res = await getFavoriteList({ page: 1, page_size: 100 })
    if (res.code === 200) {
      const item = res.data.find(f => f.product_id === parseInt(route.params.id))
      if (item) {
        isFavorited.value = true
        favoriteId.value = item.id
      }
    }
  } catch (error) {
    console.error('检查收藏状态失败', error)
  }
}

const toggleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  try {
    if (isFavorited.value) {
      const res = await deleteFavorite(favoriteId.value)
      if (res.code === 200) {
        isFavorited.value = false
        showToast('已取消收藏')
      }
    } else {
      const res = await addFavorite({ product_id: parseInt(route.params.id) })
      if (res.code === 200) {
        isFavorited.value = true
        favoriteId.value = res.data.id
        showToast('已收藏')
      }
    }
  } catch (error) {
    showToast('操作失败')
  }
}

const addToCart = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  try {
    const res = await addCart({ product_id: product.value.id, quantity: 1 })
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

const buyNow = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  const items = [{
    product_id: product.value.id,
    product_name: product.value.name,
    product_image: product.value.image,
    price: product.value.price,
    quantity: 1
  }]
  router.push({
    path: '/order/create',
    query: {
      items: JSON.stringify(items),
      single: 'true'
    }
  })
}

const loadComments = async () => {
  try {
    const res = await getCommentList({
      product_id: route.params.id,
      page: commentPage.value,
      page_size: pageSize.value
    })
    if (res.code === 200) {
      const newComments = res.data.list || []
      comments.value = [...comments.value, ...newComments]
      
      if (newComments.length < pageSize.value) {
        commentFinished.value = true
      } else {
        commentPage.value++
      }
    }
  } catch (error) {
    console.error('加载评论失败', error)
  } finally {
    commentLoading.value = false
  }
}

onMounted(() => {
  fetchProduct()
  checkFavorite()
})
</script>

<style scoped>
.product-page {
  padding-top: 46px;
  padding-bottom: 50px;
}

.banner {
  height: 300px;
}

.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.price-area {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.price {
  font-size: 24px;
  color: #f56c6c;
  font-weight: bold;
}

.market-price {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
  margin-left: 10px;
}

.tag {
  margin-left: 10px;
}

.product-name {
  font-size: 16px;
  color: #333;
  margin-top: 10px;
}

.sales-info {
  font-size: 12px;
  color: #999;
  margin-top: 10px;
  display: flex;
  gap: 20px;
}

.description {
  padding: 15px;
  background: #fff;
  font-size: 14px;
  line-height: 1.8;
  color: #666;
}

.section-title {
  padding: 12px 16px;
  font-size: 14px;
  font-weight: bold;
  color: #333;
  background: #f5f5f5;
}

.comment-item {
  width: 100%;
}

.comment-header {
  display: flex;
  align-items: center;
}

.comment-header .user-info {
  margin-left: 10px;
  flex: 1;
}

.comment-header .username {
  font-size: 14px;
  color: #333;
}

.comment-header .date {
  font-size: 12px;
  color: #999;
}

.comment-content {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}
</style>
