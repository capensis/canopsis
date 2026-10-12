import { computed } from 'vue';

import {
  ADVANCED_SEARCH_GROUPS,
  ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
  ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  PBEHAVIOR_PATTERN_PREFIX,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

import { useEntityInfosKeys, useAdvancedSearchGroupedAttributes } from './basic';
import { useAdvancedSearchEntityAttributes, useAdvancedSearchPbehaviorAttributes } from './attributes-map';

/**
 * Hook to manage advanced search attributes for entities, and pbehaviors for context explorer.
 *
 * @returns {Object} An object containing the pending state and the computed attributes.
 */
export const useEntityAdvancedSearchAttributes = () => {
  const { tc } = useI18n();

  const {
    pending: infosPending,
    entityItems: entityInfosItems,
  } = useEntityInfosKeys();

  const {
    entityInfoPropertiesWithAlias,
    entityInfoPropertyPending,
  } = useEntityInfoProperty();

  const { attributesMap: entityAttributesMap } = useAdvancedSearchEntityAttributes({ infosItems: entityInfosItems });
  const { attributesMap: pbehaviorAttributesMap } = useAdvancedSearchPbehaviorAttributes({
    prefix: PBEHAVIOR_PATTERN_PREFIX,
  });

  const wholePending = computed(() => infosPending.value || entityInfoPropertyPending.value);

  const attributesMap = computed(() => ({
    ...entityAttributesMap.value,
    ...pbehaviorAttributesMap.value,
  }));

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? tc(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS,
    aliases: entityInfoPropertiesWithAlias,
    extraGroups: [
      {
        headerKey: ADVANCED_SEARCH_GROUPS.pbehavior,
        fields: ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
      },
    ],
  });

  return {
    pending: wholePending,
    attributes,
  };
};
