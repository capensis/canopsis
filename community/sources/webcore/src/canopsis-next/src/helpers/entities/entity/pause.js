import { ENTITY_TYPES } from '@/constants';

/**
 * Returns whether an entity currently has a pause.
 *
 * @param {Object} [entity]
 * @returns {boolean}
 */
export const isEntityPaused = entity => Boolean(entity?.pause);

/**
 * An inherited pause has a source entity id. The `inherited` flag describes whether the pause
 * propagates to dependencies and therefore cannot be used to identify an inherited pause.
 *
 * @param {Object} [entity]
 * @returns {boolean}
 */
export const isEntityPauseInherited = entity => Boolean(entity?.pause?.inherited);

/**
 * Only a pause owned by the entity can be resumed directly.
 *
 * @param {Object} [entity]
 * @returns {boolean}
 */
export const isResumeAvailableForEntity = entity =>
  isEntityPaused(entity) && !isEntityPauseInherited(entity);

/**
 * Pause propagation is supported only for components and services.
 *
 * @param {Object} [entity]
 * @returns {boolean}
 */
export const isPausePropagationAvailableForEntity = entity =>
  [ENTITY_TYPES.component, ENTITY_TYPES.service].includes(entity?.type);

/**
 * Converts a pause form to the entity pause API request shape.
 *
 * @param {Object} [form]
 * @returns {{ comment: string, type: string, inherited: boolean }}
 */
export const entityPauseFormToRequest = (form = {}) => ({
  comment: form.comment ?? '',
  type: form.type ?? '',
  inherited: form.inherited ?? false,
});

/**
 * Converts a resume form to the entity resume API request shape.
 *
 * @param {Object} [form]
 * @returns {{ comment: string }}
 */
export const entityResumeFormToRequest = (form = {}) => ({
  comment: form.comment ?? '',
});
