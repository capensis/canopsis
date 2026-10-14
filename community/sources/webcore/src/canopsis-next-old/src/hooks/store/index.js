import { computed } from 'vue';

import { useComponentInstance } from '../vue';

/**
 * Normalizes a given map into an array of objects with `key` and `value` properties.
 * If the input is an array, each element is treated as both the key and value.
 * If the input is an object, it converts the entries into an array of objects.
 *
 * @param {Object|Array} map - The map to normalize, which can be either an object or an array.
 * @returns {Array.<{key: string, value: string}>} An array of objects with `key` and `value` properties.
 */
const normalizeStoreMap = (map = {}) => (
  Array.isArray(map)
    ? map.map(key => ({ key, value: key }))
    : Object.entries(map).map(([key, value]) => ({ key, value }))
);

/**
 * Retrieves the Vuex store instance from the current Vue component instance.
 * This function must be used within the Vue component's `setup()` method.
 *
 * @returns {Object} The Vuex store instance.
 * @throws {Error} If the function is not used within the `setup()` method.
 *
 * @example
 * // In a Vue component
 * import { useStore } from `./path/to/useStore`;
 *
 * export default {
 *   setup() {
 *     const store = useStore();
 *     // Now you can use the store to commit mutations, dispatch actions, or access state
 *     return {
 *       // Your other setup properties
 *     };
 *   }
 * }
 */
export const useStore = () => {
  const vm = useComponentInstance();

  if (!vm) {
    throw new Error('You must use this function within the "setup()" method');
  }

  return vm.$store;
};

/**
 * Creates hooks for accessing Vuex store module's getters and actions.
 * This function helps in creating a more modular and organized way to interact with Vuex store modules.
 *
 * @param {string} namespace - The namespace of the Vuex store module.
 * @returns {Object} An object containing the store, module, and functions to access getters and actions.
 * @throws {Error} If the provided namespace does not correspond to any module in the store.
 *
 * @example
 * // Assuming a Vuex store module is registered under `user/` namespace
 * import { useStoreModuleHooks } from `./path/to/useStoreModuleHooks`;
 *
 * export default {
 *   setup() {
 *     const { useGetters, useActions } = useStoreModuleHooks('user/');
 *     const { userName } = useGetters(['userName']);
 *     const { updateUser } = useActions({ customUpdateUser: 'updateUser' });
 *
 *     return {
 *       userName,
 *       customUpdateUser
 *     };
 *   }
 * }
 */
export const useStoreModuleHooks = (namespace) => {
  const store = useStore();
  const preparedNamespace = namespace.slice(-1) === '/' ? namespace : `${namespace}/`;
  const module = store._modulesNamespaceMap[preparedNamespace];

  if (!module) {
    throw new Error(`Incorrect module namespace = ${namespace}`);
  }

  const useGetters = getters => normalizeStoreMap(getters).reduce((acc, { key, value }) => {
    acc[key] = computed(() => store.getters[preparedNamespace + value]);

    return acc;
  }, {});

  const useActions = actions => normalizeStoreMap(actions).reduce((acc, { key, value }) => {
    acc[key] = (...args) => module.context.dispatch.apply(store, [value, ...args]);

    return acc;
  }, {});

  return {
    store,
    module,
    useGetters,
    useActions,
  };
};
