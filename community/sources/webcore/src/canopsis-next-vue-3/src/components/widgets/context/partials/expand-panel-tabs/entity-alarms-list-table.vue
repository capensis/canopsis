<template>
  <alarms-list-table-with-pagination
    :widget="widget"
    :meta="meta"
    :alarms="alarms"
    :columns="columns"
    :query="query"
    :loading="pending"
    :expandable="!resolved"
    :hide-pagination="!resolved"
    hide-actions
    @update:query="updateQuery"
  />
</template>

<script>
import { isEqual } from 'lodash';
import { mapActions } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';

import Observer from '@/services/observer';

import { useAlarmStore } from '@/stores/modules/alarm';

import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import AlarmsListTableWithPagination from '@/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue';

export default {
  inject: {
    $periodicRefresh: {
      default: () => new Observer(),
    },
  },

  components: {
    AlarmsListTableWithPagination,
  },

  props: {
    entity: {
      type: Object,
      required: true,
    },
    resolved: {
      type: Boolean,
      default: false,
    },
    columns: {
      type: Array,
      required: false,
    },
  },

  data() {
    return {
      pending: false,
      alarms: [],
      meta: {},
      query: {
        page: 1,
        limit: PAGINATION_LIMIT,
      },
    };
  },

  computed: {
    widget() {
      return generatePreparedDefaultAlarmListWidget();
    },
  },

  watch: {
    resolved: 'fetchList',
    query(query, prevQuery) {
      if (this.resolved && !isEqual(query, prevQuery)) {
        this.fetchList();
      }
    },
  },

  mounted() {
    this.fetchList();

    this.$periodicRefresh.register(this.fetchList);
  },

  beforeUnmount() {
    this.$periodicRefresh.unregister(this.fetchList);
  },

  methods: {
    ...mapActions(useAlarmStore, [
      'fetchResolvedAlarmsListWithoutStore',
      'fetchOpenAlarmsListWithoutStore',
    ]),

    updateQuery(query) {
      this.query = query;
    },

    async fetchList() {
      this.pending = true;

      if (this.resolved) {
        const params = convertWidgetQueryToRequest(this.query);

        const { data, meta } = await this.fetchResolvedAlarmsListWithoutStore({
          params: { ...params, _id: this.entity._id },
        });

        this.alarms = data;
        this.meta = meta;
      } else {
        const alarm = await this.fetchOpenAlarmsListWithoutStore({
          params: { _id: this.entity._id },
        });

        this.alarms = alarm ? [alarm] : [];
        this.meta = { total_count: this.alarms.length };
      }

      this.pending = false;
    },
  },
};
</script>
