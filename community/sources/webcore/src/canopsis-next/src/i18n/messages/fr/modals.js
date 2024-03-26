import { PBEHAVIOR_TYPE_TYPES, WIDGET_TYPES } from '@/constants';

export default {
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
      privateTitle: 'Créer une vue privée',
    },
    edit: {
      title: 'Éditer une vue',
    },
    duplicate: {
      title: 'Dupliquer une vue - {viewTitle}',
      infoMessage: 'Vous êtes en train de dupliquer une vue. Toutes les lignes et les widgets de la vue dupliquée seront copiés dans la nouvelle vue.',
    },
    success: {
      create: 'Nouvelle vue créée !',
      edit: 'Vue éditée avec succès !',
      duplicate: 'Vue dupliquée avec succès !',
      delete: 'Vue supprimée avec succès !',
    },
    fail: {
      create: 'Erreur lors de la création de la vue...',
      edit: 'Erreur lors de l\'édition de la vue...',
      duplicate: 'Échec de la duplication de la vue...',
      delete: 'Erreur lors de la suppression de la vue...',
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
      color: {
        label: 'Utiliser une couleur spéciale pour cet événement ?',
      },
    },
    errors: {
      invalid: 'Invalide',
    },
    success: {
      create: 'Comportement périodique créé avec succès !',
    },
    cancelConfirmation: 'Certaines informations ont été modifiées et ne seront pas sauvegardées. Voulez-vous vraiment quitter ce menu ?',
  },
  createPause: {
    title: 'Mettre en pause',
  },
  createAckRemove: {
    title: 'Annuler l\'acquittement',
  },
  createUnCancel: {
    title: 'Annuler l\'annulation',
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
  },
  service: {
    refreshEntities: 'Rafraîchir la liste des entités',
    editPbehaviors: 'Éditer les comportements périodiques',
    massActionsDescription: 'Vous pouvez choisir des entités pour effectuer des actions',
    actionInQueue: 'action en file d\'attente|actions en file d\'attente',
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
      [WIDGET_TYPES.map]: {
        title: 'Cartographie',
      },
      [WIDGET_TYPES.barChart]: {
        title: 'Histogramme',
      },
      [WIDGET_TYPES.lineChart]: {
        title: 'Graphique en ligne',
      },
      [WIDGET_TYPES.pieChart]: {
        title: 'Diagramme circulaire',
      },
      [WIDGET_TYPES.numbers]: {
        title: 'Nombres',
      },
      [WIDGET_TYPES.userStatistics]: {
        title: 'Statistiques des utilisateurs',
      },
      [WIDGET_TYPES.alarmStatistics]: {
        title: 'Statistiques des alarmes',
      },
      [WIDGET_TYPES.availability]: {
        title: 'Disponibilité',
      },
      chart: {
        title: 'Graphique',
      },
      report: {
        title: 'Rapport',
      },
    },
  },
  manageHistogramGroups: {
    title: {
      add: 'Ajouter un groupe',
      edit: 'Éditer un groupe',
    },
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
    duplicate: {
      title: 'Dupliquer un rôle',
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
    duplicate: {
      title: 'Dupliquer la règle SNMP',
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
      emptyTabs: 'Merci d\'ajouter un onglet',
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
  createRrule: {
    title: 'Créer une règle de récurrence',
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
    canonicalTypes: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Actif',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactif',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
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
      title: 'Éditer la consigne',
      popups: {
        success: '{instructionName} a été modifiée avec succès',
      },
    },
    duplicate: {
      title: 'Dupliquer la consigne',
      popups: {
        success: '{instructionName} a été dupliquée avec succès',
      },
    },
  },
  createRemediationConfiguration: {
    create: {
      title: 'Créer une configuration',
      popups: {
        success: '{configurationName} a été créée avec succès',
      },
    },
    edit: {
      title: 'Modifier la configuration',
      popups: {
        success: '{configurationName} a été modifiée avec succès',
      },
    },
    duplicate: {
      title: 'Dupliquer la configuration',
      popups: {
        success: '{configurationName} a été dupliquée avec succès',
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
    duplicate: {
      title: 'Dupliquer une tâche',
      popups: {
        success: '{jobName} a été dupliquée avec succès',
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
  createStateSetting: {
    create: {
      title: 'Créer une méthode de calcul d\'état',
      success: 'Méthode de calcul d\'état créée !',
    },
    edit: {
      title: 'Modifier la méthode de calcul de l\'état',
      success: 'Méthode de calcul d\'état modifiée !',
    },
    duplicate: {
      title: 'Méthode de calcul d\'état en double',
      success: 'Méthode de calcul d\'état dupliquée !',
    },
    remove: {
      success: 'Méthode de calcul d\'état supprimée !',
    },
  },
  createJunitStateSetting: {
    edit: {
      title: 'Paramètres d\'état de la suite de tests JUnit',
      success: 'Paramètre d\'état de la suite de tests JUnit modifié !',
    },
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
    dismissed: 'a rejeté vos mises à jour',
    requested: 'vous a sollicité pour une pour approbation',
    tabs: {
      updated: 'Mise à jour',
      original: 'Original',
    },
  },
  createAlarmIdleRule: {
    create: {
      title: 'Créer une règle d\'inactivité d\'alarme',
    },
    edit: {
      title: 'Modifier la règle d\'inactivité d\'alarme',
    },
    duplicate: {
      title: 'Dupliquer une règle d\'inactivité d\'alarme',
    },
  },
  createEntityIdleRule: {
    create: {
      title: 'Créer une règle d\'inactivité d\'entité',
    },
    edit: {
      title: 'Modifier la règle d\'inactivité d\'entité',
    },
    duplicate: {
      title: 'Dupliquer une règle d\'inactivité d\'entité',
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
      title: 'Archivage/Suppression des entités. Êtes vous sur ?',
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
      title: 'Modifier le filtre d\'alarme',
    },
  },
  createCorporateAlarmPattern: {
    create: {
      title: 'Créer un filtre partagé d\'alarme',
    },
    edit: {
      title: 'Modifier le filtre partagé d\'alarme',
    },
  },
  createEntityPattern: {
    create: {
      title: 'Créer un filtre d\'entité',
    },
    edit: {
      title: 'Modifier le filtre d\'entité',
    },
  },
  createCorporateEntityPattern: {
    create: {
      title: 'Créer un filtre partagé d\'entité',
    },
    edit: {
      title: 'Modifier le filtre partagé d\'entité',
    },
  },
  createPbehaviorPattern: {
    create: {
      title: 'Créer un filtre de comportement périodique',
    },
    edit: {
      title: 'Modifier le filtre de comportement périodique',
    },
  },
  createCorporatePbehaviorPattern: {
    create: {
      title: 'Créer un filtre partagé de comportement périodique',
    },
    edit: {
      title: 'Modifier le filtre partagé de comportement périodique',
    },
  },
  createMap: {
    title: 'Créer une carte',
  },
  createGeoMap: {
    create: {
      title: 'Créer une carte géographique',
    },
    edit: {
      title: 'Modifier une carte géographique',
    },
    duplicate: {
      title: 'Dupliquer une carte géographique',
    },
  },
  createFlowchartMap: {
    create: {
      title: 'Créer une carte flowchart',
    },
    edit: {
      title: 'Modifier une carte flowchart',
    },
    duplicate: {
      title: 'Dupliquer une carte flowchart',
    },
  },
  createMermaidMap: {
    create: {
      title: 'Créer un diagramme mermaid',
    },
    edit: {
      title: 'Modifier un diagramme mermaid',
    },
    duplicate: {
      title: 'Dupliquer un diagramme mermaid',
    },
  },
  createTreeOfDependenciesMap: {
    create: {
      title: 'Créer une carte d\'arbre de dépendances',
    },
    edit: {
      title: 'Modifier une carte d\'arbre de dépendances',
    },
    duplicate: {
      title: 'Dupliquer une carte d\'arbre de dépendances',
    },
    addEntity: 'Ajouter une entité',
    pinnedEntities: 'Entités épinglées',
  },
  createShareToken: {
    create: {
      title: 'Créer un jeton de partage',
    },
  },
  createWidgetTemplate: {
    create: {
      title: 'Créer un modèle de widget',
    },
    edit: {
      title: 'Modifier le modèle de widget',
    },
  },
  selectWidgetTemplateType: {
    title: 'Sélectionner le type de modèle de widget',
  },
  entityDependenciesList: {
    title: 'Diagramme de cause racine',
  },
  createDeclareTicketRule: {
    create: {
      title: 'Créer une règle de déclaration de ticket',
    },
    edit: {
      title: 'Modifier une règle de déclaration de ticket',
    },
    duplicate: {
      title: 'Dupliquer une règle de déclaration de ticket',
    },
  },
  createDeclareTicketEvent: {
    title: 'Déclarer un incident',
  },
  executeDeclareTickets: {
    title: 'Statut de la déclaration du ticket',
  },
  createAssociateTicketEvent: {
    title: 'Associer un numéro de ticket',
  },
  createAckEvent: {
    title: 'Acquitter',
  },
  createLinkRule: {
    create: {
      title: 'Créer un générateur de liens',
    },
    edit: {
      title: 'Modifier le générateur de liens',
    },
  },
  createAlarmChart: {
    [WIDGET_TYPES.barChart]: {
      create: {
        title: 'Créer un graphique à barres',
      },
      edit: {
        title: 'Modifier le graphique à barres',
      },
    },
    [WIDGET_TYPES.lineChart]: {
      create: {
        title: 'Créer un graphique en courbes',
      },
      edit: {
        title: 'Modifier le graphique en courbes',
      },
    },
    [WIDGET_TYPES.numbers]: {
      create: {
        title: 'Créer un tableau de nombres',
      },
      edit: {
        title: 'Modifier le tableau des nombres',
      },
    },
  },
  importPbehaviorException: {
    title: 'Importer des dates d\'exception',
  },
  createMaintenance: {
    enableMaintenance: 'Activer le mode maintenance',
    setup: {
      title: 'Activer le mode maintenance',
    },
    edit: {
      title: 'Modifier le mode maintenance',
    },
  },
  confirmationLeaveMaintenance: {
    title: 'Quitter le mode maintenance',
    text: 'Êtes-vous sûr de vouloir quitter le mode maintenance ?\nTous les utilisateurs pourront à nouveau se connecter à Canopsis.',
  },
  createTag: {
    create: {
      title: 'Créer un tag',
    },
    edit: {
      title: 'Modifier un tag',
    },
    duplicate: {
      title: 'Dupliquer un tag',
    },
  },
  createTheme: {
    create: {
      title: 'Créer un thème',
    },
    edit: {
      title: 'Modifier le thème',
    },
    duplicate: {
      title: 'Thème en double',
    },
  },
  archiveDisabledEntities: {
    text: 'Êtes-vous sûr de vouloir archiver les entités désactivées ?\nCette action ne peut pas être annulée',
  },
  createIcon: {
    create: {
      title: 'Icône de téléchargement',
      success: 'L\'icône a été téléchargée',
    },
    remove: {
      success: 'L\'icône a été supprimée',
    },
  },
};
