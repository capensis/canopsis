import { API_ROUTES } from '@/config';

import request from '@/services/request';

const types = {
  FETCH_CURRENT_COMPLETED: 'FETCH_CURRENT_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    current: {},
  },
  getters: {
    current: state => state.current,
  },
  mutations: {
    [types.FETCH_CURRENT_COMPLETED]: (state, current) => {
      state.current = current;
    },
  },
  actions: {
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.eventsRecord, { params });
    },

    fetchCurrent({ commit }) {
      const current = request.get(API_ROUTES.eventsRecordCurrent);

      commit(types.FETCH_CURRENT_COMPLETED, current);
    },

    fetchEventsListWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.eventsRecord}/${id}`, { params });
    },

    createExport(context, { id, eventIds = [] } = {}) {
      return request.post(`${API_ROUTES.eventsRecord}/${id}/exports`, { event_ids: eventIds });
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

    playback(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.eventsRecord}/${id}/playback`, data);
    },

    stopPlayback(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord}/${id}/playback`);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord}/${id}`);
    },

    removeEvent(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecordEvent}/${id}`);
    },
  },
};
