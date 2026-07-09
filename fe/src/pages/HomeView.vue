<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import swal from '@/plugins/swal'

const route = useRoute()
const router = useRouter()

onMounted(() => {
  if (route.query.accessDenied === 'true') {
    swal.warning('Access Denied', 'Anda tidak memiliki akses ke halaman tersebut.')
    router.replace({ query: {} })
  }
})
</script>

<template>
  <div class="mx-auto px-4">
    <h1 class="text-2xl font-bold text-gray-900 mb-4">Dashboard</h1>
    <div class="bg-white p-6 rounded-lg shadow">
      <p class="text-gray-600">Welcome to Vue Boilerplate with Vite + TypeScript</p>
    </div>
  </div>
</template>
