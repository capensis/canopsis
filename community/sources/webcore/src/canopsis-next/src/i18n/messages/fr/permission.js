import { USER_PERMISSIONS_GROUPS, USER_VIEWS_PERMISSIONS, USER_PERMISSIONS } from '@/constants';

export default {
  title: {
    /**
     * VIEWS PERMISSIONS
     */
    [USER_VIEWS_PERMISSIONS.viewActions]: 'Vue : actions',
    [USER_VIEWS_PERMISSIONS.viewGeneral]: 'Vue : général',

    /**
     * GROUPS
     */
    [USER_PERMISSIONS_GROUPS.commonviews]: 'Vues',
    [USER_PERMISSIONS_GROUPS.commonviewsPlaylist]: 'Listes de lecture',
    [USER_PERMISSIONS_GROUPS.widgets]: 'Widgets',
    [USER_PERMISSIONS_GROUPS.widgetsCommon]: 'Commun',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatistics]: 'Statistiques des alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatisticsWidgetsettings]: 'Paramètres du widget de statistiques des alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatisticsViewsettings]: 'Paramètres de vue du widget Statistiques des alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAvailability]: 'Disponibilité',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityActions]: 'Actions du widget Disponibilité',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityWidgetsettings]: 'Paramètres du widget Disponibilité',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityViewsettings]: 'Paramètres de vue d widget Disponibilité',
    [USER_PERMISSIONS_GROUPS.widgetsBarchart]: 'Histogramme',
    [USER_PERMISSIONS_GROUPS.widgetsBarchartWidgetsettings]: 'Paramètres du widget Histogramme',
    [USER_PERMISSIONS_GROUPS.widgetsBarchartViewsettings]: 'Paramètres de vue du widget Histogramme ',
    [USER_PERMISSIONS_GROUPS.widgetsCounter]: 'Compteur',
    [USER_PERMISSIONS_GROUPS.widgetsContext]: 'Explorateur de contexte',
    [USER_PERMISSIONS_GROUPS.widgetsContextViewsettings]: 'Paramètres de vue du widget Explorateur de contexte',
    [USER_PERMISSIONS_GROUPS.widgetsContextEntityactions]: 'Actions d\'entité du widget Explorateur de contexte',
    [USER_PERMISSIONS_GROUPS.widgetsContextActions]: 'Actions du widget Explorateur de contexte',
    [USER_PERMISSIONS_GROUPS.widgetsContextWidgetsettings]: 'Paramètres du widget Explorateur de contexte',
    [USER_PERMISSIONS_GROUPS.widgetsLinechart]: 'Graphique en ligne',
    [USER_PERMISSIONS_GROUPS.widgetsLinechartWidgetsettings]: 'Paramètres du widget Graphique en ligne',
    [USER_PERMISSIONS_GROUPS.widgetsLinechartViewsettings]: 'Paramètres de vue du widget Graphique en ligne',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslist]: 'Bac alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistAlarmactions]: 'Actions du widget Bac à alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistViewsettings]: 'Paramètres de vue du widget Bac à alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistActions]: 'Actions du widget Bac à alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistWidgetsettings]: 'Paramètres du widget Bac à alarmes',
    [USER_PERMISSIONS_GROUPS.widgetsMap]: 'Cartographie',
    [USER_PERMISSIONS_GROUPS.widgetsMapViewsettings]: 'Paramètres de vue widget Cartographie',
    [USER_PERMISSIONS_GROUPS.widgetsMapWidgetsettings]: 'Paramètres du widget Cartographie',
    [USER_PERMISSIONS_GROUPS.widgetsNumbers]: 'Nombres',
    [USER_PERMISSIONS_GROUPS.widgetsNumbersWidgetsettings]: 'Paramètres du widget Nombres',
    [USER_PERMISSIONS_GROUPS.widgetsNumbersViewsettings]: 'Paramètres de vue du widget Nombres',
    [USER_PERMISSIONS_GROUPS.widgetsPiechart]: 'Diagramme circulaire',
    [USER_PERMISSIONS_GROUPS.widgetsPiechartWidgetsettings]: 'Paramètres du widget Diagramme circulaire',
    [USER_PERMISSIONS_GROUPS.widgetsPiechartViewsettings]: 'Paramètres de vue du widget Diagramme circulaire',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweather]: 'Météo des Services',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherAlarmactions]: 'Actions pour les alarmes du widget Météo des Services',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherViewsettings]: 'Paramètres de vue dwidget Météo des Services',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherWidgetsettings]: 'Paramètres du widget Météo des Services',
    [USER_PERMISSIONS_GROUPS.widgetsTestingweather]: 'Scénario Junit',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatistics]: 'Statistiques Utilisateurs',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatisticsWidgetsettings]: 'Paramètres du widget Statistiques Utilisateurs',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatisticsViewsettings]: 'Paramètres de vue du widget Statistiques Utilisateurs',
    [USER_PERMISSIONS_GROUPS.api]: 'API',
    [USER_PERMISSIONS_GROUPS.apiGeneral]: 'Général',
    [USER_PERMISSIONS_GROUPS.apiRules]: 'Règles',
    [USER_PERMISSIONS_GROUPS.apiRemediation]: 'Remédiation',
    [USER_PERMISSIONS_GROUPS.apiPlanning]: 'Planification',
    [USER_PERMISSIONS_GROUPS.technical]: 'Technique',
    [USER_PERMISSIONS_GROUPS.technicalAdmin]: 'Admin',
    [USER_PERMISSIONS_GROUPS.technicalAdminCommunication]: 'Communication',
    [USER_PERMISSIONS_GROUPS.technicalAdminGeneral]: 'Général',
    [USER_PERMISSIONS_GROUPS.technicalAdminAccess]: 'Accès',
    [USER_PERMISSIONS_GROUPS.technicalExploitation]: 'Exploitation',
    [USER_PERMISSIONS_GROUPS.technicalNotification]: 'Notifications',
    [USER_PERMISSIONS_GROUPS.technicalViewsandwidgets]: 'Vues et widgets',
    [USER_PERMISSIONS_GROUPS.technicalProfile]: 'Profil',
    [USER_PERMISSIONS_GROUPS.technicalToken]: 'Jeton',

    /**
     * Business Common Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.variablesHelp]: 'Voir la liste des variables (dans tous les widgets)',
    [USER_PERMISSIONS.business.context.actions.entityComment]: 'Gérer les commentaires d\'entité (voir, créer, éditer, supprimer)',

    /**
     * Business Alarms List Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.ack]: 'Acquitter',
    [USER_PERMISSIONS.business.alarmsList.actions.fastAck]: 'Acquitter rapidement',
    [USER_PERMISSIONS.business.alarmsList.actions.ackRemove]: 'Annuler l\'acquittement',
    [USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: 'Définir un comportement périodique',
    [USER_PERMISSIONS.business.alarmsList.actions.fastPbehaviorAdd]: 'Définir un comportement périodique rapidement',
    [USER_PERMISSIONS.business.alarmsList.actions.snooze]: 'Mettre en veille',
    [USER_PERMISSIONS.business.alarmsList.actions.declareTicket]: 'Déclarer un ticket',
    [USER_PERMISSIONS.business.alarmsList.actions.associateTicket]: 'Associer un ticket',
    [USER_PERMISSIONS.business.alarmsList.actions.cancel]: 'Annuler l\'alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.unCancel]: 'Annuler la suppression de l\'alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.fastCancel]: 'Annuler rapidement l\'alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.changeState]: 'Changer et verrouiller la criticité',
    [USER_PERMISSIONS.business.alarmsList.actions.history]: 'Accéder à l\'historique des alarmes',
    [USER_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup]: 'Lier à une méta-alarme manuelle / Délier',
    [USER_PERMISSIONS.business.alarmsList.actions.comment]: 'Commenter l\'alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.filter]: 'Définir des filtres d\'alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: 'Filtrer les alarmes',
    [USER_PERMISSIONS.business.alarmsList.actions.remediationInstructionsFilter]: 'Définir des filtres par consigne de remédiation',
    [USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: 'Filtrer les alarmes par consigne de remédiation',
    [USER_PERMISSIONS.business.alarmsList.actions.links]: 'Suivre les liens d\'une alarme',
    [USER_PERMISSIONS.business.alarmsList.actions.correlation]: 'Grouper les alarmes corrélées (méta-alarmes)',
    [USER_PERMISSIONS.business.alarmsList.actions.executeInstruction]: 'Exécuter des consignes manuelles',
    [USER_PERMISSIONS.business.alarmsList.actions.category]: 'Filtrer les alarmes par catégorie',
    [USER_PERMISSIONS.business.alarmsList.actions.exportPdf]: 'Exporter en PDF',
    [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: 'Exporter la liste des alarmes en CSV',
    [USER_PERMISSIONS.business.alarmsList.actions.metaAlarmGroup]: 'Délier l\'alarme de la méta-alarme automatique',
    [USER_PERMISSIONS.business.alarmsList.actions.bookmark]: 'Ajouter / supprimer un signet',
    [USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: 'Filtrer les alarmes avec signet',

    /**
     * Business Context Explorer Permissions
     */
    [USER_PERMISSIONS.business.context.actions.createEntity]: 'Créer une entité',
    [USER_PERMISSIONS.business.context.actions.editEntity]: 'Modifier une entité',
    [USER_PERMISSIONS.business.context.actions.duplicateEntity]: 'Dupliquer une entité',
    [USER_PERMISSIONS.business.context.actions.deleteEntity]: 'Supprimer une entité',
    [USER_PERMISSIONS.business.context.actions.massEnable]: 'Activer les entités sélectionnées en masse',
    [USER_PERMISSIONS.business.context.actions.massDisable]: 'Désactiver les entités sélectionnées en masse',
    [USER_PERMISSIONS.business.context.actions.pbehavior]: 'Définir un comportement périodique',
    [USER_PERMISSIONS.business.context.actions.filter]: 'Définir les filtres des entités',
    [USER_PERMISSIONS.business.context.actions.userFilter]: 'Filtrer les entités',
    [USER_PERMISSIONS.business.context.actions.category]: 'Filtrer les entités par catégorie',
    [USER_PERMISSIONS.business.context.actions.exportAsCsv]: 'Exporter les entités en CSV',

    /**
     * Business Service Weather Permissions
     */
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAck]: 'Acquitter',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket]: 'Associer un ticket',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityDeclareTicket]: 'Déclarer un ticket',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityComment]: 'Commenter l\'alarme',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityValidate]: 'Valider les alarmes et changer leur état en critique',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: 'Invalider les alarmes et les annuler',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPause]: 'Mettre les alarmes en pause (définir le type de comportement P "Pause")',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPlay]: 'Activer les alarmes en pause (retirer le type de comportement P "Pause")',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityCancel]: 'Accéder à la liste des comportements périodiques associés aux services (dans l\'onglet des fenêtres modales des services)',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors]: 'Voir les comportements périodiques des services (dans l\'onglet des fenêtres modales des services)',
    [USER_PERMISSIONS.business.serviceWeather.actions.executeInstruction]: 'Exécuter des consignes manuelles',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityLinks]: 'Suivre les liens des alarmes',
    [USER_PERMISSIONS.business.serviceWeather.actions.moreInfos]: 'Ouvrir la fenêtre modale "Plus d\'infos"',
    [USER_PERMISSIONS.business.serviceWeather.actions.alarmsList]: 'Ouvrir la liste des alarmes disponibles pour chaque service',
    [USER_PERMISSIONS.business.serviceWeather.actions.pbehaviorList]: 'Voir les comportements périodiques des services (dans l\'onglet des fenêtres modales des entités de service)',
    [USER_PERMISSIONS.business.serviceWeather.actions.filter]: 'Définir les filtres d\'alarme',
    [USER_PERMISSIONS.business.serviceWeather.actions.userFilter]: 'Filtrer les alarmes',
    [USER_PERMISSIONS.business.serviceWeather.actions.category]: 'Filtrer les alarmes par catégorie',

    /**
     * Business Counter Permissions
     */
    [USER_PERMISSIONS.business.counter.actions.alarmsList]: 'Voir la liste des alarmes associées aux compteurs',

    /**
     * Business Testing Weather Permissions
     */
    [USER_PERMISSIONS.business.testingWeather.actions.alarmsList]: 'Ouvrir la liste des alarmes disponibles',

    /**
     * Business Map Permissions
     */
    [USER_PERMISSIONS.business.map.actions.alarmsList]: 'Voir la liste des alarmes associées aux points sur les cartes',
    [USER_PERMISSIONS.business.map.actions.filter]: 'Définir les filtres pour les points sur les cartes',
    [USER_PERMISSIONS.business.map.actions.userFilter]: 'Filtrer les points sur les cartes',
    [USER_PERMISSIONS.business.map.actions.category]: 'Filtrer les points sur les cartes par catégorie',

    /**
     * Business Bar Chart Permissions
     */
    [USER_PERMISSIONS.business.barChart.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.barChart.actions.sampling]: 'Modifier l\'échantillonnage pour les données affichées',
    [USER_PERMISSIONS.business.barChart.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.barChart.actions.userFilter]: 'Filtrer les données',

    /**
     * Business Line Chart Permissions
     */
    [USER_PERMISSIONS.business.lineChart.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.lineChart.actions.sampling]: 'Modifier l\'échantillonnage pour les données affichées',
    [USER_PERMISSIONS.business.lineChart.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.lineChart.actions.userFilter]: 'Filtrer les données',

    /**
     * Business Pie Chart Permissions
     */
    [USER_PERMISSIONS.business.pieChart.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.pieChart.actions.sampling]: 'Modifier l\'échantillonnage pour les données affichées',
    [USER_PERMISSIONS.business.pieChart.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.pieChart.actions.userFilter]: 'Filtrer les données',

    /**
     * Business Numbers Permissions
     */
    [USER_PERMISSIONS.business.numbers.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.numbers.actions.sampling]: 'Modifier l\'échantillonnage pour les données affichées',
    [USER_PERMISSIONS.business.numbers.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.numbers.actions.userFilter]: 'Filtrer les données',

    /**
     * Business User Statistics
     */
    [USER_PERMISSIONS.business.userStatistics.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.userStatistics.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.userStatistics.actions.userFilter]: 'Filtrer les données',

    /**
     * Business Alarm Statistics
     */
    [USER_PERMISSIONS.business.alarmStatistics.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.alarmStatistics.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.alarmStatistics.actions.userFilter]: 'Filtrer les données',

    /**
     * Business Availability
     */
    [USER_PERMISSIONS.business.availability.actions.interval]: 'Modifier les intervalles de temps pour les données affichées',
    [USER_PERMISSIONS.business.availability.actions.filter]: 'Définir les filtres de données',
    [USER_PERMISSIONS.business.availability.actions.userFilter]: 'Filtrer les données',
    [USER_PERMISSIONS.business.availability.actions.exportAsCsv]: 'Exporter les disponibilités en fichier CSV',

    /**
     * Technical Admin Communication
     */
    [USER_PERMISSIONS.technical.broadcastMessage]: 'Messages de diffusion',
    [USER_PERMISSIONS.technical.playlist]: 'Listes de lecture',

    /**
     * Technical Admin General
     */
    [USER_PERMISSIONS.technical.eventsRecord]: 'Enregistrements d\'événements',
    [USER_PERMISSIONS.technical.healthcheck]: 'Bilan de santé',
    [USER_PERMISSIONS.technical.healthcheckStatus]: 'Statut du bilan de santé',
    [USER_PERMISSIONS.technical.icon]: 'Paramètres - icônes',
    [USER_PERMISSIONS.technical.kpi]: 'Graphiques KPI',
    [USER_PERMISSIONS.technical.kpiCollectionSettings]: 'Paramètres de collecte des KPI',
    [USER_PERMISSIONS.technical.kpiFilters]: 'Filtres KPI',
    [USER_PERMISSIONS.technical.kpiRatingSettings]: 'Paramètres de notation des KPI',
    [USER_PERMISSIONS.technical.maintenance]: 'Mode maintenance',
    [USER_PERMISSIONS.technical.map]: 'Cartographie',
    [USER_PERMISSIONS.technical.parameters]: 'Paramètres - onglet paramètres',
    [USER_PERMISSIONS.technical.planningExceptions]: 'Dates d\'exceptions de comportement périodique',
    [USER_PERMISSIONS.technical.planningReason]: 'Raisons de comportement périodique',
    [USER_PERMISSIONS.technical.planningType]: 'Type de comportement périodique',
    [USER_PERMISSIONS.technical.remediationConfiguration]: 'Consignes - onglet consignes',
    [USER_PERMISSIONS.technical.remediationInstruction]: 'Consignes - onglet instructions',
    [USER_PERMISSIONS.technical.remediationJob]: 'Consignes - onglet tâches',
    [USER_PERMISSIONS.technical.remediationStatistic]: 'Consignes - onglet statistiques de remédiation',
    [USER_PERMISSIONS.technical.stateSetting]: 'Paramètres de calcul d\'état/sévérité',
    [USER_PERMISSIONS.technical.storageSettings]: 'Paramètres de stockage',
    [USER_PERMISSIONS.technical.tag]: 'Gestion des tags',
    [USER_PERMISSIONS.technical.techmetrics]: 'Bilan de santé - métriques des moteurs',
    [USER_PERMISSIONS.technical.widgetTemplate]: 'Paramètres - modèles de widgets',
    [USER_PERMISSIONS.technical.viewImportExport]: 'Paramètres - import / export',

    /**
     * Technical Admin Access
     */
    [USER_PERMISSIONS.technical.permission]: 'Droits',
    [USER_PERMISSIONS.technical.role]: 'Rôles',
    [USER_PERMISSIONS.technical.user]: 'Utilisateurs',

    /**
     * Technical Admin Exploitation
     */
    [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: 'Règles de déclaration de ticket',
    [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: 'Règles d\'informations dynamiques',
    [USER_PERMISSIONS.technical.exploitation.eventFilter]: 'Filtres d\'événements',
    [USER_PERMISSIONS.technical.exploitation.flappingRules]: 'Règles de bagot',
    [USER_PERMISSIONS.technical.exploitation.idleRules]: 'Règles d\'inactivité',
    [USER_PERMISSIONS.technical.exploitation.linkRule]: 'Générateur de liens',
    [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: 'Règles de méta-alarmes',
    [USER_PERMISSIONS.technical.exploitation.pbehavior]: 'Comportements périodiques',
    [USER_PERMISSIONS.technical.exploitation.resolveRules]: 'Règles de résolution',
    [USER_PERMISSIONS.technical.exploitation.scenario]: 'Scénarios',
    [USER_PERMISSIONS.technical.exploitation.snmpRule]: 'Règles SNMP',

    /**
     * Technical Admin Notification
     */
    [USER_PERMISSIONS.technical.notification.common]: 'Paramètres - paramètres de notification',
    [USER_PERMISSIONS.technical.notification.instructionStats]: 'Statistiques des consignes',

    /**
     * Technical Admin Views and widgets
     */
    [USER_PERMISSIONS.technical.privateView]: 'Vues privées',
    [USER_PERMISSIONS.technical.view]: 'Vues',

    /**
     * Technical Admin Profile
     */
    [USER_PERMISSIONS.technical.profile.theme]: 'Couleurs du thème',
    [USER_PERMISSIONS.technical.profile.corporatePattern]: 'Patterns partagés',

    /**
     * Technical Admin Token
     */
    [USER_PERMISSIONS.technical.shareToken]: 'Paramètres des jetons partagés',

    /**
     * API Permissions General
     */
    [USER_PERMISSIONS.api.general.acl]: 'Rôles, permissions, utilisateurs',
    [USER_PERMISSIONS.api.general.alarmRead]: 'Lire les alarmes',
    [USER_PERMISSIONS.api.general.alarmTag]: 'Gérer les Tags d\'alarmes',
    [USER_PERMISSIONS.api.general.alarmUpdate]: 'Mettre à jour les alarmes',
    [USER_PERMISSIONS.api.general.associativeTable]: 'Tables associatives',
    [USER_PERMISSIONS.api.general.broadcastMessage]: 'Message de diffusion',
    [USER_PERMISSIONS.api.general.theme]: 'Gérer les thèmes',
    [USER_PERMISSIONS.api.general.contextgraph]: 'Importn de référentiels',
    [USER_PERMISSIONS.api.general.corporatePattern]: 'Patterns partagés',
    [USER_PERMISSIONS.api.general.datastorageRead]: 'Lecture des paramètres de stockage des données',
    [USER_PERMISSIONS.api.general.datastorageUpdate]: 'Mise à jour des paramètres de stockage des données',
    [USER_PERMISSIONS.api.general.entity]: 'Entités',
    [USER_PERMISSIONS.api.general.entitycategory]: 'Catégories d\'entités',
    [USER_PERMISSIONS.api.general.entitycomment]: 'Commentaires d\'entités',
    [USER_PERMISSIONS.api.general.entityservice]: 'Entités de type Service',
    [USER_PERMISSIONS.api.general.event]: 'Événements',
    [USER_PERMISSIONS.api.general.exportConfigurations]: 'Exporter les configurations',
    [USER_PERMISSIONS.api.general.files]: 'Fichiers',
    [USER_PERMISSIONS.api.general.healthcheck]: 'Bilan de santé',
    [USER_PERMISSIONS.api.general.icon]: 'Icônes',
    [USER_PERMISSIONS.api.general.junit]: 'JUnit',
    [USER_PERMISSIONS.api.general.kpiFilter]: 'Filtres KPI',
    [USER_PERMISSIONS.api.general.launchEventRecording]: 'Lancer l\'enregistrement des événements',
    [USER_PERMISSIONS.api.general.maintenance]: 'Mode maintenance',
    [USER_PERMISSIONS.api.general.map]: 'Cartes',
    [USER_PERMISSIONS.api.general.messageRateStatsRead]: 'Lire le taux de messages',
    [USER_PERMISSIONS.api.general.metrics]: 'Métriques',
    [USER_PERMISSIONS.api.general.metricsSettings]: 'Paramètres des métriques',
    [USER_PERMISSIONS.api.general.notifications]: 'Paramètres de notification',
    [USER_PERMISSIONS.api.general.playlist]: 'Listes de lecture',
    [USER_PERMISSIONS.api.general.privateViewGroups]: 'Groupes de vues privées',
    [USER_PERMISSIONS.api.general.ratingSettings]: 'Paramètres de notation',
    [USER_PERMISSIONS.api.general.resendEvents]: 'Renvoyer les événements',
    [USER_PERMISSIONS.api.general.shareToken]: 'Partager les jetons',
    [USER_PERMISSIONS.api.general.stateSettings]: 'Paramètres d\'état/sévérité',
    [USER_PERMISSIONS.api.general.techmetrics]: 'Métriques techniques',
    [USER_PERMISSIONS.api.general.techmetricsSettings]: 'Paramètres des métriques techniques',
    [USER_PERMISSIONS.api.general.userInterfaceDelete]: 'Supprimer l\'interface utilisateur',
    [USER_PERMISSIONS.api.general.userInterfaceUpdate]: 'Mettre à jour l\'interface utilisateur',
    [USER_PERMISSIONS.api.general.view]: 'Vues',
    [USER_PERMISSIONS.api.general.viewgroup]: 'Groupes de vues',
    [USER_PERMISSIONS.api.general.widgetTemplate]: 'Modèles de widgets',

    /**
     * API Permissions Rules
     */
    [USER_PERMISSIONS.api.rules.action]: 'Actions',
    [USER_PERMISSIONS.api.rules.declareTicketExecution]: 'Exécuter les règles de déclaration de ticket',
    [USER_PERMISSIONS.api.rules.declareTicketRule]: 'Règles de déclaration de ticket',
    [USER_PERMISSIONS.api.rules.dynamicinfos]: 'Informations dynamiques',
    [USER_PERMISSIONS.api.rules.eventFilter]: 'Filtres d\'événements',
    [USER_PERMISSIONS.api.rules.flappingRule]: 'Règles de bagot',
    [USER_PERMISSIONS.api.rules.idleRule]: 'Règle d\'inactivité',
    [USER_PERMISSIONS.api.rules.linkRule]: 'Générateur de liens',
    [USER_PERMISSIONS.api.rules.metaalarmrule]: 'Règles de méta-alarmes',
    [USER_PERMISSIONS.api.rules.resolveRule]: 'Règles de résolution',
    [USER_PERMISSIONS.api.rules.snmpRule]: 'Règles SNMP',
    [USER_PERMISSIONS.api.rules.snmpMib]: 'MIB SNMP',

    /**
     * API Permissions Remediation
     */
    [USER_PERMISSIONS.api.remediation.execution]: 'Exécuter les consignes',
    [USER_PERMISSIONS.api.remediation.instruction]: 'Consignes',
    [USER_PERMISSIONS.api.remediation.instructionApprove]: 'Approuver les consignes',
    [USER_PERMISSIONS.api.remediation.job]: 'Tâches',
    [USER_PERMISSIONS.api.remediation.jobConfig]: 'Configurations des tâches',

    /**
     * API Permissions Planning
     */
    [USER_PERMISSIONS.api.planning.pbehavior]: 'Comportements périodiques',
    [USER_PERMISSIONS.api.planning.pbehaviorException]: 'Dates d\'exceptions de comportements périodiques',
    [USER_PERMISSIONS.api.planning.pbehaviorReason]: 'Raisons de comportements périodiques',
    [USER_PERMISSIONS.api.planning.pbehaviorType]: 'Types de comportements périodiques',
  },
};
