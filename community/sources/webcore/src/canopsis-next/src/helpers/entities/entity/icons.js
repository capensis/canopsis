import { ENTITY_TYPES_ICONS, EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

/**
 * Get icon by entity type
 *
 * @param {EntityType} type
 * @returns {string}
 */
export const getIconByEntityType = type => ENTITY_TYPES_ICONS[type] ?? 'perm_identity';

/**
 * Get icon for entity event
 *
 * @param {string} type
 */
export const getEntityEventIcon = type => EVENT_ENTITY_ICONS_BY_TYPE[type];
