import { ALARM_ADVANCED_SEARCH_GROUPS } from '@/constants';

export default {
  title: 'Advanced search',
  not: 'NOT',
  valueTypeListEmptyMessage: 'Press at least 1 symbol and press <kbd>enter</kbd> to apply',
  switchAdvancedSearchActiveToTrue: 'Switch to the advanced search',
  switchAdvancedSearchActiveToFalse: 'Switch to the simple search',
  noDataList: 'There aren\'t any items for input',
  inputPlaceholder: 'Search or filter results...',
  definedDifferent: 'The value differs from the defined one',

  groups: {
    [ALARM_ADVANCED_SEARCH_GROUPS.basic]: 'Basic',
    [ALARM_ADVANCED_SEARCH_GROUPS.messages]: 'Messages',
    [ALARM_ADVANCED_SEARCH_GROUPS.ticket]: 'Ticket',
    [ALARM_ADVANCED_SEARCH_GROUPS.dates]: 'Dates',
    [ALARM_ADVANCED_SEARCH_GROUPS.actions]: 'Actions',
    [ALARM_ADVANCED_SEARCH_GROUPS.entity]: 'Entity',
    [ALARM_ADVANCED_SEARCH_GROUPS.alias]: 'Aliases',
    [ALARM_ADVANCED_SEARCH_GROUPS.pbehavior]: 'Pbehavior',
  },

  searchForThisText: 'Press <kbd>enter</kbd> to search for this text',
  searchByText: 'Search by text',
  listDisabledMessage: 'Not possible to combine patterns\n(alarm, entity and pbehavior) with OR (only AND)',
};
