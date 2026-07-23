import { PAYLOAD_OPERATORS } from '@/constants';

/**
 * Find all variables in string.
 *
 * @param {string} text
 * @returns {RegExpMatchArray[]}
 */
export const matchPayloadVariables = (text) => {
  const match = text?.matchAll(/({{([^{}]*)}})/g);

  return (match ? Array.from(match) : []).filter(
    ([, , variable]) =>
      !PAYLOAD_OPERATORS.some((operator) => String(variable).trim().startsWith(operator)),
  );
};

/**
 * Find all operators in string.
 *
 * @param {string} text
 * @returns {RegExpMatchArray[]}
 */
export const matchPayloadOperators = (text) => {
  const match = text?.matchAll(
    /(?<variable>(?<open>{{\s?range\s?((?!}).)+?}})(?<content>((?!range|end).)+?)(?<close>{{\s?end\s?}}))/g,
  );

  return match ? Array.from(match) : [];
};

/**
 * Check is cursor inside variable.
 *
 * @param {RegExpMatchArray} group
 * @param {Number} selectionStart
 * @param {Number} selectionEnd
 * @returns {boolean}
 */
export const isSelectionIntersectWithVariable = (group, selectionStart, selectionEnd) => {
  const value = group[0];

  const startIndex = group.index;
  const endIndex = group.index + value.length;

  return (
    selectionStart >= startIndex &&
    selectionEnd > startIndex &&
    selectionStart < endIndex &&
    selectionEnd <= endIndex
  );
};

/**
 * Find selected variable.
 *
 * @param {RegExpMatchArray[]} variables
 * @param {Number} selectionStart
 * @param {Number} selectionEnd
 * @returns {RegExpMatchArray | undefined}
 */
export const findSelectedVariable = (variables = [], selectionStart, selectionEnd) =>
  variables.find((group) => isSelectionIntersectWithVariable(group, selectionStart, selectionEnd));

/**
 * Match variable by selection.
 *
 * @param {string} text
 * @param {number} selectionStart
 * @param {number} selectionEnd
 * @returns {Object}
 */
export const matchPayloadVariableBySelection = (text, selectionStart, selectionEnd) =>
  findSelectedVariable(matchPayloadVariables(text), selectionStart, selectionEnd);
