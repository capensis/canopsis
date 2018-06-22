import { ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';

export default {
  common: {
    watcher: 'observateur',
    name: 'Nom',
    description: 'Description',
    yes: 'Oui',
    no: 'Non',
    confirmation: 'Etes-vous sûr(e) ?',
    submit: 'Soumettre',
    enabled: 'Activé',
    login: 'Connexion',
    username: 'Nom d\'utilisateur',
    password: 'Mot de passe',
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
  alarmList: {
    actions: {
      ack: 'ack',
      fastAck: 'Ack rapide',
      ackRemove: 'Annuler ack',
      pbehavior: 'Comportement périodique',
      snooze: 'Snooze',
      pbehaviorList: 'Lister comportements pédioriques',
      declareTicket: 'Signaler ticket',
      associateTicket: 'Associer ticket',
      cancel: 'Annuler alarme',
      changeState: 'Changer criticité',
      moreInfos: 'Plus d\'infos',
    },
  },
  alarmListSettings: {
    alarmListSettings: 'Paramètres du bac à alarmes',
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
    selectAFilter: 'Sélectionner un filtre',
    infoPopup: 'Info popup',
    moreInfosModal: 'Fenêtre "Plus d\'infos"',
  },
  modals: {
    createEntity: {
      title: 'Créer une entitée',
      fields: {
        type: 'Types',
        manageInfos: 'Gérer Infos',
        types: {
          connector: 'connecteur',
          component: 'composant',
          resource: 'ressource',
        },
      },
    },
    moreInfos: {
      moreInfos: 'Plus d\'infos',
      defineATemplate: 'Pour définir le template de cette fenêtre, rendez-vous dans les paramètres du bac à alarmes.',
    },
    createAckEvent: {
      title: 'Ajouter un événement: ack',
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
      title: 'Ajouter un événement: snooze',
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
      title: 'Ajouter un événement: Changer état',
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
      title: 'Ajouter un événement: ackremove',
    },
    createDeclareTicket: {
      title: 'Ajouter un événement: declareticket',
    },
    createAssociateTicket: {
      title: 'Ajouter un événement: associateticket',
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
    },
  },
  tables: {
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
    contextEntities: {
      columns: {
        name: 'Nom',
        type: 'Type',
        _id: 'Id',
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
  },
  mFilterEditor: {
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
};
