import { computed } from 'vue';

import { ENTITY_FIELDS_TO_LABELS_KEYS, ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchEntityAttributes } from './attributes-map';
import { useAdvancedSearchGroupedAttributes } from './basic';

export const useEntityAdvancedSearchAttributes = () => {
  const { t } = useI18n();
  const { attributesMap } = useAdvancedSearchEntityAttributes();
  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? t(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    pending: computed(() => false),
    attributes,
  };
};
