<script setup lang="ts">
import { FormPassword, UiButton, UiCard } from '@/components/utils'
import { ref, computed, onMounted } from 'vue'
import useVuelidate from '@vuelidate/core'
import { required, minLength, sameAs } from '@vuelidate/validators'
import swal from '@/plugins/swal'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formData = ref({
  password: '',
  confirmPassword: '',
})

const rules = computed(() => ({
  password: { required, minLength: minLength(6) },
  confirmPassword: {
    required,
    sameAsPassword: sameAs(formData.value.password),
  },
}))

const v$ = useVuelidate(rules, formData)
const isLoading = ref(false)
const token = ref('')

onMounted(() => {
  const queryToken = route.query.token as string
  if (queryToken) {
    token.value = queryToken
  } else {
    swal.error('Token tidak valid', 'Token reset password tidak ditemukan di URL.')
    router.push('/login')
  }
})

async function handleResetPassword() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  if (!token.value) {
    swal.error('Token tidak valid', 'Token reset password tidak valid.')
    return
  }

  isLoading.value = true
  try {
    await authStore.resetPassword(token.value, formData.value.password)
    swal.success('Berhasil', 'Kata sandi Anda telah berhasil diubah. Silakan login kembali.')
    router.push('/login')
  } catch (error: any) {
    const message = error?.response?.data?.message || 'Gagal mengubah kata sandi.'
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
        <h2 class="mb-2 text-center text-2xl font-bold text-gray-800">Reset Password</h2>
        <p class="mb-6 text-center text-sm text-gray-500">Silakan masukkan kata sandi baru Anda.</p>
      </template>
      <form @submit.prevent="handleResetPassword">
        <FormPassword
          v-model="formData.password"
          label="Password Baru"
          placeholder="••••••••"
          class="mb-4"
          :validation="v$.password"
        />

        <FormPassword
          v-model="formData.confirmPassword"
          label="Konfirmasi Password"
          placeholder="••••••••"
          class="mb-6"
          :validation="v$.confirmPassword"
        />

        <UiButton type="submit" full-width :loading="isLoading" class="mb-4">
          Simpan Kata Sandi
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
