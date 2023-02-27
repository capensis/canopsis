import { EVENT_FILTER_PATTERN_FIELDS } from '@/constants/event-filter';

export const LINK_RULE_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
};

export const LINK_RULE_EXTERNAL_DATA_CONDITION_VALUES = {

};

export const EVENT_FILTER_EXTERNAL_DATA_CONDITION_VALUES = {
  [EVENT_FILTER_PATTERN_FIELDS.component]: {
    text: EVENT_FILTER_PATTERN_FIELDS.component,
    value: '.Event.Component',
  },
  [EVENT_FILTER_PATTERN_FIELDS.connector]: {
    text: EVENT_FILTER_PATTERN_FIELDS.connector,
    value: '.Event.Connector',
  },
  [EVENT_FILTER_PATTERN_FIELDS.connectorName]: {
    text: EVENT_FILTER_PATTERN_FIELDS.connectorName,
    value: '.Event.ConnectorName',
  },
  [EVENT_FILTER_PATTERN_FIELDS.resource]: {
    text: EVENT_FILTER_PATTERN_FIELDS.resource,
    value: '.Event.Resource',
  },
  [EVENT_FILTER_PATTERN_FIELDS.output]: {
    text: EVENT_FILTER_PATTERN_FIELDS.output,
    value: '.Event.Output',
  },
  [EVENT_FILTER_PATTERN_FIELDS.extraInfos]: {
    text: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
    value: '.Event.ExtraInfos',
  },
};
