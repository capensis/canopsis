import { useMetricsSettingsStore } from '@/stores/modules/metrics-settings'

/**
 * Metrics settings store facade.
 *
 * @returns {Object}
 */
export const useMetricsSettings = () => {
  const metricsSettingsStore = useMetricsSettingsStore()

  return {
    fetchMetricsSettingsWithoutStore: metricsSettingsStore.fetchItemWithoutStore,
    updateMetricsSettings: metricsSettingsStore.update,
  }
}
