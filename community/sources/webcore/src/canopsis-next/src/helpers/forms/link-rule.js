import { LINK_RULE_TYPES, OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import uid from '../uid';
import { removeKeyFromEntities } from '../entities';

import { filterPatternsToForm, formFilterToPatterns } from './filter';
import { externalDataToForm, formToExternalData } from './shared/external-data';

/**
 * @typedef {'alarm' | 'entity'} LinkRuleType
 */

/**
 * @typedef {Object} LinkRuleLink
 * @property {string} icon_name
 * @property {string} label
 * @property {string} url
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
  source_code: linkRule.source_code ?? '',
  links: (linkRule.links ?? []).map(linkRuleLinkToForm),
  external_data: externalDataToForm(linkRule.external_data),
  patterns: filterPatternsToForm(
    linkRule,
    [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
    [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
  ),
});

/**
 * Convert form to link rule
 *
 * @param {FilterPatternsForm} patterns
 * @param {LinkRuleLinkForm[]} links
 * @param {ExternalDataForm} externalData
 * @param {LinkRuleForm} form
 * @returns {LinkRule}
 */
export const formToLinkRule = ({ patterns, links, external_data: externalData, ...form }) => {
  const linkRule = {
    ...form,
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
    external_data: formToExternalData(externalData),
  };

  if (!linkRule.source_code) {
    linkRule.links = removeKeyFromEntities(links);
  }

  return linkRule;
};
