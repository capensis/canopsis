import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  viewsAndPages: 'Vues et pages',
  allViews: 'Toutes les vues',
  allPlaylists: 'Toutes les listes de lecture',
  errors: {
    viewsRequired: 'Au moins une page doit être sélectionnée',
  },
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Actif',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'En attente',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expiré',
  },
};
