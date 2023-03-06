export const EXTERNAL_DATA_CONDITION_TYPES = {
  select: 'select',
  regexp: 'regexp',
};

export const EXTERNAL_DATA_TYPES = {
  mongo: 'mongo',
  api: 'api',
};

export const EXTERNAL_DATA_CONDITION_FIELDS = {
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
  output: 'output',
  extraInfos: 'extra',
};

export const EXTERNAL_DATA_CONDITION_VALUES = {
  [EXTERNAL_DATA_CONDITION_FIELDS.component]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.component,
    value: '.Event.Component',
  },
  [EXTERNAL_DATA_CONDITION_FIELDS.connector]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.connector,
    value: '.Event.Connector',
  },
  [EXTERNAL_DATA_CONDITION_FIELDS.connectorName]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.connectorName,
    value: '.Event.ConnectorName',
  },
  [EXTERNAL_DATA_CONDITION_FIELDS.resource]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.resource,
    value: '.Event.Resource',
  },
  [EXTERNAL_DATA_CONDITION_FIELDS.output]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.output,
    value: '.Event.Output',
  },
  [EXTERNAL_DATA_CONDITION_FIELDS.extraInfos]: {
    text: EXTERNAL_DATA_CONDITION_FIELDS.extraInfos,
    value: '.Event.ExtraInfos',
  },
};

export const EXTERNAL_DATA_PAYLOADS_VARIABLES = {
  externalData: '.ExternalData.%reference%.Value',
};
