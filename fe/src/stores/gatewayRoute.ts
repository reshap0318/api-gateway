import { defineStore } from 'pinia'
import { required, minValue, helpers } from '@vuelidate/validators'
import { post, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables'
import type { IPermission } from './permission'

export interface IGatewayRouteServiceMini {
  id: number
  name: string
  base_path: string
  is_active: boolean
}

export interface IGatewayRoute {
  id: number
  service_id: number
  service: IGatewayRouteServiceMini
  method: string
  path_pattern: string
  permission_match_mode: 'any' | 'all'
  permissions: IPermission[]
  rate_limit_per_minute: number | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface IGatewayRoutePayload {
  id?: number
  service: number | null
  method: string
  path_pattern: string
  permission_match_mode: 'any' | 'all'
  permissions: number[]
  rate_limit_per_minute: number | null
  is_active: boolean
}

export const useGatewayRouteStore = defineStore('gatewayRoute', () => {
  const crud = useCrud<IGatewayRoute, IGatewayRoutePayload>({
    endpoint: '/routes',
    entityName: 'route',
    initialForm: {
      service: null,
      method: 'GET',
      path_pattern: '',
      permission_match_mode: 'any',
      permissions: [],
      rate_limit_per_minute: null,
      is_active: true,
    },
    formRules: {
      service: { required },
      method: { required },
      path_pattern: {
        required,
        startsWithSlash: helpers.withMessage(
          'Path harus diawali dengan /',
          (value: string) => !value || value.startsWith('/'),
        ),
      },
      permission_match_mode: { required },
      rate_limit_per_minute: {
        minValue: helpers.withMessage('Rate limit harus lebih dari 0', minValue(1)),
      },
    },
  })

  async function refreshCache(): Promise<void> {
    try {
      await post<IApiResponse<{ refreshed_at: string }>>('/gateway/cache/refresh', {})
    } catch (error: any) {
      console.error('Failed to refresh gateway cache', error)
      throw error
    }
  }

  return {
    ...crud,
    refreshCache,
  }
})
