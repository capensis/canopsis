import {
  DYNAMIC_INFO_FIELDS_TO_LABELS_KEYS,
  DYNAMIC_INFO_GROUPED_ADVANCED_SEARCH_FIELDS,
} from '@/constants';
import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchDynamicInfoAttributes } from './attributes-map';
import { useAdvancedSearchGroupedAttributes } from './basic';

export const useDynamicInfoAdvancedSearchAttributes = () => {
  const { t } = useI18n();
  const { attributesMap } = useAdvancedSearchDynamicInfoAttributes();
  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? t(DYNAMIC_INFO_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: DYNAMIC_INFO_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    attributes,
  };
};
