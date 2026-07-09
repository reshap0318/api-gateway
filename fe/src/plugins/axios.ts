declare module 'axios' {
  interface AxiosRequestConfig {
    cancelPreviousRequests?: boolean
    hideError?: boolean
    hideError400?: boolean
    hideError401?: boolean
    hideError403?: boolean
    hideError404?: boolean
    hideError409?: boolean
    hideError429?: boolean
    hideError500?: boolean
    hideError502?: boolean
    hideError503?: boolean
    hideError504?: boolean
    hideErrorNetwork?: boolean
  }
}

import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useFormError } from '@/composables/useFormError'
import swal from './swal'

export interface IApiMetadata {
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface IApiResponse<TData> {
  code: number
  message: string
  data: TData
  metadata?: IApiMetadata
}

export const ApiMetadataDefaults: IApiMetadata = {
  total: 0,
  page: 1,
  page_size: 9,
  total_pages: 1,
}

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (token && !config.headers.Authorization) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// Response interceptor
let isHandling401 = false

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const config = error.config
    const shouldHideAll = config?.hideError === true

    // Network error (no response)
    if (!error.response) {
      if (!shouldHideAll && !config?.hideErrorNetwork) {
        if (error.code === 'ECONNABORTED') {
          swal.error('Timeout', 'Request timed out. Please check your connection and try again.')
        } else if (error.message === 'Network Error') {
          swal.error(
            'Network Error',
            'Unable to connect to the server. Please check your internet connection.',
          )
        } else {
          swal.error('Error', 'An unexpected error occurred. Please try again.')
        }
      }
      return Promise.reject(error)
    }

    const status = error.response.status
    const data = error.response.data
    const message = data?.message

    if (shouldHideAll) {
      return Promise.reject(error)
    }

    switch (status) {
      case 400:
        if (!config?.hideError400) {
          swal.error('Bad Request', message || 'Invalid request. Please check your input.')
        }
        break

      case 401:
        if (!config?.hideError401) {
          if (!isHandling401) {
            isHandling401 = true
            swal
              .error('Unauthorized', message || 'Session expired. Please login again.')
              .then(() => {
                const authStore = useAuthStore()
                authStore.logout().finally(() => {
                  isHandling401 = false
                  if (!window.location.pathname.includes('/login')) {
                    window.location.href = '/login'
                  }
                })
              })
          }
        }
        break

      case 403:
        if (!config?.hideError403) {
          swal.error('Forbidden', message || 'You do not have permission to access this resource.')
        }
        break

      case 404:
        if (!config?.hideError404) {
          swal.error('Not Found', message || 'The requested resource was not found.')
        }
        break

      case 409:
        if (!config?.hideError409) {
          swal.error('Conflict', message || 'This action conflicts with existing data.')
        }
        break

      case 422:
        if (data?.errors) {
          const formError = useFormError()
          formError.set(data.errors)
        }
        break

      case 429:
        if (!config?.hideError429) {
          swal.error(
            'Too Many Requests',
            'You are making too many requests. Please wait a moment and try again.',
          )
        }
        break

      case 500:
        if (!config?.hideError500) {
          swal.error('Server Error', message || 'Internal server error. Please try again later.')
        }
        break

      case 502:
        if (!config?.hideError502) {
          swal.error(
            'Bad Gateway',
            message || 'The server is temporarily unavailable. Please try again later.',
          )
        }
        break

      case 503:
        if (!config?.hideError503) {
          swal.error(
            'Service Unavailable',
            message || 'The service is currently unavailable. Please try again later.',
          )
        }
        break

      case 504:
        if (!config?.hideError504) {
          swal.error(
            'Gateway Timeout',
            message || 'The server took too long to respond. Please try again later.',
          )
        }
        break

      default:
        if (status >= 500) {
          swal.error('Server Error', message || 'An unexpected server error occurred.')
        } else if (status >= 400) {
          swal.error(`Error ${status}`, message || 'An unexpected error occurred.')
        }
        break
    }

    return Promise.reject(error)
  },
)

// Helper methods
const get = <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> => {
  return api.get<T>(url, config)
}

const post = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.post<T>(url, data, config)
}

const put = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.put<T>(url, data, config)
}

const patch = <T = unknown>(
  url: string,
  data?: unknown,
  config?: AxiosRequestConfig,
): Promise<AxiosResponse<T>> => {
  return api.patch<T>(url, data, config)
}

const del = <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> => {
  return api.delete<T>(url, config)
}

export { get, post, put, patch, del }
export default api
