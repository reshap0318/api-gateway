<script setup lang="ts">
import { UiButton, UiTable, UiPagination, UiBadge, FormSelect } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, computed, onMounted, watch } from 'vue'
import { useGatewayRouteStore } from '@/stores/gatewayRoute'
import type { IGatewayRoute } from '@/stores/gatewayRoute'
import { useGatewayServiceStore } from '@/stores/gatewayService'
import type { IGatewayService } from '@/stores/gatewayService'
import { usePermission } from '@/composables'
import swal from '@/plugins/swal'
import { PhPlus, PhPencil, PhTrash, PhArrowsClockwise } from '@phosphor-icons/vue'
import type { TTableColumn } from '@/components/utils/types'

const routeStore = useGatewayRouteStore()
const serviceStore = useGatewayServiceStore()
const { hasAllPermissions, hasAnyPermission } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)

const refreshing = ref(false)
const allServices = ref<IGatewayService[]>([])

const columns: TTableColumn[] = [
  { title: 'Service', data: 'service', sortable: false },
  { title: 'Method', data: 'method', sortable: false },
  { title: 'Path Pattern', data: 'path_pattern', sortable: false },
  { title: 'Permissions', data: 'permissions', sortable: false },
  { title: 'Match Mode', data: 'permission_match_mode', sortable: false },
  { title: 'Rate Limit', data: 'rate_limit_per_minute', sortable: false },
  { title: 'Status', data: 'is_active', sortable: false },
  { title: 'Aksi', data: 'actions', sortable: false, class: 'text-right' },
]

const rows = computed(() => routeStore.indexData.items as unknown as Record<string, unknown>[])

// Filters (§FSD 3.4)
const serviceFilter = ref('')
const methodFilter = ref('')
const activeFilter = ref('')

const serviceOptions = computed(() => [
  { value: '', label: 'Semua Service' },
  ...allServices.value.map((s) => ({ value: String(s.id), label: s.name })),
])
const methodOptions = [
  { value: '', label: 'Semua Method' },
  ...['GET', 'POST', 'PUT', 'PATCH', 'DELETE', '*'].map((m) => ({ value: m, label: m })),
]
const activeOptions = [
  { value: '', label: 'Semua Status' },
  { value: 'true', label: 'Aktif' },
  { value: 'false', label: 'Nonaktif' },
]

watch([serviceFilter, methodFilter, activeFilter], () => fetchRoutes(1))

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(route: IGatewayRoute) {
  formModalRef.value?.show(route)
}

async function handleDelete(id: number) {
  await routeStore.remove(id)
}

function handlePageChange(page: number) {
  fetchRoutes(page)
}

function fetchRoutes(page?: number) {
  routeStore.fetchAll(page, {
    service: serviceFilter.value || undefined,
    method: methodFilter.value || undefined,
    is_active: activeFilter.value || undefined,
  })
}

async function handleRefreshCache() {
  refreshing.value = true
  try {
    await routeStore.refreshCache()
    swal.success('Berhasil', 'Cache route berhasil di-refresh.')
  } catch {
    swal.error('Gagal', 'Gagal me-refresh cache route.')
  } finally {
    refreshing.value = false
  }
}

onMounted(async () => {
  allServices.value = await serviceStore.fetchAllServices()
  fetchRoutes()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Gateway Routes</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola aturan routing (method + path pattern) dan permission per route.
        </p>
      </div>
      <div class="flex items-center gap-2">
        <UiButton
          v-if="hasAnyPermission(['route.create', 'route.edit'])"
          size="sm"
          variant="secondary"
          outline
          :loading="refreshing"
          @click="handleRefreshCache"
        >
          <template #icon>
            <PhArrowsClockwise class="w-4 h-4" />
          </template>
          Refresh Cache Now
        </UiButton>
        <UiButton v-if="hasAllPermissions(['route.create'])" size="sm" @click="openCreate">
          <template #icon>
            <PhPlus class="w-4 h-4" />
          </template>
          Tambah Route
        </UiButton>
      </div>
    </div>

    <!-- Filters -->
    <div class="mb-4 grid grid-cols-1 sm:grid-cols-3 gap-3 max-w-2xl">
      <FormSelect v-model="serviceFilter" name="service_filter" :options="serviceOptions" />
      <FormSelect v-model="methodFilter" name="method_filter" :options="methodOptions" />
      <FormSelect v-model="activeFilter" name="active_filter" :options="activeOptions" />
    </div>

    <UiTable :columns="columns" :datas="rows" :is-loading="routeStore.loading.Index">
      <template #service="{ item }">
        <span class="font-medium text-gray-800">
          {{ (item as unknown as IGatewayRoute).service?.name }}
        </span>
      </template>

      <template #method="{ item }">
        <UiBadge color="info">{{ (item as unknown as IGatewayRoute).method }}</UiBadge>
      </template>

      <template #path_pattern="{ item }">
        <code class="text-xs bg-gray-100 px-2 py-1 rounded">
          {{ (item as unknown as IGatewayRoute).service?.base_path
          }}{{ (item as unknown as IGatewayRoute).path_pattern }}
        </code>
      </template>

      <template #permissions="{ item }">
        <div class="flex flex-wrap gap-1">
          <span
            v-if="!(item as unknown as IGatewayRoute).permissions?.length"
            class="text-xs text-gray-400 italic"
          >
            Publik
          </span>
          <span
            v-for="perm in (item as unknown as IGatewayRoute).permissions"
            :key="perm.id"
            class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700"
          >
            {{ perm.name }}
          </span>
        </div>
      </template>

      <template #permission_match_mode="{ item }">
        <UiBadge
          :color="
            (item as unknown as IGatewayRoute).permission_match_mode === 'all'
              ? 'warning'
              : 'primary'
          "
        >
          {{ (item as unknown as IGatewayRoute).permission_match_mode }}
        </UiBadge>
      </template>

      <template #rate_limit_per_minute="{ item }">
        <span class="text-gray-600">
          {{ (item as unknown as IGatewayRoute).rate_limit_per_minute ?? 'Default' }}
        </span>
      </template>

      <template #is_active="{ item }">
        <UiBadge :color="(item as unknown as IGatewayRoute).is_active ? 'primary' : 'danger'">
          {{ (item as unknown as IGatewayRoute).is_active ? 'Aktif' : 'Nonaktif' }}
        </UiBadge>
      </template>

      <template #actions="{ item }">
        <div class="flex justify-end gap-1">
          <button
            v-if="hasAllPermissions(['route.edit'])"
            class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
            title="Edit"
            @click="openEdit(item as unknown as IGatewayRoute)"
          >
            <PhPencil class="w-4 h-4" />
          </button>
          <button
            v-if="hasAllPermissions(['route.delete'])"
            class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
            title="Hapus"
            :disabled="routeStore.loading.Delete"
            @click="handleDelete((item as unknown as IGatewayRoute).id)"
          >
            <PhTrash class="w-4 h-4" />
          </button>
        </div>
      </template>

      <template #empty>
        <p class="text-gray-500">Belum ada route terdaftar. Klik "Tambah Route" untuk memulai.</p>
      </template>
    </UiTable>

    <!-- Pagination -->
    <div class="mt-6 flex justify-center">
      <UiPagination
        :page="routeStore.indexData.pagination.page"
        :total-pages="routeStore.indexData.pagination.total_pages"
        @update:page="handlePageChange"
      />
    </div>
  </div>

  <FormModal ref="formModalRef" />
</template>
