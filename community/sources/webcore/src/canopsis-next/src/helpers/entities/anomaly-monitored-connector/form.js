/**
 * Convert anomaly monitored connector to form
 *
 * @param {Object} [connector={}] - Anomaly monitored connector object
 * @returns {Object} Form object
 */
export const anomalyMonitoredConnectorToForm = (connector = {}) => ({
  name: connector.name ?? '',
  enabled: connector.enabled ?? true,
});

/**
 * Convert form to anomaly monitored connector
 *
 * @param {Object} form - Form object
 * @returns {Object} Anomaly monitored connector object
 */
export const formToAnomalyMonitoredConnector = form => ({
  name: form.name,
  enabled: form.enabled,
});
