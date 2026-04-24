/**
 * @typedef {Object} PbehaviorReason
 * @property {string} name
 * @property {string} description
 * @property {boolean} hidden
 */

/**
 * @typedef {Object} PbehaviorReasonForm
 * @property {string} name
 * @property {string} description
 * @property {boolean} visible
 */

/**
 * Convert pbehavior reason data to form
 *
 * @param {PbehaviorReason} [reason = {}]
 * @return {PbehaviorReasonForm}
 */
export function pbehaviorReasonToForm(reason = {}) {
  return {
    name: reason.name ?? '',
    description: reason.description ?? '',
    hidden: !reason.visible,
  };
}

/**
 * Convert pbehavior reason form to pbehavior reason
 *
 * @param {PbehaviorReasonForm} form
 * @return {PbehaviorReason}
 */
export const formToPbehaviorReason = (form = {}) => ({
  ...form,
  hidden: !form.visible,
});
