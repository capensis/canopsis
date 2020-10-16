import { cloneDeep, pick } from 'lodash';

import { HEARTBEAT_DURATION_UNITS } from '@/constants';

/**
 * Convert heartbeat to form
 *
 * @param {Object} [heartbeat = {}]
 * @returns {Object}
 */
export function heartbeatToForm(heartbeat = {}) {
  const regex = new RegExp(`^(\\d+)(${Object.values(HEARTBEAT_DURATION_UNITS).join('|')})$`);
  const { expected_interval: expectedInterval = '' } = heartbeat;

  const [, interval = '', unit = ''] = expectedInterval.match(regex) || [];

  return {
    name: heartbeat.name || '',
    description: heartbeat.description || '',
    output: heartbeat.output || '',
    expectedInterval: {
      interval,
      unit,
    },
    pattern: heartbeat.pattern ? cloneDeep(heartbeat.pattern) : {},
  };
}

/**
 * Convert form to heartbeat
 *
 * @param {Object} form
 * @returns {Object}
 */
export function formToHeartbeat(form) {
  return {
    ...pick(form, ['name', 'description', 'output', 'pattern']),

    expected_interval: `${form.expectedInterval.interval}${form.expectedInterval.unit}`,
  };
}
