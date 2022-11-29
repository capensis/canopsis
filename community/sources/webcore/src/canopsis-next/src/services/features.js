import { mergeWith, get, has, flow, uniq, isArray, isFunction } from 'lodash';

class Features {
  constructor() {
    const features = require.context('../features/', true, /index\.js$/);

    Object.values(features
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
      }, {}))
      .map(({ key }) => features(key).default)
      .reduce((acc, plugin) => mergeWith(acc, plugin, (objValue, srcValue) => {
        if (isFunction(objValue) && isFunction(srcValue)) {
          return flow([objValue, srcValue]);
        }

        if (isArray(objValue)) {
          return uniq(objValue.concat(srcValue));
        }

        if (get(objValue, '__esModule') || get(srcValue, '__esModule')) {
          return { ...objValue, ...srcValue };
        }

        return undefined;
      }), {});
  }

  /**
   * Get features value by key
   *
   * @param {string} key
   * @param {*} [defaultValue] - default value. We must put it for arrays
   * @returns {*}
   */
  get(key, defaultValue) {
    return get(this.features, key, defaultValue);
  }

  /**
   * Check if features has value for key
   *
   * @param {string} key
   * @returns {boolean}
   */
  has(key) {
    return has(this.features, key);
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
      throw new Error(`Features feature feature in the path = ${key} is not function: ${func}`);
    }

    return func.call(context, ...args);
  }
}

export default new Features();
