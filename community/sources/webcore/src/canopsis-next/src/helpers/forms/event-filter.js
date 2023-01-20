import { cloneDeep, pick, isEmpty, omit } from 'lodash';

import {
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_ENRICHMENT_AFTER_TYPES,
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES,
  EVENT_FILTER_EXTERNAL_DATA_TYPES,
  EVENT_FILTER_TYPES,
  OLD_PATTERNS_FIELDS,
  PATTERNS_FIELDS,
} from '@/constants';

import uid from '@/helpers/uid';

import { filterPatternsToForm, formFilterToPatterns } from './filter';
import {
  exceptionsToForm,
  exceptionsToRequest,
  exdatesToForm,
  exdatesToRequest,
  formExceptionsToExceptions,
  formExdatesToExdates,
} from './planning-pbehavior';
import { formToRequest, requestToForm } from './shared/request';

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
 * @typedef {'mongo' | 'api'} EventFilterEnrichmentExternalDataType
 */

/**
 * @typedef {'select' | 'regexp'} EventFilterEnrichmentExternalDataConditionType
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
 * @typedef {
 *   Object<EventFilterEnrichmentExternalDataConditionType, Object<string, string>>
 * } EventFilterEnrichmentExternalDataCondition
 */

/**
 * @typedef {EventFilterEnrichmentExternalDataCondition} EventFilterEnrichmentExternalData
 * @property {EventFilterEnrichmentExternalDataType} type
 * @property {string} collection
 */

/**
 * @typedef {Object} EventFilterEnrichmentExternalDataConditionForm
 * @property {string} key
 * @property {EventFilterEnrichmentExternalDataConditionType} type
 * @property {string} attribute
 * @property {string} value
 */

/**
 * @typedef {Object} EventFilterEnrichmentExternalDataFormItem
 * @property {string} key
 * @property {string} reference
 * @property {EventFilterEnrichmentExternalDataType} type
 * @property {string} collection
 * @property {string} sort_by
 * @property {string} sort
 * @property {EventFilterEnrichmentExternalDataConditionForm[]} conditions
 */

/**
 * @typedef {EventFilterEnrichmentExternalDataFormItem[]} EventFilterEnrichmentExternalDataForm
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
 * @property {EventFilterEnrichmentExternalData} external_data
 */

/**
 * @typedef {PatternsForm & EventFilter} EventFilterForm
 * @property {PbehaviorExceptionForm[]} exceptions
 * @property {PbehaviorExdateForm[]} exdates
 * @property {EventFilterEnrichmentExternalDataForm} external_data
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

export const eventFilterExternalDataConditionItemToForm = (
  conditionType = EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES.select,
  attribute = '',
  value = '',
) => ({
  key: uid(),
  type: conditionType,
  attribute,
  value,
});

/**
 * Convert event filter enrichment external data condition to form
 *
 * @param {EventFilterEnrichmentExternalDataConditionType} conditionType
 * @param {EventFilterEnrichmentExternalDataCondition} condition
 * @returns {EventFilterEnrichmentExternalDataConditionForm[]}
 */
export const eventFilterExternalDataConditionToForm = (conditionType, condition) => (
  Object.entries(condition)
    .map(([attribute, value]) => eventFilterExternalDataConditionItemToForm(conditionType, attribute, value))
);

/**
 * Convert event filter enrichment external data conditions to form
 *
 * @param {EventFilterEnrichmentExternalDataConditionType[]} [conditionTypes = []]
 * @param {EventFilterEnrichmentExternalDataCondition} [item = {}]
 * @returns {EventFilterEnrichmentExternalDataConditionForm[]}
 */
export const eventFilterExternalDataConditionsToForm = (conditionTypes = [], item = {}) => {
  const conditions = conditionTypes
    .reduce((acc, conditionType) => {
      const condition = item[conditionType];

      if (condition) {
        acc.push(...eventFilterExternalDataConditionToForm(conditionType, condition));
      }

      return acc;
    }, []);

  if (!conditions.length) {
    conditions.push(eventFilterExternalDataConditionItemToForm());
  }

  return conditions;
};

export const eventFilterExternalDataItemToForm = (
  reference = '',
  item = { type: EVENT_FILTER_EXTERNAL_DATA_TYPES.mongo },
) => ({
  key: uid(),
  reference,
  type: item.type,

  request: requestToForm(item.request),
  sort_by: item.sort_by,
  sort: item.sort,
  collection: item.collection ?? '',
  conditions: eventFilterExternalDataConditionsToForm(
    Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES),
    item,
  ),
});

/**
 * Convert event filter enrichment external data to form
 *
 * @param {EventFilterEnrichmentExternalData} [externalData = {}]
 * @returns {EventFilterEnrichmentExternalDataForm}
 */
export const eventFilterEnrichmentExternalDataToForm = (externalData = {}) => (
  Object.entries(externalData).map(([reference, item]) => eventFilterExternalDataItemToForm(reference, item))
);

/**
 * Convert form to event filter enrichment external data conditions
 *
 * @param {EventFilterEnrichmentExternalDataConditionForm[]} form
 * @returns {EventFilterEnrichmentExternalDataCondition}
 */
export const formToEventFilterEnrichmentExternalDataConditions = (form = []) => (
  form.reduce((acc, { type, attribute, value }) => {
    if (!acc[type]) {
      acc[type] = {};
    }

    acc[type][attribute] = value;

    return acc;
  }, {})
);

/**
 * Convert form to event filter enrichment external data
 *
 * @param {EventFilterEnrichmentExternalDataForm} form
 * @returns {EventFilterEnrichmentExternalData}
 */
export const formToEventFilterEnrichmentExternalData = (form = []) => (
  form.reduce((acc, externalData) => {
    const { type, reference } = externalData;

    const additionalFields = type === EVENT_FILTER_EXTERNAL_DATA_TYPES.api
      ? {
        request: formToRequest(externalData.request),
      } : {
        ...pick(externalData, ['sort', 'sort_by', 'collection']),
        ...formToEventFilterEnrichmentExternalDataConditions(externalData.conditions),
      };

    acc[reference] = {
      type,

      ...additionalFields,
    };

    return acc;
  }, {})
);

/**
 * Remove old, if some patterns exist
 *
 * @param {EventFilter} eventFilter
 * @returns {FilterPatterns}
 */
export const eventFilterPatternToForm = (eventFilter) => {
  const entityPattern = eventFilter[PATTERNS_FIELDS.entity];
  const eventPattern = eventFilter[PATTERNS_FIELDS.event];

  const eventFilterWithPatterns = entityPattern?.length || eventPattern?.length
    ? omit(eventFilter, [OLD_PATTERNS_FIELDS.patterns])
    : eventFilter;

  return filterPatternsToForm(
    eventFilterWithPatterns,
    [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.event],
    [OLD_PATTERNS_FIELDS.patterns, OLD_PATTERNS_FIELDS.patterns],
  );
};

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
  priority: eventFilter.priority ?? 0,
  enabled: eventFilter.enabled ?? true,
  rrule: eventFilter.rrule ?? null,
  start: eventFilter.start,
  stop: eventFilter.stop,
  exceptions: exceptionsToForm(eventFilter.exceptions),
  exdates: exdatesToForm(eventFilter.exdates, timezone),
  config: eventFilterConfigToForm(eventFilter.config),
  external_data: !isEmpty(eventFilter.external_data)
    ? eventFilterEnrichmentExternalDataToForm(eventFilter.external_data)
    : [],
  patterns: eventFilterPatternToForm(eventFilter),
});

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
  value: eventFilterAction.value ?? '',
  description: eventFilterAction.description ?? '',
});

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
      eventFilter.config = pick(config, ['actions', 'on_success', 'on_failure']);
      break;
  }

  if (!isEmpty(externalData)) {
    eventFilter.external_data = formToEventFilterEnrichmentExternalData(externalData);
  }

  return {
    ...eventFilter,
    exdates: exdatesToRequest(formExdatesToExdates(exdates, timezone)),
    exceptions: exceptionsToRequest(formExceptionsToExceptions(exceptions)),
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.event, PATTERNS_FIELDS.entity]),
  };
};
