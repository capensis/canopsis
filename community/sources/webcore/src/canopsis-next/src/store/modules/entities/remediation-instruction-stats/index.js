import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.remediation.instructionStats,
  entityType: ENTITIES_TYPES.remediationInstructionStats,
  dataPreparer: d => d.data,
  withMeta: true,
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
