<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { PhUser, PhSignOut } from '@phosphor-icons/vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const isOpen = ref(false)

const userInitials = computed(() => {
  const name = authStore.user?.name || 'User'
  return name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

function closeOthers() {
  window.dispatchEvent(new CustomEvent('close-dropdown', { detail: 'profile' }))
}

function toggle() {
  if (isOpen.value) {
    isOpen.value = false
  } else {
    closeOthers()
    isOpen.value = true
  }
}

function close() {
  isOpen.value = false
}

function handleOutsideClose(event: Event) {
  const detail = (event as CustomEvent).detail
  if (detail !== 'profile') {
    close()
  }
}

async function handleLogout() {
  await authStore.logout()
  router.push({ name: 'Login' })
}

onMounted(() => {
  window.addEventListener('close-dropdown', handleOutsideClose)
})

onUnmounted(() => {
  window.removeEventListener('close-dropdown', handleOutsideClose)
})
</script>

<template>
  <div class="relative">
    <button
      type="button"
      class="group flex items-center gap-2.5 p-1.5 pr-3 rounded-full transition-all duration-200"
      @click.stop="toggle"
    >
      <!-- Avatar with online indicator -->
      <div class="relative">
        <div
          class="w-7 h-7 rounded-full overflow-hidden ring-2 ring-emerald-400 ring-offset-1 transition-transform duration-200 group-hover:scale-105"
        >
          <img
            v-if="authStore.user?.avatar"
            :src="authStore.user.avatar"
            class="w-full h-full object-cover"
          />
          <div
            v-else
            class="w-full h-full bg-linear-to-br from-blue-500 to-violet-500 flex items-center justify-center"
          >
            <span class="text-xs font-bold text-white leading-none">{{ userInitials }}</span>
          </div>
        </div>
        <!-- Online dot -->
        <span
          class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 bg-emerald-400 rounded-full border-2 border-white"
        />
      </div>
    </button>

    <!-- Dropdown Menu -->
    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 scale-95 -translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-1"
    >
      <div
        v-if="isOpen"
        v-click-outside="close"
        class="absolute right-0 mt-2 w-64 bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden"
        @click.stop
      >
        <!-- User Info Header -->
        <div class="px-5 py-4 bg-linear-to-br from-gray-50 to-white border-b border-gray-100">
          <div class="flex items-center gap-3">
            <!-- Avatar -->
            <div class="relative shrink-0">
              <div
                class="w-11 h-11 rounded-full overflow-hidden ring-2 ring-emerald-400 ring-offset-2"
              >
                <img
                  v-if="authStore.user?.avatar"
                  :src="authStore.user.avatar"
                  class="w-full h-full object-cover"
                />
                <div
                  v-else
                  class="w-full h-full bg-linear-to-br from-blue-500 to-violet-500 flex items-center justify-center"
                >
                  <span class="text-sm font-bold text-white leading-none">{{ userInitials }}</span>
                </div>
              </div>
              <span
                class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-emerald-400 rounded-full border-2 border-white"
              />
            </div>
            <!-- Name & Email -->
            <div class="min-w-0 flex-1">
              <p class="text-sm font-semibold text-gray-900 truncate">
                {{ authStore.user?.name }}
              </p>
              <p class="text-xs text-gray-500 truncate">
                {{ authStore.user?.email }}
              </p>
            </div>
          </div>
        </div>

        <!-- Menu Items -->
        <div class="p-2">
          <router-link
            to="/profile"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-all duration-150 group/item"
            @click="close"
          >
            <div
              class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center group-hover/item:bg-blue-100 transition-colors"
            >
              <PhUser class="w-4 h-4 text-blue-500" />
            </div>
            <span class="font-medium">Profile</span>
          </router-link>
        </div>

        <!-- Logout -->
        <div class="p-2 pt-0">
          <hr class="mb-2 border-gray-100" />
          <button
            class="flex items-center justify-center gap-2 w-full px-3 py-2.5 rounded-xl text-sm font-medium text-gray-500 hover:bg-red-50 hover:text-red-600 transition-all duration-150"
            @click="handleLogout"
          >
            <PhSignOut class="w-4 h-4" />
            <span>Log out</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
