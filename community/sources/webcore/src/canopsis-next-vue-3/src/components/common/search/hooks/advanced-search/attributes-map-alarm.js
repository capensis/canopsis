import { computed } from 'vue';

import { ALARM_FIELDS_TO_LABELS_KEYS, ALARM_GROUPED_ADVANCED_SEARCH_FIELDS } from '@/constants';
import { useI18n } from '@/hooks/i18n';

import { useAdvancedSearchAlarmAttributes } from './attributes-map';
import { useAdvancedSearchGroupedAttributes } from './basic';

export const useAlarmAdvancedSearchAttributes = () => {
  const { t } = useI18n();
  const { attributesMap } = useAdvancedSearchAlarmAttributes();
  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? t(ALARM_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: ALARM_GROUPED_ADVANCED_SEARCH_FIELDS,
  });

  return {
    pending: computed(() => false),
    attributes,
    allowOr: computed(() => true),
  };
};
