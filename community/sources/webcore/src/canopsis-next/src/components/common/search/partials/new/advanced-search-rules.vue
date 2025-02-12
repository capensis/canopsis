<template>
  <v-layout
    ref="layoutElement"
    class="c-advanced-search__groups-wrapper gap-1"
    align-center
    wrap
    @click="clickLayout"
  >
    <advanced-search-rule
      v-for="(rule, index) in rules"
      v-field="rules[index]"
      :key="rule.key"
      :attributes="attributes"
      :active="rule.key === activeKey"
      :union="index % 2 === 1"
      :first="index === 0"
      :allow-or="allowOr"
      :disabled="disabled"
      @input="update($event, index)"
      @make:active="makeActive"
      @reset:active="resetActive"
      @next="next"
      @remove="remove(index)"
    />
  </v-layout>
</template>

<script>
import { ref, provide, nextTick } from 'vue';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/new-advanced-search';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import AdvancedSearchRule from './advanced-search-rule.vue';

export default {
  components: { AdvancedSearchRule },
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
     * Removes an item from the array at the specified index.
     *
     * @param {number} index - The index of the item to remove.
     */
    const remove = index => removeItemFromArray(index);

    /**
     * Adds a new item and optionally focuses on the last input.
     *
     * @param {boolean} withoutActive - If true, does not focus on the last input.
     */
    const next = (withoutActive) => {
      add();

      if (!withoutActive) {
        nextTick(() => lastInputFocus());
      }
    };

    /**
     * Handles click events on the layout element and focuses on the last input if the target is the layout element.
     *
     * @param {Event} event - The click event.
     */
    const clickLayout = event => event.target === layoutElement.value && lastInputFocus();

    return {
      layoutElement,
      activeKey,

      makeActive,
      resetActive,
      next,
      add,
      update,
      remove,
      clickLayout,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-search__groups-wrapper > * {
  flex: 0 1 auto;
}
</style>
