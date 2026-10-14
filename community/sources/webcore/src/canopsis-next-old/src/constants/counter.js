import { ALARM_STATES } from '@/constants/alarm';

export const COUNTER_EXPORT_FILE_NAME_PREFIX = 'counter';

export const COUNTER_STATES_ICONS = {
  [ALARM_STATES.ok]: 'wb_sunny',
  [ALARM_STATES.minor]: 'person',
  [ALARM_STATES.major]: 'person',
  [ALARM_STATES.critical]: 'wb_cloudy',
};
