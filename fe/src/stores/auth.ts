import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'
import { post } from '@/plugins/axios'
import storage from '@/helpers/storage'
import type { IApiResponse } from '@/plugins/axios'
import { required, email } from '@vuelidate/validators'

export interface ILoginPayload {
  email: string
  password: string
}

export interface ILoginResponse {
  token: string
  refresh_token: string
  user: {
    id: number
    name: string
    email: string
    avatar: string | null
    created_at: string
    roles: { id: number; name: string; description: string }[]
    permissions: { id: number; name: string; description: string }[]
  }
}

export interface IRefreshTokenPayload {
  refresh_token: string
}

export interface IRefreshTokenResponse {
  token: string
  refresh_token: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getItem<string>('token'))
  const refreshToken = ref<string | null>(storage.getItem<string>('refresh_token'))
  const user = ref<ILoginResponse['user'] | null>(storage.getItem<ILoginResponse['user']>('user'))
  const isLoading = ref(false)
  const isAuthenticated = computed((): boolean => !!token.value)

  const form = reactive<ILoginPayload>({
    email: '',
    password: '',
  })

  const formRules = {
    email: {
      required,
      email,
    },
    password: {
      required,
    },
  }

  async function login(): Promise<void> {
    isLoading.value = true
    try {
      const response = await post<IApiResponse<ILoginResponse>>(
        '/auth/login',
        { email: form.email, password: form.password },
        { hideError: true },
      )
      const { token: newToken, refresh_token: newRefreshToken, user: userData } = response.data.data

      token.value = newToken
      refreshToken.value = newRefreshToken
      user.value = userData

      storage.setItem('token', newToken)
      storage.setItem('refresh_token', newRefreshToken)
      storage.setItem('user', userData)
    } finally {
      isLoading.value = false
    }
  }

  async function refreshTokenFn(): Promise<void> {
    if (!refreshToken.value) return
    try {
      const response = await post<IApiResponse<IRefreshTokenResponse>>('/auth/refresh', {
        refresh_token: refreshToken.value,
      })
      const { token: newToken, refresh_token: newRefreshToken } = response.data.data

      token.value = newToken
      refreshToken.value = newRefreshToken

      storage.setItem('token', newToken)
      storage.setItem('refresh_token', newRefreshToken)
    } catch (error) {
      console.error('Failed to refresh token', error)
      await logout()
    }
  }

  async function logout(): Promise<void> {
    try {
      if (token.value) {
        await post('/auth/logout')
      }
    } catch (error) {
      console.error('Logout error', error)
    } finally {
      token.value = null
      refreshToken.value = null
      user.value = null
      storage.clearAll()
    }
  }

  async function forgotPassword(email: string): Promise<void> {
    await post('/auth/forgot-password', { email })
  }

  async function resetPassword(token: string, new_password: string): Promise<void> {
    await post('/auth/reset-password', { token, new_password })
  }

  return {
    token,
    refreshToken,
    user,
    isLoading,
    form,
    formRules,
    login,
    refreshTokenFn,
    logout,
    forgotPassword,
    resetPassword,
    isAuthenticated,
  }
})
