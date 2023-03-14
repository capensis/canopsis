import { ALARM_METRIC_PARAMETERS, USER_METRIC_PARAMETERS, AGGREGATE_FUNCTIONS } from '@/constants';

export default {
  alarmMetrics: 'Métriques d\'alarme',
  sli: 'SLI',
  metricsNotAvailable: 'TimescaleDB ne fonctionne pas. Les métriques ne sont pas disponibles.',
  noData: 'Pas de données disponibles',
  selectMetric: 'Sélectionnez la métrique à afficher',
  customColor: 'Couleur personnalisée',
  calculationMethod: 'Méthode de calcul',
  aggregateFunctions: {
    [AGGREGATE_FUNCTIONS.sum]: 'Somme',
    [AGGREGATE_FUNCTIONS.avg]: 'Moyenne',
    [AGGREGATE_FUNCTIONS.min]: 'Min',
    [AGGREGATE_FUNCTIONS.max]: 'Max',
  },
  errors: {
    metricsMinLength: 'Au moins {count} statistiques doivent être ajoutées',
  },

  metrics: {
    parameter: 'Paramètre à comparer',
    tooltip: {
      [USER_METRIC_PARAMETERS.totalUserActivity]: '{value} temps total d\'activité',

      [ALARM_METRIC_PARAMETERS.createdAlarms]: '{value} alarmes créées',
      [ALARM_METRIC_PARAMETERS.activeAlarms]: '{value} alarmes actives',
      [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: '{value} alarmes non affichées',
      [ALARM_METRIC_PARAMETERS.instructionAlarms]: '{value} alarmes en cours de correction automatique',
      [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: '{value} alarmes sous PBehavior',
      [ALARM_METRIC_PARAMETERS.correlationAlarms]: '{value} alarmes avec corrélation',
      [ALARM_METRIC_PARAMETERS.ackAlarms]: '{value} alarmes avec ack',
      [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: '{value} alarmes actives avec acks',
      [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: '{value} alarmes avec acquittement annulé',
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: '{value} alarmes actives avec tickets',
      [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: '{value} alarmes actives sans tickets',
      [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '{value}% d\'alarmes avec correction automatique',
      [ALARM_METRIC_PARAMETERS.ratioInstructions]: '{value}% d\'alarmes avec consigne',
      [ALARM_METRIC_PARAMETERS.ratioTickets]: '{value}% d\'alarmes avec tickets créés',
      [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '{value}% d\'alarmes corrigées manuellement',
      [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '{value}% des alarmes non affichées',
      [ALARM_METRIC_PARAMETERS.averageAck]: '{value} accuser les alarmes',
      [ALARM_METRIC_PARAMETERS.averageResolve]: '{value} pour résoudre les alarmes',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: '{value} alarmes avec instructions manuelles',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: '{value} alarmes corrigées manuellement',
    },
  },

  filters: {
    helpInformation: 'Ici, les modèles de filtre pour des tranches de données supplémentaires pour les compteurs et les évaluations peuvent être ajoutés.',
  },

  ratingSettings: {
    helpInformation: 'La liste des paramètres à utiliser pour la notation.',
  },
};
