import { USER_PERMISSIONS, GROUPED_USER_PERMISSIONS_KEYS } from '@/constants';

export default {
  hideMessage: 'Got it! Hide',
  learnMore: 'Learn more on {link}',

  /**
   * Exploitation
   */
  [USER_PERMISSIONS.technical.exploitation.eventFilter]: {
    title: 'Event filter',
    message: 'The event-filter is a feature of engine-che, allowing to define rules handling events.',
  },

  [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: {
    title: 'Dynamic informations',
    message: 'The Canopsis Dynamic infos are used to add information to the alarms. This information is defined with rules indicating under which conditions information must be presented on an alarm.',
  },

  [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
    title: 'Meta alarm rule',
    message: 'Meta alarm rules can be used for grouping alarms by types and criteria (parent-child relationship, time interval, etc).',
  },

  [USER_PERMISSIONS.technical.exploitation.idleRules]: {
    title: 'Idle rules',
    message: 'Idle rules for entities and alarms can be used in order to monitor events and alarm states in order to be aware when events are not receiving or alarm state is not changed for a long time because of errors or invalid configuration.',
  },

  [USER_PERMISSIONS.technical.exploitation.flappingRules]: {
    title: 'Flapping rules',
    message: 'An alarm inherits flapping status when it oscillates from an alert to a stable state a certain number of times over a given period.',
  },

  [USER_PERMISSIONS.technical.exploitation.resolveRules]: {
    title: 'Resolve rules',
    message: 'When an alarm receives a recovery type event, it changes to the closed status. Before considering this alarm as definitively resolved, Canopsis can wait for a delay. This delay can be useful if the alarm flaps or if the user wishes keep the alarm open in case of error.',
  },

  [USER_PERMISSIONS.technical.exploitation.pbehavior]: {
    title: 'PBehaviors',
    message: 'Canopsis periodical behaviors can be used in order to define a periods when the behavior has to be changed, e.g. for  maintenance or service range.',
  },

  [USER_PERMISSIONS.technical.exploitation.scenario]: {
    title: 'Scenarios',
    message: 'The Canopsis scenarios can be used to conditionally trigger various types of actions on alarms.',
  },

  [USER_PERMISSIONS.technical.exploitation.snmpRule]: {
    title: 'SNMP rules',
    message: 'The SNMP engine allows the processing of SNMP traps retrieved by the connector snmp2canopsis.',
  },

  [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: {
    title: 'Ticket declaration rules',
    message: 'Allows you to define the rules for creating incident tickets for multiple incident management tools. These tickets are created by a driver from the WebUI, and Canopsis then takes care of contacting the incident management tool to record the ticket.',
  },

  [USER_PERMISSIONS.technical.exploitation.linkRule]: {
    title: 'Links generator',
    message: 'Allows you to define the rules for associating links with alarms',
  },

  /**
   * Admin access
   */
  [USER_PERMISSIONS.technical.permission]: {
    title: 'Rights',
    message: 'Business, View, Technical rights are applied only for UI users. API rights are applied only for API users.',
  },
  [USER_PERMISSIONS.technical.role]: {
    title: 'Roles',
    message: 'Canopsis User Role Management',
  },
  [USER_PERMISSIONS.technical.user]: {
    title: 'Users',
    message: 'Canopsis User Management',
  },

  /**
   * Admin communications
   */
  [USER_PERMISSIONS.technical.broadcastMessage]: {
    title: 'Broadcast messages',
    message: 'The Canopsis broadcasting messages can be used for displaying banners and information messages that will appear in the Canopsis interface.',
  },
  [USER_PERMISSIONS.technical.playlist]: {
    title: 'Playlists',
    message: 'Playlists can be used for the views customization which can be displayed one after another with an associated delay.',
  },

  /**
   * Admin general
   */
  [USER_PERMISSIONS.technical.parameters]: {
    title: 'Parameters',
    message: 'Various and varied parameters of the Canopsis interface',
  },
  [USER_PERMISSIONS.technical.healthcheck]: {
    title: 'Healthcheck',
    message: 'The Healthcheck feature is the dashboard with states and errors indications of all systems included to the Canopsis.',
  },
  [USER_PERMISSIONS.technical.engine]: {
    title: 'Engines',
    message: 'This page contains the information about the sequence and configuration of engines. To work properly, the chain of engines must be continuous.',
  },
  [USER_PERMISSIONS.technical.kpi]: {
    title: 'KPI',
    message: 'Allows you to present indicators in the form of graphs',
  },
  [USER_PERMISSIONS.technical.map]: {
    title: 'Maps',
    message: 'Module to define and display maps (geographic, logical, mermaid, etc.) via the "Map" widget.',
  },
  [USER_PERMISSIONS.technical.maintenance]: {
    title: 'Maintenance mode',
    message: 'Switches Canopsis to maintenance mode. All users except administrators will be logged out.',
  },
  [USER_PERMISSIONS.technical.tag]: {
    title: 'Tags management',
    message: 'Allows you to define rules for assigning tags to alarms',
  },
  [USER_PERMISSIONS.technical.storageSettings]: {
    title: 'Storage settings',
    message: 'Allows you to define data retention policies.',
  },
  [USER_PERMISSIONS.technical.stateSetting]: {
    title: 'State settings',
    message: 'Allows you to define methods for calculating the state/severity of components and/or services.\nUseful for modeling dependency trees.',
  },
  [USER_PERMISSIONS.technical.eventsRecord]: {
    title: 'Events records',
    message: 'Allows you to define and trigger the recording of events as soon as they arrive in the Canopsis data bus.',
  },

  /**
   * Grouped admin
   */
  [GROUPED_USER_PERMISSIONS_KEYS.planning]: {
    title: 'Planning',
    message: 'The Canopsis Planning Administration functionality can be used for the periodic behavior types customization.',
  },
  [GROUPED_USER_PERMISSIONS_KEYS.remediation]: {
    title: 'Instructions',
    message: 'The Canopsis Remediation feature is used for creation plans or instructions to correct situations.',
  },

  /**
   * Notifications
   */
  [USER_PERMISSIONS.technical.notification.instructionStats]: {
    title: 'Instruction rating',
    message: 'This page contains the statistics on the instructions execution. Users can rate instructions based on their performance.',
  },

  /**
   * Profile
   */
  [USER_PERMISSIONS.technical.profile.theme]: {
    title: 'Theme',
    message: 'Allows you to define graphical themes for the GUI.',
  },
};
