import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
  EXTERNAL_DATA_TABLES_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_TAGS,
  EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

export default {
  title: 'Données externes',
  add: 'Ajouter des données externes',
  empty: 'Aucune donnée externe n\'a encore été ajoutée',
  fields: {
    reference: 'Référence',
    collection: 'Collection',
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

  updatePreview: 'Mettre à jour l\'aperçu',
  loadingPreview: 'Chargement de l\'aperçu',
  selectDataType: 'Sélectionner le type de données',
  fieldsHasError: '{count} champ contient des erreurs|{count} champs contiennent des erreurs',

  andMore: 'et plus...',
  linkedRules: {
    widgets: '<strong>Widgets</strong> qui utilisent ce tableau<br><ul>{rules}</ul>',
    eventFilters: '<strong>Filtres d\'événements</strong>\n<ul>{rules}</ul>',
    links: '<strong>Links</strong>\n<ul>{rules}</ul>',
  },
  tableCanBeDeletedInConfig: 'La table ne peut être supprimée que dans le fichier de configuration',
  tableCanBeDeletedAfter: 'La table pourra être supprimée après la suppression de\n{rules}',
  tableRemovedFromConfig: 'La table est supprimée du fichier de configuration, mais elle reste utilisée dans les éléments suivants.\n<strong>Replacer la table dans le fichier de configuration ou supprimer tous les éléments qui l\'utilisent</strong>\n{rules}',
  tableEmptyColumns: 'Veuillez choisir au moins 1 colonne dans les paramètres',
  tableColumnTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.noTag]: 'Pas d\'étiquette',
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.filter]: 'Filtre',
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.context]: 'Contexte',
  },
  tableColumnDataTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string]: {
      text: '@:common.variableTypes.string',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number]: {
      text: '@:common.variableTypes.number',
      tooltip: 'Les valeurs dans le champ doivent contenir uniquement :<br>• nombres<br>• virgules<br>• points<br>• signes moins<br>• espace',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean]: {
      text: '@:common.variableTypes.boolean',
      tooltip: 'Les valeurs dans le champ doivent contenir l\'une des valeurs suivantes (insensible à la casse) :<br>• yes / no<br>• y / n<br>• oui / non<br>• true / false<br>• 1 / 0<br><br>Toutes ces valeurs seront converties en true / false',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp]: {
      text: '@:common.timestamp',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime]: {
      text: 'Date/heure',
      tooltip: 'Les valeurs dans le champ doivent être dans l\'un des formats suivants :<br>• 1990-12-31T00:00:00.000Z<br>• 1990-12-31T00:00:00:00:00<br>• 1990-12-31T00:00:00Z<br>• 1990-12-31T00:00:00+00:00',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray]: {
      text: '@:common.variableTypes.array',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp]: {
      text: '@:common.regexp',
      tooltip: 'Une nouvelle colonne Priorité sera ajoutée après ces étapes :<br>1. À toutes les valeurs ^ est ajouté au début et $ à la fin (sauf pour le .*)<br>2. Chaque valeur est notée<br><br><table style="border-collapse: collapse; width: 100%;"><thead><tr><th style="border: 1px solid currentColor;" class="pa-2">Type de valeur</th><th style="border: 1px solid currentColor;" class="pa-2">Score</th></tr></thead><tbody><tr><td style="border: 1px solid currentColor;" class="pa-2">Vide OU .* OU .+</td><td style="border: 1px solid currentColor;" class="pa-2">0</td></tr><tr><td style="border: 1px solid currentColor;" class="pa-2">Contient * OU regex général</td><td style="border: 1px solid currentColor;" class="pa-2">1</td></tr><tr><td style="border: 1px solid currentColor;" class="pa-2">Avec un seul ^ OU $</td><td style="border: 1px solid currentColor;" class="pa-2">2</td></tr><tr><td style="border: 1px solid currentColor;" class="pa-2">Valeur exacte (commence par ^ et finit par $)</td><td style="border: 1px solid currentColor;" class="pa-2">3</td></tr></tbody></table><br>3. La colonne Priorité est ajoutée avec la priorité résumée (toutes les priorités des colonnes regexp sont additionnées)',
    },
  },
  tableColumnDataTypesAdditionalChips: {
    number: {
      selectDecimalSeparator: 'Sélectionner le séparateur décimal',
      selectThousandsSeparator: 'Sélectionner le séparateur de milliers',
      decimalSeparator: 'Séparateur décimal',
      thousandsSeparator: 'Séparateur de milliers',
      decimalSeparatorDisabled: 'Est déjà utilisé comme séparateur de milliers',
      thousandsSeparatorDisabled: 'Est déjà utilisé comme séparateur décimal',
      separatorDisabledByTableSeparator: 'Est déjà utilisé comme séparateur de table',
    },
    stringArray: {
      separator: 'Séparateur',
      selectSeparator: 'Sélectionner le séparateur',
      types: {
        [EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.json]: {
          text: 'JSON',
          description: '[v1, v2, v3]',
        },
        [EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom]: {
          text: 'Analyser avec un séparateur',
          description: 'v1,v2,v3',
        },
      },
    },
    forbiddenSeparator: 'Ce séparateur ne peut pas être utilisé car il entre en conflit avec le séparateur de table',
  },
};
