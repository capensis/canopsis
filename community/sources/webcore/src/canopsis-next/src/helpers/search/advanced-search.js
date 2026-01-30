import { cloneDeep, isEmpty, isUndefined, pick } from 'lodash';

import {
  ALARM_ADVANCED_SEARCH_CHIP_TYPES,
  ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_PATTERN_FIELDS,
  PATTERN_FIELD_TYPES,
  PATTERNS_FIELDS,
  QUICK_RANGES,
  PATTERN_OPERATORS_WITHOUT_VALUE,
  ALARM_SEARCH_FIELDS_TO_COMPARISON,
  ALARM_SEARCH_NUMBER_ATTRIBUTES,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_DURATION_OPERATORS,
} from '@/constants';

import {
  formRuleToPatternRule,
  isArrayOperator,
  isDatePatternRuleField,
  isDurationPatternRuleField,
  isInfosPatternRuleField,
  isValueInfosPatternRuleField,
  patternRuleToForm,
} from '@/helpers/entities/pattern/form';
import { isPickEqual } from '@/helpers/collection';

/**
 * @typedef {'alarm_pattern' | 'entity_pattern' | 'pbehavior_pattern'} AdvancedSearchRuleItemPosition
 */

/**
 * @typedef {AdvancedSearchRuleItemPosition[]} AdvancedSearchPositions
 */

/**
 * @typedef {Object} AdvancedSearchPatterns
 * @property {PatternGroups} alarm_pattern
 * @property {PatternGroups} entity_pattern
 * @property {PatternGroups} pbehavior_pattern
 */

/**
 * @typedef {AdvancedSearchPatterns & AdvancedSearchPositions} AdvancedSearch
 * @property {string} search
 */

/**
 * @typedef {
 * 'union'
 * | 'attribute'
 * | 'operator'
 * | 'fieldType'
 * | 'dictionary'
 * | 'range'
 * | 'rangeValue'
 * | 'value'
 * | 'duration'
 * | 'text'
 * } AdvancedSearchChipType
 */

/**
 * @typedef {'OR' | 'AND'} AdvancedSearchUnion
 */

/**
 * @typedef {PatternRuleForm} AdvancedSearchFormItem
 * @property {AdvancedSearchUnion} union
 * @property {string} range
 * @property {string} text
 * @property {{ from: number, to: number }} rangeValue
 * @property {AdvancedSearchChipType[]} filled
 */

/**
 * @typedef {AdvancedSearchFormItem[]} AdvancedSearchForm
 */

/**
 * Determines the initial form item type based on the provided item and union flag.
 *
 * @param {AdvancedSearchFormItem} item - The item object to evaluate.
 * @param {boolean} [union = false] - A flag indicating whether the union type should be returned.
 * @returns {string | null}
 */
export const getInitialFormItemType = (item, union = false) => {
  if (item?.attribute || item?.union || item?.text) {
    return null;
  }

  return union ? ALARM_ADVANCED_SEARCH_CHIP_TYPES.union : ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute;
};

/**
 * Determines the next chip type for an advanced search form item based on the current type and attributes.
 *
 * @param {Object} [params = {}] - The parameters for determining the next chip type.
 * @param {string} [params.attribute] - The attribute of the rule item.
 * @param {string} [params.fieldType] - The field type of the rule item.
 * @param {string} [params.range] - The range of the rule item.
 * @param {string} [params.operator] - The operator of the rule item.
 * @param {string} [params.text] - The text of the rule item.
 * @param {boolean} [params.alias] - The alias of the rule item.
 * @param {AdvancedSearchChipType | null} [type = null] - The current chip type.
 * @returns {AdvancedSearchChipType | null} - The next chip type, or null if there is no next type.
 */
export const getNextForFormItemType = (
  { attribute, fieldType, range, operator, text, alias } = {},
  type = null,
) => {
  if (text || [ALARM_ADVANCED_SEARCH_CHIP_TYPES.union, ALARM_ADVANCED_SEARCH_CHIP_TYPES.text].includes(type)) {
    return null;
  }

  if (!attribute || !type) {
    return ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute;
  }

  switch (type) {
    case ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute:
      if (isDatePatternRuleField(attribute)) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.range;
      }

      if (isValueInfosPatternRuleField(attribute) || alias) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.fieldType;
      }

      if (attribute === ALARM_PATTERN_FIELDS.ticketData) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.dictionary;
      }

      return ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ALARM_ADVANCED_SEARCH_CHIP_TYPES.fieldType:
      if (fieldType === PATTERN_FIELD_TYPES.boolean) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.value;
      }

      if (fieldType === PATTERN_FIELD_TYPES.timestamp) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.range;
      }

      return ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ALARM_ADVANCED_SEARCH_CHIP_TYPES.dictionary:
      return ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator:
      if (isDurationPatternRuleField(attribute)) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration;
      }

      if (PATTERN_OPERATORS_WITHOUT_VALUE.includes(operator)) {
        return null;
      }

      return ALARM_ADVANCED_SEARCH_CHIP_TYPES.value;

    case ALARM_ADVANCED_SEARCH_CHIP_TYPES.range:
      if (range === QUICK_RANGES.custom.value) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValue;
      }

      return null;

    default:
      return null;
  }
};

/**
 * Generates an array of chip types for an advanced search form item, filling it based on the item's attributes.
 *
 * @param {Object} formItem - The form item for which to generate the filled array.
 * @returns {Array<AdvancedSearchChipType>} - An array of chip types representing the filled state of the form item.
 */
export const getFilledArrayForAdvancedSearchFormItem = (formItem) => {
  const filled = [];
  let type = getNextForFormItemType(formItem);

  while (type) {
    const newType = getNextForFormItemType(formItem, type);

    if (newType === type) {
      break;
    }

    if (type) {
      filled.push(type);
    }

    type = newType;
  }

  return filled;
};

/**
 * Converts an advanced search rule item into a form item.
 *
 * @param {Object} [advancedSearchRuleItem = {}] - The advanced search rule item to be converted.
 * @returns {AdvancedSearchFormItem} - The converted form item with filled attributes and range values.
 */
export const advancedSearchRuleItemToFormItem = (advancedSearchRuleItem = {}) => {
  const formItem = patternRuleToForm(advancedSearchRuleItem);

  formItem.rangeValue = {
    from: formItem.range.from,
    to: formItem.range.to,
  };
  formItem.range = formItem.range.type;
  formItem.union = null;
  formItem.text = '';

  if (formItem.dictionary) {
    formItem.attribute = [formItem.attribute, formItem.dictionary].join('.');
    formItem.dictionary = '';
  }

  if (formItem.field) {
    formItem.attribute = [formItem.attribute, formItem.field].join('.');
    formItem.field = '';
  }

  formItem.filled = getFilledArrayForAdvancedSearchFormItem(formItem);

  return formItem;
};

/**
 * Creates an advanced search union item with a specified union condition.
 *
 * @param {AdvancedSearchUnion} [union = ADVANCED_SEARCH_UNION_CONDITIONS.and] - The union condition to be applied.
 * @returns {AdvancedSearchFormItem} - The advanced search union item with the specified
 *                                     union condition and filled attributes.
 */
export const getAdvancedSearchUnionItem = (union = ADVANCED_SEARCH_UNION_CONDITIONS.and) => ({
  ...advancedSearchRuleItemToFormItem(),
  union,

  filled: [ALARM_ADVANCED_SEARCH_CHIP_TYPES.union],
});

/**
 * Converts advanced search into a form structure.
 *
 * @param {AdvancedSearch} params - The parameters for the conversion.
 * @param {string} params.search - The text search.
 * @param {AdvancedSearchPositions} params.positions - The positions of the patterns.
 * @param {AdvancedSearchPatterns} params.patterns - The advanced search patterns.
 * @returns {AdvancedSearchForm} - The form structure representing the advanced search rules.
 */
export const advancedSearchToForm = ({ search = '', positions = [], ...patterns } = {}) => {
  const clonedPatterns = cloneDeep(patterns);

  if (search) {
    const item = advancedSearchRuleItemToFormItem();

    item.text = search;
    item.filled = [ALARM_ADVANCED_SEARCH_CHIP_TYPES.text];

    return [item];
  }

  return positions.reduce((acc, key, index) => {
    if (!clonedPatterns[key][0]?.length) {
      clonedPatterns[key].shift();

      acc.push(getAdvancedSearchUnionItem(ADVANCED_SEARCH_UNION_CONDITIONS.or));
    } else if (index) {
      acc.push(getAdvancedSearchUnionItem(ADVANCED_SEARCH_UNION_CONDITIONS.and));
    }

    const formItem = advancedSearchRuleItemToFormItem(clonedPatterns[key][0]?.pop?.());

    if (key === PATTERNS_FIELDS.entity && !formItem.alias) {
      formItem.attribute = ['entity', formItem.attribute].join('.');
    } else if (key === PATTERNS_FIELDS.pbehavior) {
      formItem.attribute = ['v', formItem.attribute].join('.');
    }

    acc.push(formItem);

    return acc;
  }, []);
};

/**
 * Checks if a given field is an entity pattern field.
 *
 * @param {string} [field = ''] - The field string to be checked.
 * @returns {boolean} - Returns true if the field is an entity pattern field, otherwise false.
 */
export const isEntityPatternField = (field = '') => field.startsWith(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity);

/**
 * Checks if a given field is a pbehavior pattern field.
 *
 * @param {string} [field = ''] - The field string to be checked.
 * @returns {boolean} - Returns true if the field is an entity pattern field, otherwise false.
 */
export const isPbehaviorPatternField = (field = '') => field.startsWith(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.pbehavior);

/**
 * Checks if a given field is an alarm pattern field.
 *
 * @param {string} [field = ''] - The field string to be checked.
 * @returns {boolean} - Returns true if the field is an entity pattern field, otherwise false.
 */
export const isAlarmPatternField = (field = '') => !isEntityPatternField(field) && !isPbehaviorPatternField(field);

/**
 * Transforms a form structure into an advanced search pattern structure.
 *
 * @param {AdvancedSearchForm} [form = []] - The form array to be transformed.
 * @returns {AdvancedSearch} - The structured advanced search pattern object, including positions
 *                     and categorized pattern arrays (alarm, entity, pbehavior).
 */
export const formToAdvancedSearch = (form = []) => {
  let firstPatternKey = null;

  return form.reduce((acc, { union, rangeValue, range, filled, text, ...item }) => {
    if (text) {
      acc.search = text;

      return acc;
    }

    if (union === ADVANCED_SEARCH_UNION_CONDITIONS.and || (!union && !item.attribute)) {
      return acc;
    }

    if (union === ADVANCED_SEARCH_UNION_CONDITIONS.or) {
      acc[firstPatternKey].push([]);

      return acc;
    }

    let preparedItem = { ...item };

    if (range) {
      preparedItem.range = {
        type: range,

        ...rangeValue,
      };
    }

    let key = PATTERNS_FIELDS.alarm;

    if (isEntityPatternField(preparedItem.attribute) || item?.alias) {
      key = PATTERNS_FIELDS.entity;
      preparedItem.attribute = preparedItem.attribute.replace(/^entity\./, '');
    } else if (isPbehaviorPatternField(preparedItem.attribute)) {
      key = PATTERNS_FIELDS.pbehavior;
      preparedItem.attribute = preparedItem.attribute.replace(/^v\./, '');
    }

    if (isInfosPatternRuleField(preparedItem.attribute)) {
      const splittedAttribute = preparedItem.attribute.split('.');
      const field = splittedAttribute.pop();

      preparedItem.attribute = splittedAttribute.join('.');
      preparedItem.field = field;
    }

    preparedItem = formRuleToPatternRule(preparedItem);

    if (!firstPatternKey) {
      firstPatternKey = key;
    }

    if (!acc[key].length) {
      acc[key].push([]);
    }

    acc[key][acc[key].length - 1].push(preparedItem);
    acc.positions.push(key);

    return acc;
  }, {
    search: '',
    positions: [],
    [PATTERNS_FIELDS.alarm]: [],
    [PATTERNS_FIELDS.entity]: [],
    [PATTERNS_FIELDS.pbehavior]: [],
  });
};

/**
 * Checks if an alarm search object is empty by evaluating specific fields.
 *
 * @param {AdvancedSearch} [search = {}] - The alarm search object to evaluate.
 * @returns {boolean} - Returns true if all specified fields in the search object are empty, otherwise false.
 */
export const isEmptyAlarmSearch = (search = {}) => (
  Object.values(pick(search, ALARM_SEARCH_FIELDS_TO_COMPARISON))
    .map(isEmpty)
    .every(Boolean)
);

/**
 * Compares two alarm search objects for equality based on specific fields or their unique identifiers.
 *
 * @param {AdvancedSearch & { _id: string }} firstSearch - The first alarm search object to compare.
 * @param {AdvancedSearch & { _id: string }} secondSearch - The second alarm search object to compare.
 * @returns {boolean} - Returns true if the searches are equal based on their IDs or specified fields, otherwise false.
 */
export const isEqualAlarmSearches = (firstSearch, secondSearch) => (
  firstSearch?._id === secondSearch?._id
  || isPickEqual(firstSearch, secondSearch, ALARM_SEARCH_FIELDS_TO_COMPARISON)
);

/**
 * Determines if a rule is of a number value type based on its type and operator.
 *
 * @param {AdvancedSearchFormItem} rule - The rule object to evaluate.
 * @param {AdvancedSearchChipType} type - The type of the rule to check against.
 * @returns {boolean} - Returns true if the rule is of a number value type, otherwise false.
 */
export const isNumberValueType = (rule, type) => (
  type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.value
  && PATTERN_NUMBER_OPERATORS.includes(rule.operator)
  && (
    rule.fieldType === PATTERN_FIELD_TYPES.number
    || ALARM_SEARCH_NUMBER_ATTRIBUTES.includes(rule.attribute)
  )
);

/**
 * Checks if the given type and value represent an array item.
 *
 * @param {AdvancedSearchChipType} type - The type of the item.
 * @param {string} value - The value of the item.
 * @returns {boolean} True if the type is 'operator' and the value is an array condition, otherwise false.
 */
export const isArrayItem = (type, value) => (
  type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator && isArrayOperator(value)
);

/**
 * Checks if the given type and value represent a custom range item.
 *
 * @param {AdvancedSearchChipType} type - The type of the item.
 * @param {string} value - The value of the item.
 * @returns {boolean} True if the type is 'range' and the value is 'custom', otherwise false.
 */
export const isCustomRangeItem = (type, value) => (
  type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.range && value === QUICK_RANGES.custom.value
);

/**
 * Checks if the given type and value represent a duration item.
 *
 * @param {AdvancedSearchChipType} type - The type of the item.
 * @param {string} value - The value of the item.
 * @returns {boolean} True if the type is 'operator' and the value is a duration operator, otherwise false.
 */
export const isDurationItem = (type, value) => (
  type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator && PATTERN_DURATION_OPERATORS.includes(value)
);

/**
 * Filters a list of items based on a specified condition, including the last header item before each matching item.
 *
 * @param {Array} [items = []] - The list of items to filter.
 * @param {Function} [condition = () => true] - A function that takes an item as an argument and returns a boolean
 *                                              indicating whether the item meets the condition.
 * @returns {Array} - A new array containing items that meet the condition, along with the last header item preceding
 *                    each matching item.
 */
export const filterAdvancedSearchItems = (items = [], condition = () => true) => {
  let lastHeaderIndex;

  return items.reduce((acc, item, index) => {
    if (item.header) {
      lastHeaderIndex = index;

      return acc;
    }

    if (condition(item)) {
      if (!isUndefined(lastHeaderIndex)) {
        acc.push(items[lastHeaderIndex]);
      }

      acc.push(item);
    }

    return acc;
  }, []);
};
