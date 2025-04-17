import { PATTERN_OPERATORS, PATTERN_STRING_OPERATORS, PATTERNS_FIELDS } from '@/constants/pattern';
import { ALARM_ADVANCED_SEARCH_GROUPS, ALARM_FIELDS } from '@/constants/alarm';

export const ADVANCED_SEARCH_ITEM_TYPES = {
  field: 'field',
  condition: 'condition',
  value: 'value',
  union: 'union',
};

export const ADVANCED_SEARCH_NEXT_ITEM_TYPES = {
  [ADVANCED_SEARCH_ITEM_TYPES.field]: ADVANCED_SEARCH_ITEM_TYPES.condition,
  [ADVANCED_SEARCH_ITEM_TYPES.condition]: ADVANCED_SEARCH_ITEM_TYPES.value,
  [ADVANCED_SEARCH_ITEM_TYPES.value]: ADVANCED_SEARCH_ITEM_TYPES.union,
  [ADVANCED_SEARCH_ITEM_TYPES.union]: ADVANCED_SEARCH_ITEM_TYPES.field,
};

export const ADVANCED_SEARCH_NOT = 'NOT';

export const ADVANCED_SEARCH_UNION_CONDITIONS = {
  and: 'AND',
  or: 'OR',
};

export const ADVANCED_SEARCH_PATTERNS_PREFIXES = {
  entity: 'entity.',
  pbehavior: 'v.pbehavior_info.',
};

export const ADVANCED_SEARCH_CONDITIONS = {
  less: '<',
  more: '>',
  equal: '=',
  notEqual: '!=',
  like: 'LIKE',
  contains: 'CONTAINS',
};

export const ADVANCED_SEARCH_UNION_REGEXP_PATTERN = new RegExp(`\\s(${Object.values(ADVANCED_SEARCH_UNION_CONDITIONS).join('|')})(\\s|$)`, 'gi');

export const ADVANCED_SEARCH_UNION_FIELDS = [
  {
    value: ADVANCED_SEARCH_UNION_CONDITIONS.and,
    type: ADVANCED_SEARCH_ITEM_TYPES.union,
    text: ADVANCED_SEARCH_UNION_CONDITIONS.and,
  },
  {
    value: ADVANCED_SEARCH_UNION_CONDITIONS.or,
    type: ADVANCED_SEARCH_ITEM_TYPES.union,
    text: ADVANCED_SEARCH_UNION_CONDITIONS.or,
  },
];

export const ALARM_ADVANCED_SEARCH_CHIP_TYPES = {
  attribute: 'attribute',
  dictionary: 'dictionary',
  operator: 'operator',
  fieldType: 'fieldType',
  value: 'value',
  duration: 'duration',
  range: 'range',
  rangeValue: 'rangeValue',
  union: 'union',
  text: 'text',
};

export const ALARM_ADVANCED_SEARCH_ENTITY_OPERATORS = [
  ...PATTERN_STRING_OPERATORS,

  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.isNotOneOf,
];

export const ALARM_ADVANCED_SEARCH_GROUPS_GROUPED = {
  [ALARM_ADVANCED_SEARCH_GROUPS.basic]: [
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

export const ALARM_ADVANCED_SEARCH_ALARM_ENTITY_FIELDS = [
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

export const ALARM_ADVANCED_SEARCH_ALARM_PBEHAVIOR_INFO_FIELDS = [
  ALARM_FIELDS.pbehaviorInfoId,
  ALARM_FIELDS.pbehaviorInfoReason,
  ALARM_FIELDS.pbehaviorInfoType,
  ALARM_FIELDS.pbehaviorInfoCanonicalType,
];

export const ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME = 'advancedSearchRule';

export const ALARM_SEARCH_FIELDS_TO_COMPARISON = [
  'search',
  'positions',

  PATTERNS_FIELDS.alarm,
  PATTERNS_FIELDS.entity,
  PATTERNS_FIELDS.pbehavior,
];

export const ALARM_SEARCH_NUMBER_ATTRIBUTES = [
  ALARM_FIELDS.totalStateChanges,
  ALARM_FIELDS.entityImpactLevel,
];
