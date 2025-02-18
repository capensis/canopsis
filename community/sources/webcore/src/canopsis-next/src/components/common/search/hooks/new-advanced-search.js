import {
  ref,
  computed,
  unref,
  onBeforeMount,
  onMounted,
} from 'vue';

import {
  ADVANCED_SEARCH_VALIDATION_RULE_NAME,
  ALARM_ADVANCED_SEARCH_GROUPS,
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
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_STRING_OPERATORS,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

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

export const ENTITY_OPERATORS = [
  ...PATTERN_STRING_OPERATORS,

  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.isNotOneOf,
];

export const ADVANCED_SEARCH_GROUPED_ALARM_FIELDS = {
  [ALARM_ADVANCED_SEARCH_GROUPS.basic]: [ // TODO: rename to ADVANCED_SEARCH_ALARM_GROUPS
    ALARM_FIELDS.displayName,
    ALARM_FIELDS.connector,
    ALARM_FIELDS.connectorName,
    ALARM_FIELDS.component,
    ALARM_FIELDS.resource,
    ALARM_FIELDS.state,
    ALARM_FIELDS.status,
    ALARM_FIELDS.tags,
    ALARM_FIELDS.infos,
    ALARM_FIELDS.meta,
    ALARM_FIELDS.changeState,
    ALARM_FIELDS.totalStateChanges,
  ],
  [ALARM_ADVANCED_SEARCH_GROUPS.messages]: [
    ALARM_FIELDS.output,
    ALARM_FIELDS.longOutput,
    ALARM_FIELDS.initialOutput,
    ALARM_FIELDS.initialLongOutput,
    ALARM_FIELDS.lastComment,
    ALARM_FIELDS.lastCommentInitiator,
  ],
  [ALARM_ADVANCED_SEARCH_GROUPS.ticket]: [
    ALARM_FIELDS.ticketMessage,
    ALARM_FIELDS.ticketInitiator,
    ALARM_FIELDS.ticketValue,
    ALARM_FIELDS.ticket,
  ],
  [ALARM_ADVANCED_SEARCH_GROUPS.dates]: [
    ALARM_FIELDS.creationDate,
    ALARM_FIELDS.lastUpdateDate,
    ALARM_FIELDS.lastEventDate,
    ALARM_FIELDS.ackAt,
    ALARM_FIELDS.resolved,
    ALARM_FIELDS.activationDate,
    ALARM_FIELDS.duration,
  ],
  [ALARM_ADVANCED_SEARCH_GROUPS.actions]: [
    ALARM_FIELDS.ack,
    ALARM_FIELDS.ackBy,
    ALARM_FIELDS.ackMessage,
    ALARM_FIELDS.ackInitiator,
    ALARM_FIELDS.canceled,
    ALARM_FIELDS.canceledInitiator,
    ALARM_FIELDS.activated,
    ALARM_FIELDS.snooze,
  ],
};

export const ADVANCED_SEARCH_ALARM_ENTITY_FIELDS = [
  ALARM_FIELDS.entityId,
  ALARM_FIELDS.entityName,
  ALARM_FIELDS.entityCategoryName,
  ALARM_FIELDS.entityType,
  ALARM_FIELDS.entityComponent,
  ALARM_FIELDS.entityConnector,
  ALARM_FIELDS.entityImpactLevel,
  ALARM_FIELDS.entityInfos,
  ALARM_FIELDS.entityComponentInfos,
];

export const ADVANCED_SEARCH_ALARM_PBEHAVIOR_INFO_FIELDS = [
  ALARM_FIELDS.pbehaviorInfoName,
  ALARM_FIELDS.pbehaviorInfoReason,
  ALARM_FIELDS.pbehaviorInfoType,
  ALARM_FIELDS.pbehaviorInfoCanonicalType,
];

export const useEntityInfosKeys = () => {
  const { t } = useI18n();
  const { fetchEntityInfosKeysWithoutStore } = useService();

  const entityInfosKeys = ref([]);
  const {
    pending,
    handler,
  } = usePendingHandler(async () => {
    const { data: infos } = await fetchEntityInfosKeysWithoutStore({ params: MAX_LIMIT });

    entityInfosKeys.value = infos;
  });

  const getDefaultItemChildren = () => [
    {
      value: 'name',
      text: t('common.name'),
    }, {
      value: 'value',
      text: t('common.value'),
    },
  ];

  const items = computed(() => entityInfosKeys.value.map(({ value }) => ({
    value,
    text: value,
    items: getDefaultItemChildren(),
  })));

  onMounted(handler);

  return {
    pending,
    items,
  };
};

export const useGetEntityOptions = () => {
  const { fetchContextEntitiesListWithoutStore } = useEntity();

  const getEntityOptions = (type = Object.values(BASIC_ENTITY_TYPES)) => ({
    operators: ENTITY_OPERATORS,
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

export const useAdvancedSearchAlarmAttributes = ({ infosItems }) => {
  const { t } = useI18n();
  const { fetchAlarmTagsListWithoutStore } = useAlarmTag();
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
    [ALARM_FIELDS.connector]: getEntityOptions([BASIC_ENTITY_TYPES.connector]),
    [ALARM_FIELDS.connectorName]: getEntityOptions([BASIC_ENTITY_TYPES.connector]),
    [ALARM_FIELDS.component]: {
      ...getEntityOptions([BASIC_ENTITY_TYPES.component]),

      itemText: 'component',
    },
    [ALARM_FIELDS.resource]: getEntityOptions([BASIC_ENTITY_TYPES.resource]),
    [ALARM_FIELDS.state]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.higher,
        PATTERN_OPERATORS.lower,
      ],
      values: Object.values(ALARM_STATES).map(value => ({ value, text: t(`common.stateTypes.${value}`) })),
    },
    [ALARM_FIELDS.status]: {
      operators: [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
      ],
      values: Object.values(ALARM_STATUSES).map(value => ({ value, text: t(`common.statusTypes.${value}`) })),
    },
    [ALARM_FIELDS.tags]: {
      operators: [
        PATTERN_OPERATORS.with,
        PATTERN_OPERATORS.without,
      ],
      fetchValues: fetchAlarmTagsListWithoutStore,
      itemText: 'value',
      itemValue: 'value',
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
      itemValue: '_id',
      itemText: 'name',
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

export const useAdvancedSearchEntityAttributes = ({ infosItems }) => {
  const { t } = useI18n();
  const { fetchCategoriesListWithoutStore } = useEntityCategory();
  const { getEntityOptions } = useGetEntityOptions();

  const attributesMap = computed(() => ({
    [ALARM_FIELDS.entityId]: getEntityOptions(),
    [ALARM_FIELDS.entityName]: {
      operators: ENTITY_OPERATORS,
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
    },
    [ALARM_FIELDS.entityComponent]: getEntityOptions([ENTITY_TYPES.component]),
    [ALARM_FIELDS.entityConnector]: getEntityOptions([ENTITY_TYPES.connector]),
    [ALARM_FIELDS.entityImpactLevel]: {},
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
    [ALARM_FIELDS.pbehaviorInfoName]: {
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

export const useAdvancedSearchAttributes = ({
  allowAlarmFields,
  allowEntityFields,
  allowPbehaviorFields,
}) => {
  const { t, tc } = useI18n();

  const { pending: infosPending, items: infosItems } = useEntityInfosKeys();

  const { attributesMap: alarmAttributesMap } = useAdvancedSearchAlarmAttributes({ infosItems });
  const { attributesMap: entityAttributesMap } = useAdvancedSearchEntityAttributes({ infosItems });
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

    const result = Object.entries(ADVANCED_SEARCH_GROUPED_ALARM_FIELDS).reduce((acc, [group, items]) => {
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
      ...ADVANCED_SEARCH_ALARM_ENTITY_FIELDS.map(field => prepareItem(field, unwrappedAllowEntityFields)),

      { header: pbehaviorHeader, value: pbehaviorHeader },
      ...ADVANCED_SEARCH_ALARM_PBEHAVIOR_INFO_FIELDS.map(field => prepareItem(field, unwrappedAllowPbehaviorFields)),
    );

    return result;
  });

  return {
    pending: infosPending,
    attributes,
  };
};

export const useAdvancedSearchValidator = ({ rules }) => {
  const instance = useComponentInstance();

  /**
   * Extends the validator with a custom rule named ADVANCED_SEARCH_VALIDATION_RULE_NAME.
   * This rule is used to validate advanced search criteria based on specific conditions.
   *
   * @description
   * The ADVANCED_SEARCH_VALIDATION_RULE_NAME checks the following conditions:
   * - If the rule has an attribute and is not finished, it returns false.
   * - If the rule does not have an attribute but has a union, it checks the last two rules:
   * - If the last rule does not have an attribute and the second-to-last rule has the same key as the current rule,
   *   it returns false.
   * - Otherwise, it returns true.
   */
  const extendValidatorRule = () => (
    instance.$validator.extend(ADVANCED_SEARCH_VALIDATION_RULE_NAME, ({ rule, finished }) => {
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
