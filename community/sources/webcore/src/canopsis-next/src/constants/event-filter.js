import { PATTERN_CONDITIONS, PATTERNS_FIELDS } from './pattern';
import { EVENT_TYPES } from './event';

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
  copy: 'copy',
  setTags: 'set_tags',
  setTagsFromTemplate: 'set_tags_from_template',
  copyToEntityInfo: 'copy_to_entity_info',
  setEntityInfo: 'set_entity_info',
  setEntityInfoFromTemplate: 'set_entity_info_from_template',
  setEntityInfoFromDictionary: 'set_entity_info_from_dictionary',
};

export const EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_GROUPS = {
  setEventField: [
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
  ],
  setTags: [
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate,
  ],
  setEntityInfo: [
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary,
  ],
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
  author: 'author',
  initiator: 'initiator',
};

export const EVENT_FILTER_SOURCE_TYPES = {
  component: 'component',
  connector: 'connector',
  connectorName: 'connector_name',
  resource: 'resource',
};

export const EVENT_FILTER_FAILURE_TYPES = {
  invalidPattern: 0,
  invalidTemplate: 1,
  externalDataMongo: 2,
  externalDataApi: 3,
  other: 4,
};

export const EVENT_FILTER_SET_FIELD_ITEMS = [
  'Output',
  'LongOutput',
];

export const EVENT_FILTER_SET_TAGS_FIELDS = [
  EVENT_FILTER_PATTERN_FIELDS.output,
  EVENT_FILTER_PATTERN_FIELDS.extraInfos,
];

export const EVENT_FILTER_SET_TAGS_VALUE_PREFIXES = {
  [EVENT_FILTER_PATTERN_FIELDS.output]: 'Event.Output',
  [EVENT_FILTER_PATTERN_FIELDS.extraInfos]: 'Event.ExtraInfo.',
};

export const EVENT_FILTER_SET_TAGS_REGEX = /<value>.*<name>|<name>.*<value>/;

export const EVENT_FILTER_EVENT_EXTRA_PREFIX = 'Event.ExtraInfos.';

export const EVENT_FILTER_EXPAND_PANEL_TABS = {
  summary: 'summary',
  pattern: 'pattern',
  action: 'action',
  externalData: 'externalData',
  errors: 'errors',
};

export const EVENT_FILTER_DEFAULT_PATTERN = {
  [PATTERNS_FIELDS.event]: [[
    {
      field: EVENT_FILTER_PATTERN_FIELDS.eventType,
      cond: {
        value: EVENT_TYPES.check,
        type: PATTERN_CONDITIONS.equal,
      },
    },
  ]],
};
