import { GROUPED_USER_PERMISSIONS_KEYS, USER_PERMISSIONS } from '@/constants';

export default {
  hideMessage: 'J\'ai compris! Cacher',
  learnMore: 'En savoir plus sur {link}',

  /**
   * Exploitation
   */
  [USER_PERMISSIONS.technical.exploitation.eventFilter]: {
    title: 'Filtres d\'événements',
    message: 'Centre de filtrage et d\'enrichissement des événements et entités.',
  },

  [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: {
    title: 'Informations dynamiques',
    message: 'Centre d\'enrichissement des alarmes',
  },

  [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
    title: 'Règles de méta-alarme',
    topbarTitle: 'Méta-alarmes',
    message: 'Module de corrélation permettant de regrouper des alarmes en méta-alarmes selon des critères (relation parent/enfant, intervalle de temps, etc).',
  },

  [USER_PERMISSIONS.technical.exploitation.idleRules]: {
    title: 'Règles d\'inactivité',
    topbarTitle: 'Inactivité',
    message: 'Règles permettant de réagir en cas de "non événement". Utile pour gérer les lignes de vie des sources de données par exemple.',
  },

  [USER_PERMISSIONS.technical.exploitation.flappingRules]: {
    title: 'Règles de bagot',
    topbarTitle: 'Bagot',
    message: 'Une alarme hérite du statut bagot lorsqu\'elle oscille d\'une criticité d\'alerte à un état stable un certain nombre de fois sur une période donnée.',
  },

  [USER_PERMISSIONS.technical.exploitation.resolveRules]: {
    title: 'Règles de résolution',
    topbarTitle: 'Résolution',
    message: 'Lorsqu\'une alarme reçoit un événement de type contre alarme, elle passe dans le statut fermée.\nAvant de considérer cette alarme comme définitivement résolue, Canopsis peut attendre un délai.\nCe délai peut être utile si l\'alarme bagotte ou si l\'utilisateur souhaite conserver l\'alarme ouverte en cas d\'erreur.',
  },

  [USER_PERMISSIONS.technical.exploitation.pbehavior]: {
    title: 'Comportements périodiques',
    message: 'Les comportements périodiques de Canopsis peuvent être utilisés afin de définir une période pendant laquelle un comportement doit être modifié. Utile pour définir des périodes de maintenance, des périodes de service.',
  },

  [USER_PERMISSIONS.technical.exploitation.scenario]: {
    title: 'Scénarios',
    message: 'Les scénarios Canopsis peuvent être utilisés pour déclencher conditionnellement divers types d\'actions sur les alarmes.\nExemples : Notifications sur un chat, création d\'un ticket d\'incident, inscription d\'une alarme sur une main courante.',
  },

  [USER_PERMISSIONS.technical.exploitation.snmpRule]: {
    title: 'Règles SNMP',
    topbarTitle: 'SNMP',
    message: 'Les règles SNMP peuvent être utilisées pour traiter les traps SNMP remontées par le connecteur snmp2canopsis au travers du moteur SNMP.',
  },

  [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: {
    title: 'Règles de déclaration des tickets',
    topbarTitle: 'Déclaration de tickets',
    message: 'Permet de définir les règles de création de tickets d\'incidents à destination de plusieurs outils de gestion d\'incidents. Ces tickets sont créés par un pilote depuis la WebUI, et Canopsis se charge ensuite de contacter l\'outil de gestion d\'incident pour l\'enregistrement du ticket.',

  },

  [USER_PERMISSIONS.technical.exploitation.linkRule]: {
    title: 'Générateur de liens',
    message: 'Permet de définir les règles d\'association de liens aux alarmes',
  },

  [GROUPED_USER_PERMISSIONS_KEYS.remediation]: {
    title: 'Consignes',
    message: 'Permet de créer des consignes de remédiation, avec ou sans job, pour corriger des situations.',
  },

  /**
   * Administration - Access
   */
  [USER_PERMISSIONS.technical.permission]: {
    title: 'Droits',
    message: 'Les droits Widgets, Vues, et Technique s\'appliquent uniquement à l\'interface utilisateur. Les droits API s\'appliquent uniquement aux utilisateurs API.',
  },
  [USER_PERMISSIONS.technical.role]: {
    title: 'Rôles',
    message: 'Gestion des rôles des utilisateurs de Canopsis',
  },
  [USER_PERMISSIONS.technical.user]: {
    title: 'Utilisateurs',
    message: 'Gestion des utilisateurs de Canopsis',
  },

  /**
   * Administration - Communications
   */
  [USER_PERMISSIONS.technical.broadcastMessage]: {
    title: 'Diffusion de messages',
    message: 'Permet d\'afficher des bannières et des messages d\'information dans l\'interface de Canopsis.',
  },
  [USER_PERMISSIONS.technical.playlist]: {
    title: 'Listes de lecture',
    message: 'Permet de personnaliser l\'affichage des vues en les faisant défiler les unes après les autres avec un délai défini.',
  },
  [USER_PERMISSIONS.technical.healthcheck]: {
    title: 'Bilan de santé',
    message: 'Le Healthcheck est un tableau de bord indiquant l\'état et les erreurs de tous les composants inclus dans Canopsis.',
  },

  /**
   * Administration - Custom objects
   */
  [USER_PERMISSIONS.technical.commentTemplate]: {
    title: 'Modèles de commentaires',
    topbarTitle: 'Modèles de commentaires',
  },

  /**
   * Administration - Settings
   */
  [USER_PERMISSIONS.technical.parameters]: {
    title: 'Interface utilisateur',
  },
  [USER_PERMISSIONS.technical.viewImportExport]: {
    title: 'Import / export',
  },
  [USER_PERMISSIONS.technical.notification.common]: {
    title: 'Paramètres de notification',
    topbarTitle: 'Notifications',
  },
  [USER_PERMISSIONS.technical.engine]: {
    title: 'Engines',
    message: 'Cette page contient les informations sur la séquence et la configuration des moteurs. Pour fonctionner correctement, la chaîne des moteurs doit être continue.',
  },
  [USER_PERMISSIONS.technical.kpi]: {
    title: 'KPI',
    message: 'Permet de présenter des indicateurs sous forme de graphiques',
  },
  [USER_PERMISSIONS.technical.map]: {
    title: 'Cartographie',
    message: 'Module permettant de définir et d\'afficher des cartes (géographiques, logiques, mermaid, etc.) via le widget "Map".',
  },
  [USER_PERMISSIONS.technical.maintenance]: {
    title: 'Mode de maintenance',
    message: 'Permet de basculer Canopsis en mode maintenance. Tous les utilisateurs, à l\'exception des administrateurs, seront déconnectés.',
  },
  [USER_PERMISSIONS.technical.tag]: {
    title: 'Gestion des tags',
    message: 'Permet de définir des règles d\'attribution de tags aux alarmes',
  },
  [USER_PERMISSIONS.technical.storageSettings]: {
    title: 'Paramètres de stockage',
    topbarTitle: 'Stockage',
    message: 'Permet de définir les politiques de rétention des données.',
  },
  [USER_PERMISSIONS.technical.stateSetting]: {
    title: 'Paramètres de calcul d\'état/sévérité',
    topbarTitle: 'État',
    message: 'Permet de définir des méthodes de calcul d\'état/sévérité de composants et/ou de services.\nUtile pour modéliser des arbres de dépendances.',
  },
  [USER_PERMISSIONS.technical.eventsRecord]: {
    title: 'Enregistrements d\'événements',
    message: 'Permet de définir et de déclencher l\'enregistrement des événements dès leur arrivée dans le bus de données Canopsis.',
  },
  [USER_PERMISSIONS.technical.templateTesting]: {
    title: 'Studio templates',
  },
  [USER_PERMISSIONS.technical.externalAuthTokens]: {
    title: 'Jetons d\'authentification externes',
  },
  [USER_PERMISSIONS.technical.widgetTemplate]: {
    title: 'Modèles de widget',
  },
  [USER_PERMISSIONS.technical.externalDataTable]: {
    title: 'Données externes',
  },
  [USER_PERMISSIONS.technical.entityInfoProperty]: {
    title: 'Informations d\'entité',
  },

  /**
   * Administration - Remediation
   */
  [GROUPED_USER_PERMISSIONS_KEYS.planning]: {
    title: 'Planification',
    message: 'Permet d\'administrer la planification dans Canopsis et de personnaliser les types de comportements périodiques.',
  },

  /**
   * Profile
   */
  [USER_PERMISSIONS.technical.profile.theme]: {
    title: 'Thèmes graphiques',
    message: 'Permet de définir des thèmes graphiques pour l\'interface graphique.',
  },
};
