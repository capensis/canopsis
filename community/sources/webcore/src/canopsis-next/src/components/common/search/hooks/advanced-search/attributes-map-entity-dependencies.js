import { ENTITY_DEPENDENCIES_GROUPED_ADVANCED_SEARCH_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchGroupedAttributes } from './basic';
import { useAdvancedSearchEntityAttributes } from './attributes-map';

export const useEntityDependenciesAdvancedSearchAttributes = () => {
  const { tc } = useI18n();

  const { attributesMap: entityDependenciesAttributesMap } = useAdvancedSearchEntityAttributes();

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap: entityDependenciesAttributesMap,
    getText: (field, attrs) => attrs.text ?? tc(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: ENTITY_DEPENDENCIES_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    attributes,
  };
};
