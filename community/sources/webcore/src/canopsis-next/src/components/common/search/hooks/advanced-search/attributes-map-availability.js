import { computed } from 'vue';

import { AVAILABILITY_GROUPED_ADVANCED_SEARCH_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

import { useEntityInfosKeys, useAdvancedSearchGroupedAttributes } from './basic';
import { useAdvancedSearchEntityAttributes } from './attributes-map';

/**
 * Hook to manage advanced search attributes for entities for availability.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @returns {Object} An object containing the pending state and the computed attributes.
 */
export const useAvailabilityAdvancedSearchAttributes = () => {
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

  const wholePending = computed(() => infosPending.value || entityInfoPropertyPending.value);

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap: entityAttributesMap,
    getText: (field, attrs) => attrs.text ?? tc(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: AVAILABILITY_GROUPED_ADVANCED_SEARCH_FIELDS,
    aliases: entityInfoPropertiesWithAlias,
  });

  return {
    pending: wholePending,
    attributes,
  };
};
