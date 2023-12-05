import { isUndefined, sortBy, intersection } from 'lodash';

/**
 * Harmonize links for special category
 *
 * @param {AlarmLinks} [links = {}]
 * @param {string} [category]
 * @returns {AlarmLink[]}
 */
export const harmonizeCategoryLinks = (links = {}, category) => {
  if (isUndefined(category) || !links[category]) {
    return [];
  }

  return sortBy(links[category] ?? [], 'label');
};

/**
 * Harmonize links for all categories
 *
 * @param {AlarmLinks} [links = {}]
 * @returns {Object<string, AlarmLink[]>}
 */
export const harmonizeCategoriesLinks = (links = {}) => Object.keys(links ?? {})
  .reduce((acc, category) => {
    acc[category] = harmonizeCategoryLinks(links, category);

    return acc;
  }, {});

/**
 * Get link rule link action type
 *
 * @param {LinkRuleLink} link
 * @returns {string}
 */
export const getLinkRuleLinkActionType = (link = {}) => [link.rule_id, link.icon_name, link.label].join('.');

/**
 * Get flatten alarm links
 *
 * @param {AlarmLinks} [links = {}]
 * @returns {AlarmLink[]}
 */
export const harmonizeLinks = (links = {}) => sortBy(Object.values(links).flat(), 'label');

/**
 * Get filtered links for alarms
 *
 * @param {Alarm[]} alarms
 * @returns {AlarmLink[]}
 */
export const harmonizeAlarmsLinks = (alarms = []) => {
  const links = alarms
    .map(alarm => harmonizeLinks(alarm.links).filter(link => !!link.rule_id && !link.single));

  if (links.length === 0) {
    return [];
  }

  const linksByKeys = {};
  const alarmsRuleIds = [];

  /**
   * We are checking intersection links rule_ids in all alarms
   */
  links.forEach((itemLinks) => {
    const ruleIdsMap = {};

    itemLinks.forEach((link) => {
      const key = getLinkRuleLinkActionType(link);

      if (!linksByKeys[key]) {
        linksByKeys[key] = link;
      }

      ruleIdsMap[link.rule_id] = true;
    });

    alarmsRuleIds.push(Object.keys(ruleIdsMap));
  });

  const availableRuleIds = intersection(...alarmsRuleIds).reduce((acc, ruleId) => {
    acc[ruleId] = true;

    return acc;
  }, {});

  return sortBy(
    Object.values(linksByKeys).filter(link => availableRuleIds[link.rule_id]),
    'label',
  );
};
