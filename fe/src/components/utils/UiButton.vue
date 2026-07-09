<script setup lang="ts">
import { computed } from 'vue'
import type { UiButtonProps } from './types'

const props = withDefaults(defineProps<UiButtonProps>(), {
  type: 'button',
  disabled: false,
  loading: false,
  variant: 'primary',
  outline: false,
  size: 'md',
  rounded: 'md',
  fullWidth: false,
  loadingText: 'Loading...',
})

const solidClasses: Record<string, string> = {
  primary:
    'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 border-2 border-transparent',
  secondary:
    'bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500 border-2 border-transparent',
  danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 border-2 border-transparent',
  success:
    'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 border-2 border-transparent',
}

const outlineClasses: Record<string, string> = {
  primary:
    'bg-transparent border-2 border-blue-600 text-blue-600 hover:bg-blue-50 focus:ring-blue-500',
  secondary:
    'bg-transparent border-2 border-gray-600 text-gray-600 hover:bg-gray-50 focus:ring-gray-500',
  danger: 'bg-transparent border-2 border-red-600 text-red-600 hover:bg-red-50 focus:ring-red-500',
  success:
    'bg-transparent border-2 border-green-600 text-green-600 hover:bg-green-50 focus:ring-green-500',
}

const sizeClasses: Record<string, string> = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg',
}

const roundedClasses: Record<string, string> = {
  none: 'rounded-none',
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  full: 'rounded-full',
}

const buttonClass = computed(() => {
  const variantClass = props.outline ? outlineClasses[props.variant] : solidClasses[props.variant]

  return [
    'inline-flex items-center justify-center gap-2 font-medium transition focus:outline-none focus:ring-2',
    'disabled:cursor-not-allowed disabled:opacity-50',
    variantClass,
    sizeClasses[props.size],
    roundedClasses[props.rounded],
    props.fullWidth && 'w-full',
  ]
})
</script>

<template>
  <button :type="props.type" :disabled="props.disabled || props.loading" :class="buttonClass">
    <template v-if="!props.loading">
      <slot name="icon" />
      <slot />
    </template>
    <span v-else>{{ props.loadingText }}</span>
  </button>
</template>
