import { cloneDeep, isNumber, omit } from 'lodash';

import { BASIC_ENTITY_TYPES, ENTITIES_STATES } from '@/constants';

import { infosToArray } from './shared/common';

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
 * @property {Array} depends
 * @property {Array} impact
 * @property {number} sli_avail_state
 * @property {Array} changeable_depends
 * @property {Array} changeable_impact
 * @property {Object} infos
 * @property {number} [idle_since]
 * @property {Object} [coordinates]
 */

/**
 * @typedef {Entity} EntityForm
 * @property {Array} disabled_impact
 * @property {Array} disabled_depends
 */

/**
 *
 * @param {Entity} entity
 * @return {EntityForm}
 */
export const entityToForm = (entity = {}) => {
  const depends = entity.depends ? cloneDeep(entity.depends) : [];
  const impact = entity.impact ? cloneDeep(entity.impact) : [];
  const changeableImpact = entity.changeable_impact ?? [];
  const changeableDepends = entity.changeable_depends ?? [];

  return {
    impact,
    depends,
    name: entity.name ?? '',
    description: entity.description ?? '',
    type: entity.type ?? BASIC_ENTITY_TYPES.connector,
    enabled: entity.enabled ?? true,
    disabled_impact: impact.filter(id => !changeableImpact.includes(id)),
    disabled_depends: depends.filter(id => !changeableDepends.includes(id)),
    infos: infosToArray(entity.infos),
    impact_level: entity.impact_level,
    sli_avail_state: entity.sli_avail_state ?? ENTITIES_STATES.ok,
    coordinates: entity.coordinates ?? {
      lat: undefined,
      lng: undefined,
    },
  };
};

/**
 * Convert form to API entity object
 *
 * @param {EntityForm} form
 * @return {Entity}
 */
export const formToEntity = (form) => {
  const entity = omit(form, ['disabled_impact', 'disabled_depends', 'coordinates']);

  const disabledImpact = form.disabled_impact ?? [];
  const disabledDepends = form.disabled_depends ?? [];

  entity.impact = form.impact.filter(id => !disabledImpact.includes(id));
  entity.depends = form.depends.filter(id => !disabledDepends.includes(id));

  if (isNumber(form.coordinates.lat) && isNumber(form.coordinates.lng)) {
    entity.coordinates = form.coordinates;
  }

  return entity;
};
