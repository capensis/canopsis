import { storeToRefs } from 'pinia';

import { useAnomalyMonitoredConnectorStore } from '@/stores/modules/anomaly-monitored-connector';

/**
 * Anomaly monitored connector store facade.
 *
 * @returns {Object}
 */
export const useAnomalyMonitoredConnector = () => {
  const anomalyMonitoredConnectorStore = useAnomalyMonitoredConnectorStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(anomalyMonitoredConnectorStore);

  return {
    anomalyMonitoredConnectors: items,
    anomalyMonitoredConnectorsMeta: meta,
    anomalyMonitoredConnectorsPending: pending,
    anomalyMonitoredConnectorsFetchingParams: fetchingParams,

    fetchAnomalyMonitoredConnectorsList: anomalyMonitoredConnectorStore.fetchList,
    fetchAnomalyMonitoredConnectorsListWithPreviousParams:
      anomalyMonitoredConnectorStore.fetchListWithPreviousParams,
    fetchAnomalyMonitoredConnectorsListWithoutStore:
      anomalyMonitoredConnectorStore.fetchListWithoutStore,
    updateAnomalyMonitoredConnector: anomalyMonitoredConnectorStore.update,
    updateAnomalyMonitoredConnectorEnabled: anomalyMonitoredConnectorStore.updateEnabled,
    removeAnomalyMonitoredConnector: anomalyMonitoredConnectorStore.remove,
    bulkRemoveAnomalyMonitoredConnectors: anomalyMonitoredConnectorStore.bulkRemove,
    bulkEnableAnomalyMonitoredConnectors: anomalyMonitoredConnectorStore.bulkEnable,
    bulkDisableAnomalyMonitoredConnectors: anomalyMonitoredConnectorStore.bulkDisable,
    fetchAnomalyMonitoredConnectorStatesListWithoutStore:
      anomalyMonitoredConnectorStore.fetchStatesListWithoutStore,
    fetchAnomalyMonitoredConnectorHistoryWithoutStore:
      anomalyMonitoredConnectorStore.fetchConnectorHistoryWithoutStore,
  };
};
