export const FILTER_OPERATORS = {
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
  hasNot: 'has_not',
  isEmpty: 'is_empty',
  isNotEmpty: 'is_not_empty',

  isNull: 'is_null',
  isNotNull: 'is_not_null',

  greater: 'greater',
  less: 'less',

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
};

export const FILTER_CONDITIONS = {
  equal: 'eq',
  notEqual: 'neq',
  greater: 'gt',
  less: 'lt',
  regex: 'regex',
  hasEvery: 'has_every',
  hasOneOf: 'has_one_of',
  hasNot: 'has_not',
  isEmpty: 'is_empty',
  exist: 'exist',
  relativeTime: 'relative_time',
  absoluteTime: 'absolute_time',
};

export const FILTER_OPERATORS_FOR_ARRAY = [
  FILTER_OPERATORS.hasEvery,
  FILTER_OPERATORS.hasOneOf,
  FILTER_OPERATORS.hasNot,
  FILTER_OPERATORS.isEmpty,
  FILTER_OPERATORS.isNotEmpty,
];

export const FILTER_MONGO_OPERATORS = {
  and: '$and',
  or: '$or',
  equal: '$eq',
  notEqual: '$ne',
  in: '$in',
  notIn: '$nin',
  regex: '$regex',
  greater: '$gt',
  less: '$lt',
};

export const PATTERN_INPUT_TYPES = {
  string: 'string',
  number: 'number',
  boolean: 'boolean',
  null: 'null',
  array: 'array',
};

export const FILTER_DEFAULT_VALUES = {
  condition: FILTER_MONGO_OPERATORS.and,
  rule: {
    field: '',
    operator: '',
    input: '',
    inputType: PATTERN_INPUT_TYPES.string,
  },
  group: {
    condition: FILTER_MONGO_OPERATORS.and,
    groups: {},
    rules: {},
  },
};

export const PATTERN_OPERATORS_WITHOUT_VALUE = [
  FILTER_OPERATORS.exist,
  FILTER_OPERATORS.notExist,
  FILTER_OPERATORS.isEmpty,
  FILTER_OPERATORS.isNotEmpty,
  FILTER_OPERATORS.ticketAssociated,
  FILTER_OPERATORS.ticketNotAssociated,
  FILTER_OPERATORS.acked,
  FILTER_OPERATORS.notAcked,
  FILTER_OPERATORS.canceled,
  FILTER_OPERATORS.notCanceled,
];

export const PATTERN_RULE_TYPES = {
  infos: 'infos',
  date: 'date',
  duration: 'duration',
  number: 'number',
  string: 'string',
};

export const PATTERN_RULE_INFOS_FIELDS = {
  value: 'value',
  name: 'name',
};

export const PATTERN_ARRAY_OPERATORS = [
  FILTER_OPERATORS.hasEvery,
  FILTER_OPERATORS.hasOneOf,
  FILTER_OPERATORS.hasNot,
  FILTER_OPERATORS.isEmpty,
  FILTER_OPERATORS.isNotEmpty,
];

export const PATTERN_DURATION_OPERATORS = [
  FILTER_OPERATORS.longer,
  FILTER_OPERATORS.shorter,
];

export const PATTERN_NUMBER_OPERATORS = [
  FILTER_OPERATORS.equal,
  FILTER_OPERATORS.notEqual,
  FILTER_OPERATORS.higher,
  FILTER_OPERATORS.longer,
];

export const PATTERN_STRING_OPERATORS = [
  FILTER_OPERATORS.equal,
  FILTER_OPERATORS.contains,
  FILTER_OPERATORS.notEqual,
  FILTER_OPERATORS.notContains,
  FILTER_OPERATORS.beginsWith,
  FILTER_OPERATORS.notBeginWith,
  FILTER_OPERATORS.endsWith,
  FILTER_OPERATORS.notEndWith,
];

export const PATTERN_BOOLEAN_OPERATORS = [
  FILTER_OPERATORS.equal,
  FILTER_OPERATORS.notEqual,
];

export const PATTERN_NULL_OPERATORS = [
  FILTER_OPERATORS.equal,
  FILTER_OPERATORS.notEqual,
];

export const PATTERN_INFOS_NAME_OPERATORS = [
  FILTER_OPERATORS.exist,
  FILTER_OPERATORS.notExist,
];
