/**
 * @typedef {Object} WidgetKioskParameters
 * @property {boolean} hideActions
 * @property {boolean} hideMassSelection
 * @property {boolean} hideToolbar
 */

/**
 * Convert alarm list widget kiosk parameters to form
 *
 * @param {WidgetKioskParameters} [parameters = {}]
 * @returns {WidgetKioskParameters}
 */
export const kioskParametersToForm = (parameters = {}) => ({
  hideActions: parameters.hideActions ?? false,
  hideMassSelection: parameters.hideMassSelection ?? false,
  hideToolbar: parameters.hideToolbar ?? false,
});
