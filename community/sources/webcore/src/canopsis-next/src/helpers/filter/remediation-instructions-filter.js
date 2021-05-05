import { uniq } from 'lodash';

import { REMEDIATION_INSTRUCTION_FILTER_ALL, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

const FILTERS_TYPES = {
  include: 'include',
  exclude: 'exclude',
};

/**
 * Prepare remediation instructions filters to query
 *
 * @param {Object|Array} [filters = []]
 * @returns {Object}
 */
export function prepareRemediationInstructionsFiltersToQuery(filters = []) {
  if (!filters.length) {
    return {};
  }

  const query = {
    with_instructions: true,
  };

  const result = filters.reduce((acc, filter) => {
    const key = filter.with ? FILTERS_TYPES.include : FILTERS_TYPES.exclude;

    if (filter.all) {
      acc[key].instructions = [REMEDIATION_INSTRUCTION_FILTER_ALL];
    } else if (!acc[key].instructions.includes(REMEDIATION_INSTRUCTION_FILTER_ALL)) {
      acc[key].instructions.push(...filter.instructions);
    }

    if (filter.automatic) {
      acc[key].types.push(REMEDIATION_INSTRUCTION_TYPES.automatic);
    }

    if (filter.manual) {
      acc[key].types.push(REMEDIATION_INSTRUCTION_TYPES.manual);
    }

    return acc;
  }, {
    [FILTERS_TYPES.include]: { types: [], instructions: [] },
    [FILTERS_TYPES.exclude]: { types: [], instructions: [] },
  });

  Object.entries(result).forEach(([key, value = {}]) => {
    Object.entries(value).forEach(([secondKey, secondValue = []]) => {
      if (secondValue.length) {
        query[`${key}_${secondKey}`] = uniq(secondValue).sort();
      }
    });
  });

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
  } = userPreference.widget_preferences;

  const preparedWidgetFilters = widgetFilters
    .filter(filter => !disabledWidgetFilters.includes(filter._id));

  return [
    ...filters,
    ...preparedWidgetFilters,
  ];
}
