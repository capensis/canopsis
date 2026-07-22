/**
 * Convert metrics settings object to form.
 *
 * @param {Object} [metricsSettings={}]
 * @returns {Object}
 */
export const metricsSettingsToForm = (metricsSettings = {}) => ({
  enabled_instructions: metricsSettings.enabled_instructions ?? false,
  enabled_not_acked_metrics: metricsSettings.enabled_not_acked_metrics ?? false,
  enabled_sli_metrics: metricsSettings.enabled_sli_metrics ?? false,
})
