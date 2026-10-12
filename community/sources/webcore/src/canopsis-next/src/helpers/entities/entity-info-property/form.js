import { ENTITY_INFO_PROPERTY_TYPES } from '@/constants';

/**
 * @typedef {Object} EntityInfoProperty
 * @property {string} name
 * @property {string} [description]
 * @property {number} type
 * @property {string} [alias]
 */

/**
 * @typedef {EntityInfoProperty} EntityInfoPropertyForm
 */

/**
 * Convert entity info property object to form
 *
 * @param {EntityInfoProperty} entityInfoProperty
 * @returns {EntityInfoPropertyForm}
 */
export const entityInfoPropertyToForm = (entityInfoProperty = {}) => ({
  name: entityInfoProperty.name ?? '',
  description: entityInfoProperty.description ?? '',
  type: entityInfoProperty.type ?? ENTITY_INFO_PROPERTY_TYPES.string,
  alias: entityInfoProperty.alias ?? '',
});
