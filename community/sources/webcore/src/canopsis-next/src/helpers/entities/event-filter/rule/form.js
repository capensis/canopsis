import {
  cloneDeep,
  pick,
  isEmpty,
  omitBy,
  omit,
} from 'lodash';
import flatten from 'flat';

import {
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_ENRICHMENT_AFTER_TYPES,
  EVENT_FILTER_EVENT_EXTRA_PREFIX,
  EVENT_FILTER_TYPES,
  PATTERNS_FIELDS,
} from '@/constants';

import { uid } from '@/helpers/uid';
import {
  exceptionsToForm,
  exceptionsToRequest,
  exdatesToForm,
  exdatesToRequest,
  formExceptionsToExceptions,
  formExdatesToExdates,
} from '@/helpers/entities/pbehavior/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import {
  externalDataTemplateVariablesErrorsToForm,
  externalDataToForm,
  formToExternalData,
} from '@/helpers/entities/shared/external-data/form';

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
 * @property {number} [start]
 * @property {number} [stop]
 * @property {boolean} enabled
 * @property {EventFilterConfig} config
 * @property {string} rrule
 * @property {PbehaviorException[]} exceptions
 * @property {PbehaviorExdate[]} exdates
 * @property {ExternalData} external_data
 */

/**
 * @typedef {PatternsForm & EventFilter} EventFilterForm
 * @property {PbehaviorExceptionForm[]} exceptions
 * @property {PbehaviorExdateForm[]} exdates
 * @property {ExternalDataForm} external_data
 */

/**
 * @typedef {EventFilterEnrichmentConfig & EventFilterChangeEntityConfig} EventFilterConfigForm
 */

/**
 * Remove 'Event.ExtraInfos.' prefix from value, if it exists for dictionary type
 *
 * @param {string} [eventFilterActionValue = '']
 * @returns {string}
 */
export const eventFilterDictionaryActionValueToForm = (eventFilterActionValue = '') => (
  eventFilterActionValue.replace(EVENT_FILTER_EVENT_EXTRA_PREFIX, '')
);

/**
 * Convert event filter action to form
 *
 * @param {EventFilterAction} eventFilterAction
 * @return {EventFilterActionForm}
 */
export const eventFilterActionToForm = (eventFilterAction = {}) => ({
  key: uid(),
  type: eventFilterAction.type ?? EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
  name: eventFilterAction.name ?? '',
  value: eventFilterAction.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary
    ? eventFilterDictionaryActionValueToForm(eventFilterAction.value)
    : eventFilterAction.value,
  description: eventFilterAction.description ?? '',
});

/**
 * Convert event filter to form
 *
 * @param {EventFilterConfig | {}} [eventFilterConfig = {}]
 * @returns {EventFilterConfigForm}
 */
export const eventFilterConfigToForm = (eventFilterConfig = {}) => ({
  actions: (eventFilterConfig.actions ? cloneDeep(eventFilterConfig.actions) : []).map(eventFilterActionToForm),
  on_success: eventFilterConfig.on_success ?? EVENT_FILTER_ENRICHMENT_AFTER_TYPES.pass,
  on_failure: eventFilterConfig.on_failure ?? EVENT_FILTER_ENRICHMENT_AFTER_TYPES.pass,
  resource: eventFilterConfig.resource ?? '',
  component: eventFilterConfig.component ?? '',
  connector: eventFilterConfig.connector ?? '',
  connector_name: eventFilterConfig.connector_name ?? '',
});

/**
 * Remove old, if some patterns exist
 *
 * @param {EventFilter} eventFilter
 * @returns {FilterPatterns}
 */
export const eventFilterPatternToForm = eventFilter => filterPatternsToForm(
  eventFilter,
  [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.event],
);

/**
 * Convert event filter to form
 *
 * @param {EventFilter} [eventFilter = {}]
 * @param {string} [timezone]
 * @returns {EventFilterForm}
 */
export const eventFilterToForm = (eventFilter = {}, timezone) => ({
  _id: eventFilter._id ?? '',
  type: eventFilter.type ?? EVENT_FILTER_TYPES.drop,
  description: eventFilter.description ?? '',
  priority: eventFilter.priority,
  enabled: eventFilter.enabled ?? true,
  rrule: eventFilter.rrule ?? null,
  start: eventFilter.start,
  stop: eventFilter.stop,
  exceptions: exceptionsToForm(eventFilter.exceptions),
  exdates: exdatesToForm(eventFilter.exdates, timezone),
  config: eventFilterConfigToForm(eventFilter.config),
  external_data: externalDataToForm(eventFilter.external_data),
  patterns: eventFilterPatternToForm(eventFilter),
});

/**
 * Add 'Event.ExtraInfos.' prefix for value for dictionary type
 *
 * @param {string} [eventFilterActionValue = '']
 * @returns {string}
 */
export const formToEventFilterDictionaryActionValue = (eventFilterActionValue = '') => (
  `${EVENT_FILTER_EVENT_EXTRA_PREFIX}${eventFilterActionValue}`
);

/**
 * Convert from to event filter action fields
 *
 * @param {EventFilterActionForm} eventFilterActionForm
 * @return {EventFilterAction}
 */
export const formToEventFilterAction = eventFilterActionForm => (omit({
  ...eventFilterActionForm,
  value: eventFilterActionForm.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary
    ? formToEventFilterDictionaryActionValue(eventFilterActionForm.value)
    : eventFilterActionForm.value,
}, ['key']));

/**
 * Convert form to event filter fields
 *
 * @param {EventFilterForm} eventFilterForm
 * @param {string} timezone
 * @returns {EventFilter}
 */
export const formToEventFilter = (eventFilterForm, timezone) => {
  const {
    config,
    patterns,
    exdates,
    exceptions,
    external_data: externalData,
    ...eventFilter
  } = eventFilterForm;

  switch (eventFilterForm.type) {
    case EVENT_FILTER_TYPES.changeEntity:
      eventFilter.config = pick(config, ['resource', 'component', 'connector', 'connector_name']);
      break;
    case EVENT_FILTER_TYPES.enrichment:
      eventFilter.config = pick(config, ['on_success', 'on_failure']);
      eventFilter.config.actions = config.actions.map(formToEventFilterAction);
      break;
  }

  if (!isEmpty(externalData)) {
    eventFilter.external_data = formToExternalData(externalData);
  }

  return {
    ...eventFilter,
    exdates: exdatesToRequest(formExdatesToExdates(exdates, timezone)),
    exceptions: exceptionsToRequest(formExceptionsToExceptions(exceptions)),
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.event, PATTERNS_FIELDS.entity]),
  };
};

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @return {FlattenErrors}
 */
export const eventFilterRuleActionTemplateVariablesErrorsToForm = (errorsObject) => {
  const { value } = errorsObject;
  const actionErrors = {};

  if (!value.is_valid) {
    actionErrors.value = value.err.message;
  }

  return actionErrors;
};

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object[]} errorsArray
 * @param {EventFilterAction[]} actions
 * @return {FlattenErrors}
 */
export const eventFilterRuleActionsTemplateVariablesErrorsToForm = (errorsArray, actions) => errorsArray
  .reduce((acc, errors, index) => {
    const action = actions[index];

    acc[action.key] = eventFilterRuleActionTemplateVariablesErrorsToForm(errors);

    return acc;
  }, {});

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @param {EventFilterConfig} form
 * @return {FlattenErrors}
 */
export const eventFilterRuleConfigTemplateVariablesErrorsToForm = (errorsObject, form) => {
  const { actions, component, connector, connector_name: connectorName, resource } = errorsObject ?? {};
  const errors = {};

  if (actions) {
    errors.actions = eventFilterRuleActionsTemplateVariablesErrorsToForm(actions, form.actions);
  }

  if (component && !component.is_valid) {
    errors.component = component.err.message;
  }

  if (connector && !connector?.is_valid) {
    errors.connector = connector.err.message;
  }

  if (connectorName && !connectorName?.is_valid) {
    errors.connector_name = connectorName.err.message;
  }

  if (resource && !resource?.is_valid) {
    errors.resource = resource.err.message;
  }

  return errors;
};

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @param {EventFilterForm} form
 * @return {FlattenErrors}
 */
export const eventFilterRuleTemplateVariablesErrorsToForm = (errorsObject, form) => {
  const { external_data: externalData, config } = errorsObject;
  const errors = {};

  if (externalData) {
    errors.external_data = externalDataTemplateVariablesErrorsToForm(externalData, form.external_data);
  }

  if (config) {
    errors.config = eventFilterRuleConfigTemplateVariablesErrorsToForm(config, form.config);
  }

  return omitBy(flatten(errors), value => isEmpty(value));
};
