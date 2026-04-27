import {
  JOB_STATUS,
  JOB_LAST_RUN_STATUS,
  JOB_RULE_TYPE,
  JOB_RULE_TYPES,
  JOBS_TABS,
} from '@/constants';

export default {
  filterByStatus: 'Filtrer par statut',
  filterByLastStatus: 'Filtrer par dernier statut',
  filterByActiveState: 'Filtrer par état actif',
  ruleName: 'Nom de la règle',
  ruleNameTooltip: 'Nom de la règle au moment de la création du job',
  authTokenName: 'Nom du jeton d\'auth',
  ticketSystemName: 'Nom du système de tickets',
  ticketNumber: 'Numéro de ticket',
  ruleType: 'Type de règle',
  activeState: 'État actif',
  lastStatus: 'Dernier statut',
  startDate: 'Date de début',
  finishDate: 'Date de fin',
  failReason: 'Raison de l\'échec',
  expirationDate: 'Date d\'expiration',
  searchByRuleName: 'Rechercher par nom de règle, nom du système de tickets ou numéro de ticket',
  data: {
    request: 'Requête',
    response: 'Réponse',
    webhookFailedPrefix: 'Webhook a échoué',
  },
  status: {
    [JOB_STATUS.running]: 'En cours',
    [JOB_STATUS.paused]: 'En pause',
    [JOB_STATUS.stopped]: 'Arrêté',
  },
  lastRunStatus: {
    [JOB_LAST_RUN_STATUS.succeed]: 'Réussi',
    [JOB_LAST_RUN_STATUS.failed]: 'Échoué',
  },
  ruleTypeValues: {
    [JOB_RULE_TYPE.ticketDeclarationRule]: 'Règle de déclaration de ticket',
    [JOB_RULE_TYPE.scenario]: 'Scénario',
  },
  actions: {
    editJob: 'Modifier le job',
    repeatJob: 'Répéter le job',
    pauseJob: 'Mettre le job en pause',
    startJob: 'Démarrer le job',
  },
  types: {
    [JOB_RULE_TYPES.scenario]: 'Scénario',
    [JOB_RULE_TYPES.declareTicket]: 'Règles de déclaration de ticket',
  },
  tabs: {
    [JOBS_TABS.ticketStatus]: 'Statut des tickets',
  },
  popups: {
    updated: 'Job mis à jour',
    repeated: 'Job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> répété | {count} jobs répétés',
    restarted: 'Job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> redémarré | {count} jobs redémarrés',
    repeatFailed: 'Échec de répétition du job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> | {count} jobs n\'ont pas pu être répétés',
    restartFailed: 'Échec de redémarrage du job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> | {count} jobs n\'ont pas pu être redémarrés',
    paused: 'Job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> mis en pause | {count} jobs mis en pause',
    pauseFailed: 'Échec de mise en pause du job pour <strong>{ruleName}</strong> / <strong>Numéro de ticket {ticketNumber}</strong> | {count} jobs n\'ont pas pu être mis en pause',
  },
};
