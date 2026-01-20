import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchFlappingRulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.flappingRule, { params });
    },

    fetchIdleRulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.idleRule, { params });
    },

    fetchLinkRulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.linkRule, { params });
    },

    fetchRulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.rule, { params });
    },

    fetchPbehaviorPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.pbehavior, { params });
    },

    fetchAlarmTagPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.alarmTag, { params });
    },

    fetchWidgetFilterPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.widgetFilter, { params });
    },

    fetchServicePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.service, { params });
    },

    fetchStateSettingPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.stateSetting, { params });
    },

    fetchEventfilterPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.eventfilter, { params });
    },

    fetchScenarioPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.scenario, { params });
    },

    fetchMetaalarmrulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.metaalarmrule, { params });
    },

    fetchDeclareTicketRulePatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.declareTicketRule, { params });
    },

    fetchInstructionPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.instruction, { params });
    },

    fetchKpiFilterPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.kpiFilter, { params });
    },

    fetchDynamicInfosPatternFields(context, { params } = {}) {
      return request.get(API_ROUTES.patternFields.dynamicInfos, { params });
    },
  },
};
