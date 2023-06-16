import { isEmpty, isUndefined, omit } from 'lodash';

import { ALARMS_OPENED_VALUES, DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, QUICK_RANGES, SORT_ORDERS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { isResolvedAlarm } from '@/helpers/entities/alarm/form';
import { convertWidgetChartsToPerfDataQuery } from '@/helpers/entities/metric/query';
import { convertMultiSortToRequest } from '@/helpers/entities/shared/query';

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
 * This function converts widget with type 'AlarmsList' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertAlarmWidgetToQuery(widget) {
  const {
    liveReporting = {},
    widgetColumns = [],
    itemsPerPage,
    sort,
    mainFilter,
    opened = ALARMS_OPENED_VALUES.opened,
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
  } else if (query.opened === ALARMS_OPENED_VALUES.resolved) {
    query.tstart = QUICK_RANGES.last30Days.start;
    query.tstop = QUICK_RANGES.last30Days.stop;
  }

  if (widgetColumns.length) {
    query.active_columns = widgetColumns.map(v => v.value);
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
    isCorrelationEnabled = false,
  } = content;

  const query = {
    correlation: isCorrelationEnabled,
    filter: mainFilter,
    category,
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
 * @returns {Object}
 */
export const prepareAlarmDetailsQuery = (alarm, widget) => {
  const { sort = {}, widgetGroupColumns = [], charts = [] } = widget.parameters;
  const columns = widgetGroupColumns.length > 0
    ? widgetGroupColumns
    : DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS;

  const query = {
    _id: alarm._id,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
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
