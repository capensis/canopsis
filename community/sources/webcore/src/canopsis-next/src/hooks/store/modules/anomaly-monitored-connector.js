import { useStoreModuleHooks } from '@/hooks/store';

const useAnomalyMonitoredConnectorStoreModule = () => useStoreModuleHooks('anomalyMonitoredConnector');

/**
 * Custom hook for accessing actions related to anomaly monitored connectors from Vuex store module.
 *
 * @returns {Object} Actions to create, update, remove and fetch anomaly monitored connectors list.
 */
export const useAnomalyMonitoredConnectors = () => {
  const { useActions } = useAnomalyMonitoredConnectorStoreModule();

  const actions = useActions({
    createAnomalyMonitoredConnector: 'create',
    updateAnomalyMonitoredConnector: 'update',
    removeAnomalyMonitoredConnector: 'remove',
    updateAnomalyMonitoredConnectorEnabled: 'updateEnabled',
    fetchAnomalyMonitoredConnectorsListWithoutStore: 'fetchListWithoutStore',
    fetchAnomalyMonitoredConnectorStatesListWithoutStore: 'fetchStatesListWithoutStore',
    fetchAnomalyMonitoredConnectorHistoryWithoutStore: 'fetchConnectorHistoryWithoutStore',
  });

  return {
    ...actions,
  };
};
