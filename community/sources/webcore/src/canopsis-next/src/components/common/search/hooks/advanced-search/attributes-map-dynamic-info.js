import { DYNAMIC_INFO_FIELDS_TO_LABELS_KEYS, DYNAMIC_INFO_GROUPED_ADVANCED_SEARCH_FIELDS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchGroupedAttributes } from './basic';
import { useAdvancedSearchDynamicInfoAttributes } from './attributes-map';

/**
 * Hook to manage advanced search attributes for dynamic infos.
 *
 * @returns {Object} An object containing the computed attributes.
 */
export const useDynamicInfoAdvancedSearchAttributes = () => {
  const { tc } = useI18n();

  const { attributesMap: dynamicInfoAttributesMap } = useAdvancedSearchDynamicInfoAttributes();

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap: dynamicInfoAttributesMap,
    getText: (field, attrs) => attrs.text ?? tc(DYNAMIC_INFO_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: DYNAMIC_INFO_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    attributes,
  };
};
