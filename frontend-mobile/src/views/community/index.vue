<template>
  <div class="page-container community-page">
    <van-nav-bar title="社区分享" :left-arrow="true" @click-left="$router.back()" fixed>
      <template #right>
        <van-icon name="edit" size="18" @click="goPublish" />
      </template>
    </van-nav-bar>

    <div class="content">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div v-if="list.length > 0" class="post-list">
          <div class="post-item" v-for="item in list" :key="item.id">
            <div class="post-header">
              <van-avatar
                :src="item.avatar || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=default%20user%20avatar&image_size=square'"
                size="40"
              />
              <div class="user-info">
                <div class="username">{{ item.username || item.nickname }}</div>
                <div class="time">{{ item.created_at }}</div>
              </div>
            </div>
            <div class="post-content">{{ item.content }}</div>
            <div v-if="item.images && item.images.length > 0" class="post-images">
              <img
                v-for="(img, idx) in item.images.slice(0, 3)"
                :key="idx"
                :src="img"
                class="post-img"
              />
            </div>
            <div class="post-actions">
              <div class="action-item">
                <van-icon name="eye-o" />
                <span>{{ item.views || 0 }}</span>
              </div>
              <div class="action-item">
                <van-icon name="like-o" />
                <span>{{ item.likes || 0 }}</span>
              </div>
              <div class="action-item">
                <van-icon name="comment-o" />
                <span>{{ item.comments || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <van-empty v-else-if="!loading" description="暂无分享" />
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
    username: '小农民',
    avatar: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=farmer%20portrait&image_size=square',
    content: '今天收获了新鲜的西红柿，自家种植的，没有打农药，口感特别好！分享给大家看看~',
    images: [
      'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20red%20tomatoes%20basket&image_size=square'
    ],
    views: 325,
    likes: 45,
    comments: 12,
    created_at: '2024-01-15 10:30'
  },
  {
    id: 2,
    username: '果蔬达人',
    avatar: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=woman%20gardener&image_size=square',
    content: '推荐大家这个时令蔬菜，现在正是吃的好季节，营养丰富，价格也实惠！',
    images: [
      'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=seasonal%20vegetables%20display&image_size=square'
    ],
    views: 568,
    likes: 89,
    comments: 23,
    created_at: '2024-01-14 15:20'
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

const goPublish = () => {
  showToast('请先登录')
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.community-page {
  padding-top: 46px;
  background: #f5f5f5;
}

.content {
  padding: 10px 0;
}

.post-list {
  padding: 0 10px;
}

.post-item {
  background: #fff;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
}

.post-header {
  display: flex;
  align-items: center;
}

.user-info {
  margin-left: 10px;
}

.username {
  font-size: 15px;
  font-weight: bold;
  color: #333;
}

.time {
  font-size: 12px;
  color: #999;
  margin-top: 3px;
}

.post-content {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  margin: 12px 0;
}

.post-images {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.post-img {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 4px;
}

.post-actions {
  display: flex;
  justify-content: space-around;
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  color: #999;
}
</style>
