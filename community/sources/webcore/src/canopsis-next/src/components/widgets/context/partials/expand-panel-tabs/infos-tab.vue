<template>
  <v-layout column>
    <v-layout wrap>
      <v-flex xs4>
        <c-search-field v-model="searchingText" />
      </v-flex>
    </v-layout>
    <v-data-table
      :items="items"
      :headers="headers"
      :search="searchingText"
      item-key="name"
    />
  </v-layout>
</template>

<script>
import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

const INFOS_COLUMN_PREFIX = 'entity.infos';

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
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.description'),
          value: 'description',
        },
        {
          text: this.$t('common.value'),
          value: 'value',
        },
      ];
    },

    items() {
      return Object.entries(this.infos).map(([infoKey, info = {}]) => {
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
