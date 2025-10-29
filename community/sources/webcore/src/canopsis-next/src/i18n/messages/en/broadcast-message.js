import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  viewsAndPages: 'Views and pages',
  allViews: 'All views',
  allPlaylists: 'All playlists',
  errors: {
    viewsRequired: 'At least one page must be selected',
  },
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Active',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'Pending',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expired',
  },
};
