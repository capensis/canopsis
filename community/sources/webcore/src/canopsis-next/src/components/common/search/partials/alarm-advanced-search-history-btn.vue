<template>
  <v-menu
    :nudge-bottom="1"
    :transition="false"
    content-class="c-alarm-advanced-search__history-menu"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <c-action-btn
        :tooltip="$t('common.search')"
        :disabled="!searches.length"
        icon="history"
        v-on="on"
      />
    </template>
    <v-list ref="listElement">
      <v-list-item
        v-for="search in preparedSearches"
        :key="search._id"
        @click="select(search)"
      >
        <v-list-item-content class="pa-0">
          <alarm-advanced-search-rules
            :rules="search.rules"
            :attributes="attributes"
            disabled
          />
        </v-list-item-content>
        <v-list-item-action>
          <advanced-search-history-item-btns
            :id="search._id"
            :pinned="search.pinned"
            @remove="remove"
            @toggle-pin="togglePin"
          />
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
<script>
import { computed, ref } from 'vue';

import { advancedSearchToForm } from '@/helpers/search/alarm-advanced-search';

import AdvancedSearchHistoryItemBtns from './advanced-search-history-item-btns.vue';
import AlarmAdvancedSearchRules from './alarm-advanced-search-rules.vue';

export default {
  components: { AdvancedSearchHistoryItemBtns, AlarmAdvancedSearchRules },
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const listElement = ref(null);
    const preparedSearches = computed(() => props.searches.map(search => ({
      _id: search._id,
      pinned: search.pinned,
      rules: advancedSearchToForm(search),
    })));

    /**
     * Emits a 'select' event with the specified search configuration.
     *
     * @param {Object} search - The search configuration object to be selected.
     */
    const select = search => emit('select', search);

    /**
     * Emits a 'remove' event with the specified identifier.
     *
     * @param {string} id - The unique identifier of the item to be removed.
     */
    const remove = id => emit('remove', id);

    /**
     * Emits a 'toggle-pin' event with the specified identifier.
     *
     * @param {string} id - The unique identifier of the item whose pin status is to be toggled.
     */
    const togglePin = id => emit('toggle-pin', id);

    return {
      listElement,

      preparedSearches,

      select,
      remove,
      togglePin,
    };
  },
};
</script>

<style lang="scss">
.v-list-item
.c-alarm-advanced-search__history-menu {
  max-height: 95vh;
}
</style>
