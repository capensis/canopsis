import {
  isEmpty,
  isUndefined,
  omit,
  map,
  uniq,
  isArray,
} from 'lodash';

import {
  ALARMS_OPENED_VALUES,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  LIVE_REPORTING_QUICK_RANGES,
  SORT_ORDERS,
} from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { isResolvedAlarm } from '@/helpers/entities/alarm/form';
import { convertWidgetChartsToPerfDataQuery } from '@/helpers/entities/metric/query';
import { convertMultiSortToRequest } from '@/helpers/entities/shared/query';
import { getTemplateVariables } from '@/helpers/handlebars/variables';

/**
 *  This function converts widget.parameters.opened to query Object
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
export const getAlarmVariablesByTemplate = template => getTemplateVariables(template)
  .reduce((acc, variable) => {
    if (variable.startsWith('alarm.') || variable.startsWith('entity.')) {
      acc.push(variable.replace(/^(alarm)\./, ''));
    }

    return acc;
  }, []);

/**
 * Get all used into templates alarm variables
 *
 * @param {WidgetColumn[]} widgetColumns
 * @param {string} moreInfoTemplate
 * @param {WidgetInfoPopup[]} infoPopups
 * @param {string[]} usedAlarmProperties
 * @returns {string[]}
 */
export const convertAlarmWidgetParametersToActiveColumns = ({
  widgetColumns,
  moreInfoTemplate,
  infoPopups,
}) => {
  const activeColumns = [];

  widgetColumns.forEach(({ template, value }) => {
    activeColumns.push(value);

    if (template) {
      activeColumns.push(...getAlarmVariablesByTemplate(template));
    }
  });

  if (moreInfoTemplate) {
    activeColumns.push(...getAlarmVariablesByTemplate(moreInfoTemplate));
  }

  if (infoPopups) {
    infoPopups.forEach(({ template }) => {
      activeColumns.push(...getAlarmVariablesByTemplate(template));
    });
  }

  return uniq(activeColumns);
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
    itemsPerPage,
    sort,
    mainFilter,
    opened = ALARMS_OPENED_VALUES.opened,
    usedAlarmProperties,
  } = widget.parameters;

  const query = {
    opened,
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
    multiSortBy: [],
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

  if (sort && sort.column && sort.order) {
    query.multiSortBy.push({
      sortBy: sort.column,
      descending: sort.order === SORT_ORDERS.desc,
    });
  }

  return query;
}

/**
 * This function converts userPreference with widget type 'AlarmsList' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertAlarmUserPreferenceToQuery({ content }) {
  const {
    itemsPerPage,
    category,
    mainFilter,
    onlyBookmarks = false,
    isCorrelationEnabled = false,
  } = content;

  const query = {
    category,
    correlation: isCorrelationEnabled,
    filter: mainFilter,
    only_bookmarks: onlyBookmarks,
  };

  if (itemsPerPage) {
    query.limit = itemsPerPage;
  }

  return query;
}

/**
 * Prepare query for alarm details fetching
 *
 * @param {Alarm} alarm
 * @param {Widget} widget
 * @param {string} search
 * @returns {Object}
 */
export const prepareAlarmDetailsQuery = (alarm, widget, search) => {
  const { sort = {}, widgetGroupColumns = [], charts = [] } = widget.parameters;
  const columns = widgetGroupColumns.length > 0
    ? widgetGroupColumns
    : DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS;

  const query = {
    search,
    search_by: map(columns, 'value'),
    _id: alarm._id,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
    with_dependencies: true,
    opened: isResolvedAlarm(alarm) ? false : widget.parameters.opened,
    perf_data: convertWidgetChartsToPerfDataQuery(charts),
    steps: {
      reversed: true,
      page: 1,
      limit: PAGINATION_LIMIT,
    },
    children: {
      page: 1,
      limit: PAGINATION_LIMIT,
      multiSortBy: [],
    },
  };

  if (sort.column && sort.order && columns.some(({ value }) => value.endsWith(sort.column))) {
    query.children.multiSortBy.push({ sortBy: sort.column, descending: sort.order === SORT_ORDERS.desc });
  }

  return query;
};

/**
 * Convert alarmDetails query to request
 *
 * @param {Object} query
 * @returns {Object}
 */
export const convertAlarmDetailsQueryToRequest = query => ({
  ...query,

  children: {
    ...omit(query.children, ['multiSortBy']),

    multi_sort: convertMultiSortToRequest(query.children?.multiSortBy),
  },
});
