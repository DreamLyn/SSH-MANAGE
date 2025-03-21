import { createRouter, createWebHashHistory } from 'vue-router'
/* Layout */
import Layout from '@/layout'

/**
 * Note: 路由配置项
  }
 */

// 公共路由
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login'),
    hidden: true
  },
  {
    path: "/:pathMatch(.*)*",
    component: () => import('@/views/error/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error/401'),
    hidden: true
  },
  {
    path: '',
    component: Layout,
    redirect: '/vaults/hosts',
    children: [{
      path: 'vaults/hosts',
      name: 'Hosts',
      component: () => import('@/views/hosts/index'),
      meta: { title: 'Hosts', showMenu: true }
    }, {
      path: 'vaults/forwarding',
      name: 'Forwarding',
      component: () => import('@/views/forwarding/index'),
      meta: { title: 'Forwarding', showMenu: true }
    }, {
      path: 'terminal/:id',
      name: 'Terminal',
      component: () => import('@/views/terminal/index'),
      meta: { title: 'Terminal', showMenu: false },
    }, {
      path: 'setting',
      name: 'Setting',
      component: () => import('@/views/setting/index'),
      meta: { title: 'Setting', showMenu: true },
    }]
  },

]

// 动态路由，基于用户权限动态去加载
export const dynamicRoutes = [

]

const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoutes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
});

export default router;
