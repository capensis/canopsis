import { ref } from 'vue'

export const usePendingHandler = (callback, initialPending = false) => {
  const pending = ref(initialPending)

  const handler = async (...args) => {
    try {
      pending.value = true

      return await callback(...args)
    } finally {
      pending.value = false
    }
  }

  return {
    pending,
    handler,
  }
}
