import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} ObjectKey
 * @property {string} key
 */

/**
 * @typedef {Object} Infos
 * @property {string} name
 * @property {string} description
 * @property {string|string[]} value
 */

/**
 * @typedef {Object.<string, { description: string, value: string|string[] }>} InfosObject
 */

/**
 * Convert array with primitive values to form object
 *
 * @param {Array} array
 * @param {string} [valueKey = 'value']
 * @returns {{ key: string, [valueKey]: any }[]}
 */
export function primitiveArrayToForm(array, valueKey = 'value') {
  return array.map(value => ({ [valueKey]: value, key: uid() }));
}

/**
 * Convert form object to array with primitive values
 *
 * @param {{ key: string, [valueKey]: any }[]} array
 * @param {string} [valueKey = 'value']
 * @returns {Array}
 */
export function formToPrimitiveArray(array, valueKey = 'value') {
  return array.map(item => item[valueKey]);
}

/**
 * Default item creator for primitive array
 *
 * @returns {{value: string, key: string}}
 */
export function defaultPrimitiveArrayItemCreator() {
  return { value: '', key: uid() };
}

/**
 * Convert object infos to array
 *
 * @param {InfosObject} infos
 * @return {Infos[]}
 */
export const infosToArray = (infos = {}) => Object.entries(infos).map(([name, info]) => ({
  name,
  ...info,
}));

/**
 * Convert enabled field to form
 *
 * @param {boolean} [enabled = true]
 * @return {boolean}
 */
export const enabledToForm = (enabled = true) => enabled;
