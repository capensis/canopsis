export default {
  alarm: {
    title: 'Alarm data storage',
    titleHelp: 'When switched on, the resolved alarms data will be archived and/or deleted after the defined time period.',
    archiveAfter: 'Archive resolved alarms data after',
    deleteAfter: 'Delete resolved alarms data after',
  },
  junit: {
    title: 'JUnit data storage',
    deleteAfter: 'Delete test suites data after',
    deleteAfterHelpText: 'When switched on, the JUnit test suites data (XMLs, screenshots and videos) will be deleted after the defined time period.',
  },
  remediation: {
    title: 'Instructions data storage',
    deleteAfter: 'Delete instructions timeline data after',
    deleteAfterHelpText: 'When switched on, the instructions timelines data will be deleted after the defined time period.',
    deleteStatsAfter: 'Delete instruction statistics data after',
    deleteStatsAfterHelpText: 'When switched on, the instruction statistics will be deleted after the defined time period.',
    deleteModStatsAfter: 'Delete instructions summary data after',
    deleteModStatsAfterHelpText: 'When switched on, the instructions summary data will be deleted after the defined time period.',
  },
  entity: {
    title: 'Entities data storage',
    titleHelp: 'All disabled entities with associated alarms can be archived (moved to the separate collection) and/or deleted forever.',
    archiveDependencies: 'Remove the impacting and dependent entities as well',
    archiveDependenciesHelp: 'For connectors, all impacting and dependent components and resources will be archived or deleted forever. For components, all dependent resources will be archived or deleted forever as well.',
    archiveDisabled: 'Archive disabled entities',
  },
  entityUnlinked: {
    title: 'Unlinked entities storage',
    titleHelp: 'All unlinked connectors, components and resources without alarms and updated long time ago can be archived.',
    archiveBefore: 'Archive entities when no events received for',
    archiveUnlinked: 'Archive unlinked entities',
  },
  entityArchived: {
    title: 'Archived data storage',
    titleHelp: 'All the archived entities can be deleted forever.',
    cleanArchive: 'Clean archive',
  },
  pbehavior: {
    title: 'PBehavior data storage',
    deleteAfter: 'Delete PBehavior data after',
    deleteAfterHelpText: 'When switched on, inactive PBehaviors will be deleted after the defined time period from the last event.',
  },
  healthCheck: {
    title: 'Healthcheck data storage',
    deleteAfter: 'Delete FIFO incoming flow data after',
  },
  webhook: {
    title: 'Webhooks data storage',
    titleHelp: 'All webhook requests history is kept in logs',
    deleteAfter: 'Clear webhooks logs after',
    deleteAfterHelpText: 'All webhook requests history older than the defined time period will be deleted',
    logCredentials: 'Open auth data in logs',
    logCredentialsHelpText: 'When enabled, all credentials and auth data is written in logs in open way (not recommended). \n'
      + 'When disabled, all passwords, tokens and auth data is hidden and written as *** in logs.',
  },
  metrics: {
    title: 'Internal metrics data storage',
    titleHelp: 'When enabled, internal metrics data will be deleted after the defined time period',
    deleteAfter: 'Clear metrics storage after',
    deleteAfterHelpText: 'All internal metrics older than the defined time period will be deleted',
  },
  perfDataMetrics: {
    title: 'External metrics data storage',
    titleHelp: 'When enabled, external metrics data will be deleted after the defined time period',
    deleteAfter: 'Clear metrics storage after',
    deleteAfterHelpText: 'All external metrics older than the defined time period will be deleted',
  },
  eventFilterFailure: {
    title: 'Error messages data storage',
    titleHelp: 'When enabled, error messages data will be deleted after the defined time period. However, all errors are available in logs.',
    deleteAfter: 'Clear error messages older than',
    deleteAfterHelpText: 'All error messages older than the defined time period will be deleted',
  },
  alarmExternalTag: {
    title: 'Alarm external tags data storage',
    deleteAfter: 'Clear external tags after',
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
  eventsRecords: {
    title: 'Events recordings data storage',
    titleHelp: 'All recordings older than the defined time period will be deleted',
    deleteAfter: 'Delete events recordings after',
  },
};
