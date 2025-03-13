<template>
  <v-layout class="c-alarm-advanced-search__wrapper" align-end>
    <div
      :class="[themeClasses]"
      class="c-alarm-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
    >
      <div class="v-input__control">
        <div class="v-input__slot pb-1">
          <div class="v-text-field__slot">
            <advanced-search-rules
              v-model="rules"
              :attributes="attributes"
              :allow-or="allowOr"
              @input="inputSearch"
            />
          </div>
        </div>
      </div>
    </div>
    <v-layout>
      <advanced-search-history-btn
        :searches="searches"
        :attributes="attributes"
        @select="select"
        @remove="removeSearch"
        @toggle-pin="togglePinForSearch"
      />
      <c-action-btn
        :tooltip="$t('common.search')"
        icon="search"
        @click="submit"
      />
      <c-action-btn
        :tooltip="$t('common.clearSearch')"
        icon="clear"
        @click="reset"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { ADVANCED_SEARCH_UNION_CONDITIONS } from '@/constants';

import { uuid } from '@/helpers/uuid';
import {
  advancedSearchRuleItemToFormItem,
  formToAdvancedSearch,
  isAlarmPatternField,
  isEntityPatternField,
  isPbehaviorPatternField,
} from '@/helpers/search/alarm-advanced-search';

import { useComponentInstance } from '@/hooks/vue';

import { useAdvancedSearchAttributes, useAdvancedSearchValidator } from './hooks/alarm-advanced-search';
import AdvancedSearchRules from './partials/alarm-advanced-search-rules.vue';
import AdvancedSearchHistoryBtn from './partials/alarm-advanced-search-history-btn.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { AdvancedSearchHistoryBtn, AdvancedSearchRules },
  mixins: [Themeable],
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const instance = useComponentInstance();
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    let activeSearch = null;

    /**
     * HAS FLAGS
     */
    const hasOr = computed(() => rules.value.some(({ union }) => union === ADVANCED_SEARCH_UNION_CONDITIONS.or));
    const hasAlarmField = computed(() => rules.value.some(({ attribute }) => isAlarmPatternField(attribute)));
    const hasEntityField = computed(() => rules.value.some(({ attribute }) => isEntityPatternField(attribute)));
    const hasPbehaviorField = computed(() => rules.value.some(({ attribute }) => isPbehaviorPatternField(attribute)));

    /**
     * ALLOW FLAGS
     */
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
     * Submits the current advanced search configuration.
     *
     * @returns {Promise<void>} A promise that resolves when the submission process is complete.
     */
    const submit = async () => {
      const isValid = await instance.$validator.validateAll();

      if (isValid) {
        const newSearch = formToAdvancedSearch(rules.value);

        newSearch._id = activeSearch?._id ?? uuid();
        newSearch.pinned = activeSearch?.pinned ?? false;

        emit('submit', newSearch);
      }
    };

    /**
     * Resets the active search configuration to null.
     */
    const resetActiveSearch = () => activeSearch = null;

    /**
     * Input search handler
     */
    const inputSearch = (newRules) => {
      resetActiveSearch();

      if (newRules.length === 1 && newRules[0].text) {
        submit();
      }
    };

    /**
     * Reset the current search field errors and resets the rules to their initial state.
     */
    const reset = () => {
      instance.errors.clear();
      rules.value = [advancedSearchRuleItemToFormItem()];

      resetActiveSearch();

      emit('reset');
    };

    /**
     * Selects a search configuration and updates the active search variable.
     *
     * @param {Object} search - The search configuration object to be selected.
     * @param {Array} search.rules - An array of rules associated with the search configuration.
     */
    const select = (search) => {
      activeSearch = search;
      rules.value = [...search.rules, advancedSearchRuleItemToFormItem()];

      submit();
    };

    /**
     * Emits an event to remove a search configuration by its identifier.
     *
     * @param {string} id - The unique identifier of the search configuration to be removed.
     */
    const removeSearch = id => emit('remove:search', id);

    /**
     * Emits an event to toggle the pinned status of a search configuration.
     *
     * @param {string} id - The unique identifier of the search configuration to toggle the pin status.
     */
    const togglePinForSearch = id => emit('toggle-pin:search', id);

    useAdvancedSearchValidator({ rules });

    return {
      rules,
      attributes,
      allowOr,

      select,
      reset,
      submit,
      removeSearch,
      togglePinForSearch,
      inputSearch,
    };
  },
};
</script>

<style lang="scss">
:root {
  --alarm-advanced-search-chip-gap: 4px;
}

.c-alarm-advanced-search {
  --input-min-inline-size: 20ch;

  padding-bottom: 2px;

  & &__rule {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }
  }

  &__groups-wrapper {
    &, .layout {
      gap: var(--alarm-advanced-search-chip-gap);
    }
  }

  .v-input__append-inner {
    margin: 0;
    align-self: center;
  }

  button {
    margin-left: 4px !important;
  }

  &__chip.v-chip {
    padding: 0 8px;
    margin: 0;

    .v-chip__content {
      gap: var(--alarm-advanced-search-chip-gap);
    }

    &:has(> .v-chip__content > .v-chip) {
      padding: 0 6px !important;

      button {
        margin: 0 -2px 0 0 !important;
      }
    }

    .v-chip {
      height: 24px;
      margin: 0;

      &.theme--light {
        background: var(--v-application-background-base, #FFFFFF);
      }

      &.theme--dark {
        background: var(--v-application-background-base, #121212);
      }
    }
  }
}
</style>
