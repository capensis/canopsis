import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Active',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'Pending',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expired',
  },
};
