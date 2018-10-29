import { ENTITIES_STATES, ENTITIES_STATUSES, STATS_TYPES } from '@/constants';

export default {
  common: {
    watcher: 'Observateur',
    widget: 'Widget',
    name: 'Nom',
    description: 'Description',
    author: 'Auteur',
    yes: 'Oui',
    no: 'Non',
    confirmation: 'Etes-vous sûr(e) ?',
    submit: 'Soumettre',
    enabled: 'Activé',
    disabled: 'Désactivée',
    login: 'Connexion',
    by: 'Par',
    date: 'Date',
    comment: 'Commentaire',
    end: 'Fin',
    recursive: 'Recursif',
    states: 'Etats',
    sla: 'Sla',
    authors: 'Auteurs',
    trend: 'Tendance',
    username: 'Nom d\'utilisateur',
    password: 'Mot de passe',
    connect: 'Connexion',
    logout: 'Se déconnecter',
    title: 'Titre',
    save: 'Sauvegarder',
    label: 'Label',
    value: 'Valeur',
    add: 'Ajouter',
    delete: 'Supprimer',
    edit: 'Éditer',
    parse: 'Compiler',
    home: 'Accueil',
    entries: 'entrées',
    showing: 'Affiche',
    apply: 'Appliquer',
    to: 'à',
    of: 'sur',
    tags: 'tags',
    actionsLabel: 'Actions',
    noResults: 'Pas de résultats',
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
      statsTableSettings: 'Paramètres du tableau de stats',
    },
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
    filterEditor: 'Editeur de filtre',
    duration: 'Durée',
    tstop: 'Date de fin',
    statsSelect: 'Sélecteur de stats',
    selectAFilter: 'Sélectionner un filtre',
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
        size: {
          sm: 'Colonnes SM',
          md: 'Colonnes MD',
          lg: 'Colonnes LG',
        },
      },
    },
    moreInfosModal: 'Fenêtre "Plus d\'infos"',
    weatherTemplate: 'Template - Tuiles',
    modalTemplate: 'Template - Modal',
    entityTemplate: 'Template - Entitées',
    columnSM: 'Colonnes - Petit',
    columnMD: 'Colonnes - Moyen',
    columnLG: 'Colonnes - Large',
    contextTypeOfEntities: {
      title: 'Type d\'entitées',
      fields: {
        component: 'Composant',
        connector: 'Connecteur',
        resource: 'Ressource',
        watcher: 'Observateur',
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
      infosList: 'Infos',
      addInfos: 'Ajouter un champ info',
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
    },
    createWatcher: {
      createTitle: 'Créer un observateur',
      editTitle: 'Editer un observateur',
      displayName: 'Nom',
    },
    view: {
      title: 'Créer une vue',
      noData: 'Pas de groupe correspondant. Presser <kbd>Enter</kbd> pour en créer un nouveau',
      fields: {
        groupIds: 'Choisir un groupe, ou en créer un nouveau',
        groupTags: 'Tags du groupe',
      },
      success: 'Nouvelle vue crée',
      fail: 'Erreur dans la création de vue',
    },
    createAckEvent: {
      title: 'Ajouter un événement: Ack',
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
      title: 'Ajouter un événement: Snooze',
      fields: {
        duration: 'Durée',
      },
    },
    createCancelEvent: {
      title: 'Ajouter un événement: Annuler',
      fields: {
        output: 'Note',
      },
    },
    createChangeStateEvent: {
      title: 'Ajouter un événement: Changer l\'état',
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
      title: 'Ajoute run comportement périodique à ces éléments ?',
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
      title: 'Ajouter un événement: Annuler Ack',
    },
    createDeclareTicket: {
      title: 'Ajouter un événement: Déclarer un incident',
    },
    createAssociateTicket: {
      title: 'Ajouter un événement: Associer un numéro de ticket',
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
      },
    },
  },
  tables: {
    contextList: {
      title: 'Liste Context',
      name: 'Nom',
      id: 'Id',
      noDataText: 'Faites une recherche',
    },
    alarmGeneral: {
      title: 'Generale',
      author: 'Auteur',
      connector: 'Connecteur',
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
      connector_name: 'Nom du connecteur',
      enabled: 'Actif',
      tstart: 'Démarre',
      tstop: 'Finis',
      type_: 'Type',
      reason: 'Raison',
      rrule: 'Rrule',
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
  },
  success: {
    default: 'Action effectuée avec succès',
    createEntity: 'Entité créée avec succès',
    editEntity: 'Entité éditée avec succès',
  },
  filterEditor: {
    title: 'Editeur de filtre',
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
      connector: 'Connecteur',
      connectorName: 'Nom du connecteur',
      component: 'Composant',
      resource: 'Ressource',
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
};
