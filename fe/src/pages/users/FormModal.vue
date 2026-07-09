<script setup lang="ts">
import {
  UiModal,
  FormInput,
  FormPassword,
  FormSelect,
  FormAvatar,
  UiButton,
  UiBadge,
} from '@/components/utils'
import { computed, ref, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { useUserStore } from '@/stores/user'
import type { TUserStatus } from '@/stores/user'
import { useRoleStore } from '@/stores/role'
import type { IRole } from '@/stores/role'
import { minLength } from '@vuelidate/validators'
import { useFormError } from '@/composables/useFormError'
import swal from '@/plugins/swal'

const userStore = useUserStore()
const roleStore = useRoleStore()
const formError = useFormError()
const isVisible = ref(false)
const isEdit = computed(() => !!userStore.form.id)
const allRoles = ref<IRole[]>([])
const rolesLoading = ref(false)
const currentAvatar = ref<string | null>(null)

// User Status & Lock Panel (§FSD 3.7) — independent action, immediate effect, not bundled
// into the main form submit (backend exposes these as separate endpoints).
const currentStatus = ref<TUserStatus>('active')
const currentLockedUntil = ref<string | null>(null)
const statusSaving = ref(false)
const unlocking = ref(false)

const statusOptions: { value: TUserStatus; label: string }[] = [
  { value: 'active', label: 'Active' },
  { value: 'suspended', label: 'Suspend' },
]

const isLocked = computed(() => {
  if (!currentLockedUntil.value) return false
  return new Date(currentLockedUntil.value).getTime() > Date.now()
})

async function handleStatusChange(newStatus: TUserStatus) {
  if (!userStore.form.id || newStatus === currentStatus.value) return
  statusSaving.value = true
  try {
    const updated = await userStore.updateStatus(userStore.form.id, newStatus)
    if (updated) currentStatus.value = updated.status
    swal.success('Berhasil', 'Status user berhasil diperbarui.')
  } catch {
    swal.error('Gagal', 'Gagal memperbarui status user.')
  } finally {
    statusSaving.value = false
  }
}

async function handleUnlock() {
  if (!userStore.form.id) return
  unlocking.value = true
  try {
    const updated = await userStore.unlock(userStore.form.id)
    currentLockedUntil.value = updated?.locked_until ?? null
    swal.success('Berhasil', 'Akun berhasil dibuka kuncinya.')
  } catch {
    swal.error('Gagal', 'Gagal membuka kunci akun.')
  } finally {
    unlocking.value = false
  }
}

const dynamicRules = computed(() => {
  if (isEdit.value) {
    return {
      ...userStore.formRules,
      password: { minLength: minLength(6) },
    }
  }
  return userStore.formRules
})

const v$ = useVuelidate(dynamicRules, userStore.form)

const roleOptions = computed(() => {
  return allRoles.value.map((role) => ({
    value: role.id,
    label: role.name,
  }))
})

async function loadRoles() {
  if (allRoles.value.length > 0) return
  rolesLoading.value = true
  try {
    allRoles.value = await roleStore.fetchAllRoles()
  } finally {
    rolesLoading.value = false
  }
}

async function show(data?: {
  id?: number
  name: string
  email: string
  avatar?: string | null
  roles?: { id: number }[]
  status?: TUserStatus
  locked_until?: string | null
}) {
  if (data) {
    userStore.form.id = data.id
    userStore.form.name = data.name
    userStore.form.email = data.email
    userStore.form.password = ''
    userStore.form.password_confirmation = ''
    userStore.form.roles = data.roles?.map((r) => r.id) || []
    userStore.form.avatar = null
    currentAvatar.value = data.avatar || null
    currentStatus.value = data.status || 'active'
    currentLockedUntil.value = data.locked_until || null
  } else {
    userStore.form.id = undefined
    userStore.form.name = ''
    userStore.form.email = ''
    userStore.form.password = ''
    userStore.form.password_confirmation = ''
    userStore.form.roles = []
    userStore.form.avatar = null
    currentAvatar.value = null
    currentStatus.value = 'active'
    currentLockedUntil.value = null
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
    if (isEdit.value && userStore.form.id) {
      await userStore.update(userStore.form.id)
    } else {
      await userStore.create()
    }
    close()
  } catch {
    // error already handled by axios interceptor
  }
}

onMounted(() => {
  loadRoles()
})

defineExpose({ show, close })
</script>

<template>
  <UiModal
    v-model="isVisible"
    :title="isEdit ? 'Edit User' : 'Tambah User'"
    size="2xl"
    @close="close"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <FormAvatar v-model="userStore.form.avatar" :current-avatar="currentAvatar" />

        <FormInput
          v-model="userStore.form.name"
          name="name"
          label="Nama"
          placeholder="John Doe"
          :validation="v$.name"
        />

        <FormInput
          v-model="userStore.form.email"
          name="email"
          label="Email"
          type="email"
          placeholder="john@example.com"
          :validation="v$.email"
        />

        <FormPassword
          v-model="userStore.form.password"
          name="password"
          label="Password"
          :placeholder="isEdit ? 'Kosongkan jika tidak ingin mengubah' : 'password123'"
          :validation="v$.password"
        />

        <FormPassword
          v-model="userStore.form.password_confirmation"
          name="password_confirmation"
          label="Konfirmasi Password"
          placeholder="password123"
          :validation="v$.password_confirmation"
        />

        <FormSelect
          v-model="userStore.form.roles"
          name="roles"
          label="Roles"
          :options="roleOptions"
          placeholder="Pilih role..."
          mode="tags"
          :searchable="true"
          :loading="rolesLoading"
        />

        <!-- User Status & Lock Panel (§FSD 3.7) — edit-only, immediate effect actions -->
        <div v-if="isEdit" class="border border-gray-200 rounded-lg p-4 space-y-3">
          <h4 class="text-sm font-semibold text-gray-800">Status & Lock</h4>

          <div class="flex items-center gap-3">
            <FormSelect
              :model-value="currentStatus"
              name="status"
              label="Status"
              :options="statusOptions"
              :disabled="statusSaving"
              :classes="{ wrapper: 'flex-1' }"
              @update:model-value="(value) => handleStatusChange(value as TUserStatus)"
            />
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-sm text-gray-600">Lock:</span>
              <UiBadge v-if="isLocked" color="danger">Locked</UiBadge>
              <UiBadge v-else color="primary">Unlocked</UiBadge>
            </div>
            <UiButton
              v-if="isLocked"
              type="button"
              size="sm"
              variant="secondary"
              outline
              :loading="unlocking"
              @click="handleUnlock"
            >
              Unlock Now
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="mt-6 flex justify-end gap-2">
        <UiButton
          type="button"
          variant="secondary"
          :disabled="userStore.loading.Form"
          outline
          @click="close"
        >
          Batal
        </UiButton>
        <UiButton type="submit" :loading="userStore.loading.Form">
          {{ isEdit ? 'Perbarui' : 'Simpan' }}
        </UiButton>
      </div>
    </form>
  </UiModal>
</template>
