<script setup lang="ts">
import { FormInput, UiButton, UiCard } from '@/components/utils'
import { ref } from 'vue'
import { required, email } from '@vuelidate/validators'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import useVuelidate from '@vuelidate/core'
import swal from '@/plugins/swal'

const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  email: '',
})

const rules = {
  email: { required, email },
}

const v$ = useVuelidate(rules, formData)
const isLoading = ref(false)

async function handleForgotPassword() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  isLoading.value = true
  try {
    await authStore.forgotPassword(formData.value.email)
    swal.success('Berhasil', 'Tautan untuk mengatur ulang kata sandi telah dikirim ke email Anda.')
    router.push('/login')
  } catch (error: any) {
    const message = error?.response?.data?.message || 'Gagal mengirim email reset password.'
    swal.error('Gagal', message)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <UiCard :padded="false" :classes="{ card: 'p-8' }" class="max-w-md w-full">
      <template #header>
        <h2 class="mb-2 text-center text-2xl font-bold text-gray-800">Lupa Password</h2>
        <p class="mb-6 text-center text-sm text-gray-500">
          Masukkan email Anda dan kami akan mengirimkan tautan untuk mengatur ulang kata sandi.
        </p>
      </template>
      <form @submit.prevent="handleForgotPassword">
        <FormInput
          v-model="formData.email"
          label="Email"
          type="email"
          class="mb-6"
          placeholder="admin@example.com"
          :validation="v$.email"
        />

        <UiButton type="submit" full-width :loading="isLoading" class="mb-4">
          Kirim Tautan Reset
        </UiButton>

        <div class="text-center text-sm">
          <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
            Kembali ke Login
          </router-link>
        </div>
      </form>
    </UiCard>
  </div>
</template>
