<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
</el-breadcrumb>
</template>
<script setup>
import { ref, inject, onMounted, onBeforeUnmount } from "vue";
import { ArrowRight,Search,RefreshLeft,EditPen,Delete,Plus } from '@element-plus/icons-vue'
import moment from "moment"
import { axios } from '/src/assets/util/axios.js'
let echarts = inject("ec")
const date1 = ref(new Date())
const chart1 = ref({})
const chart2 = ref({})
onMounted(async () => {
})
onBeforeUnmount(() => {
 
})
/* methods */
</script>