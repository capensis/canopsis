import { ENTITY_TYPES } from '@/constants';

export const isEntityComponent = entity => entity.type === ENTITY_TYPES.component;

export const isEntityResource = entity => entity.type === ENTITY_TYPES.resource;

export const isEntityEventsStateSettings = entity => isEntityComponent(entity) || isEntityResource(entity);
