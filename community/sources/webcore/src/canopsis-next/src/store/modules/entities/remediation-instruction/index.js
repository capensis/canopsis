import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.remediation.instructions,
  entityType: ENTITIES_TYPES.remediationInstruction,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    fetchItemCommentsWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.instructionComments}/${id}`, { params });
    },
  },
});
