import { post, type IApiResponse } from '@/plugins/axios'

export interface IUploadedFile {
  uuid: string
  url: string
}

export async function uploadFile(file: File): Promise<IUploadedFile> {
  const formData = new FormData()
  formData.append('file', file)

  const response = await post<IApiResponse<IUploadedFile>>('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

  return response.data.data
}
