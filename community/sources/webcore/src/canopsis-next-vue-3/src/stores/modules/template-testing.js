import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

const requestList = (route, { params } = {}) => request.get(route, { params });
const requestItem = (route, { id } = {}) => request.get(`${route}/${encodeURIComponent(id)}`);
const createItem = (route, { data } = {}) => request.post(route, data);
const updateItem = (route, { id, data } = {}) => request.put(`${route}/${encodeURIComponent(id)}`, data);
const removeItem = (route, { id } = {}) => request.delete(`${route}/${encodeURIComponent(id)}`);
const validateTemplate = route => ({ data } = {}) => request.post(route, data);

export const useTemplateTestingStore = defineStore('templateTesting', () => {
  const fetchTemplateTestingDataListWithoutStore = payload => requestList(API_ROUTES.templateData, payload);
  const fetchTemplateTestingDataItemWithoutStore = payload => requestItem(API_ROUTES.templateData, payload);
  const createTemplateTestingData = payload => createItem(API_ROUTES.templateData, payload);
  const updateTemplateTestingData = payload => updateItem(API_ROUTES.templateData, payload);
  const removeTemplateTestingData = payload => removeItem(API_ROUTES.templateData, payload);

  const fetchTemplateTestingTestsListWithoutStore = payload => requestList(API_ROUTES.templateTests, payload);
  const fetchTemplateTestingTestWithoutStore = payload => requestItem(API_ROUTES.templateTests, payload);
  const createTemplateTestingTest = payload => createItem(API_ROUTES.templateTests, payload);
  const updateTemplateTestingTest = payload => updateItem(API_ROUTES.templateTests, payload);
  const removeTemplateTestingTest = payload => removeItem(API_ROUTES.templateTests, payload);

  const validateEntityServicesVariables = validateTemplate(API_ROUTES.templateValidation.entityServices);
  const validateEventFilterRulesVariables = validateTemplate(API_ROUTES.templateValidation.eventFilters);
  const validateLinkRulesVariables = validateTemplate(API_ROUTES.templateValidation.linkRules);
  const validateScenariosVariables = validateTemplate(API_ROUTES.templateValidation.scenarios);
  const validateWidgetsVariables = validateTemplate(API_ROUTES.templateValidation.widgets);
  const validateDeclareTicketRulesVariables = validateTemplate(API_ROUTES.templateValidationCat.declareTicketRules);
  const validateDynamicInfosVariables = validateTemplate(API_ROUTES.templateValidationCat.dynamicInfos);
  const validateInstructionsVariables = validateTemplate(API_ROUTES.templateValidationCat.instructions);
  const validateJobsVariables = validateTemplate(API_ROUTES.templateValidationCat.jobs);
  const validateMetaAlarmRulesVariables = validateTemplate(API_ROUTES.templateValidationCat.metaAlarmRules);
  const validateWebhookTokenRulesVariables = validateTemplate(API_ROUTES.templateValidationCat.webhookTokenRules);

  const copyEventFiltersVariables = ({ data } = {}) => request.post(API_ROUTES.copyVars.eventFilters, data);
  const copyDynamicInfosVariables = ({ data } = {}) => request.post(API_ROUTES.copyVarsCat.dynamicInfos, data);

  return {
    fetchTemplateTestingDataListWithoutStore,
    fetchTemplateTestingDataItemWithoutStore,
    createTemplateTestingData,
    updateTemplateTestingData,
    removeTemplateTestingData,

    fetchTemplateTestingTestsListWithoutStore,
    fetchTemplateTestingTestWithoutStore,
    createTemplateTestingTest,
    updateTemplateTestingTest,
    removeTemplateTestingTest,

    validateEntityServicesVariables,
    validateEventFilterRulesVariables,
    validateLinkRulesVariables,
    validateScenariosVariables,
    validateWidgetsVariables,
    validateDeclareTicketRulesVariables,
    validateDynamicInfosVariables,
    validateInstructionsVariables,
    validateJobsVariables,
    validateMetaAlarmRulesVariables,
    validateWebhookTokenRulesVariables,

    copyEventFiltersVariables,
    copyDynamicInfosVariables,
  };
});
