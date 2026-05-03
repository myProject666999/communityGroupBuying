<template>
  <div class="page-container favorite-page">
    <van-nav-bar title="我的收藏" :left-arrow="true" @click-left="$router.back()" fixed />

    <div class="favorite-content">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div v-if="favorites.length > 0" class="product-grid">
          <div
            class="product-item"
            v-for="item in favorites"
            :key="item.id"
            @click="goProduct(item.product_id)"
          >
            <div class="product-wrap">
              <van-icon
                name="star"
                class="favorite-icon"
                @click.stop="removeFavorite(item)"
              />
              <img
                :src="item.product_image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetable&image_size=square'"
                class="product-img"
              />
              <div class="product-info">
                <div class="product-name">{{ item.product_name }}</div>
                <div class="product-price">
                  <span class="price">¥{{ item.price }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <van-empty v-else-if="!loading" description="暂无收藏商品" />
      </van-list>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getFavoriteList, deleteFavorite } from '@/api/favorite'

const router = useRouter()

const favorites = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = ref(20)

const onLoad = async () => {
  try {
    const res = await getFavoriteList({
      page: page.value,
      page_size: pageSize.value
    })
    if (res.code === 200) {
      const newFavorites = res.data.list || []
      favorites.value = [...favorites.value, ...newFavorites]
      
      if (newFavorites.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (error) {
    console.error('加载收藏失败', error)
  } finally {
    loading.value = false
  }
}

const goProduct = (id) => {
  router.push(`/product/${id}`)
}

const removeFavorite = async (item) => {
  try {
    await showConfirmDialog({
      title: '取消收藏',
      message: '确定要取消收藏吗？'
    })
    const res = await deleteFavorite(item.id)
    if (res.code === 200) {
      showToast('已取消收藏')
      const index = favorites.value.findIndex(f => f.id === item.id)
      if (index > -1) {
        favorites.value.splice(index, 1)
      }
    } else {
      showToast(res.msg || '操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('操作失败')
    }
  }
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.favorite-page {
  padding-top: 46px;
  padding-bottom: 20px;
  background: #f5f5f5;
}

.favorite-content {
  padding: 10px;
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

.product-wrap {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.favorite-icon {
  position: absolute;
  top: 10px;
  right: 10px;
  color: #ff976a;
  font-size: 18px;
  z-index: 10;
}

.product-img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.product-info {
  padding: 10px;
}

.product-name {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  margin-top: 5px;
}

.price {
  font-size: 16px;
  color: #f56c6c;
  font-weight: bold;
}
</style>
