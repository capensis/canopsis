export const EXTERNAL_DATA_CONDITION_TYPES = {
  select: 'select',
  regexp: 'regexp',
};

export const EXTERNAL_DATA_TYPES = {
  mongo: 'mongo',
  api: 'api',
};

export const EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS = {
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
  output: 'output',
  extraInfos: 'extra',
};

export const EXTERNAL_DATA_DEFAULT_CONDITION_VALUES = [
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.component,
    value: '.Event.Component',
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connector,
    value: '.Event.Connector',
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connectorName,
    value: '.Event.ConnectorName',
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.resource,
    value: '.Event.Resource',
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.output,
    value: '.Event.Output',
  },
  {
    text: EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.extraInfos,
    value: '.Event.ExtraInfos',
  },
];

export const EXTERNAL_DATA_PAYLOADS_VARIABLES = {
  externalData: '.ExternalData.%reference%.Value',
};
