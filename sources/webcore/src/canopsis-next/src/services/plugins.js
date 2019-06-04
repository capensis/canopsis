import { mergeWith, get, has, isArray, isFunction, flow, uniq } from 'lodash';

class Plugins {
  constructor() {
    const plugins = require.context('@/plugins/', true, /config\.js$/);

    this.plugins = plugins.keys().map(key => plugins(key).default).reduce((acc, plugin) =>
      mergeWith(acc, plugin, (objValue, srcValue) => {
        if (isFunction(objValue) && isFunction(srcValue)) {
          return flow([objValue, srcValue]);
        }

        if (isArray(objValue)) {
          return uniq(objValue.concat(srcValue));
        }

        return undefined;
      }), {});
  }

  /**
   * Get plugins value by key
   *
   * @param {string} key
   * @returns {*}
   */
  get(key) {
    return get(this.plugins, key);
  }

  /**
   * Check if plugins has value for key
   *
   * @param {string} key
   * @returns {boolean}
   */
  has(key) {
    return has(this.plugins, key);
  }

  /**
   * Call plugins function. If we have several functions for one key we will group it by lodash.flow
   *
   * @param {string} key
   * @param {*} context
   * @param {...*} args
   * @returns {*}
   */
  call(key, context, ...args) {
    const func = this.get(key);

    if (!isFunction(func)) {
      throw new Error(`Plugins feature feature in the path = ${key} is not function: ${func}`);
    }

    return func.call(context, ...args);
  }
}

export default new Plugins();
