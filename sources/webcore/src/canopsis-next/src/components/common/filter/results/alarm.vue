<template lang="pug">
  v-data-table.elevation-1(
    :headers="headers",
    :items="items",
    :pagination.sync="pagination",
    :total-items="meta.total",
    :rows-per-page-items="rowsPerPageItems",
    :loading="pending"
  )
    template(slot="items", slot-scope="{ item }")
      td {{ item.v.connector }}
      td {{ item.v.connector_name }}
      td {{ item.v.component }}
      td {{ item.v.resource }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import filterEditorResultsMixin from '@/mixins/filter/editor/results';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

export default {
  mixins: [filterEditorResultsMixin],
  computed: {
    headers() {
      return [
        {
          text: this.$t('filterEditor.resultsTableHeaders.alarm.connector'),
          align: 'left',
          sortable: false,
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.alarm.connectorName'),
          align: 'left',
          sortable: false,
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.alarm.component'),
          align: 'left',
          sortable: false,
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.alarm.resource'),
          align: 'left',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    ...alarmMapActions({
      fetchAlarmListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      const { rowsPerPage: limit, page } = this.pagination;

      this.pending = true;
      const { alarms, total } = await this.fetchAlarmListWithoutStore({
        params: {
          limit,
          skip: limit * (page - 1),
          filter: this.filter,
        },
      });

      this.pending = false;
      this.items = alarms;
      this.meta = { total };
    },
  },
};
</script>
