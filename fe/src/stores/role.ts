import { defineStore } from 'pinia'
import { required, minLength } from '@vuelidate/validators'
import { get, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables/useCrud'
import { IPermission } from './permission'

export interface IRole {
  id: number
  name: string
  description: string
  permissions: IPermission[]
}

export interface IRolePayload {
  id?: number
  name: string
  description: string
  permissions: number[]
}

export const useRoleStore = defineStore('role', () => {
  const crud = useCrud<IRole, IRolePayload>({
    endpoint: '/roles',
    entityName: 'role',
    initialForm: { name: '', description: '', permissions: [] },
    formRules: {
      name: { required, minLength: minLength(3) },
      description: {},
    },
  })

  async function fetchAllRoles(): Promise<IRole[]> {
    try {
      const { data } = await get<IApiResponse<IRole[]>>('/roles')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all roles', error)
      return []
    }
  }

  async function fetchRolePermissions(id: number): Promise<IPermission[]> {
    try {
      const { data } = await get<IApiResponse<IPermission[]>>(`/roles/${id}/permissions`)
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch role permissions', error)
      return []
    }
  }

  return {
    ...crud,
    fetchAllRoles,
    fetchRolePermissions,
  }
})
