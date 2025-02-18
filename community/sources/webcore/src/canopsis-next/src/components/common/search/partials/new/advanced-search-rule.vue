<template>
  <v-layout
    :class="{ 'c-new-advanced-search__rule--union': union }"
    class="c-new-advanced-search__rule"
  >
    <component
      v-for="chip in chips"
      :is="chip.component"
      v-bind="chip.bind"
      :key="chip.key"
      v-on="chip.on"
    />
  </v-layout>
</template>

<script>
import { keyBy, uniq, pick } from 'lodash';
import {
  computed,
  ref,
  watch,
  nextTick,
  inject,
  onBeforeUnmount,
} from 'vue';
import { Validator } from 'vee-validate';

import {
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ADVANCED_SEARCH_CHIP_TYPES,
  PATTERN_ARRAY_OPERATORS,
  PATTERN_FIELD_TYPES,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
  PATTERN_STRING_OPERATORS,
  QUICK_RANGES,
  ADVANCED_SEARCH_VALIDATION_RULE_NAME,
} from '@/constants';

import { isArrayCondition } from '@/helpers/entities/pattern/form';
import { advancedSearchRuleItemToFormItem, getNextForFormItemType } from '@/helpers/search/new-advanced-search';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import AdvancedSearchChip from './advanced-search-chip.vue';
import AdvancedSearchRangeChip from './advanced-search-range-chip.vue';

// TODO: move to helpers
const getInitialInputTypeForRule = (rule = {}, union = false) => {
  if (rule.attribute || rule.union) {
    return null;
  }

  return union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute;
};

export default {
  components: { AdvancedSearchChip, AdvancedSearchRangeChip },
  model: {
    prop: 'rule',
    event: 'input',
  },
  props: {
    rule: {
      type: Object,
      default: () => ({}),
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    active: {
      type: Boolean,
      default: false,
    },
    union: {
      type: Boolean,
      default: false,
    },
    inputTypes: {
      type: Array,
      default: () => [
        { value: PATTERN_FIELD_TYPES.string },
        { value: PATTERN_FIELD_TYPES.number },
        { value: PATTERN_FIELD_TYPES.boolean },
        { value: PATTERN_FIELD_TYPES.stringArray },
      ],
    },
    intervalRanges: {
      type: Array,
      default: () => PATTERN_QUICK_RANGES,
    },
    allowOr: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    first: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = inject('$validator', new Validator());

    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const inputType = ref(getInitialInputTypeForRule(props.rule, props.union));
    const activeType = ref(null);

    const attributesMap = computed(() => keyBy(props.attributes, 'value'));
    const currentAttribute = computed(() => attributesMap.value[props.rule.attribute]);
    const isFinishedRule = computed(() => !inputType.value);

    const isStringFieldType = computed(() => props.rule.fieldType === PATTERN_FIELD_TYPES.string);
    const isNumberFieldType = computed(() => props.rule.fieldType === PATTERN_FIELD_TYPES.number);
    const isBooleanFieldType = computed(() => props.rule.fieldType === PATTERN_FIELD_TYPES.boolean);
    const isArrayFieldType = computed(() => props.rule.fieldType === PATTERN_FIELD_TYPES.stringArray);

    const operators = computed(() => ({
      [isStringFieldType.value]: [
        ...PATTERN_STRING_OPERATORS,

        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
      ],
      [isNumberFieldType.value]: PATTERN_NUMBER_OPERATORS,
      [isArrayFieldType.value]: PATTERN_ARRAY_OPERATORS,
    }.true ?? []));

    const preparedOperators = computed(() => (
      (currentAttribute.value?.operators ?? operators.value ?? []).map(operator => ({
        text: t(`common.operators.${operator}`),
        value: operator,
      }))
    ));

    const preparedIntervalRanges = computed(() => (
      props.intervalRanges.map(range => ({
        ...range,
        text: t(`quickRanges.types.${range.value}`),
      }))
    ));

    const preparedInputTypes = computed(() => (
      props.inputTypes.map(type => ({
        ...type,
        text: t(`common.mixedField.types.${type.value}`),
      }))
    ));

    const preparedUnionItems = computed(() => (
      Object.values(ADVANCED_SEARCH_UNION_CONDITIONS).map(value => ({
        value,
        text: value,
        disabled: value === ADVANCED_SEARCH_UNION_CONDITIONS.or && !props.allowOr,
      }))
    ));

    const preparedBooleanItems = computed(() => [
      { text: t('common.true'), value: true }, { text: t('common.false'), value: false },
    ]);

    const preparedValueItems = computed(() => ({
      [isBooleanFieldType.value]: preparedBooleanItems.value,
    }).true ?? []);

    const itemsByType = computed(() => ({
      [ADVANCED_SEARCH_CHIP_TYPES.attribute]: props.attributes,
      [ADVANCED_SEARCH_CHIP_TYPES.operator]: preparedOperators.value,
      [ADVANCED_SEARCH_CHIP_TYPES.range]: preparedIntervalRanges.value,
      [ADVANCED_SEARCH_CHIP_TYPES.fieldType]: preparedInputTypes.value,
      [ADVANCED_SEARCH_CHIP_TYPES.value]: preparedValueItems.value,
      [ADVANCED_SEARCH_CHIP_TYPES.union]: preparedUnionItems.value,
    }));

    /**
     * Navigates to the next form item type, optionally skipping the current type.
     *
     * @param {boolean} [skipType = false] - Whether to skip the current type and move to the next one.
     */
    const goToNextType = (skipType = false) => {
      inputType.value = getNextForFormItemType(props.rule, inputType.value);

      if (skipType) {
        goToNextType();

        return;
      }

      if (isFinishedRule.value) {
        emit('next', isArrayCondition(props.rule.operator));
      }
    };

    /**
     * Sets the input type for the search component.
     *
     * @param {string} type - The type of input to be set (e.g., 'attribute', 'operator').
     */
    const setInputType = type => inputType.value = type;

    /**
     * Updates the search rule's chip item based on the provided value and type.
     *
     * @param {*} value - The value to set for the specified type. This can vary depending on the type.
     * @param {string} type - The type of the chip item being updated. It should be one of the predefined
     *                        types in `ADVANCED_SEARCH_CHIP_TYPES`.
     */
    const updateChipItem = (value, type) => {
      const oldFilled = props.rule.filled ?? [];
      const typeIndex = oldFilled.indexOf(type);
      const filled = typeIndex === -1 ? oldFilled : oldFilled.slice(0, typeIndex + 1);
      const filledForRemove = typeIndex === -1 ? [] : oldFilled.slice(typeIndex + 1);

      if (type === ADVANCED_SEARCH_CHIP_TYPES.operator && isArrayCondition(value)) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.value);
      }

      if (type === ADVANCED_SEARCH_CHIP_TYPES.range && value === QUICK_RANGES.custom.value) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
      }

      updateModel({
        ...props.rule,
        ...pick(advancedSearchRuleItemToFormItem(), filledForRemove),
        [type]: value,
        filled: uniq(filled),
      });

      if ([ADVANCED_SEARCH_CHIP_TYPES.rangeValue, ADVANCED_SEARCH_CHIP_TYPES.value].includes(type)) {
        return;
      }

      setInputType(type);

      if (!isFinishedRule.value) {
        nextTick(() => goToNextType());
      }
    };

    /**
     * Updates the current search rule item with the specified value and manages the progression of filled criteria.
     *
     * @param {*} value - The value to set for the current active type. This can vary depending on the type.
     */
    const updateItem = (value) => {
      const filled = [...(props.rule.filled ?? []), inputType.value];
      const preparedRule = { ...props.rule };
      let skipType = false;

      if (inputType.value === ADVANCED_SEARCH_CHIP_TYPES.operator && isArrayCondition(value)) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.value);
        preparedRule[ADVANCED_SEARCH_CHIP_TYPES.value] = [];
        skipType = true;
      }

      if (inputType.value === ADVANCED_SEARCH_CHIP_TYPES.range && value === QUICK_RANGES.custom.value) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
        skipType = true;
      }

      updateModel({
        ...preparedRule,
        filled: uniq(filled),
        [inputType.value]: value,
      });

      if (!isFinishedRule.value) {
        nextTick(() => goToNextType(skipType));
      }
    };

    /**
     * Handles the click event on a chip and sets the active type.
     *
     * @param {string} type - The type of chip that was clicked.
     */
    const clickChip = (type) => {
      activeType.value = type;

      emit('make:active', props.rule.key);
    };

    /**
     * Handles the focus out event on a chip and resets the active type if necessary.
     *
     * @param {string} type - The type of chip that lost focus.
     */
    const focusOutChip = (type) => {
      if (type === activeType.value) {
        activeType.value = null;
        emit('reset:active', props.rule.key);
      }
    };

    /**
     * Emits a 'remove' event to notify the parent component that an item should be removed.
     */
    const remove = () => emit('remove');

    /**
     * Determines if a given type is the active type.
     *
     * @function isActiveType
     * @param {string} type - The type to check against the active type.
     * @returns {boolean} - Returns `true` if the specified type is the active type and the component is active;
     *                      otherwise, returns `false`.
     */
    const isActiveType = type => props.active && activeType.value === type;

    const getChipAttributes = ({
      input,
      first,
      closable,
      type = inputType.value,
    }) => {
      const key = `${props.rule.key}.${type}`;
      let multiple = false;
      let itemText;
      let itemValue;
      let fetchItems;

      if (type === ADVANCED_SEARCH_CHIP_TYPES.value) {
        multiple = isArrayCondition(props.rule.operator);
        itemText = currentAttribute.value?.itemText;
        itemValue = currentAttribute.value?.itemValue;
        fetchItems = currentAttribute.value?.fetchValues;
      }

      const bind = {
        disabled: props.disabled,
        items: itemsByType.value[type],
        itemText: itemText ?? 'text',
        itemValue: itemValue ?? 'value',
        allowText: first || !itemsByType.value[type]?.length,
        fetchItems,
        first,
      };

      let on = { input: updateItem };

      if (input) {
        bind.alwaysActive = true;
      } else {
        bind.active = isActiveType(type);
        bind.value = props.rule[type];
        bind.closable = closable;
        bind.multiple = multiple;
        bind.color = validator.errors.has(props.rule.key) ? 'error' : undefined;

        on = {
          input: value => updateChipItem(value, type),
          click: () => clickChip(type),
          focusout: () => focusOutChip(type),
          close: remove,
        };
      }

      return {
        key,
        component: type === ADVANCED_SEARCH_CHIP_TYPES.rangeValue
          ? 'advanced-search-range-chip'
          : 'advanced-search-chip',
        bind,
        on,
      };
    };

    const chips = computed(() => {
      const result = (props.rule.filled ?? []).map((type, index, filled) => getChipAttributes({
        type,
        first: props.first && index === 0,
        closable: index === filled.length - 1,
      }));

      if (!isFinishedRule.value && !props.rule.filled.includes(inputType.value)) {
        result.push(getChipAttributes({ input: true, first: props.first && !result.length }));
      }

      return result;
    });

    const attachValidationRule = () => validator.attach({
      name: props.rule.key,
      rules: ADVANCED_SEARCH_VALIDATION_RULE_NAME,
      getter: () => ({ rule: props.rule, finished: isFinishedRule.value }),
    });

    const detachValidationRule = () => validator.detach(props.rule.key);

    watch(() => props.union, union => (
      inputType.value = union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute
    ));

    watch(() => props.disabled, (disabled) => {
      if (disabled) {
        detachValidationRule();

        return;
      }

      attachValidationRule();
    }, { immediate: true });

    onBeforeUnmount(detachValidationRule);

    return {
      chips,
      remove,
      clickChip,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-new-advanced-search__rule {
  &:hover ::v-deep {
    .theme--light.v-chip:before {
      opacity: 0.04;
    }
    .theme--dark.v-chip:before {
      opacity: 0.08;
    }
  }

  &--union ::v-deep {
    .v-chip {
      background: var(--v-application-background-base) !important;

      &.theme--light {
        border: 1px solid var(--v-text-light-primary, rgba(0, 0, 0, 0.87));
      }

      &.theme--dark {
        border: 1px solid var(--v-text-dark-primary, #FFFFFF);
      }

      &.error {
        color: var(--v-error-base, red);
      }
    }
  }
}
</style>
