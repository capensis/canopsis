import { PATTERN_OPERATORS, PATTERN_STRING_OPERATORS, PATTERNS_FIELDS } from './pattern';
import { ALARM_PATTERN_FIELDS } from './alarm';
import { ENTITY_PATTERN_FIELDS } from './entity';
import { PBEHAVIOR_PATTERN_PREFIX, PBEHAVIOR_PATTERN_FIELDS, PBEHAVIOR_FIELDS } from './pbehavior';
import { DYNAMIC_INFO_FIELDS } from './dynamic-info';

export const REGISTER_LAST_INPUT_FOCUS_KEY = '$registerLastInputFocus';

export const ADVANCED_SEARCH_FIELDS = {
  ...PATTERNS_FIELDS,

  search: 'search_pattern',
};

export const ADVANCED_SEARCH_UNION_CONDITIONS = {
  and: 'AND',
  or: 'OR',
};

export const ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS = [
  PBEHAVIOR_PATTERN_FIELDS.name,
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
  rangeValuePeriod: 'rangeValuePeriod',
  rangeValueDate: 'rangeValueDate',
  union: 'union',
  text: 'text',
};

export const ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS = [
  ...PATTERN_STRING_OPERATORS,

  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.isNotOneOf,
];

export const ADVANCED_SEARCH_USER_OPERATORS = [
  PATTERN_OPERATORS.equal,
  PATTERN_OPERATORS.notEqual,
  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.isNotOneOf,
];

export const ALARM_ADVANCED_SEARCH_ENTITY_FIELDS = [
  ALARM_PATTERN_FIELDS.entityId,
  ALARM_PATTERN_FIELDS.entityName,
  ALARM_PATTERN_FIELDS.entityCategoryName,
  ALARM_PATTERN_FIELDS.entityType,
  ALARM_PATTERN_FIELDS.entityComponent,
  ALARM_PATTERN_FIELDS.entityConnector,
  ALARM_PATTERN_FIELDS.entityImpactLevel,
  ALARM_PATTERN_FIELDS.entityInfos,
  ALARM_PATTERN_FIELDS.entityComponentInfos,
];

export const ALARM_ADVANCED_SEARCH_PBEHAVIOR_INFO_FIELDS = [
  ALARM_PATTERN_FIELDS.pbehaviorInfoName,
  ALARM_PATTERN_FIELDS.pbehaviorInfoReason,
  ALARM_PATTERN_FIELDS.pbehaviorInfoType,
  ALARM_PATTERN_FIELDS.pbehaviorInfoCanonicalType,
];

export const ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME = 'advancedSearchRule';

export const ADVANCED_SEARCH_QUERY_FIELDS = [
  'search',
  ADVANCED_SEARCH_FIELDS.alarm,
  ADVANCED_SEARCH_FIELDS.entity,
  ADVANCED_SEARCH_FIELDS.pbehavior,
  ADVANCED_SEARCH_FIELDS.search,
];

export const ADVANCED_SEARCH_FIELDS_TO_COMPARISON = [
  ...ADVANCED_SEARCH_QUERY_FIELDS,

  'positions',
];

export const ALARM_SEARCH_NUMBER_ATTRIBUTES = [
  ALARM_PATTERN_FIELDS.totalStateChanges,
  ALARM_PATTERN_FIELDS.entityImpactLevel,
];

export const ENTITY_SEARCH_NUMBER_ATTRIBUTES = [
  ENTITY_PATTERN_FIELDS.impactLevel,
  ENTITY_PATTERN_FIELDS.impactState,
  ENTITY_PATTERN_FIELDS.koEvents,
  ENTITY_PATTERN_FIELDS.okEvents,
];

export const PBEHAVIOR_SEARCH_NUMBER_ATTRIBUTES = [
  PBEHAVIOR_FIELDS.alarmCount,
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
  alarms: 'alarms',
};

export const ALARM_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    ALARM_PATTERN_FIELDS.displayName,
    ALARM_PATTERN_FIELDS.connector,
    ALARM_PATTERN_FIELDS.connectorName,
    ALARM_PATTERN_FIELDS.component,
    ALARM_PATTERN_FIELDS.resource,
    ALARM_PATTERN_FIELDS.state,
    ALARM_PATTERN_FIELDS.status,
    ALARM_PATTERN_FIELDS.tags,
    ALARM_PATTERN_FIELDS.infos,
    ALARM_PATTERN_FIELDS.meta,
    ALARM_PATTERN_FIELDS.changeState,
    ALARM_PATTERN_FIELDS.totalStateChanges,
  ],
  [ADVANCED_SEARCH_GROUPS.alias]: [],
  [ADVANCED_SEARCH_GROUPS.messages]: [
    ALARM_PATTERN_FIELDS.output,
    ALARM_PATTERN_FIELDS.longOutput,
    ALARM_PATTERN_FIELDS.initialOutput,
    ALARM_PATTERN_FIELDS.initialLongOutput,
    ALARM_PATTERN_FIELDS.lastComment,
    ALARM_PATTERN_FIELDS.lastCommentInitiator,
  ],
  [ADVANCED_SEARCH_GROUPS.ticket]: [
    ALARM_PATTERN_FIELDS.ticketMessage,
    ALARM_PATTERN_FIELDS.ticketInitiator,
    ALARM_PATTERN_FIELDS.ticketValue,
    ALARM_PATTERN_FIELDS.ticket,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    ALARM_PATTERN_FIELDS.creationDate,
    ALARM_PATTERN_FIELDS.lastUpdateDate,
    ALARM_PATTERN_FIELDS.lastEventDate,
    ALARM_PATTERN_FIELDS.ackAt,
    ALARM_PATTERN_FIELDS.resolved,
    ALARM_PATTERN_FIELDS.activationDate,
    ALARM_PATTERN_FIELDS.duration,
  ],
  [ADVANCED_SEARCH_GROUPS.actions]: [
    ALARM_PATTERN_FIELDS.ack,
    ALARM_PATTERN_FIELDS.ackBy,
    ALARM_PATTERN_FIELDS.ackMessage,
    ALARM_PATTERN_FIELDS.ackInitiator,
    ALARM_PATTERN_FIELDS.canceled,
    ALARM_PATTERN_FIELDS.canceledInitiator,
    ALARM_PATTERN_FIELDS.activated,
    ALARM_PATTERN_FIELDS.snooze,
  ],
};

export const ENTITY_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    ENTITY_PATTERN_FIELDS.id,
    ENTITY_PATTERN_FIELDS.name,
    ENTITY_PATTERN_FIELDS.type,
    ENTITY_PATTERN_FIELDS.categoryName,
    ENTITY_PATTERN_FIELDS.component,
    ENTITY_PATTERN_FIELDS.connector,
    ENTITY_PATTERN_FIELDS.impactLevel,
    ENTITY_PATTERN_FIELDS.impactState,
    ENTITY_PATTERN_FIELDS.importSource,
    ENTITY_PATTERN_FIELDS.state,
    ENTITY_PATTERN_FIELDS.status,
    ENTITY_PATTERN_FIELDS.infos,
    ENTITY_PATTERN_FIELDS.componentInfos,
  ],
  [ADVANCED_SEARCH_GROUPS.alias]: [],
  [ADVANCED_SEARCH_GROUPS.events]: [
    ENTITY_PATTERN_FIELDS.koEvents,
    ENTITY_PATTERN_FIELDS.okEvents,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    ENTITY_PATTERN_FIELDS.idleSince,
    ENTITY_PATTERN_FIELDS.imported,
    ENTITY_PATTERN_FIELDS.lastAlarmUpdateDate,
    ENTITY_PATTERN_FIELDS.lastPbehaviorDate,
    ENTITY_PATTERN_FIELDS.lastEventDate,
  ],
};

export const PBEHAVIOR_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    PBEHAVIOR_FIELDS.name,
    PBEHAVIOR_FIELDS.author,
    PBEHAVIOR_FIELDS.enabled,
    PBEHAVIOR_FIELDS.rrule,
    PBEHAVIOR_FIELDS.reason,
    PBEHAVIOR_FIELDS.type,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    PBEHAVIOR_FIELDS.tstart,
    PBEHAVIOR_FIELDS.tstop,
    PBEHAVIOR_FIELDS.rruleEnd,
    PBEHAVIOR_FIELDS.created,
    PBEHAVIOR_FIELDS.updated,
    PBEHAVIOR_FIELDS.lastAlarmDate,
  ],
  [ADVANCED_SEARCH_GROUPS.alarms]: [
    PBEHAVIOR_FIELDS.alarmCount,
  ],
};

export const AVAILABILITY_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    ENTITY_PATTERN_FIELDS.customId,
    ENTITY_PATTERN_FIELDS.name,
    ENTITY_PATTERN_FIELDS.type,
    ENTITY_PATTERN_FIELDS.categoryName,
    ENTITY_PATTERN_FIELDS.component,
    ENTITY_PATTERN_FIELDS.connector,
    ENTITY_PATTERN_FIELDS.impactLevel,
    ENTITY_PATTERN_FIELDS.enabled,
    ENTITY_PATTERN_FIELDS.infos,
    ENTITY_PATTERN_FIELDS.componentInfos,
  ],
  [ADVANCED_SEARCH_GROUPS.alias]: [],
};

export const DYNAMIC_INFO_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    DYNAMIC_INFO_FIELDS.id,
    DYNAMIC_INFO_FIELDS.name,
    DYNAMIC_INFO_FIELDS.author,
    DYNAMIC_INFO_FIELDS.enabled,
  ],
  [ADVANCED_SEARCH_GROUPS.dates]: [
    DYNAMIC_INFO_FIELDS.created,
    DYNAMIC_INFO_FIELDS.updated,
  ],
};

export const ENTITY_DEPENDENCIES_GROUPED_ADVANCED_SEARCH_FIELDS = {
  [ADVANCED_SEARCH_GROUPS.basic]: [
    ENTITY_PATTERN_FIELDS.name,
    ENTITY_PATTERN_FIELDS.type,
  ],
};

export const ADVANCED_SEARCH_GROUPS_TO_PATTERNS = {
  [ADVANCED_SEARCH_GROUPS.basic]: ADVANCED_SEARCH_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.messages]: ADVANCED_SEARCH_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.ticket]: ADVANCED_SEARCH_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.dates]: ADVANCED_SEARCH_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.actions]: ADVANCED_SEARCH_FIELDS.alarm,
  [ADVANCED_SEARCH_GROUPS.entity]: ADVANCED_SEARCH_FIELDS.entity,
  [ADVANCED_SEARCH_GROUPS.pbehavior]: ADVANCED_SEARCH_FIELDS.pbehavior,
  [ADVANCED_SEARCH_GROUPS.search]: ADVANCED_SEARCH_FIELDS.search,
};

export const ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS = Object.entries(ADVANCED_SEARCH_GROUPS_TO_PATTERNS)
  .reduce((acc, [group, patternField]) => {
    ALARM_GROUPED_ADVANCED_SEARCH_FIELDS[group]?.forEach?.(field => acc[field] = patternField);

    return acc;
  }, {});
