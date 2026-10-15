import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

export default {
  viewsAndPages: 'Vues et pages',
  allViews: 'Toutes les vues',
  allPlaylists: 'Toutes les listes de lecture',
  closable: 'Peut être fermé par l\'utilisateur',
  closableHelp: 'Si l\'option est activée, l\'utilisateur peut fermer le message en cliquant sur le bouton Fermer\n\n'
    + 'Si le message a été fermé, il ne s\'affichera plus pour cet utilisateur',
  priorityHelp: 'S\'il y a plusieurs messages de diffusion actifs, ils seront triés selon la priorité sélectionnée (la priorité 1 en haut)\n\n'
    + 'Si les messages ont la même priorité, ils seront triés selon la date de début (le plus récent en haut)',
  errors: {
    viewsRequired: 'Au moins une page doit être sélectionnée',
  },
  statuses: {
    [BROADCAST_MESSAGES_STATUSES.active]: 'Actif',
    [BROADCAST_MESSAGES_STATUSES.pending]: 'En attente',
    [BROADCAST_MESSAGES_STATUSES.expired]: 'Expiré',
  },
};
