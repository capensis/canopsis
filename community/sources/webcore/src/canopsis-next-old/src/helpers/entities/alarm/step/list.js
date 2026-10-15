import { groupBy } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

/**
 * Get grouped steps by date
 *
 * @param {AlarmStep[]} steps
 * @return {Object.<string, AlarmStep[]>}
 */
export const groupAlarmSteps = steps => (
  groupBy(steps, step => convertDateToString(step.t, DATETIME_FORMATS.short))
);
