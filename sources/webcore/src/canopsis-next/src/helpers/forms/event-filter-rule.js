import { cloneDeep, isUndefined } from 'lodash';
import { EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES, EVENT_FILTER_RULE_TYPES } from '@/constants';

/**
 * Convert event filter rule to form
 *
 * @param {Object} [rule={}]
 * @returns {Object}
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
    externalData: rule.external_data ? cloneDeep(rule.external_data) : {},
    onSuccess: rule.on_success || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
    onFailure: rule.on_failure || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
  },
});

/**
 * Convert form.enrichmentOptions to event filter rule fields
 *
 * @param {Object} [enrichmentOptions={}]
 * @returns {Object}
 */
export const formEnrichmentOptionsToEventFilterRule = (enrichmentOptions = {}) => ({
  actions: enrichmentOptions.actions,
  external_data: enrichmentOptions.externalData,
  on_success: enrichmentOptions.onSuccess,
  on_failure: enrichmentOptions.onFailure,
});

/**
 * Convert form to event filter rule fields
 *
 * @param {Object} form
 * @returns {Object}
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
