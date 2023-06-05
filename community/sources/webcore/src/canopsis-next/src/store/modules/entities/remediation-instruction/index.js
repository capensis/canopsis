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
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.remediation.instructions, { params });
    },

    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.remediation.instructions}/${id}`);
    },

    fetchItemApprovalWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.remediation.instructions}/${id}/approval`);
    },

    updateApproval(context, { id, data }) {
      return request.put(`${API_ROUTES.remediation.instructions}/${id}/approval`, data);
    },

    rateInstruction(context, { id, data }) {
      return request.put(`${API_ROUTES.remediation.instructions}/${id}/rate`, data);
    },
  },
});
