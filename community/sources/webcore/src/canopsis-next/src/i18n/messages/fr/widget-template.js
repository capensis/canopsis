import { WIDGET_TEMPLATES_TYPES } from '@/constants';

export default {
  types: {
    [WIDGET_TEMPLATES_TYPES.alarmColumns]: 'Général : Colonnes des alarmes',
    [WIDGET_TEMPLATES_TYPES.entityColumns]: 'Général : Colonnes des entités',
    [WIDGET_TEMPLATES_TYPES.alarmMoreInfos]: 'Bac à alarmes : plus d\'infos',
    [WIDGET_TEMPLATES_TYPES.alarmExportToPdf]: 'Bac à alarmes : Modèle d\'export PDF',
    [WIDGET_TEMPLATES_TYPES.alarmQuickActions]: 'Bac à alarmes : Actions rapides (alarme unitaire)',
    [WIDGET_TEMPLATES_TYPES.alarmMassQuickActions]: 'Bac à alarmes : Actions rapides (massive)',
    [WIDGET_TEMPLATES_TYPES.alarmSortColumns]: 'Général : Colonnes de tri par défaut des alarmes',
    [WIDGET_TEMPLATES_TYPES.weatherItem]: 'Météo des services : Modèle de tuile',
    [WIDGET_TEMPLATES_TYPES.weatherModal]: 'Météo des services : Modèle de modale',
    [WIDGET_TEMPLATES_TYPES.weatherEntity]: 'Météo des services : Modèle d\'entité',
  },
  errors: {
    columnsRequired: 'Vous devez ajouter au moins une colonne.',
    quickActionsRequired: 'Vous devez ajouter au moins une action rapide.',
    sortColumnsRequired: 'Vous devez ajouter au moins une colonne de tri.',
  },
};
