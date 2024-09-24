<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>执行记录</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">执行记录</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable :prefix-icon="Search" placeholder="输入策略名查询……" @change="getList(1);current=1" style="width:30%;margin-right:5px" size="large" />
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-date-picker
          style="float:right"
          v-model="timerange"
          type="datetimerange"
          :shortcuts="shortcuts"
          range-separator="To"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          size="large"
          @change="getList(1);current=1" />
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    @expand-change="getTaskHistoryWorkload"
    @filter-change="filterTable"
    @sort-change="sortTable"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px">
    <!-- <el-table-column type="selection" width="40" /> -->
    <el-table-column prop="policy_name" label="策略名称" min-width="150">
      <template #default="scope">
      </template>
    </el-table-column>
    <el-table-column prop="policy_name" label="执行结果" width="150">
      <template #default="scope">
        <el-popover placement="right" trigger="click" :width="500">
          <template #reference>
            <el-link :underline="false" type="info">点我查看</el-link>
          </template>
          <pre>{{ scope.row.message }}</pre>
        </el-popover>
      </template>
    </el-table-column>
    <el-table-column prop="task_type" label="任务类型" :filters="taskTypeFilters" column-key="task_type" :filter-multiple="false" width="120" />
    <el-table-column label="成功/失败" :filters="successFilters" column-key="success" :filter-multiple="false" width="120">
      <template #default="scope">
        <span v-if="scope.row.success===true" style="color:#5cb87a">
          <el-icon><Check /></el-icon>
        </span>
        <span v-else-if="scope.row.success===false" style="color:#f56c6c">
          <el-icon><Close /></el-icon>
        </span>
        <span v-else style="color:#e6a23c">
          <font-awesome-icon icon="circle" class="twinkling" style="font-size:12px " />
        </span>
      </template>
    </el-table-column>
    <el-table-column label="状态" :filters="statusFilters" column-key="status" :filter-multiple="false" width="150">
      <template #default="scope">
        <el-progress v-if="scope.row.status=='running'" :percentage="50" color="#e6a23c" :show-text="false" />
        <el-progress v-else-if="scope.row.status=='finished'&&scope.row.success===true" :percentage="100" color="#5cb87a" :show-text="false" />
        <el-progress v-else :percentage="100" color="#f56c6c" :show-text="false" />
      </template>
    </el-table-column>
    <el-table-column prop="start_at" label="开始时间" sortable="custom" width="150">
      <template #default="scope">
        {{ scope.row.start_at ? moment(scope.row.start_at).format('YYYY-MM-DD HH:mm') : '' }}
      </template>
    </el-table-column>
    <el-table-column prop="finish_at" label="结束时间" sortable="custom" width="150">
      <template #default="scope">
        {{ scope.row.finish_at ? moment(scope.row.finish_at.Time).format('YYYY-MM-DD HH:mm') : '' }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="70">
      <template #default="scope">
        <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button :icon="Close" circle />
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
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @current-change="getList"
    @size-change="handleSizeChange"
    v-model:current-page="current" />
</el-card>
</template>
<script setup>
import { ArrowRight,Refresh,Close,Search } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const route = useRoute()
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const loading = ref({
  table: false
})
const timerange = ref([])
const shortcuts = [
  {
    text: '最近6小时',
    value: () => {
      return [moment().subtract(6,'hours'), moment()]
    }
  },
  {
    text: '最近1天',
    value: () => {
      return [moment().subtract(1,'days'), moment()]
    }
  },
  {
    text: '最近3天',
    value: () => {
      return [moment().subtract(3,'days'), moment()]
    }
  },
  {
    text: '最近1周',
    value: () => {
      return [moment().subtract(1,'weeks'), moment()]
    }
  },
]
const taskHistoryWorkload = ref({})
const statusFilters = ref([
  { text: 'finished', value: 'finished'},
  { text: 'running', value: 'running'},
])
const taskTypeFilters = ref([
  { text: 'backup', value: 'backup'},
  { text: 'restore', value: 'restore'},
])
const successFilters = ref([
  { text: '成功', value: '1'},
  { text: '失败', value: '0'},
])
const filterOptions = ref({})
const sort = ref({prop: 'start_at', order: 'descending'})
/* 生命周期函数 */
onBeforeMount(async () => {
  if(route.query.policy_id) {
    filterOptions.value["backup_policy_id"] = [route.query.policy_id]
  }
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=${sort.value.prop} ${sort.value.order==='descending'?'desc':'asc'}`
  if(timerange.value?.length == 2)
    url += `&range=${sort.value.prop}==${moment(timerange.value[0]).format('YYYY-MM-DD HH:mm:ss')},${moment(timerange.value[1]).format('YYYY-MM-DD HH:mm:ss')}`
  if(searchKey.value !== "") {
    url += `&policy_name=${encodeURIComponent(searchKey.value)}`
  }
  let filterParams = []
  if(Object.keys(filterOptions.value).length > 0) {
    for(let [k,v] of Object.entries(filterOptions.value)) {
      filterParams.push(`${k}==${v[0]}`)
    }
  }
  if(filterParams.length > 0) url += `&filter=${filterParams.join(",")}`
  loading.value.table = true
  let response = await axios.get(`/restic/backup/history?${url}`)
  loading.value.table = false
  list.value = response.results?.map(x => {
   return x
  })
  pageTotal.value = response.total
}
const execute = async (row) => {
  ElMessageBox.confirm(
    '确认执行此任务？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    await axios.post(`/restic/task/execute/${row.id}`)
    getList(current.value)
  }).catch((e) => {
    console.warn(e)
  })
}
const getTaskHistoryWorkload = async (row) => {
  let response = await axios.get(`/restic/db/task_history/${row.id}`)
  taskHistoryWorkload.value[row.id] = response.workloads.map(x => {
    try {
      x.status = JSON.parse(x.status)
    }
    catch {
      x.status = [x.status]
    }
    return x
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/restic/db/backup_history/${row.id}`)
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
const filterTable = async (filter) => {
  filterOptions.value = Object.assign(filterOptions.value, filter)
  for(let [k,v] of Object.entries(filterOptions.value)) {
    if(v.length === 0)
      delete filterOptions.value[k]
  }
  getList(current.value)
}
const sortTable = async (data) => {
  sort.value = data
  getList(current.value)
}
</script>