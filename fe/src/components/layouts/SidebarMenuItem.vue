<script setup lang="ts">
import { useRoute } from 'vue-router'
import { PhCaretDown } from '@phosphor-icons/vue'
import type { IMenuItem } from './SidebarMenu.vue'

const props = defineProps<{
  item: IMenuItem
  depth: number
  expandedGroups: Set<string>
  isExpanded: boolean
}>()

const emit = defineEmits<{
  toggleGroup: [key: string]
  itemClick: []
}>()

const route = useRoute()

/**
 * Build a unique key for expand/collapse tracking based on depth + label.
 * This avoids collisions when the same label appears at different depths.
 */
function groupKey(item: IMenuItem, depth: number): string {
  return `${depth}::${item.label}`
}

function isActive(to: string): boolean {
  if (to === '/') return route.path === '/'
  return route.path === to || route.path.startsWith(to + '/')
}

/**
 * Recursively check if any descendant route is currently active.
 */
function isGroupActive(item: IMenuItem): boolean {
  if (!item.children) return false
  return item.children.some(
    (child) => (child.to ? isActive(child.to) : false) || isGroupActive(child),
  )
}

function isGroupExpanded(item: IMenuItem, depth: number): boolean {
  return props.expandedGroups.has(groupKey(item, depth)) || isGroupActive(item)
}
</script>

<template>
  <!-- Title item -->
  <div v-if="item.isTitle" class="pt-4 pb-1">
    <span
      v-if="isExpanded"
      class="px-4 text-xs font-semibold text-slate-500 uppercase tracking-wider"
    >
      {{ item.label }}
    </span>
    <div v-else class="mx-4 my-2 h-px bg-white/10" />
  </div>

  <!-- Leaf item (no children) -->
  <router-link
    v-else-if="!item.children && item.to"
    :to="item.to"
    class="sidebar-item flex items-center text-slate-400 hover:text-white hover:bg-white/10 transition-all duration-200 whitespace-nowrap overflow-hidden"
    :class="[
      depth === 0
        ? isExpanded
          ? 'px-4 py-2.5 rounded-xl gap-3'
          : 'justify-center py-2.5 rounded-xl'
        : 'px-4 py-2 rounded-lg gap-3',
      isActive(item.to!)
        ? depth === 0
          ? '!text-white !bg-white/5 sidebar-item-active'
          : '!text-white !bg-blue-500/10'
        : '',
    ]"
    @click="emit('itemClick')"
  >
    <!-- Icon only on root level -->
    <component :is="item.icon" v-if="depth === 0 && item.icon" class="w-5 h-5 shrink-0" />
    <!-- Dot indicator for nested items -->
    <span
      v-if="depth > 0"
      class="w-1.5 h-1.5 rounded-full bg-current shrink-0 transition-all duration-200"
      :class="
        isActive(item.to!)
          ? 'opacity-100 text-white shadow-[0_0_8px_rgba(255,255,255,0.3)]'
          : 'opacity-50'
      "
    />
    <span v-if="isExpanded" :class="depth === 0 ? 'font-medium text-sm' : 'text-sm'">
      {{ item.label }}
    </span>
  </router-link>

  <!-- Group item (has children) -->
  <div v-else class="space-y-0.5">
    <button
      class="w-full flex items-center transition-all duration-200 whitespace-nowrap overflow-hidden"
      :class="[
        depth === 0
          ? isExpanded
            ? 'px-4 py-2.5 rounded-xl gap-3'
            : 'justify-center py-2.5 rounded-xl'
          : 'px-4 py-2 rounded-lg gap-3',
        isGroupActive(item)
          ? depth === 0
            ? 'text-white bg-white/5 sidebar-group-active'
            : 'text-blue-400 bg-blue-500/5'
          : 'text-slate-400 hover:text-white hover:bg-white/10',
      ]"
      @click="emit('toggleGroup', groupKey(item, depth))"
    >
      <!-- Icon only on root level -->
      <component :is="item.icon" v-if="depth === 0 && item.icon" class="w-5 h-5 shrink-0" />
      <!-- Dot indicator for nested group headers -->
      <span
        v-if="depth > 0"
        class="w-1.5 h-1.5 rounded-full bg-current shrink-0 transition-all duration-200"
        :class="
          isGroupActive(item)
            ? 'opacity-100 text-blue-400 shadow-[0_0_8px_rgba(96,165,250,0.5)]'
            : 'opacity-50'
        "
      />
      <span
        v-if="isExpanded"
        :class="[depth === 0 ? 'font-medium text-sm' : 'text-sm', 'flex-1 text-left']"
      >
        {{ item.label }}
      </span>
      <PhCaretDown
        v-if="isExpanded"
        class="w-4 h-4 transition-transform duration-200"
        :class="isGroupExpanded(item, depth) ? 'rotate-180' : ''"
      />
    </button>

    <!-- Children (only rendered when sidebar is expanded) -->
    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <div
        v-show="isExpanded && isGroupExpanded(item, depth)"
        class="pl-5 space-y-0.5 overflow-hidden"
      >
        <SidebarMenuItem
          v-for="(child, idx) in item.children"
          :key="idx"
          :item="child"
          :depth="depth + 1"
          :expanded-groups="expandedGroups"
          :is-expanded="isExpanded"
          @toggle-group="(key: string) => emit('toggleGroup', key)"
          @item-click="emit('itemClick')"
        />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.sidebar-item-active {
  box-shadow: inset 3px 0 0 #3b82f6;
}

.sidebar-group-active {
  box-shadow: inset 3px 0 0 #60a5fa; /* slightly lighter blue than the leaf item */
}
</style>
