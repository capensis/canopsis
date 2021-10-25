import { isFunction } from 'lodash';
import Vuex from 'vuex';

/**
 * Create mocked store module.
 *
 * @example
 *  createMockedStoreModule('info', {
 *    getters: {
 *      allowChangeSeverityToInfo: true,
 *      timezone: () => 'Timezone'
 *    },
 *    actions: {
 *      fetchAppInfo: jest.fn()
 *    }
 *  })
 *
 * @param {string} name
 * @param {Object.<string, Function | Mock>} [actions]
 * @param {Object.<string, any>} [getters]
 * @returns {Store}
 */
export const createMockedStoreModule = (name, { actions = {}, getters = {} } = {}) => new Vuex.Store({
  modules: {
    [name]: {
      namespaced: true,
      actions,
      getters: Object
        .entries(getters)
        .reduce((acc, [getterName, getterOrValue]) => {
          acc[getterName] = isFunction(getterOrValue)
            ? getterOrValue
            : () => getterOrValue;

          return acc;
        }, {}),
    },
  },
});

/**
 * Wrapper for createMockedStoreModule, for mock getters.
 *
 * @param {string} name
 * @param {Object.<string, any>} getters
 * @returns {Store}
 */
export const createMockedStoreGetters = (name, getters) => createMockedStoreModule(name, { getters });
