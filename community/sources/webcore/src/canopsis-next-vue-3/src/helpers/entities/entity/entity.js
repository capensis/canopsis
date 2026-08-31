import { ENTITY_TYPES } from '@/constants';

/**
 * Check if an entity may have a state setting.
 *
 * @param {Entity} entity
 * @returns {boolean}
 */
export const hasStateSetting = entity => [ENTITY_TYPES.component, ENTITY_TYPES.service].includes(entity?.type);
