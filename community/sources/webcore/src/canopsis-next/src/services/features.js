import { get, flow, uniq, isArray, isFunction, isObject, isUndefined, merge } from 'lodash';

/**
 * Get combined value from features by key
 *
 * @param {Object} features
 * @param {string} key
 * @returns {*}
 */
const getValue = (features, key) => features.reduce((acc, feature) => {
  const value = get(feature, key);

  if (isUndefined(value)) {
    return acc;
  }

  if (isUndefined(acc)) {
    return value;
  }

  if (typeof value !== typeof acc) {
    console.error(`Feature service: Different types of values for '${key}' key: '${acc}', '${value}'`);

    return acc;
  }

  if (isFunction(value)) {
    return flow([acc, value]);
  }

  if (isArray(value)) {
    return uniq(value.concat(acc));
  }

  if (isObject(value)) {
    return merge(acc, value);
  }

  return acc;
}, undefined);

/**
 * @class Features
 */
class Features {
  constructor() {
    const featuresFiles = require.context('../features/', true, /index\.js$/);

    this.cache = {};
    this.features = Object.values(
      featuresFiles
        .keys()
        .reduce((acc, key) => {
          const [, feature, ...rest] = key.split('/');

          if (!acc[feature] || acc[feature].length > rest.length) {
            acc[feature] = {
              key,
              length: rest.length,
            };
          }

          return acc;
        }, {}),
    ).map(({ key }) => featuresFiles(key).default);
  }

  /**
   * Get features value by key
   *
   * @param {string} key
   * @param {*} [defaultValue] - default value. We must put it for arrays
   * @returns {*}
   */
  get(key, defaultValue) {
    if (!isUndefined(this.cache[key])) {
      return this.cache[key];
    }

    const value = getValue(this.features, key);

    if (isUndefined(value)) {
      return defaultValue;
    }

    this.cache[key] = value;

    return value;
  }

  /**
   * Check if features has value for key
   *
   * @param {string} key
   * @returns {boolean}
   */
  has(key) {
    return !isUndefined(this.get(key));
  }

  /**
   * Call features function. If we have several functions for one key we will group it by lodash.flow
   *
   * @param {string} key
   * @param {*} context
   * @param {...*} args
   * @returns {*}
   */
  call(key, context, ...args) {
    const func = this.get(key);

    if (!isFunction(func)) {
      throw new Error(`Feature in the path = ${key} is not function: '${func}'`);
    }

    return func.call(context, ...args);
  }
}

export default new Features();
