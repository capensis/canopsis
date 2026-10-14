import { LINKED_RULES_COUNT, LINKED_RULES_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for generating linked rules tooltips for table items
 *
 * @param {string} [i18nPrefix='common.linkedRulesTooltip'] - i18n prefix for translation keys
 * @returns {Object} Object containing helper functions for linked rules tooltips
 */
export const useLinkedRulesTooltips = (i18nPrefix = 'common.linkedRulesTooltip') => {
  const { t } = useI18n();

  /**
   * Formats a list of linked rules of a specific type for display in a tooltip
   *
   * @param {Array} [typeLinkedRules=[]] - Array of linked rules of a specific type
   * @param {string} typeLinkedRules[].name - The name of the linked rule
   * @returns {string} HTML string containing list items with rule names, limited to the maximum count
   */
  const getLinkedRulesMessageForType = (typeLinkedRules = []) => {
    const hasMore = typeLinkedRules.length > LINKED_RULES_COUNT;

    const result = typeLinkedRules.slice(0, LINKED_RULES_COUNT);

    if (hasMore) {
      result.push({ name: t(`${i18nPrefix}.andMore`) });
    }

    return result.map(item => `<li>${item.name}</li>`).join('');
  };

  /**
   * Generates a complete message about linked rules
   *
   * Combines messages for different types of linked rules into a single HTML string for display in a tooltip.
   *
   * @param {Object} [linkedRules={}] - Object containing arrays of different types of linked rules
   * @param {Array} [linkedRules.widget=[]] - Array of widget rules linked to the item
   * @param {Array} [linkedRules.eventfilter=[]] - Array of event filter rules linked to the item
   * @param {Array} [linkedRules.linkrule=[]] - Array of link rules linked to the item
   * @param {Array} [linkedRules.scenario=[]] - Array of scenario rules linked to the item
   * @param {Array} [linkedRules.declareticketrule=[]] - Array of declare ticket rules linked to the item
   * @returns {string} HTML string containing formatted messages for all types of linked rules
   */
  const getLinkedRulesMessage = (linkedRules = {}) => {
    let message = '';

    LINKED_RULES_TYPES.forEach((ruleType) => {
      const rulesMessages = getLinkedRulesMessageForType(linkedRules[ruleType]);

      if (rulesMessages) {
        message += t(`common.linkedRulesTooltip.linkedRules.${ruleType}`, { rules: rulesMessages });
      }
    });

    return message;
  };

  return {
    getLinkedRulesMessage,
  };
};
