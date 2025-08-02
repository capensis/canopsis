import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the template validator store module hooks.
 *
 * @returns {Object} Store module hooks for the 'templateValidator' namespace
 */
export const useTemplateValidatorStoreModule = () => useStoreModuleHooks('templateValidator');

/**
 * Hook for template validation operations.
 * Provides validation functions for different rule types including declare ticket rules,
 * scenarios, and event filter rules.
 *
 * @typedef {Object} TemplateValidatorActions
 * @property {Function} validateDeclareTicketRulesVariables - Validates declare ticket rule templates
 * @property {Function} validateScenariosVariables - Validates scenario templates
 * @property {Function} validateEventFilterRulesVariables - Validates event filter rule templates
 *
 * @returns {TemplateValidatorActions} All available template validation actions
 */
export const useTemplateValidator = () => {
  const { useActions } = useTemplateValidatorStoreModule();

  const actions = useActions({
    validateDeclareTicketRulesVariables: 'validateDeclareTicketRulesVariables',
    validateScenariosVariables: 'validateScenariosVariables',
    validateEventFilterRulesVariables: 'validateEventFilterRulesVariables',
  });

  return actions;
};
