/**
 * Convert pbehavior type data to reason form
 *
 * @param {Object} type
 * @return {Object}
 */
export function pbehaviorReasonToForm(type = {}) {
  return {
    name: type.name || '',
    description: type.description || '',
  };
}
