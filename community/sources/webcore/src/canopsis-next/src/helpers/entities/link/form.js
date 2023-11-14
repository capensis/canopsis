import { omit } from 'lodash';

import {
  LINK_RULE_TYPES,
  PATTERNS_FIELDS,
  LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE,
  LINK_RULE_ACTIONS,
} from '@/constants';

import { uid } from '@/helpers/uid';
import { externalDataToForm, formToExternalData } from '@/helpers/entities/shared/external-data/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * @typedef {'alarm' | 'entity'} LinkRuleType
 */

/**
 * @typedef {'open' | 'copy'} LinkRuleAction
 */

/**
 * @typedef {Object} LinkRuleLink
 * @property {string} icon_name
 * @property {string} label
 * @property {string} url
 * @property {boolean} [hide_in_menu]
 * @property {boolean} [single]
 * @property {string} [rule_id]
 * @property {string} [category]
 * @property {LinkRuleAction} [action]
 */

/**
 * @typedef {FilterPatterns} LinkRule
 * @property {string} name
 * @property {LinkRuleType} type
 * @property {boolean} enabled
 * @property {string} [source_code]
 * @property {ExternalData} [external_data]
 * @property {LinkRuleLink[]} [links]
 */

/**
 * @typedef {LinkRuleLink & ObjectKey} LinkRuleLinkForm
 */

/**
 * @typedef {Object} LinkRuleForm
 * @property {string} name
 * @property {LinkRuleType} type
 * @property {boolean} enabled
 * @property {string} source_code
 * @property {LinkRuleLinkForm[]} links
 * @property {ExternalDataForm} external_data
 * @property {FilterPatternsForm} patterns
 */

/**
 * Convert link rule link to form
 *
 * @param {LinkRuleLink} link
 * @returns {LinkRuleLinkForm}
 */
export const linkRuleLinkToForm = (link = {}) => ({
  key: uid('link'),
  label: link.label ?? '',
  icon_name: link.icon_name ?? '',
  url: link.url ?? '',
  single: link.single ?? false,
  hide_in_menu: link.hide_in_menu ?? false,
  category: link.category ?? '',
  action: link.action ?? LINK_RULE_ACTIONS.open,
});

/**
 * Convert link rule to form
 *
 * @param {LinkRule} linkRule
 * @returns {LinkRuleForm}
 */
export const linkRuleToForm = (linkRule = {}) => ({
  name: linkRule.name ?? '',
  type: linkRule.type ?? LINK_RULE_TYPES.alarm,
  enabled: linkRule.enabled ?? true,
  source_code: linkRule.source_code ?? LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  links: (linkRule.links ?? []).map(linkRuleLinkToForm),
  external_data: externalDataToForm(linkRule.external_data),
  patterns: filterPatternsToForm(linkRule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
});

/**
 * Check if source code is default
 *
 * @param {string} [code = '']
 * @returns {boolean}
 */
export const isDefaultSourceCode = (code = '') => (
  !code || [LINK_RULE_DEFAULT_ALARM_SOURCE_CODE, LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE].includes(code)
);

/**
 * Convert link rule link to form
 *
 * @param {string} key
 * @param {LinkRuleLinkForm} form
 * @param {LinkRuleType} type
 * @returns {LinkRuleLink}
 */
export const formToLinkRuleLink = ({ key, ...form }, type = LINK_RULE_TYPES.alarm) => (
  type === LINK_RULE_TYPES.entity
    ? omit(form, ['single', 'hide_in_menu'])
    : form
);

/**
 * Convert form to link rule
 *
 * @param {FilterPatternsForm} patterns
 * @param {LinkRuleLinkForm[]} links
 * @param {string} source_code
 * @param {ExternalDataForm} externalData
 * @param {LinkRuleForm} form
 * @returns {LinkRule}
 */
export const formToLinkRule = ({ patterns, links, source_code: sourceCode, external_data: externalData, ...form }) => {
  const linkRule = {
    ...form,
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
    external_data: formToExternalData(externalData),
  };

  if (isDefaultSourceCode(sourceCode)) {
    linkRule.links = links.map(link => formToLinkRuleLink(link, form.type));
  } else {
    linkRule.source_code = sourceCode;
  }

  return linkRule;
};
