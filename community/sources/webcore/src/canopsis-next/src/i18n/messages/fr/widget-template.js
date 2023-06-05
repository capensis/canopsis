import { WIDGET_TEMPLATES_TYPES } from '@/constants';

export default {
  types: {
    [WIDGET_TEMPLATES_TYPES.alarmColumns]: 'Général : Colonnes des alarmes',
    [WIDGET_TEMPLATES_TYPES.entityColumns]: 'Général : Colonnes des entités',
    [WIDGET_TEMPLATES_TYPES.alarmMoreInfos]: 'Bac à alarmes : plus d\'infos',
    [WIDGET_TEMPLATES_TYPES.weatherItem]: 'Météo des services : Modèle de tuile',
    [WIDGET_TEMPLATES_TYPES.weatherModal]: 'Météo des services : Modèle de modale',
    [WIDGET_TEMPLATES_TYPES.weatherEntity]: 'Météo des services : Modèle d\'entité',
  },
  errors: {
    columnsRequired: 'Vous devez ajouter au moins une colonne.',
  },
};
