<template>
  <div class="page-container knowledge-page">
    <van-nav-bar title="农产品知识" :left-arrow="true" @click-left="$router.back()" fixed />

    <div class="content">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group v-if="list.length > 0">
          <van-cell
            v-for="item in list"
            :key="item.id"
            is-link
            @click="goDetail(item)"
          >
            <template #default>
              <div class="list-item">
                <img
                  :src="item.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=agriculture%20knowledge&image_size=square'"
                  class="item-img"
                />
                <div class="item-info">
                  <div class="item-title">{{ item.title }}</div>
                  <div class="item-desc">{{ item.description || item.content }}</div>
                  <div class="item-meta">
                    <span>{{ item.author || '管理员' }}</span>
                    <span>{{ item.created_at }}</span>
                    <span>阅读 {{ item.views || 0 }}</span>
                  </div>
                </div>
              </div>
            </template>
          </van-cell>
        </van-cell-group>

        <van-empty v-else-if="!loading" description="暂无内容" />
      </van-list>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const list = ref([
  {
    id: 1,
    title: '有机蔬菜的营养价值与功效',
    description: '有机蔬菜富含维生素和矿物质，对人体健康有很多益处...',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=organic%20vegetables%20nutritious&image_size=square',
    author: '农业专家',
    views: 1256,
    created_at: '2024-01-15'
  },
  {
    id: 2,
    title: '如何挑选新鲜的水果',
    description: '教你几招挑选新鲜水果的小技巧，让你每次都能买到最好的...',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20fruits%20selection&image_size=square',
    author: '营养师',
    views: 892,
    created_at: '2024-01-14'
  },
  {
    id: 3,
    title: '时令蔬菜的最佳食用季节',
    description: '了解各种蔬菜的最佳食用季节，让你的餐桌更加丰富多彩...',
    image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=seasonal%20vegetables%20chart&image_size=square',
    author: '美食博主',
    views: 1567,
    created_at: '2024-01-13'
  }
])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = ref(10)

const onLoad = async () => {
  loading.value = false
  finished.value = true
}

const goDetail = (item) => {
  router.push(`/knowledge/${item.id}`)
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.knowledge-page {
  padding-top: 46px;
  background: #f5f5f5;
}

.content {
  padding: 10px 0;
}

.list-item {
  display: flex;
  width: 100%;
}

.item-img {
  width: 100px;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.item-info {
  flex: 1;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.item-title {
  font-size: 15px;
  font-weight: bold;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-desc {
  font-size: 13px;
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-meta {
  font-size: 12px;
  color: #999;
  display: flex;
  gap: 15px;
}
</style>
