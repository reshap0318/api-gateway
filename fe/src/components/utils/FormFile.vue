<script setup lang="ts">
import { PhUploadSimple, PhX, PhFile } from '@phosphor-icons/vue'
import FormError from './FormError.vue'
import type { FormFileProps, TFileItem } from './types'

import { ref, computed } from 'vue'
import type { ValidationLike } from '@/helpers/vuelidate'

const props = withDefaults(defineProps<FormFileProps>(), {
  name: '',
  label: '',
  placeholder: 'Drag & drop files here or click to select',
  validation: undefined,
  accept: '*',
  multiple: false,
  maxSize: 0,
  disabled: false,
  classes: () => ({}),
})

const emit = defineEmits<{
  'update:modelValue': [value: File[] | null]
}>()

const files = ref<TFileItem[]>([])
const isDragOver = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const errorMessage = ref('')

const hasError = computed(() => (props.validation as ValidationLike)?.$error ?? false)

function generateId() {
  return Math.random().toString(36).substring(2, 9)
}

function createPreview(file: File): string | undefined {
  if (!file.type.startsWith('image/')) return undefined
  return URL.createObjectURL(file)
}

function addFile(file: File) {
  if (props.maxSize > 0 && file.size > props.maxSize * 1024 * 1024) {
    errorMessage.value = `File ${file.name} terlalu besar (max ${props.maxSize} MB)`
    return
  }

  errorMessage.value = ''
  const item: TFileItem = {
    id: generateId(),
    file,
    preview: createPreview(file),
  }
  files.value.push(item)
  emitUpdate()
}

function removeFile(id: string) {
  const index = files.value.findIndex((f) => f.id === id)
  if (index > -1) {
    const item = files.value[index]
    if (item.preview) {
      URL.revokeObjectURL(item.preview)
    }
    files.value.splice(index, 1)
    emitUpdate()
  }
}

function emitUpdate() {
  emit('update:modelValue', files.value.length > 0 ? files.value.map((f) => f.file) : null)
}

function handleDrop(event: DragEvent) {
  event.preventDefault()
  isDragOver.value = false

  if (props.disabled) return

  const droppedFiles = event.dataTransfer?.files
  if (!droppedFiles) return

  const newFiles = Array.from(droppedFiles)
  if (!props.multiple && newFiles.length > 0) {
    files.value = []
  }
  newFiles.forEach(addFile)
}

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const selectedFiles = input.files
  if (!selectedFiles) return

  const newFiles = Array.from(selectedFiles)
  if (!props.multiple && newFiles.length > 0) {
    files.value = []
  }
  newFiles.forEach(addFile)
  input.value = ''
}

function handleClick() {
  if (props.disabled) return
  fileInput.value?.click()
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
  isDragOver.value = true
}

function handleDragLeave() {
  isDragOver.value = false
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>

<template>
  <div :class="['w-full', props.classes.wrapper]">
    <label
      v-if="props.label"
      :class="['mb-1 block text-sm font-medium text-gray-700', props.classes.label]"
    >
      {{ props.label }}
    </label>

    <div
      :class="[
        'relative rounded-md border-2 border-dashed p-6 text-center transition cursor-pointer',
        isDragOver
          ? 'border-blue-500 bg-blue-50'
          : hasError
            ? 'border-red-500 hover:border-red-400'
            : 'border-gray-300 hover:border-gray-400',
        props.disabled && 'opacity-50 cursor-not-allowed pointer-events-none',
        props.classes.dropZone,
      ]"
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @click="handleClick"
    >
      <input
        ref="fileInput"
        type="file"
        :accept="props.accept"
        :multiple="props.multiple"
        :disabled="props.disabled"
        class="hidden"
        @change="handleFileSelect"
      />

      <PhUploadSimple :size="32" class="mx-auto text-gray-400" />
      <p class="mt-2 text-sm text-gray-600">{{ props.placeholder }}</p>
    </div>

    <!-- File List -->
    <div v-if="files.length > 0" :class="['mt-2 space-y-2', props.classes.fileItem]">
      <div
        v-for="item in files"
        :key="item.id"
        class="flex items-center gap-3 rounded-md border border-gray-200 bg-gray-50 p-2"
      >
        <div v-if="item.preview" class="h-10 w-10 shrink-0 overflow-hidden rounded">
          <img :src="item.preview" class="h-full w-full object-cover" />
        </div>
        <div v-else class="flex h-10 w-10 shrink-0 items-center justify-center rounded bg-gray-200">
          <PhFile :size="20" class="text-gray-500" />
        </div>

        <div class="flex-1 truncate">
          <p class="truncate text-sm font-medium text-gray-700">{{ item.file.name }}</p>
          <p class="text-xs text-gray-500">{{ formatFileSize(item.file.size) }}</p>
        </div>

        <button
          type="button"
          class="shrink-0 text-gray-400 hover:text-red-500 transition"
          @click.stop="removeFile(item.id)"
        >
          <PhX :size="20" />
        </button>
      </div>
    </div>

    <!-- Error Message -->
    <FormError
      :name="props.name || undefined"
      :validation="props.validation as ValidationLike"
      :classes="{ error: props.classes.error }"
    />
    <p v-if="errorMessage" class="mt-1 text-sm text-red-500">
      {{ errorMessage }}
    </p>
  </div>
</template>
