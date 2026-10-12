import { VIEW_SCREEN_MODES } from '@/constants';

export default {
  sharedViewUrl: 'Shared view url',
  shareView: 'Share view {name}',
  deleteRow: 'Delete row',
  deleteWidget: 'Delete widget',
  fullScreen: 'Full screen',
  copyWidgetId: 'Copy widget ID',
  autoHeightButton: 'If this button is selected, height will be automatically calculated.',
  selectAll: 'Select all groups and views',
  duplicateAsPrivate: 'Duplicate as private view',
  duplicateAsRegular: 'Duplicate as regular view',
  periodicRefresh: 'Periodic refresh',
  groupIds: 'Choose a group, or create a new one',
  groupTags: 'Group tags',
  noGroupsFound: 'No group corresponding. Press <kbd>enter</kbd> to create a new one',
  errors: {
    emptyTabs: 'You should create a tab',
  },
  screenMode: {
    [VIEW_SCREEN_MODES.default]: {
      title: 'Default mode',
      tooltip: 'Alt / Cmd + Shift + 1',
    },
    [VIEW_SCREEN_MODES.fullscreen]: {
      title: 'Fullscreen',
      tooltip: 'Alt / Cmd + Shift + 2',
    },
    [VIEW_SCREEN_MODES.kiosk]: {
      title: 'Kiosk only',
      tooltip: 'Alt / Cmd + Shift + 3',
    },
    [VIEW_SCREEN_MODES.kioskFullscreen]: {
      title: 'Kiosk + fullscreen',
      tooltip: 'Alt / Cmd + Shift + 4',
    },
  },
};
