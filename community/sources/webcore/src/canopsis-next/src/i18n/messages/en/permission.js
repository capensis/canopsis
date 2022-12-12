import { USER_PERMISSIONS_PREFIXES } from '@/constants';

export default {
  technical: {
    admin: 'Admin rights',
    exploitation: 'Exploitation rights',
    notification: 'Notification rights',
    profile: 'Profile rights',
  },
  business: {
    [USER_PERMISSIONS_PREFIXES.business.common]: 'Rights for common',
    [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Rights for Alarms List',
    [USER_PERMISSIONS_PREFIXES.business.context]: 'Rights for Context Explorer',
    [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Rights for Service Weather',
    [USER_PERMISSIONS_PREFIXES.business.counter]: 'Rights for Counter',
    [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Rights for Testing Weather',
    [USER_PERMISSIONS_PREFIXES.business.map]: 'Rights for Mapping',
  },
  api: {
    general: 'General',
    rules: 'Rules',
    remediation: 'Remediation',
    pbehavior: 'PBehavior',
  },
};
