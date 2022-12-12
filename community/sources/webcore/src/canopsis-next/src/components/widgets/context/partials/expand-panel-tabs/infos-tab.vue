<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-layout(row, wrap)
        v-flex(xs3)
          v-text-field(v-model="searchingText", :label="$t('context.infosSearchLabel')", dark)
      v-data-table(
        :items="items",
        :headers="headers",
        :search="searchingText",
        item-key="name"
      )
        template(#items="{ item }")
          td(v-for="column in headers", :key="column.value") {{ item | get(column.value) }}
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
