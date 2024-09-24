<template>
  <el-scrollbar>
    <div class="sidebar-logo">
      <div v-if="isCollapse" class="logo-mini"><b>A</b>-</div>
      <el-image v-else style="height:45px;margin-top:15px" src="/images/lizardrestic-logo.png" />
    </div>
    <div class="user-panel">
      <div class="pull-left image">
        <el-avatar :size="45" :src="avator" />
      </div>
      <div class="pull-left info sidebar-userinfo">
        <p style="text-align:left">{{username}}</p>
        <a><font-awesome-icon icon="circle" style="color:green" /> Online</a>
      </div>
    </div>
    <div class="user-panel sidenav" style="text-align:left;padding:8px 10px;margin-top:5px;">
      导航
    </div>
    <el-menu
        active-text-color="#ffd04b"
        background-color="#141f29"
        class="el-menu-vertical-demo"
        default-active="2"
        text-color="#fff"
        :unique-opened="true"
        :router="true">
        <el-menu-item index="/"><font-awesome-icon icon="home" />首页</el-menu-item>
        <el-menu-item index="/agent"><font-awesome-icon icon="sliders" />Agent列表</el-menu-item>
        <el-menu-item index="/repository"><font-awesome-icon icon="laptop-code" />仓库管理</el-menu-item>
        <el-menu-item index="/backup/policy"><font-awesome-icon icon="layer-group" />策略管理</el-menu-item>
        <el-menu-item index="/backup/history"><font-awesome-icon icon="list-check" />执行记录</el-menu-item>
    </el-menu>
  </el-scrollbar>
</template>

<script setup>
import { onBeforeMount, reactive, toRefs, ref, computed } from 'vue'
import { useStore } from 'vuex'
/* 变量定义 */
const store = useStore()
const avator = computed(() => {
  return store.state.avatar || "/images/avator.png"
})
const username = computed(() => {
  return store.state.username
})
const isCollapse = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  isCollapse.value = false
})
</script>