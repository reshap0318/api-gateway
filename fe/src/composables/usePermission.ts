import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function usePermission() {
  const authStore = useAuthStore()

  const permissions = computed(() => authStore.user?.permissions?.map((p) => p.name) || [])

  const hasPermission = (name: string): boolean => {
    return permissions.value.includes(name)
  }

  const hasAnyPermission = (names: string[]): boolean => {
    if (names.length === 0) return true
    return names.some((name) => permissions.value.includes(name))
  }

  const hasAllPermissions = (names: string[]): boolean => {
    if (names.length === 0) return true
    return names.every((name) => permissions.value.includes(name))
  }

  return {
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  }
}
