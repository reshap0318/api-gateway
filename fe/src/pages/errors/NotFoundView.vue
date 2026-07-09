<script setup lang="ts">
import { useRouter } from 'vue-router'
import { UiButton } from '@/components/utils'
import { PhHouse, PhWarningCircle } from '@phosphor-icons/vue'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100 px-4">
    <div class="w-full max-w-lg text-center">
      <!-- 404 Icon -->
      <div class="mb-6 flex justify-center">
        <PhWarningCircle class="h-24 w-24 text-blue-600 sm:h-32 sm:w-32" weight="duotone" />
      </div>

      <!-- Error Code -->
      <h1 class="mb-2 text-6xl font-extrabold text-gray-900 sm:text-8xl">404</h1>

      <!-- Message -->
      <h2 class="mb-4 text-xl font-semibold text-gray-800 sm:text-2xl">Halaman Tidak Ditemukan</h2>
      <p class="mb-8 text-base text-gray-500 sm:text-lg">
        Maaf, halaman yang Anda cari tidak ditemukan atau telah dipindahkan.
      </p>

      <!-- Action Button -->
      <UiButton variant="primary" size="lg" @click="goHome">
        <template #icon>
          <PhHouse class="h-5 w-5" />
        </template>
        Kembali ke Dashboard
      </UiButton>
    </div>
  </div>
</template>
