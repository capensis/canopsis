import { useTemplateTestingStore } from '@/stores/modules/template-testing';

/**
 * Template testing store facade.
 *
 * @returns {Object}
 */
export const useTemplateTesting = () => {
  const templateTestingStore = useTemplateTestingStore();

  return {
    fetchTemplateTestingDataListWithoutStore: templateTestingStore.fetchTemplateTestingDataListWithoutStore,
    fetchTemplateTestingDataItemWithoutStore: templateTestingStore.fetchTemplateTestingDataItemWithoutStore,
    createTemplateTestingData: templateTestingStore.createTemplateTestingData,
    updateTemplateTestingData: templateTestingStore.updateTemplateTestingData,
    removeTemplateTestingData: templateTestingStore.removeTemplateTestingData,

    fetchTemplateTestingTestsListWithoutStore: templateTestingStore.fetchTemplateTestingTestsListWithoutStore,
    fetchTemplateTestingTestWithoutStore: templateTestingStore.fetchTemplateTestingTestWithoutStore,
    createTemplateTestingTest: templateTestingStore.createTemplateTestingTest,
    updateTemplateTestingTest: templateTestingStore.updateTemplateTestingTest,
    removeTemplateTestingTest: templateTestingStore.removeTemplateTestingTest,

    validateEntityServicesVariables: templateTestingStore.validateEntityServicesVariables,
    validateEventFilterRulesVariables: templateTestingStore.validateEventFilterRulesVariables,
    validateLinkRulesVariables: templateTestingStore.validateLinkRulesVariables,
    validateScenariosVariables: templateTestingStore.validateScenariosVariables,
    validateWidgetsVariables: templateTestingStore.validateWidgetsVariables,
    validateDeclareTicketRulesVariables: templateTestingStore.validateDeclareTicketRulesVariables,
    validateDynamicInfosVariables: templateTestingStore.validateDynamicInfosVariables,
    validateInstructionsVariables: templateTestingStore.validateInstructionsVariables,
    validateJobsVariables: templateTestingStore.validateJobsVariables,
    validateMetaAlarmRulesVariables: templateTestingStore.validateMetaAlarmRulesVariables,
    validateWebhookTokenRulesVariables: templateTestingStore.validateWebhookTokenRulesVariables,

    copyEventFiltersVariables: templateTestingStore.copyEventFiltersVariables,
    copyDynamicInfosVariables: templateTestingStore.copyDynamicInfosVariables,
  };
};
