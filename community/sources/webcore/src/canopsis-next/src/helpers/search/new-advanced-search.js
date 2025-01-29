import { cloneDeep } from 'lodash';

import {
  ADVANCED_SEARCH_CHIP_TYPES,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_PATTERN_FIELDS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERNS_FIELDS,
  QUICK_RANGES,
} from '@/constants';

import {
  formRuleToPatternRule,
  isDatePatternRuleField,
  isDurationPatternRuleField,
  isValueInfosPatternRuleField,
  patternRuleToForm,
} from '@/helpers/entities/pattern/form';

export const getNextType = ({ attribute, fieldType, range, operator }, type) => {
  if (type === ADVANCED_SEARCH_CHIP_TYPES.union) {
    return null;
  }

  if (!attribute || !type) {
    return ADVANCED_SEARCH_CHIP_TYPES.attribute;
  }

  switch (type) {
    case ADVANCED_SEARCH_CHIP_TYPES.attribute:
      if (isDatePatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.range;
      }

      if (isValueInfosPatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.fieldType;
      }

      if (attribute === ALARM_PATTERN_FIELDS.ticketData) {
        return ADVANCED_SEARCH_CHIP_TYPES.dictionary;
      }

      return ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ADVANCED_SEARCH_CHIP_TYPES.fieldType:
      if (fieldType === PATTERN_FIELD_TYPES.boolean) {
        return ADVANCED_SEARCH_CHIP_TYPES.value;
      }

      return ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ADVANCED_SEARCH_CHIP_TYPES.dictionary:
      return ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ADVANCED_SEARCH_CHIP_TYPES.operator:
      if (isDurationPatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.duration;
      }

      if ([PATTERN_OPERATORS.isEmpty, PATTERN_OPERATORS.isNotEmpty].includes(operator)) {
        return null;
      }

      return ADVANCED_SEARCH_CHIP_TYPES.value;

    case ADVANCED_SEARCH_CHIP_TYPES.range:
      if (range === QUICK_RANGES.custom.value) {
        return ADVANCED_SEARCH_CHIP_TYPES.rangeValue;
      }

      return null;

    default:
      return null;
  }
};

export const getFilledForAdvancedSearchFormItem = (formItem) => {
  const filled = [];

  for (let type = getNextType(formItem); type; type = getNextType(formItem, type)) {
    if (type === getNextType(formItem, type)) {
      break;
    }

    filled.push(type);
  }

  return filled;
};

export const advancedSearchItemToForm = (advancedSearchItem) => {
  const form = patternRuleToForm(advancedSearchItem);

  form.filled = getFilledForAdvancedSearchFormItem(form);
  form.rangeValue = {
    from: form.range.from,
    to: form.range.to,
  };
  form.range = form.range.type;
  form.union = null;

  return form;
};

export const getAdvancedSearchUnionItem = (union = ADVANCED_SEARCH_UNION_CONDITIONS.and) => ({
  ...advancedSearchItemToForm(),
  union,

  filled: [ADVANCED_SEARCH_CHIP_TYPES.union],
});

export const advancedSearchRulesToForm = ({ positions = [], ...patterns } = {}) => {
  const clonedPatterns = cloneDeep(patterns);

  return positions.reduce((acc, key, index) => {
    if (!clonedPatterns[key][0]?.length) {
      clonedPatterns[key].shift();

      acc.push(getAdvancedSearchUnionItem(ADVANCED_SEARCH_UNION_CONDITIONS.or));
    } else if (index) {
      acc.push(getAdvancedSearchUnionItem(ADVANCED_SEARCH_UNION_CONDITIONS.and));
    }

    acc.push(advancedSearchItemToForm(clonedPatterns[key][0]?.pop?.()));

    return acc;
  }, []);
};

export const formToAdvancedSearchRules = (form = []) => {
  let firstPatternKey = null;

  return form.reduce((acc, { union, rangeValue, range, filled, ...item }) => {
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

    preparedItem = formRuleToPatternRule(preparedItem);

    let key = PATTERNS_FIELDS.alarm;

    if (preparedItem.field.startsWith('entity')) { // TODO: changed to constants
      key = PATTERNS_FIELDS.entity;
    } else if (preparedItem.field.startsWith('v.pbehavior_info')) { // TODO: changed to constants
      key = PATTERNS_FIELDS.pbehavior;
    }

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
    positions: [],
    alarm_pattern: [],
    entity_pattern: [],
    pbehavior_pattern: [],
  });
};
