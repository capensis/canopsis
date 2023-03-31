import { IDLE_RULE_ALARM_CONDITIONS, IDLE_RULE_TYPES } from '@/constants';

export default {
  timeAwaiting: 'Time awaiting',
  timeRangeAwaiting: 'Time range awaiting',
  types: {
    [IDLE_RULE_TYPES.alarm]: 'Alarm rule',
    [IDLE_RULE_TYPES.entity]: 'Entity rule',
  },
  alarmConditions: {
    [IDLE_RULE_ALARM_CONDITIONS.lastEvent]: 'No events received',
    [IDLE_RULE_ALARM_CONDITIONS.lastUpdate]: 'No state changes',
  },
};
