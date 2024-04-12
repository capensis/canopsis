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
