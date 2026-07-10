<script setup lang="ts">
import { UiTable, UiPagination, UiBadge, FormSelect } from '@/components/utils'
import { ref, computed, onMounted, watch } from 'vue'
import { useGatewayAuditLogStore } from '@/stores/gatewayAuditLog'
import type { IGatewayAuditLog } from '@/stores/gatewayAuditLog'
import { useUserStore } from '@/stores/user'
import type { IUser } from '@/stores/user'
import type { TTableColumn } from '@/components/utils/types'

const auditLogStore = useGatewayAuditLogStore()
const userStore = useUserStore()

const allUsers = ref<IUser[]>([])
const entityTypeFilter = ref('')
const actorFilter = ref('')
const fromDate = ref('')
const toDate = ref('')

const columns: TTableColumn[] = [
  { title: 'Waktu', data: 'created_at', sortable: false },
  { title: 'Entity Type', data: 'entity_type', sortable: false },
  { title: 'Entity ID', data: 'entity_id', sortable: false },
  { title: 'Action', data: 'action', sortable: false },
  { title: 'Actor', data: 'actor_name', sortable: false },
]

const rows = computed(() => auditLogStore.indexData.items as unknown as Record<string, unknown>[])

const entityTypeOptions = [
  { value: '', label: 'Semua Entity Type' },
  { value: 'service', label: 'Service' },
  { value: 'route', label: 'Route' },
]

const actorOptions = computed(() => [
  { value: '', label: 'Semua Actor' },
  ...allUsers.value.map((u) => ({ value: String(u.id), label: u.name })),
])

function actionBadgeColor(action: string): 'primary' | 'danger' | 'warning' {
  if (action === 'create') return 'primary'
  if (action === 'delete') return 'danger'
  return 'warning'
}

function formatDateTime(value: string): string {
  return new Date(value).toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatChanges(changes: string): string {
  try {
    return JSON.stringify(JSON.parse(changes), null, 2)
  } catch {
    return changes
  }
}

function fetchLogs(page?: number) {
  auditLogStore.fetchAll(page, {
    entity_type: entityTypeFilter.value || undefined,
    actor: actorFilter.value || undefined,
    from: fromDate.value || undefined,
    to: toDate.value || undefined,
  })
}

function handlePageChange(page: number) {
  fetchLogs(page)
}

watch([entityTypeFilter, actorFilter, fromDate, toDate], () => fetchLogs(1))

onMounted(async () => {
  allUsers.value = await userStore.fetchAllUsers()
  fetchLogs()
})
</script>

<template>
  <div class="mx-auto px-4">
    <!-- Header Section -->
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Audit Log</h1>
      <p class="hidden sm:block text-sm text-gray-600 mt-1">
        Histori perubahan konfigurasi Service &amp; Route (create/update/delete).
      </p>
    </div>

    <!-- Filters -->
    <div class="mb-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <FormSelect
        v-model="entityTypeFilter"
        name="entity_type_filter"
        label="Entity Type"
        :options="entityTypeOptions"
      />
      <FormSelect v-model="actorFilter" name="actor_filter" label="Actor" :options="actorOptions" />
      <div>
        <label class="mb-1 block text-sm font-medium text-gray-700">Dari Tanggal</label>
        <input
          v-model="fromDate"
          type="date"
          class="h-[42px] w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="mb-1 block text-sm font-medium text-gray-700">Sampai Tanggal</label>
        <input
          v-model="toDate"
          type="date"
          class="h-[42px] w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
        />
      </div>
    </div>

    <UiTable :columns="columns" :datas="rows" :is-loading="auditLogStore.loading.Index" expand>
      <template #created_at="{ item }">
        <span class="text-gray-600 whitespace-nowrap">
          {{ formatDateTime((item as unknown as IGatewayAuditLog).created_at) }}
        </span>
      </template>

      <template #entity_type="{ item }">
        <UiBadge color="info">{{ (item as unknown as IGatewayAuditLog).entity_type }}</UiBadge>
      </template>

      <template #action="{ item }">
        <UiBadge :color="actionBadgeColor((item as unknown as IGatewayAuditLog).action)">
          {{ (item as unknown as IGatewayAuditLog).action }}
        </UiBadge>
      </template>

      <template #actor_name="{ item }">
        <span class="text-gray-700">
          {{ (item as unknown as IGatewayAuditLog).actor_name || '-' }}
        </span>
      </template>

      <template #expandedCol="{ value }">
        <pre class="text-xs bg-gray-50 border border-gray-200 rounded-md p-3 overflow-x-auto">{{
          formatChanges((value as unknown as IGatewayAuditLog).changes)
        }}</pre>
      </template>

      <template #empty>
        <p class="text-gray-500">Belum ada aktivitas tercatat.</p>
      </template>
    </UiTable>
    <p class="mt-2 text-xs text-gray-400">Klik baris untuk melihat detail perubahan (JSON).</p>

    <!-- Pagination -->
    <div class="mt-6 flex justify-center">
      <UiPagination
        :page="auditLogStore.indexData.pagination.page"
        :total-pages="auditLogStore.indexData.pagination.total_pages"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>
