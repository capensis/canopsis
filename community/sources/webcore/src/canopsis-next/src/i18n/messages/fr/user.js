import { GROUPS_NAVIGATION_TYPES, USER_METRIC_PARAMETERS } from '@/constants';

export default {
  seeProfile: 'Voir le profil',
  selectDefaultView: 'Sélectionner une vue par défaut',
  firstName: 'Prénom',
  lastName: 'Nom',
  email: 'Email',
  language: 'Langue par défaut',
  auth: 'Type d\'auth.',
  navigationType: 'Type d\'affichage de la barre de vues',
  active: 'Session active',
  activeConnects: 'Nombre de connexions',
  navigationTypes: {
    [GROUPS_NAVIGATION_TYPES.sideBar]: 'Barre latérale',
    [GROUPS_NAVIGATION_TYPES.topBar]: 'Barre d\'entête',
  },
  metrics: {
    [USER_METRIC_PARAMETERS.totalUserActivity]: 'Durée totale de l\'activité',
  },
};
