import { uniq } from 'lodash';

import { REMEDIATION_INSTRUCTION_FILTER_ALL } from '@/constants';

/**
 * Prepare remediation instructions filters to query
 *
 * @param {Object|Array} [filters = []]
 * @returns {Object}
 */
export function prepareRemediationInstructionsFiltersToQuery(filters = []) {
  const query = {};

  const result = filters.reduce((acc, filter) => {
    const key = filter.with ? 'with' : 'without';

    if (filter.all) {
      acc[key] = [REMEDIATION_INSTRUCTION_FILTER_ALL];
    } else if (!acc[key].includes(REMEDIATION_INSTRUCTION_FILTER_ALL)) {
      acc[key].push(...filter.instructions);
    }

    return acc;
  }, { with: [], without: [] });

  if (result.with.length) {
    query.with_instructions = uniq(result.with).sort().join(',');
  }

  if (result.without.length) {
    query.without_instructions = uniq(result.without).sort().join(',');
  }

  return query;
}

/**
 * Get all enabled remediation instructions filters for user on widget
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Array}
 */
export function getRemediationInstructionsFilters(widget, userPreference) {
  const { remediationInstructionsFilters: widgetFilters = [] } = widget.parameters;
  const {
    disabledWidgetRemediationInstructionsFilters: disabledWidgetFilters = [],
    remediationInstructionsFilters: filters = [],
  } = userPreference.content;

  const preparedWidgetFilters = widgetFilters
    .filter(filter => !disabledWidgetFilters.includes(filter._id));

  return [
    ...filters,
    ...preparedWidgetFilters,
  ];
}
