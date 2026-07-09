<script setup lang="ts">
import type { UiCardProps } from './types'

const props = withDefaults(defineProps<UiCardProps>(), {
  padded: true,
  classes: () => ({}),
})

defineSlots<{
  header?: () => any
  default?: () => any
  footer?: () => any
}>()
</script>

<template>
  <div :class="['w-full', props.classes.wrapper]">
    <div :class="['overflow-hidden rounded-lg bg-white shadow-md', props.classes.card]">
      <!-- Header -->
      <div v-if="$slots.header" :class="[props.classes.header]">
        <slot name="header" />
      </div>

      <!-- Body -->
      <div :class="[props.padded && 'p-6', props.classes.body]">
        <slot />
      </div>

      <!-- Footer -->
      <div v-if="$slots.footer" :class="['bg-gray-50 px-6 py-4', props.classes.footer]">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>
