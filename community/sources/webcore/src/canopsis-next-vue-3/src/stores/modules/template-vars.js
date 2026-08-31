import { defineStore } from 'pinia';
import { ref } from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useTemplateVarsStore = defineStore('templateVars', () => {
  const items = ref({});
  const pending = ref(false);

  const fetchList = async (payload) => {
    try {
      pending.value = true;

      items.value = await request.get(API_ROUTES.templateVars, { params: payload?.params });

      return items.value;
    } finally {
      pending.value = false;
    }
  };

  const fetchTemplateVarsList = fetchList;
  const fetchEntityServicesTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecific.entityServices, { params: payload?.params });
  const fetchEventFiltersTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecific.eventFilters, { params: payload?.params });
  const fetchScenariosTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecific.scenarios, { params: payload?.params });
  const fetchLinkRulesTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecific.linkRules, { params: payload?.params });
  const fetchWidgetsTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecific.widgets, { params: payload?.params });
  const fetchDeclareTicketRulesTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.declareTicketRules, { params: payload?.params });
  const fetchDynamicInfosTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.dynamicInfos, { params: payload?.params });
  const fetchInstructionsTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.instructions, { params: payload?.params });
  const fetchJobsTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.jobs, { params: payload?.params });
  const fetchMetaAlarmRulesTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.metaAlarmRules, { params: payload?.params });
  const fetchWebhookTokenRulesTemplateVarsWithoutStore = payload => request.get(API_ROUTES.templateVarsSpecificCat.webhookTokenRules, { params: payload?.params });

  return {
    items,
    pending,

    fetchList,
    fetchTemplateVarsList,
    fetchEntityServicesTemplateVarsWithoutStore,
    fetchEventFiltersTemplateVarsWithoutStore,
    fetchScenariosTemplateVarsWithoutStore,
    fetchLinkRulesTemplateVarsWithoutStore,
    fetchWidgetsTemplateVarsWithoutStore,
    fetchDeclareTicketRulesTemplateVarsWithoutStore,
    fetchDynamicInfosTemplateVarsWithoutStore,
    fetchInstructionsTemplateVarsWithoutStore,
    fetchJobsTemplateVarsWithoutStore,
    fetchMetaAlarmRulesTemplateVarsWithoutStore,
    fetchWebhookTokenRulesTemplateVarsWithoutStore,
  };
});
