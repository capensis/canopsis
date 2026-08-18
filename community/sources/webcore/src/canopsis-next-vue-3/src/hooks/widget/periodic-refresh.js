import { computed, unref } from 'vue';

import { WIDGET_PERIODIC_REFRESH_KEY } from '@/constants';

import { durationToSeconds } from '@/helpers/date/duration';

import {
  usePeriodicRefreshBasicWatchers,
  usePeriodicRefreshChild,
  usePeriodicRefreshParent,
} from '@/hooks/periodic-refresh';

/**
 * Hook to set up periodic refresh for a widget based on options and a handler.
 *
 * @param {Object} config - The options for the periodic refresh.
 * @param {Object|import('vue').Ref<Object>} config.options - Periodic refresh settings.
 * @param {Function} config.handler - The function to be called on each refresh interval.
 */
export const useWidgetPeriodicRefresh = ({ options, handler }) => {
  // TODO: may be we should remove it later
  const { periodicRefresh } = usePeriodicRefreshParent({
    handler,
    key: WIDGET_PERIODIC_REFRESH_KEY,
  });

  const enabled = computed(() => !!unref(options)?.enabled);
  const delay = computed(() => (unref(options) ? durationToSeconds(unref(options)) : 0));

  let timer;

  const stopRefresh = () => {
    if (!timer) {
      return;
    }

    clearTimeout(timer);

    timer = null;
  };

  const startRefresh = () => {
    stopRefresh();

    if (!enabled.value) {
      return;
    }

    timer = setTimeout(tick, delay.value * 1000);
  };

  const tick = async () => {
    try {
      stopRefresh();
      await periodicRefresh.notify();
      startRefresh();
    } catch (err) {
      console.error('Error during periodic refresh:', err);
    }
  };

  usePeriodicRefreshChild({ handler: tick });

  usePeriodicRefreshBasicWatchers({
    enabled,
    delay,
    start: startRefresh,
    reset: startRefresh,
    stop: startRefresh,
  });
};
