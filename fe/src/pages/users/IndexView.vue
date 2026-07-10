<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted } from 'vue'
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
  PhShieldCheck,
  PhShieldSlash,
  PhProhibit,
  PhCheckCircle,
  PhLockKey,
  PhLockKeyOpen,
} from '@phosphor-icons/vue'

const userStore = useUserStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const statusUpdatingId = ref<number | null>(null)
const unlockingId = ref<number | null>(null)

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
  formModalRef.value?.show(user)
}

async function handleDelete(id: number) {
  await userStore.remove(id)
}

async function updateStatus(id: number, status: 'active' | 'suspended') {
  statusUpdatingId.value = id
  try {
    await userStore.updateStatus(id, status)
    await userStore.fetchAll(userStore.indexData.pagination.page)
  } finally {
    statusUpdatingId.value = null
  }
}

function handleSuspend(id: number) {
  updateStatus(id, 'suspended')
}

function handleActivate(id: number) {
  updateStatus(id, 'active')
}

function isLocked(user: IUser): boolean {
  return !!user.locked_until && new Date(user.locked_until).getTime() > Date.now()
}

async function handleUnlock(id: number) {
  unlockingId.value = id
  try {
    await userStore.unlock(id)
    await userStore.fetchAll(userStore.indexData.pagination.page)
  } finally {
    unlockingId.value = null
  }
}

function handlePageChange(page: number) {
  userStore.fetchAll(page)
}

onMounted(() => {
  userStore.fetchAll()
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
          :padded="false"
          :classes="{
            wrapper: 'h-full',
            card: `group hover:shadow-md transition-all duration-200 h-full rounded-t-2xl border-t-4 ${
              user.status === 'active' ? 'border-emerald-500' : 'border-red-500'
            }`,
            body: 'flex flex-col h-full',
          }"
        >
          <!-- Content -->
          <div class="flex flex-col items-center text-center p-5 flex-1">
            <!-- Avatar -->
            <div class="w-16 h-16 rounded-full overflow-hidden shrink-0 shadow-sm">
              <img
                v-if="user.avatar"
                :src="user.avatar"
                :alt="user.name"
                class="w-full h-full object-cover"
              />
              <div
                v-else
                :class="[
                  'flex items-center justify-center w-full h-full text-lg font-bold',
                  getAvatarBg(index),
                ]"
              >
                <span class="text-white">{{ getInitials(user.name) }}</span>
              </div>
            </div>

            <!-- Name + Email -->
            <h3 class="mt-3 text-sm font-semibold text-gray-900 truncate w-full">
              {{ user.name }}
            </h3>
            <p
              class="text-xs text-gray-500 truncate flex items-center justify-center gap-1 mt-0.5 w-full"
            >
              <PhEnvelope class="w-3.5 h-3.5 shrink-0" />
              <span class="truncate">{{ user.email }}</span>
            </p>

            <!-- Divider -->
            <div class="my-4 border-t border-gray-100 w-full"></div>

            <!-- Info Section -->
            <div class="w-full space-y-3 text-left">
              <!-- Joined Date -->
              <div class="flex items-center gap-2 text-xs">
                <div
                  class="flex items-center justify-center w-6 h-6 rounded-md bg-gray-100 text-gray-500 shrink-0"
                >
                  <PhCalendar class="w-3.5 h-3.5" />
                </div>
                <span class="text-gray-600">{{ formatDate(user.created_at) }}</span>

                <span
                  class="ml-auto inline-flex items-center justify-center w-6 h-6 rounded-full"
                  :class="
                    user.status === 'active'
                      ? 'bg-emerald-50 text-emerald-700'
                      : 'bg-red-50 text-red-600'
                  "
                  :title="user.status === 'active' ? 'Aktif' : 'Suspended'"
                >
                  <PhShieldCheck v-if="user.status === 'active'" class="w-3.5 h-3.5" />
                  <PhShieldSlash v-else class="w-3.5 h-3.5" />
                </span>

                <span
                  class="inline-flex items-center justify-center w-6 h-6 rounded-full"
                  :class="
                    isLocked(user) ? 'bg-red-50 text-red-600' : 'bg-emerald-50 text-emerald-700'
                  "
                  :title="isLocked(user) ? 'Locked' : 'Unlocked'"
                >
                  <PhLockKey class="w-3.5 h-3.5" />
                </span>
              </div>

              <!-- Roles -->
              <div class="flex items-center gap-2">
                <div
                  class="flex items-center justify-center w-6 h-6 rounded-md bg-blue-100 text-blue-600 shrink-0"
                >
                  <PhCrown class="w-3.5 h-3.5" />
                </div>
                <div class="flex flex-wrap gap-1 min-w-0 flex-1">
                  <span
                    v-for="role in user.roles?.slice(0, 2)"
                    :key="role.id"
                    class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-blue-50 text-blue-700 border border-blue-100 truncate"
                  >
                    {{ role.name }}
                  </span>
                  <span
                    v-if="user.roles && user.roles.length > 2"
                    class="inline-flex items-center px-1.5 py-0.5 rounded-md text-xs font-medium bg-gray-100 text-gray-500"
                  >
                    +{{ user.roles.length - 2 }}
                  </span>
                  <span v-if="!user.roles?.length" class="text-xs text-gray-400">Tanpa role</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer Actions -->
          <div
            class="border-t border-gray-100 bg-gray-50/80 px-3 py-2 flex items-center justify-end gap-0.5 flex-wrap"
          >
            <button
              v-if="hasAllPermissions(['user.edit']) && user.status !== 'suspended'"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-amber-600 hover:bg-amber-50 rounded-md transition-colors disabled:opacity-50"
              :disabled="statusUpdatingId === user.id"
              @click="handleSuspend(user.id)"
            >
              <PhProhibit class="w-3.5 h-3.5" />
              Suspend
            </button>
            <button
              v-if="hasAllPermissions(['user.edit']) && user.status === 'suspended'"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded-md transition-colors disabled:opacity-50"
              :disabled="statusUpdatingId === user.id"
              @click="handleActivate(user.id)"
            >
              <PhCheckCircle class="w-3.5 h-3.5" />
              Aktifkan
            </button>
            <button
              v-if="hasAllPermissions(['user.edit']) && isLocked(user)"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-violet-600 hover:bg-violet-50 rounded-md transition-colors disabled:opacity-50"
              :disabled="unlockingId === user.id"
              @click="handleUnlock(user.id)"
            >
              <PhLockKeyOpen class="w-3.5 h-3.5" />
              Unlock
            </button>
            <button
              v-if="hasAllPermissions(['user.edit'])"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
              @click="openEdit(user)"
            >
              <PhPencil class="w-3.5 h-3.5" />
              Edit
            </button>
            <button
              v-if="hasAllPermissions(['user.delete'])"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors disabled:opacity-50"
              :disabled="userStore.loading.Delete"
              @click="handleDelete(user.id)"
            >
              <PhTrash class="w-3.5 h-3.5" />
              Hapus
            </button>
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
