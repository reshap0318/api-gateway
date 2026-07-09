<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import FormError from './FormError.vue'
import type { FormSelectProps, TSelectOption } from './types'

import { computed } from 'vue'
import type { ValidationLike } from '@/helpers/vuelidate'

const props = withDefaults(defineProps<FormSelectProps>(), {
  name: '',
  label: '',
  options: () => [],
  placeholder: 'Select...',
  validation: undefined,
  searchable: false,
  mode: 'single',
  closable: true,
  disabled: false,
  loading: false,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:modelValue': [value: any]
  change: [value: any]
}>()

const internalValue = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
    emit('change', value)
  },
})
</script>

<template>
  <div :class="['w-full', props.classes.wrapper]">
    <label
      v-if="props.label"
      :class="['mb-1 block text-sm font-medium text-gray-700', props.classes.label]"
    >
      {{ props.label }}
    </label>

    <Multiselect
      v-model="internalValue"
      :options="props.options as TSelectOption[]"
      :placeholder="props.placeholder"
      :searchable="props.searchable"
      :mode="props.mode"
      :closable="props.closable"
      :disabled="props.disabled"
      :loading="props.loading"
      :caret="true"
      :append-to-body="true"
      :value-prop="'value'"
      :track-by="'label'"
      :label="'label'"
      :classes="{
        container: 'multiselect',
        containerDisabled: 'is-disabled',
        containerOpen: 'is-open',
        containerOpenTop: 'is-open-top',
        containerActive: 'is-active',
        caret: 'multiselect-caret',
        caretOpen: 'is-open',
        dropdown: 'multiselect-dropdown',
        dropdownTop: 'is-top',
        dropdownHidden: 'is-hidden',
        options: 'multiselect-options',
        option: 'multiselect-option',
        optionPointed: 'is-pointed',
        optionSelected: 'is-selected',
        optionDisabled: 'is-disabled',
        singleLabel: 'multiselect-single-label',
        multipleLabel: 'multiselect-multiple-label',
        placeholder: 'multiselect-placeholder',
        tags: 'multiselect-tags',
        tag: 'multiselect-tag',
        tagDisabled: 'is-disabled',
        tagRemove: 'multiselect-tag-remove items-left',
        tagRemoveIcon: 'multiselect-tag-remove-icon',
        search: 'multiselect-search',
        noOptions: 'multiselect-no-options',
        noResults: 'multiselect-no-results',
        group: 'multiselect-group',
        groupLabel: 'multiselect-group-label',
        groupOptions: 'multiselect-group-options',
        spinner: 'multiselect-spinner',
        clear: 'multiselect-clear',
        clearIcon: 'multiselect-clear-icon',
      }"
    />

    <FormError
      :name="props.name || undefined"
      :validation="props.validation as ValidationLike"
      :classes="{ error: props.classes.error }"
    />
  </div>
</template>

<style src="@vueform/multiselect/themes/default.css"></style>

<style scoped>
/* Container overrides */
:deep(.multiselect) {
  --ms-border-color: #d1d5db;
  --ms-border-color-active: #3b82f6;
  --ms-ring-color: rgba(59, 130, 246, 0.3);
  --ms-caret-color: #9ca3af;
  --ms-placeholder-color: #9ca3af;
  --ms-option-bg-selected: #2563eb;
  --ms-tag-bg: #dbeafe;
  --ms-tag-color: #1d4ed8;
  border-radius: 0.375rem;
  min-height: 42px;
}

:deep(.multiselect:hover) {
  --ms-border-color: #9ca3af;
}

:deep(.multiselect.is-active) {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

:deep(.multiselect.is-disabled) {
  background: #f3f4f6;
  cursor: not-allowed;
  opacity: 0.5;
}

/* Error state */
:deep(.multiselect.is-error) {
  --ms-border-color: #ef4444;
  --ms-border-color-active: #ef4444;
}

/* Placeholder */
:deep(.multiselect-placeholder) {
  color: #9ca3af;
  font-size: 0.875rem;
  padding: 0.5rem 0.75rem;
}

/* Label */
:deep(.multiselect-single-label) {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: #111827;
}

/* Caret */
:deep(.multiselect-caret) {
  background-color: #9ca3af;
  margin-right: 0.875rem;
}

:deep(.multiselect-caret.is-open) {
  transform: rotate(180deg);
}

/* Dropdown */
:deep(.multiselect-dropdown) {
  --ms-dropdown-border-color: #d1d5db;
  --ms-dropdown-radius: 0.375rem;
  z-index: 1000;
  margin-top: 0.25rem;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

/* Options */
:deep(.multiselect-option) {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  cursor: pointer;
  color: #374151;
}

:deep(.multiselect-option.is-pointed) {
  background: #eff6ff;
  color: #1f2937;
}

:deep(.multiselect-option.is-selected) {
  background: #2563eb;
  color: #fff;
}

:deep(.multiselect-option.is-pointed.is-selected) {
  background: #2563eb;
  color: #fff;
}

:deep(.multiselect-option.is-disabled) {
  background: #fff;
  color: #9ca3af;
  cursor: not-allowed;
}

/* Tags */
:deep(.multiselect-tags) {
  padding-left: 0.5rem;
  margin-top: 0.25rem;
}

:deep(.multiselect-tag) {
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.125rem 0.5rem;
  border-radius: 0.25rem;
}

/* Search */
:deep(.multiselect-search) {
  font-size: 0.875rem;
}

:deep(.multiselect-search:focus) {
  outline: none;
}

/* No options/results */
:deep(.multiselect-no-options),
:deep(.multiselect-no-results) {
  color: #6b7280;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
}

/* Clear button */
:deep(.multiselect-clear-icon) {
  background-color: #9ca3af;
}

:deep(.multiselect-clear:hover .multiselect-clear-icon) {
  background-color: #4b5563;
}
</style>
