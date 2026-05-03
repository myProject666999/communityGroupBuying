<template>
  <div class="page-container cart-page">
    <van-nav-bar title="购物车" fixed />

    <div v-if="cartList.length === 0" class="empty-cart">
      <van-empty description="购物车是空的" />
      <van-button type="primary" @click="$router.push('/home')">去逛逛</van-button>
    </div>

    <div v-else class="cart-content">
      <van-checkbox-group v-model="selectedIds" ref="checkboxGroup">
        <van-cell-group>
          <van-checkbox
            v-for="item in cartList"
            :key="item.id"
            :name="item.id"
            class="cart-item"
          >
            <div class="cart-item-content">
              <img
                :src="item.product_image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetable&image_size=square'"
                class="product-img"
                @click.stop="goProduct(item.product_id)"
              />
              <div class="product-info">
                <div class="product-name">{{ item.product_name }}</div>
                <div class="product-price">¥{{ item.price }}</div>
                <van-stepper
                  v-model="item.quantity"
                  :min="1"
                  :max="99"
                  @change="onQuantityChange(item)"
                />
              </div>
              <van-icon name="delete" class="delete-icon" @click="deleteItem(item.id)" />
            </div>
          </van-checkbox>
        </van-cell-group>
      </van-checkbox-group>
    </div>

    <van-submit-bar
      :price="totalPrice * 100"
      button-text="结算"
      @submit="onSubmit"
    >
      <van-checkbox v-model="allSelected" @change="toggleAllSelect">全选</van-checkbox>
    </van-submit-bar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getCartList, updateCart, deleteCart } from '@/api/cart'

const router = useRouter()

const cartList = ref([])
const selectedIds = ref([])
const allSelected = ref(false)
const checkboxGroup = ref(null)

const totalPrice = computed(() => {
  return cartList.value
    .filter(item => selectedIds.value.includes(item.id))
    .reduce((total, item) => total + item.price * item.quantity, 0)
})

const goProduct = (id) => {
  router.push(`/product/${id}`)
}

const toggleAllSelect = (val) => {
  if (val) {
    selectedIds.value = cartList.value.map(item => item.id)
  } else {
    selectedIds.value = []
  }
}

watch(selectedIds, () => {
  allSelected.value = selectedIds.value.length === cartList.value.length && cartList.value.length > 0
}, { deep: true })

const fetchCartList = async () => {
  try {
    const res = await getCartList()
    if (res.code === 200) {
      cartList.value = res.data
      localStorage.setItem('cartCount', cartList.value.length.toString())
    }
  } catch (error) {
    console.error('获取购物车失败', error)
  }
}

const onQuantityChange = async (item) => {
  try {
    const res = await updateCart({ id: item.id, quantity: item.quantity })
    if (res.code !== 200) {
      showToast('更新失败')
    }
  } catch (error) {
    showToast('更新失败')
  }
}

const deleteItem = async (id) => {
  try {
    await showConfirmDialog({
      title: '提示',
      message: '确定删除该商品？'
    })
    const res = await deleteCart(id)
    if (res.code === 200) {
      showToast('删除成功')
      fetchCartList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('删除失败')
    }
  }
}

const onSubmit = () => {
  if (selectedIds.value.length === 0) {
    showToast('请选择商品')
    return
  }
  const selectedItems = cartList.value.filter(item => selectedIds.value.includes(item.id))
  router.push({
    path: '/order/create',
    query: {
      cartIds: JSON.stringify(selectedIds.value),
      items: JSON.stringify(selectedItems)
    }
  })
}

onMounted(() => {
  fetchCartList()
})
</script>

<style scoped>
.cart-page {
  padding-top: 46px;
  padding-bottom: 50px;
}

.empty-cart {
  padding: 80px 20px;
  text-align: center;
}

.cart-content {
  padding-bottom: 50px;
}

.cart-item {
  padding: 10px;
}

.cart-item-content {
  display: flex;
  align-items: center;
}

.product-img {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.product-info {
  flex: 1;
  margin-left: 10px;
}

.product-name {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  font-size: 16px;
  color: #f56c6c;
  font-weight: bold;
  margin: 5px 0;
}

.delete-icon {
  padding: 10px;
  color: #999;
  font-size: 18px;
}
</style>
