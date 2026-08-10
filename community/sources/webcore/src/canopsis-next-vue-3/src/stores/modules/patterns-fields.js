import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const usePatternsFieldsStore = defineStore('patterns-fields', () => {
  const fetchByRoute =
    route =>
    ({ params } = {}) =>
      request.get(route, { params });

  return {
    fetchFlappingRulePatternFields: fetchByRoute(API_ROUTES.patternFields.flappingRule),
    fetchIdleRulePatternFields: fetchByRoute(API_ROUTES.patternFields.idleRule),
    fetchLinkRulePatternFields: fetchByRoute(API_ROUTES.patternFields.linkRule),
    fetchResolveRulePatternFields: fetchByRoute(API_ROUTES.patternFields.resolveRule),
    fetchPbehaviorPatternFields: fetchByRoute(API_ROUTES.patternFields.pbehavior),
    fetchAlarmTagPatternFields: fetchByRoute(API_ROUTES.patternFields.alarmTag),
    fetchWidgetFilterPatternFields: fetchByRoute(API_ROUTES.patternFields.widgetFilter),
    fetchServicePatternFields: fetchByRoute(API_ROUTES.patternFields.service),
    fetchStateSettingPatternFields: fetchByRoute(API_ROUTES.patternFields.stateSetting),
    fetchEventFilterPatternFields: fetchByRoute(API_ROUTES.patternFields.eventFilter),
    fetchScenarioPatternFields: fetchByRoute(API_ROUTES.patternFields.scenario),
    fetchMetaalarmrulePatternFields: fetchByRoute(API_ROUTES.patternFields.metaalarmrule),
    fetchDeclareTicketRulePatternFields: fetchByRoute(API_ROUTES.patternFields.declareTicketRule),
    fetchInstructionPatternFields: fetchByRoute(API_ROUTES.patternFields.instruction),
    fetchKpiFilterPatternFields: fetchByRoute(API_ROUTES.patternFields.kpiFilter),
    fetchDynamicInfosPatternFields: fetchByRoute(API_ROUTES.patternFields.dynamicInfos),
    fetchEventRecordPatternFields: fetchByRoute(API_ROUTES.patternFields.eventRecord),
  };
});
