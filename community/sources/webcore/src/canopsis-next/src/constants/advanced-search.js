import { PATTERN_OPERATORS, PATTERN_STRING_OPERATORS, PATTERNS_FIELDS } from './pattern';
import { ALARM_FIELDS } from './alarm';
import { ENTITY_FIELDS } from './entity';
import { PBEHAVIOR_PATTERN_PREFIX, PBEHAVIOR_PATTERN_FIELDS } from './pbehavior';

export const ADVANCED_SEARCH_UNION_CONDITIONS = {
  and: 'AND',
  or: 'OR',
};

export const ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS = [
  PBEHAVIOR_PATTERN_FIELDS.id,
  PBEHAVIOR_PATTERN_FIELDS.reason,
  PBEHAVIOR_PATTERN_FIELDS.type,
  PBEHAVIOR_PATTERN_FIELDS.canonicalType,
];

export const ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX = 'v.';

export const ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES = {
  entity: 'entity.',
  pbehavior: `${ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX}${PBEHAVIOR_PATTERN_PREFIX}`,
};

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

export const ALARM_ADVANCED_SEARCH_ENTITY_FIELDS = [
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

export const ALARM_ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS = [
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

export const ADVANCED_SEARCH_GROUPS = {
  basic: 'basic',
  messages: 'messages',
  ticket: 'ticket',
  dates: 'dates',
  actions: 'actions',
  alias: 'alias',
  entity: 'entity',
  events: 'events',
  pbehavior: 'pbehavior',
};

export const ALARM_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
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
  [ADVANCED_SEARCH_GROUPS.alias]: [],
  [ADVANCED_SEARCH_GROUPS.messages]: [
    ALARM_FIELDS.output,
    ALARM_FIELDS.longOutput,
    ALARM_FIELDS.initialOutput,
    ALARM_FIELDS.initialLongOutput,
    ALARM_FIELDS.lastComment,
    ALARM_FIELDS.lastCommentInitiator,
  ],
  [ADVANCED_SEARCH_GROUPS.ticket]: [
    ALARM_FIELDS.ticketMessage,
    ALARM_FIELDS.ticketInitiator,
    ALARM_FIELDS.ticketValue,
    ALARM_FIELDS.ticket,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    ALARM_FIELDS.creationDate,
    ALARM_FIELDS.lastUpdateDate,
    ALARM_FIELDS.lastEventDate,
    ALARM_FIELDS.ackAt,
    ALARM_FIELDS.resolved,
    ALARM_FIELDS.activationDate,
    ALARM_FIELDS.duration,
  ],
  [ADVANCED_SEARCH_GROUPS.actions]: [
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

export const ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.entity]: [
    ENTITY_FIELDS.id,
    ENTITY_FIELDS.name,
    ENTITY_FIELDS.type,
    ENTITY_FIELDS.categoryName,
    ENTITY_FIELDS.component,
    ENTITY_FIELDS.connector,
    ENTITY_FIELDS.resource,
    ENTITY_FIELDS.impactLevel,
    ENTITY_FIELDS.impactState,
    ENTITY_FIELDS.state,
    ENTITY_FIELDS.status,
    ENTITY_FIELDS.infos,
    ENTITY_FIELDS.componentInfos,
  ],
  [ADVANCED_SEARCH_GROUPS.alias]: [],
  [ADVANCED_SEARCH_GROUPS.events]: [
    ENTITY_FIELDS.koEvents,
    ENTITY_FIELDS.okEvents,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    ENTITY_FIELDS.idleSince,
    ENTITY_FIELDS.imported,
    ENTITY_FIELDS.lastUpdateDate,
    ENTITY_FIELDS.lastPbehaviorDate,
    ENTITY_FIELDS.lastEventDate,
  ],
};

export const ADVANCED_SEARCH_GROUPS_TO_PATTERNS = {
  [ADVANCED_SEARCH_GROUPS.basic]: PATTERNS_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.messages]: PATTERNS_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.ticket]: PATTERNS_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.dates]: PATTERNS_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.actions]: PATTERNS_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.entity]: PATTERNS_FIELDS.entity,
  [ADVANCED_SEARCH_GROUPS.pbehavior]: PATTERNS_FIELDS.pbehavior,
};

export const ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS = Object.entries(ADVANCED_SEARCH_GROUPS_TO_PATTERNS)
  .reduce((acc, [group, patternField]) => {
    ALARM_GROUPED_ADVANCED_SEARCH_FIELDS[group]?.forEach?.(field => acc[field] = patternField);

    return acc;
  }, {});
