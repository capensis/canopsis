import { PATTERN_TYPES } from '@/constants';

export default {
  patterns: 'Patterns',
  myPatterns: 'My patterns',
  corporatePatterns: 'Shared patterns',
  addRule: 'Add rule',
  addGroup: 'Add group',
  removeRule: 'Remove rule',
  advancedEditor: 'Advanced editor',
  simpleEditor: 'Simple editor',
  noData: 'No pattern set. Click \'@:pattern.addGroup\' button to start adding fields to the pattern',
  noDataDisabled: 'No pattern set.',
  discard: 'Discard pattern',
  oldPatternTooltip: 'Filter patterns are not migrated',
  types: {
    [PATTERN_TYPES.alarm]: 'Alarm pattern',
    [PATTERN_TYPES.entity]: 'Entity pattern',
    [PATTERN_TYPES.pbehavior]: 'Pbehavior pattern',
  },
  errors: {
    ruleRequired: 'Please add at least one rule',
    groupRequired: 'Please add at least one group',
    invalidPatterns: 'Patterns are invalid or there is a disabled pattern field',
    countOverLimit: 'The patterns you\'ve defined targets about {count} items. It can affect performance, are you sure ?',
    oldPattern: 'The current filter pattern is defined in old format. Please use the Advanced editor to view it. Filters in old format will be deprecated soon. Please create new patterns in our updated interface.',
    existExcluded: 'The rules include excluded rule.',
    required: 'At least one pattern has to be defined. Please define filter patterns for rule',
  },
};
