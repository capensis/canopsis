import { PAGINATION_LIMIT } from '@/config';
import { ALARMS_OPENED_VALUES } from '@/constants';

import { convertAlarmUserPreferenceToQuery, convertAlarmWidgetToQuery } from '@/helpers/entities/alarm/query';

describe('convertAlarmUserPreferenceToQuery', () => {
  it('should convert user preferences to a query when all fields are provided', () => {
    const userPreference = {
      content: {
        itemsPerPage: 20,
        category: 'category1',
        mainFilter: 'filter1',
        isCorrelationEnabled: true,
        onlyBookmarks: true,
      },
    };

    const result = convertAlarmUserPreferenceToQuery(userPreference);

    expect(result).toEqual({
      category: 'category1',
      filter: 'filter1',
      only_bookmarks: true,
      correlation: true,
      itemsPerPage: 20,
    });
  });

  it('should handle empty content object gracefully', () => {
    const userPreference = { content: {} };

    const result = convertAlarmUserPreferenceToQuery(userPreference);

    expect(result).toEqual({
      only_bookmarks: false,
    });
  });
});

describe('convertAlarmWidgetToQuery', () => {
  it('should convert widget with \'AlarmsList\' type to a query object with default parameters and isCorrelationEnabled', () => {
    const widget = {
      parameters: {
        liveReporting: {},
        itemsPerPage: 10,
        opened: true,
        sort: null,
        mainFilter: 'testFilter',
        usedAlarmProperties: [],
        isCorrelationEnabled: true,
      },
    };

    const result = convertAlarmWidgetToQuery(widget);

    expect(result).toEqual({
      opened: widget.parameters.opened,
      itemsPerPage: widget.parameters.itemsPerPage,
      page: 1,
      with_instructions: true,
      with_declare_tickets: true,
      with_links: true,
      sortBy: [],
      sortDesc: [],
      lockedFilter: widget.parameters.mainFilter,
      correlation: true,
    });
  });

  it('should handle empty widget.parameters gracefully', () => {
    const widget = {
      parameters: {},
    };

    const result = convertAlarmWidgetToQuery(widget);
    expect(result).toEqual({
      opened: ALARMS_OPENED_VALUES.opened,
      itemsPerPage: PAGINATION_LIMIT,
      page: 1,
      with_instructions: true,
      with_declare_tickets: true,
      with_links: true,
      sortBy: [],
      sortDesc: [],
      lockedFilter: undefined,
    });
  });
});
