import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the template validation store module.
 *
 * @returns {Object} An object containing getters and actions for the template validation.
 */
const useTemplateValidationStoreModule = () => useStoreModuleHooks('template/validation');

/**
 * Hook to access template validation store.
 *
 * @returns {Object} An object containing:
 * - Actions to validate templates for different entity types.
 */
export const useTemplateValidation = () => {
  const { useActions } = useTemplateValidationStoreModule();

  const actions = useActions({
    validateEntityServices: 'validateEntityServices',
    validateEventFilters: 'validateEventFilters',
    validateScenarios: 'validateScenarios',
    validateLinkRules: 'validateLinkRules',
    validateWidgets: 'validateWidgets',
    validateDeclareTicketRules: 'validateDeclareTicketRules',
    validateDynamicInfos: 'validateDynamicInfos',
    validateInstructions: 'validateInstructions',
    validateJobs: 'validateJobs',
    validateMetaAlarmRules: 'validateMetaAlarmRules',
    validateWebhookTokenRules: 'validateWebhookTokenRules',
  });

  return {
    ...actions,
  };
};
