import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

/**
 * @typedef {Object} PbehaviorType
 * @property {string} [_id]
 * @property {string} description
 * @property {string} icon_name
 * @property {string} name
 * @property {number} priority
 * @property {string} type
 * @property {string} color
 */

/**
 * @typedef {PbehaviorType} PbehaviorTypeForm
 */

/**
 * Convert pbehavior type data to type form
 *
 * @param {PbehaviorType} type
 * @return {PbehaviorTypeForm}
 */
export const pbehaviorTypeToForm = (type = {}) => ({
  name: type.name ?? '',
  description: type.description ?? '',
  type: type.type ?? PBEHAVIOR_TYPE_TYPES.active,
  /* 5 - next priority after default types */
  priority: type.priority ?? 5,
  icon_name: type.icon_name ?? '',
  color: type.color ?? '',
});
