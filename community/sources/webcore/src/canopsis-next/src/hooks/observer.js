import { unref, provide, onBeforeUnmount } from 'vue';

import Observer from '@/services/observer';

/**
 * Custom hook to create and provide an Observer instance.
 *
 * This hook initializes an Observer instance, provides it using the given key,
 * and ensures that all handlers are unregistered before the component is unmounted.
 *
 * @param {Object} params - The parameters object.
 * @param {Ref|string} params.key - The key used to provide the Observer instance.
 *                                  It can be a ref or a plain string.
 *
 * @returns {Object} An object containing the Observer instance.
 * @returns {Observer} return.observer - The created Observer instance.
 */
export const useObserver = ({ key }) => {
  const observer = new Observer();

  provide(unref(key), observer);

  onBeforeUnmount(() => observer.unregisterAll());

  return {
    observer,
  };
};
