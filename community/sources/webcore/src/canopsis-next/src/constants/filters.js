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

export const FILTER_INPUT_TYPES = {
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
    inputType: FILTER_INPUT_TYPES.string,
  },
  group: {
    condition: FILTER_MONGO_OPERATORS.and,
    groups: {},
    rules: {},
  },
};
