import { ENTITY_TYPES } from '@/constants';

export default {
  manageInfos: 'Gérer les informations',
  form: 'Formulaire',
  impact: 'Impacts',
  depends: 'Dépendances',
  addInformation: 'Ajouter une information',
  emptyInfos: 'Aucune information',
  availabilityState: 'État de disponibilité',
  types: {
    [ENTITY_TYPES.component]: 'Composant',
    [ENTITY_TYPES.connector]: 'Connecteur',
    [ENTITY_TYPES.resource]: 'Ressource',
    [ENTITY_TYPES.service]: 'Service',
  },
  fields: {
    categoryName: 'Nom de catégorie',
    koEvents: 'KO événements',
    okEvents: 'OK événements',
    statsKo: 'Stats KO',
    statsOk: 'Stats OK',
    idleSince: 'Inactif depuis',
    componentInfos: 'Informations sur les composants',
    alarmDisplayName: 'Nom d\'affichage de l\'alarme',
    alarmCreationDate: 'Date de création de l\'alarme',
  },
};
