import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    validateEntityServices(context, { data }) {
      return request.post(API_ROUTES.templateValidation.entityServices, data);
    },

    validateEventFilters(context, { data }) {
      return request.post(API_ROUTES.templateValidation.eventFilters, data);
    },

    validateScenarios(context, { data }) {
      return request.post(API_ROUTES.templateValidation.scenarios, data);
    },

    validateLinkRules(context, { data }) {
      return request.post(API_ROUTES.templateValidation.linkRules, data);
    },

    validateWidgets(context, { data }) {
      return request.post(API_ROUTES.templateValidation.widgets, data);
    },

    validateDeclareTicketRules(context, { data }) {
      return request.post(API_ROUTES.templateValidation.declareTicketRules, data);
    },

    validateDynamicInfos(context, { data }) {
      return request.post(API_ROUTES.templateValidation.dynamicInfos, data);
    },

    validateInstructions(context, { data }) {
      return request.post(API_ROUTES.templateValidation.instructions, data);
    },

    validateJobs(context, { data }) {
      return request.post(API_ROUTES.templateValidation.jobs, data);
    },

    validateMetaAlarmRules(context, { data }) {
      return request.post(API_ROUTES.templateValidation.metaAlarmRules, data);
    },
  },
};
