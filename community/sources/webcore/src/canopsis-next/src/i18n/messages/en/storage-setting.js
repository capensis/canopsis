export default {
  receivedFor: 'Received for',
  olderThan: 'Older than',
  alarm: {
    archiveAfter: 'Archive resolved alarms data',
    deleteAfter: 'Delete resolved alarms data',
  },
  junit: {
    title: 'JUnit',
    deleteAfter: 'Delete test suites data',
    deleteAfterHelpText: '(XMLs, screenshots, videos)',
  },
  remediation: {
    title: 'Instructions',
    deleteAfter: 'Delete instructions timeline',
    deleteStatsAfter: 'Delete instruction statistics',
    deleteModStatsAfter: 'Delete instructions summary',
  },
  entity: {
    title: 'Entities data storage',
    titleHelp: 'All disabled entities with associated alarms can be archived (moved to the separate collection) and/or deleted forever.',
    archiveDependencies: 'Remove the impacting and dependent entities as well',
    archiveDependenciesHelp: '<strong>For connectors:</strong>\n'
      + '<ul><li>all impacting and dependent <strong class="font-italic">components</strong> and <strong class="font-italic">resources</strong> will be archived or deleted forever.</li></ul>\n'
      + '<strong>For components:</strong>\n'
      + '<ul><li>all dependent <strong class="font-italic">resources</strong> will be archived or deleted forever as well.</li></ul>',
    archiveDisabled: 'Archive disabled entities',
  },
  entityUnlinked: {
    title: 'Unlinked entities',
    archiveAfter: 'Archive entities with no events',
    archiveUnlinked: 'Archive unlinked entities',
    archiveUnlinkedAfter: 'Archive unlinked entities with no events received for',
  },
  entityArchived: {
    title: 'Archived data storage',
    titleHelp: 'All the archived entities can be deleted forever.',
    cleanArchive: 'Clean archive',
  },
  pbehavior: {
    title: 'PBehavior',
    deleteAfter: 'Delete inactive PBehavior',
    deleteAfterHelpText: 'Inactive PBehaviors will be deleted after the defined time period from the last event',
  },
  healthCheck: {
    title: 'Healthcheck',
    deleteAfter: 'Delete FIFO incoming flow',
  },
  webhook: {
    title: 'Webhooks',
    deleteAfter: 'Delete webhooks requests history',
    logCredentials: 'Open auth data in logs',
    logCredentialsHelpText: 'Affects how passwords, tokens and auth data are written in logs. \n'
      + '<ul><li>enabled: in open way (not recommended)</li>'
      + '<li>disabled: hidden with ***</li></ul>',
  },
  metrics: {
    title: 'Internal metrics',
    deleteAfter: 'Delete metrics',
  },
  perfDataMetrics: {
    title: 'External metrics',
    deleteAfter: 'Delete metrics',
  },
  eventFilterFailure: {
    title: 'Event filter error messages',
    deleteAfter: 'Delete error messages',
    deleteAfterHelpText: 'All errors will still be available in logs',
  },
  alarmExternalTag: {
    title: 'Alarm external tags',
    deleteAfter: 'Delete external tags',
  },
  eventsRecords: {
    title: 'Events recordings',
    deleteAfter: 'Delete events recordings',
  },
  entityInfosLog: {
    title: 'Event filter: entity enrichment logs',
    deleteAfter: 'Delete event filter logs older than',
  },
  history: {
    scriptLaunched: 'Script launched at {launchedAt}.',
    alarm: {
      deletedCount: 'Alarms deleted: {count}.',
      archivedCount: 'Alarms archived: {count}.',
    },
    entity: {
      deletedCount: 'Entities deleted: {count}.',
      archivedCount: 'Entities archived: {count}.',
    },
    alarmExternalTag: {
      deletedCount: 'Alarm external tags deleted: {count}.',
    },
  },
};
