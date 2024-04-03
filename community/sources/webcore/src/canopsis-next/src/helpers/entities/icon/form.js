/**
 * @typedef {Object} Icon
 * @property {string} title
 * @property {string} [_id]
 * @property {string} [content]
 */

/**
 * @typedef {Object} IconForm
 * @property {string} title
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
  file: null,
});

/**
 * Convert form to request
 *
 * @param {IconForm} form
 * @return {IconRequest}
 */
export const formToRequest = ({ title, file }) => (
  file
    ? { title, file }
    : { title }
);
