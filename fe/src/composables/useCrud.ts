import { reactive, ref, type Ref, type Reactive } from 'vue'
import {
  get,
  post,
  put,
  del,
  ApiMetadataDefaults,
  type IApiResponse,
  type IApiMetadata,
} from '@/plugins/axios'
import swal from '@/plugins/swal'

type TLoadingKey = 'Index' | 'Form' | 'Delete'

export interface UseCrudOptions<TPayload> {
  endpoint: string
  entityName: string
  initialForm: TPayload
  formRules: Record<string, any>
  pageSize?: number
}

export function useCrud<
  TEntity extends { id: number | string },
  TPayload extends Record<string, any>,
>(options: UseCrudOptions<TPayload>) {
  const { endpoint, entityName, initialForm, formRules, pageSize = 9 } = options

  const indexData: Ref<{ items: TEntity[]; pagination: IApiMetadata }> = ref({
    items: [],
    pagination: { ...ApiMetadataDefaults, page_size: pageSize },
  })

  const singleData: Ref<TEntity | null> = ref(null)

  const loading: Ref<Record<TLoadingKey, boolean>> = ref({
    Index: false,
    Form: false,
    Delete: false,
  })

  const form: Reactive<TPayload> = reactive({ ...initialForm })

  // Extra filter/search query params merged into every fetchAll call. Persisted here so
  // pagination clicks (which call fetchAll(page) without params) keep the active filters.
  const activeParams: Ref<Record<string, any>> = ref({})

  function resetForm() {
    Object.keys(initialForm).forEach((key) => {
      const k = key as keyof TPayload
      const initial = initialForm[k]
      ;(form as any)[k] = Array.isArray(initial) ? [...initial] : initial
    })
  }

  function getCollection(): TEntity[] {
    return indexData.value.items
  }

  async function fetchAll(page?: number, params?: Record<string, any>): Promise<TEntity[]> {
    loading.value.Index = true
    const currentPage = page ?? indexData.value.pagination.page
    if (params) {
      activeParams.value = params
    }
    try {
      const { data } = await get<IApiResponse<TEntity[]>>(endpoint, {
        params: {
          page: currentPage,
          page_size: indexData.value.pagination.page_size,
          ...activeParams.value,
        },
      })
      indexData.value.items = data.data || []
      indexData.value.pagination = data.metadata || ApiMetadataDefaults
      return indexData.value.items
    } catch (error: any) {
      console.error(`Failed to fetch ${entityName}s`, error)
      swal.error('Gagal', `Gagal memuat daftar ${entityName}.`)
      return []
    } finally {
      loading.value.Index = false
    }
  }

  async function fetchById(id: number | string): Promise<TEntity | null> {
    try {
      const { data } = await get<IApiResponse<TEntity>>(`${endpoint}/${id}`)
      singleData.value = data.data || null
      return singleData.value
    } catch (error: any) {
      console.error(`Failed to fetch ${entityName}`, error)
      return null
    }
  }

  async function create() {
    loading.value.Form = true
    try {
      await post(endpoint, form)
      swal.success(
        'Berhasil',
        `${entityName.charAt(0).toUpperCase() + entityName.slice(1)} berhasil dibuat.`,
      )
      await fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || `Gagal membuat ${entityName}.`
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function update(id: number | string) {
    loading.value.Form = true
    try {
      await put(`${endpoint}/${id}`, form)
      swal.success(
        'Berhasil',
        `${entityName.charAt(0).toUpperCase() + entityName.slice(1)} berhasil diperbarui.`,
      )
      await fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || `Gagal memperbarui ${entityName}.`
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Form = false
    }
  }

  async function remove(id: number | string) {
    const result = await swal.warning(
      'Hapus Data',
      `Apakah Anda yakin ingin menghapus ${entityName} ini? Tindakan ini tidak dapat dibatalkan.`,
    )

    if (!result.isConfirmed) {
      return
    }

    loading.value.Delete = true
    try {
      await del(`${endpoint}/${id}`)
      swal.success(
        'Berhasil',
        `${entityName.charAt(0).toUpperCase() + entityName.slice(1)} berhasil dihapus.`,
      )
      await fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || `Gagal menghapus ${entityName}.`
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Delete = false
    }
  }

  return {
    endpoint,
    entityName,
    indexData,
    singleData,
    loading,
    form,
    formRules,
    getCollection,
    resetForm,
    fetchAll,
    fetchById,
    create,
    update,
    remove,
  }
}
