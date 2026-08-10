import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useTicketStatusJobStore = createCRUDStore({
  id: 'ticketStatusJob',
  route: API_ROUTES.ticketStatusJobs,
  withFetchingParams: true,
  actions: {
    play({ data } = {}) {
      return request.post(API_ROUTES.bulkTicketStatusJobs.play, data);
    },

    pause({ data } = {}) {
      return request.post(API_ROUTES.bulkTicketStatusJobs.pause, data);
    },
  },
});
