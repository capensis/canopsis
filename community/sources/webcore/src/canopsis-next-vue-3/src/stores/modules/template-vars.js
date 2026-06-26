import { ref } from 'vue';
import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useTemplateVarsStore = defineStore('templateVars', () => {
  const items = ref({});
  const pending = ref(true);

  /**
   * Fetch all template variables and save them to the store.
   *
   * @param {Object} [payload]
   * @param {Object} [payload.params]
   */
  const fetchList = async ({ params } = {}) => {
    try {
      items.value = {};
      pending.value = true;

      items.value = await request.get(API_ROUTES.templateVars, params);
    } catch (err) {
      console.warn(err);
    } finally {
      pending.value = false;
    }
  };

  /**
   * Fetch entity services template variables without saving them to the store.
   *
   * @param {Object} [payload]
   * @param {Object} [payload.params]
   * @return {Promise}
   */
  const fetchEntityServicesVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecific.entityServices, { params })
  );

  const fetchEventFiltersVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecific.eventFilters, { params })
  );

  const fetchScenariosVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecific.scenarios, { params })
  );

  const fetchLinkRulesVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecific.linkRules, { params })
  );

  const fetchWidgetsVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecific.widgets, { params })
  );

  const fetchDeclareTicketRulesVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.declareTicketRules, { params })
  );

  const fetchDynamicInfosVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.dynamicInfos, { params })
  );

  const fetchInstructionsVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.instructions, { params })
  );

  const fetchJobsVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.jobs, { params })
  );

  const fetchMetaAlarmRulesVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.metaAlarmRules, { params })
  );

  const fetchWebhookTokenRulesVarsWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.templateVarsSpecificCat.webhookTokenRules, { params })
  );

  return {
    items,
    pending,

    fetchList,
    fetchEntityServicesVarsWithoutStore,
    fetchEventFiltersVarsWithoutStore,
    fetchScenariosVarsWithoutStore,
    fetchLinkRulesVarsWithoutStore,
    fetchWidgetsVarsWithoutStore,
    fetchDeclareTicketRulesVarsWithoutStore,
    fetchDynamicInfosVarsWithoutStore,
    fetchInstructionsVarsWithoutStore,
    fetchJobsVarsWithoutStore,
    fetchMetaAlarmRulesVarsWithoutStore,
    fetchWebhookTokenRulesVarsWithoutStore,
  };
});
