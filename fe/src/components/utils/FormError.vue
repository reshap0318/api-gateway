<script setup lang="ts">
import { computed } from 'vue'
import { useFormError } from '@/composables/useFormError'
import { translateVuelidateError, type ValidationLike } from '@/helpers/vuelidate'

interface FormErrorClasses {
  error?: string
}

const props = defineProps<{
  name?: string
  validation?: ValidationLike
  classes?: FormErrorClasses
}>()

const formErrorStore = useFormError()

const hasVuelidateError = computed(() => props.validation?.$error ?? false)
const vuelidateMessage = computed(() => {
  if (!hasVuelidateError.value || !props.validation) return ''
  return translateVuelidateError(props.validation)
})

const hasServerError = computed(() => (props.name ? formErrorStore.has(props.name) : false))
const serverMessage = computed(() => {
  if (!hasServerError.value || !props.name) return ''
  return formErrorStore.get(props.name).join(', ')
})

const errorMessage = computed(() => {
  if (vuelidateMessage.value) return vuelidateMessage.value
  return serverMessage.value
})

const hasError = computed(() => hasVuelidateError.value || hasServerError.value)
</script>

<template>
  <p v-if="hasError && errorMessage" :class="['mt-1 text-sm text-red-500', props.classes?.error]">
    {{ errorMessage }}
  </p>
</template>
