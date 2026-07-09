<script setup lang="ts">
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted } from 'vue'
import { useRoleStore } from '@/stores/role'
import type { IRole } from '@/stores/role'
import { usePermission } from '@/composables'
import { PhPlus, PhPencil, PhTrash } from '@phosphor-icons/vue'

const roleStore = useRoleStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

const avatarColors = [
  'bg-blue-600',
  'bg-emerald-600',
  'bg-violet-600',
  'bg-amber-600',
  'bg-rose-600',
  'bg-cyan-600',
  'bg-indigo-600',
  'bg-teal-600',
]

function getInitials(name: string): string {
  const words = name.trim().split(/\s+/)
  if (words.length === 1) return words[0].charAt(0).toUpperCase()
  return (words[0].charAt(0) + words[words.length - 1].charAt(0)).toUpperCase()
}

function getAvatarColor(index: number): string {
  return avatarColors[index % avatarColors.length]
}

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(role: IRole) {
  formModalRef.value?.show(role)
}

async function handleDelete(id: number) {
  await roleStore.remove(id)
}

function handlePageChange(page: number) {
  roleStore.fetchAll(page)
}

onMounted(() => {
  roleStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Roles</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">Kelola daftar role dalam sistem.</p>
      </div>
      <UiButton v-if="hasAllPermissions(['role.create'])" size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Role
      </UiButton>
    </div>

    <!-- Loading State -->
    <div
      v-if="roleStore.loading.Index"
      class="grid gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton
        v-for="i in roleStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-48"
        rounded
      />
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="roleStore.indexData.items.length === 0"
      :icon="PhPlus"
      title="Belum ada Role"
      description="Silakan buat role baru untuk mulai mengatur hak akses sistem."
    >
      <UiButton v-if="hasAllPermissions(['role.create'])" size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Role Pertama
      </UiButton>
    </UiEmptyState>

    <!-- Data List -->
    <template v-else>
      <div class="grid gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="(role, index) in roleStore.indexData.items"
          :key="role.id"
          :classes="{
            wrapper: 'group hover:shadow-md transition-shadow h-full',
            card: 'h-full flex flex-col',
            body: 'flex flex-col flex-1 p-6',
          }"
        >
          <!-- Top: Avatar + Name + Actions -->
          <div class="flex items-center gap-3">
            <div
              :class="[
                'flex items-center justify-center w-11 h-11 rounded-full text-white text-sm font-bold shrink-0 shadow-sm',
                getAvatarColor(index),
              ]"
            >
              {{ getInitials(role.name) }}
            </div>
            <h3 class="text-lg font-semibold text-gray-900 truncate min-w-0 flex-1">
              {{ role.name }}
            </h3>
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 shrink-0"
            >
              <button
                v-if="hasAllPermissions(['role.edit'])"
                class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                title="Edit"
                @click="openEdit(role)"
              >
                <PhPencil class="w-5 h-5" />
              </button>
              <button
                v-if="hasAllPermissions(['role.delete'])"
                class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                title="Hapus"
                :disabled="roleStore.loading.Delete"
                @click="handleDelete(role.id)"
              >
                <PhTrash class="w-5 h-5" />
              </button>
            </div>
          </div>

          <!-- Description -->
          <p class="mt-2 text-sm text-gray-600 line-clamp-2">
            {{ role.description }}
          </p>

          <!-- Permissions badges (pushed to bottom) -->
          <div v-if="role.permissions?.length" class="mt-auto pt-3 flex flex-wrap gap-1.5">
            <span
              v-for="perm in role.permissions.slice(0, 5)"
              :key="perm.id"
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700"
            >
              {{ perm.name }}
            </span>
            <span
              v-if="role.permissions.length > 5"
              class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500"
            >
              +{{ role.permissions.length - 5 }} lainnya
            </span>
          </div>
        </UiCard>
      </div>

      <!-- Pagination -->
      <div class="mt-8 flex justify-center">
        <UiPagination
          :page="roleStore.indexData.pagination.page"
          :total-pages="roleStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
