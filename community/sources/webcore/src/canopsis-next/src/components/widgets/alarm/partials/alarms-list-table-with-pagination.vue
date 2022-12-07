<template lang="pug">
  alarms-list-table(
    v-model="selected",
    :widget="widget",
    :alarms="alarms",
    :total-items="meta.total_count",
    :pagination.sync="pagination",
    :editing="editing",
    :columns="columns",
    :loading="loading",
    :parent-alarm="parentAlarm",
    :is-tour-enabled="isTourEnabled",
    :hide-children="hideChildren",
    :sticky-header="stickyHeader",
    :refresh-alarms-list="refreshAlarmsList",
    :selectable="selectable",
    :expandable="expandable",
    :hide-actions="hideActions"
  )
    c-table-pagination(
      v-if="!hidePagination",
      :total-items="meta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { isEqual, pick } from 'lodash';

import { SORT_ORDERS } from '@/constants';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
    alarms: {
      type: Array,
      required: true,
    },
    meta: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    editing: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    stickyHeader: {
      type: Boolean,
      default: false,
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hidePagination: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    pagination: {
      get() {
        const { sortDir, sortKey: sortBy = null, multiSortBy = [] } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return { sortBy, descending, multiSortBy };
      },

      set(value) {
        const paginationKeys = ['sortBy', 'descending', 'multiSortBy'];
        const newPagination = pick(value, paginationKeys);
        const oldPagination = pick(this.pagination, paginationKeys);

        if (isEqual(newPagination, oldPagination)) {
          return;
        }

        const {
          sortBy = null,
          descending = false,
          multiSortBy = [],
        } = newPagination;

        const newQuery = {
          sortKey: sortBy,
          sortDir: descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          multiSortBy,
        };

        this.$emit('update:query', {
          ...this.query,
          ...newQuery,
        });
      },
    },
  },
  methods: {
    updateRecordsPerPage(limit) {
      this.$emit('update:query', {
        ...this.query,

        limit,
      });
    },

    updateQueryPage(page) {
      this.$emit('update:query', {
        ...this.query,

        page,
      });
    },
  },
};
</script>
