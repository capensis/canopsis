import { convertAlarmStateFilterToQuery } from '@/helpers/entities/alarm/query';

/**
 * This function converts widget with type 'counter' widget to query Object
 *
 * @param widget
 * @returns {{filters: *}}
 */
export function convertCounterWidgetToQuery(widget) {
  const { filters = [], parameters: { isCorrelationEnabled = false } } = widget;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    correlation: isCorrelationEnabled,
    filters: filters.map(({ _id: id }) => id),
  };
}
