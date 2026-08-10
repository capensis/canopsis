import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';
import request, { useRequestCancelling } from '@/services/request';

const DEFAULT_STATUS = {
  limit: 1,
  recording: [],
  resending: [],
};

export const useEventsRecordStore = defineStore('eventsRecord', () => {
  const currentPending = ref(false);
  const currentStatus = ref({ ...DEFAULT_STATUS });

  const recordings = computed(() => currentStatus.value?.recording ?? []);
  const resendings = computed(() => currentStatus.value?.resending ?? []);
  const limit = computed(() => currentStatus.value?.limit ?? 1);
  const recordingsById = computed(() =>
    Object.fromEntries(recordings.value.map(item => [item._id, item])),
  );
  const resendingsById = computed(() =>
    Object.fromEntries(resendings.value.map(item => [item._id, item])),
  );

  const setCurrentStatus = (status) => {
    currentStatus.value = status ?? { ...DEFAULT_STATUS };
  };

  const fetchListWithoutStore = async ({ params } = {}) => {
    const response = await request.get(API_ROUTES.eventsRecord.list, { params });

    setCurrentStatus(response.status);

    return response;
  };

  const fetchEventsListWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.eventsRecord.list}/${id}`, { params });

  const createExport = ({ id, eventIds = [], eventPattern } = {}) =>
    request.post(`${API_ROUTES.eventsRecord.list}/${id}/exports`, {
      event_ids: eventIds,
      event_pattern: eventPattern,
    });

  const fetchExport = ({ id } = {}) => request.get(`${API_ROUTES.eventsRecord.export}/${id}`);
  const playback = ({ id, data } = {}) =>
    request.post(`${API_ROUTES.eventsRecord.list}/${id}/playback`, data);
  const stopPlayback = ({ id } = {}) =>
    request.delete(`${API_ROUTES.eventsRecord.list}/${id}/playback`);
  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.eventsRecord.list}/${id}`);
  const removeEvent = ({ id } = {}) => request.delete(`${API_ROUTES.eventsRecord.event}/${id}`);
  const bulkRemoveEvent = ({ data } = {}) =>
    request.delete(API_ROUTES.eventsRecord.bulkEvent, { data });

  const fetchCurrent = () =>
    useRequestCancelling(async (source) => {
      try {
        currentPending.value = true;

        const status = await request.get(API_ROUTES.eventsRecord.current, {
          cancelToken: source.token,
        });

        setCurrentStatus(status);

        return status;
      } catch (err) {
        console.warn(err);

        return undefined;
      } finally {
        currentPending.value = false;
      }
    }, 'event-records-current');

  const startCurrent = ({ data } = {}) => request.post(API_ROUTES.eventsRecord.current, data);
  const stopCurrent = ({ id } = {}) => request.delete(`${API_ROUTES.eventsRecord.current}/${id}`);

  return {
    currentPending,
    currentStatus,
    recordings,
    resendings,
    limit,
    recordingsById,
    resendingsById,
    setCurrentStatus,
    fetchListWithoutStore,
    fetchEventsListWithoutStore,
    createExport,
    fetchExport,
    playback,
    stopPlayback,
    remove,
    removeEvent,
    bulkRemoveEvent,
    fetchCurrent,
    startCurrent,
    stopCurrent,
  };
});
