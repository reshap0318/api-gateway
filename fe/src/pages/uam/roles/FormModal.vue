<script setup lang="ts">
import { UiModal, FormInput, UiButton } from '@/components/utils'
import { computed, ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useRoleStore } from '@/stores/role'
import { usePermissionStore } from '@/stores/permission'
import type { IPermission } from '@/stores/permission'
import { useFormError } from '@/composables/useFormError'

const roleStore = useRoleStore()
const permissionStore = usePermissionStore()
const formError = useFormError()
const v$ = useVuelidate(roleStore.formRules, roleStore.form)

const isVisible = ref(false)
const isEdit = computed(() => !!roleStore.form.id)
const allPermissions = ref<IPermission[]>([])
const permissionsLoading = ref(false)

const groupedPermissions = computed(() => {
  const groups: Record<string, IPermission[]> = {}
  for (const perm of allPermissions.value) {
    const group = perm.name.includes('.') ? perm.name.split('.')[0] : 'others'
    if (!groups[group]) groups[group] = []
    groups[group].push(perm)
  }
  return groups
})

async function loadPermissions() {
  if (allPermissions.value.length > 0) return
  permissionsLoading.value = true
  try {
    allPermissions.value = await permissionStore.fetchAllPermissions()
  } finally {
    permissionsLoading.value = false
  }
}

async function show(data?: {
  id?: number
  name: string
  description: string
  permissions?: { id: number }[]
}) {
  if (data) {
    roleStore.form.id = data.id
    roleStore.form.name = data.name
    roleStore.form.description = data.description || ''
    roleStore.form.permissions = data.permissions?.map((p) => p.id) || []
  } else {
    roleStore.resetForm()
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
    if (isEdit.value && roleStore.form.id) {
      await roleStore.update(roleStore.form.id)
    } else {
      await roleStore.create()
    }
    close()
  } catch {}
}

function togglePermission(id: number) {
  const idx = roleStore.form.permissions.indexOf(id)
  if (idx === -1) {
    roleStore.form.permissions.push(id)
  } else {
    roleStore.form.permissions.splice(idx, 1)
  }
}

onMounted(() => {
  loadPermissions()
})

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit Role' : 'Tambah Role'"
    size="2xl"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormInput
          v-model="roleStore.form.name"
          name="name"
          label="Nama Role"
          placeholder="e.g. Admin"
          :validation="v$.name"
        />

        <FormInput
          v-model="roleStore.form.description"
          name="description"
          label="Deskripsi (opsional)"
          placeholder="Administrator role"
          :validation="v$.description"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Permissions</label>
          <div class="max-h-48 overflow-y-auto border border-gray-200 rounded-lg scrollbar-thin">
            <div v-if="permissionsLoading" class="p-3 text-sm text-gray-400 text-center">
              Memuat permissions...
            </div>
            <template v-else>
              <div v-for="(perms, group) in groupedPermissions" :key="group" class="p-3">
                <h4 class="text-sm font-bold text-gray-800 bg-gray-100 px-2 py-1 rounded mb-2">
                  {{ group }}
                </h4>
                <div class="grid grid-cols-2 gap-x-4 gap-y-1">
                  <label
                    v-for="perm in perms"
                    :key="perm.id"
                    class="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      :checked="roleStore.form.permissions.includes(perm.id)"
                      class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                      @change="togglePermission(perm.id)"
                    />
                    <span class="text-sm text-gray-700">{{ perm.name }}</span>
                  </label>
                </div>
              </div>
              <p v-if="allPermissions.length === 0" class="p-3 text-sm text-gray-400">
                Belum ada permission.
              </p>
            </template>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="roleStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="roleStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
