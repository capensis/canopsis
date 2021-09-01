import { cloneDeep, isUndefined, pick } from 'lodash';
import {
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES,
  EVENT_FILTER_RULE_TYPES,
} from '@/constants';

/**
 * @typedef { 'enrichment' | 'drop' | 'break' } EventFilterRuleType
 */

/**
 * @typedef {Object} EventFilterRuleSetFieldAction
 * @property {'set_field'} type
 * @property {string} name
 * @property {string|number|boolean} value
 */

/**
 * @typedef {EventFilterRuleSetFieldAction} EventFilterRuleSetFieldFromTemplateAction
 * @property {'set_field_from_template'} type
 */

/**
 * @typedef {EventFilterRuleSetFieldAction} EventFilterRuleSetEntityInfoAction
 * @property {'set_entity_info'} type
 * @property {string} description
 */

/**
 * @typedef {EventFilterRuleSetEntityInfoAction} EventFilterRuleSetEntityInfoFromTemplateAction
 * @property {'set_entity_info_from_template'} type
 */

/**
 * @typedef {Object} EventFilterRuleCopyAction
 * @property {'copy'} type
 * @property {string} from
 * @property {string} to
 */

/**
 * @typedef {Object} EventFilterRuleCopyToEntityInfoAction
 * @property {'copy_to_entity_info'} type
 * @property {string} from
 * @property {string} name
 * @property {string} description
 */

/**
 * @typedef {
 *   EventFilterRuleSetFieldAction |
 *   EventFilterRuleSetFieldFromTemplateAction |
 *   EventFilterRuleSetEntityInfoAction |
 *   EventFilterRuleSetEntityInfoFromTemplateAction |
 *   EventFilterRuleCopyAction |
 *   EventFilterRuleCopyToEntityInfoAction
 * } EventFilterRuleAction
 */

/**
 * @typedef {
 *   EventFilterRuleSetFieldAction &
 *   EventFilterRuleSetFieldFromTemplateAction &
 *   EventFilterRuleSetEntityInfoAction &
 *   EventFilterRuleSetEntityInfoFromTemplateAction &
 *   EventFilterRuleCopyAction &
 *   EventFilterRuleCopyToEntityInfoAction
 * } EventFilterRuleActionForm
 */

/**
 * @typedef {Object} EventFilterRuleEnrichment
 * @property {EventFilterRuleAction[]} actions
 * @property {Object} external_data
 * @property {string} on_success
 * @property {string} on_failure
 */

/**
 * @typedef {Object} EventFilterRuleGeneral
 * @property {string} _id
 * @property {EventFilterRuleType} type
 * @property {string} description
 * @property {Array} patterns
 * @property {number} priority
 * @property {boolean} enabled
 */

/**
 * @typedef {EventFilterRuleEnrichment & EventFilterRuleGeneral} EventFilterRule
 */

/**
 * @typedef {EventFilterRuleGeneral} EventFilterRuleGeneralForm
 */

/**
 * @typedef {EventFilterRuleEnrichment} EventFilterRuleEnrichmentForm
 */

/**
 * @typedef {Object} EventFilterRuleForm
 * @property {EventFilterRuleGeneralForm} general
 * @property {EventFilterRuleEnrichmentForm} enrichmentOptions
 */

/**
 * @typedef { 'enrichment' | 'drop' | 'break' } EventFilterRuleType
 */

/**
 * @typedef {
 *    'set_field' |
 *    'set_field_from_template' |
 *    'set_entity_info_from_template' |
 *    'copy'
 * } EventFilterRuleActionType
 */

/**
 * @typedef {Object} EventFilterRuleAction
 * @property {string} name
 * @property {EventFilterRuleActionType} type
 * @property {string} value
 */

/**
 * @typedef {Object} EventFilterRuleEnrichment
 * @property {EventFilterRuleAction[]} actions
 * @property {Object} external_data
 * @property {string} on_success
 * @property {string} on_failure
 */

/**
 * @typedef {Object} EventFilterRuleGeneral
 * @property {string} _id
 * @property {EventFilterRuleType} type
 * @property {string} description
 * @property {Array} patterns
 * @property {number} priority
 * @property {boolean} enabled
 */

/**
 * @typedef {EventFilterRuleEnrichment & EventFilterRuleGeneral} EventFilterRule
 */

/**
 * @typedef {EventFilterRuleGeneral} EventFilterRuleGeneralForm
 */

/**
 * @typedef {EventFilterRuleEnrichment} EventFilterRuleEnrichmentForm
 */

/**
 * @typedef {Object} EventFilterRuleForm
 * @property {EventFilterRuleGeneralForm} general
 * @property {EventFilterRuleEnrichmentForm} enrichmentOptions
 */

/**
 * Convert event filter rule to form
 *
 * @param {EventFilterRule} [rule={}]
 * @returns {EventFilterRuleForm}
 */
export const eventFilterRuleToForm = (rule = {}) => ({
  general: {
    _id: rule._id || '',
    type: rule.type || EVENT_FILTER_RULE_TYPES.drop,
    description: rule.description || '',
    patterns: rule.patterns ? cloneDeep(rule.patterns) : [],
    priority: rule.priority || 0,
    enabled: !isUndefined(rule.enabled) ? rule.enabled : true,
  },
  enrichmentOptions: {
    actions: rule.actions ? cloneDeep(rule.actions) : [],
    external_data: rule.external_data ? cloneDeep(rule.external_data) : {},
    on_success: rule.on_success || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
    on_failure: rule.on_failure || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
  },
});

/**
 * Convert event filter rule action to form
 *
 * @param {EventFilterRuleAction} [ruleAction = {}]
 * @return {EventFilterRuleActionForm}
 */
export const eventFilterRuleActionToForm = (ruleAction = {}) => ({
  type: ruleAction.type || EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
  name: ruleAction.name || '',
  value: !isUndefined(ruleAction.value) ? ruleAction.value : '',
  description: ruleAction.description || '',
  from: ruleAction.from || '',
  to: ruleAction.to || '',
});

/**
 * Convert form to event filter rule action fields
 *
 * @param  {EventFilterRuleActionForm} form
 * @return {EventFilterRuleAction}
 */
export const formToEventFilterRuleAction = (form) => {
  const fields = ['type'];

  switch (form.type) {
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo:
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate:
      fields.push('description');
    // eslint-disable-next-line no-fallthrough
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField:
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate:
      fields.push('name', 'value');
      break;
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy:
      fields.push('from', 'to');
      break;
    case EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo:
      fields.push('from', 'name', 'description');
      break;
    default:
  }

  return pick(form, fields);
};

/**
 * Convert form.enrichmentOptions to event filter rule fields
 *
 * @param {EventFilterRuleEnrichmentForm} [enrichmentOptions={}]
 * @returns {EventFilterRuleEnrichment}
 */
export const formEnrichmentOptionsToEventFilterRule = (enrichmentOptions = {}) => ({
  actions: enrichmentOptions.actions,
  external_data: enrichmentOptions.external_data,
  on_success: enrichmentOptions.on_success,
  on_failure: enrichmentOptions.on_failure,
});

/**
 * Convert form to event filter rule fields
 *
 * @param {EventFilterRuleForm} form
 * @returns {EventFilterRule}
 */
export const formToEventFilterRule = ({ general, enrichmentOptions }) => {
  const isEnrichmentType = general.type === EVENT_FILTER_RULE_TYPES.enrichment;

  if (isEnrichmentType) {
    return {
      ...general,
      ...formEnrichmentOptionsToEventFilterRule(enrichmentOptions),
    };
  }

  return { ...general };
};
