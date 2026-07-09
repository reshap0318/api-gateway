<script setup lang="ts">
import { UiButton, UiTable, UiPagination, UiBadge, FormInput, FormSelect } from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, computed, onMounted, watch } from 'vue'
import { useGatewayServiceStore } from '@/stores/gatewayService'
import type { IGatewayService } from '@/stores/gatewayService'
import { usePermission } from '@/composables'
import swal from '@/plugins/swal'
import {
  PhPlus,
  PhPencil,
  PhTrash,
  PhLink,
  PhMagnifyingGlass,
  PhHeartbeat,
} from '@phosphor-icons/vue'
import type { TTableColumn } from '@/components/utils/types'

const serviceStore = useGatewayServiceStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const healthChecking = ref<number | null>(null)

const columns: TTableColumn[] = [
  { title: 'Nama', data: 'name', sortable: false },
  { title: 'Base URL', data: 'base_url', sortable: false },
  { title: 'Protocol', data: 'protocol', sortable: false },
  { title: 'Rate Limit', data: 'rate_limit_per_minute', sortable: false },
  { title: 'Health Status', data: 'health_status', sortable: false },
  { title: 'Status', data: 'is_active', sortable: false },
  { title: 'Aksi', data: 'actions', sortable: false, class: 'text-right' },
]

const rows = computed(() => serviceStore.indexData.items as unknown as Record<string, unknown>[])

// Filters (§FSD 3.2) — search debounced, filters applied immediately.
const search = ref('')
const protocolFilter = ref('')
const activeFilter = ref('')
const healthStatusFilter = ref('')
let searchDebounce: ReturnType<typeof setTimeout> | undefined

const protocolOptions = [
  { value: '', label: 'Semua Protocol' },
  { value: 'http', label: 'HTTP' },
  { value: 'websocket', label: 'WebSocket' },
]
const activeOptions = [
  { value: '', label: 'Semua Status' },
  { value: 'true', label: 'Aktif' },
  { value: 'false', label: 'Nonaktif' },
]
const healthStatusOptions = [
  { value: '', label: 'Semua Health' },
  { value: 'up', label: 'Up' },
  { value: 'down', label: 'Down' },
  { value: 'unknown', label: 'Unknown' },
]

function healthBadgeColor(status: string): 'primary' | 'danger' | 'warning' {
  if (status === 'up') return 'primary'
  if (status === 'down') return 'danger'
  return 'warning'
}

function applyFilters(page = 1) {
  serviceStore.fetchAll(page, {
    search: search.value || undefined,
    protocol: protocolFilter.value || undefined,
    is_active: activeFilter.value || undefined,
    health_status: healthStatusFilter.value || undefined,
  })
}

watch(search, () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => applyFilters(1), 300)
})
watch([protocolFilter, activeFilter, healthStatusFilter], () => applyFilters(1))

function openCreate() {
  formModalRef.value?.show()
}

function openEdit(service: IGatewayService) {
  formModalRef.value?.show(service)
}

async function handleDelete(id: number) {
  await serviceStore.remove(id)
}

async function handleHealthCheck(id: number) {
  healthChecking.value = id
  try {
    await serviceStore.healthCheck(id)
  } catch {
    swal.error('Gagal', 'Gagal melakukan health check.')
  } finally {
    healthChecking.value = null
  }
}

function handlePageChange(page: number) {
  applyFilters(page)
}

onMounted(() => {
  applyFilters()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Gateway Services</h1>
        <p class="hidden sm:block text-sm text-gray-600 mt-1">
          Kelola upstream service yang bisa diakses lewat Gateway.
        </p>
      </div>
      <UiButton v-if="hasAllPermissions(['service.create'])" size="sm" @click="openCreate">
        <template #icon>
          <PhPlus class="w-4 h-4" />
        </template>
        Tambah Service
      </UiButton>
    </div>

    <!-- Search & Filters -->
    <div class="mb-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <FormInput
        v-model="search"
        name="search"
        placeholder="Cari nama service..."
        :prefix-icon="PhMagnifyingGlass"
      />
      <FormSelect v-model="protocolFilter" name="protocol_filter" :options="protocolOptions" />
      <FormSelect v-model="activeFilter" name="active_filter" :options="activeOptions" />
      <FormSelect
        v-model="healthStatusFilter"
        name="health_status_filter"
        :options="healthStatusOptions"
      />
    </div>

    <UiTable :columns="columns" :datas="rows" :is-loading="serviceStore.loading.Index">
      <template #base_url="{ item }">
        <span class="inline-flex items-center gap-1 text-gray-600">
          <PhLink class="w-3.5 h-3.5 shrink-0" />
          <span class="truncate max-w-xs">{{ (item as unknown as IGatewayService).base_url }}</span>
        </span>
      </template>

      <template #protocol="{ item }">
        <UiBadge
          :color="
            (item as unknown as IGatewayService).protocol === 'websocket' ? 'info' : 'primary'
          "
        >
          {{ (item as unknown as IGatewayService).protocol }}
        </UiBadge>
      </template>

      <template #rate_limit_per_minute="{ item }">
        <span class="text-gray-600">
          {{ (item as unknown as IGatewayService).rate_limit_per_minute ?? 'Default' }}
        </span>
      </template>

      <template #health_status="{ item }">
        <UiBadge :color="healthBadgeColor((item as unknown as IGatewayService).health_status)">
          {{ (item as unknown as IGatewayService).health_status }}
        </UiBadge>
      </template>

      <template #is_active="{ item }">
        <UiBadge :color="(item as unknown as IGatewayService).is_active ? 'primary' : 'danger'">
          {{ (item as unknown as IGatewayService).is_active ? 'Aktif' : 'Nonaktif' }}
        </UiBadge>
      </template>

      <template #actions="{ item }">
        <div class="flex justify-end gap-1">
          <button
            v-if="hasAllPermissions(['service.health-check'])"
            class="p-1.5 text-gray-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-md transition-colors disabled:opacity-50"
            title="Health Check"
            :disabled="healthChecking === (item as unknown as IGatewayService).id"
            @click="handleHealthCheck((item as unknown as IGatewayService).id)"
          >
            <PhHeartbeat class="w-4 h-4" />
          </button>
          <button
            v-if="hasAllPermissions(['service.edit'])"
            class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
            title="Edit"
            @click="openEdit(item as unknown as IGatewayService)"
          >
            <PhPencil class="w-4 h-4" />
          </button>
          <button
            v-if="hasAllPermissions(['service.delete'])"
            class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
            title="Hapus"
            :disabled="serviceStore.loading.Delete"
            @click="handleDelete((item as unknown as IGatewayService).id)"
          >
            <PhTrash class="w-4 h-4" />
          </button>
        </div>
      </template>

      <template #empty>
        <p class="text-gray-500">
          Belum ada service terdaftar. Klik "Tambah Service" untuk memulai.
        </p>
      </template>
    </UiTable>

    <!-- Pagination -->
    <div class="mt-6 flex justify-center">
      <UiPagination
        :page="serviceStore.indexData.pagination.page"
        :total-pages="serviceStore.indexData.pagination.total_pages"
        @update:page="handlePageChange"
      />
    </div>
  </div>

  <FormModal ref="formModalRef" />
</template>
