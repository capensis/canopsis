import { ENTITY_TYPES, EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

/**
 * Get icon by entity type
 *
 * @param {EntityType} type
 * @returns {string}
 */
export const getIconByEntityType = type => (
  type === ENTITY_TYPES.service
    ? '$vuetify.icons.engineering'
    : 'person'
);

/**
 * Get icon for entity event
 *
 * @param {string} type
 */
export const getEntityEventIcon = type => EVENT_ENTITY_ICONS_BY_TYPE[type];
