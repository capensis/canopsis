/**
 * @typedef {Object} Icon
 * @property {string} title
 * @property {boolean} fill_border
 * @property {string} [_id]
 * @property {string} [content]
 */

/**
 * @typedef {Object} IconForm
 * @property {string} title
 * @property {boolean} fill_border
 * @property {File} [file]
 */

/**
 * @typedef {IconForm} IconRequest
 */

/**
 * Convert icon to form
 *
 * @param {Icon} [icon = {}]
 * @return {IconForm}
 */
export const iconToForm = (icon = {}) => ({
  title: icon.title ?? '',
  fill_border: icon.fill_border ?? false,
  file: null,
});

/**
 * Convert form to request
 *
 * @param {IconForm} form
 * @return {IconRequest}
 */
export const formToRequest = ({ file, ...rest }) => (
  file
    ? { file, ...rest }
    : rest
);
