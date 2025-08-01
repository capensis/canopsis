import { INFOS_NAME_VARIABLE } from './common';

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
  extraInfos: `index .Event.ExtraInfos "${INFOS_NAME_VARIABLE}"`,
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

export const EXTERNAL_DATA_TABLE_COLUMN_TYPES = {
  noType: 0,
  filter: 1,
  context: 2,
};

export const EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS = {
  [EXTERNAL_DATA_TABLE_COLUMN_TYPES.noType]: 'grey darken-1',
  [EXTERNAL_DATA_TABLE_COLUMN_TYPES.filter]: 'warning',
  [EXTERNAL_DATA_TABLE_COLUMN_TYPES.context]: 'success',
};
export const MAX_EXTERNAL_DATA_TABLE_TOOLTIP_LINKED_RULES_COUNT = 5;
