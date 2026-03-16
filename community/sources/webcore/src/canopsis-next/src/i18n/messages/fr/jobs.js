import { JOB_STATUS, JOB_LAST_RUN_STATUS, JOB_RULE_TYPE } from '@/constants';

export default {
  filterByStatus: 'Filtrer par statut',
  filterByActiveState: 'Filtrer par état actif',
  ruleName: 'Nom de la règle',
  authTokenName: 'Nom du jeton d\'auth',
  ticketSystemName: 'Nom du système de tickets',
  ticketNumber: 'Numéro de ticket',
  ruleType: 'Type de règle',
  active: 'Actif',
  statusLabel: 'Statut',
  startDate: 'Date de début',
  finishDate: 'Date de fin',
  failReason: 'Raison de l\'échec',
  expirationDate: 'Date d\'expiration',
  searchByRuleName: 'Rechercher par nom de règle',
  webhookFailedPrefix: 'Webhook a échoué',
  status: {
    [JOB_STATUS.running]: 'En cours',
    [JOB_STATUS.paused]: 'En pause',
    [JOB_STATUS.stopped]: 'Arrêté',
    unknown: 'Inconnu',
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
    stopJob: 'Arrêter le job',
  },
  tabs: {
    instructions: 'Consignes',
    webhooks: 'Webhooks',
    ticketStatus: 'Statut des tickets',
    authToken: 'Jeton d\'auth',
  },
};
