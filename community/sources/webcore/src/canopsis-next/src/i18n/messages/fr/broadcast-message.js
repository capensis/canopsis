import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Actif',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'En attente',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expiré',
  },
};
