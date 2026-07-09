import { post, put } from '@/plugins/axios'
import { uploadFile } from '@/helpers/upload'
import swal from '@/plugins/swal'
import type { Ref, Reactive } from 'vue'
import type { IApiMetadata } from '@/plugins/axios'

interface CrudBaseReturn<TEntity, TPayload> {
  endpoint: string
  entityName: string
  indexData: Ref<{ items: TEntity[]; pagination: IApiMetadata }>
  singleData: Ref<TEntity | null>
  loading: Ref<Record<string, boolean>>
  form: Reactive<TPayload>
  formRules: Record<string, any>
  getCollection: () => TEntity[]
  resetForm: () => void
  fetchAll: (page?: number) => Promise<TEntity[]>
  fetchById: (id: number | string) => Promise<TEntity | null>
  create: () => Promise<void>
  update: (id: number | string) => Promise<void>
  remove: (id: number | string) => Promise<void>
}

export function withFile<
  TEntity extends { id: number | string },
  TPayload extends Record<string, any>,
>(crud: CrudBaseReturn<TEntity, TPayload>, fileFields: (keyof TPayload)[]) {
  async function preparePayload(excludeFields?: (keyof TPayload)[]): Promise<Record<string, any>> {
    const uploadedMap = new Map<string, string | string[]>()

    for (const field of fileFields) {
      const value = (crud.form as any)[field]
      if (!value) continue

      if (value instanceof File) {
        const uploaded = await uploadFile(value)
        uploadedMap.set(field as string, uploaded.uuid)
      } else if (Array.isArray(value) && value.length > 0 && value[0] instanceof File) {
        const files = value as File[]
        const uploads = await Promise.all(files.map((f) => uploadFile(f)))
        uploadedMap.set(
          field as string,
          uploads.map((u) => u.uuid),
        )
      }
    }

    const payload: Record<string, any> = {}

    Object.keys(crud.form).forEach((key) => {
      const k = key as keyof TPayload
      if (excludeFields?.includes(k)) return

      if (uploadedMap.has(key)) {
        payload[key] = uploadedMap.get(key)
        return
      }

      payload[key] = (crud.form as any)[k]
    })

    return payload
  }

  async function createForm(excludeFields?: (keyof TPayload)[]) {
    crud.loading.value.Form = true
    try {
      const payload = await preparePayload(excludeFields)
      await post(crud.endpoint, payload)
      swal.success(
        'Berhasil',
        `${crud.entityName.charAt(0).toUpperCase() + crud.entityName.slice(1)} berhasil dibuat.`,
      )
      await crud.fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || `Gagal membuat ${crud.entityName}.`
      swal.error('Gagal', message)
      throw error
    } finally {
      crud.loading.value.Form = false
    }
  }

  async function updateForm(id: number | string, excludeFields?: (keyof TPayload)[]) {
    crud.loading.value.Form = true
    try {
      const payload = await preparePayload(excludeFields)
      await put(`${crud.endpoint}/${id}`, payload)
      swal.success(
        'Berhasil',
        `${crud.entityName.charAt(0).toUpperCase() + crud.entityName.slice(1)} berhasil diperbarui.`,
      )
      await crud.fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || `Gagal memperbarui ${crud.entityName}.`
      swal.error('Gagal', message)
      throw error
    } finally {
      crud.loading.value.Form = false
    }
  }

  return {
    ...crud,
    createForm,
    updateForm,
  }
}
