<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'
import type { IUser } from '@/stores/user'
import { usePermission } from '@/composables'
import {
  PhPlus,
  PhPencil,
  PhTrash,
  PhEnvelope,
  PhCalendar,
  PhCrown,
  PhDotsThreeVertical,
} from '@phosphor-icons/vue'

const userStore = useUserStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const openMenuId = ref<number | null>(null)

const colorPalette = [
  'bg-blue-500',
  'bg-emerald-500',
  'bg-violet-500',
  'bg-amber-500',
  'bg-rose-500',
  'bg-cyan-500',
  'bg-indigo-500',
  'bg-teal-500',
]

function getInitials(name: string): string {
  const words = name.trim().split(/\s+/)
  if (words.length === 1) return words[0].charAt(0).toUpperCase()
  return (words[0].charAt(0) + words[words.length - 1].charAt(0)).toUpperCase()
}

function getAvatarBg(index: number): string {
  return colorPalette[index % colorPalette.length]
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(user: IUser) {
  openMenuId.value = null
  formModalRef.value?.show(user)
}

async function handleDelete(id: number) {
  openMenuId.value = null
  await userStore.remove(id)
}

function toggleMenu(userId: number, event: Event) {
  event.stopPropagation()
  openMenuId.value = openMenuId.value === userId ? null : userId
}

function handleClickOutside() {
  openMenuId.value = null
}

function handlePageChange(page: number) {
  userStore.fetchAll(page)
}

onMounted(() => {
  userStore.fetchAll()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Users</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Kelola daftar user dalam sistem.</p>
      </div>
      <UiButton v-if="hasAllPermissions(['user.create'])" size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah User
      </UiButton>
    </div>

    <!-- Loading State -->
    <div
      v-if="userStore.loading.Index"
      class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
    >
      <UiSkeleton
        v-for="i in userStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-40"
        rounded
      />
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="userStore.indexData.items.length === 0"
      :icon="PhPlus"
      title="Belum ada User"
      description="Silakan buat user baru untuk mulai mengelola akses pengguna."
    >
      <UiButton v-if="hasAllPermissions(['user.create'])" size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat User Pertama
      </UiButton>
    </UiEmptyState>

    <!-- Data List -->
    <template v-else>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <UiCard
          v-for="(user, index) in userStore.indexData.items"
          :key="user.id"
          :classes="{
            wrapper: 'h-full relative',
            card: 'group hover:shadow-md transition-all duration-200 h-full border-t-4 border-blue-500',
            body: 'p-4',
          }"
        >
          <!-- Content -->
          <div class="flex flex-col">
            <!-- Avatar + Name + Email + Action -->
            <div class="flex items-center gap-3 mb-0">
              <div class="w-10 h-10 rounded-full overflow-hidden shrink-0 shadow-sm">
                <img
                  v-if="user.avatar"
                  :src="user.avatar"
                  :alt="user.name"
                  class="w-full h-full object-cover"
                />
                <div
                  v-else
                  :class="[
                    'flex items-center justify-center w-full h-full text-sm font-bold',
                    getAvatarBg(index),
                  ]"
                >
                  <span class="text-white">{{ getInitials(user.name) }}</span>
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="text-sm font-semibold text-gray-900 truncate">
                  {{ user.name }}
                </h3>
                <p class="text-xs text-gray-500 truncate flex items-center gap-1 leading-none">
                  <PhEnvelope class="w-3 h-3 shrink-0 mt-0.5" />
                  <span class="truncate">{{ user.email }}</span>
                </p>
              </div>
              <!-- 3-Dot Menu -->
              <div class="relative shrink-0">
                <button
                  class="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                  @click="toggleMenu(user.id, $event)"
                >
                  <PhDotsThreeVertical class="w-5 h-5" />
                </button>

                <!-- Dropdown Menu -->
                <div
                  v-if="openMenuId === user.id"
                  class="absolute right-0 mt-1 w-32 bg-white rounded-lg shadow-lg border border-gray-100 py-1 z-10"
                  @click.stop
                >
                  <button
                    v-if="hasAllPermissions(['user.edit'])"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                    @click="openEdit(user)"
                  >
                    <PhPencil class="w-4 h-4 text-blue-600" />
                    <span>Edit</span>
                  </button>
                  <button
                    v-if="hasAllPermissions(['user.delete'])"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                    :disabled="userStore.loading.Delete"
                    @click="handleDelete(user.id)"
                  >
                    <PhTrash class="w-4 h-4" />
                    <span>Hapus</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Divider -->
            <div class="my-3 border-t border-gray-100"></div>

            <!-- Info Section -->
            <div class="space-y-3">
              <!-- Joined Date -->
              <div class="flex items-center gap-2 text-xs">
                <div
                  class="flex items-center justify-center w-6 h-6 rounded-md bg-gray-100 text-gray-500 shrink-0"
                >
                  <PhCalendar class="w-3.5 h-3.5" />
                </div>
                <span class="text-gray-600">{{ formatDate(user.created_at) }}</span>
              </div>

              <!-- Roles -->
              <div v-if="user.roles?.length" class="flex items-center gap-2">
                <div
                  class="flex items-center justify-center w-6 h-6 rounded-md bg-blue-100 text-blue-600 shrink-0"
                >
                  <PhCrown class="w-3.5 h-3.5" />
                </div>
                <div class="flex flex-wrap gap-1 min-w-0 flex-1">
                  <span
                    v-for="role in user.roles.slice(0, 2)"
                    :key="role.id"
                    class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-blue-50 text-blue-700 border border-blue-100 truncate"
                  >
                    {{ role.name }}
                  </span>
                  <span
                    v-if="user.roles.length > 2"
                    class="inline-flex items-center px-1.5 py-0.5 rounded-md text-xs font-medium bg-gray-100 text-gray-500"
                  >
                    +{{ user.roles.length - 2 }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </UiCard>
      </div>

      <!-- Pagination -->
      <div class="mt-6 flex justify-center">
        <UiPagination
          :page="userStore.indexData.pagination.page"
          :total-pages="userStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
