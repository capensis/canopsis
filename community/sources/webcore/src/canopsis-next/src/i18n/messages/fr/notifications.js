import { NOTIFICATION_TYPES } from '@/constants';

export default {
  tabs: {
    eventFilterFailures: 'Échecs de filtres d\'événements',
    instructionsToRate: 'Consignes à évaluer',
    instructionsToApprove: 'Consignes à approuver',
  },
  headers: {
    name: 'Nom',
    description: 'Description',
    created: 'Créé',
    lastExecutedOn: 'Dernière exécution',
    author: 'Auteur',
    actions: 'Actions',
  },
  actions: {
    rate: 'Évaluer',
    approve: 'Approuver',
    dismiss: 'Rejeter',
    refresh: 'Actualiser',
  },
  topBar: {
    types: {
      [NOTIFICATION_TYPES.instructionRate]: 'Évaluation de consigne',
      [NOTIFICATION_TYPES.eventFilterFailure]: 'Erreur de filtre d\'événement',
      [NOTIFICATION_TYPES.instructionApprove]: 'Approbation de consigne (demande)',
      [NOTIFICATION_TYPES.instructionDismiss]: 'Rejet de consigne',
    },
    seeAll: 'Voir toutes les notifications',
    noNotifications: 'Aucune notification',
  },
};
