/**
 * @typedef {Object} MetricsSettings
 * @property {boolean} enabled_manual_instructions
 * @property {boolean} enabled_not_acked_metrics
 */

/**
 * @typedef {MetricsSettings} MetricsSettingsForm
 */

/**
 * Convert data storage object to data storage form
 *
 * @param {MetricsSettings} metricsSettings
 * @return {MetricsSettingsForm}
 */
export const metricsSettingsToForm = (metricsSettings = {}) => ({
  enabled_manual_instructions: metricsSettings.enabled_manual_instructions ?? false,
  enabled_not_acked_metrics: metricsSettings.enabled_not_acked_metrics ?? false,
});
