<template>
  <div class="w-full overflow-x-auto rounded-xl border border-gray-200 bg-white shadow-sm">
    <table class="w-full text-left text-sm text-gray-900">
      <thead class="bg-slate-50 text-gray-900 border-b border-gray-200">
        <tr>
          <th
            v-for="(col, idx) in columns"
            :key="idx"
            class="px-4 py-3 font-semibold whitespace-nowrap"
            :class="[col.class, col.headerClass]"
          >
            <div
              class="flex items-center gap-1"
              :class="[col.sortable ? 'cursor-pointer hover:text-primary-800' : '']"
              @click="col.sortable ? handleSort(col.data) : null"
            >
              <slot :name="`header-${String(col.data)}`" :column="col">
                {{ col.title }}
              </slot>
              <template v-if="col.sortable">
                <svg
                  v-if="activeSorts[col.data] === 'asc'"
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4 text-primary-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 15l7-7 7 7" />
                </svg>
                <svg
                  v-else-if="activeSorts[col.data] === 'desc'"
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4 text-primary-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
                <svg
                  v-else
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4 text-gray-400 hover:text-primary-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                  />
                </svg>
              </template>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="isLoading">
          <td :colspan="columns.length" class="whitespace-nowrap px-4 py-6 text-center">
            <div class="mx-3">
              <UiSkeleton v-for="i in 2" :key="i" variant="text" />
            </div>
          </td>
        </tr>
        <tr v-else-if="!datas || datas.length === 0">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-500">
            <slot name="empty">Tidak ada data ditampilkan.</slot>
          </td>
        </tr>
        <template v-else>
          <template v-for="(row, rowIndex) in datas" :key="rowIndex">
            <tr
              class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
              :class="[row.class, expand ? 'cursor-pointer' : '']"
              @click="expandToggle(rowIndex)"
            >
              <td
                v-for="col in columns"
                :key="String(col.data)"
                class="px-4 py-4"
                :class="col.class"
              >
                <slot
                  :name="String(col.data)"
                  :value="getValue(row[col.data])"
                  :item="row"
                  :index="rowIndex"
                >
                  {{ getValue(row[col.data]) }}
                </slot>
              </td>
            </tr>
            <tr v-if="expand && expandRow[rowIndex]" class="border-b border-gray-100 bg-gray-50/50">
              <td :colspan="columns.length" class="p-4">
                <slot name="expandedCol" :value="row" :index="rowIndex"></slot>
              </td>
            </tr>
          </template>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import UiSkeleton from '@/components/utils/UiSkeleton.vue'
import { ref, computed } from 'vue'
import type { TTableColumn, TTableRow, TSortOrder } from './types'

const props = withDefaults(
  defineProps<{
    columns: TTableColumn[]
    datas?: TTableRow[]
    sorts?: Record<string, TSortOrder>
    isLoading?: boolean
    expand?: boolean
  }>(),
  {
    datas: () => [],
    sorts: () => ({}),
    isLoading: false,
    expand: false,
  },
)

const localSorts = ref<Record<string, TSortOrder>>({})

const activeSorts = computed(() => {
  return props.sorts !== undefined ? props.sorts : localSorts.value
})

const emit = defineEmits<{
  sort: [key: string, order: TSortOrder]
}>()

const expandRow = ref<Record<number, boolean>>({})

function getValue(input: unknown): unknown {
  if (input === null || input === undefined) return ''
  if (typeof input === 'number') return input.toString()
  if (typeof input === 'string') return input.trim()
  if (input) return input
  return ''
}

function expandToggle(idx: number) {
  if (props.expand) {
    expandRow.value[idx] = !expandRow.value[idx]
  }
}

function handleSort(key: string) {
  const currentSort = activeSorts.value[key]
  let nextSort: TSortOrder = 'asc'
  if (currentSort === 'asc') nextSort = 'desc'
  else if (currentSort === 'desc') nextSort = 'none'

  if (props.sorts === undefined) {
    if (nextSort === 'none') {
      delete localSorts.value[key]
    } else {
      localSorts.value[key] = nextSort
    }
  }

  emit('sort', key, nextSort)
}
</script>
