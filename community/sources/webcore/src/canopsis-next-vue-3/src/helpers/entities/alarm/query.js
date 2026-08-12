import { isArray, isEmpty, isUndefined, uniq } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  ALARMS_OPENED_VALUES,
  LIVE_REPORTING_QUICK_RANGES,
  PBEHAVIOR_ORIGINS,
  SORT_ORDERS,
} from '@/constants';
import { getTemplateVariables } from '@/helpers/handlebars/variables';
import { sanitizeHtml } from '@/helpers/html';

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
 * Gat available alarm and entity variables
 *
 * @param {string} template
 * @returns {string[]}
 */
export const getAlarmVariablesByTemplate = template => getTemplateVariables(sanitizeHtml(template))
  .reduce((acc, variable) => {
    if (variable.startsWith('alarm.') || variable.startsWith('entity.')) {
      acc.push(variable.replace(/^(alarm)\./, ''));
    }

    return acc;
  }, []);

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
    ...(template ? getAlarmVariablesByTemplate(template) : []),
  ]);

  if (moreInfoTemplate) {
    activeColumns.push(...getAlarmVariablesByTemplate(moreInfoTemplate));
  }

  (infoPopups ?? []).forEach(({ template }) => {
    activeColumns.push(...getAlarmVariablesByTemplate(template));
  });

  return uniq(activeColumns.filter(Boolean));
};

/**
 * This function converts widget with type 'AlarmsList' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertAlarmWidgetToQuery(widget) {
  const {
    liveReporting = {},
    itemsPerPage = PAGINATION_LIMIT,
    opened = ALARMS_OPENED_VALUES.opened,
    mainFilter,
    usedAlarmProperties,
    isCorrelationEnabled,
    sort = [],
  } = widget.parameters;

  const query = {
    opened,
    itemsPerPage,

    page: 1,
    with_instructions: true,
    with_tag_colors: true,
    with_declare_tickets: true,
    with_links: true,
    pbh_origin: PBEHAVIOR_ORIGINS.alarmList,
    sortBy: [],
    sortDesc: [],
    lockedFilter: mainFilter,
  };

  if (!isEmpty(liveReporting)) {
    query.tstart = liveReporting.tstart;
    query.tstop = liveReporting.tstop;
    query.time_field = liveReporting.time_field;
  } else if (query.opened === ALARMS_OPENED_VALUES.resolved) {
    query.tstart = LIVE_REPORTING_QUICK_RANGES.last30Days.start;
    query.tstop = LIVE_REPORTING_QUICK_RANGES.last30Days.stop;
  }

  const activeColumns = isArray(usedAlarmProperties)
    ? usedAlarmProperties
    : convertAlarmWidgetParametersToActiveColumns(widget.parameters);

  if (activeColumns.length) {
    query.active_columns = activeColumns;
  }

  if (isArray(sort) && sort.length > 0) {
    const lowerCasedDesc = SORT_ORDERS.desc.toLowerCase();

    query.sortBy = sort.map(column => column.sort_by).filter(Boolean);
    query.sortDesc = sort.map(column => column.sort === lowerCasedDesc);
  }

  if (!isUndefined(isCorrelationEnabled)) {
    query.correlation = isCorrelationEnabled;
  }

  return query;
}

/**
 * This function converts userPreference with widget type 'AlarmsList' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertAlarmUserPreferenceToQuery({ content = {} } = {}) {
  const {
    itemsPerPage,
    category,
    mainFilter,
    isCorrelationEnabled,
    instructionsFilter,
    tags,
    onlyBookmarks = false,
  } = content;

  const query = {
    category,
    filter: mainFilter,
    only_bookmarks: onlyBookmarks,
    instructionsFilter: instructionsFilter ?? {},
    tags,
  };

  if (!isUndefined(isCorrelationEnabled)) {
    query.correlation = isCorrelationEnabled;
  }

  if (itemsPerPage) {
    query.itemsPerPage = itemsPerPage;
  }

  return query;
}
