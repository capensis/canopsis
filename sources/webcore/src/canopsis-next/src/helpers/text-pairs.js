import uid from './uid';

/**
 * Convert text pairs array to object.
 * Example: [{ key: 'uid', text: 'Content-Type', value: 'application/json' }] => { 'Content-Type': 'application/json' }
 *
 * @param {Array} textPairs - text pairs array
 */
export function textPairsToObject(textPairs) {
  return textPairs.reduce((acc, { text, value }) => {
    acc[text] = value;

    return acc;
  }, {});
}

/**
 * Convert object to text pairs array.
 * Example: { 'Content-Type': 'application/json' } => [{ key: 'uid', text: 'Content-Type', value: 'application/json' }]
 *
 * @param {Object} object
 */
export function objectToTextPairs(object) {
  return Object.keys(object).map(text => ({ key: uid(), text, value: object[text] }));
}

/**
 * Default text pairs item creator
 */
export function defaultTextPairCreator() {
  return { key: uid(), text: '', value: '' };
}

export default {
  textPairsToObject,
  objectToTextPairs,
  defaultTextPairCreator,
};
