<template>
  <div
    :class="[themeClasses]"
    class="c-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
  >
    <div class="v-input__control">
      <div class="v-input__slot">
        <div class="v-text-field__slot">
          <advanced-search-rules
            v-model="rules"
            :attributes="attributes"
            :allow-or="allowOr"
          />
        </div>
        <div class="v-input__append-inner">
          <advanced-search-history-btn
            :searches="searches"
            :attributes="attributes"
            @select="select"
          />
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
import { computed, ref, onBeforeMount } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { ADVANCED_SEARCH_UNION_CONDITIONS } from '@/constants';

import {
  advancedSearchRuleItemToFormItem,
  formToAdvancedSearch,
  isAlarmPatternField,
  isEntityPatternField,
  isPbehaviorPatternField,
} from '@/helpers/search/new-advanced-search';

import { useComponentInstance } from '@/hooks/vue';

import { useAdvancedSearchAttributes } from './hooks/new-advanced-search';
import AdvancedSearchRules from './partials/new/advanced-search-rules.vue';
import AdvancedSearchHistoryBtn from './partials/new/advanced-search-history-btn.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { AdvancedSearchHistoryBtn, AdvancedSearchRules },
  mixins: [Themeable],
  props: {
    /* searches: {
      type: Array,
      default: () => [],
    }, */
  },
  setup(props, { emit }) {
    const searches = JSON.parse('[{"_id":"asd","positions":["alarm_pattern","alarm_pattern","pbehavior_pattern","pbehavior_pattern"],"alarm_pattern":[[{"field":"v.ack.a","cond":{"value":"c78b6ba1-39c5-44a0-86df-50b0dab36845","type":"eq"}},{"field":"v.display_name","cond":{"value":"выфвфывфы","type":"eq"}}]],"entity_pattern":[],"pbehavior_pattern":[[{"field":"v.pbehavior_info.type","cond":{"value":"c6520636-8489-4f21-91ca-2a55317931ba","type":"eq"}},{"field":"v.pbehavior_info.name","cond":{"value":"0618be38-d661-4e30-a31b-cd195af03f6f","type":"eq"}}]]}]');

    const instance = useComponentInstance();
    const rules = ref([advancedSearchRuleItemToFormItem()]);

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

    const { attributes } = useAdvancedSearchAttributes({
      allowAlarmFields,
      allowEntityFields,
      allowPbehaviorFields,
    });

    /**
     * Clears the current search field errors and resets the rules to their initial state.
     */
    const clear = () => {
      instance.errors.clear();
      rules.value = [advancedSearchRuleItemToFormItem()];
    };

    const select = (search) => {
      rules.value = [...search.rules, advancedSearchRuleItemToFormItem()];
    };

    /**
     * Validates the form and emits a 'submit' event with the advanced search criteria if valid.
     */
    const submit = async () => {
      const isValid = await instance.$validator.validateAll();

      if (isValid) {
        // console.log(JSON.stringify(formToAdvancedSearch(rules.value)));
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

        return !(!lastRule?.attribute && preLastRule?.key === rule?.key);
      }

      return true;
    });

    onBeforeMount(extendValidatorRule);

    return {
      rules,
      attributes,
      allowOr,

      searches,

      select,
      clear,
      submit,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-search { // TODO: remove new
  --v-chip-gap: 4px;
  --input-min-inline-size: 20ch;

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
