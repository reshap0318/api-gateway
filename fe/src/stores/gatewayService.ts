import { defineStore } from 'pinia'
import { required, minLength, minValue, helpers } from '@vuelidate/validators'
import { get, post, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables'

export interface IGatewayServiceHealth {
  health_status: 'unknown' | 'up' | 'down'
  health_checked_at: string | null
}

// @vuelidate/validators' `url` rule requires a TLD-style host and rejects things like
// "http://localhost:9000" — upstream services are frequently internal hosts/IPs, so a
// lenient scheme+host check (via the URL constructor) is used instead.
function isValidBaseUrl(value: string): boolean {
  if (!value) return true
  try {
    const parsed = new URL(value)
    return (parsed.protocol === 'http:' || parsed.protocol === 'https:') && !!parsed.hostname
  } catch {
    return false
  }
}

export interface IGatewayService {
  id: number
  name: string
  base_url: string
  protocol: 'http' | 'websocket'
  is_active: boolean
  rate_limit_per_minute: number | null
  health_status: 'unknown' | 'up' | 'down'
  health_checked_at: string | null
  route_count: number
  created_at: string
  updated_at: string
}

export interface IGatewayServicePayload {
  id?: number
  name: string
  base_url: string
  protocol: 'http' | 'websocket'
  rate_limit_per_minute: number | null
  is_active: boolean
}

export const useGatewayServiceStore = defineStore('gatewayService', () => {
  const crud = useCrud<IGatewayService, IGatewayServicePayload>({
    endpoint: '/services',
    entityName: 'service',
    initialForm: {
      name: '',
      base_url: '',
      protocol: 'http',
      rate_limit_per_minute: null,
      is_active: true,
    },
    formRules: {
      name: { required, minLength: minLength(3) },
      base_url: {
        required,
        validUrl: helpers.withMessage('Base URL tidak valid', isValidBaseUrl),
      },
      protocol: { required },
      rate_limit_per_minute: {
        minValue: helpers.withMessage('Rate limit harus lebih dari 0', minValue(1)),
      },
    },
  })

  async function fetchAllServices(): Promise<IGatewayService[]> {
    try {
      const { data } = await get<IApiResponse<IGatewayService[]>>('/services')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all services', error)
      return []
    }
  }

  async function healthCheck(id: number): Promise<void> {
    try {
      await post<IApiResponse<IGatewayServiceHealth>>(`/services/${id}/health-check`, {})
      await crud.fetchAll()
    } catch (error: any) {
      console.error('Failed to run health check', error)
      throw error
    }
  }

  return {
    ...crud,
    fetchAllServices,
    healthCheck,
  }
})
