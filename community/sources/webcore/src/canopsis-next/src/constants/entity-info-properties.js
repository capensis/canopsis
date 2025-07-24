export const ENTITY_INFO_PROPERTY_TYPES = {
  boolean: 0,
  number: 1,
  timestamp: 2,
  string: 3,
  string_array: 4,
};

export const ENTITY_INFO_PROPERTY_TYPE_VALUES = Object.keys(ENTITY_INFO_PROPERTY_TYPES);

export const ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS = {
  [ENTITY_INFO_PROPERTY_TYPES.boolean]: 'common.variableTypes.boolean',
  [ENTITY_INFO_PROPERTY_TYPES.number]: 'common.variableTypes.number',
  [ENTITY_INFO_PROPERTY_TYPES.timestamp]: 'common.timestamp',
  [ENTITY_INFO_PROPERTY_TYPES.string]: 'common.variableTypes.string',
  [ENTITY_INFO_PROPERTY_TYPES.string_array]: 'common.variableTypes.stringArray',
};
