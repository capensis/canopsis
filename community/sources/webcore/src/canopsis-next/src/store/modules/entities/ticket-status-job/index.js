import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.ticketStatusJobs, { params });
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.ticketStatusJobs}/${id}`, data);
    },

    play() {
      return request.post(API_ROUTES.bulkTicketStatusJobs.play);
    },

    pause() {
      return request.post(API_ROUTES.bulkTicketStatusJobs.pause);
    },

    stop() {
      return request.post(API_ROUTES.bulkTicketStatusJobs.stop);
    },
  },
};
