import { computed, unref } from 'vue';

import {
  ADVANCED_SEARCH_GROUPS,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_ADVANCED_SEARCH_ENTITY_FIELDS,
  ALARM_ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
  ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES,
  ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_GROUPED_ADVANCED_SEARCH_FIELDS,
  PBEHAVIOR_PATTERN_PREFIX,
} from '@/constants';

import {
  isAlarmPatternField,
  isAlarmEntityPatternField,
  isPbehaviorPatternField,
} from '@/helpers/search/advanced-search';

import { useI18n } from '@/hooks/i18n';
import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

import { useEntityInfosKeys, useAdvancedSearchGroupedAttributes } from './basic';
import {
  useAdvancedSearchAlarmAttributes,
  useAdvancedSearchEntityAttributes,
  useAdvancedSearchPbehaviorAttributes,
} from './attributes-map';

/**
 * Hook to manage advanced search attributes for alarms, entities, and pbehaviors for alarms list.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {Array} options.rules - The array of rules to be evaluated.
 * @returns {Object} An object containing the pending state, the computed attributes, and the computed flags.
 */
export const useAlarmAdvancedSearchAttributes = ({ rules }) => {
  const { tc } = useI18n();

  const {
    pending: infosPending,
    alarmItems: alarmInfosItems,
    entityItems: entityInfosItems,
  } = useEntityInfosKeys();

  const {
    entityInfoPropertiesWithAlias,
    entityInfoPropertyPending,
  } = useEntityInfoProperty();

  const { attributesMap: alarmAttributesMap } = useAdvancedSearchAlarmAttributes({ infosItems: alarmInfosItems });
  const { attributesMap: entityAttributesMap } = useAdvancedSearchEntityAttributes({
    infosItems: entityInfosItems,
    prefix: ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity,
  });

  const { attributesMap: pbehaviorAttributesMap } = useAdvancedSearchPbehaviorAttributes({
    prefix: `${ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX}${PBEHAVIOR_PATTERN_PREFIX}`,
  });

  const wholePending = computed(() => infosPending.value || entityInfoPropertyPending.value);

  /**
   * HAS FLAGS
   */
  const hasOr = computed(() => unref(rules).some(({ union }) => union === ADVANCED_SEARCH_UNION_CONDITIONS.or));
  const hasAlarmField = computed(() => unref(rules).some(({ attribute }) => isAlarmPatternField(attribute)));
  const hasEntityField = computed(() => unref(rules).some(({ attribute }) => isAlarmEntityPatternField(attribute)));
  const hasPbehaviorField = computed(() => (
    unref(rules).some(({ attribute }) => (
      isPbehaviorPatternField(attribute, ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX)
    ))
  ));

  /**
   * ALLOW FLAGS
   */
  const allowOr = computed(() => [
    hasEntityField.value,
    hasPbehaviorField.value,
    hasAlarmField.value,
  ].filter(Boolean).length <= 1);

  const allowAlarmFields = computed(() => !hasOr.value || (!hasEntityField.value && !hasPbehaviorField.value));
  const allowEntityFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasPbehaviorField.value));
  const allowPbehaviorFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasEntityField.value));

  const attributesMap = computed(() => ({
    ...alarmAttributesMap.value,
    ...entityAttributesMap.value,
    ...pbehaviorAttributesMap.value,
  }));

  const disallowAlarmFields = computed(() => !allowAlarmFields.value);

  const { attributes } = useAdvancedSearchGroupedAttributes({
    attributesMap,
    getText: (field, attrs) => attrs.text ?? tc(ALARM_FIELDS_TO_LABELS_KEYS[field], 2),
    groups: ALARM_GROUPED_ADVANCED_SEARCH_FIELDS,
    aliases: entityInfoPropertiesWithAlias,
    disabled: disallowAlarmFields,
    extraGroups: [
      {
        headerKey: ADVANCED_SEARCH_GROUPS.entity,
        fields: ALARM_ADVANCED_SEARCH_ENTITY_FIELDS,
        getDisabled: () => !allowEntityFields.value,
      },
      {
        headerKey: ADVANCED_SEARCH_GROUPS.pbehavior,
        fields: ALARM_ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS,
        getDisabled: () => !allowPbehaviorFields.value,
      },
    ],
  });

  return {
    pending: wholePending,
    attributes,
    hasOr,
    hasAlarmField,
    hasEntityField,
    hasPbehaviorField,
    allowOr,
    allowAlarmFields,
    allowEntityFields,
    allowPbehaviorFields,
  };
};
