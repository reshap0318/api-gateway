<script setup lang="ts">
import { FormInput, FormPassword, UiButton, UiCard } from '@/components/utils'
import useVuelidate from '@vuelidate/core'
import { useRouter } from 'vue-router'
import { PhEnvelope } from '@phosphor-icons/vue'
import swal from '@/plugins/swal'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const v$ = useVuelidate(authStore.formRules, authStore.form)

async function handleLogin() {
  const isValid = await v$.value.$validate()
  if (!isValid) return

  try {
    await authStore.login()
    swal.success('Login Berhasil')
    router.push('/')
  } catch (error: any) {
    const message =
      error?.response?.data?.message || 'Login gagal, periksa kembali email dan password'
    swal.error('Login Gagal', message)
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <UiCard :padded="false" :classes="{ card: 'p-8' }" class="max-w-md">
      <template #header>
        <h2 class="mb-6 text-center text-2xl font-bold text-gray-800">Login</h2>
      </template>
      <form @submit.prevent="handleLogin">
        <FormInput
          v-model="authStore.form.email"
          name="email"
          label="Email"
          type="email"
          class="mb-4"
          placeholder="admin@example.com"
          :validation="v$.email"
          :prefix-icon="PhEnvelope"
        />

        <FormPassword
          v-model="authStore.form.password"
          name="password"
          label="Password"
          placeholder="••••••••"
          class="mb-2"
          :validation="v$.password"
        />

        <div class="mb-6 text-right">
          <router-link
            to="/forgot-password"
            class="text-sm font-medium text-primary-600 hover:text-primary-700"
          >
            Lupa Password?
          </router-link>
        </div>

        <UiButton type="submit" full-width :loading="authStore.isLoading"> Login </UiButton>
      </form>
    </UiCard>
  </div>
</template>
