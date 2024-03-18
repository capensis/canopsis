import { GROUPS_NAVIGATION_TYPES, USER_METRIC_PARAMETERS } from '@/constants';

export default {
  seeProfile: 'Voir le profil',
  selectDefaultView: 'Sélectionner une vue par défaut',
  displayName: 'Nom d\'affichage de l\'utilisateur',
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
    [USER_METRIC_PARAMETERS.tickets]: 'Nombre de billets créés',
    [USER_METRIC_PARAMETERS.ackAlarmWithoutCancel]: 'Nombre d\'accusés de réception excluant les accusés de réception annulés',
    [USER_METRIC_PARAMETERS.averageUserSession]: 'Temps moyen d\'activité de l\'utilisateur',
    [USER_METRIC_PARAMETERS.minUserSession]: 'Durée d\'activité minimale de l\'utilisateur',
    [USER_METRIC_PARAMETERS.maxUserSession]: 'Durée d\'activité maximale de l\'utilisateur',
  },
  variables: {
    userEmail: 'Courriel de l\'utilisateur',
    userUsername: 'Nom d\'utilisateur',
    userFirstname: 'Prénom de l\'utilisateur',
    userLastname: 'Nom de famille de l\'utilisateur',
    userExternalId: 'ID externe de l\'utilisateur',
    userSource: 'Source d\'utilisateurs',
    userRole: 'Rôle d\'utilisateur',
  },
};
