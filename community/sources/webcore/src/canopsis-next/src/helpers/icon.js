import { EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

/**
 * Get icon for entity event
 *
 * @param {string} type
 */
export const getEntityEventIcon = type => EVENT_ENTITY_ICONS_BY_TYPE[type];
