<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  color?: 'primary' | 'danger' | 'info' | 'warning'
  size?: 'default' | 'sm' | 'lg'
  outline?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  size: 'default',
  outline: false,
})

const colorClasses = computed(() => {
  const baseClasses = {
    primary: {
      solid: 'bg-primary text-primary-foreground hover:bg-primary/80 focus:ring-primary',
      outline: 'border border-primary text-primary hover:bg-primary/10 focus:ring-primary',
    },
    danger: {
      solid: 'bg-red-500 text-white hover:bg-red-500/80 focus:ring-red-500',
      outline: 'border border-red-500 text-red-500 hover:bg-red-500/10 focus:ring-red-500',
    },
    info: {
      solid: 'bg-blue-500 text-white hover:bg-blue-500/80 focus:ring-blue-500',
      outline: 'border border-blue-500 text-blue-500 hover:bg-blue-500/10 focus:ring-blue-500',
    },
    warning: {
      solid: 'bg-yellow-500 text-white hover:bg-yellow-500/80 focus:ring-yellow-500',
      outline:
        'border border-yellow-500 text-yellow-500 hover:bg-yellow-500/10 focus:ring-yellow-500',
    },
  }

  const style = props.outline ? 'outline' : 'solid'
  return baseClasses[props.color][style]
})

const sizeClasses = computed(
  () =>
    ({
      default: 'px-2.5 py-0.5 text-xs',
      sm: 'px-1.5 py-0.5 text-xs',
      lg: 'px-3 py-1 text-sm',
    })[props.size],
)
</script>

<template>
  <div
    :class="[
      'inline-flex items-center rounded-full font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
      colorClasses,
      sizeClasses,
    ]"
  >
    <slot />
  </div>
</template>
