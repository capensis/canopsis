import { isObject } from 'lodash';

import uid from '@/helpers/uid';

import { PAYLOAD_VARIABLE_REGEXP } from '@/constants';

/**
 * Convert payload string to JSON with indents
 *
 * @param {string | Object} payload
 * @param {number} [indents = 4]
 * @returns {string}
 */
export const convertPayloadToJson = (payload, indents = 4) => {
  const preparedPayload = isObject(payload) ? JSON.stringify(payload) : payload;

  /**
   * Searching for all variables without quot in a string
   */
  const match = preparedPayload.matchAll(new RegExp(PAYLOAD_VARIABLE_REGEXP));

  if (match) {
    /**
     * Removing first and last symbol
     *
     * Example:
     * - in: ":{{ .Alarm.Id }},"
     * - out: "{{ .Alarm.Id }}"
     */
    const jsonVariables = Array.from(match).map((group = []) => group[1]);
    // Preparing temp variable for replace variables
    const jsonHoles = jsonVariables.map(() => `"${uid('hole_')}"`);

    // Replacing all variable on temp variable for validation
    const template = jsonVariables.reduce(
      (acc, variable, index) => acc.replace(variable, jsonHoles[index]),
      preparedPayload,
    );
    const normalizedJsonValue = JSON.stringify(JSON.parse(template), null, indents);

    // Replacing temp variable on variable
    return jsonHoles.reduce(
      (acc, hole, index) => acc.replace(hole, jsonVariables[index]),
      normalizedJsonValue,
    );
  }

  return JSON.stringify(JSON.parse(preparedPayload), null, indents);
};
