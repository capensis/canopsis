import { getCurrentInstance, watch } from 'vue';

/**
 * Hook for get current component instance
 *
 * @return {Vue}
 */
export const useComponentInstance = () => {
  const instance = getCurrentInstance();

  return instance?.proxy || instance;
};

/**
 * Hook for get current component options
 *
 * @return {Object}
 */
export const useComponentOptions = () => {
  const instance = useComponentInstance();

  return instance.$options;
};

/**
 * Hook for get current component model
 *
 * @return {Object}
 */
export const useComponentModel = () => {
  const { model } = useComponentOptions();

  return {
    event: model?.event ?? 'input',
    prop: model?.prop ?? 'value',
  };
};

/**
 * Creates a watcher that stops automatically when the comparator returns true.
 *
 * @param {import('vue').WatchSource|import('vue').WatchSource[]} value - Watcher source(s) to observe
 * @param {Function} [comparator=(v) => v] - Function (newValue, oldValue) => boolean. When true, unwatches
 * @param {import('vue').WatchOptions} [options={}] - Vue watch options (e.g. immediate, deep)
 * @return {Function} Cleanup function that stops the watcher
 */
export const watchOnce = (value, comparator = v => v, options = {}) => {
  const unwatch = watch(value, (newValue, oldValue) => {
    if (comparator(newValue, oldValue)) {
      unwatch();
    }
  }, options);

  return () => unwatch();
};
