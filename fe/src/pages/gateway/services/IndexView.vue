<script setup lang="ts">
import {
  UiCard,
  UiButton,
  UiPagination,
  UiEmptyState,
  UiSkeleton,
  FormInput,
  FormSelect,
} from '@/components/utils'
import FormModal from './FormModal.vue'
import { ref, onMounted, watch } from 'vue'
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
  PhClockCounterClockwise,
} from '@phosphor-icons/vue'

const serviceStore = useGatewayServiceStore()
const { hasAllPermissions } = usePermission()
const formModalRef = ref<InstanceType<typeof FormModal> | null>(null)
const healthChecking = ref<number | null>(null)

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

function healthLabel(status: IGatewayService['health_status']): string {
  if (status === 'up') return 'Operational'
  if (status === 'down') return 'Down'
  return 'Belum Diketahui'
}

function healthBannerClass(status: IGatewayService['health_status']): string {
  if (status === 'up') return 'bg-emerald-50 text-emerald-700'
  if (status === 'down') return 'bg-red-50 text-red-700'
  return 'bg-amber-50 text-amber-700'
}

function healthDotClass(status: IGatewayService['health_status']): string {
  if (status === 'up') return 'bg-emerald-500'
  if (status === 'down') return 'bg-red-500'
  return 'bg-amber-400'
}

function healthBorderClass(status: IGatewayService['health_status']): string {
  if (status === 'up') return 'border-emerald-500'
  if (status === 'down') return 'border-red-500'
  return 'border-amber-400'
}

function formatDateTime(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
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

    <!-- Loading State -->
    <div v-if="serviceStore.loading.Index" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <UiSkeleton
        v-for="i in serviceStore.indexData.pagination.page_size"
        :key="i"
        variant="rect"
        width="w-full"
        height="h-52"
        rounded
      />
    </div>

    <!-- Empty State -->
    <UiEmptyState
      v-else-if="serviceStore.indexData.items.length === 0"
      :icon="PhHeartbeat"
      title="Belum ada Service"
      description="Silakan daftarkan upstream service baru untuk mulai memantau statusnya."
    >
      <UiButton v-if="hasAllPermissions(['service.create'])" size="lg" @click="openCreate">
        <template #icon>
          <PhPlus class="w-5 h-5" />
        </template>
        Tambah Service Pertama
      </UiButton>
    </UiEmptyState>

    <!-- Data Grid -->
    <template v-else>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <UiCard
          v-for="service in serviceStore.indexData.items"
          :key="service.id"
          :padded="false"
          :classes="{
            wrapper: 'h-full',
            card: `group hover:shadow-md transition-all duration-200 h-full rounded-t-2xl border-t-4 ${healthBorderClass(
              service.health_status,
            )}`,
            body: 'flex flex-col h-full',
          }"
        >
          <div class="p-4 flex-1 flex flex-col">
            <!-- Name + Health Status Pill -->
            <div class="flex items-start justify-between gap-2">
              <h3 class="text-sm font-semibold text-gray-900 truncate min-w-0">
                {{ service.name }}
              </h3>
              <span
                class="shrink-0 inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium"
                :class="healthBannerClass(service.health_status)"
                :title="
                  service.health_checked_at
                    ? `Terakhir dicek: ${formatDateTime(service.health_checked_at)}`
                    : 'Belum pernah dicek'
                "
              >
                <span class="relative flex w-2 h-2">
                  <span
                    v-if="service.health_status === 'up'"
                    class="animate-ping absolute inline-flex w-full h-full rounded-full opacity-75"
                    :class="healthDotClass(service.health_status)"
                  ></span>
                  <span
                    class="relative inline-flex rounded-full w-2 h-2"
                    :class="healthDotClass(service.health_status)"
                  ></span>
                </span>
                {{ healthLabel(service.health_status) }}
              </span>
            </div>

            <!-- Divider -->
            <div class="mt-2.5 mb-2 border-t border-gray-100"></div>

            <!-- Meta Grid -->
            <div class="grid grid-cols-2 gap-y-2 gap-x-2 text-xs">
              <div>
                <p class="text-gray-400">Protocol</p>
                <p class="font-medium text-gray-700 uppercase">{{ service.protocol }}</p>
              </div>
              <div>
                <p class="text-gray-400">Rate Limit</p>
                <p class="font-medium text-gray-700">
                  {{ service.rate_limit_per_minute ?? 'Default' }}
                </p>
              </div>
              <div>
                <p class="text-gray-400">Routes</p>
                <p class="font-medium text-gray-700">{{ service.route_count }}</p>
              </div>
              <div>
                <p class="text-gray-400">Status</p>
                <p
                  class="font-medium"
                  :class="service.is_active ? 'text-emerald-600' : 'text-red-600'"
                >
                  {{ service.is_active ? 'Aktif' : 'Nonaktif' }}
                </p>
              </div>
            </div>

            <!-- URL & Last Checked (same label/value style as meta grid) -->
            <div class="mt-auto pt-2.5 grid grid-cols-2 gap-y-2 gap-x-2 text-xs">
              <div class="min-w-0">
                <p class="text-gray-400">URL</p>
                <p class="font-medium text-gray-700 truncate flex items-center gap-1">
                  <PhLink class="w-3 h-3 shrink-0" />
                  <span class="truncate">{{ service.base_url }}</span>
                </p>
              </div>
              <div class="min-w-0">
                <p class="text-gray-400">Base Path</p>
                <p class="font-medium text-gray-700 truncate">
                  <code class="truncate">{{ service.base_path }}</code>
                </p>
              </div>
              <div class="min-w-0 col-span-2">
                <p class="text-gray-400">Terakhir Dicek</p>
                <p class="font-medium text-gray-700 truncate flex items-center gap-1">
                  <PhClockCounterClockwise class="w-3 h-3 shrink-0" />
                  <span class="truncate">{{
                    service.health_checked_at
                      ? formatDateTime(service.health_checked_at)
                      : 'Belum pernah dicek'
                  }}</span>
                </p>
              </div>
            </div>
          </div>

          <!-- Footer Actions -->
          <div
            class="border-t border-gray-100 bg-gray-50/80 px-3 py-2 flex items-center justify-end gap-0.5 flex-wrap"
          >
            <button
              v-if="hasAllPermissions(['service.health-check'])"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-emerald-600 hover:bg-emerald-100 rounded-md transition-colors disabled:opacity-50"
              :disabled="healthChecking === service.id"
              @click="handleHealthCheck(service.id)"
            >
              <PhHeartbeat
                class="w-3.5 h-3.5"
                :class="{ 'animate-pulse': healthChecking === service.id }"
              />
              Cek
            </button>
            <button
              v-if="hasAllPermissions(['service.edit'])"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-blue-600 hover:bg-blue-100 rounded-md transition-colors"
              @click="openEdit(service)"
            >
              <PhPencil class="w-3.5 h-3.5" />
              Edit
            </button>
            <button
              v-if="hasAllPermissions(['service.delete'])"
              class="flex items-center gap-1 px-2 py-1.5 text-xs font-medium text-gray-500 hover:text-red-600 hover:bg-red-100 rounded-md transition-colors disabled:opacity-50"
              :disabled="serviceStore.loading.Delete"
              @click="handleDelete(service.id)"
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
          :page="serviceStore.indexData.pagination.page"
          :total-pages="serviceStore.indexData.pagination.total_pages"
          @update:page="handlePageChange"
        />
      </div>
    </template>
  </div>

  <FormModal ref="formModalRef" />
</template>
