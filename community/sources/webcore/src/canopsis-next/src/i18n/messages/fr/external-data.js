import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
} from '@/constants';

export default {
  title: 'Données externes',
  add: 'Ajouter des données externes',
  empty: 'Aucune donnée externe n\'a encore été ajoutée',
  fields: {
    reference: 'Référence',
    collection: 'Collection',
    sort: 'Trier',
    sortBy: 'Trier par',
  },
  tooltips: {
    reference: 'Sera utilisé dans les actions via <strong>.ExternalData.&lt;Reference&gt;</strong>',
  },
  types: {
    [EXTERNAL_DATA_TYPES.mongo]: 'MongoDB collection',
    [EXTERNAL_DATA_TYPES.api]: 'API',
  },
  conditionTypes: {
    [EXTERNAL_DATA_CONDITION_TYPES.select]: 'Sélectionner',
    [EXTERNAL_DATA_CONDITION_TYPES.regexp]: 'Expression régulière',
  },
  conditionValues: {
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.component]: 'Composant',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connector]: 'Connecteur',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connectorName]: 'Nom du connecteur',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.resource]: 'Ressource',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.output]: 'Sortir',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.extraInfos]: 'Informations supplémentaires',
  },
};
