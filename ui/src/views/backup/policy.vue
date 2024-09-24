<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>策略管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">策略管理</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable placeholder="输入策略名查询……" :prefix-icon="Search" @change="getList(1);current=1" style="width:25%;margin-right:5px" size="large" />
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-button-group class="pull-right">
        <el-button class="pull-right" size="large" type="primary" @click="show.add=true;form={tags:[],enable:true}" :icon="Plus">新建策略</el-button>
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
    <el-table-column prop="name" label="策略名称" min-width="160">
      <template #default="scope">
        <el-link :underline="false" type="primary" :href="`/backup/policy/${scope.row.id}/snapshot`">{{ scope.row.name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column prop="repo_name" label="仓库名称" min-width="160">
      <template #default="scope">
        <el-link :underline="false">{{ repoList[scope.row.repository_id]?.repo_name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column prop="cron" label="调度时间" min-width="160" />
    <el-table-column label="标签" min-width="180">
      <template #default="scope">
        <el-tag v-for="(item,i) in scope.row.tags" :key="i">{{ item }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="备份目录" min-width="120">
      <template #default="scope">
        <div v-html="scope.row.backup_dir" />
      </template>
    </el-table-column>
    <el-table-column label="备份主机" min-width="120">
      <template #default="scope">
        <div v-for="(item,i) in scope.row.hosts" :key="i">{{ item }}</div>
      </template>
    </el-table-column>
    <el-table-column label="是否启用" min-width="120">
      <template #default="scope">
        <el-icon v-if="scope.row.enable" class="text-green font-size18"><CircleCheckFilled /></el-icon>
        <el-icon v-else class="text-red font-size18"><CircleCloseFilled /></el-icon>
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="210">
      <template #default="scope">
        <el-button :icon="EditPen" circle @click="editOne(scope.row)" />
        <el-tooltip effect="dark" content="复制" placement="top">
          <el-button :icon="CopyDocument" circle @click="copyOne(scope.row)" />
        </el-tooltip>
        <el-popconfirm title="确认立即执行？" @confirm="doBackup(scope.row)">
          <template #reference>
            <el-button :icon="ArrowRight" circle />
          </template>
        </el-popconfirm>
        <el-tooltip effect="dark" content="执行记录" placement="top">
          <el-button :icon="Expand" circle @click="gotoHistory(scope.row)" />
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
<el-drawer v-model="show.add" direction="rtl" size="600px">
  <template #header>
    <h4 v-if="edit===false">新建策略</h4>
    <h4 v-if="edit===true">编辑策略</h4>
  </template>
  <template #default>
    <el-form ref="policy" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="策略名称" prop="name">
        <el-input v-model="form.name" size="large" clearable />
      </el-form-item>
      <el-form-item label="选择仓库" prop="repository">
        <el-select 
          v-model="form.repository" 
          placeholder="请选择仓库" 
          value-key="id" 
          clearable 
          filterable
          size="large"
          style="width:100%">
          <el-option v-for="item in repoList" :key="item.id" :label="item.repo_name" :value="item">
            <span style="float:left">{{item.repo_name}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.repo_url}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item prop="cron">
        <template #label><el-text>调度时间 
          <el-tooltip placement="top">
            <template #content>
              cron 格式：* * * * * *，分别表示秒、分、时、日、月、周
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input v-model="form.cron" size="large" clearable />
      </el-form-item>
      <el-form-item label="设置标签" prop="tags">
        <el-tag v-for="item in form.tags" :key="item" closable :disable-transitions="false" @close="handleClose(item)" size="large">{{item}}</el-tag>
        <el-input 
          v-if="inputVisible" 
          ref="inputRef" 
          v-model="inputLabel" 
          @keyup.enter="handleInputConfirm" 
          @blur="handleInputConfirm" 
          style="width:100px" />
        <el-button v-else @click="showInput">+ 添加标签</el-button>
      </el-form-item>
      <el-form-item label="备份目录/文件" prop="backup_dir">
        <el-input v-model="form.backup_dir" size="large" type="textarea" :autosize="{minRows:3}" />
      </el-form-item>
      <el-form-item label="排除路径/模式" prop="exclude">
        <el-input v-model="form.exclude" size="large" type="textarea" :autosize="{minRows:2}" />
      </el-form-item>
      <el-form-item label="保留时间" prop="retention">
        <el-input v-model="form.retention" size="large" clearable placeholder="3d（y/m/d/h）" />
      </el-form-item>
      <el-form-item label="备份主机" prop="hosts">
        <el-select 
          v-model="form.hosts" 
          placeholder="请选择主机" 
          value-key="ip" 
          clearable 
          filterable
          multiple
          size="large"
          style="width:100%">
          <el-option v-for="item in targets" :key="item.ip" :label="item.ip" :value="item">
            <span style="float:left">{{item.ip}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.system}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="是否启用" prop="enable">
        <el-switch v-model="form.enable" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(policy)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { ArrowRight,Warning,Search,Refresh,CopyDocument,Plus,Delete,EditPen,Expand } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import router from '@/router';
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
const show = ref({
  add: false
})
const edit = ref(false)
const rules = reactive({
  name: [{required: true, message: '请填写策略名称'}],
  repository: [{required: true, trigger:'blur', message: '请选择仓库'}],
  cron: [{required: true, message: '请填写调度时间'}],
  backup_dir: [{required: true, message: '请填写备份目录/文件'}],
  hosts: [{required: true, trigger:'blur', message: '请选择主机'}],
})
const repoList = ref({})
const policy = ref(null)
const targets = ref([])
const inputRef = ref(null)
const inputLabel = ref('')
const inputVisible = ref(false)
/* 生命周期函数 */
onBeforeMount(async () => {
  getRepos()
  getList(1)
  getTargets()
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value !== "") {
    url += `&search=name==${encodeURIComponent(searchKey.value)}`
  } 
  loading.value.table = true
  let response = await axios.get(`/restic/db/backup_policy?${url}`)
  loading.value.table = false
  list.value = response.results||[]
  pageTotal.value = response.total
}
const getRepos = async () => {
  let response = await axios.get(`/restic/db/repository`)
  for(let x of response.results) {
    repoList.value[x.id] = x
  }
}
const getTargets = async () => {
  targets.value = await axios.get(`/restic/targets`)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  form.value.repository = repoList.value[form.value.repository_id]
  form.value.hosts = targets.value.filter(n => row.hosts.includes(n.ip))
  edit.value = true
  show.value.add = true
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  form.value.repository = repoList.value[form.value.repository_id]
  form.value.hosts = targets.value.filter(n => row.hosts.includes(n.ip))
  edit.value = false
  show.value.add = true
}
const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value.input.focus()
  })
}
const handleClose = (tag) => {
  form.value.tags.splice(form.value.tags.indexOf(tag), 1)
}
const handleInputConfirm = () => {
  if(inputLabel.value) {
    form.value.tags.push(inputLabel.value)
  }
  inputVisible.value = false
  inputLabel.value = ''
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.repository_id = params.repository.id
      delete params.repository
      params.hosts = params.hosts.map(n => n.ip)
      console.log(params)
      await axios.post(`/restic/backup/policy`, params)
      getList(1)
      current.value = 1
      show.value.add = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/restic/backup/policy?policy_id=${row.id}`)
  getList(current.value)
}
const gotoHistory = (row) => {
  router.push({
    path: '/backup/history',
    query: {
      policy_id: row.id
    }
  })
}
const doBackup = async (row) => {
  await axios.post(`/restic/backup?policy_id=${row.id}`)
}
</script>