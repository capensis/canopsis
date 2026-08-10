import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useDeclareTicketRuleStore = createCRUDStore({
  id: 'declareTicketRule',
  route: API_ROUTES.declareTicket.rules,
  bulkRoute: API_ROUTES.declareTicket.bulkRules,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    createTestDeclareTicketExecution({ data } = {}) {
      return request.post(API_ROUTES.declareTicket.testExecution, data);
    },

    fetchDeclareTicketExecutionWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.declareTicket.declareTicketExecution}/${id}`);
    },

    fetchTestDeclareTicketExecutionWebhooksResponse({ id } = {}) {
      return request.get(`${API_ROUTES.declareTicket.testExecutionWebhooks}/${id}/response`);
    },

    bulkCreateDeclareTicketExecution({ data } = {}) {
      return request.post(API_ROUTES.declareTicket.bulkDeclareTicket, data);
    },

    fetchAssignedTicketsWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.declareTicket.alarmsAssigned, { params });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.declareTicket.bulkRulesEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.declareTicket.bulkRulesDisable, data);
    },
  },
});
