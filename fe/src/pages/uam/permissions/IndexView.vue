<script setup lang="ts">
import FormModal from './FormModal.vue'
import { UiCard, UiButton, UiPagination, UiEmptyState, UiSkeleton } from '@/components/utils'
import { ref, onMounted } from 'vue'
import { usePermissionStore, type IPermission } from '@/stores'
import { usePermission } from '@/composables'
import { PhPlus, PhPencil, PhTrash } from '@phosphor-icons/vue'

const permissionStore = usePermissionStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(permission: IPermission) {
  formModalRef.value?.show(permission)
}

async function handleDelete(id: number) {
  await permissionStore.remove(id)
}

function handlePageChange(page: number) {
  permissionStore.fetchAll(page)
}

onMounted(() => {
  permissionStore.fetchAll()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Permissions</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola daftar akses permission dalam sistem.
        </p>
      </div>
      <UiButton v-if="hasAllPermissions(['permission.create'])" size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Permission
      </UiButton>
    </div>

    <!-- Loading State -->
    <div
      v-if="permissionStore.loading.Index"
      class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
    >
      <UiSkeleton
        v-for="i in permissionStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-26"
        rounded
      />
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="permissionStore.indexData.items.length === 0"
      :icon="PhPlus"
      title="Belum ada Permission"
      description="Silakan buat permission baru untuk mulai mengatur hak akses sistem."
    >
      <UiButton v-if="hasAllPermissions(['permission.create'])" size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Buat Permission Pertama
      </UiButton>
    </UiEmptyState>

    <!-- No Search Results -->
    <!-- Data List -->
    <div v-else class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
      <UiCard
        v-for="permission in permissionStore.indexData.items"
        :key="permission.id"
        class="group hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-1">
              {{ permission.name }}
            </h3>
            <p class="text-sm text-gray-600 line-clamp-2">
              {{ permission.description }}
            </p>
          </div>
          <div
            class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
          >
            <button
              v-if="hasAllPermissions(['permission.edit'])"
              class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
              title="Edit"
              @click="openEdit(permission)"
            >
              <PhPencil class="w-5 h-5" />
            </button>
            <button
              v-if="hasAllPermissions(['permission.delete'])"
              class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
              title="Hapus"
              :disabled="permissionStore.loading.Delete"
              @click="handleDelete(permission.id)"
            >
              <PhTrash class="w-5 h-5" />
            </button>
          </div>
        </div>
      </UiCard>
    </div>

    <!-- Pagination -->
    <div class="mt-8 flex justify-center">
      <UiPagination
        :page="permissionStore.indexData.pagination.page"
        :total-pages="permissionStore.indexData.pagination.total_pages"
        @update:page="handlePageChange"
      />
    </div>
  </div>

  <FormModal ref="formModalRef" />
</template>
