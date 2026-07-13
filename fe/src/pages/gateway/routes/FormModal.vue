<script setup lang="ts">
import { UiModal, FormInput, FormSelect, UiButton } from '@/components/utils'
import { computed, ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useGatewayRouteStore } from '@/stores/gatewayRoute'
import type { IGatewayRoute } from '@/stores/gatewayRoute'
import { useGatewayServiceStore } from '@/stores/gatewayService'
import type { IGatewayService } from '@/stores/gatewayService'
import { usePermissionStore } from '@/stores/permission'
import type { IPermission } from '@/stores/permission'
import { useFormError } from '@/composables/useFormError'

const routeStore = useGatewayRouteStore()
const serviceStore = useGatewayServiceStore()
const permissionStore = usePermissionStore()
const formError = useFormError()
const v$ = useVuelidate(routeStore.formRules, routeStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!routeStore.form.id)

const allServices = ref<IGatewayService[]>([])
const allPermissions = ref<IPermission[]>([])
const optionsLoading = ref(false)

const methodOptions = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', '*'].map((m) => ({
  value: m,
  label: m === '*' ? '* (Semua Method)' : m,
}))

const matchModeOptions = [
  { value: 'any', label: 'Any — cukup salah satu permission' },
  { value: 'all', label: 'All — harus punya semua permission' },
]

const serviceOptions = computed(() =>
  allServices.value.map((s) => ({ value: s.id, label: s.name })),
)

// Preview of the real publicly-reachable URL (base_path + path_pattern), so the admin can
// see the effect of the selected Service's prefix before saving.
const fullPathPreview = computed(() => {
  const service = allServices.value.find((s) => s.id === routeStore.form.service)
  if (!service || !routeStore.form.path_pattern) return null
  return `${service.base_path}${routeStore.form.path_pattern}`
})

const permissionOptions = computed(() =>
  allPermissions.value.map((p) => ({ value: p.id, label: p.name })),
)

// FormInput's modelValue is typed as string — bridge to the numeric|null form field.
const rateLimitInput = computed({
  get: () => (routeStore.form.rate_limit_per_minute ?? '').toString(),
  set: (val: string) => {
    routeStore.form.rate_limit_per_minute = val === '' ? null : Number(val)
  },
})

async function loadOptions() {
  if (allServices.value.length > 0 && allPermissions.value.length > 0) return
  optionsLoading.value = true
  try {
    const [services, permissions] = await Promise.all([
      serviceStore.fetchAllServices(),
      permissionStore.fetchAllPermissions(),
    ])
    allServices.value = services
    allPermissions.value = permissions
  } finally {
    optionsLoading.value = false
  }
}

function show(data?: IGatewayRoute) {
  if (data) {
    routeStore.form.id = data.id
    routeStore.form.service = data.service_id
    routeStore.form.method = data.method
    routeStore.form.path_pattern = data.path_pattern
    routeStore.form.permission_match_mode = data.permission_match_mode
    routeStore.form.permissions = data.permissions?.map((p) => p.id) || []
    routeStore.form.rate_limit_per_minute = data.rate_limit_per_minute
    routeStore.form.is_active = data.is_active
  } else {
    routeStore.resetForm()
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
    if (isEdit.value && routeStore.form.id) {
      await routeStore.update(routeStore.form.id)
    } else {
      await routeStore.create()
    }
    close()
  } catch {
    // error already handled by axios interceptor
  }
}

onMounted(() => {
  loadOptions()
})

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Route' : 'Tambah Route'"
    size="2xl"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormSelect
          v-model="routeStore.form.service"
          name="service"
          label="Service"
          :options="serviceOptions"
          placeholder="Pilih service tujuan..."
          :searchable="true"
          :loading="optionsLoading"
          :validation="v$.service"
        />

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <FormSelect
            v-model="routeStore.form.method"
            name="method"
            label="Method"
            :options="methodOptions"
            :validation="v$.method"
          />

          <div class="col-span-1 sm:col-span-2">
            <FormInput
              v-model="routeStore.form.path_pattern"
              name="path_pattern"
              label="Path Pattern"
              placeholder="/user/:id atau /files/*"
              :validation="v$.path_pattern"
            />
            <p v-if="fullPathPreview" class="mt-1 text-xs text-gray-500">
              URL publik: <code>{{ fullPathPreview }}</code>
            </p>
          </div>
        </div>

        <FormSelect
          v-model="routeStore.form.permission_match_mode"
          name="permission_match_mode"
          label="Permission Match Mode"
          :options="matchModeOptions"
          :validation="v$.permission_match_mode"
        />

        <FormSelect
          v-model="routeStore.form.permissions"
          name="permissions"
          label="Permissions (kosong = publik)"
          :options="permissionOptions"
          placeholder="Pilih permission..."
          mode="tags"
          :searchable="true"
          :loading="optionsLoading"
        />

        <FormInput
          v-model="rateLimitInput"
          name="rate_limit_per_minute"
          type="number"
          label="Rate Limit per Menit (opsional override)"
          placeholder="Kosongkan untuk pakai limit Service/Global"
          :validation="v$.rate_limit_per_minute"
        />

        <label class="flex items-center gap-2 cursor-pointer">
          <input
            v-model="routeStore.form.is_active"
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
          :disabled="routeStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="routeStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
