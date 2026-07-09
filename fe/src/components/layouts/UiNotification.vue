<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { PhBell, PhCheck, PhTrash, PhCircle } from '@phosphor-icons/vue'
import { useNotificationStore } from '@/stores/notification'
import { formatTimeForHuman } from '@/helpers/date'

const notificationStore = useNotificationStore()
const isOpen = ref(false)
const scrollContainerRef = ref<HTMLElement | null>(null)
const sentinelRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const displayNotifications = computed(() =>
  notificationStore.notifications.map((n) => ({
    id: n.id,
    title: n.title,
    message: n.message,
    time: formatTimeForHuman(n.created_at),
    isRead: !notificationStore.isUnread(n),
    type: notificationStore.getNotificationType(n.type),
  })),
)

const unreadCount = computed(() => notificationStore.unreadCount)
const isLoadingMore = computed(() => notificationStore.loading.LoadMore)
const hasMoreNotifications = computed(() => notificationStore.hasMore)

function closeOthers() {
  window.dispatchEvent(new CustomEvent('close-dropdown', { detail: 'notification' }))
}

function close() {
  isOpen.value = false
}

function toggle() {
  if (isOpen.value) {
    close()
  } else {
    closeOthers()
    isOpen.value = true
    // Reset and fetch notifications when opening
    notificationStore.resetNotifications()
    notificationStore.fetchNotifications()
  }
}

async function handleMarkRead(id: number) {
  await notificationStore.markAsRead(id)
}

async function handleMarkAllRead() {
  await notificationStore.markAllAsRead()
}

async function handleDelete(id: number) {
  await notificationStore.deleteNotification(id)
}

async function handleClick(notification: { id: number; isRead: boolean }) {
  if (!notification.isRead) {
    await handleMarkRead(notification.id)
  }
}

function typeIcon(type?: string) {
  switch (type) {
    case 'success':
      return PhCheck
    case 'warning':
      return PhCircle
    case 'error':
      return PhCircle
    default:
      return PhCircle
  }
}

function typeColor(type?: string) {
  switch (type) {
    case 'success':
      return 'bg-green-100 text-green-600'
    case 'warning':
      return 'bg-yellow-100 text-yellow-600'
    case 'error':
      return 'bg-red-100 text-red-600'
    default:
      return 'bg-blue-100 text-blue-600'
  }
}

function handleOutsideClose(event: Event) {
  const detail = (event as CustomEvent).detail
  if (detail !== 'notification') {
    close()
  }
}

function setupIntersectionObserver() {
  if (observer) {
    observer.disconnect()
  }

  observer = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      if (entry.isIntersecting && hasMoreNotifications.value && !isLoadingMore.value) {
        notificationStore.loadMoreNotifications()
      }
    },
    {
      root: scrollContainerRef.value,
      rootMargin: '50px',
      threshold: 0.1,
    },
  )

  if (sentinelRef.value) {
    observer.observe(sentinelRef.value)
  }
}

watch(
  () => isOpen.value,
  async (newVal) => {
    if (newVal) {
      await nextTick()
      setupIntersectionObserver()
    } else {
      if (observer) {
        observer.disconnect()
        observer = null
      }
    }
  },
)

onMounted(() => {
  window.addEventListener('close-dropdown', handleOutsideClose)
  // Fetch unread count on mount
  notificationStore.fetchUnreadCount()
})

onUnmounted(() => {
  window.removeEventListener('close-dropdown', handleOutsideClose)
  if (observer) {
    observer.disconnect()
    observer = null
  }
})
</script>

<template>
  <div class="relative">
    <button
      type="button"
      class="relative p-2 rounded-full hover:bg-gray-100 transition-colors"
      @click.stop="toggle"
    >
      <PhBell class="h-5 w-5 text-gray-600" />
      <span
        v-if="unreadCount > 0"
        class="absolute top-1 right-1 min-w-[0.5rem] h-2 px-1 flex items-center justify-center bg-red-500 rounded-full"
      >
        <span v-if="unreadCount > 99" class="text-[0.5rem] text-white font-medium">99+</span>
      </span>
    </button>

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
        class="absolute right-0 mt-2 w-80 bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden"
        @click.stop
      >
        <!-- Header -->
        <div class="px-4 py-3 bg-linear-to-br from-gray-50 to-white border-b border-gray-100">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-gray-900">
              Notifications
              <span
                v-if="unreadCount > 0"
                class="ml-1.5 px-2 py-0.5 text-xs font-medium bg-red-100 text-red-600 rounded-full"
              >
                {{ unreadCount }}
              </span>
            </h3>
            <div class="flex items-center gap-1">
              <button
                v-if="displayNotifications.length > 0"
                class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
                title="Mark all as read"
                @click="handleMarkAllRead"
              >
                <PhCheck class="h-4 w-4" />
              </button>
              <button
                v-if="displayNotifications.length > 0"
                class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                title="Clear all"
                @click="notificationStore.deleteAllNotifications()"
              >
                <PhTrash class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Loading State (Initial Load) -->
        <div
          v-if="notificationStore.loading.Index && displayNotifications.length === 0"
          class="px-4 py-8"
        >
          <div class="flex items-center justify-center">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          </div>
        </div>

        <!-- Notification List -->
        <div v-else ref="scrollContainerRef" class="max-h-96 overflow-y-auto">
          <div v-if="displayNotifications.length === 0" class="px-4 py-8 text-center">
            <PhBell class="h-8 w-8 text-gray-300 mx-auto mb-2" />
            <p class="text-sm text-gray-500">Tidak ada notifikasi</p>
          </div>

          <template v-else>
            <button
              v-for="notification in displayNotifications"
              :key="notification.id"
              :class="[
                'w-full px-4 py-3 text-left border-b border-gray-50 hover:bg-gray-50 transition-colors',
                !notification.isRead && 'bg-blue-50/50',
              ]"
              @click="handleClick(notification)"
            >
              <div class="flex items-start gap-3">
                <!-- Icon -->
                <div
                  :class="[
                    'w-8 h-8 rounded-lg flex items-center justify-center shrink-0',
                    typeColor(notification.type),
                  ]"
                >
                  <component :is="typeIcon(notification.type)" class="h-4 w-4" />
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between gap-2">
                    <p
                      :class="[
                        'text-sm truncate',
                        !notification.isRead
                          ? 'font-semibold text-gray-900'
                          : 'font-medium text-gray-700',
                      ]"
                    >
                      {{ notification.title }}
                    </p>
                    <button
                      class="p-1 rounded text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors shrink-0"
                      @click.stop="handleDelete(notification.id)"
                    >
                      <PhTrash class="h-3.5 w-3.5" />
                    </button>
                  </div>
                  <p class="mt-0.5 text-xs text-gray-500 line-clamp-2">
                    {{ notification.message }}
                  </p>
                  <p class="mt-1 text-xs text-gray-400">
                    {{ notification.time }}
                  </p>
                </div>

                <!-- Unread indicator -->
                <span
                  v-if="!notification.isRead"
                  class="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0"
                />
              </div>
            </button>

            <!-- Load More Sentinel -->
            <div ref="sentinelRef" class="py-3">
              <div v-if="isLoadingMore" class="flex items-center justify-center">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
                <span class="ml-2 text-xs text-gray-500">Memuat...</span>
              </div>
              <div v-else-if="!hasMoreNotifications" class="text-center">
                <p class="text-xs text-gray-400">Tidak ada notifikasi lagi</p>
              </div>
            </div>
          </template>
        </div>
      </div>
    </Transition>
  </div>
</template>
