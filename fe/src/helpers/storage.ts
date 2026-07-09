const getItem = <T = unknown>(key: string, defaultValue: T | null = null): T | null => {
  try {
    const item = localStorage.getItem(key)
    if (item === null) return defaultValue
    return JSON.parse(item) as T
  } catch {
    return defaultValue
  }
}

const setItem = <T = unknown>(key: string, value: T): void => {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    console.error('Failed to save to localStorage:', error)
  }
}

const removeItem = (key: string): void => {
  try {
    localStorage.removeItem(key)
  } catch (error) {
    console.error('Failed to remove from localStorage:', error)
  }
}

const clearAll = (): void => {
  try {
    localStorage.clear()
  } catch (error) {
    console.error('Failed to clear localStorage:', error)
  }
}

export default {
  getItem,
  setItem,
  removeItem,
  clearAll,
}
