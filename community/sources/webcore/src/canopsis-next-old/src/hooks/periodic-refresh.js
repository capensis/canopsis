import {
  unref,
  inject,
  onBeforeMount,
  onBeforeUnmount,
  watch,
  onMounted,
} from 'vue';

import { VIEW_PERIODIC_REFRESH_KEY } from '@/constants';

import Observer from '@/services/observer';

import { useObserver } from '@/hooks/observer';

/**
 * Hook to use a parent handler with the periodic refresh.
 *
 * @param {Object} options - Options object.
 * @param {Function} options.handler - The handler function to register for periodic refresh.
 * @param {string} [options.key = VIEW_PERIODIC_REFRESH_KEY] - The key for periodic refresh.
 * @returns {Object} An object containing the viewPeriodicRefresh observer.
 */
export const usePeriodicRefreshParent = ({ handler, key = VIEW_PERIODIC_REFRESH_KEY }) => {
  const { observer: periodicRefresh } = useObserver({ key: unref(key) });

  onBeforeMount(() => periodicRefresh.register(handler));
  onBeforeUnmount(() => periodicRefresh.unregister(handler));

  return {
    periodicRefresh,
  };
};

/**
 * Hook to use a child handler with the periodic refresh.
 *
 * @param {Object} options - Options for the hook.
 * @param {Function} options.handler - The handler function to register.
 * @param {string} [options.key = '$periodicRefresh'] - The key to use for injection.
 * @returns {Object} An object containing the view periodic refresh instance.
 */
export const usePeriodicRefreshChild = ({ handler, key = VIEW_PERIODIC_REFRESH_KEY }) => {
  const periodicRefresh = inject(unref(key), new Observer());

  onBeforeMount(() => periodicRefresh.registerChild(handler));
  onBeforeUnmount(() => periodicRefresh.unregisterChild(handler));

  return {
    periodicRefresh,
  };
};

/**
 * Sets up watchers for periodic refresh functionality.
 *
 * @param {Object} params - The parameters for the periodic refresh.
 * @param {Ref<boolean>} params.enabled - A ref indicating whether the periodic refresh is enabled.
 * @param {Ref<number>} params.delay - A ref representing the delay interval for the periodic refresh.
 * @param {Function} [params.start] - Optional function to start the periodic refresh.
 * @param {Function} [params.reset] - Optional function to reset the periodic refresh.
 * @param {Function} [params.stop] - Optional function to stop the periodic refresh.
 */
export const usePeriodicRefreshBasicWatchers = ({ enabled, delay, start, reset, stop }) => {
  watch(enabled, (value, oldValue) => {
    if (value && !oldValue) {
      start?.();
    } else if (!value && oldValue) {
      stop?.();
    }
  });

  watch(delay, (value, oldValue) => value !== oldValue && reset?.());

  onMounted(() => unref(enabled) && stop?.());
  onBeforeUnmount(stop);
};
