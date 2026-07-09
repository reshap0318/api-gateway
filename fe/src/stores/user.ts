import { defineStore } from 'pinia'
import { required, email, minLength, helpers, requiredIf } from '@vuelidate/validators'
import { IRole } from './role'
import { IPermission } from './permission'
import { useCrud, withFile } from '@/composables'
import { get, put, post, type IApiResponse } from '@/plugins/axios'

export type TUserStatus = 'active' | 'suspended'

export interface IUser {
  id: number
  email: string
  name: string
  avatar: string | null
  status: TUserStatus
  locked_until: string | null
  created_at: string
  roles: IRole[]
  permissions: IPermission[]
}

export interface IUserPayload {
  id?: number
  name: string
  email: string
  password: string
  password_confirmation: string
  roles: number[]
  avatar: File | null
}

export const useUserStore = defineStore('user', () => {
  const crud = useCrud<IUser, IUserPayload>({
    endpoint: '/users',
    entityName: 'user',
    initialForm: {
      name: '',
      email: '',
      password: '',
      password_confirmation: '',
      roles: [],
      avatar: null,
    },
    formRules: {
      name: { required, minLength: minLength(2) },
      email: { required, email },
      password: { required, minLength: minLength(6) },
      password_confirmation: {},
    },
    pageSize: 12,
  })

  const userCrud = withFile<IUser, IUserPayload>(crud, ['avatar'])

  const formRules = {
    ...crud.formRules,
    password_confirmation: {
      requiredIf: requiredIf(() => !!crud.form.password),
      sameAsPassword: helpers.withMessage(
        'Password tidak cocok',
        (value: string) => value === crud.form.password,
      ),
    },
    roles: {
      required: helpers.withMessage('Role wajib dipilih', (value: number[]) => value.length > 0),
    },
  }

  async function create() {
    try {
      await userCrud.createForm(['id'])
    } catch (error: any) {
      console.error('Failed to create user', error)
      throw error
    }
  }

  async function update(id: number) {
    try {
      const excludeFields: (keyof IUserPayload)[] = ['id']
      if (!crud.form.password) {
        excludeFields.push('password', 'password_confirmation')
      }
      await userCrud.updateForm(id, excludeFields)
    } catch (error: any) {
      console.error('Failed to update user', error)
      throw error
    }
  }

  async function fetchAllUsers(): Promise<IUser[]> {
    try {
      const { data } = await get<IApiResponse<IUser[]>>('/users')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all users', error)
      return []
    }
  }

  // Status (§2.25) and Lock (§2.26/§2.27) are independent mechanisms with their own
  // endpoints — not bundled into the main update() form submit.
  async function updateStatus(id: number, status: TUserStatus): Promise<IUser | null> {
    try {
      const { data } = await put<IApiResponse<IUser>>(`/users/${id}/status`, { status })
      return data.data
    } catch (error: any) {
      console.error('Failed to update user status', error)
      throw error
    }
  }

  async function unlock(id: number): Promise<IUser | null> {
    try {
      const { data } = await post<IApiResponse<IUser>>(`/users/${id}/unlock`, {})
      return data.data
    } catch (error: any) {
      console.error('Failed to unlock user', error)
      throw error
    }
  }

  return {
    ...userCrud,
    formRules,
    create,
    update,
    fetchAllUsers,
    updateStatus,
    unlock,
  }
})
