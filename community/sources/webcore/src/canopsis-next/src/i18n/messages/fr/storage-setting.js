export default {
  receivedFor: 'Reçu pour',
  olderThan: 'Plus ancien que',
  alarm: {
    archiveAfter: 'Archiver les données d\'alarmes résolues',
    deleteAfter: 'Supprimer les données d\'alarmes résolues',
  },
  junit: {
    title: 'JUnit',
    deleteAfter: 'Supprimer les données des suites de tests',
    deleteAfterHelpText: '(XMLs, captures d\'écran, vidéos)',
  },
  remediation: {
    title: 'Instructions',
    deleteAfter: 'Supprimer la chronologie des instructions',
    deleteStatsAfter: 'Supprimer les statistiques d\'instructions',
    deleteModStatsAfter: 'Supprimer le résumé des instructions',
  },
  entity: {
    title: 'Stockage des données d\'entités',
    titleHelp: 'Toutes les entités désactivées avec des alarmes associées peuvent être archivées (déplacées vers une collection séparée) et/ou supprimées définitivement.',
    archiveDependencies: 'Supprimer également les entités impactantes et dépendantes',
    archiveDependenciesHelp: '<strong>Pour les connecteurs :</strong>\n'
      + '<ul><li>tous les <strong class="font-italic">composants</strong> et <strong class="font-italic">ressources</strong> impactants et dépendants seront archivés ou supprimés définitivement.</li></ul>\n'
      + '<strong>Pour les composants :</strong>\n'
      + '<ul><li>toutes les <strong class="font-italic">ressources</strong> dépendantes seront également archivées ou supprimées définitivement.</li></ul>',
    archiveDisabled: 'Archiver les entités désactivées',
  },
  entityUnlinked: {
    title: 'Entités non liées',
    archiveAfter: 'Archiver les entités sans événements',
    archiveUnlinked: 'Archiver les entités non liées',
    archiveUnlinkedAfter: 'Archiver les entités non liées sans événements reçus depuis',
  },
  entityArchived: {
    title: 'Stockage des données archivées',
    titleHelp: 'Toutes les entités archivées peuvent être supprimées définitivement.',
    cleanArchive: 'Nettoyer l\'archive',
  },
  pbehavior: {
    title: 'PBehavior',
    deleteAfter: 'Supprimer les PBehavior inactifs',
    deleteAfterHelpText: 'Les PBehaviors inactifs seront supprimés après la période définie depuis le dernier événement',
  },
  healthCheck: {
    title: 'Vérification de santé',
    deleteAfter: 'Supprimer le flux entrant FIFO',
  },
  webhook: {
    title: 'Webhooks',
    deleteAfter: 'Supprimer l\'historique des requêtes webhooks',
    logCredentials: 'Ouvrir les données d\'authentification dans les logs',
    logCredentialsHelpText: 'Affecte la façon dont les mots de passe, tokens et données d\'authentification sont écrits dans les logs. \n'
      + '<ul><li>activé : de manière ouverte (non recommandé)</li>'
      + '<li>désactivé : masqué avec ***</li></ul>',
  },
  metrics: {
    title: 'Métriques internes',
    deleteAfter: 'Supprimer les métriques',
  },
  perfDataMetrics: {
    title: 'Métriques externes',
    deleteAfter: 'Supprimer les métriques',
  },
  eventFilterFailure: {
    title: 'Messages d\'erreur du filtre d\'événements',
    deleteAfter: 'Supprimer les messages d\'erreur',
    deleteAfterHelpText: 'Toutes les erreurs seront toujours disponibles dans les logs',
  },
  alarmExternalTag: {
    title: 'Tags externes d\'alarme',
    deleteAfter: 'Supprimer les tags externes',
  },
  eventsRecords: {
    title: 'Enregistrements d\'événements',
    deleteAfter: 'Supprimer les enregistrements d\'événements',
  },
  history: {
    scriptLaunched: 'Script lancé à {launchedAt}.',
    alarm: {
      deletedCount: 'Alarmes supprimées : {count}.',
      archivedCount: 'Alarmes archivées : {count}.',
    },
    entity: {
      deletedCount: 'Entités supprimées : {count}.',
      archivedCount: 'Entités archivées : {count}.',
    },
    alarmExternalTag: {
      deletedCount: 'Tags externes d\'alarme supprimés : {count}.',
    },
  },
};
