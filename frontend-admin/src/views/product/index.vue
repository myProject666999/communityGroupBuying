<template>
  <div class="product-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>农产品列表</span>
          <el-button type="primary" @click="handleAdd">新增商品</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="商品名称">
          <el-input v-model="searchForm.name" placeholder="请输入商品名称" clearable />
        </el-form-item>
        <el-form-item label="商品分类">
          <el-select v-model="searchForm.category_id" placeholder="请选择分类" clearable>
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="商品名称" min-width="200" />
        <el-table-column prop="category_name" label="分类" width="120" />
        <el-table-column label="图片" width="100">
          <template #default="{ row }">
            <el-image
              :src="row.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=default%20product&image_size=square'"
              style="width: 60px; height: 60px"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column prop="price" label="售价" width="100">
          <template #default="{ row }">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100" />
        <el-table-column prop="sales" label="销量" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="toggleStatus(row)">
              {{ row.status === 1 ? '下架' : '上架' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="售价" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="市场价" prop="market_price">
          <el-input-number v-model="form.market_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入商品描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新增商品')
const formRef = ref(null)

const searchForm = reactive({
  name: '',
  category_id: null
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const tableData = ref([
  { id: 1, name: '新鲜有机青菜 500g', category_name: '蔬菜', price: 9.9, stock: 100, sales: 256, status: 1 },
  { id: 2, name: '红富士苹果 5斤装', category_name: '水果', price: 29.9, stock: 80, sales: 189, status: 1 },
  { id: 3, name: '农家土鸡蛋 30枚', category_name: '蛋类', price: 39.9, stock: 50, sales: 423, status: 1 }
])

const categories = ref([
  { id: 1, name: '蔬菜' },
  { id: 2, name: '水果' },
  { id: 3, name: '肉类' },
  { id: 4, name: '海鲜' },
  { id: 5, name: '粮油' },
  { id: 6, name: '蛋类' }
])

const form = reactive({
  id: null,
  name: '',
  category_id: null,
  price: 0,
  market_price: 0,
  stock: 0,
  status: 1,
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择商品分类', trigger: 'change' }],
  price: [{ required: true, message: '请输入售价', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/product/list',
      method: 'get',
      params: {
        ...searchForm,
        page: pagination.page,
        page_size: pagination.page_size
      }
    })
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取商品列表失败', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const res = await request({
      url: '/admin/category/list',
      method: 'get'
    })
    if (res.code === 200 && res.data.length > 0) {
      categories.value = res.data
    }
  } catch (error) {
    console.error('获取分类列表失败', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.category_id = null
  pagination.page = 1
  fetchData()
}

const handleAdd = () => {
  dialogTitle.value = '新增商品'
  form.id = null
  form.name = ''
  form.category_id = null
  form.price = 0
  form.market_price = 0
  form.stock = 0
  form.status = 1
  form.description = ''
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑商品'
  form.id = row.id
  form.name = row.name
  form.category_id = row.category_id
  form.price = row.price
  form.market_price = row.market_price
  form.stock = row.stock
  form.status = row.status
  form.description = row.description
  dialogVisible.value = true
}

const toggleStatus = (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const text = newStatus === 1 ? '上架' : '下架'
  ElMessageBox.confirm(`确定要${text}该商品吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await request({
        url: `/admin/product/status/${row.id}`,
        method: 'put',
        data: { status: newStatus }
      })
      if (res.code === 200) {
        ElMessage.success(`${text}成功`)
        fetchData()
      }
    } catch (error) {
      ElMessage.error(`${text}失败`)
    }
  }).catch(() => {})
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await request({
        url: `/admin/product/delete/${row.id}`,
        method: 'delete'
      })
      if (res.code === 200) {
        ElMessage.success('删除成功')
        fetchData()
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const url = form.id ? `/admin/product/update/${form.id}` : '/admin/product/add'
    const method = form.id ? 'put' : 'post'
    const res = await request({ url, method, data: form })
    if (res.code === 200) {
      ElMessage.success(form.id ? '编辑成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchData()
  fetchCategories()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
