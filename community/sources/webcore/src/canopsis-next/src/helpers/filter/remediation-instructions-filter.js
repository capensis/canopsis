import { isBoolean, uniq } from 'lodash';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

/**
 * @typedef {Object} RemediationInstructionFilterQueryParameters
 * @property {boolean} [with_instructions]
 * @property {string[]} [include_instructions]
 * @property {string[]} [exclude_instructions]
 * @property {string[]} [include_instructions_types]
 * @property {string[]} [exclude_instructions_types]
 */

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
  const query = {};

  query.has_running_execution = filters.reduce((acc, { has_running: hasRunning }) => {
    if (isBoolean(hasRunning) && acc !== false) {
      return hasRunning;
    }

    return acc;
  }, undefined);

  if (!filters.length) {
    return query;
  }

  const result = filters.reduce((acc, filter) => {
    const key = filter.with
      ? PARAMETERS_FILTERS_TYPES.include
      : PARAMETERS_FILTERS_TYPES.exclude;

    if (
      acc[key].instruction_types.includes(REMEDIATION_INSTRUCTION_TYPES.auto)
      && acc[key].instruction_types.includes(REMEDIATION_INSTRUCTION_TYPES.manual)
    ) {
      return acc;
    }

    if (filter.all) {
      acc[key].instruction_types = [REMEDIATION_INSTRUCTION_TYPES.auto, REMEDIATION_INSTRUCTION_TYPES.manual];
    }

    if (filter.auto) {
      acc[key].instruction_types.push(REMEDIATION_INSTRUCTION_TYPES.auto);
    }

    if (filter.manual) {
      acc[key].instruction_types.push(REMEDIATION_INSTRUCTION_TYPES.manual);
    }

    const instructionsIds = filter.instructions
      ? filter.instructions.map(({ _id }) => _id)
      : [];

    acc[key].instructions.push(...instructionsIds);

    return acc;
  }, {
    [PARAMETERS_FILTERS_TYPES.include]: { instruction_types: [], instructions: [] },
    [PARAMETERS_FILTERS_TYPES.exclude]: { instruction_types: [], instructions: [] },
  });

  Object.entries(result).forEach(([filterTypeKey, filterTypeRules = {}]) => {
    Object.entries(filterTypeRules).forEach(([filterTypeRuleKey, filterTypeRule = []]) => {
      if (filterTypeRule.length) {
        query[`${filterTypeKey}_${filterTypeRuleKey}`] = uniq(filterTypeRule).sort();
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
  } = userPreference.content;

  const preparedWidgetFilters = widgetFilters
    .filter(filter => !disabledWidgetFilters.includes(filter._id));

  return [
    ...filters,
    ...preparedWidgetFilters,
  ];
}
