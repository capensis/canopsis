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

    createExport(context, { id } = {}) {
      return request.post(`${API_ROUTES.eventsRecord}/${id}/exports`, { event_ids: [] }); // TODO: remove event_ids
    },

    fetchExport(context, { id } = {}) {
      return request.get(`${API_ROUTES.eventsRecordExport}/${id}`);
    },

    start(context, { data } = {}) {
      return request.post(API_ROUTES.eventsRecordCurrent, data);
    },

    stop() {
      return request.delete(API_ROUTES.eventsRecordCurrent);
    },

    remove() {
      // return request.remove(`${API_ROUTES.eventsRecording}/${id}`);
    },
  },
};
