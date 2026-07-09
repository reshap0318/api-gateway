// UiButton types
export interface UiButtonProps {
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  variant?: 'primary' | 'secondary' | 'danger' | 'success'
  outline?: boolean
  size?: 'sm' | 'md' | 'lg'
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full'
  fullWidth?: boolean
  loadingText?: string
}

// UiCard types
export interface UiCardClasses {
  wrapper?: string
  card?: string
  header?: string
  body?: string
  footer?: string
}

export interface UiCardProps {
  padded?: boolean
  classes?: UiCardClasses
}

// UiModal types
export type TModalSize =
  | 'sm'
  | 'md'
  | 'lg'
  | 'xl'
  | '2xl'
  | '3xl'
  | '4xl'
  | '5xl'
  | '6xl'
  | '7xl'
  | 'full'

export interface UiModalClasses {
  container?: string
  header?: string
  body?: string
  footer?: string
}

export interface UiModalProps {
  modelValue: boolean
  title?: string
  size?: TModalSize
  persistent?: boolean
  classes?: UiModalClasses
}

// UiEmptyState types
export interface UiEmptyStateClasses {
  wrapper?: string
  icon?: string
  title?: string
  description?: string
}

export interface UiEmptyStateProps {
  icon?: object
  title: string
  description?: string
  variant?: 'default' | 'compact'
  classes?: UiEmptyStateClasses
}

// UiPagination types
export interface UiPaginationClasses {
  wrapper?: string
  button?: string
  buttonActive?: string
  buttonDisabled?: string
  ellipsis?: string
}

export interface UiPaginationProps {
  page: number
  totalPages: number
  maxVisible?: number
  classes?: UiPaginationClasses
}

// FormInput types
export interface FormInputClasses {
  wrapper?: string
  label?: string
  input?: string
  error?: string
}

export interface FormInputProps {
  modelValue: string
  name?: string
  label?: string
  type?: string
  placeholder?: string
  validation?: object
  prefixIcon?: object
  suffixIcon?: object
  iconSize?: number
  classes?: FormInputClasses
}

// FormPassword types
export interface FormPasswordClasses {
  wrapper?: string
  label?: string
  input?: string
  error?: string
}

export interface FormPasswordProps {
  modelValue: string
  name?: string
  label?: string
  placeholder?: string
  validation?: object
  classes?: FormPasswordClasses
}

// FormSelect types
export interface TSelectOption {
  value: string | number
  label: string
  [key: string]: any
}

export interface FormSelectClasses {
  wrapper?: string
  label?: string
  error?: string
}

export interface FormSelectProps {
  modelValue: any
  name?: string
  label?: string
  options?: TSelectOption[] | string[]
  placeholder?: string
  validation?: object
  searchable?: boolean
  mode?: 'single' | 'multiple' | 'tags'
  closable?: boolean
  disabled?: boolean
  loading?: boolean
  classes?: FormSelectClasses
}

// FormAvatar types
export interface FormAvatarClasses {
  wrapper?: string
  label?: string
  preview?: string
  uploadButton?: string
}

export interface FormAvatarProps {
  modelValue: File | null
  currentAvatar?: string | null
  label?: string
  accept?: string
  maxSize?: number
  disabled?: boolean
  classes?: FormAvatarClasses
}

// FormFile types
export interface TFileItem {
  id: string
  file: File
  preview?: string
}

export interface FormFileClasses {
  wrapper?: string
  label?: string
  dropZone?: string
  error?: string
  fileItem?: string
}

export interface FormFileProps {
  modelValue: FileList | File[] | null
  name?: string
  label?: string
  placeholder?: string
  validation?: object
  accept?: string
  multiple?: boolean
  maxSize?: number
  disabled?: boolean
  classes?: FormFileClasses
}

// UiSkeleton types
export type TSkeletonVariant = 'text' | 'circle' | 'rect' | 'card'

export interface UiSkeletonClasses {
  wrapper?: string
  bar?: string
}

export interface UiSkeletonProps {
  variant?: TSkeletonVariant
  width?: string
  height?: string
  rounded?: boolean
  count?: number
  classes?: UiSkeletonClasses
}

export type TSortOrder = 'asc' | 'desc' | 'none'

// UiTable types
export interface TTableColumn {
  title: string
  data: string
  class?: string
  headerClass?: string
  sortable?: boolean
}

export interface TTableRow {
  id?: number | string
  class?: string
  [key: string]: unknown
}
