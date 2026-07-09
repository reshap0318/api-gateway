<script setup lang="ts">
import type { UiEmptyStateProps } from './types'

const props = withDefaults(defineProps<UiEmptyStateProps>(), {
  icon: undefined,
  description: '',
  variant: 'default',
  classes: () => ({}),
})
</script>

<template>
  <div
    :class="[
      'text-center bg-white rounded-2xl border-2 border-dashed border-gray-300 shadow-sm',
      props.classes.wrapper,
      props.variant === 'compact' ? 'py-12' : 'py-24',
    ]"
  >
    <div
      v-if="props.icon"
      :class="[
        'inline-flex items-center justify-center rounded-full mb-4',
        props.classes.icon,
        props.variant === 'compact' ? 'w-16 h-16' : 'w-20 h-20',
        'bg-indigo-100 text-indigo-600',
      ]"
    >
      <component :is="props.icon" :class="props.variant === 'compact' ? 'w-8 h-8' : 'w-10 h-10'" />
    </div>
    <h3
      :class="[
        'font-semibold text-gray-900 mb-2',
        props.classes.title,
        props.variant === 'compact' ? 'text-xl' : 'text-2xl',
      ]"
    >
      {{ props.title }}
    </h3>
    <p v-if="props.description" :class="['text-gray-600 mb-6', props.classes.description]">
      {{ props.description }}
    </p>
    <slot />
  </div>
</template>
