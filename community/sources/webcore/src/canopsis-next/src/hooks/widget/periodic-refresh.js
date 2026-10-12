import { computed, unref } from 'vue';

import { WIDGET_PERIODIC_REFRESH_KEY } from '@/constants';

import { durationToSeconds } from '@/helpers/date/duration';

import {
  usePeriodicRefreshParent,
  usePeriodicRefreshChild,
  usePeriodicRefreshBasicWatchers,
} from '@/hooks/periodic-refresh';

/**
 * Hook to set up periodic refresh for a widget based on options and a handler.
 *
 * @param {Object} options - The options for the periodic refresh.
 * @param {Object} options.enabled - Flag to enable/disable periodic refresh.
 * @param {number} options.value - The value for the refresh interval.
 * @param {DurationUnit} options.unit - The unit of the refresh interval.
 * @param {Function} handler - The function to be called on each refresh interval.
 */
export const useWidgetPeriodicRefresh = ({ options, handler }) => {
  const { periodicRefresh } = usePeriodicRefreshParent({ handler, key: WIDGET_PERIODIC_REFRESH_KEY });

  const unwrappedOptions = unref(options) || {};
  const enabled = computed(() => !!unwrappedOptions.enabled);
  const delay = computed(() => (unwrappedOptions ? durationToSeconds(unwrappedOptions) : 0));

  let timer;

  const stopRefresh = () => {
    if (!timer) {
      return;
    }

    clearTimeout(timer);

    timer = null;
  };

  const tick = async () => {
    try {
      stopRefresh();
      await periodicRefresh.notify();
      // eslint-disable-next-line no-use-before-define
      startRefresh();
    } catch (err) {
      console.error('Error during periodic refresh:', err);
    }
  };

  const startRefresh = () => {
    stopRefresh();

    if (!enabled.value) {
      return;
    }

    timer = setTimeout(tick, delay.value * 1000);
  };

  usePeriodicRefreshChild({ handler: tick });

  /**
   * Periodic refresh watchers and lifecycle listeners
   */
  usePeriodicRefreshBasicWatchers({
    enabled,
    delay,
    start: startRefresh,
    reset: startRefresh,
    stop: startRefresh,
  });
};
