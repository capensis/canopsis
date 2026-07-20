import { useTechMetricsStore } from '@/stores/modules/tech-metrics'

/**
 * Tech metrics store facade.
 *
 * @returns {Object}
 */
export const useTechMetrics = () => {
  const techMetricsStore = useTechMetricsStore()

  return {
    createTechMetricsExport: techMetricsStore.createTechMetricsExport,
    fetchTechMetricsExport: techMetricsStore.fetchTechMetricsExport,
    fetchTechMetricsSettings: techMetricsStore.fetchTechMetricsSettings,
    updateTechMetricsSettings: techMetricsStore.updateTechMetricsSettings,
  }
}
