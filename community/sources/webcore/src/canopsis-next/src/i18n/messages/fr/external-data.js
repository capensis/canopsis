import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
  EXTERNAL_DATA_TABLES_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_TYPES,
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
    [EXTERNAL_DATA_TYPES.api]: 'API',
    [EXTERNAL_DATA_TYPES.table]: 'Table',
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
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.output]: 'Output',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.extraInfos]: 'Informations supplémentaires',
  },

  tableTypes: {
    [EXTERNAL_DATA_TABLES_TYPES.mongo]: 'MongoDB',
    [EXTERNAL_DATA_TABLES_TYPES.postgres]: 'PostgreSQL',
  },

  tableNameTooltip: 'Symboles pris en charge : lettres latines, "_", chiffres (pas au début)',

  importFileDescription: 'La première ligne doit contenir les noms de colonnes',
  exportTableStructure: 'Exporter la structure de la table',

  tableField: 'Collection / Table',

  tableCanBeDeletedInConfig: 'La table ne peut être supprimée que dans le fichier de configuration',
  tableCanBeDeletedAfter: 'La table pourra être supprimée après la suppression de\n{rules}',
  tableRemovedFromConfig: 'La table est supprimée du fichier de configuration, mais elle reste utilisée dans les éléments suivants.\n<strong>Replacer la table dans le fichier de configuration ou supprimer tous les éléments qui l\'utilisent</strong>\n{rules}',
  tableEmptyColumns: 'Veuillez choisir au moins 1 colonne dans les paramètres',
  tableColumnTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.noType]: 'Aucun type',
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.filter]: 'Filtre',
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.context]: 'Contexte',
  },
};
