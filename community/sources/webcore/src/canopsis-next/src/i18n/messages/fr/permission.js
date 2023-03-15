import { USER_PERMISSIONS_PREFIXES } from '@/constants';

export default {
  technical: {
    admin: 'Droits d\'administration',
    exploitation: 'Droits d\'exploitation',
    notification: 'Droits de notification',
    profile: 'Droits de profil',
  },
  business: {
    [USER_PERMISSIONS_PREFIXES.business.common]: 'Droits communs',
    [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Droits pour le widget : Bac à alarmes',
    [USER_PERMISSIONS_PREFIXES.business.context]: 'Droits pour le widget : Explorateur de contexte',
    [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Droits pour le widget: Météo des services',
    [USER_PERMISSIONS_PREFIXES.business.counter]: 'Droits pour le widget : Compteur',
    [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Droits pour le widget : Scénario des tests',
    [USER_PERMISSIONS_PREFIXES.business.map]: 'Droits pour le widget : Cartographie',
    [USER_PERMISSIONS_PREFIXES.business.barChart]: 'Droits pour le widget : Diagramme à bandes',
    [USER_PERMISSIONS_PREFIXES.business.lineChart]: 'Droits pour le widget : Graphique en ligne',
    [USER_PERMISSIONS_PREFIXES.business.pieChart]: 'Droits pour le widget : Diagramme circulaire',
    [USER_PERMISSIONS_PREFIXES.business.numbers]: 'Droits pour le widget : Nombres',
  },
  api: {
    general: 'Général',
    rules: 'Règles',
    remediation: 'Remédiation',
    pbehavior: 'PBehavior',
  },
};
