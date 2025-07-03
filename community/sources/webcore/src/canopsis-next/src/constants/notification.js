export const NOTIFICATION_TYPES = {
  instructionRate: 0,
  eventFilterFailure: 1,
  instructionApprove: 2,
  instructionDismiss: 3,
};

export const DEFAULT_NOTIFICATION_TOP_BAR_LIMIT = 3;

export const NOTIFICATIONS_PAGE_TABS_KEYS = {
  instructionsToApprove: 'instructions-to-approve',
  instructionsToRate: 'instructions-to-rate',
  eventFilterFailures: 'event-filter-failures',
};

export const NOTIFICATIONS_PAGE_TABS_KEYS_BY_TYPE = {
  [NOTIFICATION_TYPES.instructionRate]: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate,
  [NOTIFICATION_TYPES.eventFilterFailure]: NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures,
  [NOTIFICATION_TYPES.instructionApprove]: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove,
  [NOTIFICATION_TYPES.instructionDismiss]: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove,
};
