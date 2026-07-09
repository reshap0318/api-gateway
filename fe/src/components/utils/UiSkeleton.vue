<script setup lang="ts">
import type { UiSkeletonProps } from './types'

const props = withDefaults(defineProps<UiSkeletonProps>(), {
  variant: 'text',
  width: 'w-full',
  height: 'h-4',
  rounded: false,
  count: 1,
  classes: () => ({}),
})
</script>

<template>
  <div :class="['animate-pulse', props.classes.wrapper]">
    <template v-if="props.variant === 'circle'">
      <div
        v-for="i in props.count"
        :key="i"
        :class="['bg-gray-200 rounded-full', props.height, props.width, props.classes.bar]"
      />
    </template>

    <template v-else-if="props.variant === 'rect'">
      <div
        v-for="i in props.count"
        :key="i"
        :class="[
          'bg-gray-200',
          props.height,
          props.width,
          props.rounded && 'rounded-lg',
          props.classes.bar,
        ]"
      />
    </template>

    <template v-else-if="props.variant === 'card'">
      <div
        v-for="i in props.count"
        :key="i"
        :class="['bg-white rounded-lg shadow-md p-4 space-y-3', props.classes.bar]"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-gray-200 shrink-0" />
          <div class="space-y-2 flex-1">
            <div class="h-4 bg-gray-200 rounded w-1/3" />
            <div class="h-3 bg-gray-200 rounded w-2/3" />
          </div>
        </div>
        <div class="space-y-2">
          <div class="h-3 bg-gray-200 rounded w-full" />
          <div class="h-3 bg-gray-200 rounded w-5/6" />
        </div>
      </div>
    </template>

    <template v-else>
      <div
        v-for="i in props.count"
        :key="i"
        :class="[
          'bg-gray-200 rounded',
          props.height,
          props.width,
          props.classes.bar,
          i < props.count && 'mb-2',
        ]"
      />
    </template>
  </div>
</template>
