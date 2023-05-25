/**
 * @typedef {Object} PbehaviorReason
 * @property {string} name
 * @property {string} description
 * @property {boolean} hidden
 */

/**
 * Convert pbehavior type data to reason form
 *
 * @param {PbehaviorReason} [type = {}]
 * @return {PbehaviorReason}
 */
export function pbehaviorReasonToForm(type = {}) {
  return {
    name: type.name ?? '',
    description: type.description ?? '',
    hidden: type.hidden ?? false,
  };
}
