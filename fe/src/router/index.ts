import type { RouteMeta, RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

interface CustomRouteMeta extends RouteMeta {
  guest?: boolean
  requiresAuth?: boolean
  permissions?: string[]
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/LoginView.vue'),
    meta: { guest: true },
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/pages/auth/ForgotPasswordView.vue'),
    meta: { guest: true },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/pages/auth/ResetPasswordView.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: DefaultLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/pages/HomeView.vue'),
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/pages/users/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['user.index'] },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/pages/profile/IndexView.vue'),
      },
      {
        path: 'uam/permissions',
        name: 'Permissions',
        component: () => import('@/pages/uam/permissions/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['permission.index'] },
      },
      {
        path: 'uam/roles',
        name: 'Roles',
        component: () => import('@/pages/uam/roles/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['role.index'] },
      },
      {
        path: 'gateway/services',
        name: 'GatewayServices',
        component: () => import('@/pages/gateway/services/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['service.index'] },
      },
      {
        path: 'gateway/routes',
        name: 'GatewayRoutes',
        component: () => import('@/pages/gateway/routes/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['route.index'] },
      },
      {
        path: 'audit-logs',
        name: 'AuditLogs',
        component: () => import('@/pages/gateway/audit-logs/IndexView.vue'),
        meta: { requiresAuth: true, permissions: ['audit.index'] },
      },
    ],
  },
  // Catch-all route for 404 - must be last
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/errors/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const authStore = useAuthStore()
  const token = authStore.token

  if (to.meta.requiresAuth && !token) {
    return { name: 'Login' }
  }

  if (to.meta.guest && token) {
    return { name: 'Home' }
  }

  const meta = to.meta as CustomRouteMeta
  if (meta.permissions && token) {
    const userPermissions = authStore.user?.permissions?.map((p) => p.name) || []
    const hasAccess = meta.permissions.some((perm) => userPermissions.includes(perm))

    if (!hasAccess) {
      return { name: 'Home', query: { accessDenied: 'true' } }
    }
  }

  return true
})

export default router
