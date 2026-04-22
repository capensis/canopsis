import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.anomalyMonitoredConnector,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    /**
     * Update enabled status for an anomaly monitored connector.
     *
     * @param {ActionContext} context
     * @param {Object} options
     * @param {string} options.id - Connector id
     * @param {Object} [options.data] - Data to update (enabled)
     * @returns {Promise<AxiosPromise>}
     */
    updateEnabled(context, { id, data }) {
      return request.patch(
        `${API_ROUTES.anomalyMonitoredConnector}/${encodeURIComponent(id)}`,
        data,
      );
    },

    /**
     * Fetch anomaly monitored connector states (paginated) without storing them.
     *
     * @param {ActionContext} context
     * @param {Object} [options]
     * @returns {Promise<AxiosPromise>}
     */
    fetchStatesListWithoutStore(context, options) {
      return request.get(API_ROUTES.anomalyMonitoredConnectorStates, options);
    },

    /**
     * Fetch time-series history for an anomaly monitored connector.
     *
     * @param {ActionContext} context
     * @param {Object} options
     * @param {string} options.id - Connector id
     * @param {Object} [options.params] - Query params (from, to)
     * @returns {Promise<AxiosPromise>}
     */
    fetchConnectorHistoryWithoutStore(context, { id, params } = {}) {
      return request.get(
        `${API_ROUTES.anomalyMonitoredConnector}/${encodeURIComponent(id)}/history`,
        { params },
      );
    },
  },
});
