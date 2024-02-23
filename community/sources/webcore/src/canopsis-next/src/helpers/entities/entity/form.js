import { isNumber, omit } from 'lodash';

import { ENTITY_TYPES, BASIC_ENTITY_TYPES, ENTITIES_STATES } from '@/constants';

import { infosToArray } from '@/helpers/entities/shared/form';

/**
 * @typedef { 'service' | 'connector' | 'resource' | 'component' } EntityType
 */

/**
 * @typedef {Object} Entity
 * @property {string} [_id]
 * @property {string} name
 * @property {string} description
 * @property {string} type
 * @property {number} impact_level
 * @property {boolean} enabled
 * @property {number} sli_avail_state
 * @property {Object} infos
 * @property {number} [idle_since]
 * @property {Object} [coordinates]
 */

/**
 * @typedef {Entity} EntityForm
 */

/**
 * Check if entity is component
 *
 * @param {EntityType} type
 * @returns {boolean}
 */
export const isEntityComponentType = type => type === ENTITY_TYPES.component;

/**
 *
 * @param {Entity} entity
 * @return {EntityForm}
 */
export const entityToForm = (entity = {}) => ({
  _id: entity._id ?? '',
  name: entity.name ?? '',
  description: entity.description ?? '',
  type: entity.type ?? BASIC_ENTITY_TYPES.connector,
  enabled: entity.enabled ?? true,
  infos: infosToArray(entity.infos),
  impact_level: entity.impact_level,
  sli_avail_state: entity.sli_avail_state ?? ENTITIES_STATES.ok,
  coordinates: entity.coordinates ?? {
    lat: undefined,
    lng: undefined,
  },
});

/**
 * Convert form to API entity object
 *
 * @param {EntityForm} form
 * @return {Entity}
 */
export const formToEntity = (form) => {
  const entity = omit(form, ['coordinates', '_id']);

  if (isNumber(form.coordinates.lat) && isNumber(form.coordinates.lng)) {
    entity.coordinates = form.coordinates;
  }

  return entity;
};
