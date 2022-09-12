import { cloneDeep, pick } from 'lodash';

import {
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_ENRICHMENT_AFTER_TYPES,
  EVENT_FILTER_TYPES,
  OLD_PATTERNS_FIELDS,
  PATTERNS_FIELDS,
} from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

/**
 * @typedef { 'enrichment' | 'drop' | 'break' | 'change_entity' } EventFilterType
 */

/**
 * @typedef {
 *  'set_field' |
 *  'set_field_from_template' |
 *  'set_entity_info' |
 *  'set_entity_info_from_template' |
 *  'copy' |
 *  'copy_to_entity_info'
 * } EventFilterActionType
 */

/**
 * @typedef {Object} EventFilterAction
 * @property {string} type
 * @property {string} name
 * @property {string} value
 * @property {string} description
 */

/**
 * @typedef {EventFilterAction} EventFilterActionForm
 */

/**
 * @typedef {Object} EventFilterEnrichmentConfig
 * @property {EventFilterAction[]} actions
 * @property {string} on_success
 * @property {string} on_failure
 */

/**
 * @typedef {Object} EventFilterChangeEntityConfig
 * @property {string} resource
 * @property {string} component
 * @property {string} connector
 * @property {string} connector_name
 */

/**
 * @typedef {EventFilterEnrichmentConfig | EventFilterChangeEntityConfig} EventFilterConfig
 */

/**
 * @typedef {FilterPatterns} EventFilter
 * @property {string} _id
 * @property {EventFilterType} type
 * @property {string} description
 * @property {number} priority
 * @property {boolean} enabled
 * @property {EventFilterConfig} config
 */

/**
 * @typedef {PatternsForm & EventFilter} EventFilterForm
 */

/**
 * @typedef {EventFilterEnrichmentConfig & EventFilterChangeEntityConfig} EventFilterConfigForm
 */

/**
 * Convert event filter to form
 *
 * @param {EventFilterConfig | {}} [eventFilterConfig = {}]
 * @returns {EventFilterConfigForm}
 */
export const eventFilterConfigToForm = (eventFilterConfig = {}) => ({
  actions: eventFilterConfig.actions ? cloneDeep(eventFilterConfig.actions) : [],
  on_success: eventFilterConfig.on_success ?? EVENT_FILTER_ENRICHMENT_AFTER_TYPES.pass,
  on_failure: eventFilterConfig.on_failure ?? EVENT_FILTER_ENRICHMENT_AFTER_TYPES.pass,
  resource: eventFilterConfig.resource ?? '',
  component: eventFilterConfig.component ?? '',
  connector: eventFilterConfig.connector ?? '',
  connector_name: eventFilterConfig.connector_name ?? '',
});

/**
 * Convert event filter to form
 *
 * @param {EventFilter | {}} [eventFilter = {}]
 * @returns {EventFilterForm}
 */
export const eventFilterToForm = (eventFilter = {}) => ({
  _id: eventFilter._id ?? '',
  type: eventFilter.type ?? EVENT_FILTER_TYPES.drop,
  description: eventFilter.description ?? '',
  priority: eventFilter.priority ?? 0,
  enabled: eventFilter.enabled ?? true,
  config: eventFilterConfigToForm(eventFilter.config),
  patterns: filterPatternsToForm(
    eventFilter,
    [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.event],
    [OLD_PATTERNS_FIELDS.patterns],
  ),
});

/**
 * Convert event filter action to form
 *
 * @param {EventFilterAction} eventFilterAction
 * @return {EventFilterActionForm}
 */
export const eventFilterActionToForm = (eventFilterAction = {}) => ({
  type: eventFilterAction.type ?? EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
  name: eventFilterAction.name ?? '',
  value: eventFilterAction.value ?? '',
  description: eventFilterAction.description ?? '',
});

/**
 * Convert form to event filter fields
 *
 * @param {EventFilterForm} eventFilterForm
 * @returns {EventFilter}
 */
export const formToEventFilter = (eventFilterForm) => {
  const { config, patterns, ...eventFilter } = eventFilterForm;

  switch (eventFilterForm.type) {
    case EVENT_FILTER_TYPES.changeEntity:
      eventFilter.config = pick(config, ['resource', 'component', 'connector', 'connector_name']);
      break;
    case EVENT_FILTER_TYPES.enrichment:
      eventFilter.config = pick(config, ['actions', 'on_success', 'on_failure']);
      break;
  }

  return {
    ...eventFilter,
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.event, PATTERNS_FIELDS.entity]),
  };
};