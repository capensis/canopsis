<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-layout(row, wrap)
        v-flex(xs3)
          v-text-field(v-model="searchingText", :label="$t('context.moreInfos.infosSearchLabel')", dark)
      v-data-table(
        :items="items",
        :headers="headers",
        :search="searchingText",
        item-key="item.name"
      )
        template(slot="items", slot-scope="props")
          td(
            v-for="column in headers",
            :key="column.value"
          ) {{ props.item | get(column.value) }}
</template>

<script>
import alarmColumnFiltersMixin from '@/mixins/entities/alarm-column-filters';

const INFOS_COLUMN_PREFIX = 'entity.infos';

export default {
  mixins: [alarmColumnFiltersMixin],
  props: {
    infos: {
      type: Object,
      required: true,
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
        const valueColumnFilter = this.columnFiltersMap[`${INFOS_COLUMN_PREFIX}.${infoKey}.value`];

        return {
          name: infoKey,
          description: info.description,
          value: this.$options.filters.get(info, 'value', valueColumnFilter),
        };
      });
    },
  },
};
</script>
