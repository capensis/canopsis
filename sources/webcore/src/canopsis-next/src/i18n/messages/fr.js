import { ENTITIES_STATES, ENTITIES_STATUSES, STATS_TYPES, STATS_CRITICITY } from '@/constants';

export default {
  common: {
    undefined: 'Non définie',
    entity: 'Entitée',
    watcher: 'Observateur',
    pbehaviors: 'Comportements périodiques',
    widget: 'Widget',
    addWidget: 'Ajouter un widget',
    refresh: 'Rafraîchir',
    toggleEditView: 'Activer/Désactiver le mode édition',
    name: 'Nom',
    description: 'Description',
    author: 'Auteur',
    submit: 'Envoyer',
    cancel: 'Annuler',
    options: 'Options',
    type: 'Type',
    quitEditing: 'Quitter le mode d\'édition',
    enabled: 'Activé(e)',
    disabled: 'Désactivé(e)',
    login: 'Connexion',
    yes: 'Oui',
    no: 'Non',
    default: 'Défaut',
    confirmation: 'Etes-vous sûr(e) ?',
    parameters: 'Paramètres',
    by: 'Par',
    date: 'Date',
    comment: 'Commentaire',
    end: 'Fin',
    recursive: 'Recursif',
    states: 'Etats',
    sla: 'Sla',
    authors: 'Auteurs',
    stat: 'Statistique',
    trend: 'Tendance',
    users: 'Utilisateurs',
    roles: 'Roles',
    rights: 'Droits',
    username: 'Nom d\'utilisateur',
    password: 'Mot de passe',
    connect: 'Connexion',
    optionnal: 'Optionnel',
    logout: 'Se déconnecter',
    title: 'Titre',
    save: 'Sauvegarder',
    label: 'Label',
    field: 'Champs',
    value: 'Valeur',
    add: 'Ajouter',
    create: 'Créer',
    delete: 'Supprimer',
    show: 'Afficher',
    edit: 'Éditer',
    parse: 'Compiler',
    home: 'Accueil',
    step: 'Etape',
    entries: 'entrées',
    showing: 'Affiche',
    apply: 'Appliquer',
    to: 'à',
    of: 'sur',
    tags: 'tags',
    actionsLabel: 'Actions',
    noResults: 'Pas de résultats',
    exploitation: 'Exploitation',
    administration: 'Administration',
    actions: {
      close: 'Fermer',
      acknowledge: 'Acquitter',
      acknowledgeAndReport: 'Acquitter et signaler un incident',
      saveChanges: 'Sauvegarder',
      reportIncident: 'Signaler un incident',
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
  user: {
    firstName: 'Prénom',
    lastName: 'Nom',
    role: 'Role',
    defaultView: 'Vue par défaut',
  },
  context: {
    impacts: 'Impacts',
    dependencies: 'Dépendances',
    moreInfos: {
      type: 'Type',
      lastActiveDate: 'Dernière Date d\'Activité',
    },
  },
  search: {
    advancedSearch: '<span>Aide sur la recherche avancée :</span>\n' +
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
  },
  entities: {
    watcher: 'observateurs',
    entities: 'entités',
  },
  login: {
    errors: {
      incorrectEmailOrPassword: 'Mot de passe/Email incorrect',
    },
  },
  alarmList: {
    actions: {
      titles: {
        ack: 'Ack',
        fastAck: 'Ack rapide',
        ackRemove: 'Annuler ack',
        pbehavior: 'Comportement périodique',
        snooze: 'Snooze',
        pbehaviorList: 'Lister comportements pédioriques',
        declareTicket: 'Déclarer un incident',
        associateTicket: 'Associer ticket',
        cancel: 'Annuler alarme',
        changeState: 'Changer criticité',
        moreInfos: 'Plus d\'infos',
      },
      iconsTitles: {
        ack: 'Ack',
        declareTicket: 'Déclarer un incident',
        canceled: 'Annulée',
        snooze: 'Snooze',
        pbehaviors: 'Comportement périodique',
      },
      iconsFields: {
        ticketNumber: 'Numéro de ticket',
      },
    },
  },
  weather: {
    moreInfos: 'Plus d\'infos',
  },
  pbehaviors: {
    connector: 'Connecteur',
    connectorName: 'Nom du connecteur',
    isEnabled: 'Est actif',
    begins: 'Débute',
    ends: 'Se termine',
    type: 'Type',
    reason: 'Raison',
    rrule: 'Rrule',
  },
  settings: {
    titles: {
      alarmListSettings: 'Paramètres du bac à alarmes',
      contextTableSettings: 'Paramètres de l\'explorateur de contexte',
      weatherSettings: 'Paramètres de la météo des services',
      statsHistogramSettings: 'Paramètres de l\'histogramme',
      statsCurvesSettings: 'Paramètres de courbes de stats',
      statsTableSettings: 'Paramètres du tableau de stats',
      statsCalendarSettings: 'Paramètres du calendrier',
      statsNumberSettings: 'Paramètres du compteur de stat',
    },
    advancedSettings: 'Paramètres avancés',
    widgetTitle: 'Titre du widget',
    columnName: 'Nom de la colonne',
    defaultSortColumn: 'Colonne de tri par défaut',
    columnNames: 'Nom des colonnes',
    periodicRefresh: 'Rafraichissement périodique',
    defaultNumberOfElementsPerPage: 'Nombre d\'élement/page par défaut',
    elementsPerPage: 'Élements par page',
    filterOnOpenResolved: 'Filtre sur Open/Resolved',
    open: 'Ouverte',
    resolved: 'Résolue',
    filters: 'Filtres',
    filterEditor: 'Ésditeur de filtre',
    duration: 'Durée',
    tstop: 'Date de fin',
    periodsNumber: 'Nombre d\'étapes',
    statName: 'Nom de la statistique',
    statsSelect: {
      title: 'Sélecteur de statistique',
      required: 'Veuillez sélectionner au moins une statistique',
    },
    yesNoMode: 'Mode Oui/Non',
    selectAFilter: 'Selectionner un filtre',
    criticityLevels: 'Niveaux de criticité',
    colorsSelector: {
      title: 'Sélecteur de couleur',
      statsCriticity: {
        [STATS_CRITICITY.ok]: 'ok',
        [STATS_CRITICITY.minor]: 'minor',
        [STATS_CRITICITY.major]: 'major',
        [STATS_CRITICITY.critical]: 'critical',
      },
    },
    statsNumbers: {
      title: 'Cellule de stats',
      yesNoMode: 'Mode Oui/Non',
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
    weatherTemplate: 'Template - Tuiles',
    modalTemplate: 'Template - Modal',
    entityTemplate: 'Template - Entitées',
    columnSM: 'Colonnes - Petit',
    columnMD: 'Colonnes - Moyen',
    columnLG: 'Colonnes - Large',
    height: 'Hauteur',
    margin: {
      title: 'Marges',
      top: 'Marge - Haut',
      right: 'Marge - Droite',
      bottom: 'Marge - Bas',
      left: 'Marge - Gauche',
    },
    contextTypeOfEntities: {
      title: 'Type d\'entitées',
      fields: {
        component: 'Composant',
        connector: 'Connecteur',
        resource: 'Ressource',
        watcher: 'Observateur',
      },
    },
    statSelector: {
      error: {
        alreadyExist: 'Une statistique avec ce nom existe déjà.',
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
    considerPbehaviors: {
      title: 'Prendre en compte les comportements périodiques ?',
    },
    serviceWeatherModalTypes: {
      title: 'Type de modal',
      fields: {
        moreInfo: 'Plus d\'infos',
        alarmList: 'Bac à alarmes',
      },
    },
  },
  modals: {
    contextInfos: {
      title: 'Infos sur l\'entité',
    },
    createEntity: {
      createTitle: 'Créer une entitée',
      editTitle: 'Editer une entitée',
      duplicateTitle: 'Dupliquer une entitée',
      infosList: 'Infos',
      addInfos: 'Ajouter un champ info',
      noInfos: 'No infos',
      fields: {
        type: 'Type',
        manageInfos: 'Gérer Infos',
        form: 'Formulaire',
        impact: 'Impacts',
        depends: 'Dépendances',
        types: {
          connector: 'connecteur',
          component: 'composant',
          resource: 'ressource',
        },
      },
      success: {
        create: 'Entité créée avec succès !',
        edit: 'Entité editée avec succès !',
        duplicate: 'Entité dupliquée avec succès !',
      },
    },
    createWatcher: {
      createTitle: 'Créer un observateur',
      editTitle: 'Editer un observateur',
      duplicateTitle: 'Dupliquer un observateur',
      displayName: 'Nom',
      success: {
        create: 'Observateur créé avec succès !',
        edit: 'Observateur edité avec succès !',
        duplicate: 'Observateur dupliqué avec succès !',
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
        title: 'Dupliquer une vue',
        infoMessage: 'Vous êtes en train de dupliquer une vue. Toutes les lignes et les widgets de la vue dupliquée seront copiés dans la nouvelle vue.',
      },
      noData: 'Aucun groupe correspondant. Appuyez sur <kbd>enter</kbd> pour en créer un nouveau.',
      fields: {
        groupIds: 'Choisissez une groupe, ou créez-en un nouveau',
        groupTags: 'Labels de groupes',
      },
      success: {
        create: 'Nouvelle vue créée !',
        edit: 'Vue éditée avec succès !',
      },
      fail: {
        create: 'Erreur dans la création de la vue...',
        edit: 'Erreur dans l\'édition de la vue...',
      },
    },
    createAckEvent: {
      title: 'Acquitter',
      tooltips: {
        ackResources: 'Voulez-vous acquitter les ressources liées ?',
      },
      fields: {
        ticket: 'Numéro du ticket',
        output: 'Note',
        ackResources: 'Ack ressources',
      },
    },
    createSnoozeEvent: {
      title: 'Snooze',
      fields: {
        duration: 'Durée',
      },
    },
    createCancelEvent: {
      title: 'Annuler',
      fields: {
        output: 'Note',
      },
    },
    createChangeStateEvent: {
      title: 'Changer l\'état',
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
      title: 'Ajouter un comportement périodique à ces éléments ?',
      fields: {
        name: 'Nom',
        start: 'Début',
        stop: 'Fin',
        reason: 'Raison',
        type: 'Type',
        rRuleQuestion: 'Ajouter une rrule à ce comportement périodique ?',
      },
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
    },
    liveReporting: {
      editLiveReporting: 'Suivi personnalisé',
      dateInterval: 'Interval de dates',
      today: 'Aujourd\'hui',
      yesterday: 'Hier',
      last7Days: '7 derniers jours',
      last30Days: '30 derniers jours',
      thisMonth: 'Ce mois',
      lastMonth: 'Mois dernier',
      custom: 'Personnalisé',
      tstart: 'Démarre',
      tstop: 'Finis',
    },
    moreInfos: {
      moreInfos: 'Plus d\'infos',
      defineATemplate: 'Pour définir le template de cette fenêtre, rendez-vous dans les paramètres du bac à alarmes.',
    },
    watcher: {
      criticity: 'Criticity',
      organization: 'Organization',
      numberOk: 'Nombre Ok',
      numberKo: 'Nombre Ko',
      state: 'State',
      name: 'Nom',
      org: 'Org',
      noData: 'Pas de données',
      ticketing: 'Ticketing',
      application_crit_label: 'Criticité',
      product_line: 'Ligne produit',
      service_period: 'Plage de surveillance',
      isInCarat: 'Cartographic repository',
      application_label: 'Description',
      target_platform: 'Environnement',
      scenario_label: 'Label',
      scenario_probe_name: 'Sonde',
      scenario_calendar: 'Intervals d\'éxécution',
    },
    filter: {
      create: {
        title: 'Créer un filtre',
      },
      edit: {
        title: 'Editer',
      },
      fields: {
        title: 'Nom',
      },
    },
    colorPicker: {
      title: 'Sélecteur de couleur',
    },
    textEditor: {
      title: 'Éditeur de texte',
    },
    widgetCreation: {
      title: 'Sélectionnez un widget',
      types: {
        alarmList: {
          title: 'Bac à alarmes',
        },
        context: {
          title: 'Explorateur de contexte',
        },
        weather: {
          title: 'Météo de services',
        },
        statsHistogram: {
          title: 'Histogramme de statistiques',
        },
        statsCurves: {
          title: 'Courbes de statistiques',
        },
        statsTable: {
          title: 'Tableau de statistiques',
        },
        statsCalendar: {
          title: 'Calendrier',
        },
        statsNumber: {
          title: 'Compteur de statistiques',
        },
      },
    },
    manageHistogramGroups: {
      title: {
        add: 'Ajouter un groupe',
        edit: 'Editer un groupe',
      },
    },
    addStat: {
      title: {
        add: 'Ajouter une statistique',
        edit: 'Editer une statistique',
      },
    },
    group: {
      create: {
        title: 'Créer un groupe',
      },
      edit: {
        title: 'Editer un groupe',
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
      title: 'Créer un utilisateur',
      fields: {
        username: 'Nom d\'utilisateur',
        firstName: 'Prénom',
        lastName: 'Nom',
        email: 'Email',
        password: 'Mot de passe',
        language: 'Langue de l\'interface par défaut',
        enabled: 'Actif',
      },
    },
    editUser: {
      title: 'Editer un utilisateur',
    },
    createRole: {
      title: 'Créer un role',
    },
    editRole: {
      title: 'Editer un role',
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
        title: 'Editer une règle',
        success: 'Règle éditée avec succès !',
      },
      priority: 'Priorité',
      editPattern: 'Editer le pattern',
      advanced: 'Avancée',
      addAField: 'Ajouter un champs',
      simpleEditor: 'Editeur simple',
      advancedEditor: 'Editeur avancé',
      comparisonRules: 'Règles de comparaison',
      enrichmentOptions: 'Options d\'enrichissement',
      editActions: 'Editer les actions',
      addAction: 'Ajouter une action',
      actions: 'Actions',
      from: 'Depuis',
      to: 'Vers',
      externalData: 'Données externes',
      onSuccess: 'En cas de succès',
      onFailure: 'En cas d\'échec',
    },
    viewTab: {
      create: {
        title: 'Ajouter un onglet',
      },
      edit: {
        title: 'Editer l\'onglet',
      },
      fields: {
        title: 'Titre',
      },
    },
  },
  tables: {
    noData: 'Aucune donnée',
    contextList: {
      title: 'Liste Context',
      name: 'Nom',
      type: 'Type',
      id: 'Id',
    },
    alarmGeneral: {
      title: 'Generale',
      author: 'Auteur',
      connector: 'Connecteur',
      connectorName: 'Nom du connecteur',
      component: 'Composant',
      resource: 'Ressource',
      output: 'Output',
      lastUpdateDate: 'Date de dernière modification',
      creationDate: 'Date de création',
      duration: 'Durée',
      state: 'État',
      status: 'Status',
      extraDetails: 'Détails supplémentaires',
    },
    /**
     * This object for pbehavior fields from database
     */
    pbehaviorList: {
      name: 'Nom',
      author: 'Auteur',
      connector: 'Connecteur',
      connectorName: 'Nom du connecteur',
      enabled: 'Actif',
      tstart: 'Démarre',
      tstop: 'Finis',
      type_: 'Type',
      reason: 'Raison',
      rrule: 'Rrule',
    },
    rolesList: {
      name: 'Nom',
      actions: 'Actions',
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
      [ENTITIES_STATES.minor]: 'Minor',
      [ENTITIES_STATES.major]: 'Major',
      [ENTITIES_STATES.critical]: 'Critical',
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
    admin: {
      users: {
        columns: {
          username: 'Nom d\'utilisateur',
          role: 'Role',
          enabled: 'Actif',
        },
      },
    },
  },
  rRule: {
    advancedHint: 'Séparer les nombres par une virgule',
    textLabel: 'Rrule',
    stringLabel: 'Résumé',
    tabs: {
      simple: 'Simple',
      advanced: 'Avancé',
    },
    errors: {
      main: 'La Rrule choisis n\'est pas valide. Nous vous recommandons de la modifier avant de sauvegarder',
    },
    fields: {
      freq: 'Fréquence',
      until: 'Jusqu\'à',
      byweekday: 'Par jour de la semaine',
      count: 'Répéter',
      interval: 'Interval',
      wkst: 'Semaine de début',
      bymonth: 'Par mois',
      bysetpos: {
        label: 'Par position',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, positifs ou négatifs. Chaque entier correspondra à la n-ième occurence de la règle dans l\'interval de fréquence. Par exemple, une \'bysetpos\' de -1 combinée à une fréquence mensuelle, et une \'byweekday\' de (Lundi, Mardi, Mercredi, Jeudi, Vendredi), va nous donner le dernier jour travaillé de chaque mois',
      },
      bymonthday: {
        label: 'Par jour du mois',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours du mois auquel appliquer la récurrence.',
      },
      byyearday: {
        label: 'Par jour de l\'année',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours de l\'année auquel appliquer la récurrence.',
      },
      byweekno: {
        label: 'Par semaine n°',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux numéros de semaine a laquelle appliquer la récurrence. Les numéros de semaines sont ceux de ISO8601, la première semaine de l\'année étant celle contenant au moins 4 jours de cette année',
      },
      byhour: {
        label: 'Par heure',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux heures auquel appliquer la récurrence.',
      },
      byminute: {
        label: 'Par minute',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux minutes auquel appliquer la récurrence.',
      },
      bysecond: {
        label: 'Par seconde',
        tooltip: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux secondes auquel appliquer la récurrence.',
      },
    },
  },
  errors: {
    default: 'Une erreur s\'est produite...',
    lineNotEmpty: 'Cette ligne n\'est pas vide',
    JSONNotValid: 'JSON non valide..',
  },
  calendar: {
    today: 'Aujourd\'hui',
    month: 'Mois',
    week: 'Semaine',
    day: 'Jour',
  },
  success: {
    default: 'Action effectuée avec succès',
    createEntity: 'Entité créée avec succès',
    editEntity: 'Entité éditée avec succès',
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
    resultsTableHeaders: {
      alarm: {
        connector: 'Connecteur',
        connectorName: 'Nom du connecteur',
        component: 'Composant',
        resource: 'Ressource',
      },
      entity: {
        id: 'ID',
        name: 'Nom',
        type: 'Type',
      },
      errors: {
        invalidJSON: 'JSON non valide',
      },
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
      [STATS_TYPES.ackTimeSla.value]: 'Taux d\'Ack conforme Sla',
      [STATS_TYPES.resolveTimeSla.value]: 'Taux de résolution conforme Sla',
      [STATS_TYPES.timeInState.value]: 'Proportion du temps dans l\'état',
      [STATS_TYPES.stateRate.value]: 'Taux à cet état',
      [STATS_TYPES.mtbf.value]: 'Temps moyen entre pannes',
      [STATS_TYPES.currentState.value]: 'Etat courant',
      [STATS_TYPES.ongoingAlarms.value]: 'Nombre d\'alarmes en cours pendant la période',
      [STATS_TYPES.currentOngoingAlarms.value]: 'Nombre d\'alarmes actuellement en cours',
    },
  },
  eventFilter: {
    title: 'Filtre d\'événements',
    type: 'Type',
    pattern: 'Pattern',
    priority: 'Priorité',
    enabled: 'Activé',
    actions: 'Actions',
    actionsRequired: 'Veuillez ajouter au moins une action',
  },
  layout: {
    sideBar: {
      buttons: {
        edit: 'Activer/Désactiver le mode d\'édition',
        create: 'Créer une vue',
        settings: 'Paramètres',
      },
    },
  },
  parameters: {
    interfaceLanguage: 'Langage de l\'interface',
  },
};
