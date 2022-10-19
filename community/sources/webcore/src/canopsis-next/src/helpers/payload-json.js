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
    const arrayMatch = Array.from(match);
    const jsonVariables = arrayMatch.map((group = []) => group[0]);
    // Preparing temp variable for replace variables
    const jsonHoles = arrayMatch.map(({ index }) => {
      const previousSymbol = preparedPayload[index - 1];
      const hole = uid('hole_');

      return previousSymbol === '"' ? hole : `"${hole}"`;
    });

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

/**
 * Match variable by selection
 *
 * @param {string} text
 * @param {number} selectionStart
 * @param {number} selectionEnd
 * @returns {Object}
 */
export const matchPayloadVariableBySelection = (text, selectionStart, selectionEnd) => {
  const match = text?.matchAll(/({{([^{}]*)}})/g);

  return match && Array.from(match).find((group) => {
    const value = group[0];

    const startIndex = group.index;
    const endIndex = group.index + value.length;

    return selectionStart >= startIndex
      && selectionEnd > startIndex
      && selectionStart < endIndex
      && selectionEnd <= endIndex;
  });
};
