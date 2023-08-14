import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.declareTicket.rules,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.declareTicket.bulkRules, { data });
    },

    createTestDeclareTicketExecution(context, { data }) {
      return request.post(API_ROUTES.declareTicket.testExecution, data);
    },

    fetchDeclareTicketExecutionWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.declareTicket.declareTicketExecution}/${id}`);
    },

    fetchTestDeclareTicketExecutionWebhooksResponse(context, { id }) {
      return request.get(`${API_ROUTES.declareTicket.testExecutionWebhooks}/${id}/response`);
    },

    bulkCreateDeclareTicketExecution(context, { data }) {
      return request.post(API_ROUTES.declareTicket.bulkDeclareTicket, data);
    },

    fetchAssignedTicketsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.declareTicket.alarmsAssigned, { params });
    },
  },
});
