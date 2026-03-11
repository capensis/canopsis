import { computed, unref } from 'vue';

import {
  ADVANCED_SEARCH_USER_OPERATORS,
  ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
  ALARM_EVENT_INITIATORS,
  ALARM_PATTERN_FIELDS,
  ALARM_STATES,
  ALARM_STATUSES,
  BASIC_ENTITY_TYPES,
  ENTITY_TYPES,
  ENTITY_PATTERN_FIELDS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_DATE_OPERATORS,
  PBEHAVIOR_TYPE_TYPES,
  PBEHAVIOR_FIELDS,
  PATTERN_ALARM_TAG_LABEL_OPERATORS,
  DYNAMIC_INFO_FIELDS,
} from '@/constants';

import { addPrefixToAttributesMap, getNumberMinValueAttributes } from '@/helpers/search/advanced-search';

import { useI18n } from '@/hooks/i18n';
import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';
import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label';
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
import { useUser } from '@/hooks/store/modules/user';
import { useEntityCategory } from '@/hooks/store/modules/entity-category';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { useDynamicInfo } from '@/hooks/store/modules/dynamic-info';

import { useGetEntityOptions } from './basic';

/**
 * Hook to manage advanced search attributes for alarms.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {Array} options.infosItems - An array of information items used in the search.
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for alarms.
 * @property {Function} attributesMap[].fetchValues - Function to fetch values for the attribute.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 * @property {string} attributesMap[].itemValue - The key used to identify the value in fetched data.
 * @property {string} attributesMap[].itemText - The key used to identify the display text in fetched data.
 * @property {Array} [attributesMap[].values] - Predefined values for the attribute, if applicable.
 */
export const useAdvancedSearchAlarmAttributes = ({ infosItems }) => {
  const { t } = useI18n();
  const { fetchAlarmTagsListWithoutStore } = useAlarmTag();
  const { fetchAlarmTagsLabelsListWithoutStore } = useAlarmTagLabel();
  const { fetchMetaAlarmRulesListWithoutStore } = useMetaAlarmRule();
  const { fetchUsersListWithoutStore } = useUser();
  const { getEntityOptions } = useGetEntityOptions();

  const STRING_WITH_EXIST_AND_ONE_OF_OPERATORS = [
    ...ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,

    PATTERN_OPERATORS.exist,
  ];

  const DURATION_OPTIONS = { operators: PATTERN_DURATION_OPERATORS };

  const INITIATOR_OPTIONS = {
    operators: ADVANCED_SEARCH_USER_OPERATORS,
    values: Object.values(ALARM_EVENT_INITIATORS).map(initiator => ({ value: initiator, text: initiator })),
  };

  const attributesMap = computed(() => ({
    /**
     * Basic
     */
    [ALARM_PATTERN_FIELDS.displayName]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.connector]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.connector]),

      itemText: 'connector_type',
      itemValue: 'connector_type',
    },
    [ALARM_PATTERN_FIELDS.connectorName]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.connector]),

      itemValue: 'name',
    },
    [ALARM_PATTERN_FIELDS.component]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.component]),

      itemText: 'component',
    },
    [ALARM_PATTERN_FIELDS.resource]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.resource]),

      itemText: 'name',
      itemValue: 'name',
    },
    [ALARM_PATTERN_FIELDS.state]: {
      operators: PATTERN_NUMBER_OPERATORS,
      values: Object.values(ALARM_STATES).map(value => ({ value, text: t(`common.stateTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ALARM_PATTERN_FIELDS.status]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
      ],
      values: Object.values(ALARM_STATUSES).map(value => ({ value, text: t(`common.statusTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ALARM_PATTERN_FIELDS.tags]: {
      operators: [
        PATTERN_OPERATORS.with,
        PATTERN_OPERATORS.without,
        PATTERN_OPERATORS.withLabel,
        PATTERN_OPERATORS.withoutLabel,
      ],
      fetchValues: ({ params }, { operator } = {}) => (
        PATTERN_ALARM_TAG_LABEL_OPERATORS.includes(operator)
          ? fetchAlarmTagsLabelsListWithoutStore({ params })
          : fetchAlarmTagsListWithoutStore({ params })
      ),
      itemText: ({ operator } = {}) => (
        PATTERN_ALARM_TAG_LABEL_OPERATORS.includes(operator)
          ? '_id'
          : 'value'
      ),
      itemValue: ({ operator } = {}) => (
        PATTERN_ALARM_TAG_LABEL_OPERATORS.includes(operator)
          ? '_id'
          : 'value'
      ),
    },
    [ALARM_PATTERN_FIELDS.infos]: {
      items: unref(infosItems),
    },
    [ALARM_PATTERN_FIELDS.meta]: {
      operators: [
        PATTERN_OPERATORS.isMetaAlarm,
        PATTERN_OPERATORS.isNotMetaAlarm,
        PATTERN_OPERATORS.ruleIs,
      ],
      fetchValues: fetchMetaAlarmRulesListWithoutStore,
      itemText: 'name',
      itemValue: '_id',
    },
    [ALARM_PATTERN_FIELDS.changeState]: {
      operators: PATTERN_EXISTS_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.totalStateChanges]: getNumberMinValueAttributes(1),

    /**
     * Messages
     */
    [ALARM_PATTERN_FIELDS.output]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.longOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.initialOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.initialLongOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.lastComment]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.lastCommentInitiator]: INITIATOR_OPTIONS,

    /**
     * Ticket
     */
    [ALARM_PATTERN_FIELDS.ticketMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.ticketInitiator]: INITIATOR_OPTIONS,
    [ALARM_PATTERN_FIELDS.ticketValue]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.ticket]: {
      operators: [
        PATTERN_OPERATORS.ticketAssociated,
        PATTERN_OPERATORS.ticketNotAssociated,
      ],
    },

    /**
     * Date
     */
    [ALARM_PATTERN_FIELDS.creationDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.lastUpdateDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.lastEventDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.ackAt]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.resolved]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.activationDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.duration]: DURATION_OPTIONS,

    /**
     * Actions
     */
    [ALARM_PATTERN_FIELDS.ack]: {
      operators: [
        PATTERN_OPERATORS.acked,
        PATTERN_OPERATORS.notAcked,
      ],
    },
    [ALARM_PATTERN_FIELDS.ackBy]: {
      operators: ADVANCED_SEARCH_USER_OPERATORS,
      fetchValues: fetchUsersListWithoutStore,
      itemValue: 'display_name',
      itemText: 'display_name',
    },
    [ALARM_PATTERN_FIELDS.ackMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_PATTERN_FIELDS.ackInitiator]: INITIATOR_OPTIONS,
    [ALARM_PATTERN_FIELDS.canceled]: {
      operators: [
        PATTERN_OPERATORS.canceled,
        PATTERN_OPERATORS.notCanceled,
      ],
    },
    [ALARM_PATTERN_FIELDS.canceledInitiator]: INITIATOR_OPTIONS,
    [ALARM_PATTERN_FIELDS.activated]: {
      operators: [
        PATTERN_OPERATORS.activated,
        PATTERN_OPERATORS.inactive,
      ],
    },
    [ALARM_PATTERN_FIELDS.snooze]: {
      operators: [
        PATTERN_OPERATORS.snoozed,
        PATTERN_OPERATORS.notSnoozed,
      ],
    },
  }));

  return {
    attributesMap,
  };
};

/**
 * Hook to manage alarm advanced search attributes for entities.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {Array} options.infosItems - An array of information items used in the search.
 * @param {string} [options.prefix = ''] - The prefix to use for the attributes.
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for entities.
 * @property {Function} attributesMap[].fetchValues - Function to fetch values for the attribute.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 * @property {string} attributesMap[].itemValue - The key used to identify the value in fetched data.
 * @property {string} attributesMap[].itemText - The key used to identify the display text in fetched data.
 * @property {Array} [attributesMap[].values] - Predefined values for the attribute, if applicable.
 */
export const useAdvancedSearchEntityAttributes = ({ infosItems, prefix = '' } = {}) => {
  const { t } = useI18n();
  const { fetchCategoriesListWithoutStore } = useEntityCategory();
  const { getEntityOptions } = useGetEntityOptions();

  const attributesMap = computed(() => ({
    /**
     * Basic
     */
    [ENTITY_PATTERN_FIELDS.id]: getEntityOptions(),
    [ENTITY_PATTERN_FIELDS.customId]: getEntityOptions(),
    [ENTITY_PATTERN_FIELDS.name]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.category]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
      ],
      fetchValues: fetchCategoriesListWithoutStore,
      itemValue: '_id',
      itemText: 'name',
    },
    [ENTITY_PATTERN_FIELDS.type]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
      ],
      values: Object.values(ENTITY_TYPES).map(type => ({ value: type, text: t(`entity.types.${type}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ENTITY_PATTERN_FIELDS.component]: getEntityOptions([ENTITY_TYPES.component]),
    [ENTITY_PATTERN_FIELDS.connector]: getEntityOptions([ENTITY_TYPES.connector]),
    [ENTITY_PATTERN_FIELDS.resource]: getEntityOptions([ENTITY_TYPES.resource]),
    [ENTITY_PATTERN_FIELDS.impactLevel]: getNumberMinValueAttributes(1),
    [ENTITY_PATTERN_FIELDS.impactState]: getNumberMinValueAttributes(0),
    [ENTITY_PATTERN_FIELDS.importSource]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.state]: {
      operators: PATTERN_NUMBER_OPERATORS,
      values: Object.values(ALARM_STATES).map(value => ({ value, text: t(`common.stateTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ENTITY_PATTERN_FIELDS.status]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
      ],
      values: Object.values(ALARM_STATUSES).map(value => ({ value, text: t(`common.statusTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ENTITY_PATTERN_FIELDS.infos]: {
      items: unref(infosItems) ?? [],
    },
    [ENTITY_PATTERN_FIELDS.componentInfos]: {
      items: unref(infosItems) ?? [],
    },
    [ENTITY_PATTERN_FIELDS.enabled]: {
      operators: [
        PATTERN_OPERATORS.enabled,
        PATTERN_OPERATORS.disabled,
      ],
    },

    /**
     * Events
     */
    [ENTITY_PATTERN_FIELDS.koEvents]: getNumberMinValueAttributes(0),
    [ENTITY_PATTERN_FIELDS.okEvents]: getNumberMinValueAttributes(0),

    /**
     * Dates
     */
    [ENTITY_PATTERN_FIELDS.idleSince]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.imported]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.lastUpdateDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.lastPbehaviorDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.lastEventDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [ENTITY_PATTERN_FIELDS.lastAlarmUpdateDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
  }));

  const attributesMapWithPrefix = computed(() => {
    const unwrappedPrefix = unref(prefix);

    return unwrappedPrefix ? addPrefixToAttributesMap(attributesMap.value, unwrappedPrefix) : attributesMap.value;
  });

  return {
    attributesMap: attributesMapWithPrefix,
  };
};

/**
 * Hook to manage advanced search attributes for pbehaviors.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {string} [options.prefix = ''] - The prefix to use for the attributes.
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for pbehaviors.
 * @property {Function} attributesMap[].fetchValues - Function to fetch values for the attribute.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 * @property {string} attributesMap[].text - Display text for the attribute.
 * @property {Array} [attributesMap[].values] - Predefined values for the attribute, if applicable.
 */
export const useAdvancedSearchPbehaviorAttributes = ({ prefix = '' } = {}) => {
  const { t } = useI18n();
  const { fetchPbehaviorsListWithoutStore } = usePbehavior();
  const { fetchPbehaviorReasonsListWithoutStore } = usePbehaviorReason();
  const { fetchPbehaviorTypesListWithoutStore } = usePbehaviorType();
  const { fetchUsersListWithoutStore } = useUser();

  const BASE_OPERATORS = [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual];
  const BASE_OPTIONS = {
    operators: BASE_OPERATORS,
    itemValue: '_id',
    itemText: 'name',
  };

  /**
   * We are using PBEHAVIOR_FIELDS here instead of PBEHAVIOR_PATTERN_FIELDS because
   * we are not using the pattern prefix in all places
   */
  const attributesMap = computed(() => ({
    [PBEHAVIOR_FIELDS.name]: {
      ...BASE_OPTIONS,
      itemValue: 'name',
      text: t('pbehavior.pbehaviorName'),
      fetchValues: fetchPbehaviorsListWithoutStore,
    },
    [PBEHAVIOR_FIELDS.author]: {
      operators: ADVANCED_SEARCH_USER_OPERATORS,
      text: t('common.author'),
      fetchValues: fetchUsersListWithoutStore,
      itemValue: '_id',
      itemText: 'display_name',
    },
    [PBEHAVIOR_FIELDS.enabled]: {
      operators: [
        PATTERN_OPERATORS.enabled,
        PATTERN_OPERATORS.disabled,
      ],
    },
    [PBEHAVIOR_FIELDS.rrule]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.reason]: {
      ...BASE_OPTIONS,
      text: t('pbehavior.pbehaviorReason'),
      fetchValues: fetchPbehaviorReasonsListWithoutStore,
    },
    [PBEHAVIOR_FIELDS.type]: {
      ...BASE_OPTIONS,
      text: t('pbehavior.pbehaviorType'),
      fetchValues: fetchPbehaviorTypesListWithoutStore,
    },
    [PBEHAVIOR_FIELDS.canonicalType]: {
      operators: BASE_OPERATORS,
      text: t('pbehavior.pbehaviorCanonicalType'),
      values: Object.values(PBEHAVIOR_TYPE_TYPES).map(type => ({
        value: type,
        text: t(`pbehavior.types.types.${type}`),
      })),
    },
    [PBEHAVIOR_FIELDS.tstart]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.tstop]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.rruleEnd]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.created]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.updated]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.lastAlarmDate]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [PBEHAVIOR_FIELDS.alarmCount]: getNumberMinValueAttributes(0),
  }));

  const attributesMapWithPrefix = computed(() => {
    const unwrappedPrefix = unref(prefix);

    return unwrappedPrefix ? addPrefixToAttributesMap(attributesMap.value, unwrappedPrefix) : attributesMap.value;
  });

  return {
    attributesMap: attributesMapWithPrefix,
  };
};

/**
 * Hook to manage advanced search attributes for dynamic infos.
 *
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for dynamic infos.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 */
export const useAdvancedSearchDynamicInfoAttributes = () => {
  const { fetchUsersListWithoutStore } = useUser();
  const { fetchDynamicInfosListWithoutStore } = useDynamicInfo();

  const attributesMap = computed(() => ({
    [DYNAMIC_INFO_FIELDS.id]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
      fetchValues: fetchDynamicInfosListWithoutStore,
      itemValue: '_id',
      itemText: 'name',
    },
    [DYNAMIC_INFO_FIELDS.name]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [DYNAMIC_INFO_FIELDS.author]: {
      operators: ADVANCED_SEARCH_USER_OPERATORS,
      fetchValues: fetchUsersListWithoutStore,
      itemValue: '_id',
      itemText: 'display_name',
    },
    [DYNAMIC_INFO_FIELDS.enabled]: {
      operators: [
        PATTERN_OPERATORS.enabled,
        PATTERN_OPERATORS.disabled,
      ],
    },
    [DYNAMIC_INFO_FIELDS.created]: {
      operators: PATTERN_DATE_OPERATORS,
    },
    [DYNAMIC_INFO_FIELDS.updated]: {
      operators: PATTERN_DATE_OPERATORS,
    },
  }));

  return {
    attributesMap,
  };
};
