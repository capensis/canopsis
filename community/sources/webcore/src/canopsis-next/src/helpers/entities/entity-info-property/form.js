import { ENTITY_INFO_PROPERTY_TYPES } from '@/constants/entity-info-properties';

/**
 * @typedef {Object} EntityInfoProperty
 * @property {string} infos_key
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
  infos_key: entityInfoProperty.infos_key ?? '',
  description: entityInfoProperty.description ?? '',
  type: entityInfoProperty.type ?? ENTITY_INFO_PROPERTY_TYPES.string, // Default to string
  alias: entityInfoProperty.alias ?? '',
});
