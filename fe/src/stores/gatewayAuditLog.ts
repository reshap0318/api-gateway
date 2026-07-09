import { defineStore } from 'pinia'
import { useCrud } from '@/composables'

export interface IGatewayAuditLog {
  id: number
  entity_type: string
  entity_id: number
  action: string
  actor_user_id: number
  actor_name: string
  changes: string
  created_at: string
}

// Audit Log is read-only (no create/update/delete from the UI) — useCrud is still used for
// its fetchAll/pagination/loading plumbing, but the payload type is intentionally empty.
export type IGatewayAuditLogPayload = Record<string, never>

export const useGatewayAuditLogStore = defineStore('gatewayAuditLog', () => {
  const crud = useCrud<IGatewayAuditLog, IGatewayAuditLogPayload>({
    endpoint: '/audit-logs',
    entityName: 'audit log',
    initialForm: {} as IGatewayAuditLogPayload,
    formRules: {},
  })

  return { ...crud }
})
