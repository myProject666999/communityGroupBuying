<template>
  <div class="page-container address-page">
    <van-nav-bar
      title="收货地址"
      :left-arrow="true"
      @click-left="$router.back()"
      fixed
    >
      <template #right>
        <van-icon name="plus" size="20" @click="goAdd" />
      </template>
    </van-nav-bar>

    <div class="address-content">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group v-if="addressList.length > 0">
          <van-cell
            v-for="item in addressList"
            :key="item.id"
            is-link
            @click="selectAddress(item)"
          >
            <template #default>
              <div class="address-item">
                <div class="address-header">
                  <span class="name">{{ item.name }}</span>
                  <span class="phone">{{ item.phone }}</span>
                  <van-tag v-if="item.is_default" type="primary" size="mini">默认</van-tag>
                </div>
                <div class="address-detail">
                  {{ item.province }}{{ item.city }}{{ item.district }}{{ item.detail }}
                </div>
                <div class="address-actions">
                  <van-button
                    size="mini"
                    type="primary"
                    @click.stop="setDefault(item)"
                    v-if="!item.is_default"
                  >设为默认</van-button>
                  <van-button
                    size="mini"
                    @click.stop="editAddress(item)"
                  >编辑</van-button>
                  <van-button
                    size="mini"
                    type="danger"
                    @click.stop="deleteAddress(item)"
                  >删除</van-button>
                </div>
              </div>
            </template>
          </van-cell>
        </van-cell-group>

        <van-empty v-else-if="!loading" description="暂无收货地址">
          <van-button type="primary" @click="goAdd">添加地址</van-button>
        </van-empty>
      </van-list>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getAddressList, deleteAddress, setDefaultAddress } from '@/api/address'

const route = useRoute()
const router = useRouter()

const addressList = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = ref(20)

const onLoad = async () => {
  try {
    const res = await getAddressList()
    if (res.code === 200) {
      addressList.value = res.data
      finished.value = true
    }
  } catch (error) {
    console.error('加载地址失败', error)
  } finally {
    loading.value = false
  }
}

const goAdd = () => {
  router.push('/address/add')
}

const editAddress = (item) => {
  router.push(`/address/edit/${item.id}`)
}

const selectAddress = (item) => {
  if (route.query.from === 'order') {
    router.back()
  }
}

const setDefault = async (item) => {
  try {
    const res = await setDefaultAddress(item.id)
    if (res.code === 200) {
      showToast('设置成功')
      onLoad()
    } else {
      showToast(res.msg || '设置失败')
    }
  } catch (error) {
    showToast('设置失败')
  }
}

const deleteAddress = async (item) => {
  try {
    await showConfirmDialog({
      title: '删除确认',
      message: '确定要删除该地址吗？'
    })
    const res = await deleteAddress(item.id)
    if (res.code === 200) {
      showToast('删除成功')
      onLoad()
    } else {
      showToast(res.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      showToast('删除失败')
    }
  }
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.address-page {
  padding-top: 46px;
  padding-bottom: 20px;
}

.address-content {
  padding: 10px;
}

.address-item {
  width: 100%;
}

.address-header {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
}

.name {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.phone {
  font-size: 14px;
  color: #666;
  margin-left: 10px;
}

.address-detail {
  font-size: 13px;
  color: #999;
  margin-bottom: 10px;
}

.address-actions {
  display: flex;
  gap: 10px;
}
</style>
