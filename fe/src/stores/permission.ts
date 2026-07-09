import { defineStore } from 'pinia'
import { required, minLength } from '@vuelidate/validators'
import { get, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables/useCrud'

export interface IPermission {
  id: number
  name: string
  description: string
}

export interface IPermissionPayload {
  id?: number
  name: string
  description: string
}

export const usePermissionStore = defineStore('permission', () => {
  const crud = useCrud<IPermission, IPermissionPayload>({
    endpoint: '/permissions',
    entityName: 'permission',
    initialForm: { name: '', description: '' },
    formRules: {
      name: { required, minLength: minLength(3) },
      description: {},
    },
  })

  async function fetchAllPermissions(): Promise<IPermission[]> {
    try {
      const { data } = await get<IApiResponse<IPermission[]>>('/permissions')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all permissions', error)
      return []
    }
  }

  return {
    ...crud,
    fetchAllPermissions,
  }
})
