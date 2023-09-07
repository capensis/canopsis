import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.remediation.instructionStats,
}, {
  actions: {
    fetchSummaryWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/summary`, { params });
    },

    fetchChangesWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/changes`, { params });
    },

    fetchExecutionsWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.instructionStats}/${id}/executions`, { params });
    },

    fetchCommentsWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.instructionComments}/${id}`, { params });
    },
  },
});
