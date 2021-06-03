import uid from '@/helpers/uid';

import { PAYLOAD_VARIABLE_REGEXP } from '@/constants';

/**
 * Convert payload string to JSON with indents
 *
 * @param {string} payload
 * @param {number} [indents]
 * @returns {string}
 */
export const convertPayloadToJson = (payload, indents) => {
  // Searching for all variables without quot in a string
  const match = payload.matchAll(new RegExp(PAYLOAD_VARIABLE_REGEXP));

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
      payload,
    );
    const normalizedJsonValue = JSON.stringify(JSON.parse(template), null, indents);

    // Replacing temp variable on variable
    return jsonHoles.reduce(
      (acc, hole, index) => acc.replace(hole, jsonVariables[index]),
      normalizedJsonValue,
    );
  }

  return JSON.stringify(JSON.parse(payload), null, indents);
};