import { USER_PERMISSIONS_PREFIXES, USERS_PERMISSIONS } from '@/constants';

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
  },
  api: {
    general: 'Général',
    rules: 'Règles',
    remediation: 'Remédiation',
    pbehavior: 'PBehavior',
  },
  permissions: {
    /**
     * Business Common Permissions
     */
    [USERS_PERMISSIONS.business.alarmsList.actions.variablesHelp]: {
      name: 'Accès à la liste des variables disponibles',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des variables dans la liste des alarmes et la météo du service',
    },

    /**
     * Business Alarms List Permissions
     */
    [USERS_PERMISSIONS.business.alarmsList.actions.ack]: {
      name: 'Droits sur la liste des alarmes : ack',
      description: 'Les utilisateurs disposant de cette autorisation peuvent acquitter les alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.fastAck]: {
      name: 'Droits sur la liste des alarmes : acquittement rapide',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer un acquittement rapide des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.ackRemove]: {
      name: 'Droits sur la liste des alarmes : annuler ack',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler la confirmation',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: {
      name: 'Droits sur la liste d\'alarmes : pcomportement action',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Comportement périodique" et modifier les comportements PB pour les alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.snooze]: {
      name: 'Droits sur la liste des alarmes : répéter l\'alarme',
      description: 'Users with this permission can snooze alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.declareTicket]: {
      name: 'Droits sur liste d\'alarmes : déclarer un ticket',
      description: 'Les utilisateurs avec cette permission peuvent faire la déclaration des tickets',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.associateTicket]: {
      name: 'Droits sur liste d\'alarmes : ticket associé',
      description: 'Les utilisateurs disposant de cette autorisation peuvent associer un ticket',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.cancel]: {
      name: 'Droits sur la liste des alarmes : annulation',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler les alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.fastCancel]: {
      name: 'Droits sur la liste des alarmes : annulation rapide de l\'alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer une annulation rapide des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.changeState]: {
      name: 'Droits sur la liste des alarmes : modifier l\'état',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les états des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.history]: {
      name: 'Rights on alarm list: history',
      description: 'Les utilisateurs disposant de cette autorisation peuvent afficher l\'historique des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup]: {
      name: 'Droits sur la liste des alarmes : actions manuelles sur les méta-alarmes',
      description: 'Les utilisateurs disposant de cette autorisation peuvent appliquer des règles de méta-alarmes manuelles et des alarmes de groupe',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.comment]: {
      name: 'Droits sur la liste des alarmes : Accès à l\'action \'Commentaire\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent commenter les alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: {
      name: 'Droits sur la liste des alarmes : afficher les filtres d\'alarmes',
      description: 'Les utilisateurs disposant de cette autorisation peuvent afficher la liste des filtres disponibles dans la liste des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.editFilter]: {
      name: 'Droits sur la liste des alarmes : modifier les filtres d\'alarme',
      description: 'Users with this permission can edit filters for alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.addFilter]: {
      name: 'Droits sur la liste des alarmes : ajouter des filtres d\'alarme',
      description: 'Users with this permission can add filters for alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.userFilter]: {
      name: 'Droits sur la liste des alarmes : afficher les filtres d\'alarme',
      description: 'Le filtre d\'alarme est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters]: {
      name: 'Rights on alarm list: Access to view filters by remediation instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir et appliquer la liste des filtres créés par des instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter]: {
      name: 'Droits sur la liste des alarmes : accès aux filtres d\'édition par les instructions de correction',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres à l\'aide d\'instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter]: {
      name: 'Droits sur la liste des alarmes : accès à l\'ajout de filtres par instructions de correction',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres à l\'aide d\'instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: {
      name: 'Droits sur la liste des alarmes : Accès aux filtres par instructions de remédiation',
      description: 'Le filtre par instructions est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.links]: {
      name: 'Droits sur la liste des alarmes : Accès aux liens',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder et suivre les liens dans la liste des alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.correlation]: {
      name: 'Droits sur la liste des alarmes : Accès au regroupement des alarmes corrélées',
      description: 'Les utilisateurs disposant de cette autorisation peuvent activer le regroupement des alarmes corrélées',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.executeInstruction]: {
      name: 'Droits sur liste d\'alarmes : Accès aux exécutions d\'instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des instructions pour corriger les alarmes',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.category]: {
      name: 'Droits sur la liste des alarmes : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer la liste des alarmes par catégorie',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: {
      name: 'Droits sur la liste des alarmes : Accès à l\'exportation des alarmes au format CSV',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exporter des alarmes vers CSV',
    },

    /**
     * Business Context Explorer Permissions
     */
    [USERS_PERMISSIONS.business.context.actions.createEntity]: {
      name: 'Droits sur l\'explorateur de contexte : créer une entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent créer de nouvelles entités',
    },
    [USERS_PERMISSIONS.business.context.actions.editEntity]: {
      name: 'Droits sur l\'explorateur de contexte : modifier l\'entité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les entités',
    },
    [USERS_PERMISSIONS.business.context.actions.duplicateEntity]: {
      name: 'Droits sur l\'explorateur de contexte : entité en double',
      description: 'Les utilisateurs disposant de cette autorisation peuvent dupliquer des entités',
    },
    [USERS_PERMISSIONS.business.context.actions.deleteEntity]: {
      name: 'Droits sur l\'explorateur de contexte : supprimer l\'entité',
      description: 'Users with this permission can delete entities',
    },
    [USERS_PERMISSIONS.business.context.actions.pbehaviorAdd]: {
      name: 'Droits sur l\'explorateur de contexte : action PBehavior',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Comportement périodique" et modifier les comportements PB pour les entités',
    },
    [USERS_PERMISSIONS.business.context.actions.massEnable]: {
      name: 'Droits sur le contexte : activer l\'action en masse',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Activer" pour les entités',
    },
    [USERS_PERMISSIONS.business.context.actions.massDisable]: {
      name: 'Droits sur le contexte : action de désactivation en masse',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à l\'action "Désactiver" pour les entités',
    },
    [USERS_PERMISSIONS.business.context.actions.listFilters]: {
      name: 'Droits sur l\'explorateur de contexte : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles dans l\'explorateur de contexte',
    },
    [USERS_PERMISSIONS.business.context.actions.editFilter]: {
      name: 'Droits sur l\'explorateur de contexte : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres d\'entité',
    },
    [USERS_PERMISSIONS.business.context.actions.addFilter]: {
      name: 'Droits sur l\'explorateur de contexte : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres sur les entités affichées dans l\'explorateur de contexte',
    },
    [USERS_PERMISSIONS.business.context.actions.userFilter]: {
      name: 'Droits sur l\'explorateur de contexte : afficher les filtres',
      description: 'Le filtre d\'entité est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.business.context.actions.category]: {
      name: 'Droits sur l\'explorateur de contexte : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les entités par catégorie',
    },
    [USERS_PERMISSIONS.business.context.actions.exportAsCsv]: {
      name: 'Droits sur l\'explorateur de contexte : Exporter au format csv',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exporter des entités sous forme de fichier CSV',
    },

    /**
     * Business Service Weather Permissions
     */
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityAck]: {
      name: 'Service météo : accès à l\'action \'Acquittement\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent acquitter les alarmes',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket]: {
      name: 'Service météo : accès à l\'action \'Associate Ticket\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent associer des tickets aux alarmes',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityComment]: {
      name: 'Service météo : accès à l\'action \'Commentaire\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des commentaires',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityValidate]: {
      name: 'Service météo : accès à l\'action \'Valider\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent valider les alarmes et modifier leur état en critique',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: {
      name: 'Service météo : accès à l\'action \'Invalider\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent invalider les alarmes et les annuler',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityPause]: {
      name: 'Service météo : accès à l\'action \'Pause\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent suspendre les alarmes (définir le type de PBehavior "Pause")',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityPlay]: {
      name: 'Service météo : accès à l\'action \'Jouer\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent activer les alarmes en pause (supprimer le type PBehavior "Pause")',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityCancel]: {
      name: 'Service météo : accès à l\'action \'Annuler\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent annuler les alarmes',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors]: {
      name: 'Service météo : Accès à la gestion des comportements',
      description: 'Les utilisateurs disposant de cette permission peuvent accéder à la liste des PBehaviors associés aux services (dans le sous-onglet des fenêtres modales des services)',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.executeInstruction]: {
      name: 'Météo de service : accès pour exécuter l\'instruction',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des instructions pour les alarmes',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks]: {
      name: 'Service météo : Accès aux liens',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir les liens associés aux alarmes',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: {
      name: 'Service météo : accès au modal \'Plus d\'infos\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder à la fenêtre modale "Plus d\'infos"',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: {
      name: 'Service météo : Accès à \'Liste des alarmes\' modal',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ouvrir la liste des alarmes disponibles pour chaque service',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList]: {
      name: '',
      description: '',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.listFilters]: {
      name: 'Droits sur le service météo : Afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.editFilter]: {
      name: 'Droits sur le service météo : Modifier le filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres appliqués',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.addFilter]: {
      name: 'Droits sur le service météo : Ajouter un filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.userFilter]: {
      name: 'Droits sur le service météo : Afficher le filtre',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.category]: {
      name: 'Droits sur le service météo : Filtrer par catégorie',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les services par catégorie',
    },

    /**
     * Business Counter Permissions
     */
    [USERS_PERMISSIONS.business.counter.actions.alarmsList]: {
      name: 'Compteur : Accès au modal \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux compteurs',
    },

    /**
     * Business Testing Weather Permissions
     */
    [USERS_PERMISSIONS.business.testingWeather.actions.alarmsList]: {
      name: 'Test météo : Accès au modal \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux tests météorologiques',
    },

    /**
     * Business Testing Weather Permissions
     */
    [USERS_PERMISSIONS.business.map.actions.alarmsList]: {
      name: 'Droits sur les cartes : Accès au modal \'Liste des alarmes\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des alarmes associées aux points sur les cartes',
    },
    [USERS_PERMISSIONS.business.map.actions.listFilters]: {
      name: 'Droits sur les cartes : Afficher le filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.map.actions.editFilter]: {
      name: 'Droits sur les cartes : modifier le filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres des cartes',
    },
    [USERS_PERMISSIONS.business.map.actions.addFilter]: {
      name: 'Droits sur les cartes : Ajouter un filtre',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres pour les cartes',
    },
    [USERS_PERMISSIONS.business.map.actions.userFilter]: {
      name: 'Droits sur les cartes : Afficher le filtre',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.business.map.actions.category]: {
      name: 'Droits sur les cartes : Accès à l\'action \'Catégorie\'',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les points par catégories',
    },

    /**
     * Business Bar Chart Permissions
     */
    [USERS_PERMISSIONS.business.barChart.actions.interval]: {
      name: 'Diagramme à barres : intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les intervalles de temps pour les données affichées',
    },
    [USERS_PERMISSIONS.business.barChart.actions.sampling]: {
      name: 'Diagramme à barres : échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USERS_PERMISSIONS.business.barChart.actions.listFilters]: {
      name: 'Graphique à barres : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.barChart.actions.editFilter]: {
      name: 'Graphique à barres : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USERS_PERMISSIONS.business.barChart.actions.addFilter]: {
      name: 'Diagramme à barres : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USERS_PERMISSIONS.business.barChart.actions.userFilter]: {
      name: 'Graphique à barres : afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Technical Line Chart Permissions
     */
    [USERS_PERMISSIONS.business.lineChart.actions.interval]: {
      name: 'Graphique en courbes : intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.sampling]: {
      name: 'Graphique en courbes : échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.listFilters]: {
      name: 'Graphique linéaire : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.editFilter]: {
      name: 'Graphique linéaire : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.addFilter]: {
      name: 'Graphique linéaire : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.userFilter]: {
      name: 'Graphique en courbes : afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Technical Pie Chart Permissions
     */
    [USERS_PERMISSIONS.business.pieChart.actions.interval]: {
      name: 'Piechart : intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.sampling]: {
      name: 'Diagramme à secteurs : échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.listFilters]: {
      name: 'Piechart : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.editFilter]: {
      name: 'Piechart : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.addFilter]: {
      name: 'Piechart : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.userFilter]: {
      name: 'Piechart : afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Technical Pie Chart Permissions
     */
    [USERS_PERMISSIONS.business.numbers.actions.interval]: {
      name: 'Nombres : intervalle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'intervalle de temps pour les données affichées',
    },
    [USERS_PERMISSIONS.business.numbers.actions.sampling]: {
      name: 'Chiffres : échantillonnage',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier l\'échantillonnage des données affichées',
    },
    [USERS_PERMISSIONS.business.numbers.actions.listFilters]: {
      name: 'Numéros : afficher les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent voir la liste des filtres disponibles',
    },
    [USERS_PERMISSIONS.business.numbers.actions.editFilter]: {
      name: 'Nombres : modifier les filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent modifier les filtres',
    },
    [USERS_PERMISSIONS.business.numbers.actions.addFilter]: {
      name: 'Chiffres : ajouter des filtres',
      description: 'Les utilisateurs disposant de cette autorisation peuvent ajouter des filtres',
    },
    [USERS_PERMISSIONS.business.numbers.actions.userFilter]: {
      name: 'Chiffres : afficher les filtres',
      description: 'Le filtre est affiché pour les utilisateurs disposant de cette autorisation',
    },

    /**
     * Technical General Permissions
     */
    [USERS_PERMISSIONS.technical.view]: {
      name: 'Vues',
      description: 'Cette permission définit l\'accès à la liste des vues',
    },
    [USERS_PERMISSIONS.technical.role]: {
      name: 'Les rôles',
      description: 'Cette autorisation définit l\'accès à la liste des rôles',
    },
    [USERS_PERMISSIONS.technical.permission]: {
      name: 'Droits',
      description: 'Cette autorisation définit l\'accès à la liste des droits',
    },
    [USERS_PERMISSIONS.technical.user]: {
      name: 'Utilisatrices',
      description: 'Cette autorisation définit l\'accès à la liste des utilisateurs',
    },
    [USERS_PERMISSIONS.technical.parameters]: {
      name: 'Paramètres',
      description: 'Cette autorisation définit l\'accès aux réglages et paramètres de Canopsis',
    },
    [USERS_PERMISSIONS.technical.broadcastMessage]: {
      name: 'Diffuser des messages',
      description: 'Cette autorisation définit l\'accès au panneau d\'administration des messages de diffusion',
    },
    [USERS_PERMISSIONS.technical.playlist]: {
      name: 'Listes de lecture',
      description: 'Cette autorisation définit l\'accès aux paramètres des Playlists',
    },
    [USERS_PERMISSIONS.technical.planning]: {
      name: 'Planification',
      description: 'Cette autorisation définit l\'accès aux paramètres Planning et PBehavior',
    },
    [USERS_PERMISSIONS.technical.planningType]: {
      name: 'Type de planification',
      description: 'Cette permission définit l\'accès aux types PBehavior',
    },
    [USERS_PERMISSIONS.technical.planningReason]: {
      name: 'Raison de la planification',
      description: 'Cette autorisation définit l\'accès aux raisons PBehavior',
    },
    [USERS_PERMISSIONS.technical.planningExceptions]: {
      name: 'Planification des dates d\'exceptions',
      description: 'Cette autorisation définit l\'accès aux dates d\'exception pour PBehaviors',
    },
    [USERS_PERMISSIONS.technical.remediation]: {
      name: 'Remédiation',
      description: 'Cette autorisation définit l\'accès au panneau d\'administration de la correction',
    },
    [USERS_PERMISSIONS.technical.remediationInstruction]: {
      name: 'Instruction de correction',
      description: 'Cette autorisation définit l\'accès à la liste des Instructions',
    },
    [USERS_PERMISSIONS.technical.remediationJob]: {
      name: 'Travail de remédiation',
      description: 'Cette autorisation définit l\'accès à la liste des Jobs',
    },
    [USERS_PERMISSIONS.technical.remediationConfiguration]: {
      name: 'Configuration de la correction',
      description: 'Cette autorisation définit l\'accès à la configuration de la correction',
    },
    [USERS_PERMISSIONS.technical.remediationStatistic]: {
      name: 'Statistiques de remédiation',
      description: 'Cette autorisation définit l\'accès aux statistiques de remédiation',
    },
    [USERS_PERMISSIONS.technical.healthcheck]: {
      name: 'Bilan de santé',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité Healthcheck',
    },
    [USERS_PERMISSIONS.technical.techmetrics]: {
      name: 'Métriques techniques',
      description: 'Cette autorisation définit l\'accès aux métriques Tech',
    },
    [USERS_PERMISSIONS.technical.engine]: {
      name: 'Moteurs',
      description: '',
    },
    [USERS_PERMISSIONS.technical.healthcheckStatus]: {
      name: 'État de la vérification de l\'état',
      description: 'L\'état de la vérification de l\'état du système est affiché dans l\'en-tête pour les utilisateurs disposant de cette autorisation',
    },
    [USERS_PERMISSIONS.technical.kpi]: {
      name: 'KPI',
      description: 'Cette autorisation définit l\'accès aux métriques KPI',
    },
    [USERS_PERMISSIONS.technical.kpiFilters]: {
      name: 'Filtres KPI',
      description: 'Cette autorisation définit l\'accès aux filtres pour les métriques KPI',
    },
    [USERS_PERMISSIONS.technical.kpiRatingSettings]: {
      name: 'Paramètres d\'évaluation des KPI',
      description: 'Cette autorisation définit l\'accès aux paramètres d\'évaluation des KPI',
    },
    [USERS_PERMISSIONS.technical.kpiCollectionSettings]: {
      name: 'Paramètres de collecte de KPI',
      description: 'Cette autorisation définit l\'accès aux paramètres de KPI Collection',
    },
    [USERS_PERMISSIONS.technical.map]: {
      name: 'Éditeur de carte',
      description: 'Cette permission définit l\'accès à l\'éditeur de carte',
    },
    [USERS_PERMISSIONS.technical.shareToken]: {
      name: 'Partager le jeton',
      description: 'Cette autorisation définit l\'accès aux paramètres des jetons partagés',
    },
    [USERS_PERMISSIONS.technical.widgetTemplate]: {
      name: 'Modèles de widgets',
      description: 'Cette autorisation définit l\'accès aux modèles de widget',
    },

    /**
     * Technical Exploitation Permissions
     */
    [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
      name: 'Exploitation : Filtres d\'événements',
      description: 'Cette autorisation définit l\'accès aux filtres d\'événements',
    },
    [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
      name: 'Exploitation : Pcomportements',
      description: 'Cette autorisation définit l\'accès aux événements PBehavior',
    },
    [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
      name: 'Exploitation : règles SNMP',
      description: 'Cette permission définit l\'accès aux règles SNMP',
    },
    [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      name: 'Exploitation : règles d\'information dynamiques',
      description: 'Cette permission définit l\'accès à la fonctionnalité d\'infos dynamiques',
    },
    [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      name: 'Exploitation : règles de méta-alarme',
      description: 'Cette autorisation définit l\'accès aux règles de méta-alarme et de corrélation',
    },
    [USERS_PERMISSIONS.technical.exploitation.scenario]: {
      name: 'Exploitation : Scénarios',
      description: 'Cette permission définit l\'accès à la fonctionnalité des scénarios',
    },
    [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
      name: 'Exploitation : règles d\'inactivité',
      description: 'Cette autorisation définit l\'accès aux règles d\'inactivité',
    },
    [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
      name: 'Exploitation : règles de battement',
      description: 'Cette permission définit l\'accès aux règles de battement',
    },
    [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
      name: 'Exploitation : résoudre les règles',
      description: 'Cette autorisation définit l\'accès aux règles de résolution',
    },
    [USERS_PERMISSIONS.technical.exploitation.declareTicketRule]: {
      name: 'Exploitation : Déclarer les règles du ticket',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité de déclaration de ticket',
    },
    [USERS_PERMISSIONS.technical.exploitation.linkRule]: {
      name: 'Exploitation : règles de liaison',
      description: 'Cette autorisation définit l\'accès aux liens et les règles de lien',
    },

    /**
     * Technical Notification Permissions
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      name: 'Notifications : statistiques des instructions',
      description: 'Cette permission définit l\'accès aux notifications associées aux statistiques d\'instructions',
    },

    /**
     * Technical Profile Permissions
     */
    [USERS_PERMISSIONS.technical.profile.corporatePattern]: {
      name: 'Profil : Modèles d\'entreprise',
      description: 'Cette autorisation définit l\'accès à la fonctionnalité des modèles d\'entreprise',
    },

    /**
     * API Permissions
     */
    [USERS_PERMISSIONS.api.general.acl]: {
      name: 'Rôles, autorisations, utilisateurs',
      description: 'Accès à la route de l\'API vers les rôles, autorisations et utilisateurs CRUD',
    },
    [USERS_PERMISSIONS.api.general.appInfoRead]: {
      name: 'Lire les informations sur l\'application',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.alarmRead]: {
      name: 'Lire les alarmes',
      description: 'Accès à la route API pour lire les alarmes',
    },
    [USERS_PERMISSIONS.api.general.alarmUpdate]: {
      name: 'Mettre à jour les alarmes',
      description: 'Accès à la route API pour mettre à jour les alarmes',
    },
    [USERS_PERMISSIONS.api.general.entity]: {
      name: 'Entité',
      description: 'Accès à la route API vers les entités CRUD',
    },
    [USERS_PERMISSIONS.api.general.entityservice]: {
      name: 'Service d\'entité',
      description: 'Accès à la route API vers les services CRUD',
    },
    [USERS_PERMISSIONS.api.general.entitycategory]: {
      name: 'Catégories d\'entités',
      description: 'Accès à la route API vers les catégories d\'entités CRUD',
    },
    [USERS_PERMISSIONS.api.general.event]: {
      name: 'Événement',
      description: 'Accès à la route API pour les événements',
    },
    [USERS_PERMISSIONS.api.general.view]: {
      name: 'Vues',
      description: 'Accès à la route API vers les vues CRUD',
    },
    [USERS_PERMISSIONS.api.general.viewgroup]: {
      name: 'Afficher les groupes',
      description: 'Accès à la route de l\'API vers les groupes de vues CRUD',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceUpdate]: {
      name: 'Mettre à jour l\'interface utilisateur',
      description: 'Accès à la route API pour mettre à jour l\'interface utilisateur',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceDelete]: {
      name: 'Supprimer l\'interface utilisateur',
      description: 'Accès à la route API pour supprimer l\'interface utilisateur',
    },
    [USERS_PERMISSIONS.api.general.datastorageRead]: {
      name: 'Paramètres de stockage de données lus',
      description: 'Accès à la route API pour lire les paramètres de stockage de données',
    },
    [USERS_PERMISSIONS.api.general.datastorageUpdate]: {
      name: 'Mise à jour des paramètres de stockage des données',
      description: 'Accès à la route API pour modifier les paramètres de stockage des données',
    },
    [USERS_PERMISSIONS.api.general.associativeTable]: {
      name: 'Tableau associatif',
      description: 'Accès à la route API avec un stockage de données associé (modèles d\'informations dynamiques, etc.)',
    },
    [USERS_PERMISSIONS.api.general.stateSettings]: {
      name: 'Paramètres d\'état',
      description: 'Accès aux paramètres d\'acheminement de l\'API vers l\'état',
    },
    [USERS_PERMISSIONS.api.general.files]: {
      name: 'Déposer',
      description: 'Accès à la route API vers les fichiers CRUD',
    },
    [USERS_PERMISSIONS.api.general.healthcheck]: {
      name: 'Bilan de santé',
      description: 'Accès à la route de l\'API pour la vérification de l\'état',
    },
    [USERS_PERMISSIONS.api.general.techmetrics]: {
      name: 'Métriques techniques',
      description: 'Accès à la route de l\'API vers les métriques techniques',
    },
    [USERS_PERMISSIONS.api.general.contextgraph]: {
      name: 'Importation de graphiques de contexte',
      description: 'Accès à la route API pour l\'import du graphe de contexte',
    },
    [USERS_PERMISSIONS.api.general.broadcastMessage]: {
      name: 'Message diffusé',
      description: 'Accès à la route API pour les messages diffusés',
    },
    [USERS_PERMISSIONS.api.general.junit]: {
      name: 'JUnit',
      description: 'Accès à la route de l\'API vers l\'API JUnit',
    },
    [USERS_PERMISSIONS.api.general.notifications]: {
      name: 'Paramètres de notification',
      description: 'Accès à la route API pour les paramètres de notification',
    },
    [USERS_PERMISSIONS.api.general.metrics]: {
      name: 'Métrique',
      description: 'Accès à la route API pour les métriques',
    },
    [USERS_PERMISSIONS.api.general.metricsSettings]: {
      name: 'Paramètres des métriques',
      description: 'Accès à la route API pour les paramètres de métrique',
    },
    [USERS_PERMISSIONS.api.general.ratingSettings]: {
      name: 'Paramètres de notation',
      description: 'Accès à la route API pour les paramètres de notation',
    },
    [USERS_PERMISSIONS.api.general.filter]: {
      name: 'Filtres KPI',
      description: 'Accès à la route API vers les filtres KPI',
    },
    [USERS_PERMISSIONS.api.general.corporatePattern]: {
      name: 'Modèles d\'entreprise',
      description: 'Accès à la route API pour les modèles d\'entreprise',
    },
    [USERS_PERMISSIONS.api.general.exportConfigurations]: {
      name: 'Configurations d\'exportation',
      description: 'Accès à la route de l\'API pour exporter la configuration',
    },
    [USERS_PERMISSIONS.api.general.map]: {
      name: 'Carte',
      description: 'Accès à la route API vers les cartes CRUD',
    },
    [USERS_PERMISSIONS.api.general.shareToken]: {
      name: 'Partager le jeton',
      description: 'Accès à la route API vers les jetons partagés CRUD',
    },
    [USERS_PERMISSIONS.api.general.declareTicketExecution]: {
      name: 'Exécuter déclarer les règles de ticket',
      description: 'Accès à la route de l\'API pour exécuter les règles de déclaration de ticket',
    },
    [USERS_PERMISSIONS.api.general.widgetTemplate]: {
      name: 'Modèles de widgets',
      description: 'Accès à la route API vers les modèles de widgets CRUD',
    },

    [USERS_PERMISSIONS.api.rules.action]: {
      name: 'Actions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer des actions CRUD via l\'API',
    },
    [USERS_PERMISSIONS.api.rules.dynamicinfos]: {
      name: 'Informations dynamiques',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD informations dynamiques par API',
    },
    [USERS_PERMISSIONS.api.rules.eventFilter]: {
      name: 'Filtre d\'événement',
      description: 'Les utilisateurs disposant de cette autorisation peuvent filtrer les événements CRUD par API',
    },
    [USERS_PERMISSIONS.api.rules.idleRule]: {
      name: 'Règle d\'inactivité',
      description: 'Les utilisateurs disposant de cette autorisation peuvent utiliser les règles d\'inactivité CRUD via l\'API',
    },
    [USERS_PERMISSIONS.api.rules.metaalarmrule]: {
      name: 'Règle de méta-alarme',
      description: 'Les utilisateurs disposant de cette autorisation peuvent appliquer les règles d\'alarme de méta CRUD via l\'API',
    },
    [USERS_PERMISSIONS.api.rules.playlist]: {
      name: 'Listes de lecture',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder aux listes de lecture CRUD par API',
    },
    [USERS_PERMISSIONS.api.rules.flappingRule]: {
      name: 'Règle de battement',
      description: 'Les utilisateurs disposant de cette autorisation peuvent appliquer les règles de battement CRUD par API',
    },
    [USERS_PERMISSIONS.api.rules.resolveRule]: {
      name: 'Résoudre la règle',
      description: 'Les utilisateurs disposant de cette autorisation peuvent résoudre les règles CRUD par API',
    },
    [USERS_PERMISSIONS.api.rules.snmpRule]: {
      name: 'Règle SNMP',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD SNMP règles par API',
    },
    [USERS_PERMISSIONS.api.rules.snmpMib]: {
      name: 'MIB SNMP',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD SNMP MIB par API',
    },
    [USERS_PERMISSIONS.api.rules.declareTicketRule]: {
      name: 'Déclarer la règle de ticket',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD déclarer des règles de ticket par API',
    },
    [USERS_PERMISSIONS.api.rules.linkRule]: {
      name: 'Règle de lien',
      description: 'Les utilisateurs disposant de cette autorisation peuvent créer des liens CRUD et des règles de lien par API',
    },

    [USERS_PERMISSIONS.api.remediation.instruction]: {
      name: 'Instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent recevoir des instructions CRUD via l\'API',
    },
    [USERS_PERMISSIONS.api.remediation.jobConfig]: {
      name: 'Configurations de tâches',
      description: 'Les utilisateurs disposant de cette autorisation peuvent configurer les tâches CRUD par API',
    },
    [USERS_PERMISSIONS.api.remediation.job]: {
      name: 'Emplois',
      description: 'Les utilisateurs disposant de cette autorisation peuvent effectuer des tâches CRUD via l\'API',
    },
    [USERS_PERMISSIONS.api.remediation.execution]: {
      name: 'Exécute les instructions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent exécuter des instructions via l\'API',
    },
    [USERS_PERMISSIONS.api.remediation.instructionApprove]: {
      name: 'Instruction approuver',
      description: 'Les utilisateurs disposant de cette autorisation peuvent approuver les instructions via l\'API',
    },
    [USERS_PERMISSIONS.api.remediation.messageRateStatsRead]: {
      name: 'Statistiques sur le débit des messages',
      description: 'Les utilisateurs disposant de cette autorisation peuvent accéder aux statistiques de taux de messages par API',
    },

    [USERS_PERMISSIONS.api.pbehavior.pbehavior]: {
      name: 'PBehaviors',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD PBehavior dates des événements par API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorException]: {
      name: 'PBehavior exceptions',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD PBehavior dates d\'exceptions par API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorReason]: {
      name: 'Raisons comportementales',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD PBehavior raisons dates par API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorType]: {
      name: 'Types de PBehavior',
      description: 'Les utilisateurs disposant de cette autorisation peuvent CRUD PBehavior types dates par API',
    },
  },
};
