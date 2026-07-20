import { useMessageRateStatsStore } from '@/stores/modules/message-rate-stats'

/**
 * Message rate stats store facade.
 *
 * @returns {Object}
 */
export const useMessageRateStats = () => {
  const messageRateStatsStore = useMessageRateStatsStore()

  return {
    fetchMessageRateStatsWithoutStore: messageRateStatsStore.fetchMessageRateStatsWithoutStore,
  }
}
