<template>
  <v-layout
    :class="{ 'c-new-advanced-search__group--union': union }"
    class="c-new-advanced-search__group"
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
import { keyBy, uniq, difference, pick } from 'lodash';
import {
  computed,
  ref,
  watch,
  nextTick,
  onMounted,
  onBeforeUnmount,
  inject,
} from 'vue';
import { Validator } from 'vee-validate';

import {
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ADVANCED_SEARCH_CHIP_TYPES,
  PATTERN_ARRAY_OPERATORS,
  PATTERN_CONDITIONS,
  PATTERN_FIELD_TYPES,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
  PATTERN_STRING_OPERATORS,
  QUICK_RANGES,
} from '@/constants';

import { isArrayCondition } from '@/helpers/entities/pattern/form';
import { advancedSearchRuleItemToFormItem, getNextForFormItemType } from '@/helpers/search/new-advanced-search';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import AdvancedSearchChip from './advanced-search-chip.vue';
import AdvancedSearchRangeChip from './advanced-search-range-chip.vue';

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
    activeKey: {
      type: String,
      default: '',
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
    const { updateField, updateModel } = useModelField(props, emit);
    const activeType = ref(props.union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute);

    const attributesMap = computed(() => keyBy(props.attributes, 'value'));
    const currentAttribute = computed(() => attributesMap.value[props.rule.attribute]);
    const isFinishedGroup = computed(() => !activeType.value);
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
      { text: 'True', value: 'true' }, { text: 'False', value: 'false' },
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

    const goToNextType = (steps = 1) => {
      for (let i = 0; i < steps; i += 1) {
        activeType.value = getNextForFormItemType(props.rule, activeType.value);
      }
    };

    const selectChipItem = (value, type) => {
      const newRule = { ...props.rule, [type]: value };
      const oldRuleNextType = getNextForFormItemType(props.rule, type);
      let nextType = getNextForFormItemType(newRule, type);
      if (!nextType) {
        if (oldRuleNextType) {
          const typesForClear = [oldRuleNextType];
          updateModel({
            ...props.rule,
            ...pick(advancedSearchRuleItemToFormItem(), typesForClear),
            [type]: value,
            filled: difference(props.rule.filled, typesForClear),
          });

          return;
        }

        updateField(type, value);

        return;
      }

      activeType.value = nextType;

      const typesForClear = [nextType];

      while (nextType = getNextForFormItemType(newRule, nextType)) {
        typesForClear.push(nextType);
      }

      const filled = difference(props.rule.filled, typesForClear);

      if (type === ADVANCED_SEARCH_CHIP_TYPES.operator && isArrayCondition(value)) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.value);
      }

      if (type === ADVANCED_SEARCH_CHIP_TYPES.range && value === QUICK_RANGES.custom.value) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
      }

      updateModel({
        ...props.rule,
        ...pick(advancedSearchRuleItemToFormItem(), typesForClear),
        [type]: value,
        filled: uniq(filled),
      });
    };

    const selectItem = (value) => {
      const filled = [...(props.rule.filled ?? []), activeType.value];
      const preparedRule = { ...props.rule };
      let actionTypeSteps = 1;

      if (
        activeType.value === ADVANCED_SEARCH_CHIP_TYPES.operator
        && [
          PATTERN_CONDITIONS.hasNot,
          PATTERN_CONDITIONS.hasOneOf,
          PATTERN_CONDITIONS.isOneOf,
          PATTERN_CONDITIONS.isNotOneOf,
          PATTERN_CONDITIONS.hasEvery,
        ].includes(value)
      ) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.value);
        preparedRule[ADVANCED_SEARCH_CHIP_TYPES.value] = [];
        actionTypeSteps = 2;
      }

      if (activeType.value === ADVANCED_SEARCH_CHIP_TYPES.range && value === QUICK_RANGES.custom.value) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
        actionTypeSteps = 2;
      }

      updateModel({
        ...preparedRule,
        filled: uniq(filled),
        [activeType.value]: value,
      });

      if (!isFinishedGroup.value) {
        nextTick(() => goToNextType(actionTypeSteps));
      }
    };

    const clickChip = key => emit('click:chip', key);

    const remove = () => emit('remove');

    const chips = computed(() => {
      const result = (props.rule.filled ?? []).map((type, index, filled) => {
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

        return {
          key,
          component: type === ADVANCED_SEARCH_CHIP_TYPES.rangeValue ? 'advanced-search-range-chip' : 'advanced-search-chip',
          bind: {
            value: props.rule[type],
            active: key === props.activeKey,
            items: itemsByType.value[type],
            allowText: !itemsByType.value[type]?.length,
            closable: index === filled.length - 1 && !props.union,
            disabled: props.disabled,
            itemText: itemText ?? 'text',
            itemValue: itemValue ?? 'value',
            fetchItems,
            multiple,
            color: validator.errors.has(props.rule.key) ? 'error' : undefined,
          },
          on: {
            input: value => selectChipItem(value, type),
            click: () => clickChip(key),
            focusout: () => emit('focusout'),
            next: () => setTimeout(() => emit('next', getNextForFormItemType(
              props.rule,
              multiple
                ? getNextForFormItemType(props.rule, type)
                : type,
            ))), // TODO: refactor
            close: remove,
          },
        };
      });

      if (!isFinishedGroup.value && !props.rule.filled.includes(activeType.value)) {
        const type = activeType.value;
        const key = `${props.rule.key}.${type}`;

        let itemText;
        let itemValue;
        let fetchItems;

        if (type === ADVANCED_SEARCH_CHIP_TYPES.value) {
          itemText = currentAttribute.value?.itemText;
          itemValue = currentAttribute.value?.itemValue;
          fetchItems = currentAttribute.value?.fetchValues;
        }

        result.push({
          key,
          component: type === ADVANCED_SEARCH_CHIP_TYPES.rangeValue ? 'advanced-search-range-chip' : 'advanced-search-chip',
          bind: {
            first: props.first && !result.length,
            value: type === ADVANCED_SEARCH_CHIP_TYPES.rangeValue ? props.rule[type] : undefined,
            active: key === props.activeKey,
            disabled: props.disabled,
            alwaysActive: true,
            items: itemsByType.value[type],
            itemText: itemText ?? 'text',
            itemValue: itemValue ?? 'value',
            fetchItems,
            allowText: !itemsByType.value[type]?.length,
          },
          on: {
            input: selectItem,
            next: () => setTimeout(() => emit('next', getNextForFormItemType(props.rule, type))), // TODO: refactor
          },
        });
      }

      return result;
    });

    watch(() => props.union, (union) => {
      activeType.value = union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute;
    });

    onMounted(() => {
      validator.attach({
        name: props.rule.key,
        rules: 'required:true',
        getter: () => isFinishedGroup.value,
      });
    });

    onBeforeUnmount(() => {
      validator.detach(props.rule.key);
    });

    return {
      activeType,
      chips,
      remove,
      selectChipItem,
      selectItem,
      clickChip,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-new-advanced-search__group {
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
    }
  }
}
</style>
