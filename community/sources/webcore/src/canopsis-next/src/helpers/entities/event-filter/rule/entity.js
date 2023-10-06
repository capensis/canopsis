import { EVENT_FILTER_TYPES } from '@/constants';

/**
 * Check event filter rule type is enrichment
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isEnrichmentEventFilterRuleType = type => type === EVENT_FILTER_TYPES.enrichment;

/**
 * Check event filter rule type is change entity
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isChangeEntityEventFilterRuleType = type => type === EVENT_FILTER_TYPES.changeEntity;
