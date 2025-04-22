import {
  ref,
  computed,
  inject,
  unref,
  watch,
  onBeforeMount,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import { Validator } from 'vee-validate';

import {
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_ADVANCED_SEARCH_CHIP_TYPES,
  ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME,
  ALARM_ADVANCED_SEARCH_GROUPS,
  ALARM_ADVANCED_SEARCH_ENTITY_OPERATORS,
  ALARM_ADVANCED_SEARCH_GROUPS_GROUPED,
  ALARM_ADVANCED_SEARCH_ALARM_ENTITY_FIELDS,
  ALARM_ADVANCED_SEARCH_ALARM_PBEHAVIOR_INFO_FIELDS,
  ALARM_EVENT_INITIATORS,
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_STATES,
  ALARM_STATUSES,
  BASIC_ENTITY_TYPES,
  ENTITY_TYPES,
  MAX_LIMIT,
  PATTERN_DURATION_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_FIELD_TYPES,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_STRING_OPERATORS,
  PBEHAVIOR_TYPE_TYPES,
  PATTERN_ALARM_TAG_LABEL_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
} from '@/constants';

import { deepKeyBy } from '@/helpers/array';
import { getOperatorsByFieldType } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useEntity } from '@/hooks/store/modules/entity';
import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';
import { useService } from '@/hooks/store/modules/service';
import { usePendingHandler } from '@/hooks/query/pending';
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
import { useUser } from '@/hooks/store/modules/user';
import { useEntityCategory } from '@/hooks/store/modules/entity-category';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { useComponentInstance } from '@/hooks/vue';
import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label';

/**
 * Hook to manage fetching and processing entity information keys for advanced search.
 *
 * @returns {Object} An object containing:
 * @property {boolean} pending - A reactive boolean indicating the fetch operation's pending state.
 * @property {Array} items - A computed array of entity information keys formatted for search configurations.
 */
export const useEntityInfosKeys = () => {
  const { t } = useI18n();
  const { fetchEntityInfosKeysWithoutStore } = useService();

  const entityInfosKeys = ref([]);
  const {
    pending,
    handler,
  } = usePendingHandler(async () => {
    const { data: infos } = await fetchEntityInfosKeysWithoutStore({ params: { limit: MAX_LIMIT } });

    entityInfosKeys.value = infos;
  });

  /**
   * Generates default child items for a given chip text prefix.
   *
   * @param {string} chipTextPrefix - The prefix to be used for chip text.
   * @returns {Array<Object>} An array of objects, each containing:
   *   - {string} value: The value of the item.
   *   - {string} text: The translated text for the item.
   *   - {string} chipText: The complete chip text including the prefix.
   */
  const getDefaultItemChildren = (chipTextPrefix = '') => [
    PATTERN_RULE_INFOS_FIELDS.name,
    PATTERN_RULE_INFOS_FIELDS.value,
  ].map((value) => {
    const text = t(`common.${value}`);
    const result = { value, text };

    if (chipTextPrefix) {
      result.chipText = `${chipTextPrefix}.${text}`;
    }

    if (value === PATTERN_RULE_INFOS_FIELDS.name) {
      result.operators = PATTERN_EXISTS_OPERATORS;
    }

    if (value === PATTERN_RULE_INFOS_FIELDS.name) {
      result.operators = PATTERN_EXISTS_OPERATORS;
    }

    return result;
  });

  /**
   * Retrieves items formatted for search configurations, with optional chip text prefix.
   *
   * @param {string} chipTextPrefix - The prefix to be used for chip text.
   * @returns {Array<Object>} An array of objects, each containing:
   *   - {string} value: The value of the item.
   *   - {string} chipText: The complete chip text including the prefix.
   *   - {string} text: The text for the item.
   *   - {Array<Object>} items: The default child items for the chip text.
   */
  const getItems = chipTextPrefix => entityInfosKeys.value.map(({ value }) => {
    const chipText = chipTextPrefix ? `${chipTextPrefix}.${value}` : undefined;

    return {
      value,
      chipText,
      text: value,
      items: getDefaultItemChildren(chipText),
    };
  });

  const alarmItems = computed(() => getItems(t(ALARM_FIELDS_TO_LABELS_KEYS[ALARM_FIELDS.infos])));
  const entityItems = computed(() => getItems(t(ALARM_FIELDS_TO_LABELS_KEYS[ALARM_FIELDS.entityInfos])));

  onMounted(handler);

  return {
    pending,
    alarmItems,
    entityItems,
    getItems,
  };
};

/**
 * Hook to retrieve entity options for advanced search configurations.
 *
 * @returns {Object} An object containing the `getEntityOptions` function.
 * @property {Function} getEntityOptions - Function to get entity options.
 * @property {Array} getEntityOptions.operators - Array of operators applicable to the entity.
 * @property {Function} getEntityOptions.fetchValues - Function to fetch entity values with parameters.
 * @property {string} getEntityOptions.itemText - The key used to identify the display text in fetched data.
 * @property {string} getEntityOptions.itemValue - The key used to identify the value in fetched data.
 */
export const useGetEntityOptions = () => {
  const { fetchContextEntitiesListWithoutStore } = useEntity();

  const getEntityOptions = (type = Object.values(BASIC_ENTITY_TYPES)) => ({
    operators: ALARM_ADVANCED_SEARCH_ENTITY_OPERATORS,
    fetchValues: ({ params }) => fetchContextEntitiesListWithoutStore({
      params: { ...params, type },
    }),
    itemText: 'name',
    itemValue: '_id',
  });

  return {
    getEntityOptions,
  };
};

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

  const STRING_WITH_ONE_OF_OPERATORS = [
    ...PATTERN_STRING_OPERATORS,

    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ];
  const STRING_WITH_EXIST_AND_ONE_OF_OPERATORS = [
    ...STRING_WITH_ONE_OF_OPERATORS,

    PATTERN_OPERATORS.exist,
  ];

  const USER_OPERATORS = [
    PATTERN_OPERATORS.equal,
    PATTERN_OPERATORS.notEqual,
    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ];

  const DURATION_OPTIONS = { operators: PATTERN_DURATION_OPERATORS };

  const INITIATOR_OPTIONS = {
    operators: USER_OPERATORS,
    values: Object.values(ALARM_EVENT_INITIATORS).map(initiator => ({ value: initiator, text: initiator })),
  };

  const attributesMap = computed(() => ({
    /**
     * Basic
     */
    [ALARM_FIELDS.displayName]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
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
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialOutput]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialLongOutput]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
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
      operators: USER_OPERATORS,
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
 * Hook to manage advanced search attributes for entities.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {Array} options.infosItems - An array of information items used in the search.
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for entities.
 * @property {Function} attributesMap[].fetchValues - Function to fetch values for the attribute.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 * @property {string} attributesMap[].itemValue - The key used to identify the value in fetched data.
 * @property {string} attributesMap[].itemText - The key used to identify the display text in fetched data.
 * @property {Array} [attributesMap[].values] - Predefined values for the attribute, if applicable.
 */
export const useAdvancedSearchEntityAttributes = ({ infosItems }) => {
  const { t } = useI18n();
  const { fetchCategoriesListWithoutStore } = useEntityCategory();
  const { getEntityOptions } = useGetEntityOptions();

  const attributesMap = computed(() => ({
    [ALARM_FIELDS.entityId]: getEntityOptions(),
    [ALARM_FIELDS.entityName]: {
      operators: ALARM_ADVANCED_SEARCH_ENTITY_OPERATORS,
    },
    [ALARM_FIELDS.entityCategoryName]: {
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
    [ALARM_FIELDS.entityType]: {
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
    [ALARM_FIELDS.entityComponent]: getEntityOptions([ENTITY_TYPES.component]),
    [ALARM_FIELDS.entityConnector]: getEntityOptions([ENTITY_TYPES.connector]),
    [ALARM_FIELDS.entityImpactLevel]: {
      operators: PATTERN_NUMBER_OPERATORS,
    },
    [ALARM_FIELDS.entityInfos]: {
      items: unref(infosItems),
    },
    [ALARM_FIELDS.entityComponentInfos]: {
      items: unref(infosItems),
    },
  }));

  return {
    attributesMap,
  };
};

/**
 * Hook to manage advanced search attributes for pbehaviors.
 *
 * @returns {Object} An object containing the computed `attributesMap`.
 * @property {Object} attributesMap - A map of attribute configurations for pbehaviors.
 * @property {Function} attributesMap[].fetchValues - Function to fetch values for the attribute.
 * @property {Array} attributesMap[].operators - Array of operators applicable to the attribute.
 * @property {string} attributesMap[].text - Display text for the attribute.
 * @property {Array} [attributesMap[].values] - Predefined values for the attribute, if applicable.
 */
export const useAdvancedSearchPbehaviorAttributes = () => {
  const { t } = useI18n();
  const { fetchPbehaviorsListWithoutStore } = usePbehavior();
  const { fetchPbehaviorReasonsListWithoutStore } = usePbehaviorReason();
  const { fetchPbehaviorTypesListWithoutStore } = usePbehaviorType();

  const BASE_OPERATORS = [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual];
  const BASE_OPTIONS = {
    operators: BASE_OPERATORS,
    itemValue: '_id',
    itemText: 'name',
  };

  const attributesMap = computed(() => ({
    [ALARM_FIELDS.pbehaviorInfoId]: {
      ...BASE_OPTIONS,
      text: t('pbehavior.pbehaviorName'),
      fetchValues: fetchPbehaviorsListWithoutStore,
    },
    [ALARM_FIELDS.pbehaviorInfoReason]: {
      ...BASE_OPTIONS,
      text: t('pbehavior.pbehaviorReason'),
      fetchValues: fetchPbehaviorReasonsListWithoutStore,
    },
    [ALARM_FIELDS.pbehaviorInfoType]: {
      ...BASE_OPTIONS,
      text: t('pbehavior.pbehaviorType'),
      fetchValues: fetchPbehaviorTypesListWithoutStore,
    },
    [ALARM_FIELDS.pbehaviorInfoCanonicalType]: {
      operators: BASE_OPERATORS,
      text: t('pbehavior.pbehaviorCanonicalType'),
      values: Object.values(PBEHAVIOR_TYPE_TYPES).map(type => ({
        value: type,
        text: t(`pbehavior.types.types.${type}`),
      })),
    },
  }));

  return {
    attributesMap,
  };
};

/**
 * Hook to manage advanced search attributes for alarms, entities, and pbehaviors.
 *
 * @param {Object} options - Options for configuring the advanced search attributes.
 * @param {boolean} options.allowAlarmFields - Flag to allow alarm fields in the search.
 * @param {boolean} options.allowEntityFields - Flag to allow entity fields in the search.
 * @param {boolean} options.allowPbehaviorFields - Flag to allow pbehavior fields in the search.
 * @returns {Object} An object containing the pending state and the computed attributes.
 */
export const useAdvancedSearchAttributes = ({
  allowAlarmFields,
  allowEntityFields,
  allowPbehaviorFields,
}) => {
  const { t, tc } = useI18n();

  const {
    pending: infosPending,
    alarmItems: alarmInfosItems,
    entityItems: entityInfosItems,
  } = useEntityInfosKeys();

  const { attributesMap: alarmAttributesMap } = useAdvancedSearchAlarmAttributes({ infosItems: alarmInfosItems });
  const { attributesMap: entityAttributesMap } = useAdvancedSearchEntityAttributes({ infosItems: entityInfosItems });
  const { attributesMap: pbehaviorAttributesMap } = useAdvancedSearchPbehaviorAttributes();

  const attributesMap = computed(() => ({
    ...alarmAttributesMap.value,
    ...entityAttributesMap.value,
    ...pbehaviorAttributesMap.value,
  }));

  const prepareItem = (field, allow) => {
    const attributes = attributesMap.value[field];

    return {
      ...attributesMap.value[field],
      value: field,
      text: attributes.text ?? tc(ALARM_FIELDS_TO_LABELS_KEYS[field], 2),
      disabled: !allow,
    };
  };

  const attributes = computed(() => {
    const unwrappedAllowAlarmFields = unref(allowAlarmFields);
    const unwrappedAllowEntityFields = unref(allowEntityFields);
    const unwrappedAllowPbehaviorFields = unref(allowPbehaviorFields);

    const result = Object.entries(ALARM_ADVANCED_SEARCH_GROUPS_GROUPED).reduce((acc, [group, items]) => {
      const header = t(`advancedSearch.groups.${group}`);

      acc.push(
        { header, value: header },
        ...items.map(field => prepareItem(field, unwrappedAllowAlarmFields)),
      );

      return acc;
    }, []);

    const entityHeader = t(`advancedSearch.groups.${ALARM_ADVANCED_SEARCH_GROUPS.entity}`);
    const pbehaviorHeader = t(`advancedSearch.groups.${ALARM_ADVANCED_SEARCH_GROUPS.pbehavior}`);

    result.push(
      { header: entityHeader, value: entityHeader },
      ...ALARM_ADVANCED_SEARCH_ALARM_ENTITY_FIELDS.map(field => (
        prepareItem(field, unwrappedAllowEntityFields)
      )),

      { header: pbehaviorHeader, value: pbehaviorHeader },
      ...ALARM_ADVANCED_SEARCH_ALARM_PBEHAVIOR_INFO_FIELDS.map(field => (
        prepareItem(field, unwrappedAllowPbehaviorFields)
      )),
    );

    return result;
  });

  return {
    pending: infosPending,
    attributes,
  };
};

/**
 * Hook to extend the validator with a custom rule for advanced search validation.
 *
 * @param {Object} options - Options for the validator.
 * @param {Array} options.rules - The array of rules to be validated.
 * @returns {Object} An object containing the `extendValidatorRule` function.
 */
export const useAdvancedSearchValidator = ({ rules }) => {
  const instance = useComponentInstance();

  /**
   * Extends the validator with a custom rule named ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME.
   * This rule is used to validate advanced search criteria based on specific conditions.
   *
   * @description
   * The ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME checks the following conditions:
   * - If the rule has an attribute and is not finished, it returns false.
   * - If the rule does not have an attribute but has a union, it checks the last two rules:
   * - If the last rule does not have an attribute and the second-to-last rule has the same key as the current rule,
   *   it returns false.
   * - Otherwise, it returns true.
   */
  const extendValidatorRule = () => (
    instance.$validator.extend(ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME, ({ rule, finished }) => {
      const unwrappedRules = unref(rules);

      if (rule.attribute && !finished) {
        return false;
      }

      if (!rule.attribute && rule.union) {
        const lastRule = unwrappedRules.at(-1);
        const preLastRule = unwrappedRules.at(-2);

        return !(!lastRule?.attribute && preLastRule?.key === rule?.key);
      }

      return true;
    })
  );

  onBeforeMount(extendValidatorRule);

  return {
    extendValidatorRule,
  };
};

/**
 * A Vue composition function that manages active items for an advanced search rule.
 * It computes various properties based on the rule's field type, operators, interval ranges,
 * input types, and union conditions.
 *
 * @param {Object} options - Options for configuring the active items.
 * @param {Object} [options.rule = {}] - The rule object to evaluate.
 * @param {Array} [options.attributes = []] - List of available attributes for the rule.
 * @param {Array} [options.intervalRanges = []] - List of interval ranges for range-based rules.
 * @param {Array} [options.inputTypes = []] - List of input types for the rule.
 * @param {boolean} [options.allowOr = false] - Determines if the 'OR' union condition is allowed.
 * @returns {Object} - An object containing the current attribute and items by type.
 */
export const useAdvancedSearchRuleActiveItems = ({
  rule = {},
  attributes = [],
  intervalRanges = [],
  inputTypes = [],
  allowOr = false,
} = {}) => {
  const { t } = useI18n();

  const attributesMap = computed(() => deepKeyBy(unref(attributes), 'value'));
  const currentAttribute = computed(() => attributesMap.value[unref(rule).attribute]);

  /**
   * FIELD TYPE
   */
  const isBooleanFieldType = computed(() => unref(rule).fieldType === PATTERN_FIELD_TYPES.boolean);

  /**
   * OPERATORS ITEMS
   */
  const operators = computed(() => {
    const unwrappedRule = unref(rule);

    let result = getOperatorsByFieldType(unwrappedRule.fieldType);

    if (unwrappedRule.fieldType === PATTERN_FIELD_TYPES.string) {
      result = [
        ...result,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
      ];
    }

    return result;
  });

  const preparedOperators = computed(() => (
    (currentAttribute.value?.operators ?? operators.value ?? []).map(operator => ({
      text: t(`common.operators.${operator}`),
      value: operator,
    }))
  ));

  /**
   * INTERVAL ITEMS
   */
  const preparedIntervalRanges = computed(() => (
    unref(intervalRanges).map(range => ({
      ...range,
      text: t(`quickRanges.types.${range.value}`),
    }))
  ));

  /**
   * INPUT TYPES ITEMS
   */
  const preparedInputTypes = computed(() => (
    unref(inputTypes).map(type => ({
      ...type,
      text: t(`common.mixedField.types.${type.value}`),
    }))
  ));

  /**
   * UNION ITEMS
   */
  const preparedUnionItems = computed(() => (
    Object.values(ADVANCED_SEARCH_UNION_CONDITIONS).map(value => ({
      value,
      text: value,
      disabled: value === ADVANCED_SEARCH_UNION_CONDITIONS.or && !unref(allowOr),
    }))
  ));

  /**
   * BOOLEAN ITEMS
   */
  const preparedBooleanItems = computed(() => [
    { text: t('common.true'), value: true }, { text: t('common.false'), value: false },
  ]);

  /**
   * VALUE ITEMS
   */
  const preparedValueItems = computed(() => (
    currentAttribute.value?.values
      ?? { [isBooleanFieldType.value]: preparedBooleanItems.value }.true
      ?? []
  ));

  /**
   * ITEMS BY TYPE
   */
  const itemsByType = computed(() => ({
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute]: unref(attributes),
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator]: preparedOperators.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.range]: preparedIntervalRanges.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.fieldType]: preparedInputTypes.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.value]: preparedValueItems.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.union]: preparedUnionItems.value,
  }));

  return {
    attributesMap,
    currentAttribute,
    itemsByType,
  };
};

/**
 * A Vue composition function that manages the attachment and detachment of validation rules
 * for an advanced search rule. It uses the injected `$validator` from `vee-validate` to
 * dynamically attach or detach validation rules based on the `disabled` state.
 *
 * @param {Object} options - Options for configuring the validator attachment.
 * @param {Object} [options.rule = {}] - The rule object to be validated.
 * @param {boolean} [options.disabled = false] - Determines if the validation should be disabled.
 * @param {boolean} [options.isFinishedRule = false] - Indicates if the rule is finished, affecting validation logic.
 * @returns {Object} - An object containing the `validator` instance.
 */
export const useAttachAdvancedSearchRuleValidator = ({
  rule = {},
  disabled = false,
  isFinishedRule = false,
} = {}) => {
  const validator = inject('$validator', new Validator());

  /**
   * Attaches a validation rule to the validator instance.
   * The rule is identified by its key and uses a predefined validation rule name.
   * The getter function provides the rule and its finished state for validation.
   */
  const attachValidationRule = () => validator.attach({
    name: unref(rule).key,
    rules: ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME,
    getter: () => ({ rule: unref(rule), finished: unref(isFinishedRule) }),
  });

  /**
   * Detaches a validation rule from the validator instance.
   * The rule is identified by its key, ensuring it is no longer validated.
   */
  const detachValidationRule = () => validator.detach(unref(rule).key);

  watch(() => unref(disabled), (newDisabled) => {
    if (newDisabled) {
      detachValidationRule();

      return;
    }

    attachValidationRule();
  }, { immediate: true });

  onBeforeUnmount(detachValidationRule);

  return {
    validator,
  };
};
