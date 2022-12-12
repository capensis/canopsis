import { ENTITY_TYPES } from '@/constants';

export default {
  manageInfos: 'Gérer les informations',
  form: 'Formulaire',
  impact: 'Impacts',
  depends: 'Dépendances',
  addInformation: 'Ajouter une information',
  emptyInfos: 'Aucune information',
  availabilityState: 'État de disponibilité',
  okEvents: 'OK événements',
  koEvents: 'KO événements',
  types: {
    [ENTITY_TYPES.component]: 'Composant',
    [ENTITY_TYPES.connector]: 'Connecteur',
    [ENTITY_TYPES.resource]: 'Ressource',
    [ENTITY_TYPES.service]: 'Service',
  },
};
