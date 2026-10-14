import { convertAlarmStateFilterToQuery } from '@/helpers/entities/alarm/query';

/**
 * Convert StatsCalendar widget to query.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const convertStatsCalendarWidgetToQuery = (widget) => {
  const { filters = [], parameters: { considerPbehaviors = false } = {} } = widget;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    considerPbehaviors,
    filters,
    time_field: 't',
  };
};
