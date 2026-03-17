import { PBEHAVIOR_FIELDS_TO_LABELS_KEYS, PBEHAVIOR_GROUPED_ADVANCED_SEARCH_FIELDS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchGroupedAttributes } from './basic';
import { useAdvancedSearchPbehaviorAttributes } from './attributes-map';

/**
 * Hook to manage advanced search attributes for pbehaviors.
 *
 * @returns {Object} An object containing the computed attributes.
 * @property {Array} attributes - Array of pbehavior attributes with headers for advanced search field dropdowns.
 */
export const usePbehaviorAdvancedSearchAttributes = () => {
  const { t } = useI18n();

  const { attributesMap: pbehaviorAttributesMap } = useAdvancedSearchPbehaviorAttributes();

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap: pbehaviorAttributesMap,
    getText: field => t(PBEHAVIOR_FIELDS_TO_LABELS_KEYS[field]),
    groups: PBEHAVIOR_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    attributes,
  };
};
