/**
 * @typedef {Object} MetricsSettings
 * @property {boolean} enabled_instructions
 * @property {boolean} enabled_not_acked_metrics
 * @property {boolean} enabled_sli_metrics
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
  enabled_instructions: metricsSettings.enabled_instructions ?? false,
  enabled_not_acked_metrics: metricsSettings.enabled_not_acked_metrics ?? false,
  enabled_sli_metrics: metricsSettings.enabled_sli_metrics ?? false,
});
