import { JOB_STATE, JOB_RUN_STATUS, JOB_RULE_TYPE } from '@/constants/jobs';

export default {
  filterByStatus: 'Filtrer par statut',
  ruleName: 'Nom de la règle',
  authTokenName: 'Nom du jeton d\'auth',
  ticketSystemName: 'Nom du système de tickets',
  ruleType: 'Type de règle',
  active: 'Actif',
  statusLabel: 'Statut',
  startDate: 'Date de début',
  finishDate: 'Date de fin',
  failReason: 'Raison de l\'échec',
  expirationDate: 'Date d\'expiration',
  searchByRuleName: 'Rechercher par nom de règle',
  status: {
    [JOB_STATE.running]: 'En cours',
    [JOB_STATE.paused]: 'En pause',
    [JOB_STATE.stopped]: 'Arrêté',
    unknown: 'Inconnu',
  },
  activeState: {
    [JOB_STATE.running]: 'Actif',
    [JOB_STATE.paused]: 'En pause',
    [JOB_STATE.stopped]: 'Arrêté',
  },
  runStatus: {
    [JOB_RUN_STATUS.succeed]: 'Réussi',
    [JOB_RUN_STATUS.failed]: 'Échoué',
    inProgress: 'En cours',
  },
  ruleTypeValues: {
    [JOB_RULE_TYPE.ticketDeclarationRule]: 'Règle de déclaration de ticket',
    [JOB_RULE_TYPE.scenario]: 'Scénario',
  },
  tabs: {
    instructions: 'Consignes',
    webhooks: 'Webhooks',
    ticketStatus: 'Statut des tickets',
    authToken: 'Jeton d\'auth',
  },
  actions: {
    start: 'Démarrer',
    stop: 'Arrêter',
    resume: 'Reprendre',
    pause: 'Mettre en pause',
    edit: 'Modifier',
  },
};
