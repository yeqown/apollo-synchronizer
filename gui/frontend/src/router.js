import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'dashboard',
        component: () => import('./views/Dashboard.vue')
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        component: () => import('./views/Dashboard.vue')
    },
    {
        path: '/about',
        name: 'about',
        component: () => import('./views/About.vue')
    },
    {
        path: '/synchronize',
        name: 'synchronize',
        component: () => import('./views/Synchronize.vue')
    },
    {
        path: '/setting',
        name: 'setting',
        component: () => import('./views/Setting.vue')
    },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router
