import { ROUTES_NAMES } from './common';

export const DEFAULT_BROADCAST_MESSAGE_COLOR = '#e75e40';

export const BROADCAST_MESSAGES_STATUSES = {
  active: 0,
  pending: 1,
  expired: 2,
};

export const BROADCAST_MESSAGE_VIEWS = {
  login: 'login',
  exploitation: 'exploitation',
  administration: 'administration',
  notifications: 'notifications',
  profile: 'profile',
  allViews: 'all-views',
  allPlaylists: 'all-playlists',
};

export const ROUTES_NAMES_TO_BROADCAST_MESSAGES = {
  [ROUTES_NAMES.login]: BROADCAST_MESSAGE_VIEWS.login,
  [ROUTES_NAMES.home]: BROADCAST_MESSAGE_VIEWS.login,
  [ROUTES_NAMES.error]: BROADCAST_MESSAGE_VIEWS.login,

  [ROUTES_NAMES.exploitationDeclareTicketRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationDynamicInfos]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationPbehaviors]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationEventFilters]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationSnmpRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationMetaAlarmRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationScenarios]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationIdleRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationFlappingRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationResolveRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationLinkRules]: BROADCAST_MESSAGE_VIEWS.exploitation,
  [ROUTES_NAMES.exploitationExternalDataTables]: BROADCAST_MESSAGE_VIEWS.exploitation,

  [ROUTES_NAMES.adminRights]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminUsers]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminRoles]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminParameters]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminBroadcastMessages]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminPlaylists]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminPlanning]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminRemediation]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminHealthcheck]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminKPI]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminMaps]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminTags]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminStorageSettings]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminStateSettings]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminEventsRecords]: BROADCAST_MESSAGE_VIEWS.administration,
  [ROUTES_NAMES.adminJobs]: BROADCAST_MESSAGE_VIEWS.administration,

  [ROUTES_NAMES.notifications]: BROADCAST_MESSAGE_VIEWS.notifications,

  [ROUTES_NAMES.profilePatterns]: BROADCAST_MESSAGE_VIEWS.profile,
  [ROUTES_NAMES.profileThemes]: BROADCAST_MESSAGE_VIEWS.profile,

  [ROUTES_NAMES.view]: BROADCAST_MESSAGE_VIEWS.allViews,
  [ROUTES_NAMES.viewKiosk]: BROADCAST_MESSAGE_VIEWS.allViews,
  [ROUTES_NAMES.alarms]: BROADCAST_MESSAGE_VIEWS.allViews,

  [ROUTES_NAMES.playlist]: BROADCAST_MESSAGE_VIEWS.allPlaylists,
};
