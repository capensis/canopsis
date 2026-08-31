import { computed, provide, inject, ref, unref, watch, onMounted } from 'vue';

import {
  ADVANCED_SEARCH_GROUPS,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES,
  ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
  ALARM_ADVANCED_SEARCH_CHIP_TYPES,
  ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME,
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  BASIC_ENTITY_TYPES,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
  DEFAULT_PATTERN_FIELD_TYPES,
  ENTITY_PATTERN_FIELD_TYPES,
} from '@/constants';

import { deepKeyBy } from '@/helpers/array';
import { getOperatorsByFieldType } from '@/helpers/entities/pattern/form';
import { getNextForFormItemType } from '@/helpers/search/advanced-search';

import { useField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useDynamicInfo } from '@/hooks/store/modules/dynamic-info';
import { useEntity } from '@/hooks/store/modules/entity';
import { useService } from '@/hooks/store/modules/service';

export const ADVANCED_SEARCH_VALIDATION_CONTEXT = Symbol('advancedSearchValidationContext');

/**
 * Hook to manage fetching and processing entity information keys for advanced search.
 *
 * @returns {Object} An object containing:
 * @property {boolean} pending - A reactive boolean indicating the fetch operation's pending state.
 * @property {Array} items - A computed array of entity information keys formatted for search configurations.
 */
export const useEntityInfosKeys = () => {
  const { t } = useI18n();
  const { fetchEntityInfosKeysWithoutStore } = useService();
  const { fetchDynamicInfosKeysWithoutStore } = useDynamicInfo();

  const entityInfosKeys = ref([]);
  const dynamicInfosKeys = ref([]);
  const { pending, handler } = usePendingHandler(async () => {
    const [{ data: entityInfos }, { data: dynamicInfos }] = await Promise.all([
      fetchEntityInfosKeysWithoutStore({ params: { paginate: false } }),
      fetchDynamicInfosKeysWithoutStore({ params: { paginate: false } }),
    ]);

    entityInfosKeys.value = entityInfos;
    dynamicInfosKeys.value = dynamicInfos;
  });

  /**
   * Generates default child items for a given chip text prefix.
   *
   * @param {string} chipTextPrefix - The prefix to be used for chip text.
   * @param {Array<Object>} inputTypes - The input types of the item.
   * @param {string} definedType - The defined type of the item.
   * @returns {Array<Object>} An array of objects, each containing:
   *   - {string} value: The value of the item.
   *   - {string} text: The translated text for the item.
   *   - {string} chipText: The complete chip text including the prefix.
   */
  const getDefaultItemChildren = (chipTextPrefix = '', definedType, inputTypes) => [PATTERN_RULE_INFOS_FIELDS.name, PATTERN_RULE_INFOS_FIELDS.value].map((value) => {
      const text = t(`common.${value}`);
      const result = { value, text };

      if (chipTextPrefix) {
        result.chipText = `${chipTextPrefix}.${text}`;
      }

      if (value === PATTERN_RULE_INFOS_FIELDS.name) {
        result.operators = PATTERN_EXISTS_OPERATORS;
      } else {
        result.definedType = definedType;
        result.inputTypes = inputTypes;
      }

      return result;
    });

  /**
   * Retrieves items formatted for search configurations, with optional chip text prefix.
   *
   * @param {string} chipTextPrefix - The prefix to be used for chip text.
   * @param {Array<Object>} keys - The keys to be used for the items.
   * @returns {Array<Object>} An array of objects, each containing:
   *   - {string} value: The value of the item.
   *   - {string} chipText: The complete chip text including the prefix.
   *   - {string} text: The text for the item.
   *   - {Array<Object>} items: The default child items for the chip text.
   */
  const getItems = (chipTextPrefix, keys, fieldTypes) => keys.map(({ value, type }) => {
      const chipText = chipTextPrefix ? `${chipTextPrefix}.${value}` : undefined;
      const definedType = ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES[type];

      return {
        value,
        chipText,
        definedType,
        text: value,
        items: getDefaultItemChildren(chipText, definedType, fieldTypes),
      };
    });

  const alarmItems = computed(() => getItems(
      t(ALARM_FIELDS_TO_LABELS_KEYS[ALARM_FIELDS.infos]),
      dynamicInfosKeys.value,
      DEFAULT_PATTERN_FIELD_TYPES,
    ),
  );

  const entityItems = computed(() => getItems(
      t(ALARM_FIELDS_TO_LABELS_KEYS[ALARM_FIELDS.entityInfos]),
      entityInfosKeys.value,
      ENTITY_PATTERN_FIELD_TYPES,
    ),
  );

  onMounted(handler);

  return {
    pending,
    alarmItems,
    entityItems,
    getItems,
  };
};

/**
 * Hook to retrieve entity options for advanced search configurations.
 *
 * @returns {Object} An object containing the `getEntityOptions` function.
 * @property {Function} getEntityOptions - Function to get entity options.
 * @property {Array} getEntityOptions.operators - Array of operators applicable to the entity.
 * @property {Function} getEntityOptions.fetchValues - Function to fetch entity values with parameters.
 * @property {string} getEntityOptions.itemText - The key used to identify the display text in fetched data.
 * @property {string} getEntityOptions.itemValue - The key used to identify the value in fetched data.
 */
export const useGetEntityOptions = () => {
  const { fetchContextEntitiesListWithoutStore } = useEntity();

  const getEntityOptions = (type = Object.values(BASIC_ENTITY_TYPES)) => ({
    operators: ADVANCED_SEARCH_STRING_WITH_ONE_OF_OPERATORS,
    fetchValues: ({ params }) => fetchContextEntitiesListWithoutStore({
        params: { ...params, type },
      }),
    itemText: '_id',
    itemValue: '_id',
  });

  return {
    getEntityOptions,
  };
};

/**
 * Hook to process grouped attributes into a flat list with headers.
 *
 * @param {Object} options - Options for processing grouped attributes.
 * @param {Object} options.groups - Object mapping group keys to arrays of field names.
 * @param {Object} options.attributesMap - Map of field name to attribute config (operators, fetchValues, etc).
 * @param {Function} options.getText - (field, attributes) => string. Provides display text for the dropdown.
 * @param {string} [options.aliasGroupKey=ADVANCED_SEARCH_GROUPS.alias] - Key of the alias group.
 * @param {Array} [options.aliases] - Alias items (can be reactive).
 * @param {boolean} [options.disabled=false] - Default disabled state (can be reactive).
 * @param {Array<{headerKey: string, fields: string[], getDisabled?: () => boolean}>} [options.extraGroups=[]] -
 *   Additional groups to append (e.g. entity, pbehavior for alarm).
 * @returns {Object} An object containing the computed `attributes` array.
 * @property {Array} attributes - Array of attributes with headers for advanced search field dropdowns.
 */
export const useAdvancedSearchGroupedAttributes = ({
  groups,
  attributesMap,
  getText,
  aliasGroupKey = ADVANCED_SEARCH_GROUPS.alias,
  aliases,
  disabled = false,
  extraGroups = [],
}) => {
  const { t } = useI18n();
  /**
   * Builds an attribute item from the map for the advanced search dropdown.
   *
   * @param {string} field - Field name key in the attributes map.
   * @param {boolean} isDisabled - Whether the attribute should be disabled.
   * @param {Object} map - Attributes map (field name -> config).
   * @returns {Object|null} Attribute item with operators, text, value; null if field not in map.
   */
  const prepareFromMap = (field, isDisabled, map) => {
    const attributes = map[field];

    if (!attributes) {
      return null;
    }

    const prepared = {
      ...attributes,
      disabled: isDisabled,
      text: getText(field, attributes),
    };

    return { ...prepared, value: prepared.value ?? field };
  };

  const attributes = computed(() => {
    const unwrappedDisabled = unref(disabled);
    const unwrappedAliases = unref(aliases) ?? [];
    const map = unref(attributesMap) ?? {};

    const result = Object.entries(groups).reduce((acc, [group, groupItems]) => {
      const header = t(`advancedSearch.groups.${group}`);

      if (group === aliasGroupKey) {
        if (unwrappedAliases.length) {
          acc.push(
            { header, value: header },
            ...unwrappedAliases.map(item => ({
              alias: true,
              value: item.alias,
              text: item.alias,
              inputTypes: ENTITY_PATTERN_FIELD_TYPES,
              definedType: ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES[item.type],
              original: item,
              disabled: unwrappedDisabled,
            })),
          );
        }

        return acc;
      }

      acc.push(
        { header, value: header },
        ...groupItems.map(item => prepareFromMap(item, unwrappedDisabled, map)).filter(Boolean),
      );

      return acc;
    }, []);

    for (const { headerKey, fields, getDisabled } of extraGroups) {
      const isDisabled = getDisabled ? getDisabled() : unwrappedDisabled;

      result.push(
        {
          header: t(`advancedSearch.groups.${headerKey}`),
          value: t(`advancedSearch.groups.${headerKey}`),
        },
        ...fields.map(field => prepareFromMap(field, isDisabled, map)).filter(Boolean),
      );
    }

    return result;
  });

  return {
    attributes,
  };
};

/**
 * Hook to extend the validator with a custom rule for advanced search validation.
 *
 * @param {Object} options - Options for the validator.
 * @param {Array} options.rules - The array of rules to be validated.
 * @returns {Object} An object containing the `extendValidatorRule` function.
 */
export const useAdvancedSearchValidator = ({
  rules,
  name = 'advancedSearch',
  disabled = false,
} = {}) => {
  const ruleErrors = ref({});
  const touched = ref(false);

  const isFinishedRule = (rule) => {
    const lastFilledType = rule.filled?.at(-1);

    return !getNextForFormItemType(rule, lastFilledType);
  };

  const getRuleErrors = () => {
    const unwrappedRules = unref(rules);

    if (unref(disabled)) {
      return {};
    }

    return unwrappedRules.reduce((acc, rule, index) => {
      if (rule.attribute && !isFinishedRule(rule)) {
        acc[rule.key] = true;
      }

      if (!rule.attribute && rule.union) {
        const lastRule = unwrappedRules.at(-1);
        const preLastRule = unwrappedRules.at(-2);

        if (!lastRule?.attribute && preLastRule?.key === rule?.key) {
          acc[rule.key] = true;
        }
      }

      if (!rule.attribute && !rule.union && index !== unwrappedRules.length - 1) {
        acc[rule.key] = true;
      }

      return acc;
    }, {});
  };

  const validateRules = () => {
    ruleErrors.value = getRuleErrors();

    return !Object.keys(ruleErrors.value).length || 'Invalid advanced search';
  };

  const { validate, setValue, resetField } = useField(name, validateRules, {
    initialValue: unref(rules),
    validateOnValueUpdate: false,
  });

  watch(
    rules,
    (newRules) => {
      ruleErrors.value = {};
      touched.value = false;

      setValue(newRules, false);
    },
    { deep: true },
  );

  const errors = {
    clear: () => {
      ruleErrors.value = {};
      touched.value = false;
      resetField({ value: unref(rules), errors: [] });
    },
    has: key => touched.value && !!ruleErrors.value[key],
  };

  const validateAll = async () => {
    touched.value = true;
    setValue(unref(rules), false);

    const { valid } = await validate();

    return valid;
  };

  provide(ADVANCED_SEARCH_VALIDATION_CONTEXT, {
    errors,
  });

  return {
    errors,
    validateAll,
    extendValidatorRule: () => {},
  };
};

/**
 * A Vue composition function that manages active items for an advanced search rule.
 * It computes various properties based on the rule's field type, operators, interval ranges,
 * input types, and union conditions.
 *
 * @param {Object} options - Options for configuring the active items.
 * @param {Object} [options.rule = {}] - The rule object to evaluate.
 * @param {Array} [options.attributes = []] - List of available attributes for the rule.
 * @param {Array} [options.intervalRanges = []] - List of interval ranges for range-based rules.
 * @param {boolean} [options.allowOr = false] - Determines if the 'OR' union condition is allowed.
 * @returns {Object} - An object containing the current attribute and items by type.
 */
export const useAdvancedSearchRuleActiveItems = ({
  rule = {},
  attributes = [],
  intervalRanges = [],
  allowOr = false,
} = {}) => {
  const { t } = useI18n();

  const attributesMap = computed(() => deepKeyBy(unref(attributes), 'value'));
  const currentAttribute = computed(() => attributesMap.value[unref(rule).attribute]);

  /**
   * FIELD TYPE
   */
  const isBooleanFieldType = computed(() => unref(rule).fieldType === PATTERN_FIELD_TYPES.boolean);

  /**
   * OPERATORS ITEMS
   */
  const operators = computed(() => {
    const unwrappedRule = unref(rule);

    let result = getOperatorsByFieldType(unwrappedRule.fieldType);

    if (unwrappedRule.fieldType === PATTERN_FIELD_TYPES.string) {
      result = [...result, PATTERN_OPERATORS.isOneOf, PATTERN_OPERATORS.isNotOneOf];
    }

    return result;
  });

  const preparedOperators = computed(() => (currentAttribute.value?.operators ?? operators.value ?? []).map(operator => ({
      text: t(`common.operators.${operator}`),
      value: operator,
    })),
  );

  /**
   * INTERVAL ITEMS
   */
  const preparedIntervalRanges = computed(() => unref(intervalRanges).map(range => ({
      ...range,
      text: t(`quickRanges.types.${range.value}`),
    })),
  );

  /**
   * INPUT TYPES ITEMS
   */
  const preparedInputTypes = computed(() => {
    const { definedType, inputTypes = DEFAULT_PATTERN_FIELD_TYPES } = unref(currentAttribute) ?? {};

    return unref(inputTypes).map(type => ({
      ...type,
      text: t(`common.mixedField.types.${type.value}`),
      defined: definedType === type.value,
    }));
  });

  /**
   * UNION ITEMS
   */
  const preparedUnionItems = computed(() => Object.values(ADVANCED_SEARCH_UNION_CONDITIONS).map(value => ({
      value,
      text: value,
      disabled: value === ADVANCED_SEARCH_UNION_CONDITIONS.or && !unref(allowOr),
    })),
  );

  /**
   * BOOLEAN ITEMS
   */
  const preparedBooleanItems = computed(() => [
    { text: t('common.true'), value: true },
    { text: t('common.false'), value: false },
  ]);

  /**
   * VALUE ITEMS
   */
  const preparedValueItems = computed(
    () => currentAttribute.value?.values ??
      { [isBooleanFieldType.value]: preparedBooleanItems.value }.true ??
      [],
  );

  /**
   * ITEMS BY TYPE
   */
  const itemsByType = computed(() => ({
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute]: unref(attributes),
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.operator]: preparedOperators.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.range]: preparedIntervalRanges.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.fieldType]: preparedInputTypes.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.value]: preparedValueItems.value,
    [ALARM_ADVANCED_SEARCH_CHIP_TYPES.union]: preparedUnionItems.value,
  }));

  return {
    attributesMap,
    currentAttribute,
    itemsByType,
  };
};

/**
 * A Vue composition function that manages the attachment and detachment of validation rules
 * for an advanced search rule. It uses the injected `$validator` from `vee-validate` to
 * dynamically attach or detach validation rules based on the `disabled` state.
 *
 * @returns {Object} - An object containing the `validator` instance.
 */
export const useAttachAdvancedSearchRuleValidator = () => {
  const context = inject(ADVANCED_SEARCH_VALIDATION_CONTEXT, null);

  return {
    validator: {
      errors: {
        has: key => context?.errors.has(key),
      },
      ruleName: ALARM_ADVANCED_SEARCH_VALIDATION_RULE_NAME,
    },
  };
};
