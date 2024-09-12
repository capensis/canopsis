import { isBoolean } from 'lodash';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { mapIds } from '@/helpers/array';

/**
 * @typedef {Object} RemediationInstructionFilterQuery
 * @property {boolean} [running]
 * @property {string[]} [include]
 * @property {string[]} [exclude]
 * @property {string[]} [include_types]
 * @property {string[]} [exclude_types]
 */

/**
 * @typedef {Object} RemediationInstructionFilterQueryParameters
 * @property {RemediationInstructionFilterQuery[]} [instructions]
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

  if (!filters.length) {
    return query;
  }

  query.instructions = filters.map((filter) => {
    const instructionQuery = {};

    const key = filter.with ? PARAMETERS_FILTERS_TYPES.include : PARAMETERS_FILTERS_TYPES.exclude;
    const typesKey = `${key}_types`;

    if (filter.all) {
      instructionQuery[typesKey] = [
        REMEDIATION_INSTRUCTION_TYPES.auto,
        REMEDIATION_INSTRUCTION_TYPES.manual,
        REMEDIATION_INSTRUCTION_TYPES.simpleManual,
      ];
    } else if (filter.auto || filter.manual) {
      instructionQuery[typesKey] = [];

      if (filter.auto) {
        instructionQuery[typesKey].push(REMEDIATION_INSTRUCTION_TYPES.auto);
      }

      if (filter.manual) {
        instructionQuery[typesKey].push(
          REMEDIATION_INSTRUCTION_TYPES.manual,
          REMEDIATION_INSTRUCTION_TYPES.simpleManual,
        );
      }
    }

    if (filter.instructions?.length) {
      instructionQuery[key] = mapIds(filter.instructions);
    }

    if (isBoolean(filter.running)) {
      instructionQuery.running = filter.running;
    }

    return instructionQuery;
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
