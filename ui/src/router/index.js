import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      meta: {
        description: "首页",
      },
      component: () => import("../views/home.vue")
    },
    {
      path: '/agent',
      name: 'agent',
      meta: {
        description: "Agent列表",
      },
      component: () => import("../views/agent.vue")
    },
    {
      path: '/repository',
      name: 'repository',
      meta: {
        description: "仓库管理",
      },
      component: () => import("../views/repository.vue")
    },
    {
      path: '/backup/policy',
      name: 'policy',
      meta: {
        description: "策略管理",
      },
      component: () => import("../views/backup/policy.vue")
    },
    {
      path: '/backup/policy/:id/snapshot',
      name: 'snapshot',
      meta: {
        description: "备份快照",
      },
      component: () => import("../views/backup/snapshot.vue")
    },
    {
      path: '/backup/history',
      name: '',
      meta: {
        description: "执行记录",
      },
      component: () => import("../views/backup/history.vue")
    },
  ]
})

export default router
