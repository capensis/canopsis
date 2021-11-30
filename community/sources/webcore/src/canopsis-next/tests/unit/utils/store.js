import { isFunction } from 'lodash';
import Vuex from 'vuex';
/**
 * @typedef {Object} Module
 * @property {string} name
 * @property {Object.<string, Function | Mock>} [actions]
 * @property {Object} [state]
 * @property {Object.<string, any>} [getters]
 */

const convertMockedGettersToStore = (getters = {}) => Object
  .entries(getters)
  .reduce((acc, [getterName, getterOrValue]) => {
    acc[getterName] = isFunction(getterOrValue)
      ? getterOrValue
      : () => getterOrValue;

    return acc;
  }, {});
/**
 * Create mocked store module.
 *
 * @example
 *  createMockedStoreModules({
 *    name: 'info',
 *    getters: {
 *      allowChangeSeverityToInfo: true,
 *      timezone: () => 'Timezone'
 *    },
 *    actions: {
 *      fetchAppInfo: jest.fn()
 *    }
 *  })
 *
 * @param {Module[]} modules
 * @returns {Store}
 */
export const createMockedStoreModules = modules => new Vuex.Store({
  modules: modules.reduce((acc, { name, actions = {}, getters, state }) => {
    acc[name] = {
      namespaced: true,
      state,
      actions,
      getters: convertMockedGettersToStore(getters),
    };

    return acc;
  }, {}),
});

/**
 * Wrapper for createMockedStoreModule, for mock getters.
 *
 * @param {string} name
 * @param {Object.<string, any>} getters
 * @returns {Store}
 */
export const createMockedStoreGetters = ({ name, ...getters }) => createMockedStoreModules([{ name, getters }]);
