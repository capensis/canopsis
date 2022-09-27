import { merge } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  STATS_TYPES,
  STATS_CRITICITY,
  QUICK_RANGES,
  TOURS,
  BROADCAST_MESSAGES_STATUSES,
  USER_PERMISSIONS_PREFIXES,
  PBEHAVIOR_RRULE_PERIODS_RANGES,
  WIDGET_TYPES,
  ACTION_TYPES,
  ENTITY_TYPES,
  TEST_SUITE_STATUSES,
  SIDE_BARS,
  STATE_SETTING_METHODS,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
  IDLE_RULE_TYPES,
  IDLE_RULE_ALARM_CONDITIONS,
  USERS_PERMISSIONS,
  ALARMS_OPENED_VALUES,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_ENGINES_NAMES,
  GROUPS_NAVIGATION_TYPES,
  ALARM_METRIC_PARAMETERS,
  USER_METRIC_PARAMETERS,
  EVENT_FILTER_TYPES,
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  PATTERN_FIELD_TYPES,
  PBEHAVIOR_TYPE_TYPES,
  SCENARIO_TRIGGERS,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import featureService from '@/services/features';

export default merge({
  common: {
    ok: 'Ok',
    undefined: 'Non défini',
    entity: 'Entité | Entités',
    service: 'Service',
    widget: 'Widget',
    addWidget: 'Ajouter un widget',
    addTab: 'Ajouter un onglet',
    shareLink: 'Créer un lien de partage',
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
    parameter: 'Paramètre | Paramètres',
    by: 'Par',
    date: 'Date',
    comment: 'Commentaire | Commentaires',
    lastComment: 'Dernier commentaire',
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
    user: 'Utilisateur | Utilisateurs',
    role: 'Rôle | Rôles',
    import: 'Importer',
    export: 'Exporter',
    profile: 'Profil',
    username: 'Identifiant utilisateur',
    password: 'Mot de passe',
    authKey: 'Clé d\'authentification',
    widgetId: 'Identifiant du widget',
    connect: 'Connexion',
    optional: 'Optionnel',
    logout: 'Se déconnecter',
    title: 'Titre',
    save: 'Sauvegarder',
    label: 'Label',
    field: 'Champ',
    value: 'Valeur',
    limit: 'Limite',
    add: 'Ajouter',
    create: 'Créer',
    delete: 'Supprimer',
    show: 'Afficher',
    edit: 'Éditer',
    duplicate: 'Dupliquer',
    play: 'Lecture',
    copyLink: 'Copier le lien',
    parse: 'Analyser',
    home: 'Accueil',
    step: 'Étape',
    paginationItems: 'Affiche {first} à {last} sur {total} Entrées',
    apply: 'Appliquer',
    from: 'Depuis',
    to: 'Vers',
    tags: 'Étiquettes',
    actionsLabel: 'Actions',
    noResults: 'Pas de résultats',
    result: 'Résultat',
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
    icon: 'Icône',
    fullscreen: 'Plein écran',
    interval: 'Période',
    status: 'Statut',
    unit: 'Unité',
    delay: 'Intervalle',
    begin: 'Début',
    timezone: 'Fuseau horaire',
    reason: 'Raison',
    or: 'OU',
    and: 'ET',
    priority: 'Priorité',
    clear: 'Nettoyer',
    deleteAll: 'Tout supprimer',
    payload: 'Payload',
    note: 'Note',
    output: 'Output',
    displayName: 'Afficher un nom',
    created: 'Date de création',
    updated: 'Date de dernière modification',
    expired: 'Date d\'expiration',
    lastEventDate: 'Date du dernier événement',
    pattern: 'Modèle | Modèles',
    correlation: 'Corrélation',
    periods: 'Périodes',
    range: 'Plage',
    duration: 'Durée',
    previous: 'Précédent',
    next: 'Suivant',
    eventPatterns: 'Modèles des événements',
    alarmPatterns: 'Modèles des alarmes',
    entityPatterns: 'Modèles des entités',
    pbehaviorPatterns: 'Modèles de comportement',
    totalEntityPatterns: 'Total des modèles d\'entité',
    serviceWeatherPatterns: 'Modèles météorologiques de service',
    addFilter: 'Ajouter un filtre',
    id: 'Identifiant',
    reset: 'Réinitialiser',
    selectColor: 'Sélectionner la couleur',
    triggers: 'Déclencheurs',
    disableDuringPeriods: 'Désactiver pendant les pauses',
    retryDelay: 'Intervalle de tentatives',
    retryUnit: 'Unité d\'essai',
    retryCount: 'Nombre de tentatives après échec',
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
    error: 'Erreur | Erreurs',
    failures: 'Échecs',
    skipped: 'Ignoré',
    current: 'Actuel',
    average: 'Moyenne',
    information: 'Information | Informations',
    file: 'Fichier',
    group: 'Groupe | Groupes',
    view: 'Vue | Vues',
    tab: 'Onglet | Onglets',
    access: 'Accès',
    communication: 'Communication | Communications',
    general: 'Général',
    notification: 'Notification | Notifications',
    dismiss: 'Rejeter',
    approve: 'Approuver',
    summary: 'Résumé',
    recurrence: 'Récurrence',
    statistics: 'Statistiques',
    action: 'Action | Actions',
    minimal: 'Minimal',
    optimal: 'Optimal',
    graph: 'Graphique | Graphiques',
    systemStatus: 'État du système',
    downloadAsPng: 'Télécharger en PNG',
    rating: 'Notation | Notations',
    sampling: 'Échantillonnage',
    parametersToDisplay: '{count} paramètres à afficher',
    uptime: 'Uptime',
    maintenance: 'Maintenance',
    downtime: 'Downtime',
    toTheTop: 'Jusqu\'au sommet',
    time: 'Temps',
    lastModifiedOn: 'Dernière modification le',
    lastModifiedBy: 'Dernière modification par',
    exportAsCsv: 'Export as csv',
    criteria: 'Critères',
    ratingSettings: 'Paramètres d\'évaluation',
    pbehavior: 'Comportement périodique | Comportements périodiques',
    searchBy: 'Recherché par',
    dictionary: 'Dictionnaire',
    condition: 'Condition | Conditions',
    template: 'Template',
    pbehaviorList: 'Lister les comportements périodiques',
    canceled: 'Annulé',
    snoozed: 'En attente',
    impact: 'Impact | Impacts',
    depend: 'Depend | Depends',
    componentInfo: 'Component info | Component infos',
    connector: 'Type de connecteur',
    connectorName: 'Nom du connecteur',
    component: 'Composant',
    resource: 'Ressource',
    extraDetails: 'Détails supplémentaires',
    acked: 'Acked',
    ackedAt: 'Acked at',
    ackedBy: 'Acked by',
    resolvedAt: 'Resolved at',
    extraInfo: 'Extra info | Extra infos',
    custom: 'Personnalisé',
    eventType: 'Type d\'événement',
    sourceType: 'Type de Source',
    cycleDependency: 'Dépendance au cycle',
    checkPattern: 'Motif à carreaux',
    itemFound: '{count} article trouvé | {count} articles trouvés',
    canonicalType: 'Type canonique',
    instructions: 'Des instructions',
    playlist: 'Liste de lecture | Listes de lecture',
    actions: {
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
      [EVENT_ENTITY_TYPES.executeInstruction]: 'Exécuter la consigne',
    },
    acknowledge: 'Acquitter',
    acknowledgeAndDeclareTicket: 'Acquitter et déclarer un ticket',
    acknowledgeAndAssociateTicket: 'Acquitter et associer un ticket',
    saveChanges: 'Sauvegarder',
    reportIncident: 'Signaler un incident',
    times: {
      second: 'seconde | secondes',
      minute: 'minute | minutes',
      hour: 'heure | heures',
      day: 'jour | jours',
      week: 'semaine | semaines',
      month: 'mois | mois',
      year: 'année | années',
    },
    timeFrequencies: {
      secondly: 'Par seconde',
      minutely: 'Par minute',
      hourly: 'Par heure',
      daily: 'Quotidien',
      weekly: 'Hebdomadiare',
      monthly: 'Mensuel',
      yearly: 'Annuel',
    },
    weekDays: {
      monday: 'Lundi',
      tuesday: 'Mardi',
      wednesday: 'Mercredi',
      thursday: 'Jeudi',
      friday: 'Vendredi',
      saturday: 'Samedi',
      sunday: 'Dimanche',
    },
    months: {
      january: 'Janvier',
      february: 'Février',
      march: 'Mars',
      april: 'Avril',
      may: 'Mai',
      june: 'Juin',
      july: 'Juillet',
      august: 'Août',
      september: 'Septembre',
      october: 'Octobre',
      november: 'Novembre',
      december: 'Décembre',
    },
    stateTypes: {
      [ENTITIES_STATES.ok]: 'Ok',
      [ENTITIES_STATES.minor]: 'Mineur',
      [ENTITIES_STATES.major]: 'Majeur',
      [ENTITIES_STATES.critical]: 'Critique',
    },
    statusTypes: {
      [ENTITIES_STATUSES.closed]: 'Fermée',
      [ENTITIES_STATUSES.ongoing]: 'En cours',
      [ENTITIES_STATUSES.flapping]: 'Bagot',
      [ENTITIES_STATUSES.stealthy]: 'Furtive',
      [ENTITIES_STATUSES.cancelled]: 'Annulée',
      [ENTITIES_STATUSES.noEvents]: 'Pas d\'événements',
    },
    operators: {
      [PATTERN_OPERATORS.equal]: 'Égal',
      [PATTERN_OPERATORS.contains]: 'Contient',
      [PATTERN_OPERATORS.notEqual]: 'Inégal',
      [PATTERN_OPERATORS.notContains]: 'Ne contient pas',
      [PATTERN_OPERATORS.beginsWith]: 'Commence par',
      [PATTERN_OPERATORS.notBeginWith]: 'Ne commence pas par',
      [PATTERN_OPERATORS.endsWith]: 'Se termine par',
      [PATTERN_OPERATORS.notEndWith]: 'Ne se termine pas par',
      [PATTERN_OPERATORS.exist]: 'Exister',
      [PATTERN_OPERATORS.notExist]: 'N\'existe pas',

      [PATTERN_OPERATORS.hasEvery]: 'A chaque',
      [PATTERN_OPERATORS.hasOneOf]: 'A l\'un des',
      [PATTERN_OPERATORS.isOneOf]: 'Fait partie de',
      [PATTERN_OPERATORS.hasNot]: 'N\'a pas',
      [PATTERN_OPERATORS.isNotOneOf]: 'N\'est pas l\'un des',
      [PATTERN_OPERATORS.isEmpty]: 'Est vide',
      [PATTERN_OPERATORS.isNotEmpty]: 'N\'est pas vide',

      [PATTERN_OPERATORS.higher]: 'Plus haut que',
      [PATTERN_OPERATORS.lower]: 'Plus bas que',

      [PATTERN_OPERATORS.longer]: 'Plus long',
      [PATTERN_OPERATORS.shorter]: 'Plus court',

      [PATTERN_OPERATORS.ticketAssociated]: 'Le billet est associé',
      [PATTERN_OPERATORS.ticketNotAssociated]: 'Le billet n\'est pas associé',

      [PATTERN_OPERATORS.canceled]: 'Annulé',
      [PATTERN_OPERATORS.notCanceled]: 'Non annulé',

      [PATTERN_OPERATORS.snoozed]: 'En attente',
      [PATTERN_OPERATORS.notSnoozed]: 'Non mis en attente',

      [PATTERN_OPERATORS.acked]: 'Acquis',
      [PATTERN_OPERATORS.notAcked]: 'Non confirmé',

      [PATTERN_OPERATORS.isGrey]: 'Carrelage gris',
      [PATTERN_OPERATORS.isNotGrey]: 'Pas de carreaux gris',
    },
    entityEventTypes: {
      [EVENT_ENTITY_TYPES.ack]: 'Acquitter',
      [EVENT_ENTITY_TYPES.ackRemove]: 'Suppression d\'acquittement',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Associer un ticket',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Déclarer un incident',
      [EVENT_ENTITY_TYPES.cancel]: 'Annuler',
      [EVENT_ENTITY_TYPES.uncancel]: 'Uncancel',
      [EVENT_ENTITY_TYPES.changeState]: 'Changer d\'état',
      [EVENT_ENTITY_TYPES.check]: 'Vérifier',
      [EVENT_ENTITY_TYPES.comment]: 'Commenter l\'alarme',
      [EVENT_ENTITY_TYPES.snooze]: 'Snooze',
    },
    scenarioTriggers: {
      [SCENARIO_TRIGGERS.create]: {
        text: 'Création d\'alarme',
      },
      [SCENARIO_TRIGGERS.statedec]: {
        text: 'Diminution de la criticité',
      },
      [SCENARIO_TRIGGERS.changestate]: {
        text: 'Changement et verrouillage de la criticité',
      },
      [SCENARIO_TRIGGERS.stateinc]: {
        text: 'Augmentation de la criticité',
      },
      [SCENARIO_TRIGGERS.changestatus]: {
        text: 'Changement de statut (flapping, bagot, ...)',
      },
      [SCENARIO_TRIGGERS.ack]: {
        text: 'Acquittement d\'une alarme',
      },
      [SCENARIO_TRIGGERS.ackremove]: {
        text: 'Suppression de l\'acquittement d\'une alarme',
      },
      [SCENARIO_TRIGGERS.cancel]: {
        text: 'Annulation d\'une alarme',
      },
      [SCENARIO_TRIGGERS.uncancel]: {
        text: 'Annulation de l\'annulation d\'une alarme',
        helpText: 'L\'annulation ne peut se faire que par un événement posté sur l\'API',
      },
      [SCENARIO_TRIGGERS.comment]: {
        text: 'Commentaire sur une alarme',
      },
      [SCENARIO_TRIGGERS.done]: {
        text: 'Alarme en statut "done"',
        helpText: 'Ne peut s\'obtenir que par un événement posté sur l\'API',
      },
      [SCENARIO_TRIGGERS.declareticket]: {
        text: 'Déclaration de ticket depuis l\'interface graphique',
      },
      [SCENARIO_TRIGGERS.declareticketwebhook]: {
        text: 'Déclaration de ticket depuis un webhook',
      },
      [SCENARIO_TRIGGERS.assocticket]: {
        text: 'Association de ticket sur une alarme',
      },
      [SCENARIO_TRIGGERS.snooze]: {
        text: 'Mise en veille d\'une alarme',
      },
      [SCENARIO_TRIGGERS.unsnooze]: {
        text: 'Sortie de veille d\'une alarme',
      },
      [SCENARIO_TRIGGERS.resolve]: {
        text: 'Résolution d\'une alarme',
      },
      [SCENARIO_TRIGGERS.activate]: {
        text: 'Activation d\'une alarme',
      },
      [SCENARIO_TRIGGERS.pbhenter]: {
        text: 'Comportement périodique démarré',
      },
      [SCENARIO_TRIGGERS.pbhleave]: {
        text: 'Comportement périodique terminé',
      },
      [SCENARIO_TRIGGERS.instructionfail]: {
        text: 'Consigne manuelle en erreur',
      },
      [SCENARIO_TRIGGERS.autoinstructionfail]: {
        text: 'Consigne automatique en erreur',
      },
      [SCENARIO_TRIGGERS.instructionjobfail]: {
        text: 'Job de remédiation en erreur',
      },
      [SCENARIO_TRIGGERS.instructionjobcomplete]: {
        text: 'Job de remédiation terminé',
      },
      [SCENARIO_TRIGGERS.instructioncomplete]: {
        text: 'Consigne manuelle terminée',
      },
      [SCENARIO_TRIGGERS.autoinstructioncomplete]: {
        text: 'Consigne automatique terminée',
      },
    },
  },
  variableTypes: {
    string: 'Chaîne de caractères',
    number: 'Nombre',
    boolean: 'Booléen',
    null: 'Nul',
    array: 'Tableau',
  },
  context: {
    impacts: 'Impacts',
    dependencies: 'Dépendances',
    noEventsFilter: 'Aucun filtre d\'événements',
    impactChain: 'Chaîne d\'impact',
    impactDepends: 'Impacts/Dépendances',
    treeOfDependencies: 'Arbre des dépendances',
    infosSearchLabel: 'Rechercher une info',
    eventStatisticsMessage: '{ok} OK événements\n{ko} KO événements',
    eventStatistics: 'Statistiques d\'événement',
    actions: {
      titles: {
        editEntity: 'Éditer l\'entité',
        duplicateEntity: 'Dupliquer l\'entité',
        deleteEntity: 'Supprimer l\'entité',
        pbehavior: 'Comportement périodique',
        variablesHelp: 'Liste des variables disponibles',
        massEnable: 'Activer les entités',
        massDisable: 'Désactiver les entités',
      },
    },
    fab: {
      common: 'Ajouter une nouvelle entité',
      addService: 'Ajouter une nouvelle entité de service',
    },
    popups: {
      massDeleteWarning: 'La suppression en masse ne peut pas être appliquée pour certains des éléments sélectionnés, ils ne seront donc pas supprimés.',
    },
  },
  search: {
    alarmAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
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
    contextAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
      + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
      + '<p>Le "-" avant la recherche est obligatoire</p>\n'
      + '<p>Opérateurs:\n'
      + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
      + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
      + '<dl><dt>Exemples :</dt><dt>- Name = "name_1"</dt>\n'
      + '    <dd>Entités dont le names est "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
      + '    <dd>Entités dont le names est "name_1" et la types est "service"</dd><dt>- infos.custom.value="Custom value" OR Type="resource"</dt>\n'
      + '    <dd>Entités dont le infos.custom.value est "Custom value" ou la type est "resource"</dd><dt>- infos.custom.value LIKE 1 OR infos.custom.value LIKE 2</dt>\n'
      + '    <dd>Entités dont le infos.custom.value contient 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
      + '    <dd>Entités dont le name n\'est pas "name_1"</dd>\n'
      + '</dl>',
    dynamicInfoAdvancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
      + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
      + '<p>Le "-" avant la recherche est obligatoire</p>\n'
      + '<p>Opérateurs:\n'
      + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
      + '<p>Pour effectuer une recherche dans les "modèles", utilisez le mot-clé "pattern" comme &lt;NomColonne&gt;</p>\n'
      + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
      + '<dl><dt>Examples :</dt><dt>- description = "testdyninfo"</dt>\n'
      + '    <dd>Règles d\'Informations dynamiques dont la description est "testdyninfo"</dd><dt>- pattern = "SEARCHPATTERN1"</dt>\n'
      + '    <dd>Règles d\'Informations dynamiques dont un des modèles est "SEARCHPATTERN1"</dd><dt>- pattern LIKE "SEARCHPATTERN2"</dt>\n'
      + '    <dd>Règles d\'Informations dynamiques dont un des modèles contient "SEARCHPATTERN2"</dd>'
      + '</dl>',
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
      incorrectEmailOrPassword: 'Mot de passe ou Email incorrect',
    },
  },
  alarmList: {
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
        changeState: 'Changer et verrouiller la criticité',
        variablesHelp: 'Liste des variables disponibles',
        history: 'Historique',
        groupRequest: 'Proposition de regroupement pour méta-alarmes',
        manualMetaAlarmGroup: 'Gestion manuelle des méta-alarmes',
        manualMetaAlarmUngroup: 'Dissocier l\'alarme de la méta-alarme manuelle',
        comment: 'Commenter l\'alarme',
        executeInstruction: 'Exécuter la consigne "{instructionName}"',
        resumeInstruction: 'Reprendre la consigne "{instructionName}"',
      },
      iconsTitles: {
        ack: 'Acquitter',
        declareTicket: 'Déclarer un incident',
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
        [EVENT_ENTITY_TYPES.instructionJobStart]: 'L\'exécution d\'une tâche de remédiation a été démarrée',
        [EVENT_ENTITY_TYPES.instructionJobComplete]: 'L\'exécution de la tâche de remédiation est terminée',
        [EVENT_ENTITY_TYPES.instructionJobAbort]: 'L\'exécution de la tâche de remédiation a été abandonnée',
        [EVENT_ENTITY_TYPES.instructionJobFail]: 'L\'exécution de la tâche de remédiation a échouée',
        [EVENT_ENTITY_TYPES.autoInstructionStart]: 'La consigne automatique a été lancée',
        [EVENT_ENTITY_TYPES.autoInstructionComplete]: 'La consigne automatique a été complétée',
        [EVENT_ENTITY_TYPES.autoInstructionFail]: 'La consigne automatique a échoué',
        [EVENT_ENTITY_TYPES.autoInstructionAlreadyRunning]: 'La consigne automatique a déjà été lancée pour une autre alarme',
        [EVENT_ENTITY_TYPES.junitTestSuiteUpdate]: 'La suite de tests a été mise à jour',
        [EVENT_ENTITY_TYPES.junitTestCaseUpdate]: 'Le cas de test a été mis à jour',
      },
    },
    tabs: {
      moreInfos: 'Plus d\'infos',
      timeLine: 'Chronologie',
      alarmsChildren: 'Alarmes liées',
      trackSource: 'Source de la piste',
      impactChain: 'Chaîne d\'impact',
      entityGantt: 'Diagramme de Gantt',
    },
    moreInfos: {
      defineATemplate: 'Pour définir le template de cette fenêtre, rendez-vous dans les paramètres du bac à alarmes.',
    },
    infoPopup: 'Info popup',
    tooltips: {
      priority: 'Le paramètre de priorité est dérivé de la gravité de l\'alarme multipliée par le niveau d\'impact de l\'entité sur laquelle l\'alarme est déclenchée',
      hasInstruction: 'Au moins une consigne de remédiation est attachée à cette alarme',
      hasManualInstructionInRunning: 'Consigne manuelle en cours',
      hasAutoInstructionInRunning: 'Consigne automatique en cours',
      allAutoInstructionExecuted: 'Toutes les consignes automatiques ont été exécutées',
      awaitingInstructionComplete: 'En attente de la fin de la consigne pour terminer',
      autoInstructionsFailed: 'Les instructions automatiques ont échoué',
    },
    metrics: {
      [ALARM_METRIC_PARAMETERS.createdAlarms]: 'Nombre d\'alarmes créées',
      [ALARM_METRIC_PARAMETERS.activeAlarms]: 'Nombre d\'alarmes actives',
      [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: 'Nombre d\'alarmes non affichées',
      [ALARM_METRIC_PARAMETERS.instructionAlarms]: 'Nombre d\'alarmes en cours de remédiation automatique',
      [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: 'Nombre d\'alarmes avec PBehavior',
      [ALARM_METRIC_PARAMETERS.correlationAlarms]: 'Nombre d\'alarmes corrélées',
      [ALARM_METRIC_PARAMETERS.ackAlarms]: 'Nombre d\'alarmes avec acquittement',
      [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: 'Number of active alarms with acks',
      [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: 'Nombre d\'accusés de réception annulés',
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: 'Nombre d\'alarmes actives avec acks',
      [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: 'Nombre d\'alarmes actives sans tickets',
      [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '% d\'alarmes corrélées',
      [ALARM_METRIC_PARAMETERS.ratioInstructions]: '% d\'alarmes avec remédiation automatique',
      [ALARM_METRIC_PARAMETERS.ratioTickets]: '% d\'alarmes avec tickets créés',
      [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '% d\'alarmes corrigées manuellement',
      [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '% d\'alarmes non affichées',
      [ALARM_METRIC_PARAMETERS.averageAck]: 'Délai moyen d\'acquittement des alarmes',
      [ALARM_METRIC_PARAMETERS.averageResolve]: 'Temps moyen pour résoudre les alarmes',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: 'Nombre d\'alarmes corrigées manuellement',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: 'Nombre d\'alarmes avec instructions manuelles',
    },
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
    lastAlarmDate: 'Date de la dernière alarme',
    massRemove: 'Supprimer les comportements',
    massEnable: 'Activer les comportements',
    massDisable: 'Désactiver les comportements',
    searchHelp: '<span>Aide sur la recherche avancée :</span>\n'
      + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
      + '<p>Le "-" avant la recherche est obligatoire</p>\n'
      + '<p>Opérateurs : <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
      + '<p>Pour effectuer une recherche dans les "patterns", utilisez le mot-clé "pattern" comme &lt;NomColonne&gt;</p>\n'
      + '<p>Les types de valeurs : String entre doubles guillemets, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
      + '<dl>'
      + '  <dt>Examples :</dt>'
      + '  <dt>- name = "name_1"</dt>\n'
      + '  <dd>Le nom du comportement périodique est "name_1"</dd>\n'
      + '  <dt>- rrule = "rrule_1"</dt>\n'
      + '  <dd>La règle de comportement périodique est "rrule_1"</dd>\n'
      + '  <dt>- filter = "filter_1"</dt>\n'
      + '  <dd>Le filtre de comportement périodique est "filter_1"</dd>\n'
      + '  <dt>- type.name = "type_name_1"</dt>\n'
      + '  <dd>Le nom du type de comportement périodique est "type_name_1"</dd>\n'
      + '  <dt>- reason.name = "reason_name_1"</dt>\n'
      + '  <dd>Le nom de la raison du comportement périodique est "reason_name_1"</dd>'
      + '</dl>',
    tabs: {
      entities: 'Entités',
    },
  },
  settings: {
    titles: {
      [SIDE_BARS.alarmSettings]: 'Paramètres du bac à alarmes',
      [SIDE_BARS.contextSettings]: 'Paramètres de l\'explorateur de contexte',
      [SIDE_BARS.serviceWeatherSettings]: 'Paramètres de la météo des services',
      [SIDE_BARS.statsCalendarSettings]: 'Paramètres du calendrier',
      [SIDE_BARS.textSettings]: 'Paramètres du widget de texte',
      [SIDE_BARS.counterSettings]: 'Paramètres du widget de compteur',
      [SIDE_BARS.testingWeatherSettings]: 'Paramètres du widget scénario des tests',
    },
    openedTypes: {
      [ALARMS_OPENED_VALUES.opened]: 'Alarmes ouvertes',
      [ALARMS_OPENED_VALUES.resolved]: 'Toutes les alarmes résolues',
      [ALARMS_OPENED_VALUES.all]: 'Alarmes ouvertes et récemment résolues',
    },
    advancedSettings: 'Paramètres avancés',
    widgetTitle: 'Titre du widget',
    columnName: 'Nom de la colonne',
    defaultSortColumn: 'Colonne de tri par défaut',
    sortColumnNoData: 'Appuyez sur <kbd>enter</kbd> pour en créer une nouvelle',
    columnNames: 'Nom des colonnes',
    exportColumnNames: 'Nom des colonnes à exporter',
    groupColumnNames: 'Nom des colonnes des méta-alarmes',
    trackColumnNames: 'Suivre les colonnes de source d\'alarme',
    treeOfDependenciesColumnNames: 'Nom des colonnes pour l\'arborescence des dépendances',
    orderBy: 'Trier par',
    periodicRefresh: 'Rafraichissement périodique',
    defaultNumberOfElementsPerPage: 'Nombre d\'élements par page par défaut',
    elementsPerPage: 'Élements par page',
    filterOnOpenResolved: 'Filtre sur Ouverte/Résolue',
    open: 'Ouverte',
    resolved: 'Résolue',
    filters: 'Filtres',
    filterEditor: 'Éditeur de filtre',
    isAckNoteRequired: 'Champ \'Note\' requis lors d\'un acquittement ?',
    isSnoozeNoteRequired: 'Champ \'Note\' requis lors d\'une mise en veille ?',
    linksCategoriesAsList: 'Afficher les liens sous forme de liste ?',
    linksCategoriesLimit: 'Nombre d\'éléments de catégorie',
    isMultiAckEnabled: 'Acquittement multiple',
    isMultiDeclareTicketEnabled: 'Déclarer un ticket multiple',
    fastAckOutput: 'Commentaire d\'acquittement rapide',
    isHtmlEnabledOnTimeLine: 'HTML activé dans la chronologie ?',
    isCorrelationEnabled: 'Corrélation activée ?',
    duration: 'Durée',
    tstop: 'Date de fin',
    periodsNumber: 'Nombre d\'étapes',
    yesNoMode: 'Mode Oui/Non',
    selectAFilter: 'Sélectionner un filtre',
    exportAsCsv: 'Exporter le widget sous forme de fichier csv',
    criticityLevels: 'Niveaux de criticité',
    isPriorityEnabled: 'Afficher la priorité',
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
    infoPopup: {
      title: 'Info popup',
      fields: {
        column: 'Colonne',
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
      withTemplate: 'Modèle personnalisé',
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
    receiveByApi: 'Réponse de l\'API',
    serverStorage: 'Stockage serveur',
    filenameRecognition: 'Reconnaissance du nom de fichier',
    resultDirectory: 'Stockage des résultats de test',
    screenshotDirectories: {
      title: 'Paramètres de stockage des captures d\'écran',
      helpText: 'Définir où les captures d\'écran sont stockées',
    },
    screenshotMask: {
      title: 'Règle de nommage des fichiers capture d\'écran',
      helpText: '<dl>'
        + '<dt>Définissez la règle de nommage des fichiers dont les captures d\'écran sont créées à l\'aide des variables suivantes:<dt>\n'
        + '<dd>- nom du cas de test %test_case%</dd>\n'
        + '<dd>- date (YYYY, MM, DD)</dd>\n'
        + '<dd>- temps d\'exécution (hh, mm, ss)</dd>'
        + '</dl>',
    },
    videoDirectories: {
      title: 'Paramètres de stockage vidéo',
      helpText: 'Définir où la vidéo est stockée',
    },
    videoMask: {
      title: 'Règle de nommage des fichiers vidéo',
      helpText: '<dl>'
        + '<dt>Définissez la règle de nommage des fichiers dont les vidéos sont créées à l\'aide des variables suivantes:<dt>\n'
        + '<dd>- nom du cas de test %test_case%</dd>\n'
        + '<dd>- date (YYYY, MM, DD)</dd>\n'
        + '<dd>- temps d\'exécution (hh, mm, ss)</dd>'
        + '</dl>',
    },
    stickyHeader: 'En-tête collant',
    reportFileRegexp: {
      title: 'Masque de fichier de rapport',
      helpText: '<dl>'
        + '<dt>Définir le nom de fichier regexp de quel rapport:<dt>\n'
        + '<dd>Par exemple:</dd>\n'
        + '<dd>"^(?P&lt;name&gt;\\\\w+)_(.+)\\\\.xml$"</dd>\n'
        + '</dl>',
    },
    density: {
      title: 'Vue par défaut',
      comfort: 'Vue confort',
      compact: 'Vue compacte',
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
        groupTags: 'Étiquettes de groupe',
      },
      success: {
        create: 'Nouvelle vue créée !',
        edit: 'Vue éditée avec succès !',
        duplicate: 'Afficher dupliqué avec succès !',
        delete: 'Vue supprimée avec succès !',
      },
      fail: {
        create: 'Erreur lors de la création de la vue...',
        edit: 'Erreur lors de l\'édition de la vue...',
        duplicate: 'Échec de la duplication de la vue...',
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
        ackResources: 'Acquitter les ressources',
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
      title: 'Mise en veille',
      fields: {
        duration: 'Durée',
      },
    },
    createCancelEvent: {
      title: 'Annuler',
    },
    createGroupRequestEvent: {
      title: 'Proposition de regroupement pour méta-alarmes',
    },
    createGroupEvent: {
      title: 'Créer une méta-alarme',
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
            noEnding: 'Sans fin',
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
        noAckItems: 'Il y a {count} élément sans accusé de réception. L\'événement acquitté pour l\'élément sera envoyé avant. | Il y a {count} éléments sans accusé de réception. Les événements acquittés pour les éléments seront envoyés avant.',
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
      addInfoPopup: {
        title: 'Ajouter une popup d\'info',
      },
    },
    variablesHelp: {
      variables: 'Variables',
      copyToClipboard: 'Copier dans le presse-papier',
    },
    service: {
      actionPending: 'action en attente | actions en attente',
      refreshEntities: 'Rafraîchir la liste des entités',
      editPbehaviors: 'Éditer les comportements périodiques',
      entity: {
        tabs: {
          info: 'Info',
          treeOfDependencies: 'Arbre de dépendances',
        },
      },
    },
    createFilter: {
      create: {
        title: 'Créer un filtre',
      },
      edit: {
        title: 'Éditer un filtre',
      },
      duplicate: {
        title: 'Dupliquer un filtre',
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
        [WIDGET_TYPES.statsCalendar]: {
          title: 'Calendrier',
        },
        [WIDGET_TYPES.text]: {
          title: 'Texte',
        },
        [WIDGET_TYPES.counter]: {
          title: 'Compteur',
        },
        [WIDGET_TYPES.testingWeather]: {
          title: 'Scénarios Junit',
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
      prefixTitle: '{prefix} - bac à alarmes',
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
    createEventFilter: {
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
        success: 'Règle supprimée avec succès !',
      },
      editPattern: 'Éditer le modèle',
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
    createSnmpRule: {
      create: {
        title: 'Créer une règle SNMP',
      },
      edit: {
        title: 'Modifier la règle SNMP',
      },
    },
    selectView: {
      title: 'Sélectionner une vue',
    },
    selectViewTab: {
      title: 'Sélectionnez l\'onglet',
    },
    createDynamicInfo: {
      alarmUpdate: 'La règle mettra à jour les alarmes existantes !',
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
        emptyInfos: 'Au moins une information doit être ajoutée.',
      },
      steps: {
        infos: {
          title: 'Informations',
        },
        patterns: {
          title: 'Modèles',
          alarmPatterns: 'Modèles des alarmes',
          entityPatterns: 'Modèles des entités',
          validationError: 'Au moins un modèle est requis. Merci d\'ajouter un modèle sur les alarmes et/ou un modèle sur les événements',
        },
      },
    },
    createDynamicInfoInformation: {
      create: {
        title: 'Ajouter une information à la règle d\'information dynamique',
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
    },
    createPlaylist: {
      create: {
        title: 'Créer une liste de lecture',
      },
      edit: {
        title: 'Éditer une liste de lecture',
      },
      duplicate: {
        title: 'Dupliquer une liste de lecture',
      },
      errors: {
        emptyTabs: 'Merci de ajouter un onglet',
      },
      fields: {
        interval: 'Période',
        unit: 'Unité',
      },
      groups: 'Groupe',
      manageTabs: 'Gérer les onglets',
    },
    pbehaviorPlanning: {
      title: 'Comportements périodiques',
    },
    selectExceptionsLists: {
      title: 'Choisissez la liste des exceptions',
    },
    createRrule: {
      title: 'Créer une récurrence',
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
      title: 'Créer une raison',
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
      emptyExdates: 'Aucune date d\'exception ajoutée pour le moment',
    },
    createManualMetaAlarm: {
      title: 'Gestion manuelle des méta-alarmes',
      noData: 'Aucune méta-alarme correspondante. Appuyez sur <kbd>Entrée</kbd> pour en créer un nouveau',
      fields: {
        metaAlarm: 'Méta-alarme manuelle',
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
      fields: {
        host: 'Hôte',
        token: 'Jeton d\'autorisation',
      },
    },
    createRemediationJob: {
      create: {
        title: 'Créer une tâche',
        popups: {
          success: '{jobName} a été créée avec succès',
        },
      },
      edit: {
        title: 'Éditer une tâche',
        popups: {
          success: '{jobName} a été modifiée avec succès',
        },
      },
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
      title: 'Évaluer cette consigne "{name}"',
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
      title: 'Configuration d\'état du test JUnit',
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
    remediationInstructionApproval: {
      title: 'Approbation des consignes',
      requested: 'demandé pour approbation',
      tabs: {
        updated: 'Mise à jour',
        original: 'Original',
      },
    },
    createAlarmIdleRule: {
      create: {
        title: 'Créer une règle d\'alarme',
      },
      edit: {
        title: 'Modifier la règle d\'alarme',
      },
      duplicate: {
        title: 'Règle d\'alarme en double',
      },
    },
    createEntityIdleRule: {
      create: {
        title: 'Créer une règle d\'entité',
      },
      edit: {
        title: 'Modifier la règle d\'entité',
      },
      duplicate: {
        title: 'Règle d\'entité en double',
      },
    },
    createAlarmStatusRule: {
      flapping: {
        create: {
          title: 'Créer une règle de bagot',
        },
        edit: {
          title: 'Modifier la règle de bagot',
        },
        duplicate: {
          title: 'Dupliquer la règle de bagot',
        },
      },
      resolve: {
        create: {
          title: 'Créer une règle de résolution',
        },
        edit: {
          title: 'Modifier la règle de résolution',
        },
        duplicate: {
          title: 'Dupliquer la règle de résolution',
        },
      },
    },
    webSocketError: {
      title: 'Erreur de connexion WebSocket',
      text: '<p>Les Websockets ne sont pas disponibles, les fonctionnalités suivantes sont donc restreintes:</p>'
        + '<p>'
        + '<ul>'
        + '<li>En-tête Healthcheck</li>'
        + '<li>Graphique du réseau Healthcheck</li>'
        + '<li>Messages diffusés actifs</li>'
        + '<li>Sessions d\'utilisateurs actifs</li>'
        + '<li>Exécution de la remédiation</li>'
        + '</ul>'
        + '</p>'
        + '<p>Veuillez vérifier la configuration de votre serveur.</p>',
      shortText: '<p>Les Websockets ne sont pas disponibles, les fonctionnalités suivantes sont donc restreintes:</p>'
        + '<p>'
        + '<ul>'
        + '<li>Messages diffusés actifs</li>'
        + '<li>Sessions d\'utilisateurs actifs</li>'
        + '</ul>'
        + '</p>'
        + '<p>Veuillez vérifier la configuration de votre serveur.</p>',
    },
    confirmationPhrase: {
      phrase: 'Phrase',
      updateStorageSettings: {
        title: 'Changement de politique de stockage. Êtes vous sur ?',
        text: 'Vous êtes sur le point d\'enregistrer une politique d\'archivage et/ou de suppression de données.\n'
          + '<strong>Les opérations qui en découleront, notamment la suppression de données, seront irreversibles.</strong>',
        phraseText: 'Merci de recopier le texte qui suit pour confirmer:',
        phrase: 'modifier la politique de stockage',
      },
      cleanStorage: {
        title: 'Archivage/Suppression des entités désactivées. Êtes vous sur ?',
        text: 'Vous êtes sur le point d\'archiver et/ou de supprimer des données.\n'
          + '<strong>Les opérations de suppression sont irreversibles.</strong>',
        phraseText: 'Merci de recopier le texte qui suit pour confirmer:',
        phrase: 'archiver ou supprimer',
      },
    },
    pbehaviorsCalendar: {
      title: 'Comportements périodiques',
      entity: {
        title: 'Comportements périodiques - {name}',
      },
    },
    createAlarmPattern: {
      create: {
        title: 'Créer un filtre d\'alarme',
      },
      edit: {
        title: 'Modifier le modèle d\'alarme',
      },
    },
    createCorporateAlarmPattern: {
      create: {
        title: 'Créer un filtre d\'alarme partagé',
      },
      edit: {
        title: 'Modifier le filtre d\'alarme partagé',
      },
    },
    createEntityPattern: {
      create: {
        title: 'Créer un filtre d\'entité',
      },
      edit: {
        title: 'Modifier le modèle d\'entité',
      },
    },
    createCorporateEntityPattern: {
      create: {
        title: 'Créer un filtre d\'entité partagée',
      },
      edit: {
        title: 'Modifier le filtre d\'entité partagée',
      },
    },
    createPbehaviorPattern: {
      create: {
        title: 'Créer un filtre de comportement',
      },
      edit: {
        title: 'Modifier le modèle de comportement',
      },
    },
    createCorporatePbehaviorPattern: {
      create: {
        title: 'Créer un filtre de comportement partagé',
      },
      edit: {
        title: 'Modifier le filtre de comportement partagé',
      },
    },
    createShareToken: {
      create: {
        title: 'Créer un jeton de partage',
      },
    },
  },
  tables: {
    noData: 'Aucune donnée',
    contextEntities: {
      columns: {
        name: 'Nom',
        type: 'Type',
        _id: 'Id',
      },
    },
    noColumns: {
      message: 'Veuillez sélectionner au moins une colonne',
    },
    broadcastMessages: {
      statuses: {
        [BROADCAST_MESSAGES_STATUSES.active]: 'Actif',
        [BROADCAST_MESSAGES_STATUSES.pending]: 'En attente',
        [BROADCAST_MESSAGES_STATUSES.expired]: 'Expiré',
      },
    },
  },
  recurrenceRule: {
    advancedHint: 'Séparer les nombres par une virgule',
    freq: 'Fréquence',
    until: 'Jusqu\'à',
    byweekday: 'Par jour de la semaine',
    count: 'Répéter',
    interval: 'Intervalle',
    wkst: 'Semaine de début',
    bymonth: 'Par mois',
    bysetpos: 'Par position',
    bymonthday: 'Par jour du mois',
    byyearday: 'Par jour de l\'année',
    byweekno: 'Par semaine n°',
    byhour: 'Par heure',
    byminute: 'Par minute',
    bysecond: 'Par seconde',
    tabs: {
      simple: 'Simple',
      advanced: 'Avancé',
    },
    errors: {
      main: 'La récurrence choisie n\'est pas valide. Nous vous recommandons de la modifier avant de sauvegarder',
    },
    periodsRanges: {
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek]: 'Cette semaine',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek]: 'Semaine prochaine',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks]: 'Deux prochaines semaines',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth]: 'Ce mois',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth]: 'Le mois prochain',
    },
    tooltips: {
      bysetpos: 'Si renseigné, doit être un ou plusieurs nombres entiers, positifs ou négatifs. Chaque entier correspondra à la ènième occurence de la règle dans l\'intervalle de fréquence. Par exemple, une \'bysetpos\' de -1 combinée à une fréquence mensuelle, et une \'byweekday\' de (lundi, mardi, mercredi, jeudi, vendredi), va nous donner le dernier jour travaillé de chaque mois',
      bymonthday: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours du mois auxquels s\'appliquera la récurrence.',
      byyearday: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux jours de l\'année auxquels  s\'appliquera la récurrence.',
      byweekno: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux numéros de semaine auxquelles s\'appliquera la récurrence. Les numéros de semaines sont ceux de ISO8601, la première semaine de l\'année étant celle contenant au moins 4 jours de cette année.',
      byhour: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux heures auxquelles s\'appliquera la récurrence.',
      byminute: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux minutes auxquelles s\'appliquera la récurrence.',
      bysecond: 'Si renseigné, doit être un ou plusieurs nombres entiers, correspondant aux secondes auxquelles s\'appliquera la récurrence.',
    },
  },
  errors: {
    default: 'Une erreur s\'est produite...',
    lineNotEmpty: 'Cette ligne n\'est pas vide',
    JSONNotValid: 'JSON non valide',
    versionNotFound: 'Erreur dans la récupération du numéro de version...',
    statsRequestProblem: 'Erreur dans la récupération des statistiques',
    statsWrongEditionError: "Les widgets de statistiques ne sont pas disponibles dans l'édition 'community' de Canopsis",
    socketConnectionProblem: 'Problème de connexion aux websockets',
    endDateLessOrEqualStartDate: 'La date de fin doit se situer après la date de début',
    unknownWidgetType: 'Type de widget inconnu: {type}',
    unique: 'Le champ doit être unique',
  },
  warnings: {
    authTokenExpired: 'Le jeton d\'authentification a expiré',
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
    linkCopied: 'Lien copié dans le presse-papier',
    authKeyCopied: 'Clé d\'authentification copiée dans le presse-papier',
    widgetIdCopied: 'Identifant du widget copié dans le presse-papier',
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
      mixFilters: 'Mélanger les filtres',
    },
    buttons: {
      list: 'Gérer les filtres',
    },
  },
  stats: {
    types: {
      [STATS_TYPES.alarmsCreated.value]: 'Alarmes créées',
      [STATS_TYPES.alarmsResolved.value]: 'Alarmes résolues',
      [STATS_TYPES.alarmsCanceled.value]: 'Alarmes annulées',
      [STATS_TYPES.alarmsAcknowledged.value]: 'Alarmes acquittées',
      [STATS_TYPES.ackTimeSla.value]: 'Taux d\'acquittement conforme SLA',
      [STATS_TYPES.resolveTimeSla.value]: 'Taux de résolution conforme SLA',
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
    externalData: 'Données externes',
    actionsRequired: 'Veuillez ajouter au moins une action',
    configRequired: 'No configuration defined. Please add at least one config parameter',
    idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
    editPattern: 'Éditer le modèle',
    advanced: 'Avancée',
    addAField: 'Ajouter un champ',
    simpleEditor: 'Éditeur simple',
    field: 'Champ',
    value: 'Valeur',
    advancedEditor: 'Éditeur avancé',
    comparisonRules: 'Règles de comparaison',
    editActions: 'Éditer les actions',
    addAction: 'Ajouter une action',
    editAction: 'Éditer une action',
    actions: 'Actions',
    onSuccess: 'En cas de succès',
    onFailure: 'En cas d\'échec',
    configuration: 'Configuration',
    resource: 'ID de ressource ou modèle',
    component: 'ID de composant ou modèle',
    connector: 'ID ou modèle de connecteur',
    connectorName: 'Nom ou modèle de connecteur',
    types: {
      [EVENT_FILTER_TYPES.drop]: 'Drop',
      [EVENT_FILTER_TYPES.break]: 'Break',
      [EVENT_FILTER_TYPES.enrichment]: 'Enrichment',
      [EVENT_FILTER_TYPES.changeEntity]: 'Change entity',
    },
    tooltips: {
      addValueRuleField: 'Ajouter une règle',
      editValueRuleField: 'Éditer la règle',
      addObjectRuleField: 'Ajouter un groupe de règles',
      editObjectRuleField: 'Éditer le groupe de règles',
      removeRuleField: 'Supprimer le groupe/la règle',
      copyFromHelp: '<p>Les variables accessibles sont: <strong>Event</strong></p>'
        + '<i>Quelques exemples:</i> <span>"Event.ExtraInfos.datecustom"</span>',
    },
    actionsTypes: {
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy]: {
        text: 'Copier une valeur d\'un champ d\'événement à un autre',
        message: 'Cette action est utilisée pour copier la valeur d\'un contrôle dans un événement.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description.\n- valeur : le nom du champ dont la valeur doit être copiée. Il peut s\'agir d\'un champ d\'événement, d\'un sous-groupe d\'une expression régulière ou d\'une donnée externe.\n- nom : le nom du champ événement dans lequel la valeur doit être copiée.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo]: {
        text: 'Copier une valeur d\'un champ d\'un événement vers une information d\'une entité',
        message: 'Cette action est utilisée pour copier la valeur du champ d\'un événement dans le champ d\'une entité. Notez que l\'entité doit être ajoutée à l\'événement en premier.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description.\n- nom : le nom du champ d\'une entité.\n- valeur : le nom du champ dont la valeur doit être copiée. Il peut s\'agir d\'un champ d\'événement, d\'un sous-groupe d\'une expression régulière ou d\'une donnée externe.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo]: {
        text: 'Définir une information d\'une entité sur une constante',
        message: 'Cette action permet de définir les informations dynamiques d\'une entité correspondant à l\'événement. Notez que l\'entité doit être ajoutée à l\'événement en premier.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description.\n- nom : le nom du champ.\n- valeur : la valeur d\'un champ.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate]: {
        text: 'Définir une chaîne d\'informations sur une entité à l\'aide d\'un modèle',
        message: 'Cette action permet de modifier les informations dynamiques d\'une entité correspondant à l\'événement. Notez que l\'entité doit être ajoutée à l\'événement.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description\n- nom : le nom du champ.\n- valeur : le modèle utilisé pour déterminer la valeur de la donnée.\nDes modèles {{.Event.NomDuChamp}}, des expressions régulières ou des données externes peuvent être utilisés.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField]: {
        text: 'Définir un champ d\'un événement sur une constante',
        message: 'Cette action peut être utilisée pour modifier un champ de l\'événement.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description\n- nom : le nom du champ.\n- valeur : la nouvelle valeur du champ.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate]: {
        text: 'Définir un champ de chaîne d\'un événement à l\'aide d\'un modèle',
        message: 'Cette action vous permet de modifier un champ d\'événement à partir d\'un modèle.',
        description: 'Les paramètres de l\'action sont :\n- description (optionnel) : la description\n- nom : le nom du champ.\n- valeur : le modèle utilisé pour déterminer la valeur du champ.\n  Des modèles {{.Event.NomDuChamp}}, des expressions régulières ou des données externes peuvent être utilisés.',
      },
    },
  },
  metaAlarmRule: {
    outputTemplate: 'Modèle de message',
    thresholdType: 'Type de seuil',
    thresholdRate: 'Taux de déclenchement',
    thresholdCount: 'Seuil de déclenchement',
    timeInterval: 'Intervalle de temps',
    valuePath: 'Chemin de valeur | Chemins de valeur',
    autoResolve: 'Résolution automatique',
    idHelp: 'Si ce champ n\'est pas renseigné, un identifiant unique sera généré automatiquement à la création de la règle',
    corelId: 'Identifiant de corrélation',
    corelIdHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
      + '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelStatus: 'Statut de corrélation',
    corelStatusHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
      + '<i>Quelques exemples:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelParent: 'Corrélation parent',
    corelChild: 'Corrélation enfant',
    outputTemplateHelp: '<p>Les variables accessibles sont:</p>\n'
      + '<p><strong>.Count</strong>: Le nombre d\'alarmes conséquences attachées à la méta-alarme.</p>'
      + '<p><strong>.Children</strong>: L\'ensemble des variables de la dernière alarme conséquence attachée à la méta-alarme.</p>'
      + '<p><strong>.Rule</strong>: Les informations administratives de la méta-alarme en elle-même.</p>'
      + '<p>Quelques exemples:</p>'
      + '<p><strong>{{ .Count }} conséquences;</strong> Message de la dernière alarme conséquence : <strong>{{ .Children.Alarm.Value.State.Message }};</strong> Règle : <strong>{{ .Rule.Name }};</strong></p>'
      + '<p>Un message informatif statique</p>'
      + '<p>Corrélé par la règle <strong>{{ .Rule.Name }}</strong></p>',
    removeConfirmationText: 'Lors de la suppression d\'une règle de méta-alarme, toutes les méta-alarmes correspondantes seront également supprimées.\n'
      + 'Êtes-vous sûr de continuer?',
    errors: {
      noValuePaths: 'Vous devez ajouter au moins un chemin de valeur',
    },
  },
  layout: {
    sideBar: {
      buttons: {
        edit: 'Activer/Désactiver le mode d\'édition',
        create: 'Créer une vue',
        settings: 'Paramètres',
      },
      loggedUsersCount: 'Sessions actives',
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
      importExportViews: 'Import/Export',
      stateSettings: 'Paramètres d\'état',
      storageSettings: 'Paramètres de stockage',
      notificationsSettings: 'Paramètres des notifications',
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
    copyWidgetId: 'Copier l\'identifiant du widget',
    autoHeightButton: 'Si ce bouton est sélectionné, la hauteur sera calculée automatiquement.',
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
      ip_or_fqdn: 'Le champ doit être une adresse IP ou un nom de domaine valide',
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
      size: 'Le champ doit avoir un poids inférieur à {1} Ko',
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
    grey: 'Gris',
    primaryIcon: 'Icône principale',
    secondaryIcon: 'Icône secondaire',
    massActions: 'Actions de masse',
    cannotBeApplied: 'Cette action ne peut pas être appliquée',
    actions: {
      [WEATHER_ACTIONS_TYPES.entityAck]: 'Acquitter',
      [WEATHER_ACTIONS_TYPES.entityValidate]: 'Valider',
      [WEATHER_ACTIONS_TYPES.entityInvalidate]: 'Invalider',
      [WEATHER_ACTIONS_TYPES.entityPause]: 'Pause',
      [WEATHER_ACTIONS_TYPES.entityPlay]: 'Supprimer la pause',
      [WEATHER_ACTIONS_TYPES.entityCancel]: 'Annuler',
      [WEATHER_ACTIONS_TYPES.entityAssocTicket]: 'Associer un ticket',
      [WEATHER_ACTIONS_TYPES.entityComment]: 'Commenter l\'alarme',
      [WEATHER_ACTIONS_TYPES.executeInstruction]: 'Exécuter l\'instruction',
      [WEATHER_ACTIONS_TYPES.declareTicket]: 'Déclarer un incident',
    },
    iconTypes: {
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactif',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Entretien',
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
      notification: 'Droits de notification',
      profile: 'Droits de profil',
    },
    business: {
      [USER_PERMISSIONS_PREFIXES.business.common]: 'Droits communs',
      [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Droits pour le widget : Bac à alarmes',
      [USER_PERMISSIONS_PREFIXES.business.context]: 'Droits pour le widget : Explorateur de contexte',
      [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Droits pour le widget: Météo des services',
      [USER_PERMISSIONS_PREFIXES.business.counter]: 'Droits pour le widget : Compteur',
      [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Droits pour le widget : Scénario des tests',
    },
    api: {
      general: 'Général',
      rules: 'Règles',
      remediation: 'Remédiation',
      pbehavior: 'PBehavior',
    },
  },

  pbehavior: {
    periodsCalendar: 'Calendrier avec périodes',
    buttons: {
      addFilter: 'Ajouter un filtre',
      editFilter: 'Modifier le filtre',
      addRRule: 'Ajouter une règle de récurrence',
      editRrule: 'Modifier la règle de récurrence',
    },
  },

  pbehaviorExceptions: {
    title: 'Dates d\'exception',
    create: 'Ajouter une date d\'exception',
    choose: 'Sélectionnez la liste d\'exclusion',
    usingException: 'Ne peut pas être supprimé car il est en cours d\'utilisation.',
    emptyExceptions: 'Aucune exception ajoutée pour le moment.',
  },

  pbehaviorTypes: {
    usingType: 'Le type ne peut être supprimé car il est en cours d\'utilisation.',
    defaultType: 'Le type par défaut ne peut pas être modifié.',
    types: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Actif',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactif',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Entretien',
    },
  },

  pbehaviorReasons: {
    usingReason: 'La raison ne peut pas être supprimée car elle est en cours d\'utilisation.',
  },

  planning: {
    tabs: {
      type: 'Type',
      reason: 'Raison',
      exceptions: 'Dates d\'exception',
    },
  },

  healthcheck: {
    metricsUnavailable: 'Les métriques ne sont pas collectées',
    notRunning: '{name} n\'est pas disponible',
    queueOverflow: 'Débordement de file d\'attente',
    lackOfInstances: 'Manque d\'instances',
    diffInstancesConfig: 'Configuration des instances non valide',
    queueLength: 'Longueur de la file d\'attente {queueLength}/{maxQueueLength}',
    instancesCount: 'Instances {instances}/{minInstances}',
    activeInstances: 'Seules {instances} sont actives sur {minInstances}. Le nombre optimal d\'instances est de {optimalInstances}.',
    queueOverflowed: 'La file d\'attente est débordée : {queueLength} messages sur {maxQueueLength}.\nVeuillez vérifier les instances.',
    engineDown: '{name} est en panne, le système n\'est pas opérationnel.\nVeuillez vérifier le journal ou redémarrer le service.',
    engineDownOrSlow: '{name} est en panne ou répond trop lentement, le système n\'est pas opérationnel.\nVeuillez vérifier le journal ou redémarrer l\'instance.',
    timescaleDown: '{name} est en panne, les métriques et les KPI ne sont pas collectés.\nVeuillez vérifier le journal ou redémarrer l\'instance.',
    invalidEnginesOrder: 'Configuration des moteurs non valide',
    invalidInstancesConfiguration: 'Configuration des instances non valide : les instances du moteur lisent ou écrivent dans différentes files d\'attente.\nVeuillez vérifier les instances.',
    chainConfigurationInvalid: 'La configuration de la chaîne des moteurs n\'est pas valide.\nReportez-vous ci-dessous pour la séquence correcte des moteurs :',
    queueLimit: 'Limite de longueur de file d\'attente',
    defineQueueLimit: 'Définir la limite de longueur de file d\'attente des moteurs',
    notifyUsersQueueLimit: 'Les utilisateurs peuvent être avertis lorsque la limite de longueur de file d\'attente est dépassée',
    numberOfInstances: 'Nombre d\'instances',
    notifyUsersNumberOfInstances: 'Les utilisateurs peuvent être avertis lorsque le nombre d\'instances actives est inférieur à la valeur minimale. Le nombre optimal d\'instances est affiché lorsque l\'état du moteur n\'est pas disponible.',
    messagesHistory: 'Historique de traitement des messages FIFO',
    messagesLastHour: 'Traitement des messages FIFO pour la dernière heure',
    messagesPerHour: 'messages/heure',
    unknown: 'Cet état du système n\'est pas disponible',
    systemStatusChipError: 'Le système n\'est pas opérationnel',
    systemStatusServerError: 'La configuration du système n\'est pas valide, veuillez contacter l\'administrateur',
    systemsOperational: 'Tous les systèmes sont opérationnels',
    validation: {
      max_value: 'Le champ doit être égal ou inférieur au nombre optimal d\'instances',
      min_value: 'Le champ doit être égal ou supérieur au nombre minimal d\'instances',
    },
    nodes: {
      [HEALTHCHECK_SERVICES_NAMES.mongo]: {
        name: 'MongoDB',
        edgeLabel: 'Vérification de l\'état',
      },

      [HEALTHCHECK_SERVICES_NAMES.rabbit]: {
        name: 'RabbitMQ',
        edgeLabel: 'Vérification de l\'état',
      },

      [HEALTHCHECK_SERVICES_NAMES.redis]: {
        name: 'Redis',
        edgeLabel: 'Données FIFO\nÉtat de Redis',
      },

      [HEALTHCHECK_SERVICES_NAMES.events]: {
        name: 'Événements',
      },

      [HEALTHCHECK_SERVICES_NAMES.api]: {
        name: 'Canopsis API',
      },

      [HEALTHCHECK_SERVICES_NAMES.enginesChain]: {
        name: 'Chaîne des moteurs',
      },

      [HEALTHCHECK_SERVICES_NAMES.healthcheck]: {
        name: 'Healthcheck',
      },

      [HEALTHCHECK_ENGINES_NAMES.webhook]: {
        name: 'Webhook',
        description: 'Gère les webhooks',
      },

      [HEALTHCHECK_ENGINES_NAMES.fifo]: {
        name: 'FIFO',
        edgeLabel: 'État de RabbitMQ\nFlux entrant des KPIs',
        description: 'Gère la file d\'attente des événements et des alarmes',
      },

      [HEALTHCHECK_ENGINES_NAMES.axe]: {
        name: 'AXE',
        description: 'Crée des alarmes et effectue des actions avec elles',
      },

      [HEALTHCHECK_ENGINES_NAMES.che]: {
        name: 'CHE',
        description: 'Applique les filtres d\'événements et les entités créées',
      },

      [HEALTHCHECK_ENGINES_NAMES.pbehavior]: {
        name: 'Pbehavior',
        description: 'Vérifie si l\'alarme est sous PBehavior',
      },

      [HEALTHCHECK_ENGINES_NAMES.action]: {
        name: 'Action',
        description: 'Déclenche le lancement des actions',
      },

      [HEALTHCHECK_ENGINES_NAMES.service]: {
        name: 'Service',
        description: 'Met à jour les compteurs et génère service-events',
      },

      [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: {
        name: 'Infos dynamiques',
        description: 'Ajoute des informations dynamiques à l\'alarme',
      },

      [HEALTHCHECK_ENGINES_NAMES.correlation]: {
        name: 'Corrélation',
        description: 'Gère la corrélation',
      },

      [HEALTHCHECK_ENGINES_NAMES.remediation]: {
        name: 'Remédiation',
        description: 'Déclenche les consignes',
      },
    },
  },

  remediation: {
    tabs: {
      configurations: 'Configurations',
      jobs: 'Tâches',
      statistics: 'Statistiques de remédiation',
    },
  },

  remediationInstructions: {
    usingInstruction: 'Ne peut pas être supprimée, car en cours d\'utilisation',
    addStep: 'Ajouter une étape',
    addOperation: 'Ajouter une opération',
    addJob: 'Ajouter une tâche',
    endpoint: 'Point de terminaison',
    job: 'Tâche | Tâches',
    listJobs: 'Liste des tâches',
    endpointAvatar: 'EP',
    workflow: 'Si cette étape échoue :',
    jobWorkflow: 'Comportement si cette tâche échoue:',
    remainingStep: 'Continuer avec les étapes restantes',
    remainingJob: 'Continuer avec la tâche restante',
    timeToComplete: 'Temps d\'exécution (estimation)',
    emptySteps: 'Aucune étape ajoutée pour le moment',
    emptyOperations: 'Aucune opération ajoutée pour le moment',
    emptyJobs: 'Aucune tâche ajoutée pour le moment',
    timeoutAfterExecution: 'Délai d\'attente après l\'exécution de la consigne',
    requestApproval: 'Demande d\'approbation',
    type: 'Type de consigne',
    approvalPending: 'En attente d\'approbation',
    needApprove: 'Une approbation est nécessaire',
    types: {
      [REMEDIATION_INSTRUCTION_TYPES.manual]: 'Manuel',
      [REMEDIATION_INSTRUCTION_TYPES.auto]: 'Automatique',
    },
    tooltips: {
      endpoint: 'Le point de terminaison doit être une question qui appelle une réponse Oui/Non',
    },
    table: {
      rating: 'Évaluation',
      monthExecutions: '№ d\'exécutions\nce mois-ci',
      lastExecutedOn: 'Dernière exécution le',
    },
    errors: {
      runningInstruction: 'Les changements ne peuvent pas être enregistrés car la consigne est en cours d\'exécution. Voulez-vous stopper l\'exécution de la consigne et ainsi enregistrer les changements ?',
      operationRequired: 'Veuillez ajouter au moins une opération',
      stepRequired: 'Veuillez ajouter au moins une étape',
      jobRequired: 'Veuillez ajouter au moins une tâche',
    },
  },

  remediationJobs: {
    addJobs: 'Ajouter {count} tâche | Ajouter {count} tâches',
    usingJob: 'La tâche ne peut être supprimée, car elle est en cours d\'utilisation',
    table: {
      configuration: 'Configuration',
      jobId: 'Identifiant de la tâche',
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
    queueNumber: '{number} {name} travaux sont dans la file d\'attente',
    popups: {
      success: '{instructionName} a été exécutée avec succès',
      failed: '{instructionName} a échoué. Veuillez faire remonter ce problème',
      connectionError: 'Il y a un problème de connexion. Veuillez cliquer sur le bouton d\'actualisation ou recharger la page.',
      wasAborted: '{instructionName} a été abandonnée',
      wasPaused: 'La consigne {instructionName} sur l\'alarme {alarmName} a été interrompue à {date}. Vous pouvez la reprendre manuellement.',
    },
    jobs: {
      title: 'Tâches attribuées :',
      startedAt: 'Date de déclenchement\n(par Canopsis)',
      launchedAt: 'Date de lancement\n(par l\'ordonnanceur)',
      completedAt: 'Fin de traitement\n(par l\'ordonnanceur)',
      waitAlert: 'L\'ordonnanceur ne répond pas, veuillez contacter votre administrateur',
      skip: 'Ignorer la tâche',
      await: 'Attendre',
      failedReason: 'Raison de l\'échec',
      output: 'Retour',
    },
  },

  remediationInstructionsFilters: {
    button: 'Créer un filtre sur les consignes de remédiation',
    with: 'Avec les consignes sélectionnées',
    without: 'Sans les consignes sélectionnées',
    selectAll: 'Tout sélectionner',
    selectedInstructions: 'Consignes sélectionnées',
    selectedInstructionsHelp: 'Les consignes du type sélectionné sont exclues de la liste',
    chip: {
      with: 'AVEC',
      without: 'SANS',
      all: 'TOUT',
    },
  },

  remediationInstructionStats: {
    alarmsTimeline: 'Chronologie des alarmes',
    alarmId: 'Identifiant de l\'alarme',
    executedOn: 'Exécuté sur',
    lastExecutedOn: 'Dernière exécution le',
    modifiedOn: 'Dernière modification le',
    averageCompletionTime: 'Temps moyen\nd\'achèvement',
    executionCount: 'Nombre\nd\'exécutions',
    totalExecutions: 'Total des exécutions',
    successfulExecutions: 'Exécutions réussies',
    alarmStates: 'Alarmes affectées par l\'état',
    okAlarmStates: 'Nombre de résultats\nÉtats OK',
    notAvailable: 'Indisponible',
    instructionChanged: 'La consigne a été modifiée',
    alarmResolvedDate: 'Date de résolution de l\'alarme',
    showFailedExecutions: 'Afficher les exécutions d\'instructions ayant échoué',
    actions: {
      needRate: 'Notez-le!',
      rate: 'Évaluer',
    },
  },

  remediationPatterns: {
    tabs: {
      pbehaviorTypes: {
        title: 'Types de comportements périodiques',
        fields: {
          activeOnTypes: 'Actif sur les types',
          disabledOnTypes: 'Désactivé sur les types',
        },
      },
    },
  },

  remediationJob: {
    configuration: 'Configuration',
    jobId: 'Identifiant de la tâche',
    query: 'Requête',
    multipleExecutions: 'Autoriser l\'exécution parallèle',
    retryAmount: 'Montant de la nouvelle tentative',
    retryInterval: 'Intervalle de relance',
    addPayload: 'Ajouter un payload',
    deletePayload: 'Supprimer le payload',
    payloadHelp: '<p>Les variables accessibles sont: <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
      + '<i>Quelques exemples:</i>'
      + '<pre>{\n  resource: "{{ .Alarm.Value.Resource }}",\n  entity: "{{ .Entity.ID }}"\n}</pre>',
    errors: {
      invalidJSON: 'JSON non valide',
    },
  },

  remediationStatistic: {
    remediation: 'Remédiation',
    fields: {
      all: 'Tout',
    },
    labels: {
      remediated: 'Corrigé',
      notRemediated: 'Non corrigé',
    },
    tooltips: {
      remediated: '{value} alarmes corrigées',
      assigned: '{value} alarmes avec instructions',
    },
  },

  scenario: {
    triggers: 'Déclencheurs',
    emitTrigger: 'Émettre un déclencheur',
    withAuth: 'Avez-vous besoin de champs d\'authentification ?',
    emptyResponse: 'Réponse vide',
    isRegexp: 'La valeur peut être une expression régulière',
    headerKey: 'Clé d\'en-tête',
    headerValue: 'Valeur d\'en-tête',
    key: 'Clé',
    skipVerify: 'Ne pas vérifier les certificats HTTPS',
    headers: 'En-têtes',
    declareTicket: 'Déclarer un ticket',
    workflow: 'Comportement si cette action ne correspond pas :',
    remainingAction: 'Continuer avec les actions restantes',
    addAction: 'Ajouter une action',
    emptyActions: 'Aucune action ajoutée pour le moment',
    output: 'Format d\'action de sortie',
    forwardAuthor: 'Transmettre l\'auteur à l\'étape suivante',
    urlHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>'
      + '<i>Quelques exemples :</i>'
      + '<pre>"https://exampleurl.com?resource={{ .Alarm.Value.Resource }}"</pre>'
      + '<pre>"https://exampleurl.com?entity_id={{ .Entity.ID }}"</pre>'
      + '<pre>"https://exampleurl.com?children_count={{ len .Children }}"</pre>'
      + '<pre>"https://exampleurl.com?children={{ range .Children }}{{ .ID }}{{ end }}"</pre>',
    outputHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong> et <strong>.Entity</strong></p>'
      + '<i>Quelques exemples:</i>'
      + '<pre>Resource - {{ .Alarm.Value.Resource }}. Entity - {{ .Entity.ID }}.</pre>',
    payloadHelp: '<p>Les variables accessibles sont : <strong>.Alarm</strong>, <strong>.Entity</strong> et <strong>.Children</strong></p>'
      + '<i>Quelques exemples:</i>'
      + '<pre>{\n'
      + '  resource: "{{ .Alarm.Value.Resource }}",\n'
      + '  entity: "{{ .Entity.ID }}",\n'
      + '  children_count: "{{ len .Children }}",\n'
      + '  children: {{ range .Children }}{{ .ID }}{{ end }}\n'
      + '}</pre>',
    actions: {
      [ACTION_TYPES.snooze]: 'Mettre en veille',
      [ACTION_TYPES.pbehavior]: 'Définir un comportement périodique',
      [ACTION_TYPES.changeState]: 'Changer l\'état (Change et vérouille la criticité)',
      [ACTION_TYPES.ack]: 'Acquitter',
      [ACTION_TYPES.ackremove]: 'Supprimer l\'acquittement',
      [ACTION_TYPES.assocticket]: 'Associer un ticket',
      [ACTION_TYPES.cancel]: 'Annuler',
      [ACTION_TYPES.webhook]: 'Webhook',
    },
    tabs: {
      pattern: 'Modèle',
    },
    errors: {
      actionRequired: 'Veuillez ajouter au moins une action',
      priorityExist: 'La priorité du scénario actuel est déjà utilisée. Voulez-vous changer la priorité actuelle du scénario en {priority} ?',
    },
  },

  mixedField: {
    types: {
      [PATTERN_FIELD_TYPES.string]: '@:variableTypes.string',
      [PATTERN_FIELD_TYPES.number]: '@:variableTypes.number',
      [PATTERN_FIELD_TYPES.boolean]: '@:variableTypes.boolean',
      [PATTERN_FIELD_TYPES.null]: '@:variableTypes.null',
      [PATTERN_FIELD_TYPES.stringArray]: '@:variableTypes.array',
    },
  },

  entity: {
    manageInfos: 'Gérer les informations',
    form: 'Formulaire',
    impact: 'Impacts',
    depends: 'Dépendances',
    addInformation: 'Ajouter une information',
    emptyInfos: 'Aucune information',
    availabilityState: 'État de haute disponibilité',
    types: {
      [ENTITY_TYPES.component]: 'Composant',
      [ENTITY_TYPES.connector]: 'Connecteur',
      [ENTITY_TYPES.resource]: 'Ressource',
      [ENTITY_TYPES.service]: 'Service',
    },
  },

  service: {
    outputTemplate: 'Modèle de message',
    createCategory: 'Ajouter une catégorie',
    createCategoryHelp: 'Appuyez sur <kbd>enter</kbd> pour enregistrer',
    availabilityState: '@:entity.availabilityState',
  },

  users: {
    seeProfile: 'Voir le profil',
    selectDefaultView: 'Sélectionner une vue par défaut',
    username: 'Identifiant utilisateur',
    firstName: 'Prénom',
    lastName: 'Nom',
    email: 'Email',
    role: 'Rôle',
    enabled: 'Actif',
    password: 'Mot de passe',
    language: 'Langue par défaut',
    auth: 'Type d\'auth.',
    navigationType: 'Type d\'affichage de la barre de vues',
    navigationTypes: {
      [GROUPS_NAVIGATION_TYPES.sideBar]: 'Barre latérale',
      [GROUPS_NAVIGATION_TYPES.topBar]: 'Barre d\'entête',
    },
    metrics: {
      [USER_METRIC_PARAMETERS.totalUserActivity]: 'Durée totale de l\'activité',
    },
  },

  testSuite: {
    xmlFeed: 'Flux XML',
    hostname: 'Nom d\'hôte',
    lastUpdate: 'Dernière mise à jour',
    timeTaken: 'Temps passé',
    totalTests: 'Total des tests',
    disabledTests: 'Tests désactivés',
    copyMessage: 'Copier le message système',
    systemError: 'Erreur système',
    systemErrorMessage: 'Message d\'erreur système',
    systemOut: 'Retour système',
    systemOutMessage: 'Message de retour du système',
    compareWithHistorical: 'Comparer avec les données historiques',
    className: 'Nom du test',
    line: 'Ligne',
    failureMessage: 'Message d\'échec',
    noData: 'Aucun message système trouvé dans le formulaire XML',
    tabs: {
      globalMessages: 'Messages globaux',
      gantt: 'Gantt',
      details: 'Détails',
      screenshots: 'Captures d\'écran',
      videos: 'Vidéos',
    },
    statuses: {
      [TEST_SUITE_STATUSES.passed]: 'Passé',
      [TEST_SUITE_STATUSES.skipped]: 'Ignoré',
      [TEST_SUITE_STATUSES.error]: 'En erreur',
      [TEST_SUITE_STATUSES.failed]: 'Échoué',
      [TEST_SUITE_STATUSES.total]: 'Temps total pris',
    },
    popups: {
      systemMessageCopied: 'Message système copié dans le presse-papier',
    },
  },

  stateSetting: {
    worstLabel: 'Le pire critère :',
    worstHelpText: 'Canopsis compte l\'état pour chaque critère défini. L\'état final de la suite de tests JUnit est considéré comme le pire des états résultants.',
    criterion: 'Critère',
    serviceState: 'État du service',
    methods: {
      [STATE_SETTING_METHODS.worst]: 'Pire',
      [STATE_SETTING_METHODS.worstOfShare]: 'Pire des états',
    },
    states: {
      minor: 'Mineur',
      major: 'Majeur',
      critical: 'Critique',
    },
  },

  storageSettings: {
    alarm: {
      title: 'Stockage des données d\'alarme',
      titleHelp: 'Lorsque ces options sont activées, les données d\'alarmes résolues sont archivées et/ou supprimées après la période de temps définie.',
      archiveAfter: 'Archiver les données d\'alarmes résolues après',
      deleteAfter: 'Supprimer les données d\'alarmes résolues après',
    },
    junit: {
      title: 'Stockage de données JUnit',
      deleteAfter: 'Supprimer les données des suites de tests après',
      deleteAfterHelpText: 'Lorsque cette option est activée, les données des suites de tests JUnit (XML, captures d\'écran et vidéos) sont supprimées après la période définie.',
    },
    remediation: {
      title: 'Stockage des données de consigne',
      accumulateAfter: 'Accumuler les statistiques des consignes après',
      deleteAfter: 'Supprimer les données des consignes après',
      deleteAfterHelpText: 'Lorsque cette option est activée, les données statistiques des consignes sont supprimées après la période de temps définie.',
    },
    entity: {
      title: 'Stockage des données des entités',
      titleHelp: 'Toutes les entités désactivées avec des alarmes associées peuvent être archivées (déplacées dans la collection séparée) et/ou supprimées pour toujours.',
      archiveEntity: 'Archiver les entités désactivées',
      deleteEntity: 'Supprimer définitivement les entités désactivées de l\'archive',
      archiveDependencies: 'Supprimer également les entités impactantes et dépendantes',
      archiveDependenciesHelp: 'Pour les connecteurs, tous les composants et toutes les ressources impactants et dépendants seront archivés ou supprimés pour toujours. Pour les composants, toutes les ressources dépendantes seront également archivées ou supprimées pour toujours.',
      cleanStorage: 'Archiver ou Supprimer les entités désactivées',
    },
    pbehavior: {
      title: 'Stockage des données de comportements périodiques',
      deleteAfter: 'Supprimer les données de comportements périodiques après',
      deleteAfterHelpText: 'Lorsque cette option est activée, les comportements périodiques inactifs sont supprimés après la période de temps définie à partir du dernier événement.',
    },
    healthCheck: {
      title: 'Stockage des données du bilan de santé',
      deleteAfter: 'Supprimer les données de flux entrant FIFO après',
    },
    history: {
      scriptLaunched: 'Script lancé à {launchedAt}.',
      alarm: {
        deletedCount: 'Alarmes supprimées : {count}.',
        archivedCount: 'Alarmes archivées : {count}.',
      },
      entity: {
        deletedCount: 'Entités supprimées : {count}.',
        archivedCount: 'Entités archivées : {count}.',
      },
    },
  },

  notificationSettings: {
    instruction: {
      header: 'Instructions',
      rate: 'Notifications "Évaluer les consignes"',
      rateFrequency: 'La fréquence',
      duration: 'Intervalle de temps',
    },
  },

  quickRanges: {
    title: 'Valeurs usuelles',
    timeField: 'Champ de temps',
    types: {
      [QUICK_RANGES.custom.value]: 'Personnalisé',
      [QUICK_RANGES.last15Minutes.value]: '15 dernières minutes',
      [QUICK_RANGES.last30Minutes.value]: '30 dernières minutes',
      [QUICK_RANGES.last1Hour.value]: 'Dernière heure',
      [QUICK_RANGES.last3Hour.value]: '3 dernières heures',
      [QUICK_RANGES.last6Hour.value]: '6 dernières heures',
      [QUICK_RANGES.last12Hour.value]: '12 dernières heures',
      [QUICK_RANGES.last24Hour.value]: '24 dernières heures',
      [QUICK_RANGES.last2Days.value]: '2 derniers jours',
      [QUICK_RANGES.last7Days.value]: '7 derniers jours',
      [QUICK_RANGES.last30Days.value]: '30 derniers jours',
      [QUICK_RANGES.last1Year.value]: 'Dernière année',
      [QUICK_RANGES.yesterday.value]: 'Hier',
      [QUICK_RANGES.previousWeek.value]: 'Dernière semaine',
      [QUICK_RANGES.previousMonth.value]: 'Dernier mois',
      [QUICK_RANGES.today.value]: 'Aujourd\'hui',
      [QUICK_RANGES.todaySoFar.value]: 'Aujourd\'hui jusqu\'à maintenant',
      [QUICK_RANGES.thisWeek.value]: 'Cette semaine',
      [QUICK_RANGES.thisWeekSoFar.value]: 'Cette semaine jusqu\'à maintenant',
      [QUICK_RANGES.thisMonth.value]: 'Ce mois',
      [QUICK_RANGES.thisMonthSoFar.value]: 'Ce mois jusqu\'à maintenant',
    },
  },

  idleRules: {
    timeAwaiting: 'Temps d\'attente',
    timeRangeAwaiting: 'Plage de temps en attente',
    types: {
      [IDLE_RULE_TYPES.alarm]: 'Règle d\'alarme',
      [IDLE_RULE_TYPES.entity]: 'Règle d\'entité',
    },
    alarmConditions: {
      [IDLE_RULE_ALARM_CONDITIONS.lastEvent]: 'Aucun événement reçu',
      [IDLE_RULE_ALARM_CONDITIONS.lastUpdate]: 'Aucun changement d\'état',
    },
  },

  alarmStatusRules: {
    frequencyLimit: 'Nombre d\'oscillations',
  },

  icons: {
    noEvents: 'Aucun événement reçu pendant {duration} par certaines dépendances',
  },

  pageHeaders: {
    hideMessage: 'J\'ai compris! Cacher',
    learnMore: 'En savoir plus sur {link}',

    /**
     * Exploitation
     */
    [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
      title: 'Filtres d\'événements',
      message: 'Le filtre d\'événements est une fonctionnalité du moteur CHE qui permet la définition de règles de traitement des événements.',
    },

    [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      title: 'Informations dynamiques',
      message: 'Les informations dynamiques de Canopsis sont utilisées pour enrichire les alarmes. Ces enrichissements sont définis par des règles indiquant dans quelles conditions les informations doivent être présentées sur une alarme.',
    },

    [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      title: 'Règles de méta-alarme',
      message: 'Les règles de méta-alarme peuvent être utilisées pour grouper les alarmes par types et critères (relation parent/enfant, intervalle de temps, etc).',
    },

    [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
      title: 'Règles d\'inactivité',
      message: 'Les règles d\'inactivité des entités ou des alarmes peuvent être utilisées pour surveiller les événements et les états d\'alarme afin de savoir si des événements ne sont pas reçus ou si l\'état d\'alarme n\'est pas modifié pendant une longue période en raison d\'erreurs ou d\'une configuration non valide.',
    },

    [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
      title: 'Règles de bagot',
      // message: '', // TODO: need to put description
    },

    [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
      title: 'Règles de résolution',
      // message: '', // TODO: need to put description
    },

    [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
      title: 'Comportements périodiques',
      message: 'Les comportements périodiques de Canopsis peuvent être utilisés afin de définir une période pendant laquelle le comportement doit être modifié, par exemple pour la maintenance ou le type de services.',
    },

    [USERS_PERMISSIONS.technical.exploitation.scenario]: {
      title: 'Scénarios',
      message: 'Les scénarios Canopsis peuvent être utilisés pour déclencher conditionnellement divers types d\'actions sur les alarmes.',
    },

    [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
      title: 'Règles SNMP',
      message: 'Les règles SNMP peuvent être utilisées pour traiter les traps SNMP remontées par le connecteur snmp2canopsis au travers du moteur SNMP.',
    },

    /**
     * Admin access
     */
    [USERS_PERMISSIONS.technical.permission]: {
      title: 'Droits',
    },
    [USERS_PERMISSIONS.technical.role]: {
      title: 'Rôles',
    },
    [USERS_PERMISSIONS.technical.user]: {
      title: 'Utilisateurs',
    },

    /**
     * Admin communications
     */
    [USERS_PERMISSIONS.technical.broadcastMessage]: {
      title: 'Diffusion de messages',
      message: 'La diffusion de messages peut être utilisée pour afficher les bannières et les messages d\'information qui apparaîtront dans l\'interface de Canopsis.',
    },
    [USERS_PERMISSIONS.technical.playlist]: {
      title: 'Listes de lecture',
      message: 'Les listes de lecture peuvent être utilisées pour la personnalisation des vues qui peuvent être affichées les unes après les autres avec un délai associé.',
    },
    [USERS_PERMISSIONS.technical.healthcheck]: {
      title: 'Bilan de santé',
      message: 'La fonction Healthcheck est le tableau de bord avec des indications d\'états et d\'erreurs de tous les systèmes inclus dans Canopsis.',
    },
    [USERS_PERMISSIONS.technical.engine]: {
      title: 'Engines',
      message: 'This page contains the information about the sequence and configuration of engines. To work properly, the chain of engines must be continuous.',
    },
    [USERS_PERMISSIONS.technical.kpi]: {
      title: 'KPI',
      message: '', // TODO: add correct message
    },

    /**
     * Admin general
     */
    [USERS_PERMISSIONS.technical.parameters]: {
      title: 'Paramètres',
    },
    [USERS_PERMISSIONS.technical.planning]: {
      title: 'Planification',
      message: 'La fonctionnalité d\'administration de la planification de Canopsis peut être utilisée pour la personnalisation des types de comportements périodiques.',
    },
    [USERS_PERMISSIONS.technical.remediation]: {
      title: 'Consignes',
      message: 'La fonction de remédiation de Canopsis peut être utilisée pour créer des plans ou des consignes visant à corriger des situations.',
    },

    /**
     * Notifications
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      title: 'Évaluation des consignes',
      message: 'Cette page contient les statistiques sur l\'exécution des consignes. Les utilisateurs peuvent noter les consignes en fonction de leurs performances.',
    },
  },

  userInterface: {
    title: 'Interface utilisateur',
    appTitle: 'Titre de l\'application',
    language: 'Langue par défaut',
    footer: 'Page d\'identification : pied de page',
    description: 'Page d\'identification : description',
    logo: 'Logo',
    infoPopupTimeout: 'Délai d\'affichage pour les popups d\'informations',
    errorPopupTimeout: 'Délai d\'affichage pour les popups d\'erreurs',
    allowChangeSeverityToInfo: 'Autorise le changement de criticité en Info',
    maxMatchedItems: 'Seuil d\'éléments avant avertissement',
    checkCountRequestTimeout: 'Délai d\'expiration de la requête',
    tooltips: {
      maxMatchedItems: 'Avertit l\'utilisateur lorsque le nombre d\'éléments correspondant aux modèles est supérieur à cette valeur',
      checkCountRequestTimeout: 'Définit le délai d\'expiration (en secondes) de la demande d\'éléments correspondants',
    },
  },

  kpi: {
    alarmMetrics: 'Métriques d\'alarme',
    sli: 'SLI',
    metricsNotAvailable: 'TimescaleDB ne fonctionne pas. Les métriques ne sont pas disponibles.',
    noData: 'Pas de données disponibles',
  },

  kpiMetrics: {
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
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: '{value} alarmes actives avec acks',
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

  kpiFilters: {
    helpInformation: 'Ici, les modèles de filtre pour des tranches de données supplémentaires pour les compteurs et les évaluations peuvent être ajoutés.',
  },

  kpiRatingSettings: {
    helpInformation: 'La liste des paramètres à utiliser pour la notation.',
  },

  snmpRule: {
    oid: 'OID',
    module: 'Sélectionnez un module MIB',
    output: 'Message',
    resource: 'Ressource',
    component: 'Composant',
    connectorName: 'Nom du connecteur',
    state: 'Criticité',
    toCustom: 'Personnaliser',
    writeTemplate: 'Écrire un modèle',
    defineVar: 'Définir la variable SNMP correspondante',
    moduleMibObjects: 'Champ d\'association des variables SNMP',
    regex: 'Expression régulière',
    formatter: 'Format (groupe de capture avec \\x)',
    uploadMib: 'Envoyer un fichier MIB',
    addSnmpRule: 'Ajouter une règle SNMP',
  },

  pattern: {
    patterns: 'Modèles',
    myPatterns: 'Mes modèles',
    corporatePatterns: 'Modèles partagés',
    addRule: 'Ajouter une règle',
    addGroup: 'Ajouter un groupe',
    removeRule: 'Supprimer la règle',
    advancedEditor: 'Éditeur avancé',
    simpleEditor: 'Éditeur simple',
    noData: 'Aucun modèle. Cliquez sur \'@:pattern.addGroup\' pour ajouter des champs au modèle',
    noDataDisabled: 'Aucun modèle.',
    discard: 'Jeter le motif',
    types: {
      [PATTERN_TYPES.alarm]: 'Modèle d\'alarme',
      [PATTERN_TYPES.entity]: 'Modèle d\'entité',
      [PATTERN_TYPES.pbehavior]: 'Modèle de comportement',
    },
    errors: {
      ruleRequired: 'Veuillez ajouter au moins une règle',
      groupRequired: 'Veuillez ajouter au moins un groupe',
      invalidPatterns: 'Les modèles ne sont pas valides ou il y a un champ de modèle désactivé',
      countOverLimit: 'Le modèle que vous avez défini cible {count} éléments. Cela peut affecter les performances, en êtes-vous sûr ?',
      oldPattern: 'Le modèle de filtre actuel est défini dans l\'ancien format. Veuillez utiliser l\'éditeur avancé pour l\'afficher. Les filtres dans l\'ancien format seront bientôt obsolètes. Veuillez créer de nouveaux modèles dans notre interface mise à jour.',
      existExcluded: 'Les règles incluent la règle exclue.',
    },
  },

  filter: {
    oldPattern: 'Ancien format de motif',
  },
}, featureService.get('i18n.fr'));
