import { computed } from 'vue';

import {
  ADVANCED_SEARCH_GROUPS,
  ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS,
  PBEHAVIOR_PATTERN_PREFIX,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

import {
  useAdvancedSearchEntityAttributes,
  useAdvancedSearchPbehaviorAttributes,
} from './attributes-map';
import { useEntityInfosKeys, useAdvancedSearchGroupedAttributes } from './basic';

export const useEntityAdvancedSearchAttributes = ({
  groups = ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS,
} = {}) => {
  const { t } = useI18n();
  const { pending: infosPending, entityItems: infosItems } = useEntityInfosKeys();
  const { entityInfoPropertiesWithAlias, entityInfoPropertiesPending } = useEntityInfoProperty();

  const { attributesMap: entityAttributesMap } = useAdvancedSearchEntityAttributes({ infosItems });
  const { attributesMap: pbehaviorAttributesMap } = useAdvancedSearchPbehaviorAttributes({
    prefix: PBEHAVIOR_PATTERN_PREFIX,
  });

  const attributesMap = computed(() => ({
    ...entityAttributesMap.value,
    ...pbehaviorAttributesMap.value,
  }));

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? t(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2),
    groups,
    aliases: entityInfoPropertiesWithAlias,
    extraGroups: [
      {
        headerKey: ADVANCED_SEARCH_GROUPS.pbehavior,
        fields: ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
      },
    ],
  });

  return {
    pending: computed(() => infosPending.value || entityInfoPropertiesPending.value),
    attributes,
  };
};
