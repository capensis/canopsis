export const FILTER_OPERATORS = {
  equal: 'equal',
  notEqual: 'not equal',
  in: 'in',
  notIn: 'not in',
  beginsWith: 'begins with',
  doesntBeginWith: 'doesn\'t begin with',
  contains: 'contains',
  doesntContains: 'doesn\'t contain',
  endsWith: 'ends with',
  doesntEndWith: 'doesn\'t end with',
  isEmpty: 'is empty',
  isEmptyArray: 'is empty array',
  isNotEmpty: 'is not empty',
  isNull: 'is null',
  isNotNull: 'is not null',
  greater: 'greater',
  less: 'less',
};

export const FILTER_OPERATORS_FOR_ARRAY = [FILTER_OPERATORS.in, FILTER_OPERATORS.notIn];

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
