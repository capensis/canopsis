<template>
  <div
    :class="[themeClasses]"
    class="c-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
  >
    <div class="v-input__control">
      <div class="v-input__slot">
        <div class="v-text-field__slot">
          <v-layout
            ref="layoutElement"
            class="c-advanced-search__groups-wrapper gap-1"
            align-center
            wrap
            @click="clickLayout"
          >
            <advanced-search-rule
              v-for="(rule, index) in rules"
              v-model="rules[index]"
              :key="rule.key"
              :attributes="items"
              :active="rule.key === activeKey"
              :union="index % 2 === 1"
              :first="index === 0"
              :allow-or="allowOr"
              @input="update($event, index)"
              @make:active="makeActive"
              @reset:active="resetActive"
              @next="nextRule"
              @remove="remove(index)"
            />
          </v-layout>
        </div>
        <div class="v-input__append-inner">
          <v-menu bottom>
            <template #activator="{ on }">
              <c-action-btn
                :tooltip="$t('common.search')"
                icon="history"
                v-on="on"
              />
            </template>
          </v-menu>
        </div>
      </div>
    </div>
    <v-layout>
      <c-action-btn
        :tooltip="$t('common.search')"
        icon="search"
        @click="submit"
      />
      <c-action-btn
        :tooltip="$t('common.clearSearch')"
        icon="clear"
        @click="clear"
      />
    </v-layout>
  </div>
</template>

<script>
import { computed, ref, set, provide, nextTick, onBeforeMount } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { ADVANCED_SEARCH_CHIP_TYPES, ADVANCED_SEARCH_UNION_CONDITIONS } from '@/constants';

import {
  advancedSearchToForm,
  advancedSearchRuleItemToFormItem,
  formToAdvancedSearch,
  isAlarmPatternField,
  isEntityPatternField,
  isPbehaviorPatternField,
} from '@/helpers/search/new-advanced-search';

import { useComponentInstance } from '@/hooks/vue';

import { useAdvancedSearchAttributes } from './hooks/new-advanced-search';
import AdvancedSearchRule from './partials/new/advanced-search-rule.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { AdvancedSearchRule },
  mixins: [Themeable],
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const instance = useComponentInstance();

    const layoutElement = ref(null);
    const rules = ref([advancedSearchRuleItemToFormItem()]);
    const activeKey = ref();

    const hasOr = computed(() => rules.value.some(({ union }) => union === ADVANCED_SEARCH_UNION_CONDITIONS.or));
    const hasAlarmField = computed(() => rules.value.some(({ attribute }) => isAlarmPatternField(attribute)));
    const hasEntityField = computed(() => rules.value.some(({ attribute }) => isEntityPatternField(attribute)));
    const hasPbehaviorField = computed(() => rules.value.some(({ attribute }) => isPbehaviorPatternField(attribute)));

    const allowOr = computed(() => [
      hasEntityField.value,
      hasPbehaviorField.value,
      hasAlarmField.value,
    ].filter(Boolean).length <= 1);

    const allowAlarmFields = computed(() => !hasOr.value || (!hasEntityField.value && !hasPbehaviorField.value));
    const allowEntityFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasPbehaviorField.value));
    const allowPbehaviorFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasEntityField.value));

    const { attributes: items } = useAdvancedSearchAttributes({
      allowAlarmFields,
      allowEntityFields,
      allowPbehaviorFields,
    });

    const addNewRule = () => {
      const newRule = advancedSearchRuleItemToFormItem();

      rules.value.push(advancedSearchRuleItemToFormItem());

      return newRule;
    };

    const makeActive = (key) => {
      activeKey.value = key;
    };

    const resetActive = (key) => {
      if (key === activeKey.value) {
        activeKey.value = null;
      }
    };

    let lastInputFocus = () => {};

    const registerLastInputFocus = focus => lastInputFocus = focus;

    provide('$registerLastInputFocus', registerLastInputFocus);

    const clickLayout = event => event.target === layoutElement.value && lastInputFocus();

    const nextRule = (withoutActive) => {
      const newRule = addNewRule();

      if (!withoutActive) {
        makeActive(newRule.key);
        nextTick(() => lastInputFocus());
      }
    };

    /**
     * Updates a rule at the specified index with the provided value and clears any existing errors.
     *
     * @param {*} value - The new value to set for the rule at the specified index. This can be any
     *                    data type that represents a valid rule.
     * @param {number} index - The index of the rule to be updated. This should be a valid index within
     *                         the `rules` array.
     */
    const update = (value, index) => {
      instance.errors.clear();
      set(rules.value, index, value);
    };

    /**
     * Removes a rule from the rules array at the specified index, or resets it if it's the last rule.
     *
     * @param {number} index - The index of the rule to be removed or reset. This should be a valid
     *                         index within the `rules` array.
     */
    const remove = (index) => {
      if (index === rules.value.length - 1) {
        set(rules.value, index, advancedSearchRuleItemToFormItem());

        return;
      }

      rules.value.splice(index, index === rules.value.length - 2 ? 1 : 2);
    };

    /**
     * Clears the current search field errors and resets the rules to their initial state.
     */
    const clear = () => {
      instance.errors.clear();
      rules.value = [advancedSearchRuleItemToFormItem()];
    };

    /**
     * Validates the form and emits a 'submit' event with the advanced search criteria if valid.
     */
    const submit = async () => {
      const isValid = await instance.$validator.validateAll();

      if (isValid) {
        console.log('SUBMIT');
        emit('submit', formToAdvancedSearch(rules.value));
      }
    };

    const extendValidatorRule = () => instance.$validator.extend('advancedSearchRule', ({ rule, finished }) => {
      if (rule.attribute && !finished) {
        return false;
      }

      if (!rule.attribute && rule.union) {
        const lastRule = rules.value.at(-1);
        const preLastRule = rules.value.at(-2);

        return !(!lastRule.attribute && preLastRule.key === rule.key);
      }

      return true;
    });

    onBeforeMount(extendValidatorRule);

    return {
      layoutElement,
      allowOr,
      rules,
      activeKey,

      items,

      update,
      makeActive,
      resetActive,
      nextRule,
      remove,
      submit,
      clear,
      clickLayout,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-search { // TODO: remove new
  --v-chip-gap: 4px;
  --input-min-inline-size: 20ch;

  &__groups-wrapper > * {
    flex: 0 1 auto;
  }

  &::v-deep {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }

    .layout {
      padding: var(--v-chip-gap) 0;
      gap: var(--v-chip-gap);
    }

    .v-chip {
      padding: 0 8px;
      margin: 0;

      &:has(> .v-chip__content > .v-chip) {
        padding: 0 6px !important;

        button {
          margin: 0 -2px 0 0 !important;
        }
      }

      &__content {
        gap: var(--v-chip-gap);
      }

      .v-chip {
        height: 24px;

        &.theme--light {
          background: var(--v-application-background-base, #FFFFFF);
        }

        &.theme--dark {
          background: var(--v-application-background-base, #121212);
        }
      }
    }

    button {
      margin-left: 4px !important;
    }
  }
}
</style>
