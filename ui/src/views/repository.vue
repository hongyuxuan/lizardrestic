<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>仓库管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">仓库管理</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable placeholder="输入仓库名查询……" :prefix-icon="Search" @change="getList(1);current=1" style="width:25%;margin-right:5px" size="large" />
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-button-group class="pull-right">
        <el-button class="pull-right" size="large" type="primary" @click="show=true;form={}" :icon="Plus">新建仓库</el-button>
      </el-button-group>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="repo_name" label="仓库名称" min-width="160" />
    <el-table-column prop="repo_url" label="仓库地址" min-width="200" />
    <el-table-column prop="files_count" label="文件数量" min-width="120" />
    <el-table-column prop="total_size" label="仓库数据量" min-width="120" />
    <el-table-column prop="snapshot_count" label="快照数量" min-width="120" />
    <el-table-column prop="Option" label="操作" width="100">
      <template #default="scope">
        <el-tooltip effect="dark" content="复制" placement="top">
          <el-button :icon="CopyDocument" circle @click="copyOne(scope.row)" />
        </el-tooltip>
        <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button :icon="Delete" circle />
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
  <el-pagination 
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[10, 30, 50, 100]"
    layout="total, sizes, prev, pager, fnext, jumper" 
    :total="pageTotal"
    @size-change="handleSizeChange"
    @current-change="getList"
    v-model:current-page="current" />
</el-card>
<el-drawer v-model="show" direction="rtl" size="600px">
  <template #header>
    <h4>新建仓库</h4>
  </template>
  <template #default>
    <el-form ref="repo" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="仓库名称" prop="repo_name">
        <el-input v-model="form.repo_name" size="large" />
      </el-form-item>
      <el-form-item prop="repo_url">
        <template #label><el-text>仓库地址 
          <el-tooltip placement="top">
            <template #content>
              仓库地址仅支持兼容 S3 协议的对象存储，例如 Minio 或 AWS S3<br>对象存储需提前创建好桶
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input v-model="form.repo_url" size="large" placeholder="s3:" />
      </el-form-item>
      <el-form-item label="S3AccessKey" prop="s3_access_key">
        <el-input v-model="form.s3_access_key" size="large" />
      </el-form-item>
      <el-form-item label="S3SecretKey" prop="s3_secret_key">
        <el-input v-model="form.s3_secret_key" size="large" />
      </el-form-item>
      <el-form-item label="仓库密码" prop="password">
        <el-input v-model="form.password" size="large" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(repo)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { ArrowRight,Warning,Search,Refresh,CopyDocument,Plus,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const loading = ref({
  table: false,
})
const searchKey = ref("")
const form = ref({})
const show = ref(false)
const rules = reactive({
  repo_name: [{required: true, message: '请填写仓库名称'}],
  repo_url: [{required: true, message: '请填写仓库地址'}],
  s3_access_key: [{required: true, message: '请填写S3AccessKey'}],
  s3_secret_key: [{required: true, message: '请填写S3SecretKey'}],
  password: [{required: true, message: '请填写仓库密码'}],
})
const repo = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value !== "") {
    url += `&search=repo_name==${encodeURIComponent(searchKey.value)}`
  } 
  loading.value.table = true
  let response = await axios.get(`/restic/db/repository?${url}`)
  loading.value.table = false
  list.value = response.results||[]
  pageTotal.value = response.total
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/restic/db/repository/${row.id}`)
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(!params.repo_url.startsWith("s3:")) {
        ElMessage.warning('仓库地址必须以s3:开头')
        return
      }
      await axios.post(`/restic/repository`, params)
      getList(1)
      current.value = 1
      show.value = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
</script>