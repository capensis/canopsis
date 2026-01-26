<template>
  <c-new-advanced-search-field
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :allow-or="allowOr"
    @submit="submit"
    @reset="reset"
    @remove:search="removeSearch"
    @toggle-pin:search="togglePinForSearch"
  />
</template>

<script>
import { ref } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/alarm-advanced-search';

import { useAlarmAdvancedSearchAttributes } from './hooks/alarm-advanced-search';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [Themeable],
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    const { attributes, allowOr } = useAlarmAdvancedSearchAttributes({ rules });

    /**
     * Submits the current advanced search configuration.
     *
     * @returns {Promise<void>} A promise that resolves when the submission process is complete.
     */
    const submit = newSearch => emit('submit', newSearch);

    /**
     * Emits an event to reset the advanced search configuration.
     */
    const reset = () => emit('reset');

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

    return {
      rules,
      attributes,
      allowOr,

      reset,
      submit,
      removeSearch,
      togglePinForSearch,
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
