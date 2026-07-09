<script setup lang="ts">
import { UiNotification, UiUserDropdown } from './index'

defineProps<{
  showHamburger: boolean
}>()

const emit = defineEmits<{
  toggleSidebar: []
}>()
</script>

<template>
  <header class="sticky top-0 z-30 bg-white shadow-sm">
    <div class="flex items-center justify-between px-4 h-12">
      <!-- Left: Hamburger -->
      <div class="flex items-center gap-4">
        <button
          v-if="showHamburger"
          class="p-2 rounded-md hover:bg-gray-100 transition-colors"
          @click="emit('toggleSidebar')"
        >
          <slot name="menu-icon" />
        </button>
      </div>

      <!-- Right: Notification & Profile -->
      <div class="flex items-center gap-2">
        <!-- Notification Bell -->
        <UiNotification />

        <!-- User Dropdown -->
        <UiUserDropdown />
      </div>
    </div>
  </header>
</template>
