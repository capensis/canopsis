<template>
  <v-layout
    :class="{ 'c-alarm-advanced-search__rule--union': union }"
    class="c-alarm-advanced-search__rule"
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
import { uniq, pick } from 'lodash';
import {
  computed,
  ref,
  watch,
  nextTick,
  toRef,
} from 'vue';

import { ALARM_ADVANCED_SEARCH_CHIP_TYPES, PATTERN_FIELD_TYPES, PATTERN_QUICK_RANGES } from '@/constants';

import { isArrayCondition } from '@/helpers/entities/pattern/form';
import {
  advancedSearchRuleItemToFormItem,
  getInitialFormItemType,
  getNextForFormItemType,
  isArrayItem,
  isCustomRangeItem,
  isDurationItem,
  isNumberValueType,
} from '@/helpers/search/alarm-advanced-search';

import { useModelField } from '@/hooks/form/model-field';

import { useAdvancedSearchRuleActiveItems, useAttachAdvancedSearchRuleValidator } from '../hooks/alarm-advanced-search';

import AlarmAdvancedSearchChip from './alarm-advanced-search-chip.vue';
import AlarmAdvancedSearchRangeChip from './alarm-advanced-search-range-chip.vue';
import AlarmAdvancedSearchDurationChip from './alarm-advanced-search-duration-chip.vue';

export default {
  components: {
    AlarmAdvancedSearchChip,
    AlarmAdvancedSearchRangeChip,
    AlarmAdvancedSearchDurationChip,
  },
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
    focusOnMount: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const inputType = ref(getInitialFormItemType(props.rule, props.union));
    const activeType = ref(null);

    const isFinishedRule = computed(() => !inputType.value);

    const { attributesMap, currentAttribute, itemsByType } = useAdvancedSearchRuleActiveItems({
      rule: toRef(props, 'rule'),
      attributes: toRef(props, 'attributes'),
      intervalRanges: toRef(props, 'intervalRanges'),
      inputTypes: toRef(props, 'inputTypes'),
      allowOr: toRef(props, 'allowOr'),
    });

    const { validator } = useAttachAdvancedSearchRuleValidator({
      isFinishedRule,
      rule: toRef(props, 'rule'),
      disabled: toRef(props, 'disabled'),
    });

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

      if (isFinishedRule.value && !props.rule.text) {
        emit('next');
      }
    };

    /**
     * Sets the input type for the search component.
     *
     * @param {string} type - The type of input to be set (e.g., 'attribute', 'operator').
     */
    const setInputType = type => inputType.value = type;

    /**
     * Determines if the given value is considered text based on the current attribute and attributes map.
     *
     * @param {string} value - The value to be checked.
     * @returns {boolean} - Returns true if the value is considered text, false otherwise.
     */
    const isText = value => !props.union && !currentAttribute.value && !attributesMap.value[value];

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
     * Updates the search rule's chip item based on the provided value and type.
     *
     * @param {*} value - The value to set for the specified type. This can vary depending on the type.
     * @param {string} type - The type of the chip item being updated. It should be one of the predefined
     *                        types in `ALARM_ADVANCED_SEARCH_CHIP_TYPES`.
     */
    const updateChipItem = (value, type) => {
      const oldFilled = props.rule.filled ?? [];
      const typeIndex = oldFilled.indexOf(type);
      const filled = typeIndex === -1 ? oldFilled : oldFilled.slice(0, typeIndex + 1);
      const filledForRemove = typeIndex === -1 ? [] : oldFilled.slice(typeIndex + 1);
      let skipType = false;

      if (isArrayItem(type, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.value);
        skipType = true;
      }

      if (isCustomRangeItem(type, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
        skipType = true;
      }

      if (isDurationItem(type, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration);
        skipType = true;
      }

      updateModel({
        ...props.rule,
        ...pick(advancedSearchRuleItemToFormItem(), filledForRemove),
        [type]: value,
        filled: uniq(filled),
      });

      if (
        [
          ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValue,
          ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration,
          ALARM_ADVANCED_SEARCH_CHIP_TYPES.value,
        ].includes(type)
      ) {
        return;
      }

      setInputType(type);

      if (!isFinishedRule.value) {
        nextTick(() => goToNextType(skipType));
      }
    };

    /**
     * Updates the current search rule item with the specified value and manages the progression of filled criteria.
     *
     * @param {*} value - The value to set for the current active type. This can vary depending on the type.
     */
    const updateItem = (value) => {
      if (isText(value)) {
        setInputType(ALARM_ADVANCED_SEARCH_CHIP_TYPES.text);
      }

      const filled = [...(props.rule.filled ?? []), inputType.value];
      const preparedRule = { ...props.rule };
      let skipType = false;

      if (isArrayItem(inputType.value, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.value);
        preparedRule[ALARM_ADVANCED_SEARCH_CHIP_TYPES.value] = [];
        skipType = true;
      }

      if (isCustomRangeItem(inputType.value, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValue);
        skipType = true;
      }

      if (isDurationItem(inputType.value, value)) {
        filled.push(ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration);
        skipType = true;
      }

      updateModel({
        ...preparedRule,
        filled: uniq(filled),
        [inputType.value]: value,
      });

      nextTick(() => goToNextType(skipType));
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

    /**
     * Generates the attributes and event handlers for a chip component based on the provided parameters.
     *
     * @param {Object} options - The options for configuring the chip attributes.
     * @param {boolean} options.input - Determines if the chip should always be active.
     * @param {boolean} options.first - Indicates if this is the first chip, affecting text allowance.
     * @param {boolean} options.closable - Specifies if the chip can be closed.
     * @param {string} [options.type = inputType.value] - The type of the chip, defaulting to the current input type.
     * @returns {Object} - An object containing the key, component type, binding attributes,
     *                     and event handlers for the chip.
     */
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

      if (type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.value) {
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
        number: isNumberValueType(props.rule, type),
        fetchItems,
        first,
        focusOnMount: props.focusOnMount,
      };

      let on = { input: updateItem };

      if (input) {
        bind.alwaysActive = true;
      } else {
        bind.alwaysActive = multiple && !props.rule[type]?.length;
        bind.active = isActiveType(type);
        bind.value = props.rule[type];
        bind.closable = closable;
        bind.multiple = multiple;
        bind.color = validator.errors.has(props.rule.key) ? 'error' : undefined;

        on = {
          input: value => updateChipItem(value, isText(value) ? ALARM_ADVANCED_SEARCH_CHIP_TYPES.text : type),
          click: () => clickChip(type),
          focusout: () => focusOutChip(type),
          close: remove,
        };
      }

      let component = 'alarm-advanced-search-chip';

      if (type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.rangeValue) {
        component = 'alarm-advanced-search-range-chip';
      } else if (type === ALARM_ADVANCED_SEARCH_CHIP_TYPES.duration) {
        component = 'alarm-advanced-search-duration-chip';
      }

      return {
        key,
        component,
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

    watch(() => props.union, union => (
      inputType.value = union ? ALARM_ADVANCED_SEARCH_CHIP_TYPES.union : ALARM_ADVANCED_SEARCH_CHIP_TYPES.attribute
    ));

    return {
      chips,
      remove,
      clickChip,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-alarm-advanced-search__rule {
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
