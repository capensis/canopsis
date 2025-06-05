import { NOTIFICATION_TYPES } from '@/constants';

export default {
  tabs: {
    instructionsToApprove: 'Instructions approvals',
    instructionsToRate: 'Instructions rating',
    eventFilterFailures: 'Event filter errors',
  },
  headers: {
    name: 'Name',
    description: 'Description',
    created: 'Created',
    lastExecutedOn: 'Last Executed On',
    author: 'Author',
    actions: 'Actions',
  },
  actions: {
    rate: 'Rate',
    approve: 'Approve',
    dismiss: 'Dismiss',
    refresh: 'Refresh',
  },
  topBar: {
    types: {
      [NOTIFICATION_TYPES.instructionRate]: 'Instruction rate',
      [NOTIFICATION_TYPES.eventFilterFailure]: 'Event filter failure',
      [NOTIFICATION_TYPES.instructionDismiss]: 'Instruction dismiss',
      [NOTIFICATION_TYPES.instructionApprove]: 'Instruction approval (request)',
    },
  },
};
