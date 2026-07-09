<script setup lang="ts">
import FormError from './FormError.vue'
import { PhEye, PhEyeSlash } from '@phosphor-icons/vue'
import type { FormPasswordProps } from './types'

import { ref, computed } from 'vue'
import type { ValidationLike } from '@/helpers/vuelidate'

const props = withDefaults(defineProps<FormPasswordProps>(), {
  name: '',
  label: '',
  placeholder: '',
  validation: undefined,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const showPassword = ref(false)

const inputType = computed(() => (showPassword.value ? 'text' : 'password'))

const hasError = computed(() => (props.validation as ValidationLike)?.$error ?? false)

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
      <input
        :id="props.label"
        :type="inputType"
        :value="props.modelValue"
        :placeholder="props.placeholder"
        :class="[
          'w-full rounded-md border px-3 py-2 outline-none transition',
          'border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500',
          hasError && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          'pr-10',
          props.classes.input,
        ]"
        @input="onInput"
      />

      <button
        type="button"
        class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500 hover:text-gray-700"
        @click="showPassword = !showPassword"
      >
        <PhEye v-if="!showPassword" :size="20" weight="regular" />
        <PhEyeSlash v-else :size="20" weight="regular" />
      </button>
    </div>

    <FormError
      :name="props.name || undefined"
      :validation="props.validation as ValidationLike"
      :classes="{ error: props.classes.error }"
    />
  </div>
</template>
