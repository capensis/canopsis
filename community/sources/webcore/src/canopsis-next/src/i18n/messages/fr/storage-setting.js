export default {
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
  webhook: {
    title: 'Stockage de données Webhooks',
    titleHelp: 'L\'historique de toutes les demandes de webhook est conservé dans des journaux',
    deleteAfter: 'Effacer les journaux des webhooks après',
    deleteAfterHelpText: 'Tous les historiques de demandes de webhook antérieurs à la période définie seront supprimés',
    sanitizeRequest: 'Ouvrir les données d\'authentification dans les journaux',
    sanitizeRequestHelpText: 'Lorsqu\'il est activé, toutes les informations d\'identification et les données d\'authentification sont écrites dans les journaux de manière ouverte (non recommandé). \n'
      + 'Lorsqu\'il est désactivé, tous les mots de passe, jetons et données d\'authentification sont masqués et écrits sous la forme *** dans les journaux.',
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
};
