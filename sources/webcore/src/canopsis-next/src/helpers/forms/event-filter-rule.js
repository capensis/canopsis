import { cloneDeep, isUndefined } from 'lodash';
import { EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES, EVENT_FILTER_RULE_TYPES } from '@/constants';

/**
 * Convert event filter rule to form
 *
 * @param {Object} [rule={}]
 * @returns {Object}
 */
export function eventFilterRuleToForm(rule = {}) {
  return {
    general: {
      _id: rule._id || '',
      type: rule.type || EVENT_FILTER_RULE_TYPES.drop,
      description: rule.description || '',
      pattern: rule.pattern ? cloneDeep(rule.pattern) : {},
      priority: rule.priority || 0,
      enabled: !isUndefined(rule.enabled) ? rule.enabled : true,
    },
    enrichmentOptions: {
      actions: rule.actions ? cloneDeep(rule.actions) : [],
      externalData: rule.externalData ? cloneDeep(rule.externalData) : {},
      onSuccess: rule.on_success || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
      onFailure: rule.on_failure || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
    },
  };
}

/**
 * Convert form.enrichmentOptions to event filter rule fields
 *
 * @param {Object} [enrichmentOptions={}]
 * @returns {Object}
 */
export function formEnrichmentOptionsToEventFilterRule(enrichmentOptions = {}) {
  return {
    actions: enrichmentOptions.actions,
    external_data: enrichmentOptions.externalData,
    on_success: enrichmentOptions.onSuccess,
    on_failure: enrichmentOptions.onFailure,
  };
}
