<script setup lang="ts">
import { UiModal, FormInput, FormSelect, UiButton } from '@/components/utils'
import { computed, ref } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useGatewayServiceStore } from '@/stores/gatewayService'
import type { IGatewayService } from '@/stores/gatewayService'
import { useFormError } from '@/composables/useFormError'

const serviceStore = useGatewayServiceStore()
const formError = useFormError()
const v$ = useVuelidate(serviceStore.formRules, serviceStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!serviceStore.form.id)

const protocolOptions = [
  { value: 'http', label: 'HTTP (REST)' },
  { value: 'websocket', label: 'WebSocket' },
]

// FormInput's modelValue is typed as string — bridge to the numeric|null form field.
const rateLimitInput = computed({
  get: () => (serviceStore.form.rate_limit_per_minute ?? '').toString(),
  set: (val: string) => {
    serviceStore.form.rate_limit_per_minute = val === '' ? null : Number(val)
  },
})

function show(data?: IGatewayService) {
  if (data) {
    serviceStore.form.id = data.id
    serviceStore.form.name = data.name
    serviceStore.form.base_url = data.base_url
    serviceStore.form.protocol = data.protocol
    serviceStore.form.rate_limit_per_minute = data.rate_limit_per_minute
    serviceStore.form.is_active = data.is_active
  } else {
    serviceStore.resetForm()
  }
  v$.value.$reset()
  formError.clear()
  isVisible.value = true
}

function close() {
  isVisible.value = false
}

async function handleSubmit() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    if (isEdit.value && serviceStore.form.id) {
      await serviceStore.update(serviceStore.form.id)
    } else {
      await serviceStore.create()
    }
    close()
  } catch {
    // error already handled by axios interceptor
  }
}

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Service' : 'Tambah Service'"
    size="lg"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="serviceStore.form.name"
          name="name"
          label="Nama Service"
          placeholder="e.g. User Service"
          :validation="v$.name"
        />

        <FormInput
          v-model="serviceStore.form.base_url"
          name="base_url"
          label="Base URL"
          placeholder="http://localhost:9000"
          :validation="v$.base_url"
        />

        <FormSelect
          v-model="serviceStore.form.protocol"
          name="protocol"
          label="Protocol"
          :options="protocolOptions"
          :validation="v$.protocol"
        />

        <FormInput
          v-model="rateLimitInput"
          name="rate_limit_per_minute"
          type="number"
          label="Rate Limit per Menit (opsional)"
          placeholder="Kosongkan untuk pakai default global"
          :validation="v$.rate_limit_per_minute"
        />

        <label class="flex items-center gap-2 cursor-pointer">
          <input
            v-model="serviceStore.form.is_active"
            type="checkbox"
            class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span class="text-sm font-medium text-gray-700">Aktif</span>
        </label>
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="serviceStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="serviceStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
