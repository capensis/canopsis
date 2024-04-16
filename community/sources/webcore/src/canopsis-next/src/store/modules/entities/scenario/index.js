import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.scenario.scenarios,
}, {
  actions: {
    createTestScenarioExecution(context, { data }) {
      return request.post(API_ROUTES.scenario.testExecution, data);
    },

    fetchTestScenarioExecutionWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.scenario.testExecution}/${id}`);
    },

    fetchTestScenarioExecutionWebhooksResponse(context, { id }) {
      return request.get(`${API_ROUTES.scenario.testExecutionWebhooks}/${id}/response`);
    },
  },
});
