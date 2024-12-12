import { USER_PERMISSIONS_PREFIXES, USER_PERMISSIONS } from '@/constants';

export default {
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
    [USER_PERMISSIONS_PREFIXES.business.map]: 'Droits pour le widget : Cartographie',
    [USER_PERMISSIONS_PREFIXES.business.barChart]: 'Droits pour le widget : Histogramme',
    [USER_PERMISSIONS_PREFIXES.business.lineChart]: 'Droits pour le widget : Graphique en ligne',
    [USER_PERMISSIONS_PREFIXES.business.pieChart]: 'Droits pour le widget : Diagramme circulaire',
    [USER_PERMISSIONS_PREFIXES.business.numbers]: 'Droits pour le widget : Nombres',
    [USER_PERMISSIONS_PREFIXES.business.userStatistics]: 'Droits pour le widget : Statistiques des utilisateurs',
    [USER_PERMISSIONS_PREFIXES.business.alarmStatistics]: 'Droits pour le widget : Statistiques des alarmes',
    [USER_PERMISSIONS_PREFIXES.business.availability]: 'Droits pour le widget : Disponibilité',
  },
  api: {
    general: 'Général',
    rules: 'Règles',
    remediation: 'Remédiation',
    pbehavior: 'Comportements périodiques',
    eventsRecord: 'Enregistrement des événements',
  },
  permissions: {
    /**
     * Business Common Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.variablesHelp]: {
      name: 'Accès à la liste des variables disponibles',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des variables dans la liste des alarmes et la météo du service',
    },
    [USER_PERMISSIONS.business.context.actions.entityCommentsList]: {
      name: 'Accès à la liste des commentaires des entités',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des commentaires de l\'entité',
    },
    [USER_PERMISSIONS.business.context.actions.createEntityComment]: {
      name: 'Accès à la création de commentaires d\'entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent créer des commentaires d\'entité',
    },
    [USER_PERMISSIONS.business.context.actions.editEntityComment]: {
      name: 'Accès à l\'édition des commentaires de l\'entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les commentaires de l\'entité',
    },

    /**
     * Business Alarms List Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.ack]: {
      name: 'Droits sur le bac à alarmes : ack',
      description: 'Les utilisateurs disposant de cette autorisation peuvent acquitter les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.fastAck]: {
      name: 'Droits sur le bac à alarmes : acquittement rapide',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer un acquittement rapide des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.ackRemove]: {
      name: 'Droits sur le bac à alarmes : annuler ack',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler un acquittement',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: {
      name: 'Droits sur le bac à alarmes : ajouter un comportement périodique',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Comportement périodique" et modifier les comportements périodiques pour les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.fastPbehaviorAdd]: {
      name: 'Droits sur la liste d\'alarmes : Ajouter un comportement périodique rapidemment',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Définir un comportement périodique rapidemment"',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.snooze]: {
      name: 'Droits sur le bac à alarmes : mettre en veille une alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent mettre en veille les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.declareTicket]: {
      name: 'Droits sur liste d\'alarmes : déclarer un ticket',
      description: 'Les utilisateurs avec cette permission peuvent déclarer des tickets',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.associateTicket]: {
      name: 'Droits sur liste d\'alarmes : associer un ticket',
      description: 'Les utilisateurs disposant de cette autorisation peuvent associer un ticket à une alarme',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.cancel]: {
      name: 'Droits sur le bac à alarmes : annuler une alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.fastCancel]: {
      name: 'Droits sur le bac à alarmes : annulation rapide d\'une alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer une annulation rapide des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.changeState]: {
      name: 'Droits sur le bac à alarmes : modifier la criticité d\'une alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les criticités des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.history]: {
      name: 'Droits sur le bac à alarmes : historique',
      description: 'Les utilisateurs disposant de cette autorisation peuvent afficher l\'historique des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup]: {
      name: 'Droits sur le bac à alarmes : actions manuelles sur les méta-alarmes',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer des méta-alarmes manuelles',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.comment]: {
      name: 'Droits sur le bac à alarmes : commenter une alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent commenter les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.listFilters]: {
      name: 'Droits sur le bac à alarmes : afficher les filtres d\'alarmes',
      description: 'Les utilisateurs disposant de cette autorisation peuvent afficher la liste des filtres disponibles dans la liste des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.editFilter]: {
      name: 'Droits sur le bac à alarmes : modifier les filtres d\'alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent éditer les filtres disponibles dans la liste des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.addFilter]: {
      name: 'Droits sur le bac à alarmes : ajouter des filtres d\'alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres dans la liste des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: {
      name: 'Droits sur le bac à alarmes : afficher les filtres d\'alarme',
      description: 'Le filtre d\'alarme est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters]: { // TODO: change to _remediationInstructionsFilter
      name: 'Droits sur le bac à alarmes : accès aux filtres de remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent lister les filtres de remédiation',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter]: { // TODO: change to _remediationInstructionsFilter
      name: 'Droits sur le bac à alarmes : accès à l\'édition des filtres de remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres de remédiation',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter]: { // TODO: change to _remediationInstructionsFilter
      name: 'Droits sur le bac à alarmes : accès à l\'ajout de filtres de remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres de remédiation',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: {
      name: 'Droits sur le bac à alarmes : Accès aux filtres de remédiation',
      description: 'Le filtre par instructions est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.links]: {
      name: 'Droits sur le bac à alarmes : Accès aux liens',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder et suivre les liens dans la liste des alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.correlation]: {
      name: 'Droits sur le bac à alarmes : Accès au regroupement des alarmes corrélées',
      description: 'Les utilisateurs disposant de cette autorisation peuvent activer le regroupement des alarmes corrélées',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.executeInstruction]: {
      name: 'Droits sur liste des alarmes : Accès aux exécutions de remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des remédisation pour corriger les alarmes',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.category]: {
      name: 'Droits sur le bac à alarmes : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer la liste des alarmes par catégorie',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: {
      name: 'Droits sur le bac à alarmes : Accès à l\'export des alarmes au format CSV',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exporter des alarmes au format CSV',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.addBookmark]: {
      name: 'Droits sur la liste des alarmes : Accès à l\'ajout de favoris aux alarmes',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter un signet à l\'alarme',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.removeBookmark]: {
      name: 'Droits sur la liste des alarmes : Accès à la suppression du signet de l\'alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent supprimer le signet de l\'alarme',
    },
    [USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: {
      name: 'Droits sur la liste des alarmes : Accès au filtrage des alarmes uniquement par favoris',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les alarmes uniquement en fonction des signets',
    },

    /**
     * Business Context Explorer Permissions
     */
    [USER_PERMISSIONS.business.context.actions.createEntity]: {
      name: 'Droits sur l\'explorateur de contexte : créer une entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent créer de nouvelles entités',
    },
    [USER_PERMISSIONS.business.context.actions.editEntity]: {
      name: 'Droits sur l\'explorateur de contexte : modifier l\'entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les entités',
    },
    [USER_PERMISSIONS.business.context.actions.duplicateEntity]: {
      name: 'Droits sur l\'explorateur de contexte : dupliquer une entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent dupliquer des entités',
    },
    [USER_PERMISSIONS.business.context.actions.deleteEntity]: {
      name: 'Droits sur l\'explorateur de contexte : supprimer l\'entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent supprimer des entités',
    },
    [USER_PERMISSIONS.business.context.actions.pbehaviorAdd]: {
      name: 'Droits sur l\'explorateur de contexte : action ajouter un comportement périodique',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Comportement périodique" et modifier les comportements PB pour les entités',
    },
    [USER_PERMISSIONS.business.context.actions.massEnable]: {
      name: 'Droits sur l\'explorateur de contexte : action d\'activation en masse',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer une action en masse pour activer les entités sélectionnées',
    },
    [USER_PERMISSIONS.business.context.actions.massDisable]: {
      name: 'Droits sur l\'explorateur de contexte : action de désactivation en masse',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer une action en masse pour désactiver les entités sélectionnées',
    },
    [USER_PERMISSIONS.business.context.actions.listFilters]: {
      name: 'Droits sur l\'explorateur de contexte : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles dans l\'explorateur de contexte',
    },
    [USER_PERMISSIONS.business.context.actions.editFilter]: {
      name: 'Droits sur l\'explorateur de contexte : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres d\'entité',
    },
    [USER_PERMISSIONS.business.context.actions.addFilter]: {
      name: 'Droits sur l\'explorateur de contexte : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres sur les entités affichées dans l\'explorateur de contexte',
    },
    [USER_PERMISSIONS.business.context.actions.userFilter]: {
      name: 'Droits sur l\'explorateur de contexte : afficher les filtres',
      description: 'Le filtre d\'entité est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.context.actions.category]: {
      name: 'Droits sur l\'explorateur de contexte : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les entités par catégorie',
    },
    [USER_PERMISSIONS.business.context.actions.exportAsCsv]: {
      name: 'Droits sur l\'explorateur de contexte : Exporter au format csv',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exporter des entités sous forme de fichier CSV',
    },

    /**
     * Business Service Weather Permissions
     */
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAck]: {
      name: 'Droits sur la météo des services : Acquitter',
      description: 'Les utilisateurs disposant de cette autorisation peuvent acquitter les alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket]: {
      name: 'Droits sur la météo des services : Associer un ticket',
      description: 'Les utilisateurs disposant de cette autorisation peuvent associer des tickets aux alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityDeclareTicket]: {
      name: 'Droits sur la météo des services : Déclarer un ticket',
      description: 'Les utilisateurs disposant de cette autorisation peuvent déclarer des tickets aux alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityComment]: {
      name: 'Droits sur la météo des services : Commenter',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des commentaires',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityValidate]: {
      name: 'Droits sur la météo des services : Valider (pouce en l\'air)',
      description: 'Les utilisateurs disposant de cette autorisation peuvent valider les alarmes et modifier leur état en critique',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: {
      name: 'Droits sur la météo des services : Invalider (pouce en bas)',
      description: 'Les utilisateurs disposant de cette autorisation peuvent invalider les alarmes et les annuler',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPause]: {
      name: 'Droits sur la météo des services : Mettre en pause',
      description: 'Les utilisateurs disposant de cette autorisation peuvent suspendre les alarmes (définir le type de PBehavior "Pause")',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPlay]: {
      name: 'Droits sur la météo des services : Supprimer une pause',
      description: 'Les utilisateurs disposant de cette autorisation peuvent activer les alarmes en pause (supprimer le type PBehavior "Pause")',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityCancel]: {
      name: 'Droits sur la météo des services : Annuler',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler les alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors]: {
      name: 'Droits sur la météo des services : Gérer les comportements périodiques',
      description: 'Les utilisateurs disposant de cette permission peuvent accéder à la liste des PBehaviors associés aux services (dans le sous-onglet des fenêtres modales des services)',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.executeInstruction]: {
      name: 'Droits sur la météo des services : Déclencher une remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des remédiations pour les alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.entityLinks]: {
      name: 'Droits sur la météo des services : Accès aux liens',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir les liens associés aux alarmes',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.moreInfos]: {
      name: 'Droits sur la météo des services : Accès à la fenêtre \'Plus d\'infos\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à la fenêtre modale "Plus d\'infos"',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.alarmsList]: {
      name: 'Droits sur la météo des services : Accès à la fenêtre \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ouvrir la liste des alarmes disponibles pour chaque service',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.pbehaviorList]: {
      name: 'Droits sur la météo des services : Accès à la liste des comportements périodiques du service',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à la liste de tous les PBehaviors des services (dans le sous-onglet des fenêtres modales des entités de service)',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.listFilters]: {
      name: 'Droits sur la météo des services : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.editFilter]: {
      name: 'Droits sur la météo des services : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres appliqués',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.addFilter]: {
      name: 'Droits sur la météo des services : Ajouter un filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.userFilter]: {
      name: 'Droits sur la météo des services : Utiliser les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.serviceWeather.actions.category]: {
      name: 'Droits sur la météo des services : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les services par catégorie',
    },

    /**
     * Business Counter Permissions
     */
    [USER_PERMISSIONS.business.counter.actions.alarmsList]: {
      name: 'Droits sur les Compteurs : Accès à la fenêtre \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux compteurs',
    },

    /**
     * Business Testing Weather Permissions
     */
    [USER_PERMISSIONS.business.testingWeather.actions.alarmsList]: {
      name: 'Droits sur les scénarios de test : Accès à la fenêtre \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux tests météorologiques',
    },

    /**
     * Mapping Permissions
     */
    [USER_PERMISSIONS.business.map.actions.alarmsList]: {
      name: 'Droits sur la cartographie : Accès à la fenêtre \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux points sur les cartes',
    },
    [USER_PERMISSIONS.business.map.actions.listFilters]: {
      name: 'Droits sur la cartographie : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.map.actions.editFilter]: {
      name: 'Droits sur la cartographie : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres des cartes',
    },
    [USER_PERMISSIONS.business.map.actions.addFilter]: {
      name: 'Droits sur la cartographie : Ajouter un filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres pour les cartes',
    },
    [USER_PERMISSIONS.business.map.actions.userFilter]: {
      name: 'Droits sur la cartographie : Afficher le filtre',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.map.actions.category]: {
      name: 'Droits sur la cartographie : Accès à l\'action \'Catégorie\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les points par catégories',
    },

    /**
     * Business Bar Chart Permissions
     */
    [USER_PERMISSIONS.business.barChart.actions.interval]: {
      name: 'Droits sur les histogrammes : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les intervalles de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.barChart.actions.sampling]: {
      name: 'Droits sur les histogrammes : Échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USER_PERMISSIONS.business.barChart.actions.listFilters]: {
      name: 'Droits sur les histogrammes : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.barChart.actions.editFilter]: {
      name: 'Droits sur les histogrammes : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.barChart.actions.addFilter]: {
      name: 'Droits sur les histogrammes : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.barChart.actions.userFilter]: {
      name: 'Droits sur les histogrammes : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business Line Chart Permissions
     */
    [USER_PERMISSIONS.business.lineChart.actions.interval]: {
      name: 'Droits sur les graphiques en ligne : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier',
    },
    [USER_PERMISSIONS.business.lineChart.actions.sampling]: {
      name: 'Droits sur les graphiques en ligne : Échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USER_PERMISSIONS.business.lineChart.actions.listFilters]: {
      name: 'Droits sur les graphiques en ligne : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.lineChart.actions.editFilter]: {
      name: 'Droits sur les graphiques en ligne : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.lineChart.actions.addFilter]: {
      name: 'Droits sur les graphiques en ligne : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.lineChart.actions.userFilter]: {
      name: 'Droits sur les graphiques en ligne : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business Pie Chart Permissions
     */
    [USER_PERMISSIONS.business.pieChart.actions.interval]: {
      name: 'Droits sur les diagrammes circulaires : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.pieChart.actions.sampling]: {
      name: 'Droits sur les diagrammes circulaires : Échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USER_PERMISSIONS.business.pieChart.actions.listFilters]: {
      name: 'Droits sur les diagrammes circulaires : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.pieChart.actions.editFilter]: {
      name: 'Droits sur les diagrammes circulaires : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.pieChart.actions.addFilter]: {
      name: 'Droits sur les diagrammes circulaires : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.pieChart.actions.userFilter]: {
      name: 'Droits sur les diagrammes circulaires : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business Numbers Permissions
     */
    [USER_PERMISSIONS.business.numbers.actions.interval]: {
      name: 'Droits sur les nombres : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.numbers.actions.sampling]: {
      name: 'Droits sur les nombres : Échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USER_PERMISSIONS.business.numbers.actions.listFilters]: {
      name: 'Droitr sur les nombres : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.numbers.actions.editFilter]: {
      name: 'Droits sur les nombres : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.numbers.actions.addFilter]: {
      name: 'Droits sur les nombres : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.numbers.actions.userFilter]: {
      name: 'Droits sur les nombres : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business User Statistics
     */
    [USER_PERMISSIONS.business.userStatistics.actions.interval]: {
      name: 'Droits sur les statistiques des utilisateurs : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.userStatistics.actions.listFilters]: {
      name: 'Droits sur les statistiques des utilisateurs : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.userStatistics.actions.editFilter]: {
      name: 'Droits sur les statistiques des utilisateurs : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.userStatistics.actions.addFilter]: {
      name: 'Droits sur les statistiques des utilisateurs : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.userStatistics.actions.userFilter]: {
      name: 'Droits sur les statistiques des utilisateurs : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business Alarm Statistics
     */
    [USER_PERMISSIONS.business.alarmStatistics.actions.interval]: {
      name: 'Droits sur les statistiques des alarmes : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.alarmStatistics.actions.listFilters]: {
      name: 'Droits sur les statistiques des alarmes : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.alarmStatistics.actions.editFilter]: {
      name: 'Droits sur les statistiques des alarmes : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.alarmStatistics.actions.addFilter]: {
      name: 'Droits sur les statistiques des alarmes : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.alarmStatistics.actions.userFilter]: {
      name: 'Droits sur les statistiques des alarmes : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Business Availability
     */
    [USER_PERMISSIONS.business.availability.actions.interval]: {
      name: 'Droits à la disponibilité : Intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USER_PERMISSIONS.business.availability.actions.listFilters]: {
      name: 'Droits à la disponibilité : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USER_PERMISSIONS.business.availability.actions.editFilter]: {
      name: 'Droits à la disponibilité : Modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USER_PERMISSIONS.business.availability.actions.addFilter]: {
      name: 'Droits à la disponibilité : Ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USER_PERMISSIONS.business.availability.actions.userFilter]: {
      name: 'Droits à la disponibilité : Afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.business.availability.actions.exportAsCsv]: {
      name: 'Droits de disponibilité : Exporter au format CSV',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exporter les disponibilités sous forme de fichier CSV',
    },

    /**
     * Technical General Permissions
     */
    [USER_PERMISSIONS.technical.view]: {
      name: 'Vues',
      description: 'Cette permission définit l\'accès à la liste des vues',
    },
    [USER_PERMISSIONS.technical.privateView]: {
      name: 'Vues privées',
      description: 'Cette autorisation définit l\'accès à la liste des vues privées',
    },
    [USER_PERMISSIONS.technical.role]: {
      name: 'Rôles',
      description: 'Cette autorisation définit l\'accès à la liste des rôles',
    },
    [USER_PERMISSIONS.technical.permission]: {
      name: 'Droits',
      description: 'Cette autorisation définit l\'accès à la liste des droits',
    },
    [USER_PERMISSIONS.technical.user]: {
      name: 'Utilisateurs',
      description: 'Cette autorisation définit l\'accès à la liste des utilisateurs',
    },
    [USER_PERMISSIONS.technical.parameters]: {
      name: 'Paramètres',
      description: 'Cette autorisation définit l\'accès aux réglages et paramètres de Canopsis',
    },
    [USER_PERMISSIONS.technical.broadcastMessage]: {
      name: 'Diffuser des messages',
      description: 'Cette autorisation définit l\'accès au panneau d\'administration des messages de diffusion',
    },
    [USER_PERMISSIONS.technical.playlist]: {
      name: 'Listes de lecture',
      description: 'Cette autorisation définit l\'accès aux paramètres des Playlists',
    },
    [USER_PERMISSIONS.technical.planningType]: {
      name: 'Planification : Types de comportements périodiques',
      description: 'Cette permission définit l\'accès aux types PBehavior',
    },
    [USER_PERMISSIONS.technical.planningReason]: {
      name: 'Planification : Raisons',
      description: 'Cette autorisation définit l\'accès aux raisons PBehavior',
    },
    [USER_PERMISSIONS.technical.planningExceptions]: {
      name: 'Planification : dates d\'exceptions',
      description: 'Cette autorisation définit l\'accès aux dates d\'exception pour PBehaviors',
    },
    [USER_PERMISSIONS.technical.remediationInstruction]: {
      name: 'Remédiation : consignes',
      description: 'Cette autorisation définit l\'accès à la liste des consignes',
    },
    [USER_PERMISSIONS.technical.remediationJob]: {
      name: 'Remédiation : Jobs',
      description: 'Cette autorisation définit l\'accès à la liste des Jobs',
    },
    [USER_PERMISSIONS.technical.remediationConfiguration]: {
      name: 'Remédiation : configuration',
      description: 'Cette autorisation définit l\'accès à la configuration de la remédiation',
    },
    [USER_PERMISSIONS.technical.remediationStatistic]: {
      name: 'Remédiation : statistiques',
      description: 'Cette autorisation définit l\'accès aux statistiques de remédiation',
    },
    [USER_PERMISSIONS.technical.healthcheck]: {
      name: 'Bilan de santé',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité Healthcheck',
    },
    [USER_PERMISSIONS.technical.techmetrics]: {
      name: 'Métriques techniques',
      description: 'Cette autorisation définit l\'accès aux métriques Tech',
    },
    [USER_PERMISSIONS.technical.engine]: {
      name: 'Moteurs',
      description: 'Cette autorisation définit l\'accès à la configuration des moteurs',
    },
    [USER_PERMISSIONS.technical.healthcheckStatus]: {
      name: 'Statut du bilan de santé',
      description: 'Statut du bilan de santé dans l\'entête  pour les utilisateurs disposant de cette autorisation',
    },
    [USER_PERMISSIONS.technical.kpi]: {
      name: 'KPI',
      description: 'Cette autorisation définit l\'accès aux métriques KPI',
    },
    [USER_PERMISSIONS.technical.kpiFilters]: {
      name: 'KPI : filtres',
      description: 'Cette autorisation définit l\'accès aux filtres pour les métriques KPI',
    },
    [USER_PERMISSIONS.technical.kpiRatingSettings]: {
      name: 'KPI : rating',
      description: 'Cette autorisation définit l\'accès aux paramètres d\'évaluation des KPI',
    },
    [USER_PERMISSIONS.technical.kpiCollectionSettings]: {
      name: 'Paramètres de la collection KPI',
      description: 'Cette autorisation définit l\'accès aux paramètres de la collection KPI',
    },
    [USER_PERMISSIONS.technical.map]: {
      name: 'Cartographie : Éditeur de carte',
      description: 'Cette permission définit l\'accès à l\'éditeur de carte',
    },
    [USER_PERMISSIONS.technical.shareToken]: {
      name: 'Jetons partagés',
      description: 'Cette autorisation définit l\'accès aux paramètres des jetons partagés',
    },
    [USER_PERMISSIONS.technical.widgetTemplate]: {
      name: 'Modèles de widgets',
      description: 'Cette autorisation définit l\'accès aux modèles de widgets',
    },
    [USER_PERMISSIONS.technical.maintenance]: {
      name: 'Mode Maintenance',
      description: 'Cette autorisation définit l\'accès au mode Maintenance',
    },
    [USER_PERMISSIONS.technical.tag]: {
      name: 'Gestion des tags',
      description: 'Cette autorisation définit l\'accès à la gestion des tags',
    },
    [USER_PERMISSIONS.technical.eventsRecord]: {
      name: 'Enregistrement des événements',
      description: 'Cette autorisation définit l\'accès aux enregistrements d\'événements',
    },

    /**
     * Technical Exploitation Permissions
     */
    [USER_PERMISSIONS.technical.exploitation.eventFilter]: {
      name: 'Exploitation : Filtres d\'événements',
      description: 'Cette autorisation définit l\'accès aux filtres d\'événements',
    },
    [USER_PERMISSIONS.technical.exploitation.pbehavior]: {
      name: 'Exploitation : Comportements périodiques',
      description: 'Cette autorisation définit l\'accès aux Comportements périodiques',
    },
    [USER_PERMISSIONS.technical.exploitation.snmpRule]: {
      name: 'Exploitation : règles SNMP',
      description: 'Cette permission définit l\'accès aux règles SNMP',
    },
    [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      name: 'Exploitation : Règles d\'informations dynamiques',
      description: 'Cette permission définit l\'accès à la fonctionnalité d\'infos dynamiques',
    },
    [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      name: 'Exploitation : Règles de méta-alarme',
      description: 'Cette autorisation définit l\'accès aux règles de méta-alarme et de corrélation',
    },
    [USER_PERMISSIONS.technical.exploitation.scenario]: {
      name: 'Exploitation : Scénarios',
      description: 'Cette permission définit l\'accès à la fonctionnalité des scénarios',
    },
    [USER_PERMISSIONS.technical.exploitation.idleRules]: {
      name: 'Exploitation : Règles d\'inactivité',
      description: 'Cette autorisation définit l\'accès aux règles d\'inactivité',
    },
    [USER_PERMISSIONS.technical.exploitation.flappingRules]: {
      name: 'Exploitation : Règles de bagot',
      description: 'Cette permission définit l\'accès aux règles de bagot',
    },
    [USER_PERMISSIONS.technical.exploitation.resolveRules]: {
      name: 'Exploitation : Règles de résolution',
      description: 'Cette autorisation définit l\'accès aux règles de résolution',
    },
    [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: {
      name: 'Exploitation : Règles de déclaration de tickets',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité de déclaration de ticket',
    },
    [USER_PERMISSIONS.technical.exploitation.linkRule]: {
      name: 'Exploitation : Règles de génération des liens',
      description: 'Cette autorisation définit l\'accès aux liens et les règles de liens',
    },

    /**
     * Technical Notification Permissions
     */
    [USER_PERMISSIONS.technical.notification.instructionStats]: {
      name: 'Notifications : Statistiques des remédiations',
      description: 'Cette permission définit l\'accès aux notifications associées aux statistiques dre remédiation',
    },

    /**
     * Technical Profile Permissions
     */
    [USER_PERMISSIONS.technical.profile.corporatePattern]: {
      name: 'Profil : Modèles partagés',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité des modèles partagés',
    },
    [USER_PERMISSIONS.technical.profile.theme]: {
      name: 'Thèmes graphiques',
      description: 'Cette autorisation définit l\'accès aux thèmes',
    },

    /**
     * API Permissions
     */
    [USER_PERMISSIONS.api.general.acl]: {
      name: 'Rôles, autorisations, utilisateurs',
      description: 'Accès à la route de l\'API vers les rôles, autorisations et utilisateurs CRUD',
    },
    [USER_PERMISSIONS.api.general.appInfoRead]: {
      name: 'Lire les informations sur l\'application (appInfo)',
      description: 'Accès à la route API pour lire les informations sur l\'application',
    },
    [USER_PERMISSIONS.api.general.alarmRead]: {
      name: 'Lire les alarmes',
      description: 'Accès à la route API pour lire les alarmes',
    },
    [USER_PERMISSIONS.api.general.alarmUpdate]: {
      name: 'Mettre à jour les alarmes',
      description: 'Accès à la route API pour mettre à jour les alarmes',
    },
    [USER_PERMISSIONS.api.general.entity]: {
      name: 'Entité',
      description: 'Accès à la route API vers les entités CRUD',
    },
    [USER_PERMISSIONS.api.general.entityservice]: {
      name: 'Entités type service',
      description: 'Accès à la route API vers les services CRUD',
    },
    [USER_PERMISSIONS.api.general.entitycategory]: {
      name: 'Catégories d\'entités',
      description: 'Accès à la route API vers les catégories d\'entités CRUD',
    },
    [USER_PERMISSIONS.api.general.event]: {
      name: 'Événement',
      description: 'Accès à la route API pour les événements',
    },
    [USER_PERMISSIONS.api.general.view]: {
      name: 'Vues',
      description: 'Accès à la route API vers les vues CRUD',
    },
    [USER_PERMISSIONS.api.general.viewgroup]: {
      name: 'Afficher les groupes',
      description: 'Accès à la route de l\'API vers les groupes de vues CRUD',
    },
    [USER_PERMISSIONS.api.general.privateViewGroups]: {
      name: 'Groupes de vue privée',
      description: 'Accès à la route API vers les groupes de vues privées CRUD',
    },
    [USER_PERMISSIONS.api.general.userInterfaceUpdate]: {
      name: 'Mettre à jour les paramètres de l\'interface utilisateur',
      description: 'Accès à la route API pour mettre à jour l\'interface utilisateur',
    },
    [USER_PERMISSIONS.api.general.userInterfaceDelete]: {
      name: 'Supprimer les paramètres de l\'interface utilisateur',
      description: 'Accès à la route API pour supprimer l\'interface utilisateur',
    },
    [USER_PERMISSIONS.api.general.datastorageRead]: {
      name: 'Paramètres de stockage de données : Lecture',
      description: 'Accès à la route API pour lire les paramètres de stockage de données',
    },
    [USER_PERMISSIONS.api.general.datastorageUpdate]: {
      name: 'Paramètres de stockage de données : Mise à jour',
      description: 'Accès à la route API pour modifier les paramètres de stockage des données',
    },
    [USER_PERMISSIONS.api.general.associativeTable]: {
      name: 'Tableau associatif',
      description: 'Accès à la route API avec un stockage de données associé (modèles d\'informations dynamiques, etc.)',
    },
    [USER_PERMISSIONS.api.general.stateSettings]: {
      name: 'Paramètres de criticités',
      description: 'Accès aux paramètres de criticités',
    },
    [USER_PERMISSIONS.api.general.files]: {
      name: 'Gestion des fichiers : Déposer',
      description: 'Accès à la route API vers les fichiers CRUD',
    },
    [USER_PERMISSIONS.api.general.healthcheck]: {
      name: 'Bilan de santé',
      description: 'Accès à la route de l\'API pour la vérification de l\'état',
    },
    [USER_PERMISSIONS.api.general.techmetrics]: {
      name: 'Métriques techniques',
      description: 'Accès à la route de l\'API vers les métriques techniques',
    },
    [USER_PERMISSIONS.api.general.contextgraph]: {
      name: 'Import du context-graph',
      description: 'Accès à la route API pour l\'import du graphe de contexte',
    },
    [USER_PERMISSIONS.api.general.broadcastMessage]: {
      name: 'Diffusion de messages',
      description: 'Accès à la route API pour les messages diffusés',
    },
    [USER_PERMISSIONS.api.general.junit]: {
      name: 'JUnit',
      description: 'Accès à la route de l\'API vers l\'API JUnit',
    },
    [USER_PERMISSIONS.api.general.notifications]: {
      name: 'Paramètres de notification',
      description: 'Accès à la route API pour les paramètres de notification',
    },
    [USER_PERMISSIONS.api.general.metrics]: {
      name: 'Métrique',
      description: 'Accès à la route API pour les métriques',
    },
    [USER_PERMISSIONS.api.general.metricsSettings]: {
      name: 'Paramètres des métriques',
      description: 'Accès à la route API pour les paramètres de métriques',
    },
    [USER_PERMISSIONS.api.general.ratingSettings]: {
      name: 'Paramètres de notation',
      description: 'Accès à la route API pour les paramètres de notation',
    },
    [USER_PERMISSIONS.api.general.filter]: {
      name: 'Filtres KPI',
      description: 'Accès à la route API vers les filtres KPI',
    },
    [USER_PERMISSIONS.api.general.corporatePattern]: {
      name: 'Modèles partagés',
      description: 'Accès à la route API pour les modèles partagés',
    },
    [USER_PERMISSIONS.api.general.exportConfigurations]: {
      name: 'Configurations d\'export',
      description: 'Accès à la route de l\'API pour exporter la configuration',
    },
    [USER_PERMISSIONS.api.general.map]: {
      name: 'Cartographie',
      description: 'Accès à la route API vers les cartes CRUD',
    },
    [USER_PERMISSIONS.api.general.shareToken]: {
      name: 'Jetons partagés',
      description: 'Accès à la route API vers les jetons partagés CRUD',
    },
    [USER_PERMISSIONS.api.general.declareTicketExecution]: {
      name: 'Déclaration de tickets',
      description: 'Accès à la route de l\'API pour exécuter les règles de déclaration de ticket',
    },
    [USER_PERMISSIONS.api.general.widgetTemplate]: {
      name: 'Modèles de widgets',
      description: 'Accès à la route API vers les modèles de widgets CRUD',
    },
    [USER_PERMISSIONS.api.general.maintenance]: {
      name: 'Mode Maintenance',
      description: 'Accès à l\'API route vers le mode maintenance',
    },
    [USER_PERMISSIONS.api.general.theme]: {
      name: 'Thèmes graphiques',
      description: 'Accès à l\'API route vers les thèmes',
    },

    [USER_PERMISSIONS.api.rules.action]: {
      name: 'Actions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les actions via l\'API',
    },
    [USER_PERMISSIONS.api.rules.dynamicinfos]: {
      name: 'Informations dynamiques',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les informations dynamiques par API',
    },
    [USER_PERMISSIONS.api.rules.eventFilter]: {
      name: 'Filtres d\'événement',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les filtres des événements par API',
    },
    [USER_PERMISSIONS.api.rules.idleRule]: {
      name: 'Règles d\'inactivité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les règles d\'inactivité via l\'API',
    },
    [USER_PERMISSIONS.api.rules.metaalarmrule]: {
      name: 'Règles de méta-alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les règles d\'alarme de méta via l\'API',
    },
    [USER_PERMISSIONS.api.rules.playlist]: {
      name: 'Listes de lecture',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les listes de lecture par API',
    },
    [USER_PERMISSIONS.api.rules.flappingRule]: {
      name: 'Règles de bagot',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les règles de bagot par API',
    },
    [USER_PERMISSIONS.api.rules.resolveRule]: {
      name: 'Règles de résolution',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les règles de résolution CRUD par API',
    },
    [USER_PERMISSIONS.api.rules.snmpRule]: {
      name: 'Règles SNMP',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer SNMP règles par API',
    },
    [USER_PERMISSIONS.api.rules.snmpMib]: {
      name: 'MIB SNMP',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer SNMP MIB par API',
    },
    [USER_PERMISSIONS.api.rules.declareTicketRule]: {
      name: 'Règles de déclaration de tickets',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les règles de déclaration de tickets par API',
    },
    [USER_PERMISSIONS.api.rules.linkRule]: {
      name: 'Règle de lien',
      description: 'Les utilisateurs disposant de cette autorisation peuvent créer des liens CRUD et des règles de lien par API',
    },

    [USER_PERMISSIONS.api.remediation.instruction]: {
      name: 'Instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les instructions par API',
    },
    [USER_PERMISSIONS.api.remediation.jobConfig]: {
      name: 'Configurations de remédiation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les configurations de remédiation par API',
    },
    [USER_PERMISSIONS.api.remediation.job]: {
      name: 'Jobs',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les jobs de remédiation par API',
    },
    [USER_PERMISSIONS.api.remediation.execution]: {
      name: 'Éxécuter les instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des instructions via l\'API',
    },
    [USER_PERMISSIONS.api.remediation.instructionApprove]: {
      name: 'Approuver les Instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent approuver les instructions via l\'API',
    },
    [USER_PERMISSIONS.api.remediation.messageRateStatsRead]: {
      name: 'Statistiques sur le débit des messages',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder aux statistiques de taux de messages par API',
    },

    [USER_PERMISSIONS.api.pbehavior.pbehavior]: {
      name: 'Comportements périodiques',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les comportements périodiques par API',
    },
    [USER_PERMISSIONS.api.pbehavior.pbehaviorException]: {
      name: 'Comportements périodiques : Dates d\'exceptions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les dates d\'exceptions par API',
    },
    [USER_PERMISSIONS.api.pbehavior.pbehaviorReason]: {
      name: 'Comportements périodiques : Raisons',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les raisons de comportements périodiques par API',
    },
    [USER_PERMISSIONS.api.pbehavior.pbehaviorType]: {
      name: 'Comportements périodiques : Types',
      description: 'Les utilisateurs disposant de cette autorisation peuvent gérer les types de comportements périodiques par API',
    },

    [USER_PERMISSIONS.api.eventsRecord.launchEventRecording]: {
      name: 'Lancement d\'un enregistrement d\'événements',
      description: 'Accès à la route API pour lancer et récupérer les enregistrements d\'événements',
    },
    [USER_PERMISSIONS.api.eventsRecord.resendEvents]: {
      name: 'Renvoyer les événements',
      description: 'Accès à la route API pour renvoyer les événements à partir des enregistrements d\'événements',
    },
  },
};
