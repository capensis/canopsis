<template>
  <v-layout
    ref="layoutElement"
    :class="{ 'c-alarm-advanced-search__groups-wrapper--disabled': disabled }"
    class="c-alarm-advanced-search__groups-wrapper gap-1"
    align-center
    wrap
    @mouseup="mouseupLayout"
  >
    <alarm-advanced-search-rule
      v-for="(rule, index) in rules"
      :key="rule.key"
      :rule="rule"
      :attributes="attributes"
      :active="rule.key === activeKey"
      :union="index % 2 === 1"
      :first="index === 0"
      :allow-or="allowOr"
      :disabled="disabled"
      :focus-on-mount="getFocusOnMount(index)"
      @input="update($event, index)"
      @make:active="makeActive"
      @reset:active="resetActive"
      @next="next($event, index)"
      @remove="remove(index)"
    />
  </v-layout>
</template>

<script>
import { ref, provide, nextTick } from 'vue';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/alarm-advanced-search';
import { isArrayCondition } from '@/helpers/entities/pattern/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import AlarmAdvancedSearchRule from './alarm-advanced-search-rule.vue';

export default {
  components: { AlarmAdvancedSearchRule },
  model: {
    prop: 'rules',
    event: 'input',
  },
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    allowOr: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const {
      addItemIntoArray,
      updateItemInArray,
      removeItemFromArray,
    } = useArrayModelField(props, emit);

    const layoutElement = ref(null);
    const activeKey = ref();

    /**
     * A no-operation function to store the last input focus.
     *
     * @type {Function}
     */
    let lastInputFocus = () => {};

    /**
     * Provides a method to register the last input focus function.
     *
     * @param {Function} focus - The function to set as the last input focus.
     */
    provide('$registerLastInputFocus', focus => lastInputFocus = focus);

    /**
     * Sets the active key to the specified key.
     *
     * @param {string} key - The key to set as active.
     */
    const makeActive = key => activeKey.value = key;

    /**
     * Resets the active key if it matches the specified key.
     *
     * @param {string} key - The key to check against the active key.
     */
    const resetActive = (key) => {
      if (key === activeKey.value) {
        activeKey.value = null;
      }
    };

    /**
     * Adds a new item into the array using a form item.
     */
    const add = () => addItemIntoArray(advancedSearchRuleItemToFormItem());

    /**
     * Updates an item in the array at the specified index with the given value.
     *
     * @param {*} value - The new value to update the item with.
     * @param {number} index - The index of the item to update.
     */
    const update = (value, index) => updateItemInArray(index, value);

    /**
     * Removes a rule from the rules array at the specified index.
     * If there is only one rule in the array, it replaces by new rule.
     *
     * @param {number} index - The index of the rule to be removed or updated.
     */
    const remove = (index) => {
      if (props.rules.length === 1) {
        updateItemInArray(index, advancedSearchRuleItemToFormItem());

        return;
      }

      removeItemFromArray(index);
    };

    /**
     * Advances to the next rule in the rules array. If the current index is the last one,
     * it adds a new rule. Optionally focuses on the last input if `withoutActive` is false.
     *
     * @param {boolean} withoutActive - Determines whether to focus on the last input after adding a new rule.
     * @param {number} index - The current index in the rules array.
     */
    const next = (withoutActive, index) => {
      if (index !== props.rules.length - 1) {
        return;
      }

      add();

      nextTick(() => lastInputFocus());
    };

    /**
     * Handles mouseup events on the layout element and focuses on the last input if the target is the layout element.
     *
     * @param {Event} event - The mouseup event.
     */
    const mouseupLayout = event => event.target === layoutElement.value && lastInputFocus();

    const getFocusOnMount = (index) => {
      if (!index && !props.rules[index]?.attribute) {
        return false;
      }

      return !isArrayCondition(props.rules[index - 1]?.operator);
    };

    return {
      layoutElement,
      activeKey,

      makeActive,
      resetActive,
      next,
      add,
      update,
      remove,
      mouseupLayout,
      getFocusOnMount,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-alarm-advanced-search__groups-wrapper {
  > * {
    flex: 0 1 auto;
  }

  &--disabled::v-deep .v-chip {
    cursor: pointer !important;
  }
}
</style>
