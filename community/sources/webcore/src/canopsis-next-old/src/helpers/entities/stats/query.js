import { convertAlarmStateFilterToQuery } from '@/helpers/entities/alarm/query';

/**
 * This function converts widget with type 'StatsCalendar' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertStatsCalendarWidgetToQuery(widget) {
  const { filters = [], parameters: { considerPbehaviors = false } } = widget;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    considerPbehaviors,
    filters,
    time_field: 't',
  };
}
