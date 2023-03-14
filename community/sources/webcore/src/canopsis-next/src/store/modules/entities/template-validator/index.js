import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    validateDeclareTicketRulesVariables(context, { data }) {
      return request.post(API_ROUTES.templateValidator.declareTicketRules, data);
    },

    validateScenariosVariables(context, { data }) {
      return request.post(API_ROUTES.templateValidator.scenarios, data);
    },
  },
};
