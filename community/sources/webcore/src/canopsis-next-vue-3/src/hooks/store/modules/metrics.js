import { storeToRefs } from 'pinia'

import { useMetricsStore } from '@/stores/modules/metrics'

/**
 * Metrics store facade.
 *
 * @returns {Object}
 */
export const useMetrics = () => {
  const metricsStore = useMetricsStore()
  const {
    widgetMetrics,
    groupMetrics,
    vectorMetrics,
    availabilityMetrics,
    getListByWidgetId,
    getPendingByWidgetId,
    getGroupMetricsListByWidgetId,
    getGroupMetricsPendingByWidgetId,
    getGroupMetricsMetaByWidgetId,
    getVectorMetricsListByWidgetId,
    getVectorMetricsPendingByWidgetId,
    getVectorMetricsMetaByWidgetId,
    getAvailabilityListByWidgetId,
    getAvailabilityPendingByWidgetId,
    getAvailabilityMetaByWidgetId,
  } = storeToRefs(metricsStore)

  return {
    metricsByWidget: widgetMetrics,
    groupMetricsByWidget: groupMetrics,
    vectorMetricsByWidget: vectorMetrics,
    availabilityMetricsByWidget: availabilityMetrics,
    getMetricsListByWidgetId: getListByWidgetId,
    getMetricsPendingByWidgetId: getPendingByWidgetId,
    getAggregatedMetricsListByWidgetId: getListByWidgetId,
    getAggregatedMetricsPendingByWidgetId: getPendingByWidgetId,
    getGroupMetricsListByWidgetId,
    getGroupMetricsPendingByWidgetId,
    getGroupMetricsMetaByWidgetId,
    getVectorMetricsListByWidgetId,
    getVectorMetricsPendingByWidgetId,
    getVectorMetricsMetaByWidgetId,
    getAvailabilityListByWidgetId,
    getAvailabilityPendingByWidgetId,
    getAvailabilityMetaByWidgetId,
    fetchSliMetricsWithoutStore: metricsStore.fetchSliMetricsWithoutStore,
    fetchRatingMetricsWithoutStore: metricsStore.fetchRatingMetricsWithoutStore,
    fetchAlarmsMetricsWithoutStore: metricsStore.fetchAlarmsMetricsWithoutStore,
    createKpiAlarmExport: metricsStore.createKpiAlarmExport,
    createKpiAlarmAggregateExport: metricsStore.createKpiAlarmAggregateExport,
    createRemediationExport: metricsStore.createRemediationExport,
    createKpiRatingExport: metricsStore.createKpiRatingExport,
    createKpiSliExport: metricsStore.createKpiSliExport,
    fetchMetricExport: metricsStore.fetchMetricExport,
    fetchExternalMetricsWithoutStore: metricsStore.fetchExternalMetricsWithoutStore,
    fetchEntityAlarmsMetricsWithoutStore: metricsStore.fetchEntityAlarmsMetricsWithoutStore,
    fetchEntityAggregateMetricsWithoutStore: metricsStore.fetchEntityAggregateMetricsWithoutStore,
    fetchAggregatedMetrics: metricsStore.fetchAggregatedMetrics,
    fetchAggregatedMetricsList: metricsStore.fetchAggregatedMetrics,
    fetchAggregatedMetricsWithoutStore: metricsStore.fetchAggregatedMetricsWithoutStore,
    fetchGroupMetricsList: metricsStore.fetchGroupMetricsList,
    fetchGroupMetrics: metricsStore.fetchGroupMetrics,
    createGroupMetricsExport: metricsStore.createGroupMetricsExport,
    fetchVectorMetricsList: metricsStore.fetchVectorMetricsList,
    fetchVectorMetrics: metricsStore.fetchVectorMetrics,
    fetchAvailabilityList: metricsStore.fetchAvailabilityList,
    fetchAvailability: metricsStore.fetchAvailability,
    fetchAvailabilityWithoutStore: metricsStore.fetchAvailabilityWithoutStore,
    fetchAvailabilityHistoryWithoutStore: metricsStore.fetchAvailabilityHistoryWithoutStore,
    fetchAvailabilityHistoryExport: metricsStore.fetchAvailabilityHistoryExport,
    createAvailabilityExport: metricsStore.createAvailabilityExport,
    fetchAvailabilityExport: metricsStore.fetchAvailabilityExport,
  }
}
