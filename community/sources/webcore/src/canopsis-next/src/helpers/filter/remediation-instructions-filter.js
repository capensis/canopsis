/**
 * @typedef {Object} RemediationInstructionFilterQueryParameters
 * @property {boolean} [with_instructions]
 * @property {string[]} [include_instructions]
 * @property {string[]} [exclude_instructions]
 * @property {string[]} [include_types]
 * @property {string[]} [exclude_types]
 */

import { uniq } from 'lodash';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

const PARAMETERS_FILTERS_TYPES = {
  include: 'include',
  exclude: 'exclude',
};

/**
 * Prepare remediation instructions filters to query
 *
 * @param {RemediationInstructionFilter[]} [filters = []]
 * @returns {RemediationInstructionFilterQueryParameters}
 */
export function prepareRemediationInstructionsFiltersToQuery(filters = []) {
  const query = {
    with_instructions: true,
  };

  if (!filters.length) {
    return query;
  }

  const result = filters.reduce((acc, filter) => {
    const key = filter.with ? PARAMETERS_FILTERS_TYPES.include : PARAMETERS_FILTERS_TYPES.exclude;

    if (
      acc[key].types.includes(REMEDIATION_INSTRUCTION_TYPES.autom)
      && acc[key].types.includes(REMEDIATION_INSTRUCTION_TYPES.manual)
    ) {
      return acc;
    }

    if (filter.all) {
      acc[key].types = [REMEDIATION_INSTRUCTION_TYPES.autom, REMEDIATION_INSTRUCTION_TYPES.manual];
    }

    if (filter.auto) {
      acc[key].types.push(REMEDIATION_INSTRUCTION_TYPES.autom);
    }

    if (filter.manual) {
      acc[key].types.push(REMEDIATION_INSTRUCTION_TYPES.manual);
    }

    acc[key].instructions.push(...(filter.instructions || []).map(({ _id }) => _id));

    return acc;
  }, {
    [PARAMETERS_FILTERS_TYPES.include]: { types: [], instructions: [] },
    [PARAMETERS_FILTERS_TYPES.exclude]: { types: [], instructions: [] },
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
 * @param {Widget} widget
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
