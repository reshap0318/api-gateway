<script setup lang="ts">
import SidebarMenuItem from './SidebarMenuItem.vue'
import { ref, computed } from 'vue'

export interface IMenuItem {
  icon?: unknown
  label: string
  to?: string
  children?: IMenuItem[]
  isTitle?: boolean
  permission?: string[]
}

const props = defineProps<{
  appName: string
  menuItems: IMenuItem[]
  isOpen: boolean
  isCollapsed: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const expandedGroups = ref<Set<string>>(new Set())
const isHovered = ref(false)

/**
 * Sidebar is visually expanded when:
 * - it is NOT collapsed, OR
 * - it IS collapsed but the user is hovering over it
 */
const isExpanded = computed(() => !props.isCollapsed || isHovered.value)

function toggleGroup(key: string) {
  if (expandedGroups.value.has(key)) {
    expandedGroups.value.delete(key)
  } else {
    expandedGroups.value.add(key)
  }
}

const handleItemClick = () => {
  if (window.innerWidth < 768) {
    emit('close')
  }
}

const handleMouseEnter = () => {
  if (props.isCollapsed) {
    isHovered.value = true
  }
}

const handleMouseLeave = () => {
  isHovered.value = false
}
</script>

<template>
  <!-- Overlay for mobile -->
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 md:hidden"
      @click="emit('close')"
    />
  </Teleport>

  <!-- Sidebar -->
  <aside
    :class="[
      'fixed top-0 left-0 z-50 h-full transition-all duration-300 ease-in-out',
      'bg-linear-to-b from-slate-900 via-slate-800 to-slate-900',
      'border-r border-white/5 overflow-hidden',
      // Width: expanded = w-64, collapsed = w-16
      isExpanded ? 'w-64' : 'w-16',
      // Mobile: slide in/out based on isOpen
      isOpen ? 'translate-x-0' : '-translate-x-full',
      // Desktop (md+): always visible, no translate
      'md:translate-x-0',
      // Shadow: overlay shadow when collapsed+hovered, normal shadow otherwise
      isCollapsed && isHovered ? 'shadow-[4px_0_24px_rgba(0,0,0,0.4)]' : 'shadow-2xl',
    ]"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <!-- Logo / Header -->
    <div
      class="flex items-center h-12 border-b border-white/10 transition-all duration-300"
      :class="isExpanded ? 'px-5' : 'justify-center px-2'"
    >
      <router-link to="/" class="flex items-center group overflow-hidden">
        <span
          v-if="isExpanded"
          class="text-lg font-bold bg-linear-to-r from-white to-slate-300 bg-clip-text text-transparent whitespace-nowrap"
        >
          {{ appName }}
        </span>
        <span v-else class="text-lg font-bold text-white">
          {{ appName.charAt(0) }}
        </span>
      </router-link>
    </div>

    <!-- Menu -->
    <nav
      class="sidebar-nav space-y-1 mt-2 overflow-y-auto overflow-x-hidden transition-all duration-300"
      :class="isExpanded ? 'p-3' : 'p-2'"
      style="max-height: calc(100vh - 3.5rem)"
    >
      <SidebarMenuItem
        v-for="(item, index) in menuItems"
        :key="index"
        :item="item"
        :depth="0"
        :expanded-groups="expandedGroups"
        :is-expanded="isExpanded"
        @toggle-group="toggleGroup"
        @item-click="handleItemClick"
      />
    </nav>
  </aside>
</template>

<style scoped>
.sidebar-item-active {
  box-shadow: inset 3px 0 0 #3b82f6;
}

.sidebar-nav {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}

.sidebar-nav:hover {
  scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar-nav::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: transparent;
  border-radius: 4px;
  transition: background 0.2s ease;
}

.sidebar-nav:hover::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

.sidebar-nav:hover::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.35);
}
</style>
