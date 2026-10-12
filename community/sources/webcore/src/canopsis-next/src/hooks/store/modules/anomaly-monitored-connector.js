import { useStoreModuleHooks } from '@/hooks/store';

const useAnomalyMonitoredConnectortoreModule = () => useStoreModuleHooks('anomalyMonitoredConnector');

/**
 * Custom hook for accessing actions related to anomaly monitored connectors from Vuex store module.
 *
 * @returns {Object} Actions for anomaly monitored connectors CRUD, bulk enabled toggles, and fetches.
 */
export const useAnomalyMonitoredConnector = () => {
  const { useActions } = useAnomalyMonitoredConnectortoreModule();

  const actions = useActions({
    createAnomalyMonitoredConnector: 'create',
    updateAnomalyMonitoredConnector: 'update',
    removeAnomalyMonitoredConnector: 'remove',
    bulkRemoveAnomalyMonitoredConnectors: 'bulkRemove',
    bulkEnableAnomalyMonitoredConnectors: 'bulkEnable',
    bulkDisableAnomalyMonitoredConnectors: 'bulkDisable',
    updateAnomalyMonitoredConnectorEnabled: 'updateEnabled',
    fetchAnomalyMonitoredConnectorsListWithoutStore: 'fetchListWithoutStore',
    fetchAnomalyMonitoredConnectorStatesListWithoutStore: 'fetchStatesListWithoutStore',
    fetchAnomalyMonitoredConnectorHistoryWithoutStore: 'fetchConnectorHistoryWithoutStore',
  });

  return {
    ...actions,
  };
};
