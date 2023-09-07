/**
 * @typedef {Object} TextWidgetParameters
 * @property {string} template
 */

/**
 * Convert text widget parameters to form
 *
 * @param {TextWidgetParameters} parameters
 * @return {TextWidgetParameters}
 */
export const textWidgetParametersToForm = (parameters = {}) => ({
  template: parameters.template ?? '',
});
