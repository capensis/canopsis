<template>
  <c-layout class="c-advanced-search__wrapper" align="end">
    <div
      :class="{ 'v-input--is-loading': pending }"
      class="c-advanced-search mt-0 pt-2 v-input v-input--hide-details v-text-field v-text-field--single-line v-select v-autocomplete"
    >
      <div class="v-input__control">
        <div class="v-input__slot pb-1">
          <div class="v-text-field__slot" @mouseup="focusSearchInput">
            <advanced-search-rules
              ref="rulesElement"
              :rules="rules"
              :attributes="attributes"
              :allow-or="allowOr"
              :allow-focus-on-mount="allowFocusOnMount"
              :disabled="disabled"
              @update:rules="inputSearch"
            />
          </div>
        </div>
        <v-progress-linear v-if="pending" height="2" color="primary" indeterminate />
      </div>
    </div>
    <c-layout>
      <advanced-search-history-btn
        v-if="withHistory"
        :searches="searches"
        :attributes="attributes"
        @select="select"
        @remove="removeSearch"
        @toggle-pin="togglePinForSearch"
      />
      <c-action-btn :tooltip="$t('common.search')" icon="search" @click="submit" />
      <c-action-btn :tooltip="$t('common.clearSearch')" icon="clear" @click="reset" />
    </c-layout>
  </c-layout>
</template>

<script>
import { ref, toRef, inject, onMounted, onBeforeUnmount } from 'vue';

import { uid } from '@/helpers/uid';
import {
  advancedSearchRuleItemToFormItem,
  formToAdvancedSearch,
} from '@/helpers/search/advanced-search';

import { useAdvancedSearchValidator } from './hooks/advanced-search';
import AdvancedSearchRules from './partials/advanced-search-rules.vue';
import AdvancedSearchHistoryBtn from './partials/advanced-search-history-btn.vue';

export default {
  components: { AdvancedSearchHistoryBtn, AdvancedSearchRules },

  props: {
    rules: {
      type: Array,
      default: () => [advancedSearchRuleItemToFormItem()],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
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
    alarmPattern: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'advancedSearch',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:rules', 'submit', 'reset', 'remove:search', 'toggle-pin:search'],

  setup(props, { emit }) {
    const { validateAll, errors } = useAdvancedSearchValidator({
      rules: toRef(props, 'rules'),
      name: toRef(props, 'name'),
      disabled: toRef(props, 'disabled'),
    });
    const allowFocusOnMount = ref(true);
    const rulesElement = ref(null);

    let activeSearch = null;

    /**
     * Temporarily disables the allow focus on mount flag and re-enables it after a delay.
     */
    const runToggleAllowFocusOnMount = () => {
      allowFocusOnMount.value = false;

      setTimeout(() => (allowFocusOnMount.value = true), 1000);
    };

    const updateModel = (rules) => {
      emit('update:rules', rules);
    };

    /**
     * Submits the current advanced search configuration.
     */
    const submit = async () => {
      const isValid = await validateAll();

      if (isValid) {
        const newSearch = formToAdvancedSearch(props.rules, props.alarmPattern);

        newSearch._id = activeSearch?._id ?? uid('advanced-search');
        newSearch.pinned = activeSearch?.pinned ?? false;
        newSearch.rules = props.rules;

        emit('submit', newSearch);
      }
    };

    const resetActiveSearch = () => (activeSearch = null);

    const inputSearch = (newRules) => {
      resetActiveSearch();
      updateModel(newRules);
    };

    const reset = () => {
      errors.clear();
      updateModel([advancedSearchRuleItemToFormItem()]);

      resetActiveSearch();

      emit('reset');
    };

    const select = (search) => {
      activeSearch = search;
      updateModel([...search.rules, advancedSearchRuleItemToFormItem()]);

      runToggleAllowFocusOnMount();

      submit();
    };

    const removeSearch = (id) => emit('remove:search', id);
    const togglePinForSearch = (id) => emit('toggle-pin:search', id);
    const focusSearchInput = (event) => {
      if (props.disabled || event.target.closest('button, input, .v-chip, .v-list-item')) {
        return;
      }

      rulesElement.value?.focusLastInput?.();
    };

    const registerSelectAdvancedSearch = inject('$registerSelectAdvancedSearch', () => {});

    onMounted(() => registerSelectAdvancedSearch(select));
    onBeforeUnmount(() => registerSelectAdvancedSearch(null));

    return {
      rulesElement,
      allowFocusOnMount,

      select,
      reset,
      submit,
      removeSearch,
      togglePinForSearch,
      inputSearch,
      focusSearchInput,
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

  flex: 1 1 auto;
  min-width: 280px;

  .v-input__control,
  .v-input__slot,
  .v-text-field__slot {
    width: 100%;
  }

  .v-input__control {
    display: block;
    grid-area: auto;
  }

  .v-input__slot {
    display: flex;
    align-items: flex-end;
    min-height: 40px;
    cursor: text;
    border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.42);
  }

  .v-text-field__slot {
    display: flex;
    align-items: center;
    min-width: 0;
  }

  & &__rule {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
      border: 0;
      outline: 0;
      background: transparent;
      color: currentColor;

      &::placeholder {
        color: rgba(var(--v-theme-on-surface), 0.6);
        opacity: 1;
      }
    }
  }

  &__wrapper {
    width: 100%;
  }

  &__groups-wrapper {
    &,
    .layout,
    .c-layout {
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
      background: rgb(var(--v-theme-background));
    }
  }
}
</style>
