import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  viewsAndPages: 'Views and pages',
  allViews: 'All views',
  allPlaylists: 'All playlists',
  closable: 'Can be closed by user',
  closableHelp: 'If it\'s on, user can close the message by clicking Close button\n\n'
    + 'If message was closed, it won\'t appear anymore for this user',
  priorityHelp: 'If there will be more than 1 active broadcast message, they will be sorted according to the selected priority (1st on top)\n\n'
    + 'If the messages have the same priority, they will be sorted according to the start date (newest on top)',
  errors: {
    viewsRequired: 'At least one page must be selected',
  },
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Active',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'Pending',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expired',
  },
};
