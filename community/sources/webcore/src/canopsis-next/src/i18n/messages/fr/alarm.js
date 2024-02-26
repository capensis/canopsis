import { EVENT_ENTITY_TYPES, ALARM_METRIC_PARAMETERS } from '@/constants';

export default {
  liveReporting: 'Définir un intervalle de dates',
  ackAuthor: 'Auteur de l\'acquittement',
  lastCommentAuthor: 'Auteur du dernier commentaire',
  lastCommentMessage: 'Message du dernier commentaire',
  metaAlarm: 'Meta-alarmes',
  acknowledge: 'Acquitter',
  ackResources: 'Acquitter les ressources',
  ackResourcesQuestion: 'Voulez-vous acquitter les ressources liées ?',
  actionsRequired: 'Des actions sont requises',
  acknowledgeAndDeclareTicket: 'Acquitter et déclarer un ticket',
  acknowledgeAndAssociateTicket: 'Acquitter et associer un ticket',
  advancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
    + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
    + '<p>Le "-" avant la recherche est obligatoire</p>\n'
    + '<p>Opérateurs:\n'
    + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
    + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
    + '<dl><dt>Exemples :</dt><dt>- Connector = "connector_1"</dt>\n'
    + '    <dd>Alarmes dont le connecteur est "connector_1"</dd><dt>- Connector="connector_1" AND Resource="resource_3"</dt>\n'
    + '    <dd>Alarmes dont le connecteur est "connector_1" et la ressource est "resource_3"</dd><dt>- Connector="connector_1" OR Resource="resource_3"</dt>\n'
    + '    <dd>Alarmes dont le connecteur est "connector_1" ou la ressource est "resource_3"</dd><dt>- Connector LIKE 1 OR Connector LIKE 2</dt>\n'
    + '    <dd>Alarmes dont le connecteur contient 1 ou 2</dd><dt>- NOT Connector = "connector_1"</dt>\n'
    + '    <dd>Alarmes dont le connecteur n\'est pas "connector_1"</dd>\n'
    + '</dl>',
  otherTickets: 'D\'autres tickets sont disponibles dans le panneau d\'expansion',
  noAlarmFound: 'Aucune alarme n\'a été trouvée selon les modèles définis',
  associateTicketResources: 'Ticket associé pour les ressources',
  followLink: 'Suivez le lien "{title}"',
  hasBookmark: 'L\'alarme a un signet',
  filterByBookmark: 'Filtrer par signet',
  popups: {
    exportFailed: 'Impossible d\'exporter la liste des alarmes au format CSV',
    addBookmarkSuccess: 'Le signet a été ajouté',
    addBookmarkFailed: 'Quelque chose s\'est mal passé. Le signet n\'a pas été ajouté',
    removeBookmarkSuccess: 'Le signet a été supprimé',
    removeBookmarkFailed: 'Quelque chose s\'est mal passé. Le signet n\'a pas été supprimé',
  },
  actions: {
    titles: {
      ack: 'Acquitter',
      fastAck: 'Acquitter rapidement',
      ackRemove: 'Annuler l\'acquittement',
      pbehavior: 'Définir un comportement périodique',
      snooze: 'Mettre en veille',
      declareTicket: 'Déclarer un incident',
      associateTicket: 'Associer un ticket',
      cancel: 'Annuler l\'alarme',
      unCancel: 'Annuler l\'annulation de l\'alarme',
      fastCancel: 'Annulation rapide',
      changeState: 'Changer et verrouiller la criticité',
      variablesHelp: 'Liste des variables disponibles',
      history: 'Historique',
      createManualMetaAlarm: 'Gestion manuelle des méta-alarmes',
      removeAlarmsFromManualMetaAlarm: 'Dissocier l\'alarme de la méta-alarme manuelle',
      removeAlarmsFromAutoMetaAlarm: 'Dissocier l\'alarme de la méta-alarme',
      comment: 'Commenter l\'alarme',
      exportPdf: 'Exporter l\'alarme au format PDF',
      addBookmark: 'Ajouter un signet',
      removeBookmark: 'Supprimer le signet',
    },
    iconsTitles: {
      ack: 'Acquitter',
      declareTicket: 'Déclarer un ticket',
      canceled: 'Annulé',
      snooze: 'Mettre en veille',
      pbehaviors: 'Définir un comportement périodique',
      grouping: 'Meta-alarmes',
      comment: 'Commentaire',
    },
    iconsFields: {
      ticketNumber: 'Numéro de ticket',
      parents: 'Causes',
      children: 'Conséquences',
    },
  },
  timeLine: {
    titlePaths: {
      by: 'par',
    },
    stateCounter: {
      header: 'Criticités compressées (depuis le dernier changement de statut)',
      stateIncreased: 'Criticité augmentée',
      stateDecreased: 'Criticité diminuée',
    },
    types: {
      [EVENT_ENTITY_TYPES.ack]: 'Acquittement',
      [EVENT_ENTITY_TYPES.ackRemove]: 'Suppression d\'acquittement',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Association d\'un ticket',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Le ticket est déclaré avec succès',
      [EVENT_ENTITY_TYPES.declareTicketFail]: 'La déclaration du ticket a échoué',
      [EVENT_ENTITY_TYPES.webhookStart]: 'Webhook a démarré',
      [EVENT_ENTITY_TYPES.webhookComplete]: 'Webhook exécuté avec succès',
      [EVENT_ENTITY_TYPES.webhookFail]: 'Webhook échoué',
      [EVENT_ENTITY_TYPES.snooze]: 'Alarme mise en veille',
      [EVENT_ENTITY_TYPES.unsooze]: 'Alarme sortie de veille',
      [EVENT_ENTITY_TYPES.pbhenter]: 'Comportement périodique activé',
      [EVENT_ENTITY_TYPES.pbhleave]: 'Comportement périodique désactivé',
      [EVENT_ENTITY_TYPES.cancel]: 'Alarme annulée',
      [EVENT_ENTITY_TYPES.comment]: 'Alarme commentée',
      [EVENT_ENTITY_TYPES.metaalarmattach]: 'Alarme liée à la méta alarme',
    },
  },
  tabs: {
    moreInfos: 'Plus d\'infos',
    timeLine: 'Chronologie',
    charts: 'Graphiques',
    alarmsChildren: 'Alarmes liées',
    trackSource: 'Cause racine',
    impactChain: 'Chaîne d\'impact',
    entityGantt: 'Diagramme de Gantt',
    ticketsDeclared: 'Tickets déclarés',
  },
  moreInfos: {
    defineATemplate: 'Pour définir le template de cette fenêtre, rendez-vous dans les paramètres du bac à alarmes.',
  },
  infoPopup: 'Info popup',
  tooltips: {
    priority: 'La priorité est égale à la sévérité multipliée par le niveau d\'impact de l\'entité sur laquelle l\'alarme est déclenchée',
    runningManualInstructions: 'Consigne manuelle <strong>{title}</strong> en cours | Consignes manuelles <strong>{title}</strong> en cours',
    runningAutoInstructions: 'Consigne automatique <strong>{title}</strong> en cours | Consignes automatiques <strong>{title}</strong> en cours',
    successfulManualInstructions: 'Consigne manuelle <strong>{title}</strong> réussie | Consignes manuelles <strong>{title}</strong> réussies',
    successfulAutoInstructions: 'Consigne automatique <strong>{title}</strong> réussies | Consignes automatiques <strong>{title}</strong> réussies',
    failedManualInstructions: 'Consigne manuelle <strong>{title}</strong> en échec | Consignes manuelles <strong>{title}</strong> en échec',
    failedAutoInstructions: 'Consigne automatique <strong>{title}</strong> en échec | Consignes automatiques <strong>{title}</strong> en échec',
    hasManualInstruction: 'Il y a une consigne manuelle associée | Il y a des consignes manuelles associées',
    resetChangeColumns: 'Réinitialiser l\'ordre/le redimensionnement des colonnes',
    startChangeColumns: 'Commencer à modifier l\'ordre/le redimensionnement des colonnes',
    finishChangeColumns: 'Terminer la modification de l\'ordre/du redimensionnement des colonnes',
  },
  metrics: {
    [ALARM_METRIC_PARAMETERS.createdAlarms]: 'Nombre d\'alarmes créées',
    [ALARM_METRIC_PARAMETERS.activeAlarms]: 'Nombre d\'alarmes actives',
    [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: 'Nombre d\'alarmes non affichées',
    [ALARM_METRIC_PARAMETERS.instructionAlarms]: 'Nombre d\'alarmes en cours de remédiation automatique',
    [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: 'Nombre d\'alarmes avec comportement périodique',
    [ALARM_METRIC_PARAMETERS.correlationAlarms]: 'Nombre d\'alarmes corrélées',
    [ALARM_METRIC_PARAMETERS.ackAlarms]: 'Nombre d\'alarmes avec acquittement',
    [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: 'Nombre d\'alarmes actives avec acquittement',
    [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: 'Nombre d\'alarmes avec acquittement annulé',
    [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: 'Nombre d\'alarmes actives avec tickets',
    [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: 'Nombre d\'alarmes actives sans tickets',
    [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '% d\'alarmes corrélées',
    [ALARM_METRIC_PARAMETERS.ratioInstructions]: '% d\'alarmes avec remédiation automatique',
    [ALARM_METRIC_PARAMETERS.ratioTickets]: '% d\'alarmes avec tickets créés',
    [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '% d\'alarmes remédiées manuellement',
    [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '% d\'alarmes non affichées',
    [ALARM_METRIC_PARAMETERS.averageAck]: 'Délai moyen d\'acquittement des alarmes',
    [ALARM_METRIC_PARAMETERS.averageResolve]: 'Délai moyen de résolution des alarmes',
    [ALARM_METRIC_PARAMETERS.timeToAck]: 'Délai d\'acquittement des alarmes',
    [ALARM_METRIC_PARAMETERS.timeToResolve]: 'Délai de résolution des alarmes',
    [ALARM_METRIC_PARAMETERS.minAck]: 'Temps minimum pour acquitter les alarmes',
    [ALARM_METRIC_PARAMETERS.maxAck]: 'Délai maximum pour acquitter les alarmes',
    [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: 'Nombre d\'alarmes remédiées manuellement',
    [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: 'Nombre d\'alarmes avec remédiation manuelle',
    [ALARM_METRIC_PARAMETERS.notAckedAlarms]: 'Nombre d\'alarmes non acquittées',
    [ALARM_METRIC_PARAMETERS.notAckedInHourAlarms]: 'Nombre d\'alarmes non acquittées avec durée 1-4h',
    [ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms]: 'Nombre d\'alarmes non acquittées avec durée 4-24h',
    [ALARM_METRIC_PARAMETERS.notAckedInDayAlarms]: 'Nombre d\'alarmes non acquittées de plus de 24h',
    [ALARM_METRIC_PARAMETERS.minResolve]: 'Temps minimum pour résoudre les alarmes',
    [ALARM_METRIC_PARAMETERS.maxResolve]: 'Temps max pour résoudre les alarmes',
  },
  fields: {
    displayName: 'Nom simplifié (DisplayName)',
    assignedInstructions: 'Consignes assignées',
    initialOutput: 'Sortie initiale longue',
    initialLongOutput: 'Sortie longue initiale',
    lastComment: 'Dernier commentaire',
    lastCommentInitiator: 'Initiateur du dernier commentaire',
    ackBy: 'Acquitté par',
    ackMessage: 'Message de l\'acquittement',
    ackInitiator: 'Origine de l\'acquittement',
    stateMessage: 'Message d\'état',
    statusMessage: 'Message de statut',
    totalStateChanges: 'Nombre de changements d\'état',
    stateValue: 'Valeur d\'état',
    statusValue: 'Valeur de statut',
    ackAt: 'Acquitté à',
    stateAt: 'État changé à',
    statusAt: 'Statut a changé à',
    resolved: 'Résolue à',
    activationDate: 'Date d\'activation',
    currentStateDuration: 'Durée de l\'état actuel',
    snoozeDuration: 'Durée de mise en veille',
    pbhInactiveDuration: 'Durée d\'inactivité (comportement périodique)',
    activeDuration: 'Durée active',
    eventsCount: 'Nombre d\'événements',
    extraDetails: 'Détails supplémentaires',
    alarmInfos: 'Informations sur l\'alarme',
    ticketAuthor: 'Auteur du ticket',
    ticketId: 'ID du ticket',
    ticketMessage: 'Message du ticket',
    ticketInitiator: 'Initiateur du ticket',
    ticketCreatedAt: 'Ticket créé à',
    ticketData: 'Données du ticket',
    entityId: 'ID d\'entité',
    entityName: 'Nom de l\'entité',
    entityCategoryName: 'Nom de la catégorie d\'entité',
    entityType: 'Type d\'entité',
    entityComponent: 'Composant d\'entité',
    entityConnector: 'Connecteur d\'entité',
    entityImpactLevel: 'Niveau d\'impact de l\'entité',
    entityKoEvents: 'Événements d\'entité KO',
    entityOkEvents: 'Événements d\'entité OK',
    entityInfos: 'Informations sur l\'entité',
    entityComponentInfos: 'Informations sur les composants de l\'entité',
    entityLastPbehaviorDate: 'Date du dernier comportement de l\'entité',
    openedChildren: 'Conséquences ouvertes',
    closedChildren: 'Conséquences fermées',
    canceledInitiator: 'Initiateur annulé',
    changeState: 'Changer d\'état',
  },
};
