import { computed, unref } from 'vue';

import {
  ADVANCED_SEARCH_USER_OPERATORS,
  ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
  ALARM_EVENT_INITIATORS,
  ALARM_FIELDS,
  ALARM_STATES,
  ALARM_STATUSES,
  BASIC_ENTITY_TYPES,
  ENTITY_TYPES,
  ENTITY_FIELDS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PBEHAVIOR_TYPE_TYPES,
  PBEHAVIOR_FIELDS,
  PATTERN_ALARM_TAG_LABEL_OPERATORS,
  DYNAMIC_INFO_FIELDS,
} from '@/constants';

import { addPrefixToAttributesMap } from '@/helpers/search/advanced-search';

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
    [ALARM_FIELDS.displayName]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.connector]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.connector]),

      itemText: 'connector_type',
      itemValue: 'connector_type',
    },
    [ALARM_FIELDS.connectorName]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.connector]),

      itemValue: 'name',
    },
    [ALARM_FIELDS.component]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.component]),

      itemText: 'component',
    },
    [ALARM_FIELDS.resource]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.resource]),

      itemValue: 'name',
    },
    [ALARM_FIELDS.state]: {
      operators: PATTERN_NUMBER_OPERATORS,
      values: Object.values(ALARM_STATES).map(value => ({ value, text: t(`common.stateTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ALARM_FIELDS.status]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
      ],
      values: Object.values(ALARM_STATUSES).map(value => ({ value, text: t(`common.statusTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ALARM_FIELDS.tags]: {
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
    [ALARM_FIELDS.infos]: {
      items: unref(infosItems),
    },
    [ALARM_FIELDS.meta]: {
      operators: [
        PATTERN_OPERATORS.isMetaAlarm,
        PATTERN_OPERATORS.isNotMetaAlarm,
        PATTERN_OPERATORS.ruleIs,
      ],
      fetchValues: fetchMetaAlarmRulesListWithoutStore,
    },
    [ALARM_FIELDS.changeState]: {
      operators: PATTERN_EXISTS_OPERATORS,
    },
    [ALARM_FIELDS.totalStateChanges]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },

    /**
     * Messages
     */
    [ALARM_FIELDS.output]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.longOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialLongOutput]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.lastComment]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.lastCommentInitiator]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },

    /**
     * Ticket
     */
    [ALARM_FIELDS.ticketMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ticketInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.ticketValue]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ticket]: {
      operators: [
        PATTERN_OPERATORS.ticketAssociated,
        PATTERN_OPERATORS.ticketNotAssociated,
      ],
    },

    /**
     * Date
     */
    [ALARM_FIELDS.creationDate]: {},
    [ALARM_FIELDS.lastUpdateDate]: {},
    [ALARM_FIELDS.lastEventDate]: {},
    [ALARM_FIELDS.ackAt]: {},
    [ALARM_FIELDS.resolved]: {},
    [ALARM_FIELDS.activationDate]: {},
    [ALARM_FIELDS.duration]: DURATION_OPTIONS,

    /**
     * Actions
     */
    [ALARM_FIELDS.ack]: {
      operators: [
        PATTERN_OPERATORS.acked,
        PATTERN_OPERATORS.notAcked,
      ],
    },
    [ALARM_FIELDS.ackBy]: {
      operators: ADVANCED_SEARCH_USER_OPERATORS,
      fetchValues: fetchUsersListWithoutStore,
      itemValue: 'display_name',
      itemText: 'display_name',
    },
    [ALARM_FIELDS.ackMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ackInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.canceled]: {
      operators: [
        PATTERN_OPERATORS.canceled,
        PATTERN_OPERATORS.notCanceled,
      ],
    },
    [ALARM_FIELDS.canceledInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.activated]: {
      operators: [
        PATTERN_OPERATORS.activated,
        PATTERN_OPERATORS.inactive,
      ],
    },
    [ALARM_FIELDS.snooze]: {
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
    [ENTITY_FIELDS.id]: getEntityOptions(),
    [ENTITY_FIELDS.name]: {
      operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    },
    [ENTITY_FIELDS.categoryName]: {
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
    [ENTITY_FIELDS.type]: {
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
    [ENTITY_FIELDS.categoryName]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
      ],
      fetchValues: fetchCategoriesListWithoutStore,
    },
    [ENTITY_FIELDS.component]: getEntityOptions([ENTITY_TYPES.component]),
    [ENTITY_FIELDS.connector]: getEntityOptions([ENTITY_TYPES.connector]),
    [ENTITY_FIELDS.resource]: getEntityOptions([ENTITY_TYPES.resource]),
    [ENTITY_FIELDS.impactLevel]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },
    [ENTITY_FIELDS.impactState]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },
    [ENTITY_FIELDS.state]: {
      operators: PATTERN_NUMBER_OPERATORS,
      values: Object.values(ALARM_STATES).map(value => ({ value, text: t(`common.stateTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ENTITY_FIELDS.status]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
      ],
      values: Object.values(ALARM_STATUSES).map(value => ({ value, text: t(`common.statusTypes.${value}`) })),
      itemText: 'text',
      itemValue: 'value',
    },
    [ENTITY_FIELDS.infos]: {
      items: unref(infosItems) ?? [],
    },
    [ENTITY_FIELDS.componentInfos]: {
      items: unref(infosItems) ?? [],
    },
    [ENTITY_FIELDS.enabled]: {
      operators: [
        PATTERN_OPERATORS.enabled,
        PATTERN_OPERATORS.disabled,
      ],
    },

    /**
     * Events
     */
    [ENTITY_FIELDS.koEvents]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },
    [ENTITY_FIELDS.okEvents]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },

    /**
     * Dates
     */
    [ENTITY_FIELDS.idleSince]: {},
    [ENTITY_FIELDS.imported]: {},
    [ENTITY_FIELDS.lastUpdateDate]: {},
    [ENTITY_FIELDS.lastPbehaviorDate]: {},
    [ENTITY_FIELDS.lastEventDate]: {},
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
      text: t('pbehavior.pbehaviorName'),
      fetchValues: fetchPbehaviorsListWithoutStore,
    },
    [PBEHAVIOR_FIELDS.author]: {
      operators: ADVANCED_SEARCH_USER_OPERATORS,
      text: t('common.author'),
      fetchValues: fetchUsersListWithoutStore,
      itemValue: 'display_name',
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
    [PBEHAVIOR_FIELDS.tstart]: {},
    [PBEHAVIOR_FIELDS.tstop]: {},
    [PBEHAVIOR_FIELDS.rruleEnd]: {},
    [PBEHAVIOR_FIELDS.created]: {},
    [PBEHAVIOR_FIELDS.updated]: {},
    [PBEHAVIOR_FIELDS.lastAlarmDate]: {},
    [PBEHAVIOR_FIELDS.alarmCount]: {
      operators: PATTERN_NUMBER_OPERATORS,
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
      itemValue: 'display_name',
      itemText: 'display_name',
    },
    [DYNAMIC_INFO_FIELDS.enabled]: {
      operators: [
        PATTERN_OPERATORS.enabled,
        PATTERN_OPERATORS.disabled,
      ],
    },
    [DYNAMIC_INFO_FIELDS.created]: {},
    [DYNAMIC_INFO_FIELDS.updated]: {},
  }));

  return {
    attributesMap,
  };
};
