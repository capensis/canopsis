import {
  cloneDeep,
  isEmpty,
  isString,
  isNaN,
  isUndefined,
  mapKeys,
  pick,
  omit,
} from 'lodash';

import {
  ALARM_ADVANCED_SEARCH_CHIP_TYPES,
  ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS,
  ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX,
  ALARM_PATTERN_FIELDS,
  ENTITY_SEARCH_NUMBER_ATTRIBUTES,
  PATTERN_CONDITIONS,
  PATTERN_FIELD_TYPES,
  ADVANCED_SEARCH_FIELDS,
  PBEHAVIOR_PATTERN_PREFIX,
  PATTERN_OPERATORS_WITHOUT_VALUE,
  ADVANCED_SEARCH_FIELDS_TO_COMPARISON,
  ADVANCED_SEARCH_QUERY_FIELDS,
  ALARM_SEARCH_NUMBER_ATTRIBUTES,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_DATE_OPERATORS,
  PATTERN_OPERATORS,
  PBEHAVIOR_SEARCH_NUMBER_ATTRIBUTES,
} from '@/constants';

import { uuid } from '@/helpers/uuid';
import {
  formRuleToPatternRule,
  isArrayOperator,
  isDatePatternRuleField,
  isDurationPatternRuleField,
  isInfosPatternRuleField,
  isIntervalDateOperator,
  isValueInfosPatternRuleField,
  patternRuleToForm,
} from '@/helpers/entities/pattern/form';
import { isPickEqual } from '@/helpers/collection';
import { formToPrimitiveArray } from '@/helpers/entities/shared/form';

/**
 * Adds prefix to all keys of an attributes map.
 *
 * @param {Object} [attributesMap = {}] - The attributes map to transform.
 * @param {string} [prefix = ''] - The prefix to add to each key.
 * @returns {Object} - New object with prefixed keys.
 */
export const addPrefixToAttributesMap = (attributesMap = {}, prefix = '') => (
  mapKeys(attributesMap, (value, key) => `${prefix}${key}`)
);

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
 * | 'rangeValueDate'
 * | 'rangeValuePeriod'
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
 * @property {{ from: number, to: number }} rangeValueDate
 * @property {{ from: string, to: string }} rangeValuePeriod
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
 * Returns the type of range value item based on the operator.
 *
 * @param {string} value - The operator of the item.
 * @returns {AdvancedSearchChipType} The type of range value item.
 */
export const getRangeValueItemType = value => ({
  [PATTERN_OPERATORS.inRangeDates]: ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValueDate,
  [PATTERN_OPERATORS.inRangePeriod]: ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValuePeriod,
}[value]);

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
  { attribute, fieldType, operator, text, alias } = {},
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
      if (isDatePatternRuleField(attribute)) {
        return isIntervalDateOperator(operator)
          ? ALARM_ADVANCED_SEARCH_CHIP_TYPES.range
          : getRangeValueItemType(operator);
      }

      if (isDurationPatternRuleField(attribute)) {
        return ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration;
      }

      if (PATTERN_OPERATORS_WITHOUT_VALUE.includes(operator)) {
        return null;
      }

      return ALARM_ADVANCED_SEARCH_CHIP_TYPES.value;

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
  const rangeValueKey = getRangeValueItemType(formItem.operator);

  if (rangeValueKey) {
    formItem[rangeValueKey] = {
      from: formItem.range.from,
      to: formItem.range.to,
    };
  }

  if (isArrayOperator(formItem.operator)) {
    formItem.value = formToPrimitiveArray(formItem.value);
  }

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

    if (key === ADVANCED_SEARCH_FIELDS.entity && !formItem.alias) {
      formItem.attribute = `${ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity}${formItem.attribute}`;
    } else if (key === ADVANCED_SEARCH_FIELDS.pbehavior) {
      formItem.attribute = `${ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX}${formItem.attribute}`;
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
export const isAlarmEntityPatternField = (field = '') => field.startsWith(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity);

/**
 * Checks if a given field is a pbehavior pattern field.
 *
 * @param {string} [field = ''] - The field string to be checked.
 * @returns {boolean} - Returns true if the field is an entity pattern field, otherwise false.
 */
export const isPbehaviorPatternField = (field = '', prefix = '') => field.startsWith(`${prefix}${PBEHAVIOR_PATTERN_PREFIX}`);

/**
 * Checks if a given field is an alarm pattern field.
 *
 * @param {string} [field = ''] - The field string to be checked.
 * @returns {boolean} - Returns true if the field is an entity pattern field, otherwise false.
 */
export const isAlarmPatternField = (field = '') => field && !isAlarmEntityPatternField(field) && !isPbehaviorPatternField(field);

/**
 * Transforms a form structure into an advanced search pattern structure.
 *
 * @param {AdvancedSearchForm} [form = []] - The form array to be transformed.
 * @param {boolean} [alarmPattern = true] - Whether to include the alarm pattern.
 * @returns {AdvancedSearch} - The structured advanced search pattern object, including positions
 *                     and categorized pattern arrays (alarm, entity, pbehavior).
 */
export const formToAdvancedSearch = (form = [], alarmPattern = false) => {
  let firstPatternKey = null;

  return form.reduce((acc, { union, range, filled, text, ...item }) => {
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

    let preparedItem = omit(item, ['rangeValueDate', 'rangeValuePeriod']);

    if (range) {
      const rangeValueKey = getRangeValueItemType(item.operator);

      preparedItem.range = {
        type: range,

        ...item[rangeValueKey],
      };
    }

    const pbehaviorPrefix = alarmPattern ? ALARM_ADVANCED_SEARCH_PBEHAVIOR_PATTERN_PREFIX : '';

    let key = alarmPattern ? ADVANCED_SEARCH_FIELDS.alarm : ADVANCED_SEARCH_FIELDS.search;

    if (alarmPattern) {
      if (isAlarmEntityPatternField(preparedItem.attribute) || item?.alias) {
        key = ADVANCED_SEARCH_FIELDS.entity;
        preparedItem.attribute = preparedItem.attribute.replace(/^entity\./, '');
      } else if (isPbehaviorPatternField(preparedItem.attribute, pbehaviorPrefix)) {
        key = ADVANCED_SEARCH_FIELDS.pbehavior;
        preparedItem.attribute = preparedItem.attribute.replace(pbehaviorPrefix, '');
      }
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
    [ADVANCED_SEARCH_FIELDS.alarm]: [],
    [ADVANCED_SEARCH_FIELDS.entity]: [],
    [ADVANCED_SEARCH_FIELDS.pbehavior]: [],
    [ADVANCED_SEARCH_FIELDS.search]: [],
  });
};

/**
 * Checks if an advanced search object is empty by evaluating specific fields.
 *
 * @param {AdvancedSearch} [search = {}] - The alarm search object to evaluate.
 * @returns {boolean} - Returns true if all specified fields in the search object are empty, otherwise false.
 */
export const isEmptyAdvancedSearch = (search = {}) => (
  Object.values(pick(search, ADVANCED_SEARCH_FIELDS_TO_COMPARISON))
    .map(isEmpty)
    .every(Boolean)
);

/**
 * Checks if an advanced search pattern is empty.
 *
 * @param {Object|Array|string} [pattern = {}] - The pattern object, array, or JSON string to evaluate.
 * @returns {boolean} - Returns true if the pattern is empty, otherwise false.
 */
export const isEmptyAdvancedSearchPattern = (pattern = {}) => isEmpty(pattern) || pattern === '[]';

/**
 * Prepares a query object from a base query and an advanced search object.
 * Spreads the base query, sets search text and page to 1, and adds non-empty pattern fields as JSON strings.
 * Accepts a string as search for plain text search (treated as { search }).
 *
 * @param {Object} query - The base query object to merge
 * @param {AdvancedSearch|string} search - The advanced search object or plain search string
 * @returns {Object} - The new query object ready to be applied
 */
export const prepareQueryWithAdvancedSearch = (query = {}, search = {}) => {
  const newQuery = omit(query, ADVANCED_SEARCH_QUERY_FIELDS);

  newQuery.page = 1;

  ADVANCED_SEARCH_QUERY_FIELDS.forEach((field) => {
    const value = search?.[field];

    if (!isEmptyAdvancedSearchPattern(value)) {
      newQuery[field] = isString(value) ? value : JSON.stringify(value);
    }
  });

  return newQuery;
};

/**
 * Prepares a query object with all advanced search fields removed and page reset to 1.
 *
 * @param {Object} [query = {}] - The source query object
 * @returns {Object} - A new query object without search, alarm_pattern, entity_pattern, pbehavior_pattern,
 *                    search_pattern, and with page set to 1
 */
export const prepareQueryWithoutAdvancedSearch = (query = {}) => ({
  ...omit(query, ADVANCED_SEARCH_QUERY_FIELDS),

  page: 1,
});

/**
 * Compares two alarm search objects for equality based on specific fields or their unique identifiers.
 *
 * @param {AdvancedSearch & { _id: string }} firstSearch - The first alarm search object to compare.
 * @param {AdvancedSearch & { _id: string }} secondSearch - The second alarm search object to compare.
 * @returns {boolean} - Returns true if the searches are equal based on their IDs or specified fields, otherwise false.
 */
export const isEqualAdvancedSearches = (firstSearch, secondSearch) => (
  firstSearch?._id === secondSearch?._id
  || isPickEqual(firstSearch, secondSearch, ADVANCED_SEARCH_FIELDS_TO_COMPARISON)
);

/**
 * Merges a search into the saved searches array: updates existing equal search or prepends the new one.
 *
 * @param {Array<AdvancedSearch & { _id: string }>} savedSearches - The current saved searches
 * @param {AdvancedSearch & { _id: string }} search - The search to merge
 * @returns {Array<AdvancedSearch & { _id: string }>} - The updated searches array
 */
export const mergeSearchIntoSavedSearches = (savedSearches, search) => {
  let found = false;
  const updatedSearches = savedSearches.map((value) => {
    if (isEqualAdvancedSearches(value, search)) {
      found = true;

      return { ...search, pinned: value.pinned || search.pinned };
    }

    return value;
  });

  return found ? updatedSearches : [search, ...savedSearches];
};

/**
 * Creates an advanced search object from a field/value pair.
 * Used by alarm column cells to apply a filter from a chip click.
 *
 * @param {string} field - The advanced search field name (e.g. entity.component)
 * @param {*} value - The value to filter by
 * @returns {Object} Advanced search object with _id, pinned, and rules
 */
export const createAdvancedSearchFromAlarmFieldValue = (field, value) => {
  const patternField = ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS[field];
  const preparedField = field
    .replace(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity, '')
    .replace(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.pbehavior, '');

  const pattern = [[{ field: preparedField, cond: { value, type: PATTERN_CONDITIONS.equal } }]];

  return {
    _id: uuid(),
    pinned: false,
    rules: advancedSearchToForm({
      positions: [patternField],
      [patternField]: pattern,
    }),
  };
};

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
    || ENTITY_SEARCH_NUMBER_ATTRIBUTES.includes(rule.attribute)
    || PBEHAVIOR_SEARCH_NUMBER_ATTRIBUTES.includes(rule.attribute)
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
 * Checks if the given type and value represent a range for date item.
 *
 * @param {AdvancedSearchChipType} type - The type of the item.
 * @param {string} value - The value of the item.
 * @returns {boolean} True if the type is 'range' and the value is an interval date operator, otherwise false.
 */
export const isRangeItem = (type, value) => (
  type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.range && PATTERN_DATE_OPERATORS.includes(value)
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

/**
 * Returns input attributes and preparer for number fields with a minimum value constraint.
 *
 * @param {number} [min = 0] - The minimum allowed value.
 * @returns {Object} Object containing inputAttributes (min, step, inputmode) and inputPreparer function
 *                  that normalizes the value to min when invalid or below minimum,
 *                  or returns empty string when value is empty.
 */
export const getNumberMinValueAttributes = (min = 0) => ({
  operators: PATTERN_NUMBER_OPERATORS,
  inputAttributes: {
    min,
    step: 1,
    inputmode: 'numeric',
  },
  inputPreparer: (value) => {
    if (!value && value !== min) {
      return '';
    }

    const num = Number(value);

    return (isNaN(num) || num < min) ? min : num;
  },
});
