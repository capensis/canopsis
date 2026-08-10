<template>
  <v-menu content-class="c-advanced-search__history-menu" location="bottom">
    <template #activator="{ props }">
      <c-action-btn
        :tooltip="$t('common.search')"
        :disabled="!searches.length"
        icon="history"
        v-bind="props"
      />
    </template>
    <v-list ref="listElement">
      <v-list-item v-for="search in preparedSearches" :key="search._id" @click="select(search)">
        <template #title>
          <advanced-search-rules :rules="search.rules" :attributes="attributes" disabled />
        </template>
        <template #append>
          <advanced-search-history-item-btns
            :id="search._id"
            :pinned="search.pinned"
            @remove="remove"
            @toggle-pin="togglePin"
          />
        </template>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { computed, ref } from 'vue';

import { advancedSearchToForm } from '@/helpers/search/advanced-search';

import AdvancedSearchHistoryItemBtns from './advanced-search-history-item-btns.vue';
import AdvancedSearchRules from './advanced-search-rules.vue';

export default {
  components: { AdvancedSearchHistoryItemBtns, AdvancedSearchRules },

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

  emits: ['select', 'remove', 'toggle-pin'],

  setup(props, { emit }) {
    const listElement = ref(null);
    const preparedSearches = computed(() =>
      props.searches.map(search => ({
        _id: search._id,
        pinned: search.pinned,
        rules: search.rules ?? advancedSearchToForm(search),
      })),
    );

    const select = search => emit('select', search);
    const remove = id => emit('remove', id);
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
.c-advanced-search__history-menu {
  max-height: 95vh;
}
</style>
