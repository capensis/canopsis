export const EVENT_FILTER_TYPES = {
  drop: 'drop',
  break: 'break',
  enrichment: 'enrichment',
  changeEntity: 'change_entity',
};

export const EVENT_FILTER_ENRICHMENT_AFTER_TYPES = {
  pass: 'pass',
  break: 'break',
  drop: 'drop',
};

export const EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES = {
  setField: 'set_field',
  setFieldFromTemplate: 'set_field_from_template',
  setEntityInfoFromTemplate: 'set_entity_info_from_template',
  copy: 'copy',
  setEntityInfo: 'set_entity_info',
  copyToEntityInfo: 'copy_to_entity_info',
};

export const EVENT_FILTER_PATTERN_FIELDS = {
  sourceType: 'source_type',
  eventType: 'event_type',
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
  output: 'output',
  extraInfos: 'extra',
  state: 'state',
  longOutput: 'long_output',
};

export const EVENT_FILTER_SOURCE_TYPES = {
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
};
