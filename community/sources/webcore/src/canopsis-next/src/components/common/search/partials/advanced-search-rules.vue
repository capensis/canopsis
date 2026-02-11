<template>
  <v-layout
    ref="layoutElement"
    :class="{ 'c-advanced-search__groups-wrapper--disabled': disabled }"
    class="c-advanced-search__groups-wrapper gap-1"
    align-center
    wrap
    @mouseup="mouseupLayout"
  >
    <advanced-search-rule
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
      @next="next(index)"
      @remove="remove(index)"
    />
  </v-layout>
</template>

<script>
import { ref, nextTick } from 'vue';

import { REGISTER_LAST_INPUT_FOCUS_KEY } from '@/constants';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/advanced-search';
import { isArrayOperator } from '@/helpers/entities/pattern/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import { useLastInputFocus } from '../hooks/focus';
import { useActiveKey } from '../hooks/active-key';

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
      required: true,
      default: () => [],
    },
    attributes: {
      type: Array,
      required: true,
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
    allowFocusOnMount: {
      type: Boolean,
      default: true,
    },
  },
  setup(props, { emit }) {
    const {
      updateModel,
      addItemIntoArray,
      updateItemInArray,
      removeItemsFromArray,
    } = useArrayModelField(props, emit);

    const layoutElement = ref(null);
    const { activeKey, makeActive, resetActive } = useActiveKey();
    const { focusRegister } = useLastInputFocus(REGISTER_LAST_INPUT_FOCUS_KEY);

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
      /**
       * If we are removing pre last item we should remove only last and replace the prelast by empty rule
       *
       * @example
       * We are removing the first item:
       * [`_id` `equal` `someId`], [<empty rule>] => [<empty rule>]
       */
      if (index === props.rules.length - 2) {
        const newRules = props.rules.reduce((acc, rule, ruleIndex) => {
          if (ruleIndex === index) {
            acc.push(advancedSearchRuleItemToFormItem());
          } else if (ruleIndex !== index + 1) {
            acc.push(rule);
          }

          return acc;
        }, []);

        updateModel(newRules);

        return;
      }

      if (props.rules.length === 1 || index === props.rules.length - 1) {
        updateItemInArray(index, advancedSearchRuleItemToFormItem());

        return;
      }

      /**
       * If we are removing the item not from the end and we have more then 2 items
       *
       * @example
       * We are removing the first item:
       * [`_id` `equal` `someId1`], [`AND`], [`_id` `equal` `someId2`], [`OR`], [<empty rule>] =>
       * [`_id` `equal` `someId2`], [`OR`], [<empty rule>]
       */
      removeItemsFromArray([index, index % 2 === 0 ? index + 1 : index - 1]);
    };

    /**
     * Advances to the next rule in the rules array. If the current index is the last one, it adds a new rule.
     *
     * @param {number} index - The current index in the rules array.
     */
    const next = (index) => {
      if (index !== props.rules.length - 1) {
        return;
      }

      add();

      nextTick(() => focusRegister.call());
    };

    /**
     * Handles mouseup events on the layout element and focuses on the last input if the target is the layout element.
     *
     * @param {Event} event - The mouseup event.
     */
    const mouseupLayout = event => event.target === layoutElement.value && focusRegister.call();

    /**
     * Determines whether a rule should receive focus when mounted based on its position
     * and the previous rule's condition.
     *
     * @param {number} index - The index of the rule in the rules array
     * @returns {boolean} - Returns true if the rule should receive focus, false otherwise
     */
    const getFocusOnMount = (index) => {
      if (!index && !props.rules[index]?.attribute) {
        return false;
      }

      return props.allowFocusOnMount && !isArrayOperator(props.rules[index - 1]?.operator);
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
.c-advanced-search__groups-wrapper {
  > * {
    flex: 0 1 auto;
  }

  &--disabled::v-deep .v-chip {
    cursor: pointer !important;
  }
}
</style>
