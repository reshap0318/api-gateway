import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { get, put, type IApiResponse } from '@/plugins/axios'
import { required, email, minLength, helpers, requiredIf } from '@vuelidate/validators'
import { uploadFile } from '@/helpers/upload'
import storage from '@/helpers/storage'
import swal from '@/plugins/swal'
import { useAuthStore } from './auth'
import type { IRole } from './role'
import type { IPermission } from './permission'

export interface IProfile {
  id: number
  email: string
  name: string
  avatar: string | null
  created_at: string
  roles: IRole[]
  permissions: IPermission[]
}

export interface IProfilePayload {
  name: string
  email: string
  password: string
  password_confirmation: string
  avatar: File | null
}

export const useProfileStore = defineStore('profile', () => {
  const profile = ref<IProfile | null>(null)
  const loading = ref<Record<string, boolean>>({
    Fetch: false,
    Update: false,
  })

  const form = reactive<IProfilePayload>({
    name: '',
    email: '',
    password: '',
    password_confirmation: '',
    avatar: null,
  })

  const formRules = {
    name: { required, minLength: minLength(2) },
    email: { required, email },
    password: { minLength: minLength(6) },
    password_confirmation: {
      requiredIf: requiredIf(() => !!form.password),
      sameAsPassword: helpers.withMessage(
        'Konfirmasi password tidak cocok',
        (value: string) => !form.password || value === form.password,
      ),
    },
  }

  async function fetchProfile() {
    loading.value.Fetch = true
    try {
      const { data } = await get<IApiResponse<IProfile>>('/me')
      profile.value = data.data || null
      if (profile.value) {
        form.name = profile.value.name
        form.email = profile.value.email
      }
    } catch (error: any) {
      console.error('Failed to fetch profile', error)
      swal.error('Gagal', 'Gagal memuat data profile.')
    } finally {
      loading.value.Fetch = false
    }
  }

  async function updateProfile() {
    loading.value.Update = true
    try {
      const payload: Record<string, any> = {
        name: form.name,
        email: form.email,
      }

      if (form.avatar) {
        const uploaded = await uploadFile(form.avatar)
        payload.avatar = uploaded.uuid
      }

      if (form.password) {
        payload.password = form.password
        payload.password_confirmation = form.password_confirmation
      }

      const { data } = await put<IApiResponse<IProfile>>('/me', payload)
      profile.value = data.data || null

      const authStore = useAuthStore()
      if (profile.value && authStore.user) {
        authStore.user = {
          ...authStore.user,
          name: profile.value.name,
          email: profile.value.email,
          avatar: profile.value.avatar,
        }
        storage.setItem('user', authStore.user)
      }

      form.password = ''
      form.password_confirmation = ''
      form.avatar = null

      swal.success('Berhasil', 'Profile berhasil diperbarui.')
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memperbarui profile.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Update = false
    }
  }

  return {
    profile,
    loading,
    form,
    formRules,
    fetchProfile,
    updateProfile,
  }
})
