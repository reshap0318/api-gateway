<script setup lang="ts">
import SidebarMenu from '@/components/layouts/SidebarMenu.vue'
import TopBar from '@/components/layouts/TopBar.vue'
import { ref, computed } from 'vue'
import {
  PhList,
  PhHouse,
  PhShieldCheck,
  PhUsers,
  PhGlobe,
  PhClockCounterClockwise,
} from '@phosphor-icons/vue'
import type { IMenuItem } from '@/components/layouts/SidebarMenu.vue'
import { usePermission } from '@/composables'

const sidebarOpen = ref(false)
const sidebarCollapsed = ref(false)
const appName = import.meta.env.VITE_APP_NAME || 'Admin'

const { hasAnyPermission } = usePermission()

const menuItems = computed<IMenuItem[]>(() => {
  const items: IMenuItem[] = [
    { icon: PhHouse, label: 'Dashboard', to: '/' },
    { isTitle: true, label: 'Management' },
    { icon: PhUsers, label: 'Users', to: '/users', permission: ['user.index'] },
    {
      icon: PhShieldCheck,
      label: 'UAM',
      children: [
        { label: 'Roles', to: '/uam/roles', permission: ['role.index'] },
        { label: 'Permissions', to: '/uam/permissions', permission: ['permission.index'] },
      ],
    },
    {
      icon: PhGlobe,
      label: 'API Gateway',
      children: [
        { label: 'Services', to: '/gateway/services', permission: ['service.index'] },
        { label: 'Routes', to: '/gateway/routes', permission: ['route.index'] },
      ],
    },
    {
      icon: PhClockCounterClockwise,
      label: 'Audit Log',
      to: '/audit-logs',
      permission: ['audit.index'],
    },
  ]

  return filterMenuByPermission(items)
})

function filterMenuByPermission(items: IMenuItem[]): IMenuItem[] {
  const filtered = items
    .filter((item) => {
      if (item.isTitle) return true
      if (!item.permission || item.permission.length === 0) return true
      return hasAnyPermission(item.permission)
    })
    .map((item) => {
      if (item.children) {
        return { ...item, children: filterMenuByPermission(item.children) }
      }
      return item
    })
    .filter((item) => !(item.children && item.children.length === 0))

  // Remove title items that have no visible non-title items following them
  return filtered.filter((item, index) => {
    if (!item.isTitle) return true
    const next = filtered[index + 1]
    return next !== undefined && !next.isTitle
  })
}

const toggleSidebar = () => {
  if (window.innerWidth >= 768) {
    sidebarCollapsed.value = !sidebarCollapsed.value
  } else {
    sidebarOpen.value = !sidebarOpen.value
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Sidebar -->
    <SidebarMenu
      :app-name="appName"
      :menu-items="menuItems"
      :is-open="sidebarOpen"
      :is-collapsed="sidebarCollapsed"
      @close="sidebarOpen = false"
    />

    <!-- Main Content -->
    <div
      class="transition-all duration-300 ease-in-out"
      :class="sidebarCollapsed ? 'md:ml-16' : 'md:ml-64'"
    >
      <!-- Top Bar -->
      <TopBar :show-hamburger="true" @toggle-sidebar="toggleSidebar">
        <template #menu-icon>
          <PhList class="w-6 h-6" />
        </template>
      </TopBar>

      <!-- Page Content -->
      <main class="p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>
