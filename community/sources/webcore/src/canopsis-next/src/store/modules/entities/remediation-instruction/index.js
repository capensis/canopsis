import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.remediation.instructions,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
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
