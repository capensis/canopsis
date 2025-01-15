import { unref, onMounted, onBeforeUnmount } from 'vue';

import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';

/**
 * Creates a polling mechanism to fetch events record current data at a specified interval.
 *
 * @param {number} [interval = 10000] - The interval in milliseconds at which to fetch the data.
 * @returns {Object} An object with a method to start polling for events record current data.
 */
export const useEventRecordCurrentPolling = (interval = 10000) => {
  let timer = null;

  /**
   * Fetches events record current data using the `fetchEventsRecordCurrent` method.
   */
  const { fetchEventsRecordCurrent } = useEventsRecordCurrent();

  /**
   * Starts the polling mechanism by fetching data at the specified interval.
   */
  const startPolling = async () => {
    clearTimeout(timer);

    await fetchEventsRecordCurrent();

    timer = setTimeout(startPolling, unref(interval));
  };

  /**
   * Stops the polling mechanism.
   */
  const stopPolling = () => clearTimeout(timer);

  onMounted(startPolling);
  onBeforeUnmount(stopPolling);

  return {
    fetchEventsRecordCurrent: startPolling,
  };
};
