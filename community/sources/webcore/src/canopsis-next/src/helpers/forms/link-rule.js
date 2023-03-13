import { omit } from 'lodash';

import {
  LINK_RULE_TYPES,
  OLD_PATTERNS_FIELDS,
  PATTERNS_FIELDS,
  LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  LINK_RULE_DEFAULT_ENTITY_SOURCE_CODE,
} from '@/constants';

import uid from '../uid';

import { filterPatternsToForm, formFilterToPatterns } from './filter';
import { externalDataToForm, formToExternalData } from './shared/external-data';
import { enabledToForm } from './shared/common';

/**
 * @typedef {'alarm' | 'entity'} LinkRuleType
 */

/**
 * @typedef {Object} LinkRuleLink
 * @property {string} icon_name
 * @property {string} label
 * @property {string} url
 * @property {boolean} [single]
 * @property {string} [rule_id]
 * @property {string} [category]
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
 * @typedef {LinkRuleLink} LinkRuleLinkForm
 * @property {string} key
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
  category: link.category ?? '',
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
  enabled: enabledToForm(linkRule.enabled),
  source_code: linkRule.source_code ?? LINK_RULE_DEFAULT_ALARM_SOURCE_CODE,
  links: (linkRule.links ?? []).map(linkRuleLinkToForm),
  external_data: externalDataToForm(linkRule.external_data),
  patterns: filterPatternsToForm(
    linkRule,
    [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
    [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
  ),
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
    ? omit(form, ['single'])
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
