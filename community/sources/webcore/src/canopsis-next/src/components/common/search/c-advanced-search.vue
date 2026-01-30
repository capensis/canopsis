<template>
  <v-layout class="c-advanced-search__wrapper" align-end>
    <div
      :class="[themeClasses]"
      class="c-advanced-search mt-0 pt-2 v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
    >
      <div class="v-input__control">
        <div class="v-input__slot pb-1">
          <div class="v-text-field__slot">
            <advanced-search-rules
              v-field="rules"
              :attributes="attributes"
              :allow-or="allowOr"
              :allow-focus-on-mount="allowFocusOnMount"
              @input="inputSearch"
            />
          </div>
        </div>
      </div>
    </div>
    <v-layout>
      <advanced-search-history-btn
        v-if="withHistory"
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
import { ref, toRef } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { uuid } from '@/helpers/uuid';
import { advancedSearchRuleItemToFormItem, formToAdvancedSearch } from '@/helpers/search/advanced-search';

import { useValidator } from '@/hooks/validator/validator';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import { useAdvancedSearchValidator } from './hooks/advanced-search';
import AdvancedSearchRules from './partials/advanced-search-rules.vue';
import AdvancedSearchHistoryBtn from './partials/advanced-search-history-btn.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { AdvancedSearchHistoryBtn, AdvancedSearchRules },
  mixins: [Themeable],
  model: {
    prop: 'rules',
    event: 'input',
  },
  props: {
    rules: {
      type: Array,
      default: () => [advancedSearchRuleItemToFormItem()],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    allowOr: {
      type: Boolean,
      default: true,
    },
    searches: {
      type: Array,
      default: () => [],
    },
    withHistory: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const { updateModel } = useArrayModelField(props, emit);
    const allowFocusOnMount = ref(true);

    let activeSearch = null;

    /**
     * Temporarily disables the allow focus on mount flag and re-enables it after a delay.
     * This is used to prevent unwanted focus behavior during certain operations.
     */
    const runToggleAllowFocusOnMount = () => {
      allowFocusOnMount.value = false;

      setTimeout(() => allowFocusOnMount.value = true, 1000);
    };

    /**
     * Submits the current advanced search configuration.
     *
     * @returns {Promise<void>} A promise that resolves when the submission process is complete.
     */
    const submit = async () => {
      const isValid = await validator.validateAll();

      if (isValid) {
        const newSearch = formToAdvancedSearch(props.rules);

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
      updateModel(newRules);

      if (newRules.length === 1 && newRules[0].text) {
        submit();
      }
    };

    /**
     * Reset the current search field errors and resets the rules to their initial state.
     */
    const reset = () => {
      validator.errors.clear();
      updateModel([advancedSearchRuleItemToFormItem()]);

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
      updateModel([...search.rules, advancedSearchRuleItemToFormItem()]);

      runToggleAllowFocusOnMount();

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

    useAdvancedSearchValidator({ rules: toRef(props, 'rules') });

    return {
      allowFocusOnMount,

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
  --advanced-search-chip-gap: 4px;
}

.c-advanced-search {
  --input-min-inline-size: 30ch;

  & &__rule {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }
  }

  &__groups-wrapper {
    &, .layout {
      gap: var(--advanced-search-chip-gap);
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
      gap: var(--advanced-search-chip-gap);
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
