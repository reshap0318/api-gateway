import { ref } from 'vue'

const errors = ref<Record<string, string[]>>({})

export function useFormError() {
  function set(newErrors: Record<string, string[]>) {
    errors.value = newErrors
  }

  function has(field: string) {
    return field in errors.value
  }

  function get(field: string) {
    return errors.value[field] || []
  }

  function clear(field?: string) {
    if (field) {
      delete errors.value[field]
    } else {
      errors.value = {}
    }
  }

  return { errors, set, has, get, clear }
}
