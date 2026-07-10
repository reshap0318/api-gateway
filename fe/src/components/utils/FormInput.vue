<script setup lang="ts">
import FormError from './FormError.vue'
import type { FormInputProps } from './types'

import { computed } from 'vue'
import type { ValidationLike } from '@/helpers/vuelidate'

const props = withDefaults(defineProps<FormInputProps>(), {
  type: 'text',
  name: '',
  label: '',
  placeholder: '',
  validation: undefined,
  prefixIcon: undefined,
  suffixIcon: undefined,
  iconSize: 20,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const hasError = computed(() => (props.validation as ValidationLike)?.$error ?? false)

const hasPrefix = computed(() => !!props.prefixIcon)
const hasSuffix = computed(() => !!props.suffixIcon)

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>

<template>
  <div :class="['w-full', props.classes.wrapper]">
    <label
      v-if="props.label"
      :for="props.label"
      :class="['mb-1 block text-sm font-medium text-gray-700', props.classes.label]"
    >
      {{ props.label }}
    </label>

    <div class="relative">
      <component
        :is="props.prefixIcon"
        v-if="hasPrefix"
        :size="props.iconSize"
        weight="regular"
        class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
      />

      <input
        :id="props.label"
        :type="props.type"
        :value="props.modelValue"
        :placeholder="props.placeholder"
        :class="[
          'w-full rounded-md border px-3 py-2 outline-none transition bg-white',
          'border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500',
          hasError && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          hasPrefix && 'pl-10',
          hasSuffix && 'pr-10',
          props.classes.input,
        ]"
        @input="onInput"
      />

      <component
        :is="props.suffixIcon"
        v-if="hasSuffix"
        :size="props.iconSize"
        weight="regular"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
      />
    </div>

    <FormError
      :name="props.name || undefined"
      :validation="props.validation as ValidationLike"
      :classes="{ error: props.classes.error }"
    />
  </div>
</template>
