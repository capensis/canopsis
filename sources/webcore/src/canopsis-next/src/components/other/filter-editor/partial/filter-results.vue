<template lang="pug">
  v-data-table.elevation-1(
  :headers="headers",
  :items="items",
  :pagination.sync="pagination",
  :total-items="meta.total",
  :rows-per-page-items="rowsPerPageItems",
  :loading="pending"
  )
    template(slot="items", slot-scope="props")
      td {{props.item.v.connector}}
      td {{props.item.v.connector_name}}
      td {{props.item.v.component}}
      td {{props.item.v.resource}}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import EventBus from '@/event-bus';
import { PAGINATION_PER_PAGE_VALUES } from '@/config';

const { mapActions: alarmsMapActions } = createNamespacedHelpers('alarm');

export default {
  props: {
    filter: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  data() {
    return {
      items: [],
      meta: {},
      pending: false,
      pagination: {
        page: 1,
        rowsPerPage: 5,
      },
      headers: [
        {
          text: this.$t('filterEditor.resultsTableHeaders.connector'),
          align: 'left',
          sortable: false,
          value: 'connector',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.connectorName'),
          align: 'left',
          sortable: false,
          value: 'connector_name',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.component'),
          align: 'left',
          sortable: false,
          value: 'component',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.resource'),
          align: 'left',
          sortable: false,
          value: 'resource',
        },
      ],
    };
  },
  computed: {
    rowsPerPageItems() {
      return PAGINATION_PER_PAGE_VALUES;
    },
  },
  watch: {
    pagination(value, oldValue) {
      if (value.page !== oldValue.page || value.rowsPerPage !== oldValue.rowsPerPage) {
        this.fetchList();
      }
    },
  },
  created() {
    EventBus.$on('filter-editor:results:fetch', this.fetchList);
  },
  beforeDestroy() {
    EventBus.$off('filter-editor:results:fetch', this.fetchList);
  },
  methods: {
    ...alarmsMapActions({
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
