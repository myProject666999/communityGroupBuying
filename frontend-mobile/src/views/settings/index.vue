<template>
  <div class="page-container settings-page">
    <van-nav-bar title="个人设置" :left-arrow="true" @click-left="$router.back()" fixed />

    <van-cell-group>
      <van-cell title="头像" is-link>
        <template #default>
          <van-avatar
            :src="userInfo.avatar || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=default%20user%20avatar&image_size=square'"
            size="36"
          />
        </template>
      </van-cell>
      <van-cell title="昵称" :value="userInfo.nickname || '未设置'" is-link @click="showNicknameDialog = true" />
      <van-cell title="手机号" :value="userInfo.phone || '未设置'" is-link @click="showPhoneDialog = true" />
    </van-cell-group>

    <van-cell-group>
      <van-cell title="修改密码" is-link @click="showPasswordDialog = true" />
    </van-cell-group>

    <van-dialog
      v-model:show="showNicknameDialog"
      title="修改昵称"
      show-cancel-button
      @confirm="updateNickname"
    >
      <van-field
        v-model="editForm.nickname"
        placeholder="请输入昵称"
        :border="false"
      />
    </van-dialog>

    <van-dialog
      v-model:show="showPhoneDialog"
      title="修改手机号"
      show-cancel-button
      @confirm="updatePhone"
    >
      <van-field
        v-model="editForm.phone"
        placeholder="请输入手机号"
        :border="false"
      />
    </van-dialog>

    <van-dialog
      v-model:show="showPasswordDialog"
      title="修改密码"
      show-cancel-button
      @confirm="updatePassword"
    >
      <van-field
        v-model="passwordForm.old_password"
        type="password"
        placeholder="请输入原密码"
        :border="true"
        class="password-field"
      />
      <van-field
        v-model="passwordForm.new_password"
        type="password"
        placeholder="请输入新密码"
        :border="true"
        class="password-field"
      />
      <van-field
        v-model="passwordForm.confirm_password"
        type="password"
        placeholder="请再次输入新密码"
        :border="false"
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { showToast } from 'vant'
import { useUserStore } from '@/store/user'
import { getUserInfo, updateUserInfo, updatePassword as apiUpdatePassword } from '@/api/user'

const userStore = useUserStore()

const userInfo = ref({})
const showNicknameDialog = ref(false)
const showPhoneDialog = ref(false)
const showPasswordDialog = ref(false)

const editForm = reactive({
  nickname: '',
  phone: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const fetchUserInfo = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 200) {
      userInfo.value = res.data
      userStore.setUserInfo(res.data)
    }
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
}

const updateNickname = async () => {
  if (!editForm.nickname) {
    showToast('请输入昵称')
    return
  }
  
  try {
    const res = await updateUserInfo({ nickname: editForm.nickname })
    if (res.code === 200) {
      showToast('修改成功')
      userInfo.value.nickname = editForm.nickname
      showNicknameDialog.value = false
    } else {
      showToast(res.msg || '修改失败')
    }
  } catch (error) {
    showToast('修改失败')
  }
}

const updatePhone = async () => {
  if (!editForm.phone) {
    showToast('请输入手机号')
    return
  }
  
  try {
    const res = await updateUserInfo({ phone: editForm.phone })
    if (res.code === 200) {
      showToast('修改成功')
      userInfo.value.phone = editForm.phone
      showPhoneDialog.value = false
    } else {
      showToast(res.msg || '修改失败')
    }
  } catch (error) {
    showToast('修改失败')
  }
}

const updatePassword = async () => {
  if (!passwordForm.old_password || !passwordForm.new_password || !passwordForm.confirm_password) {
    showToast('请填写完整信息')
    return
  }
  
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    showToast('两次密码输入不一致')
    return
  }
  
  try {
    const res = await apiUpdatePassword(passwordForm)
    if (res.code === 200) {
      showToast('修改成功')
      showPasswordDialog.value = false
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
    } else {
      showToast(res.msg || '修改失败')
    }
  } catch (error) {
    showToast('修改失败')
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.settings-page {
  padding-top: 46px;
  background: #f5f5f5;
}

.password-field {
  margin-bottom: 1px;
}
</style>
