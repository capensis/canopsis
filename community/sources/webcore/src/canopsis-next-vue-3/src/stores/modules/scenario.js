import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useScenarioStore = createCRUDStore({
  id: 'scenario',
  route: API_ROUTES.scenario.scenarios,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    createTestScenarioExecution({ data } = {}) {
      return request.post(API_ROUTES.scenario.testExecution, data);
    },

    fetchTestScenarioExecutionWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.scenario.testExecution}/${id}`);
    },

    fetchTestScenarioExecutionWebhooksResponse({ id } = {}) {
      return request.get(`${API_ROUTES.scenario.testExecutionWebhooks}/${id}/response`);
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.scenario.bulkEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.scenario.bulkDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.scenario.bulk, { data });
    },
  },
});
