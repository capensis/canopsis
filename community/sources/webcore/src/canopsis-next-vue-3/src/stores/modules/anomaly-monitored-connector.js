import { API_ROUTES } from '@/config'

import request from '@/services/request'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useAnomalyMonitoredConnectorStore = createCRUDStore({
  id: 'anomalyMonitoredConnector',
  route: API_ROUTES.anomalyMonitoredConnector,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    updateEnabled({ id, data } = {}) {
      return request.patch(
        `${API_ROUTES.anomalyMonitoredConnector}/${encodeURIComponent(id)}`,
        data,
      )
    },

    fetchStatesListWithoutStore(options) {
      return request.get(API_ROUTES.anomalyMonitoredConnectorStates, options)
    },

    fetchConnectorHistoryWithoutStore({ id, params } = {}) {
      return request.get(
        `${API_ROUTES.anomalyMonitoredConnector}/${encodeURIComponent(id)}/history`,
        { params },
      )
    },

    bulkRemove({ data = [] } = {}) {
      return request.delete(API_ROUTES.bulkAnomalyMonitoredConnector, {
        data: data.map(({ id }) => ({ _id: id })),
      })
    },

    bulkEnable({ data = [] } = {}) {
      return request.patch(
        API_ROUTES.bulkAnomalyMonitoredConnector,
        data.map(({ id }) => ({ id, enabled: true })),
      )
    },

    bulkDisable({ data = [] } = {}) {
      return request.patch(
        API_ROUTES.bulkAnomalyMonitoredConnector,
        data.map(({ id }) => ({ id, enabled: false })),
      )
    },
  },
})
