/**
 * @typedef {Object} MetaAlarmLink
 * @property {string} metaAlarm
 * @property {string} comment
 * @property {boolean} auto_resolve
 */

/**
 * @typedef {Object} MetaAlarmLinkRequest
 * @property {string} comment
 * @property {boolean} auto_resolve
 * @property {string} [id]
 * @property {string} [name]
 */

/**
 * Convert meta alarm link event to form
 *
 * @param {MetaAlarmLink} [link = {}]
 * @returns {MetaAlarmLinkForm}
 */
export const metaAlarmLinkToForm = (link = {}) => ({
  metaAlarm: link.metaAlarm ?? null,
  comment: link.comment ?? '',
  auto_resolve: link.auto_resolve ?? false,
});

/**
 * Convert form to meta alarm link event request
 *
 * @param {Object} metaAlarm
 * @param {MetaAlarmLink} form
 * @returns {MetaAlarmLinkRequest}
 */
export const formToMetaAlarmLinkRequest = ({ metaAlarm, ...form } = {}) => ({
  ...form,

  [metaAlarm?.noData ? 'name' : 'id']: metaAlarm?._id,
});
