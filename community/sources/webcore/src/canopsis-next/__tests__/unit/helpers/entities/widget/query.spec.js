import { PAGINATION_LIMIT } from '@/config';
import { ALARMS_OPENED_VALUES, WIDGET_TYPES } from '@/constants';

import * as alarmQuery from '@/helpers/entities/alarm/query';
import { prepareWidgetQuery } from '@/helpers/entities/widget/query';

describe('prepareWidgetQuery', () => {
  const defaultWidget = {
    type: WIDGET_TYPES.alarmList,
    parameters: {
      isCorrelationEnabled: undefined,
    },
  };
  const defaultUserPreference = {
    content: {
      isCorrelationEnabled: undefined,
    },
  };

  const defaultWidgetQuery = {
    opened: ALARMS_OPENED_VALUES.opened,
    itemsPerPage: PAGINATION_LIMIT,
    page: 1,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
    sortBy: [],
    sortDesc: [],
    lockedFilter: undefined,
  };

  const defaultUserPreferenceQuery = {
    only_bookmarks: false,
  };

  const defaultResult = {
    ...defaultWidgetQuery,
    ...defaultUserPreferenceQuery,
  };

  it('should convert widget and userPreference to query objects and merge them', () => {
    vi.spyOn(alarmQuery, 'convertAlarmWidgetToQuery').mockReturnValueOnce(defaultWidgetQuery);
    vi.spyOn(alarmQuery, 'convertAlarmUserPreferenceToQuery').mockReturnValueOnce(
      defaultUserPreferenceQuery,
    );

    const result = prepareWidgetQuery(defaultWidget, defaultUserPreference);

    expect(result).toEqual(defaultResult);
  });

  it('should convert widget and userPreference (with isCorrelationEnabled) to query objects and merge them', () => {
    vi.spyOn(alarmQuery, 'convertAlarmWidgetToQuery').mockReturnValueOnce(defaultWidgetQuery);
    vi.spyOn(alarmQuery, 'convertAlarmUserPreferenceToQuery').mockReturnValueOnce({
      ...defaultUserPreferenceQuery,
      correlation: true,
    });

    const result = prepareWidgetQuery(defaultWidget, defaultUserPreference);

    expect(result).toEqual({
      ...defaultResult,

      correlation: true,
    });
  });

  it('should convert widget (with isCorrelationEnabled) and userPreference to query objects and merge them', () => {
    vi.spyOn(alarmQuery, 'convertAlarmWidgetToQuery').mockReturnValueOnce({
      ...defaultWidgetQuery,
      correlation: true,
    });
    vi.spyOn(alarmQuery, 'convertAlarmUserPreferenceToQuery').mockReturnValueOnce(
      defaultUserPreferenceQuery,
    );

    const result = prepareWidgetQuery(defaultWidget, defaultUserPreference);

    expect(result).toEqual({
      ...defaultResult,

      correlation: true,
    });
  });

  it('should convert widget (with isCorrelationEnabled) and userPreference (with disabled isCorrelationEnabled) to query objects and merge them', () => {
    vi.spyOn(alarmQuery, 'convertAlarmWidgetToQuery').mockReturnValueOnce({
      ...defaultWidgetQuery,
      correlation: true,
    });
    vi.spyOn(alarmQuery, 'convertAlarmUserPreferenceToQuery').mockReturnValueOnce({
      ...defaultUserPreferenceQuery,
      correlation: null,
    });

    const result = prepareWidgetQuery(defaultWidget, defaultUserPreference);

    expect(result).toEqual({
      ...defaultResult,

      correlation: null,
    });
  });
});
