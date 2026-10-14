<template>
  <c-layout column>
    <c-layout wrap>
      <v-col cols="4" class="pa-0">
        <c-search-field v-model="searchingText" />
      </v-col>
    </c-layout>
    <v-data-table
      :items="items"
      :headers="headers"
      :search="searchingText"
      item-value="name"
    />
  </c-layout>
</template>

<script>
import { INFOS_COLUMN_PREFIX } from '@/constants';

import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

export default {
  mixins: [widgetColumnsFiltersMixin],

  props: {
    infos: {
      type: Object,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },

  data() {
    return {
      searchingText: '',
    };
  },

  computed: {
    headers() {
      return [
        {
          title: this.$t('common.name'),
          key: 'name',
          value: 'name',
        },
        {
          title: this.$t('common.description'),
          key: 'description',
          value: 'description',
        },
        {
          title: this.$t('common.value'),
          key: 'value',
          value: 'value',
        },
      ];
    },

    items() {
      return Object.entries(this.infos ?? {}).map(([infoKey, info = {}]) => {
        const valueColumnFilter = this.columnsFiltersMap[`${INFOS_COLUMN_PREFIX}.${infoKey}.value`] ?? String;

        return {
          name: infoKey,
          description: info.description,
          value: valueColumnFilter(info.value),
        };
      });
    },
  },
};
</script>
