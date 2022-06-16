import { merge } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  STATS_TYPES,
  STATS_CRITICITY,
  STATS_QUICK_RANGES,
  TOURS,
  BROADCAST_MESSAGES_STATUSES,
  USER_PERMISSIONS_PREFIXES,
  REMEDIATION_CONFIGURATION_TYPES,
  PBEHAVIOR_RRULE_PERIODS_RANGES,
  ENGINES_NAMES,
  WIDGET_TYPES,
  SCENARIO_ACTION_TYPES,
  ENTITY_TYPES,
  TEST_SUITE_STATUSES,
  SIDE_BARS,
  STATE_SETTING_METHODS,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
} from '@/constants';

import featureService from '@/services/features';

export default merge({
  common: {
    ok: 'Ok',
    undefined: 'Non défini',
    entity: 'Entité',
    service: 'Service',
    pbehaviors: 'Comportements périodiques',
    widget: 'Widget',
    addWidget: 'Ajouter un widget',
    addTab: 'Ajouter un onglet',
    addPbehavior: 'Ajouter un comportement périodique',
    refresh: 'Rafraîchir',
    toggleEditView: 'Activer/Désactiver le mode édition',
    toggleEditViewSubtitle: 'Si vous souhaitez enregistrer les positions des widgets, vous devez désactiver le mode d\'édition pour cela',
    name: 'Nom',
    description: 'Description',
    author: 'Auteur',
    submit: 'Soumettre',
    cancel: 'Annuler',
    continue: 'Continuer',
    stop: 'Arrêter',
    options: 'Options',
    type: 'Type',
    quitEditing: 'Quitter le mode d\'édition',
    enabled: 'Activé(e)',
    disabled: 'Désactivé(e)',
    login: 'Connexion',
    yes: 'Oui',
    no: 'Non',
    default: 'Défaut',
    confirmation: 'Êtes-vous sûr(e) ?',
    parameters: 'Paramètres',
    by: 'Par',
    date: 'Date',
    comment: 'Commentaire | Commentaires',
    end: 'Fin',
    start: 'Début',
    message: 'Message',
    preview: 'Aperçu',
    recursive: 'Récursif',
    select: 'Sélectionner',
    states: 'Criticités',
    state: 'Criticité',
    sla: 'SLA',
    authors: 'Auteurs',
    stat: 'Statistique',
    trend: 'Tendance',
    users: 'Utilisateurs',
    roles: 'Rôles',
    import: 'Importation',
    export: 'Exportation',
    rights: 'Droits',
    profile: 'Profil',
    username: 'Identifiant utilisateur',
    password: 'Mot de passe',
    authKey: 'Auth. key',
    widgetId: 'Widget id',
    connect: 'Connexion',
    optional: 'Optionnel',
    logout: 'Se déconnecter',
    title: 'Titre',
    save: 'Sauvegarder',
    label: 'Label',
    field: 'Champs',
    value: 'Valeur',
    limit: 'Limite',
    add: 'Ajouter',
    create: 'Créer',
    delete: 'Supprimer',
    show: 'Afficher',
    edit: 'Éditer',
    duplicate: 'Dupliquer',
    play: 'Play',
    copyLink: 'Copier le lien',
    parse: 'Compiler',
    home: 'Accueil',
    step: 'Étape',
    paginationItems: 'Affiche {first} à {last} sur {total} Entrées',
    apply: 'Appliquer',
    from: 'Depuis',
    to: 'Vers',
    tags: 'tags',
    actionsLabel: 'Actions',
    noResults: 'Pas de résultats',
    exploitation: 'Exploitation',
    administration: 'Administration',
    forbidden: 'Accès refusé',
    notFound: 'Introuvable',
    search: 'Recherche',
    filters: 'Filtres',
    filter: 'Filtre',
    emptyObject: 'Objet vide',
    startDate: 'Date de début',
    endDate: 'Date de fin',
    links: 'Liens',
    stack: 'Pile',
    edition: 'Édition',
    broadcastMessages: 'Diffuser des messages',
    playlists: 'Playlists',
    planning: 'Planification',
    remediation: 'Remédiation',
    instructions: 'Consignes',
    metaAlarmRule: 'Meta alarm rule',
    dynamicInfo: 'Informations dynamiques',
    icon: 'Icône',
    fullscreen: 'Plein écran',
    interval: 'Période',
    status: 'Statut',
    unit: 'Unité',
    delay: 'Intervalle',
    begin: 'Commencer',
    timezone: 'Fuseau horaire',
    reason: 'Raison',
    or: 'OU',
    and: 'ET',
    priority: 'Priorité',
    clear: 'Clair',
    deleteAll: 'Tout supprimer',
    payload: 'Payload',
    output: 'Note',
    created: 'Date de création',
    updated: 'Date de dernière modification',
    pattern: 'Pattern | Patterns',
    correlation: 'Corrélation',
    periods: 'Périodes',
    range: 'Gamme',
    duration: 'Durée',
    engines: 'Engines',
    previous: 'Précédent',
    next: 'Suivant',
    eventPatterns: 'Patterns des événements',
    alarmPatterns: 'Patterns des alarmes',
    entityPatterns: 'Pattern des entités',
    totalEntityPatterns: 'Total des modèles d\'entité',
    addFilter: 'Ajouter un filtre',
    id: 'Id',
    reset: 'Réinitialiser',
    selectColor: 'Sélectionnez la couleur',
    triggers: 'Triggers',
    disableDuringPeriods: 'Désactiver pendant les pauses',
    retryDelay: 'Intervalle',
    retryUnit: 'Unit',
    retryCount: 'Nombre d\'essais après échec',
    ticket: 'Ticket',
    method: 'Méthode',
    url: 'URL',
    category: 'Catégorie',
    infos: 'Infos',
    impactLevel: 'Niveau d\'impact',
    loadMore: 'Charger plus',
    initiator: 'Initiateur',
    download: 'Télécharger',
    percent: 'Pourcent | Pourcentages',
    tests: 'Tests',
    total: 'Total',
    errors: 'Erreurs',
    failures: 'Failures',
    skipped: 'Ignoré',
    current: 'Actuel',
    average: 'Moyenne',
    information: 'Information | Informations',
    file: 'Déposer',
    group: 'Grouper | Groupes',
    view: 'Vue | Vues',
    tab: 'Onglet | Onglets',
    access: 'Accès',
    communication: 'Communication | Communication',
    general: 'Général',
    actions: {
      close: 'Fermer',
      acknowledgeAndDeclareTicket: 'Acquitter et déclarer un ticket',
      acknowledgeAndAssociateTicket: 'Acquitter et associer un ticket',
      saveChanges: 'Sauvegarder',
      reportIncident: 'Signaler un incident',
      [EVENT_ENTITY_TYPES.ack]: 'Acquitter',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Déclarer un incident',
      [EVENT_ENTITY_TYPES.validate]: 'Valider',
      [EVENT_ENTITY_TYPES.invalidate]: 'Invalider',
      [EVENT_ENTITY_TYPES.pause]: 'Pause',
      [EVENT_ENTITY_TYPES.play]: 'Supprimer la pause',
      [EVENT_ENTITY_TYPES.cancel]: 'Annuler',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Associer un ticket',
      [EVENT_ENTITY_TYPES.comment]: 'Commenter l\'alarme',
    },
    times: {
      second: 'seconde | secondes',
      minute: 'minute | minutes',
      hour: 'heure | heures',
      day: 'jour | jours',
      week: 'semaine | semaines',
      month: 'mois | mois',
      year: 'année | années',
    },
  },
  variableTypes: {
    string: 'String',
    number: 'Nombre',
    boolean: 'Booléen',
    null: 'Nul',
    array: 'Array',
  },
  user: {
    role: 'Rôle',
    defaultView: 'Vue par défaut',
    seeProfile: 'Voir le profil',
    selectDefaultView: 'Sélectionner une vue par défaut',
  },
  context: {
    impacts: 'Impacts',
    dependencies: 'Dépendances',
    expandPanel: {
      infos: 'Les informations',
      type: 'Type',
      enabled: 'Activé',
      disabled: 'Désactivé',
      infosSearchLabel: 'Rechercher une info',
      tabs: {
        main: 'Principal',
        pbehaviors: 'Comportements périodiques',
        impactDepends: 'Impacts/Dépendances',
        infos: 'Infos',
        treeOfDependencies: 'Arbre des dépendances',
        impactChain: 'Chaîne d\'impact',
      },
    },
    actions: {
      titles: {
        editEntity: 'Éditer l\'entité',
        duplicateEntity: 'Dupliquer l\'entité',
        deleteEntity: 'Supprimer l\'entité',
        pbehavior: 'Comportement périodique',
        variablesHelp: 'Liste des variables disponibles',
      },
    },
    entityInfo: {
      valueAsList: 'Changer le type de valeur en liste',
    },
    fab: {
      common: 'Ajouter une nouvelle entité',
      addService: 'Ajouter une nouvelle entité de service',
    },
  },
  search: {
    alarmAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n' +
    '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n' +
    '<p>Le "-" avant la recherche est obligatoire</p>\n' +
    '<p>Opérateurs:\n' +
    '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n' +
    '<p>Les types de valeurs : String entre doubles guillemets, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n' +
    '<dl><dt>Exemples :</dt><dt>- Connector = "connector_1"</dt>\n' +
    '    <dd>Alarmes dont le connecteur est "connector_1"</dd><dt>- Connector="connector_1" AND Resource="resource_3"</dt>\n' +
    '    <dd>Alarmes dont le connecteur est "connector_1" et la ressource est "resource_3"</dd><dt>- Connector="connector_1" OR Resource="resource_3"</dt>\n' +
    '    <dd>Alarmes dont le connecteur est "connector_1" ou la ressource est "resource_3"</dd><dt>- Connector LIKE 1 OR Connector LIKE 2</dt>\n' +
    '    <dd>Alarmes dont le connecteur contient 1 ou 2</dd><dt>- NOT Connector = "connector_1"</dt>\n' +
    '    <dd>Alarmes dont le connecteur n\'est pas "connector_1"</dd>\n' +
    '</dl>',
    contextAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n' +
      '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n' +
      '<p>Le "-" avant la recherche est obligatoire</p>\n' +
      '<p>Opérateurs:\n' +
      '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n' +
      '<p>Les types de valeurs : String entre doubles guillemets, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n' +
      '<dl><dt>Exemples :</dt><dt>- Name = "name_1"</dt>\n' +
      '    <dd>Entités dont le names est "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n' +
      '    <dd>Entités dont le names est "name_1" et la types est "service"</dd><dt>- infos.custom.value="Custom value" OR Type="resource"</dt>\n' +
      '    <dd>Entités dont le infos.custom.value est "Custom value" ou la type est "resource"</dd><dt>- infos.custom.value LIKE 1 OR infos.custom.value LIKE 2</dt>\n' +
      '    <dd>Entités dont le infos.custom.value contient 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n' +
      '    <dd>Entités dont le name n\'est pas "name_1"</dd>\n' +
      '</dl>',
    dynamicInfoAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n' +
      '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n' +
      '<p>Le "-" avant la recherche est obligatoire</p>\n' +
      '<p>Opérateurs:\n' +
      '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n' +
      '<p>Pour effectuer une recherche dans les "patterns", utilisez le mot-clé "pattern" comme &lt;NomColonne&gt;</p>\n' +
      '<p>Les types de valeurs : String entre doubles guillemets, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n' +
      '<dl><dt>Examples :</dt><dt>- description = "testdyninfo"</dt>\n' +
      '    <dd>Règles d\'Informations dynamiques dont la description est "testdyninfo"</dd><dt>- pattern = "SEARCHPATTERN1"</dt>\n' +
      '    <dd>Règles d\'Informations dynamiques dont un des patterns est "SEARCHPATTERN1"</dd><dt>- pattern LIKE "SEARCHPATTERN2"</dt>\n' +
      '    <dd>Règles d\'Informations dynamiques dont un des patterns contient "SEARCHPATTERN2"</dd>' +
      '</dl>',
    submit: 'Rechercher',
    clear: 'Ne plus appliquer cette recherche',
  },
  login: {
    base: 'Standard',
    LDAP: 'LDAP',
    loginWithCAS: 'Se connecter avec CAS',
    loginWithSAML: 'Se connecter avec SAML',
    documentation: 'Documentation',
    forum: 'Forum',
    website: 'Canopsis.com',
    connectionProtocols: 'Modes de connexion',
    errors: {
      incorrectEmailOrPassword: 'Mot de passe / Email incorrect',
    },
  },
  alarmList: {
    actions: {
      titles: {
        ack: 'Ack',
        fastAck: 'Ack rapide',
        ackRemove: 'Annuler ack',
        pbehavior: 'Comportement périodique',
        snooze: 'Mettre en veille',
        pbehaviorList: 'Lister les comportements périodiques',
        declareTicket: 'Déclarer un incident',
        associateTicket: 'Associer un ticket',
        cancel: 'Annuler l\'alarme',
        changeState: 'Changer et verrouiller la criticité',
        variablesHelp: 'Liste des variables disponibles',
        history: 'Historique',
        groupRequest: 'Proposition de regroupement pour meta alarmes',
        manualMetaAlarmGroup: 'Gestion manuelle des méta-alarmes',
        manualMetaAlarmUngroup: 'Dissocier l\'alarme de la méta-alarme manuelle',
        comment: 'Commenter l\'alarme',
        executeInstruction: 'Exécuter la consigne "{instructionName}"',
        resumeInstruction: 'Reprendre la consigne "{instructionName}"',
      },
      iconsTitles: {
        ack: 'Ack',
        declareTicket: 'Déclarer un incident',
        canceled: 'Annulé',
        snooze: 'Mettre en veille',
        pbehaviors: 'Comportement périodique',
        grouping: 'Meta alarmes',
        comment: 'Commentaire',
      },
      iconsFields: {
        ticketNumber: 'Numéro de ticket',
        causes: 'Causes',
        consequences: 'Conséquences',
        rule: 'Règle | Règles',
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
        [EVENT_ENTITY_TYPES.stateinc]: 'Augmentation de la criticité',
        [EVENT_ENTITY_TYPES.statedec]: 'Diminution de la criticité',
        [EVENT_ENTITY_TYPES.statusinc]: 'Augmentation du statut',
        [EVENT_ENTITY_TYPES.statusdec]: 'Diminution du statut',
        [EVENT_ENTITY_TYPES.assocTicket]: 'Association d\'un ticket',
        [EVENT_ENTITY_TYPES.declareTicket]: 'Déclaration d\'un ticket',
        [EVENT_ENTITY_TYPES.snooze]: 'Alarme mise en veille',
        [EVENT_ENTITY_TYPES.unsooze]: 'Alarme sortie de veille',
        [EVENT_ENTITY_TYPES.changeState]: 'Changement et verrouillage de la criticité',
        [EVENT_ENTITY_TYPES.pbhenter]: 'Comportement périodique activé',
        [EVENT_ENTITY_TYPES.pbhleave]: 'Comportement périodique désactivé',
        [EVENT_ENTITY_TYPES.cancel]: 'Alarme annulée',
        [EVENT_ENTITY_TYPES.comment]: 'Alarme commentée',
        [EVENT_ENTITY_TYPES.metaalarmattach]: 'Alarme liée à la méta alarme',
        [EVENT_ENTITY_TYPES.instructionStart]: 'L\'exécution de la consigne a été déclenchée',
        [EVENT_ENTITY_TYPES.instructionPause]: 'L\'exécution de la consigne a été mise en pause',
        [EVENT_ENTITY_TYPES.instructionResume]: 'L\'exécution de la consigne a été reprise',
        [EVENT_ENTITY_TYPES.instructionComplete]: 'L\'exécution de la consigne a été terminée',
        [EVENT_ENTITY_TYPES.instructionAbort]: 'L\'exécution de la consigne a été abandonnée',
        [EVENT_ENTITY_TYPES.instructionFail]: 'L\'exécution de la consigne a échoué',
        [EVENT_ENTITY_TYPES.instructionJobStart]: 'L\'exécution d\'un job de remédiation a été démarrée',
        [EVENT_ENTITY_TYPES.instructionJobComplete]: 'L\'exécution du job de remédiation est terminée',
        [EVENT_ENTITY_TYPES.instructionJobAbort]: 'L\'exécution du job de remédiation a été abandonnée',
        [EVENT_ENTITY_TYPES.instructionJobFail]: 'L\'exécution du job de remédiation a échouée',
        [EVENT_ENTITY_TYPES.junitTestSuiteUpdate]: 'La suite de tests a été mise à jour',
        [EVENT_ENTITY_TYPES.junitTestCaseUpdate]: 'Le cas de test a été mis à jour',
      },
    },
    tabs: {
      moreInfos: 'Plus d\'infos',
      timeLine: 'Chronologie',
      alarmsConsequences: 'Alarmes liées',
      alarmsCauses: 'Causes des alarmes',
      trackSource: 'Source de la piste',
      impactChain: 'Chaîne d\'impact',
      entityGantt: 'Diagramme de Gantt',
    },
    moreInfos: {
      defineATemplate: 'Pour définir le template de cette fenêtre, rendez-vous dans les paramètres du bac à alarmes.',
    },
    infoPopup: 'Info popup',
    instructionInfoPopup: 'Au moins une consigne est attachée à cette alarme',
    priorityPopup: 'Le paramètre de priorité est dérivé de la gravité de l\'alarme multipliée par le niveau d\'impact de l\'entité sur laquelle l\'alarme est déclenchée',
  },
  weather: {
    moreInfos: 'Plus d\'infos',
  },
  pbehaviors: {
    connector: 'Type de connecteur',
    connectorName: 'Nom du connecteur',
    isEnabled: 'Est actif',
    begins: 'Débute',
    ends: 'Se termine',
    type: 'Type',
    reason: 'Raison',
    rrule: 'Récurrence',
    status: 'Statut',
    created: 'Date de création',
    updated: 'Date de dernière modification',
    tabs: {
      filter: 'Filtre',
      comments: 'Commentaires',
      entities: 'Entités',
    },
  },
  settings: {
    titles: {
      [SIDE_BARS.alarmSettings]: 'Paramètres du bac à alarmes',
      [SIDE_BARS.contextSettings]: 'Paramètres de l\'explorateur de contexte',
      [SIDE_BARS.serviceWeatherSettings]: 'Paramètres de la météo des services',
      [SIDE_BARS.statsHistogramSettings]: 'Paramètres de l\'histogramme',
      [SIDE_BARS.statsCurvesSettings]: 'Paramètres de courbes de stats',
      [SIDE_BARS.statsTableSettings]: 'Paramètres du tableau de stats',
      [SIDE_BARS.statsCalendarSettings]: 'Paramètres du calendrier',
      [SIDE_BARS.statsNumberSettings]: 'Paramètres du compteur de stats',
      [SIDE_BARS.statsParetoSettings]: 'Paramètres du diagramme de Pareto',
      [SIDE_BARS.textSettings]: 'Paramètres du widget de texte',
      [SIDE_BARS.counterSettings]: 'Paramètres du widget de compteur',
      [SIDE_BARS.testingWeatherSettings]: 'Tester la météo',
    },
    advancedSettings: 'Paramètres avancés',
    widgetTitle: 'Titre du widget',
    columnName: 'Nom de la colonne',
    defaultSortColumn: 'Colonne de tri par défaut',
    sortColumnNoData: 'Appuyez sur <kbd>enter</kbd> pour en créer une nouvelle',
    columnNames: 'Nom des colonnes',
    exportColumnNames: 'Nom des colonnes des exporter',
    groupColumnNames: 'Nom des colonnes des meta alarmes',
    trackColumnNames: 'Suivre les colonnes de source d\'alarme',
    treeOfDependenciesColumnNames: 'Noms de colonne pour l\'arborescence des dépendances',
    orderBy: 'Trier par',
    periodicRefresh: 'Rafraichissement périodique',
    defaultNumberOfElementsPerPage: 'Nombre d\'élements par page par défaut',
    elementsPerPage: 'Élements par page',
    filterOnOpenResolved: 'Filtre sur Open/Resolved',
    open: 'Ouverte',
    resolved: 'Résolue',
    filters: 'Filtres',
    filterEditor: 'Éditeur de filtre',
    isAckNoteRequired: 'Champ \'Note\' requis lors d\'un ack ?',
    isSnoozeNoteRequired: 'Champ \'Note\' requis lorsque d\'un snooze ?',
    linksCategoriesAsList: 'Afficher les liens sous forme de liste ?',
    linksCategoriesLimit: 'Nombre d\'éléments de catégorie',
    isMultiAckEnabled: 'Ack multiple',
    fastAckOutput: 'Commentaire d\'Ack rapide',
    isHtmlEnabledOnTimeLine: 'HTML activé dans la chronologie ?',
    isCorrelationEnabled: 'Corrélation activée ?',
    duration: 'Durée',
    tstop: 'Date de fin',
    periodsNumber: 'Nombre d\'étapes',
    statName: 'Nom de la statistique',
    stats: 'Statistiques',
    statsSelect: {
      title: 'Sélecteur de statistique',
      required: 'Veuillez sélectionner au moins une statistique',
      draggable: 'Essayez de faire glisser un élément',
    },
    yesNoMode: 'Mode Oui/Non',
    selectAFilter: 'Sélectionner un filtre',
    exportAsCsv: 'Exporter le widget sous forme de fichier csv',
    criticityLevels: 'Niveaux de criticité',
    exportCsv: {
      title: 'Exporter CSV',
      fields: {
        separator: 'Séparateur',
        datetimeFormat: 'Format date/heure',
      },
    },
    colorsSelector: {
      title: 'Sélecteur de couleur',
      statsCriticity: {
        [STATS_CRITICITY.ok]: 'ok',
        [STATS_CRITICITY.minor]: 'mineur',
        [STATS_CRITICITY.major]: 'majeur',
        [STATS_CRITICITY.critical]: 'critique',
      },
    },
    statsDateInterval: {
      monthPeriodInfo: "Avec une période 'au mois', les dates de début/fin de calcul des statistiques seront arrondies au 1er jour du mois, à 00:00 UTC",
      fields: {
        quickRanges: 'Valeurs usuelles',
      },
      quickRanges: {
        [STATS_QUICK_RANGES.custom.value]: 'Personnalisé',
        [STATS_QUICK_RANGES.last2Days.value]: '2 derniers jours',
        [STATS_QUICK_RANGES.last7Days.value]: '7 derniers jours',
        [STATS_QUICK_RANGES.last30Days.value]: '30 derniers jours',
        [STATS_QUICK_RANGES.last1Year.value]: 'Dernière année',
        [STATS_QUICK_RANGES.yesterday.value]: 'Hier',
        [STATS_QUICK_RANGES.previousWeek.value]: 'Dernière semaine',
        [STATS_QUICK_RANGES.previousMonth.value]: 'Dernier mois',
        [STATS_QUICK_RANGES.today.value]: 'Aujourd\'hui',
        [STATS_QUICK_RANGES.todaySoFar.value]: 'Aujourd\'hui jusqu\'à maintenant',
        [STATS_QUICK_RANGES.thisWeek.value]: 'Cette semaine',
        [STATS_QUICK_RANGES.thisWeekSoFar.value]: 'Cette semaine jusqu\'à maintenant',
        [STATS_QUICK_RANGES.thisMonth.value]: 'Ce mois',
        [STATS_QUICK_RANGES.thisMonthSoFar.value]: 'Ce mois jusqu\'à maintenant',
        [STATS_QUICK_RANGES.last1Hour.value]: 'Dernière heure',
        [STATS_QUICK_RANGES.last3Hour.value]: '3 dernières heures',
        [STATS_QUICK_RANGES.last6Hour.value]: '6 dernières heures',
        [STATS_QUICK_RANGES.last12Hour.value]: '12 dernières heures',
        [STATS_QUICK_RANGES.last24Hour.value]: '24 dernières heures',
      },
    },
    statsNumbers: {
      title: 'Cellule de stats',
      yesNoMode: 'Mode Oui/Non',
      defaultStat: 'Défaut : Alarmes créées',
      sortOrder: 'Sens de tri',
      displayMode: 'Mode d\'affichage',
      selectAColor: 'Sélectionner une couleur',
    },
    infoPopup: {
      title: 'Info popup',
      fields: {
        column: 'Column',
        template: 'Template',
      },
    },
    rowGridSize: {
      title: 'Taille du widget',
      noData: 'Aucune ligne correspondante. Appuyez sur <kbd>enter</kbd> pour en créer une nouvelle',
      fields: {
        row: 'Ligne',
      },
    },
    moreInfosModal: 'Fenêtre "Plus d\'infos"',
    expandGridRangeSize: 'Largeur-position "Plus d\'infos / chronologie"',
    weatherTemplate: 'Template - Tuiles',
    modalTemplate: 'Template - Modale',
    entityTemplate: 'Template - Entités',
    blockTemplate: 'Template - Tuiles',
    columnSM: 'Colonnes - Petit',
    columnMD: 'Colonnes - Moyen',
    columnLG: 'Colonnes - Large',
    limit: 'Limite',
    height: 'Hauteur',
    margin: {
      title: 'Marges',
      top: 'Marge - Haut',
      right: 'Marge - Droite',
      bottom: 'Marge - Bas',
      left: 'Marge - Gauche',
    },
    contextTypeOfEntities: {
      title: 'Type d\'entité',
      fields: {
        [ENTITY_TYPES.component]: 'Composant',
        [ENTITY_TYPES.connector]: 'Type de connecteur',
        [ENTITY_TYPES.resource]: 'Ressource',
        [ENTITY_TYPES.service]: 'Service',
      },
    },
    statSelector: {
      error: {
        alreadyExist: 'Une statistique portant ce nom existe déjà.',
      },
    },
    statsGroups: {
      title: 'Groupes de statistiques',
      manageGroups: 'Ajouter un groupe',
      required: 'Veuillez créer au moins un groupe',
    },
    statsColor: {
      title: 'Couleurs des statistiques',
      pickColor: 'Sélectionner une couleur',
    },
    statsAnnotationLine: {
      title: 'Ligne repère',
      enabled: 'Activée',
      value: 'Valeur',
      label: 'Label',
      pickLineColor: 'Couleur de la ligne',
      pickLabelColor: 'Couleur du label',
    },
    statsPointsStyles: {
      title: 'Forme des points',
    },
    considerPbehaviors: {
      title: 'Prendre en compte les comportements périodiques ?',
    },
    serviceWeatherModalTypes: {
      title: 'Type de modale',
      fields: {
        moreInfo: 'Plus d\'infos',
        alarmList: 'Bac à alarmes',
        both: 'Les deux',
      },
    },
    templateEditor: 'Modèle',
    columns: {
      isHtml: 'Est-ce du HTML ?',
      isState: 'Affiché comme une criticité ?',
    },
    liveReporting: {
      title: 'Suivi personnalisé',
    },
    counterLevels: {
      title: 'Niveaux',
      fields: {
        counter: 'Compteur',
      },
    },
    counters: 'Compteurs',
    remediationInstructionsFilters: 'Filtres de consignes',
    colorIndicator: {
      title: 'Indicateur de couleur',
      fields: {
        displayAsSeverity: 'Afficher comme gravité',
        displayAsPriority: 'Afficher en priorité',
      },
    },
    receiveByApi: 'Recevoir par l\'API',
    serverStorage: 'Stockage serveur',
    filenameRecognition: 'Reconnaissance du nom de fichier',
    resultDirectory: 'Stockage des résultats de test',
    screenshotDirectories: {
      title: 'Paramètres de stockage des captures d\'écran',
      helpText: 'Définir où les captures d\'écran sont stockées',
    },
    screenshotMask: {
      title: 'Masque de nom de fichier de captures d\'écran',
      helpText: '<dl>' +
        '<dt>Définissez le masque de nom de fichier dont les captures d\'écran sont créées à l\'aide des variables suivantes:<dt>\n' +
        '<dd>- nom du cas de test %test_case%</dd>\n' +
        '<dd>- date (YYYY, MM, DD)</dd>\n' +
        '<dd>- temps d\'exécution (hh, mm, ss)</dd>' +
        '</dl>',
    },
    videoDirectories: {
      title: 'Paramètres de stockage vidéo',
      helpText: 'Définir où la vidéo est stockée',
    },
    videoMask: {
      title: 'Masque de nom de fichier de vidéos',
      helpText: '<dl>' +
        '<dt>Définissez le masque de nom de fichier dont les vidéos sont créées à l\'aide des variables suivantes:<dt>\n' +
        '<dd>- nom du cas de test %test_case%</dd>\n' +
        '<dd>- date (YYYY, MM, DD)</dd>\n' +
        '<dd>- temps d\'exécution (hh, mm, ss)</dd>' +
        '</dl>',
    },
  },
  modals: {
    common: {
      titleButtons: {
        minimizeTooltip: 'Vous avez déjà réduit la fenêtre modale',
      },
    },
    contextInfos: {
      title: 'Infos sur l\'entité',
    },
    createEntity: {
      create: {
        title: 'Créer une entité',
      },
      edit: {
        title: 'Éditer une entité',
      },
      duplicate: {
        title: 'Dupliquer une entité',
      },
      success: {
        create: 'Entité créée avec succès !',
        edit: 'Entité éditée avec succès !',
        duplicate: 'Entité dupliquée avec succès !',
      },
    },
    createService: {
      create: {
        title: 'Créer un service',
      },
      edit: {
        title: 'Éditer un service',
      },
      duplicate: {
        title: 'Dupliquer un service',
      },
      success: {
        create: 'Service créé avec succès !',
        edit: 'Service édité avec succès !',
        duplicate: 'Service dupliqué avec succès !',
      },
    },
    createEntityInfo: {
      create: {
        title: 'Ajouter une information',
      },
      edit: {
        title: 'Éditer une information',
      },
    },
    view: {
      select: {
        title: 'Sélectionner une vue',
      },
      create: {
        title: 'Créer une vue',
      },
      edit: {
        title: 'Éditer une vue',
      },
      duplicate: {
        title: 'Dupliquer une vue - {viewTitle}',
        infoMessage: 'Vous êtes en train de dupliquer une vue. Toutes les lignes et les widgets de la vue dupliquée seront copiés dans la nouvelle vue.',
      },
      noData: 'Aucun groupe correspondant. Appuyez sur <kbd>enter</kbd> pour en créer un nouveau.',
      fields: {
        periodicRefresh: 'Rafraichissement périodique',
        groupIds: 'Choisissez une groupe, ou créez-en un nouveau',
        groupTags: 'Labels de groupes',
      },
      success: {
        create: 'Nouvelle vue créée !',
        edit: 'Vue éditée avec succès !',
        delete: 'Vue supprimée avec succès !',
      },
      fail: {
        create: 'Erreur lors de la création de la vue...',
        edit: 'Erreur lors de  l\'édition de la vue...',
        delete: 'Erreur lors de la suppression de la vue...',
      },
    },
    createEvent: {
      fields: {
        output: 'Note',
      },
    },
    createAckEvent: {
      title: 'Acquitter',
      tooltips: {
        ackResources: 'Voulez-vous acquitter les ressources liées ?',
      },
      fields: {
        ticket: 'Numéro du ticket',
        output: 'Note',
        ackResources: 'Ack ressources',
      },
    },
    confirmAckWithTicket: {
      continueAndAssociateTicket: 'Continuer et associer un ticket',
      infoMessage: `Un numéro de ticket a été renseigné.
        Peut-être souhaitiez-vous associer un ticket à cet incident.
        Si tel est le cas, cliquez sur le bouton "Continuer et associer un ticket".
        Pour continuer l'action d'acquittement sans prendre en compte le numéro de ticket,
        cliquez sur le bouton "Continuer".`,
    },
    createSnoozeEvent: {
      title: 'Snooze',
      fields: {
        duration: 'Durée',
      },
    },
    createCancelEvent: {
      title: 'Annuler',
    },
    createGroupRequestEvent: {
      title: 'Proposition de regroupement pour meta alarmes',
    },
    createGroupEvent: {
      title: 'Créer une meta alarme',
    },
    createChangeStateEvent: {
      title: 'Changer la сriticité',
      states: {
        ok: 'Info',
        minor: 'Mineur',
        major: 'Majeur',
        critical: 'Critique',
      },
      fields: {
        output: 'Note',
      },
    },
    createPbehavior: {
      create: {
        title: 'Ajouter un comportement périodique',
      },
      edit: {
        title: 'Éditer un comportement périodique',
      },
      duplicate: {
        title: 'Dupliquer un comportement périodique',
      },
      steps: {
        general: {
          title: 'Paramètres généraux',
          dates: 'Dates',
          fields: {
            enabled: 'Activé',
            name: 'Nom',
            reason: 'Raison',
            type: 'Type',
            start: 'Début',
            stop: 'Fin',
            fullDay: 'Toute la journée',
            noEnding: 'Pas de fin',
            startOnTrigger: 'Démarrer sur déclencheur',
          },
        },
        filter: {
          title: 'Filtre',
        },
        rrule: {
          title: 'Règle de récurrence',
          exdate: 'Dates d\'exclusion',
          buttons: {
            addExdate: 'Ajouter une date d\'exclusion',
          },
          fields: {
            rRuleQuestion: 'Ajouter une règle de récurrence au comportement périodique ?',
          },
        },
        comments: {
          title: 'Commentaires',
          buttons: {
            addComment: 'Ajouter un commentaire',
          },
          fields: {
            message: 'Message',
          },
        },
      },
      errors: {
        invalid: 'Invalide',
      },
      success: {
        create: 'Comportement périodique créé avec succès ! Celui-ci peut mettre jusqu\'à 60 sec pour apparaître dans l\'interface',
      },
      cancelConfirmation: 'Certaines informations ont été modifiées et ne seront pas sauvegardées. Voulez-vous vraiment quitter ce menu ?',
    },
    createPause: {
      title: 'Mettre en pause',
      comment: 'Commentaire',
      reason: 'Raison',
    },
    createAckRemove: {
      title: 'Annuler l\'acquittement',
    },
    createDeclareTicket: {
      title: 'Déclarer un incident',
    },
    createAssociateTicket: {
      title: 'Associer un numéro de ticket',
      fields: {
        ticket: 'Numéro du ticket',
      },
      alerts: {
        noAckItems: 'Il y a {count} élément sans accusé de réception. L\'événement Ack pour l\'élément sera envoyé avant. | Il y a {count} éléments sans accusé de réception. Les événements Ack pour les éléments seront envoyés avant.',
      },
    },
    liveReporting: {
      editLiveReporting: 'Suivi personnalisé',
      dateInterval: 'Intervalle de dates',
      today: 'Aujourd\'hui',
      yesterday: 'Hier',
      last7Days: '7 derniers jours',
      last30Days: '30 derniers jours',
      thisMonth: 'Ce mois-ci',
      lastMonth: 'Le mois dernier',
      custom: 'Personnalisé',
      tstart: 'Démarre le',
      tstop: 'Finit le',
    },
    infoPopupSetting: {
      title: 'Info popup',
      add: 'Ajouter',
      column: 'Colonne',
      template: 'Template',
      addInfoPopup: {
        title: 'Ajouter une popup d\'info',
      },
    },
    variablesHelp: {
      variables: 'Variables',
      copyToClipboard: 'Copier dans le presse-papier',
    },
    service: {
      actionPending: 'action(s) en attente',
      refreshEntities: 'Rafraîchir la liste des entités',
      editPbehaviors: 'Éditer les pbehaviors',
      entity: {
        tabs: {
          info: 'Info',
          treeOfDependencies: 'Arbre de dépendances',
        },
      },
    },
    filter: {
      create: {
        title: 'Créer un filtre',
      },
      edit: {
        title: 'Éditer un filtre',
      },
      fields: {
        title: 'Nom',
      },
      emptyFilters: 'Aucun filtre ajouté pour le moment',
    },
    colorPicker: {
      title: 'Sélecteur de couleur',
    },
    textEditor: {
      title: 'Éditeur de texte',
    },
    createWidget: {
      title: 'Sélectionnez un widget',
      types: {
        [WIDGET_TYPES.alarmList]: {
          title: 'Bac à alarmes',
        },
        [WIDGET_TYPES.context]: {
          title: 'Explorateur de contexte',
        },
        [WIDGET_TYPES.serviceWeather]: {
          title: 'Météo des services',
        },
        [WIDGET_TYPES.statsHistogram]: {
          title: 'Histogramme des statistiques',
        },
        [WIDGET_TYPES.statsCurves]: {
          title: 'Courbes de statistiques',
        },
        [WIDGET_TYPES.statsTable]: {
          title: 'Tableau de statistiques',
        },
        [WIDGET_TYPES.statsCalendar]: {
          title: 'Calendrier',
        },
        [WIDGET_TYPES.statsNumber]: {
          title: 'Compteur de statistiques',
        },
        [WIDGET_TYPES.statsPareto]: {
          title: 'Diagramme de Pareto',
        },
        [WIDGET_TYPES.text]: {
          title: 'Texte',
        },
        [WIDGET_TYPES.counter]: {
          title: 'Compteur',
        },
        [WIDGET_TYPES.testingWeather]: {
          title: 'Tester la météo',
        },
      },
    },
    manageHistogramGroups: {
      title: {
        add: 'Ajouter un groupe',
        edit: 'Éditer un groupe',
      },
    },
    addStat: {
      title: {
        add: 'Ajouter une statistique',
        edit: 'Éditer une statistique',
      },
      slaRequired: "La paramètre 'SLA' est obligatoire",
    },
    group: {
      create: {
        title: 'Créer un groupe',
      },
      edit: {
        title: 'Éditer un groupe',
      },
      fields: {
        name: 'Nom',
      },
      errors: {
        isNotEmpty: 'Ce groupe n\'est pas vide',
      },
    },
    alarmsList: {
      title: 'Bac à alarmes',
    },
    createUser: {
      create: {
        title: 'Créer un utilisateur',
      },
      edit: {
        title: 'Éditer un utilisateur',
      },
    },
    createRole: {
      create: {
        title: 'Créer un rôle',
      },
      edit: {
        title: 'Éditer un rôle',
      },
    },
    eventFilterRule: {
      create: {
        title: 'Créer une règle',
        success: 'Règle créée avec succès !',
      },
      duplicate: {
        title: 'Dupliquer une règle',
        success: 'Règle créée avec succès !',
      },
      edit: {
        title: 'Éditer une règle',
        success: 'Règle éditée avec succès !',
      },
      remove: {
        success: 'La règle a bien été supprimée !',
      },
      priority: 'Priorité',
      editPattern: 'Éditer le pattern',
      advanced: 'Avancée',
      addAField: 'Ajouter un champ',
      simpleEditor: 'Éditeur simple',
      field: 'Champ',
      value: 'Valeur',
      advancedEditor: 'Éditeur avancé',
      comparisonRules: 'Règles de comparaison',
      enrichmentOptions: 'Options d\'enrichissement',
      editActions: 'Éditer les actions',
      addAction: 'Ajouter une action',
      editAction: 'Éditer une action',
      actions: 'Actions',
      externalData: 'Données externes',
      onSuccess: 'En cas de succès',
      onFailure: 'En cas d\'échec',
      tooltips: {
        addValueRuleField: 'Ajouter une règle',
        editValueRuleField: 'Éditer la règle',
        addObjectRuleField: 'Ajouter un groupe',
        editObjectRuleField: 'Éditer le groupe',
        removeRuleField: 'Supprimer le groupe/la règle',
        copyFromHelp: '<p>Les variables accessibles sont: <strong>Event</strong></p>' +
          '<i>Quelques exemples:</i> <span>"Event.ExtraInfos.datecustom"</span>',
      },
      actionsTypes: {
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy]: {
          text: 'Copier une valeur d\'un champ d\'événement à un autre',
          message: 'Cette action est utilisée pour copier la valeur d\'un contrôle dans un événement.',
          description: 'Les paramètres de l\'action sont :\n- depuis : le nom du champ dont la valeur doit être copiée. Il peut s\'agir d\'un champ d\'événement, d\'un sous-groupe d\'une expression régulière ou d\'une donnée externe.\n- vers : le nom du champ événement dans lequel la valeur doit être copiée.',
        },
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo]: {
          text: 'Copier une valeur d\'un champ d\'un événement vers une information d\'une entité',
          message: 'Cette action est utilisée pour copier la valeur du champ d\'un événement dans le champ d\'une entité. Notez que l\'entité doit être ajoutée à l\'événement en premier.',
          description: 'Les paramètres de l\'action sont :\n- nom : le nom du champ d\'une entité.\n- description (optionnel) : la description.\n- depuis : le nom du champ dont la valeur doit être copiée. Il peut s\'agir d\'un champ d\'événement, d\'un sous-groupe d\'une expression régulière ou d\'une donnée externe.',
        },
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo]: {
          text: 'Définir une information d\'une entité sur une constante',
          message: 'Cette action permet de définir les informations dynamiques d\'une entité correspondant à l\'événement. Notez que l\'entité doit être ajoutée à l\'événement en premier.',
          description: 'Les paramètres de l\'action sont :\n- nom : le nom du champ.\n- description (optionnel) : la description.\n- valeur : la valeur d\'un champ.',
        },
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate]: {
          text: 'Définir une chaîne d\'informations sur une entité à l\'aide d\'un modèle',
          message: 'Cette action permet de modifier les informations dynamiques d\'une entité correspondant à l\'événement. Notez que l\'entité doit être ajoutée à l\'événement.',
          description: 'Les paramètres de l\'action sont :\n- nom : le nom du champ.\n- description (optionnel) : la description\n- valeur : le modèle utilisé pour déterminer la valeur de la donnée.\nDes modèles {{.Event.NomDuChamp}}, des expressions régulières ou des données externes peuvent être utilisés.',
        },
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField]: {
          text: 'Définir un champ d\'un événement sur une constante',
          message: 'Cette action peut être utilisée pour modifier un champ de l\'événement.',
          description: 'Les paramètres de l\'action sont :\n- nom : le nom du champ.\n- valeur : la nouvelle valeur du champ.',
        },
        [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate]: {
          text: 'Définir un champ de chaîne d\'un événement à l\'aide d\'un modèle',
          message: 'Cette action vous permet de modifier un champ d\'événement à partir d\'un modèle.',
          description: 'Les paramètres de l\'action sont :\n- nom : le nom du champ.\n- valeur : le modèle utilisé pour déterminer la valeur du champ.\n  Des modèles {{.Event.NomDuChamp}}, des expressions régulières ou des données externes peuvent être utilisés.',
        },
      },
    },
    metaAlarmRule: {
      create: {
        title: 'Créer une règle',
        success: 'Règle créée avec succès !',
      },
      duplicate: {
        title: 'Dupliquer une règle',
        success: 'Règle créée avec succès !',
      },
      edit: {
        title: 'Éditer une règle',
        success: 'Règle éditée avec succès !',
      },
      remove: {
        success: 'Rule successfully removed !',
      },
      editPattern: 'Éditer le pattern',
      actions: 'Actions',
    },
    viewTab: {
      create: {
        title: 'Ajouter un onglet',
      },
      edit: {
        title: 'Éditer l\'onglet',
      },
      duplicate: {
        title: 'Dupliquer l\'onglet',
      },
      fields: {
        title: 'Titre',
      },
    },
    statsDateInterval: {
      title: 'Stats - Intervalle de dates',
      fields: {
        periodValue: 'Période',
        periodUnit: 'Unité',
      },
      errors: {
        endDateLessOrEqualStartDate: 'La date de fin doit se situer après la date de début',
      },
      info: {
        monthPeriodUnit: 'Les statistiques calculées seront situées entre {start} et {stop}',
      },
    },
    createSnmpRule: {
      create: {
        title: 'Créer une règle SNMP',
      },
      edit: {
        title: 'Modifier la règle SNMP',
      },
      fields: {
        oid: {
          title: 'OID',
          labels: {
            module: 'Sélectionnez un module MIB',
          },
        },
        output: {
          title: 'Message',
        },
        resource: {
          title: 'Ressource',
        },
        component: {
          title: 'Composant',
        },
        connectorName: {
          title: 'Nom du connecteur',
        },
        state: {
          title: 'Criticité',
          labels: {
            toCustom: 'Pour personnaliser',
            defineVar: 'Définir la variable SNMP correspondante',
            writeTemplate: 'Écrire un modèle',
          },
        },
        moduleMibObjects: {
          vars: 'Champ de match vars SNMP',
          regex: 'Regex',
          formatter: 'Format (groupe de capture avec \\x)',
        },
      },
    },
    selectViewTab: {
      title: 'Sélectionnez l\'onglet',
    },
    createHeartbeat: {
      create: {
        title: 'Créer un heartbeat',
        success: 'Heartbeat créé avec succès !',
      },
      edit: {
        title: 'Modifier le heartbeat',
        success: 'Heartbeat modifié avec succès !',
      },
      duplicate: {
        title: 'Dupliquer un heartbeat',
      },
      remove: {
        success: 'Heartbeat supprimé avec succès !',
      },
      massRemove: {
        success: 'Heartbeats supprimés avec succès !',
      },
      patternRequired: 'Un pattern est requis',
    },
    createDynamicInfo: {
      create: {
        title: 'Créer une information dynamique',
        success: 'Information dynamique créé avec succès !',
      },
      edit: {
        title: 'Éditer une information dynamique',
        success: 'Information dynamique éditée avec succès !',
      },
      duplicate: {
        title: 'Dupliquer une information dynamique',
      },
      remove: {
        success: 'Information dynamique supprimée avec succès !',
      },
      errors: {
        invalid: 'Invalide',
      },
      steps: {
        general: {
          fields: {
            id: 'Id',
            name: 'Nom',
            description: 'Description',
          },
        },
        infos: {
          title: 'Informations',
          validationError: 'Toutes les valeurs doivent être saisies',
        },
        patterns: {
          title: 'Patterns',
          alarmPatterns: 'Patterns des alarmes',
          entityPatterns: 'Patterns des entités',
          validationError: 'Au moins un pattern est requis. Merci d\'ajouter un pattern sur les alarmes et/ou un pattern sur les événements',
        },
      },
    },
    createDynamicInfoInformation: {
      create: {
        title: 'Ajouter une information à la règle d\'information dynamique',
      },
      fields: {
        name: 'Nom',
        value: 'Valeur',
      },
    },
    dynamicInfoTemplatesList: {
      title: 'Modèles d\'informations dynamiques',
    },
    createDynamicInfoTemplate: {
      create: {
        title: 'Créer un modèle d\'informations dynamiques',
      },
      edit: {
        title: 'Éditer un modèle d\'informations dynamiques',
      },
      fields: {
        names: 'Attributs',
      },
      buttons: {
        addName: 'Ajouter un attribut',
      },
      errors: {
        noNames: 'Vous devez ajouter au moins 1 attribut',
      },
      emptyNames: 'Aucun nom ajouté pour le moment',
    },
    importExportViews: {
      title: 'Vues d\'importation / exportation',
      groups: 'Groupes',
      views: 'Vues',
      result: 'Résultat',
    },
    createBroadcastMessage: {
      create: {
        title: 'Créer un message de diffusion',
      },
      edit: {
        title: 'Modifier le message diffusé',
      },
      defaultMessage: 'Votre message ici',
    },
    createCommentEvent: {
      title: 'Ajouter un commentaire',
      fields: {
        comment: 'Commentaire',
      },
    },
    createPlaylist: {
      create: {
        title: 'Créer playlist',
      },
      edit: {
        title: 'Éditée playlist',
      },
      duplicate: {
        title: 'Dupliquer une playlist',
      },
      errors: {
        emptyTabs: 'Merci de ajouter un onglet',
      },
      fields: {
        interval: 'Période',
        unit: 'Unité',
      },
      groups: 'Groupe',
      result: 'Résultat',
      manageTabs: 'Gérer les onglets',
    },
    pbehaviorPlanning: {
      title: 'Comportement périodiques',
    },
    selectExceptionsLists: {
      title: 'Choisissez la liste des exceptions',
    },
    createRrule: {
      title: 'Créer un récurrence',
    },
    createPbehaviorType: {
      title: 'Créer un type',
      iconNameHint: 'Entrez le nom d\'une icône à partir de material.io',
      errors: {
        iconName: 'Le nom est invalide',
      },
      fields: {
        name: 'Nom',
        description: 'Description',
        type: 'Type',
        priority: 'Priorité',
        iconName: 'Nom de l\'icône',
        isSpecialColor: 'Utiliser une couleur spéciale pour le type ?',
      },
    },
    pbehaviorRecurrentChangesConfirmation: {
      title: 'Modifier',
      fields: {
        selected: 'Seulement période sélectionnée',
        all: 'Toutes les périodes',
      },
    },
    createPbehaviorReason: {
      title: 'Créer un reason',
      fields: {
        name: 'Nom',
        description: 'Description',
      },
    },
    createPbehaviorException: {
      title: 'Créer une liste d\'exceptions',
      addDate: 'Ajouter une date',
      fields: {
        name: 'Nom',
        description: 'Description',
      },
      emptyExdates: 'Aucun exdate ajouté pour le moment',
    },
    createManualMetaAlarm: {
      title: 'Gestion manuelle des méta-alarmes',
      noData: 'Aucune méta-alarme correspondante. Appuyez sur <kbd>Entrée</kbd> pour en créer un nouveau',
      fields: {
        metaAlarm: 'Méta-alarme manuelle',
        output: 'Note',
      },
    },
    createRemediationInstruction: {
      create: {
        title: 'Créer une consigne',
        popups: {
          success: '{instructionName} a été créée avec succès',
        },
      },
      edit: {
        title: 'Éditer une consigne',
        popups: {
          success: '{instructionName} a été modifiée avec succès',
        },
      },
    },
    createRemediationConfiguration: {
      create: {
        title: 'Créer une configuration',
        popups: {
          success: '{configurationName} a été créé avec succès',
        },
      },
      edit: {
        title: 'Modifier la configuration',
        popups: {
          success: '{configurationName} a été modifié avec succès',
        },
      },
      types: {
        [REMEDIATION_CONFIGURATION_TYPES.rundeck]: 'Rundeck',
        [REMEDIATION_CONFIGURATION_TYPES.awx]: 'Awx',
      },
      fields: {
        host: 'Hôte',
        token: 'Jeton d\'autorisation',
      },
    },
    createRemediationJob: {
      create: {
        title: 'Créer un job',
        popups: {
          success: '{jobName} a été créé avec succès',
        },
      },
      edit: {
        title: 'Éditer un job',
        popups: {
          success: '{jobName} a été modifié avec succès',
        },
      },
      fields: {
        configuration: 'Configuration',
        jobId: 'Job ID',
      },
      errors: {
        invalidJSON: 'JSON non valide',
      },
      payloadHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>' +
        '<i>Quelques exemples:</i>' +
        '<pre>{\n  resource: "{{ .Alarm.Value.Resource }}",\n  entity: "{{ .Entity.ID }}"\n}</pre>',
    },
    clickOutsideConfirmation: {
      title: 'Êtes-vous sûr(e) ?',
      text: 'Les modifications ne seront pas enregistrées. Êtes-vous sûr(e) ?',
      buttons: {
        save: 'Sauvegarder',
        dontSave: 'Ne pas sauvegarder',
        backToForm: 'Retour au formulaire',
      },
    },
    patterns: {
      title: 'Attribuer des modèles',
    },
    rateInstruction: {
      title: 'Évaluer cette consigne',
      text: 'Dans quelle mesure cette consigne a-t-elle été utile ?',
    },
    createScenario: {
      create: {
        title: 'Créer un scénario',
        success: 'Scénario créé !',
      },
      edit: {
        title: 'Modifier le scénario',
        success: 'Scénario modifié !',
      },
      duplicate: {
        title: 'Dupliquer scénario',
        success: 'Scénario dupliqué !',
      },
      remove: {
        success: 'Scénario supprimé !',
      },
    },
    serviceDependencies: {
      impacts: {
        title: 'Impacts pour {name}',
      },
      dependencies: {
        title: 'Dépendances pour {name}',
      },
    },
    stateSetting: {
      title: 'JUnit test suite state settings',
    },
    defineStorage: {
      title: 'Définir le stockage des résultats',
      field: {
        placeholder: 'Entrez le chemin d\'accès au dossier de résultats',
      },
    },
    defineXMLStorage: {
      title: 'Définir le stockage XML',
      field: {
        placeholder: 'Saisissez le chemin d\'accès au dossier XML',
      },
    },
    defineScreenshotStorage: {
      title: 'Définir le stockage des captures d\'écran',
      field: {
        placeholder: 'Saisissez le chemin d\'accès au dossier des captures d\'écran',
      },
    },
    defineVideoStorage: {
      title: 'Définir le stockage vidéo',
      field: {
        placeholder: 'Saisissez le chemin d\'accès au dossier vidéo',
      },
    },
  },
  tables: {
    noData: 'Aucune donnée',
    alarmGeneral: {
      author: 'Auteur',
      connector: 'Type de connecteur',
      connectorName: 'Nom du connecteur',
      component: 'Composant',
      resource: 'Ressource',
      output: 'Message',
      lastUpdateDate: 'Date de dernière modification',
      creationDate: 'Date de création',
      duration: 'Durée',
      state: 'Criticité',
      status: 'Statut',
      extraDetails: 'Détails supplémentaires',
    },
    alarmStatus: {
      [ENTITIES_STATUSES.off]: 'Fermée',
      [ENTITIES_STATUSES.ongoing]: 'En cours',
      [ENTITIES_STATUSES.flapping]: 'Bagot',
      [ENTITIES_STATUSES.stealthy]: 'Furtive',
      [ENTITIES_STATUSES.cancelled]: 'Annulée',
    },
    alarmStates: {
      [ENTITIES_STATES.ok]: 'Info',
      [ENTITIES_STATES.minor]: 'Mineur',
      [ENTITIES_STATES.major]: 'Majeur',
      [ENTITIES_STATES.critical]: 'Critique',
    },
    contextEntities: {
      columns: {
        name: 'Nom',
        type: 'Type',
        _id: 'Id',
      },
    },
    noColumns: {
      message: 'Veuillez sélectionner au moins 1 colonne',
    },
    broadcastMessages: {
      statuses: {
        [BROADCAST_MESSAGES_STATUSES.active]: 'Actif',
        [BROADCAST_MESSAGES_STATUSES.pending]: 'En attente',
        [BROADCAST_MESSAGES_STATUSES.expired]: 'Expiré',
      },
    },
  },
  rRule: {
    advancedHint: 'Séparer les nombres par une virgule',
    textLabel: 'Récurrence',
    stringLabel: 'Résumé',
    tabs: {
      simple: 'Simple',
      advanced: 'Avancé',
    },
    errors: {
      main: 'La récurrence choisie n\'est pas valide. Nous vous recommandons de la modifier avant de sauvegarder',
    },
    periodsRanges: {
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek]: 'Cette semaine',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek]: 'Prochaine semaine',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks]: 'Prochaines 2 semaines',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth]: 'Ce mois',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth]: 'Le mois prochain',
    },
    fields: {
      freq: 'Fréquence',
      until: 'Jusqu\'à',
      byweekday: 'Par jour de la semaine',
      count: 'Répéter',
      interval: 'Intervalle',
      wkst: 'Semaine de début',
      bymonth: 'Par mois',
      bysetpos: {
        label: 'Par position',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, positifs ou négatifs. Chaque entier correspondra à la ènième occurence de la règle dans l\'intervalle de fréquence. Par exemple, une \'bysetpos\' de -1 combinée à une fréquence mensuelle, et une \'byweekday\' de (lundi, mardi, mercredi, jeudi, vendredi), va nous donner le dernier jour travaillé de chaque mois',
      },
      bymonthday: {
        label: 'Par jour du mois',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours du mois auxquels s\'appliquera la récurrence.',
      },
      byyearday: {
        label: 'Par jour de l\'année',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours de l\'année auxquels  s\'appliquera la récurrence.',
      },
      byweekno: {
        label: 'Par semaine n°',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux numéros de semaine auxquelles s\'appliquera la récurrence. Les numéros de semaines sont ceux de ISO8601, la première semaine de l\'année étant celle contenant au moins 4 jours de cette année.',
      },
      byhour: {
        label: 'Par heure',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux heures auxquelles s\'appliquera la récurrence.',
      },
      byminute: {
        label: 'Par minute',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux minutes auxquelles s\'appliquera la récurrence.',
      },
      bysecond: {
        label: 'Par seconde',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux secondes auxquelles s\'appliquera la récurrence.',
      },
    },
  },
  errors: {
    default: 'Une erreur s\'est produite...',
    lineNotEmpty: 'Cette ligne n\'est pas vide',
    JSONNotValid: 'JSON non valide',
    versionNotFound: 'Erreur dans la récupération du numéro de version...',
    statsRequestProblem: 'Erreur dans la récupération des statistiques',
    statsWrongEditionError: "Les widgets de statistiques ne sont pas disponibles dans l'édition 'core' de Canopsis",
  },
  calendar: {
    today: 'Aujourd\'hui',
    month: 'Mois',
    week: 'Semaine',
    day: 'Jour',
    pbehaviorPlanningLegend: {
      title: 'Légende',
      noData: 'Il n\'y a pas de dates d\'exception sur le calendrier',
    },
  },
  success: {
    default: 'Action effectuée avec succès',
    createEntity: 'Entité créée avec succès',
    editEntity: 'Entité éditée avec succès',
    pathCopied: 'Chemin copié dans le presse-papier',
    linkCopied: 'Lien copié dans le presse-papiers',
    authKeyCopied: 'Clé d\'authentification copiée dans le presse-papiers',
    widgetIdCopied: 'Widget id copié dans le presse-papier',
  },
  filterEditor: {
    title: 'Éditeur de filtre',
    tabs: {
      visualEditor: 'Éditeur visuel',
      advancedEditor: 'Éditeur avancé',
      results: 'Résultats',
    },
    buttons: {
      addRule: 'Ajouter une règle',
      addGroup: 'Ajouter un groupe',
      deleteGroup: 'Supprimer un groupe',
    },
    hints: {
      alarm: {
        service: 'Service',
        connector: 'Type de connecteur',
        connectorName: 'Nom du connecteur',
        component: 'Composant',
        resource: 'Ressource',
      },
    },
    errors: {
      cantParseToVisualEditor: 'Nous ne pouvons pas analyser ce filtre dans l\'éditeur visuel',
      invalidJSON: 'JSON non valide',
      required: 'Merci d\'ajouter au moins une règle valide',
    },
  },
  filterSelector: {
    defaultFilter: 'Filtre par défaut',
    fields: {
      mixFilters: 'Mix de filtres',
    },
    buttons: {
      list: 'Gérer les filtres',
    },
  },
  validator: {
    unique: 'Le champ doit être unique',
  },
  stats: {
    types: {
      [STATS_TYPES.alarmsCreated.value]: 'Alarmes créées',
      [STATS_TYPES.alarmsResolved.value]: 'Alarmes résolues',
      [STATS_TYPES.alarmsCanceled.value]: 'Alarmes annulées',
      [STATS_TYPES.alarmsAcknowledged.value]: 'Alarmes acquittées',
      [STATS_TYPES.ackTimeSla.value]: 'Taux d\'Ack conforme Sla',
      [STATS_TYPES.resolveTimeSla.value]: 'Taux de résolution conforme Sla',
      [STATS_TYPES.timeInState.value]: 'Proportion du temps dans la сriticité',
      [STATS_TYPES.stateRate.value]: 'Taux à cette сriticité',
      [STATS_TYPES.mtbf.value]: 'Temps moyen entre pannes',
      [STATS_TYPES.currentState.value]: 'Criticité courante',
      [STATS_TYPES.ongoingAlarms.value]: 'Nombre d\'alarmes en cours pendant la période',
      [STATS_TYPES.currentOngoingAlarms.value]: 'Nombre d\'alarmes actuellement en cours',
      [STATS_TYPES.currentOngoingAlarmsWithAck.value]: 'Nombre d\'alarmes acquittées actuellement en cours',
      [STATS_TYPES.currentOngoingAlarmsWithoutAck.value]: 'Nombre d\'alarmes non acquittées actuellement en cours',
    },
  },
  eventFilter: {
    title: 'Filtre d\'événements',
    externalDatas: 'Données externes',
    actionsRequired: 'Veuillez ajouter au moins une action',
    idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
  },
  metaAlarmRule: {
    outputTemplate: 'Modèle de sortie',
    thresholdType: 'Type de seuil',
    thresholdRate: 'Taux de déclenchement',
    thresholdCount: 'Seuil de déclenchement',
    timeInterval: 'Intervalle de temps',
    valuePath: 'Chemin de valeur | Chemins de valeur',
    autoResolve: 'Résolution automatique',
    idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
    corelId: 'Corel ID',
    corelIdHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>' +
      '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelStatus: 'Corel statut',
    corelStatusHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>' +
      '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelParent: 'Corel parent',
    corelChild: 'Corel enfant',
    outputTemplateHelp: '<p>Les variables accessibles sont:</p>\n' +
      '<p><strong>.Count</strong>: Le nombre d\'alarmes conséquences attachées à la méta alarme.</p>' +
      '<p><strong>.Children</strong>: L\'ensemble des variables de la dernière alarme conséquence attachée à la méta alarme.</p>' +
      '<p><strong>.Rule</strong>: Les informations administratives de la méta alarme en elle-même.</p>' +
      '<p>Quelques exemples:</p>' +
      '<p><strong>{{ .Count }} conséquences;</strong> Message de la dernière alarme conséquence : <strong>{{ .Children.Alarm.Value.State.Message }};</strong> Règle : <strong>{{ .Rule.Name }};</strong></p>' +
      '<p>Un message informatif statique</p>' +
      '<p>Corrélé par la règle <strong>{{ .Rule.Name }}</strong></p>',
    errors: {
      noValuePaths: 'Vous devez ajouter au moins 1 chemin de valeur',
    },
  },
  snmpRules: {
    title: 'Règles SNMP',
    uploadMib: 'Envoyer un fichier MIB',
    addSnmpRule: 'Ajouter une règle SNMP',
  },
  layout: {
    sideBar: {
      buttons: {
        edit: 'Activer/Désactiver le mode d\'édition',
        create: 'Créer une vue',
        settings: 'Paramètres',
      },
      activeSessions: 'Sessions actives',
      ordering: {
        popups: {
          success: 'Les groupes ont été réorganisés',
          error: 'Plusieurs groupes n\'ont pas été réorganisés',
          periodicRefreshWasPaused: 'Le rafraîchissement périodique est mis en pause pendant l\'édition du menu',
          periodicRefreshWasResumed: 'Reprise du rafraîchissement périodique',
        },
      },
    },
  },
  parameters: {
    tabs: {
      parameters: 'Paramètres',
      importExportViews: 'Importation/Exportation',
      stateSettings: 'Paramètres d\'état',
      storageSettings: 'Paramètres de stockage',
    },
    interfaceLanguage: 'Langue de l\'interface',
    groupsNavigationType: {
      title: 'Type d\'affichage de la barre de vues',
      items: {
        sideBar: 'Barre latérale',
        topBar: 'Barre d\'entête',
      },
    },
    userInterfaceForm: {
      title: 'Interface utilisateur',
      fields: {
        appTitle: 'Titre de l\'application',
        language: 'Langue par défaut',
        footer: 'Page d\'identification : pied de page',
        description: 'Page d\'identification : description',
        logo: 'Logo',
        infoPopupTimeout: 'Timeout pour les popup d\'informations',
        errorPopupTimeout: 'Timeout pour les popup d\'erreurs',
        popupTimeoutUnit: 'Unité',
        allowChangeSeverityToInfo: 'Allow change severity to info',
        maxMatchedItems: 'Articles correspondants au maximum',
      },
      tooltips: {
        maxMatchedItems: 'il doit avertir l\'utilisateur lorsque le nombre d\'éléments correspondant aux modèles est supérieur à cette valeur',
      },
    },
  },
  view: {
    errors: {
      emptyTabs: 'Merci de créer un onglet',
    },
    deleteRow: 'Supprimer la ligne',
    deleteWidget: 'Supprimer le widget',
    fullScreen: 'Plein écran',
    fullScreenShortcut: 'Alt + Entrée / Command + Entrée',
    copyWidgetId: 'Copier l\'ID du widget',
    autoHeightButton: 'Si ce bouton est sélectionné, la hauteur sera calculée automatiquement.',
  },
  patternsList: {
    noData: 'Aucun pattern. Cliquez sur \'Ajouter\' pour ajouter des champs au pattern',
    noDataDisabled: 'Aucun pattern.',
  },
  validation: {
    messages: {
      _default: "Le champ n'est pas valide",
      after: 'Le champ doit être postérieur à {1}',
      after_with_inclusion: 'Le champ doit être postérieur ou égal à {1}',
      alpha: 'Le champ ne peut contenir que des lettres',
      alpha_dash: 'Le champ ne peut contenir que des caractères alpha-numériques, tirets ou soulignés',
      alpha_num: 'Le champ ne peut contenir que des caractères alpha-numériques',
      alpha_spaces: 'Le champ ne peut contenir que des lettres ou des espaces',
      before: 'Le champ doit être antérieur à {1}',
      before_with_inclusion: 'Le champ doit être antérieur ou égal à {1}',
      between: 'Le champ doit être compris entre {1} et {2}',
      confirmed: 'Le champ ne correspond pas à {1}',
      credit_card: 'Le champ est invalide',
      date_between: 'Le champ doit être situé entre {1} et {2}',
      date_format: 'Le champ doit être au format {1}',
      decimal: 'Le champ doit être un nombre et peut contenir {1} décimales',
      digits: 'Le champ doit être un nombre entier de {1} chiffres',
      dimensions: 'Le champ doit avoir une taille de {1} pixels par {2} pixels',
      email: 'Le champ doit être une adresse e-mail valide',
      excluded: 'Le champ doit être une valeur valide',
      ext: 'Le champ doit être un fichier valide',
      image: 'Le champ doit être une image',
      included: 'Le champ doit être une valeur valide',
      integer: 'Le champ doit être un entier',
      ip: 'Le champ doit être une adresse IP',
      ip_or_fqdn: 'The field must be a valid ip address or FQDN',
      length: 'Le champ doit contenir {1} caractères',
      max: 'Le champ ne peut pas contenir plus de {length} caractères',
      max_value: 'Le champ doit avoir une valeur de {max} ou moins',
      mimes: 'Le champ doit avoir un type MIME valide',
      min: 'Le champ doit contenir au minimum {1} caractères',
      min_value: 'Le champ doit avoir une valeur de {1} ou plus',
      numeric: 'Le champ ne peut contenir que des chiffres',
      regex: 'Le champ est invalide',
      required: 'Le champ est obligatoire',
      required_if: 'Le champ est obligatoire lorsque {1} possède cette valeur',
      size: 'Le champ doit avoir un poids inférieur à {1}KB',
      url: 'Le champ n\'est pas une URL valide',
    },
    custom: {
      tstop: {
        after: 'La date de fin doit être postérieure à {1}',
      },
      logo: {
        size: 'La taille {0} doit être inférieure à {1} Ko.',
      },
    },
  },
  home: {
    popups: {
      info: {
        noAccessToDefaultView: 'Accès refusé à la vue par défaut. Redirection vers la vue par défaut de votre rôle.',
        notSelectedRoleDefaultView: 'Pas de vue par défaut sélectionnée pour votre rôle.',
        noAccessToRoleDefaultView: 'Accès refusé à la vue par défaut de votre rôle.',
      },
    },
  },
  serviceWeather: {
    seeAlarms: 'Voir les alarmes',
  },
  heartbeat: {
    title: 'Heartbeats',
    table: {
      fields: {
        id: 'ID',
        expectedInterval: 'Interval',
      },
    },
  },
  contextGeneralTable: {
    addSelection: 'Ajouter une sélection',
  },
  liveReporting: {
    button: 'Définir un intervalle de dates',
  },
  tours: {
    [TOURS.alarmsExpandPanel]: {
      step1: 'Détails',
      step2: 'Onglet plus d\'infos (N\'apparaît que s\'il existe une configuration)',
      step3: 'Onglet chronologie',
    },
  },
  handlebars: {
    requestHelper: {
      errors: {
        timeout: 'Délai d\'attente dépassé',
        unauthorized: 'Accès non autorisé',
        other: 'Erreur de récupération des données',
      },
    },
  },
  importExportViews: {
    selectAll: 'Sélectionnez tous les groupes et vues',
  },
  playlist: {
    player: {
      tooltips: {
        fullscreen: 'Les actions sont désactivées en mode plein écran',
      },
    },
  },

  permissions: {
    technical: {
      admin: 'Droits d\'administration',
      exploitation: 'Droits d\'exploitation',
    },
    business: {
      [USER_PERMISSIONS_PREFIXES.business.common]: 'Droits communs',
      [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Droits pour le widget Bac à alarmes',
      [USER_PERMISSIONS_PREFIXES.business.context]: 'Droits pour le widget Explorateur de contexte',
      [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Droits pour le widget Météo des services',
      [USER_PERMISSIONS_PREFIXES.business.counter]: 'Droits pour le widget Compteur',
      [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Droits pour le widget Météo des testing',
    },
  },

  pbehavior: {
    buttons: {
      addFilter: 'Ajouter un filtre',
      editFilter: 'Modifier le filtre',
      addRRule: 'Ajouter une règle de récurrence',
      editRrule: 'Modifier la règle de récurrence',
    },
    alerts: {
      countOverLimit: 'Le filtre que vous avez défini cible {count} entités. Cela peut affecter les performances, en êtes-vous sûr?',
      countRequestError: 'Le calcul du nombre d\'entités ciblées par le filtre s\'est terminée avec une erreur. Il se peut que ce nombre dépasse la limite conseillée et que cela affecte les performances, êtes-vous sûr?',
    },
  },

  pbehaviorExceptions: {
    title: 'Dates d\'exception',
    create: 'Ajouter une date d\'exception',
    choose: 'Sélectionnez la liste d\'exclusion',
    emptyExceptions: 'Aucune exception ajoutée pour le moment',
  },

  pbehaviorTypes: {
    usingType: 'Le type ne peut être supprimé, car il est en cours d\'utilisation',
    defaultType: 'Le type utilise la valeur par défaut, car il ne peut pas être modifié',
  },

  pbehaviorReasons: {
    usingReason: 'La raison ne peut pas être supprimée, car elle est en cours d\'utilisation',
  },

  planning: {
    tabs: {
      type: 'Type',
      reason: 'Raison',
      exceptions: 'Dates d\'exception',
    },
  },

  engines: {
    [ENGINES_NAMES.event]: {
      title: 'Event',
      description: 'Provient de la ressource',
    },

    [ENGINES_NAMES.webhook]: {
      title: 'Webhook',
      description: 'Gère les webhooks',
    },
    [ENGINES_NAMES.fifo]: {
      title: 'FIFO',
      description: 'Gère la file d\'attente des événements et des alarmes',
    },
    [ENGINES_NAMES.axe]: {
      title: 'AXE',
      description: 'Crée des alarmes et effectue des actions avec elles',
    },
    [ENGINES_NAMES.che]: {
      title: 'CHE',
      description: 'Applique les filtres d\'événements et les entités créées',
    },
    [ENGINES_NAMES.pbehavior]: {
      title: 'Pbehavior',
      description: 'Vérifie si l\'alarme est sous PBehavior',
    },
    [ENGINES_NAMES.action]: {
      title: 'Action',
      description: 'Déclenche le lancement des actions',
    },
    [ENGINES_NAMES.service]: {
      title: 'Service',
      description: 'Met à jour les compteurs et génère service-events',
    },
    [ENGINES_NAMES.dynamicInfo]: {
      title: 'Dynamic infos',
      description: 'Ajoute des informations dynamiques à l\'alarme',
    },
    [ENGINES_NAMES.correlation]: {
      title: 'Correlation',
      description: 'Gère la corrélation',
    },
    [ENGINES_NAMES.heartbeat]: {
      title: 'Heartbeat',
      description: 'Génère une alarme si un type d\'évènement ne se produit plus',
    },
  },

  remediation: {
    tabs: {
      instructions: 'Consignes',
      configurations: 'Configurations',
      jobs: 'Jobs',
    },
  },

  remediationInstructions: {
    usingInstruction: 'Ne peut pas être supprimée, car en cours d\'utilisation',
    addStep: 'Ajouter une étape',
    addOperation: 'Ajouter une opération',
    addEndpoint: 'Ajouter un point de terminaison',
    endpoint: 'Point de terminaison',
    endpointAvatar: 'EP',
    workflow: 'Si cette étape échoue :',
    remainingStep: 'Continuer avec les étapes restantes',
    timeToComplete: 'Temps d\'exécution (estimation)',
    emptySteps: 'Aucune étape ajoutée pour le moment',
    emptyOperations: 'Aucune opération ajoutée pour le moment',
    tooltips: {
      endpoint: 'Le point de terminaison doit être une question qui appelle une réponse Oui / Non',
    },
    table: {
      rating: 'Évaluation',
      lastModifiedOn: 'Dernière modification le',
      averageTimeCompletion: 'Temps moyen\nd\'exécution',
      monthExecutions: '№ d\'exécutions\nce mois-ci',
      lastExecutedBy: 'Dernière exécution par',
      lastExecutedOn: 'Dernière exécution le',
    },
    errors: {
      runningInstruction: 'Les changements ne peuvent pas être enregistrés car la consigne est en cours d\'exécution. Voulez-vous stopper l\'exécution de la consigne et ainsi enregistrer les changements ?',
      operationRequired: 'Veuillez ajouter au moins une opération',
      stepRequired: 'Veuillez ajouter au moins une étape',
    },
  },

  remediationJobs: {
    addJobs: 'Ajouter {count} job | Ajouter {count} jobs',
    usingJob: 'Le job ne peut être supprimé, car il est en cours d\'utilisation',
    table: {
      configuration: 'Configuration',
      jobId: 'Job ID',
    },
  },

  remediationConfigurations: {
    usingConfiguration: 'Ne peut pas être supprimée, car en cours d\'utilisation',
    table: {
      host: 'Hôte',
    },
  },

  remediationInstructionExecute: {
    timeToComplete: '{duration} pour terminer',
    completedAt: 'Terminé à {time}',
    failedAt: 'Échec à {time}',
    startedAt: 'Commencé à {time}\n(Date de lancement Canopsis)',
    closeConfirmationText: 'Souhaitez-vous reprendre cette consigne plus tard ?',
    popups: {
      success: '{instructionName} a été exécutée avec succès',
      failed: '{instructionName} a échoué. Veuillez faire remonter ce problème',
      connectionError: 'Il y a un problème de connexion. Veuillez cliquer sur le bouton d\'actualisation ou recharger la page.',
      wasPaused: 'La consigne {instructionName} sur l\'alarme {alarmName} a été interrompue à {date}. Vous pouvez la reprendre manuellement.',
    },
    jobs: {
      title: 'Jobs attribués :',
      startedAt: 'Date de déclenchement\n(par Canopsis)',
      launchedAt: 'Date de lancement\n(par l\'ordonnanceur)',
      completedAt: 'Fin de traitement\n(par l\'ordonnanceur)',
      waitAlert: 'L\'exécuteur de jobs ne répond pas, veuillez contacter votre administrateur',
      skip: 'Ignorer le job',
      await: 'Attendre',
      failedReason: 'Raison d\'échec',
      output: 'Output',
    },
  },

  remediationInstructionsFilters: {
    button: 'Créer un filtre de consignes',
    fields: {
      with: 'Avec les consignes sélectionnées',
      without: 'Sans les consignes sélectionnées',
      selectAll: 'Tout sélectionner',
      selectedInstructions: 'Consignes sélectionnées',
    },
    chip: {
      with: 'AVEC',
      without: 'SANS',
      all: 'TOUT',
    },
  },

  remediationPatterns: {
    tabs: {
      pbehaviorTypes: {
        title: 'Pbehavior types',
        fields: {
          activeOnTypes: 'Actif sur les types',
          disabledOnTypes: 'Désactivé sur les types',
        },
      },
    },
  },

  scenario: {
    title: 'Scénarios',
    headers: 'En-têtes',
    declareTicket: 'Déclarer un ticket',
    workflow: 'Workflow si cette action ne correspond pas :',
    remainingAction: 'Continuer avec les actions restantes',
    addAction: 'Ajouter une action',
    emptyActions: 'Aucune action ajoutée pour le moment',
    urlHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>' +
      '<i>Quelques exemples :</i>' +
      '<pre>"https://exampleurl.com?resource={{ .Alarm.Value.Resource }}"</pre>' +
      '<pre>"https://exampleurl.com?entity_id={{ .Entity.ID }}"</pre>' +
      '<pre>"https://exampleurl.com?children_count={{ len .Children }}"</pre>' +
      '<pre>"https://exampleurl.com?children={{ range .Children }}{{ .ID }}{{ end }}"</pre>',
    outputHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong> et <strong>.Entity</strong></p>' +
      '<i>Quelques exemples:</i>' +
      '<pre>Resource - {{ .Alarm.Value.Resource }}. Entity - {{ .Entity.ID }}.</pre>',
    payloadHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>' +
      '<i>Quelques exemples:</i>' +
      '<pre>{\n' +
      '  resource: "{{ .Alarm.Value.Resource }}",\n' +
      '  entity: "{{ .Entity.ID }}",\n' +
      '  children_count: "{{ len .Children }}",\n' +
      '  children: {{ range .Children }}{{ .ID }}{{ end }}\n' +
      '}</pre>',
    actions: {
      [SCENARIO_ACTION_TYPES.snooze]: 'Snooze',
      [SCENARIO_ACTION_TYPES.pbehavior]: 'Pbehavior',
      [SCENARIO_ACTION_TYPES.changeState]: 'Change state (Change and lock severity)',
      [SCENARIO_ACTION_TYPES.ack]: 'Acknowledge',
      [SCENARIO_ACTION_TYPES.ackremove]: 'Acknowledge remove',
      [SCENARIO_ACTION_TYPES.assocticket]: 'Associate ticket',
      [SCENARIO_ACTION_TYPES.cancel]: 'Cancel',
      [SCENARIO_ACTION_TYPES.webhook]: 'Webhook',
    },
    fields: {
      triggers: 'Triggers',
      emitTrigger: 'Émettre un trigger',
      withAuth: 'Avez-vous besoin de champs d\'authentification ?',
      emptyResponse: 'Réponse vide',
      isRegexp: 'La valeur peut être une RegExp',
      headerKey: "Clé d'en-tête",
      headerValue: "Valeur d'en-tête",
      key: 'Clé',
      skipVerify: 'Ne pas vérifier les certificats HTTPS',
    },
    tabs: {
      pattern: 'Pattern',
    },
    errors: {
      actionRequired: 'Veuillez ajouter au moins une action',
    },
  },

  mixedField: {
    types: {
      string: '@:variableTypes.string',
      number: '@:variableTypes.number',
      boolean: '@:variableTypes.boolean',
      null: '@:variableTypes.null',
      array: '@:variableTypes.array',
    },
  },

  entity: {
    types: {
      connector: 'type de connecteur',
      component: 'composant',
      resource: 'ressource',
    },
    fields: {
      type: 'Type',
      manageInfos: 'Gérer les informations',
      form: 'Formulaire',
      impact: 'Impacts',
      depends: 'Dépendances',
    },
    manageInfos: {
      title: 'Informations',
      createTitle: 'Ajouter une information',
      emptyInfos: 'Aucune information',
    },
  },

  service: {
    fields: {
      category: 'Catégorie',
      name: 'Nom',
      outputTemplate: 'Modèle de sortie',
      createCategory: 'Ajouter une catégorie',
      createCategoryHelp: 'Appuyez sur <kbd>enter</kbd> pour enregistrer',
    },
  },

  users: {
    table: {
      username: 'Identifiant utilisateur',
      firstName: 'Prénom',
      lastName: 'Nom',
      role: 'Rôle',
      enabled: 'Actif',
      auth: 'Auth',
    },
    fields: {
      username: 'Identifiant utilisateur',
      firstName: 'Prénom',
      lastName: 'Nom',
      email: 'Email',
      password: 'Mot de passe',
      role: 'Rôle',
      language: 'Langue de l\'interface par défaut',
    },
  },

  testSuite: {
    xmlFeed: 'Flux XML',
    hostname: 'Nom d\'hôte',
    lastUpdate: 'Dernière mise à jour',
    timeTaken: 'Temps pris',
    totalTests: 'Total des tests',
    disabledTests: 'Tests désactivés',
    copyMessage: 'Copier le message système',
    systemError: 'Erreur système',
    systemErrorMessage: 'Message d\'erreur système',
    systemOut: 'Système hors',
    systemOutMessage: 'Message de sortie du système',
    compareWithHistorical: 'Comparer avec les données historiques',
    className: 'Nom du cours',
    line: 'Ligne',
    failureMessage: 'Message d\'échec',
    noData: 'Aucun message système trouvé dans XML',
    tabs: {
      summary: 'Résumé',
      globalMessages: 'Messages globaux',
      gantt: 'Gantt',
      details: 'Des détails',
      screenshots: 'Captures d\'écran',
      videos: 'Vidéos',
    },
    statuses: {
      [TEST_SUITE_STATUSES.passed]: 'Passé',
      [TEST_SUITE_STATUSES.skipped]: 'Ignoré',
      [TEST_SUITE_STATUSES.error]: 'Erreur',
      [TEST_SUITE_STATUSES.failed]: 'Manqué',
    },
    popups: {
      systemMessageCopied: 'Message système copié dans le presse-papiers',
    },
  },

  stateSetting: {
    worstLabel: 'The worst of:',
    worstHelpText: 'Canopsis compte l\'état pour chaque critère défini. L\'état final de la suite de tests JUnit est considéré comme le pire des états résultants.',
    criterion: 'Criterion',
    serviceState: 'État du service',
    methods: {
      [STATE_SETTING_METHODS.worst]: 'Pire',
      [STATE_SETTING_METHODS.worstOfShare]: 'Pire de part',
    },
    states: {
      minor: 'Mineur',
      major: 'Majeur',
      critical: 'Critique',
    },
  },

  storageSetting: {
    junit: {
      title: 'Stockage de données JUnit',
      deleteAfter: 'Supprimer les données des suites de tests après',
      deleteAfterHelpText: 'Lorsqu\'elles sont activées, les données des suites de tests JUnit (XML, captures d\'écran et vidéos) seront supprimées après la période définie.',
    },
    history: {
      junit: 'Script lancé à {launchedAt}',
    },
  },
}, featureService.get('i18n.fr'));
