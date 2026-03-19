import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.ticketStatusJobs,
  withFetchingParams: true,
}, {
  actions: {
    play(context, { data } = {}) {
      return request.post(API_ROUTES.bulkTicketStatusJobs.play, data);
    },

    pause(context, { data } = {}) {
      return request.post(API_ROUTES.bulkTicketStatusJobs.pause, data);
    },
  },
});
