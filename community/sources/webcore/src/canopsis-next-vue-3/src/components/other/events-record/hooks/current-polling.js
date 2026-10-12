import { onBeforeUnmount, onMounted, unref } from 'vue';

import { useEventsRecord } from '@/hooks/store/modules/events-record';

export const useEventRecordCurrentPolling = (interval = 10000) => {
  let timer = null;

  const { fetchCurrentEventsRecord } = useEventsRecord();

  const startPolling = async () => {
    clearTimeout(timer);

    await fetchCurrentEventsRecord();

    timer = setTimeout(startPolling, unref(interval));
  };

  const stopPolling = () => clearTimeout(timer);

  onMounted(startPolling);
  onBeforeUnmount(stopPolling);

  return {
    fetchCurrentEventsRecord,
  };
};
