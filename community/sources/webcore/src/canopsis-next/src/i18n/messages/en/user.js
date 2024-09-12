import { GROUPS_NAVIGATION_TYPES, USER_METRIC_PARAMETERS } from '@/constants';

export default {
  seeProfile: 'See profile',
  selectDefaultView: 'Select default view',
  displayName: 'User display name',
  firstName: 'First name',
  lastName: 'Last name',
  email: 'Email',
  language: 'User interface language',
  auth: 'Auth',
  navigationType: 'Groups navigation type',
  active: 'Session active',
  activeConnects: 'Connections count',
  navigationTypes: {
    [GROUPS_NAVIGATION_TYPES.sideBar]: 'Side bar',
    [GROUPS_NAVIGATION_TYPES.topBar]: 'Top bar',
  },
  metrics: {
    [USER_METRIC_PARAMETERS.totalUserActivity]: 'Total activity time',
    [USER_METRIC_PARAMETERS.tickets]: 'Number of tickets created',
    [USER_METRIC_PARAMETERS.ackAlarmWithoutCancel]: 'Number of acks excluding canceled acks',
    [USER_METRIC_PARAMETERS.averageUserSession]: 'Average user activity time',
    [USER_METRIC_PARAMETERS.minUserSession]: 'Min user activity time',
    [USER_METRIC_PARAMETERS.maxUserSession]: 'Max user activity time',
  },
  variables: {
    userEmail: 'User email',
    userUsername: 'User username',
    userFirstname: 'User first name',
    userLastname: 'User last name',
    userExternalId: 'User external ID',
    userSource: 'User source',
    userRole: 'User role',
  },
};
