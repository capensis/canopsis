import {
  EVENT_FILTER_SET_TAGS_VALUE_PREFIXES,
  EVENT_FILTER_SET_TAGS_FIELDS,
  EVENT_FILTER_TYPES,
  PATTERN_OPERATORS,
} from '@/constants';

/**
 * Check event filter rule type is enrichment
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isEnrichmentEventFilterRuleType = type => type === EVENT_FILTER_TYPES.enrichment;

/**
 * Check event filter rule type is change entity
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isChangeEntityEventFilterRuleType = type => type === EVENT_FILTER_TYPES.changeEntity;

/**
 * Get set tags items from pattern in selector items format
 *
 * @param {Pattern} pattern
 * @return {{ text: string, value: string }[]}
 */
export const getSetTagsItemsFromPattern = (pattern = {}) => {
  const { groups = [] } = pattern;

  return groups.reduce((acc, group) => {
    const rules = group.rules
      .filter(({ attribute, operator }) => (
        operator === PATTERN_OPERATORS.regexp && EVENT_FILTER_SET_TAGS_FIELDS.includes(attribute)
      ))
      .map(({ attribute, dictionary, value }) => ({
        text: dictionary || attribute,
        valueForValidation: value,
        value: `${EVENT_FILTER_SET_TAGS_VALUE_PREFIXES[attribute]}${dictionary}`,
      }));

    acc.push(...rules);

    return acc;
  }, []);
};
