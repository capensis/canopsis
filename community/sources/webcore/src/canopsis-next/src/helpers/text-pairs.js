import uid from './uid';

/**
 * @typedef {Object} TextPairObject
 *
 * @property {string} text
 * @property {string} value
 * @property {string} key
 */

/**
 * Convert text pair to object
 *
 * @param {TextPairObject} [textPair={}]
 * @returns {TextPairObject}
 */
export const textPairToForm = (textPair = {}) => ({
  key: textPair.key ?? uid(),
  text: textPair.text ?? '',
  value: textPair.value ?? '',
});

/**
 * Convert text pairs array to object.
 * Example: [{ key: 'uid', text: 'Content-Type', value: 'application/json' }] => { 'Content-Type': 'application/json' }
 *
 * @param {TextPairObject[]} textPairs - text pairs array
 * @returns {Object}
 */
export const textPairsToObject = textPairs => textPairs.reduce((acc, { text, value }) => {
  acc[text] = value;

  return acc;
}, {});

/**
 * Convert object to text pairs array.
 * Example: { 'Content-Type': 'application/json' } => [{ key: 'uid', text: 'Content-Type', value: 'application/json' }]
 *
 * @param {Object} object
 */
export const objectToTextPairs = (object = {}) => Object.keys(object)
  .map(text => ({ key: uid(), text, value: object[text] }));
