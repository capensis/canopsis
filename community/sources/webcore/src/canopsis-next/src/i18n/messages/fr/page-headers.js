import { USERS_PERMISSIONS } from '@/constants';

export default {
  hideMessage: 'J\'ai compris! Cacher',
  learnMore: 'En savoir plus sur {link}',

  /**
   * Exploitation
   */
  [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
    title: 'Filtres d\'événements',
    message: 'Le filtre d\'événements est une fonctionnalité du moteur CHE qui permet la définition de règles de traitement des événements.',
  },

  [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
    title: 'Informations dynamiques',
    message: 'Les informations dynamiques de Canopsis sont utilisées pour enrichire les alarmes. Ces enrichissements sont définis par des règles indiquant dans quelles conditions les informations doivent être présentées sur une alarme.',
  },

  [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
    title: 'Règles de méta-alarme',
    message: 'Les règles de méta-alarme peuvent être utilisées pour grouper les alarmes par types et critères (relation parent/enfant, intervalle de temps, etc).',
  },

  [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
    title: 'Règles d\'inactivité',
    message: 'Les règles d\'inactivité des entités ou des alarmes peuvent être utilisées pour surveiller les événements et les états d\'alarme afin de savoir si des événements ne sont pas reçus ou si l\'état d\'alarme n\'est pas modifié pendant une longue période en raison d\'erreurs ou d\'une configuration non valide.',
  },

  [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
    title: 'Règles de bagot',
    message: 'Une alarme hérite du statut bagot lorsqu\'elle oscille d\'une criticité d\'alerte à un état stable un certain nombre de fois sur une période donnée.',
  },

  [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
    title: 'Règles de résolution',
    message: 'Lorsqu\'une alarme reçoit un événement de type contre alarme, elle passe dans le statut fermée.\nAvant de considérer cette alarme comme définitivement résolue, Canopsis peut attendre un délai.\nCe délai peut être utile si l\'alarme bagotte ou si l\'utilisateur souhaite conserver l\'alarme ouverte en cas d\'erreur.',
  },

  [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
    title: 'Comportements périodiques',
    message: 'Les comportements périodiques de Canopsis peuvent être utilisés afin de définir une période pendant laquelle le comportement doit être modifié, par exemple pour la maintenance ou le type de services.',
  },

  [USERS_PERMISSIONS.technical.exploitation.scenario]: {
    title: 'Scénarios',
    message: 'Les scénarios Canopsis peuvent être utilisés pour déclencher conditionnellement divers types d\'actions sur les alarmes.',
  },

  [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
    title: 'Règles SNMP',
    message: 'Les règles SNMP peuvent être utilisées pour traiter les traps SNMP remontées par le connecteur snmp2canopsis au travers du moteur SNMP.',
  },

  [USERS_PERMISSIONS.technical.exploitation.declareTicketRule]: {
    title: 'Règles de déclaration des tickets',
    // message: '', // TODO: need to put description
  },

  [USERS_PERMISSIONS.technical.exploitation.linkRule]: {
    title: 'Générateur de liens',
    // message: '', // TODO: need to put description
  },

  /**
   * Admin access
   */
  [USERS_PERMISSIONS.technical.permission]: {
    title: 'Droits',
  },
  [USERS_PERMISSIONS.technical.role]: {
    title: 'Rôles',
  },
  [USERS_PERMISSIONS.technical.user]: {
    title: 'Utilisateurs',
  },

  /**
   * Admin communications
   */
  [USERS_PERMISSIONS.technical.broadcastMessage]: {
    title: 'Diffusion de messages',
    message: 'La diffusion de messages peut être utilisée pour afficher les bannières et les messages d\'information qui apparaîtront dans l\'interface de Canopsis.',
  },
  [USERS_PERMISSIONS.technical.playlist]: {
    title: 'Listes de lecture',
    message: 'Les listes de lecture peuvent être utilisées pour la personnalisation des vues qui peuvent être affichées les unes après les autres avec un délai associé.',
  },

  /**
   * Admin general
   */
  [USERS_PERMISSIONS.technical.parameters]: {
    title: 'Paramètres',
  },
  [USERS_PERMISSIONS.technical.planning]: {
    title: 'Planification',
    message: 'La fonctionnalité d\'administration de la planification de Canopsis peut être utilisée pour la personnalisation des types de comportements périodiques.',
  },
  [USERS_PERMISSIONS.technical.remediation]: {
    title: 'Consignes',
    message: 'La fonction de remédiation de Canopsis peut être utilisée pour créer des plans ou des consignes visant à corriger des situations.',
  },
  [USERS_PERMISSIONS.technical.healthcheck]: {
    title: 'Bilan de santé',
    message: 'La fonction Healthcheck est le tableau de bord avec des indications d\'états et d\'erreurs de tous les systèmes inclus dans Canopsis.',
  },
  [USERS_PERMISSIONS.technical.engine]: {
    title: 'Engines',
    message: 'Cette page contient les informations sur la séquence et la configuration des moteurs. Pour fonctionner correctement, la chaîne des moteurs doit être continue.',
  },
  [USERS_PERMISSIONS.technical.kpi]: {
    title: 'KPI',
    // message: '', // TODO: add correct message
  },
  [USERS_PERMISSIONS.technical.map]: {
    title: 'Cartographie',
    // message: '', // TODO: add correct message
  },
  [USERS_PERMISSIONS.technical.maintenance]: {
    title: 'Mode de Maintenance',
    // message: '', // TODO: add correct message
  },
  [USERS_PERMISSIONS.technical.tag]: {
    title: 'Gestion des Tags',
    // message: '', // TODO: add correct message
  },
  [USERS_PERMISSIONS.technical.storageSettings]: {
    title: 'Paramètres de stockage',
    // message: '', // TODO: add correct message
  },
  [USERS_PERMISSIONS.technical.stateSetting]: {
    title: 'Paramètres d\'état',
    // message: '', // TODO: add correct message
  },

  /**
   * Notifications
   */
  [USERS_PERMISSIONS.technical.notification.instructionStats]: {
    title: 'Évaluation des consignes',
    message: 'Cette page contient les statistiques sur l\'exécution des consignes. Les utilisateurs peuvent noter les consignes en fonction de leurs performances.',
  },

  /**
   * Profile
   */
  [USERS_PERMISSIONS.technical.profile.theme]: {
    title: 'Thèmes graphiques',
    // message: '', // TODO: add correct message
  },
};
