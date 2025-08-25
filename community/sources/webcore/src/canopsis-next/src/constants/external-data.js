export const EXTERNAL_DATA_CONDITION_TYPES = {
  select: 'select',
  regexp: 'regexp',
};

export const EXTERNAL_DATA_TYPES = {
  api: 'api',
  table: 'table',
};

export const EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS = {
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
  output: 'output',
  extraInfos: 'extra',
};

export const EXTERNAL_DATA_PAYLOADS_VARIABLES = {
  component: '.Event.Component',
  connector: '.Event.Connector',
  connectorName: '.Event.ConnectorName',
  resource: '.Event.Resource',
  output: '.Event.Output',
  extraInfos: 'index .Event.ExtraInfos "%infos_name%"',
  externalData: '.ExternalData.%reference%',
  regexp: '.RegexMatch.%field%.%name%',
};

export const ACTION_COPY_PAYLOAD_VARIABLES = {
  connector: 'Event.Connector',
  connectorName: 'Event.ConnectorName',
  component: 'Event.Component',
  resource: 'Event.Resource',
  output: 'Event.Output',
  extraInfos: 'Event.ExtraInfos.',
  regexMatch: 'RegexMatch.',
  externalData: 'ExternalData.',
};

export const EXTERNAL_DATA_DEFAULT_CONDITION_VALUES = [
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.component,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.component,
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connector,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.connector,
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connectorName,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.connectorName,
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.resource,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.resource,
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.output,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.output,
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.extraInfos,
    value: EXTERNAL_DATA_PAYLOADS_VARIABLES.extraInfos,
  },
];

export const EXTERNAL_DATA_TABLES_TYPES = {
  mongo: 0,
  postgres: 1,
};

export const EXTERNAL_DATA_TABLE_COLUMN_TAGS = {
  noType: 0,
  filter: 1,
  context: 2,
};

export const EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS = {
  [EXTERNAL_DATA_TABLE_COLUMN_TAGS.noType]: 'grey darken-1',
  [EXTERNAL_DATA_TABLE_COLUMN_TAGS.filter]: 'warning',
  [EXTERNAL_DATA_TABLE_COLUMN_TAGS.context]: 'success',
};
export const MAX_EXTERNAL_DATA_TABLE_TOOLTIP_LINKED_RULES_COUNT = 5;

export const EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES = {
  string: 1,
  boolean: 2,
  number: 3,
  stringArray: 4,
  datetime: 5,
  timestamp: 6,
};

export const EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_DECIMAL_SEPARATOR = {
  comma: 'comma',
  dot: 'dot',
};

export const EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_THOUSANDS_SEPARATOR = {
  ...EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_DECIMAL_SEPARATOR,

  space: 'space',
};

export const EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES = {
  json: 1,
  custom: 2,
};

export const EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS = {
  comma: ',',
  semicolon: ';',
};

export const EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR = 'custom';
