<template>
  <v-menu bottom>
    <template #activator="{ on }">
      <c-action-btn
        :tooltip="$t('common.search')"
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
        <v-list-item-content>
          <advanced-search-rules
            :rules="search.rules"
            :attributes="attributes"
            disabled
          />
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
<script>
import { computed, ref } from 'vue';

import { advancedSearchToForm } from '@/helpers/search/new-advanced-search';

import AdvancedSearchRules from './advanced-search-rules.vue';

export default {
  components: { AdvancedSearchRules },
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

    return {
      listElement,

      preparedSearches,

      select,
    };
  },
};
</script>
