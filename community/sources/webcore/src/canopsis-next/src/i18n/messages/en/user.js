import { GROUPS_NAVIGATION_TYPES, USER_METRIC_PARAMETERS } from '@/constants';

export default {
  seeProfile: 'See profile',
  selectDefaultView: 'Select default view',
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
  },
};
