import { ENTITY_TYPES } from '@/constants';

export const isEntityComponent = entity => entity.type === ENTITY_TYPES.component;

export const isEntityResource = entity => entity.type === ENTITY_TYPES.resource;

/**
 * Check if entity may have state setting
 *
 * @param {Entity} entity
 * @return {boolean}
 */
export const hasStateSetting = entity => [ENTITY_TYPES.component, ENTITY_TYPES.service].includes(entity?.type);
