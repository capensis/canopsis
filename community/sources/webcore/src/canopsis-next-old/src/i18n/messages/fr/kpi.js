import { ALARM_METRIC_PARAMETERS, USER_METRIC_PARAMETERS, AGGREGATE_FUNCTIONS } from '@/constants';

export default {
  alarmMetrics: 'Métriques d\'alarme',
  sli: 'SLI',
  metricsNotAvailable: 'TimescaleDB ne fonctionne pas. Les métriques ne sont pas disponibles.',
  noData: 'Pas de données disponibles',
  selectMetric: 'Sélectionnez la métrique à afficher',
  addMetricMask: 'Ajoutez des métriques par masque, par ex. cpu*',
  displayedLabel: 'Étiquette affichée',
  customColor: 'Couleur personnalisée',
  calculationMethod: 'Méthode de calcul',
  periodTrend: '{count} pour la période\n{from} - {to}',
  largeCountOfMetrics: 'La liste des métriques à afficher est tronquée.',
  onlyDisplayed: 'Seules {count} métriques sont affichées.',
  autoAdd: 'Ajout automatique',
  addExternal: 'Ajouter externe',
  tabs: {
    collectionSettings: 'Paramètres d\'évaluation',
    ratingSettings: 'Paramètres d\'évaluation',
  },

  aggregateFunctions: {
    [AGGREGATE_FUNCTIONS.last]: 'Dernier',
    [AGGREGATE_FUNCTIONS.sum]: 'Somme',
    [AGGREGATE_FUNCTIONS.avg]: 'Moyenne',
    [AGGREGATE_FUNCTIONS.min]: 'Min',
    [AGGREGATE_FUNCTIONS.max]: 'Max',
  },

  errors: {
    metricsMinLength: 'Au moins {count} statistiques doivent être ajoutées',
    emptyMetrics: 'Aucune donnée à afficher',
  },

  popups: {
    exportFailed: 'Échec de l\'exportation du graphique au format CSV',
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
      [ALARM_METRIC_PARAMETERS.timeToAck]: '{value} accuser les alarmes',
      [ALARM_METRIC_PARAMETERS.timeToResolve]: '{value} pour résoudre les alarmes',
      [ALARM_METRIC_PARAMETERS.minAck]: '{value} - temps minimum pour acquitter les alarmes',
      [ALARM_METRIC_PARAMETERS.maxAck]: '{value} - temps maximum pour acquitter les alarmes',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: '{value} alarmes corrigées manuellement',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: '{value} alarmes avec instructions manuelles',
      [ALARM_METRIC_PARAMETERS.notAckedAlarms]: '{value} alarmes non acquittées',
      [ALARM_METRIC_PARAMETERS.notAckedInHourAlarms]: '{value} alarmes non acquittées avec une durée de 1-4h',
      [ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms]: '{value} alarmes non acquittées d\'une durée de 4-24h',
      [ALARM_METRIC_PARAMETERS.notAckedInDayAlarms]: '{value} alarmes non acquittées datant de plus de 24h',
    },
  },

  filters: {
    helpInformation: 'Ici, les modèles de filtre pour des tranches de données supplémentaires pour les compteurs et les évaluations peuvent être ajoutés.',
  },

  ratingSettings: {
    helpInformation: 'La liste des paramètres à utiliser pour la notation.',
  },

  collectionSetting: {
    basicMetrics: 'Métriques de base',
    optionalMetrics: 'Métriques facultatives',
    manualInstructions: 'Nombre d\'alarmes avec instructions manuelles',
    notAckedMetrics: 'Nombre d\'alarmes actives non acquittées de différentes durées',
    sliDuration: 'Durée SLI',
  },

  statisticsWidgets: {
    metrics: {
      [ALARM_METRIC_PARAMETERS.createdAlarms]: 'Nombre moyen d\'alarmes créées',
      [ALARM_METRIC_PARAMETERS.activeAlarms]: 'Nombre moyen d\'alarmes actives',
      [ALARM_METRIC_PARAMETERS.instructionAlarms]: 'Nombre moyen d\'alarmes avec remédiation automatique',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: 'Nombre moyen d\'alarmes actives avec remédiations manuelles',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: 'Nombre moyen d\'alarmes remédiées manuellement',
      [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: 'Nombre moyen d\'alarmes non affichées',
      [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: 'Nombre moyen d\'alarmes avec PBehavior',
      [ALARM_METRIC_PARAMETERS.correlationAlarms]: 'Nombre moyen d\'alarmes avec corrélation',
      [ALARM_METRIC_PARAMETERS.notAckedAlarms]: 'Nombre moyen d\'alarmes non acquittées',
      [ALARM_METRIC_PARAMETERS.notAckedInHourAlarms]: 'Nombre moyen d\'alarmes non acquittées avec une durée de 1 à 4h',
      [ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms]: 'Nombre moyen d\'alarmes non acquittées d\'une durée de 4 à 24h',
      [ALARM_METRIC_PARAMETERS.notAckedInDayAlarms]: 'Nombre moyen d\'alarmes non acquittées datant de plus de 24h',
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: 'Nombre moyen d\'alarmes actives avec tickets créés',
      [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: 'Nombre moyen d\'alarmes actives sans ticket créé',
    },
  },
};
