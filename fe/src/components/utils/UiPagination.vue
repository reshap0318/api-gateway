<script setup lang="ts">
import { computed } from 'vue'
import { PhCaretLeft, PhCaretRight } from '@phosphor-icons/vue'
import type { UiPaginationProps } from './types'

const props = withDefaults(defineProps<UiPaginationProps>(), {
  maxVisible: 5,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:page': [page: number]
}>()

const pages = computed(() => {
  const total = props.totalPages
  const current = props.page
  const max = props.maxVisible

  if (total <= max) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const half = Math.floor(max / 2)
  let start = current - half
  let end = current + half

  if (start < 1) {
    start = 1
    end = max
  }

  if (end > total) {
    end = total
    start = total - max + 1
  }

  const result: (number | '...')[] = []

  if (start > 1) {
    result.push(1)
    if (start > 2) result.push('...')
  }

  for (let i = start; i <= end; i++) {
    result.push(i)
  }

  if (end < total) {
    if (end < total - 1) result.push('...')
    result.push(total)
  }

  return result
})

function goTo(p: number) {
  if (p >= 1 && p <= props.totalPages && p !== props.page) {
    emit('update:page', p)
  }
}
</script>

<template>
  <nav
    v-if="props.totalPages > 1"
    :class="['flex items-center justify-center gap-1', props.classes.wrapper]"
    aria-label="Pagination"
  >
    <!-- Previous -->
    <button
      :class="[
        'inline-flex items-center justify-center w-9 h-9 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors disabled:opacity-40 disabled:cursor-not-allowed',
        props.classes.button,
      ]"
      :disabled="props.page <= 1"
      @click="goTo(props.page - 1)"
    >
      <PhCaretLeft class="w-4 h-4" />
    </button>

    <!-- Page numbers -->
    <template v-for="(p, idx) in pages" :key="idx">
      <span
        v-if="p === '...'"
        :class="[
          'inline-flex items-center justify-center w-9 h-9 text-sm text-gray-400 select-none',
          props.classes.ellipsis,
        ]"
      >
        …
      </span>
      <button
        v-else
        :class="[
          'inline-flex items-center justify-center w-9 h-9 rounded-lg text-sm font-medium transition-colors',
          p === props.page
            ? ['bg-blue-600 text-white shadow-sm', props.classes.buttonActive]
            : ['text-gray-600 hover:bg-gray-100 hover:text-gray-900', props.classes.button],
        ]"
        @click="goTo(p as number)"
      >
        {{ p }}
      </button>
    </template>

    <!-- Next -->
    <button
      :class="[
        'inline-flex items-center justify-center w-9 h-9 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors disabled:opacity-40 disabled:cursor-not-allowed',
        props.classes.button,
      ]"
      :disabled="props.page >= props.totalPages"
      @click="goTo(props.page + 1)"
    >
      <PhCaretRight class="w-4 h-4" />
    </button>
  </nav>
</template>
