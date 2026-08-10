import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useRemediationInstructionStatsStore = createCRUDStore({
  id: 'remediationInstructionStats',
  route: API_ROUTES.remediation.instructionStats,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    fetchSummaryWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/summary`, { params });
    },

    fetchChangesWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/changes`, { params });
    },

    fetchExecutionsWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/executions`, { params });
    },

    fetchCommentsWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.remediation.instructionComments}/${id}`, { params });
    },
  },
});
