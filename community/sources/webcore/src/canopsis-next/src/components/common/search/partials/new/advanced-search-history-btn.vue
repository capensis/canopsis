<template>
  <v-menu
    :nudge-bottom="1"
    :transition="false"
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
        class="c-new-advanced-search__history__item"
        @click="select(search)"
      >
        <v-list-item-content>
          <advanced-search-rules
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

import { advancedSearchToForm } from '@/helpers/search/new-advanced-search';

import AdvancedSearchHistoryItemBtns from '../advanced-search-history-item-btns.vue';

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
  setup(props, { emit }) {
    const listElement = ref(null);
    const preparedSearches = computed(() => props.searches.map(search => ({
      _id: search._id,
      rules: advancedSearchToForm(search),
    })));

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
