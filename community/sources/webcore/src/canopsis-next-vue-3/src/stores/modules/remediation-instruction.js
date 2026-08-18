import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useRemediationInstructionStore = createCRUDStore({
  id: 'remediationInstruction',
  route: API_ROUTES.remediation.instructions,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    fetchItemWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.remediation.instructions}/${id}`);
    },

    fetchItemApprovalWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.remediation.instructions}/${id}/approval`);
    },

    updateApproval({ id, data } = {}) {
      return request.put(`${API_ROUTES.remediation.instructions}/${id}/approval`, data);
    },

    rateInstruction({ id, data } = {}) {
      return request.put(`${API_ROUTES.remediation.instructions}/${id}/rate`, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.remediation.bulkInstructions, { data });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.remediation.bulkInstructionsEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.remediation.bulkInstructionsDisable, data);
    },
  },
});
