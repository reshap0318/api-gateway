<script setup lang="ts">
import { PhX } from '@phosphor-icons/vue'
import type { TModalSize, UiModalProps } from './types'

const props = withDefaults(defineProps<UiModalProps>(), {
  title: '',
  size: 'md',
  persistent: false,
  classes: () => ({}),
})

const sizeClasses: Record<TModalSize, string> = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  '3xl': 'max-w-3xl',
  '4xl': 'max-w-4xl',
  '5xl': 'max-w-5xl',
  '6xl': 'max-w-6xl',
  '7xl': 'max-w-7xl',
  full: 'max-w-full',
}

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
}>()

function close() {
  emit('update:modelValue', false)
  emit('close')
}

function closeFromBackdrop() {
  if (props.persistent) return
  close()
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="props.modelValue" class="fixed inset-0 z-50 overflow-y-auto">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/50 transition-opacity" @click="closeFromBackdrop" />

        <!-- Modal Panel -->
        <div class="flex min-h-full items-center justify-center p-4 text-center sm:p-0">
          <Transition
            enter-active-class="transition duration-200 ease-out"
            enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to-class="opacity-100 translate-y-0 sm:scale-100"
            leave-active-class="transition duration-150 ease-in"
            leave-from-class="opacity-100 translate-y-0 sm:scale-100"
            leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <div
              v-if="props.modelValue"
              :class="[
                'relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 w-full',
                sizeClasses[props.size],
                props.classes.container,
              ]"
            >
              <!-- Header -->
              <div
                :class="[
                  'flex items-center justify-between border-b border-gray-200 px-6 py-4',
                  props.classes.header,
                ]"
              >
                <slot name="header">
                  <h3 class="text-lg font-semibold text-gray-900">{{ props.title }}</h3>
                </slot>
                <button
                  type="button"
                  class="text-gray-400 hover:text-gray-500 transition-colors"
                  @click="close"
                >
                  <PhX class="h-5 w-5" />
                </button>
              </div>

              <!-- Body -->
              <div :class="['px-6 py-4', props.classes.body]">
                <slot />
              </div>

              <!-- Footer -->
              <div
                v-if="$slots.footer"
                :class="[
                  'border-t border-gray-200 bg-gray-50 px-6 py-4 flex justify-end gap-3',
                  props.classes.footer,
                ]"
              >
                <slot name="footer" />
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
