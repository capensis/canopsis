import sanitizeHTML from 'sanitize-html';
import { isString } from 'lodash';


import { SANITIZE_TEXT_EDITOR_OPTIONS } from '@/config';

/**
 * Sanitize helper with default options
 *
 * @param {string} text
 * @param {Object} [options = SANITIZE_TEXT_EDITOR_OPTIONS]
 */
export function sanitize(text, options = SANITIZE_TEXT_EDITOR_OPTIONS) {
  try {
    return sanitizeHTML(text, options);
  } catch (err) {
    console.warn(err);

    return '';
  }
}

/**
 * Default sanitizer. We use it for maps and etc.
 *
 * @param {string} text
 * @returns {string}
 */
export function defaultSanitizer(text) {
  return text && isString(text) ? sanitize(text) : text;
}

/**
 * Getter for default sanitizer for a array of strings/Objects
 *
 * @param {string} fieldKey
 * @returns {Function}
 */
export function getDefaultSanitizerForArray(fieldKey) {
  return (items) => {
    if (fieldKey) {
      return items.map(item => ({ ...item, [fieldKey]: defaultSanitizer(item[fieldKey]) }));
    }

    return items.map(defaultSanitizer);
  };
}
