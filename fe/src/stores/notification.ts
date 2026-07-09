import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  get,
  patch,
  del,
  type IApiResponse,
  type IApiMetadata,
  ApiMetadataDefaults,
} from '@/plugins/axios'
import swal from '@/plugins/swal'
import { formatTimeForHuman } from '@/helpers/date'

export interface INotification {
  id: number
  user_id: number
  type: string
  title: string
  message: string
  data: string | null
  read_at: string | null
  created_at: string
}

export interface INotificationFilters {
  is_read?: string
  type?: string
}

export interface IUnreadCountResponse {
  unread_count: number
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<INotification[]>([])
  const unreadCount = ref(0)
  const pagination = ref<IApiMetadata>({ ...ApiMetadataDefaults })
  const loading = ref<Record<string, boolean>>({
    Index: false,
    LoadMore: false,
    UnreadCount: false,
    MarkRead: false,
    MarkAllRead: false,
    Delete: false,
  })
  const hasMore = ref(true)

  async function fetchNotifications(page?: number, filters?: INotificationFilters) {
    loading.value.Index = true
    const currentPage = page ?? 1
    try {
      const params: Record<string, any> = {
        page: currentPage,
        page_size: pagination.value.page_size,
      }
      if (filters?.is_read !== undefined) {
        params.is_read = filters.is_read
      }
      if (filters?.type) {
        params.type = filters.type
      }

      const { data } = await get<IApiResponse<INotification[]>>('/notifications', { params })
      notifications.value = data.data || []
      pagination.value = data.metadata || ApiMetadataDefaults
      hasMore.value = currentPage < (data.metadata?.total_pages ?? 1)
      return notifications.value
    } catch (error: any) {
      console.error('Failed to fetch notifications', error)
      return []
    } finally {
      loading.value.Index = false
    }
  }

  async function loadMoreNotifications(filters?: INotificationFilters) {
    if (loading.value.LoadMore || loading.value.Index || !hasMore.value) {
      return
    }

    loading.value.LoadMore = true
    const nextPage = pagination.value.page + 1
    try {
      const params: Record<string, any> = {
        page: nextPage,
        page_size: pagination.value.page_size,
      }
      if (filters?.is_read !== undefined) {
        params.is_read = filters.is_read
      }
      if (filters?.type) {
        params.type = filters.type
      }

      const { data } = await get<IApiResponse<INotification[]>>('/notifications', { params })
      const newNotifications = data.data || []
      notifications.value = [...notifications.value, ...newNotifications]
      pagination.value = data.metadata || ApiMetadataDefaults
      hasMore.value = nextPage < (data.metadata?.total_pages ?? 1)
    } catch (error: any) {
      console.error('Failed to load more notifications', error)
    } finally {
      loading.value.LoadMore = false
    }
  }

  async function fetchUnreadCount() {
    loading.value.UnreadCount = true
    try {
      const { data } = await get<IApiResponse<IUnreadCountResponse>>('/notifications/unread-count')
      unreadCount.value = data.data?.unread_count ?? 0
      return unreadCount.value
    } catch (error: any) {
      console.error('Failed to fetch unread count', error)
      return 0
    } finally {
      loading.value.UnreadCount = false
    }
  }

  async function markAsRead(id: number) {
    loading.value.MarkRead = true
    try {
      await patch(`/notifications/${id}/read`)
      // Update local state
      const notification = notifications.value.find((n) => n.id === id)
      if (notification && !notification.read_at) {
        notification.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (error: any) {
      console.error('Failed to mark notification as read', error)
    } finally {
      loading.value.MarkRead = false
    }
  }

  async function markAllAsRead() {
    loading.value.MarkAllRead = true
    try {
      await patch('/notifications/mark-all-read')
      // Update local state
      notifications.value.forEach((n) => {
        if (!n.read_at) {
          n.read_at = new Date().toISOString()
        }
      })
      unreadCount.value = 0
    } catch (error: any) {
      console.error('Failed to mark all notifications as read', error)
    } finally {
      loading.value.MarkAllRead = false
    }
  }

  async function deleteAllNotifications() {
    loading.value.Delete = true
    try {
      await del('/notifications')
      const unreadBefore = notifications.value.filter((n) => n.read_at === null).length
      notifications.value = []
      unreadCount.value = Math.max(0, unreadCount.value - unreadBefore)
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus semua notifikasi.'
      swal.error('Gagal', message)
    } finally {
      loading.value.Delete = false
    }
  }

  async function deleteNotification(id: number) {
    loading.value.Delete = true
    try {
      await del(`/notifications/${id}`)
      // Update local state
      const wasUnread = notifications.value.find((n) => n.id === id)?.read_at === null
      notifications.value = notifications.value.filter((n) => n.id !== id)
      if (wasUnread) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal menghapus notifikasi.'
      swal.error('Gagal', message)
    } finally {
      loading.value.Delete = false
    }
  }

  function isUnread(notification: INotification): boolean {
    return notification.read_at === null
  }

  function getNotificationType(type: string): 'info' | 'success' | 'warning' | 'error' {
    if (type.includes('completed') || type.includes('success')) return 'success'
    if (type.includes('warning') || type.includes('alert')) return 'warning'
    if (type.includes('error') || type.includes('failed')) return 'error'
    return 'info'
  }

  function resetNotifications() {
    notifications.value = []
    pagination.value = { ...ApiMetadataDefaults }
    hasMore.value = true
  }

  return {
    notifications,
    unreadCount,
    pagination,
    hasMore,
    loading,
    fetchNotifications,
    loadMoreNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteAllNotifications,
    deleteNotification,
    isUnread,
    formatTimeForHuman,
    getNotificationType,
    resetNotifications,
  }
})
