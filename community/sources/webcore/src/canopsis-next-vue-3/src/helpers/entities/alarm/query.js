import { isUndefined, uniq } from 'lodash';

import { ALARMS_OPENED_VALUES } from '@/constants';
import { getTemplateVariables } from '@/helpers/handlebars/variables';

/**
 * This function converts widget.parameters.opened to query Object
 *
 * @param {Object} parameters
 * @returns {{ opened: boolean }}
 */
export function convertAlarmStateFilterToQuery({ parameters }) {
  const { opened = ALARMS_OPENED_VALUES.opened } = parameters;
  const query = {};

  if (!isUndefined(opened)) {
    query.opened = opened;
  }

  return query;
}

/**
 * Collect alarm properties used by widget columns and templates.
 *
 * @param {Object} parameters
 * @returns {string[]}
 */
export const convertAlarmWidgetParametersToActiveColumns = ({
  widgetColumns = [],
  moreInfoTemplate,
  infoPopups,
} = {}) => {
  const activeColumns = widgetColumns.flatMap(({ value, template }) => [
    value,
    ...(template ? getTemplateVariables(template) : []),
  ]);

  if (moreInfoTemplate) {
    activeColumns.push(...getTemplateVariables(moreInfoTemplate));
  }

  (infoPopups ?? []).forEach(({ template }) => {
    activeColumns.push(...getTemplateVariables(template));
  });

  return uniq(activeColumns.filter(Boolean));
};
