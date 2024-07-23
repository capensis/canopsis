import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.eventsRecord, { params });
    },

    fetchCurrentWithoutStore() {
      return request.get(API_ROUTES.eventsRecordCurrent);
    },

    fetchEventsListWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.eventsRecord}/${id}`, { params });
    },

    start(context, { id } = {}) {
      return request.post(`${API_ROUTES.eventsRecord}/${id}/playback`);
    },

    stop(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord}/${id}/playback`);
    },

    remove() {
      // return request.remove(`${API_ROUTES.eventsRecording}/${id}`);
    },
  },
};
