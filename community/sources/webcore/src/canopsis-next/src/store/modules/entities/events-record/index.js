import { API_ROUTES } from '@/config';

import request from '@/services/request';

import currentModule from './current';

export default {
  namespaced: true,
  modules: {
    current: currentModule,
  },
  actions: {
    async fetchListWithoutStore({ dispatch }, { params } = {}) {
      const response = await request.get(API_ROUTES.eventsRecord.list, { params });

      dispatch('current/setCurrent', response.status);

      return response;
    },

    fetchEventsListWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.eventsRecord.list}/${id}`, { params });
    },

    createExport(context, { id, eventIds = [] } = {}) {
      return request.post(`${API_ROUTES.eventsRecord.list}/${id}/exports`, { event_ids: eventIds });
    },

    fetchExport(context, { id } = {}) {
      return request.get(`${API_ROUTES.eventsRecord.export}/${id}`);
    },

    playback({ dispatch }, { id, data } = {}) {
      try {
        dispatch('current/setCurrentResending', true);

        return request.post(`${API_ROUTES.eventsRecord.list}/${id}/playback`, data);
      } catch (err) {
        dispatch('current/setCurrentResending', false);

        throw err;
      }
    },

    async stopPlayback({ dispatch }, { id } = {}) {
      const response = await request.delete(`${API_ROUTES.eventsRecord.list}/${id}/playback`);

      dispatch('current/reset');

      return response;
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord.list}/${id}`);
    },

    removeEvent(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord.event}/${id}`);
    },

    bulkRemoveEvent(context, { data }) {
      return request.delete(API_ROUTES.eventsRecord.bulkEvent, { data });
    },
  },
};
