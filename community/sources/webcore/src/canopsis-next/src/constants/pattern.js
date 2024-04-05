import { QUICK_RANGES } from './common';

export const PATTERN_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
  pbehavior: 'pbehavior',
};

export const PATTERN_OPERATORS = {
  equal: 'equal',
  contains: 'contains',
  notEqual: 'not_equal',
  notContains: 'not_contains',

  beginsWith: 'begins_with',
  notBeginWith: 'not_begin_with',
  endsWith: 'ends_with',
  notEndWith: 'not_end_with',

  exist: 'exist',
  notExist: 'not_exist',

  hasEvery: 'has_every',
  hasOneOf: 'has_one_of',
  isOneOf: 'is_one_of',
  hasNot: 'has_not',
  isNotOneOf: 'is_not_one_of',
  isEmpty: 'is_empty',
  isNotEmpty: 'is_not_empty',

  higher: 'higher_than',
  lower: 'lower_than',

  longer: 'longer',
  shorter: 'shorter',

  ticketAssociated: 'ticket_associated',
  ticketNotAssociated: 'ticket_not_associated',

  canceled: 'canceled',
  notCanceled: 'not_canceled',

  snoozed: 'snoozed',
  notSnoozed: 'not_snoozed',

  acked: 'acked',
  notAcked: 'not_acked',

  isGrey: 'is_grey',
  isNotGrey: 'is_not_grey',

  with: 'with',
  without: 'without',

  activated: 'activated',
  inactive: 'inactive',

  isMetaAlarm: 'is_meta_alarm',
  isNotMetaAlarm: 'is_not_meta_alarm',

  regexp: 'regexp',
};

export const PATTERN_CONDITIONS = {
  equal: 'eq',
  notEqual: 'neq',
  greater: 'gt',
  less: 'lt',
  regexp: 'regexp',
  hasEvery: 'has_every',
  isOneOf: 'is_one_of',
  hasOneOf: 'has_one_of',
  hasNot: 'has_not',
  isNotOneOf: 'is_not_one_of',
  isEmpty: 'is_empty',
  exist: 'exist',
  relativeTime: 'relative_time',
  absoluteTime: 'absolute_time',
  contains: 'contain',
  notContains: 'not_contain',
  beginsWith: 'begin_with',
  notBeginWith: 'not_begin_with',
  endsWith: 'end_with',
  notEndWith: 'not_end_with',
};

export const PATTERN_FIELD_TYPES = {
  string: 'string',
  number: 'int',
  boolean: 'bool',
  null: 'null',
  stringArray: 'string_array',
};

export const PATTERN_OPERATORS_WITHOUT_VALUE = [
  PATTERN_OPERATORS.exist,
  PATTERN_OPERATORS.notExist,
  PATTERN_OPERATORS.isEmpty,
  PATTERN_OPERATORS.isNotEmpty,
  PATTERN_OPERATORS.ticketAssociated,
  PATTERN_OPERATORS.ticketNotAssociated,
  PATTERN_OPERATORS.acked,
  PATTERN_OPERATORS.notAcked,
  PATTERN_OPERATORS.snoozed,
  PATTERN_OPERATORS.notSnoozed,
  PATTERN_OPERATORS.canceled,
  PATTERN_OPERATORS.notCanceled,
  PATTERN_OPERATORS.isNotGrey,
  PATTERN_OPERATORS.isGrey,
  PATTERN_OPERATORS.activated,
  PATTERN_OPERATORS.inactive,
  PATTERN_OPERATORS.isMetaAlarm,
  PATTERN_OPERATORS.isNotMetaAlarm,
];

export const PATTERN_RULE_TYPES = {
  infos: 'infos',
  extraInfos: 'extraInfos',
  date: 'date',
  duration: 'duration',
  string: 'string',
  number: 'number',
  object: 'object',
};

export const PATTERN_RULE_INFOS_FIELDS = {
  value: 'value',
  name: 'name',
};

export const PATTERN_ARRAY_OPERATORS = [
  PATTERN_OPERATORS.hasEvery,
  PATTERN_OPERATORS.hasOneOf,
  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.hasNot,
  PATTERN_OPERATORS.isNotOneOf,
  PATTERN_OPERATORS.isEmpty,
  PATTERN_OPERATORS.isNotEmpty,
];

export const PATTERN_DURATION_OPERATORS = [
  PATTERN_OPERATORS.longer,
  PATTERN_OPERATORS.shorter,
];

export const PATTERN_NUMBER_OPERATORS = [
  PATTERN_OPERATORS.equal,
  PATTERN_OPERATORS.notEqual,
  PATTERN_OPERATORS.higher,
  PATTERN_OPERATORS.lower,
];

export const PATTERN_STRING_OPERATORS = [
  PATTERN_OPERATORS.equal,
  PATTERN_OPERATORS.notEqual,
  PATTERN_OPERATORS.contains,
  PATTERN_OPERATORS.notContains,
  PATTERN_OPERATORS.beginsWith,
  PATTERN_OPERATORS.notBeginWith,
  PATTERN_OPERATORS.endsWith,
  PATTERN_OPERATORS.notEndWith,
  PATTERN_OPERATORS.regexp,
];

export const PATTERN_BOOLEAN_OPERATORS = [
  PATTERN_OPERATORS.equal,
  PATTERN_OPERATORS.notEqual,
];

export const PATTERN_NULL_OPERATORS = [
  PATTERN_OPERATORS.equal,
  PATTERN_OPERATORS.notEqual,
];

export const PATTERN_EXISTS_OPERATORS = [
  PATTERN_OPERATORS.exist,
  PATTERN_OPERATORS.notExist,
];

export const PATTERN_QUICK_RANGES = [
  QUICK_RANGES.last15Minutes,
  QUICK_RANGES.last30Minutes,
  QUICK_RANGES.last1Hour,
  QUICK_RANGES.last3Hour,
  QUICK_RANGES.last3Hour,
  QUICK_RANGES.last6Hour,
  QUICK_RANGES.last12Hour,
  QUICK_RANGES.last24Hour,
  QUICK_RANGES.last2Days,
  QUICK_RANGES.last7Days,
  QUICK_RANGES.last30Days,
  QUICK_RANGES.last1Year,
  QUICK_RANGES.custom,
];

export const PATTERN_CUSTOM_ITEM_VALUE = Symbol('custom');

export const PATTERN_TABS = {
  patterns: 'patterns',
  corporatePatterns: 'corporatePatterns',
};

export const PATTERN_EDITOR_TABS = {
  simple: 'simple',
  advanced: 'advanced',
};

export const PATTERNS_FIELDS = {
  alarm: 'alarm_pattern',
  entity: 'entity_pattern',
  pbehavior: 'pbehavior_pattern',
  event: 'event_pattern',
  totalEntity: 'total_entity_pattern',
  serviceWeather: 'weather_service_pattern',
};
