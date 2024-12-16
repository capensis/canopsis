import { metaAlarmRuleInfosItemToForm, metaAlarmRuleTagsToForm } from '@/helpers/entities/meta-alarm/rule/form';
import { removeKeyFromEntities } from '@/helpers/array';

/**
 * @typedef {Object} MetaAlarmLink
 * @property {string} metaAlarm
 * @property {string} comment
 * @property {string} component
 * @property {string} resource
 * @property {boolean} auto_resolve
 * @property {MetaAlarmRuleTags} tags
 * @property {MetaAlarmRuleInfosItem[]} infos
 */

/**
 * @typedef {MetaAlarmLink} MetaAlarmLinkForm
 * @property {MetaAlarmRuleInfosItemForm[]} infos
 */

/**
 * @typedef {Object} MetaAlarmLinkRequest
 * @property {string} comment
 * @property {string} component
 * @property {string} resource
 * @property {boolean} auto_resolve
 * @property {MetaAlarmRuleTags} tags
 * @property {MetaAlarmRuleInfosItem[]} infos
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
  component: link.component ?? '',
  resource: link.resource ?? '',
  auto_resolve: link.auto_resolve ?? false,
  tags: metaAlarmRuleTagsToForm(link.tags || {}),
  infos: (link.infos || []).map(metaAlarmRuleInfosItemToForm),
});

/**
 * Convert form to meta alarm link event request
 *
 * @param {MetaAlarmRuleInfosItemForm[]} infos
 * @param {string | Object} metaAlarm
 * @param {MetaAlarmLinkForm} form
 * @returns {{infos: Array}}
 */
export const formToMetaAlarmLinkRequest = ({ infos, metaAlarm, ...form } = {}) => {
  const metaAlarmLinkRequest = {
    ...form,

    infos: removeKeyFromEntities(infos),
  };

  if (metaAlarm?._id) {
    metaAlarmLinkRequest.id = metaAlarm?._id;
  } else {
    metaAlarmLinkRequest.name = metaAlarm;
  }

  return metaAlarmLinkRequest;
};
