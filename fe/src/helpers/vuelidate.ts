interface ErrorObject {
  $message: string
  $params?: Record<string, unknown>
  $validator?: string
}

export interface ValidationLike {
  $error: boolean
  $errors: ErrorObject[]
}

export const translations: Record<string, string> = {
  required: 'This field is required',
  minLength: 'Minimum {min} characters required',
}

export const getErrorMessage = (field: ValidationLike): string => {
  if (!field.$error || field.$errors.length === 0) return ''
  return field.$errors[0].$message
}

export const translateVuelidateError = (field: ValidationLike): string => {
  if (!field.$error || field.$errors.length === 0) return ''

  const messages = field.$errors.map((e) => {
    if (e.$validator && translations[e.$validator]) {
      return interpolateTemplate(translations[e.$validator], e.$params)
    }
    return e.$message
  })

  return messages.join(', ')
}

function interpolateTemplate(template: string, params?: Record<string, unknown>): string {
  if (!params) return template
  return template.replace(/\{(\w+)\}/g, (_, key) => {
    const value = params[key]
    return value !== undefined ? String(value) : `{${key}}`
  })
}
