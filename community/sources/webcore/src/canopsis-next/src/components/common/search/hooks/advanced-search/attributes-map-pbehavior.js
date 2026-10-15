import {
  PBEHAVIOR_FIELDS_TO_LABELS_KEYS,
  PBEHAVIOR_GROUPED_ADVANCED_SEARCH_FIELDS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchPbehaviorAttributes } from './attributes-map';
import { useAdvancedSearchGroupedAttributes } from './basic';

export const usePbehaviorAdvancedSearchAttributes = () => {
  const { t } = useI18n();
  const { attributesMap } = useAdvancedSearchPbehaviorAttributes();
  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? t(PBEHAVIOR_FIELDS_TO_LABELS_KEYS[field]),
    groups: PBEHAVIOR_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    attributes,
  };
};
