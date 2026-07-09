import Swal from 'sweetalert2'
import type { SweetAlertOptions, SweetAlertResult } from 'sweetalert2'

const defaultConfig: SweetAlertOptions = {
  confirmButtonText: 'Confirm',
  cancelButtonText: 'Cancel',
  showCancelButton: false,
  reverseButtons: true,
  width: '32rem', // w-96
  padding: '1.5rem', // p-6
  backdrop: `
    rgba(0, 0, 0, 0.5)
  `,
  customClass: {
    popup: 'rounded-xl shadow-2xl',
    container: 'backdrop-blur-sm',
    title: 'text-xl font-bold text-gray-900 dark:text-gray-100',
    htmlContainer: 'text-gray-600 dark:text-gray-300',
    closeButton: 'hover:bg-gray-100 dark:hover:bg-gray-800',
    icon: 'mb-4',
    input:
      'border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary focus:border-transparent',
    inputLabel: 'text-gray-700 dark:text-gray-300 mb-2',
    validationMessage: 'text-danger mt-2',
    actions: 'flex gap-2 mt-6',
    confirmButton:
      'px-6 py-2.5 rounded-lg font-semibold transition-all hover:opacity-90 active:scale-95 bg-blue-600 text-white',
    cancelButton:
      'px-6 py-2.5 rounded-lg font-semibold transition-all hover:opacity-90 active:scale-95 bg-gray-200 text-gray-700',
    loader: 'text-primary',
  },
  buttonsStyling: false,
}

const success = (title: string, text?: string): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    icon: 'success',
    title,
    text,
  })
}

const error = (title: string, text?: string): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    icon: 'error',
    title,
    text,
    customClass: {
      ...defaultConfig.customClass,
      confirmButton:
        'px-6 py-2.5 rounded-lg font-semibold transition-all hover:opacity-90 active:scale-95 bg-red-600 text-white',
    },
  })
}

const warning = (title: string, text?: string): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    icon: 'warning',
    title,
    text,
    showCancelButton: true,
    confirmButtonText: 'Yes',
    customClass: {
      ...defaultConfig.customClass,
      confirmButton:
        'px-6 py-2.5 rounded-lg font-semibold transition-all hover:opacity-90 active:scale-95 bg-yellow-500 text-white',
    },
  })
}

const info = (title: string, text?: string): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    icon: 'info',
    title,
    text,
  })
}

const loading = (title: string, text?: string): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    title,
    text,
    allowOutsideClick: false,
    allowEscapeKey: false,
    didOpen: () => {
      Swal.showLoading()
    },
  })
}

const custom = (options: SweetAlertOptions): Promise<SweetAlertResult> => {
  return Swal.fire({
    ...defaultConfig,
    ...options,
    customClass: {
      ...defaultConfig.customClass,
      ...options.customClass,
    },
  } as SweetAlertOptions)
}

const close = (): void => {
  Swal.close()
}

export default {
  success,
  error,
  warning,
  info,
  loading,
  custom,
  close,
}
