<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/backup/policy' }">策略管理</el-breadcrumb-item>
  <el-breadcrumb-item>{{ policy.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">备份快照</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-select 
          v-model="host" 
          placeholder="请选择主机" 
          size="large"
          @change="getList"
          style="width:20%;margin-right:5px">
          <el-option v-for="item in policy.hosts" :key="item" :label="item" :value="item">
            <span style="float:left">{{item}}</span>
          </el-option>
        </el-select>
        <el-input v-model="searchTag" clearable placeholder="过滤标签……" :prefix-icon="Search" @change="getList();current=1" style="width:25%;margin-right:5px" size="large" />
        <el-checkbox label="最近10条" :value="10" style="vertical-align: middle;" v-model="latest" @change="getList()" />
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-button-group class="pull-right">
        <el-popconfirm title="确认删除？" @confirm="deleteBatch">
          <template #reference>
            <el-button class="pull-right" size="large" type="danger" :icon="DeleteFilled" :loading="loading.delete">批量删除</el-button>
          </template>
        </el-popconfirm>
      </el-button-group>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    @selection-change="select"
    style="width:100%;margin-top:10px">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="short_id" label="快照ID" width="120" />
    <el-table-column prop="repo_name" label="快照时间" width="200">
      <template #default="scope">
        {{ moment(scope.row.time).format('YYYY-MM-DD HH:mm:ss.SSS') }}
      </template>
    </el-table-column>
    <el-table-column prop="hostname" label="主机" min-width="120" />
    <el-table-column label="标签" min-width="180">
      <template #default="scope">
        <el-tag v-for="(item,i) in scope.row.tags" :key="i">{{ item }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="备份目录" min-width="120">
      <template #default="scope">
        {{ scope.row.paths.join(", ") }}
      </template>
    </el-table-column>
    <el-table-column prop="size" label="快照大小" min-width="120" />
    <el-table-column prop="Option" label="操作" width="170">
      <template #default="scope">
        <el-tooltip content="恢复快照" placement="top">
          <el-button :icon="RefreshLeft" circle @click="snapshot=scope.row;show.restore=true" />
        </el-tooltip>
        <el-tooltip content="列出快照文件" placement="top">
          <el-button circle :icon="Menu" @click="snapshot=scope.row;show.ls=true" />
        </el-tooltip>
        <el-tooltip content="搜索快照文件" placement="top">
          <el-button circle :icon="Search" @click="snapshot=scope.row;show.find=true" />
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
    layout="total" 
    :total="pageTotal" />
</el-card>
<el-dialog v-model="show.ls" title="列出快照文件" width="800px">
  <el-row>
    <el-col :span="24">
      <el-input v-model="dir" style="width:40%;margin-bottom:10px" size="large" placeholder="输入目录，按回车提交" @change="lsSnapshot" />
    </el-col>
  </el-row>
  <pre v-loading="loading.ls" style="min-height:70px">{{ lsContent }}</pre>
</el-dialog>
<el-dialog v-model="show.find" title="搜索快照文件" width="800px">
  <el-row>
    <el-col :span="24">
      <el-input v-model="pattern" style="width:40%;margin-bottom:10px" size="large" placeholder="输入关键词，可使用通配符，按回车提交" @change="findSnapshot" />
    </el-col>
  </el-row>
  <pre v-loading="loading.find" style="min-height:70px">{{ findContent }}</pre>
</el-dialog>
<el-drawer v-model="show.restore" direction="rtl" size="600px">
  <template #header>
    <h4>恢复快照</h4>
  </template>
  <template #default>
    <el-form ref="re" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="恢复主机" prop="host">
        <el-select 
          v-model="form.host" 
          placeholder="请选择主机" 
          clearable 
          filterable
          size="large"
          style="width:100%">
          <el-option v-for="item in policy.hosts" :key="item" :label="item" :value="item">
            <span style="float:left">{{item}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="恢复目录" prop="target">
        <el-input v-model="form.target" size="large" />
      </el-form-item>
      <el-form-item label="排除文件/目录" prop="target">
        <el-input v-model="form.exclude" size="large" placeholder="支持通配符模式" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="restore(re)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { ArrowRight,RefreshLeft,Search,Refresh,Menu,DeleteFilled,Delete,EditPen } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 变量定义 */
const route = useRoute()
const list = ref([])
const pageTotal = ref(0)
const loading = ref({
  table: false,
  delete: false,
  ls: false,
  find: false
})
const searchTag = ref("")
const latest = ref(true)
const policy = ref({})
const host = ref("")
const selected = ref([])
const show = ref({
  ls: false,
  find: false,
  restore: false
})
const snapshot = ref({})
const dir = ref("")
const pattern = ref("")
const lsContent = ref("")
const findContent = ref("")
const form = ref({})
const rules = reactive({
  host: [{required: true, message: '请选择恢复主机'}],
  target: [{required: true, trigger:'blur', message: '请填写恢复目录'}],
})
const re = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getPolicy()
});
/* methods */
const getList = async () => {
  let url = `host=${host.value}`
  if(searchTag.value !== "") {
    url += `&tag=${searchTag.value}`
  } 
  if(latest.value === true) {
    url += `&latest=10`
  }
  loading.value.table = true
  let response = await axios.get(`/restic/backup/policy/${policy.value.id}/snapshots?${url}`)
  list.value = response.reverse()
  loading.value.table = false
  pageTotal.value = list.value.length
}
const getPolicy = async () => {
  policy.value = await axios.get(`/restic/db/backup_policy/${route.params.id}`)
}
const select = (val) => {
  selected.value = val
}
const deleteBatch = async () => {
  loading.value.delete = true
  await axios.delete(`/restic/snapshots`, {
    data: {
      repo_url: policy.value.repository.repo_url,
      ids: selected.value.map(n => n.short_id)
    }
  })
  ElMessage.success({message: '批量删除成功'})
  loading.value.delete = false
  getList()
}
const lsSnapshot = async () => {
  let url = `/restic/snapshots/ls?repo_url=${policy.value.repository.repo_url}&snapshot_id=${snapshot.value.short_id}&host=${host.value}&dir=${dir.value}`
  if(searchTag.value !== "") {
    url += `&tag=${searchTag.value}`
  }
  loading.value.ls = true
  lsContent.value = await axios.get(url)
  loading.value.ls = false
}
const findSnapshot = async () => {
  let url = `/restic/snapshots/find?repo_url=${policy.value.repository.repo_url}&snapshot_id=${snapshot.value.short_id}&host=${host.value}&pattern=${pattern.value}`
  if(searchTag.value !== "") {
    url += `&tag=${searchTag.value}`
  }
  loading.value.find = true
  findContent.value = await axios.get(url)
  loading.value.find = false
}
const restore = async (row) => {
  await axios.post(`/restic/backup/restore`, {
    "policy_id": policy.value.id,
    "host": form.value.host,
    "target": form.value.target,
    "snapshot_id": snapshot.value.short_id,
    "exclude": form.value.exclude,
  })
  show.value.restore = false
}
</script>